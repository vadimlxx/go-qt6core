package qstringbuilder

// Package qstringbuilder предоставляет Go-биндинг для Qt 6.11 QtCore/QStringBuilder.
//
// TODO (на будущее, не реализовано):
//   - operator% и operator+ (QT_USE_QSTRINGBUILDER) — в Go нет перегрузки операторов;
//     используйте Concat / AppendPart / Append.
//   - QConcatenable<T>-специализации — внутренний механизм C++, не экспортируется.
//   - ExactSize-оптимизация — внутренняя деталь C++; в Go материализация всегда даёт точный
//     результат, resize делается в C++.
//   - QStringBuilder::size() и ::isNull() — не экспортируются (в Qt семантика "приблизительная",
//     верхняя граница); используйте ToQString().Size() / ToQByteArray().Size().
//   - toUpper()/toLower() как методы builder'а — тип результата зависит от flavor;
//     используйте ToQString().ToUpper() / ToLower() или ToQByteArray().ToUpper() / ToLower().
//   - operator std::string_view / std::u16string_view — неприменимо в Go.
//   - QDebug operator<< — неприменимо.
//   - QUtf8StringView / QAnyStringView как части — в Qt 6.11 не входят в QStringBuilder;
//     при передаче таких объектов — panic.
//   - QStringBuilder<A, B> с более чем двумя операндами — покрывается Concat(...) или цепочкой Append.

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC -Wno-sfinae-incomplete
#cgo LDFLAGS: -lstdc++
#include <stdint.h>
#include <stddef.h>
#include "qstringbuilder_c.h"
*/
import "C"

import (
	"fmt"
	"runtime"
	"unsafe"

	"github.com/vadimlxx/go-qt6core/qbytearray"
	"github.com/vadimlxx/go-qt6core/qbytearrayview"
	"github.com/vadimlxx/go-qt6core/qchar"
	"github.com/vadimlxx/go-qt6core/qlatin1char"
	"github.com/vadimlxx/go-qt6core/qlatin1stringview"
	"github.com/vadimlxx/go-qt6core/qstring"
	"github.com/vadimlxx/go-qt6core/qstringview"
)

// part is an internal representation of one operand in a QStringBuilder.
type part struct {
	kind  C.int           // QSB_KIND_* constant
	tag   C.int           // 16-bit / 8-bit value for single-character parts, else 0
	data  unsafe.Pointer  // borrowed pointer into keep-alive source
	size  C.ptrdiff_t     // element count (bytes or UTF-16 code units)
	keep  any             // source Go object (kept alive during C call)
}

// QStringBuilder is a Go equivalent of Qt's QStringBuilder<A, B> template.
//
// It accumulates a sequence of Parts and materializes them into *qstring.QString or
// *qbytearray.QByteArray on demand. See the package doc for the full semantics.
//
// A QStringBuilder is not safe for concurrent mutation. Materialization (ToQString /
// ToQByteArray / AppendToQString / AppendToQByteArray) is read-only and may be called
// concurrently as long as no concurrent Append is in progress.
type QStringBuilder struct {
	parts      []part
	qsFlavored bool
}

// ---- construction --------------------------------------------------------

// Concat creates a QStringBuilder from the given operands. Operands may be any of the
// supported types (see package doc); passing an unsupported type panics.
//
// Equivalent to Qt's `a % b % c % ...` chain.
func Concat(operands ...any) *QStringBuilder {
	b := &QStringBuilder{}
	for _, op := range operands {
		b.addOperand(op)
	}
	return b
}

// FromString returns a QStringBuilder holding a UTF-8 string (QByteArray-flavored).
// Equivalent to Qt's `QByteArray(s)` as a QStringBuilder operand.
func FromString(s string) *QStringBuilder {
	b := &QStringBuilder{}
	b.addOperand(s)
	return b
}

// FromBytes returns a QStringBuilder holding raw bytes (QByteArray-flavored).
func FromBytes(bs []byte) *QStringBuilder {
	b := &QStringBuilder{}
	b.addOperand(bs)
	return b
}

// FromByte returns a QStringBuilder holding a single Latin-1 byte (QByteArray-flavored,
// mirrors Qt's `char` operand).
func FromByte(c byte) *QStringBuilder {
	b := &QStringBuilder{}
	b.addOperand(c)
	return b
}

// FromQChar returns a QStringBuilder holding one UTF-16 code unit (QString-flavored,
// mirrors Qt's `QChar` operand).
func FromQChar(c qchar.QChar) *QStringBuilder {
	b := &QStringBuilder{}
	b.addOperand(c)
	return b
}

// FromQLatin1Char returns a QStringBuilder holding one Latin-1 character (QString-flavored,
// mirrors Qt's `QLatin1Char` operand).
func FromQLatin1Char(c qlatin1char.QLatin1Char) *QStringBuilder {
	b := &QStringBuilder{}
	b.addOperand(c)
	return b
}

// FromUint16 returns a QStringBuilder holding one char16_t code unit (QString-flavored,
// mirrors Qt's `char16_t` operand).
func FromUint16(c uint16) *QStringBuilder {
	b := &QStringBuilder{}
	b.addOperand(c)
	return b
}

// FromUtf16 returns a QStringBuilder holding UTF-16 code units (QString-flavored,
// mirrors Qt's `const char16_t*` operand).
func FromUtf16(units []uint16) *QStringBuilder {
	b := &QStringBuilder{}
	b.addOperand(units)
	return b
}

// FromQString wraps an existing *qstring.QString (QString-flavored).
// The view keeps q alive; do not Close q while the builder (or any materialization) is in use.
func FromQString(q *qstring.QString) *QStringBuilder {
	b := &QStringBuilder{}
	b.addOperand(q)
	return b
}

// FromQStringView wraps a QStringView (QString-flavored).
func FromQStringView(v qstringview.QStringView) *QStringBuilder {
	b := &QStringBuilder{}
	b.addOperand(v)
	return b
}

// FromQLatin1StringView wraps a QLatin1StringView (QString-flavored).
func FromQLatin1StringView(v qlatin1stringview.QLatin1StringView) *QStringBuilder {
	b := &QStringBuilder{}
	b.addOperand(v)
	return b
}

// FromQByteArray wraps an existing *qbytearray.QByteArray (QByteArray-flavored).
// The view keeps ba alive; do not Close ba while the builder (or any materialization) is in use.
func FromQByteArray(ba *qbytearray.QByteArray) *QStringBuilder {
	b := &QStringBuilder{}
	b.addOperand(ba)
	return b
}

// FromQByteArrayView wraps a QByteArrayView (QByteArray-flavored).
func FromQByteArrayView(v qbytearrayview.QByteArrayView) *QStringBuilder {
	b := &QStringBuilder{}
	b.addOperand(v)
	return b
}

// addOperand is the single dispatcher for every supported type.
func (b *QStringBuilder) addOperand(op any) {
	switch v := op.(type) {
	case string:
		if len(v) > 0 {
			b.parts = append(b.parts, part{
				kind: C.QSB_KIND_UTF8_STRING,
				data: unsafe.Pointer(unsafe.StringData(v)),
				size: C.ptrdiff_t(len(v)),
				keep: v,
			})
		}
	case []byte:
		if len(v) > 0 {
			b.parts = append(b.parts, part{
				kind: C.QSB_KIND_UTF8_STRING,
				data: unsafe.Pointer(&v[0]),
				size: C.ptrdiff_t(len(v)),
				keep: v,
			})
		}
	case byte:
		b.parts = append(b.parts, part{
			kind: C.QSB_KIND_QLATIN1CHAR,
			tag:  C.int(v),
		})
	case uint16:
		b.parts = append(b.parts, part{
			kind: C.QSB_KIND_QCHAR,
			tag:  C.int(v),
		})
		b.qsFlavored = true
	case []uint16:
		if len(v) > 0 {
			b.parts = append(b.parts, part{
				kind: C.QSB_KIND_QSTRINGVIEW,
				data: unsafe.Pointer(&v[0]),
				size: C.ptrdiff_t(len(v)),
				keep: v,
			})
			b.qsFlavored = true
		}
	case qchar.QChar:
		b.parts = append(b.parts, part{
			kind: C.QSB_KIND_QCHAR,
			tag:  C.int(uint16(v)),
		})
		b.qsFlavored = true
	case qlatin1char.QLatin1Char:
		b.parts = append(b.parts, part{
			kind: C.QSB_KIND_QLATIN1CHAR,
			tag:  C.int(byte(v)),
		})
		b.qsFlavored = true
	case qlatin1stringview.QLatin1StringView:
		if !v.IsNull() && v.Size() > 0 {
			b.parts = append(b.parts, part{
				kind: C.QSB_KIND_QLATIN1STRINGVIEW,
				data: v.DataUnsafe(),
				size: C.ptrdiff_t(v.Size()),
				keep: v,
			})
			b.qsFlavored = true
		}
	case qstringview.QStringView:
		if !v.IsNull() && v.Size() > 0 {
			b.parts = append(b.parts, part{
				kind: C.QSB_KIND_QSTRINGVIEW,
				data: v.Utf16Unsafe(),
				size: C.ptrdiff_t(v.Size()),
				keep: v,
			})
			b.qsFlavored = true
		}
	case *qstring.QString:
		if v != nil && !v.IsNull() {
			b.parts = append(b.parts, part{
				kind: C.QSB_KIND_QSTRING,
				data: v.UnsafePtr(),
				keep: v,
			})
			b.qsFlavored = true
		}
	case *qbytearray.QByteArray:
		if v != nil && !v.IsNull() {
			b.parts = append(b.parts, part{
				kind: C.QSB_KIND_QBYTEARRAY,
				data: v.UnsafePtr(),
				keep: v,
			})
		}
	case qbytearrayview.QByteArrayView:
		if !v.IsNull() && v.Size() > 0 {
			b.parts = append(b.parts, part{
				kind: C.QSB_KIND_QBYTEARRAYVIEW,
				data: v.DataUnsafe(),
				size: C.ptrdiff_t(v.Size()),
				keep: v,
			})
		}
	case *QStringBuilder:
		if v != nil {
			b.parts = append(b.parts, v.parts...)
			if v.qsFlavored {
				b.qsFlavored = true
			}
		}
	case nil:
		// ignore nil operands
	default:
		panic(fmt.Sprintf("qstringbuilder: unsupported operand type %T", op))
	}
}

// ---- queries -------------------------------------------------------------

// Len returns the number of accumulated parts (empty parts are not stored).
func (b *QStringBuilder) Len() int { return len(b.parts) }

// IsQStringFlavored reports whether materialization produces a QString (true) or QByteArray
// (false). Equivalent to Qt's ConvertTo selection: QString wins if any part is QString-flavored.
func (b *QStringBuilder) IsQStringFlavored() bool { return b.qsFlavored }

// ---- concatenation -------------------------------------------------------

// Append concatenates another QStringBuilder's parts into b (in place) and returns b for chaining.
func (b *QStringBuilder) Append(other *QStringBuilder) *QStringBuilder {
	if other != nil {
		b.parts = append(b.parts, other.parts...)
		if other.qsFlavored {
			b.qsFlavored = true
		}
	}
	return b
}

// AppendPart is a convenience wrapper around addOperand for chaining.
func (b *QStringBuilder) AppendPart(op any) *QStringBuilder {
	b.addOperand(op)
	return b
}

// ---- materialization -----------------------------------------------------

// buildCArray materializes b.parts into a freshly allocated C array.
// The returned slice and the parts slice must remain alive during the C call.
func (b *QStringBuilder) buildCArray() (*C.qsb_part_t, []any) {
	n := len(b.parts)
	if n == 0 {
		return nil, nil
	}
	arr := make([]C.qsb_part_t, n)
	keep := make([]any, 0, n)
	for i, p := range b.parts {
		arr[i].kind = p.kind
		arr[i].tag = p.tag
		// p.data is unsafe.Pointer (Go pointer). Convert to raw integer so the C
		// struct contains no Go pointers — this avoids cgo's
		// "Go pointer to Go pointer" check. The pointed-to memory stays alive
		// because p.data is still an unsafe.Pointer in b.parts (GC-visible),
		// and we additionally KeepAlive(b) after the C call.
		arr[i].data = C.uintptr_t(uintptr(p.data))
		arr[i].size = p.size
		if p.keep != nil {
			keep = append(keep, p.keep)
		}
	}
	return &arr[0], keep
}

// AppendToQString appends b's parts to q (in place). Equivalent to Qt's `q += builder`.
// The builder is not consumed; q must not be nil.
func (b *QStringBuilder) AppendToQString(q *qstring.QString) {
	if q == nil {
		panic("qstringbuilder: AppendToQString with nil target")
	}
	n := len(b.parts)
	if n == 0 {
		return
	}
	arr, keep := b.buildCArray()
	C.qstringbuilder_append_to_qstring(q.UnsafePtr(), arr, C.size_t(n))
	runtime.KeepAlive(q)
	runtime.KeepAlive(b)
	for _, k := range keep {
		runtime.KeepAlive(k)
	}
}

// AppendToQByteArray appends b's parts to ba (in place). Equivalent to Qt's `ba += builder`.
// The builder is not consumed; ba must not be nil.
func (b *QStringBuilder) AppendToQByteArray(ba *qbytearray.QByteArray) {
	if ba == nil {
		panic("qstringbuilder: AppendToQByteArray with nil target")
	}
	n := len(b.parts)
	if n == 0 {
		return
	}
	arr, keep := b.buildCArray()
	C.qstringbuilder_append_to_qbytearray(ba.UnsafePtr(), arr, C.size_t(n))
	runtime.KeepAlive(ba)
	runtime.KeepAlive(b)
	for _, k := range keep {
		runtime.KeepAlive(k)
	}
}

// ToQString materializes b as a fresh *qstring.QString.
// Equivalent to Qt's implicit conversion to QString (works for any flavor).
func (b *QStringBuilder) ToQString() *qstring.QString {
	q := qstring.New()
	b.AppendToQString(q)
	return q
}

// ToQByteArray materializes b as a fresh *qbytearray.QByteArray (UTF-8 encoded).
// For QByteArray-flavored builders this is exact; for QString-flavored builders the contents
// are UTF-8 encoded (matches Qt's `builder.toUtf8()`).
func (b *QStringBuilder) ToQByteArray() *qbytearray.QByteArray {
	ba := qbytearray.New()
	b.AppendToQByteArray(ba)
	return ba
}

// ToUtf8 is equivalent to ToQString().Utf8().
// Only meaningful for QString-flavored builders; on QByteArray-flavored builders it returns
// the UTF-8 encoding of the (already byte-oriented) contents.
func (b *QStringBuilder) ToUtf8() *qbytearray.QByteArray { return b.ToQString().Utf8() }

// ToLatin1 is equivalent to ToQString().Latin1().
func (b *QStringBuilder) ToLatin1() *qbytearray.QByteArray { return b.ToQString().Latin1() }

// ToLocal8Bit is equivalent to ToQString().Local8Bit().
func (b *QStringBuilder) ToLocal8Bit() *qbytearray.QByteArray { return b.ToQString().Local8Bit() }
