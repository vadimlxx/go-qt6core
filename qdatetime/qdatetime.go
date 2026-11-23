package qdatetime

/*
QDateTime — handle над C++ QDateTime; владение через runtime.AddCleanup + Close().

# Семантика

QDateTime использует COW (Q_DECLARE_SHARED), но в Go представлен как handle.
Копирование Go-указателя не создаёт новый объект — для независимой копии
используйте Clone().

# Владение

C++ объект освобождается одним из двух способов:
  - Close() — немедленно, детерминированно;
  - runtime.AddCleanup — когда *QDateTime становится недостижим.

Close() опционален, идемпотентен и безопасен на nil. Если Close() и cleanup
гоняются, владение захватывает ровно один участник (атомарный Swap(nil)),
двойного free нет.

# Use after Close

После Close() любой метод, кроме Close(), паникует
("qdatetime: use after Close (or nil receiver)").

# Интеграция с QLocale

handle() возвращает *C.go_qdatetime_t, который QLocale использует напрямую
(см. QLocale.ToStringDateTime*).

# Что нового по сравнению с прошлой версией

  - FromStringCal / FromStringBaseYear / FromStringBaseYearCal —
    QDateTime::fromString с QCalendar и/или baseYear.
  - TimeZoneAbbreviation() *qtext.QString — QDateTime::timeZoneAbbreviation.
  - IsDaylightTime() bool — QDateTime::isDaylightTime.
  - Free-функции QDateStartOfDay / QDateEndOfDay / QDateStartOfDayZone /
    QDateEndOfDayZone — обёртки над QDate::startOfDay/endOfDay.
  - XxxString() варианты для ToString/ToStringFormat/ToStringFormatCal.

# Что всё ещё не реализовано (TODO)

  - TransitionResolution в сигнатурах конструкторов и сеттеров.
  - YearRange.
  - std::chrono interop (fromStdTimePoint / fromStdLocalTime / toStdSysMilliseconds /
    toStdSysSeconds / addDuration) — покрыто через FromMSecsSinceEpoch/AddMSecs.
  - QDataStream / QDebug операторы, qHash.
  - Darwin: fromCFDate/toCFDate, fromNSDate/toNSDate.
*/

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC
#cgo LDFLAGS: -lstdc++
#include "qdatetime_c.h"
*/
import "C"

import (
	"runtime"
	"sync/atomic"
	"unsafe"

	"github.com/vadimlxx/go-qt6core/qdata"
	"github.com/vadimlxx/go-qt6core/qtext"
)

// qdatetimeState holds the C++ QDateTime pointer in an atomic slot so that
// Close() and the runtime cleanup can race safely.
type qdatetimeState struct {
	p atomic.Pointer[byte]
}

func (s *qdatetimeState) cpp() *C.go_qdatetime_t {
	if bp := s.p.Load(); bp != nil {
		return (*C.go_qdatetime_t)(unsafe.Pointer(bp))
	}
	return nil
}

func (s *qdatetimeState) claim() *C.go_qdatetime_t {
	if bp := s.p.Swap(nil); bp != nil {
		return (*C.go_qdatetime_t)(unsafe.Pointer(bp))
	}
	return nil
}

// QDateTime wraps a heap-allocated Qt6 QDateTime (COW handle).
type QDateTime struct {
	state   *qdatetimeState
	cleanup runtime.Cleanup
}

func wrapQDateTime(p *C.go_qdatetime_t) *QDateTime {
	if p == nil {
		return nil
	}
	st := &qdatetimeState{}
	st.p.Store((*byte)(unsafe.Pointer(p)))
	dt := &QDateTime{state: st}
	dt.cleanup = runtime.AddCleanup(dt, func(st *qdatetimeState) {
		if cpp := st.claim(); cpp != nil {
			C.qdatetime_delete(cpp)
		}
	}, st)
	return dt
}

func (dt *QDateTime) handle() *C.go_qdatetime_t {
	if p, ok := dt.tryHandle(); ok {
		return p
	}
	panic("qdatetime: use after Close (or nil receiver)")
}

func (dt *QDateTime) tryHandle() (*C.go_qdatetime_t, bool) {
	if dt == nil || dt.state == nil {
		return nil, false
	}
	if p := dt.state.cpp(); p != nil {
		return p, true
	}
	return nil, false
}

// Close releases the underlying C++ object immediately and cancels the
// cleanup. Close is optional, idempotent, and nil-safe.
// Qt: QDateTime::~QDateTime().
func (dt *QDateTime) Close() {
	if dt == nil || dt.state == nil {
		return
	}
	dt.cleanup.Stop()
	if p := dt.state.claim(); p != nil {
		C.qdatetime_delete(p)
	}
}

// qdateToJD / qdateFromJD bridge value-typed QDate and int64 julian day.
func qdateToJD(d QDate) int64    { return d.ToJulianDay() }
func qdateFromJD(jd int64) QDate { return NewQDateFromJulianDay(jd) }

// qtimeToMSecs / qtimeFromMSecs bridge value-typed QTime and int32 msecs.
func qtimeToMSecs(t QTime) int32    { return int32(t.MsecsSinceStartOfDay()) }
func qtimeFromMSecs(ms int32) QTime { return NewQTimeFromMSecsSinceStartOfDay(int(ms)) }

// cstr returns the internal data pointer of s (nil for empty string).
func cstr(s string) *C.char {
	if len(s) == 0 {
		return nil
	}
	return (*C.char)(unsafe.Pointer(unsafe.StringData(s)))
}

// ---- Construction --------------------------------------------------------

// NewQDateTime — Qt: QDateTime().
func NewQDateTime() *QDateTime { return wrapQDateTime(C.qdatetime_new_default()) }

// NewQDateTimeFromDateTime — Qt: QDateTime(QDate, QTime).
func NewQDateTimeFromDateTime(date QDate, time QTime) *QDateTime {
	return wrapQDateTime(C.qdatetime_new_from_date_time(
		C.int64_t(qdateToJD(date)), C.int32_t(qtimeToMSecs(time))))
}

// NewQDateTimeFromString — Qt: QDateTime::fromString(QStringView, Qt::DateFormat).
func NewQDateTimeFromString(str string, format qdata.DateFormat) *QDateTime {
	dt := wrapQDateTime(C.qdatetime_from_string_utf8(cstr(str), C.ptrdiff_t(len(str)), C.int(format)))
	runtime.KeepAlive(str)
	return dt
}

// NewQDateTimeFromStringFormat — Qt: QDateTime::fromString(QStringView, QStringView).
func NewQDateTimeFromStringFormat(str, format string) *QDateTime {
	dt := wrapQDateTime(C.qdatetime_from_string_format_utf8(
		cstr(str), C.ptrdiff_t(len(str)),
		cstr(format), C.ptrdiff_t(len(format))))
	runtime.KeepAlive(str)
	runtime.KeepAlive(format)
	return dt
}

// NewQDateTimeFromStringCal — Qt: QDateTime::fromString(QStringView, QStringView, QCalendar).
// Uses QLocale::DefaultTwoDigitBaseYear (1900).
func NewQDateTimeFromStringCal(str, format string, cal QCalendar) *QDateTime {
	dt := wrapQDateTime(C.qdatetime_from_string_format_cal(
		cstr(str), C.ptrdiff_t(len(str)),
		cstr(format), C.ptrdiff_t(len(format)),
		cal.v))
	runtime.KeepAlive(str)
	runtime.KeepAlive(format)
	return dt
}

// NewQDateTimeFromStringBaseYear — Qt: QDateTime::fromString(QStringView, QStringView, int).
func NewQDateTimeFromStringBaseYear(str, format string, baseYear int) *QDateTime {
	dt := wrapQDateTime(C.qdatetime_from_string_format_base_year(
		cstr(str), C.ptrdiff_t(len(str)),
		cstr(format), C.ptrdiff_t(len(format)),
		C.int(baseYear)))
	runtime.KeepAlive(str)
	runtime.KeepAlive(format)
	return dt
}

// NewQDateTimeFromStringBaseYearCal — Qt: QDateTime::fromString(QStringView, QStringView, int, QCalendar).
func NewQDateTimeFromStringBaseYearCal(str, format string, baseYear int, cal QCalendar) *QDateTime {
	dt := wrapQDateTime(C.qdatetime_from_string_format_base_year_cal(
		cstr(str), C.ptrdiff_t(len(str)),
		cstr(format), C.ptrdiff_t(len(format)),
		C.int(baseYear), cal.v))
	runtime.KeepAlive(str)
	runtime.KeepAlive(format)
	return dt
}

// Clone — Qt: QDateTime(const QDateTime&).
func (dt *QDateTime) Clone() *QDateTime {
	r := wrapQDateTime(C.qdatetime_clone(dt.handle()))
	runtime.KeepAlive(dt)
	return r
}

// Swap — Qt: QDateTime::swap(QDateTime&).
func (dt *QDateTime) Swap(other *QDateTime) {
	if op, ok := other.tryHandle(); ok {
		C.qdatetime_swap(dt.handle(), op)
		runtime.KeepAlive(dt)
		runtime.KeepAlive(other)
	}
}

// ---- Static constructors -------------------------------------------------

// NewQDateTimeCurrentDateTime — Qt: QDateTime::currentDateTime().
func NewQDateTimeCurrentDateTime() *QDateTime {
	return wrapQDateTime(C.qdatetime_current_date_time())
}

// NewQDateTimeCurrentDateTimeUtc — Qt: QDateTime::currentDateTimeUtc().
func NewQDateTimeCurrentDateTimeUtc() *QDateTime {
	return wrapQDateTime(C.qdatetime_current_date_time_utc())
}

// NewQDateTimeFromMSecsSinceEpoch — Qt: QDateTime::fromMSecsSinceEpoch(qint64).
func NewQDateTimeFromMSecsSinceEpoch(msecs int64) *QDateTime {
	return wrapQDateTime(C.qdatetime_from_msecs_since_epoch(C.int64_t(msecs)))
}

// NewQDateTimeFromSecsSinceEpoch — Qt: QDateTime::fromSecsSinceEpoch(qint64).
func NewQDateTimeFromSecsSinceEpoch(secs int64) *QDateTime {
	return wrapQDateTime(C.qdatetime_from_secs_since_epoch(C.int64_t(secs)))
}

// QDateTimeCurrentMSecsSinceEpoch — Qt: QDateTime::currentMSecsSinceEpoch().
func QDateTimeCurrentMSecsSinceEpoch() int64 {
	return int64(C.qdatetime_current_msecs_since_epoch())
}

// QDateTimeCurrentSecsSinceEpoch — Qt: QDateTime::currentSecsSinceEpoch().
func QDateTimeCurrentSecsSinceEpoch() int64 {
	return int64(C.qdatetime_current_secs_since_epoch())
}

// ---- Accessors -----------------------------------------------------------

// IsNull — Qt: QDateTime::isNull() const.
func (dt *QDateTime) IsNull() bool {
	r := bool(C.qdatetime_is_null(dt.handle()))
	runtime.KeepAlive(dt)
	return r
}

// IsValid — Qt: QDateTime::isValid() const.
func (dt *QDateTime) IsValid() bool {
	r := bool(C.qdatetime_is_valid(dt.handle()))
	runtime.KeepAlive(dt)
	return r
}

// IsDaylightTime — Qt: QDateTime::isDaylightTime() const.
func (dt *QDateTime) IsDaylightTime() bool {
	r := bool(C.qdatetime_is_daylight_time(dt.handle()))
	runtime.KeepAlive(dt)
	return r
}

// Date — Qt: QDateTime::date() const.
func (dt *QDateTime) Date() QDate {
	r := qdateFromJD(int64(C.qdatetime_date(dt.handle())))
	runtime.KeepAlive(dt)
	return r
}

// Time — Qt: QDateTime::time() const.
func (dt *QDateTime) Time() QTime {
	r := qtimeFromMSecs(int32(C.qdatetime_time(dt.handle())))
	runtime.KeepAlive(dt)
	return r
}

// TimeSpec — Qt: QDateTime::timeSpec() const.
func (dt *QDateTime) TimeSpec() qdata.TimeSpec {
	r := qdata.TimeSpec(C.qdatetime_time_spec(dt.handle()))
	runtime.KeepAlive(dt)
	return r
}

// OffsetFromUtc — Qt: QDateTime::offsetFromUtc() const.
func (dt *QDateTime) OffsetFromUtc() int {
	r := int(C.qdatetime_offset_from_utc(dt.handle()))
	runtime.KeepAlive(dt)
	return r
}

// ToMSecsSinceEpoch — Qt: QDateTime::toMSecsSinceEpoch() const.
func (dt *QDateTime) ToMSecsSinceEpoch() int64 {
	r := int64(C.qdatetime_to_msecs_since_epoch(dt.handle()))
	runtime.KeepAlive(dt)
	return r
}

// ToSecsSinceEpoch — Qt: QDateTime::toSecsSinceEpoch() const.
func (dt *QDateTime) ToSecsSinceEpoch() int64 {
	r := int64(C.qdatetime_to_secs_since_epoch(dt.handle()))
	runtime.KeepAlive(dt)
	return r
}

// TimeZoneAbbreviation — Qt: QDateTime::timeZoneAbbreviation() const.
func (dt *QDateTime) TimeZoneAbbreviation() *qtext.QString {
	var out unsafe.Pointer
	C.qdatetime_time_zone_abbreviation(dt.handle(), &out)
	runtime.KeepAlive(dt)
	return qtext.QStringFromUnsafePtr(out)
}

// TimeZoneAbbreviationString — Go-string variant of TimeZoneAbbreviation.
func (dt *QDateTime) TimeZoneAbbreviationString() string {
	s := dt.TimeZoneAbbreviation()
	if s == nil {
		return ""
	}
	defer s.Close()
	return s.GoString()
}

// ---- Setters -------------------------------------------------------------

// SetDate — Qt: QDateTime::setDate(QDate).
func (dt *QDateTime) SetDate(date QDate) {
	C.qdatetime_set_date(dt.handle(), C.int64_t(qdateToJD(date)))
	runtime.KeepAlive(dt)
}

// SetTime — Qt: QDateTime::setTime(QTime).
func (dt *QDateTime) SetTime(time QTime) {
	C.qdatetime_set_time(dt.handle(), C.int32_t(qtimeToMSecs(time)))
	runtime.KeepAlive(dt)
}

// SetMSecsSinceEpoch — Qt: QDateTime::setMSecsSinceEpoch(qint64).
func (dt *QDateTime) SetMSecsSinceEpoch(msecs int64) {
	C.qdatetime_set_msecs_since_epoch(dt.handle(), C.int64_t(msecs))
	runtime.KeepAlive(dt)
}

// SetSecsSinceEpoch — Qt: QDateTime::setSecsSinceEpoch(qint64).
func (dt *QDateTime) SetSecsSinceEpoch(secs int64) {
	C.qdatetime_set_secs_since_epoch(dt.handle(), C.int64_t(secs))
	runtime.KeepAlive(dt)
}

// ---- Arithmetic ----------------------------------------------------------

// AddDays — Qt: QDateTime::addDays(qint64) const.
func (dt *QDateTime) AddDays(days int64) *QDateTime {
	r := wrapQDateTime(C.qdatetime_add_days(dt.handle(), C.int64_t(days)))
	runtime.KeepAlive(dt)
	return r
}

// AddMonths — Qt: QDateTime::addMonths(int) const.
func (dt *QDateTime) AddMonths(months int) *QDateTime {
	r := wrapQDateTime(C.qdatetime_add_months(dt.handle(), C.int(months)))
	runtime.KeepAlive(dt)
	return r
}

// AddYears — Qt: QDateTime::addYears(int) const.
func (dt *QDateTime) AddYears(years int) *QDateTime {
	r := wrapQDateTime(C.qdatetime_add_years(dt.handle(), C.int(years)))
	runtime.KeepAlive(dt)
	return r
}

// AddSecs — Qt: QDateTime::addSecs(qint64) const.
func (dt *QDateTime) AddSecs(secs int64) *QDateTime {
	r := wrapQDateTime(C.qdatetime_add_secs(dt.handle(), C.int64_t(secs)))
	runtime.KeepAlive(dt)
	return r
}

// AddMSecs — Qt: QDateTime::addMSecs(qint64) const.
func (dt *QDateTime) AddMSecs(msecs int64) *QDateTime {
	r := wrapQDateTime(C.qdatetime_add_msecs(dt.handle(), C.int64_t(msecs)))
	runtime.KeepAlive(dt)
	return r
}

// DaysTo — Qt: QDateTime::daysTo(const QDateTime&) const.
func (dt *QDateTime) DaysTo(other *QDateTime) int64 {
	op, ok := other.tryHandle()
	if !ok {
		return 0
	}
	r := int64(C.qdatetime_days_to(dt.handle(), op))
	runtime.KeepAlive(dt)
	runtime.KeepAlive(other)
	return r
}

// SecsTo — Qt: QDateTime::secsTo(const QDateTime&) const.
func (dt *QDateTime) SecsTo(other *QDateTime) int64 {
	op, ok := other.tryHandle()
	if !ok {
		return 0
	}
	r := int64(C.qdatetime_secs_to(dt.handle(), op))
	runtime.KeepAlive(dt)
	runtime.KeepAlive(other)
	return r
}

// MsecsTo — Qt: QDateTime::msecsTo(const QDateTime&) const.
func (dt *QDateTime) MsecsTo(other *QDateTime) int64 {
	op, ok := other.tryHandle()
	if !ok {
		return 0
	}
	r := int64(C.qdatetime_msecs_to(dt.handle(), op))
	runtime.KeepAlive(dt)
	runtime.KeepAlive(other)
	return r
}

// ---- Strings -------------------------------------------------------------

// ToString — Qt: QDateTime::toString(Qt::DateFormat) const.
func (dt *QDateTime) ToString(format qdata.DateFormat) *qtext.QString {
	var out unsafe.Pointer
	C.qdatetime_to_string(dt.handle(), C.int(format), &out)
	runtime.KeepAlive(dt)
	return qtext.QStringFromUnsafePtr(out)
}

// ToStringString — Go-string variant of ToString.
func (dt *QDateTime) ToStringString(format qdata.DateFormat) string {
	s := dt.ToString(format)
	if s == nil {
		return ""
	}
	defer s.Close()
	return s.GoString()
}

// ToStringFormat — Qt: QDateTime::toString(QStringView) const.
func (dt *QDateTime) ToStringFormat(format string) *qtext.QString {
	var out unsafe.Pointer
	C.qdatetime_to_string_format(dt.handle(), cstr(format), C.ptrdiff_t(len(format)), &out)
	runtime.KeepAlive(dt)
	runtime.KeepAlive(format)
	return qtext.QStringFromUnsafePtr(out)
}

// ToStringFormatString — Go-string variant.
func (dt *QDateTime) ToStringFormatString(format string) string {
	s := dt.ToStringFormat(format)
	if s == nil {
		return ""
	}
	defer s.Close()
	return s.GoString()
}

// ToStringFormatCal — Qt: QDateTime::toString(QStringView, QCalendar) const.
func (dt *QDateTime) ToStringFormatCal(format string, cal QCalendar) *qtext.QString {
	var out unsafe.Pointer
	C.qdatetime_to_string_format_cal(dt.handle(),
		cstr(format), C.ptrdiff_t(len(format)), cal.v, &out)
	runtime.KeepAlive(dt)
	runtime.KeepAlive(format)
	return qtext.QStringFromUnsafePtr(out)
}

// ToStringFormatCalString — Go-string variant.
func (dt *QDateTime) ToStringFormatCalString(format string, cal QCalendar) string {
	s := dt.ToStringFormatCal(format, cal)
	if s == nil {
		return ""
	}
	defer s.Close()
	return s.GoString()
}

// ---- Comparison ----------------------------------------------------------

// Equals — Qt: QDateTime::operator==(const QDateTime&).
func (dt *QDateTime) Equals(other *QDateTime) bool {
	op, ok := other.tryHandle()
	if !ok {
		return false
	}
	r := bool(C.qdatetime_equals(dt.handle(), op))
	runtime.KeepAlive(dt)
	runtime.KeepAlive(other)
	return r
}

// ---- Conversions ---------------------------------------------------------

// ToLocalTime — Qt: QDateTime::toLocalTime() const.
func (dt *QDateTime) ToLocalTime() *QDateTime {
	r := wrapQDateTime(C.qdatetime_to_local_time(dt.handle()))
	runtime.KeepAlive(dt)
	return r
}

// ToUTC — Qt: QDateTime::toUTC() const.
func (dt *QDateTime) ToUTC() *QDateTime {
	r := wrapQDateTime(C.qdatetime_to_utc(dt.handle()))
	runtime.KeepAlive(dt)
	return r
}

// ToOffsetFromUtc — Qt: QDateTime::toOffsetFromUtc(int) const.
func (dt *QDateTime) ToOffsetFromUtc(offsetSeconds int) *QDateTime {
	r := wrapQDateTime(C.qdatetime_to_offset_from_utc(dt.handle(), C.int(offsetSeconds)))
	runtime.KeepAlive(dt)
	return r
}

// ---- QTimeZone integration -----------------------------------------------

// TimeZone — Qt: QDateTime::timeZone() const.
func (dt *QDateTime) TimeZone() *QTimeZone {
	r := wrapQTimeZone(C.qdatetime_time_zone(dt.handle()))
	runtime.KeepAlive(dt)
	return r
}

// SetTimeZone — Qt: QDateTime::setTimeZone(const QTimeZone&).
func (dt *QDateTime) SetTimeZone(tz *QTimeZone) {
	if tz == nil || tz.p == nil {
		return
	}
	C.qdatetime_set_time_zone(dt.handle(), tz.p)
	runtime.KeepAlive(dt)
	runtime.KeepAlive(tz)
}

// ToTimeZone — Qt: QDateTime::toTimeZone(const QTimeZone&) const.
func (dt *QDateTime) ToTimeZone(tz *QTimeZone) *QDateTime {
	if tz == nil || tz.p == nil {
		return nil
	}
	r := wrapQDateTime(C.qdatetime_to_time_zone(dt.handle(), tz.p))
	runtime.KeepAlive(dt)
	runtime.KeepAlive(tz)
	return r
}

// CurrentDateTimeZone — Qt: QDateTime::currentDateTime(const QTimeZone&).
func CurrentDateTimeZone(tz *QTimeZone) *QDateTime {
	if tz == nil || tz.p == nil {
		return nil
	}
	r := wrapQDateTime(C.qdatetime_current_date_time_zone(tz.p))
	runtime.KeepAlive(tz)
	return r
}

// FromMSecsSinceEpochZone — Qt: QDateTime::fromMSecsSinceEpoch(qint64, const QTimeZone&).
func FromMSecsSinceEpochZone(msecs int64, tz *QTimeZone) *QDateTime {
	if tz == nil || tz.p == nil {
		return nil
	}
	r := wrapQDateTime(C.qdatetime_from_msecs_since_epoch_zone(C.int64_t(msecs), tz.p))
	runtime.KeepAlive(tz)
	return r
}

// FromSecsSinceEpochZone — Qt: QDateTime::fromSecsSinceEpoch(qint64, const QTimeZone&).
func FromSecsSinceEpochZone(secs int64, tz *QTimeZone) *QDateTime {
	if tz == nil || tz.p == nil {
		return nil
	}
	r := wrapQDateTime(C.qdatetime_from_secs_since_epoch_zone(C.int64_t(secs), tz.p))
	runtime.KeepAlive(tz)
	return r
}

// ---- QDate::startOfDay / endOfDay (free functions) -----------------------
// These wrap QDate methods. They live here (not in qdate.go) so that
// QDateTime can be returned without QDate importing QDateTime types.
// Once qdate.go adds these methods, these free functions can be deprecated.

// QDateStartOfDay — Qt: QDate::startOfDay() const.
func QDateStartOfDay(d QDate) *QDateTime {
	return wrapQDateTime(C.qdate_start_of_day(C.int64_t(qdateToJD(d))))
}

// QDateEndOfDay — Qt: QDate::endOfDay() const.
func QDateEndOfDay(d QDate) *QDateTime {
	return wrapQDateTime(C.qdate_end_of_day(C.int64_t(qdateToJD(d))))
}

// QDateStartOfDayZone — Qt: QDate::startOfDay(const QTimeZone&) const.
func QDateStartOfDayZone(d QDate, tz *QTimeZone) *QDateTime {
	if tz == nil || tz.p == nil {
		return nil
	}
	r := wrapQDateTime(C.qdate_start_of_day_zone(C.int64_t(qdateToJD(d)), tz.p))
	runtime.KeepAlive(tz)
	return r
}

// QDateEndOfDayZone — Qt: QDate::endOfDay(const QTimeZone&) const.
func QDateEndOfDayZone(d QDate, tz *QTimeZone) *QDateTime {
	if tz == nil || tz.p == nil {
		return nil
	}
	r := wrapQDateTime(C.qdate_end_of_day_zone(C.int64_t(qdateToJD(d)), tz.p))
	runtime.KeepAlive(tz)
	return r
}
