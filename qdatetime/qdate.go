package qdatetime

/*
	QDate — биндинг Qt 6.11 QtCore/QDate как value-типа.

	QDate хранит юлианский день (int64), trivially copyable, без COW. В Go это
	QDate (int64); присваивание создаёт независимую копию. Указателей,
	финализаторов и Close() нет.

	Важно:
	- NewQDate() возвращает null-дату (Qt QDate(), юлианский день INT64_MIN);
	- нулевое значение QDate{} == 0 — валидная дата (1 января 4713 г. до н.э.),
		отличная от null.

	Реализовано:
	- конструкторы: NewQDate, NewQDateFromYMD, NewQDateFromYMDCal,
		NewQDateFromJulianDay, NewQDateCurrentDate, NewQDateFromString,
		NewQDateFromStringFormat, NewQDateFromStringFormatCal,
		NewQDateFromStringBaseYear, NewQDateFromStringBaseYearCal, NewQDateFromTime;
	- геттеры и календарные варианты, арифметика, строки и сеттеры;
	- startOfDay/endOfDay: методы StartOfDay, StartOfDaySpec, StartOfDayZone,
		EndOfDay, EndOfDaySpec, EndOfDayZone (возвращают *QDateTime);
	- конвертация в time.Time и обратно (ToTime / NewQDateFromTime).

	Не реализовано:
	- fromStdSysDays / toStdSysDays (std::chrono::sys_days);
	- addDuration(std::chrono::days);
	- QDataStream/qHash/QDebug-интеграция;
	- Darwin: fromCFDate/toCFDate, fromNSDate/toNSDate.

	Константы Qt::DateFormat берутся из qdata; собственных enum у QDate нет.
*/

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC -Wno-sfinae-incomplete
#cgo LDFLAGS: -lstdc++
#include <stdint.h>
#include <stdbool.h>
#include <stddef.h>
#include "qdatetime_c.h"
*/
import "C"

import (
	"runtime"
	"time"
	"unsafe"

	"github.com/vadimlxx/go-qt6core/qdata"
	"github.com/vadimlxx/go-qt6core/qtext"
)

// QDate is a value type mirroring Qt6 QDate (single int64 Julian day).
//
// Qt: QDate.
type QDate int64

// NewQDate returns a null QDate.
// Qt: QDate::QDate().
func NewQDate() QDate { return QDate(C.qdate_new_null()) }

// NewQDateFromYMD returns QDate(y, m, d).
// Qt: QDate::QDate(int y, int m, int d).
func NewQDateFromYMD(y, m, d int) QDate {
	return QDate(C.qdate_new_from_ymd(C.int(y), C.int(m), C.int(d)))
}

// NewQDateFromYMDCal returns QDate(y, m, d, cal).
// Qt: QDate::QDate(int y, int m, int d, QCalendar cal).
func NewQDateFromYMDCal(y, m, d int, cal QCalendar) QDate {
	return QDate(C.qdate_from_ymd_cal(C.int(y), C.int(m), C.int(d), cal.v))
}

// NewQDateFromJulianDay returns QDate::fromJulianDay(jd).
// Qt: QDate::fromJulianDay(qint64 jd).
func NewQDateFromJulianDay(jd int64) QDate {
	return QDate(C.qdate_from_julian_day(C.int64_t(jd)))
}

// NewQDateCurrentDate returns the current date.
// Qt: QDate::currentDate().
func NewQDateCurrentDate() QDate { return QDate(C.qdate_current_date()) }

// NewQDateFromString parses str using the given Qt date format.
// Returns a null QDate on failure.
// Qt: QDate::fromString(const QString &string, Qt::DateFormat format).
func NewQDateFromString(str string, format qdata.DateFormat) QDate {
	var p *C.char
	if len(str) > 0 {
		p = (*C.char)(unsafe.Pointer(unsafe.StringData(str)))
	}
	d := QDate(C.qdate_from_string(p, C.ptrdiff_t(len(str)), C.int(format)))
	runtime.KeepAlive(str)
	return d
}

// NewQDateFromStringFormat parses str using the given format string.
// Returns a null QDate on failure.
// Qt: QDate::fromString(const QString &string, const QString &format).
func NewQDateFromStringFormat(str, format string) QDate {
	var ps, pf *C.char
	if len(str) > 0 {
		ps = (*C.char)(unsafe.Pointer(unsafe.StringData(str)))
	}
	if len(format) > 0 {
		pf = (*C.char)(unsafe.Pointer(unsafe.StringData(format)))
	}
	d := QDate(C.qdate_from_string_format(
		ps, C.ptrdiff_t(len(str)),
		pf, C.ptrdiff_t(len(format)),
	))
	runtime.KeepAlive(str)
	runtime.KeepAlive(format)
	return d
}

// NewQDateFromStringFormatCal parses str using format and a calendar.
// Returns a null QDate on failure.
// Qt: QDate::fromString(const QString &string, const QString &format, QCalendar cal).
func NewQDateFromStringFormatCal(str, format string, cal QCalendar) QDate {
	var ps, pf *C.char
	if len(str) > 0 {
		ps = (*C.char)(unsafe.Pointer(unsafe.StringData(str)))
	}
	if len(format) > 0 {
		pf = (*C.char)(unsafe.Pointer(unsafe.StringData(format)))
	}
	d := QDate(C.qdate_from_string_format_cal(
		ps, C.ptrdiff_t(len(str)),
		pf, C.ptrdiff_t(len(format)),
		cal.v,
	))
	runtime.KeepAlive(str)
	runtime.KeepAlive(format)
	return d
}

// NewQDateFromStringBaseYear parses str using format and a two-digit-year base.
// Returns a null QDate on failure.
// Qt: QDate::fromString(const QString &string, const QString &format, int baseYear).
func NewQDateFromStringBaseYear(str, format string, baseYear int) QDate {
	var ps, pf *C.char
	if len(str) > 0 {
		ps = (*C.char)(unsafe.Pointer(unsafe.StringData(str)))
	}
	if len(format) > 0 {
		pf = (*C.char)(unsafe.Pointer(unsafe.StringData(format)))
	}
	d := QDate(C.qdate_from_string_format_base_year(
		ps, C.ptrdiff_t(len(str)),
		pf, C.ptrdiff_t(len(format)),
		C.int(baseYear),
	))
	runtime.KeepAlive(str)
	runtime.KeepAlive(format)
	return d
}

// NewQDateFromStringBaseYearCal parses str using format, base year and calendar.
// Returns a null QDate on failure.
// Qt: QDate::fromString(const QString &string, const QString &format, int baseYear, QCalendar cal).
func NewQDateFromStringBaseYearCal(str, format string, baseYear int, cal QCalendar) QDate {
	var ps, pf *C.char
	if len(str) > 0 {
		ps = (*C.char)(unsafe.Pointer(unsafe.StringData(str)))
	}
	if len(format) > 0 {
		pf = (*C.char)(unsafe.Pointer(unsafe.StringData(format)))
	}
	d := QDate(C.qdate_from_string_format_base_year_cal(
		ps, C.ptrdiff_t(len(str)),
		pf, C.ptrdiff_t(len(format)),
		C.int(baseYear),
		cal.v,
	))
	runtime.KeepAlive(str)
	runtime.KeepAlive(format)
	return d
}

// NewQDateFromTime returns the date part of t.
// Qt: QDate(int year, int month, int day) with the fields taken from time.Time.
func NewQDateFromTime(t time.Time) QDate {
	y, m, d := t.Date()
	return NewQDateFromYMD(y, int(m), d)
}

// IsNull reports whether the date is null (invalid).
// Qt: QDate::isNull() const.
func (d QDate) IsNull() bool { return bool(C.qdate_is_null(C.int64_t(d))) }

// IsValid reports whether the date is valid.
// Qt: QDate::isValid() const.
func (d QDate) IsValid() bool { return bool(C.qdate_is_valid(C.int64_t(d))) }

// Year returns the year in the Gregorian calendar.
// Qt: QDate::year() const.
func (d QDate) Year() int { return int(C.qdate_year(C.int64_t(d))) }

// Month returns the month (1..12).
// Qt: QDate::month() const.
func (d QDate) Month() int { return int(C.qdate_month(C.int64_t(d))) }

// Day returns the day of the month (1..31).
// Qt: QDate::day() const.
func (d QDate) Day() int { return int(C.qdate_day(C.int64_t(d))) }

// DayOfWeek returns the day of the week (1 = Monday .. 7 = Sunday).
// Qt: QDate::dayOfWeek() const.
func (d QDate) DayOfWeek() int { return int(C.qdate_day_of_week(C.int64_t(d))) }

// DayOfYear returns the day of the year (1..366).
// Qt: QDate::dayOfYear() const.
func (d QDate) DayOfYear() int { return int(C.qdate_day_of_year(C.int64_t(d))) }

// DaysInMonth returns the number of days in the date's month.
// Qt: QDate::daysInMonth() const.
func (d QDate) DaysInMonth() int { return int(C.qdate_days_in_month(C.int64_t(d))) }

// DaysInYear returns the number of days in the date's year.
// Qt: QDate::daysInYear() const.
func (d QDate) DaysInYear() int { return int(C.qdate_days_in_year(C.int64_t(d))) }

// WeekNumber returns the ISO 8601 week number; if yearNum is non-nil it also
// receives the week-numbering year.
// Qt: QDate::weekNumber(int *yearNumber = nullptr) const.
func (d QDate) WeekNumber(yearNum *int) int {
	var yn C.int
	r := int(C.qdate_week_number(C.int64_t(d), &yn))
	if yearNum != nil {
		*yearNum = int(yn)
	}
	return r
}

// YearCal returns the year in the given calendar.
// Qt: QDate::year(QCalendar cal) const.
func (d QDate) YearCal(cal QCalendar) int {
	return int(C.qdate_year_cal(C.int64_t(d), cal.v))
}

// MonthCal returns the month in the given calendar.
// Qt: QDate::month(QCalendar cal) const.
func (d QDate) MonthCal(cal QCalendar) int {
	return int(C.qdate_month_cal(C.int64_t(d), cal.v))
}

// DayCal returns the day in the given calendar.
// Qt: QDate::day(QCalendar cal) const.
func (d QDate) DayCal(cal QCalendar) int {
	return int(C.qdate_day_cal(C.int64_t(d), cal.v))
}

// AddDays returns the date shifted by the given number of days.
// Qt: QDate::addDays(qint64 days) const.
func (d QDate) AddDays(days int64) QDate {
	return QDate(C.qdate_add_days(C.int64_t(d), C.int64_t(days)))
}

// AddMonths returns the date shifted by the given number of months.
// Qt: QDate::addMonths(int months) const.
func (d QDate) AddMonths(months int) QDate {
	return QDate(C.qdate_add_months(C.int64_t(d), C.int(months)))
}

// AddYears returns the date shifted by the given number of years.
// Qt: QDate::addYears(int years) const.
func (d QDate) AddYears(years int) QDate {
	return QDate(C.qdate_add_years(C.int64_t(d), C.int(years)))
}

// AddMonthsCal returns the date shifted by months in the given calendar.
// Qt: QDate::addMonths(int months, QCalendar cal) const.
func (d QDate) AddMonthsCal(months int, cal QCalendar) QDate {
	return QDate(C.qdate_add_months_cal(C.int64_t(d), C.int(months), cal.v))
}

// AddYearsCal returns the date shifted by years in the given calendar.
// Qt: QDate::addYears(int years, QCalendar cal) const.
func (d QDate) AddYearsCal(years int, cal QCalendar) QDate {
	return QDate(C.qdate_add_years_cal(C.int64_t(d), C.int(years), cal.v))
}

// DaysTo returns the number of days from d to other.
// Qt: QDate::daysTo(QDate d) const.
func (d QDate) DaysTo(other QDate) int64 {
	return int64(C.qdate_days_to(C.int64_t(d), C.int64_t(other)))
}

// ToString formats the date using the given Qt date format.
// Qt: QDate::toString(Qt::DateFormat format) const.
func (d QDate) ToString(format qdata.DateFormat) *qtext.QString {
	var out unsafe.Pointer
	C.qdate_to_string(C.int64_t(d), C.int(format), &out)
	return qtext.QStringFromUnsafePtr(out)
}

// ToStringFormat formats the date using a format string.
// Qt: QDate::toString(const QString &format) const.
func (d QDate) ToStringFormat(format string) *qtext.QString {
	var p *C.char
	if len(format) > 0 {
		p = (*C.char)(unsafe.Pointer(unsafe.StringData(format)))
	}
	var out unsafe.Pointer
	C.qdate_to_string_format(C.int64_t(d), p, C.ptrdiff_t(len(format)), &out)
	runtime.KeepAlive(format)
	return qtext.QStringFromUnsafePtr(out)
}

// ToStringFormatCal formats the date using a format string and a calendar.
// Qt: QDate::toString(const QString &format, QCalendar cal) const.
func (d QDate) ToStringFormatCal(format string, cal QCalendar) *qtext.QString {
	var p *C.char
	if len(format) > 0 {
		p = (*C.char)(unsafe.Pointer(unsafe.StringData(format)))
	}
	var out unsafe.Pointer
	C.qdate_to_string_format_cal(C.int64_t(d), p, C.ptrdiff_t(len(format)), cal.v, &out)
	runtime.KeepAlive(format)
	return qtext.QStringFromUnsafePtr(out)
}

// SetDate changes the date to y-m-day; returns false for an invalid date.
// Qt: QDate::setDate(int year, int month, int day).
func (d *QDate) SetDate(y, m, day int) bool {
	return bool(C.qdate_set_date((*C.go_qdate_t)(unsafe.Pointer(d)),
		C.int(y), C.int(m), C.int(day)))
}

// SetDateCal changes the date in the given calendar.
// Qt: QDate::setDate(int year, int month, int day, QCalendar cal).
func (d *QDate) SetDateCal(y, m, day int, cal QCalendar) bool {
	return bool(C.qdate_set_date_cal((*C.go_qdate_t)(unsafe.Pointer(d)),
		C.int(y), C.int(m), C.int(day), cal.v))
}

// GetDate returns the year, month and day of the date.
// Qt: QDate::getDate(int *year, int *month, int *day) const.
func (d QDate) GetDate() (y, m, day int) {
	var cy, cm, cd C.int
	C.qdate_get_date(C.int64_t(d), &cy, &cm, &cd)
	return int(cy), int(cm), int(cd)
}

// ToJulianDay returns the Julian day number.
// Qt: QDate::toJulianDay() const.
func (d QDate) ToJulianDay() int64 { return int64(C.qdate_to_julian_day(C.int64_t(d))) }

// QDateIsValid mirrors QDate::isValid(y, m, d).
// Qt: QDate::isValid(int y, int m, int d) static.
func QDateIsValid(y, m, d int) bool {
	return bool(C.qdate_static_is_valid(C.int(y), C.int(m), C.int(d)))
}

// QDateIsLeapYear mirrors QDate::isLeapYear(year).
// Qt: QDate::isLeapYear(int year) static.
func QDateIsLeapYear(year int) bool {
	return bool(C.qdate_static_is_leap_year(C.int(year)))
}

// ---- day boundaries -----------------------------------------------------

// StartOfDay returns the first moment of the day in local time.
// Qt: QDate::startOfDay() const.
func (d QDate) StartOfDay() *QDateTime {
	return wrapQDateTime(C.qdate_start_of_day(C.int64_t(d)))
}

// StartOfDaySpec returns the first moment of the day for an explicit time spec.
// Qt: QDate::startOfDay(Qt::TimeSpec spec, int offsetSeconds) const.
func (d QDate) StartOfDaySpec(spec qdata.TimeSpec, offsetSeconds int) *QDateTime {
	return wrapQDateTime(C.qdate_start_of_day_spec(C.int64_t(d), C.int(spec), C.int(offsetSeconds)))
}

// StartOfDayZone returns the first moment of the day in the given zone.
// Qt: QDate::startOfDay(const QTimeZone &zone) const.
func (d QDate) StartOfDayZone(tz *QTimeZone) *QDateTime {
	if tz == nil || tz.p == nil {
		panic("qdate: StartOfDayZone with nil or closed QTimeZone")
	}
	r := wrapQDateTime(C.qdate_start_of_day_zone(C.int64_t(d), tz.p))
	runtime.KeepAlive(tz)
	return r
}

// EndOfDay returns the last moment of the day in local time (23:59:59.999).
// Qt: QDate::endOfDay() const.
func (d QDate) EndOfDay() *QDateTime {
	return wrapQDateTime(C.qdate_end_of_day(C.int64_t(d)))
}

// EndOfDaySpec returns the last moment of the day for an explicit time spec.
// Qt: QDate::endOfDay(Qt::TimeSpec spec, int offsetSeconds) const.
func (d QDate) EndOfDaySpec(spec qdata.TimeSpec, offsetSeconds int) *QDateTime {
	return wrapQDateTime(C.qdate_end_of_day_spec(C.int64_t(d), C.int(spec), C.int(offsetSeconds)))
}

// EndOfDayZone returns the last moment of the day in the given zone.
// Qt: QDate::endOfDay(const QTimeZone &zone) const.
func (d QDate) EndOfDayZone(tz *QTimeZone) *QDateTime {
	if tz == nil || tz.p == nil {
		panic("qdate: EndOfDayZone with nil or closed QTimeZone")
	}
	r := wrapQDateTime(C.qdate_end_of_day_zone(C.int64_t(d), tz.p))
	runtime.KeepAlive(tz)
	return r
}

// ---- time.Time interop --------------------------------------------------

// ToTime returns the date at 00:00:00 UTC. A null/invalid QDate yields the zero
// time.Time.
// Qt: QDate::year()/month()/day() combined into a Go time.Time (no direct Qt
// counterpart — QDate has no std::chrono conversion for time.Time).
func (d QDate) ToTime() time.Time {
	if !d.IsValid() {
		return time.Time{}
	}
	return time.Date(d.Year(), time.Month(d.Month()), d.Day(), 0, 0, 0, 0, time.UTC)
}
