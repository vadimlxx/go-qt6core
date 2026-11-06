package qstringview

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC -Wno-sfinae-incomplete
#cgo LDFLAGS: -lstdc++
#include <stdint.h>
#include <stdbool.h>
#include <stddef.h>
#include "qstringview_c.h"
*/
import "C"

import (
	"runtime"
	"unicode/utf16"
	"unsafe"

	"github.com/vadimlxx/go-qt6core/qbytearray"
	"github.com/vadimlxx/go-qt6core/qchar"
)

// QStringView is a read-only, non-owning view over a UTF-16 buffer.
// Mirrors Qt6 QStringView (primitive value type).
//
// The zero value is a valid null view (size=0, data=nil).
// Use FromUtf16 / FromQChars / FromString to construct from real data.
type QStringView struct {
	data      unsafe.Pointer // const char16_t*
	size      int
	keepAlive any // holds the underlying []uint16 / []qchar.QChar / string
}

// ---- construction --------------------------------------------------------

// Empty returns a null QStringView (matches QStringView()).
func Empty() QStringView { return QStringView{} }

// NewUnsafe builds a view directly from a raw UTF-16 pointer and size.
// Caller is solely responsible for keeping the underlying data alive
// and unmutated for as long as the view (or any of its slices) is used.
func NewUnsafe(data unsafe.Pointer, size int) QStringView {
	if size < 0 {
		panic("qstringview: NewUnsafe negative size")
	}
	return QStringView{data: data, size: size}
}

// FromUtf16 returns a view over units. The Go slice is kept alive by
// the view; mutations of units are visible through the view.
func FromUtf16(units []uint16) QStringView {
	if len(units) == 0 {
		return Empty()
	}
	return QStringView{
		data:      unsafe.Pointer(&units[0]),
		size:      len(units),
		keepAlive: units,
	}
}

// FromQChars returns a view over chars. Zero-copy: qchar.QChar is a
// 16-bit type in Go, so []qchar.QChar shares layout with []uint16.
// The slice is kept alive by the view.
func FromQChars(chars []qchar.QChar) QStringView {
	if len(chars) == 0 {
		return Empty()
	}
	return QStringView{
		data:      unsafe.Pointer(&chars[0]),
		size:      len(chars),
		keepAlive: chars,
	}
}

// FromString converts s (UTF-8) into UTF-16 and returns a view over the
// resulting units. The conversion allocates; the view owns and keeps the
// converted buffer alive. For zero-copy from existing UTF-16 data, use
// FromUtf16 directly.
func FromString(s string) QStringView {
	if len(s) == 0 {
		return Empty()
	}
	units := utf16.Encode([]rune(s))
	return QStringView{
		data:      unsafe.Pointer(&units[0]),
		size:      len(units),
		keepAlive: units,
	}
}

// NewUnsafeWithOwner builds a view over raw UTF-16 data and attaches an
// owner reference that keeps the underlying buffer alive for as long as
// the view (or its slices) is referenced.
//
// Advanced usage: intended for sibling go-qt6core subpackages
// (e.g. qstring) that need a zero-copy view over their internal storage.
// The owner is typically the *QString that provided data.
func NewUnsafeWithOwner(data unsafe.Pointer, size int, owner any) QStringView {
	if size < 0 {
		panic("qstringview: NewUnsafeWithOwner negative size")
	}
	if size > 0 && data == nil {
		panic("qstringview: NewUnsafeWithOwner nil data with non-zero size")
	}
	return QStringView{data: data, size: size, keepAlive: owner}
}

// ---- accessors -----------------------------------------------------------

func (v QStringView) Size() int    { return v.size }
func (v QStringView) Len() int     { return v.size }
func (v QStringView) IsEmpty() bool { return v.size == 0 }
func (v QStringView) IsNull() bool  { return v.data == nil }

// Utf16Unsafe returns the raw UTF-16 data pointer.
// Valid only while the view (and its keepAlive target) is alive and unmutated.
func (v QStringView) Utf16Unsafe() unsafe.Pointer { return v.data }

// Utf16 returns a copy of the view's code units.
func (v QStringView) Utf16() []uint16 {
	if v.size == 0 {
		return nil
	}
	out := make([]uint16, v.size)
	copy(out, unsafe.Slice((*uint16)(v.data), v.size))
	runtime.KeepAlive(v)
	return out
}

// GoString decodes the UTF-16 code units into a Go string.
func (v QStringView) GoString() string {
	if v.size == 0 {
		return ""
	}
	units := unsafe.Slice((*uint16)(v.data), v.size)
	s := string(utf16.Decode(units))
	runtime.KeepAlive(v)
	return s
}

// At returns the i-th QChar. Panics if i is out of range.
func (v QStringView) At(i int) qchar.QChar {
	if i < 0 || i >= v.size {
		panic("qstringview: At out of range")
	}
	r := C.qstringview_at((*C.uint16_t)(v.data), C.ptrdiff_t(v.size), C.ptrdiff_t(i))
	runtime.KeepAlive(v)
	return qchar.New(uint16(r))
}

// Front returns the first QChar. Panics if the view is empty.
func (v QStringView) Front() qchar.QChar {
	if v.size == 0 {
		panic("qstringview: Front on empty view")
	}
	r := C.qstringview_front((*C.uint16_t)(v.data), C.ptrdiff_t(v.size))
	runtime.KeepAlive(v)
	return qchar.New(uint16(r))
}

// Back returns the last QChar. Panics if the view is empty.
func (v QStringView) Back() qchar.QChar {
	if v.size == 0 {
		panic("qstringview: Back on empty view")
	}
	r := C.qstringview_back((*C.uint16_t)(v.data), C.ptrdiff_t(v.size))
	runtime.KeepAlive(v)
	return qchar.New(uint16(r))
}

// ---- predicates ----------------------------------------------------------

func (v QStringView) IsUpper() bool {
	r := C.qstringview_is_upper((*C.uint16_t)(v.data), C.ptrdiff_t(v.size))
	runtime.KeepAlive(v)
	return bool(r)
}
func (v QStringView) IsLower() bool {
	r := C.qstringview_is_lower((*C.uint16_t)(v.data), C.ptrdiff_t(v.size))
	runtime.KeepAlive(v)
	return bool(r)
}
func (v QStringView) IsRightToLeft() bool {
	r := C.qstringview_is_right_to_left((*C.uint16_t)(v.data), C.ptrdiff_t(v.size))
	runtime.KeepAlive(v)
	return bool(r)
}
func (v QStringView) IsValidUtf16() bool {
	r := C.qstringview_is_valid_utf16((*C.uint16_t)(v.data), C.ptrdiff_t(v.size))
	runtime.KeepAlive(v)
	return bool(r)
}

// ---- slicing (return new views over the same buffer) --------------------

func (v QStringView) sliceFromC(
	call func(data *C.uint16_t, size C.ptrdiff_t, od **C.uint16_t, os *C.ptrdiff_t),
) QStringView {
	var outData *C.uint16_t
	var outSize C.ptrdiff_t
	call((*C.uint16_t)(v.data), C.ptrdiff_t(v.size), &outData, &outSize)
	runtime.KeepAlive(v)
	return QStringView{
		data:      unsafe.Pointer(outData),
		size:      int(outSize),
		keepAlive: v.keepAlive,
	}
}

// Left clamps n (Qt semantics: n < 0 or n > size → whole view).
func (v QStringView) Left(n int) QStringView {
	return v.sliceFromC(func(d *C.uint16_t, s C.ptrdiff_t, od **C.uint16_t, os *C.ptrdiff_t) {
		C.qstringview_left(d, s, C.ptrdiff_t(n), od, os)
	})
}

// Right clamps n (Qt semantics).
func (v QStringView) Right(n int) QStringView {
	return v.sliceFromC(func(d *C.uint16_t, s C.ptrdiff_t, od **C.uint16_t, os *C.ptrdiff_t) {
		C.qstringview_right(d, s, C.ptrdiff_t(n), od, os)
	})
}

// Mid clamps pos/n (Qt semantics; n = -1 → до конца).
func (v QStringView) Mid(pos, n int) QStringView {
	return v.sliceFromC(func(d *C.uint16_t, s C.ptrdiff_t, od **C.uint16_t, os *C.ptrdiff_t) {
		C.qstringview_mid(d, s, C.ptrdiff_t(pos), C.ptrdiff_t(n), od, os)
	})
}

func (v QStringView) First(n int) QStringView {
	if n < 0 || n > v.size {
		panic("qstringview: First out of range")
	}
	return v.sliceFromC(func(d *C.uint16_t, s C.ptrdiff_t, od **C.uint16_t, os *C.ptrdiff_t) {
		C.qstringview_first(d, s, C.ptrdiff_t(n), od, os)
	})
}

func (v QStringView) Last(n int) QStringView {
	if n < 0 || n > v.size {
		panic("qstringview: Last out of range")
	}
	return v.sliceFromC(func(d *C.uint16_t, s C.ptrdiff_t, od **C.uint16_t, os *C.ptrdiff_t) {
		C.qstringview_last(d, s, C.ptrdiff_t(n), od, os)
	})
}

func (v QStringView) Sliced(pos int) QStringView {
	if pos < 0 || pos > v.size {
		panic("qstringview: Sliced out of range")
	}
	return v.sliceFromC(func(d *C.uint16_t, s C.ptrdiff_t, od **C.uint16_t, os *C.ptrdiff_t) {
		C.qstringview_sliced(d, s, C.ptrdiff_t(pos), od, os)
	})
}

func (v QStringView) SlicedN(pos, n int) QStringView {
	if pos < 0 || n < 0 || pos+n > v.size {
		panic("qstringview: SlicedN out of range")
	}
	return v.sliceFromC(func(d *C.uint16_t, s C.ptrdiff_t, od **C.uint16_t, os *C.ptrdiff_t) {
		C.qstringview_sliced_n(d, s, C.ptrdiff_t(pos), C.ptrdiff_t(n), od, os)
	})
}

func (v QStringView) Chopped(n int) QStringView {
	if n < 0 || n > v.size {
		panic("qstringview: Chopped out of range")
	}
	return v.sliceFromC(func(d *C.uint16_t, s C.ptrdiff_t, od **C.uint16_t, os *C.ptrdiff_t) {
		C.qstringview_chopped(d, s, C.ptrdiff_t(n), od, os)
	})
}

func (v QStringView) Trimmed() QStringView {
	return v.sliceFromC(func(d *C.uint16_t, s C.ptrdiff_t, od **C.uint16_t, os *C.ptrdiff_t) {
		C.qstringview_trimmed(d, s, od, os)
	})
}

// Truncate shrinks the view to length n in place.
func (v *QStringView) Truncate(n int) {
	if n < 0 || n > v.size {
		panic("qstringview: Truncate out of range")
	}
	v.size = n
}

// Chop removes n code units from the end in place.
func (v *QStringView) Chop(n int) {
	if n < 0 || n > v.size {
		panic("qstringview: Chop out of range")
	}
	v.size -= n
}

// Slice replaces the receiver with Sliced(pos).
func (v *QStringView) Slice(pos int) { *v = v.Sliced(pos) }

// SliceN replaces the receiver with SlicedN(pos, n).
func (v *QStringView) SliceN(pos, n int) { *v = v.SlicedN(pos, n) }

// ---- search --------------------------------------------------------------

func csInt(cs bool) C.int {
	if cs {
		return 0
	}
	return 1
}

func (v QStringView) IndexOf(c qchar.QChar, from int, caseSensitive bool) int {
	r := C.qstringview_index_of_qchar((*C.uint16_t)(v.data), C.ptrdiff_t(v.size),
		C.uint16_t(c.Unicode()), C.ptrdiff_t(from), csInt(caseSensitive))
	runtime.KeepAlive(v)
	return int(r)
}

func (v QStringView) IndexOfView(needle QStringView, from int, caseSensitive bool) int {
	r := C.qstringview_index_of_view((*C.uint16_t)(v.data), C.ptrdiff_t(v.size),
		(*C.uint16_t)(needle.data), C.ptrdiff_t(needle.size),
		C.ptrdiff_t(from), csInt(caseSensitive))
	runtime.KeepAlive(v)
	runtime.KeepAlive(needle)
	return int(r)
}

func (v QStringView) LastIndexOf(c qchar.QChar, from int, caseSensitive bool) int {
	r := C.qstringview_last_index_of_qchar((*C.uint16_t)(v.data), C.ptrdiff_t(v.size),
		C.uint16_t(c.Unicode()), C.ptrdiff_t(from), csInt(caseSensitive))
	runtime.KeepAlive(v)
	return int(r)
}

func (v QStringView) LastIndexOfView(needle QStringView, from int, caseSensitive bool) int {
	r := C.qstringview_last_index_of_view((*C.uint16_t)(v.data), C.ptrdiff_t(v.size),
		(*C.uint16_t)(needle.data), C.ptrdiff_t(needle.size),
		C.ptrdiff_t(from), csInt(caseSensitive))
	runtime.KeepAlive(v)
	runtime.KeepAlive(needle)
	return int(r)
}

func (v QStringView) Count(c qchar.QChar, caseSensitive bool) int {
	r := C.qstringview_count_qchar((*C.uint16_t)(v.data), C.ptrdiff_t(v.size),
		C.uint16_t(c.Unicode()), csInt(caseSensitive))
	runtime.KeepAlive(v)
	return int(r)
}

func (v QStringView) CountView(needle QStringView, caseSensitive bool) int {
	r := C.qstringview_count_view((*C.uint16_t)(v.data), C.ptrdiff_t(v.size),
		(*C.uint16_t)(needle.data), C.ptrdiff_t(needle.size), csInt(caseSensitive))
	runtime.KeepAlive(v)
	runtime.KeepAlive(needle)
	return int(r)
}

func (v QStringView) Contains(c qchar.QChar, caseSensitive bool) bool {
	r := C.qstringview_contains_qchar((*C.uint16_t)(v.data), C.ptrdiff_t(v.size),
		C.uint16_t(c.Unicode()), csInt(caseSensitive))
	runtime.KeepAlive(v)
	return bool(r)
}

func (v QStringView) ContainsView(needle QStringView, caseSensitive bool) bool {
	r := C.qstringview_contains_view((*C.uint16_t)(v.data), C.ptrdiff_t(v.size),
		(*C.uint16_t)(needle.data), C.ptrdiff_t(needle.size), csInt(caseSensitive))
	runtime.KeepAlive(v)
	runtime.KeepAlive(needle)
	return bool(r)
}

func (v QStringView) StartsWith(c qchar.QChar, caseSensitive bool) bool {
	r := C.qstringview_starts_with_qchar((*C.uint16_t)(v.data), C.ptrdiff_t(v.size),
		C.uint16_t(c.Unicode()), csInt(caseSensitive))
	runtime.KeepAlive(v)
	return bool(r)
}

func (v QStringView) StartsWithView(prefix QStringView, caseSensitive bool) bool {
	r := C.qstringview_starts_with_view((*C.uint16_t)(v.data), C.ptrdiff_t(v.size),
		(*C.uint16_t)(prefix.data), C.ptrdiff_t(prefix.size), csInt(caseSensitive))
	runtime.KeepAlive(v)
	runtime.KeepAlive(prefix)
	return bool(r)
}

func (v QStringView) EndsWith(c qchar.QChar, caseSensitive bool) bool {
	r := C.qstringview_ends_with_qchar((*C.uint16_t)(v.data), C.ptrdiff_t(v.size),
		C.uint16_t(c.Unicode()), csInt(caseSensitive))
	runtime.KeepAlive(v)
	return bool(r)
}

func (v QStringView) EndsWithView(suffix QStringView, caseSensitive bool) bool {
	r := C.qstringview_ends_with_view((*C.uint16_t)(v.data), C.ptrdiff_t(v.size),
		(*C.uint16_t)(suffix.data), C.ptrdiff_t(suffix.size), csInt(caseSensitive))
	runtime.KeepAlive(v)
	runtime.KeepAlive(suffix)
	return bool(r)
}

// ---- compare -------------------------------------------------------------

func (v QStringView) Compare(other QStringView, caseSensitive bool) int {
	r := C.qstringview_compare_view((*C.uint16_t)(v.data), C.ptrdiff_t(v.size),
		(*C.uint16_t)(other.data), C.ptrdiff_t(other.size), csInt(caseSensitive))
	runtime.KeepAlive(v)
	runtime.KeepAlive(other)
	return int(r)
}

func (v QStringView) CompareQChar(c qchar.QChar, caseSensitive bool) int {
	r := C.qstringview_compare_qchar((*C.uint16_t)(v.data), C.ptrdiff_t(v.size),
		C.uint16_t(c.Unicode()), csInt(caseSensitive))
	runtime.KeepAlive(v)
	return int(r)
}

// ---- conversions to owning QByteArray -----------------------------------

// ToUtf8 returns a fresh *qbytearray.QByteArray with UTF-8 encoding.
func (v QStringView) ToUtf8() *qbytearray.QByteArray {
	var out unsafe.Pointer
	C.qstringview_to_utf8((*C.uint16_t)(v.data), C.ptrdiff_t(v.size), &out)
	runtime.KeepAlive(v)
	return qbytearray.AdoptUnsafePtr(out)
}

// ToLatin1 returns a fresh *qbytearray.QByteArray with Latin-1 encoding
// (chars outside Latin-1 become '?').
func (v QStringView) ToLatin1() *qbytearray.QByteArray {
	var out unsafe.Pointer
	C.qstringview_to_latin1((*C.uint16_t)(v.data), C.ptrdiff_t(v.size), &out)
	runtime.KeepAlive(v)
	return qbytearray.AdoptUnsafePtr(out)
}

// ToLocal8Bit returns a fresh *qbytearray.QByteArray with local 8-bit encoding.
func (v QStringView) ToLocal8Bit() *qbytearray.QByteArray {
	var out unsafe.Pointer
	C.qstringview_to_local8bit((*C.uint16_t)(v.data), C.ptrdiff_t(v.size), &out)
	runtime.KeepAlive(v)
	return qbytearray.AdoptUnsafePtr(out)
}

// ---- numeric conversion --------------------------------------------------

func (v QStringView) ToShort(base int) (int16, bool) {
	var ok C.bool
	r := C.qstringview_to_short((*C.uint16_t)(v.data), C.ptrdiff_t(v.size), &ok, C.int(base))
	runtime.KeepAlive(v)
	return int16(r), bool(ok)
}
func (v QStringView) ToUShort(base int) (uint16, bool) {
	var ok C.bool
	r := C.qstringview_to_ushort((*C.uint16_t)(v.data), C.ptrdiff_t(v.size), &ok, C.int(base))
	runtime.KeepAlive(v)
	return uint16(r), bool(ok)
}
func (v QStringView) ToInt(base int) (int32, bool) {
	var ok C.bool
	r := C.qstringview_to_int((*C.uint16_t)(v.data), C.ptrdiff_t(v.size), &ok, C.int(base))
	runtime.KeepAlive(v)
	return int32(r), bool(ok)
}
func (v QStringView) ToUInt(base int) (uint32, bool) {
	var ok C.bool
	r := C.qstringview_to_uint((*C.uint16_t)(v.data), C.ptrdiff_t(v.size), &ok, C.int(base))
	runtime.KeepAlive(v)
	return uint32(r), bool(ok)
}
func (v QStringView) ToLong(base int) (int64, bool) {
	var ok C.bool
	r := C.qstringview_to_long((*C.uint16_t)(v.data), C.ptrdiff_t(v.size), &ok, C.int(base))
	runtime.KeepAlive(v)
	return int64(r), bool(ok)
}
func (v QStringView) ToULong(base int) (uint64, bool) {
	var ok C.bool
	r := C.qstringview_to_ulong((*C.uint16_t)(v.data), C.ptrdiff_t(v.size), &ok, C.int(base))
	runtime.KeepAlive(v)
	return uint64(r), bool(ok)
}
func (v QStringView) ToLongLong(base int) (int64, bool) {
	var ok C.bool
	r := C.qstringview_to_longlong((*C.uint16_t)(v.data), C.ptrdiff_t(v.size), &ok, C.int(base))
	runtime.KeepAlive(v)
	return int64(r), bool(ok)
}
func (v QStringView) ToULongLong(base int) (uint64, bool) {
	var ok C.bool
	r := C.qstringview_to_ulonglong((*C.uint16_t)(v.data), C.ptrdiff_t(v.size), &ok, C.int(base))
	runtime.KeepAlive(v)
	return uint64(r), bool(ok)
}
func (v QStringView) ToFloat() (float32, bool) {
	var ok C.bool
	r := C.qstringview_to_float((*C.uint16_t)(v.data), C.ptrdiff_t(v.size), &ok)
	runtime.KeepAlive(v)
	return float32(r), bool(ok)
}
func (v QStringView) ToDouble() (float64, bool) {
	var ok C.bool
	r := C.qstringview_to_double((*C.uint16_t)(v.data), C.ptrdiff_t(v.size), &ok)
	runtime.KeepAlive(v)
	return float64(r), bool(ok)
}

// MaxSize mirrors QStringView::maxSize().
func MaxSize() int { return int(C.qstringview_max_size()) }
