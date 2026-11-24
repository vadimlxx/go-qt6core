/* Order: version check first, then real Qt headers. */
#include "../internal/qt6core/check_qt_version.h"

#include <QDate>
#include <QTime>
#include <QDateTime>
#include <QTimeZone>
#include <QString>
#include <QCalendar>

#include "qdatetime_c.h"

namespace {

/* Value-type conversions. */
inline QDate D(go_qdate_t d) { return QDate::fromJulianDay(d); }
inline go_qdate_t toDateGo(const QDate& d) { return d.toJulianDay(); }

inline QTime T(go_qtime_t t) {
    if (t < 0) return QTime();
    return QTime::fromMSecsSinceStartOfDay(t);
}
inline go_qtime_t toTimeGo(const QTime& t) {
    if (t.isNull()) return -1;
    return t.msecsSinceStartOfDay();
}

/* Go strings are NOT null-terminated. Use the (data, len) overload. */
inline QString U8(const char* data, ptrdiff_t len) {
    return QString::fromUtf8(data, static_cast<qsizetype>(len));
}

/* Single conversion point for the int-based enums crossing the C ABI. */
inline Qt::DateFormat DF(int v) noexcept {
    return static_cast<Qt::DateFormat>(v);
}
inline QDateTime::TransitionResolution TR(int v) noexcept {
    return static_cast<QDateTime::TransitionResolution>(v);
}
inline Qt::TimeSpec TS(int v) noexcept {
    return static_cast<Qt::TimeSpec>(v);
}

/* QDate::startOfDay(Qt::TimeSpec, int) is deprecated since Qt 6.9; build the
 * equivalent QTimeZone and use the non-deprecated overload. */
inline QTimeZone zoneFromSpec(int spec, int offsetSeconds) {
    switch (TS(spec)) {
    case Qt::LocalTime:     return QTimeZone(QTimeZone::LocalTime);
    case Qt::UTC:           return QTimeZone(QTimeZone::UTC);
    case Qt::OffsetFromUTC: return QTimeZone::fromSecondsAheadOfUtc(offsetSeconds);
    default:                return QTimeZone();
    }
}

} // namespace

extern "C" {

/* ==== QDate ==== */

go_qdate_t qdate_new_null(void)                    { return toDateGo(QDate()); }
go_qdate_t qdate_new_from_ymd(int y, int m, int d) { return toDateGo(QDate(y, m, d)); }
go_qdate_t qdate_from_ymd_cal(int y, int m, int d, go_qcalendar_t cal) {
    return toDateGo(QDate(y, m, d, cal));
}
go_qdate_t qdate_from_julian_day(int64_t jd) { return toDateGo(QDate::fromJulianDay(jd)); }
go_qdate_t qdate_current_date(void)          { return toDateGo(QDate::currentDate()); }

int64_t qdate_to_julian_day(go_qdate_t d) { return D(d).toJulianDay(); }
bool qdate_is_null(go_qdate_t d)          { return D(d).isNull(); }
bool qdate_is_valid(go_qdate_t d)         { return D(d).isValid(); }
int  qdate_year(go_qdate_t d)             { return D(d).year(); }
int  qdate_month(go_qdate_t d)            { return D(d).month(); }
int  qdate_day(go_qdate_t d)              { return D(d).day(); }
int  qdate_day_of_week(go_qdate_t d)      { return D(d).dayOfWeek(); }
int  qdate_day_of_year(go_qdate_t d)      { return D(d).dayOfYear(); }
int  qdate_days_in_month(go_qdate_t d)    { return D(d).daysInMonth(); }
int  qdate_days_in_year(go_qdate_t d)     { return D(d).daysInYear(); }
int  qdate_week_number(go_qdate_t d, int* yn) { return D(d).weekNumber(yn); }

int qdate_year_cal(go_qdate_t d, go_qcalendar_t cal)  { return D(d).year(cal); }
int qdate_month_cal(go_qdate_t d, go_qcalendar_t cal) { return D(d).month(cal); }
int qdate_day_cal(go_qdate_t d, go_qcalendar_t cal)   { return D(d).day(cal); }

go_qdate_t qdate_add_days(go_qdate_t d, int64_t days) {
    return toDateGo(D(d).addDays(days));
}
go_qdate_t qdate_add_months(go_qdate_t d, int months) {
    return toDateGo(D(d).addMonths(months));
}
go_qdate_t qdate_add_years(go_qdate_t d, int years) {
    return toDateGo(D(d).addYears(years));
}
go_qdate_t qdate_add_months_cal(go_qdate_t d, int months, go_qcalendar_t cal) {
    return toDateGo(D(d).addMonths(months, cal));
}
go_qdate_t qdate_add_years_cal(go_qdate_t d, int years, go_qcalendar_t cal) {
    return toDateGo(D(d).addYears(years, cal));
}
int64_t qdate_days_to(go_qdate_t d, go_qdate_t other) {
    return D(d).daysTo(D(other));
}

/* In-place setters: write the mutated QDate back to *d. */
bool qdate_set_date(go_qdate_t* d, int y, int m, int day) {
    QDate qd = QDate::fromJulianDay(*d);
    bool ok = qd.setDate(y, m, day);
    *d = qd.toJulianDay();
    return ok;
}
bool qdate_set_date_cal(go_qdate_t* d, int y, int m, int day, go_qcalendar_t cal) {
    QDate qd = QDate::fromJulianDay(*d);
    bool ok = qd.setDate(y, m, day, cal);
    *d = qd.toJulianDay();
    return ok;
}
void qdate_get_date(go_qdate_t d, int* y, int* m, int* day) {
    D(d).getDate(y, m, day);
}

void qdate_to_string(go_qdate_t d, int format, void** out_qstr) {
    *out_qstr = new QString(D(d).toString(DF(format)));
}
void qdate_to_string_format(go_qdate_t d, const char* fmt, ptrdiff_t fmt_len,
                            void** out_qstr) {
    *out_qstr = new QString(D(d).toString(U8(fmt, fmt_len)));
}
void qdate_to_string_format_cal(go_qdate_t d, const char* fmt, ptrdiff_t fmt_len,
                                go_qcalendar_t cal, void** out_qstr) {
    *out_qstr = new QString(D(d).toString(U8(fmt, fmt_len), cal));
}

go_qdate_t qdate_from_string(const char* str, ptrdiff_t str_len, int format) {
    return toDateGo(QDate::fromString(U8(str, str_len), DF(format)));
}
go_qdate_t qdate_from_string_format(const char* str, ptrdiff_t str_len,
                                    const char* fmt, ptrdiff_t fmt_len) {
    return toDateGo(QDate::fromString(U8(str, str_len), U8(fmt, fmt_len)));
}
go_qdate_t qdate_from_string_format_cal(const char* str, ptrdiff_t str_len,
                                        const char* fmt, ptrdiff_t fmt_len,
                                        go_qcalendar_t cal) {
    return toDateGo(QDate::fromString(U8(str, str_len), U8(fmt, fmt_len), cal));
}
go_qdate_t qdate_from_string_format_base_year(const char* str, ptrdiff_t str_len,
                                              const char* fmt, ptrdiff_t fmt_len,
                                              int base_year) {
    return toDateGo(QDate::fromString(U8(str, str_len), U8(fmt, fmt_len), base_year));
}
go_qdate_t qdate_from_string_format_base_year_cal(const char* str, ptrdiff_t str_len,
                                                  const char* fmt, ptrdiff_t fmt_len,
                                                  int base_year, go_qcalendar_t cal) {
    return toDateGo(QDate::fromString(U8(str, str_len), U8(fmt, fmt_len), base_year, cal));
}

bool qdate_static_is_valid(int y, int m, int d) { return QDate::isValid(y, m, d); }
bool qdate_static_is_leap_year(int year)        { return QDate::isLeapYear(year); }

go_qdatetime_t* qdate_start_of_day(go_qdate_t date) {
    return new QDateTime(D(date).startOfDay());
}
go_qdatetime_t* qdate_end_of_day(go_qdate_t date) {
    return new QDateTime(D(date).endOfDay());
}
go_qdatetime_t* qdate_start_of_day_spec(go_qdate_t date, int spec, int offset_seconds) {
    return new QDateTime(D(date).startOfDay(zoneFromSpec(spec, offset_seconds)));
}
go_qdatetime_t* qdate_end_of_day_spec(go_qdate_t date, int spec, int offset_seconds) {
    return new QDateTime(D(date).endOfDay(zoneFromSpec(spec, offset_seconds)));
}
go_qdatetime_t* qdate_start_of_day_zone(go_qdate_t date, const go_qtimezone_t tz) {
    return new QDateTime(D(date).startOfDay(*tz));
}
go_qdatetime_t* qdate_end_of_day_zone(go_qdate_t date, const go_qtimezone_t tz) {
    return new QDateTime(D(date).endOfDay(*tz));
}

/* ==== QTime ==== */

go_qtime_t qtime_default(void)                         { return toTimeGo(QTime()); }
go_qtime_t qtime_from_hms(int h, int m, int s, int ms) { return toTimeGo(QTime(h, m, s, ms)); }
go_qtime_t qtime_current_time(void)                    { return toTimeGo(QTime::currentTime()); }
go_qtime_t qtime_from_msecs_since_start_of_day(int msecs) {
    if (msecs < 0) return -1;
    return toTimeGo(QTime::fromMSecsSinceStartOfDay(msecs));
}

go_qtime_t qtime_from_string_utf8(const char* str, ptrdiff_t str_len, int format) {
    return toTimeGo(QTime::fromString(U8(str, str_len), DF(format)));
}
go_qtime_t qtime_from_string_format_utf8(const char* str, ptrdiff_t str_len,
                                         const char* format, ptrdiff_t format_len) {
    return toTimeGo(QTime::fromString(U8(str, str_len), U8(format, format_len)));
}

bool qtime_is_null(go_qtime_t t)                  { return T(t).isNull(); }
bool qtime_is_valid(go_qtime_t t)                 { return T(t).isValid(); }
int  qtime_hour(go_qtime_t t)                     { return T(t).hour(); }
int  qtime_minute(go_qtime_t t)                   { return T(t).minute(); }
int  qtime_second(go_qtime_t t)                   { return T(t).second(); }
int  qtime_msec(go_qtime_t t)                     { return T(t).msec(); }
int  qtime_msecs_since_start_of_day(go_qtime_t t) { return T(t).msecsSinceStartOfDay(); }

go_qtime_t qtime_add_secs(go_qtime_t t, int secs)   { return toTimeGo(T(t).addSecs(secs)); }
go_qtime_t qtime_add_msecs(go_qtime_t t, int msecs) { return toTimeGo(T(t).addMSecs(msecs)); }
int qtime_secs_to(go_qtime_t t, go_qtime_t other)   { return T(t).secsTo(T(other)); }
int qtime_msecs_to(go_qtime_t t, go_qtime_t other)  { return T(t).msecsTo(T(other)); }

void qtime_to_string(go_qtime_t t, int format, void** out_qstr) {
    *out_qstr = new QString(T(t).toString(DF(format)));
}
void qtime_to_string_format(go_qtime_t t, const char* format, ptrdiff_t format_len,
                            void** out_qstr) {
    *out_qstr = new QString(T(t).toString(U8(format, format_len)));
}

bool qtime_set_hms(go_qtime_t* t, int h, int m, int s, int ms) {
    QTime qt = T(*t);
    bool ok = qt.setHMS(h, m, s, ms);
    *t = toTimeGo(qt);
    return ok;
}

bool qtime_static_is_valid(int h, int m, int s, int ms) { return QTime::isValid(h, m, s, ms); }

/* ==== QDateTime ==== */

go_qdatetime_t* qdatetime_new_default(void) { return new QDateTime(); }

go_qdatetime_t* qdatetime_new_from_date_time(go_qdate_t date, go_qtime_t time) {
    return new QDateTime(D(date), T(time));
}
go_qdatetime_t* qdatetime_new_from_date_time_resolve(go_qdate_t date, go_qtime_t time,
                                                     int resolve) {
    return new QDateTime(D(date), T(time), TR(resolve));
}
go_qdatetime_t* qdatetime_new_from_date_time_zone(go_qdate_t date, go_qtime_t time,
                                                  const go_qtimezone_t tz) {
    return new QDateTime(D(date), T(time), *tz);
}
go_qdatetime_t* qdatetime_new_from_date_time_zone_resolve(go_qdate_t date, go_qtime_t time,
                                                          const go_qtimezone_t tz,
                                                          int resolve) {
    return new QDateTime(D(date), T(time), *tz, TR(resolve));
}

go_qdatetime_t* qdatetime_new_copy(const go_qdatetime_t* other) {
    return new QDateTime(*other);
}
go_qdatetime_t* qdatetime_clone(const go_qdatetime_t* other) {
    return new QDateTime(*other);
}
void qdatetime_delete(go_qdatetime_t* dt) { delete dt; }
void qdatetime_swap(go_qdatetime_t* self, go_qdatetime_t* other) {
    if (self && other) self->swap(*other);
}

bool qdatetime_is_null(const go_qdatetime_t* dt)          { return dt->isNull(); }
bool qdatetime_is_valid(const go_qdatetime_t* dt)         { return dt->isValid(); }
bool qdatetime_is_daylight_time(const go_qdatetime_t* dt) { return dt->isDaylightTime(); }

go_qdate_t qdatetime_date(const go_qdatetime_t* dt) { return dt->date().toJulianDay(); }
go_qtime_t qdatetime_time(const go_qdatetime_t* dt) { return toTimeGo(dt->time()); }
int qdatetime_time_spec(const go_qdatetime_t* dt)   { return static_cast<int>(dt->timeSpec()); }
go_qtimezone_t qdatetime_time_representation(const go_qdatetime_t* dt) {
    return new QTimeZone(dt->timeRepresentation());
}
int qdatetime_offset_from_utc(const go_qdatetime_t* dt) { return dt->offsetFromUtc(); }
int64_t qdatetime_to_msecs_since_epoch(const go_qdatetime_t* dt) { return dt->toMSecsSinceEpoch(); }
int64_t qdatetime_to_secs_since_epoch(const go_qdatetime_t* dt)  { return dt->toSecsSinceEpoch(); }
void qdatetime_time_zone_abbreviation(const go_qdatetime_t* dt, void** out_qstr) {
    *out_qstr = new QString(dt->timeZoneAbbreviation());
}

void qdatetime_set_date(go_qdatetime_t* dt, go_qdate_t date) { dt->setDate(D(date)); }
void qdatetime_set_date_resolve(go_qdatetime_t* dt, go_qdate_t date, int resolve) {
    dt->setDate(D(date), TR(resolve));
}
void qdatetime_set_time(go_qdatetime_t* dt, go_qtime_t time) { dt->setTime(T(time)); }
void qdatetime_set_time_resolve(go_qdatetime_t* dt, go_qtime_t time, int resolve) {
    dt->setTime(T(time), TR(resolve));
}
void qdatetime_set_msecs_since_epoch(go_qdatetime_t* dt, int64_t msecs) {
    dt->setMSecsSinceEpoch(msecs);
}
void qdatetime_set_secs_since_epoch(go_qdatetime_t* dt, int64_t secs) {
    dt->setSecsSinceEpoch(secs);
}

go_qdatetime_t* qdatetime_add_days(const go_qdatetime_t* dt, int64_t days) {
    return new QDateTime(dt->addDays(days));
}
go_qdatetime_t* qdatetime_add_months(const go_qdatetime_t* dt, int months) {
    return new QDateTime(dt->addMonths(months));
}
go_qdatetime_t* qdatetime_add_years(const go_qdatetime_t* dt, int years) {
    return new QDateTime(dt->addYears(years));
}
go_qdatetime_t* qdatetime_add_secs(const go_qdatetime_t* dt, int64_t secs) {
    return new QDateTime(dt->addSecs(secs));
}
go_qdatetime_t* qdatetime_add_msecs(const go_qdatetime_t* dt, int64_t msecs) {
    return new QDateTime(dt->addMSecs(msecs));
}

int64_t qdatetime_days_to(const go_qdatetime_t* dt, const go_qdatetime_t* other) {
    return dt->daysTo(*other);
}
int64_t qdatetime_secs_to(const go_qdatetime_t* dt, const go_qdatetime_t* other) {
    return dt->secsTo(*other);
}
int64_t qdatetime_msecs_to(const go_qdatetime_t* dt, const go_qdatetime_t* other) {
    return dt->msecsTo(*other);
}

go_qdatetime_t* qdatetime_current_date_time(void) {
    return new QDateTime(QDateTime::currentDateTime());
}
go_qdatetime_t* qdatetime_current_date_time_utc(void) {
    return new QDateTime(QDateTime::currentDateTimeUtc());
}
go_qdatetime_t* qdatetime_current_date_time_zone(const go_qtimezone_t tz) {
    return new QDateTime(QDateTime::currentDateTime(*tz));
}
go_qdatetime_t* qdatetime_from_msecs_since_epoch(int64_t msecs) {
    return new QDateTime(QDateTime::fromMSecsSinceEpoch(msecs));
}
go_qdatetime_t* qdatetime_from_msecs_since_epoch_zone(int64_t msecs, const go_qtimezone_t tz) {
    return new QDateTime(QDateTime::fromMSecsSinceEpoch(msecs, *tz));
}
go_qdatetime_t* qdatetime_from_secs_since_epoch(int64_t secs) {
    return new QDateTime(QDateTime::fromSecsSinceEpoch(secs));
}
go_qdatetime_t* qdatetime_from_secs_since_epoch_zone(int64_t secs, const go_qtimezone_t tz) {
    return new QDateTime(QDateTime::fromSecsSinceEpoch(secs, *tz));
}
int64_t qdatetime_current_msecs_since_epoch(void) { return QDateTime::currentMSecsSinceEpoch(); }
int64_t qdatetime_current_secs_since_epoch(void)  { return QDateTime::currentSecsSinceEpoch(); }

void qdatetime_to_string(const go_qdatetime_t* dt, int format, void** out_qstr) {
    *out_qstr = new QString(dt->toString(DF(format)));
}
void qdatetime_to_string_format(const go_qdatetime_t* dt, const char* format,
                                ptrdiff_t format_len, void** out_qstr) {
    *out_qstr = new QString(dt->toString(U8(format, format_len)));
}
void qdatetime_to_string_format_cal(const go_qdatetime_t* dt, const char* format,
                                    ptrdiff_t format_len, go_qcalendar_t cal,
                                    void** out_qstr) {
    *out_qstr = new QString(dt->toString(U8(format, format_len), cal));
}

go_qdatetime_t* qdatetime_from_string_utf8(const char* str, ptrdiff_t str_len, int format) {
    return new QDateTime(QDateTime::fromString(U8(str, str_len), DF(format)));
}
go_qdatetime_t* qdatetime_from_string_format_utf8(const char* str, ptrdiff_t str_len,
                                                  const char* format, ptrdiff_t format_len) {
    return new QDateTime(QDateTime::fromString(U8(str, str_len), U8(format, format_len)));
}
go_qdatetime_t* qdatetime_from_string_format_cal(const char* str, ptrdiff_t str_len,
                                                 const char* format, ptrdiff_t format_len,
                                                 go_qcalendar_t cal) {
    return new QDateTime(QDateTime::fromString(
        U8(str, str_len), U8(format, format_len), cal));
}
go_qdatetime_t* qdatetime_from_string_format_base_year(const char* str, ptrdiff_t str_len,
                                                       const char* format, ptrdiff_t format_len,
                                                       int base_year) {
    return new QDateTime(QDateTime::fromString(
        U8(str, str_len), U8(format, format_len), base_year));
}
go_qdatetime_t* qdatetime_from_string_format_base_year_cal(const char* str, ptrdiff_t str_len,
                                                           const char* format, ptrdiff_t format_len,
                                                           int base_year, go_qcalendar_t cal) {
    return new QDateTime(QDateTime::fromString(
        U8(str, str_len), U8(format, format_len), base_year, cal));
}

bool qdatetime_equals(const go_qdatetime_t* dt, const go_qdatetime_t* other) {
    return *dt == *other;
}

go_qdatetime_t* qdatetime_to_local_time(const go_qdatetime_t* dt) {
    return new QDateTime(dt->toLocalTime());
}
go_qdatetime_t* qdatetime_to_utc(const go_qdatetime_t* dt) {
    return new QDateTime(dt->toUTC());
}
go_qdatetime_t* qdatetime_to_offset_from_utc(const go_qdatetime_t* dt, int offset_seconds) {
    return new QDateTime(dt->toOffsetFromUtc(offset_seconds));
}

go_qtimezone_t qdatetime_time_zone(const go_qdatetime_t* dt) {
    return new QTimeZone(dt->timeZone());
}
void qdatetime_set_time_zone(go_qdatetime_t* dt, const go_qtimezone_t tz) {
    dt->setTimeZone(*tz);
}
void qdatetime_set_time_zone_resolve(go_qdatetime_t* dt, const go_qtimezone_t tz,
                                     int resolve) {
    dt->setTimeZone(*tz, TR(resolve));
}
go_qdatetime_t* qdatetime_to_time_zone(const go_qdatetime_t* dt, const go_qtimezone_t tz) {
    return new QDateTime(dt->toTimeZone(*tz));
}

} // extern "C"
