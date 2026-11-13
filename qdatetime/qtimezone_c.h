#ifndef GO_QT6CORE_QDATETIME_QTIMEZONE_C_H
#define GO_QT6CORE_QDATETIME_QTIMEZONE_C_H

#include <stdint.h>
#include <stdbool.h>
#include <stddef.h>
#include "qtimezone_fwd.h"
#include "qdatetime_c.h"


#ifdef __cplusplus
class QTimeZone;
typedef QTimeZone* go_qtimezone_t;
extern "C" {
#else
typedef struct go_qtimezone_s go_qtimezone_s;
typedef go_qtimezone_s* go_qtimezone_t;
#endif


/* --- Lifecycle --- */
go_qtimezone_t qtimezone_new_default(void);
go_qtimezone_t qtimezone_new_local(void);
go_qtimezone_t qtimezone_new_utc(void);
go_qtimezone_t qtimezone_new_from_offset(int offset_seconds);
go_qtimezone_t qtimezone_new_from_iana(const char* id, ptrdiff_t id_len);
go_qtimezone_t qtimezone_new_copy(const go_qtimezone_t other);
go_qtimezone_t qtimezone_clone(const go_qtimezone_t self);
void           qtimezone_delete(go_qtimezone_t self);

/* --- Validity / identification --- */
bool qtimezone_is_valid(const go_qtimezone_t self);
void qtimezone_id(const go_qtimezone_t self, void** out_qba);

/* --- TimeSpec / offset --- */
int qtimezone_time_spec(const go_qtimezone_t self);
int qtimezone_fixed_seconds_ahead_of_utc(const go_qtimezone_t self);
bool qtimezone_is_utc_or_fixed_offset(const go_qtimezone_t self);

/* --- Offset queries (take QDateTime*) --- */
int qtimezone_offset_from_utc(const go_qtimezone_t self, const go_qdatetime_t* at);
int qtimezone_standard_time_offset(const go_qtimezone_t self, const go_qdatetime_t* at);
int qtimezone_daylight_time_offset(const go_qtimezone_t self, const go_qdatetime_t* at);
bool qtimezone_is_daylight_time(const go_qtimezone_t self, const go_qdatetime_t* at);
bool qtimezone_has_daylight_time(const go_qtimezone_t self);
bool qtimezone_has_transitions(const go_qtimezone_t self);

/* --- Display names --- */
void qtimezone_display_name_at(const go_qtimezone_t self, const go_qdatetime_t* at,
                               int name_type, const void* locale, void** out_qstr);
void qtimezone_display_name_type(const go_qtimezone_t self, int time_type, int name_type,
                                 const void* locale, void** out_qstr);
void qtimezone_abbreviation(const go_qtimezone_t self, const go_qdatetime_t* at, void** out_qstr);

/* --- Transitions (return new QDateTime* via out-param) --- */
void qtimezone_next_transition(const go_qtimezone_t self, const go_qdatetime_t* after,
                               int64_t* out_offset, int64_t* out_std_offset,
                               int64_t* out_dst_offset, void** out_at_utc, void** out_abbr);
void qtimezone_previous_transition(const go_qtimezone_t self, const go_qdatetime_t* before,
                                   int64_t* out_offset, int64_t* out_std_offset,
                                   int64_t* out_dst_offset, void** out_at_utc, void** out_abbr);

/* --- Static --- */
void qtimezone_system_id(void** out_qba);
go_qtimezone_t qtimezone_system_zone(void);
go_qtimezone_t qtimezone_static_utc(void);
bool qtimezone_is_time_zone_id_available(const char* id, ptrdiff_t id_len);

/* --- Conversions --- */
go_qtimezone_t qtimezone_as_backend_zone(const go_qtimezone_t self);
void qtimezone_territory(const go_qtimezone_t self, int* out_territory);

#ifdef __cplusplus
}
#endif
#endif
