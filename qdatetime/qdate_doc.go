// Package qdatetime provides the QDate binding: a value-type wrapper around
// Qt 6.11 QtCore/QDate.
//
// QDate stores a Julian day as int64. In Go it is QDate (int64): trivially
// copyable, no COW, no pointers, no destructor, no Close.
//
// # Null vs zero
//
//	QDate{}            Julian day 0       valid date (1 January 4713 BCE)
//	NewQDate()         null QDate         mirrors Qt QDate() (Julian day INT64_MIN)
//	NewQDateFromYMD()  explicit date
//
// The Go zero value is therefore not the same as Qt's default-constructed
// QDate. Use IsNull / IsValid to distinguish.
//
// # Type mapping
//
//	Qt                                  Go
//	QDate                               qdatetime.QDate
//	Qt::DateFormat                      qdata.DateFormat
//	QCalendar                           qdatetime.QCalendar
//	QString                             *qtext.QString
//	time.Time                           time.Time
//
// # Constructors
//
//	NewQDate()                                    QDate()
//	NewQDateFromYMD(y, m, d)                      QDate(y, m, d)
//	NewQDateFromYMDCal(y, m, d, cal)              QDate(y, m, d, cal)
//	NewQDateFromJulianDay(jd)                     QDate::fromJulianDay(jd)
//	NewQDateCurrentDate()                         QDate::currentDate()
//	NewQDateFromString(s, format)                 QDate::fromString(s, Qt::DateFormat)
//	NewQDateFromStringFormat(s, format)           QDate::fromString(s, QString)
//	NewQDateFromStringFormatCal(s, format, cal)   QDate::fromString(s, QString, QCalendar)
//	NewQDateFromStringBaseYear(s, f, year)        QDate::fromString(s, QString, int)
//	NewQDateFromStringBaseYearCal(s, f, y, cal)   QDate::fromString(s, QString, int, QCalendar)
//	NewQDateFromTime(time.Time)                   fields taken from time.Time
//
// # Accessors and arithmetic
//
//	Year, Month, Day, DayOfWeek, DayOfYear, DaysInMonth, DaysInYear, WeekNumber,
//	IsNull, IsValid, YearCal, MonthCal, DayCal, ToJulianDay,
//	AddDays, AddMonths, AddYears, AddMonthsCal, AddYearsCal, DaysTo
//
// # Strings and setters
//
//	ToString, ToStringFormat, ToStringFormatCal,
//	SetDate, SetDateCal, GetDate,
//	QDateIsValid, QDateIsLeapYear
//
// # Day boundaries and time.Time
//
//	StartOfDay() / StartOfDaySpec(spec, offset) / StartOfDayZone(tz)
//	EndOfDay()   / EndOfDaySpec(spec, offset)   / EndOfDayZone(tz)
//	ToTime() time.Time
//
// StartOfDay/EndOfDay return a *QDateTime. ToTime returns midnight UTC; a
// null/invalid date yields the zero time.Time.
//
// # Not implemented yet
//
//   - fromStdSysDays / toStdSysDays (std::chrono::sys_days).
//   - addDuration(std::chrono::days).
//   - QDataStream / qHash / QDebug integration.
//   - Darwin: fromCFDate/toCFDate, fromNSDate/toNSDate.
//   - Comparison operators (<, <=, ==, ...): use ToJulianDay or the underlying
//     int64 value.
//
// The C bridge is shared by QDate, QTime and QDateTime in qdatetime_c.h /
// qdatetime_c.cpp, mirroring QtCore's qdatetime.h / qdatetime.cpp.
//
// This package is a subpackage of github.com/vadimlxx/go-qt6core.
package qdatetime
