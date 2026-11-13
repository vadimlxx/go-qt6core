#ifndef GO_QT6CORE_QDATETIME_QDATETIME_C_H
#define GO_QT6CORE_QDATETIME_QDATETIME_C_H

#include <stdint.h>
#include <stdbool.h>
#include <stddef.h>
#include "qdate_c.h"
#include "qtime_c.h"
#include "qcalendar_c.h"
#include "qtimezone_fwd.h"

#ifdef __cplusplus
extern "C" {
#endif

typedef struct QDateTime go_qdatetime_t;

go_qdatetime_t* qdatetime_new_default(void);
go_qdatetime_t* qdatetime_new_from_date_time(go_qdate_t date, go_qtime_t time);
go_qdatetime_t* qdatetime_new_copy(const go_qdatetime_t* other);
go_qdatetime_t* qdatetime_clone(const go_qdatetime_t* other);
void            qdatetime_delete(go_qdatetime_t* dt);

bool qdatetime_is_null(const go_qdatetime_t* dt);
bool qdatetime_is_valid(const go_qdatetime_t* dt);

go_qdate_t qdatetime_date(const go_qdatetime_t* dt);
go_qtime_t qdatetime_time(const go_qdatetime_t* dt);
int        qdatetime_time_spec(const go_qdatetime_t* dt);
int        qdatetime_offset_from_utc(const go_qdatetime_t* dt);
int64_t    qdatetime_to_msecs_since_epoch(const go_qdatetime_t* dt);
int64_t    qdatetime_to_secs_since_epoch(const go_qdatetime_t* dt);

void qdatetime_set_date(go_qdatetime_t* dt, go_qdate_t date);
void qdatetime_set_time(go_qdatetime_t* dt, go_qtime_t time);
void qdatetime_set_msecs_since_epoch(go_qdatetime_t* dt, int64_t msecs);
void qdatetime_set_secs_since_epoch(go_qdatetime_t* dt, int64_t secs);

go_qdatetime_t* qdatetime_add_days(const go_qdatetime_t* dt, int64_t days);
go_qdatetime_t* qdatetime_add_months(const go_qdatetime_t* dt, int months);
go_qdatetime_t* qdatetime_add_years(const go_qdatetime_t* dt, int years);
go_qdatetime_t* qdatetime_add_secs(const go_qdatetime_t* dt, int64_t secs);
go_qdatetime_t* qdatetime_add_msecs(const go_qdatetime_t* dt, int64_t msecs);

int64_t qdatetime_days_to(const go_qdatetime_t* dt, const go_qdatetime_t* other);
int64_t qdatetime_secs_to(const go_qdatetime_t* dt, const go_qdatetime_t* other);
int64_t qdatetime_msecs_to(const go_qdatetime_t* dt, const go_qdatetime_t* other);

go_qdatetime_t* qdatetime_current_date_time(void);
go_qdatetime_t* qdatetime_current_date_time_utc(void);
go_qdatetime_t* qdatetime_from_msecs_since_epoch(int64_t msecs);
go_qdatetime_t* qdatetime_from_secs_since_epoch(int64_t secs);
int64_t         qdatetime_current_msecs_since_epoch(void);
int64_t         qdatetime_current_secs_since_epoch(void);

void qdatetime_to_string(const go_qdatetime_t* dt, int format, void** out_qstr);
void qdatetime_to_string_format(const go_qdatetime_t* dt, const char* format,
                                ptrdiff_t format_len, void** out_qstr);
void qdatetime_to_string_format_cal(const go_qdatetime_t* dt, const char* format,
                                    ptrdiff_t format_len, go_qcalendar_t cal,
                                    void** out_qstr);

go_qdatetime_t* qdatetime_from_string_utf8(const char* str, ptrdiff_t str_len, int format);
go_qdatetime_t* qdatetime_from_string_format_utf8(const char* str, ptrdiff_t str_len,
                                                  const char* format, ptrdiff_t format_len);

bool qdatetime_equals(const go_qdatetime_t* dt, const go_qdatetime_t* other);

go_qdatetime_t* qdatetime_to_local_time(const go_qdatetime_t* dt);
go_qdatetime_t* qdatetime_to_utc(const go_qdatetime_t* dt);
go_qdatetime_t* qdatetime_to_offset_from_utc(const go_qdatetime_t* dt, int offset_seconds);

/* --- QTimeZone integration --- */
go_qtimezone_t  qdatetime_time_zone(const go_qdatetime_t* dt);
void            qdatetime_set_time_zone(go_qdatetime_t* dt, const go_qtimezone_t tz);
go_qdatetime_t* qdatetime_to_time_zone(const go_qdatetime_t* dt, const go_qtimezone_t tz);
go_qdatetime_t* qdatetime_current_date_time_zone(const go_qtimezone_t tz);
go_qdatetime_t* qdatetime_from_msecs_since_epoch_zone(int64_t msecs, const go_qtimezone_t tz);
go_qdatetime_t* qdatetime_from_secs_since_epoch_zone(int64_t secs, const go_qtimezone_t tz);

#ifdef __cplusplus
}
#endif
#endif
