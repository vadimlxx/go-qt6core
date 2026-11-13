// Package qdatetime provides Go bindings for Qt 6.11 QtCore date/time-related types.
//
// QTimeZone uses COW (Q_DECLARE_SHARED) and is NOT trivially copyable. In Go it is a
// HANDLE: a struct holding a pointer to a heap-allocated C++ QTimeZone. All accessors go
// directly to the C++ object. Use Clone() for an independent copy, Close() to free
// immediately; otherwise a finalizer handles it.
//
// The Go zero value QTimeZone{} is invalid (nil pointer). Use NewQTimeZone() for a null
// QTimeZone that mirrors Qt's QTimeZone().
//
// # Internal representation (ShortData)
//
// QTimeZone internally uses a union of ShortData (for Qt::UTC, Qt::LocalTime and
// Qt::OffsetFromUTC) and a pointer to QTimeZonePrivate (for full IANA zones). This is
// handled transparently by Qt; the Go binding never exposes ShortData. Use TimeSpec() to
// query which mode is active, and IsUTCOrFixedOffset() as a convenience check.
//
// # Construction
//
//	NewQTimeZone()                    — null QTimeZone (Qt QTimeZone())
//	NewQTimeZoneLocal()               — QTimeZone(QTimeZone::LocalTime)
//	NewQTimeZoneUTC()                 — QTimeZone(QTimeZone::UTC)
//	NewQTimeZoneFromOffset(seconds)   — QTimeZone(int offsetSeconds)
//	NewQTimeZoneFromIANA(id)          — QTimeZone(QByteArray ianaId)
//	SystemTimeZone()                  — QTimeZone::systemTimeZone()
//	QTimeZoneUTC()                    — QTimeZone::utc()
//
// # Accessors
//
//	IsValid() bool
//	Id() *qtext.QByteArray
//	TimeSpec() TimeSpec
//	FixedSecondsAheadOfUtc() int
//	IsUTCOrFixedOffset() bool
//	Territory() Country
//
// # Offset queries (take *QDateTime)
//
//	OffsetFromUtc(at) int
//	StandardTimeOffset(at) int
//	DaylightTimeOffset(at) int
//	IsDaylightTime(at) bool
//	HasDaylightTime() bool
//	HasTransitions() bool
//
// # Display names (return fresh *qtext.QString)
//
//	DisplayNameAt(at, nameType, locale) *qtext.QString
//	DisplayNameType(timeType, nameType, locale) *qtext.QString
//	Abbreviation(at) *qtext.QString
//
// # Transitions
//
//	NextTransition(after) (offset, stdOffset, dstOffset int64, atUtc *QDateTime, abbr *qtext.QString)
//	PreviousTransition(before) (offset, stdOffset, dstOffset int64, atUtc *QDateTime, abbr *qtext.QString)
//
// # Conversions
//
//	AsBackendZone() *QTimeZone
//
// # Static utilities
//
//	SystemTimeZoneId() *qtext.QByteArray   — QTimeZone::systemTimeZoneId()
//	IsTimeZoneIdAvailable(id) bool         — QTimeZone::isTimeZoneIdAvailable()
//
// # Enum types in this package
//
//	time_type.go  — TimeType (QTimeZone::TimeType: StandardTime, DaylightTime, GenericTime)
//	name_type.go  — NameType (QTimeZone::NameType: DefaultName, LongName, ShortName, OffsetName)
//
// # Not yet implemented
//
//   - QTimeZone::transitions(from, to) OffsetDataList — requires QList<OffsetData>.
//   - QTimeZone::availableTimeZoneIds() — requires QList<QByteArray>.
//   - QTimeZone::availableTimeZoneIds(territory) / (offsetSeconds) — same reason.
//   - QTimeZone::ianaIdToWindowsId(), windowsIdToIanaIds(), windowsIdToDefaultIanaId() —
//     require QList<QByteArray> / windows id tables.
//   - QTimeZone::hasAlternativeName().
//   - fromCFTimeZone / toCFTimeZone, fromNSTimeZone / toNSTimeZone (Darwin).
//   - std::chrono interop: fromStdTimeZonePtr, std::chrono::zoned_time integration.
//   - QDataStream / QDebug operators.
//   - Deprecated country() (use Territory() instead).
//   - QTimeZone(zoneId, offsetSeconds, name, abbreviation, territory, comment) — the
//     custom-zone constructor with explicit metadata.
//
// # QDateTime integration (requires additions to qdatetime_c.h/.cpp)
//
// To make QDateTime interoperate with QTimeZone, add these methods to QDateTime
// (they are declared in qdatetime.h but not yet bound in the C ABI):
//
//	QTimeZone timeZone() const;
//	void setTimeZone(const QTimeZone&, TransitionResolution = LegacyBehavior);
//	QDateTime toTimeZone(const QTimeZone&) const;
//	static QDateTime currentDateTime(const QTimeZone&);
//	static QDateTime fromMSecsSinceEpoch(qint64, const QTimeZone&);
//	static QDateTime fromSecsSinceEpoch(qint64, const QTimeZone&);
//	QDateTime startOfDay(const QTimeZone&) const;   // on QDate
//	QDateTime endOfDay(const QTimeZone&) const;     // on QDate
//
// Their C ABI wrappers should follow the existing pattern: take `go_qdatetime_t*` and/or
// `go_qtimezone_t`, return `go_qdatetime_t*` (new heap object) or accept an out-parameter.
//
// # Requirements
//
// Qt >= 6.11, C++17, pkg-config Qt6Core.
package qdatetime
