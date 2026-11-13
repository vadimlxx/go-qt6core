#ifndef GO_QT6CORE_QDATETIME_QTIME_C_H
#define GO_QT6CORE_QDATETIME_QTIME_C_H

#include <stdint.h>
#include <stdbool.h>
#include <stddef.h>

#ifdef __cplusplus
extern "C" {
#endif

/* go_qtime_t: -1 = null QTime, 0..86399999 = valid msecs since start of day. */
typedef int32_t go_qtime_t;

go_qtime_t qtime_default(void);
go_qtime_t qtime_from_hms(int h, int m, int s, int ms);
go_qtime_t qtime_current_time(void);
go_qtime_t qtime_from_string_utf8(const char* str, ptrdiff_t str_len, int format);
go_qtime_t qtime_from_string_format_utf8(const char* str, ptrdiff_t str_len,
                                         const char* format, ptrdiff_t format_len);
go_qtime_t qtime_from_msecs_since_start_of_day(int msecs);

bool qtime_is_null(go_qtime_t t);
bool qtime_is_valid(go_qtime_t t);
int  qtime_hour(go_qtime_t t);
int  qtime_minute(go_qtime_t t);
int  qtime_second(go_qtime_t t);
int  qtime_msec(go_qtime_t t);
int  qtime_msecs_since_start_of_day(go_qtime_t t);

go_qtime_t qtime_add_secs(go_qtime_t t, int secs);
go_qtime_t qtime_add_msecs(go_qtime_t t, int msecs);
int        qtime_secs_to(go_qtime_t t, go_qtime_t other);
int        qtime_msecs_to(go_qtime_t t, go_qtime_t other);

void qtime_to_string(go_qtime_t t, int format, void** out_qstr);
void qtime_to_string_format(go_qtime_t t, const char* format, ptrdiff_t format_len,
                            void** out_qstr);

bool qtime_set_hms(go_qtime_t* t, int h, int m, int s, int ms);

bool qtime_static_is_valid(int h, int m, int s, int ms);

#ifdef __cplusplus
}
#endif
#endif