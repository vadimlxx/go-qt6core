/* C ABI for Qt 6.11 QtCore/QChar. */
#ifndef QT6CORE_QCHAR_C_H
#define QT6CORE_QCHAR_C_H

#include <stdint.h>
#include <stdbool.h>
#include "qstring_c.h"

#ifdef __cplusplus
extern "C" {
#endif

/* --- Constructors / accessors --- */
uint16_t qchar_unicode(uint16_t self);
uint8_t  qchar_to_latin1(uint16_t self);

/* --- Instance predicates --- */
bool qchar_is_null(uint16_t self);
bool qchar_is_print(uint16_t self);
bool qchar_is_space(uint16_t self);
bool qchar_is_mark(uint16_t self);
bool qchar_is_letter(uint16_t self);
bool qchar_is_number(uint16_t self);
bool qchar_is_letter_or_number(uint16_t self);
bool qchar_is_digit(uint16_t self);
bool qchar_is_lower(uint16_t self);
bool qchar_is_upper(uint16_t self);
bool qchar_is_title_case(uint16_t self);
bool qchar_is_punct(uint16_t self);
bool qchar_is_symbol(uint16_t self);
bool qchar_is_non_character(uint16_t self);
bool qchar_is_high_surrogate(uint16_t self);
bool qchar_is_low_surrogate(uint16_t self);
bool qchar_is_surrogate(uint16_t self);
bool qchar_has_mirrored(uint16_t self);

/* --- Instance classification --- */
int      qchar_category(uint16_t self);
int      qchar_direction(uint16_t self);
int      qchar_joining_type(uint16_t self);
uint8_t  qchar_combining_class(uint16_t self);
uint16_t qchar_mirrored_char(uint16_t self);
int      qchar_decomposition_tag(uint16_t self);
int      qchar_digit_value(uint16_t self);
int      qchar_script(uint16_t self);
int      qchar_unicode_version(uint16_t self);

/* --- Instance case conversion (return QChar, i.e. uint16) --- */
uint16_t qchar_to_lower(uint16_t self);
uint16_t qchar_to_upper(uint16_t self);
uint16_t qchar_to_title(uint16_t self);
uint16_t qchar_to_case_folded(uint16_t self);

/* --- Static helpers (UCS-4) --- */
int      qchar_static_category(uint32_t ucs4);
int      qchar_static_direction(uint32_t ucs4);
int      qchar_static_joining_type(uint32_t ucs4);
uint8_t  qchar_static_combining_class(uint32_t ucs4);
uint32_t qchar_static_mirrored_char(uint32_t ucs4);
bool     qchar_static_has_mirrored(uint32_t ucs4);
int      qchar_static_decomposition_tag(uint32_t ucs4);
int      qchar_static_digit_value(uint32_t ucs4);
uint32_t qchar_static_to_lower(uint32_t ucs4);
uint32_t qchar_static_to_upper(uint32_t ucs4);
uint32_t qchar_static_to_title(uint32_t ucs4);
uint32_t qchar_static_to_case_folded(uint32_t ucs4);
int      qchar_static_script(uint32_t ucs4);
int      qchar_static_unicode_version(uint32_t ucs4);
int      qchar_static_current_unicode_version(void);

bool     qchar_static_is_print(uint32_t ucs4);
bool     qchar_static_is_space(uint32_t ucs4);
bool     qchar_static_is_mark(uint32_t ucs4);
bool     qchar_static_is_letter(uint32_t ucs4);
bool     qchar_static_is_number(uint32_t ucs4);
bool     qchar_static_is_letter_or_number(uint32_t ucs4);
bool     qchar_static_is_digit(uint32_t ucs4);
bool     qchar_static_is_lower(uint32_t ucs4);
bool     qchar_static_is_upper(uint32_t ucs4);
bool     qchar_static_is_title_case(uint32_t ucs4);
bool     qchar_static_is_punct(uint32_t ucs4);
bool     qchar_static_is_symbol(uint32_t ucs4);
bool     qchar_static_is_non_character(uint32_t ucs4);
bool     qchar_static_is_high_surrogate(uint32_t ucs4);
bool     qchar_static_is_low_surrogate(uint32_t ucs4);
bool     qchar_static_is_surrogate(uint32_t ucs4);
bool     qchar_requires_surrogates(uint32_t ucs4);

/* --- Surrogates --- */
uint16_t qchar_high_surrogate(uint32_t ucs4);
uint16_t qchar_low_surrogate(uint32_t ucs4);
uint32_t qchar_surrogate_to_ucs4(uint16_t high, uint16_t low);

/* --- Static constructors --- */
uint16_t qchar_from_latin1(char c);
uint16_t qchar_from_ucs2(uint16_t ucs2);

/* QChar::fromUcs4 returns a small array (1 or 2 QChars). We expose both
 * code units; the second is 0 if the input fit into a single QChar. */
void qchar_from_ucs4(uint32_t ucs4, uint16_t out[2]);

/* --- Decomposition (returns new QString* via out pointer) --- */
QString* qchar_decomposition(uint32_t ucs4);

#ifdef __cplusplus
}
#endif

#endif
