package quuid

/*
 Краткое описание реализации:

   QUuid в Qt — value-type из 16 байт (data1: uint32, data2/data3: uint16,
   data4[8]: uint8). В Go обёрнут через непрозрачный указатель на
   heap-аллоцированный C++ QUuid; владение — финализатор и явный Close().

   Id128Bytes — [16]byte, соответствует QUuid::Id128Bytes. Порядок байтов
   при конвертации в/из 16-байтного массива задаётся qdata.Endian
   (BigEndian=0, LittleEndian=1), совпадает с QSysInfo::Endian.

   Функции, принимающие/возвращающие QString и QByteArray, имеют две
   формы:
     - с *qtext.QString / *qtext.QByteArray (тонкая обёртка);
     - с Go string / []byte (конвертация на стороне Go).
   Так пользователь может не тянуть qtext, если ему достаточно string/[]byte,
   и наоборот — использовать уже готовые Qt-объекты.

   Передача строк в C сделана через пару (ptr, size) и unsafe.StringData —
   без C.CString / C.free, как в qstring_c.h. Это убирает зависимость от
   <stdlib.h> в cgo-преамбуле и лишний malloc.

 Что сознательно не вошло (см. quuid_doc.go):

   - Конструктор QUuid(Qt::Initialization) — служебный.
   - Конструктор QUuid(Id128Bytes, QSysInfo::Endian) — заменён
     NewQUuidFromBytes(Id128Bytes, qdata.Endian).
   - fromBytes(const void*, QSysInfo::Endian) — заменён на []byte-версию.
   - toUInt128 / fromUInt128 — зависят от __int128, не на всех платформах.
   - toCFUUID / fromCFUUID, toNSUUID / fromNSUUID — Darwin-only.
   - GUID-конвертации — Windows-only.
   - QDataStream / QDebug операторы — не идиоматичны в Go.
   - qHash(QUuid) — не критично, Go умеет хэшировать [16]byte.
   - operator<, operator> — заменены на Compare.
*/

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC -Wno-sfinae-incomplete
#cgo LDFLAGS: -lstdc++
#include <stdint.h>
#include <stdbool.h>
#include <stddef.h>
#include "quuid_c.h"
*/
import "C"

import (
	"runtime"
	"sync/atomic"
	"unsafe"

	"github.com/vadimlxx/go-qt6core/qdata"
	"github.com/vadimlxx/go-qt6core/qtext"
)

// Id128Bytes mirrors QUuid::Id128Bytes (16 bytes).
type Id128Bytes [16]byte

type quuidState struct {
	p atomic.Pointer[byte]
}

func (s *quuidState) cpp() *C.QUuid {
	if bp := s.p.Load(); bp != nil {
		return (*C.QUuid)(unsafe.Pointer(bp))
	}
	return nil
}

func (s *quuidState) claim() *C.QUuid {
	if bp := s.p.Swap(nil); bp != nil {
		return (*C.QUuid)(unsafe.Pointer(bp))
	}
	return nil
}

// QUuid wraps a Qt6 QUuid object.
type QUuid struct {
	// p *C.QUuid
	state   *quuidState
	cleanup runtime.Cleanup
}

func wrapQUuid(p unsafe.Pointer) *QUuid {
	if p == nil {
		return nil
	}

	st := &quuidState{}
	st.p.Store((*byte)(p))
	q := &QUuid{state: st}
	q.cleanup = runtime.AddCleanup(q, func(st *quuidState) {
		if cpp := st.claim(); cpp != nil {
			C.quuid_delete(cpp)
		}
	}, st)

	return q
}

func (q *QUuid) handle() *C.QUuid {
	if cpp, ok := q.tryHandle(); ok {
		return cpp
	}
	panic("quuid: use after Close (or nil receiver)")
}

func (q *QUuid) tryHandle() (*C.QUuid, bool) {
	if q == nil || q.state == nil {
		return nil, false
	}
	if cpp := q.state.cpp(); cpp != nil {
		return cpp, true
	}
	return nil, false
}

// Close releases the underlying C++ object immediately and disables the finalizer.
// Qt: QUuid::~QUuid().
func (q *QUuid) Close() {
	if q == nil || q.state == nil {
		return
	}

	q.cleanup.Stop()
	if cpp := q.state.claim(); cpp != nil {
		C.quuid_delete(cpp)
	}
}

// ---- construction --------------------------------------------------------

// NewQUuid returns a null UUID (all zero bytes).
// Qt: QUuid().
func NewQUuid(arg ...any) *QUuid {
	if len(arg) == 0 || arg[0] == nil {
		return wrapQUuid(unsafe.Pointer(C.quuid_new_empty()))
	}

	if len(arg) == 1 {
		switch s := arg[0].(type) {
		case string:
			return NewQUuidFromString(s)
		case *qtext.QString:
			return NewQUuidFromQString(s)
		case []byte:
			return NewQUuidFromRfc4122(s)
		case *QUuid:
			return s.Clone()
		}
	}

	if len(arg) == 11 {
		l,  ok1 := arg[0].(uint32)
		w1, ok2 := arg[1].(uint16)
		w2, ok3 := arg[2].(uint16)
		b1, ok4 := arg[3].(uint8)
		b2, ok5 := arg[4].(uint8)
		b3, ok6 := arg[5].(uint8)
		b4, ok7 := arg[6].(uint8)
		b5, ok8 := arg[7].(uint8)
		b6, ok9 := arg[8].(uint8)
		b7, ok10 := arg[9].(uint8)
		b8, ok11 := arg[10].(uint8)

		if ok1 && ok2 && ok3 && ok4 && ok5 && ok6 && ok7 && ok8 && ok9 && ok10 && ok11 {
			return NewQUuidFromUint32(l, w1, w2, b1, b2, b3, b4, b5, b6, b7, b8)
		}
	}

	panic("NewQUuid invalid argument")
}

// NewQUuidFromUint32 builds a QUuid from the individual fields.
// Qt: QUuid(uint l, ushort w1, ushort w2, uchar b1, ..., uchar b8).
func NewQUuidFromUint32(l uint32, w1, w2 uint16, b1, b2, b3, b4, b5, b6, b7, b8 uint8) *QUuid {
	return wrapQUuid(unsafe.Pointer(C.quuid_new_from_uint32(
		C.uint32_t(l), C.uint16_t(w1), C.uint16_t(w2),
		C.uint8_t(b1), C.uint8_t(b2), C.uint8_t(b3),
		C.uint8_t(b4), C.uint8_t(b5), C.uint8_t(b6),
		C.uint8_t(b7), C.uint8_t(b8))))
}

// NewQUuidFromString parses a UUID from a Go string (UTF-8, with explicit length).
// Qt: QUuid::fromString(QAnyStringView) — with QString::fromUtf8.
func NewQUuidFromString(s string) *QUuid {
	var p *C.char
	if len(s) > 0 {
		p = (*C.char)(unsafe.Pointer(unsafe.StringData(s)))
	}
	q := wrapQUuid(unsafe.Pointer(C.quuid_new_from_string(p, C.size_t(len(s)))))
	runtime.KeepAlive(s)
	return q
}

// NewQUuidFromQString parses a UUID from a QString.
// Qt: QUuid::fromString(QAnyStringView).
func NewQUuidFromQString(s *qtext.QString) *QUuid {
	if s == nil {
		return NewQUuid()
	}
	bs := s.Utf8()
	defer bs.Close()
	return NewQUuidFromString(bs.String())
}

// NewQUuidFromBytes builds a UUID from 16 raw bytes with the given byte order.
// Qt: QUuid(Id128Bytes, QSysInfo::Endian).
func NewQUuidFromBytes(b Id128Bytes, order qdata.Endian) *QUuid {
	return wrapQUuid(unsafe.Pointer(C.quuid_new_from_bytes((*C.uint8_t)(unsafe.Pointer(&b[0])),
		C.int(order))))
}

// NewQUuidFromRfc4122 builds a UUID from RFC 4122 big-endian bytes.
// Qt: QUuid::fromRfc4122(QByteArrayView).
func NewQUuidFromRfc4122(b []byte) *QUuid {
	var p *C.uint8_t
	if len(b) > 0 {
		p = (*C.uint8_t)(unsafe.Pointer(&b[0]))
	}
	q := wrapQUuid(unsafe.Pointer(C.quuid_from_rfc4122(p, C.size_t(len(b)))))
	runtime.KeepAlive(b)
	return q
}

// Clone returns a deep copy of the UUID.
// Qt: QUuid(const QUuid&) — copy constructor.
func (q *QUuid) Clone() *QUuid {
	return wrapQUuid(unsafe.Pointer(C.quuid_clone(q.handle())))
}

// ---- comparison ----------------------------------------------------------

// Equal reports whether q and other represent the same UUID.
// Qt: operator==(const QUuid&).
func (q *QUuid) Equal(other *QUuid) bool {
	return bool(C.quuid_equal(q.handle(), other.handle()))
}

// Compare returns -1, 0, or +1 like strcmp (Qt's strong ordering).
// Qt: operator<=>(const QUuid&) (Qt 6.8+) / operator< / operator>.
func (q *QUuid) Compare(other *QUuid) int {
	return int(C.quuid_compare(q.handle(), other.handle()))
}

// ---- conversion ----------------------------------------------------------

// ToString returns the textual representation in the given format.
// Qt: QUuid::toString(StringFormat) → QString (returned here as string).
func (q *QUuid) ToString(format StringFormat) string {
	p := C.quuid_to_string(q.handle(), C.int(format))
	if p == nil {
		return ""
	}
	defer C.quuid_free_string(p)
	return C.GoString(p)
}

// ToQString returns the textual representation as a QString.
// Qt: QUuid::toString(StringFormat).
func (q *QUuid) ToQString(format StringFormat) *qtext.QString {
	return qtext.NewQString(q.ToString(format))
}

// ToByteArray returns the textual representation as UTF-8 bytes.
// Qt: QUuid::toByteArray(StringFormat) → QByteArray (returned here as []byte).
func (q *QUuid) ToByteArray(format StringFormat) []byte {
	var n C.size_t
	p := C.quuid_to_byte_array(q.handle(), C.int(format), &n)
	if p == nil || n == 0 {
		return nil
	}
	defer C.quuid_free_string(p)
	return C.GoBytes(unsafe.Pointer(p), C.int(n))
}

// ToQByteArray returns the textual representation as a QByteArray.
// Qt: QUuid::toByteArray(StringFormat).
func (q *QUuid) ToQByteArray(format StringFormat) *qtext.QByteArray {
	qs := qtext.NewQString(q.ToString(format))
	return qs.Utf8()
}

// ToBytes returns the raw 16 bytes in the requested byte order.
// Qt: QUuid::toBytes(QSysInfo::Endian) → Id128Bytes.
func (q *QUuid) ToBytes(order qdata.Endian) Id128Bytes {
	var out Id128Bytes
	C.quuid_to_bytes(q.handle(), C.int(order), (*C.uint8_t)(unsafe.Pointer(&out[0])))
	return out
}

// ToRfc4122 returns the RFC 4122 big-endian bytes.
// Qt: QUuid::toRfc4122() → QByteArray (returned here as []byte).
func (q *QUuid) ToRfc4122() []byte {
	var n C.size_t
	p := C.quuid_to_rfc4122(q.handle(), &n)
	if p == nil || n == 0 {
		return nil
	}
	defer C.quuid_free_string(p)
	return C.GoBytes(unsafe.Pointer(p), C.int(n))
}

// ---- properties ----------------------------------------------------------

// Variant returns the UUID variant.
// Qt: QUuid::variant() const.
func (q *QUuid) Variant() Variant { return Variant(C.quuid_variant(q.handle())) }

// Version returns the UUID version (or VerUnknown).
// Qt: QUuid::version() const.
func (q *QUuid) Version() Version { return Version(C.quuid_version(q.handle())) }

// IsNull reports whether the UUID is all-zero.
// Qt: QUuid::isNull() const.
func (q *QUuid) IsNull() bool { return bool(C.quuid_is_null(q.handle())) }

// ---- static creators ----------------------------------------------------

// CreateUuid returns a new random (version 4) UUID.
// Qt: QUuid::createUuid().
func CreateUuid() *QUuid {
	return wrapQUuid(unsafe.Pointer(C.quuid_create_uuid()))
}

// CreateUuidV3 returns a deterministic name-based (MD5) UUID.
// Qt: QUuid::createUuidV3(QUuid, QByteArrayView).
func CreateUuidV3(ns *QUuid, data []byte) *QUuid {
	var p *C.uint8_t
	if len(data) > 0 {
		p = (*C.uint8_t)(unsafe.Pointer(&data[0]))
	}
	q := wrapQUuid(unsafe.Pointer(C.quuid_create_uuid_v3(ns.handle(), p, C.size_t(len(data)))))
	runtime.KeepAlive(data)
	return q
}

// CreateUuidV3String is the QString overload of CreateUuidV3.
// Qt: QUuid::createUuidV3(const QUuid&, const QString&).
func CreateUuidV3String(ns *QUuid, s string) *QUuid {
	return CreateUuidV3(ns, []byte(s))
}

// CreateUuidV5 returns a deterministic name-based (SHA-1) UUID.
// Qt: QUuid::createUuidV5(QUuid, QByteArrayView).
func CreateUuidV5(ns *QUuid, data []byte) *QUuid {
	var p *C.uint8_t
	if len(data) > 0 {
		p = (*C.uint8_t)(unsafe.Pointer(&data[0]))
	}
	q := wrapQUuid(unsafe.Pointer(C.quuid_create_uuid_v5(ns.handle(), p, C.size_t(len(data)))))
	runtime.KeepAlive(data)
	return q
}

// CreateUuidV5String is the QString overload of CreateUuidV5.
// Qt: QUuid::createUuidV5(const QUuid&, const QString&).
func CreateUuidV5String(ns *QUuid, s string) *QUuid {
	return CreateUuidV5(ns, []byte(s))
}

// CreateUuidV7 returns a new time-ordered (version 7) UUID.
// Qt: QUuid::createUuidV7().
func CreateUuidV7() *QUuid {
	return wrapQUuid(unsafe.Pointer(C.quuid_create_uuid_v7()))
}

// ---- static parsers -----------------------------------------------------

// FromString parses a UUID from a Go string (UTF-8).
// Qt: QUuid::fromString(QAnyStringView) — with QString::fromUtf8.
func FromString(s string) *QUuid { return NewQUuidFromString(s) }

// FromQString parses a UUID from a QString.
// Qt: QUuid::fromString(QAnyStringView).
func FromQString(s *qtext.QString) *QUuid { return NewQUuidFromQString(s) }

// FromBytes builds a UUID from 16 bytes in the given order.
// Qt: QUuid::fromBytes(const void*, QSysInfo::Endian).
func FromBytes(b Id128Bytes, order qdata.Endian) *QUuid {
	return NewQUuidFromBytes(b, order)
}

// FromRfc4122 builds a UUID from RFC 4122 big-endian bytes.
// Qt: QUuid::fromRfc4122(QByteArrayView).
func FromRfc4122(b []byte) *QUuid { return NewQUuidFromRfc4122(b) }
