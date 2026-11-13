#include "../internal/qt6core/check_qt_version.h"

#include <QDateTime>
#include <QDate>
#include <QTime>
#include <QString>
#include <QTimeZone>

#include "qdatetime_c.h"
#include "qcalendar_c.h"

namespace {
inline QDateTime*       DT(go_qdatetime_t* p)       { return reinterpret_cast<QDateTime*>(p); }
inline const QDateTime* DT(const go_qdatetime_t* p) { return reinterpret_cast<const QDateTime*>(p); }
inline go_qdatetime_t*  toGo(QDateTime* p)          { return reinterpret_cast<go_qdatetime_t*>(p); }

inline QDate D(go_qdate_t d) { return QDate::fromJulianDay(d); }
inline QTime T(go_qtime_t t) {
    if (t < 0) return QTime();
    return QTime::fromMSecsSinceStartOfDay(t);
}
inline go_qtime_t toTimeGo(const QTime& t) {
    if (t.isNull()) return -1;
    return t.msecsSinceStartOfDay();
}

/* Go strings are NOT null-terminated. Use (data, len). */
inline QString U8(const char* data, ptrdiff_t len) {
    return QString::fromUtf8(data, static_cast<qsizetype>(len));
}
}

extern "C" {

go_qdatetime_t* qdatetime_new_default(void) { return toGo(new QDateTime()); }

go_qdatetime_t* qdatetime_new_from_date_time(go_qdate_t date, go_qtime_t time) {
    return toGo(new QDateTime(D(date), T(time)));
}

go_qdatetime_t* qdatetime_new_copy(const go_qdatetime_t* other) {
    return toGo(new QDateTime(*DT(other)));
}
go_qdatetime_t* qdatetime_clone(const go_qdatetime_t* other) {
    return toGo(new QDateTime(*DT(other)));
}
void qdatetime_delete(go_qdatetime_t* dt) { delete DT(dt); }

bool qdatetime_is_null(const go_qdatetime_t* dt)  { return DT(dt)->isNull(); }
bool qdatetime_is_valid(const go_qdatetime_t* dt) { return DT(dt)->isValid(); }

go_qdate_t qdatetime_date(const go_qdatetime_t* dt) {
    return DT(dt)->date().toJulianDay();
}
go_qtime_t qdatetime_time(const go_qdatetime_t* dt) {
    return toTimeGo(DT(dt)->time());
}
int qdatetime_time_spec(const go_qdatetime_t* dt) {
    return static_cast<int>(DT(dt)->timeSpec());
}
int qdatetime_offset_from_utc(const go_qdatetime_t* dt) {
    return DT(dt)->offsetFromUtc();
}
int64_t qdatetime_to_msecs_since_epoch(const go_qdatetime_t* dt) {
    return DT(dt)->toMSecsSinceEpoch();
}
int64_t qdatetime_to_secs_since_epoch(const go_qdatetime_t* dt) {
    return DT(dt)->toSecsSinceEpoch();
}

void qdatetime_set_date(go_qdatetime_t* dt, go_qdate_t date) {
    DT(dt)->setDate(D(date));
}
void qdatetime_set_time(go_qdatetime_t* dt, go_qtime_t time) {
    DT(dt)->setTime(T(time));
}
void qdatetime_set_msecs_since_epoch(go_qdatetime_t* dt, int64_t msecs) {
    DT(dt)->setMSecsSinceEpoch(msecs);
}
void qdatetime_set_secs_since_epoch(go_qdatetime_t* dt, int64_t secs) {
    DT(dt)->setSecsSinceEpoch(secs);
}

go_qdatetime_t* qdatetime_add_days(const go_qdatetime_t* dt, int64_t days) {
    return toGo(new QDateTime(DT(dt)->addDays(days)));
}
go_qdatetime_t* qdatetime_add_months(const go_qdatetime_t* dt, int months) {
    return toGo(new QDateTime(DT(dt)->addMonths(months)));
}
go_qdatetime_t* qdatetime_add_years(const go_qdatetime_t* dt, int years) {
    return toGo(new QDateTime(DT(dt)->addYears(years)));
}
go_qdatetime_t* qdatetime_add_secs(const go_qdatetime_t* dt, int64_t secs) {
    return toGo(new QDateTime(DT(dt)->addSecs(secs)));
}
go_qdatetime_t* qdatetime_add_msecs(const go_qdatetime_t* dt, int64_t msecs) {
    return toGo(new QDateTime(DT(dt)->addMSecs(msecs)));
}

int64_t qdatetime_days_to(const go_qdatetime_t* dt, const go_qdatetime_t* other) {
    return DT(dt)->daysTo(*DT(other));
}
int64_t qdatetime_secs_to(const go_qdatetime_t* dt, const go_qdatetime_t* other) {
    return DT(dt)->secsTo(*DT(other));
}
int64_t qdatetime_msecs_to(const go_qdatetime_t* dt, const go_qdatetime_t* other) {
    return DT(dt)->msecsTo(*DT(other));
}

go_qdatetime_t* qdatetime_current_date_time(void) {
    return toGo(new QDateTime(QDateTime::currentDateTime()));
}
go_qdatetime_t* qdatetime_current_date_time_utc(void) {
    return toGo(new QDateTime(QDateTime::currentDateTimeUtc()));
}
go_qdatetime_t* qdatetime_from_msecs_since_epoch(int64_t msecs) {
    return toGo(new QDateTime(QDateTime::fromMSecsSinceEpoch(msecs)));
}
go_qdatetime_t* qdatetime_from_secs_since_epoch(int64_t secs) {
    return toGo(new QDateTime(QDateTime::fromSecsSinceEpoch(secs)));
}
int64_t qdatetime_current_msecs_since_epoch(void) { return QDateTime::currentMSecsSinceEpoch(); }
int64_t qdatetime_current_secs_since_epoch(void)  { return QDateTime::currentSecsSinceEpoch(); }

void qdatetime_to_string(const go_qdatetime_t* dt, int format, void** out_qstr) {
    *out_qstr = new QString(DT(dt)->toString(static_cast<Qt::DateFormat>(format)));
}
void qdatetime_to_string_format(const go_qdatetime_t* dt, const char* format,
                                ptrdiff_t format_len, void** out_qstr) {
    *out_qstr = new QString(DT(dt)->toString(U8(format, format_len)));
}
void qdatetime_to_string_format_cal(const go_qdatetime_t* dt, const char* format,
                                    ptrdiff_t format_len, go_qcalendar_t cal,
                                    void** out_qstr) {
    *out_qstr = new QString(DT(dt)->toString(U8(format, format_len), cal));
}

go_qdatetime_t* qdatetime_from_string_utf8(const char* str, ptrdiff_t str_len, int format) {
    return toGo(new QDateTime(QDateTime::fromString(
        U8(str, str_len), static_cast<Qt::DateFormat>(format))));
}
go_qdatetime_t* qdatetime_from_string_format_utf8(const char* str, ptrdiff_t str_len,
                                                  const char* format, ptrdiff_t format_len) {
    return toGo(new QDateTime(QDateTime::fromString(
        U8(str, str_len), U8(format, format_len))));
}

bool qdatetime_equals(const go_qdatetime_t* dt, const go_qdatetime_t* other) {
    return *DT(dt) == *DT(other);
}

go_qdatetime_t* qdatetime_to_local_time(const go_qdatetime_t* dt) {
    return toGo(new QDateTime(DT(dt)->toLocalTime()));
}
go_qdatetime_t* qdatetime_to_utc(const go_qdatetime_t* dt) {
    return toGo(new QDateTime(DT(dt)->toUTC()));
}
go_qdatetime_t* qdatetime_to_offset_from_utc(const go_qdatetime_t* dt, int offset_seconds) {
    return toGo(new QDateTime(DT(dt)->toOffsetFromUtc(offset_seconds)));
}

/* --- QTimeZone integration --- */
go_qtimezone_t qdatetime_time_zone(const go_qdatetime_t* dt) {
	return reinterpret_cast<go_qtimezone_t>(new QTimeZone(DT(dt)->timeZone()));
}
void qdatetime_set_time_zone(go_qdatetime_t* dt, const go_qtimezone_t tz) {
	auto* z = reinterpret_cast<QTimeZone*>(tz);
	DT(dt)->setTimeZone(*z);
}
go_qdatetime_t* qdatetime_to_time_zone(const go_qdatetime_t* dt, const go_qtimezone_t tz) {
	auto* z = reinterpret_cast<QTimeZone*>(tz);
	return toGo(new QDateTime(DT(dt)->toTimeZone(*z)));
}
go_qdatetime_t* qdatetime_current_date_time_zone(const go_qtimezone_t tz) {
	auto* z = reinterpret_cast<QTimeZone*>(tz);
	return toGo(new QDateTime(QDateTime::currentDateTime(*z)));
}
go_qdatetime_t* qdatetime_from_msecs_since_epoch_zone(int64_t msecs, const go_qtimezone_t tz) {
	auto* z = reinterpret_cast<QTimeZone*>(tz);
	return toGo(new QDateTime(QDateTime::fromMSecsSinceEpoch(msecs, *z)));
}
go_qdatetime_t* qdatetime_from_secs_since_epoch_zone(int64_t secs, const go_qtimezone_t tz) {
	auto* z = reinterpret_cast<QTimeZone*>(tz);
	return toGo(new QDateTime(QDateTime::fromSecsSinceEpoch(secs, *z)));
}

} // extern "C"
