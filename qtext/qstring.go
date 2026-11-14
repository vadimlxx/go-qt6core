package qtext

/*
 Краткое описание реализации:

   QString — handle над C++ QString; владение через finalizer + Close().
   Мутирующие методы идут напрямую в C++ (COW detach — внутри Qt). Строки
   Go и []byte передаются как UTF-8 через unsafe.StringData / &b[0] без
   промежуточной копии; runtime.KeepAlive гарантирует, что аргумент не
   будет собран до возврата из C-вызова.

   В Go нет перегрузок, поэтому для разных C++-перегрузок одной и той же
   функции используются суффиксы:
     - QChar, QString, String — тип аргумента (QChar / QString / Go string);
     - Int, UInt, LL, ULL, Float, Double — перегрузки по целочисленному типу;
     - CS — перегрузка с Qt::CaseSensitivity;
     - WithMatch — вариант с out-параметром QRegularExpressionMatch*;
     - Range — replace(qsizetype, qsizetype, ...).
   Каждый метод снабжён комментарием // Qt: ... , показывающим точный
   C++-прототип, которому он соответствует.

   Реализовано:

     - arg: ArgLl, ArgUll, ArgDouble, ArgQChar, ArgQString, ArgString.
     - compare: CompareQString, CompareQChar, CompareString, LocaleAwareCompare.
     - normalization: Normalized, NormalizedWithVersion.
     - split: SplitQString, SplitQChar, SplitQRegularExpression.
     - section: SectionQChar, SectionQString, SectionQRegularExpression.
     - escape: ToHtmlEscaped.
     - repeat: Repeated.
     - low-level: Capacity, IsSharedWith, ResizeForOverwrite.
     - remove с cs: RemoveQChar, RemoveQString, RemoveString,
       RemoveQRegularExpression.
     - replace с cs: ReplaceQCharQCharCS, ReplaceQCharQStringCS,
       ReplaceQStringQStringCS, ReplaceStringStringCS;
       ReplaceQRegularExpression, ReplaceQRegularExpressionString.
     - поиск с регуляркой: IndexOfQRegularExpression(+WithMatch),
       LastIndexOfQRegularExpression(+WithMatch),
       ContainsQRegularExpression(+WithMatch), CountQRegularExpression.
     - setNum/number для int/uint/float: SetNumInt, SetNumUInt, SetNumFloat,
       QStringNumberInt, QStringNumberUInt.
     - FromLocal8Bit constructor.
     - ToUcs4 — копия []uint32.

   Владение split-результатом: C++ выделяет внешний массив QString* и
   каждый QString* отдельно. Go оборачивает каждый внутренний указатель
   в *QString с финализатором; внешний массив освобождается через
   qstring_free_string_array. Элементы []*QString живут независимо.

   SplitBehavior берётся из qdata (Qt::SplitBehavior — общий тип для
   QString/QByteArray/QStringView).

 Что сознательно не вошло (пока):

   - arg() с цепочкой %1..%99 из одного вызова — используйте явные
     ArgLl/ArgQString/...
   - Vasprintf / Asprintf / Sprintf — не идиоматичны в Go.
   - toWCharArray / fromWCharArray — wchar_t разной ширины.
   - toStdString / toStdU16String / toStdU32String — нужен std::string
     interop (отдельный cgo-слой).
   - toPercentEncoding / fromPercentEncoding.
   - Reverse-iterators (rbegin/rend/crbegin/crend) — неприменимы в Go.
   - QChar*-конструкторы и setRawData/fromRawData — опасны для GC.
   - QLatin1StringView-overloads compare/startsWith/… — добавляются
     по запросу; сейчас все строковые параметры принимаются как Go string.
   - operator+ / operator+= — заменяются Clone().Append*().
   - Null-конструктор QString(nullptr) — NewQString() даёт empty non-null.
*/

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC -Wno-sfinae-incomplete
#cgo LDFLAGS: -lstdc++
#include <stdint.h>
#include <stdbool.h>
#include <stddef.h>
#include "qstring_c.h"
*/
import "C"

import (
	"runtime"
	"unsafe"

	"github.com/vadimlxx/go-qt6core/qdata"
)

// QString wraps a Qt6 QString (UTF-16 COW string).
type QString struct {
	p *C.QString
}

// wrapQString takes ownership of the raw C++ QString pointer.
func wrapQString(p unsafe.Pointer) *QString {
	if p == nil {
		return nil
	}
	q := &QString{p: (*C.QString)(p)}
	runtime.SetFinalizer(q, finalizeQString)
	return q
}

// finalizeQString is the runtime finalizer.
func finalizeQString(q *QString) {
	if q.p != nil {
		C.qstring_delete(q.p)
		q.p = nil
	}
}

// Close releases the underlying C++ object immediately and disables the finalizer.
// Qt: QString::~QString() — explicit Close maps to the destructor call.
func (q *QString) Close() {
	if q.p != nil {
		C.qstring_delete(q.p)
		q.p = nil
		runtime.SetFinalizer(q, nil)
	}
}

// ---- construction --------------------------------------------------------

// NewEmptyQString returns an empty, non-null QString.
// Qt: QString() — default constructor.
func NewEmptyQString() *QString {
	return wrapQString(unsafe.Pointer(C.qstring_new_empty()))
}

// QStringFromString decodes s as UTF-8 into a new QString.
// Qt: QString::fromUtf8(const char *, qsizetype) (also reachable via the implicit
//     QString(const char*) when QT_ASCII_CAST is enabled).
func QStringFromString(s string) *QString {
	var p *C.char
	if len(s) > 0 {
		p = (*C.char)(unsafe.Pointer(unsafe.StringData(s)))
	}
	q := wrapQString(unsafe.Pointer(C.qstring_new_from_utf8(p, C.ptrdiff_t(len(s)))))
	runtime.KeepAlive(s)
	return q
}

// QStringFromUtf8 decodes b as UTF-8 into a new QString.
// Qt: QString::fromUtf8(QByteArrayView).
func QStringFromUtf8(b []byte) *QString {
	var p *C.char
	if len(b) > 0 {
		p = (*C.char)(unsafe.Pointer(&b[0]))
	}
	q := wrapQString(unsafe.Pointer(C.qstring_new_from_utf8(p, C.ptrdiff_t(len(b)))))
	runtime.KeepAlive(b)
	return q
}

// QStringFromUtf16 constructs a QString from UTF-16 code units.
// Qt: QString::fromUtf16(const char16_t *, qsizetype).
func QStringFromUtf16(units []uint16) *QString {
	var p *C.uint16_t
	if len(units) > 0 {
		p = (*C.uint16_t)(unsafe.Pointer(&units[0]))
	}
	q := wrapQString(unsafe.Pointer(C.qstring_new_from_utf16(p, C.ptrdiff_t(len(units)))))
	runtime.KeepAlive(units)
	return q
}

// QStringFromLatin1 decodes b as Latin-1 into a new QString.
// Qt: QString::fromLatin1(QByteArrayView).
func QStringFromLatin1(b []byte) *QString {
	var p *C.char
	if len(b) > 0 {
		p = (*C.char)(unsafe.Pointer(&b[0]))
	}
	q := wrapQString(unsafe.Pointer(C.qstring_new_from_latin1(p, C.ptrdiff_t(len(b)))))
	runtime.KeepAlive(b)
	return q
}

// QStringFromLocal8Bit decodes b using the local 8-bit encoding.
// Qt: QString::fromLocal8Bit(QByteArrayView).
func QStringFromLocal8Bit(b []byte) *QString {
	var p *C.char
	if len(b) > 0 {
		p = (*C.char)(unsafe.Pointer(&b[0]))
	}
	q := wrapQString(unsafe.Pointer(C.qstring_new_from_local8bit(p, C.ptrdiff_t(len(b)))))
	runtime.KeepAlive(b)
	return q
}

// QStringFromQByteArray decodes ba as UTF-8.
// Qt: QString::fromUtf8(const QByteArray&) (also reachable via QString(const QByteArray&)).
func QStringFromQByteArray(ba *QByteArray) *QString {
	if ba == nil || ba.p == nil {
		return NewEmptyQString()
	}
	q := wrapQString(unsafe.Pointer(C.qstring_new_from_qbytearray(ba.p)))
	runtime.KeepAlive(ba)
	return q
}

// QStringFromQStringView copies v into a new QString.
// Qt: QString(QStringView) — explicit constructor.
func QStringFromQStringView(v QStringView) *QString {
	if v.IsNull() {
		return NewEmptyQString()
	}
	q := wrapQString(unsafe.Pointer(C.qstring_new_from_utf16(
		(*C.uint16_t)(v.Utf16Unsafe()), C.ptrdiff_t(v.Size()))))
	runtime.KeepAlive(v)
	return q
}

// NewQString is a convenience dispatcher that picks the matching constructor
// based on the concrete argument type.
// Qt: no single counterpart — dispatches to QString::fromUtf8 / fromUtf16 /
//     fromLatin1 / QString(QStringView) etc. based on the Go argument type.
func NewQString(arg ...any) *QString {
	if len(arg) == 0 || arg[0] == nil {
		return NewEmptyQString()
	}
	switch v := arg[0].(type) {
	case string:
		return QStringFromString(v)
	case []byte:
		return QStringFromUtf8(v)
	case []uint16:
		return QStringFromUtf16(v)
	case QByteArray:
		return QStringFromQByteArray(&v)
	case *QByteArray:
		return QStringFromQByteArray(v)
	case QStringView:
		return QStringFromQStringView(v)
	}
	panic("NewQString invalid argument")
}

// Clone returns a shallow COW copy.
// Qt: QString(const QString&) — copy constructor.
func (q *QString) Clone() *QString {
	return wrapQString(unsafe.Pointer(C.qstring_new_copy(q.p)))
}

// ---- accessors -----------------------------------------------------------

// Size returns the number of UTF-16 code units.
// Qt: QString::size() const.
func (q *QString) Size() int { return int(C.qstring_size(q.p)) }

// Len is an alias for Size.
// Qt: QString::length() const (which forwards to size()).
func (q *QString) Len() int { return q.Size() }

// Capacity returns the allocated capacity in code units.
// Qt: QString::capacity() const.
func (q *QString) Capacity() int { return int(C.qstring_capacity(q.p)) }

// IsEmpty reports whether the string has zero length.
// Qt: QString::isEmpty() const.
func (q *QString) IsEmpty() bool { return bool(C.qstring_is_empty(q.p)) }

// IsNull reports whether the string is null (default-constructed or from a null pointer).
// Qt: QString::isNull() const.
func (q *QString) IsNull() bool { return bool(C.qstring_is_null(q.p)) }

// IsDetached reports whether the string is not shared with any other QString.
// Qt: QString::isDetached() const.
func (q *QString) IsDetached() bool { return bool(C.qstring_is_detached(q.p)) }

// Detach ensures the string is not shared with any other QString.
// Qt: QString::detach().
func (q *QString) Detach() { C.qstring_detach(q.p) }

// IsSharedWith reports whether q and other share the same data.
// Qt: QString::isSharedWith(const QString&) const.
func (q *QString) IsSharedWith(other *QString) bool {
	r := bool(C.qstring_is_shared_with(q.p, other.p))
	runtime.KeepAlive(other)
	return r
}

// Utf16Unsafe returns the raw internal UTF-16 data pointer.
// Qt: QString::utf16() const — no copy.
// Caller must ensure the string is neither mutated nor Close()d while the pointer is in use.
func (q *QString) Utf16Unsafe() unsafe.Pointer {
	return unsafe.Pointer(C.qstring_utf16(q.p))
}

// Utf16 returns a copy of the contents as UTF-16 code units.
// Qt: QString::utf16() const — with a copy into a Go slice.
func (q *QString) Utf16() []uint16 {
	n := q.Size()
	if n == 0 {
		return nil
	}
	p := C.qstring_utf16(q.p)
	out := make([]uint16, n)
	copy(out, unsafe.Slice((*uint16)(unsafe.Pointer(p)), n))
	runtime.KeepAlive(q)
	return out
}

// GoString decodes the contents as UTF-8 into a Go string.
// Qt: no direct counterpart — implemented via QString::toUtf8().
func (q *QString) GoString() string {
	ba := q.Utf8()
	defer ba.Close()
	return ba.String()
}

// String implements fmt.Stringer.
// Qt: no direct counterpart — same as GoString.
func (q *QString) String() string { return q.GoString() }

// Utf8 returns a fresh *QByteArray with UTF-8 encoding.
// Qt: QString::toUtf8() const &.
func (q *QString) Utf8() *QByteArray {
	var out *C.QByteArray
	C.qstring_to_utf8_qba(q.p, &out)
	runtime.KeepAlive(q)
	return wrapQByteArray(unsafe.Pointer(out))
}

// Latin1 returns a fresh *QByteArray with Latin-1 encoding.
// Qt: QString::toLatin1() const &.
func (q *QString) Latin1() *QByteArray {
	var out *C.QByteArray
	C.qstring_to_latin1_qba(q.p, &out)
	runtime.KeepAlive(q)
	return wrapQByteArray(unsafe.Pointer(out))
}

// Local8Bit returns a fresh *QByteArray with the local 8-bit encoding.
// Qt: QString::toLocal8Bit() const &.
func (q *QString) Local8Bit() *QByteArray {
	var out *C.QByteArray
	C.qstring_to_local8bit_qba(q.p, &out)
	runtime.KeepAlive(q)
	return wrapQByteArray(unsafe.Pointer(out))
}

// At returns the i-th QChar. Panics if i is out of range.
// Qt: QString::at(qsizetype) const.
func (q *QString) At(i int) QChar {
	if i < 0 || i >= q.Size() {
		panic("qstring: At out of range")
	}
	r := C.qstring_at(q.p, C.ptrdiff_t(i))
	runtime.KeepAlive(q)
	return NewQChar(uint16(r))
}

// Front returns the first QChar. Panics if empty.
// Qt: QString::front() const.
func (q *QString) Front() QChar {
	if q.IsEmpty() {
		panic("qstring: Front on empty QString")
	}
	r := C.qstring_front(q.p)
	runtime.KeepAlive(q)
	return NewQChar(uint16(r))
}

// Back returns the last QChar. Panics if empty.
// Qt: QString::back() const.
func (q *QString) Back() QChar {
	if q.IsEmpty() {
		panic("qstring: Back on empty QString")
	}
	r := C.qstring_back(q.p)
	runtime.KeepAlive(q)
	return NewQChar(uint16(r))
}

// ToQStringView returns a zero-copy QStringView over the internal data.
// Qt: QStringView(const QString&) — implicit conversion via the QStringView constructor.
func (q *QString) ToQStringView() QStringView {
	v := newQStringViewWithOwner(q.Utf16Unsafe(), q.Size(), q)
	runtime.KeepAlive(q)
	return v
}

// ---- modifiers -----------------------------------------------------------

// Clear makes the string null.
// Qt: QString::clear().
func (q *QString) Clear() { C.qstring_clear(q.p) }

// Resize changes the string length, content of the new tail is undefined.
// Qt: QString::resize(qsizetype).
func (q *QString) Resize(n int) { C.qstring_resize(q.p, C.ptrdiff_t(n)) }

// ResizeFill changes the string length, filling new code units with c.
// Qt: QString::resize(qsizetype, QChar).
func (q *QString) ResizeFill(n int, c QChar) {
	C.qstring_resize_fill(q.p, C.ptrdiff_t(n), C.uint16_t(c.Unicode()))
}

// ResizeForOverwrite changes the length without initializing the tail.
// Qt: QString::resizeForOverwrite(qsizetype).
func (q *QString) ResizeForOverwrite(n int) { C.qstring_resize_for_overwrite(q.p, C.ptrdiff_t(n)) }

// Reserve pre-allocates storage for at least n code units.
// Qt: QString::reserve(qsizetype).
func (q *QString) Reserve(n int) { C.qstring_reserve(q.p, C.ptrdiff_t(n)) }

// Squeeze releases unused capacity.
// Qt: QString::squeeze().
func (q *QString) Squeeze() { C.qstring_squeeze(q.p) }

// Truncate cuts the string to length pos.
// Qt: QString::truncate(qsizetype).
func (q *QString) Truncate(pos int) { C.qstring_truncate(q.p, C.ptrdiff_t(pos)) }

// Chop removes n code units from the end.
// Qt: QString::chop(qsizetype).
func (q *QString) Chop(n int) { C.qstring_chop(q.p, C.ptrdiff_t(n)) }

// Fill replaces contents with n copies of c. n < 0 means "pad to current size".
// Qt: QString::fill(QChar, qsizetype).
func (q *QString) Fill(c QChar, n int) {
	C.qstring_fill(q.p, C.uint16_t(c.Unicode()), C.ptrdiff_t(n))
}

// Append appends a single character.
// Qt: QString::append(QChar).
func (q *QString) Append(c QChar) {
	C.qstring_append_qchar(q.p, C.uint16_t(c.Unicode()))
}

// AppendQString appends another QString.
// Qt: QString::append(const QString&).
func (q *QString) AppendQString(other *QString) {
	C.qstring_append_qstring(q.p, other.p)
	runtime.KeepAlive(other)
}

// AppendString appends a Go string (UTF-8).
// Qt: QString::append(const QString&) (UTF-8 decoded; Qt also exposes append(QUtf8StringView)).
func (q *QString) AppendString(s string) {
	var p *C.char
	if len(s) > 0 {
		p = (*C.char)(unsafe.Pointer(unsafe.StringData(s)))
	}
	C.qstring_append_utf8(q.p, p, C.ptrdiff_t(len(s)))
	runtime.KeepAlive(s)
}

// Prepend prepends a single character.
// Qt: QString::prepend(QChar).
func (q *QString) Prepend(c QChar) {
	C.qstring_prepend_qchar(q.p, C.uint16_t(c.Unicode()))
}

// PrependQString prepends another QString.
// Qt: QString::prepend(const QString&).
func (q *QString) PrependQString(other *QString) {
	C.qstring_prepend_qstring(q.p, other.p)
	runtime.KeepAlive(other)
}

// PrependString prepends a Go string (UTF-8).
// Qt: QString::prepend(const QString&) (UTF-8 decoded; Qt also exposes prepend(QUtf8StringView)).
func (q *QString) PrependString(s string) {
	var p *C.char
	if len(s) > 0 {
		p = (*C.char)(unsafe.Pointer(unsafe.StringData(s)))
	}
	C.qstring_prepend_utf8(q.p, p, C.ptrdiff_t(len(s)))
	runtime.KeepAlive(s)
}

// Insert inserts a single character at position i.
// Qt: QString::insert(qsizetype, QChar).
func (q *QString) Insert(i int, c QChar) {
	C.qstring_insert_qchar(q.p, C.ptrdiff_t(i), C.uint16_t(c.Unicode()))
}

// InsertQString inserts another QString at position i.
// Qt: QString::insert(qsizetype, const QString&).
func (q *QString) InsertQString(i int, other *QString) {
	C.qstring_insert_qstring(q.p, C.ptrdiff_t(i), other.p)
	runtime.KeepAlive(other)
}

// InsertString inserts a Go string (UTF-8) at position i.
// Qt: QString::insert(qsizetype, const QString&) (UTF-8 decoded; Qt also exposes insert(qsizetype, QUtf8StringView)).
func (q *QString) InsertString(i int, s string) {
	var p *C.char
	if len(s) > 0 {
		p = (*C.char)(unsafe.Pointer(unsafe.StringData(s)))
	}
	C.qstring_insert_utf8(q.p, C.ptrdiff_t(i), p, C.ptrdiff_t(len(s)))
	runtime.KeepAlive(s)
}

// Remove removes n code units starting at i.
// Qt: QString::remove(qsizetype, qsizetype).
func (q *QString) Remove(i, n int) { C.qstring_remove(q.p, C.ptrdiff_t(i), C.ptrdiff_t(n)) }

// RemoveAt removes the code unit at position i.
// Qt: QString::removeAt(qsizetype).
func (q *QString) RemoveAt(i int) { C.qstring_remove_at(q.p, C.ptrdiff_t(i)) }

// RemoveFirst removes the first code unit (no-op if empty).
// Qt: QString::removeFirst().
func (q *QString) RemoveFirst() { C.qstring_remove_first(q.p) }

// RemoveLast removes the last code unit (no-op if empty).
// Qt: QString::removeLast().
func (q *QString) RemoveLast() { C.qstring_remove_last(q.p) }

// RemoveQChar removes all occurrences of c with the given case sensitivity.
// Qt: QString::remove(QChar, Qt::CaseSensitivity).
func (q *QString) RemoveQChar(c QChar, caseSensitive bool) {
	C.qstring_remove_qchar(q.p, C.uint16_t(c.Unicode()), csInt(caseSensitive))
}

// RemoveQString removes all occurrences of another QString.
// Qt: QString::remove(const QString&, Qt::CaseSensitivity).
func (q *QString) RemoveQString(other *QString, caseSensitive bool) {
	C.qstring_remove_qstring(q.p, other.p, csInt(caseSensitive))
	runtime.KeepAlive(other)
}

// RemoveString removes all occurrences of a Go string (UTF-8 decoded).
// Qt: QString::remove(const QString&, Qt::CaseSensitivity) (UTF-8 decoded; Qt also exposes remove(QLatin1StringView, cs)).
func (q *QString) RemoveString(s string, caseSensitive bool) {
	var p *C.char
	if len(s) > 0 {
		p = (*C.char)(unsafe.Pointer(unsafe.StringData(s)))
	}
	C.qstring_remove_utf8(q.p, p, C.ptrdiff_t(len(s)), csInt(caseSensitive))
	runtime.KeepAlive(s)
}

// ReplaceQCharQChar replaces a single character with another.
// Qt: QString::replace(QChar, QChar).
func (q *QString) ReplaceQCharQChar(before, after QChar) {
	C.qstring_replace_qchar_qchar(q.p, C.uint16_t(before.Unicode()), C.uint16_t(after.Unicode()))
}

// ReplaceQCharQString replaces a single character with a QString.
// Qt: QString::replace(QChar, const QString&).
func (q *QString) ReplaceQCharQString(before QChar, after *QString) {
	C.qstring_replace_qchar_qstring(q.p, C.uint16_t(before.Unicode()), after.p)
	runtime.KeepAlive(after)
}

// ReplaceQCharString replaces a single character with a Go string (UTF-8).
// Qt: QString::replace(QChar, const QString&) (UTF-8 decoded; Qt also exposes replace(QChar, QLatin1StringView)).
func (q *QString) ReplaceQCharString(before QChar, after string) {
	var p *C.char
	if len(after) > 0 {
		p = (*C.char)(unsafe.Pointer(unsafe.StringData(after)))
	}
	C.qstring_replace_qchar_utf8(q.p, C.uint16_t(before.Unicode()), p, C.ptrdiff_t(len(after)))
	runtime.KeepAlive(after)
}

// ReplaceQStringQChar replaces a QString with a single QChar.
// Qt: QString::replace(const QString&, const QString&) — resolves via implicit QChar→QString conversion;
//     Qt has no direct (const QString&, QChar) overload.
func (q *QString) ReplaceQStringQChar(before *QString, after QChar) {
	C.qstring_replace_qstring_qchar(q.p, before.p, C.uint16_t(after.Unicode()))
	runtime.KeepAlive(before)
}

// ReplaceQStringQString replaces one QString with another.
// Qt: QString::replace(const QString&, const QString&).
func (q *QString) ReplaceQStringQString(before, after *QString) {
	C.qstring_replace_qstring_qstring(q.p, before.p, after.p)
	runtime.KeepAlive(before)
	runtime.KeepAlive(after)
}

// ReplaceQStringString replaces a QString with a Go string (UTF-8).
// Qt: QString::replace(const QString&, const QString&) (UTF-8 decoded; Qt also exposes replace(const QString&, QLatin1StringView)).
func (q *QString) ReplaceQStringString(before *QString, after string) {
	var p *C.char
	if len(after) > 0 {
		p = (*C.char)(unsafe.Pointer(unsafe.StringData(after)))
	}
	C.qstring_replace_qstring_utf8(q.p, before.p, p, C.ptrdiff_t(len(after)))
	runtime.KeepAlive(before)
	runtime.KeepAlive(after)
}

// ReplaceStringQChar replaces a Go string (UTF-8) with a single QChar.
// Qt: QString::replace(const QString&, const QString&) — resolves via implicit QChar→QString conversion;
//     Qt has no direct (const QString&, QChar) overload.
func (q *QString) ReplaceStringQChar(before string, after QChar) {
	var p *C.char
	if len(before) > 0 {
		p = (*C.char)(unsafe.Pointer(unsafe.StringData(before)))
	}
	C.qstring_replace_utf8_qchar(q.p, p, C.ptrdiff_t(len(before)), C.uint16_t(after.Unicode()))
	runtime.KeepAlive(before)
}

// ReplaceStringQString replaces a Go string (UTF-8) with a QString.
// Qt: QString::replace(const QString&, const QString&) (UTF-8 decoded; Qt also exposes replace(QLatin1StringView, const QString&)).
func (q *QString) ReplaceStringQString(before string, after *QString) {
	var p *C.char
	if len(before) > 0 {
		p = (*C.char)(unsafe.Pointer(unsafe.StringData(before)))
	}
	C.qstring_replace_utf8_qstring(q.p, p, C.ptrdiff_t(len(before)), after.p)
	runtime.KeepAlive(before)
	runtime.KeepAlive(after)
}

// ReplaceStringString replaces one Go string (UTF-8) with another.
// Qt: QString::replace(const QString&, const QString&) (UTF-8 decoded; Qt also exposes replace(QLatin1StringView, QLatin1StringView)).
func (q *QString) ReplaceStringString(before, after string) {
	var pb, pa *C.char
	if len(before) > 0 { pb = (*C.char)(unsafe.Pointer(unsafe.StringData(before))) }
	if len(after) > 0  { pa = (*C.char)(unsafe.Pointer(unsafe.StringData(after)))  }
	C.qstring_replace_utf8_utf8(q.p, pb, C.ptrdiff_t(len(before)), pa, C.ptrdiff_t(len(after)))
	runtime.KeepAlive(before)
	runtime.KeepAlive(after)
}

// ReplaceRangeString replaces n code units at position i with a Go string (UTF-8).
// Qt: QString::replace(qsizetype, qsizetype, const QString&) (UTF-8 decoded).
func (q *QString) ReplaceRangeString(i, n int, after string) {
	var p *C.char
	if len(after) > 0 {
		p = (*C.char)(unsafe.Pointer(unsafe.StringData(after)))
	}
	C.qstring_replace_range_utf8(q.p, C.ptrdiff_t(i), C.ptrdiff_t(n), p, C.ptrdiff_t(len(after)))
	runtime.KeepAlive(after)
}

// ReplaceQCharQCharCS replaces a single character with another, with case sensitivity.
// Qt: QString::replace(QChar, QChar, Qt::CaseSensitivity).
func (q *QString) ReplaceQCharQCharCS(before, after QChar, caseSensitive bool) {
	C.qstring_replace_qchar_qchar_cs(q.p, C.uint16_t(before.Unicode()),
		C.uint16_t(after.Unicode()), csInt(caseSensitive))
}

// ReplaceQCharQStringCS replaces a single character with a QString, with case sensitivity.
// Qt: QString::replace(QChar, const QString&, Qt::CaseSensitivity).
func (q *QString) ReplaceQCharQStringCS(before QChar, after *QString, caseSensitive bool) {
	C.qstring_replace_qchar_qstring_cs(q.p, C.uint16_t(before.Unicode()), after.p,
		csInt(caseSensitive))
	runtime.KeepAlive(after)
}

// ReplaceQStringQStringCS replaces one QString with another, with case sensitivity.
// Qt: QString::replace(const QString&, const QString&, Qt::CaseSensitivity).
func (q *QString) ReplaceQStringQStringCS(before, after *QString, caseSensitive bool) {
	C.qstring_replace_qstring_qstring_cs(q.p, before.p, after.p, csInt(caseSensitive))
	runtime.KeepAlive(before)
	runtime.KeepAlive(after)
}

// ReplaceStringStringCS replaces one Go string (UTF-8) with another, with case sensitivity.
// Qt: QString::replace(const QString&, const QString&, Qt::CaseSensitivity) (UTF-8 decoded).
func (q *QString) ReplaceStringStringCS(before, after string, caseSensitive bool) {
	var pb, pa *C.char
	if len(before) > 0 { pb = (*C.char)(unsafe.Pointer(unsafe.StringData(before))) }
	if len(after) > 0  { pa = (*C.char)(unsafe.Pointer(unsafe.StringData(after)))  }
	C.qstring_replace_utf8_utf8_cs(q.p, pb, C.ptrdiff_t(len(before)), pa,
		C.ptrdiff_t(len(after)), csInt(caseSensitive))
	runtime.KeepAlive(before)
	runtime.KeepAlive(after)
}

// ---- slicing (return new owned QString) ---------------------------------

// Left returns a new QString with up to n leading code units.
// Qt: QString::left(qsizetype) const &.
func (q *QString) Left(n int) *QString {
	return wrapQString(unsafe.Pointer(C.qstring_left(q.p, C.ptrdiff_t(n))))
}

// Right returns a new QString with up to n trailing code units.
// Qt: QString::right(qsizetype) const &.
func (q *QString) Right(n int) *QString {
	return wrapQString(unsafe.Pointer(C.qstring_right(q.p, C.ptrdiff_t(n))))
}

// Mid returns a new QString starting at pos with n code units (n < 0 means "to end").
// Qt: QString::mid(qsizetype, qsizetype) const &.
func (q *QString) Mid(pos, n int) *QString {
	return wrapQString(unsafe.Pointer(C.qstring_mid(q.p, C.ptrdiff_t(pos), C.ptrdiff_t(n))))
}

// First returns a new QString with the first n code units. Panics if n is out of [0, Size()].
// Qt: QString::first(qsizetype) const &.
func (q *QString) First(n int) *QString {
	if n < 0 || n > q.Size() {
		panic("qstring: First out of range")
	}
	return wrapQString(unsafe.Pointer(C.qstring_first(q.p, C.ptrdiff_t(n))))
}

// Last returns a new QString with the last n code units. Panics if n is out of [0, Size()].
// Qt: QString::last(qsizetype) const &.
func (q *QString) Last(n int) *QString {
	if n < 0 || n > q.Size() {
		panic("qstring: Last out of range")
	}
	return wrapQString(unsafe.Pointer(C.qstring_last(q.p, C.ptrdiff_t(n))))
}

// Sliced returns a new QString from pos to the end. Panics if pos is out of [0, Size()].
// Qt: QString::sliced(qsizetype) const &.
func (q *QString) Sliced(pos int) *QString {
	if pos < 0 || pos > q.Size() {
		panic("qstring: Sliced out of range")
	}
	return wrapQString(unsafe.Pointer(C.qstring_sliced(q.p, C.ptrdiff_t(pos))))
}

// SlicedN returns a new QString of n code units starting at pos.
// Panics if pos or n is out of range.
// Qt: QString::sliced(qsizetype, qsizetype) const &.
func (q *QString) SlicedN(pos, n int) *QString {
	if pos < 0 || n < 0 || pos+n > q.Size() {
		panic("qstring: SlicedN out of range")
	}
	return wrapQString(unsafe.Pointer(C.qstring_sliced_n(q.p, C.ptrdiff_t(pos), C.ptrdiff_t(n))))
}

// Chopped returns a new QString with the last n code units removed.
// Panics if n is out of [0, Size()].
// Qt: QString::chopped(qsizetype) const &.
func (q *QString) Chopped(n int) *QString {
	if n < 0 || n > q.Size() {
		panic("qstring: Chopped out of range")
	}
	return wrapQString(unsafe.Pointer(C.qstring_chopped(q.p, C.ptrdiff_t(n))))
}

// Repeated returns a new QString with contents repeated times times.
// Qt: QString::repeated(qsizetype) const.
func (q *QString) Repeated(times int) *QString {
	if times < 0 {
		panic("qstring: Repeated negative times")
	}
	return wrapQString(unsafe.Pointer(C.qstring_repeated(q.p, C.ptrdiff_t(times))))
}

// ---- search --------------------------------------------------------------

// IndexOfQChar returns the index of the first occurrence of c (or -1).
// Qt: QString::indexOf(QChar, qsizetype, Qt::CaseSensitivity) const.
func (q *QString) IndexOfQChar(c QChar, from int, caseSensitive bool) int {
	return int(C.qstring_index_of_qchar(q.p, C.uint16_t(c.Unicode()),
		C.ptrdiff_t(from), csInt(caseSensitive)))
}

// IndexOfQString returns the index of the first occurrence of other (or -1).
// Qt: QString::indexOf(const QString&, qsizetype, Qt::CaseSensitivity) const.
func (q *QString) IndexOfQString(other *QString, from int, caseSensitive bool) int {
	r := C.qstring_index_of_qstring(q.p, other.p, C.ptrdiff_t(from), csInt(caseSensitive))
	runtime.KeepAlive(other)
	return int(r)
}

// IndexOfString returns the index of the first occurrence of a Go string (or -1).
// Qt: QString::indexOf(const QString&, qsizetype, Qt::CaseSensitivity) const (UTF-8 decoded;
//     Qt also exposes indexOf(QLatin1StringView, from, cs) and indexOf(QStringView, from, cs)).
func (q *QString) IndexOfString(s string, from int, caseSensitive bool) int {
	var p *C.char
	if len(s) > 0 {
		p = (*C.char)(unsafe.Pointer(unsafe.StringData(s)))
	}
	r := C.qstring_index_of_utf8(q.p, p, C.ptrdiff_t(len(s)), C.ptrdiff_t(from), csInt(caseSensitive))
	runtime.KeepAlive(s)
	return int(r)
}

// LastIndexOfQChar returns the index of the last occurrence of c at or before from (or -1).
// Qt: QString::lastIndexOf(QChar, qsizetype, Qt::CaseSensitivity) const.
func (q *QString) LastIndexOfQChar(c QChar, from int, caseSensitive bool) int {
	return int(C.qstring_last_index_of_qchar(q.p, C.uint16_t(c.Unicode()),
		C.ptrdiff_t(from), csInt(caseSensitive)))
}

// LastIndexOfQString returns the index of the last occurrence of other at or before from (or -1).
// Qt: QString::lastIndexOf(const QString&, qsizetype, Qt::CaseSensitivity) const.
func (q *QString) LastIndexOfQString(other *QString, from int, caseSensitive bool) int {
	r := C.qstring_last_index_of_qstring(q.p, other.p, C.ptrdiff_t(from), csInt(caseSensitive))
	runtime.KeepAlive(other)
	return int(r)
}

// LastIndexOfString returns the index of the last occurrence of a Go string at or before from (or -1).
// Qt: QString::lastIndexOf(const QString&, qsizetype, Qt::CaseSensitivity) const (UTF-8 decoded;
//     Qt also exposes lastIndexOf(QLatin1StringView, from, cs)).
func (q *QString) LastIndexOfString(s string, from int, caseSensitive bool) int {
	var p *C.char
	if len(s) > 0 {
		p = (*C.char)(unsafe.Pointer(unsafe.StringData(s)))
	}
	r := C.qstring_last_index_of_utf8(q.p, p, C.ptrdiff_t(len(s)), C.ptrdiff_t(from),
		csInt(caseSensitive))
	runtime.KeepAlive(s)
	return int(r)
}

// ContainsQChar reports whether q contains c.
// Qt: QString::contains(QChar, Qt::CaseSensitivity) const.
func (q *QString) ContainsQChar(c QChar, caseSensitive bool) bool {
	return bool(C.qstring_contains_qchar(q.p, C.uint16_t(c.Unicode()), csInt(caseSensitive)))
}

// ContainsQString reports whether q contains another QString.
// Qt: QString::contains(const QString&, Qt::CaseSensitivity) const.
func (q *QString) ContainsQString(other *QString, caseSensitive bool) bool {
	r := C.qstring_contains_qstring(q.p, other.p, csInt(caseSensitive))
	runtime.KeepAlive(other)
	return bool(r)
}

// ContainsString reports whether q contains a Go string (UTF-8 decoded).
// Qt: QString::contains(const QString&, Qt::CaseSensitivity) const (UTF-8 decoded;
//     Qt also exposes contains(QLatin1StringView, cs) and contains(QStringView, cs)).
func (q *QString) ContainsString(s string, caseSensitive bool) bool {
	var p *C.char
	if len(s) > 0 {
		p = (*C.char)(unsafe.Pointer(unsafe.StringData(s)))
	}
	r := C.qstring_contains_utf8(q.p, p, C.ptrdiff_t(len(s)), csInt(caseSensitive))
	runtime.KeepAlive(s)
	return bool(r)
}

// CountQChar returns the number of non-overlapping occurrences of c.
// Qt: QString::count(QChar, Qt::CaseSensitivity) const.
func (q *QString) CountQChar(c QChar, caseSensitive bool) int {
	return int(C.qstring_count_qchar(q.p, C.uint16_t(c.Unicode()), csInt(caseSensitive)))
}

// CountQString returns the number of non-overlapping occurrences of another QString.
// Qt: QString::count(const QString&, Qt::CaseSensitivity) const.
func (q *QString) CountQString(other *QString, caseSensitive bool) int {
	r := C.qstring_count_qstring(q.p, other.p, csInt(caseSensitive))
	runtime.KeepAlive(other)
	return int(r)
}

// CountString returns the number of non-overlapping occurrences of a Go string (UTF-8).
// Qt: QString::count(const QString&, Qt::CaseSensitivity) const (UTF-8 decoded;
//     Qt also exposes count(QStringView, cs)).
func (q *QString) CountString(s string, caseSensitive bool) int {
	var p *C.char
	if len(s) > 0 {
		p = (*C.char)(unsafe.Pointer(unsafe.StringData(s)))
	}
	r := C.qstring_count_utf8(q.p, p, C.ptrdiff_t(len(s)), csInt(caseSensitive))
	runtime.KeepAlive(s)
	return int(r)
}

// StartsWithQChar reports whether q starts with c.
// Qt: QString::startsWith(QChar, Qt::CaseSensitivity) const.
func (q *QString) StartsWithQChar(c QChar, caseSensitive bool) bool {
	return bool(C.qstring_starts_with_qchar(q.p, C.uint16_t(c.Unicode()), csInt(caseSensitive)))
}

// StartsWithQString reports whether q starts with another QString.
// Qt: QString::startsWith(const QString&, Qt::CaseSensitivity) const.
func (q *QString) StartsWithQString(other *QString, caseSensitive bool) bool {
	r := C.qstring_starts_with_qstring(q.p, other.p, csInt(caseSensitive))
	runtime.KeepAlive(other)
	return bool(r)
}

// StartsWithString reports whether q starts with a Go string (UTF-8 decoded).
// Qt: QString::startsWith(const QString&, Qt::CaseSensitivity) const (UTF-8 decoded;
//     Qt also exposes startsWith(QLatin1StringView, cs) and startsWith(QStringView, cs)).
func (q *QString) StartsWithString(s string, caseSensitive bool) bool {
	var p *C.char
	if len(s) > 0 {
		p = (*C.char)(unsafe.Pointer(unsafe.StringData(s)))
	}
	r := C.qstring_starts_with_utf8(q.p, p, C.ptrdiff_t(len(s)), csInt(caseSensitive))
	runtime.KeepAlive(s)
	return bool(r)
}

// EndsWithQChar reports whether q ends with c.
// Qt: QString::endsWith(QChar, Qt::CaseSensitivity) const.
func (q *QString) EndsWithQChar(c QChar, caseSensitive bool) bool {
	return bool(C.qstring_ends_with_qchar(q.p, C.uint16_t(c.Unicode()), csInt(caseSensitive)))
}

// EndsWithQString reports whether q ends with another QString.
// Qt: QString::endsWith(const QString&, Qt::CaseSensitivity) const.
func (q *QString) EndsWithQString(other *QString, caseSensitive bool) bool {
	r := C.qstring_ends_with_qstring(q.p, other.p, csInt(caseSensitive))
	runtime.KeepAlive(other)
	return bool(r)
}

// EndsWithString reports whether q ends with a Go string (UTF-8 decoded).
// Qt: QString::endsWith(const QString&, Qt::CaseSensitivity) const (UTF-8 decoded;
//     Qt also exposes endsWith(QLatin1StringView, cs) and endsWith(QStringView, cs)).
func (q *QString) EndsWithString(s string, caseSensitive bool) bool {
	var p *C.char
	if len(s) > 0 {
		p = (*C.char)(unsafe.Pointer(unsafe.StringData(s)))
	}
	r := C.qstring_ends_with_utf8(q.p, p, C.ptrdiff_t(len(s)), csInt(caseSensitive))
	runtime.KeepAlive(s)
	return bool(r)
}

// ---- QRegularExpression search ------------------------------------------

// IndexOfQRegularExpression returns the index of the first match of re (or -1).
// Qt: QString::indexOf(const QRegularExpression&, qsizetype, QRegularExpressionMatch*) const — with match=nullptr.
func (q *QString) IndexOfQRegularExpression(re *QRegularExpression, from int) int {
	return q.IndexOfQRegularExpressionWithMatch(re, from, nil)
}

// IndexOfQRegularExpressionWithMatch is like IndexOfQRegularExpression, but
// additionally fills match (if non-nil) with the match result.
// Qt: QString::indexOf(const QRegularExpression&, qsizetype, QRegularExpressionMatch*) const.
func (q *QString) IndexOfQRegularExpressionWithMatch(re *QRegularExpression, from int,
	match *QRegularExpressionMatch) int {
	if re == nil || re.p == nil {
		panic("qstring: IndexOfQRegularExpressionWithMatch nil regexp")
	}
	var cMatch *C.QRegularExpressionMatch
	if match != nil && match.p != nil {
		cMatch = match.p
	}
	r := C.qstring_index_of_qregexp(q.p, re.p, C.ptrdiff_t(from), cMatch)
	runtime.KeepAlive(re)
	runtime.KeepAlive(match)
	return int(r)
}

// LastIndexOfQRegularExpression returns the index of the last match of re at or before from (or -1).
// Qt: QString::lastIndexOf(const QRegularExpression&, qsizetype, QRegularExpressionMatch*) const — with match=nullptr.
func (q *QString) LastIndexOfQRegularExpression(re *QRegularExpression, from int) int {
	return q.LastIndexOfQRegularExpressionWithMatch(re, from, nil)
}

// LastIndexOfQRegularExpressionWithMatch is like LastIndexOfQRegularExpression,
// but additionally fills match (if non-nil) with the match result.
// Qt: QString::lastIndexOf(const QRegularExpression&, qsizetype, QRegularExpressionMatch*) const.
func (q *QString) LastIndexOfQRegularExpressionWithMatch(re *QRegularExpression, from int,
	match *QRegularExpressionMatch) int {
	if re == nil || re.p == nil {
		panic("qstring: LastIndexOfQRegularExpressionWithMatch nil regexp")
	}
	var cMatch *C.QRegularExpressionMatch
	if match != nil && match.p != nil {
		cMatch = match.p
	}
	r := C.qstring_last_index_of_qregexp(q.p, re.p, C.ptrdiff_t(from), cMatch)
	runtime.KeepAlive(re)
	runtime.KeepAlive(match)
	return int(r)
}

// ContainsQRegularExpression reports whether re matches anywhere in q.
// Qt: QString::contains(const QRegularExpression&, QRegularExpressionMatch*) const — with match=nullptr.
func (q *QString) ContainsQRegularExpression(re *QRegularExpression) bool {
	return q.ContainsQRegularExpressionWithMatch(re, nil)
}

// ContainsQRegularExpressionWithMatch is like ContainsQRegularExpression, but
// additionally fills match (if non-nil) with the match result.
// Qt: QString::contains(const QRegularExpression&, QRegularExpressionMatch*) const.
func (q *QString) ContainsQRegularExpressionWithMatch(re *QRegularExpression,
	match *QRegularExpressionMatch) bool {
	if re == nil || re.p == nil {
		panic("qstring: ContainsQRegularExpressionWithMatch nil regexp")
	}
	var cMatch *C.QRegularExpressionMatch
	if match != nil && match.p != nil {
		cMatch = match.p
	}
	r := C.qstring_contains_qregexp(q.p, re.p, cMatch)
	runtime.KeepAlive(re)
	runtime.KeepAlive(match)
	return bool(r)
}

// CountQRegularExpression returns the number of times re matches q.
// Note: Qt counts OVERLAPPING matches: \d+ applied to "123" yields 3
// (matching "123", "23" and "3"). Use a word-bounded pattern such as
// \b\d+\b to get the number of distinct digit runs.
// Qt: QString::count(const QRegularExpression&) const.
func (q *QString) CountQRegularExpression(re *QRegularExpression) int {
	if re == nil || re.p == nil {
		panic("qstring: CountQRegularExpression nil regexp")
	}
	r := C.qstring_count_qregexp(q.p, re.p)
	runtime.KeepAlive(re)
	return int(r)
}

// SectionQRegularExpression returns the section delimited by matches of re,
// from start to end inclusive. end == -1 means "last section".
// Qt: QString::section(const QRegularExpression&, qsizetype, qsizetype, SectionFlags) const.
func (q *QString) SectionQRegularExpression(re *QRegularExpression, start, end int,
	flags SectionFlag) *QString {
	if re == nil || re.p == nil {
		panic("qstring: SectionQRegularExpression nil regexp")
	}
	r := wrapQString(unsafe.Pointer(C.qstring_section_qregexp(q.p, re.p,
		C.ptrdiff_t(start), C.ptrdiff_t(end), C.int(flags))))
	runtime.KeepAlive(re)
	return r
}

// ReplaceQRegularExpression replaces all matches of re with after.
// Qt: QString::replace(const QRegularExpression&, const QString&).
func (q *QString) ReplaceQRegularExpression(re *QRegularExpression, after *QString) {
	if re == nil || re.p == nil {
		panic("qstring: ReplaceQRegularExpression nil regexp")
	}
	if after == nil || after.p == nil {
		panic("qstring: ReplaceQRegularExpression nil after")
	}
	C.qstring_replace_qregexp_qstring(q.p, re.p, after.p)
	runtime.KeepAlive(re)
	runtime.KeepAlive(after)
}

// ReplaceQRegularExpressionString replaces all matches of re with a Go string (UTF-8).
// Qt: QString::replace(const QRegularExpression&, const QString&) (UTF-8 decoded).
func (q *QString) ReplaceQRegularExpressionString(re *QRegularExpression, after string) {
	qa := QStringFromString(after)
	defer qa.Close()
	q.ReplaceQRegularExpression(re, qa)
}

// RemoveQRegularExpression removes all matches of re.
// Qt: QString::remove(const QRegularExpression&).
func (q *QString) RemoveQRegularExpression(re *QRegularExpression) {
	if re == nil || re.p == nil {
		panic("qstring: RemoveQRegularExpression nil regexp")
	}
	C.qstring_remove_qregexp(q.p, re.p)
	runtime.KeepAlive(re)
}

// SplitQRegularExpression splits q by matches of sep.
// Qt: QString::split(const QRegularExpression&, Qt::SplitBehavior) const.
func (q *QString) SplitQRegularExpression(sep *QRegularExpression,
	behavior qdata.SplitBehavior) []*QString {
	if sep == nil || sep.p == nil {
		panic("qstring: SplitQRegularExpression nil sep")
	}
	var count C.ptrdiff_t
	arr := C.qstring_split_qregexp(q.p, sep.p, C.int(behavior), &count)
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

// ---- case / trim / escape -----------------------------------------------

// ToLower returns a lowercased copy.
// Qt: QString::toLower() const &.
func (q *QString) ToLower() *QString {
	return wrapQString(unsafe.Pointer(C.qstring_to_lower(q.p)))
}

// ToUpper returns an uppercased copy.
// Qt: QString::toUpper() const &.
func (q *QString) ToUpper() *QString {
	return wrapQString(unsafe.Pointer(C.qstring_to_upper(q.p)))
}

// ToCaseFolded returns a case-folded copy (for case-insensitive comparison).
// Qt: QString::toCaseFolded() const &.
func (q *QString) ToCaseFolded() *QString {
	return wrapQString(unsafe.Pointer(C.qstring_to_case_folded(q.p)))
}

// ToHtmlEscaped returns a copy with HTML metacharacters escaped.
// Qt: QString::toHtmlEscaped() const.
func (q *QString) ToHtmlEscaped() *QString {
	return wrapQString(unsafe.Pointer(C.qstring_to_html_escaped(q.p)))
}

// IsLower reports whether the string contains no uppercase letters.
// Qt: QString::isLower() const.
func (q *QString) IsLower() bool { return bool(C.qstring_is_lower(q.p)) }

// IsUpper reports whether the string contains no lowercase letters.
// Qt: QString::isUpper() const.
func (q *QString) IsUpper() bool { return bool(C.qstring_is_upper(q.p)) }

// Trimmed returns a copy with leading and trailing whitespace removed.
// Qt: QString::trimmed() const &.
func (q *QString) Trimmed() *QString {
	return wrapQString(unsafe.Pointer(C.qstring_trimmed(q.p)))
}

// Simplified returns a copy with internal whitespace collapsed and trimmed.
// Qt: QString::simplified() const &.
func (q *QString) Simplified() *QString {
	return wrapQString(unsafe.Pointer(C.qstring_simplified(q.p)))
}

// LeftJustified returns a copy padded on the right to width characters.
// Qt: QString::leftJustified(qsizetype, QChar, bool) const.
func (q *QString) LeftJustified(width int, fill QChar, truncate bool) *QString {
	return wrapQString(unsafe.Pointer(C.qstring_left_justified(
		q.p, C.ptrdiff_t(width), C.uint16_t(fill.Unicode()), C.bool(truncate))))
}

// RightJustified returns a copy padded on the left to width characters.
// Qt: QString::rightJustified(qsizetype, QChar, bool) const.
func (q *QString) RightJustified(width int, fill QChar, truncate bool) *QString {
	return wrapQString(unsafe.Pointer(C.qstring_right_justified(
		q.p, C.ptrdiff_t(width), C.uint16_t(fill.Unicode()), C.bool(truncate))))
}

// ---- normalization -------------------------------------------------------

// Normalized returns a Unicode-normalized copy using the default Unicode version
// (QChar::Unicode_Unassigned).
// Qt: QString::normalized(QString::NormalizationForm, QChar::UnicodeVersion) const — version = QChar::Unicode_Unassigned.
func (q *QString) Normalized(mode NormalizationForm) *QString {
	return q.NormalizedWithVersion(mode, 0)
}

// NormalizedWithVersion returns a Unicode-normalized copy using the given Unicode version.
// version 0 corresponds to QChar::Unicode_Unassigned.
// Qt: QString::normalized(QString::NormalizationForm, QChar::UnicodeVersion) const.
func (q *QString) NormalizedWithVersion(mode NormalizationForm, version int) *QString {
	return wrapQString(unsafe.Pointer(C.qstring_normalized(q.p, C.int(mode), C.int(version))))
}

// ---- compare ------------------------------------------------------------

// CompareQString compares q with another QString.
// Returns negative, zero, or positive like strcmp.
// Qt: QString::compare(const QString&, Qt::CaseSensitivity) const.
func (q *QString) CompareQString(other *QString, caseSensitive bool) int {
	r := C.qstring_compare_qstring(q.p, other.p, csInt(caseSensitive))
	runtime.KeepAlive(other)
	return int(r)
}

// CompareQChar compares q with a single QChar.
// Qt: QString::compare(QChar, Qt::CaseSensitivity) const.
func (q *QString) CompareQChar(c QChar, caseSensitive bool) int {
	return int(C.qstring_compare_qchar(q.p, C.uint16_t(c.Unicode()), csInt(caseSensitive)))
}

// CompareString compares q with a Go string (UTF-8 decoded).
// Qt: QString::compare(const QString&, Qt::CaseSensitivity) const (UTF-8 decoded;
//     Qt also exposes compare(QLatin1StringView, cs) and compare(QStringView, cs)).
func (q *QString) CompareString(s string, caseSensitive bool) int {
	var p *C.char
	if len(s) > 0 {
		p = (*C.char)(unsafe.Pointer(unsafe.StringData(s)))
	}
	r := C.qstring_compare_utf8(q.p, p, C.ptrdiff_t(len(s)), csInt(caseSensitive))
	runtime.KeepAlive(s)
	return int(r)
}

// LocaleAwareCompare compares q with another QString in a locale-aware manner.
// Qt: QString::localeAwareCompare(const QString&) const.
func (q *QString) LocaleAwareCompare(other *QString) int {
	r := C.qstring_locale_aware_compare_qstring(q.p, other.p)
	runtime.KeepAlive(other)
	return int(r)
}

// ---- numeric ------------------------------------------------------------

// ToShort parses q as a signed 16-bit integer in the given base.
// Qt: QString::toShort(bool*, int) const.
func (q *QString) ToShort(base int) (int16, bool) {
	var ok C.bool
	r := C.qstring_to_short(q.p, &ok, C.int(base))
	runtime.KeepAlive(q)
	return int16(r), bool(ok)
}

// ToUShort parses q as an unsigned 16-bit integer in the given base.
// Qt: QString::toUShort(bool*, int) const.
func (q *QString) ToUShort(base int) (uint16, bool) {
	var ok C.bool
	r := C.qstring_to_ushort(q.p, &ok, C.int(base))
	runtime.KeepAlive(q)
	return uint16(r), bool(ok)
}

// ToInt parses q as a signed 32-bit integer in the given base.
// Qt: QString::toInt(bool*, int) const.
func (q *QString) ToInt(base int) (int32, bool) {
	var ok C.bool
	r := C.qstring_to_int(q.p, &ok, C.int(base))
	runtime.KeepAlive(q)
	return int32(r), bool(ok)
}

// ToUInt parses q as an unsigned 32-bit integer in the given base.
// Qt: QString::toUInt(bool*, int) const.
func (q *QString) ToUInt(base int) (uint32, bool) {
	var ok C.bool
	r := C.qstring_to_uint(q.p, &ok, C.int(base))
	runtime.KeepAlive(q)
	return uint32(r), bool(ok)
}

// ToLong parses q as a signed long in the given base.
// Qt: QString::toLong(bool*, int) const.
func (q *QString) ToLong(base int) (int64, bool) {
	var ok C.bool
	r := C.qstring_to_long(q.p, &ok, C.int(base))
	runtime.KeepAlive(q)
	return int64(r), bool(ok)
}

// ToULong parses q as an unsigned long in the given base.
// Qt: QString::toULong(bool*, int) const.
func (q *QString) ToULong(base int) (uint64, bool) {
	var ok C.bool
	r := C.qstring_to_ulong(q.p, &ok, C.int(base))
	runtime.KeepAlive(q)
	return uint64(r), bool(ok)
}

// ToLongLong parses q as a signed 64-bit integer in the given base.
// Qt: QString::toLongLong(bool*, int) const.
func (q *QString) ToLongLong(base int) (int64, bool) {
	var ok C.bool
	r := C.qstring_to_longlong(q.p, &ok, C.int(base))
	runtime.KeepAlive(q)
	return int64(r), bool(ok)
}

// ToULongLong parses q as an unsigned 64-bit integer in the given base.
// Qt: QString::toULongLong(bool*, int) const.
func (q *QString) ToULongLong(base int) (uint64, bool) {
	var ok C.bool
	r := C.qstring_to_ulonglong(q.p, &ok, C.int(base))
	runtime.KeepAlive(q)
	return uint64(r), bool(ok)
}

// ToFloat parses q as a float.
// Qt: QString::toFloat(bool*) const.
func (q *QString) ToFloat() (float32, bool) {
	var ok C.bool
	r := C.qstring_to_float(q.p, &ok)
	runtime.KeepAlive(q)
	return float32(r), bool(ok)
}

// ToDouble parses q as a double.
// Qt: QString::toDouble(bool*) const.
func (q *QString) ToDouble() (float64, bool) {
	var ok C.bool
	r := C.qstring_to_double(q.p, &ok)
	runtime.KeepAlive(q)
	return float64(r), bool(ok)
}

// SetNumLL sets q to the textual representation of a signed 64-bit integer.
// Qt: QString::setNum(qlonglong, int).
func (q *QString) SetNumLL(n int64, base int) {
	C.qstring_set_num_ll(q.p, C.int64_t(n), C.int(base))
}

// SetNumULL sets q to the textual representation of an unsigned 64-bit integer.
// Qt: QString::setNum(qulonglong, int).
func (q *QString) SetNumULL(n uint64, base int) {
	C.qstring_set_num_ull(q.p, C.uint64_t(n), C.int(base))
}

// SetNumInt sets q to the textual representation of a signed 32-bit integer.
// Qt: QString::setNum(int, int).
func (q *QString) SetNumInt(n int32, base int) {
	C.qstring_set_num_int(q.p, C.int32_t(n), C.int(base))
}

// SetNumUInt sets q to the textual representation of an unsigned 32-bit integer.
// Qt: QString::setNum(uint, int).
func (q *QString) SetNumUInt(n uint32, base int) {
	C.qstring_set_num_uint(q.p, C.uint32_t(n), C.int(base))
}

// SetNumFloat sets q to the textual representation of a float.
// Qt: QString::setNum(float, char, int).
func (q *QString) SetNumFloat(v float32, format byte, precision int) {
	C.qstring_set_num_float(q.p, C.float(v), C.char(format), C.int(precision))
}

// SetNumDouble sets q to the textual representation of a double.
// Qt: QString::setNum(double, char, int).
func (q *QString) SetNumDouble(v float64, format byte, precision int) {
	C.qstring_set_num_double(q.p, C.double(v), C.char(format), C.int(precision))
}

// QStringNumberLL returns a QString with the textual representation of a signed 64-bit integer.
// Qt: QString::number(qlonglong, int).
func QStringNumberLL(n int64, base int) *QString {
	return wrapQString(unsafe.Pointer(C.qstring_number_ll(C.int64_t(n), C.int(base))))
}

// QStringNumberULL returns a QString with the textual representation of an unsigned 64-bit integer.
// Qt: QString::number(qulonglong, int).
func QStringNumberULL(n uint64, base int) *QString {
	return wrapQString(unsafe.Pointer(C.qstring_number_ull(C.uint64_t(n), C.int(base))))
}

// QStringNumberInt returns a QString with the textual representation of a signed 32-bit integer.
// Qt: QString::number(int, int).
func QStringNumberInt(n int32, base int) *QString {
	return wrapQString(unsafe.Pointer(C.qstring_number_int(C.int32_t(n), C.int(base))))
}

// QStringNumberUInt returns a QString with the textual representation of an unsigned 32-bit integer.
// Qt: QString::number(uint, int).
func QStringNumberUInt(n uint32, base int) *QString {
	return wrapQString(unsafe.Pointer(C.qstring_number_uint(C.uint32_t(n), C.int(base))))
}

// QStringNumberDouble returns a QString with the textual representation of a double.
// Qt: QString::number(double, char, int).
func QStringNumberDouble(v float64, format byte, precision int) *QString {
	return wrapQString(unsafe.Pointer(C.qstring_number_double(C.double(v), C.char(format),
		C.int(precision))))
}

// ---- static / predicates ------------------------------------------------

// QStringMaxSize returns the maximum QString length in code units.
// Qt: QString::maxSize().
func QStringMaxSize() int { return int(C.qstring_max_size()) }

// IsRightToLeft reports whether q has right-to-left directionality.
// Qt: QString::isRightToLeft() const.
func (q *QString) IsRightToLeft() bool { return bool(C.qstring_is_right_to_left(q.p)) }

// IsValidUtf16 reports whether q contains valid UTF-16 (no unpaired surrogates).
// Qt: QString::isValidUtf16() const.
func (q *QString) IsValidUtf16() bool { return bool(C.qstring_is_valid_utf16(q.p)) }

// ---- arg (chainable, mirrors QString::arg) ------------------------------

// ArgLl substitutes the lowest-numbered %n placeholder with a signed 64-bit integer.
// Qt: QString::arg(qlonglong, int, int, QChar) const.
func (q *QString) ArgLl(n int64, fieldWidth, base int, fillChar QChar) *QString {
	return wrapQString(unsafe.Pointer(C.qstring_arg_ll(q.p, C.int64_t(n), C.int(fieldWidth),
		C.int(base), C.uint16_t(fillChar.Unicode()))))
}

// ArgUll substitutes the lowest-numbered %n placeholder with an unsigned 64-bit integer.
// Qt: QString::arg(qulonglong, int, int, QChar) const.
func (q *QString) ArgUll(n uint64, fieldWidth, base int, fillChar QChar) *QString {
	return wrapQString(unsafe.Pointer(C.qstring_arg_ull(q.p, C.uint64_t(n), C.int(fieldWidth),
		C.int(base), C.uint16_t(fillChar.Unicode()))))
}

// ArgDouble substitutes the lowest-numbered %n placeholder with a double.
// format and precision follow Qt conventions: 'g' (default) and -1 for shortest.
// Qt: QString::arg(double, int, char, int, QChar) const.
func (q *QString) ArgDouble(v float64, fieldWidth int, format byte, precision int,
	fillChar QChar) *QString {
	return wrapQString(unsafe.Pointer(C.qstring_arg_double(q.p, C.double(v), C.int(fieldWidth),
		C.char(format), C.int(precision), C.uint16_t(fillChar.Unicode()))))
}

// ArgQChar substitutes the lowest-numbered %n placeholder with a single character.
// Qt: QString::arg(QChar, int, QChar) const.
func (q *QString) ArgQChar(c QChar, fieldWidth int, fillChar QChar) *QString {
	return wrapQString(unsafe.Pointer(C.qstring_arg_qchar(q.p, C.uint16_t(c.Unicode()),
		C.int(fieldWidth), C.uint16_t(fillChar.Unicode()))))
}

// ArgQString substitutes the lowest-numbered %n placeholder with another QString.
// Qt: QString::arg(const QString&, int, QChar) const.
func (q *QString) ArgQString(other *QString, fieldWidth int, fillChar QChar) *QString {
	r := wrapQString(unsafe.Pointer(C.qstring_arg_qstring(q.p, other.p,
		C.int(fieldWidth), C.uint16_t(fillChar.Unicode()))))
	runtime.KeepAlive(other)
	return r
}

// ArgString substitutes the lowest-numbered %n placeholder with a Go string (UTF-8 decoded).
// Qt: QString::arg(const QString&, int, QChar) const (UTF-8 decoded; Qt also exposes
//     arg(QStringView, int, QChar) and arg(QLatin1StringView, int, QChar)).
func (q *QString) ArgString(s string, fieldWidth int, fillChar QChar) *QString {
	var p *C.char
	if len(s) > 0 {
		p = (*C.char)(unsafe.Pointer(unsafe.StringData(s)))
	}
	r := wrapQString(unsafe.Pointer(C.qstring_arg_utf8(q.p, p, C.ptrdiff_t(len(s)),
		C.int(fieldWidth), C.uint16_t(fillChar.Unicode()))))
	runtime.KeepAlive(s)
	return r
}

// ---- split --------------------------------------------------------------

// SplitQString splits q by sep. behavior controls empty-part handling.
// The returned slice owns each *QString independently; the outer array is
// released automatically.
// Qt: QString::split(const QString&, Qt::SplitBehavior, Qt::CaseSensitivity) const.
func (q *QString) SplitQString(sep *QString, behavior qdata.SplitBehavior,
	caseSensitive bool) []*QString {
	var count C.ptrdiff_t
	arr := C.qstring_split_qstring(q.p, sep.p, C.int(behavior), csInt(caseSensitive), &count)
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

// SplitQChar splits q by a single character.
// Qt: QString::split(QChar, Qt::SplitBehavior, Qt::CaseSensitivity) const.
func (q *QString) SplitQChar(sep QChar, behavior qdata.SplitBehavior,
	caseSensitive bool) []*QString {
	var count C.ptrdiff_t
	arr := C.qstring_split_qchar(q.p, C.uint16_t(sep.Unicode()), C.int(behavior),
		csInt(caseSensitive), &count)
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

// SectionQChar returns the section delimited by sep, from start to end inclusive.
// Qt: QString::section(QChar, qsizetype, qsizetype, SectionFlags) const.
func (q *QString) SectionQChar(sep QChar, start, end int, flags SectionFlag) *QString {
	return wrapQString(unsafe.Pointer(C.qstring_section_qchar(q.p, C.uint16_t(sep.Unicode()),
		C.ptrdiff_t(start), C.ptrdiff_t(end), C.int(flags))))
}

// SectionQString returns the section delimited by sep, from start to end inclusive.
// Qt: QString::section(const QString&, qsizetype, qsizetype, SectionFlags) const.
func (q *QString) SectionQString(sep *QString, start, end int, flags SectionFlag) *QString {
	r := wrapQString(unsafe.Pointer(C.qstring_section_qstring(q.p, sep.p,
		C.ptrdiff_t(start), C.ptrdiff_t(end), C.int(flags))))
	runtime.KeepAlive(sep)
	return r
}

// ---- toUcs4 -------------------------------------------------------------

// ToUcs4 returns the string as a slice of Unicode code points (32-bit).
// Qt: QString::toUcs4() const.
func (q *QString) ToUcs4() []uint32 {
	var count C.ptrdiff_t
	arr := C.qstring_to_ucs4(q.p, &count)
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

// ---- convenience constructors ------------------------------------------

// QStringFromQChar returns a QString of length 1 containing c.
// Qt: QString(QChar) — explicit constructor.
func QStringFromQChar(c QChar) *QString {
	return QStringFromUtf16([]uint16{uint16(c)})
}

// QStringFromRepeat returns a QString of length n with every code unit set to c.
// Qt: QString(qsizetype, QChar) — explicit constructor.
func QStringFromRepeat(n int, c QChar) *QString {
	s := NewQString()
	s.ResizeFill(n, c)
	return s
}
