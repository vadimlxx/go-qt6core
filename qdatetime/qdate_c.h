#ifndef GO_QT6CORE_QDATETIME_QDATE_C_H
#define GO_QT6CORE_QDATETIME_QDATE_C_H

#include <stdint.h>
#include <stdbool.h>
#include <stddef.h>
#include "qcalendar_c.h"

#ifdef __cplusplus
extern "C" {
#endif

typedef int64_t go_qdate_t;

go_qdate_t qdate_default(void);
go_qdate_t qdate_from_ymd(int y, int m, int d);
go_qdate_t qdate_from_ymd_cal(int y, int m, int d, go_qcalendar_t cal);
go_qdate_t qdate_from_julian_day(int64_t jd);
go_qdate_t qdate_current_date(void);

go_qdate_t qdate_from_string_utf8(const char* str, ptrdiff_t str_len, int format);
go_qdate_t qdate_from_string_format_utf8(const char* str, ptrdiff_t str_len,
                                         const char* format, ptrdiff_t format_len);

bool qdate_is_null(go_qdate_t d);
bool qdate_is_valid(go_qdate_t d);
int  qdate_year(go_qdate_t d);
int  qdate_month(go_qdate_t d);
int  qdate_day(go_qdate_t d);
int  qdate_day_of_week(go_qdate_t d);
int  qdate_day_of_year(go_qdate_t d);
int  qdate_days_in_month(go_qdate_t d);
int  qdate_days_in_year(go_qdate_t d);
int  qdate_week_number(go_qdate_t d, int* year_num);

int  qdate_year_cal(go_qdate_t d, go_qcalendar_t cal);
int  qdate_month_cal(go_qdate_t d, go_qcalendar_t cal);
int  qdate_day_cal(go_qdate_t d, go_qcalendar_t cal);

go_qdate_t qdate_add_days(go_qdate_t d, int64_t days);
go_qdate_t qdate_add_months(go_qdate_t d, int months);
go_qdate_t qdate_add_years(go_qdate_t d, int years);
go_qdate_t qdate_add_months_cal(go_qdate_t d, int months, go_qcalendar_t cal);
go_qdate_t qdate_add_years_cal(go_qdate_t d, int years, go_qcalendar_t cal);
int64_t    qdate_days_to(go_qdate_t d, go_qdate_t other);

void qdate_to_string(go_qdate_t d, int format, void** out_qstr);
void qdate_to_string_format(go_qdate_t d, const char* format, ptrdiff_t format_len,
                            void** out_qstr);
void qdate_to_string_format_cal(go_qdate_t d, const char* format, ptrdiff_t format_len,
                                go_qcalendar_t cal, void** out_qstr);

bool qdate_set_date(go_qdate_t* d, int y, int m, int day);
bool qdate_set_date_cal(go_qdate_t* d, int y, int m, int day, go_qcalendar_t cal);
void qdate_get_date(go_qdate_t d, int* y, int* m, int* day);

int64_t qdate_to_julian_day(go_qdate_t d);

bool qdate_static_is_valid(int y, int m, int d);
bool qdate_static_is_leap_year(int year);

#ifdef __cplusplus
}
#endif
#endif