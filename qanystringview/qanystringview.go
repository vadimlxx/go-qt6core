package qanystringview

/*
 что не реализовано:
 visit() — паттерн visitor из C++ не переносится в Go; используйте Kind() + As*().
 arg() — вариативный шаблон QString::arg; появится с QString::arg.
 Одиночные символьные конструкторы QAnyStringView(Char) — используйте FromString(string(c)).
 fromArray() — используйте FromBytes / FromString.
 Неявные конструкторы с capacity-hint из QString/QByteArray — используйте FromQString / FromQByteArray.
 Операторы сравнения (==, <) — используйте Compare / Equal.
 Неявные конверсии в QString / QByteArray — используйте ToString / As*.
 operator std::string_view — используйте Bytes().
 Итераторы (begin/end/rbegin/rend) — используйте конкретные view-типы.
 STL container compatibility — не нужно в Go.
 detects_US_ASCII_at_compile_time — неактуально для Go.
 */

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC -Wno-sfinae-incomplete
#cgo LDFLAGS: -lstdc++
#include <stdint.h>
#include <stdbool.h>
#include <stddef.h>
#include "qanystringview_c.h"
*/
import "C"

import (
	"runtime"
	"unsafe"

	"github.com/vadimlxx/go-qt6core/qbytearray"
	"github.com/vadimlxx/go-qt6core/qchar"
	"github.com/vadimlxx/go-qt6core/qlatin1stringview"
	"github.com/vadimlxx/go-qt6core/qstring"
	"github.com/vadimlxx/go-qt6core/qstringview"
	"github.com/vadimlxx/go-qt6core/qutf8stringview"
)

// Kind identifies the encoding stored in a QAnyStringView.
// Values match Qt's internal tag: 0=Utf8, 1=Latin1, 2=Utf16.
type Kind int

const (
	KindUtf8   Kind = 0
	KindLatin1 Kind = 1
	KindUtf16  Kind = 2
)

func (k Kind) String() string {
	switch k {
	case KindUtf8:
		return "Utf8"
	case KindLatin1:
		return "Latin1"
	case KindUtf16:
		return "Utf16"
	}
	return "Kind(?)"
}

// QAnyStringView is a type-erased, read-only string view over UTF-8, Latin-1, or UTF-16 data.
// Mirrors Qt6 QAnyStringView (Q_PRIMITIVE_TYPE).
//
// The zero value is a valid null view (KindUtf8, size=0, data=nil). Prefer Empty() over
// a bare QAnyStringView{} literal to route through the C++ default constructor.
type QAnyStringView struct {
	data      unsafe.Pointer
	size      int
	kind      Kind
	keepAlive any
}

// ---- construction --------------------------------------------------------

// Empty returns a null QAnyStringView by invoking the C++ default constructor
// QAnyStringView(). This stays correct even if Qt changes the internal representation
// of an empty view.
func Empty() QAnyStringView {
	var data unsafe.Pointer
	var size C.ptrdiff_t
	var tag C.int
	C.qanystringview_default(&data, &size, &tag)
	return QAnyStringView{data: data, size: int(size), kind: Kind(tag)}
}

// FromBytes returns a view over b as UTF-8 (KindUtf8). The Go slice is kept alive.
func FromBytes(b []byte) QAnyStringView {
	if len(b) == 0 {
		return Empty()
	}
	return QAnyStringView{data: unsafe.Pointer(&b[0]), size: len(b), kind: KindUtf8, keepAlive: b}
}

// FromString returns a view over s as UTF-8 (KindUtf8). Go strings are immutable.
func FromString(s string) QAnyStringView {
	if len(s) == 0 {
		return Empty()
	}
	return QAnyStringView{data: unsafe.Pointer(unsafe.StringData(s)), size: len(s), kind: KindUtf8, keepAlive: s}
}

// FromQByteArray returns a zero-copy UTF-8 view over q's internal buffer. The view keeps q alive.
func FromQByteArray(q *qbytearray.QByteArray) QAnyStringView {
	if q == nil || q.IsNull() {
		return Empty()
	}
	return QAnyStringView{data: q.ConstDataUnsafe(), size: q.Size(), kind: KindUtf8, keepAlive: q}
}

// FromQUtf8StringView wraps a QUtf8StringView (KindUtf8). The original view is kept alive.
func FromQUtf8StringView(v qutf8stringview.QUtf8StringView) QAnyStringView {
	if v.IsNull() {
		return Empty()
	}
	return QAnyStringView{data: v.DataUnsafe(), size: v.Size(), kind: KindUtf8, keepAlive: v}
}

// FromQLatin1StringView wraps a QLatin1StringView (KindLatin1).
func FromQLatin1StringView(v qlatin1stringview.QLatin1StringView) QAnyStringView {
	if v.IsNull() {
		return Empty()
	}
	return QAnyStringView{data: v.DataUnsafe(), size: v.Size(), kind: KindLatin1, keepAlive: v}
}

// FromQStringView wraps a QStringView (KindUtf16).
func FromQStringView(v qstringview.QStringView) QAnyStringView {
	if v.IsNull() {
		return Empty()
	}
	return QAnyStringView{data: v.Utf16Unsafe(), size: v.Size(), kind: KindUtf16, keepAlive: v}
}

// NewUnsafe builds a QAnyStringView from raw data + size + kind. Advanced.
//
// Caller is solely responsible for keeping data alive and unmutated.
// Prefer the typed From* constructors; use NewUnsafe only when you already have
// a raw pointer whose lifetime you control.
func NewUnsafe(data unsafe.Pointer, size int, kind Kind) QAnyStringView {
	if size < 0 {
		panic("qanystringview: NewUnsafe negative size")
	}
	if size > 0 && data == nil {
		panic("qanystringview: NewUnsafe nil data with non-zero size")
	}
	if kind != KindUtf8 && kind != KindLatin1 && kind != KindUtf16 {
		panic("qanystringview: NewUnsafe invalid kind")
	}
	return QAnyStringView{data: data, size: size, kind: kind}
}

// ---- introspection -------------------------------------------------------

func (v QAnyStringView) Kind() Kind    { return v.kind }
func (v QAnyStringView) Size() int     { return v.size }
func (v QAnyStringView) Len() int      { return v.size }
func (v QAnyStringView) IsEmpty() bool { return v.size == 0 }
func (v QAnyStringView) IsNull() bool  { return v.data == nil }

// SizeBytes returns the byte size (size × 2 for Utf16, size for Utf8/Latin1).
func (v QAnyStringView) SizeBytes() int {
	if v.kind == KindUtf16 {
		return v.size * 2
	}
	return v.size
}

func (v QAnyStringView) IsUtf8() bool   { return v.kind == KindUtf8 }
func (v QAnyStringView) IsLatin1() bool { return v.kind == KindLatin1 }
func (v QAnyStringView) IsUtf16() bool  { return v.kind == KindUtf16 }

// DataUnsafe returns the raw data pointer (typed per Kind).
// Valid only while the view (and its keepAlive target) is alive.
func (v QAnyStringView) DataUnsafe() unsafe.Pointer { return v.data }

// Bytes returns a copy of the raw bytes (for Utf8/Latin1: the byte string; for Utf16:
// the underlying UTF-16 bytes in native endianness). Use As*() + Go conversion for typed access.
func (v QAnyStringView) Bytes() []byte {
	if v.size == 0 {
		return nil
	}
	n := v.SizeBytes()
	out := make([]byte, n)
	copy(out, unsafe.Slice((*byte)(v.data), n))
	runtime.KeepAlive(v)
	return out
}

// ---- element access ------------------------------------------------------

// Front returns the first character as a QChar. Panics if empty.
func (v QAnyStringView) Front() qchar.QChar {
	if v.size == 0 {
		panic("qanystringview: Front on empty view")
	}
	r := C.qanystringview_front(v.data, C.ptrdiff_t(v.size), C.int(v.kind))
	runtime.KeepAlive(v)
	return qchar.New(uint16(r))
}

// Back returns the last character as a QChar. Panics if empty.
func (v QAnyStringView) Back() qchar.QChar {
	if v.size == 0 {
		panic("qanystringview: Back on empty view")
	}
	r := C.qanystringview_back(v.data, C.ptrdiff_t(v.size), C.int(v.kind))
	runtime.KeepAlive(v)
	return qchar.New(uint16(r))
}

// ---- slicing -------------------------------------------------------------

// Left clamps n (Qt semantics; n < 0 or n > size → whole view).
func (v QAnyStringView) Left(n int) QAnyStringView {
	var od unsafe.Pointer
	var os C.ptrdiff_t
	var ot C.int
	C.qanystringview_left(v.data, C.ptrdiff_t(v.size), C.int(v.kind), C.ptrdiff_t(n), &od, &os, &ot)
	runtime.KeepAlive(v)
	return QAnyStringView{data: od, size: int(os), kind: Kind(ot), keepAlive: v.keepAlive}
}

// Right clamps n (Qt semantics).
func (v QAnyStringView) Right(n int) QAnyStringView {
	var od unsafe.Pointer
	var os C.ptrdiff_t
	var ot C.int
	C.qanystringview_right(v.data, C.ptrdiff_t(v.size), C.int(v.kind), C.ptrdiff_t(n), &od, &os, &ot)
	runtime.KeepAlive(v)
	return QAnyStringView{data: od, size: int(os), kind: Kind(ot), keepAlive: v.keepAlive}
}

// Mid clamps pos/n (Qt semantics; n = -1 → to the end). May return a null view.
func (v QAnyStringView) Mid(pos, n int) QAnyStringView {
	var od unsafe.Pointer
	var os C.ptrdiff_t
	var ot C.int
	C.qanystringview_mid(v.data, C.ptrdiff_t(v.size), C.int(v.kind),
		C.ptrdiff_t(pos), C.ptrdiff_t(n), &od, &os, &ot)
	runtime.KeepAlive(v)
	return QAnyStringView{data: od, size: int(os), kind: Kind(ot), keepAlive: v.keepAlive}
}

// First returns the first n characters. Panics on out-of-range (Qt release builds omit the check).
func (v QAnyStringView) First(n int) QAnyStringView {
	if n < 0 || n > v.size {
		panic("qanystringview: First out of range")
	}
	var od unsafe.Pointer
	var os C.ptrdiff_t
	var ot C.int
	C.qanystringview_first(v.data, C.ptrdiff_t(v.size), C.int(v.kind), C.ptrdiff_t(n), &od, &os, &ot)
	runtime.KeepAlive(v)
	return QAnyStringView{data: od, size: int(os), kind: Kind(ot), keepAlive: v.keepAlive}
}

// Last returns the last n characters. Panics on out-of-range.
func (v QAnyStringView) Last(n int) QAnyStringView {
	if n < 0 || n > v.size {
		panic("qanystringview: Last out of range")
	}
	var od unsafe.Pointer
	var os C.ptrdiff_t
	var ot C.int
	C.qanystringview_last(v.data, C.ptrdiff_t(v.size), C.int(v.kind), C.ptrdiff_t(n), &od, &os, &ot)
	runtime.KeepAlive(v)
	return QAnyStringView{data: od, size: int(os), kind: Kind(ot), keepAlive: v.keepAlive}
}

// Sliced returns the suffix starting at pos. Panics on out-of-range.
func (v QAnyStringView) Sliced(pos int) QAnyStringView {
	if pos < 0 || pos > v.size {
		panic("qanystringview: Sliced out of range")
	}
	var od unsafe.Pointer
	var os C.ptrdiff_t
	var ot C.int
	C.qanystringview_sliced(v.data, C.ptrdiff_t(v.size), C.int(v.kind), C.ptrdiff_t(pos), &od, &os, &ot)
	runtime.KeepAlive(v)
	return QAnyStringView{data: od, size: int(os), kind: Kind(ot), keepAlive: v.keepAlive}
}

// SlicedN returns the substring [pos, pos+n). Panics on out-of-range.
func (v QAnyStringView) SlicedN(pos, n int) QAnyStringView {
	if pos < 0 || n < 0 || pos+n > v.size {
		panic("qanystringview: SlicedN out of range")
	}
	var od unsafe.Pointer
	var os C.ptrdiff_t
	var ot C.int
	C.qanystringview_sliced_n(v.data, C.ptrdiff_t(v.size), C.int(v.kind),
		C.ptrdiff_t(pos), C.ptrdiff_t(n), &od, &os, &ot)
	runtime.KeepAlive(v)
	return QAnyStringView{data: od, size: int(os), kind: Kind(ot), keepAlive: v.keepAlive}
}

// Chopped removes the last n characters. Panics on out-of-range.
func (v QAnyStringView) Chopped(n int) QAnyStringView {
	if n < 0 || n > v.size {
		panic("qanystringview: Chopped out of range")
	}
	var od unsafe.Pointer
	var os C.ptrdiff_t
	var ot C.int
	C.qanystringview_chopped(v.data, C.ptrdiff_t(v.size), C.int(v.kind), C.ptrdiff_t(n), &od, &os, &ot)
	runtime.KeepAlive(v)
	return QAnyStringView{data: od, size: int(os), kind: Kind(ot), keepAlive: v.keepAlive}
}

// Truncate shrinks the view to length n in place.
func (v *QAnyStringView) Truncate(n int) {
	if n < 0 || n > v.size {
		panic("qanystringview: Truncate out of range")
	}
	v.size = n
}

// Chop removes n code units/bytes from the end in place.
func (v *QAnyStringView) Chop(n int) {
	if n < 0 || n > v.size {
		panic("qanystringview: Chop out of range")
	}
	v.size -= n
}

// Slice replaces the receiver with Sliced(pos).
func (v *QAnyStringView) Slice(pos int) { *v = v.Sliced(pos) }

// SliceN replaces the receiver with SlicedN(pos, n).
func (v *QAnyStringView) SliceN(pos, n int) { *v = v.SlicedN(pos, n) }

// ---- dispatch to concrete type ------------------------------------------

// AsQStringView returns the underlying QStringView and true if Kind is Utf16,
// or (qstringview.Empty(), false) otherwise.
func (v QAnyStringView) AsQStringView() (qstringview.QStringView, bool) {
	if v.kind != KindUtf16 {
		return qstringview.Empty(), false
	}
	if v.data == nil {
		return qstringview.Empty(), true
	}
	return qstringview.NewUnsafeWithOwner(v.data, v.size, v.keepAlive), true
}

// AsQLatin1StringView returns the underlying QLatin1StringView and true if Kind is Latin1.
func (v QAnyStringView) AsQLatin1StringView() (qlatin1stringview.QLatin1StringView, bool) {
	if v.kind != KindLatin1 {
		return qlatin1stringview.Empty(), false
	}
	if v.data == nil {
		return qlatin1stringview.Empty(), true
	}
	return qlatin1stringview.NewUnsafeWithOwner(v.data, v.size, v.keepAlive), true
}

// AsQUtf8StringView returns the underlying QUtf8StringView and true if Kind is Utf8.
func (v QAnyStringView) AsQUtf8StringView() (qutf8stringview.QUtf8StringView, bool) {
	if v.kind != KindUtf8 {
		return qutf8stringview.Empty(), false
	}
	if v.data == nil {
		return qutf8stringview.Empty(), true
	}
	return qutf8stringview.NewUnsafeWithOwner(v.data, v.size, v.keepAlive), true
}

// ---- comparison ----------------------------------------------------------

// Compare returns Qt's int comparison result (negative/zero/positive).
// Equivalent to the static QAnyStringView::compare(lhs, rhs, cs).
func Compare(lhs, rhs QAnyStringView, caseSensitive bool) int {
	csInt := C.int(0)
	if !caseSensitive {
		csInt = 1
	}
	r := C.qanystringview_compare(lhs.data, C.ptrdiff_t(lhs.size), C.int(lhs.kind),
		rhs.data, C.ptrdiff_t(rhs.size), C.int(rhs.kind), csInt)
	runtime.KeepAlive(lhs)
	runtime.KeepAlive(rhs)
	return int(r)
}

// Equal returns true if lhs == rhs (case-sensitive).
// Equivalent to the static QAnyStringView::equal(lhs, rhs).
func Equal(lhs, rhs QAnyStringView) bool {
	r := C.qanystringview_equal(lhs.data, C.ptrdiff_t(lhs.size), C.int(lhs.kind),
		rhs.data, C.ptrdiff_t(rhs.size), C.int(rhs.kind))
	runtime.KeepAlive(lhs)
	runtime.KeepAlive(rhs)
	return bool(r)
}

// ---- conversion ----------------------------------------------------------

// ToString returns a fresh *qstring.QString by decoding the contents according to Kind.
// Equivalent to QAnyStringView::toString() (dispatches to QString::fromUtf8/fromLatin1/fromUtf16).
func (v QAnyStringView) ToString() *qstring.QString {
	var out unsafe.Pointer
	C.qanystringview_to_qstring(v.data, C.ptrdiff_t(v.size), C.int(v.kind), &out)
	runtime.KeepAlive(v)
	return qstring.AdoptUnsafePtr(out)
}
