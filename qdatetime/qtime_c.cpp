#include "../internal/qt6core/check_qt_version.h"

#include <QTime>
#include <QString>

#include "qtime_c.h"

namespace {

inline QTime T(go_qtime_t t) {
    if (t < 0) return QTime();
    return QTime::fromMSecsSinceStartOfDay(t);
}
inline go_qtime_t toGo(const QTime& t) {
    if (t.isNull()) return -1;
    return t.msecsSinceStartOfDay();
}

/* Go strings are NOT null-terminated. Use (data, len). */
inline QString U8(const char* data, ptrdiff_t len) {
    return QString::fromUtf8(data, static_cast<qsizetype>(len));
}

}

extern "C" {

go_qtime_t qtime_default(void)                         { return toGo(QTime()); }
go_qtime_t qtime_from_hms(int h, int m, int s, int ms) { return toGo(QTime(h, m, s, ms)); }
go_qtime_t qtime_current_time(void)                    { return toGo(QTime::currentTime()); }

go_qtime_t qtime_from_string_utf8(const char* str, ptrdiff_t str_len, int format) {
    return toGo(QTime::fromString(U8(str, str_len),
                                  static_cast<Qt::DateFormat>(format)));
}
go_qtime_t qtime_from_string_format_utf8(const char* str, ptrdiff_t str_len,
                                         const char* format, ptrdiff_t format_len) {
    return toGo(QTime::fromString(U8(str, str_len), U8(format, format_len)));
}
go_qtime_t qtime_from_msecs_since_start_of_day(int msecs) {
    if (msecs < 0) return -1;
    return toGo(QTime::fromMSecsSinceStartOfDay(msecs));
}

bool qtime_is_null(go_qtime_t t)                  { return T(t).isNull(); }
bool qtime_is_valid(go_qtime_t t)                 { return T(t).isValid(); }
int  qtime_hour(go_qtime_t t)                     { return T(t).hour(); }
int  qtime_minute(go_qtime_t t)                   { return T(t).minute(); }
int  qtime_second(go_qtime_t t)                   { return T(t).second(); }
int  qtime_msec(go_qtime_t t)                     { return T(t).msec(); }
int  qtime_msecs_since_start_of_day(go_qtime_t t) { return T(t).msecsSinceStartOfDay(); }

go_qtime_t qtime_add_secs(go_qtime_t t, int secs)   { return toGo(T(t).addSecs(secs));  }
go_qtime_t qtime_add_msecs(go_qtime_t t, int msecs) { return toGo(T(t).addMSecs(msecs)); }
int qtime_secs_to(go_qtime_t t, go_qtime_t other)   { return T(t).secsTo(T(other));      }
int qtime_msecs_to(go_qtime_t t, go_qtime_t other)  { return T(t).msecsTo(T(other));     }

void qtime_to_string(go_qtime_t t, int format, void** out_qstr) {
    *out_qstr = new QString(T(t).toString(static_cast<Qt::DateFormat>(format)));
}
void qtime_to_string_format(go_qtime_t t, const char* format, ptrdiff_t format_len,
                            void** out_qstr) {
    *out_qstr = new QString(T(t).toString(U8(format, format_len)));
}

bool qtime_set_hms(go_qtime_t* t, int h, int m, int s, int ms) {
    QTime qt = T(*t);
    bool ok = qt.setHMS(h, m, s, ms);
    *t = toGo(qt);
    return ok;
}

bool qtime_static_is_valid(int h, int m, int s, int ms) { return QTime::isValid(h, m, s, ms); }

} // extern "C"