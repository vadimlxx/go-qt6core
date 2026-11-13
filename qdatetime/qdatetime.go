package qdatetime

/*
 QDateTime — биндинг Qt 6.11 QtCore/QDateTime как HANDLE.

 QDateTime использует COW (Q_DECLARE_SHARED) и НЕ является trivially copyable.
 В Go хранится указатель на C++ QDateTime. Все геттеры/сеттеры идут напрямую в C++.
 Для получения независимой копии используйте Clone().

 Семантика значений vs ссылок:
   - QDateTime — ссылочный тип (handle). Копирование Go-указателя не создаёт новый объект.
   - Для независимой копии — Clone().
   - Go-нулевое значение QDateTime{} невалидно. Используйте NewQDateTime().

 Управление памятью:
   - Конструкторы и Clone создают НОВЫЙ heap-объект.
   - Close() освобождает немедленно; иначе срабатывает finalizer.

 Соглашения:
   - Конструкторы: NewQDateTime (null), NewQDateTimeFromDateTime(date, time).
   - Геттеры: Date, Time, TimeSpec, OffsetFromUtc, ToMSecsSinceEpoch, ToSecsSinceEpoch,
     IsNull, IsValid.
   - Сеттеры: SetDate, SetTime, SetMSecsSinceEpoch, SetSecsSinceEpoch.
   - Арифметика: AddDays, AddMonths, AddYears, AddSecs, AddMSecs,
     DaysTo, SecsTo, MsecsTo.
   - Статические: NewQDateTimeCurrentDateTime, NewQDateTimeCurrentDateTimeUtc,
     NewQDateTimeFromMSecsSinceEpoch, NewQDateTimeFromSecsSinceEpoch,
     QDateTimeCurrentMSecsSinceEpoch, QDateTimeCurrentSecsSinceEpoch.
   - Строки: ToString, ToStringFormat, ToStringFormatCal,
     NewQDateTimeFromString, NewQDateTimeFromStringFormat.
   - Сравнение: Equals.
   - Конвертация: ToLocalTime, ToUTC, ToOffsetFromUtc.

 Что не реализовано (см. qdatetime_doc.go):
   - Методы с QTimeZone (timeZone, setTimeZone, toTimeZone, startOfDay/endOfDay с зоной,
     currentDateTime(zone), fromMSecsSinceEpoch(zone)).
   - TransitionResolution в сигнатурах конструкторов и сеттеров (пока только enum).
   - YearRange.
   - std::chrono interop.
   - QDataStream / QDebug операторы, qHash.
   - fromCFDate/toCFDate, fromNSDate/toNSDate (Darwin).
   - Deprecated overloads (toTimeSpec, setTimeSpec, setOffsetFromUtc).
   - QDateTime(QDate, QTime, QCalendar) — в Qt 6.11 такого конструктора нет.
     Для календарно-зависимого форматирования используйте ToStringFormatCal.
*/

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC -Wno-sfinae-incomplete
#cgo LDFLAGS: -lstdc++
#include "qdatetime_c.h"
*/
import "C"

import (
	"runtime"
	"unsafe"

	"github.com/vadimlxx/go-qt6core/qtext"
)

// QDateTime wraps a heap-allocated Qt6 QDateTime (COW handle).
type QDateTime struct {
	p *C.go_qdatetime_t
}

func wrapQDateTime(p *C.go_qdatetime_t) *QDateTime {
	if p == nil {
		return nil
	}
	dt := &QDateTime{p: p}
	runtime.SetFinalizer(dt, finalizeQDateTime)
	return dt
}

func finalizeQDateTime(dt *QDateTime) {
	if dt.p != nil {
		C.qdatetime_delete(dt.p)
		dt.p = nil
	}
}

// Close releases the underlying C++ object immediately and disables the finalizer.
func (dt *QDateTime) Close() {
	if dt.p != nil {
		C.qdatetime_delete(dt.p)
		dt.p = nil
		runtime.SetFinalizer(dt, nil)
	}
}

// ---- Construction --------------------------------------------------------

// NewQDateTime returns a null QDateTime (mirrors Qt QDateTime()).
func NewQDateTime() *QDateTime { return wrapQDateTime(C.qdatetime_new_default()) }

// NewQDateTimeFromDateTime returns QDateTime(date, time).
// If either date or time is null (see QDate.IsNull / QTime.IsNull), the resulting
// QDateTime is null.
func NewQDateTimeFromDateTime(date QDate, time QTime) *QDateTime {
	return wrapQDateTime(C.qdatetime_new_from_date_time(C.int64_t(date), C.int32_t(time)))
}

// NewQDateTimeFromString returns QDateTime::fromString(str, format).
func NewQDateTimeFromString(str string, format DateFormat) *QDateTime {
	var p *C.char
	if len(str) > 0 {
		p = (*C.char)(unsafe.Pointer(unsafe.StringData(str)))
	}
	dt := wrapQDateTime(C.qdatetime_from_string_utf8(p, C.ptrdiff_t(len(str)), C.int(format)))
	runtime.KeepAlive(str)
	return dt
}

// NewQDateTimeFromStringFormat returns QDateTime::fromString(str, formatQString).
func NewQDateTimeFromStringFormat(str, format string) *QDateTime {
	var ps, pf *C.char
	if len(str) > 0 {
		ps = (*C.char)(unsafe.Pointer(unsafe.StringData(str)))
	}
	if len(format) > 0 {
		pf = (*C.char)(unsafe.Pointer(unsafe.StringData(format)))
	}
	dt := wrapQDateTime(C.qdatetime_from_string_format_utf8(
		ps, C.ptrdiff_t(len(str)),
		pf, C.ptrdiff_t(len(format)),
	))
	runtime.KeepAlive(str)
	runtime.KeepAlive(format)
	return dt
}

// Clone returns a heap copy of dt (independent C++ object).
func (dt *QDateTime) Clone() *QDateTime {
	return wrapQDateTime(C.qdatetime_clone(dt.p))
}

// ---- Static constructors -------------------------------------------------

// NewQDateTimeCurrentDateTime returns QDateTime::currentDateTime().
func NewQDateTimeCurrentDateTime() *QDateTime {
	return wrapQDateTime(C.qdatetime_current_date_time())
}

// NewQDateTimeCurrentDateTimeUtc returns QDateTime::currentDateTimeUtc().
func NewQDateTimeCurrentDateTimeUtc() *QDateTime {
	return wrapQDateTime(C.qdatetime_current_date_time_utc())
}

// NewQDateTimeFromMSecsSinceEpoch returns QDateTime::fromMSecsSinceEpoch(msecs).
func NewQDateTimeFromMSecsSinceEpoch(msecs int64) *QDateTime {
	return wrapQDateTime(C.qdatetime_from_msecs_since_epoch(C.int64_t(msecs)))
}

// NewQDateTimeFromSecsSinceEpoch returns QDateTime::fromSecsSinceEpoch(secs).
func NewQDateTimeFromSecsSinceEpoch(secs int64) *QDateTime {
	return wrapQDateTime(C.qdatetime_from_secs_since_epoch(C.int64_t(secs)))
}

// QDateTimeCurrentMSecsSinceEpoch returns QDateTime::currentMSecsSinceEpoch().
func QDateTimeCurrentMSecsSinceEpoch() int64 {
	return int64(C.qdatetime_current_msecs_since_epoch())
}

// QDateTimeCurrentSecsSinceEpoch returns QDateTime::currentSecsSinceEpoch().
func QDateTimeCurrentSecsSinceEpoch() int64 {
	return int64(C.qdatetime_current_secs_since_epoch())
}

// ---- Accessors -----------------------------------------------------------

func (dt *QDateTime) IsNull() bool  { return bool(C.qdatetime_is_null(dt.p)) }
func (dt *QDateTime) IsValid() bool { return bool(C.qdatetime_is_valid(dt.p)) }
func (dt *QDateTime) Date() QDate   { return QDate(C.qdatetime_date(dt.p)) }
func (dt *QDateTime) Time() QTime   { return QTime(C.qdatetime_time(dt.p)) }
func (dt *QDateTime) TimeSpec() TimeSpec {
	return TimeSpec(C.qdatetime_time_spec(dt.p))
}
func (dt *QDateTime) OffsetFromUtc() int {
	return int(C.qdatetime_offset_from_utc(dt.p))
}
func (dt *QDateTime) ToMSecsSinceEpoch() int64 {
	return int64(C.qdatetime_to_msecs_since_epoch(dt.p))
}
func (dt *QDateTime) ToSecsSinceEpoch() int64 {
	return int64(C.qdatetime_to_secs_since_epoch(dt.p))
}

// ---- Setters -------------------------------------------------------------

func (dt *QDateTime) SetDate(date QDate) {
	C.qdatetime_set_date(dt.p, C.int64_t(date))
}
func (dt *QDateTime) SetTime(time QTime) {
	C.qdatetime_set_time(dt.p, C.int32_t(time))
}
func (dt *QDateTime) SetMSecsSinceEpoch(msecs int64) {
	C.qdatetime_set_msecs_since_epoch(dt.p, C.int64_t(msecs))
}
func (dt *QDateTime) SetSecsSinceEpoch(secs int64) {
	C.qdatetime_set_secs_since_epoch(dt.p, C.int64_t(secs))
}

// ---- Arithmetic (each returns a new *QDateTime) --------------------------

func (dt *QDateTime) AddDays(days int64) *QDateTime {
	return wrapQDateTime(C.qdatetime_add_days(dt.p, C.int64_t(days)))
}
func (dt *QDateTime) AddMonths(months int) *QDateTime {
	return wrapQDateTime(C.qdatetime_add_months(dt.p, C.int(months)))
}
func (dt *QDateTime) AddYears(years int) *QDateTime {
	return wrapQDateTime(C.qdatetime_add_years(dt.p, C.int(years)))
}
func (dt *QDateTime) AddSecs(secs int64) *QDateTime {
	return wrapQDateTime(C.qdatetime_add_secs(dt.p, C.int64_t(secs)))
}
func (dt *QDateTime) AddMSecs(msecs int64) *QDateTime {
	return wrapQDateTime(C.qdatetime_add_msecs(dt.p, C.int64_t(msecs)))
}
func (dt *QDateTime) DaysTo(other *QDateTime) int64 {
	return int64(C.qdatetime_days_to(dt.p, other.p))
}
func (dt *QDateTime) SecsTo(other *QDateTime) int64 {
	return int64(C.qdatetime_secs_to(dt.p, other.p))
}
func (dt *QDateTime) MsecsTo(other *QDateTime) int64 {
	return int64(C.qdatetime_msecs_to(dt.p, other.p))
}

// ---- Strings -------------------------------------------------------------

func (dt *QDateTime) ToString(format DateFormat) *qtext.QString {
	var out unsafe.Pointer
	C.qdatetime_to_string(dt.p, C.int(format), &out)
	return qtext.QStringFromUnsafePtr(out)
}

func (dt *QDateTime) ToStringFormat(format string) *qtext.QString {
	var p *C.char
	if len(format) > 0 {
		p = (*C.char)(unsafe.Pointer(unsafe.StringData(format)))
	}
	var out unsafe.Pointer
	C.qdatetime_to_string_format(dt.p, p, C.ptrdiff_t(len(format)), &out)
	runtime.KeepAlive(format)
	return qtext.QStringFromUnsafePtr(out)
}

func (dt *QDateTime) ToStringFormatCal(format string, cal QCalendar) *qtext.QString {
	var p *C.char
	if len(format) > 0 {
		p = (*C.char)(unsafe.Pointer(unsafe.StringData(format)))
	}
	var out unsafe.Pointer
	C.qdatetime_to_string_format_cal(dt.p, p, C.ptrdiff_t(len(format)), cal.v, &out)
	runtime.KeepAlive(format)
	return qtext.QStringFromUnsafePtr(out)
}

// ---- Comparison ----------------------------------------------------------

// Equals reports whether dt and other hold the same datetime.
func (dt *QDateTime) Equals(other *QDateTime) bool {
	return bool(C.qdatetime_equals(dt.p, other.p))
}

// ---- Conversions (each returns a new *QDateTime) -------------------------

func (dt *QDateTime) ToLocalTime() *QDateTime {
	return wrapQDateTime(C.qdatetime_to_local_time(dt.p))
}
func (dt *QDateTime) ToUTC() *QDateTime {
	return wrapQDateTime(C.qdatetime_to_utc(dt.p))
}
func (dt *QDateTime) ToOffsetFromUtc(offsetSeconds int) *QDateTime {
	return wrapQDateTime(C.qdatetime_to_offset_from_utc(dt.p, C.int(offsetSeconds)))
}
