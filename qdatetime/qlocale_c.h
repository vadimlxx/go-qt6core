/* C ABI for Qt 6.11 QtCore/QLocale. */
#ifndef GO_QT6CORE_QDATETIME_QLOCALE_C_H
#define GO_QT6CORE_QDATETIME_QLOCALE_C_H

#include <stddef.h>
#include <stdint.h>
#include <stdbool.h>
#include "../qtext/qstring_c.h"   /* QString* (opaque in C mode) */
#include "qdate_c.h"
#include "qcalendar_c.h"

/* go_qlocale_t: handle (heap QLocale*). We own this typedef. */
#ifdef __cplusplus
class QLocale;
class QDateTime;
typedef QLocale*  go_qlocale_t;
typedef QDateTime go_qdatetime_t;   /* same typedef as qdatetime_c.h / qtimezone_c.h */
extern "C" {
#else
typedef struct go_qlocale_s go_qlocale_s;
typedef go_qlocale_s*      go_qlocale_t;
typedef struct QDateTime   go_qdatetime_t;
#endif

// #ifdef __cplusplus
// class QLocale;
// class QDateTime;
// typedef QLocale*    go_qlocale_t;
// typedef QDateTime   go_qdatetime_t;
// #else
// typedef struct go_qlocale_s   go_qlocale_s;
// typedef struct QDateTime      QDateTime;     /* opaque */
// typedef go_qlocale_s*         go_qlocale_t;
// typedef QDateTime             go_qdatetime_t;
// #endif

// #ifdef __cplusplus
// extern "C" {
// #endif


/* ---- Lifecycle ---- */
go_qlocale_t qlocale_new_default(void);
go_qlocale_t qlocale_new_from_name(const char* name, ptrdiff_t len);
go_qlocale_t qlocale_new_from_language_territory(int language, int territory);
go_qlocale_t qlocale_new_from_language_script_territory(int language, int script, int territory);
go_qlocale_t qlocale_new_copy(const go_qlocale_t other);
void         qlocale_delete(go_qlocale_t self);
void         qlocale_swap(go_qlocale_t self, go_qlocale_t other);

/* ---- Value accessors ---- */
int qlocale_language(const go_qlocale_t self);
int qlocale_script(const go_qlocale_t self);
int qlocale_territory(const go_qlocale_t self);
int qlocale_measurement_system(const go_qlocale_t self);
int qlocale_text_direction(const go_qlocale_t self);
int qlocale_first_day_of_week(const go_qlocale_t self);
int qlocale_number_options(const go_qlocale_t self);
void qlocale_set_number_options(go_qlocale_t self, int options);

/* ---- String representation (return new QString*; caller owns) ---- */
QString* qlocale_name(const go_qlocale_t self, int separator);
QString* qlocale_bcp47_name(const go_qlocale_t self, int separator);
QString* qlocale_native_language_name(const go_qlocale_t self);
QString* qlocale_native_territory_name(const go_qlocale_t self);

/* ---- Number parsing ---- */
bool qlocale_to_short(const go_qlocale_t self, const char* str, ptrdiff_t len, int16_t* out);
bool qlocale_to_ushort(const go_qlocale_t self, const char* str, ptrdiff_t len, uint16_t* out);
bool qlocale_to_int(const go_qlocale_t self, const char* str, ptrdiff_t len, int32_t* out);
bool qlocale_to_uint(const go_qlocale_t self, const char* str, ptrdiff_t len, uint32_t* out);
bool qlocale_to_long(const go_qlocale_t self, const char* str, ptrdiff_t len, int64_t* out);
bool qlocale_to_ulong(const go_qlocale_t self, const char* str, ptrdiff_t len, uint64_t* out);
bool qlocale_to_longlong(const go_qlocale_t self, const char* str, ptrdiff_t len, int64_t* out);
bool qlocale_to_ulonglong(const go_qlocale_t self, const char* str, ptrdiff_t len, uint64_t* out);
bool qlocale_to_float(const go_qlocale_t self, const char* str, ptrdiff_t len, float* out);
bool qlocale_to_double(const go_qlocale_t self, const char* str, ptrdiff_t len, double* out);

/* ---- Number formatting ---- */
QString* qlocale_to_string_ll(const go_qlocale_t self, int64_t value);
QString* qlocale_to_string_ull(const go_qlocale_t self, uint64_t value);
QString* qlocale_to_string_ll_w(const go_qlocale_t self, int64_t value, int fieldWidth, uint32_t fillChar);
QString* qlocale_to_string_ull_w(const go_qlocale_t self, uint64_t value, int fieldWidth, uint32_t fillChar);
QString* qlocale_to_string_double(const go_qlocale_t self, double value, char format, int precision);
QString* qlocale_to_string_float(const go_qlocale_t self, float value, char format, int precision);

/* ---- Locale-specific strings ---- */
QString* qlocale_decimal_point(const go_qlocale_t self);
QString* qlocale_group_separator(const go_qlocale_t self);
QString* qlocale_percent(const go_qlocale_t self);
QString* qlocale_zero_digit(const go_qlocale_t self);
QString* qlocale_negative_sign(const go_qlocale_t self);
QString* qlocale_positive_sign(const go_qlocale_t self);
QString* qlocale_exponential(const go_qlocale_t self);

/* ---- Month/day names ---- */
QString* qlocale_month_name(const go_qlocale_t self, int month, int format);
QString* qlocale_standalone_month_name(const go_qlocale_t self, int month, int format);
QString* qlocale_day_name(const go_qlocale_t self, int day, int format);
QString* qlocale_standalone_day_name(const go_qlocale_t self, int day, int format);

/* ---- AM/PM ---- */
QString* qlocale_am_text(const go_qlocale_t self);
QString* qlocale_pm_text(const go_qlocale_t self);

/* ---- Case conversion ---- */
QString* qlocale_to_upper(const go_qlocale_t self, const char* str, ptrdiff_t len);
QString* qlocale_to_lower(const go_qlocale_t self, const char* str, ptrdiff_t len);

/* ---- Currency ---- */
QString* qlocale_currency_symbol(const go_qlocale_t self, int format);
QString* qlocale_to_currency_string_ll(const go_qlocale_t self, int64_t value, const QString* symbol);
QString* qlocale_to_currency_string_ull(const go_qlocale_t self, uint64_t value, const QString* symbol);
QString* qlocale_to_currency_string_double(const go_qlocale_t self, double value, const QString* symbol, int precision);

/* ---- Data size ---- */
QString* qlocale_formatted_data_size(const go_qlocale_t self, int64_t bytes, int precision, int format);

/* ---- Quotation ---- */
QString* qlocale_quote_string(const go_qlocale_t self, const char* str, ptrdiff_t len, int style);

/* ---- Separated list (builds QStringList from array of QString*) ---- */
QString* qlocale_create_separated_list(const go_qlocale_t self, const QString** strings, ptrdiff_t count);

/* ---- Collation ---- */
go_qlocale_t qlocale_collation(const go_qlocale_t self);

/* ---- UILanguages (heap QString** array; free with qlocale_free_string_array) ---- */
QString** qlocale_ui_languages(const go_qlocale_t self, int separator, ptrdiff_t* out_count);
void      qlocale_free_string_array(QString** arr);

/* ---- Weekdays (heap int*; free with qlocale_free_int_array) ---- */
int* qlocale_weekdays(const go_qlocale_t self, ptrdiff_t* out_count);
void qlocale_free_int_array(int* arr);

/* =================================================================
 * Date/Time-aware API
 *
 * QDate passed as int64 julian day (QDate::toJulianDay()).
 * QTime passed as int32 msecs-since-start-of-day (QTime::msecsSinceStartOfDay()).
 * QDateTime passed as `const go_qdatetime_t* dt` (pointer to QDateTime).
 * QCalendar passed as `go_qcalendar_t cal` (value type; from qcalendar_c.h).
 * ================================================================= */

/* toString(QDate, FormatType) const */
QString* qlocale_to_string_date_fmt(const go_qlocale_t self, int64_t julian_day, int format);
/* toString(QDate, QStringView) const */
QString* qlocale_to_string_date_str(const go_qlocale_t self, int64_t julian_day,
                                    const char* fmt, ptrdiff_t fmt_len);
/* toString(QDate, FormatType, QCalendar) const */
QString* qlocale_to_string_date_fmt_cal(const go_qlocale_t self, int64_t julian_day,
                                        int format, go_qcalendar_t cal);
/* toString(QDate, QStringView, QCalendar) const */
QString* qlocale_to_string_date_str_cal(const go_qlocale_t self, int64_t julian_day,
                                        const char* fmt, ptrdiff_t fmt_len,
                                        go_qcalendar_t cal);

/* toString(QTime, FormatType) const */
QString* qlocale_to_string_time_fmt(const go_qlocale_t self, int32_t msecs, int format);
/* toString(QTime, QStringView) const */
QString* qlocale_to_string_time_str(const go_qlocale_t self, int32_t msecs,
                                    const char* fmt, ptrdiff_t fmt_len);

/* toString(QDateTime, FormatType) const */
QString* qlocale_to_string_datetime_fmt(const go_qlocale_t self, const go_qdatetime_t* dt, int format);
/* toString(QDateTime, QStringView) const */
QString* qlocale_to_string_datetime_str(const go_qlocale_t self, const go_qdatetime_t* dt,
                                        const char* fmt, ptrdiff_t fmt_len);
/* toString(QDateTime, FormatType, QCalendar) const */
QString* qlocale_to_string_datetime_fmt_cal(const go_qlocale_t self, const go_qdatetime_t* dt,
                                            int format, go_qcalendar_t cal);
/* toString(QDateTime, QStringView, QCalendar) const */
QString* qlocale_to_string_datetime_str_cal(const go_qlocale_t self, const go_qdatetime_t* dt,
                                            const char* fmt, ptrdiff_t fmt_len,
                                            go_qcalendar_t cal);

/* dateFormat / timeFormat / dateTimeFormat */
QString* qlocale_date_format(const go_qlocale_t self, int format);
QString* qlocale_time_format(const go_qlocale_t self, int format);
QString* qlocale_datetime_format(const go_qlocale_t self, int format);

/* toDate: returns julian day; QDate() (null) is represented as INT64_MIN */
int64_t qlocale_to_date_fmt(const go_qlocale_t self, const char* str, ptrdiff_t len,
                            int format, int base_year);
int64_t qlocale_to_date_str(const go_qlocale_t self, const char* str, ptrdiff_t len,
                            const char* fmt, ptrdiff_t fmt_len, int base_year);
int64_t qlocale_to_date_fmt_cal(const go_qlocale_t self, const char* str, ptrdiff_t len,
                                int format, go_qcalendar_t cal, int base_year);
int64_t qlocale_to_date_str_cal(const go_qlocale_t self, const char* str, ptrdiff_t len,
                                const char* fmt, ptrdiff_t fmt_len,
                                go_qcalendar_t cal, int base_year);

/* toTime: returns msecs since start-of-day; QTime() (null) is represented as -1 */
int32_t qlocale_to_time_fmt(const go_qlocale_t self, const char* str, ptrdiff_t len, int format);
int32_t qlocale_to_time_str(const go_qlocale_t self, const char* str, ptrdiff_t len,
                            const char* fmt, ptrdiff_t fmt_len);

/* toDateTime: returns new go_qdatetime_t* (caller owns) or NULL on allocation failure */
go_qdatetime_t* qlocale_to_datetime_fmt(const go_qlocale_t self, const char* str, ptrdiff_t len,
                                        int format, int base_year);
go_qdatetime_t* qlocale_to_datetime_str(const go_qlocale_t self, const char* str, ptrdiff_t len,
                                        const char* fmt, ptrdiff_t fmt_len, int base_year);
go_qdatetime_t* qlocale_to_datetime_fmt_cal(const go_qlocale_t self, const char* str, ptrdiff_t len,
                                            int format, go_qcalendar_t cal, int base_year);
go_qdatetime_t* qlocale_to_datetime_str_cal(const go_qlocale_t self, const char* str, ptrdiff_t len,
                                            const char* fmt, ptrdiff_t fmt_len,
                                            go_qcalendar_t cal, int base_year);

/* ---- Static conversion methods ---- */
QString* qlocale_language_to_code(int language, int codeTypes);
int      qlocale_code_to_language(const char* code, ptrdiff_t len, int codeTypes);
QString* qlocale_territory_to_code(int territory);
int      qlocale_code_to_territory(const char* code, ptrdiff_t len);
QString* qlocale_script_to_code(int script);
int      qlocale_code_to_script(const char* code, ptrdiff_t len);
QString* qlocale_language_to_string(int language);
QString* qlocale_territory_to_string(int territory);
QString* qlocale_script_to_string(int script);

/* ---- Static default / system ---- */
void         qlocale_set_default(const go_qlocale_t locale);
go_qlocale_t qlocale_c(void);
go_qlocale_t qlocale_system(void);

/* ---- MatchingLocales (heap go_qlocale_t* array; free with qlocale_free_locale_array) ---- */
go_qlocale_t* qlocale_matching_locales(int language, int script, int territory, ptrdiff_t* out_count);
void          qlocale_free_locale_array(go_qlocale_t* arr);

/* ---- CountriesForLanguage (heap int*; free with qlocale_free_int_array) ---- */
int* qlocale_countries_for_language(int language, ptrdiff_t* out_count);

#ifdef __cplusplus
}
#endif
#endif
