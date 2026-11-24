/* C ABI for Qt 6.11 QtCore/QDate, QTime and QDateTime.
 *
 * QtCore declares all three classes in qdatetime.h and implements them in
 * qdatetime.cpp, so this bridge follows the same layout: one header and one
 * implementation file for the whole trio (qdate_c.* / qtime_c.* are gone).
 *
 * Type model:
 *   - go_qdate_t      — int64_t Julian day (value type, no heap);
 *   - go_qtime_t      — int32_t msecs since start of day, -1 = null (value type);
 *   - go_qdatetime_t  — QDateTime* (heap handle, owned by the caller);
 *   - go_qcalendar_t  — QCalendar value (needed by the QDate calendar variants);
 *   - go_qymd_t       — QCalendar::YearMonthDay value;
 *   - go_qtimezone_t  — QTimeZone* (declared in qtimezone_c.h).
 *
 * qtimezone_c.h is included here because QDate/QDateTime expose QTimeZone
 * parameters (startOfDay/endOfDay, setTimeZone, toTimeZone). It forward-declares
 * go_qtimezone_t/go_qdatetime_t and does not include this header, so there is no
 * include cycle.
 */
#ifndef GO_QT6CORE_QDATETIME_QDATETIME_C_H
#define GO_QT6CORE_QDATETIME_QDATETIME_C_H

#include <stdint.h>
#include <stdbool.h>
#include <stddef.h>

#include "../qtext/qstring_c.h"   /* QString* (opaque in C mode) */
#include "qtimezone_c.h"           /* go_qtimezone_t, go_qdatetime_t */

/* ---- Value types ---- */

/* go_qdate_t: Julian day as int64_t.
 *   QDate{}   → Julian day 0        → valid date
 *   null      → Julian day INT64_MIN (qdate_new_null) */
typedef int64_t go_qdate_t;

/* go_qtime_t: -1 = null QTime, 0..86399999 = valid msecs since start of day. */
typedef int32_t go_qtime_t;

/* go_qcalendar_t: QCalendar by value.
 * In C++ mode this pulls in <QCalendar>; in C mode it is an opaque struct with
 * the same memory layout (single const void*). */
#ifdef __cplusplus
#include <QCalendar>
typedef QCalendar go_qcalendar_t;
#else
typedef struct go_qcalendar_s { const void* d_ptr; } go_qcalendar_t;
#endif

/* go_qymd_t mirrors QCalendar::YearMonthDay. */
typedef struct go_qymd_s {
    int year;
    int month;
    int day;
} go_qymd_t;

#ifdef __cplusplus
extern "C" {
#endif

/* ==== QDate ==== */

/* Lifecycle (value type — no heap, no Close).
 * Qt: QDate(), QDate(int y, int m, int d), QDate(int y, int m, int d, QCalendar),
 *     QDate::fromJulianDay(qint64), QDate::currentDate(). */
go_qdate_t qdate_new_null(void);
go_qdate_t qdate_new_from_ymd(int y, int m, int d);
go_qdate_t qdate_from_ymd_cal(int y, int m, int d, go_qcalendar_t cal);
go_qdate_t qdate_from_julian_day(int64_t jd);
go_qdate_t qdate_current_date(void);

/* Accessors.
 * Qt: QDate::toJulianDay/isNull/isValid/year/month/day/dayOfWeek/dayOfYear/
 *     daysInMonth/daysInYear/weekNumber. */
int64_t qdate_to_julian_day(go_qdate_t d);
bool    qdate_is_null(go_qdate_t d);
bool    qdate_is_valid(go_qdate_t d);
int     qdate_year(go_qdate_t d);
int     qdate_month(go_qdate_t d);
int     qdate_day(go_qdate_t d);
int     qdate_day_of_week(go_qdate_t d);
int     qdate_day_of_year(go_qdate_t d);
int     qdate_days_in_month(go_qdate_t d);
int     qdate_days_in_year(go_qdate_t d);
int     qdate_week_number(go_qdate_t d, int* year_num);

/* Calendar-aware accessors.
 * Qt: QDate::year(QCalendar) / month(QCalendar) / day(QCalendar). */
int qdate_year_cal(go_qdate_t d, go_qcalendar_t cal);
int qdate_month_cal(go_qdate_t d, go_qcalendar_t cal);
int qdate_day_cal(go_qdate_t d, go_qcalendar_t cal);

/* Arithmetic.
 * Qt: QDate::addDays/addMonths/addYears/addMonths(cal)/addYears(cal)/daysTo. */
go_qdate_t qdate_add_days(go_qdate_t d, int64_t days);
go_qdate_t qdate_add_months(go_qdate_t d, int months);
go_qdate_t qdate_add_years(go_qdate_t d, int years);
go_qdate_t qdate_add_months_cal(go_qdate_t d, int months, go_qcalendar_t cal);
go_qdate_t qdate_add_years_cal(go_qdate_t d, int years, go_qcalendar_t cal);
int64_t    qdate_days_to(go_qdate_t d, go_qdate_t other);

/* Mutators.
 * Qt: QDate::setDate/setDate(cal)/getDate. */
bool qdate_set_date(go_qdate_t* d, int y, int m, int day);
bool qdate_set_date_cal(go_qdate_t* d, int y, int m, int day, go_qcalendar_t cal);
void qdate_get_date(go_qdate_t d, int* y, int* m, int* day);

/* String conversion.
 * Qt: QDate::toString(Qt::DateFormat), toString(QString), toString(QString, QCalendar). */
void qdate_to_string(go_qdate_t d, int format, void** out_qstr);
void qdate_to_string_format(go_qdate_t d, const char* fmt, ptrdiff_t fmt_len, void** out_qstr);
void qdate_to_string_format_cal(go_qdate_t d, const char* fmt, ptrdiff_t fmt_len,
                                go_qcalendar_t cal, void** out_qstr);

/* Parsing (returns a null QDate on failure).
 * Qt: QDate::fromString(QString, Qt::DateFormat), fromString(QString, QString),
 *     fromString(QString, QString, QCalendar),
 *     fromString(QString, QString, int baseYear),
 *     fromString(QString, QString, int baseYear, QCalendar). */
go_qdate_t qdate_from_string(const char* str, ptrdiff_t str_len, int format);
go_qdate_t qdate_from_string_format(const char* str, ptrdiff_t str_len,
                                    const char* fmt, ptrdiff_t fmt_len);
go_qdate_t qdate_from_string_format_cal(const char* str, ptrdiff_t str_len,
                                        const char* fmt, ptrdiff_t fmt_len, go_qcalendar_t cal);
go_qdate_t qdate_from_string_format_base_year(const char* str, ptrdiff_t str_len,
                                              const char* fmt, ptrdiff_t fmt_len, int base_year);
go_qdate_t qdate_from_string_format_base_year_cal(const char* str, ptrdiff_t str_len,
                                                  const char* fmt, ptrdiff_t fmt_len,
                                                  int base_year, go_qcalendar_t cal);

/* Static helpers.
 * Qt: QDate::isValid(y,m,d), QDate::isLeapYear(year). */
bool qdate_static_is_valid(int y, int m, int d);
bool qdate_static_is_leap_year(int year);

/* Day boundaries (new in Qt 6; require QDateTime/QTimeZone).
 * Qt: QDate::startOfDay(), endOfDay(),
 *     startOfDay(Qt::TimeSpec, int offsetSeconds), endOfDay(Qt::TimeSpec, int),
 *     startOfDay(const QTimeZone &), endOfDay(const QTimeZone &). */
go_qdatetime_t* qdate_start_of_day(go_qdate_t date);
go_qdatetime_t* qdate_end_of_day(go_qdate_t date);
go_qdatetime_t* qdate_start_of_day_spec(go_qdate_t date, int spec, int offset_seconds);
go_qdatetime_t* qdate_end_of_day_spec(go_qdate_t date, int spec, int offset_seconds);
go_qdatetime_t* qdate_start_of_day_zone(go_qdate_t date, const go_qtimezone_t tz);
go_qdatetime_t* qdate_end_of_day_zone(go_qdate_t date, const go_qtimezone_t tz);

/* ==== QTime ==== */

/* Lifecycle (value type — no heap).
 * Qt: QTime(), QTime(int h, int m, int s, int ms), QTime::currentTime(),
 *     QTime::fromMSecsSinceStartOfDay(int). */
go_qtime_t qtime_default(void);
go_qtime_t qtime_from_hms(int h, int m, int s, int ms);
go_qtime_t qtime_current_time(void);
go_qtime_t qtime_from_msecs_since_start_of_day(int msecs);

/* Parsing.
 * Qt: QTime::fromString(QString, Qt::DateFormat), fromString(QString, QString). */
go_qtime_t qtime_from_string_utf8(const char* str, ptrdiff_t str_len, int format);
go_qtime_t qtime_from_string_format_utf8(const char* str, ptrdiff_t str_len,
                                         const char* format, ptrdiff_t format_len);

/* Accessors.
 * Qt: QTime::isNull/isValid/hour/minute/second/msec/msecsSinceStartOfDay. */
bool qtime_is_null(go_qtime_t t);
bool qtime_is_valid(go_qtime_t t);
int  qtime_hour(go_qtime_t t);
int  qtime_minute(go_qtime_t t);
int  qtime_second(go_qtime_t t);
int  qtime_msec(go_qtime_t t);
int  qtime_msecs_since_start_of_day(go_qtime_t t);

/* Arithmetic.
 * Qt: QTime::addSecs/addMSecs/secsTo/msecsTo. */
go_qtime_t qtime_add_secs(go_qtime_t t, int secs);
go_qtime_t qtime_add_msecs(go_qtime_t t, int msecs);
int        qtime_secs_to(go_qtime_t t, go_qtime_t other);
int        qtime_msecs_to(go_qtime_t t, go_qtime_t other);

/* String conversion.
 * Qt: QTime::toString(Qt::DateFormat), toString(QString). */
void qtime_to_string(go_qtime_t t, int format, void** out_qstr);
void qtime_to_string_format(go_qtime_t t, const char* format, ptrdiff_t format_len,
                            void** out_qstr);

/* Mutators / static.
 * Qt: QTime::setHMS, QTime::isValid(h,m,s,ms). */
bool qtime_set_hms(go_qtime_t* t, int h, int m, int s, int ms);
bool qtime_static_is_valid(int h, int m, int s, int ms);

/* ==== QDateTime ==== */

/* Construction and copying.
 * Qt: QDateTime(), QDateTime(QDate, QTime),
 *     QDateTime(QDate, QTime, TransitionResolution),
 *     QDateTime(QDate, QTime, QTimeZone, TransitionResolution),
 *     QDateTime(const QDateTime &), QDateTime::swap. */
go_qdatetime_t* qdatetime_new_default(void);
go_qdatetime_t* qdatetime_new_from_date_time(go_qdate_t date, go_qtime_t time);
go_qdatetime_t* qdatetime_new_from_date_time_resolve(go_qdate_t date, go_qtime_t time,
                                                     int resolve);
go_qdatetime_t* qdatetime_new_from_date_time_zone(go_qdate_t date, go_qtime_t time,
                                                  const go_qtimezone_t tz);
go_qdatetime_t* qdatetime_new_from_date_time_zone_resolve(go_qdate_t date, go_qtime_t time,
                                                          const go_qtimezone_t tz,
                                                          int resolve);
go_qdatetime_t* qdatetime_new_copy(const go_qdatetime_t* other);
go_qdatetime_t* qdatetime_clone(const go_qdatetime_t* other);
void            qdatetime_delete(go_qdatetime_t* dt);
void            qdatetime_swap(go_qdatetime_t* self, go_qdatetime_t* other);

/* Status.
 * Qt: QDateTime::isNull/isValid/isDaylightTime. */
bool qdatetime_is_null(const go_qdatetime_t* dt);
bool qdatetime_is_valid(const go_qdatetime_t* dt);
bool qdatetime_is_daylight_time(const go_qdatetime_t* dt);

/* Accessors.
 * Qt: QDateTime::date/time/timeSpec/timeRepresentation/offsetFromUtc/
 *     toMSecsSinceEpoch/toSecsSinceEpoch/timeZoneAbbreviation. */
go_qdate_t qdatetime_date(const go_qdatetime_t* dt);
go_qtime_t qdatetime_time(const go_qdatetime_t* dt);
int        qdatetime_time_spec(const go_qdatetime_t* dt);
go_qtimezone_t qdatetime_time_representation(const go_qdatetime_t* dt);
int        qdatetime_offset_from_utc(const go_qdatetime_t* dt);
int64_t    qdatetime_to_msecs_since_epoch(const go_qdatetime_t* dt);
int64_t    qdatetime_to_secs_since_epoch(const go_qdatetime_t* dt);
void       qdatetime_time_zone_abbreviation(const go_qdatetime_t* dt, void** out_qstr);

/* Setters.
 * Qt: QDateTime::setDate(date, TransitionResolution),
 *     setTime(time, TransitionResolution), setTimeZone(zone, TransitionResolution),
 *     setMSecsSinceEpoch, setSecsSinceEpoch. */
void qdatetime_set_date(go_qdatetime_t* dt, go_qdate_t date);
void qdatetime_set_date_resolve(go_qdatetime_t* dt, go_qdate_t date, int resolve);
void qdatetime_set_time(go_qdatetime_t* dt, go_qtime_t time);
void qdatetime_set_time_resolve(go_qdatetime_t* dt, go_qtime_t time, int resolve);
void qdatetime_set_msecs_since_epoch(go_qdatetime_t* dt, int64_t msecs);
void qdatetime_set_secs_since_epoch(go_qdatetime_t* dt, int64_t secs);

/* Arithmetic.
 * Qt: QDateTime::addDays/addMonths/addYears/addSecs/addMSecs. */
go_qdatetime_t* qdatetime_add_days(const go_qdatetime_t* dt, int64_t days);
go_qdatetime_t* qdatetime_add_months(const go_qdatetime_t* dt, int months);
go_qdatetime_t* qdatetime_add_years(const go_qdatetime_t* dt, int years);
go_qdatetime_t* qdatetime_add_secs(const go_qdatetime_t* dt, int64_t secs);
go_qdatetime_t* qdatetime_add_msecs(const go_qdatetime_t* dt, int64_t msecs);

/* Distances.
 * Qt: QDateTime::daysTo/secsTo/msecsTo. */
int64_t qdatetime_days_to(const go_qdatetime_t* dt, const go_qdatetime_t* other);
int64_t qdatetime_secs_to(const go_qdatetime_t* dt, const go_qdatetime_t* other);
int64_t qdatetime_msecs_to(const go_qdatetime_t* dt, const go_qdatetime_t* other);

/* Static constructors.
 * Qt: QDateTime::currentDateTime(), currentDateTimeUtc(), currentDateTime(zone),
 *     fromMSecsSinceEpoch([zone]), fromSecsSinceEpoch([zone]),
 *     currentMSecsSinceEpoch(), currentSecsSinceEpoch(). */
go_qdatetime_t* qdatetime_current_date_time(void);
go_qdatetime_t* qdatetime_current_date_time_utc(void);
go_qdatetime_t* qdatetime_current_date_time_zone(const go_qtimezone_t tz);
go_qdatetime_t* qdatetime_from_msecs_since_epoch(int64_t msecs);
go_qdatetime_t* qdatetime_from_msecs_since_epoch_zone(int64_t msecs, const go_qtimezone_t tz);
go_qdatetime_t* qdatetime_from_secs_since_epoch(int64_t secs);
go_qdatetime_t* qdatetime_from_secs_since_epoch_zone(int64_t secs, const go_qtimezone_t tz);
int64_t         qdatetime_current_msecs_since_epoch(void);
int64_t         qdatetime_current_secs_since_epoch(void);

/* String conversion.
 * Qt: QDateTime::toString(Qt::DateFormat), toString(QString),
 *     toString(QString, QCalendar). */
void qdatetime_to_string(const go_qdatetime_t* dt, int format, void** out_qstr);
void qdatetime_to_string_format(const go_qdatetime_t* dt, const char* format,
                                ptrdiff_t format_len, void** out_qstr);
void qdatetime_to_string_format_cal(const go_qdatetime_t* dt, const char* format,
                                    ptrdiff_t format_len, go_qcalendar_t cal,
                                    void** out_qstr);

/* Parsing.
 * Qt: QDateTime::fromString(QString, Qt::DateFormat),
 *     fromString(QString, QString[, int baseYear[, QCalendar]]). */
go_qdatetime_t* qdatetime_from_string_utf8(const char* str, ptrdiff_t str_len, int format);
go_qdatetime_t* qdatetime_from_string_format_utf8(const char* str, ptrdiff_t str_len,
                                                  const char* format, ptrdiff_t format_len);
go_qdatetime_t* qdatetime_from_string_format_cal(const char* str, ptrdiff_t str_len,
                                                 const char* format, ptrdiff_t format_len,
                                                 go_qcalendar_t cal);
go_qdatetime_t* qdatetime_from_string_format_base_year(const char* str, ptrdiff_t str_len,
                                                       const char* format, ptrdiff_t format_len,
                                                       int base_year);
go_qdatetime_t* qdatetime_from_string_format_base_year_cal(const char* str, ptrdiff_t str_len,
                                                           const char* format, ptrdiff_t format_len,
                                                           int base_year, go_qcalendar_t cal);

/* Comparison.
 * Qt: operator==(const QDateTime &, const QDateTime &). */
bool qdatetime_equals(const go_qdatetime_t* dt, const go_qdatetime_t* other);

/* Time zone conversions.
 * Qt: QDateTime::toLocalTime/toUTC/toOffsetFromUtc/toTimeZone,
 *     timeZone(), setTimeZone(), setTimeZone(zone, TransitionResolution). */
go_qdatetime_t* qdatetime_to_local_time(const go_qdatetime_t* dt);
go_qdatetime_t* qdatetime_to_utc(const go_qdatetime_t* dt);
go_qdatetime_t* qdatetime_to_offset_from_utc(const go_qdatetime_t* dt, int offset_seconds);
go_qtimezone_t  qdatetime_time_zone(const go_qdatetime_t* dt);
void            qdatetime_set_time_zone(go_qdatetime_t* dt, const go_qtimezone_t tz);
void            qdatetime_set_time_zone_resolve(go_qdatetime_t* dt, const go_qtimezone_t tz,
                                                int resolve);
go_qdatetime_t* qdatetime_to_time_zone(const go_qdatetime_t* dt, const go_qtimezone_t tz);

#ifdef __cplusplus
}
#endif

#endif /* GO_QT6CORE_QDATETIME_QDATETIME_C_H */
