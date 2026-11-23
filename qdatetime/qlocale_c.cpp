#include "../internal/qt6core/check_qt_version.h"

#include <QLocale>
#include <QString>
#include <QStringList>
#include <QDate>
#include <QTime>
#include <QDateTime>
#include <QCalendar>
#include <QtCore/qnamespace.h>

#include "qlocale_c.h"

namespace {

/* Enum conversions — concentrated here to keep function bodies cast-free. */
inline QLocale::Language            Lang(int v)     { return static_cast<QLocale::Language>(v); }
inline QLocale::Script              Scr(int v)      { return static_cast<QLocale::Script>(v); }
inline QLocale::Territory           Terr(int v)     { return static_cast<QLocale::Territory>(v); }
inline QLocale::FormatType          Fmt(int v)      { return static_cast<QLocale::FormatType>(v); }
inline QLocale::NumberOptions       NumOpts(int v)  { return static_cast<QLocale::NumberOptions>(v); }
inline QLocale::CurrencySymbolFormat CurFmt(int v)  { return static_cast<QLocale::CurrencySymbolFormat>(v); }
inline QLocale::DataSizeFormats     DataFmt(int v)  { return static_cast<QLocale::DataSizeFormats>(v); }
inline QLocale::QuotationStyle      QuoteStyle(int v){ return static_cast<QLocale::QuotationStyle>(v); }
inline QLocale::TagSeparator        TagSep(int v)   { return static_cast<QLocale::TagSeparator>(static_cast<char>(v)); }
inline QLocale::LanguageCodeTypes   LangCodes(int v){ return static_cast<QLocale::LanguageCodeTypes>(v); }

/* Handle helpers (mirror qdatetime_c.cpp). */
inline const QDateTime* DT(const go_qdatetime_t* p) { return reinterpret_cast<const QDateTime*>(p); }
inline QDateTime*       DT(go_qdatetime_t* p)       { return reinterpret_cast<QDateTime*>(p); }

inline QString U8(const char* d, ptrdiff_t n) {
    return QString::fromUtf8(d, static_cast<int>(n));
}

inline QDate JD(int64_t jd)   { return QDate::fromJulianDay(static_cast<qint64>(jd)); }
inline QTime MS(int32_t ms)   { return QTime::fromMSecsSinceStartOfDay(static_cast<int>(ms)); }

} // namespace

extern "C" {

/* =========== Lifecycle =========== */

go_qlocale_t qlocale_new_default() { return new QLocale(); }
go_qlocale_t qlocale_new_from_name(const char* name, ptrdiff_t len) {
    return new QLocale(U8(name, len));
}
go_qlocale_t qlocale_new_from_language_territory(int language, int territory) {
    return new QLocale(Lang(language), Terr(territory));
}
go_qlocale_t qlocale_new_from_language_script_territory(int language, int script, int territory) {
    return new QLocale(Lang(language), Scr(script), Terr(territory));
}
go_qlocale_t qlocale_new_copy(const go_qlocale_t other) {
    return other ? new QLocale(*other) : new QLocale();
}
void qlocale_delete(go_qlocale_t self) { delete self; }
void qlocale_swap(go_qlocale_t self, go_qlocale_t other) {
    if (self && other) self->swap(*other);
}

/* =========== Value accessors =========== */

int qlocale_language(const go_qlocale_t self)           { return static_cast<int>(self->language()); }
int qlocale_script(const go_qlocale_t self)             { return static_cast<int>(self->script()); }
int qlocale_territory(const go_qlocale_t self)          { return static_cast<int>(self->territory()); }
int qlocale_measurement_system(const go_qlocale_t self) { return static_cast<int>(self->measurementSystem()); }
int qlocale_text_direction(const go_qlocale_t self)     { return static_cast<int>(self->textDirection()); }
int qlocale_first_day_of_week(const go_qlocale_t self)  { return static_cast<int>(self->firstDayOfWeek()); }
int qlocale_number_options(const go_qlocale_t self)     { return static_cast<int>(self->numberOptions()); }
void qlocale_set_number_options(go_qlocale_t self, int options) {
    self->setNumberOptions(NumOpts(options));
}

/* =========== String representation =========== */

QString* qlocale_name(const go_qlocale_t self, int separator) {
    return new QString(self->name(TagSep(separator)));
}
QString* qlocale_bcp47_name(const go_qlocale_t self, int separator) {
    return new QString(self->bcp47Name(TagSep(separator)));
}
QString* qlocale_native_language_name(const go_qlocale_t self) {
    return new QString(self->nativeLanguageName());
}
QString* qlocale_native_territory_name(const go_qlocale_t self) {
    return new QString(self->nativeTerritoryName());
}

/* =========== Number parsing =========== */

bool qlocale_to_short(const go_qlocale_t self, const char* str, ptrdiff_t len, int16_t* out) {
    bool ok; *out = static_cast<int16_t>(self->toShort(U8(str, len), &ok)); return ok;
}
bool qlocale_to_ushort(const go_qlocale_t self, const char* str, ptrdiff_t len, uint16_t* out) {
    bool ok; *out = static_cast<uint16_t>(self->toUShort(U8(str, len), &ok)); return ok;
}
bool qlocale_to_int(const go_qlocale_t self, const char* str, ptrdiff_t len, int32_t* out) {
    bool ok; *out = static_cast<int32_t>(self->toInt(U8(str, len), &ok)); return ok;
}
bool qlocale_to_uint(const go_qlocale_t self, const char* str, ptrdiff_t len, uint32_t* out) {
    bool ok; *out = static_cast<uint32_t>(self->toUInt(U8(str, len), &ok)); return ok;
}
bool qlocale_to_long(const go_qlocale_t self, const char* str, ptrdiff_t len, int64_t* out) {
    bool ok; *out = static_cast<int64_t>(self->toLong(U8(str, len), &ok)); return ok;
}
bool qlocale_to_ulong(const go_qlocale_t self, const char* str, ptrdiff_t len, uint64_t* out) {
    bool ok; *out = static_cast<uint64_t>(self->toULong(U8(str, len), &ok)); return ok;
}
bool qlocale_to_longlong(const go_qlocale_t self, const char* str, ptrdiff_t len, int64_t* out) {
    bool ok; *out = static_cast<int64_t>(self->toLongLong(U8(str, len), &ok)); return ok;
}
bool qlocale_to_ulonglong(const go_qlocale_t self, const char* str, ptrdiff_t len, uint64_t* out) {
    bool ok; *out = static_cast<uint64_t>(self->toULongLong(U8(str, len), &ok)); return ok;
}
bool qlocale_to_float(const go_qlocale_t self, const char* str, ptrdiff_t len, float* out) {
    bool ok; *out = self->toFloat(U8(str, len), &ok); return ok;
}
bool qlocale_to_double(const go_qlocale_t self, const char* str, ptrdiff_t len, double* out) {
    bool ok; *out = self->toDouble(U8(str, len), &ok); return ok;
}

/* =========== Number formatting =========== */

QString* qlocale_to_string_ll(const go_qlocale_t self, int64_t value) {
    return new QString(self->toString(static_cast<qlonglong>(value)));
}
QString* qlocale_to_string_ull(const go_qlocale_t self, uint64_t value) {
    return new QString(self->toString(static_cast<qulonglong>(value)));
}
QString* qlocale_to_string_ll_w(const go_qlocale_t self, int64_t value, int fieldWidth, uint32_t fillChar) {
    return new QString(self->toString(static_cast<qlonglong>(value), fieldWidth,
                                      static_cast<char32_t>(fillChar)));
}
QString* qlocale_to_string_ull_w(const go_qlocale_t self, uint64_t value, int fieldWidth, uint32_t fillChar) {
    return new QString(self->toString(static_cast<qulonglong>(value), fieldWidth,
                                      static_cast<char32_t>(fillChar)));
}
QString* qlocale_to_string_double(const go_qlocale_t self, double value, char format, int precision) {
    return new QString(self->toString(value, format, precision));
}
QString* qlocale_to_string_float(const go_qlocale_t self, float value, char format, int precision) {
    return new QString(self->toString(value, format, precision));
}

/* =========== Locale-specific strings =========== */

QString* qlocale_decimal_point(const go_qlocale_t self)   { return new QString(self->decimalPoint()); }
QString* qlocale_group_separator(const go_qlocale_t self) { return new QString(self->groupSeparator()); }
QString* qlocale_percent(const go_qlocale_t self)         { return new QString(self->percent()); }
QString* qlocale_zero_digit(const go_qlocale_t self)      { return new QString(self->zeroDigit()); }
QString* qlocale_negative_sign(const go_qlocale_t self)   { return new QString(self->negativeSign()); }
QString* qlocale_positive_sign(const go_qlocale_t self)   { return new QString(self->positiveSign()); }
QString* qlocale_exponential(const go_qlocale_t self)     { return new QString(self->exponential()); }

/* =========== Month/day names =========== */

QString* qlocale_month_name(const go_qlocale_t self, int month, int format) {
    return new QString(self->monthName(month, Fmt(format)));
}
QString* qlocale_standalone_month_name(const go_qlocale_t self, int month, int format) {
    return new QString(self->standaloneMonthName(month, Fmt(format)));
}
QString* qlocale_day_name(const go_qlocale_t self, int day, int format) {
    return new QString(self->dayName(day, Fmt(format)));
}
QString* qlocale_standalone_day_name(const go_qlocale_t self, int day, int format) {
    return new QString(self->standaloneDayName(day, Fmt(format)));
}

/* =========== AM/PM =========== */

QString* qlocale_am_text(const go_qlocale_t self) { return new QString(self->amText()); }
QString* qlocale_pm_text(const go_qlocale_t self) { return new QString(self->pmText()); }

/* =========== Case conversion =========== */

QString* qlocale_to_upper(const go_qlocale_t self, const char* str, ptrdiff_t len) {
    return new QString(self->toUpper(U8(str, len)));
}
QString* qlocale_to_lower(const go_qlocale_t self, const char* str, ptrdiff_t len) {
    return new QString(self->toLower(U8(str, len)));
}

/* =========== Currency =========== */

QString* qlocale_currency_symbol(const go_qlocale_t self, int format) {
    return new QString(self->currencySymbol(CurFmt(format)));
}
QString* qlocale_to_currency_string_ll(const go_qlocale_t self, int64_t value, const QString* symbol) {
    return new QString(self->toCurrencyString(static_cast<qlonglong>(value),
                                               symbol ? *symbol : QString()));
}
QString* qlocale_to_currency_string_ull(const go_qlocale_t self, uint64_t value, const QString* symbol) {
    return new QString(self->toCurrencyString(static_cast<qulonglong>(value),
                                               symbol ? *symbol : QString()));
}
QString* qlocale_to_currency_string_double(const go_qlocale_t self, double value,
                                           const QString* symbol, int precision) {
    return new QString(self->toCurrencyString(value, symbol ? *symbol : QString(), precision));
}

/* =========== Data size =========== */

QString* qlocale_formatted_data_size(const go_qlocale_t self, int64_t bytes, int precision, int format) {
    return new QString(self->formattedDataSize(static_cast<qint64>(bytes), precision, DataFmt(format)));
}

/* =========== Quotation =========== */

QString* qlocale_quote_string(const go_qlocale_t self, const char* str, ptrdiff_t len, int style) {
    return new QString(self->quoteString(U8(str, len), QuoteStyle(style)));
}

/* =========== Separated list =========== */

QString* qlocale_create_separated_list(const go_qlocale_t self, const QString** strings, ptrdiff_t count) {
    QStringList list;
    for (ptrdiff_t i = 0; i < count; ++i) if (strings[i]) list.append(*strings[i]);
    return new QString(self->createSeparatedList(list));
}

/* =========== Collation =========== */

go_qlocale_t qlocale_collation(const go_qlocale_t self) {
    return new QLocale(self->collation());
}

/* =========== UILanguages =========== */

QString** qlocale_ui_languages(const go_qlocale_t self, int separator, ptrdiff_t* out_count) {
    QStringList list = self->uiLanguages(TagSep(separator));
    const qsizetype n = list.size();
    *out_count = static_cast<ptrdiff_t>(n);
    QString** arr = new QString*[n > 0 ? static_cast<size_t>(n) : 1];
    for (qsizetype i = 0; i < n; ++i) arr[i] = new QString(list[i]);
    return arr;
}
void qlocale_free_string_array(QString** arr) { delete[] arr; }

/* =========== Weekdays =========== */

int* qlocale_weekdays(const go_qlocale_t self, ptrdiff_t* out_count) {
    QList<Qt::DayOfWeek> list = self->weekdays();
    const qsizetype n = list.size();
    *out_count = static_cast<ptrdiff_t>(n);
    int* arr = new int[n > 0 ? static_cast<size_t>(n) : 1];
    for (qsizetype i = 0; i < n; ++i) arr[i] = static_cast<int>(list[i]);
    return arr;
}
void qlocale_free_int_array(int* arr) { delete[] arr; }

/* =========== toString(date/time/datetime, ...) =========== */

QString* qlocale_to_string_date_fmt(const go_qlocale_t self, int64_t jd, int format) {
    return new QString(self->toString(JD(jd), Fmt(format)));
}
QString* qlocale_to_string_date_str(const go_qlocale_t self, int64_t jd,
                                    const char* fmt, ptrdiff_t fmt_len) {
    return new QString(self->toString(JD(jd), U8(fmt, fmt_len)));
}
QString* qlocale_to_string_date_fmt_cal(const go_qlocale_t self, int64_t jd,
                                        int format, go_qcalendar_t cal) {
    if (!cal.isValid()) return new QString(self->toString(JD(jd), Fmt(format)));
    return new QString(self->toString(JD(jd), Fmt(format), cal));
}
QString* qlocale_to_string_date_str_cal(const go_qlocale_t self, int64_t jd,
                                        const char* fmt, ptrdiff_t fmt_len,
                                        go_qcalendar_t cal) {
    if (!cal.isValid()) return new QString(self->toString(JD(jd), U8(fmt, fmt_len)));
    return new QString(self->toString(JD(jd), U8(fmt, fmt_len), cal));
}

QString* qlocale_to_string_time_fmt(const go_qlocale_t self, int32_t ms, int format) {
    return new QString(self->toString(MS(ms), Fmt(format)));
}
QString* qlocale_to_string_time_str(const go_qlocale_t self, int32_t ms,
                                    const char* fmt, ptrdiff_t fmt_len) {
    return new QString(self->toString(MS(ms), U8(fmt, fmt_len)));
}

QString* qlocale_to_string_datetime_fmt(const go_qlocale_t self, const go_qdatetime_t* dt, int format) {
    if (!dt) return new QString();
    return new QString(self->toString(*DT(dt), Fmt(format)));
}
QString* qlocale_to_string_datetime_str(const go_qlocale_t self, const go_qdatetime_t* dt,
                                        const char* fmt, ptrdiff_t fmt_len) {
    if (!dt) return new QString();
    return new QString(self->toString(*DT(dt), U8(fmt, fmt_len)));
}
QString* qlocale_to_string_datetime_fmt_cal(const go_qlocale_t self, const go_qdatetime_t* dt,
                                            int format, go_qcalendar_t cal) {
    if (!dt) return new QString();
    if (!cal.isValid()) return new QString(self->toString(*DT(dt), Fmt(format)));
    return new QString(self->toString(*DT(dt), Fmt(format), cal));
}
QString* qlocale_to_string_datetime_str_cal(const go_qlocale_t self, const go_qdatetime_t* dt,
                                            const char* fmt, ptrdiff_t fmt_len,
                                            go_qcalendar_t cal) {
    if (!dt) return new QString();
    if (!cal.isValid()) return new QString(self->toString(*DT(dt), U8(fmt, fmt_len)));
    return new QString(self->toString(*DT(dt), U8(fmt, fmt_len), cal));
}

QString* qlocale_date_format(const go_qlocale_t self, int format) {
    return new QString(self->dateFormat(Fmt(format)));
}
QString* qlocale_time_format(const go_qlocale_t self, int format) {
    return new QString(self->timeFormat(Fmt(format)));
}
QString* qlocale_datetime_format(const go_qlocale_t self, int format) {
    return new QString(self->dateTimeFormat(Fmt(format)));
}

/* =========== toDate/toTime/toDateTime =========== */

int64_t qlocale_to_date_fmt(const go_qlocale_t self, const char* str, ptrdiff_t len,
                            int format, int base_year) {
    return static_cast<int64_t>(self->toDate(U8(str, len), Fmt(format), base_year).toJulianDay());
}
int64_t qlocale_to_date_str(const go_qlocale_t self, const char* str, ptrdiff_t len,
                            const char* fmt, ptrdiff_t fmt_len, int base_year) {
    return static_cast<int64_t>(self->toDate(U8(str, len), U8(fmt, fmt_len), base_year).toJulianDay());
}
int64_t qlocale_to_date_fmt_cal(const go_qlocale_t self, const char* str, ptrdiff_t len,
                                int format, go_qcalendar_t cal, int base_year) {
    if (!cal.isValid()) return static_cast<int64_t>(self->toDate(U8(str, len), Fmt(format), base_year).toJulianDay());
    return static_cast<int64_t>(self->toDate(U8(str, len), Fmt(format), cal, base_year).toJulianDay());
}
int64_t qlocale_to_date_str_cal(const go_qlocale_t self, const char* str, ptrdiff_t len,
                                const char* fmt, ptrdiff_t fmt_len,
                                go_qcalendar_t cal, int base_year) {
    if (!cal.isValid()) return static_cast<int64_t>(self->toDate(U8(str, len), U8(fmt, fmt_len), base_year).toJulianDay());
    return static_cast<int64_t>(self->toDate(U8(str, len), U8(fmt, fmt_len), cal, base_year).toJulianDay());
}

int32_t qlocale_to_time_fmt(const go_qlocale_t self, const char* str, ptrdiff_t len, int format) {
    QTime t = self->toTime(U8(str, len), Fmt(format));
    return t.isValid() ? static_cast<int32_t>(t.msecsSinceStartOfDay()) : -1;
}
int32_t qlocale_to_time_str(const go_qlocale_t self, const char* str, ptrdiff_t len,
                            const char* fmt, ptrdiff_t fmt_len) {
    QTime t = self->toTime(U8(str, len), U8(fmt, fmt_len));
    return t.isValid() ? static_cast<int32_t>(t.msecsSinceStartOfDay()) : -1;
}

go_qdatetime_t* qlocale_to_datetime_fmt(const go_qlocale_t self, const char* str, ptrdiff_t len,
                                        int format, int base_year) {
    return new QDateTime(self->toDateTime(U8(str, len), Fmt(format), base_year));
}
go_qdatetime_t* qlocale_to_datetime_str(const go_qlocale_t self, const char* str, ptrdiff_t len,
                                        const char* fmt, ptrdiff_t fmt_len, int base_year) {
    return new QDateTime(self->toDateTime(U8(str, len), U8(fmt, fmt_len), base_year));
}
go_qdatetime_t* qlocale_to_datetime_fmt_cal(const go_qlocale_t self, const char* str, ptrdiff_t len,
                                            int format, go_qcalendar_t cal, int base_year) {
    if (!cal.isValid()) return new QDateTime(self->toDateTime(U8(str, len), Fmt(format), base_year));
    return new QDateTime(self->toDateTime(U8(str, len), Fmt(format), cal, base_year));
}
go_qdatetime_t* qlocale_to_datetime_str_cal(const go_qlocale_t self, const char* str, ptrdiff_t len,
                                            const char* fmt, ptrdiff_t fmt_len,
                                            go_qcalendar_t cal, int base_year) {
    if (!cal.isValid()) return new QDateTime(self->toDateTime(U8(str, len), U8(fmt, fmt_len), base_year));
    return new QDateTime(self->toDateTime(U8(str, len), U8(fmt, fmt_len), cal, base_year));
}

/* =========== Static conversion =========== */

QString* qlocale_language_to_code(int language, int codeTypes) {
    return new QString(QLocale::languageToCode(Lang(language), LangCodes(codeTypes)));
}
int qlocale_code_to_language(const char* code, ptrdiff_t len, int codeTypes) {
    return static_cast<int>(QLocale::codeToLanguage(U8(code, len), LangCodes(codeTypes)));
}
QString* qlocale_territory_to_code(int territory) {
    return new QString(QLocale::territoryToCode(Terr(territory)));
}
int qlocale_code_to_territory(const char* code, ptrdiff_t len) {
    return static_cast<int>(QLocale::codeToTerritory(U8(code, len)));
}
QString* qlocale_script_to_code(int script) {
    return new QString(QLocale::scriptToCode(Scr(script)));
}
int qlocale_code_to_script(const char* code, ptrdiff_t len) {
    return static_cast<int>(QLocale::codeToScript(U8(code, len)));
}
QString* qlocale_language_to_string(int language) {
    return new QString(QLocale::languageToString(Lang(language)));
}
QString* qlocale_territory_to_string(int territory) {
    return new QString(QLocale::territoryToString(Terr(territory)));
}
QString* qlocale_script_to_string(int script) {
    return new QString(QLocale::scriptToString(Scr(script)));
}

/* =========== Static default / system =========== */

void qlocale_set_default(const go_qlocale_t locale) {
    if (locale) QLocale::setDefault(*locale);
}
go_qlocale_t qlocale_c(void)      { return new QLocale(QLocale::c()); }
go_qlocale_t qlocale_system(void) { return new QLocale(QLocale::system()); }

/* =========== MatchingLocales =========== */

go_qlocale_t* qlocale_matching_locales(int language, int script, int territory, ptrdiff_t* out_count) {
    QList<QLocale> list = QLocale::matchingLocales(Lang(language), Scr(script), Terr(territory));
    const qsizetype n = list.size();
    *out_count = static_cast<ptrdiff_t>(n);
    go_qlocale_t* arr = new go_qlocale_t[n > 0 ? static_cast<size_t>(n) : 1];
    for (qsizetype i = 0; i < n; ++i) arr[i] = new QLocale(list[i]);
    return arr;
}
void qlocale_free_locale_array(go_qlocale_t* arr) { delete[] arr; }

} // extern "C"