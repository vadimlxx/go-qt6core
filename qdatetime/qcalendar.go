package qdatetime

/*
QCalendar — биндинг Qt 6.11 QtCore/QCalendar как VALUE-тип.

QCalendar в Qt — value-тип: одно поле const QCalendarBackend *d_ptr (8 байт на LP64,
4 на ILP32), trivially copyable, без COW, без user-defined destructor. В Go хранится
как плоская структура из тех же байт — полная копия. Присваивание Go-переменной создаёт
независимую копию (value-семантика). Указателей, финализаторов и Close() НЕТ.

ВАЖНО: Qt-конструктор QCalendar() возвращает ВАЛИДНЫЙ Gregorian-календарь
(d_ptr != nil), а Go-нулевое значение QCalendar{} имеет d_ptr == nil — это
НЕВАЛИДНЫЙ календарь:
  - NewQCalendar()                       → QCalendar()      — валидный Gregorian
  - NewQCalendarFromSystem(sys)          → QCalendar(System)
  - NewQCalendarFromSystemID(id)         → QCalendar(SystemId)
  - NewQCalendarFromAny(data, size, tag) → QCalendar(QAnyStringView)
  - QCalendar{}                          → d_ptr == nil → IsValid() == false

# Соглашения

  - Геттеры/запросы: IsValid, IsGregorian, IsLunar, IsLuniSolar, IsSolar,
    IsProleptic, HasYearZero, MaximumDaysInMonth, MinimumDaysInMonth,
    MaximumMonthsInYear, DaysInMonth, DaysInYear, MonthsInYear, IsDateValid,
    IsLeapYear.
  - Name() → *qtext.QString; NameString() → string.
  - QDate-конверсии: DateFromParts(year, month, day), DateFromPartsYMD(ymd),
    MatchCenturyToWeekday(ymd, dow), PartsFromDate(date), DayOfWeek(date).
  - Имена месяцев/дней: MonthName, StandaloneMonthName, WeekDayName,
    StandaloneWeekDayName (все принимают *QLocale).
  - Форматирование: DateTimeToString(format, dt, dateOnly, timeOnly, locale).
  - Статическое: AvailableCalendars() []string.

# YearMonthDay

Go-тип YearMonthDay имеет поля Year, Month, Day (int). Соответствует
QCalendar::YearMonthDay. Нулевое значение (0,0,0) — невалидно; для явного
"не задано" используйте (-2147483648, -2147483648, -2147483648) или проверяйте
IsValid().

# Что НЕ реализовано (TODO)

  - QDataStream / QDebug операторы.
  - qHash(QCalendar).
  - QCalendar(SystemId) для id > System::Last возвращает invalid (не UB).
*/

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC
#cgo LDFLAGS: -lstdc++
#include "qcalendar_c.h"
*/
import "C"

import (
	"runtime"
	"unsafe"

	"github.com/vadimlxx/go-qt6core/qtext"
)

// YearMonthDay mirrors QCalendar::YearMonthDay.
type YearMonthDay struct {
	Year  int
	Month int
	Day   int
}

// QCalendar is a value type mirroring Qt6 QCalendar (single pointer).
// The Go zero value is an INVALID calendar (d_ptr == nil).
type QCalendar struct {
	v C.go_qcalendar_t
}

// handle returns the underlying C value (for QLocale integration).
func (q QCalendar) handle() C.go_qcalendar_t { return q.v }

// ---- Construction --------------------------------------------------------

// NewQCalendarFromSystem — Qt: QCalendar(System). Invalid for out-of-range values.
func NewQCalendarFromSystem(system System) QCalendar {
	return QCalendar{v: C.qcalendar_from_system(C.int(system))}
}

// NewQCalendarFromSystemID — Qt: QCalendar(SystemId). Invalid for ids outside
// [0, System::Last] (returns invalid calendar instead of UB).
func NewQCalendarFromSystemID(id SystemID) QCalendar {
	return QCalendar{v: C.qcalendar_from_systemid(C.uintptr_t(id))}
}

// NewQCalendarFromAny — Qt: QCalendar(QAnyStringView).
// tag: 0 = UTF-8, 1 = Latin-1, 2 = UTF-16.
func NewQCalendarFromAny(data unsafe.Pointer, size, tag int) QCalendar {
	return QCalendar{v: C.qcalendar_from_any(data, C.ptrdiff_t(size), C.int(tag))}
}

// NewQCalendar dispatches by argument type:
//
//	NewQCalendar()              → QCalendar()      — valid Gregorian
//	NewQCalendar(System)        → QCalendar(System)
//	NewQCalendar(SystemID)      → QCalendar(SystemId)
//	NewQCalendar(ptr, n, tag)   → QCalendar(QAnyStringView)
func NewQCalendar(arg ...any) QCalendar {
	if len(arg) == 0 || arg[0] == nil {
		return QCalendar{v: C.qcalendar_default()}
	}
	if len(arg) == 1 {
		switch v := arg[0].(type) {
		case System:
			return NewQCalendarFromSystem(v)
		case SystemID:
			return NewQCalendarFromSystemID(v)
		}
	}
	if len(arg) == 3 {
		data, ok1 := arg[0].(unsafe.Pointer)
		size, ok2 := arg[1].(int)
		tag, ok3 := arg[2].(int)
		if ok1 && ok2 && ok3 {
			return NewQCalendarFromAny(data, size, tag)
		}
	}
	panic("NewQCalendar: unsupported argument combination")
}

// ---- Validity ------------------------------------------------------------

// IsValid — Qt: QCalendar::isValid() const.
func (q QCalendar) IsValid() bool { return bool(C.qcalendar_is_valid(q.v)) }

// ---- Date queries --------------------------------------------------------

// DaysInMonth — Qt: QCalendar::daysInMonth(int, int) const.
func (q QCalendar) DaysInMonth(month int, year int) int {
	return int(C.qcalendar_days_in_month(q.v, C.int(month), C.int(year)))
}

// DaysInYear — Qt: QCalendar::daysInYear(int) const.
func (q QCalendar) DaysInYear(year int) int {
	return int(C.qcalendar_days_in_year(q.v, C.int(year)))
}

// MonthsInYear — Qt: QCalendar::monthsInYear(int) const.
func (q QCalendar) MonthsInYear(year int) int {
	return int(C.qcalendar_months_in_year(q.v, C.int(year)))
}

// IsDateValid — Qt: QCalendar::isDateValid(int, int, int) const.
func (q QCalendar) IsDateValid(year int, month int, day int) bool {
	return bool(C.qcalendar_is_date_valid(q.v, C.int(year), C.int(month), C.int(day)))
}

// ---- Leap years ----------------------------------------------------------

// IsLeapYear — Qt: QCalendar::isLeapYear(int) const.
func (q QCalendar) IsLeapYear(year int) bool {
	return bool(C.qcalendar_is_leap_year(q.v, C.int(year)))
}

// ---- Calendar properties -------------------------------------------------

// IsGregorian — Qt: QCalendar::isGregorian() const.
func (q QCalendar) IsGregorian() bool { return bool(C.qcalendar_is_gregorian(q.v)) }

// IsLunar — Qt: QCalendar::isLunar() const.
func (q QCalendar) IsLunar() bool { return bool(C.qcalendar_is_lunar(q.v)) }

// IsLuniSolar — Qt: QCalendar::isLuniSolar() const.
func (q QCalendar) IsLuniSolar() bool { return bool(C.qcalendar_is_luni_solar(q.v)) }

// IsSolar — Qt: QCalendar::isSolar() const.
func (q QCalendar) IsSolar() bool { return bool(C.qcalendar_is_solar(q.v)) }

// IsProleptic — Qt: QCalendar::isProleptic() const.
func (q QCalendar) IsProleptic() bool { return bool(C.qcalendar_is_proleptic(q.v)) }

// HasYearZero — Qt: QCalendar::hasYearZero() const.
func (q QCalendar) HasYearZero() bool { return bool(C.qcalendar_has_year_zero(q.v)) }

// MaximumDaysInMonth — Qt: QCalendar::maximumDaysInMonth() const.
func (q QCalendar) MaximumDaysInMonth() int {
	return int(C.qcalendar_maximum_days_in_month(q.v))
}

// MinimumDaysInMonth — Qt: QCalendar::minimumDaysInMonth() const.
func (q QCalendar) MinimumDaysInMonth() int {
	return int(C.qcalendar_minimum_days_in_month(q.v))
}

// MaximumMonthsInYear — Qt: QCalendar::maximumMonthsInYear() const.
func (q QCalendar) MaximumMonthsInYear() int {
	return int(C.qcalendar_maximum_months_in_year(q.v))
}

// Name — Qt: QCalendar::name() const.
func (q QCalendar) Name() *qtext.QString {
	var out unsafe.Pointer
	C.qcalendar_name(q.v, &out)
	return qtext.QStringFromUnsafePtr(out)
}

// NameString — Go-string variant of Name.
func (q QCalendar) NameString() string {
	s := q.Name()
	if s == nil {
		return ""
	}
	defer s.Close()
	return s.GoString()
}

// ---- QDate conversions ---------------------------------------------------

// DateFromParts — Qt: QCalendar::dateFromParts(int, int, int) const.
func (q QCalendar) DateFromParts(year, month, day int) QDate {
	return qdateFromJD(int64(C.qcalendar_date_from_parts(q.v, C.int(year), C.int(month), C.int(day))))
}

// DateFromPartsYMD — Qt: QCalendar::dateFromParts(const YearMonthDay&) const.
func (q QCalendar) DateFromPartsYMD(ymd YearMonthDay) QDate {
	return qdateFromJD(int64(C.qcalendar_date_from_ymd(q.v, C.go_qymd_t{
		year: C.int(ymd.Year), month: C.int(ymd.Month), day: C.int(ymd.Day),
	})))
}

// MatchCenturyToWeekday — Qt: QCalendar::matchCenturyToWeekday(const YearMonthDay&, int) const.
func (q QCalendar) MatchCenturyToWeekday(ymd YearMonthDay, dow int) QDate {
	return qdateFromJD(int64(C.qcalendar_match_century_to_weekday(q.v, C.go_qymd_t{
		year: C.int(ymd.Year), month: C.int(ymd.Month), day: C.int(ymd.Day),
	}, C.int(dow))))
}

// PartsFromDate — Qt: QCalendar::partsFromDate(QDate) const.
func (q QCalendar) PartsFromDate(date QDate) YearMonthDay {
	y := C.qcalendar_parts_from_date(q.v, C.int64_t(qdateToJD(date)))
	return YearMonthDay{Year: int(y.year), Month: int(y.month), Day: int(y.day)}
}

// DayOfWeek — Qt: QCalendar::dayOfWeek(QDate) const.
func (q QCalendar) DayOfWeek(date QDate) int {
	return int(C.qcalendar_day_of_week(q.v, C.int64_t(qdateToJD(date))))
}

// ---- Month / weekday names ----------------------------------------------

// MonthName — Qt: QCalendar::monthName(const QLocale&, int, int, FormatType) const.
func (q QCalendar) MonthName(locale *QLocale, month, year int, format FormatType) *qtext.QString {
	var lp C.go_qlocale_t
	if locale != nil {
		lp, _ = locale.tryHandle()
	}
	var out unsafe.Pointer
	C.qcalendar_month_name(q.v, lp, C.int(month), C.int(year), C.int(format), &out)
	runtime.KeepAlive(locale)
	return qtext.QStringFromUnsafePtr(out)
}

// MonthNameString — Go-string variant of MonthName.
func (q QCalendar) MonthNameString(locale *QLocale, month, year int, format FormatType) string {
	s := q.MonthName(locale, month, year, format)
	if s == nil {
		return ""
	}
	defer s.Close()
	return s.GoString()
}

// StandaloneMonthName — Qt: QCalendar::standaloneMonthName(const QLocale&, int, int, FormatType) const.
func (q QCalendar) StandaloneMonthName(locale *QLocale, month, year int, format FormatType) *qtext.QString {
	var lp C.go_qlocale_t
	if locale != nil {
		lp, _ = locale.tryHandle()
	}
	var out unsafe.Pointer
	C.qcalendar_standalone_month_name(q.v, lp, C.int(month), C.int(year), C.int(format), &out)
	runtime.KeepAlive(locale)
	return qtext.QStringFromUnsafePtr(out)
}

// StandaloneMonthNameString — Go-string variant.
func (q QCalendar) StandaloneMonthNameString(locale *QLocale, month, year int, format FormatType) string {
	s := q.StandaloneMonthName(locale, month, year, format)
	if s == nil {
		return ""
	}
	defer s.Close()
	return s.GoString()
}

// WeekDayName — Qt: QCalendar::weekDayName(const QLocale&, int, FormatType) const.
func (q QCalendar) WeekDayName(locale *QLocale, day int, format FormatType) *qtext.QString {
	var lp C.go_qlocale_t
	if locale != nil {
		lp, _ = locale.tryHandle()
	}
	var out unsafe.Pointer
	C.qcalendar_week_day_name(q.v, lp, C.int(day), C.int(format), &out)
	runtime.KeepAlive(locale)
	return qtext.QStringFromUnsafePtr(out)
}

// WeekDayNameString — Go-string variant.
func (q QCalendar) WeekDayNameString(locale *QLocale, day int, format FormatType) string {
	s := q.WeekDayName(locale, day, format)
	if s == nil {
		return ""
	}
	defer s.Close()
	return s.GoString()
}

// StandaloneWeekDayName — Qt: QCalendar::standaloneWeekDayName(const QLocale&, int, FormatType) const.
func (q QCalendar) StandaloneWeekDayName(locale *QLocale, day int, format FormatType) *qtext.QString {
	var lp C.go_qlocale_t
	if locale != nil {
		lp, _ = locale.tryHandle()
	}
	var out unsafe.Pointer
	C.qcalendar_standalone_week_day_name(q.v, lp, C.int(day), C.int(format), &out)
	runtime.KeepAlive(locale)
	return qtext.QStringFromUnsafePtr(out)
}

// StandaloneWeekDayNameString — Go-string variant.
func (q QCalendar) StandaloneWeekDayNameString(locale *QLocale, day int, format FormatType) string {
	s := q.StandaloneWeekDayName(locale, day, format)
	if s == nil {
		return ""
	}
	defer s.Close()
	return s.GoString()
}

// ---- Date-time formatting -----------------------------------------------

// DateTimeToString — Qt: QCalendar::dateTimeToString(QStringView, const QDateTime&,
// QDate, QTime, const QLocale&) const.
func (q QCalendar) DateTimeToString(format string, dt *QDateTime, dateOnly QDate,
	timeOnly QTime, locale *QLocale) *qtext.QString {
	var dtp *C.go_qdatetime_t
	if dt != nil {
		dtp, _ = dt.tryHandle()
	}
	var lp C.go_qlocale_t
	if locale != nil {
		lp, _ = locale.tryHandle()
	}
	var out unsafe.Pointer
	C.qcalendar_date_time_to_string(q.v, cstr(format), C.ptrdiff_t(len(format)),
		dtp, C.int64_t(qdateToJD(dateOnly)), C.int32_t(qtimeToMSecs(timeOnly)), lp, &out)
	runtime.KeepAlive(dt)
	runtime.KeepAlive(locale)
	runtime.KeepAlive(format)
	return qtext.QStringFromUnsafePtr(out)
}

// DateTimeToStringString — Go-string variant.
func (q QCalendar) DateTimeToStringString(format string, dt *QDateTime, dateOnly QDate,
	timeOnly QTime, locale *QLocale) string {
	s := q.DateTimeToString(format, dt, dateOnly, timeOnly, locale)
	if s == nil {
		return ""
	}
	defer s.Close()
	return s.GoString()
}

// ---- Available calendars -------------------------------------------------

// AvailableCalendars — Qt: static QCalendar::availableCalendars().
// Returns a fresh []string (each element is a Go string).
func AvailableCalendars() []string {
	var count C.ptrdiff_t
	arr := C.qcalendar_available_calendars(&count)
	if arr == nil {
		return nil
	}
	defer C.qcalendar_free_string_array(arr)
	n := int(count)
	if n == 0 {
		return nil
	}
	src := unsafe.Slice((**C.QString)(unsafe.Pointer(arr)), n)
	out := make([]string, n)
	for i := 0; i < n; i++ {
		s := qtext.QStringFromUnsafePtr(unsafe.Pointer(src[i]))
		out[i] = s.GoString()
		s.Close()
	}
	return out
}