package qtext

/*
	QString — handle над C++ QString; владение через runtime.AddCleanup + Close().

	# Семантика

	Каждый *QString — это handle на C++ QString в куче. Копирование (Clone)
	дёшево (COW refcount++), мутирующие методы detach-ят внутри Qt. Go-строки и
	[]byte передаются в C как UTF-8 без промежуточной копии; runtime.KeepAlive
	гарантирует, что Go-аргумент не будет собран до возврата из C.

	# Владение

	C++ объект освобождается одним из двух способов:
	- Close() — немедленно, детерминированно;
	- runtime.AddCleanup — когда *QString становится недостижим.

	Close() опционален: если его не вызывать, cleanup удалит объект при сборке
	мусора. Close() идемпотентен и безопасен на nil. Если Close() и cleanup
	гоняются, владение захватывает ровно один участник (атомарный Swap(nil)),
	двойного free нет.

	Close() предпочтителен в долгоживущих процессах для ограничения пиковой
	памяти.

	# Null vs empty

	Qt различает null (QString(), IsNull()==true) и empty-but-not-null
	(например, QString("abc").left(0), IsNull()==false, IsEmpty()==true).

	- NewEmptyQString()              → null QString (Qt: QString()).
	- NewQStringFromString("")       → null QString (fromUtf8(nullptr, 0)).
	- NewQStringFromUtf8(nil)        → null QString.
	- NewQStringFromUtf16(nil)       → null QString.
	- Left(0) / Mid(0,0) / ...       → результат Qt (обычно empty-not-null).

	Используйте IsNull() для различения; IsEmpty() — для проверки нулевой длины
	независимо от null-ности.

	# Use after Close

	После Close() любой метод, кроме Close(), паникует ("qstring: use after
	Close (or nil receiver)"). Внутри все методы используют q.handle().
	Nil-safe методы (IsSharedWith, ToPercentEncoding) используют tryHandle().

	# Соглашения по именам

	Все конструкторы QString начинаются с NewQString* (например,
	NewQStringFromString, NewQStringFromUtf16, NewQStringNumberLL). Аргументные
	суффиксы: QChar, QString, String — тип аргумента; Int, UInt, LL, ULL, Float,
	Double — целочисленные перегрузки; CS — перегрузка с Qt::CaseSensitivity;
	WithMatch — вариант с out-параметром QRegularExpressionMatch*.

	Что сознательно не вошло: см. qstring_doc.go.
*/

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC
#cgo LDFLAGS: -lstdc++
#include <stdint.h>
#include <stdbool.h>
#include <stddef.h>
#include "qstring_c.h"
*/
import "C"

import (
	"runtime"
	"sync/atomic"
	"unsafe"

	"github.com/vadimlxx/go-qt6core/qdata"
)

// qstringState holds the C++ QString pointer in an atomic slot so that
// Close() and the runtime cleanup can race safely: whoever swaps the slot
// to nil owns the C++ object and is responsible for deleting it.
//
// Why *byte and not *C.QString: cgo exposes C.QString as an incomplete Go
// type (forward-declared only under C mode), and atomic.Pointer[T] requires
// T to be a complete type. The pointer value is unchanged; only its Go-side
// static type differs. Cast back to *C.QString at C-call boundaries.
//
// State is a separate allocation (not embedded in QString) because the
// cleanup function must not capture *QString — doing so would keep q
// reachable forever and the cleanup would never run.
type qstringState struct {
	p atomic.Pointer[byte]
}

// cpp returns the underlying *C.QString, or nil if the slot was claimed
// (closed) already. Read-only, does not take ownership.
func (s *qstringState) cpp() *C.QString {
	if bp := s.p.Load(); bp != nil {
		return (*C.QString)(unsafe.Pointer(bp))
	}
	return nil
}

// claim atomically sets the slot to nil and returns the *C.QString that
// was there (if any). Exactly one caller gets a non-nil result for a
// given stored pointer; subsequent claimers see nil. This is what makes
// Close() vs cleanup safe against double free.
func (s *qstringState) claim() *C.QString {
	if bp := s.p.Swap(nil); bp != nil {
		return (*C.QString)(unsafe.Pointer(bp))
	}
	return nil
}

// QString wraps a Qt6 QString (UTF-16 COW string).
//
// Ownership: the underlying C++ QString is released either by Close()
// (eager) or by runtime.AddCleanup (safety net). Close is optional but
// recommended for long-running applications to bound peak memory.
type QString struct {
	state   *qstringState
	cleanup runtime.Cleanup
}

// wrapQString takes ownership of the raw C++ QString pointer.
//
// A cleanup is registered with runtime.AddCleanup: if Close is never
// called, the C++ object is deleted once the *QString becomes unreachable.
// The cleanup receives the *qstringState (not *QString) as its arg — this
// breaks the q → cleanup → q cycle that would otherwise prevent the
// cleanup from ever firing.
func wrapQString(p unsafe.Pointer) *QString {
	if p == nil {
		return nil
	}
	st := &qstringState{}
	st.p.Store((*byte)(p))
	q := &QString{state: st}
	q.cleanup = runtime.AddCleanup(q, func(st *qstringState) {
		if cpp := st.claim(); cpp != nil {
			C.qstring_delete(cpp)
		}
	}, st)
	return q
}

// handle returns the underlying C++ pointer, panicking if q is nil or closed.
// All methods except Close() (and nil-safe methods like IsSharedWith,
// ToPercentEncoding) go through this guard.
func (q *QString) handle() *C.QString {
	if cpp, ok := q.tryHandle(); ok {
		return cpp
	}
	panic("qstring: use after Close (or nil receiver)")
}

// tryHandle is the non-panicking variant of handle. It returns nil, false
// if q is nil, has no state, or has already been closed. Use this when a
// nil/closed *QString is a legitimate input (e.g. IsSharedWith(nil)).
func (q *QString) tryHandle() (*C.QString, bool) {
	if q == nil || q.state == nil {
		return nil, false
	}
	if cpp := q.state.cpp(); cpp != nil {
		return cpp, true
	}
	return nil, false
}

// Close releases the underlying C++ object immediately and cancels the
// cleanup. Close is optional, idempotent, and nil-safe.
//
// If Close races with the runtime cleanup, exactly one of them claims
// the C++ object (via an atomic swap) and deletes it; the other sees nil
// and does nothing. No double free.
//
// Qt: QString::~QString().
func (q *QString) Close() {
	if q == nil || q.state == nil {
		return
	}
	// Best-effort cancel: prevents the cleanup from being scheduled if it
	// hasn't run yet. Stop() returns no value in Go 1.24; if it races with
	// an already-running cleanup, the atomic claim below still makes the
	// outcome safe.
	q.cleanup.Stop()
	if cpp := q.state.claim(); cpp != nil {
		C.qstring_delete(cpp)
	}
}

// ---- construction --------------------------------------------------------

// NewEmptyQString returns a null QString (Qt: QString()).
// Qt: QString() — default constructor.
func NewEmptyQString() *QString {
	return wrapQString(unsafe.Pointer(C.qstring_new_empty()))
}

// NewQStringFromString decodes s as UTF-8 into a new QString.
// Qt: QString::fromUtf8(const char *, qsizetype).
// Empty input yields a null QString (Qt: fromUtf8(nullptr, 0)).
func NewQStringFromString(s string) *QString {
	var p *C.char
	if len(s) > 0 {
		p = (*C.char)(unsafe.Pointer(unsafe.StringData(s)))
	}
	q := wrapQString(unsafe.Pointer(C.qstring_new_from_utf8(p, C.ptrdiff_t(len(s)))))
	runtime.KeepAlive(s)
	return q
}

// NewQStringFromUtf8 decodes b as UTF-8.
// Qt: QString::fromUtf8(QByteArrayView).
func NewQStringFromUtf8(b []byte) *QString {
	var p *C.char
	if len(b) > 0 {
		p = (*C.char)(unsafe.Pointer(&b[0]))
	}
	q := wrapQString(unsafe.Pointer(C.qstring_new_from_utf8(p, C.ptrdiff_t(len(b)))))
	runtime.KeepAlive(b)
	return q
}

// NewQStringFromUtf16 constructs a QString from UTF-16 code units.
// Qt: QString::fromUtf16(const char16_t *, qsizetype).
func NewQStringFromUtf16(units []uint16) *QString {
	var p *C.uint16_t
	if len(units) > 0 {
		p = (*C.uint16_t)(unsafe.Pointer(&units[0]))
	}
	q := wrapQString(unsafe.Pointer(C.qstring_new_from_utf16(p, C.ptrdiff_t(len(units)))))
	runtime.KeepAlive(units)
	return q
}

// NewQStringFromLatin1 decodes b as Latin-1.
// Qt: QString::fromLatin1(QByteArrayView).
func NewQStringFromLatin1(b []byte) *QString {
	var p *C.char
	if len(b) > 0 {
		p = (*C.char)(unsafe.Pointer(&b[0]))
	}
	q := wrapQString(unsafe.Pointer(C.qstring_new_from_latin1(p, C.ptrdiff_t(len(b)))))
	runtime.KeepAlive(b)
	return q
}

// NewQStringFromLocal8Bit decodes b using the local 8-bit encoding.
// Qt: QString::fromLocal8Bit(QByteArrayView).
func NewQStringFromLocal8Bit(b []byte) *QString {
	var p *C.char
	if len(b) > 0 {
		p = (*C.char)(unsafe.Pointer(&b[0]))
	}
	q := wrapQString(unsafe.Pointer(C.qstring_new_from_local8bit(p, C.ptrdiff_t(len(b)))))
	runtime.KeepAlive(b)
	return q
}

// NewQStringFromQByteArray decodes ba as UTF-8.
// Qt: QString::fromUtf8(const QByteArray&).
// A nil ba produces a null QString.
func NewQStringFromQByteArray(ba *QByteArray) *QString {
	if ba == nil {
		return NewEmptyQString()
	}
	bp, ok := ba.tryHandle()
	if !ok {
		return NewEmptyQString()
	}
	q := wrapQString(unsafe.Pointer(C.qstring_new_from_qbytearray(bp)))
	runtime.KeepAlive(ba)
	return q
}

// NewQStringFromQStringView copies v into a new QString (via fromUtf16).
// Qt: QString(QStringView) — explicit constructor.
func NewQStringFromQStringView(v QStringView) *QString {
	if v.IsNull() {
		return NewEmptyQString()
	}
	q := wrapQString(unsafe.Pointer(C.qstring_new_from_utf16(
		(*C.uint16_t)(v.Utf16Unsafe()), C.ptrdiff_t(v.Size()))))
	runtime.KeepAlive(v)
	return q
}

// NewQStringFromUcs4 constructs a QString from Unicode code points.
// Qt: QString::fromUcs4(const char32_t *unicode, qsizetype size).
func NewQStringFromUcs4(codepoints []rune) *QString {
	var p *C.uint32_t
	if len(codepoints) > 0 {
		p = (*C.uint32_t)(unsafe.Pointer(&codepoints[0]))
	}
	q := wrapQString(unsafe.Pointer(C.qstring_new_from_ucs4(p, C.ptrdiff_t(len(codepoints)))))
	runtime.KeepAlive(codepoints)
	return q
}

// NewQStringFromPercentEncoding decodes a percent-encoded QByteArray.
// A nil ba produces a null QString.
//
// Qt: QString::fromPercentEncoding was deprecated in Qt 6.0 and removed in
// 6.11. This function reproduces its behaviour via
// QString::fromUtf8(QByteArray::fromPercentEncoding(...)).
func NewQStringFromPercentEncoding(ba *QByteArray, percent byte) *QString {
	if ba == nil {
		return NewEmptyQString()
	}
	bp, ok := ba.tryHandle()
	if !ok {
		return NewEmptyQString()
	}
	q := wrapQString(unsafe.Pointer(C.qstring_from_percent_encoding(bp, C.char(percent))))
	runtime.KeepAlive(ba)
	return q
}

// ---- convenience constructors ------------------------------------------

// NewQStringFromQChar — Qt: QString(QChar) (implicit constructor; implemented via fromUtf16).
func NewQStringFromQChar(c QChar) *QString {
	return NewQStringFromUtf16([]uint16{uint16(c)})
}

// NewQStringFromRepeat — Qt: QString(qsizetype, QChar).
func NewQStringFromRepeat(n int, c QChar) *QString {
	s := NewEmptyQString()
	s.ResizeFill(n, c)
	return s
}

// NewQString dispatches to the matching constructor based on the concrete
// Go argument type. Qt: no single counterpart.
func NewQString(arg ...any) *QString {
	if len(arg) == 0 || arg[0] == nil {
		return NewEmptyQString()
	}

	if len(arg) == 1 {
		switch v := arg[0].(type) {
		case string:
			return NewQStringFromString(v)
		case []byte:
			return NewQStringFromUtf8(v)
		case []uint16:
			return NewQStringFromUtf16(v)
		case []rune:
			return NewQStringFromUcs4(v)
		case *QByteArray:
			return NewQStringFromQByteArray(v)
		case QStringView:
			return NewQStringFromQStringView(v)
		case QChar:
			return NewQStringFromQChar(v)
		case *QString:
			return v.Clone()
		}
	}

	if len(arg) == 2 {
		n, ok1 := arg[0].(int)
		c, ok2 := arg[1].(QChar)
		if ok1 && ok2 {
			return NewQStringFromRepeat(n, c)
		}
	}

	panic("NewQString invalid argument")
}

// ---- instance methods ---------------------------------------------------

// Clone returns a shallow COW copy.
// Qt: QString(const QString&) — copy constructor.
func (q *QString) Clone() *QString {
	p := q.handle()
	r := wrapQString(unsafe.Pointer(C.qstring_new_copy(p)))
	runtime.KeepAlive(q)
	return r
}

// Swap swaps q and other (no data copy, only COW refcount handles).
// Qt: QString::swap(QString &other).
func (q *QString) Swap(other *QString) {
	p := q.handle()
	op := other.handle()
	C.qstring_swap(p, op)
	runtime.KeepAlive(q)
	runtime.KeepAlive(other)
}

// Size — number of UTF-16 code units.
// Qt: QString::size() const.
func (q *QString) Size() int {
	p := q.handle()
	r := int(C.qstring_size(p))
	runtime.KeepAlive(q)
	return r
}

// Len is an alias for Size.
// Qt: QString::length() const.
func (q *QString) Len() int { return q.Size() }

// Capacity — allocated capacity in code units.
// Qt: QString::capacity() const.
func (q *QString) Capacity() int {
	p := q.handle()
	r := int(C.qstring_capacity(p))
	runtime.KeepAlive(q)
	return r
}

// IsEmpty — zero length.
// Qt: QString::isEmpty() const.
func (q *QString) IsEmpty() bool {
	p := q.handle()
	r := bool(C.qstring_is_empty(p))
	runtime.KeepAlive(q)
	return r
}

// IsNull — null QString (default-constructed or from null ptr).
// Qt: QString::isNull() const.
func (q *QString) IsNull() bool {
	p := q.handle()
	r := bool(C.qstring_is_null(p))
	runtime.KeepAlive(q)
	return r
}

// IsDetached — not shared with any other QString.
// Qt: QString::isDetached() const.
func (q *QString) IsDetached() bool {
	p := q.handle()
	r := bool(C.qstring_is_detached(p))
	runtime.KeepAlive(q)
	return r
}

// Detach ensures q is not shared with any other QString.
// Qt: QString::detach().
func (q *QString) Detach() {
	p := q.handle()
	C.qstring_detach(p)
	runtime.KeepAlive(q)
}

// IsSharedWith reports whether q and other share data.
// Qt: QString::isSharedWith(const QString&) const.
// Returns false if other is nil or closed.
func (q *QString) IsSharedWith(other *QString) bool {
	p := q.handle()
	op, ok := other.tryHandle()
	if !ok {
		runtime.KeepAlive(q)
		return false
	}
	r := bool(C.qstring_is_shared_with(p, op))
	runtime.KeepAlive(q)
	runtime.KeepAlive(other)
	return r
}

// Utf16Unsafe returns the raw internal UTF-16 data pointer (no copy).
// Qt: QString::utf16() const.
// The pointer is valid only while q is alive and not mutated.
func (q *QString) Utf16Unsafe() unsafe.Pointer {
	p := q.handle()
	r := unsafe.Pointer(C.qstring_utf16(p))
	runtime.KeepAlive(q)
	return r
}

// Utf16 returns a copy of the contents as UTF-16 code units.
// Qt: QString::utf16() const (copied).
func (q *QString) Utf16() []uint16 {
	p := q.handle()
	n := int(C.qstring_size(p))
	if n == 0 {
		runtime.KeepAlive(q)
		return nil
	}
	src := C.qstring_utf16(p)
	out := make([]uint16, n)
	copy(out, unsafe.Slice((*uint16)(unsafe.Pointer(src)), n))
	runtime.KeepAlive(q)
	return out
}

// GoString decodes the contents as UTF-8 into a Go string.
// Qt: no direct counterpart — via QString::toUtf8().
func (q *QString) GoString() string {
	ba := q.Utf8()
	defer ba.Close()
	return ba.String()
}

// String implements fmt.Stringer.
func (q *QString) String() string { return q.GoString() }

// At returns the i-th QChar. Panics if i is out of range.
// Qt: QString::at(qsizetype) const.
func (q *QString) At(i int) QChar {
	p := q.handle()
	if i < 0 || i >= int(C.qstring_size(p)) {
		runtime.KeepAlive(q)
		panic("qstring: At out of range")
	}
	r := C.qstring_at(p, C.ptrdiff_t(i))
	runtime.KeepAlive(q)
	return NewQChar(uint16(r))
}

// Front returns the first QChar. Panics if empty.
// Qt: QString::front() const.
func (q *QString) Front() QChar {
	p := q.handle()
	if bool(C.qstring_is_empty(p)) {
		runtime.KeepAlive(q)
		panic("qstring: Front on empty QString")
	}
	r := C.qstring_front(p)
	runtime.KeepAlive(q)
	return NewQChar(uint16(r))
}

// Back returns the last QChar. Panics if empty.
// Qt: QString::back() const.
func (q *QString) Back() QChar {
	p := q.handle()
	if bool(C.qstring_is_empty(p)) {
		runtime.KeepAlive(q)
		panic("qstring: Back on empty QString")
	}
	r := C.qstring_back(p)
	runtime.KeepAlive(q)
	return NewQChar(uint16(r))
}

// ToQStringView returns a zero-copy QStringView that keeps q alive.
// Qt: QStringView(const QString&) — implicit conversion.
func (q *QString) ToQStringView() QStringView {
	p := q.handle()
	v := newQStringViewWithOwner(unsafe.Pointer(C.qstring_utf16(p)), int(C.qstring_size(p)), q)
	runtime.KeepAlive(q)
	return v
}

// Utf8 returns a fresh *QByteArray with UTF-8 encoding.
// Qt: QString::toUtf8() const &.
func (q *QString) Utf8() *QByteArray {
	p := q.handle()
	var out *C.QByteArray
	C.qstring_to_utf8_qba(p, &out)
	runtime.KeepAlive(q)
	return wrapQByteArray(unsafe.Pointer(out))
}

// Latin1 returns a fresh *QByteArray with Latin-1 encoding.
// Qt: QString::toLatin1() const &.
func (q *QString) Latin1() *QByteArray {
	p := q.handle()
	var out *C.QByteArray
	C.qstring_to_latin1_qba(p, &out)
	runtime.KeepAlive(q)
	return wrapQByteArray(unsafe.Pointer(out))
}

// Local8Bit returns a fresh *QByteArray with the local 8-bit encoding.
// Qt: QString::toLocal8Bit() const &.
func (q *QString) Local8Bit() *QByteArray {
	p := q.handle()
	var out *C.QByteArray
	C.qstring_to_local8bit_qba(p, &out)
	runtime.KeepAlive(q)
	return wrapQByteArray(unsafe.Pointer(out))
}

// ToPercentEncoding returns a percent-encoded QByteArray.
// exclude/include may be nil or closed (treated as empty QByteArray).
//
// Qt: QString::toPercentEncoding was deprecated in Qt 6.0 and removed in
// 6.11. This method reproduces its behaviour via
// QString::toUtf8().toPercentEncoding(...).
func (q *QString) ToPercentEncoding(exclude, include *QByteArray, percent byte) *QByteArray {
	p := q.handle()
	var excl, incl *C.QByteArray
	if exclude != nil {
		if bp, ok := exclude.tryHandle(); ok {
			excl = bp
		}
	}
	if include != nil {
		if bp, ok := include.tryHandle(); ok {
			incl = bp
		}
	}
	var out *C.QByteArray
	C.qstring_to_percent_encoding(p, excl, incl, C.char(percent), &out)
	runtime.KeepAlive(q)
	runtime.KeepAlive(exclude)
	runtime.KeepAlive(include)
	return wrapQByteArray(unsafe.Pointer(out))
}

// Clear — make the string null.
// Qt: QString::clear().
func (q *QString) Clear() {
	p := q.handle()
	C.qstring_clear(p)
	runtime.KeepAlive(q)
}

// Resize — change length, new tail content undefined.
// Qt: QString::resize(qsizetype).
func (q *QString) Resize(n int) {
	p := q.handle()
	C.qstring_resize(p, C.ptrdiff_t(n))
	runtime.KeepAlive(q)
}

// ResizeFill — change length, fill new units with c.
// Qt: QString::resize(qsizetype, QChar).
func (q *QString) ResizeFill(n int, c QChar) {
	p := q.handle()
	C.qstring_resize_fill(p, C.ptrdiff_t(n), C.uint16_t(c.Unicode()))
	runtime.KeepAlive(q)
}

// ResizeForOverwrite — change length without initialising the tail.
// Qt: QString::resizeForOverwrite(qsizetype).
func (q *QString) ResizeForOverwrite(n int) {
	p := q.handle()
	C.qstring_resize_for_overwrite(p, C.ptrdiff_t(n))
	runtime.KeepAlive(q)
}

// Reserve — pre-allocate storage for at least n units.
// Qt: QString::reserve(qsizetype).
func (q *QString) Reserve(n int) {
	p := q.handle()
	C.qstring_reserve(p, C.ptrdiff_t(n))
	runtime.KeepAlive(q)
}

// Squeeze — release unused capacity.
// Qt: QString::squeeze().
func (q *QString) Squeeze() {
	p := q.handle()
	C.qstring_squeeze(p)
	runtime.KeepAlive(q)
}

// Truncate — cut to length pos.
// Qt: QString::truncate(qsizetype).
func (q *QString) Truncate(pos int) {
	p := q.handle()
	C.qstring_truncate(p, C.ptrdiff_t(pos))
	runtime.KeepAlive(q)
}

// Chop — remove n code units from the end.
// Qt: QString::chop(qsizetype).
func (q *QString) Chop(n int) {
	p := q.handle()
	C.qstring_chop(p, C.ptrdiff_t(n))
	runtime.KeepAlive(q)
}

// Fill — replace contents with n copies of c (n<0 → pad to current size).
// Qt: QString::fill(QChar, qsizetype).
func (q *QString) Fill(c QChar, n int) {
	p := q.handle()
	C.qstring_fill(p, C.uint16_t(c.Unicode()), C.ptrdiff_t(n))
	runtime.KeepAlive(q)
}

// Append appends a single QChar.
// Qt: QString::append(QChar).
func (q *QString) Append(c QChar) {
	p := q.handle()
	C.qstring_append_qchar(p, C.uint16_t(c.Unicode()))
	runtime.KeepAlive(q)
}

// AppendQString appends another QString.
// Qt: QString::append(const QString&).
func (q *QString) AppendQString(other *QString) {
	p := q.handle()
	op := other.handle()
	C.qstring_append_qstring(p, op)
	runtime.KeepAlive(q)
	runtime.KeepAlive(other)
}

// AppendString appends a Go string (UTF-8).
// Qt: QString::append(const QString&) (string decoded as UTF-8).
func (q *QString) AppendString(s string) {
	p := q.handle()
	var sp *C.char
	if len(s) > 0 {
		sp = (*C.char)(unsafe.Pointer(unsafe.StringData(s)))
	}
	C.qstring_append_utf8(p, sp, C.ptrdiff_t(len(s)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(s)
}

// Prepend prepends a single QChar.
// Qt: QString::prepend(QChar).
func (q *QString) Prepend(c QChar) {
	p := q.handle()
	C.qstring_prepend_qchar(p, C.uint16_t(c.Unicode()))
	runtime.KeepAlive(q)
}

// PrependQString prepends another QString.
// Qt: QString::prepend(const QString&).
func (q *QString) PrependQString(other *QString) {
	p := q.handle()
	op := other.handle()
	C.qstring_prepend_qstring(p, op)
	runtime.KeepAlive(q)
	runtime.KeepAlive(other)
}

// PrependString prepends a Go string (UTF-8).
// Qt: QString::prepend(const QString&).
func (q *QString) PrependString(s string) {
	p := q.handle()
	var sp *C.char
	if len(s) > 0 {
		sp = (*C.char)(unsafe.Pointer(unsafe.StringData(s)))
	}
	C.qstring_prepend_utf8(p, sp, C.ptrdiff_t(len(s)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(s)
}

// Insert inserts a single QChar at position i.
// Qt: QString::insert(qsizetype, QChar).
func (q *QString) Insert(i int, c QChar) {
	p := q.handle()
	C.qstring_insert_qchar(p, C.ptrdiff_t(i), C.uint16_t(c.Unicode()))
	runtime.KeepAlive(q)
}

// InsertQString inserts another QString at position i.
// Qt: QString::insert(qsizetype, const QString&).
func (q *QString) InsertQString(i int, other *QString) {
	p := q.handle()
	op := other.handle()
	C.qstring_insert_qstring(p, C.ptrdiff_t(i), op)
	runtime.KeepAlive(q)
	runtime.KeepAlive(other)
}

// InsertString inserts a Go string (UTF-8) at position i.
// Qt: QString::insert(qsizetype, const QString&).
func (q *QString) InsertString(i int, s string) {
	p := q.handle()
	var sp *C.char
	if len(s) > 0 {
		sp = (*C.char)(unsafe.Pointer(unsafe.StringData(s)))
	}
	C.qstring_insert_utf8(p, C.ptrdiff_t(i), sp, C.ptrdiff_t(len(s)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(s)
}

// Remove removes n code units starting at i.
// Qt: QString::remove(qsizetype, qsizetype).
func (q *QString) Remove(i, n int) {
	p := q.handle()
	C.qstring_remove(p, C.ptrdiff_t(i), C.ptrdiff_t(n))
	runtime.KeepAlive(q)
}

// RemoveAt removes the code unit at position i.
// Qt: QString::removeAt(qsizetype).
func (q *QString) RemoveAt(i int) {
	p := q.handle()
	C.qstring_remove_at(p, C.ptrdiff_t(i))
	runtime.KeepAlive(q)
}

// RemoveFirst removes the first code unit (no-op if empty).
// Qt: QString::removeFirst().
func (q *QString) RemoveFirst() {
	p := q.handle()
	C.qstring_remove_first(p)
	runtime.KeepAlive(q)
}

// RemoveLast removes the last code unit (no-op if empty).
// Qt: QString::removeLast().
func (q *QString) RemoveLast() {
	p := q.handle()
	C.qstring_remove_last(p)
	runtime.KeepAlive(q)
}

// RemoveQChar removes all occurrences of c.
// Qt: QString::remove(QChar, Qt::CaseSensitivity).
func (q *QString) RemoveQChar(c QChar, caseSensitive bool) {
	p := q.handle()
	C.qstring_remove_qchar(p, C.uint16_t(c.Unicode()), csInt(caseSensitive))
	runtime.KeepAlive(q)
}

// RemoveQString removes all occurrences of other.
// Qt: QString::remove(const QString&, Qt::CaseSensitivity).
func (q *QString) RemoveQString(other *QString, caseSensitive bool) {
	p := q.handle()
	op := other.handle()
	C.qstring_remove_qstring(p, op, csInt(caseSensitive))
	runtime.KeepAlive(q)
	runtime.KeepAlive(other)
}

// RemoveString removes all occurrences of a Go string (UTF-8).
// Qt: QString::remove(const QString&, Qt::CaseSensitivity).
func (q *QString) RemoveString(s string, caseSensitive bool) {
	p := q.handle()
	var sp *C.char
	if len(s) > 0 {
		sp = (*C.char)(unsafe.Pointer(unsafe.StringData(s)))
	}
	C.qstring_remove_utf8(p, sp, C.ptrdiff_t(len(s)), csInt(caseSensitive))
	runtime.KeepAlive(q)
	runtime.KeepAlive(s)
}

// ReplaceQCharQChar — Qt: QString::replace(QChar, QChar).
func (q *QString) ReplaceQCharQChar(before, after QChar) {
	p := q.handle()
	C.qstring_replace_qchar_qchar(p, C.uint16_t(before.Unicode()),
		C.uint16_t(after.Unicode()))
	runtime.KeepAlive(q)
}

// ReplaceQCharQString — Qt: QString::replace(QChar, const QString&).
func (q *QString) ReplaceQCharQString(before QChar, after *QString) {
	p := q.handle()
	ap := after.handle()
	C.qstring_replace_qchar_qstring(p, C.uint16_t(before.Unicode()), ap)
	runtime.KeepAlive(q)
	runtime.KeepAlive(after)
}

// ReplaceQCharString — Qt: QString::replace(QChar, const QString&) (after decoded as UTF-8).
func (q *QString) ReplaceQCharString(before QChar, after string) {
	p := q.handle()
	var ap *C.char
	if len(after) > 0 {
		ap = (*C.char)(unsafe.Pointer(unsafe.StringData(after)))
	}
	C.qstring_replace_qchar_utf8(p, C.uint16_t(before.Unicode()), ap, C.ptrdiff_t(len(after)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(after)
}

// ReplaceQStringQChar — Qt: QString::replace(const QString&, QChar, Qt::CaseSensitivity)
// with cs = Qt::CaseSensitive. (Qt does have this overload since 6.0.)
func (q *QString) ReplaceQStringQChar(before *QString, after QChar) {
	p := q.handle()
	bp := before.handle()
	C.qstring_replace_qstring_qchar(p, bp, C.uint16_t(after.Unicode()))
	runtime.KeepAlive(q)
	runtime.KeepAlive(before)
}

// ReplaceQStringQString — Qt: QString::replace(const QString&, const QString&).
func (q *QString) ReplaceQStringQString(before, after *QString) {
	p := q.handle()
	bp := before.handle()
	ap := after.handle()
	C.qstring_replace_qstring_qstring(p, bp, ap)
	runtime.KeepAlive(q)
	runtime.KeepAlive(before)
	runtime.KeepAlive(after)
}

// ReplaceQStringString — Qt: QString::replace(const QString&, const QString&) (after UTF-8 decoded).
func (q *QString) ReplaceQStringString(before *QString, after string) {
	p := q.handle()
	bp := before.handle()
	var ap *C.char
	if len(after) > 0 {
		ap = (*C.char)(unsafe.Pointer(unsafe.StringData(after)))
	}
	C.qstring_replace_qstring_utf8(p, bp, ap, C.ptrdiff_t(len(after)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(before)
	runtime.KeepAlive(after)
}

// ReplaceStringQChar — Qt: QString::replace(const QString&, QChar, Qt::CaseSensitivity)
// with cs = Qt::CaseSensitive (before UTF-8 decoded).
func (q *QString) ReplaceStringQChar(before string, after QChar) {
	p := q.handle()
	var bp *C.char
	if len(before) > 0 {
		bp = (*C.char)(unsafe.Pointer(unsafe.StringData(before)))
	}
	C.qstring_replace_utf8_qchar(p, bp, C.ptrdiff_t(len(before)),
		C.uint16_t(after.Unicode()))
	runtime.KeepAlive(q)
	runtime.KeepAlive(before)
}

// ReplaceStringQString — Qt: QString::replace(const QString&, const QString&) (before UTF-8 decoded).
func (q *QString) ReplaceStringQString(before string, after *QString) {
	p := q.handle()
	ap := after.handle()
	var bp *C.char
	if len(before) > 0 {
		bp = (*C.char)(unsafe.Pointer(unsafe.StringData(before)))
	}
	C.qstring_replace_utf8_qstring(p, bp, C.ptrdiff_t(len(before)), ap)
	runtime.KeepAlive(q)
	runtime.KeepAlive(before)
	runtime.KeepAlive(after)
}

// ReplaceStringString — Qt: QString::replace(const QString&, const QString&) (both UTF-8 decoded).
func (q *QString) ReplaceStringString(before, after string) {
	p := q.handle()
	var bp, ap *C.char
	if len(before) > 0 {
		bp = (*C.char)(unsafe.Pointer(unsafe.StringData(before)))
	}
	if len(after) > 0 {
		ap = (*C.char)(unsafe.Pointer(unsafe.StringData(after)))
	}
	C.qstring_replace_utf8_utf8(p, bp, C.ptrdiff_t(len(before)),
		ap, C.ptrdiff_t(len(after)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(before)
	runtime.KeepAlive(after)
}

// ReplaceRangeString — Qt: QString::replace(qsizetype, qsizetype, const QString&) (UTF-8 decoded).
func (q *QString) ReplaceRangeString(i, n int, after string) {
	p := q.handle()
	var ap *C.char
	if len(after) > 0 {
		ap = (*C.char)(unsafe.Pointer(unsafe.StringData(after)))
	}
	C.qstring_replace_range_utf8(p, C.ptrdiff_t(i), C.ptrdiff_t(n),
		ap, C.ptrdiff_t(len(after)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(after)
}

// ReplaceQCharQCharCS — Qt: QString::replace(QChar, QChar, Qt::CaseSensitivity).
func (q *QString) ReplaceQCharQCharCS(before, after QChar, caseSensitive bool) {
	p := q.handle()
	C.qstring_replace_qchar_qchar_cs(p, C.uint16_t(before.Unicode()),
		C.uint16_t(after.Unicode()), csInt(caseSensitive))
	runtime.KeepAlive(q)
}

// ReplaceQCharQStringCS — Qt: QString::replace(QChar, const QString&, Qt::CaseSensitivity).
func (q *QString) ReplaceQCharQStringCS(before QChar, after *QString, caseSensitive bool) {
	p := q.handle()
	ap := after.handle()
	C.qstring_replace_qchar_qstring_cs(p, C.uint16_t(before.Unicode()), ap,
		csInt(caseSensitive))
	runtime.KeepAlive(q)
	runtime.KeepAlive(after)
}

// ReplaceQStringQStringCS — Qt: QString::replace(const QString&, const QString&, Qt::CaseSensitivity).
func (q *QString) ReplaceQStringQStringCS(before, after *QString, caseSensitive bool) {
	p := q.handle()
	bp := before.handle()
	ap := after.handle()
	C.qstring_replace_qstring_qstring_cs(p, bp, ap, csInt(caseSensitive))
	runtime.KeepAlive(q)
	runtime.KeepAlive(before)
	runtime.KeepAlive(after)
}

// ReplaceStringStringCS — Qt: QString::replace(const QString&, const QString&, Qt::CaseSensitivity) (UTF-8 decoded).
func (q *QString) ReplaceStringStringCS(before, after string, caseSensitive bool) {
	p := q.handle()
	var bp, ap *C.char
	if len(before) > 0 {
		bp = (*C.char)(unsafe.Pointer(unsafe.StringData(before)))
	}
	if len(after) > 0 {
		ap = (*C.char)(unsafe.Pointer(unsafe.StringData(after)))
	}
	C.qstring_replace_utf8_utf8_cs(p, bp, C.ptrdiff_t(len(before)),
		ap, C.ptrdiff_t(len(after)), csInt(caseSensitive))
	runtime.KeepAlive(q)
	runtime.KeepAlive(before)
	runtime.KeepAlive(after)
}

// Left — Qt: QString::left(qsizetype) const &.
func (q *QString) Left(n int) *QString {
	p := q.handle()
	r := wrapQString(unsafe.Pointer(C.qstring_left(p, C.ptrdiff_t(n))))
	runtime.KeepAlive(q)
	return r
}

// Right — Qt: QString::right(qsizetype) const &.
func (q *QString) Right(n int) *QString {
	p := q.handle()
	r := wrapQString(unsafe.Pointer(C.qstring_right(p, C.ptrdiff_t(n))))
	runtime.KeepAlive(q)
	return r
}

// Mid — Qt: QString::mid(qsizetype, qsizetype) const &.
func (q *QString) Mid(pos, n int) *QString {
	p := q.handle()
	r := wrapQString(unsafe.Pointer(C.qstring_mid(p, C.ptrdiff_t(pos), C.ptrdiff_t(n))))
	runtime.KeepAlive(q)
	return r
}

// First — panics if n is out of [0, Size()]. Qt: QString::first(qsizetype).
func (q *QString) First(n int) *QString {
	p := q.handle()
	if n < 0 || n > int(C.qstring_size(p)) {
		runtime.KeepAlive(q)
		panic("qstring: First out of range")
	}
	r := wrapQString(unsafe.Pointer(C.qstring_first(p, C.ptrdiff_t(n))))
	runtime.KeepAlive(q)
	return r
}

// Last — panics if n is out of [0, Size()]. Qt: QString::last(qsizetype).
func (q *QString) Last(n int) *QString {
	p := q.handle()
	if n < 0 || n > int(C.qstring_size(p)) {
		runtime.KeepAlive(q)
		panic("qstring: Last out of range")
	}
	r := wrapQString(unsafe.Pointer(C.qstring_last(p, C.ptrdiff_t(n))))
	runtime.KeepAlive(q)
	return r
}

// Sliced — panics if pos out of range. Qt: QString::sliced(qsizetype).
func (q *QString) Sliced(pos int) *QString {
	p := q.handle()
	if pos < 0 || pos > int(C.qstring_size(p)) {
		runtime.KeepAlive(q)
		panic("qstring: Sliced out of range")
	}
	r := wrapQString(unsafe.Pointer(C.qstring_sliced(p, C.ptrdiff_t(pos))))
	runtime.KeepAlive(q)
	return r
}

// SlicedN — panics if pos or n out of range. Qt: QString::sliced(qsizetype, qsizetype).
func (q *QString) SlicedN(pos, n int) *QString {
	p := q.handle()
	if pos < 0 || n < 0 || pos+n > int(C.qstring_size(p)) {
		runtime.KeepAlive(q)
		panic("qstring: SlicedN out of range")
	}
	r := wrapQString(unsafe.Pointer(C.qstring_sliced_n(p, C.ptrdiff_t(pos), C.ptrdiff_t(n))))
	runtime.KeepAlive(q)
	return r
}

// Chopped — panics if n out of [0, Size()]. Qt: QString::chopped(qsizetype).
func (q *QString) Chopped(n int) *QString {
	p := q.handle()
	if n < 0 || n > int(C.qstring_size(p)) {
		runtime.KeepAlive(q)
		panic("qstring: Chopped out of range")
	}
	r := wrapQString(unsafe.Pointer(C.qstring_chopped(p, C.ptrdiff_t(n))))
	runtime.KeepAlive(q)
	return r
}

// Repeated — panics if times < 0. Qt: QString::repeated(qsizetype) const.
func (q *QString) Repeated(times int) *QString {
	p := q.handle()
	if times < 0 {
		runtime.KeepAlive(q)
		panic("qstring: Repeated negative times")
	}
	r := wrapQString(unsafe.Pointer(C.qstring_repeated(p, C.ptrdiff_t(times))))
	runtime.KeepAlive(q)
	return r
}

// ---- search --------------------------------------------------------------

// IndexOfQChar — Qt: QString::indexOf(QChar, qsizetype, Qt::CaseSensitivity).
func (q *QString) IndexOfQChar(c QChar, from int, caseSensitive bool) int {
	p := q.handle()
	r := int(C.qstring_index_of_qchar(p, C.uint16_t(c.Unicode()),
		C.ptrdiff_t(from), csInt(caseSensitive)))
	runtime.KeepAlive(q)
	return r
}

// IndexOfQString — Qt: QString::indexOf(const QString&, qsizetype, Qt::CaseSensitivity).
func (q *QString) IndexOfQString(other *QString, from int, caseSensitive bool) int {
	p := q.handle()
	op := other.handle()
	r := int(C.qstring_index_of_qstring(p, op, C.ptrdiff_t(from), csInt(caseSensitive)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(other)
	return r
}

// IndexOfString — Qt: QString::indexOf(const QString&, qsizetype, Qt::CaseSensitivity) (UTF-8 decoded).
func (q *QString) IndexOfString(s string, from int, caseSensitive bool) int {
	p := q.handle()
	var sp *C.char
	if len(s) > 0 {
		sp = (*C.char)(unsafe.Pointer(unsafe.StringData(s)))
	}
	r := int(C.qstring_index_of_utf8(p, sp, C.ptrdiff_t(len(s)),
		C.ptrdiff_t(from), csInt(caseSensitive)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(s)
	return r
}

// LastIndexOfQChar — Qt: QString::lastIndexOf(QChar, qsizetype, Qt::CaseSensitivity).
func (q *QString) LastIndexOfQChar(c QChar, from int, caseSensitive bool) int {
	p := q.handle()
	r := int(C.qstring_last_index_of_qchar(p, C.uint16_t(c.Unicode()),
		C.ptrdiff_t(from), csInt(caseSensitive)))
	runtime.KeepAlive(q)
	return r
}

// LastIndexOfQString — Qt: QString::lastIndexOf(const QString&, qsizetype, Qt::CaseSensitivity).
func (q *QString) LastIndexOfQString(other *QString, from int, caseSensitive bool) int {
	p := q.handle()
	op := other.handle()
	r := int(C.qstring_last_index_of_qstring(p, op, C.ptrdiff_t(from),
		csInt(caseSensitive)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(other)
	return r
}

// LastIndexOfString — Qt: QString::lastIndexOf(const QString&, qsizetype, Qt::CaseSensitivity) (UTF-8 decoded).
func (q *QString) LastIndexOfString(s string, from int, caseSensitive bool) int {
	p := q.handle()
	var sp *C.char
	if len(s) > 0 {
		sp = (*C.char)(unsafe.Pointer(unsafe.StringData(s)))
	}
	r := int(C.qstring_last_index_of_utf8(p, sp, C.ptrdiff_t(len(s)),
		C.ptrdiff_t(from), csInt(caseSensitive)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(s)
	return r
}

// ContainsQChar — Qt: QString::contains(QChar, Qt::CaseSensitivity).
func (q *QString) ContainsQChar(c QChar, caseSensitive bool) bool {
	p := q.handle()
	r := bool(C.qstring_contains_qchar(p, C.uint16_t(c.Unicode()), csInt(caseSensitive)))
	runtime.KeepAlive(q)
	return r
}

// ContainsQString — Qt: QString::contains(const QString&, Qt::CaseSensitivity).
func (q *QString) ContainsQString(other *QString, caseSensitive bool) bool {
	p := q.handle()
	op := other.handle()
	r := bool(C.qstring_contains_qstring(p, op, csInt(caseSensitive)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(other)
	return r
}

// ContainsString — Qt: QString::contains(const QString&, Qt::CaseSensitivity) (UTF-8 decoded).
func (q *QString) ContainsString(s string, caseSensitive bool) bool {
	p := q.handle()
	var sp *C.char
	if len(s) > 0 {
		sp = (*C.char)(unsafe.Pointer(unsafe.StringData(s)))
	}
	r := bool(C.qstring_contains_utf8(p, sp, C.ptrdiff_t(len(s)), csInt(caseSensitive)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(s)
	return r
}

// CountQChar — Qt: QString::count(QChar, Qt::CaseSensitivity).
func (q *QString) CountQChar(c QChar, caseSensitive bool) int {
	p := q.handle()
	r := int(C.qstring_count_qchar(p, C.uint16_t(c.Unicode()), csInt(caseSensitive)))
	runtime.KeepAlive(q)
	return r
}

// CountQString — Qt: QString::count(const QString&, Qt::CaseSensitivity).
func (q *QString) CountQString(other *QString, caseSensitive bool) int {
	p := q.handle()
	op := other.handle()
	r := int(C.qstring_count_qstring(p, op, csInt(caseSensitive)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(other)
	return r
}

// CountString — Qt: QString::count(const QString&, Qt::CaseSensitivity) (UTF-8 decoded).
func (q *QString) CountString(s string, caseSensitive bool) int {
	p := q.handle()
	var sp *C.char
	if len(s) > 0 {
		sp = (*C.char)(unsafe.Pointer(unsafe.StringData(s)))
	}
	r := int(C.qstring_count_utf8(p, sp, C.ptrdiff_t(len(s)), csInt(caseSensitive)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(s)
	return r
}

// StartsWithQChar — Qt: QString::startsWith(QChar, Qt::CaseSensitivity).
func (q *QString) StartsWithQChar(c QChar, caseSensitive bool) bool {
	p := q.handle()
	r := bool(C.qstring_starts_with_qchar(p, C.uint16_t(c.Unicode()),
		csInt(caseSensitive)))
	runtime.KeepAlive(q)
	return r
}

// StartsWithQString — Qt: QString::startsWith(const QString&, Qt::CaseSensitivity).
func (q *QString) StartsWithQString(other *QString, caseSensitive bool) bool {
	p := q.handle()
	op := other.handle()
	r := bool(C.qstring_starts_with_qstring(p, op, csInt(caseSensitive)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(other)
	return r
}

// StartsWithString — Qt: QString::startsWith(const QString&, Qt::CaseSensitivity) (UTF-8 decoded).
func (q *QString) StartsWithString(s string, caseSensitive bool) bool {
	p := q.handle()
	var sp *C.char
	if len(s) > 0 {
		sp = (*C.char)(unsafe.Pointer(unsafe.StringData(s)))
	}
	r := bool(C.qstring_starts_with_utf8(p, sp, C.ptrdiff_t(len(s)), csInt(caseSensitive)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(s)
	return r
}

// EndsWithQChar — Qt: QString::endsWith(QChar, Qt::CaseSensitivity).
func (q *QString) EndsWithQChar(c QChar, caseSensitive bool) bool {
	p := q.handle()
	r := bool(C.qstring_ends_with_qchar(p, C.uint16_t(c.Unicode()), csInt(caseSensitive)))
	runtime.KeepAlive(q)
	return r
}

// EndsWithQString — Qt: QString::endsWith(const QString&, Qt::CaseSensitivity).
func (q *QString) EndsWithQString(other *QString, caseSensitive bool) bool {
	p := q.handle()
	op := other.handle()
	r := bool(C.qstring_ends_with_qstring(p, op, csInt(caseSensitive)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(other)
	return r
}

// EndsWithString — Qt: QString::endsWith(const QString&, Qt::CaseSensitivity) (UTF-8 decoded).
func (q *QString) EndsWithString(s string, caseSensitive bool) bool {
	p := q.handle()
	var sp *C.char
	if len(s) > 0 {
		sp = (*C.char)(unsafe.Pointer(unsafe.StringData(s)))
	}
	r := bool(C.qstring_ends_with_utf8(p, sp, C.ptrdiff_t(len(s)), csInt(caseSensitive)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(s)
	return r
}

// ---- QRegularExpression search ------------------------------------------

// IndexOfQRegularExpression — Qt: QString::indexOf(const QRegularExpression&, qsizetype, QRegularExpressionMatch*)
// with match=nullptr.
func (q *QString) IndexOfQRegularExpression(re *QRegularExpression, from int) int {
	return q.IndexOfQRegularExpressionWithMatch(re, from, nil)
}

// IndexOfQRegularExpressionWithMatch — Qt: QString::indexOf(const QRegularExpression&, qsizetype, QRegularExpressionMatch*).
func (q *QString) IndexOfQRegularExpressionWithMatch(re *QRegularExpression, from int,
	match *QRegularExpressionMatch) int {
	p := q.handle()
	rep := re.handle()
	var cm *C.QRegularExpressionMatch
	if match != nil && match.p != nil {
		cm = match.p
	}
	r := int(C.qstring_index_of_qregexp(p, rep, C.ptrdiff_t(from), cm))
	runtime.KeepAlive(q)
	runtime.KeepAlive(re)
	runtime.KeepAlive(match)
	return r
}

// LastIndexOfQRegularExpression — Qt: QString::lastIndexOf(const QRegularExpression&, qsizetype, QRegularExpressionMatch*)
// with match=nullptr.
func (q *QString) LastIndexOfQRegularExpression(re *QRegularExpression, from int) int {
	return q.LastIndexOfQRegularExpressionWithMatch(re, from, nil)
}

// LastIndexOfQRegularExpressionWithMatch — Qt: QString::lastIndexOf(const QRegularExpression&, qsizetype, QRegularExpressionMatch*).
func (q *QString) LastIndexOfQRegularExpressionWithMatch(re *QRegularExpression, from int,
	match *QRegularExpressionMatch) int {
	p := q.handle()
	rep := re.handle()
	var cm *C.QRegularExpressionMatch
	if match != nil && match.p != nil {
		cm = match.p
	}
	r := int(C.qstring_last_index_of_qregexp(p, rep, C.ptrdiff_t(from), cm))
	runtime.KeepAlive(q)
	runtime.KeepAlive(re)
	runtime.KeepAlive(match)
	return r
}

// ContainsQRegularExpression — Qt: QString::contains(const QRegularExpression&, QRegularExpressionMatch*)
// with match=nullptr.
func (q *QString) ContainsQRegularExpression(re *QRegularExpression) bool {
	return q.ContainsQRegularExpressionWithMatch(re, nil)
}

// ContainsQRegularExpressionWithMatch — Qt: QString::contains(const QRegularExpression&, QRegularExpressionMatch*).
func (q *QString) ContainsQRegularExpressionWithMatch(re *QRegularExpression,
	match *QRegularExpressionMatch) bool {
	p := q.handle()
	rep := re.handle()
	var cm *C.QRegularExpressionMatch
	if match != nil && match.p != nil {
		cm = match.p
	}
	r := bool(C.qstring_contains_qregexp(p, rep, cm))
	runtime.KeepAlive(q)
	runtime.KeepAlive(re)
	runtime.KeepAlive(match)
	return r
}

// CountQRegularExpression — Qt: QString::count(const QRegularExpression&).
// Match-overlap semantics depend on the linked Qt build (some use non-overlapping
// globalMatch, others iterate overlapping matches). Verify against your Qt.
func (q *QString) CountQRegularExpression(re *QRegularExpression) int {
	p := q.handle()
	rep := re.handle()
	r := int(C.qstring_count_qregexp(p, rep))
	runtime.KeepAlive(q)
	runtime.KeepAlive(re)
	return r
}

// ReplaceQRegularExpression — Qt: QString::replace(const QRegularExpression&, const QString&).
func (q *QString) ReplaceQRegularExpression(re *QRegularExpression, after *QString) {
	p := q.handle()
	rep := re.handle()
	ap := after.handle()
	C.qstring_replace_qregexp_qstring(p, rep, ap)
	runtime.KeepAlive(q)
	runtime.KeepAlive(re)
	runtime.KeepAlive(after)
}

// ReplaceQRegularExpressionString — Qt: QString::replace(const QRegularExpression&, const QString&) (after UTF-8 decoded).
func (q *QString) ReplaceQRegularExpressionString(re *QRegularExpression, after string) {
	qa := NewQStringFromString(after)
	defer qa.Close()
	q.ReplaceQRegularExpression(re, qa)
}

// RemoveQRegularExpression — Qt: QString::remove(const QRegularExpression&).
func (q *QString) RemoveQRegularExpression(re *QRegularExpression) {
	p := q.handle()
	rep := re.handle()
	C.qstring_remove_qregexp(p, rep)
	runtime.KeepAlive(q)
	runtime.KeepAlive(re)
}

// ---- case / trim / escape -----------------------------------------------

// ToLower — Qt: QString::toLower() const &.
func (q *QString) ToLower() *QString {
	p := q.handle()
	r := wrapQString(unsafe.Pointer(C.qstring_to_lower(p)))
	runtime.KeepAlive(q)
	return r
}

// ToUpper — Qt: QString::toUpper() const &.
func (q *QString) ToUpper() *QString {
	p := q.handle()
	r := wrapQString(unsafe.Pointer(C.qstring_to_upper(p)))
	runtime.KeepAlive(q)
	return r
}

// ToCaseFolded — Qt: QString::toCaseFolded() const &.
func (q *QString) ToCaseFolded() *QString {
	p := q.handle()
	r := wrapQString(unsafe.Pointer(C.qstring_to_case_folded(p)))
	runtime.KeepAlive(q)
	return r
}

// ToHtmlEscaped — Qt: QString::toHtmlEscaped() const.
func (q *QString) ToHtmlEscaped() *QString {
	p := q.handle()
	r := wrapQString(unsafe.Pointer(C.qstring_to_html_escaped(p)))
	runtime.KeepAlive(q)
	return r
}

// IsLower — Qt: QString::isLower() const.
func (q *QString) IsLower() bool {
	p := q.handle()
	r := bool(C.qstring_is_lower(p))
	runtime.KeepAlive(q)
	return r
}

// IsUpper — Qt: QString::isUpper() const.
func (q *QString) IsUpper() bool {
	p := q.handle()
	r := bool(C.qstring_is_upper(p))
	runtime.KeepAlive(q)
	return r
}

// Trimmed — Qt: QString::trimmed() const &.
func (q *QString) Trimmed() *QString {
	p := q.handle()
	r := wrapQString(unsafe.Pointer(C.qstring_trimmed(p)))
	runtime.KeepAlive(q)
	return r
}

// Simplified — Qt: QString::simplified() const &.
func (q *QString) Simplified() *QString {
	p := q.handle()
	r := wrapQString(unsafe.Pointer(C.qstring_simplified(p)))
	runtime.KeepAlive(q)
	return r
}

// LeftJustified — Qt: QString::leftJustified(qsizetype, QChar, bool) const.
func (q *QString) LeftJustified(width int, fill QChar, truncate bool) *QString {
	p := q.handle()
	r := wrapQString(unsafe.Pointer(C.qstring_left_justified(p,
		C.ptrdiff_t(width), C.uint16_t(fill.Unicode()), C.bool(truncate))))
	runtime.KeepAlive(q)
	return r
}

// RightJustified — Qt: QString::rightJustified(qsizetype, QChar, bool) const.
func (q *QString) RightJustified(width int, fill QChar, truncate bool) *QString {
	p := q.handle()
	r := wrapQString(unsafe.Pointer(C.qstring_right_justified(p,
		C.ptrdiff_t(width), C.uint16_t(fill.Unicode()), C.bool(truncate))))
	runtime.KeepAlive(q)
	return r
}

// ---- normalization -------------------------------------------------------

// Normalized — Qt: QString::normalized(NormalizationForm, QChar::UnicodeVersion)
// with version = QChar::Unicode_Unassigned (0).
func (q *QString) Normalized(mode NormalizationForm) *QString {
	return q.NormalizedWithVersion(mode, 0)
}

// NormalizedWithVersion — Qt: QString::normalized(NormalizationForm, QChar::UnicodeVersion).
func (q *QString) NormalizedWithVersion(mode NormalizationForm, version int) *QString {
	p := q.handle()
	r := wrapQString(unsafe.Pointer(C.qstring_normalized(p, C.int(mode), C.int(version))))
	runtime.KeepAlive(q)
	return r
}

// ---- compare ------------------------------------------------------------

// CompareQString — Qt: QString::compare(const QString&, Qt::CaseSensitivity).
func (q *QString) CompareQString(other *QString, caseSensitive bool) int {
	p := q.handle()
	op := other.handle()
	r := int(C.qstring_compare_qstring(p, op, csInt(caseSensitive)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(other)
	return r
}

// CompareQChar — Qt: QString::compare(QChar, Qt::CaseSensitivity).
func (q *QString) CompareQChar(c QChar, caseSensitive bool) int {
	p := q.handle()
	r := int(C.qstring_compare_qchar(p, C.uint16_t(c.Unicode()), csInt(caseSensitive)))
	runtime.KeepAlive(q)
	return r
}

// CompareString — Qt: QString::compare(const QString&, Qt::CaseSensitivity) (UTF-8 decoded).
func (q *QString) CompareString(s string, caseSensitive bool) int {
	p := q.handle()
	var sp *C.char
	if len(s) > 0 {
		sp = (*C.char)(unsafe.Pointer(unsafe.StringData(s)))
	}
	r := int(C.qstring_compare_utf8(p, sp, C.ptrdiff_t(len(s)), csInt(caseSensitive)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(s)
	return r
}

// LocaleAwareCompare — Qt: QString::localeAwareCompare(const QString&) const.
func (q *QString) LocaleAwareCompare(other *QString) int {
	p := q.handle()
	op := other.handle()
	r := int(C.qstring_locale_aware_compare_qstring(p, op))
	runtime.KeepAlive(q)
	runtime.KeepAlive(other)
	return r
}

// ---- numeric ------------------------------------------------------------

// ToShort — Qt: QString::toShort(bool*, int).
func (q *QString) ToShort(base int) (int16, bool) {
	p := q.handle()
	var ok C.bool
	r := int16(C.qstring_to_short(p, &ok, C.int(base)))
	runtime.KeepAlive(q)
	return r, bool(ok)
}

// ToUShort — Qt: QString::toUShort(bool*, int).
func (q *QString) ToUShort(base int) (uint16, bool) {
	p := q.handle()
	var ok C.bool
	r := uint16(C.qstring_to_ushort(p, &ok, C.int(base)))
	runtime.KeepAlive(q)
	return r, bool(ok)
}

// ToInt — Qt: QString::toInt(bool*, int).
func (q *QString) ToInt(base int) (int32, bool) {
	p := q.handle()
	var ok C.bool
	r := int32(C.qstring_to_int(p, &ok, C.int(base)))
	runtime.KeepAlive(q)
	return r, bool(ok)
}

// ToUInt — Qt: QString::toUInt(bool*, int).
func (q *QString) ToUInt(base int) (uint32, bool) {
	p := q.handle()
	var ok C.bool
	r := uint32(C.qstring_to_uint(p, &ok, C.int(base)))
	runtime.KeepAlive(q)
	return r, bool(ok)
}

// ToLong — Qt: QString::toLong(bool*, int).
func (q *QString) ToLong(base int) (int64, bool) {
	p := q.handle()
	var ok C.bool
	r := int64(C.qstring_to_long(p, &ok, C.int(base)))
	runtime.KeepAlive(q)
	return r, bool(ok)
}

// ToULong — Qt: QString::toULong(bool*, int).
func (q *QString) ToULong(base int) (uint64, bool) {
	p := q.handle()
	var ok C.bool
	r := uint64(C.qstring_to_ulong(p, &ok, C.int(base)))
	runtime.KeepAlive(q)
	return r, bool(ok)
}

// ToLongLong — Qt: QString::toLongLong(bool*, int).
func (q *QString) ToLongLong(base int) (int64, bool) {
	p := q.handle()
	var ok C.bool
	r := int64(C.qstring_to_longlong(p, &ok, C.int(base)))
	runtime.KeepAlive(q)
	return r, bool(ok)
}

// ToULongLong — Qt: QString::toULongLong(bool*, int).
func (q *QString) ToULongLong(base int) (uint64, bool) {
	p := q.handle()
	var ok C.bool
	r := uint64(C.qstring_to_ulonglong(p, &ok, C.int(base)))
	runtime.KeepAlive(q)
	return r, bool(ok)
}

// ToFloat — Qt: QString::toFloat(bool*).
func (q *QString) ToFloat() (float32, bool) {
	p := q.handle()
	var ok C.bool
	r := float32(C.qstring_to_float(p, &ok))
	runtime.KeepAlive(q)
	return r, bool(ok)
}

// ToDouble — Qt: QString::toDouble(bool*).
func (q *QString) ToDouble() (float64, bool) {
	p := q.handle()
	var ok C.bool
	r := float64(C.qstring_to_double(p, &ok))
	runtime.KeepAlive(q)
	return r, bool(ok)
}

// SetNumLL — Qt: QString::setNum(qlonglong, int).
func (q *QString) SetNumLL(n int64, base int) {
	p := q.handle()
	C.qstring_set_num_ll(p, C.int64_t(n), C.int(base))
	runtime.KeepAlive(q)
}

// SetNumULL — Qt: QString::setNum(qulonglong, int).
func (q *QString) SetNumULL(n uint64, base int) {
	p := q.handle()
	C.qstring_set_num_ull(p, C.uint64_t(n), C.int(base))
	runtime.KeepAlive(q)
}

// SetNumInt — Qt: QString::setNum(int, int).
func (q *QString) SetNumInt(n int32, base int) {
	p := q.handle()
	C.qstring_set_num_int(p, C.int32_t(n), C.int(base))
	runtime.KeepAlive(q)
}

// SetNumUInt — Qt: QString::setNum(uint, int).
func (q *QString) SetNumUInt(n uint32, base int) {
	p := q.handle()
	C.qstring_set_num_uint(p, C.uint32_t(n), C.int(base))
	runtime.KeepAlive(q)
}

// SetNumFloat — Qt: QString::setNum(float, char, int).
func (q *QString) SetNumFloat(v float32, format byte, precision int) {
	p := q.handle()
	C.qstring_set_num_float(p, C.float(v), C.char(format), C.int(precision))
	runtime.KeepAlive(q)
}

// SetNumDouble — Qt: QString::setNum(double, char, int).
func (q *QString) SetNumDouble(v float64, format byte, precision int) {
	p := q.handle()
	C.qstring_set_num_double(p, C.double(v), C.char(format), C.int(precision))
	runtime.KeepAlive(q)
}

// NewQStringNumberLL — Qt: QString::number(qlonglong, int).
func NewQStringNumberLL(n int64, base int) *QString {
	return wrapQString(unsafe.Pointer(C.qstring_number_ll(C.int64_t(n), C.int(base))))
}

// NewQStringNumberULL — Qt: QString::number(qulonglong, int).
func NewQStringNumberULL(n uint64, base int) *QString {
	return wrapQString(unsafe.Pointer(C.qstring_number_ull(C.uint64_t(n), C.int(base))))
}

// NewQStringNumberInt — Qt: QString::number(int, int).
func NewQStringNumberInt(n int32, base int) *QString {
	return wrapQString(unsafe.Pointer(C.qstring_number_int(C.int32_t(n), C.int(base))))
}

// NewQStringNumberUInt — Qt: QString::number(uint, int).
func NewQStringNumberUInt(n uint32, base int) *QString {
	return wrapQString(unsafe.Pointer(C.qstring_number_uint(C.uint32_t(n), C.int(base))))
}

// NewQStringNumberDouble — Qt: QString::number(double, char, int).
func NewQStringNumberDouble(v float64, format byte, precision int) *QString {
	return wrapQString(unsafe.Pointer(C.qstring_number_double(C.double(v),
		C.char(format), C.int(precision))))
}

// ---- static / predicates ------------------------------------------------

// QStringMaxSize — Qt: QString::maxSize(). Not a constructor: returns int.
func QStringMaxSize() int { return int(C.qstring_max_size()) }

// IsRightToLeft — Qt: QString::isRightToLeft() const.
func (q *QString) IsRightToLeft() bool {
	p := q.handle()
	r := bool(C.qstring_is_right_to_left(p))
	runtime.KeepAlive(q)
	return r
}

// IsValidUtf16 — Qt: QString::isValidUtf16() const.
func (q *QString) IsValidUtf16() bool {
	p := q.handle()
	r := bool(C.qstring_is_valid_utf16(p))
	runtime.KeepAlive(q)
	return r
}

// ---- arg ----------------------------------------------------------------

// ArgLl — Qt: QString::arg(qlonglong, int, int, QChar) const.
func (q *QString) ArgLl(n int64, fieldWidth, base int, fillChar QChar) *QString {
	p := q.handle()
	r := wrapQString(unsafe.Pointer(C.qstring_arg_ll(p, C.int64_t(n),
		C.int(fieldWidth), C.int(base), C.uint16_t(fillChar.Unicode()))))
	runtime.KeepAlive(q)
	return r
}

// ArgUll — Qt: QString::arg(qulonglong, int, int, QChar) const.
func (q *QString) ArgUll(n uint64, fieldWidth, base int, fillChar QChar) *QString {
	p := q.handle()
	r := wrapQString(unsafe.Pointer(C.qstring_arg_ull(p, C.uint64_t(n),
		C.int(fieldWidth), C.int(base), C.uint16_t(fillChar.Unicode()))))
	runtime.KeepAlive(q)
	return r
}

// ArgDouble — Qt: QString::arg(double, int, char, int, QChar) const.
func (q *QString) ArgDouble(v float64, fieldWidth int, format byte, precision int,
	fillChar QChar) *QString {
	p := q.handle()
	r := wrapQString(unsafe.Pointer(C.qstring_arg_double(p, C.double(v),
		C.int(fieldWidth), C.char(format), C.int(precision),
		C.uint16_t(fillChar.Unicode()))))
	runtime.KeepAlive(q)
	return r
}

// ArgQChar — Qt: QString::arg(QChar, int, QChar) const.
func (q *QString) ArgQChar(c QChar, fieldWidth int, fillChar QChar) *QString {
	p := q.handle()
	r := wrapQString(unsafe.Pointer(C.qstring_arg_qchar(p, C.uint16_t(c.Unicode()),
		C.int(fieldWidth), C.uint16_t(fillChar.Unicode()))))
	runtime.KeepAlive(q)
	return r
}

// ArgQString — Qt: QString::arg(const QString&, int, QChar) const.
func (q *QString) ArgQString(other *QString, fieldWidth int, fillChar QChar) *QString {
	p := q.handle()
	op := other.handle()
	r := wrapQString(unsafe.Pointer(C.qstring_arg_qstring(p, op,
		C.int(fieldWidth), C.uint16_t(fillChar.Unicode()))))
	runtime.KeepAlive(q)
	runtime.KeepAlive(other)
	return r
}

// ArgString — Qt: QString::arg(const QString&, int, QChar) const (UTF-8 decoded).
func (q *QString) ArgString(s string, fieldWidth int, fillChar QChar) *QString {
	p := q.handle()
	var sp *C.char
	if len(s) > 0 {
		sp = (*C.char)(unsafe.Pointer(unsafe.StringData(s)))
	}
	r := wrapQString(unsafe.Pointer(C.qstring_arg_utf8(p, sp, C.ptrdiff_t(len(s)),
		C.int(fieldWidth), C.uint16_t(fillChar.Unicode()))))
	runtime.KeepAlive(q)
	runtime.KeepAlive(s)
	return r
}

// ---- split --------------------------------------------------------------

// SplitQString — Qt: QString::split(const QString&, Qt::SplitBehavior, Qt::CaseSensitivity).
func (q *QString) SplitQString(sep *QString, behavior qdata.SplitBehavior,
	caseSensitive bool) []*QString {
	p := q.handle()
	sp := sep.handle()
	var count C.ptrdiff_t
	arr := C.qstring_split_qstring(p, sp, C.int(behavior),
		csInt(caseSensitive), &count)
	runtime.KeepAlive(q)
	runtime.KeepAlive(sep)
	if arr == nil {
		return nil
	}
	defer C.qstring_free_string_array(arr)
	n := int(count)
	if n == 0 {
		return nil
	}
	ptrs := unsafe.Slice(arr, n)
	result := make([]*QString, n)
	for i := 0; i < n; i++ {
		result[i] = wrapQString(unsafe.Pointer(ptrs[i]))
	}
	return result
}

// SplitQChar — Qt: QString::split(QChar, Qt::SplitBehavior, Qt::CaseSensitivity).
func (q *QString) SplitQChar(sep QChar, behavior qdata.SplitBehavior,
	caseSensitive bool) []*QString {
	p := q.handle()
	var count C.ptrdiff_t
	arr := C.qstring_split_qchar(p, C.uint16_t(sep.Unicode()),
		C.int(behavior), csInt(caseSensitive), &count)
	runtime.KeepAlive(q)
	if arr == nil {
		return nil
	}
	defer C.qstring_free_string_array(arr)
	n := int(count)
	if n == 0 {
		return nil
	}
	ptrs := unsafe.Slice(arr, n)
	result := make([]*QString, n)
	for i := 0; i < n; i++ {
		result[i] = wrapQString(unsafe.Pointer(ptrs[i]))
	}
	return result
}

// SplitQRegularExpression — Qt: QString::split(const QRegularExpression&, Qt::SplitBehavior).
func (q *QString) SplitQRegularExpression(re *QRegularExpression,
	behavior qdata.SplitBehavior) []*QString {
	p := q.handle()
	rep := re.handle()
	var count C.ptrdiff_t
	arr := C.qstring_split_qregexp(p, rep, C.int(behavior), &count)
	runtime.KeepAlive(q)
	runtime.KeepAlive(re)
	if arr == nil {
		return nil
	}
	defer C.qstring_free_string_array(arr)
	n := int(count)
	if n == 0 {
		return nil
	}
	ptrs := unsafe.Slice(arr, n)
	result := make([]*QString, n)
	for i := 0; i < n; i++ {
		result[i] = wrapQString(unsafe.Pointer(ptrs[i]))
	}
	return result
}

// ---- section ------------------------------------------------------------

// SectionQChar — Qt: QString::section(QChar, qsizetype, qsizetype, SectionFlags).
func (q *QString) SectionQChar(sep QChar, start, end int, flags SectionFlag) *QString {
	p := q.handle()
	r := wrapQString(unsafe.Pointer(C.qstring_section_qchar(p,
		C.uint16_t(sep.Unicode()), C.ptrdiff_t(start), C.ptrdiff_t(end), C.int(flags))))
	runtime.KeepAlive(q)
	return r
}

// SectionQString — Qt: QString::section(const QString&, qsizetype, qsizetype, SectionFlags).
func (q *QString) SectionQString(sep *QString, start, end int, flags SectionFlag) *QString {
	p := q.handle()
	sp := sep.handle()
	r := wrapQString(unsafe.Pointer(C.qstring_section_qstring(p, sp,
		C.ptrdiff_t(start), C.ptrdiff_t(end), C.int(flags))))
	runtime.KeepAlive(q)
	runtime.KeepAlive(sep)
	return r
}

// SectionQRegularExpression — Qt: QString::section(const QRegularExpression&, qsizetype, qsizetype, SectionFlags).
func (q *QString) SectionQRegularExpression(re *QRegularExpression, start, end int,
	flags SectionFlag) *QString {
	p := q.handle()
	rep := re.handle()
	r := wrapQString(unsafe.Pointer(C.qstring_section_qregexp(p, rep,
		C.ptrdiff_t(start), C.ptrdiff_t(end), C.int(flags))))
	runtime.KeepAlive(q)
	runtime.KeepAlive(re)
	return r
}

// ---- toUcs4 -------------------------------------------------------------

// ToUcs4 — Qt: QString::toUcs4() const.
func (q *QString) ToUcs4() []uint32 {
	p := q.handle()
	var count C.ptrdiff_t
	arr := C.qstring_to_ucs4(p, &count)
	runtime.KeepAlive(q)
	if arr == nil {
		return nil
	}
	defer C.qstring_free_ucs4_array(arr)
	n := int(count)
	if n == 0 {
		return nil
	}
	src := unsafe.Slice((*C.uint32_t)(unsafe.Pointer(arr)), n)
	result := make([]uint32, n)
	for i := 0; i < n; i++ {
		result[i] = uint32(src[i])
	}
	return result
}
