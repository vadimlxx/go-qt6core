#ifndef GO_QT6CORE_QDATETIME_QCALENDAR_C_H
#define GO_QT6CORE_QDATETIME_QCALENDAR_C_H

#include <stddef.h>
#include <stdint.h>
#include <stdbool.h>

/* FIX: qstring_c.h lives in ../qtext/, not next to us. */
#include "../qtext/qstring_c.h"

/* qdatetime_c.h defines go_qdate_t, go_qtime_t, go_qcalendar_t and go_qymd_t.
 * It does not include us, so there is no cycle. go_qlocale_t is still
 * forward-declared here (QCalendar only uses QLocale through pointers). */
#include "qdatetime_c.h"

#ifdef __cplusplus
class QLocale;
typedef QLocale*  go_qlocale_t;
extern "C" {
#else
typedef struct go_qlocale_s go_qlocale_s;
typedef go_qlocale_s* go_qlocale_t;
#endif

/* ---- Construction ---- */
go_qcalendar_t qcalendar_default(void);                    /* QCalendar() — Gregorian */
go_qcalendar_t qcalendar_from_system(int system);          /* QCalendar(System) */
go_qcalendar_t qcalendar_from_systemid(uintptr_t id);      /* QCalendar(SystemId) */
go_qcalendar_t qcalendar_from_any(const void* data, ptrdiff_t size, int tag); /* QCalendar(QAnyStringView) */

/* ---- Validity ---- */
bool qcalendar_is_valid(go_qcalendar_t self);

/* ---- Date queries ---- */
int  qcalendar_days_in_month(go_qcalendar_t self, int month, int year);
int  qcalendar_days_in_year(go_qcalendar_t self, int year);
int  qcalendar_months_in_year(go_qcalendar_t self, int year);
bool qcalendar_is_date_valid(go_qcalendar_t self, int year, int month, int day);

/* ---- Leap years ---- */
bool qcalendar_is_leap_year(go_qcalendar_t self, int year);

/* ---- Calendar properties ---- */
bool qcalendar_is_gregorian(go_qcalendar_t self);
bool qcalendar_is_lunar(go_qcalendar_t self);
bool qcalendar_is_luni_solar(go_qcalendar_t self);
bool qcalendar_is_solar(go_qcalendar_t self);
bool qcalendar_is_proleptic(go_qcalendar_t self);
bool qcalendar_has_year_zero(go_qcalendar_t self);
int  qcalendar_maximum_days_in_month(go_qcalendar_t self);
int  qcalendar_minimum_days_in_month(go_qcalendar_t self);
int  qcalendar_maximum_months_in_year(go_qcalendar_t self);
void qcalendar_name(go_qcalendar_t self, void** out_qstr);

/* ---- QDate conversions ---- */
go_qdate_t qcalendar_date_from_parts(go_qcalendar_t self, int year, int month, int day);
go_qdate_t qcalendar_date_from_ymd(go_qcalendar_t self, go_qymd_t ymd);
go_qdate_t qcalendar_match_century_to_weekday(go_qcalendar_t self, go_qymd_t ymd, int dow);
go_qymd_t  qcalendar_parts_from_date(go_qcalendar_t self, go_qdate_t date);
int        qcalendar_day_of_week(go_qcalendar_t self, go_qdate_t date);

/* ---- Month / weekday names ---- */
void qcalendar_month_name(go_qcalendar_t self, const go_qlocale_t locale,
                          int month, int year, int format, void** out_qstr);
void qcalendar_standalone_month_name(go_qcalendar_t self, const go_qlocale_t locale,
                                     int month, int year, int format, void** out_qstr);
void qcalendar_week_day_name(go_qcalendar_t self, const go_qlocale_t locale,
                             int day, int format, void** out_qstr);
void qcalendar_standalone_week_day_name(go_qcalendar_t self, const go_qlocale_t locale,
                                        int day, int format, void** out_qstr);

/* ---- Date-time formatting ---- */
void qcalendar_date_time_to_string(go_qcalendar_t self, const char* format, ptrdiff_t format_len,
                                   const go_qdatetime_t* dt, go_qdate_t date_only,
                                   go_qtime_t time_only, const go_qlocale_t locale,
                                   void** out_qstr);

/* ---- Available calendars ---- */
QString** qcalendar_available_calendars(ptrdiff_t* out_count);
void      qcalendar_free_string_array(QString** arr);

#ifdef __cplusplus
}
#endif
#endif /* GO_QT6CORE_QDATETIME_QCALENDAR_C_H */