// Package qdatetime provides Go bindings for Qt 6.11 QtCore date/time types:
// QDate, QTime and QDateTime.
//
// QDate and QTime are value types (int64 Julian day, int32 msecs since
// midnight). QDateTime is a heap-allocated Qt class represented by a Go handle:
// the C++ object is released by Close() or by runtime.AddCleanup when the
// wrapper becomes unreachable. Close is optional, idempotent and nil-safe.
//
// # Type mapping
//
//	Qt                                  Go
//	QDateTime                           *qdatetime.QDateTime
//	QDate                               qdatetime.QDate
//	QTime                               qdatetime.QTime
//	QTimeZone                           *qdatetime.QTimeZone
//	QCalendar                           qdatetime.QCalendar
//	Qt::TimeSpec                        qdata.TimeSpec
//	Qt::DateFormat                      qdata.DateFormat
//	QDateTime::TransitionResolution     qdatetime.TransitionResolution
//	QString                             *qtext.QString / string
//	time.Time                           time.Time
//
// QDateTime::TransitionResolution already has its own file
// (transition_resolution.go); this package does not declare a second copy.
//
// # Constructors
//
//	NewQDateTime()                                      QDateTime()
//	NewQDateTimeFromDateTime(date, time)                QDateTime(QDate, QTime)
//	NewQDateTimeFromDateTimeResolve(d, t, resolve)      QDateTime(QDate, QTime, TransitionResolution)
//	NewQDateTimeFromDateTimeZone(d, t, tz)              QDateTime(QDate, QTime, QTimeZone)
//	NewQDateTimeFromDateTimeZoneResolve(d, t, tz, r)    QDateTime(QDate, QTime, QTimeZone, TransitionResolution)
//	NewQDateTimeFromString(s, format)                   QDateTime::fromString(s, Qt::DateFormat)
//	NewQDateTimeFromStringFormat(s, format)             QDateTime::fromString(s, QString)
//	NewQDateTimeFromStringCal(s, format, cal)           QDateTime::fromString(s, QString, QCalendar)
//	NewQDateTimeFromStringBaseYear(s, f, year)          QDateTime::fromString(s, QString, int)
//	NewQDateTimeFromStringBaseYearCal(s, f, y, cal)     QDateTime::fromString(s, QString, int, QCalendar)
//	NewQDateTimeCurrentDateTime()                       QDateTime::currentDateTime()
//	NewQDateTimeCurrentDateTimeUtc()                    QDateTime::currentDateTimeUtc()
//	NewQDateTimeFromMSecsSinceEpoch(ms)                 QDateTime::fromMSecsSinceEpoch(ms)
//	NewQDateTimeFromSecsSinceEpoch(s)                   QDateTime::fromSecsSinceEpoch(s)
//	NewQDateTimeFromTime(time.Time)                     instant taken from time.Time
//
// Zone-aware factories are also available: CurrentDateTimeZone,
// FromMSecsSinceEpochZone, FromSecsSinceEpochZone.
//
// # Accessors
//
//	IsNull, IsValid, IsDaylightTime, Date, Time, TimeSpec, TimeRepresentation,
//	OffsetFromUtc, ToMSecsSinceEpoch, ToSecsSinceEpoch,
//	TimeZoneAbbreviation, TimeZoneAbbreviationString, TimeZone, Equals
//
// # Setters
//
//	SetDate, SetDateResolve, SetTime, SetTimeResolve,
//	SetMSecsSinceEpoch, SetSecsSinceEpoch,
//	SetTimeZone, SetTimeZoneResolve
//
// # Arithmetic and distances
//
//	AddDays, AddMonths, AddYears, AddSecs, AddMSecs,
//	DaysTo, SecsTo, MsecsTo
//
// # Time zones
//
//	ToLocalTime, ToUTC, ToOffsetFromUtc, ToTimeZone
//
// # Strings
//
//	ToString, ToStringString, ToStringFormat, ToStringFormatString,
//	ToStringFormatCal, ToStringFormatCalString
//
// # time.Time
//
//	ToTime() time.Time
//	NewQDateTimeFromTime(time.Time)
//
// ToTime returns the absolute instant in UTC with millisecond precision; a
// null/invalid datetime yields the zero time.Time. Conversion back preserves the
// instant (the time zone is not round-tripped).
//
// # Day boundaries
//
// QDate provides StartOfDay / EndOfDay (plus the Spec and Zone variants), which
// return *QDateTime. Package-level compatibility wrappers QDateStartOfDay,
// QDateEndOfDay, QDateStartOfDayZone and QDateEndOfDayZone delegate to them.
//
// # Not implemented yet
//
//   - std::chrono interop: fromStdTimePoint / fromStdLocalTime /
//     toStdSysMilliseconds / toStdSysSeconds / addDuration.
//   - QDataStream / QDebug operators and qHash.
//   - Darwin: fromCFDate/toCFDate, fromNSDate/toNSDate.
//   - Comparison operators: use Equals/DaysTo/SecsTo/MsecsTo.
//
// # C bridge
//
// QDate, QTime and QDateTime share one C ABI (qdatetime_c.h) and one
// implementation (qdatetime_c.cpp), mirroring QtCore's qdatetime.h /
// qdatetime.cpp. The Go files and tests stay separate per class
// (qdate.go, qtime.go, qdatetime.go).
//
// This package is a subpackage of github.com/vadimlxx/go-qt6core.
package qdatetime
