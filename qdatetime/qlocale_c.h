/* C ABI for Qt 6.11 QtCore/QLocale.
 *
 * QLocale is exposed as a HANDLE (QLocale*). Each handle owns a heap-allocated C++ QLocale.
 * All accessors and mutators go directly to the C++ object.
 *
 * String outputs use void** out-parameter; the Go side wraps the result via
 * qtext.QStringFromUnsafePtr. This avoids including qtext headers here.
 */
#ifndef GO_QT6CORE_QDATETIME_QLOCALE_C_H
#define GO_QT6CORE_QDATETIME_QLOCALE_C_H

#include <stddef.h>
#include <stdint.h>
#include <stdbool.h>

#ifdef __cplusplus
class QLocale;
typedef QLocale* go_qlocale_t;
extern "C" {
#else
typedef struct go_qlocale_s go_qlocale_s;
typedef go_qlocale_s* go_qlocale_t;
#endif

/* --- Lifecycle --- */
go_qlocale_t qlocale_new_default(void);
go_qlocale_t qlocale_new_from_name(const char* name, ptrdiff_t len);
go_qlocale_t qlocale_new_from_language_country(int language, int country);
go_qlocale_t qlocale_new_copy(const go_qlocale_t other);
go_qlocale_t qlocale_clone(const go_qlocale_t self);
void         qlocale_delete(go_qlocale_t self);

/* --- Accessors (values) --- */
int qlocale_language(const go_qlocale_t self);
int qlocale_script(const go_qlocale_t self);
int qlocale_country(const go_qlocale_t self);
int qlocale_measurement_system(const go_qlocale_t self);
int qlocale_text_direction(const go_qlocale_t self);

/* --- String representation (out_qstr = new QString*) --- */
void qlocale_name(const go_qlocale_t self, void** out_qstr);
void qlocale_bcp47_name(const go_qlocale_t self, void** out_qstr);
void qlocale_native_language_name(const go_qlocale_t self, void** out_qstr);
void qlocale_native_territory_name(const go_qlocale_t self, void** out_qstr);

/* --- Number parsing --- */
bool qlocale_to_int64(const go_qlocale_t self, const char* str, ptrdiff_t len, int64_t* out_value);
bool qlocale_to_double(const go_qlocale_t self, const char* str, ptrdiff_t len, double* out_value);

/* --- Number formatting --- */
void qlocale_to_string_int64(const go_qlocale_t self, int64_t value, void** out_qstr);
void qlocale_to_string_uint64(const go_qlocale_t self, uint64_t value, void** out_qstr);
void qlocale_to_string_double(const go_qlocale_t self, double value, char format, int precision, void** out_qstr);

/* --- Locale-specific strings --- */
void qlocale_decimal_point(const go_qlocale_t self, void** out_qstr);
void qlocale_group_separator(const go_qlocale_t self, void** out_qstr);
void qlocale_percent(const go_qlocale_t self, void** out_qstr);
void qlocale_zero_digit(const go_qlocale_t self, void** out_qstr);
void qlocale_negative_sign(const go_qlocale_t self, void** out_qstr);
void qlocale_positive_sign(const go_qlocale_t self, void** out_qstr);
void qlocale_exponential(const go_qlocale_t self, void** out_qstr);

/* --- Month / day names --- */
void qlocale_month_name(const go_qlocale_t self, int month, int format, void** out_qstr);
void qlocale_standalone_month_name(const go_qlocale_t self, int month, int format, void** out_qstr);
void qlocale_day_name(const go_qlocale_t self, int day, int format, void** out_qstr);
void qlocale_standalone_day_name(const go_qlocale_t self, int day, int format, void** out_qstr);

/* --- AM/PM --- */
void qlocale_am_text(const go_qlocale_t self, void** out_qstr);
void qlocale_pm_text(const go_qlocale_t self, void** out_qstr);

/* --- Case conversion --- */
void qlocale_to_upper(const go_qlocale_t self, const char* str, ptrdiff_t len, void** out_qstr);
void qlocale_to_lower(const go_qlocale_t self, const char* str, ptrdiff_t len, void** out_qstr);

/* --- Static methods --- */
void qlocale_set_default(const go_qlocale_t locale);
go_qlocale_t qlocale_c(void);
go_qlocale_t qlocale_system(void);

#ifdef __cplusplus
}
#endif
#endif