#include "../internal/qt6core/check_qt_version.h"

#include <QTimeZone>
#include <QByteArray>
#include <QString>
#include <QDateTime>
#include <QLocale>

#include "qtimezone_c.h"
#include "qdatetime_c.h"

namespace {
//inline QTimeZone*       Z(go_qtimezone_t p)       { return reinterpret_cast<QTimeZone*>(p); }
inline const QTimeZone* Z(const go_qtimezone_t p) { return reinterpret_cast<const QTimeZone*>(p); }
inline go_qtimezone_t   toGo(QTimeZone* p)        { return reinterpret_cast<go_qtimezone_t>(p); }

/* Go strings are NOT null-terminated. Use (data, len). */
inline QByteArray BA(const char* d, ptrdiff_t n) {
    return QByteArray(d, static_cast<qsizetype>(n));
}

inline const QDateTime* DT(const go_qdatetime_t* p) {
    return reinterpret_cast<const QDateTime*>(p);
}
}

extern "C" {

/* --- Lifecycle --- */
go_qtimezone_t qtimezone_new_default(void) { return toGo(new QTimeZone()); }
go_qtimezone_t qtimezone_new_local(void)   { return toGo(new QTimeZone(QTimeZone::LocalTime)); }
go_qtimezone_t qtimezone_new_utc(void)     { return toGo(new QTimeZone(QTimeZone::UTC)); }
go_qtimezone_t qtimezone_new_from_offset(int offset_seconds) {
	return toGo(new QTimeZone(QTimeZone::fromSecondsAheadOfUtc(offset_seconds)));
}
go_qtimezone_t qtimezone_new_from_iana(const char* id, ptrdiff_t id_len) {
    return toGo(new QTimeZone(BA(id, id_len)));
}
go_qtimezone_t qtimezone_new_copy(const go_qtimezone_t other) { return toGo(new QTimeZone(*Z(other))); }
go_qtimezone_t qtimezone_clone(const go_qtimezone_t self)     { return toGo(new QTimeZone(*Z(self))); }
void           qtimezone_delete(go_qtimezone_t self)          { delete Z(self); }

/* --- Validity / identification --- */
bool qtimezone_is_valid(const go_qtimezone_t self) { return Z(self)->isValid(); }
void qtimezone_id(const go_qtimezone_t self, void** out_qba) {
    *out_qba = new QByteArray(Z(self)->id());
}

/* --- TimeSpec / offset --- */
int qtimezone_time_spec(const go_qtimezone_t self) {
    return static_cast<int>(Z(self)->timeSpec());
}
int qtimezone_fixed_seconds_ahead_of_utc(const go_qtimezone_t self) {
    return Z(self)->fixedSecondsAheadOfUtc();
}
bool qtimezone_is_utc_or_fixed_offset(const go_qtimezone_t self) {
    return Z(self)->isUtcOrFixedOffset();
}

/* --- Offset queries --- */
int qtimezone_offset_from_utc(const go_qtimezone_t self, const go_qdatetime_t* at) {
    return Z(self)->offsetFromUtc(*DT(at));
}
int qtimezone_standard_time_offset(const go_qtimezone_t self, const go_qdatetime_t* at) {
    return Z(self)->standardTimeOffset(*DT(at));
}
int qtimezone_daylight_time_offset(const go_qtimezone_t self, const go_qdatetime_t* at) {
    return Z(self)->daylightTimeOffset(*DT(at));
}
bool qtimezone_is_daylight_time(const go_qtimezone_t self, const go_qdatetime_t* at) {
    return Z(self)->isDaylightTime(*DT(at));
}
bool qtimezone_has_daylight_time(const go_qtimezone_t self) {
    return Z(self)->hasDaylightTime();
}
bool qtimezone_has_transitions(const go_qtimezone_t self) {
    return Z(self)->hasTransitions();
}

/* --- Display names --- */
void qtimezone_display_name_at(const go_qtimezone_t self, const go_qdatetime_t* at,
                               int name_type, const void* locale, void** out_qstr) {
    const QLocale* loc = static_cast<const QLocale*>(locale);
    *out_qstr = new QString(Z(self)->displayName(
        *DT(at), static_cast<QTimeZone::NameType>(name_type), loc ? *loc : QLocale()));
}
void qtimezone_display_name_type(const go_qtimezone_t self, int time_type, int name_type,
                                 const void* locale, void** out_qstr) {
    const QLocale* loc = static_cast<const QLocale*>(locale);
    *out_qstr = new QString(Z(self)->displayName(
        static_cast<QTimeZone::TimeType>(time_type),
        static_cast<QTimeZone::NameType>(name_type),
        loc ? *loc : QLocale()));
}
void qtimezone_abbreviation(const go_qtimezone_t self, const go_qdatetime_t* at, void** out_qstr) {
    *out_qstr = new QString(Z(self)->abbreviation(*DT(at)));
}

/* --- Transitions --- */
void qtimezone_next_transition(const go_qtimezone_t self, const go_qdatetime_t* after,
                               int64_t* out_offset, int64_t* out_std_offset,
                               int64_t* out_dst_offset, void** out_at_utc, void** out_abbr) {
    const QTimeZone::OffsetData od = Z(self)->nextTransition(*DT(after));
    *out_offset     = od.offsetFromUtc;
    *out_std_offset = od.standardTimeOffset;
    *out_dst_offset = od.daylightTimeOffset;
    *out_at_utc     = new QDateTime(od.atUtc);
    *out_abbr       = new QString(od.abbreviation);
}
void qtimezone_previous_transition(const go_qtimezone_t self, const go_qdatetime_t* before,
                                   int64_t* out_offset, int64_t* out_std_offset,
                                   int64_t* out_dst_offset, void** out_at_utc, void** out_abbr) {
    const QTimeZone::OffsetData od = Z(self)->previousTransition(*DT(before));
    *out_offset     = od.offsetFromUtc;
    *out_std_offset = od.standardTimeOffset;
    *out_dst_offset = od.daylightTimeOffset;
    *out_at_utc     = new QDateTime(od.atUtc);
    *out_abbr       = new QString(od.abbreviation);
}

/* --- Static --- */
void qtimezone_system_id(void** out_qba) {
    *out_qba = new QByteArray(QTimeZone::systemTimeZoneId());
}
go_qtimezone_t qtimezone_system_zone(void) { return toGo(new QTimeZone(QTimeZone::systemTimeZone())); }
go_qtimezone_t qtimezone_static_utc(void)  { return toGo(new QTimeZone(QTimeZone::utc())); }
bool qtimezone_is_time_zone_id_available(const char* id, ptrdiff_t id_len) {
    return QTimeZone::isTimeZoneIdAvailable(BA(id, id_len));
}

/* --- Conversions --- */
go_qtimezone_t qtimezone_as_backend_zone(const go_qtimezone_t self) {
    return toGo(new QTimeZone(Z(self)->asBackendZone()));
}
void qtimezone_territory(const go_qtimezone_t self, int* out_territory) {
    *out_territory = static_cast<int>(Z(self)->territory());
}

} // extern "C"
