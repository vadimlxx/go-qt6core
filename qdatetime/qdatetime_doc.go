// Package qdatetime provides Go bindings for Qt 6.11 QtCore date/time-related types.
//
// This file documents QDate, QTime, and QDateTime. For other classes in the same package,
// see their dedicated files:
//
//   - qlocale_doc.go      — QLocale
//   - qcalendar_doc.go    — QCalendar
//   - qinternal_doc.go    — QInternal, QKeyCombination
//
// # QDate
//
// QDate is a value type in Qt: a single int64 field holding the Julian day number. It is
// trivially copyable, has no COW, and no user-defined destructor. In Go it is stored as a
// plain int64 value, so assignment copies the data; no pointers, no finalizers, no Close().
//
// Null vs zero:
//
//	QDate{}          → Julian day 0    → valid date (1 January 4713 BCE)
//	NewQDate()       → null QDate      → mirrors Qt's QDate() (Julian day INT64_MIN)
//	NewQDateFromYMD(y,m,d)             → explicit date
//
// Note that the Go zero value and Qt's default-constructed QDate are NOT the same: the Go
// zero value is a valid date, while NewQDate() returns the null date. Use IsNull() / IsValid()
// to test, and IsValid() on the zero value returns true.
//
// Construction:
//
//	NewQDate()                          — null QDate (Qt QDate())
//	NewQDateFromYMD(y, m, d)            — QDate(y, m, d) (Gregorian)
//	NewQDateFromYMDCal(y, m, d, cal)    — QDate(y, m, d, cal)
//	NewQDateFromJulianDay(jd)           — QDate::fromJulianDay(jd)
//	NewQDateCurrentDate()               — QDate::currentDate()
//	NewQDateFromString(s, format)       — QDate::fromString(s, Qt::DateFormat)
//	NewQDateFromStringFormat(s, format) — QDate::fromString(s, QString format)
//
// Accessors (Gregorian-optimized):
//
//	Year(), Month(), Day(), DayOfWeek(), DayOfYear(), DaysInMonth(), DaysInYear()
//	WeekNumber(yearNum *int) int        — ISO 8601 week; if yearNum != nil, the ISO year
//	IsNull(), IsValid()
//
// Calendar-aware accessors (take a QCalendar):
//
//	YearCal(cal), MonthCal(cal), DayCal(cal)
//
// Arithmetic (each returns a new QDate; the receiver is unchanged):
//
//	AddDays(days int64)                 — QDate::addDays
//	AddMonths(months int)               — QDate::addMonths (Gregorian)
//	AddYears(years int)                 — QDate::addYears (Gregorian)
//	AddMonthsCal(months, cal)           — QDate::addMonths(months, cal)
//	AddYearsCal(years, cal)             — QDate::addYears(years, cal)
//	DaysTo(other QDate) int64           — QDate::daysTo
//
// In-place setters:
//
//	SetDate(y, m, d) bool               — QDate::setDate (Gregorian); returns true on success
//	SetDateCal(y, m, d, cal) bool       — QDate::setDate(y, m, d, cal)
//	GetDate() (y, m, d int)             — QDate::getDate
//
// Strings (each returns a fresh *qtext.QString; caller must Close() it):
//
//	ToString(format)                    — QDate::toString(Qt::DateFormat)
//	ToStringFormat(format)              — QDate::toString(QString format)
//	ToStringFormatCal(format, cal)      — QDate::toString(QString format, QCalendar cal)
//
// Static:
//
//	QDateIsValid(y, m, d)               — QDate::isValid(y, m, d)
//	QDateIsLeapYear(year)               — QDate::isLeapYear(year)
//
// Conversion:
//
//	(d).ToJulianDay() int64             — QDate::toJulianDay
//
// # QTime
//
// QTime is a value type in Qt: a single int field. Internally Qt stores the milliseconds
// since the start of day, with -1 used as the null sentinel. In Go it is stored as int32;
// the null value is -1, valid values are 0..86399999 (msecs since midnight). Assignment
// copies the value; no pointers, no finalizers, no Close().
//
// Null vs zero:
//
//	QTime{}          → msecs = 0       → valid time 00:00:00.000 (midnight)
//	NewQTime()       → null QTime      → mirrors Qt's QTime() (msecs = -1)
//	NewQTimeFromHMS(h,m,s,ms)          → explicit time
//
// Note that the Go zero value and Qt's default-constructed QTime are NOT the same: the Go
// zero value is midnight, while NewQTime() returns the null time. Use IsNull() / IsValid().
//
// Construction:
//
//	NewQTime()                             — null QTime (Qt QTime())
//	NewQTimeFromHMS(h, m, s, ms)           — QTime(h, m, s, ms)
//	NewQTimeCurrentTime()                  — QTime::currentTime()
//	NewQTimeFromString(s, format)          — QTime::fromString(s, Qt::DateFormat)
//	NewQTimeFromStringFormat(s, format)    — QTime::fromString(s, QString format)
//	NewQTimeFromMSecsSinceStartOfDay(ms)   — QTime::fromMSecsSinceStartOfDay (ms < 0 → null)
//
// Accessors:
//
//	Hour(), Minute(), Second(), Msec()
//	MsecsSinceStartOfDay() int             — 0 for null QTime (matches Qt, loses null flag)
//	IsNull(), IsValid()
//
// Arithmetic (each returns a new QTime):
//
//	AddSecs(secs int)                      — QTime::addSecs
//	AddMSecs(msecs int)                    — QTime::addMSecs
//	SecsTo(other QTime) int                — QTime::secsTo
//	MsecsTo(other QTime) int               — QTime::msecsTo
//
// In-place setter:
//
//	SetHMS(h, m, s, ms) bool               — QTime::setHMS; returns true on success
//
// Strings (each returns a fresh *qtext.QString; caller must Close() it):
//
//	ToString(format)                       — QTime::toString(Qt::DateFormat)
//	ToStringFormat(format)                 — QTime::toString(QString format)
//
// Static:
//
//	QTimeIsValid(h, m, s, ms)              — QTime::isValid(h, m, s, ms)
//
// # QDateTime
//
// QDateTime uses COW (Q_DECLARE_SHARED) and is NOT trivially copyable. In Go it is a
// HANDLE: a struct holding a pointer to a heap-allocated C++ QDateTime. All accessors go
// directly to the C++ object. Use Clone() for an independent copy, Close() to free
// immediately; otherwise a finalizer handles it.
//
// The Go zero value QDateTime{} is invalid (nil pointer). Use NewQDateTime() for a null
// QDateTime that mirrors Qt's QDateTime().
//
// Construction:
//
//	NewQDateTime()                              — null QDateTime (Qt QDateTime())
//	NewQDateTimeFromDateTime(date, time)        — QDateTime(QDate, QTime)
//	(dt).Clone() *QDateTime                     — heap copy
//	NewQDateTimeCurrentDateTime()               — QDateTime::currentDateTime()
//	NewQDateTimeCurrentDateTimeUtc()            — QDateTime::currentDateTimeUtc()
//	NewQDateTimeFromMSecsSinceEpoch(ms)         — QDateTime::fromMSecsSinceEpoch
//	NewQDateTimeFromSecsSinceEpoch(sec)         — QDateTime::fromSecsSinceEpoch
//	NewQDateTimeFromString(s, format)           — QDateTime::fromString(s, Qt::DateFormat)
//	NewQDateTimeFromStringFormat(s, format)     — QDateTime::fromString(s, QString format)
//
// Note on calendar-aware construction: QDateTime has NO constructor taking QCalendar in
// Qt 6.11. QCalendar is used only for formatting (ToStringFormatCal). To construct a
// QDateTime from a calendar-aware QDate, use NewQDateFromYMDCal followed by
// NewQDateTimeFromDateTime.
//
// Accessors:
//
//	Date() QDate                                — QDateTime::date
//	Time() QTime                                — QDateTime::time
//	TimeSpec() TimeSpec                         — QDateTime::timeSpec
//	OffsetFromUtc() int                         — QDateTime::offsetFromUtc
//	ToMSecsSinceEpoch() int64                   — QDateTime::toMSecsSinceEpoch
//	ToSecsSinceEpoch() int64                    — QDateTime::toSecsSinceEpoch
//	IsNull(), IsValid()
//
// In-place setters:
//
//	SetDate(date QDate)                         — QDateTime::setDate
//	SetTime(time QTime)                         — QDateTime::setTime
//	SetMSecsSinceEpoch(ms int64)                — QDateTime::setMSecsSinceEpoch
//	SetSecsSinceEpoch(sec int64)                — QDateTime::setSecsSinceEpoch
//
// Arithmetic (each returns a new *QDateTime; receiver unchanged):
//
//	AddDays(days int64)                         — QDateTime::addDays
//	AddMonths(months int)                       — QDateTime::addMonths
//	AddYears(years int)                         — QDateTime::addYears
//	AddSecs(secs int64)                         — QDateTime::addSecs
//	AddMSecs(msecs int64)                       — QDateTime::addMSecs
//	DaysTo(other *QDateTime) int64              — QDateTime::daysTo
//	SecsTo(other *QDateTime) int64              — QDateTime::secsTo
//	MsecsTo(other *QDateTime) int64             — QDateTime::msecsTo
//
// Strings (each returns a fresh *qtext.QString; caller must Close() it):
//
//	ToString(format)                            — QDateTime::toString(Qt::DateFormat)
//	ToStringFormat(format)                      — QDateTime::toString(QString format)
//	ToStringFormatCal(format, cal)              — QDateTime::toString(QString format, QCalendar)
//
// Comparison:
//
//	(dt).Equals(other *QDateTime) bool          — operator==
//
// Conversions (each returns a new *QDateTime):
//
//	ToLocalTime() *QDateTime                    — QDateTime::toLocalTime
//	ToUTC() *QDateTime                          — QDateTime::toUTC
//	ToOffsetFromUtc(offsetSeconds int) *QDateTime — QDateTime::toOffsetFromUtc
//
// Static utilities:
//
//	QDateTimeCurrentMSecsSinceEpoch() int64     — QDateTime::currentMSecsSinceEpoch
//	QDateTimeCurrentSecsSinceEpoch() int64      — QDateTime::currentSecsSinceEpoch
//
// # Enum types in this package
//
// These enums are declared in dedicated files:
//
//	date_format.go            — DateFormat      (Qt::DateFormat)
//	time_spec.go              — TimeSpec        (Qt::TimeSpec)
//	transition_resolution.go  — TransitionResolution (QDateTime::TransitionResolution)
//	year_range.go             — YearRange       (QDateTime::YearRange)
//
// # Dependencies
//
//   - qtext.QString    — all ToString* methods return *qtext.QString.
//   - QCalendar        — same package; used in *Cal methods.
//
// # Not yet implemented
//
// QDate:
//   - startOfDay / endOfDay (require QDateTime).
//   - fromStdSysDays / toStdSysDays (C++ std::chrono).
//   - QDataStream operators, qHash, fromCFString / toCFString (Darwin).
//
// QTime:
//   - QDataStream operators, qHash.
//
// QDateTime:
//   - Methods requiring QTimeZone: timeZone, setTimeZone, toTimeZone,
//     startOfDay/endOfDay(zone), currentDateTime(zone), fromMSecsSinceEpoch(zone),
//     fromSecsSinceEpoch(zone), fromStdZonedTime.
//   - TransitionResolution integrated into constructor/setter signatures (currently only
//     the enum is exposed).
//   - YearRange used anywhere (currently only the enum is exposed).
//   - std::chrono interop: fromStdTimePoint, fromStdLocalTime, toStdSysMilliseconds,
//     toStdSysSeconds, addDuration.
//   - QDataStream / QDebug operators, qHash.
//   - fromCFDate / toCFDate, fromNSDate / toNSDate (Darwin).
//   - Deprecated overloads: QDateTime(QDate, QTime, Qt::TimeSpec, int), setTimeSpec,
//     setOffsetFromUtc, toTimeSpec, fromMSecsSinceEpoch(msecs, TimeSpec, offset),
//     fromSecsSinceEpoch(secs, TimeSpec, offset).
//
// # Requirements
//
// Qt >= 6.11, C++17, pkg-config Qt6Core.
package qdatetime
