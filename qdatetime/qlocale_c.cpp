#include "../internal/qt6core/check_qt_version.h"

#include <QLocale>
#include <QString>

#include "qlocale_c.h"

extern "C" {

/* --- Lifecycle --- */
go_qlocale_t qlocale_new_default() { return new QLocale(); }

go_qlocale_t qlocale_new_from_name(const char* name, ptrdiff_t len) {
    return new QLocale(QString::fromUtf8(name, static_cast<int>(len)));
}

go_qlocale_t qlocale_new_from_language_country(int language, int country) {
    return new QLocale(static_cast<QLocale::Language>(language),
                       static_cast<QLocale::Country>(country));
}

go_qlocale_t qlocale_new_copy(const go_qlocale_t other) { return new QLocale(*other); }
go_qlocale_t qlocale_clone(const go_qlocale_t self) { return new QLocale(*self); }
void         qlocale_delete(go_qlocale_t self) { delete self; }

/* --- Accessors --- */
int qlocale_language(const go_qlocale_t self) { return static_cast<int>(self->language()); }
int qlocale_script(const go_qlocale_t self) { return static_cast<int>(self->script()); }
int qlocale_country(const go_qlocale_t self) { return static_cast<int>(self->territory()); }
int qlocale_measurement_system(const go_qlocale_t self) {
    return static_cast<int>(self->measurementSystem());
}
int qlocale_text_direction(const go_qlocale_t self) {
    return static_cast<int>(self->textDirection());
}

/* --- String representation --- */
void qlocale_name(const go_qlocale_t self, void** out_qstr) {
    *out_qstr = new QString(self->name());
}
void qlocale_bcp47_name(const go_qlocale_t self, void** out_qstr) {
    *out_qstr = new QString(self->bcp47Name());
}
void qlocale_native_language_name(const go_qlocale_t self, void** out_qstr) {
    *out_qstr = new QString(self->nativeLanguageName());
}
void qlocale_native_territory_name(const go_qlocale_t self, void** out_qstr) {
    *out_qstr = new QString(self->nativeTerritoryName());
}

/* --- Number parsing --- */
bool qlocale_to_int64(const go_qlocale_t self, const char* str, ptrdiff_t len, int64_t* out_value) {
    bool ok;
    qlonglong v = self->toLongLong(QString::fromUtf8(str, static_cast<int>(len)), &ok);
    *out_value = static_cast<int64_t>(v);
    return ok;
}
bool qlocale_to_double(const go_qlocale_t self, const char* str, ptrdiff_t len, double* out_value) {
    bool ok;
    double v = self->toDouble(QString::fromUtf8(str, static_cast<int>(len)), &ok);
    *out_value = v;
    return ok;
}

/* --- Number formatting --- */
void qlocale_to_string_int64(const go_qlocale_t self, int64_t value, void** out_qstr) {
    *out_qstr = new QString(self->toString(static_cast<qlonglong>(value)));
}
void qlocale_to_string_uint64(const go_qlocale_t self, uint64_t value, void** out_qstr) {
    *out_qstr = new QString(self->toString(static_cast<qulonglong>(value)));
}
void qlocale_to_string_double(const go_qlocale_t self, double value, char format, int precision,
                              void** out_qstr) {
    *out_qstr = new QString(self->toString(value, format, precision));
}

/* --- Locale-specific strings --- */
void qlocale_decimal_point(const go_qlocale_t self, void** out_qstr) {
    *out_qstr = new QString(self->decimalPoint());
}
void qlocale_group_separator(const go_qlocale_t self, void** out_qstr) {
    *out_qstr = new QString(self->groupSeparator());
}
void qlocale_percent(const go_qlocale_t self, void** out_qstr) {
    *out_qstr = new QString(self->percent());
}
void qlocale_zero_digit(const go_qlocale_t self, void** out_qstr) {
    *out_qstr = new QString(self->zeroDigit());
}
void qlocale_negative_sign(const go_qlocale_t self, void** out_qstr) {
    *out_qstr = new QString(self->negativeSign());
}
void qlocale_positive_sign(const go_qlocale_t self, void** out_qstr) {
    *out_qstr = new QString(self->positiveSign());
}
void qlocale_exponential(const go_qlocale_t self, void** out_qstr) {
    *out_qstr = new QString(self->exponential());
}

/* --- Month / day names --- */
void qlocale_month_name(const go_qlocale_t self, int month, int format, void** out_qstr) {
    *out_qstr = new QString(self->monthName(month, static_cast<QLocale::FormatType>(format)));
}
void qlocale_standalone_month_name(const go_qlocale_t self, int month, int format, void** out_qstr) {
    *out_qstr = new QString(self->standaloneMonthName(month,
                                                      static_cast<QLocale::FormatType>(format)));
}
void qlocale_day_name(const go_qlocale_t self, int day, int format, void** out_qstr) {
    *out_qstr = new QString(self->dayName(day, static_cast<QLocale::FormatType>(format)));
}
void qlocale_standalone_day_name(const go_qlocale_t self, int day, int format, void** out_qstr) {
    *out_qstr = new QString(self->standaloneDayName(day,
                                                    static_cast<QLocale::FormatType>(format)));
}

/* --- AM/PM --- */
void qlocale_am_text(const go_qlocale_t self, void** out_qstr) {
    *out_qstr = new QString(self->amText());
}
void qlocale_pm_text(const go_qlocale_t self, void** out_qstr) {
    *out_qstr = new QString(self->pmText());
}

/* --- Case conversion --- */
void qlocale_to_upper(const go_qlocale_t self, const char* str, ptrdiff_t len, void** out_qstr) {
    QString s = QString::fromUtf8(str, static_cast<int>(len));
    *out_qstr = new QString(self->toUpper(s));
}
void qlocale_to_lower(const go_qlocale_t self, const char* str, ptrdiff_t len, void** out_qstr) {
    QString s = QString::fromUtf8(str, static_cast<int>(len));
    *out_qstr = new QString(self->toLower(s));
}

/* --- Static methods --- */
void qlocale_set_default(const go_qlocale_t locale) { QLocale::setDefault(*locale); }
go_qlocale_t qlocale_c() { return new QLocale(QLocale::c()); }
go_qlocale_t qlocale_system() { return new QLocale(QLocale::system()); }

} // extern "C"