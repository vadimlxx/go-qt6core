#ifndef GO_QT6CORE_QDATETIME_QDATE_C_H
#define GO_QT6CORE_QDATETIME_QDATE_C_H

#include <stdint.h>
#include <stdbool.h>
#include <stddef.h>
#include "../qtext/qstring_c.h"   /* QString* (opaque in C mode) */

/* go_qdate_t: a value type (Julian day as int64_t).
 *
 *   QDate{}     in Go  → Julian day 0      → valid date
 *   QDate null  in Go  → Julian day INT64_MIN (returned by qdate_new_null)
 */
typedef int64_t go_qdate_t;

/* go_qcalendar_t: a VALUE type (single pointer, 8 bytes LP64).
 *
 * Defined HERE — not in qcalendar_c.h — to break the include cycle:
 * qdate_c.h uses go_qcalendar_t in its _cal variants, and qcalendar_c.h
 * returns go_qdate_t from its dateFromParts / partsFromDate methods.
 * Since both are value types (passed by value), neither can be forward-
 * declared; one must be fully defined before the other is seen. qdate
 * is the foundational type (QCalendar returns QDate), so it owns the
 * go_qcalendar_t definition. qcalendar_c.h gets it via #include "qdate_c.h".
 *
 * In C++ mode this pulls in <QCalendar>; in C mode it defines an opaque
 * struct with the same memory layout as QCalendar (single const void*).
 */
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

/* ---- Lifecycle (value type — no heap, no Close) ---- */
go_qdate_t qdate_new_null(void);                 /* QDate() — null, jd = INT64_MIN */
go_qdate_t qdate_new_from_ymd(int y, int m, int d);          /* QDate(y,m,d) */
go_qdate_t qdate_from_ymd_cal(int y, int m, int d, go_qcalendar_t cal); /* QDate(y,m,d,cal) */
go_qdate_t qdate_from_julian_day(int64_t jd);   /* QDate::fromJulianDay */
go_qdate_t qdate_current_date(void);            /* QDate::currentDate */

/* ---- Accessors ---- */
int64_t qdate_to_julian_day(go_qdate_t d);      /* QDate::toJulianDay */
bool    qdate_is_null(go_qdate_t d);            /* QDate::isNull */
bool    qdate_is_valid(go_qdate_t d);           /* QDate::isValid */

int qdate_year(go_qdate_t d);                   /* QDate::year() */
int qdate_month(go_qdate_t d);                  /* QDate::month() */
int qdate_day(go_qdate_t d);                    /* QDate::day() */
int qdate_day_of_week(go_qdate_t d);            /* QDate::dayOfWeek() */
int qdate_day_of_year(go_qdate_t d);            /* QDate::dayOfYear() */
int qdate_days_in_month(go_qdate_t d);          /* QDate::daysInMonth() */
int qdate_days_in_year(go_qdate_t d);           /* QDate::daysInYear() */
int qdate_week_number(go_qdate_t d, int* year_num); /* QDate::weekNumber(int*) */

/* ---- Calendar-aware accessors ---- */
int qdate_year_cal(go_qdate_t d, go_qcalendar_t cal);   /* QDate::year(QCalendar) */
int qdate_month_cal(go_qdate_t d, go_qcalendar_t cal);  /* QDate::month(QCalendar) */
int qdate_day_cal(go_qdate_t d, go_qcalendar_t cal);    /* QDate::day(QCalendar) */

/* ---- Arithmetic ---- */
go_qdate_t qdate_add_days(go_qdate_t d, int64_t days);    /* QDate::addDays */
go_qdate_t qdate_add_months(go_qdate_t d, int months);    /* QDate::addMonths */
go_qdate_t qdate_add_years(go_qdate_t d, int years);      /* QDate::addYears */
go_qdate_t qdate_add_months_cal(go_qdate_t d, int months, go_qcalendar_t cal);
go_qdate_t qdate_add_years_cal(go_qdate_t d, int years, go_qcalendar_t cal);
int64_t    qdate_days_to(go_qdate_t d, go_qdate_t other); /* QDate::daysTo */

/* ---- Mutators ---- */
bool qdate_set_date(go_qdate_t* d, int y, int m, int day);                    /* QDate::setDate */
bool qdate_set_date_cal(go_qdate_t* d, int y, int m, int day, go_qcalendar_t cal);
void qdate_get_date(go_qdate_t d, int* y, int* m, int* day);                  /* QDate::getDate */

/* ---- String conversion ---- */
void qdate_to_string(go_qdate_t d, int format, void** out_qstr);              /* QDate::toString(Qt::DateFormat) */
void qdate_to_string_format(go_qdate_t d, const char* fmt, ptrdiff_t fmt_len, void** out_qstr);
void qdate_to_string_format_cal(go_qdate_t d, const char* fmt, ptrdiff_t fmt_len,
                                go_qcalendar_t cal, void** out_qstr);

/* QDate::fromString — returns null QDate on parse failure (check qdate_is_valid) */
go_qdate_t qdate_from_string(const char* str, ptrdiff_t str_len, int format);
go_qdate_t qdate_from_string_format(const char* str, ptrdiff_t str_len,
                                    const char* fmt, ptrdiff_t fmt_len);
go_qdate_t qdate_from_string_format_cal(const char* str, ptrdiff_t str_len,
                                        const char* fmt, ptrdiff_t fmt_len, go_qcalendar_t cal);

/* ---- Static ---- */
bool qdate_static_is_valid(int y, int m, int d);   /* QDate::isValid(y,m,d) */
bool qdate_static_is_leap_year(int year);          /* QDate::isLeapYear */

#ifdef __cplusplus
}
#endif
#endif /* GO_QT6CORE_QDATETIME_QDATE_C_H */