// Package quuid provides Go bindings for Qt 6.11 QtCore QUuid.
//
// QUuid is a 16-byte value type (data1 uint32, data2/data3 uint16, data4 [8]byte)
// identifying a Universally Unique Identifier. In Go, each *QUuid is a handle
// to a heap-allocated C++ QUuid; call Close() for immediate release, otherwise
// a finalizer handles cleanup.
//
// The raw byte representation is exported as Id128Bytes = [16]byte, mirroring
// QUuid::Id128Bytes. Byte order for conversions is taken from qdata.Endian
// (BigEndian = 0, LittleEndian = 1), which mirrors QSysInfo::Endian.
//
// Construction:
//   - NewQUuid()                              → QUuid()
//   - NewQUuidFromUint32(l,w1,w2,b1..b8)      → QUuid(uint, ushort, ushort, uchar×8)
//   - NewQUuidFromString(s)                   → QUuid::fromString(QAnyStringView) [fromUtf8]
//   - NewQUuidFromQString(*qtext.QString)     → QUuid::fromString(QAnyStringView)
//   - NewQUuidFromBytes(Id128Bytes, order)    → QUuid(Id128Bytes, QSysInfo::Endian)
//   - NewQUuidFromRfc4122([]byte)             → QUuid::fromRfc4122(QByteArrayView)
//   - Clone()                                 → copy constructor
//
// Conversion:
//   - ToString(format) string                 → QUuid::toString(StringFormat) → QString
//   - ToQString(format) *qtext.QString        → QUuid::toString(StringFormat)
//   - ToByteArray(format) []byte              → QUuid::toByteArray(StringFormat) → QByteArray
//   - ToQByteArray(format) *qtext.QByteArray  → QUuid::toByteArray(StringFormat)
//   - ToBytes(order) Id128Bytes               → QUuid::toBytes(QSysInfo::Endian)
//   - ToRfc4122() []byte                      → QUuid::toRfc4122() → QByteArray
//
// Properties:
//   - Variant() Variant                       → QUuid::variant()
//   - Version() Version                       → QUuid::version()
//   - IsNull() bool                           → QUuid::isNull()
//
// Comparison:
//   - Equal(other) bool                       → operator==
//   - Compare(other) int                      → operator<=> / operator< / operator>
//
// Static creators (also as package-level functions):
//   - CreateUuid()                            → QUuid::createUuid()
//   - CreateUuidV3(ns, []byte)                → QUuid::createUuidV3(QUuid, QByteArrayView)
//   - CreateUuidV3String(ns, string)          → QUuid::createUuidV3(const QUuid&, const QString&)
//   - CreateUuidV5(ns, []byte)                → QUuid::createUuidV5(QUuid, QByteArrayView)
//   - CreateUuidV5String(ns, string)          → QUuid::createUuidV5(const QUuid&, const QString&)
//   - CreateUuidV7()                          → QUuid::createUuidV7()
//
// Static parsers (also as package-level functions):
//   - FromString(s) / FromQString(qs)         → QUuid::fromString(QAnyStringView)
//   - FromBytes(Id128Bytes, order)            → QUuid::fromBytes(const void*, QSysInfo::Endian)
//   - FromRfc4122([]byte)                     → QUuid::fromRfc4122(QByteArrayView)
//
// Not implemented
//
//   - QUuid(Qt::Initialization) — internal default-construct helper.
//   - QUuid(Id128Bytes, QSysInfo::Endian) — replaced by NewQUuidFromBytes.
//   - fromBytes(const void*, QSysInfo::Endian) — replaced by FromBytes([16]byte, Endian).
//   - toUInt128 / fromUInt128 — require __int128; not portable across all targets.
//   - toCFUUID / fromCFUUID, toNSUUID / fromNSUUID — Darwin-only.
//   - GUID conversion operators — Windows-only.
//   - QDataStream operator<< / operator>> — not idiomatic in Go.
//   - QDebug operator<< — not idiomatic in Go.
//   - qHash(QUuid) — Go can hash [16]byte directly.
//   - operator< / operator> — use Compare instead.
//
// This file is part of package quuid of github.com/vadimlxx/go-qt6core.
// Requires Qt >= 6.11 and a C++17 compiler; links via pkg-config Qt6Core.
package quuid
