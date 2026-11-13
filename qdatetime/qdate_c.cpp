#include "../internal/qt6core/check_qt_version.h"

#include <QDate>
#include <QString>
#include <QCalendar>

#include "qdate_c.h"
#include "qcalendar_c.h"

namespace {
inline QDate D(go_qdate_t d) { return QDate::fromJulianDay(d); }
inline go_qdate_t toGo(const QDate& d) { return d.toJulianDay(); }

/* Go strings are NOT null-terminated. Use the (data, len) overload of
 * QString::fromUtf8 everywhere a Go string crosses the C ABI. */
inline QString U8(const char* data, ptrdiff_t len) {
    return QString::fromUtf8(data, static_cast<qsizetype>(len));
}
}

extern "C" {

go_qdate_t qdate_default(void)                              { return toGo(QDate()); }
go_qdate_t qdate_from_ymd(int y, int m, int d)              { return toGo(QDate(y, m, d)); }
go_qdate_t qdate_from_ymd_cal(int y, int m, int d, go_qcalendar_t cal) {
    return toGo(QDate(y, m, d, cal));
}
go_qdate_t qdate_from_julian_day(int64_t jd)                { return toGo(QDate::fromJulianDay(jd)); }
go_qdate_t qdate_current_date(void)                         { return toGo(QDate::currentDate()); }

go_qdate_t qdate_from_string_utf8(const char* str, ptrdiff_t str_len, int format) {
    return toGo(QDate::fromString(U8(str, str_len),
                                  static_cast<Qt::DateFormat>(format)));
}
go_qdate_t qdate_from_string_format_utf8(const char* str, ptrdiff_t str_len,
                                         const char* format, ptrdiff_t format_len) {
    return toGo(QDate::fromString(U8(str, str_len), U8(format, format_len)));
}

bool qdate_is_null(go_qdate_t d)              { return D(d).isNull(); }
bool qdate_is_valid(go_qdate_t d)             { return D(d).isValid(); }
int  qdate_year(go_qdate_t d)                 { return D(d).year(); }
int  qdate_month(go_qdate_t d)                { return D(d).month(); }
int  qdate_day(go_qdate_t d)                  { return D(d).day(); }
int  qdate_day_of_week(go_qdate_t d)          { return D(d).dayOfWeek(); }
int  qdate_day_of_year(go_qdate_t d)          { return D(d).dayOfYear(); }
int  qdate_days_in_month(go_qdate_t d)        { return D(d).daysInMonth(); }
int  qdate_days_in_year(go_qdate_t d)         { return D(d).daysInYear(); }
int  qdate_week_number(go_qdate_t d, int* yn) { return D(d).weekNumber(yn); }

int qdate_year_cal (go_qdate_t d, go_qcalendar_t cal) { return D(d).year(cal);  }
int qdate_month_cal(go_qdate_t d, go_qcalendar_t cal) { return D(d).month(cal); }
int qdate_day_cal  (go_qdate_t d, go_qcalendar_t cal) { return D(d).day(cal);   }

go_qdate_t qdate_add_days(go_qdate_t d, int64_t days) {
    return toGo(D(d).addDays(days));
}
go_qdate_t qdate_add_months(go_qdate_t d, int months) {
    return toGo(D(d).addMonths(months));
}
go_qdate_t qdate_add_years(go_qdate_t d, int years) {
    return toGo(D(d).addYears(years));
}
go_qdate_t qdate_add_months_cal(go_qdate_t d, int months, go_qcalendar_t cal) {
    return toGo(D(d).addMonths(months, cal));
}
go_qdate_t qdate_add_years_cal(go_qdate_t d, int years, go_qcalendar_t cal) {
    return toGo(D(d).addYears(years, cal));
}
int64_t qdate_days_to(go_qdate_t d, go_qdate_t other) {
    return D(d).daysTo(D(other));
}

void qdate_to_string(go_qdate_t d, int format, void** out_qstr) {
    *out_qstr = new QString(D(d).toString(static_cast<Qt::DateFormat>(format)));
}
void qdate_to_string_format(go_qdate_t d, const char* format, ptrdiff_t format_len,
                            void** out_qstr) {
    *out_qstr = new QString(D(d).toString(U8(format, format_len)));
}
void qdate_to_string_format_cal(go_qdate_t d, const char* format, ptrdiff_t format_len,
                                go_qcalendar_t cal, void** out_qstr) {
    *out_qstr = new QString(D(d).toString(U8(format, format_len), cal));
}

/* In-place setters: write the mutated QDate back to *d. */
bool qdate_set_date(go_qdate_t* d, int y, int m, int day) {
    QDate qd = QDate::fromJulianDay(*d);
    bool ok = qd.setDate(y, m, day);
    *d = qd.toJulianDay();
    return ok;
}
bool qdate_set_date_cal(go_qdate_t* d, int y, int m, int day, go_qcalendar_t cal) {
    QDate qd = QDate::fromJulianDay(*d);
    bool ok = qd.setDate(y, m, day, cal);
    *d = qd.toJulianDay();
    return ok;
}
void qdate_get_date(go_qdate_t d, int* y, int* m, int* day) {
    D(d).getDate(y, m, day);
}

int64_t qdate_to_julian_day(go_qdate_t d) { return D(d).toJulianDay(); }

bool qdate_static_is_valid(int y, int m, int d) { return QDate::isValid(y, m, d); }
bool qdate_static_is_leap_year(int year)        { return QDate::isLeapYear(year); }

} // extern "C"