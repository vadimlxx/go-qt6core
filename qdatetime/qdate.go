package qdatetime

/*
 QDate — биндинг Qt 6.11 QtCore/QDate как VALUE-тип.

 QDate хранит юлианский день (int64, 8 байт), trivially copyable, без COW.
 В Go хранится как int64. Присваивание создаёт независимую копию.
 Указателей, финализаторов и Close() нет.

 ВАЖНО:
   - NewQDate() возвращает null-дату (Qt QDate() → юлианский день INT64_MIN).
   - Go-нулевое значение QDate{} = 0 — это валидная дата (1 января 4713 г. до н.э.),
     отличная от null. Для null используйте NewQDate(); для «нулевой» даты — QDate{}.

 Соглашения:
   - Конструкторы: NewQDate(), NewQDateFromYMD, NewQDateFromYMDCal, NewQDateFromJulianDay,
     NewQDateCurrentDate, NewQDateFromString, NewQDateFromStringFormat.
   - Геттеры: Year, Month, Day, DayOfWeek, DayOfYear, DaysInMonth, DaysInYear, WeekNumber,
     IsNull, IsValid.
   - С QCalendar: YearCal, MonthCal, DayCal.
   - Арифметика: AddDays, AddMonths, AddYears, AddMonthsCal, AddYearsCal, DaysTo.
   - Строки: ToString, ToStringFormat, ToStringFormatCal → *qtext.QString.
   - Сеттеры: SetDate, SetDateCal, GetDate.
   - Статические: QDateIsValid, QDateIsLeapYear, ToJulianDay.

 Что не реализовано (см. qdatetime_doc.go):
   - startOfDay/endOfDay (требуют QDateTime).
   - fromStdSysDays/toStdSysDays (C++ chrono).
   - QDataStream операторы, qHash, toCFString/fromCFString.
*/

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC -Wno-sfinae-incomplete
#cgo LDFLAGS: -lstdc++
#include "qdate_c.h"
*/
import "C"

import (
	"runtime"
	"unsafe"

	"github.com/vadimlxx/go-qt6core/qtext"
)

// QDate is a value type mirroring Qt6 QDate (single int64 Julian day).
type QDate int64

// NewQDate returns a null QDate (mirrors Qt QDate()).
func NewQDate() QDate { return QDate(C.qdate_default()) }

// NewQDateFromYMD returns QDate(y, m, d).
func NewQDateFromYMD(y, m, d int) QDate {
	return QDate(C.qdate_from_ymd(C.int(y), C.int(m), C.int(d)))
}

// NewQDateFromYMDCal returns QDate(y, m, d, cal).
func NewQDateFromYMDCal(y, m, d int, cal QCalendar) QDate {
	return QDate(C.qdate_from_ymd_cal(C.int(y), C.int(m), C.int(d), cal.v))
}

// NewQDateFromJulianDay returns QDate::fromJulianDay(jd).
func NewQDateFromJulianDay(jd int64) QDate {
	return QDate(C.qdate_from_julian_day(C.int64_t(jd)))
}

// NewQDateCurrentDate returns QDate::currentDate().
func NewQDateCurrentDate() QDate { return QDate(C.qdate_current_date()) }

// NewQDateFromString returns QDate::fromString(str, format).
func NewQDateFromString(str string, format DateFormat) QDate {
	var p *C.char
	if len(str) > 0 {
		p = (*C.char)(unsafe.Pointer(unsafe.StringData(str)))
	}
	d := QDate(C.qdate_from_string_utf8(p, C.ptrdiff_t(len(str)), C.int(format)))
	runtime.KeepAlive(str)
	return d
}

// NewQDateFromStringFormat returns QDate::fromString(str, formatQString).
func NewQDateFromStringFormat(str, format string) QDate {
	var ps, pf *C.char
	if len(str) > 0 {
		ps = (*C.char)(unsafe.Pointer(unsafe.StringData(str)))
	}
	if len(format) > 0 {
		pf = (*C.char)(unsafe.Pointer(unsafe.StringData(format)))
	}
	d := QDate(C.qdate_from_string_format_utf8(
		ps, C.ptrdiff_t(len(str)),
		pf, C.ptrdiff_t(len(format)),
	))
	runtime.KeepAlive(str)
	runtime.KeepAlive(format)
	return d
}

// ---- Accessors -----------------------------------------------------------

func (d QDate) IsNull() bool       { return bool(C.qdate_is_null(C.int64_t(d))) }
func (d QDate) IsValid() bool      { return bool(C.qdate_is_valid(C.int64_t(d))) }
func (d QDate) Year() int          { return int(C.qdate_year(C.int64_t(d))) }
func (d QDate) Month() int         { return int(C.qdate_month(C.int64_t(d))) }
func (d QDate) Day() int           { return int(C.qdate_day(C.int64_t(d))) }
func (d QDate) DayOfWeek() int     { return int(C.qdate_day_of_week(C.int64_t(d))) }
func (d QDate) DayOfYear() int     { return int(C.qdate_day_of_year(C.int64_t(d))) }
func (d QDate) DaysInMonth() int   { return int(C.qdate_days_in_month(C.int64_t(d))) }
func (d QDate) DaysInYear() int    { return int(C.qdate_days_in_year(C.int64_t(d))) }

func (d QDate) WeekNumber(yearNum *int) int {
	var yn C.int
	r := int(C.qdate_week_number(C.int64_t(d), &yn))
	if yearNum != nil {
		*yearNum = int(yn)
	}
	return r
}

func (d QDate) YearCal(cal QCalendar) int {
	return int(C.qdate_year_cal(C.int64_t(d), cal.v))
}
func (d QDate) MonthCal(cal QCalendar) int {
	return int(C.qdate_month_cal(C.int64_t(d), cal.v))
}
func (d QDate) DayCal(cal QCalendar) int {
	return int(C.qdate_day_cal(C.int64_t(d), cal.v))
}

// ---- Arithmetic ----------------------------------------------------------

func (d QDate) AddDays(days int64) QDate {
	return QDate(C.qdate_add_days(C.int64_t(d), C.int64_t(days)))
}
func (d QDate) AddMonths(months int) QDate {
	return QDate(C.qdate_add_months(C.int64_t(d), C.int(months)))
}
func (d QDate) AddYears(years int) QDate {
	return QDate(C.qdate_add_years(C.int64_t(d), C.int(years)))
}
func (d QDate) AddMonthsCal(months int, cal QCalendar) QDate {
	return QDate(C.qdate_add_months_cal(C.int64_t(d), C.int(months), cal.v))
}
func (d QDate) AddYearsCal(years int, cal QCalendar) QDate {
	return QDate(C.qdate_add_years_cal(C.int64_t(d), C.int(years), cal.v))
}
func (d QDate) DaysTo(other QDate) int64 {
	return int64(C.qdate_days_to(C.int64_t(d), C.int64_t(other)))
}

// ---- Strings -------------------------------------------------------------

func (d QDate) ToString(format DateFormat) *qtext.QString {
	var out unsafe.Pointer
	C.qdate_to_string(C.int64_t(d), C.int(format), &out)
	return qtext.QStringFromUnsafePtr(out)
}

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

// ---- Setters -------------------------------------------------------------

func (d *QDate) SetDate(y, m, day int) bool {
	return bool(C.qdate_set_date((*C.go_qdate_t)(unsafe.Pointer(d)),
		C.int(y), C.int(m), C.int(day)))
}
func (d *QDate) SetDateCal(y, m, day int, cal QCalendar) bool {
	return bool(C.qdate_set_date_cal((*C.go_qdate_t)(unsafe.Pointer(d)),
		C.int(y), C.int(m), C.int(day), cal.v))
}
func (d QDate) GetDate() (y, m, day int) {
	var cy, cm, cd C.int
	C.qdate_get_date(C.int64_t(d), &cy, &cm, &cd)
	return int(cy), int(cm), int(cd)
}

// ---- Static --------------------------------------------------------------

// ToJulianDay returns the Julian day number.
func (d QDate) ToJulianDay() int64 { return int64(C.qdate_to_julian_day(C.int64_t(d))) }

// QDateIsValid mirrors QDate::isValid(y, m, d).
func QDateIsValid(y, m, d int) bool {
	return bool(C.qdate_static_is_valid(C.int(y), C.int(m), C.int(d)))
}

// QDateIsLeapYear mirrors QDate::isLeapYear(year).
func QDateIsLeapYear(year int) bool {
	return bool(C.qdate_static_is_leap_year(C.int(year)))
}
