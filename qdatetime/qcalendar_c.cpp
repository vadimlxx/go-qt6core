#include "../internal/qt6core/check_qt_version.h"

#include <QCalendar>
#include <QString>
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
} // namespace

extern "C" {

go_qcalendar_t qcalendar_default() {
    return QCalendar();
}

go_qcalendar_t qcalendar_from_system(int system) {
    if (system < 0 || system > static_cast<int>(QCalendar::System::Last))
        return invalidCalendar();
    return QCalendar(static_cast<QCalendar::System>(system));
}

go_qcalendar_t qcalendar_from_systemid(uintptr_t id) {
    // QCalendar::SystemId(size_t) / SystemId(System) are private (friend
    // QCalendarBackend only). Since SystemId for enum-valued systems stores
    // exactly size_t(System), we reconstruct via System directly.
    if (id > static_cast<uintptr_t>(QCalendar::System::Last))
        return invalidCalendar();
    return QCalendar(static_cast<QCalendar::System>(id));
}

go_qcalendar_t qcalendar_from_any(const void* data, ptrdiff_t size, int tag) {
    return QCalendar(makeAny(data, size, tag));
}

/* --- Validity --- */
bool qcalendar_is_valid(go_qcalendar_t self) { return self.isValid(); }

/* --- Date queries --- */
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

/* --- Leap years --- */
bool qcalendar_is_leap_year(go_qcalendar_t self, int year) {
    return self.isLeapYear(year);
}

/* --- Calendar properties --- */
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

} // extern "C"