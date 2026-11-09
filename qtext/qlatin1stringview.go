package qtext

/*
 QLatin1String — type alias

 что не вошло:

 toString() — требует qstring (не написан).

 arg() — varargs template, возвращает QString.

 tokenize() — template, возвращает QStringTokenizer.

 toWCharArray() — платформо-зависимый wchar_t.

 toUcs4() — возвращает QList<uint> (нет generic-контейнера).

 QLatin1StringView ↔ QString сравнения (compare(QString) и др.) — вместе с qstring.

 QRegularExpression-перегрузки — вместе с qregularexpression.

 QUtf8StringView-перегрузка compare — вместе с qutf8stringview.
*/

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC -Wno-sfinae-incomplete
#cgo LDFLAGS: -lstdc++
#include <stdint.h>
#include <stdbool.h>
#include <stddef.h>
#include "qlatin1stringview_c.h"
*/
import "C"

import (
	"runtime"
	"unsafe"
)

// QLatin1StringView is a read-only, non-owning view over a Latin-1 buffer.
// Mirrors Qt6 QLatin1StringView (Q_RELOCATABLE_TYPE; trivial value type).
//
// The zero value is a valid null view (size=0, data=nil).
// Use FromBytes / FromString / FromQByteArray to construct from real data.
type QLatin1StringView struct {
	data      unsafe.Pointer // const char*
	size      int
	keepAlive any // holds the underlying []byte / string / *qbytearray.QByteArray
}

// QLatin1String is a synonym for QLatin1StringView (Qt6 name alias).
// Both names refer to the same type in Qt 6.x; see qlatin1stringview.h.
type QLatin1String = QLatin1StringView

// ---- construction --------------------------------------------------------

// Empty returns a null QLatin1StringView (matches QLatin1StringView()).
func EmptyQLatin1StringView() QLatin1StringView { return QLatin1StringView{} }

// NewUnsafe builds a view directly from a raw Latin-1 byte pointer and size.
// Caller is solely responsible for keeping the underlying data alive and
// unmutated for as long as the view (or any of its slices) is used.
func NewQLatin1StringViewUnsafe(data unsafe.Pointer, size int) QLatin1StringView {
	if size < 0 {
		panic("qlatin1stringview: NewUnsafe negative size")
	}
	if size > 0 && data == nil {
		panic("qlatin1stringview: NewUnsafe nil data with non-zero size")
	}
	return QLatin1StringView{data: data, size: size}
}

// FromBytes returns a view over b interpreted as raw Latin-1 bytes.
// The Go slice is kept alive by the view; mutations of b are visible
// through the view.
func QLatin1StringViewFromBytes(b []byte) QLatin1StringView {
	if len(b) == 0 {
		return QLatin1StringView{}
	}
	return QLatin1StringView{
		data:      unsafe.Pointer(&b[0]),
		size:      len(b),
		keepAlive: b,
	}
}

// FromString returns a view over s interpreted as raw Latin-1 bytes.
// For ASCII strings this is exact. For UTF-8 strings with non-ASCII
// code points, each UTF-8 byte becomes one Latin-1 character (the same
// way QByteArray → QLatin1StringView behaves in Qt).
//
// For proper UTF-8 → Latin-1 transcoding, encode manually or wait for
// the qstring package.
func QLatin1StringViewFromString(s string) QLatin1StringView {
	if len(s) == 0 {
		return QLatin1StringView{}
	}
	return QLatin1StringView{
		data:      unsafe.Pointer(unsafe.StringData(s)),
		size:      len(s),
		keepAlive: s,
	}
}

// FromQByteArray returns a zero-copy view into q's internal buffer,
// interpreting the bytes as Latin-1. The view keeps q alive. Do NOT
// mutate or Close q while the view exists.
func QLatin1StringViewFromQByteArray(q *QByteArray) QLatin1StringView {
	if q == nil || q.IsNull() {
		return QLatin1StringView{}
	}
	return QLatin1StringView{
		data:      q.ConstDataUnsafe(),
		size:      q.Size(),
		keepAlive: q,
	}
}

// ---- accessors -----------------------------------------------------------

func (v QLatin1StringView) Size() int    { return v.size }
func (v QLatin1StringView) Len() int     { return v.size }
func (v QLatin1StringView) IsEmpty() bool { return v.size == 0 }
func (v QLatin1StringView) IsNull() bool  { return v.data == nil }

// DataUnsafe returns the raw Latin-1 data pointer.
// Valid only while the view (and its keepAlive target) is alive and unmutated.
func (v QLatin1StringView) DataUnsafe() unsafe.Pointer { return v.data }

// Bytes returns a copy of the view's bytes.
func (v QLatin1StringView) Bytes() []byte {
	if v.size == 0 {
		return nil
	}
	out := make([]byte, v.size)
	copy(out, unsafe.Slice((*byte)(v.data), v.size))
	runtime.KeepAlive(v)
	return out
}

// GoString returns the view's bytes as a Go string (Latin-1 bytes
// taken as-is, no transcoding).
func (v QLatin1StringView) GoString() string {
	var s string
	if v.size > 0 {
		s = string(unsafe.Slice((*byte)(v.data), v.size))
	}
	runtime.KeepAlive(v)
	return s
}

// At returns the i-th qlatin1char.QLatin1Char. Panics if i is out of range.
func (v QLatin1StringView) At(i int) QLatin1Char {
	if i < 0 || i >= v.size {
		panic("qlatin1stringview: At out of range")
	}
	r := C.qlatin1stringview_at((*C.char)(v.data), C.ptrdiff_t(v.size), C.ptrdiff_t(i))
	runtime.KeepAlive(v)
	return NewQLatin1Char(byte(r))
}

// Front returns the first qlatin1char.QLatin1Char. Panics if empty.
func (v QLatin1StringView) Front() QLatin1Char {
	if v.size == 0 {
		panic("qlatin1stringview: Front on empty view")
	}
	r := C.qlatin1stringview_front((*C.char)(v.data), C.ptrdiff_t(v.size))
	runtime.KeepAlive(v)
	return NewQLatin1Char(byte(r))
}

// Back returns the last qlatin1char.QLatin1Char. Panics if empty.
func (v QLatin1StringView) Back() QLatin1Char {
	if v.size == 0 {
		panic("qlatin1stringview: Back on empty view")
	}
	r := C.qlatin1stringview_back((*C.char)(v.data), C.ptrdiff_t(v.size))
	runtime.KeepAlive(v)
	return NewQLatin1Char(byte(r))
}

// ---- slicing (return new views over the same buffer) --------------------

func (v QLatin1StringView) sliceFromC(
	call func(data *C.char, size C.ptrdiff_t, od **C.char, os *C.ptrdiff_t),
) QLatin1StringView {
	var outData *C.char
	var outSize C.ptrdiff_t
	call((*C.char)(v.data), C.ptrdiff_t(v.size), &outData, &outSize)
	runtime.KeepAlive(v)
	return QLatin1StringView{
		data:      unsafe.Pointer(outData),
		size:      int(outSize),
		keepAlive: v.keepAlive,
	}
}

// Mid clamps pos/n (Qt semantics; n = -1 → до конца).
func (v QLatin1StringView) Mid(pos, n int) QLatin1StringView {
	return v.sliceFromC(func(d *C.char, s C.ptrdiff_t, od **C.char, os *C.ptrdiff_t) {
		C.qlatin1stringview_mid(d, s, C.ptrdiff_t(pos), C.ptrdiff_t(n), od, os)
	})
}

// Left clamps n (Qt semantics).
func (v QLatin1StringView) Left(n int) QLatin1StringView {
	return v.sliceFromC(func(d *C.char, s C.ptrdiff_t, od **C.char, os *C.ptrdiff_t) {
		C.qlatin1stringview_left(d, s, C.ptrdiff_t(n), od, os)
	})
}

// Right clamps n (Qt semantics).
func (v QLatin1StringView) Right(n int) QLatin1StringView {
	return v.sliceFromC(func(d *C.char, s C.ptrdiff_t, od **C.char, os *C.ptrdiff_t) {
		C.qlatin1stringview_right(d, s, C.ptrdiff_t(n), od, os)
	})
}

func (v QLatin1StringView) Sliced(pos int) QLatin1StringView {
	if pos < 0 || pos > v.size {
		panic("qlatin1stringview: Sliced out of range")
	}
	return v.sliceFromC(func(d *C.char, s C.ptrdiff_t, od **C.char, os *C.ptrdiff_t) {
		C.qlatin1stringview_sliced(d, s, C.ptrdiff_t(pos), od, os)
	})
}

func (v QLatin1StringView) SlicedN(pos, n int) QLatin1StringView {
	if pos < 0 || n < 0 || pos+n > v.size {
		panic("qlatin1stringview: SlicedN out of range")
	}
	return v.sliceFromC(func(d *C.char, s C.ptrdiff_t, od **C.char, os *C.ptrdiff_t) {
		C.qlatin1stringview_sliced_n(d, s, C.ptrdiff_t(pos), C.ptrdiff_t(n), od, os)
	})
}

func (v QLatin1StringView) First(n int) QLatin1StringView {
	if n < 0 || n > v.size {
		panic("qlatin1stringview: First out of range")
	}
	return v.sliceFromC(func(d *C.char, s C.ptrdiff_t, od **C.char, os *C.ptrdiff_t) {
		C.qlatin1stringview_first(d, s, C.ptrdiff_t(n), od, os)
	})
}

func (v QLatin1StringView) Last(n int) QLatin1StringView {
	if n < 0 || n > v.size {
		panic("qlatin1stringview: Last out of range")
	}
	return v.sliceFromC(func(d *C.char, s C.ptrdiff_t, od **C.char, os *C.ptrdiff_t) {
		C.qlatin1stringview_last(d, s, C.ptrdiff_t(n), od, os)
	})
}

func (v QLatin1StringView) Chopped(n int) QLatin1StringView {
	if n < 0 || n > v.size {
		panic("qlatin1stringview: Chopped out of range")
	}
	return v.sliceFromC(func(d *C.char, s C.ptrdiff_t, od **C.char, os *C.ptrdiff_t) {
		C.qlatin1stringview_chopped(d, s, C.ptrdiff_t(n), od, os)
	})
}

func (v QLatin1StringView) Trimmed() QLatin1StringView {
	return v.sliceFromC(func(d *C.char, s C.ptrdiff_t, od **C.char, os *C.ptrdiff_t) {
		C.qlatin1stringview_trimmed(d, s, od, os)
	})
}

// Truncate shrinks the view to length n in place.
func (v *QLatin1StringView) Truncate(n int) {
	if n < 0 || n > v.size {
		panic("qlatin1stringview: Truncate out of range")
	}
	v.size = n
}

// Chop removes n bytes from the end in place.
func (v *QLatin1StringView) Chop(n int) {
	if n < 0 || n > v.size {
		panic("qlatin1stringview: Chop out of range")
	}
	v.size -= n
}

// Slice replaces the receiver with Sliced(pos).
func (v *QLatin1StringView) Slice(pos int) { *v = v.Sliced(pos) }

// SliceN replaces the receiver with SlicedN(pos, n).
func (v *QLatin1StringView) SliceN(pos, n int) { *v = v.SlicedN(pos, n) }

// ---- compare -------------------------------------------------------------

func (v QLatin1StringView) CompareView(other QLatin1StringView, caseSensitive bool) int {
	r := C.qlatin1stringview_compare_view(
		(*C.char)(v.data), C.ptrdiff_t(v.size),
		(*C.char)(other.data), C.ptrdiff_t(other.size),
		csInt(caseSensitive))
	runtime.KeepAlive(v)
	runtime.KeepAlive(other)
	return int(r)
}

func (v QLatin1StringView) CompareQStringView(other QStringView, caseSensitive bool) int {
	r := C.qlatin1stringview_compare_qstringview(
		(*C.char)(v.data), C.ptrdiff_t(v.size),
		(*C.uint16_t)(other.Utf16Unsafe()), C.ptrdiff_t(other.Size()),
		csInt(caseSensitive))
	runtime.KeepAlive(v)
	runtime.KeepAlive(other)
	return int(r)
}

func (v QLatin1StringView) CompareQChar(c QChar, caseSensitive bool) int {
	r := C.qlatin1stringview_compare_qchar(
		(*C.char)(v.data), C.ptrdiff_t(v.size),
		C.uint16_t(c.Unicode()), csInt(caseSensitive))
	runtime.KeepAlive(v)
	return int(r)
}

// ---- startsWith / endsWith ----------------------------------------------

func (v QLatin1StringView) StartsWithView(prefix QLatin1StringView, caseSensitive bool) bool {
	r := C.qlatin1stringview_starts_with_view(
		(*C.char)(v.data), C.ptrdiff_t(v.size),
		(*C.char)(prefix.data), C.ptrdiff_t(prefix.size),
		csInt(caseSensitive))
	runtime.KeepAlive(v)
	runtime.KeepAlive(prefix)
	return bool(r)
}

func (v QLatin1StringView) StartsWithQStringView(prefix QStringView, caseSensitive bool) bool {
	r := C.qlatin1stringview_starts_with_qstringview(
		(*C.char)(v.data), C.ptrdiff_t(v.size),
		(*C.uint16_t)(prefix.Utf16Unsafe()), C.ptrdiff_t(prefix.Size()),
		csInt(caseSensitive))
	runtime.KeepAlive(v)
	runtime.KeepAlive(prefix)
	return bool(r)
}

func (v QLatin1StringView) StartsWithQChar(c QChar, caseSensitive bool) bool {
	r := C.qlatin1stringview_starts_with_qchar(
		(*C.char)(v.data), C.ptrdiff_t(v.size),
		C.uint16_t(c.Unicode()), csInt(caseSensitive))
	runtime.KeepAlive(v)
	return bool(r)
}

func (v QLatin1StringView) EndsWithView(suffix QLatin1StringView, caseSensitive bool) bool {
	r := C.qlatin1stringview_ends_with_view(
		(*C.char)(v.data), C.ptrdiff_t(v.size),
		(*C.char)(suffix.data), C.ptrdiff_t(suffix.size),
		csInt(caseSensitive))
	runtime.KeepAlive(v)
	runtime.KeepAlive(suffix)
	return bool(r)
}

func (v QLatin1StringView) EndsWithQStringView(suffix QStringView, caseSensitive bool) bool {
	r := C.qlatin1stringview_ends_with_qstringview(
		(*C.char)(v.data), C.ptrdiff_t(v.size),
		(*C.uint16_t)(suffix.Utf16Unsafe()), C.ptrdiff_t(suffix.Size()),
		csInt(caseSensitive))
	runtime.KeepAlive(v)
	runtime.KeepAlive(suffix)
	return bool(r)
}

func (v QLatin1StringView) EndsWithQChar(c QChar, caseSensitive bool) bool {
	r := C.qlatin1stringview_ends_with_qchar(
		(*C.char)(v.data), C.ptrdiff_t(v.size),
		C.uint16_t(c.Unicode()), csInt(caseSensitive))
	runtime.KeepAlive(v)
	return bool(r)
}

// ---- indexOf / lastIndexOf ----------------------------------------------

func (v QLatin1StringView) IndexOfView(other QLatin1StringView, from int, caseSensitive bool) int {
	r := C.qlatin1stringview_index_of_view(
		(*C.char)(v.data), C.ptrdiff_t(v.size),
		(*C.char)(other.data), C.ptrdiff_t(other.size),
		C.ptrdiff_t(from), csInt(caseSensitive))
	runtime.KeepAlive(v)
	runtime.KeepAlive(other)
	return int(r)
}

func (v QLatin1StringView) IndexOfQStringView(other QStringView, from int, caseSensitive bool) int {
	r := C.qlatin1stringview_index_of_qstringview(
		(*C.char)(v.data), C.ptrdiff_t(v.size),
		(*C.uint16_t)(other.Utf16Unsafe()), C.ptrdiff_t(other.Size()),
		C.ptrdiff_t(from), csInt(caseSensitive))
	runtime.KeepAlive(v)
	runtime.KeepAlive(other)
	return int(r)
}

func (v QLatin1StringView) IndexOfQChar(c QChar, from int, caseSensitive bool) int {
	r := C.qlatin1stringview_index_of_qchar(
		(*C.char)(v.data), C.ptrdiff_t(v.size),
		C.uint16_t(c.Unicode()), C.ptrdiff_t(from), csInt(caseSensitive))
	runtime.KeepAlive(v)
	return int(r)
}

func (v QLatin1StringView) LastIndexOfViewFrom(other QLatin1StringView, from int, caseSensitive bool) int {
	r := C.qlatin1stringview_last_index_of_view(
		(*C.char)(v.data), C.ptrdiff_t(v.size),
		(*C.char)(other.data), C.ptrdiff_t(other.size),
		C.ptrdiff_t(from), csInt(caseSensitive))
	runtime.KeepAlive(v)
	runtime.KeepAlive(other)
	return int(r)
}

// LastIndexOfView mirrors the Qt overload without a `from` parameter
// (default from = size()).
func (v QLatin1StringView) LastIndexOfView(other QLatin1StringView, caseSensitive bool) int {
	return v.LastIndexOfViewFrom(other, v.size, caseSensitive)
}

func (v QLatin1StringView) LastIndexOfQStringViewFrom(other QStringView, from int, caseSensitive bool) int {
	r := C.qlatin1stringview_last_index_of_qstringview(
		(*C.char)(v.data), C.ptrdiff_t(v.size),
		(*C.uint16_t)(other.Utf16Unsafe()), C.ptrdiff_t(other.Size()),
		C.ptrdiff_t(from), csInt(caseSensitive))
	runtime.KeepAlive(v)
	runtime.KeepAlive(other)
	return int(r)
}

func (v QLatin1StringView) LastIndexOfQStringView(other QStringView, caseSensitive bool) int {
	return v.LastIndexOfQStringViewFrom(other, v.size, caseSensitive)
}

func (v QLatin1StringView) LastIndexOfQCharFrom(c QChar, from int, caseSensitive bool) int {
	r := C.qlatin1stringview_last_index_of_qchar(
		(*C.char)(v.data), C.ptrdiff_t(v.size),
		C.uint16_t(c.Unicode()), C.ptrdiff_t(from), csInt(caseSensitive))
	runtime.KeepAlive(v)
	return int(r)
}

// LastIndexOfQChar mirrors the Qt overload without a `from` parameter
// (default from = -1, which Qt clamps to size()-1).
func (v QLatin1StringView) LastIndexOfQChar(c QChar, caseSensitive bool) int {
	return v.LastIndexOfQCharFrom(c, -1, caseSensitive)
}

// ---- contains ------------------------------------------------------------

func (v QLatin1StringView) ContainsView(other QLatin1StringView, caseSensitive bool) bool {
	r := C.qlatin1stringview_contains_view(
		(*C.char)(v.data), C.ptrdiff_t(v.size),
		(*C.char)(other.data), C.ptrdiff_t(other.size),
		csInt(caseSensitive))
	runtime.KeepAlive(v)
	runtime.KeepAlive(other)
	return bool(r)
}

func (v QLatin1StringView) ContainsQStringView(other QStringView, caseSensitive bool) bool {
	r := C.qlatin1stringview_contains_qstringview(
		(*C.char)(v.data), C.ptrdiff_t(v.size),
		(*C.uint16_t)(other.Utf16Unsafe()), C.ptrdiff_t(other.Size()),
		csInt(caseSensitive))
	runtime.KeepAlive(v)
	runtime.KeepAlive(other)
	return bool(r)
}

func (v QLatin1StringView) ContainsQChar(c QChar, caseSensitive bool) bool {
	r := C.qlatin1stringview_contains_qchar(
		(*C.char)(v.data), C.ptrdiff_t(v.size),
		C.uint16_t(c.Unicode()), csInt(caseSensitive))
	runtime.KeepAlive(v)
	return bool(r)
}

// ---- count ---------------------------------------------------------------

func (v QLatin1StringView) CountView(other QLatin1StringView, caseSensitive bool) int {
	r := C.qlatin1stringview_count_view(
		(*C.char)(v.data), C.ptrdiff_t(v.size),
		(*C.char)(other.data), C.ptrdiff_t(other.size),
		csInt(caseSensitive))
	runtime.KeepAlive(v)
	runtime.KeepAlive(other)
	return int(r)
}

func (v QLatin1StringView) CountQStringView(other QStringView, caseSensitive bool) int {
	r := C.qlatin1stringview_count_qstringview(
		(*C.char)(v.data), C.ptrdiff_t(v.size),
		(*C.uint16_t)(other.Utf16Unsafe()), C.ptrdiff_t(other.Size()),
		csInt(caseSensitive))
	runtime.KeepAlive(v)
	runtime.KeepAlive(other)
	return int(r)
}

func (v QLatin1StringView) CountQChar(c QChar, caseSensitive bool) int {
	r := C.qlatin1stringview_count_qchar(
		(*C.char)(v.data), C.ptrdiff_t(v.size),
		C.uint16_t(c.Unicode()), csInt(caseSensitive))
	runtime.KeepAlive(v)
	return int(r)
}

// ---- numeric conversion --------------------------------------------------

func (v QLatin1StringView) ToShort(base int) (int16, bool) {
	var ok C.bool
	r := C.qlatin1stringview_to_short((*C.char)(v.data), C.ptrdiff_t(v.size), &ok, C.int(base))
	runtime.KeepAlive(v)
	return int16(r), bool(ok)
}
func (v QLatin1StringView) ToUShort(base int) (uint16, bool) {
	var ok C.bool
	r := C.qlatin1stringview_to_ushort((*C.char)(v.data), C.ptrdiff_t(v.size), &ok, C.int(base))
	runtime.KeepAlive(v)
	return uint16(r), bool(ok)
}
func (v QLatin1StringView) ToInt(base int) (int32, bool) {
	var ok C.bool
	r := C.qlatin1stringview_to_int((*C.char)(v.data), C.ptrdiff_t(v.size), &ok, C.int(base))
	runtime.KeepAlive(v)
	return int32(r), bool(ok)
}
func (v QLatin1StringView) ToUInt(base int) (uint32, bool) {
	var ok C.bool
	r := C.qlatin1stringview_to_uint((*C.char)(v.data), C.ptrdiff_t(v.size), &ok, C.int(base))
	runtime.KeepAlive(v)
	return uint32(r), bool(ok)
}
func (v QLatin1StringView) ToLong(base int) (int64, bool) {
	var ok C.bool
	r := C.qlatin1stringview_to_long((*C.char)(v.data), C.ptrdiff_t(v.size), &ok, C.int(base))
	runtime.KeepAlive(v)
	return int64(r), bool(ok)
}
func (v QLatin1StringView) ToULong(base int) (uint64, bool) {
	var ok C.bool
	r := C.qlatin1stringview_to_ulong((*C.char)(v.data), C.ptrdiff_t(v.size), &ok, C.int(base))
	runtime.KeepAlive(v)
	return uint64(r), bool(ok)
}
func (v QLatin1StringView) ToLongLong(base int) (int64, bool) {
	var ok C.bool
	r := C.qlatin1stringview_to_longlong((*C.char)(v.data), C.ptrdiff_t(v.size), &ok, C.int(base))
	runtime.KeepAlive(v)
	return int64(r), bool(ok)
}
func (v QLatin1StringView) ToULongLong(base int) (uint64, bool) {
	var ok C.bool
	r := C.qlatin1stringview_to_ulonglong((*C.char)(v.data), C.ptrdiff_t(v.size), &ok, C.int(base))
	runtime.KeepAlive(v)
	return uint64(r), bool(ok)
}
func (v QLatin1StringView) ToFloat() (float32, bool) {
	var ok C.bool
	r := C.qlatin1stringview_to_float((*C.char)(v.data), C.ptrdiff_t(v.size), &ok)
	runtime.KeepAlive(v)
	return float32(r), bool(ok)
}
func (v QLatin1StringView) ToDouble() (float64, bool) {
	var ok C.bool
	r := C.qlatin1stringview_to_double((*C.char)(v.data), C.ptrdiff_t(v.size), &ok)
	runtime.KeepAlive(v)
	return float64(r), bool(ok)
}

// ---- conversions to owning QByteArray -----------------------------------

// ToUtf8 returns a fresh *QByteArray with UTF-8 encoding
// (mirrors QLatin1StringView::toUtf8()).
func (v QLatin1StringView) ToUtf8() *QByteArray {
	p := C.qlatin1stringview_to_utf8((*C.char)(v.data), C.ptrdiff_t(v.size))
	runtime.KeepAlive(v)
	return wrapQByteArray(p)
}

// MaxSize mirrors QLatin1StringView::maxSize().
func QLatin1StringViewMaxSize() int { return int(C.qlatin1stringview_max_size()) }

// NewUnsafeWithOwner builds a view over raw Latin-1 data and attaches an owner reference
// keeping the underlying buffer alive for as long as the view (or its slices) is referenced.
//
// Advanced usage: intended for sibling go-qt6core subpackages (e.g. qanystringview).
func NewQLatin1StringViewUnsafeWithOwner(data unsafe.Pointer, size int, owner any) QLatin1StringView {
	if size < 0 {
		panic("qlatin1stringview: NewUnsafeWithOwner negative size")
	}
	if size > 0 && data == nil {
		panic("qlatin1stringview: NewUnsafeWithOwner nil data with non-zero size")
	}
	return QLatin1StringView{data: data, size: size, keepAlive: owner}
}
