package qbytearrayview

/*
 В Qt есть first() (синоним front, deprecated) и first(qsizetype n) (subview). В Go перегрузки нет → используем Front()/Back() для байт, First(n)/Last(n) для subviews — как в qbytearray.

 Что сознательно не вошло
 STL-итераторы (begin/end/rbegin/rend) — не идиоматичны в Go; Bytes() — замена.

 operator std::string_view — не идиоматичен; String()/Bytes().

 Неявные конструкторы (const char*, std::string_view, контейнеры) — Go не поддерживает перегрузки; явные FromBytes/FromString/FromQByteArray покрывают 100% случаев.

 Операторы сравнения с QChar (Q_DECLARE_STRONGLY_ORDERED(..., QChar)) — добавятся вместе с пакетом qchar-bridge, если понадобится (сравнение view с одиночным QChar — редкий edge case).

 qChecksum(QByteArrayView, ...) — свободная функция из qbytearrayalgorithms.h; может быть добавлена отдельно, если понадобится.

 QByteArrayView::first() / last() (deprecated aliases) — не включены: их поведение = Front()/Back().
*/

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC -Wno-sfinae-incomplete
#cgo LDFLAGS: -lstdc++
#include <stdint.h>
#include <stdbool.h>
#include <stddef.h>
#include "qbytearrayview_c.h"
*/
import "C"

import (
	"runtime"
	"unsafe"

	"github.com/vadimlxx/go-qt6core/qbytearray"
)

// QByteArrayView is a read-only, non-owning view over a byte buffer.
// Mirrors Qt6 QByteArrayView (primitive value type).
//
// The zero value is a valid null view (size=0, data=nil).
// Use FromBytes / FromString / FromQByteArray to construct from real data.
type QByteArrayView struct {
	data      unsafe.Pointer
	size      int
	keepAlive any // keeps the underlying Go object ([]byte, string, *qbytearray.QByteArray) alive
}

// ---- construction --------------------------------------------------------

// Empty returns a null QByteArrayView (matches QByteArrayView()).
func Empty() QByteArrayView { return QByteArrayView{} }

// NewUnsafe builds a view directly from a raw pointer and size.
// The caller is solely responsible for keeping the underlying data alive
// and unmutated for as long as the view (or any of its slices) is used.
func NewUnsafe(data unsafe.Pointer, size int) QByteArrayView {
	if size < 0 {
		panic("qbytearrayview: NewUnsafe negative size")
	}
	return QByteArrayView{data: data, size: size}
}

// FromBytes returns a view over b. The Go slice b is kept alive by the view;
// mutations of b (e.g. b[0] = x) are visible through the view — same semantics
// as Qt QByteArrayView.
func FromBytes(b []byte) QByteArrayView {
	if len(b) == 0 {
		return Empty()
	}
	return QByteArrayView{
		data:      unsafe.Pointer(&b[0]),
		size:      len(b),
		keepAlive: b,
	}
}

// FromString returns a view over the bytes of s.
// Go strings are immutable, so the view is effectively read-only.
func FromString(s string) QByteArrayView {
	if len(s) == 0 {
		return Empty()
	}
	return QByteArrayView{
		data:      unsafe.Pointer(unsafe.StringData(s)),
		size:      len(s),
		keepAlive: s,
	}
}

// FromQByteArray returns a zero-copy view into q's internal buffer.
// The view keeps q alive. Do NOT mutate or Close q while the view exists.
func FromQByteArray(q *qbytearray.QByteArray) QByteArrayView {
	if q == nil || q.IsNull() {
		return Empty()
	}
	return QByteArrayView{
		data:      q.ConstDataUnsafe(),
		size:      q.Size(),
		keepAlive: q,
	}
}

// ---- accessors -----------------------------------------------------------

func (v QByteArrayView) Size() int  { return v.size }
func (v QByteArrayView) Len() int   { return v.size }
func (v QByteArrayView) IsEmpty() bool { return v.size == 0 }
func (v QByteArrayView) IsNull() bool  { return v.data == nil }

// DataUnsafe returns the underlying data pointer.
// Valid only while the view (and its keepAlive target) is alive.
func (v QByteArrayView) DataUnsafe() unsafe.Pointer { return v.data }

// Bytes returns a copy of the view's bytes.
func (v QByteArrayView) Bytes() []byte {
	if v.size == 0 {
		return nil
	}
	out := make([]byte, v.size)
	copy(out, unsafe.Slice((*byte)(v.data), v.size))
	runtime.KeepAlive(v)
	return out
}

// String returns a copy of the view's bytes as a Go string.
func (v QByteArrayView) String() string {
	var s string
	if v.size > 0 {
		s = string(unsafe.Slice((*byte)(v.data), v.size))
	}
	runtime.KeepAlive(v)
	return s
}

// At returns the i-th byte. Panics if i is out of range.
func (v QByteArrayView) At(i int) byte {
	if i < 0 || i >= v.size {
		panic("qbytearrayview: At out of range")
	}
	r := C.qbytearrayview_at((*C.char)(v.data), C.ptrdiff_t(v.size), C.ptrdiff_t(i))
	runtime.KeepAlive(v)
	return byte(r)
}

// Front returns the first byte. Panics if the view is empty.
func (v QByteArrayView) Front() byte {
	if v.size == 0 {
		panic("qbytearrayview: Front on empty view")
	}
	r := C.qbytearrayview_front((*C.char)(v.data), C.ptrdiff_t(v.size))
	runtime.KeepAlive(v)
	return byte(r)
}

// Back returns the last byte. Panics if the view is empty.
func (v QByteArrayView) Back() byte {
	if v.size == 0 {
		panic("qbytearrayview: Back on empty view")
	}
	r := C.qbytearrayview_back((*C.char)(v.data), C.ptrdiff_t(v.size))
	runtime.KeepAlive(v)
	return byte(r)
}

func (v QByteArrayView) IsValidUtf8() bool {
	r := C.qbytearrayview_is_valid_utf8((*C.char)(v.data), C.ptrdiff_t(v.size))
	runtime.KeepAlive(v)
	return bool(r)
}

// ---- search --------------------------------------------------------------

func (v QByteArrayView) IndexOf(c byte, from int) int {
	r := C.qbytearrayview_index_of_char((*C.char)(v.data), C.ptrdiff_t(v.size),
		C.char(c), C.ptrdiff_t(from))
	runtime.KeepAlive(v)
	return int(r)
}

func (v QByteArrayView) IndexOfBytes(needle []byte, from int) int {
	var p *C.char
	if len(needle) > 0 {
		p = (*C.char)(unsafe.Pointer(&needle[0]))
	}
	r := C.qbytearrayview_index_of_view((*C.char)(v.data), C.ptrdiff_t(v.size),
		p, C.ptrdiff_t(len(needle)), C.ptrdiff_t(from))
	runtime.KeepAlive(v)
	runtime.KeepAlive(needle)
	return int(r)
}

func (v QByteArrayView) LastIndexOf(c byte, from int) int {
	r := C.qbytearrayview_last_index_of_char((*C.char)(v.data), C.ptrdiff_t(v.size),
		C.char(c), C.ptrdiff_t(from))
	runtime.KeepAlive(v)
	return int(r)
}

func (v QByteArrayView) LastIndexOfBytes(needle []byte, from int) int {
	var p *C.char
	if len(needle) > 0 {
		p = (*C.char)(unsafe.Pointer(&needle[0]))
	}
	r := C.qbytearrayview_last_index_of_view((*C.char)(v.data), C.ptrdiff_t(v.size),
		p, C.ptrdiff_t(len(needle)), C.ptrdiff_t(from))
	runtime.KeepAlive(v)
	runtime.KeepAlive(needle)
	return int(r)
}

func (v QByteArrayView) Count(c byte) int {
	r := C.qbytearrayview_count_char((*C.char)(v.data), C.ptrdiff_t(v.size), C.char(c))
	runtime.KeepAlive(v)
	return int(r)
}

func (v QByteArrayView) CountBytes(needle []byte) int {
	var p *C.char
	if len(needle) > 0 {
		p = (*C.char)(unsafe.Pointer(&needle[0]))
	}
	r := C.qbytearrayview_count_view((*C.char)(v.data), C.ptrdiff_t(v.size),
		p, C.ptrdiff_t(len(needle)))
	runtime.KeepAlive(v)
	runtime.KeepAlive(needle)
	return int(r)
}

func (v QByteArrayView) Contains(c byte) bool {
	r := C.qbytearrayview_contains_char((*C.char)(v.data), C.ptrdiff_t(v.size), C.char(c))
	runtime.KeepAlive(v)
	return bool(r)
}

func (v QByteArrayView) ContainsBytes(needle []byte) bool {
	var p *C.char
	if len(needle) > 0 {
		p = (*C.char)(unsafe.Pointer(&needle[0]))
	}
	r := C.qbytearrayview_contains_view((*C.char)(v.data), C.ptrdiff_t(v.size),
		p, C.ptrdiff_t(len(needle)))
	runtime.KeepAlive(v)
	runtime.KeepAlive(needle)
	return bool(r)
}

func (v QByteArrayView) StartsWith(c byte) bool {
	r := C.qbytearrayview_starts_with_char((*C.char)(v.data), C.ptrdiff_t(v.size), C.char(c))
	runtime.KeepAlive(v)
	return bool(r)
}

func (v QByteArrayView) StartsWithBytes(prefix []byte) bool {
	var p *C.char
	if len(prefix) > 0 {
		p = (*C.char)(unsafe.Pointer(&prefix[0]))
	}
	r := C.qbytearrayview_starts_with_view((*C.char)(v.data), C.ptrdiff_t(v.size),
		p, C.ptrdiff_t(len(prefix)))
	runtime.KeepAlive(v)
	runtime.KeepAlive(prefix)
	return bool(r)
}

func (v QByteArrayView) EndsWith(c byte) bool {
	r := C.qbytearrayview_ends_with_char((*C.char)(v.data), C.ptrdiff_t(v.size), C.char(c))
	runtime.KeepAlive(v)
	return bool(r)
}

func (v QByteArrayView) EndsWithBytes(suffix []byte) bool {
	var p *C.char
	if len(suffix) > 0 {
		p = (*C.char)(unsafe.Pointer(&suffix[0]))
	}
	r := C.qbytearrayview_ends_with_view((*C.char)(v.data), C.ptrdiff_t(v.size),
		p, C.ptrdiff_t(len(suffix)))
	runtime.KeepAlive(v)
	runtime.KeepAlive(suffix)
	return bool(r)
}

// Compare returns Qt's compare() result. cs == true → case-sensitive
// (matches Qt's default Qt::CaseSensitive).
func (v QByteArrayView) Compare(other []byte, cs bool) int {
	var p *C.char
	if len(other) > 0 {
		p = (*C.char)(unsafe.Pointer(&other[0]))
	}
	csInt := C.int(0)
	if !cs {
		csInt = 1
	}
	r := C.qbytearrayview_compare((*C.char)(v.data), C.ptrdiff_t(v.size),
		p, C.ptrdiff_t(len(other)), csInt)
	runtime.KeepAlive(v)
	runtime.KeepAlive(other)
	return int(r)
}

// ---- slicing (return new views sharing the same buffer) -----------------

func viewOut(v QByteArrayView, cdata *C.char, csize C.ptrdiff_t) QByteArrayView {
	// not used directly; helper for compiler
	_, _ = cdata, csize
	return v
}

func (v QByteArrayView) sliceFromC(
	call func(data *C.char, size C.ptrdiff_t, outData **C.char, outSize *C.ptrdiff_t),
) QByteArrayView {
	var outData *C.char
	var outSize C.ptrdiff_t
	call((*C.char)(v.data), C.ptrdiff_t(v.size), &outData, &outSize)
	runtime.KeepAlive(v)
	return QByteArrayView{
		data:      unsafe.Pointer(outData),
		size:      int(outSize),
		keepAlive: v.keepAlive,
	}
}

func (v QByteArrayView) Left(n int) QByteArrayView {
	return v.sliceFromC(func(d *C.char, s C.ptrdiff_t, od **C.char, os *C.ptrdiff_t) {
		C.qbytearrayview_left(d, s, C.ptrdiff_t(n), od, os)
	})
}
func (v QByteArrayView) Right(n int) QByteArrayView {
	return v.sliceFromC(func(d *C.char, s C.ptrdiff_t, od **C.char, os *C.ptrdiff_t) {
		C.qbytearrayview_right(d, s, C.ptrdiff_t(n), od, os)
	})
}
func (v QByteArrayView) Mid(pos, n int) QByteArrayView {
	return v.sliceFromC(func(d *C.char, s C.ptrdiff_t, od **C.char, os *C.ptrdiff_t) {
		C.qbytearrayview_mid(d, s, C.ptrdiff_t(pos), C.ptrdiff_t(n), od, os)
	})
}
func (v QByteArrayView) First(n int) QByteArrayView {
	if n < 0 || n > v.size {
		panic("qbytearrayview: First out of range")
	}
	return v.sliceFromC(func(d *C.char, s C.ptrdiff_t, od **C.char, os *C.ptrdiff_t) {
		C.qbytearrayview_first(d, s, C.ptrdiff_t(n), od, os)
	})
}
func (v QByteArrayView) Last(n int) QByteArrayView {
	if n < 0 || n > v.size {
		panic("qbytearrayview: Last out of range")
	}
	return v.sliceFromC(func(d *C.char, s C.ptrdiff_t, od **C.char, os *C.ptrdiff_t) {
		C.qbytearrayview_last(d, s, C.ptrdiff_t(n), od, os)
	})
}
func (v QByteArrayView) Sliced(pos int) QByteArrayView {
	if pos < 0 || pos > v.size {
		panic("qbytearrayview: Sliced out of range")
	}
	return v.sliceFromC(func(d *C.char, s C.ptrdiff_t, od **C.char, os *C.ptrdiff_t) {
		C.qbytearrayview_sliced(d, s, C.ptrdiff_t(pos), od, os)
	})
}
func (v QByteArrayView) SlicedN(pos, n int) QByteArrayView {
	if pos < 0 || n < 0 || pos+n > v.size {
		panic("qbytearrayview: SlicedN out of range")
	}
	return v.sliceFromC(func(d *C.char, s C.ptrdiff_t, od **C.char, os *C.ptrdiff_t) {
		C.qbytearrayview_sliced_n(d, s, C.ptrdiff_t(pos), C.ptrdiff_t(n), od, os)
	})
}
func (v QByteArrayView) Chopped(n int) QByteArrayView {
	if n < 0 || n > v.size {
		panic("qbytearrayview: Chopped out of range")
	}
	return v.sliceFromC(func(d *C.char, s C.ptrdiff_t, od **C.char, os *C.ptrdiff_t) {
		C.qbytearrayview_chopped(d, s, C.ptrdiff_t(n), od, os)
	})
}
func (v QByteArrayView) Trimmed() QByteArrayView {
	return v.sliceFromC(func(d *C.char, s C.ptrdiff_t, od **C.char, os *C.ptrdiff_t) {
		C.qbytearrayview_trimmed(d, s, od, os)
	})
}

// Truncate shrinks the view to length n. Mutates the receiver in place,
// mirroring QByteArrayView::truncate(qsizetype).
func (v *QByteArrayView) Truncate(n int) {
	if n < 0 || n > v.size {
		panic("qbytearrayview: Truncate out of range")
	}
	v.size = n
}

// Chop removes n bytes from the end. Mutates the receiver in place.
func (v *QByteArrayView) Chop(n int) {
	if n < 0 || n > v.size {
		panic("qbytearrayview: Chop out of range")
	}
	v.size -= n
}

// Slice replaces the receiver with Sliced(pos).
func (v *QByteArrayView) Slice(pos int) { *v = v.Sliced(pos) }

// SliceN replaces the receiver with SlicedN(pos, n).
func (v *QByteArrayView) SliceN(pos, n int) { *v = v.SlicedN(pos, n) }

// ---- numeric conversion --------------------------------------------------

func (v QByteArrayView) ToShort(base int) (int16, bool) {
	var ok C.bool
	r := C.qbytearrayview_to_short((*C.char)(v.data), C.ptrdiff_t(v.size), &ok, C.int(base))
	runtime.KeepAlive(v)
	return int16(r), bool(ok)
}
func (v QByteArrayView) ToUShort(base int) (uint16, bool) {
	var ok C.bool
	r := C.qbytearrayview_to_ushort((*C.char)(v.data), C.ptrdiff_t(v.size), &ok, C.int(base))
	runtime.KeepAlive(v)
	return uint16(r), bool(ok)
}
func (v QByteArrayView) ToInt(base int) (int32, bool) {
	var ok C.bool
	r := C.qbytearrayview_to_int((*C.char)(v.data), C.ptrdiff_t(v.size), &ok, C.int(base))
	runtime.KeepAlive(v)
	return int32(r), bool(ok)
}
func (v QByteArrayView) ToUInt(base int) (uint32, bool) {
	var ok C.bool
	r := C.qbytearrayview_to_uint((*C.char)(v.data), C.ptrdiff_t(v.size), &ok, C.int(base))
	runtime.KeepAlive(v)
	return uint32(r), bool(ok)
}
func (v QByteArrayView) ToLongLong(base int) (int64, bool) {
	var ok C.bool
	r := C.qbytearrayview_to_longlong((*C.char)(v.data), C.ptrdiff_t(v.size), &ok, C.int(base))
	runtime.KeepAlive(v)
	return int64(r), bool(ok)
}
func (v QByteArrayView) ToULongLong(base int) (uint64, bool) {
	var ok C.bool
	r := C.qbytearrayview_to_ulonglong((*C.char)(v.data), C.ptrdiff_t(v.size), &ok, C.int(base))
	runtime.KeepAlive(v)
	return uint64(r), bool(ok)
}
func (v QByteArrayView) ToFloat() (float32, bool) {
	var ok C.bool
	r := C.qbytearrayview_to_float((*C.char)(v.data), C.ptrdiff_t(v.size), &ok)
	runtime.KeepAlive(v)
	return float32(r), bool(ok)
}
func (v QByteArrayView) ToDouble() (float64, bool) {
	var ok C.bool
	r := C.qbytearrayview_to_double((*C.char)(v.data), C.ptrdiff_t(v.size), &ok)
	runtime.KeepAlive(v)
	return float64(r), bool(ok)
}

// ---- conversion to owning QByteArray ------------------------------------

// ToByteArray returns a fresh, owned *qbytearray.QByteArray with a copy
// of the view's bytes (mirrors QByteArrayView::toByteArray()).
func (v QByteArrayView) ToByteArray() *qbytearray.QByteArray {
	var out unsafe.Pointer
	C.qbytearrayview_to_bytearray((*C.char)(v.data), C.ptrdiff_t(v.size), &out)
	runtime.KeepAlive(v)
	return qbytearray.AdoptUnsafePtr(out)
}

// ---- static --------------------------------------------------------------

// MaxSize mirrors QByteArrayView::maxSize().
func MaxSize() int { return int(C.qbytearrayview_max_size()) }
