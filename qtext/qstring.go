package qtext

/*
 Panic vs clamp:

   - At / First / Last / Sliced / SlicedN / Chopped — паникуют при выходе
     за границы (Qt в release-сборке отключает Q_ASSERT → UB).
   - Left / Right / Mid — клампят значения (семантика Qt), без паники;
     соответствует QByteArray-биндингу.

 Что сознательно не вошло:

   - arg() — variadic template, требует много перегрузок.
   - split() / splitRef() — возвращают QStringList / QList<QStringView>
     (нужны специализации контейнеров).
   - QRegularExpression / QRegularExpressionMatch — отдельный пакет.
   - localeAwareCompare() — требует QLocale.
   - toHtmlEscaped() — cross-cutting escaping.
   - toPercentEncoding() / fromPercentEncoding() — вместе с qurl или отдельно.
   - toStdString() / toStdU16String() / toStdU32String() — std::basic_string interop.
   - toWCharArray() / toUcs4() — платформо-зависимые / контейнер-зависимые.
   - vasprintf / asprintf / sprintf — не идиоматичны в Go.
   - Reverse-iterators (rbegin/rend/crbegin/crend).
   - QChar*-конструкторы и setRawData/fromRawData.
   - QLatin1StringView-перегрузки compare/startsWith/… — добавляются
     по мере необходимости (импортных циклов больше нет).
   - operator+ / operator+= — заменяются Clone().Append*().
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

	qt6core "github.com/vadimlxx/go-qt6core/internal/qt6core"
)

// QString wraps a Qt6 QString (UTF-16 COW string).
// Each *QString owns a heap-allocated C++ QString; the finalizer (and
// explicit Close) frees it.
type QString struct {
	p *C.QString
}

func wrapQString(p unsafe.Pointer) *QString {
	if p == nil {
		return nil
	}
	q := &QString{p: (*C.QString)(p)}
	runtime.SetFinalizer(q, finalizeQString)
	return q
}

func finalizeQString(q *QString) {
	if q.p != nil {
		C.qstring_delete(q.p)
		q.p = nil
	}
}

// Close releases the underlying C++ object immediately and disables finalizer.
func (q *QString) Close() {
	if q.p != nil {
		C.qstring_delete(q.p)
		q.p = nil
		runtime.SetFinalizer(q, nil)
	}
}

// ---- construction --------------------------------------------------------

// New returns an empty QString.
func NewQString() *QString {
	return wrapQString(unsafe.Pointer(C.qstring_new_empty()))
}

// FromString constructs a QString by decoding s as UTF-8.
func QStringFromString(s string) *QString {
	return wrapQString(qt6core.QStringFromUtf8(s))
}

// FromUtf8 constructs a QString by decoding b as UTF-8.
func QStringFromUtf8(b []byte) *QString {
	return wrapQString(qt6core.QStringFromUtf8Bytes(b))
}

// FromUtf16 constructs a QString from UTF-16 code units.
func QStringFromUtf16(units []uint16) *QString {
	return wrapQString(qt6core.QStringFromUtf16(units))
}

// FromLatin1 constructs a QString by decoding b as Latin-1.
func QStringFromLatin1(b []byte) *QString {
	return wrapQString(qt6core.QStringFromLatin1(b))
}

// FromQByteArray constructs a QString by decoding ba as UTF-8
// (mirrors QString::fromUtf8(const QByteArray&)).
func QStringFromQByteArray(ba *QByteArray) *QString {
	return wrapQString(qt6core.QStringFromQByteArray(ba.UnsafePtr()))
}

// FromQStringView constructs a QString by copying the view's contents.
func QStringFromQStringView(v QStringView) *QString {
	if v.IsNull() {
		return wrapQString(unsafe.Pointer(C.qstring_new_empty()))
	}
	return wrapQString(unsafe.Pointer(C.qstring_new_from_qstringview(
		(*C.uint16_t)(v.Utf16Unsafe()), C.ptrdiff_t(v.Size()))))
}

// Clone returns a shallow COW copy (Qt refcount++; cheap).
func (q *QString) Clone() *QString {
	return wrapQString(unsafe.Pointer(C.qstring_new_copy(q.p)))
}

// ---- accessors -----------------------------------------------------------

func (q *QString) Size() int  { return int(C.qstring_size(q.p)) }
func (q *QString) Len() int   { return q.Size() }
func (q *QString) IsEmpty() bool    { return bool(C.qstring_is_empty(q.p)) }
func (q *QString) IsNull() bool     { return bool(C.qstring_is_null(q.p)) }
func (q *QString) IsDetached() bool { return bool(C.qstring_is_detached(q.p)) }
func (q *QString) Detach()          { C.qstring_detach(q.p) }

// Utf16Unsafe returns the raw internal UTF-16 data pointer.
// Valid until the QString is mutated, resized, or Close()d.
func (q *QString) Utf16Unsafe() unsafe.Pointer {
	return unsafe.Pointer(C.qstring_utf16(q.p))
}

// Utf16 returns a copy of the contents as UTF-16 code units.
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
func (q *QString) GoString() string {
	s := qt6core.QStringToUtf8String(unsafe.Pointer(q.p))
	runtime.KeepAlive(q)
	return s
}

// String implements fmt.Stringer.
func (q *QString) String() string { return q.GoString() }

// Utf8 returns a fresh *QByteArray with UTF-8 encoding.
func (q *QString) Utf8() *QByteArray {
	var out *C.QByteArray
	C.qstring_to_utf8_qba(q.p, &out)
	runtime.KeepAlive(q)
	return QByteArrayFromUnsafePtr(unsafe.Pointer(out))
}

// Latin1 returns a fresh *QByteArray with Latin-1 encoding
// (chars outside Latin-1 become '?').
func (q *QString) Latin1() *QByteArray {
	var out *C.QByteArray
	C.qstring_to_latin1_qba(q.p, &out)
	runtime.KeepAlive(q)
	return QByteArrayFromUnsafePtr(unsafe.Pointer(out))
}

// Local8Bit returns a fresh *QByteArray with local 8-bit encoding.
func (q *QString) Local8Bit() *QByteArray {
	var out *C.QByteArray
	C.qstring_to_local8bit_qba(q.p, &out)
	runtime.KeepAlive(q)
	return QByteArrayFromUnsafePtr(unsafe.Pointer(out))
}

// At returns the i-th QChar. Panics if i is out of range.
func (q *QString) At(i int) QChar {
	if i < 0 || i >= q.Size() {
		panic("qstring: At out of range")
	}
	r := C.qstring_at(q.p, C.ptrdiff_t(i))
	runtime.KeepAlive(q)
	return NewQChar(uint16(r))
}

// Front returns the first QChar. Panics if empty.
func (q *QString) Front() QChar {
	if q.IsEmpty() {
		panic("qstring: Front on empty QString")
	}
	r := C.qstring_front(q.p)
	runtime.KeepAlive(q)
	return NewQChar(uint16(r))
}

// Back returns the last QChar. Panics if empty.
func (q *QString) Back() QChar {
	if q.IsEmpty() {
		panic("qstring: Back on empty QString")
	}
	r := C.qstring_back(q.p)
	runtime.KeepAlive(q)
	return NewQChar(uint16(r))
}

// ToQStringView returns a zero-copy QStringView over the internal data.
func (q *QString) ToQStringView() QStringView {
	v := newQStringViewWithOwner(q.Utf16Unsafe(), q.Size(), q)
	runtime.KeepAlive(q)
	return v
}

// ---- modifiers -----------------------------------------------------------

func (q *QString) Clear()                    { C.qstring_clear(q.p) }
func (q *QString) Resize(n int)              { C.qstring_resize(q.p, C.ptrdiff_t(n)) }
func (q *QString) ResizeFill(n int, c QChar) { C.qstring_resize_fill(q.p, C.ptrdiff_t(n), C.uint16_t(c.Unicode())) }
func (q *QString) Reserve(n int)             { C.qstring_reserve(q.p, C.ptrdiff_t(n)) }
func (q *QString) Squeeze()                  { C.qstring_squeeze(q.p) }
func (q *QString) Truncate(pos int)          { C.qstring_truncate(q.p, C.ptrdiff_t(pos)) }
func (q *QString) Chop(n int)                { C.qstring_chop(q.p, C.ptrdiff_t(n)) }
func (q *QString) Fill(c QChar, n int) { C.qstring_fill(q.p, C.uint16_t(c.Unicode()), C.ptrdiff_t(n)) }

func (q *QString) Append(c QChar) {
	C.qstring_append_qchar(q.p, C.uint16_t(c.Unicode()))
}
func (q *QString) AppendQString(other *QString) {
	C.qstring_append_qstring(q.p, other.p)
	runtime.KeepAlive(other)
}
func (q *QString) AppendString(s string) {
	var p *C.char
	if len(s) > 0 {
		p = (*C.char)(unsafe.Pointer(unsafe.StringData(s)))
	}
	C.qstring_append_utf8(q.p, p, C.ptrdiff_t(len(s)))
	runtime.KeepAlive(s)
}

func (q *QString) Prepend(c QChar) {
	C.qstring_prepend_qchar(q.p, C.uint16_t(c.Unicode()))
}
func (q *QString) PrependQString(other *QString) {
	C.qstring_prepend_qstring(q.p, other.p)
	runtime.KeepAlive(other)
}
func (q *QString) PrependString(s string) {
	var p *C.char
	if len(s) > 0 {
		p = (*C.char)(unsafe.Pointer(unsafe.StringData(s)))
	}
	C.qstring_prepend_utf8(q.p, p, C.ptrdiff_t(len(s)))
	runtime.KeepAlive(s)
}

func (q *QString) Insert(i int, c QChar) {
	C.qstring_insert_qchar(q.p, C.ptrdiff_t(i), C.uint16_t(c.Unicode()))
}
func (q *QString) InsertQString(i int, other *QString) {
	C.qstring_insert_qstring(q.p, C.ptrdiff_t(i), other.p)
	runtime.KeepAlive(other)
}
func (q *QString) InsertString(i int, s string) {
	var p *C.char
	if len(s) > 0 {
		p = (*C.char)(unsafe.Pointer(unsafe.StringData(s)))
	}
	C.qstring_insert_utf8(q.p, C.ptrdiff_t(i), p, C.ptrdiff_t(len(s)))
	runtime.KeepAlive(s)
}

func (q *QString) Remove(i, n int) { C.qstring_remove(q.p, C.ptrdiff_t(i), C.ptrdiff_t(n)) }
func (q *QString) RemoveAt(i int)  { C.qstring_remove_at(q.p, C.ptrdiff_t(i)) }
func (q *QString) RemoveFirst()    { C.qstring_remove_first(q.p) }
func (q *QString) RemoveLast()     { C.qstring_remove_last(q.p) }

func (q *QString) ReplaceQCharQChar(before, after QChar) {
	C.qstring_replace_qchar_qchar(q.p, C.uint16_t(before.Unicode()), C.uint16_t(after.Unicode()))
}
func (q *QString) ReplaceQCharQString(before QChar, after *QString) {
	C.qstring_replace_qchar_qstring(q.p, C.uint16_t(before.Unicode()), after.p)
	runtime.KeepAlive(after)
}
func (q *QString) ReplaceQCharString(before QChar, after string) {
	var p *C.char
	if len(after) > 0 { p = (*C.char)(unsafe.Pointer(unsafe.StringData(after))) }
	C.qstring_replace_qchar_utf8(q.p, C.uint16_t(before.Unicode()), p, C.ptrdiff_t(len(after)))
	runtime.KeepAlive(after)
}
func (q *QString) ReplaceQStringQChar(before *QString, after QChar) {
	C.qstring_replace_qstring_qchar(q.p, before.p, C.uint16_t(after.Unicode()))
	runtime.KeepAlive(before)
}
func (q *QString) ReplaceQStringQString(before, after *QString) {
	C.qstring_replace_qstring_qstring(q.p, before.p, after.p)
	runtime.KeepAlive(before)
	runtime.KeepAlive(after)
}
func (q *QString) ReplaceQStringString(before *QString, after string) {
	var p *C.char
	if len(after) > 0 { p = (*C.char)(unsafe.Pointer(unsafe.StringData(after))) }
	C.qstring_replace_qstring_utf8(q.p, before.p, p, C.ptrdiff_t(len(after)))
	runtime.KeepAlive(before)
	runtime.KeepAlive(after)
}
func (q *QString) ReplaceStringQChar(before string, after QChar) {
	var p *C.char
	if len(before) > 0 { p = (*C.char)(unsafe.Pointer(unsafe.StringData(before))) }
	C.qstring_replace_utf8_qchar(q.p, p, C.ptrdiff_t(len(before)), C.uint16_t(after.Unicode()))
	runtime.KeepAlive(before)
}
func (q *QString) ReplaceStringQString(before string, after *QString) {
	var p *C.char
	if len(before) > 0 { p = (*C.char)(unsafe.Pointer(unsafe.StringData(before))) }
	C.qstring_replace_utf8_qstring(q.p, p, C.ptrdiff_t(len(before)), after.p)
	runtime.KeepAlive(before)
	runtime.KeepAlive(after)
}
func (q *QString) ReplaceStringString(before, after string) {
	var pb, pa *C.char
	if len(before) > 0 { pb = (*C.char)(unsafe.Pointer(unsafe.StringData(before))) }
	if len(after) > 0  { pa = (*C.char)(unsafe.Pointer(unsafe.StringData(after)))  }
	C.qstring_replace_utf8_utf8(q.p, pb, C.ptrdiff_t(len(before)), pa, C.ptrdiff_t(len(after)))
	runtime.KeepAlive(before)
	runtime.KeepAlive(after)
}
func (q *QString) ReplaceRangeString(i, n int, after string) {
	var p *C.char
	if len(after) > 0 { p = (*C.char)(unsafe.Pointer(unsafe.StringData(after))) }
	C.qstring_replace_range_utf8(q.p, C.ptrdiff_t(i), C.ptrdiff_t(n), p, C.ptrdiff_t(len(after)))
	runtime.KeepAlive(after)
}

// ---- slicing (return new owned QString) ---------------------------------

func (q *QString) Left(n int) *QString {
	return wrapQString(unsafe.Pointer(C.qstring_left(q.p, C.ptrdiff_t(n))))
}
func (q *QString) Right(n int) *QString {
	return wrapQString(unsafe.Pointer(C.qstring_right(q.p, C.ptrdiff_t(n))))
}
func (q *QString) Mid(pos, n int) *QString {
	return wrapQString(unsafe.Pointer(C.qstring_mid(q.p, C.ptrdiff_t(pos), C.ptrdiff_t(n))))
}
func (q *QString) First(n int) *QString {
	if n < 0 || n > q.Size() {
		panic("qstring: First out of range")
	}
	return wrapQString(unsafe.Pointer(C.qstring_first(q.p, C.ptrdiff_t(n))))
}
func (q *QString) Last(n int) *QString {
	if n < 0 || n > q.Size() {
		panic("qstring: Last out of range")
	}
	return wrapQString(unsafe.Pointer(C.qstring_last(q.p, C.ptrdiff_t(n))))
}
func (q *QString) Sliced(pos int) *QString {
	if pos < 0 || pos > q.Size() {
		panic("qstring: Sliced out of range")
	}
	return wrapQString(unsafe.Pointer(C.qstring_sliced(q.p, C.ptrdiff_t(pos))))
}
func (q *QString) SlicedN(pos, n int) *QString {
	if pos < 0 || n < 0 || pos+n > q.Size() {
		panic("qstring: SlicedN out of range")
	}
	return wrapQString(unsafe.Pointer(C.qstring_sliced_n(q.p, C.ptrdiff_t(pos), C.ptrdiff_t(n))))
}
func (q *QString) Chopped(n int) *QString {
	if n < 0 || n > q.Size() {
		panic("qstring: Chopped out of range")
	}
	return wrapQString(unsafe.Pointer(C.qstring_chopped(q.p, C.ptrdiff_t(n))))
}

// ---- search --------------------------------------------------------------

func (q *QString) IndexOfQChar(c QChar, from int, caseSensitive bool) int {
	return int(C.qstring_index_of_qchar(q.p, C.uint16_t(c.Unicode()), C.ptrdiff_t(from), csInt(caseSensitive)))
}
func (q *QString) IndexOfQString(other *QString, from int, caseSensitive bool) int {
	r := C.qstring_index_of_qstring(q.p, other.p, C.ptrdiff_t(from), csInt(caseSensitive))
	runtime.KeepAlive(other)
	return int(r)
}
func (q *QString) IndexOfString(s string, from int, caseSensitive bool) int {
	var p *C.char
	if len(s) > 0 { p = (*C.char)(unsafe.Pointer(unsafe.StringData(s))) }
	r := C.qstring_index_of_utf8(q.p, p, C.ptrdiff_t(len(s)), C.ptrdiff_t(from), csInt(caseSensitive))
	runtime.KeepAlive(s)
	return int(r)
}

func (q *QString) LastIndexOfQChar(c QChar, from int, caseSensitive bool) int {
	return int(C.qstring_last_index_of_qchar(q.p, C.uint16_t(c.Unicode()), C.ptrdiff_t(from), csInt(caseSensitive)))
}
func (q *QString) LastIndexOfQString(other *QString, from int, caseSensitive bool) int {
	r := C.qstring_last_index_of_qstring(q.p, other.p, C.ptrdiff_t(from), csInt(caseSensitive))
	runtime.KeepAlive(other)
	return int(r)
}
func (q *QString) LastIndexOfString(s string, from int, caseSensitive bool) int {
	var p *C.char
	if len(s) > 0 { p = (*C.char)(unsafe.Pointer(unsafe.StringData(s))) }
	r := C.qstring_last_index_of_utf8(q.p, p, C.ptrdiff_t(len(s)), C.ptrdiff_t(from), csInt(caseSensitive))
	runtime.KeepAlive(s)
	return int(r)
}

func (q *QString) ContainsQChar(c QChar, caseSensitive bool) bool {
	return bool(C.qstring_contains_qchar(q.p, C.uint16_t(c.Unicode()), csInt(caseSensitive)))
}
func (q *QString) ContainsQString(other *QString, caseSensitive bool) bool {
	r := C.qstring_contains_qstring(q.p, other.p, csInt(caseSensitive))
	runtime.KeepAlive(other)
	return bool(r)
}
func (q *QString) ContainsString(s string, caseSensitive bool) bool {
	var p *C.char
	if len(s) > 0 { p = (*C.char)(unsafe.Pointer(unsafe.StringData(s))) }
	r := C.qstring_contains_utf8(q.p, p, C.ptrdiff_t(len(s)), csInt(caseSensitive))
	runtime.KeepAlive(s)
	return bool(r)
}

func (q *QString) CountQChar(c QChar, caseSensitive bool) int {
	return int(C.qstring_count_qchar(q.p, C.uint16_t(c.Unicode()), csInt(caseSensitive)))
}
func (q *QString) CountQString(other *QString, caseSensitive bool) int {
	r := C.qstring_count_qstring(q.p, other.p, csInt(caseSensitive))
	runtime.KeepAlive(other)
	return int(r)
}
func (q *QString) CountString(s string, caseSensitive bool) int {
	var p *C.char
	if len(s) > 0 { p = (*C.char)(unsafe.Pointer(unsafe.StringData(s))) }
	r := C.qstring_count_utf8(q.p, p, C.ptrdiff_t(len(s)), csInt(caseSensitive))
	runtime.KeepAlive(s)
	return int(r)
}

func (q *QString) StartsWithQChar(c QChar, caseSensitive bool) bool {
	return bool(C.qstring_starts_with_qchar(q.p, C.uint16_t(c.Unicode()), csInt(caseSensitive)))
}
func (q *QString) StartsWithQString(other *QString, caseSensitive bool) bool {
	r := C.qstring_starts_with_qstring(q.p, other.p, csInt(caseSensitive))
	runtime.KeepAlive(other)
	return bool(r)
}
func (q *QString) StartsWithString(s string, caseSensitive bool) bool {
	var p *C.char
	if len(s) > 0 { p = (*C.char)(unsafe.Pointer(unsafe.StringData(s))) }
	r := C.qstring_starts_with_utf8(q.p, p, C.ptrdiff_t(len(s)), csInt(caseSensitive))
	runtime.KeepAlive(s)
	return bool(r)
}

func (q *QString) EndsWithQChar(c QChar, caseSensitive bool) bool {
	return bool(C.qstring_ends_with_qchar(q.p, C.uint16_t(c.Unicode()), csInt(caseSensitive)))
}
func (q *QString) EndsWithQString(other *QString, caseSensitive bool) bool {
	r := C.qstring_ends_with_qstring(q.p, other.p, csInt(caseSensitive))
	runtime.KeepAlive(other)
	return bool(r)
}
func (q *QString) EndsWithString(s string, caseSensitive bool) bool {
	var p *C.char
	if len(s) > 0 { p = (*C.char)(unsafe.Pointer(unsafe.StringData(s))) }
	r := C.qstring_ends_with_utf8(q.p, p, C.ptrdiff_t(len(s)), csInt(caseSensitive))
	runtime.KeepAlive(s)
	return bool(r)
}

// ---- case / trim --------------------------------------------------------

func (q *QString) ToLower() *QString {
	return wrapQString(unsafe.Pointer(C.qstring_to_lower(q.p)))
}
func (q *QString) ToUpper() *QString {
	return wrapQString(unsafe.Pointer(C.qstring_to_upper(q.p)))
}
func (q *QString) ToCaseFolded() *QString {
	return wrapQString(unsafe.Pointer(C.qstring_to_case_folded(q.p)))
}
func (q *QString) IsLower() bool { return bool(C.qstring_is_lower(q.p)) }
func (q *QString) IsUpper() bool { return bool(C.qstring_is_upper(q.p)) }

func (q *QString) Trimmed() *QString {
	return wrapQString(unsafe.Pointer(C.qstring_trimmed(q.p)))
}
func (q *QString) Simplified() *QString {
	return wrapQString(unsafe.Pointer(C.qstring_simplified(q.p)))
}
func (q *QString) LeftJustified(width int, fill QChar, truncate bool) *QString {
	return wrapQString(unsafe.Pointer(C.qstring_left_justified(
		q.p, C.ptrdiff_t(width), C.uint16_t(fill.Unicode()), C.bool(truncate))))
}
func (q *QString) RightJustified(width int, fill QChar, truncate bool) *QString {
	return wrapQString(unsafe.Pointer(C.qstring_right_justified(
		q.p, C.ptrdiff_t(width), C.uint16_t(fill.Unicode()), C.bool(truncate))))
}

// ---- compare ------------------------------------------------------------

func (q *QString) CompareQString(other *QString, caseSensitive bool) int {
	r := C.qstring_compare_qstring(q.p, other.p, csInt(caseSensitive))
	runtime.KeepAlive(other)
	return int(r)
}
func (q *QString) CompareQChar(c QChar, caseSensitive bool) int {
	return int(C.qstring_compare_qchar(q.p, C.uint16_t(c.Unicode()), csInt(caseSensitive)))
}
func (q *QString) CompareString(s string, caseSensitive bool) int {
	var p *C.char
	if len(s) > 0 { p = (*C.char)(unsafe.Pointer(unsafe.StringData(s))) }
	r := C.qstring_compare_utf8(q.p, p, C.ptrdiff_t(len(s)), csInt(caseSensitive))
	runtime.KeepAlive(s)
	return int(r)
}

// ---- numeric ------------------------------------------------------------

func (q *QString) ToShort(base int) (int16, bool) {
	var ok C.bool
	r := C.qstring_to_short(q.p, &ok, C.int(base))
	runtime.KeepAlive(q)
	return int16(r), bool(ok)
}
func (q *QString) ToUShort(base int) (uint16, bool) {
	var ok C.bool
	r := C.qstring_to_ushort(q.p, &ok, C.int(base))
	runtime.KeepAlive(q)
	return uint16(r), bool(ok)
}
func (q *QString) ToInt(base int) (int32, bool) {
	var ok C.bool
	r := C.qstring_to_int(q.p, &ok, C.int(base))
	runtime.KeepAlive(q)
	return int32(r), bool(ok)
}
func (q *QString) ToUInt(base int) (uint32, bool) {
	var ok C.bool
	r := C.qstring_to_uint(q.p, &ok, C.int(base))
	runtime.KeepAlive(q)
	return uint32(r), bool(ok)
}
func (q *QString) ToLong(base int) (int64, bool) {
	var ok C.bool
	r := C.qstring_to_long(q.p, &ok, C.int(base))
	runtime.KeepAlive(q)
	return int64(r), bool(ok)
}
func (q *QString) ToULong(base int) (uint64, bool) {
	var ok C.bool
	r := C.qstring_to_ulong(q.p, &ok, C.int(base))
	runtime.KeepAlive(q)
	return uint64(r), bool(ok)
}
func (q *QString) ToLongLong(base int) (int64, bool) {
	var ok C.bool
	r := C.qstring_to_longlong(q.p, &ok, C.int(base))
	runtime.KeepAlive(q)
	return int64(r), bool(ok)
}
func (q *QString) ToULongLong(base int) (uint64, bool) {
	var ok C.bool
	r := C.qstring_to_ulonglong(q.p, &ok, C.int(base))
	runtime.KeepAlive(q)
	return uint64(r), bool(ok)
}
func (q *QString) ToFloat() (float32, bool) {
	var ok C.bool
	r := C.qstring_to_float(q.p, &ok)
	runtime.KeepAlive(q)
	return float32(r), bool(ok)
}
func (q *QString) ToDouble() (float64, bool) {
	var ok C.bool
	r := C.qstring_to_double(q.p, &ok)
	runtime.KeepAlive(q)
	return float64(r), bool(ok)
}

func (q *QString) SetNumLL(n int64, base int) {
	C.qstring_set_num_ll(q.p, C.int64_t(n), C.int(base))
}
func (q *QString) SetNumULL(n uint64, base int) {
	C.qstring_set_num_ull(q.p, C.uint64_t(n), C.int(base))
}
func (q *QString) SetNumDouble(v float64, format byte, precision int) {
	C.qstring_set_num_double(q.p, C.double(v), C.char(format), C.int(precision))
}

func QStringNumberLL(n int64, base int) *QString {
	return wrapQString(unsafe.Pointer(C.qstring_number_ll(C.int64_t(n), C.int(base))))
}
func QStringNumberULL(n uint64, base int) *QString {
	return wrapQString(unsafe.Pointer(C.qstring_number_ull(C.uint64_t(n), C.int(base))))
}
func QStringNumberDouble(v float64, format byte, precision int) *QString {
	return wrapQString(unsafe.Pointer(C.qstring_number_double(C.double(v), C.char(format), C.int(precision))))
}

// ---- static / predicates ------------------------------------------------

// MaxSize mirrors QString::maxSize().
func QStringMaxSize() int { return int(C.qstring_max_size()) }

func (q *QString) IsRightToLeft() bool { return bool(C.qstring_is_right_to_left(q.p)) }
func (q *QString) IsValidUtf16() bool  { return bool(C.qstring_is_valid_utf16(q.p)) }

// FromQChar constructs a QString of length 1 containing c
// (mirrors QString(QChar c)).
func QStringFromQChar(c QChar) *QString {
	return QStringFromUtf16([]uint16{uint16(c)})
}

// FromRepeat constructs a QString of length n with every character set to c
// (mirrors QString(qsizetype size, QChar c)).
func QStringFromRepeat(n int, c QChar) *QString {
	s := NewQString()
	s.ResizeFill(n, c)
	return s
}
