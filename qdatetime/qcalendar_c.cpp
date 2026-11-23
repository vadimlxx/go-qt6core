#include "../internal/qt6core/check_qt_version.h"

#include <QCalendar>
#include <QString>
#include <QStringList>
#include <QDate>
#include <QTime>
#include <QDateTime>
#include <QLocale>
#include <QAnyStringView>
#include <QStringView>
#include <QLatin1StringView>
#include <QUtf8StringView>

#include "qcalendar_c.h"

namespace {
inline QAnyStringView makeAny(const void* data, ptrdiff_t size, int tag) noexcept {
    switch (tag) {
    case 0:
        return QAnyStringView(QUtf8StringView(static_cast<const char*>(data), size));
    case 1:
        return QAnyStringView(QLatin1StringView(static_cast<const char*>(data), size));
    case 2:
        return QAnyStringView(QStringView(static_cast<const char16_t*>(data), size));
    default:
        return QAnyStringView();
    }
}

inline QCalendar invalidCalendar() {
    return QCalendar(QAnyStringView());
}

inline QCalendar::FormatType Fmt(int v) {
    return static_cast<QLocale::FormatType>(v);
}

inline QCalendar::YearMonthDay YMD(go_qymd_t y) {
    return QCalendar::YearMonthDay(y.year, y.month, y.day);
}
inline go_qymd_t toYMD(const QCalendar::YearMonthDay& y) {
    return go_qymd_t{y.year, y.month, y.day};
}

inline QDateTime* DT(const go_qdatetime_t* p) { return reinterpret_cast<QDateTime*>(const_cast<go_qdatetime_t*>(p)); }
inline QDate D(go_qdate_t jd)                 { return QDate::fromJulianDay(jd); }
inline QTime T(go_qtime_t ms) {
    if (ms < 0) return QTime();
    return QTime::fromMSecsSinceStartOfDay(ms);
}

inline QString U8(const char* data, ptrdiff_t len) {
    return QString::fromUtf8(data, static_cast<qsizetype>(len));
}
} // namespace

extern "C" {

/* ---- Construction ---- */
go_qcalendar_t qcalendar_default() { return QCalendar(); }

go_qcalendar_t qcalendar_from_system(int system) {
    if (system < 0 || system > static_cast<int>(QCalendar::System::Last))
        return invalidCalendar();
    return QCalendar(static_cast<QCalendar::System>(system));
}

go_qcalendar_t qcalendar_from_systemid(uintptr_t id) {
    if (id > static_cast<uintptr_t>(QCalendar::System::Last))
        return invalidCalendar();
    return QCalendar(static_cast<QCalendar::System>(id));
}

go_qcalendar_t qcalendar_from_any(const void* data, ptrdiff_t size, int tag) {
    return QCalendar(makeAny(data, size, tag));
}

/* ---- Validity ---- */
bool qcalendar_is_valid(go_qcalendar_t self) { return self.isValid(); }

/* ---- Date queries ---- */
int qcalendar_days_in_month(go_qcalendar_t self, int month, int year) {
    return self.daysInMonth(month, year);
}
int qcalendar_days_in_year(go_qcalendar_t self, int year) {
    return self.daysInYear(year);
}
int qcalendar_months_in_year(go_qcalendar_t self, int year) {
    return self.monthsInYear(year);
}
bool qcalendar_is_date_valid(go_qcalendar_t self, int year, int month, int day) {
    return self.isDateValid(year, month, day);
}

/* ---- Leap years ---- */
bool qcalendar_is_leap_year(go_qcalendar_t self, int year) {
    return self.isLeapYear(year);
}

/* ---- Calendar properties ---- */
bool qcalendar_is_gregorian(go_qcalendar_t self)  { return self.isGregorian(); }
bool qcalendar_is_lunar(go_qcalendar_t self)      { return self.isLunar(); }
bool qcalendar_is_luni_solar(go_qcalendar_t self) { return self.isLuniSolar(); }
bool qcalendar_is_solar(go_qcalendar_t self)      { return self.isSolar(); }
bool qcalendar_is_proleptic(go_qcalendar_t self)  { return self.isProleptic(); }
bool qcalendar_has_year_zero(go_qcalendar_t self) { return self.hasYearZero(); }
int  qcalendar_maximum_days_in_month(go_qcalendar_t self) { return self.maximumDaysInMonth(); }
int  qcalendar_minimum_days_in_month(go_qcalendar_t self) { return self.minimumDaysInMonth(); }
int  qcalendar_maximum_months_in_year(go_qcalendar_t self) { return self.maximumMonthsInYear(); }
void qcalendar_name(go_qcalendar_t self, void** out_qstr) {
    *out_qstr = new QString(self.name());
}

/* ---- QDate conversions ---- */
go_qdate_t qcalendar_date_from_parts(go_qcalendar_t self, int year, int month, int day) {
    return self.dateFromParts(year, month, day).toJulianDay();
}
go_qdate_t qcalendar_date_from_ymd(go_qcalendar_t self, go_qymd_t ymd) {
    return self.dateFromParts(YMD(ymd)).toJulianDay();
}
go_qdate_t qcalendar_match_century_to_weekday(go_qcalendar_t self, go_qymd_t ymd, int dow) {
    return self.matchCenturyToWeekday(YMD(ymd), dow).toJulianDay();
}
go_qymd_t qcalendar_parts_from_date(go_qcalendar_t self, go_qdate_t date) {
    return toYMD(self.partsFromDate(D(date)));
}
int qcalendar_day_of_week(go_qcalendar_t self, go_qdate_t date) {
    return self.dayOfWeek(D(date));
}

/* ---- Month / weekday names ---- */
void qcalendar_month_name(go_qcalendar_t self, const go_qlocale_t locale,
                          int month, int year, int format, void** out_qstr) {
    if (!locale) { *out_qstr = new QString(); return; }
    *out_qstr = new QString(self.monthName(*locale, month, year, Fmt(format)));
}
void qcalendar_standalone_month_name(go_qcalendar_t self, const go_qlocale_t locale,
                                     int month, int year, int format, void** out_qstr) {
    if (!locale) { *out_qstr = new QString(); return; }
    *out_qstr = new QString(self.standaloneMonthName(*locale, month, year, Fmt(format)));
}
void qcalendar_week_day_name(go_qcalendar_t self, const go_qlocale_t locale,
                             int day, int format, void** out_qstr) {
    if (!locale) { *out_qstr = new QString(); return; }
    *out_qstr = new QString(self.weekDayName(*locale, day, Fmt(format)));
}
void qcalendar_standalone_week_day_name(go_qcalendar_t self, const go_qlocale_t locale,
                                        int day, int format, void** out_qstr) {
    if (!locale) { *out_qstr = new QString(); return; }
    *out_qstr = new QString(self.standaloneWeekDayName(*locale, day, Fmt(format)));
}

/* ---- Date-time formatting ---- */
void qcalendar_date_time_to_string(go_qcalendar_t self, const char* format, ptrdiff_t format_len,
                                   const go_qdatetime_t* dt, go_qdate_t date_only,
                                   go_qtime_t time_only, const go_qlocale_t locale,
                                   void** out_qstr) {
    QDateTime dtVal;
    if (dt) dtVal = *DT(dt);
    QLocale loc = locale ? *locale : QLocale();
    *out_qstr = new QString(self.dateTimeToString(
        U8(format, format_len), dtVal, D(date_only), T(time_only), loc));
}

/* ---- Available calendars ---- */
QString** qcalendar_available_calendars(ptrdiff_t* out_count) {
    QStringList list = QCalendar::availableCalendars();
    const qsizetype n = list.size();
    *out_count = static_cast<ptrdiff_t>(n);
    QString** arr = new QString*[n > 0 ? static_cast<size_t>(n) : 1];
    for (qsizetype i = 0; i < n; ++i) arr[i] = new QString(list[i]);
    return arr;
}
void qcalendar_free_string_array(QString** arr) { delete[] arr; }

} // extern "C"