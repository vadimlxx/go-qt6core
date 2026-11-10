package qtext

/*
 At/Front/Back → byte	В Qt storage_type = char; Go byte = uint8. Байтовая семантика корректна — это UTF-8 view, не QChar

 Go-panic на out-of-range (At/First/Last/Sliced/SlicedN/Chopped/Truncate/Chop)	Qt release отключает Q_ASSERT → UB. Left/Right/Mid — пробрасывают клампы Qt.

 что не вошло:

 arg() — variadic template.

 Reverse-iterators (rbegin/rend/crbegin/crend).

 operator std::string_view / operator std::u8string_view — не идиоматичны в Go.

 QBasicUtf8StringView<true> (char8_t-based variant) — отдельный биндинг, если нужен.

 fromArray() — тривиально выражается через FromBytes.

 Конструкторы из контейнеров — покрываются From*-хелперами.

 Симметричные методы в qstringview/qlatin1stringview/qstring (избегаем цикла импортов).
*/

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC -Wno-sfinae-incomplete
#cgo LDFLAGS: -lstdc++
#include <stdint.h>
#include <stdbool.h>
#include <stddef.h>
#include "qutf8stringview_c.h"
*/
import "C"

import (
	"runtime"
	"unsafe"
)

// QUtf8StringView is a read-only, non-owning view over a UTF-8 buffer.
// Mirrors Qt6 QUtf8StringView (= QBasicUtf8StringView<false>; Q_PRIMITIVE_TYPE).
//
// The zero value is a valid null view (size=0, data=nil).
type QUtf8StringView struct {
	data      unsafe.Pointer // const char*
	size      int
	keepAlive any // holds the underlying []byte / string / *QByteArray / etc.
}

// ---- construction --------------------------------------------------------

// NewUnsafe builds a view directly from a raw UTF-8 byte pointer and size.
// Caller is solely responsible for keeping the underlying data alive and
// unmutated for as long as the view (or any of its slices) is used.
// func NewUnsafe(data unsafe.Pointer, size int) QUtf8StringView {
// 	if size < 0 {
// 		panic("qutf8stringview: NewUnsafe negative size")
// 	}
// 	if size > 0 && data == nil {
// 		panic("qutf8stringview: NewUnsafe nil data with non-zero size")
// 	}
// 	return QUtf8StringView{data: data, size: size}
// }

// FromBytes returns a zero-copy view over b, interpreted as UTF-8.
// The Go slice b is kept alive by the view.
// Use IsValidUtf8() to verify the buffer contains valid UTF-8.
func QUtf8StringViewFromBytes(b []byte) QUtf8StringView {
	if len(b) == 0 {
		return QUtf8StringView{}
	}
	return QUtf8StringView{
		data:      unsafe.Pointer(&b[0]),
		size:      len(b),
		keepAlive: b,
	}
}

// FromString returns a zero-copy view over s interpreted as UTF-8.
// Go strings are immutable, so the view is effectively read-only.
// Use IsValidUtf8() to verify the buffer contains valid UTF-8.
func QUtf8StringViewFromString(s string) QUtf8StringView {
	if len(s) == 0 {
		return QUtf8StringView{}
	}
	return QUtf8StringView{
		data:      unsafe.Pointer(unsafe.StringData(s)),
		size:      len(s),
		keepAlive: s,
	}
}

// FromQByteArray returns a zero-copy view over q's internal buffer,
// interpreting the bytes as UTF-8. The view keeps q alive.
// Do NOT mutate or Close q while the view exists.
func QUtf8StringViewFromQByteArray(q *QByteArray) QUtf8StringView {
	if q == nil || q.IsNull() {
		return QUtf8StringView{}
	}
	return QUtf8StringView{
		data:      q.ConstDataUnsafe(),
		size:      q.Size(),
		keepAlive: q,
	}
}

// FromQString returns a UTF-8 view over q's contents.
//
// This allocates a temporary UTF-8 QByteArray internally and keeps it alive
// via the view. Use q.ToQStringView() for zero-copy UTF-16 access instead.
func QUtf8StringViewFromQString(q *QString) QUtf8StringView {
	if q == nil || q.IsNull() {
		return QUtf8StringView{}
	}
	ba := q.Utf8()
	return QUtf8StringView{
		data:      ba.ConstDataUnsafe(),
		size:      ba.Size(),
		keepAlive: ba,
	}
}

func NewQUtf8StringView(arg ...any) QUtf8StringView {
	if len(arg) == 0 || arg[0] == nil {
		return QUtf8StringView{}
	}

	switch arg[0].(type) {
	case []byte:
		return QUtf8StringViewFromBytes(arg[0].([]byte))
	case string:
		return QUtf8StringViewFromString(arg[0].(string))
	case QByteArray:
		ba := arg[0].(QByteArray)
		return QUtf8StringViewFromQByteArray(&ba)
	case *QByteArray:
		return QUtf8StringViewFromQByteArray(arg[0].(*QByteArray))
	case QString:
		s := arg[0].(QString)
		return QUtf8StringViewFromQString(&s)
	case *QString:
		return QUtf8StringViewFromQString(arg[0].(*QString))
	}

	panic("NewQUtf8StringView invalid argumentes")
}

// ---- accessors -----------------------------------------------------------

func (v QUtf8StringView) Size() int     { return v.size }
func (v QUtf8StringView) Len() int      { return v.size }
func (v QUtf8StringView) IsEmpty() bool { return v.size == 0 }
func (v QUtf8StringView) IsNull() bool  { return v.data == nil }

// DataUnsafe returns the raw UTF-8 data pointer.
// Valid only while the view (and its keepAlive target) is alive and unmutated.
func (v QUtf8StringView) DataUnsafe() unsafe.Pointer { return v.data }

// Bytes returns a copy of the view's bytes.
func (v QUtf8StringView) Bytes() []byte {
	if v.size == 0 {
		return nil
	}
	out := make([]byte, v.size)
	copy(out, unsafe.Slice((*byte)(v.data), v.size))
	runtime.KeepAlive(v)
	return out
}

// String returns a copy of the view's bytes as a Go string.
func (v QUtf8StringView) String() string {
	var s string
	if v.size > 0 {
		s = string(unsafe.Slice((*byte)(v.data), v.size))
	}
	runtime.KeepAlive(v)
	return s
}

// At returns the i-th byte. Panics if i is out of range.
func (v QUtf8StringView) At(i int) byte {
	if i < 0 || i >= v.size {
		panic("qutf8stringview: At out of range")
	}
	r := C.qutf8stringview_at((*C.char)(v.data), C.ptrdiff_t(v.size), C.ptrdiff_t(i))
	runtime.KeepAlive(v)
	return byte(r)
}

// Front returns the first byte. Panics if empty.
func (v QUtf8StringView) Front() byte {
	if v.size == 0 {
		panic("qutf8stringview: Front on empty view")
	}
	r := C.qutf8stringview_front((*C.char)(v.data), C.ptrdiff_t(v.size))
	runtime.KeepAlive(v)
	return byte(r)
}

// Back returns the last byte. Panics if empty.
func (v QUtf8StringView) Back() byte {
	if v.size == 0 {
		panic("qutf8stringview: Back on empty view")
	}
	r := C.qutf8stringview_back((*C.char)(v.data), C.ptrdiff_t(v.size))
	runtime.KeepAlive(v)
	return byte(r)
}

// IsValidUtf8 reports whether the view's bytes form valid UTF-8.
func (v QUtf8StringView) IsValidUtf8() bool {
	r := C.qutf8stringview_is_valid_utf8((*C.char)(v.data), C.ptrdiff_t(v.size))
	runtime.KeepAlive(v)
	return bool(r)
}

// ---- slicing (return new views over the same buffer) --------------------

func (v QUtf8StringView) sliceFromC(
	call func(data *C.char, size C.ptrdiff_t, od **C.char, os *C.ptrdiff_t),
) QUtf8StringView {
	var outData *C.char
	var outSize C.ptrdiff_t
	call((*C.char)(v.data), C.ptrdiff_t(v.size), &outData, &outSize)
	runtime.KeepAlive(v)
	return QUtf8StringView{
		data:      unsafe.Pointer(outData),
		size:      int(outSize),
		keepAlive: v.keepAlive,
	}
}

// Left clamps n (Qt semantics).
func (v QUtf8StringView) Left(n int) QUtf8StringView {
	return v.sliceFromC(func(d *C.char, s C.ptrdiff_t, od **C.char, os *C.ptrdiff_t) {
		C.qutf8stringview_left(d, s, C.ptrdiff_t(n), od, os)
	})
}

// Right clamps n (Qt semantics).
func (v QUtf8StringView) Right(n int) QUtf8StringView {
	return v.sliceFromC(func(d *C.char, s C.ptrdiff_t, od **C.char, os *C.ptrdiff_t) {
		C.qutf8stringview_right(d, s, C.ptrdiff_t(n), od, os)
	})
}

// Mid clamps pos/n (Qt semantics; n = -1 → до конца).
func (v QUtf8StringView) Mid(pos, n int) QUtf8StringView {
	return v.sliceFromC(func(d *C.char, s C.ptrdiff_t, od **C.char, os *C.ptrdiff_t) {
		C.qutf8stringview_mid(d, s, C.ptrdiff_t(pos), C.ptrdiff_t(n), od, os)
	})
}

func (v QUtf8StringView) Sliced(pos int) QUtf8StringView {
	if pos < 0 || pos > v.size {
		panic("qutf8stringview: Sliced out of range")
	}
	return v.sliceFromC(func(d *C.char, s C.ptrdiff_t, od **C.char, os *C.ptrdiff_t) {
		C.qutf8stringview_sliced(d, s, C.ptrdiff_t(pos), od, os)
	})
}

func (v QUtf8StringView) SlicedN(pos, n int) QUtf8StringView {
	if pos < 0 || n < 0 || pos+n > v.size {
		panic("qutf8stringview: SlicedN out of range")
	}
	return v.sliceFromC(func(d *C.char, s C.ptrdiff_t, od **C.char, os *C.ptrdiff_t) {
		C.qutf8stringview_sliced_n(d, s, C.ptrdiff_t(pos), C.ptrdiff_t(n), od, os)
	})
}

func (v QUtf8StringView) First(n int) QUtf8StringView {
	if n < 0 || n > v.size {
		panic("qutf8stringview: First out of range")
	}
	return v.sliceFromC(func(d *C.char, s C.ptrdiff_t, od **C.char, os *C.ptrdiff_t) {
		C.qutf8stringview_first(d, s, C.ptrdiff_t(n), od, os)
	})
}

func (v QUtf8StringView) Last(n int) QUtf8StringView {
	if n < 0 || n > v.size {
		panic("qutf8stringview: Last out of range")
	}
	return v.sliceFromC(func(d *C.char, s C.ptrdiff_t, od **C.char, os *C.ptrdiff_t) {
		C.qutf8stringview_last(d, s, C.ptrdiff_t(n), od, os)
	})
}

func (v QUtf8StringView) Chopped(n int) QUtf8StringView {
	if n < 0 || n > v.size {
		panic("qutf8stringview: Chopped out of range")
	}
	return v.sliceFromC(func(d *C.char, s C.ptrdiff_t, od **C.char, os *C.ptrdiff_t) {
		C.qutf8stringview_chopped(d, s, C.ptrdiff_t(n), od, os)
	})
}

// Truncate shrinks the view to length n in place.
func (v *QUtf8StringView) Truncate(n int) {
	if n < 0 || n > v.size {
		panic("qutf8stringview: Truncate out of range")
	}
	v.size = n
}

// Chop removes n bytes from the end in place.
func (v *QUtf8StringView) Chop(n int) {
	if n < 0 || n > v.size {
		panic("qutf8stringview: Chop out of range")
	}
	v.size -= n
}

// Slice replaces the receiver with Sliced(pos).
func (v *QUtf8StringView) Slice(pos int) { *v = v.Sliced(pos) }

// SliceN replaces the receiver with SlicedN(pos, n).
func (v *QUtf8StringView) SliceN(pos, n int) { *v = v.SlicedN(pos, n) }

// ---- compare (cs = true → CaseSensitive, false → CaseInsensitive) --------

func (v QUtf8StringView) Compare(other QUtf8StringView, caseSensitive bool) int {
	r := C.qutf8stringview_compare_view(
		(*C.char)(v.data), C.ptrdiff_t(v.size),
		(*C.char)(other.data), C.ptrdiff_t(other.size),
		csInt(caseSensitive))
	runtime.KeepAlive(v)
	runtime.KeepAlive(other)
	return int(r)
}

func (v QUtf8StringView) CompareQByteArray(other *QByteArray, caseSensitive bool) int {
	r := C.qutf8stringview_compare_qbytearray((*C.char)(v.data), C.ptrdiff_t(v.size), other.p, csInt(caseSensitive))
	runtime.KeepAlive(v)
	runtime.KeepAlive(other)
	return int(r)
}

func (v QUtf8StringView) CompareQChar(c QChar, caseSensitive bool) int {
	r := C.qutf8stringview_compare_qchar(
		(*C.char)(v.data), C.ptrdiff_t(v.size), C.uint16_t(c.Unicode()), csInt(caseSensitive))
	runtime.KeepAlive(v)
	return int(r)
}

func (v QUtf8StringView) CompareQStringView(other QStringView, caseSensitive bool) int {
	r := C.qutf8stringview_compare_qstringview(
		(*C.char)(v.data), C.ptrdiff_t(v.size),
		(*C.uint16_t)(other.Utf16Unsafe()), C.ptrdiff_t(other.Size()),
		csInt(caseSensitive))
	runtime.KeepAlive(v)
	runtime.KeepAlive(other)
	return int(r)
}

func (v QUtf8StringView) CompareQLatin1StringView(other QLatin1StringView, caseSensitive bool) int {
	r := C.qutf8stringview_compare_qlatin1stringview(
		(*C.char)(v.data), C.ptrdiff_t(v.size),
		(*C.char)(other.DataUnsafe()), C.ptrdiff_t(other.Size()),
		csInt(caseSensitive))
	runtime.KeepAlive(v)
	runtime.KeepAlive(other)
	return int(r)
}

// ---- equal (always case-sensitive, as in Qt) ----------------------------

func (v QUtf8StringView) Equal(other QUtf8StringView) bool {
	r := C.qutf8stringview_equal_view(
		(*C.char)(v.data), C.ptrdiff_t(v.size),
		(*C.char)(other.data), C.ptrdiff_t(other.size))
	runtime.KeepAlive(v)
	runtime.KeepAlive(other)
	return bool(r)
}

func (v QUtf8StringView) EqualQByteArray(other *QByteArray) bool {
	r := C.qutf8stringview_equal_qbytearray((*C.char)(v.data), C.ptrdiff_t(v.size), other.p)
	runtime.KeepAlive(v)
	runtime.KeepAlive(other)
	return bool(r)
}

func (v QUtf8StringView) EqualQChar(c QChar) bool {
	r := C.qutf8stringview_equal_qchar(
		(*C.char)(v.data), C.ptrdiff_t(v.size), C.uint16_t(c.Unicode()))
	runtime.KeepAlive(v)
	return bool(r)
}

func (v QUtf8StringView) EqualQStringView(other QStringView) bool {
	r := C.qutf8stringview_equal_qstringview(
		(*C.char)(v.data), C.ptrdiff_t(v.size),
		(*C.uint16_t)(other.Utf16Unsafe()), C.ptrdiff_t(other.Size()))
	runtime.KeepAlive(v)
	runtime.KeepAlive(other)
	return bool(r)
}

func (v QUtf8StringView) EqualQLatin1StringView(other QLatin1StringView) bool {
	r := C.qutf8stringview_equal_qlatin1stringview(
		(*C.char)(v.data), C.ptrdiff_t(v.size),
		(*C.char)(other.DataUnsafe()), C.ptrdiff_t(other.Size()))
	runtime.KeepAlive(v)
	runtime.KeepAlive(other)
	return bool(r)
}

// ---- conversion to owning QString ---------------------------------------

// ToString returns a fresh *QString by decoding the view's bytes
// as UTF-8 (mirrors QUtf8StringView::toString()).
func (v QUtf8StringView) ToString() *QString {
	qs := wrapQString(unsafe.Pointer(C.qutf8stringview_to_qstring((*C.char)(v.data), C.ptrdiff_t(v.size))))
	runtime.KeepAlive(v)
	return qs
}

// MaxSize mirrors QUtf8StringView::maxSize().
func QUtf8StringViewMaxSize() int { return int(C.qutf8stringview_max_size()) }

// NewUnsafeWithOwner builds a view over raw UTF-8 data and attaches an owner reference
// keeping the underlying buffer alive for as long as the view (or its slices) is referenced.
//
// Advanced usage: intended for sibling go-qt6core subpackages (e.g. qanystringview).
func NewUnsafeWithOwner(data unsafe.Pointer, size int, owner any) QUtf8StringView {
	if size < 0 {
		panic("qutf8stringview: NewUnsafeWithOwner negative size")
	}
	if size > 0 && data == nil {
		panic("qutf8stringview: NewUnsafeWithOwner nil data with non-zero size")
	}
	return QUtf8StringView{data: data, size: size, keepAlive: owner}
}
