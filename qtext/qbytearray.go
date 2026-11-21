package qtext

/*
	QByteArray — handle над C++ QByteArray; владение через runtime.AddCleanup + Close().

	# Семантика

	Каждый *QByteArray — handle на C++ QByteArray в куче. Копирование (Clone)
	дёшево (COW refcount++), мутирующие методы detach-ят внутри Qt. Go-срезы
	передаются в C без промежуточной копии; runtime.KeepAlive гарантирует, что
	Go-аргумент не будет собран до возврата из C.

	# Владение

	C++ объект освобождается одним из двух способов:
	- Close() — немедленно, детерминированно;
	- runtime.AddCleanup — когда *QByteArray становится недостижим.

	Close() опционален: если его не вызывать, cleanup удалит объект при сборке
	мусора. Close() идемпотентен и безопасен на nil. Если Close() и cleanup
	гоняются, владение захватывает ровно один участник (атомарный Swap(nil)),
	двойного free нет.

	# Use after Close

	После Close() любой метод, кроме Close(), паникует ("qbytearray: use after
	Close (or nil receiver)"). Nil-safe методы (IsSharedWith, ToPercentEncoding)
	используют tryHandle().

	# Соглашения по именам

	Все конструкторы QByteArray начинаются с NewQByteArray* (например,
	NewQByteArrayFromBytes, NewQByteArrayFromFill, NewQByteArrayNumberLL).

	# Что сознательно не вошло

	- fromRawData / setRawData — требуют runtime.Pinner + явного владения.
		Используйте NewQByteArrayFromBytes (копирует данные).
	- Итераторы (begin/end/rbegin/rend) — не идиоматичны в Go; используйте Bytes().
	- operator[] — Go не поддерживает перегрузку; используйте At/SetAt.
	- Неявные конверсии (operator const char*, std::string_view) —
		используйте Bytes()/String().
	- fromCFData / toNSData / fromEcmaUint8Array — платформо-специфично.
*/

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC
#cgo LDFLAGS: -lstdc++
#include <stdint.h>
#include <stdbool.h>
#include <stddef.h>
#include "qbytearray_c.h"
*/
import "C"

import (
	"runtime"
	"sync/atomic"
	"unsafe"
)

// qbytearrayState holds the C++ QByteArray pointer in an atomic slot so
// that Close() and the runtime cleanup can race safely: whoever swaps the
// slot to nil owns the C++ object and is responsible for deleting it.
//
// Why *byte and not *C.QByteArray: cgo exposes C.QByteArray as an
// incomplete Go type (forward-declared only under C mode), and
// atomic.Pointer[T] requires T to be a complete type. The pointer value is
// unchanged; only its Go-side static type differs. Cast back to
// *C.QByteArray at C-call boundaries.
//
// State is a separate allocation (not embedded in QByteArray) because the
// cleanup function must not capture *QByteArray — doing so would keep q
// reachable forever and the cleanup would never run.
type qbytearrayState struct {
	p atomic.Pointer[byte]
}

// cpp returns the underlying *C.QByteArray, or nil if the slot was claimed
// (closed) already. Read-only, does not take ownership.
func (s *qbytearrayState) cpp() *C.QByteArray {
	if bp := s.p.Load(); bp != nil {
		return (*C.QByteArray)(unsafe.Pointer(bp))
	}
	return nil
}

// claim atomically sets the slot to nil and returns the *C.QByteArray that
// was there (if any). Exactly one caller gets a non-nil result for a given
// stored pointer; subsequent claimers see nil. This is what makes Close()
// vs cleanup safe against double free.
func (s *qbytearrayState) claim() *C.QByteArray {
	if bp := s.p.Swap(nil); bp != nil {
		return (*C.QByteArray)(unsafe.Pointer(bp))
	}
	return nil
}

// QByteArray wraps a Qt6 QByteArray.
// Zero value is not usable; use NewQByteArray / NewQByteArrayFromBytes / etc.
type QByteArray struct {
	state   *qbytearrayState
	cleanup runtime.Cleanup
}

// wrapQByteArray wraps an opaque C++ QByteArray* into *QByteArray.
// Accepts unsafe.Pointer for symmetry with wrapQString and to allow the
// C-ABI to return either *C.QByteArray or go_qt6core_qbytearray_t*.
// Takes ownership: the cleanup (or Close) frees the C++ object.
func wrapQByteArray(p unsafe.Pointer) *QByteArray {
	if p == nil {
		return nil
	}
	st := &qbytearrayState{}
	st.p.Store((*byte)(p))
	q := &QByteArray{state: st}
	q.cleanup = runtime.AddCleanup(q, func(st *qbytearrayState) {
		if cpp := st.claim(); cpp != nil {
			C.qbytearray_delete(cpp)
		}
	}, st)
	return q
}

// handle returns the underlying C++ pointer, panicking if q is nil or closed.
// All methods except Close() (and nil-safe methods like IsSharedWith,
// ToPercentEncoding) go through this guard.
func (q *QByteArray) handle() *C.QByteArray {
	if cpp, ok := q.tryHandle(); ok {
		return cpp
	}
	panic("qbytearray: use after Close (or nil receiver)")
}

// tryHandle is the non-panicking variant of handle. It returns nil, false
// if q is nil, has no state, or has already been closed.
func (q *QByteArray) tryHandle() (*C.QByteArray, bool) {
	if q == nil || q.state == nil {
		return nil, false
	}
	if cpp := q.state.cpp(); cpp != nil {
		return cpp, true
	}
	return nil, false
}

// ---- construction --------------------------------------------------------

// NewQByteArrayFromBytes constructs a QByteArray copying b.
// Qt: QByteArray(const char *, qsizetype) — copied.
func NewQByteArrayFromBytes(b []byte) *QByteArray {
	var p *C.char
	if len(b) > 0 {
		p = (*C.char)(unsafe.Pointer(&b[0]))
	}
	q := wrapQByteArray(unsafe.Pointer(C.qbytearray_new_from_data(p, C.ptrdiff_t(len(b)))))
	runtime.KeepAlive(b)
	return q
}

// NewQByteArrayFromFill returns a QByteArray of size n filled with c.
// Qt: QByteArray(qsizetype, char).
func NewQByteArrayFromFill(n int, c byte) *QByteArray {
	return wrapQByteArray(unsafe.Pointer(
		C.qbytearray_new_fill(C.ptrdiff_t(n), C.char(c))))
}

// NewQByteArrayFromQCharDecomposition returns the Unicode decomposition of c
// as UTF-8 bytes (QChar::decomposition().toUtf8()).
func NewQByteArrayFromQCharDecomposition(c QChar) *QByteArray {
	s := c.Decompose()
	if s == nil {
		return NewQByteArrayFromBytes(nil)
	}
	defer s.Close()
	return s.Utf8()
}

// NewQByteArrayFromQCharDecompositionLatin1 is NewQByteArrayFromQCharDecomposition
// with Latin-1 encoding.
func NewQByteArrayFromQCharDecompositionLatin1(c QChar) *QByteArray {
	s := c.Decompose()
	if s == nil {
		return NewQByteArrayFromBytes(nil)
	}
	defer s.Close()
	return s.Latin1()
}

// NewQByteArray dispatches to the matching constructor based on the
// concrete Go argument types. Qt: no single counterpart.
//
//	NewQByteArray()                 → empty QByteArray
//	NewQByteArray("hello")          → from string (UTF-8)
//	NewQByteArray([]byte{...})      → copy of slice
//	NewQByteArray(n int, c byte)    → n bytes filled with c
//	NewQByteArray(QChar)            → QChar::decomposition().toUtf8()
//	NewQByteArray(*QByteArray)      → COW copy (Clone)
func NewQByteArray(args ...any) *QByteArray {
	switch len(args) {
	case 0:
		return wrapQByteArray(unsafe.Pointer(C.qbytearray_new_empty()))
	case 1:
		switch v := args[0].(type) {
		case []byte:
			return NewQByteArrayFromBytes(v)
		case string:
			return NewQByteArrayFromBytes([]byte(v))
		case QChar:
			return NewQByteArrayFromQCharDecomposition(v)
		case *QByteArray:
			return v.Clone()
		}
	case 2:
		n, ok1 := args[0].(int)
		c, ok2 := args[1].(byte)
		if ok1 && ok2 {
			return NewQByteArrayFromFill(n, c)
		}
	}
	panic("NewQByteArray: unsupported argument combination")
}

// NewQByteArrayNumberLL — Qt: QByteArray::number(qlonglong, int).
func NewQByteArrayNumberLL(n int64, base int) *QByteArray {
	return wrapQByteArray(unsafe.Pointer(C.qbytearray_number_ll(C.int64_t(n), C.int(base))))
}

// NewQByteArrayNumberULL — Qt: QByteArray::number(qulonglong, int).
func NewQByteArrayNumberULL(n uint64, base int) *QByteArray {
	return wrapQByteArray(unsafe.Pointer(C.qbytearray_number_ull(C.uint64_t(n), C.int(base))))
}

// NewQByteArrayNumberDouble — Qt: QByteArray::number(double, char, int).
func NewQByteArrayNumberDouble(v float64, format byte, precision int) *QByteArray {
	return wrapQByteArray(unsafe.Pointer(
		C.qbytearray_number_double(C.double(v), C.char(format), C.int(precision))))
}

// NewQByteArrayFromBase64 — Qt: QByteArray::fromBase64(QByteArrayView, Base64Options).
func NewQByteArrayFromBase64(b []byte, options Base64Options) *QByteArray {
	src := NewQByteArrayFromBytes(b)
	defer src.Close()
	return wrapQByteArray(unsafe.Pointer(C.qbytearray_from_base64(src.handle(), C.int(options))))
}

// NewQByteArrayFromBase64Encoding mirrors QByteArray::fromBase64Encoding.
// Qt: QByteArray::fromBase64Encoding(QByteArrayView, Base64Options).
func NewQByteArrayFromBase64Encoding(b []byte, options Base64Options) (*QByteArray, Base64DecodingStatus) {
	var p *C.char
	if len(b) > 0 {
		p = (*C.char)(unsafe.Pointer(&b[0]))
	}
	var status C.int
	r := wrapQByteArray(unsafe.Pointer(
		C.qbytearray_from_base64_encoding(p, C.ptrdiff_t(len(b)), C.int(options), &status)))
	runtime.KeepAlive(b)
	return r, Base64DecodingStatus(status)
}

// NewQByteArrayFromHex — Qt: QByteArray::fromHex(QByteArrayView).
func NewQByteArrayFromHex(b []byte) *QByteArray {
	src := NewQByteArrayFromBytes(b)
	defer src.Close()
	return wrapQByteArray(unsafe.Pointer(C.qbytearray_from_hex(src.handle())))
}

// NewQByteArrayFromPercentEncoding — Qt: QByteArray::fromPercentEncoding(QByteArrayView, char).
func NewQByteArrayFromPercentEncoding(b []byte, percent byte) *QByteArray {
	src := NewQByteArrayFromBytes(b)
	defer src.Close()
	return wrapQByteArray(unsafe.Pointer(C.qbytearray_from_percent(src.handle(), C.char(percent))))
}

// ---- lifecycle -----------------------------------------------------------

// Close releases the underlying C++ object immediately and cancels the
// cleanup. Close is optional, idempotent, and nil-safe.
//
// If Close races with the runtime cleanup, exactly one of them claims the
// C++ object and deletes it; the other sees nil and does nothing. No
// double free.
// Qt: QByteArray::~QByteArray().
func (q *QByteArray) Close() {
	if q == nil || q.state == nil {
		return
	}
	q.cleanup.Stop()
	if cpp := q.state.claim(); cpp != nil {
		C.qbytearray_delete(cpp)
	}
}

// Clone returns a shallow COW copy (refcount++; cheap).
// Qt: QByteArray(const QByteArray&) — copy constructor.
func (q *QByteArray) Clone() *QByteArray {
	r := wrapQByteArray(unsafe.Pointer(C.qbytearray_new_copy(q.handle())))
	runtime.KeepAlive(q)
	return r
}

// ---- basic accessors -----------------------------------------------------

// Size — Qt: QByteArray::size() const.
func (q *QByteArray) Size() int {
	n := C.qbytearray_size(q.handle())
	runtime.KeepAlive(q)
	return int(n)
}

// Len is an alias for Size. Qt: QByteArray::length() const.
func (q *QByteArray) Len() int { return q.Size() }

// Capacity — Qt: QByteArray::capacity() const.
func (q *QByteArray) Capacity() int {
	n := C.qbytearray_capacity(q.handle())
	runtime.KeepAlive(q)
	return int(n)
}

// IsEmpty — Qt: QByteArray::isEmpty() const.
func (q *QByteArray) IsEmpty() bool {
	v := C.qbytearray_is_empty(q.handle())
	runtime.KeepAlive(q)
	return bool(v)
}

// IsNull — Qt: QByteArray::isNull() const.
func (q *QByteArray) IsNull() bool {
	v := C.qbytearray_is_null(q.handle())
	runtime.KeepAlive(q)
	return bool(v)
}

// IsDetached — Qt: QByteArray::isDetached() const.
func (q *QByteArray) IsDetached() bool {
	v := C.qbytearray_is_detached(q.handle())
	runtime.KeepAlive(q)
	return bool(v)
}

// Detach — Qt: QByteArray::detach().
func (q *QByteArray) Detach() {
	C.qbytearray_detach(q.handle())
	runtime.KeepAlive(q)
}

// IsSharedWith reports whether q and other share the same internal buffer.
// Returns false if other is nil or closed.
// Qt: QByteArray::isSharedWith(const QByteArray&) const.
func (q *QByteArray) IsSharedWith(other *QByteArray) bool {
	p := q.handle()
	op, ok := other.tryHandle()
	if !ok {
		runtime.KeepAlive(q)
		return false
	}
	v := C.qbytearray_is_shared_with(p, op)
	runtime.KeepAlive(q)
	runtime.KeepAlive(other)
	return bool(v)
}

// QByteArrayMaxSize mirrors QByteArray::maxSize(). Not a constructor.
func QByteArrayMaxSize() int { return int(C.qbytearray_max_size()) }

// Bytes returns a copy of q's bytes.
// Qt: no direct counterpart — via QByteArray::constData() copied into Go.
func (q *QByteArray) Bytes() []byte {
	n := C.qbytearray_size(q.handle())
	if n == 0 {
		runtime.KeepAlive(q)
		return nil
	}
	data := C.qbytearray_data(q.handle())
	b := C.GoBytes(unsafe.Pointer(data), C.int(n))
	runtime.KeepAlive(q)
	return b
}

// String returns a copy of q's bytes as a Go string.
// Qt: no direct counterpart — via QByteArray::constData() copied into Go.
func (q *QByteArray) String() string {
	n := C.qbytearray_size(q.handle())
	if n == 0 {
		runtime.KeepAlive(q)
		return ""
	}
	data := C.qbytearray_data(q.handle())
	s := C.GoStringN(data, C.int(n))
	runtime.KeepAlive(q)
	return s
}

// ---- modifiers -----------------------------------------------------------

// Resize — Qt: QByteArray::resize(qsizetype).
func (q *QByteArray) Resize(n int) {
	C.qbytearray_resize(q.handle(), C.ptrdiff_t(n))
	runtime.KeepAlive(q)
}

// ResizeFill — Qt: QByteArray::resize(qsizetype, char).
func (q *QByteArray) ResizeFill(n int, c byte) {
	C.qbytearray_resize_fill(q.handle(), C.ptrdiff_t(n), C.char(c))
	runtime.KeepAlive(q)
}

// ResizeForOverwrite — Qt: QByteArray::resizeForOverwrite(qsizetype).
func (q *QByteArray) ResizeForOverwrite(n int) {
	C.qbytearray_resize_for_overwrite(q.handle(), C.ptrdiff_t(n))
	runtime.KeepAlive(q)
}

// Fill — Qt: QByteArray::fill(char, qsizetype).
func (q *QByteArray) Fill(c byte, n int) {
	C.qbytearray_fill(q.handle(), C.char(c), C.ptrdiff_t(n))
	runtime.KeepAlive(q)
}

// Clear — Qt: QByteArray::clear().
func (q *QByteArray) Clear() {
	C.qbytearray_clear(q.handle())
	runtime.KeepAlive(q)
}

// Reserve — Qt: QByteArray::reserve(qsizetype).
func (q *QByteArray) Reserve(n int) {
	C.qbytearray_reserve(q.handle(), C.ptrdiff_t(n))
	runtime.KeepAlive(q)
}

// Squeeze — Qt: QByteArray::squeeze().
func (q *QByteArray) Squeeze() {
	C.qbytearray_squeeze(q.handle())
	runtime.KeepAlive(q)
}

// Truncate — Qt: QByteArray::truncate(qsizetype).
func (q *QByteArray) Truncate(pos int) {
	C.qbytearray_truncate(q.handle(), C.ptrdiff_t(pos))
	runtime.KeepAlive(q)
}

// Chop — Qt: QByteArray::chop(qsizetype).
func (q *QByteArray) Chop(n int) {
	C.qbytearray_chop(q.handle(), C.ptrdiff_t(n))
	runtime.KeepAlive(q)
}

// NullTerminate appends a '\0' to the end of q in place (if not already present).
// Qt: QByteArray::append('\0') / QByteArray::constData() contract.
func (q *QByteArray) NullTerminate() {
	C.qbytearray_null_terminate(q.handle())
	runtime.KeepAlive(q)
}

// NullTerminated returns a new QByteArray with a trailing '\0' appended.
// Qt: no direct counterpart — via QByteArray + '\0'.
func (q *QByteArray) NullTerminated() *QByteArray {
	r := wrapQByteArray(unsafe.Pointer(C.qbytearray_null_terminated(q.handle())))
	runtime.KeepAlive(q)
	return r
}

// ---- element access ------------------------------------------------------

// At returns the i-th byte. Panics if i is out of range.
// Qt: QByteArray::at(qsizetype) const.
func (q *QByteArray) At(i int) byte {
	n := q.Size()
	if i < 0 || i >= n {
		panic("qbytearray: At out of range")
	}
	c := C.qbytearray_at(q.handle(), C.ptrdiff_t(i))
	runtime.KeepAlive(q)
	return byte(c)
}

// SetAt sets the i-th byte. Panics if i is out of range.
// Qt: QByteArray::operator[](qsizetype).
func (q *QByteArray) SetAt(i int, c byte) {
	n := q.Size()
	if i < 0 || i >= n {
		panic("qbytearray: SetAt out of range")
	}
	C.qbytearray_set_at(q.handle(), C.ptrdiff_t(i), C.char(c))
	runtime.KeepAlive(q)
}

// Front returns the first byte. Panics if empty.
// Qt: QByteArray::front() const.
func (q *QByteArray) Front() byte {
	if q.IsEmpty() {
		panic("qbytearray: Front on empty QByteArray")
	}
	c := C.qbytearray_front(q.handle())
	runtime.KeepAlive(q)
	return byte(c)
}

// Back returns the last byte. Panics if empty.
// Qt: QByteArray::back() const.
func (q *QByteArray) Back() byte {
	if q.IsEmpty() {
		panic("qbytearray: Back on empty QByteArray")
	}
	c := C.qbytearray_back(q.handle())
	runtime.KeepAlive(q)
	return byte(c)
}

// ---- search --------------------------------------------------------------

// IndexOf — Qt: QByteArray::indexOf(char, qsizetype) const.
func (q *QByteArray) IndexOf(c byte, from int) int {
	n := C.qbytearray_index_of_char(q.handle(), C.char(c), C.ptrdiff_t(from))
	runtime.KeepAlive(q)
	return int(n)
}

// IndexOfBytes — Qt: QByteArray::indexOf(QByteArrayView, qsizetype) const.
func (q *QByteArray) IndexOfBytes(sub []byte, from int) int {
	var p *C.char
	if len(sub) > 0 {
		p = (*C.char)(unsafe.Pointer(&sub[0]))
	}
	n := C.qbytearray_index_of_data(q.handle(), p, C.ptrdiff_t(len(sub)), C.ptrdiff_t(from))
	runtime.KeepAlive(q)
	runtime.KeepAlive(sub)
	return int(n)
}

// LastIndexOf — Qt: QByteArray::lastIndexOf(char, qsizetype) const.
func (q *QByteArray) LastIndexOf(c byte, from int) int {
	n := C.qbytearray_last_index_of_char(q.handle(), C.char(c), C.ptrdiff_t(from))
	runtime.KeepAlive(q)
	return int(n)
}

// LastIndexOfBytes — Qt: QByteArray::lastIndexOf(QByteArrayView, qsizetype) const.
func (q *QByteArray) LastIndexOfBytes(sub []byte, from int) int {
	var p *C.char
	if len(sub) > 0 {
		p = (*C.char)(unsafe.Pointer(&sub[0]))
	}
	n := C.qbytearray_last_index_of_data(q.handle(), p, C.ptrdiff_t(len(sub)), C.ptrdiff_t(from))
	runtime.KeepAlive(q)
	runtime.KeepAlive(sub)
	return int(n)
}

// Count — Qt: QByteArray::count(char) const.
func (q *QByteArray) Count(c byte) int {
	n := C.qbytearray_count_char(q.handle(), C.char(c))
	runtime.KeepAlive(q)
	return int(n)
}

// CountBytes — Qt: QByteArray::count(QByteArrayView) const.
func (q *QByteArray) CountBytes(sub []byte) int {
	var p *C.char
	if len(sub) > 0 {
		p = (*C.char)(unsafe.Pointer(&sub[0]))
	}
	n := C.qbytearray_count_data(q.handle(), p, C.ptrdiff_t(len(sub)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(sub)
	return int(n)
}

// Contains — Qt: QByteArray::contains(char) const.
func (q *QByteArray) Contains(c byte) bool {
	v := C.qbytearray_contains_char(q.handle(), C.char(c))
	runtime.KeepAlive(q)
	return bool(v)
}

// ContainsBytes — Qt: QByteArray::contains(QByteArrayView) const.
func (q *QByteArray) ContainsBytes(sub []byte) bool {
	var p *C.char
	if len(sub) > 0 {
		p = (*C.char)(unsafe.Pointer(&sub[0]))
	}
	v := C.qbytearray_contains_data(q.handle(), p, C.ptrdiff_t(len(sub)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(sub)
	return bool(v)
}

// Compare returns Qt's compare() result. cs == true → case-sensitive.
// Qt: QByteArray::compare(QByteArrayView, Qt::CaseSensitivity) const.
func (q *QByteArray) Compare(other []byte, cs bool) int {
	var p *C.char
	if len(other) > 0 {
		p = (*C.char)(unsafe.Pointer(&other[0]))
	}
	csInt := C.int(0)
	if !cs {
		csInt = 1
	}
	r := C.qbytearray_compare(q.handle(), p, C.ptrdiff_t(len(other)), csInt)
	runtime.KeepAlive(q)
	runtime.KeepAlive(other)
	return int(r)
}

// ---- slicing -------------------------------------------------------------

// Left — Qt: QByteArray::left(qsizetype) const &.
func (q *QByteArray) Left(n int) *QByteArray {
	r := wrapQByteArray(unsafe.Pointer(C.qbytearray_left(q.handle(), C.ptrdiff_t(n))))
	runtime.KeepAlive(q)
	return r
}

// Right — Qt: QByteArray::right(qsizetype) const &.
func (q *QByteArray) Right(n int) *QByteArray {
	r := wrapQByteArray(unsafe.Pointer(C.qbytearray_right(q.handle(), C.ptrdiff_t(n))))
	runtime.KeepAlive(q)
	return r
}

// Mid — Qt: QByteArray::mid(qsizetype, qsizetype) const &.
func (q *QByteArray) Mid(i, length int) *QByteArray {
	r := wrapQByteArray(unsafe.Pointer(C.qbytearray_mid(q.handle(), C.ptrdiff_t(i), C.ptrdiff_t(length))))
	runtime.KeepAlive(q)
	return r
}

// First — panics if n is out of [0, Size()]. Qt: QByteArray::first(qsizetype).
func (q *QByteArray) First(n int) *QByteArray {
	if n < 0 || n > q.Size() {
		panic("qbytearray: First out of range")
	}
	r := wrapQByteArray(unsafe.Pointer(C.qbytearray_first(q.handle(), C.ptrdiff_t(n))))
	runtime.KeepAlive(q)
	return r
}

// Last — panics if n is out of [0, Size()]. Qt: QByteArray::last(qsizetype).
func (q *QByteArray) Last(n int) *QByteArray {
	if n < 0 || n > q.Size() {
		panic("qbytearray: Last out of range")
	}
	r := wrapQByteArray(unsafe.Pointer(C.qbytearray_last(q.handle(), C.ptrdiff_t(n))))
	runtime.KeepAlive(q)
	return r
}

// Sliced — panics if pos out of range. Qt: QByteArray::sliced(qsizetype).
func (q *QByteArray) Sliced(pos int) *QByteArray {
	if pos < 0 || pos > q.Size() {
		panic("qbytearray: Sliced out of range")
	}
	r := wrapQByteArray(unsafe.Pointer(C.qbytearray_sliced(q.handle(), C.ptrdiff_t(pos))))
	runtime.KeepAlive(q)
	return r
}

// SlicedN — panics if pos or n out of range. Qt: QByteArray::sliced(qsizetype, qsizetype).
func (q *QByteArray) SlicedN(pos, n int) *QByteArray {
	sz := q.Size()
	if pos < 0 || n < 0 || pos+n > sz {
		panic("qbytearray: SlicedN out of range")
	}
	r := wrapQByteArray(unsafe.Pointer(
		C.qbytearray_sliced_n(q.handle(), C.ptrdiff_t(pos), C.ptrdiff_t(n))))
	runtime.KeepAlive(q)
	return r
}

// Chopped — panics if n out of [0, Size()]. Qt: QByteArray::chopped(qsizetype).
func (q *QByteArray) Chopped(n int) *QByteArray {
	if n < 0 || n > q.Size() {
		panic("qbytearray: Chopped out of range")
	}
	r := wrapQByteArray(unsafe.Pointer(C.qbytearray_chopped(q.handle(), C.ptrdiff_t(n))))
	runtime.KeepAlive(q)
	return r
}

// ---- predicates ----------------------------------------------------------

// StartsWith — Qt: QByteArray::startsWith(char) const.
func (q *QByteArray) StartsWith(c byte) bool {
	v := C.qbytearray_starts_with_char(q.handle(), C.char(c))
	runtime.KeepAlive(q)
	return bool(v)
}

// StartsWithBytes — Qt: QByteArray::startsWith(QByteArrayView) const.
func (q *QByteArray) StartsWithBytes(prefix []byte) bool {
	var p *C.char
	if len(prefix) > 0 {
		p = (*C.char)(unsafe.Pointer(&prefix[0]))
	}
	v := C.qbytearray_starts_with_data(q.handle(), p, C.ptrdiff_t(len(prefix)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(prefix)
	return bool(v)
}

// EndsWith — Qt: QByteArray::endsWith(char) const.
func (q *QByteArray) EndsWith(c byte) bool {
	v := C.qbytearray_ends_with_char(q.handle(), C.char(c))
	runtime.KeepAlive(q)
	return bool(v)
}

// EndsWithBytes — Qt: QByteArray::endsWith(QByteArrayView) const.
func (q *QByteArray) EndsWithBytes(suffix []byte) bool {
	var p *C.char
	if len(suffix) > 0 {
		p = (*C.char)(unsafe.Pointer(&suffix[0]))
	}
	v := C.qbytearray_ends_with_data(q.handle(), p, C.ptrdiff_t(len(suffix)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(suffix)
	return bool(v)
}

// IsUpper — Qt: QByteArray::isUpper() const.
func (q *QByteArray) IsUpper() bool {
	v := C.qbytearray_is_upper(q.handle())
	runtime.KeepAlive(q)
	return bool(v)
}

// IsLower — Qt: QByteArray::isLower() const.
func (q *QByteArray) IsLower() bool {
	v := C.qbytearray_is_lower(q.handle())
	runtime.KeepAlive(q)
	return bool(v)
}

// IsValidUtf8 — Qt: QByteArray::isValidUtf8() const.
func (q *QByteArray) IsValidUtf8() bool {
	v := C.qbytearray_is_valid_utf8(q.handle())
	runtime.KeepAlive(q)
	return bool(v)
}

// ---- transforms ----------------------------------------------------------

// ToLower — Qt: QByteArray::toLower() const &.
func (q *QByteArray) ToLower() *QByteArray {
	r := wrapQByteArray(unsafe.Pointer(C.qbytearray_to_lower(q.handle())))
	runtime.KeepAlive(q)
	return r
}

// ToUpper — Qt: QByteArray::toUpper() const &.
func (q *QByteArray) ToUpper() *QByteArray {
	r := wrapQByteArray(unsafe.Pointer(C.qbytearray_to_upper(q.handle())))
	runtime.KeepAlive(q)
	return r
}

// Trimmed — Qt: QByteArray::trimmed() const &.
func (q *QByteArray) Trimmed() *QByteArray {
	r := wrapQByteArray(unsafe.Pointer(C.qbytearray_trimmed(q.handle())))
	runtime.KeepAlive(q)
	return r
}

// Simplified — Qt: QByteArray::simplified() const &.
func (q *QByteArray) Simplified() *QByteArray {
	r := wrapQByteArray(unsafe.Pointer(C.qbytearray_simplified(q.handle())))
	runtime.KeepAlive(q)
	return r
}

// LeftJustified — Qt: QByteArray::leftJustified(qsizetype, char, bool) const.
func (q *QByteArray) LeftJustified(width int, fill byte, truncate bool) *QByteArray {
	r := wrapQByteArray(unsafe.Pointer(
		C.qbytearray_left_justified(q.handle(), C.ptrdiff_t(width), C.char(fill), C.bool(truncate))))
	runtime.KeepAlive(q)
	return r
}

// RightJustified — Qt: QByteArray::rightJustified(qsizetype, char, bool) const.
func (q *QByteArray) RightJustified(width int, fill byte, truncate bool) *QByteArray {
	r := wrapQByteArray(unsafe.Pointer(
		C.qbytearray_right_justified(q.handle(), C.ptrdiff_t(width), C.char(fill), C.bool(truncate))))
	runtime.KeepAlive(q)
	return r
}

// Repeated — Qt: QByteArray::repeated(qsizetype) const.
func (q *QByteArray) Repeated(times int) *QByteArray {
	r := wrapQByteArray(unsafe.Pointer(C.qbytearray_repeated(q.handle(), C.ptrdiff_t(times))))
	runtime.KeepAlive(q)
	return r
}

// ---- append/prepend/insert ----------------------------------------------

// Append — Qt: QByteArray::append(char).
func (q *QByteArray) Append(c byte) {
	C.qbytearray_append_char(q.handle(), C.char(c))
	runtime.KeepAlive(q)
}

// AppendBytes — Qt: QByteArray::append(QByteArrayView).
func (q *QByteArray) AppendBytes(b []byte) {
	var p *C.char
	if len(b) > 0 {
		p = (*C.char)(unsafe.Pointer(&b[0]))
	}
	C.qbytearray_append_data(q.handle(), p, C.ptrdiff_t(len(b)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(b)
}

// AppendString is a convenience wrapper around AppendBytes.
func (q *QByteArray) AppendString(s string) { q.AppendBytes([]byte(s)) }

// AppendCount — Qt: QByteArray::append(qsizetype, char).
func (q *QByteArray) AppendCount(count int, c byte) {
	C.qbytearray_append_count(q.handle(), C.ptrdiff_t(count), C.char(c))
	runtime.KeepAlive(q)
}

// AppendQByteArray — Qt: QByteArray::append(const QByteArray&).
func (q *QByteArray) AppendQByteArray(other *QByteArray) {
	C.qbytearray_append_ba(q.handle(), other.handle())
	runtime.KeepAlive(q)
	runtime.KeepAlive(other)
}

// Prepend — Qt: QByteArray::prepend(char).
func (q *QByteArray) Prepend(c byte) {
	C.qbytearray_prepend_char(q.handle(), C.char(c))
	runtime.KeepAlive(q)
}

// PrependBytes — Qt: QByteArray::prepend(QByteArrayView).
func (q *QByteArray) PrependBytes(b []byte) {
	var p *C.char
	if len(b) > 0 {
		p = (*C.char)(unsafe.Pointer(&b[0]))
	}
	C.qbytearray_prepend_data(q.handle(), p, C.ptrdiff_t(len(b)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(b)
}

// PrependString is a convenience wrapper around PrependBytes.
func (q *QByteArray) PrependString(s string) { q.PrependBytes([]byte(s)) }

// PrependCount — Qt: QByteArray::prepend(qsizetype, char).
func (q *QByteArray) PrependCount(count int, c byte) {
	C.qbytearray_prepend_count(q.handle(), C.ptrdiff_t(count), C.char(c))
	runtime.KeepAlive(q)
}

// PrependQByteArray — Qt: QByteArray::prepend(const QByteArray&).
func (q *QByteArray) PrependQByteArray(other *QByteArray) {
	C.qbytearray_prepend_ba(q.handle(), other.handle())
	runtime.KeepAlive(q)
	runtime.KeepAlive(other)
}

// InsertBytes — Qt: QByteArray::insert(qsizetype, QByteArrayView).
func (q *QByteArray) InsertBytes(i int, b []byte) {
	if i < 0 || i > q.Size() {
		panic("qbytearray: InsertBytes out of range")
	}
	var p *C.char
	if len(b) > 0 {
		p = (*C.char)(unsafe.Pointer(&b[0]))
	}
	C.qbytearray_insert_data(q.handle(), C.ptrdiff_t(i), p, C.ptrdiff_t(len(b)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(b)
}

// InsertString is a convenience wrapper around InsertBytes.
func (q *QByteArray) InsertString(i int, s string) { q.InsertBytes(i, []byte(s)) }

// InsertByte — Qt: QByteArray::insert(qsizetype, char).
func (q *QByteArray) InsertByte(i int, c byte) {
	if i < 0 || i > q.Size() {
		panic("qbytearray: InsertByte out of range")
	}
	C.qbytearray_insert_char(q.handle(), C.ptrdiff_t(i), C.char(c))
	runtime.KeepAlive(q)
}

// InsertCount — Qt: QByteArray::insert(qsizetype, qsizetype, char).
func (q *QByteArray) InsertCount(i, count int, c byte) {
	if i < 0 || i > q.Size() {
		panic("qbytearray: InsertCount out of range")
	}
	C.qbytearray_insert_count(q.handle(), C.ptrdiff_t(i), C.ptrdiff_t(count), C.char(c))
	runtime.KeepAlive(q)
}

// ---- remove/replace ------------------------------------------------------

// Remove — Qt: QByteArray::remove(qsizetype, qsizetype).
func (q *QByteArray) Remove(i, n int) {
	if i < 0 || n < 0 || i+n > q.Size() {
		panic("qbytearray: Remove out of range")
	}
	C.qbytearray_remove(q.handle(), C.ptrdiff_t(i), C.ptrdiff_t(n))
	runtime.KeepAlive(q)
}

// RemoveAt — Qt: QByteArray::removeAt(qsizetype).
func (q *QByteArray) RemoveAt(pos int) {
	C.qbytearray_remove_at(q.handle(), C.ptrdiff_t(pos))
	runtime.KeepAlive(q)
}

// RemoveFirst — Qt: QByteArray::removeFirst().
func (q *QByteArray) RemoveFirst() {
	C.qbytearray_remove_first(q.handle())
	runtime.KeepAlive(q)
}

// RemoveLast — Qt: QByteArray::removeLast().
func (q *QByteArray) RemoveLast() {
	C.qbytearray_remove_last(q.handle())
	runtime.KeepAlive(q)
}

// ReplaceChar — Qt: QByteArray::replace(char, char).
func (q *QByteArray) ReplaceChar(before, after byte) {
	C.qbytearray_replace_char_char(q.handle(), C.char(before), C.char(after))
	runtime.KeepAlive(q)
}

// ReplaceRange — Qt: QByteArray::replace(qsizetype, qsizetype, QByteArrayView).
func (q *QByteArray) ReplaceRange(i, n int, with []byte) {
	if i < 0 || n < 0 || i+n > q.Size() {
		panic("qbytearray: ReplaceRange out of range")
	}
	var p *C.char
	if len(with) > 0 {
		p = (*C.char)(unsafe.Pointer(&with[0]))
	}
	C.qbytearray_replace_range(q.handle(), C.ptrdiff_t(i), C.ptrdiff_t(n), p, C.ptrdiff_t(len(with)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(with)
}

// ReplaceAll — Qt: QByteArray::replace(QByteArrayView, QByteArrayView).
func (q *QByteArray) ReplaceAll(before, after []byte) {
	var pb, pa *C.char
	if len(before) > 0 {
		pb = (*C.char)(unsafe.Pointer(&before[0]))
	}
	if len(after) > 0 {
		pa = (*C.char)(unsafe.Pointer(&after[0]))
	}
	C.qbytearray_replace_data_data(q.handle(), pb, C.ptrdiff_t(len(before)), pa, C.ptrdiff_t(len(after)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(before)
	runtime.KeepAlive(after)
}

// ---- numeric conversion --------------------------------------------------

// ToShort — Qt: QByteArray::toShort(bool*, int) const.
func (q *QByteArray) ToShort(base int) (int16, bool) {
	var ok C.bool
	v := C.qbytearray_to_short(q.handle(), &ok, C.int(base))
	runtime.KeepAlive(q)
	return int16(v), bool(ok)
}

// ToUShort — Qt: QByteArray::toUShort(bool*, int) const.
func (q *QByteArray) ToUShort(base int) (uint16, bool) {
	var ok C.bool
	v := C.qbytearray_to_ushort(q.handle(), &ok, C.int(base))
	runtime.KeepAlive(q)
	return uint16(v), bool(ok)
}

// ToInt — Qt: QByteArray::toInt(bool*, int) const.
func (q *QByteArray) ToInt(base int) (int32, bool) {
	var ok C.bool
	v := C.qbytearray_to_int(q.handle(), &ok, C.int(base))
	runtime.KeepAlive(q)
	return int32(v), bool(ok)
}

// ToUInt — Qt: QByteArray::toUInt(bool*, int) const.
func (q *QByteArray) ToUInt(base int) (uint32, bool) {
	var ok C.bool
	v := C.qbytearray_to_uint(q.handle(), &ok, C.int(base))
	runtime.KeepAlive(q)
	return uint32(v), bool(ok)
}

// ToLongLong — Qt: QByteArray::toLongLong(bool*, int) const.
func (q *QByteArray) ToLongLong(base int) (int64, bool) {
	var ok C.bool
	v := C.qbytearray_to_longlong(q.handle(), &ok, C.int(base))
	runtime.KeepAlive(q)
	return int64(v), bool(ok)
}

// ToULongLong — Qt: QByteArray::toULongLong(bool*, int) const.
func (q *QByteArray) ToULongLong(base int) (uint64, bool) {
	var ok C.bool
	v := C.qbytearray_to_ulonglong(q.handle(), &ok, C.int(base))
	runtime.KeepAlive(q)
	return uint64(v), bool(ok)
}

// ToFloat — Qt: QByteArray::toFloat(bool*) const.
func (q *QByteArray) ToFloat() (float32, bool) {
	var ok C.bool
	v := C.qbytearray_to_float(q.handle(), &ok)
	runtime.KeepAlive(q)
	return float32(v), bool(ok)
}

// ToDouble — Qt: QByteArray::toDouble(bool*) const.
func (q *QByteArray) ToDouble() (float64, bool) {
	var ok C.bool
	v := C.qbytearray_to_double(q.handle(), &ok)
	runtime.KeepAlive(q)
	return float64(v), bool(ok)
}

// SetNumLL — Qt: QByteArray::setNum(qlonglong, int).
func (q *QByteArray) SetNumLL(n int64, base int) {
	C.qbytearray_set_num_ll(q.handle(), C.int64_t(n), C.int(base))
	runtime.KeepAlive(q)
}

// SetNumULL — Qt: QByteArray::setNum(qulonglong, int).
func (q *QByteArray) SetNumULL(n uint64, base int) {
	C.qbytearray_set_num_ull(q.handle(), C.uint64_t(n), C.int(base))
	runtime.KeepAlive(q)
}

// SetNumDouble — Qt: QByteArray::setNum(double, char, int).
func (q *QByteArray) SetNumDouble(v float64, format byte, precision int) {
	C.qbytearray_set_num_double(q.handle(), C.double(v), C.char(format), C.int(precision))
	runtime.KeepAlive(q)
}

// ---- encoding ------------------------------------------------------------

// ToBase64 — Qt: QByteArray::toBase64(Base64Options) const.
func (q *QByteArray) ToBase64(options Base64Options) *QByteArray {
	r := wrapQByteArray(unsafe.Pointer(C.qbytearray_to_base64(q.handle(), C.int(options))))
	runtime.KeepAlive(q)
	return r
}

// ToHex — Qt: QByteArray::toHex(char) const.
func (q *QByteArray) ToHex(separator byte) *QByteArray {
	r := wrapQByteArray(unsafe.Pointer(C.qbytearray_to_hex(q.handle(), C.char(separator))))
	runtime.KeepAlive(q)
	return r
}

// ToPercentEncoding returns a percent-encoded QByteArray.
// exclude/include may be nil or empty (treated as empty QByteArray).
// Qt: QByteArray::toPercentEncoding(const QByteArray&, const QByteArray&, char) const.
func (q *QByteArray) ToPercentEncoding(exclude, include []byte, percent byte) *QByteArray {
	exc := NewQByteArrayFromBytes(exclude)
	inc := NewQByteArrayFromBytes(include)
	defer exc.Close()
	defer inc.Close()

	var ep, ip *C.QByteArray
	if bp, ok := exc.tryHandle(); ok {
		ep = bp
	}
	if bp, ok := inc.tryHandle(); ok {
		ip = bp
	}

	r := wrapQByteArray(unsafe.Pointer(
		C.qbytearray_to_percent(q.handle(), ep, ip, C.char(percent))))
	runtime.KeepAlive(q)
	return r
}

// ---- compression ---------------------------------------------------------

// Compress — Qt: QByteArray::compress(int) const &.
func (q *QByteArray) Compress(level int) *QByteArray {
	r := wrapQByteArray(unsafe.Pointer(C.qbytearray_compress(q.handle(), C.int(level))))
	runtime.KeepAlive(q)
	return r
}

// Uncompress — Qt: QByteArray::uncompress() const &.
func (q *QByteArray) Uncompress() *QByteArray {
	r := wrapQByteArray(unsafe.Pointer(C.qbytearray_uncompress(q.handle())))
	runtime.KeepAlive(q)
	return r
}

// ---- split ---------------------------------------------------------------

// Split mirrors QByteArray::split(char), returning []*QByteArray.
// Each element is a separately owned handle.
// Qt: QByteArray::split(char) const.
func (q *QByteArray) Split(sep byte) []*QByteArray {
	var n C.ptrdiff_t
	arr := C.qbytearray_split(q.handle(), C.char(sep), &n)
	runtime.KeepAlive(q)
	if arr == nil || n == 0 {
		return nil
	}
	defer C.qbytearray_split_free_array(arr)

	cSlice := unsafe.Slice(arr, int(n))
	out := make([]*QByteArray, int(n))
	for i, p := range cSlice {
		out[i] = wrapQByteArray(unsafe.Pointer(p))
	}
	return out
}
