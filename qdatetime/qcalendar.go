package qdatetime

/*
 QCalendar — биндинг Qt 6.11 QtCore/QCalendar как VALUE-тип.

 QCalendar в Qt — value-тип: одно поле const QCalendarBackend *d_ptr (8 байт на LP64,
 4 на ILP32), trivially copyable, без COW, без user-defined destructor. В Go хранится
 как плоская структура из тех же байт — полная копия. Присваивание Go-переменной создаёт
 независимую копию (value-семантика). Указателей, финализаторов и Close() нет.

 ВАЖНО: Qt-конструктор QCalendar() возвращает ВАЛИДНЫЙ Gregorian-календарь (d_ptr != nil),
 а Go-нулевое значение QCalendar{} имеет d_ptr == nil — это НЕВАЛИДНЫЙ календарь:
   - NewQCalendar()                      → QCalendar()      — валидный Gregorian
   - NewQCalendarFromSystem(sys)         → QCalendar(System)
   - NewQCalendarFromSystemID(id)        → QCalendar(SystemId)
   - NewQCalendarFromAny(data, size, tag) → QCalendar(QAnyStringView)
   - QCalendar{}                         → d_ptr == nil → IsValid() == false

 Все методы — чистые геттеры/запросы; мутирующих методов у QCalendar нет.

 Соглашения:
   - Геттеры: IsValid(), IsGregorian(), IsLunar(), IsLuniSolar(), IsSolar(), IsProleptic(),
     HasYearZero(), MaximumDaysInMonth(), MinimumDaysInMonth(), MaximumMonthsInYear().
   - Запросы: DaysInMonth(month, year), DaysInYear(year), MonthsInYear(year),
     IsDateValid(year, month, day), IsLeapYear(year).
   - Name() → *qtext.QString (владелец — вызывающий, нужно Close()).

 Что не реализовано (см. qcalendar_doc.go):

   - Методы, требующие QDate / QDateTime / QTime (пока нет в этом пакете):
       DateFromParts, PartsFromDate, MatchCenturyToWeekday, DayOfWeek,
       MonthName, StandaloneMonthName, WeekDayName, StandaloneWeekDayName,
       DateTimeToString.
   - AvailableCalendars() — нужен QStringList в qtext.
   - QCalendar::YearMonthDay — не выставлен как Go-тип.
*/

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC -Wno-sfinae-incomplete
#cgo LDFLAGS: -lstdc++
#include "qcalendar_c.h"
*/
import "C"

import (
	"unsafe"

	"github.com/vadimlxx/go-qt6core/qtext"
)

// QCalendar is a value type mirroring Qt6 QCalendar (single pointer).
// The Go zero value is an INVALID calendar (d_ptr == nil).
type QCalendar struct {
	v C.go_qcalendar_t
}

// ---- Construction --------------------------------------------------------

// NewQCalendarFromSystem returns QCalendar(system). For system values outside
// [0, System::Last], the result is an invalid calendar.
func NewQCalendarFromSystem(system System) QCalendar {
	return QCalendar{v: C.qcalendar_from_system(C.int(system))}
}

// NewQCalendarFromSystemID returns QCalendar(SystemId). For ids outside
// [0, System::Last], the result is an invalid calendar.
func NewQCalendarFromSystemID(id SystemID) QCalendar {
	return QCalendar{v: C.qcalendar_from_systemid(C.uintptr_t(id))}
}

// NewQCalendarFromAny returns QCalendar(QAnyStringView).
// tag: 0 = UTF-8, 1 = Latin-1, 2 = UTF-16.
func NewQCalendarFromAny(data unsafe.Pointer, size, tag int) QCalendar {
	return QCalendar{v: C.qcalendar_from_any(data, C.ptrdiff_t(size), C.int(tag))}
}

// NewQCalendar returns QCalendar() — the Gregorian calendar, optimised.
func NewQCalendar(arg ...any) QCalendar {
	if len(arg) == 0 || arg == nil {
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
		tag,  ok3 := arg[2].(int)

		if ok1 && ok2 && ok3 {
			return NewQCalendarFromAny(data, size, tag)
		}
	}

	panic("NewQCalendar unknown argument")
}

// ---- Validity ------------------------------------------------------------

// IsValid reports whether the calendar is valid (d_ptr != nil).
// The Go zero value QCalendar{} is not valid.
func (q QCalendar) IsValid() bool { return bool(C.qcalendar_is_valid(q.v)) }

// ---- Date queries --------------------------------------------------------

func (q QCalendar) DaysInMonth(month int, year int) int {
	return int(C.qcalendar_days_in_month(q.v, C.int(month), C.int(year)))
}

func (q QCalendar) DaysInYear(year int) int {
	return int(C.qcalendar_days_in_year(q.v, C.int(year)))
}

func (q QCalendar) MonthsInYear(year int) int {
	return int(C.qcalendar_months_in_year(q.v, C.int(year)))
}

func (q QCalendar) IsDateValid(year int, month int, day int) bool {
	return bool(C.qcalendar_is_date_valid(q.v, C.int(year), C.int(month), C.int(day)))
}

// ---- Leap years ----------------------------------------------------------

func (q QCalendar) IsLeapYear(year int) bool {
	return bool(C.qcalendar_is_leap_year(q.v, C.int(year)))
}

// ---- Calendar properties -------------------------------------------------

func (q QCalendar) IsGregorian() bool { return bool(C.qcalendar_is_gregorian(q.v)) }
func (q QCalendar) IsLunar() bool     { return bool(C.qcalendar_is_lunar(q.v)) }
func (q QCalendar) IsLuniSolar() bool { return bool(C.qcalendar_is_luni_solar(q.v)) }
func (q QCalendar) IsSolar() bool     { return bool(C.qcalendar_is_solar(q.v)) }
func (q QCalendar) IsProleptic() bool { return bool(C.qcalendar_is_proleptic(q.v)) }
func (q QCalendar) HasYearZero() bool { return bool(C.qcalendar_has_year_zero(q.v)) }

func (q QCalendar) MaximumDaysInMonth() int {
	return int(C.qcalendar_maximum_days_in_month(q.v))
}
func (q QCalendar) MinimumDaysInMonth() int {
	return int(C.qcalendar_minimum_days_in_month(q.v))
}
func (q QCalendar) MaximumMonthsInYear() int {
	return int(C.qcalendar_maximum_months_in_year(q.v))
}

// Name returns a fresh *qtext.QString. Caller must Close() it.
func (q QCalendar) Name() *qtext.QString {
	var out unsafe.Pointer
	C.qcalendar_name(q.v, &out)
	return qtext.QStringFromUnsafePtr(out)
}
