/* qtext/qchar_c.h
 *
 * Single C ABI header for QLatin1Char and QChar, mirroring the layout of
 * Qt 6.11 QtCore's qchar.h (both classes in one header).
 *
 * Declarations only. Definitions live in qtext_cxx_impl.cpp.
 *
 * Why not `inline` definitions here: cgo compiles the Go preamble as C,
 * so it only sees the extern "C" declarations. The definitions must be
 * emitted by a real C++ translation unit (qtext_cxx_impl.cpp) — plain
 * non-inline `extern "C"` functions, so the linker is guaranteed to find
 * them. (Inline functions are only emitted when ODR-used; a TU that just
 * includes the header emits nothing and the link fails.)
 */
#ifndef QT6CORE_QCHAR_C_H
#define QT6CORE_QCHAR_C_H

#include <stdint.h>
#include <stdbool.h>
#include "qstring_c.h"  /* QString (opaque in C mode) */

#ifdef __cplusplus
extern "C" {
#endif

/* ================================================================
 * QLatin1Char
 * ================================================================ */

uint8_t  qlatin1char_to_latin1(uint8_t self);
uint16_t qlatin1char_unicode  (uint8_t self);

bool qlatin1char_compares_equal(uint8_t a, uint8_t b);
int  qlatin1char_compare        (uint8_t a, uint8_t b);

bool qlatin1char_compares_equal_byte(uint8_t self, uint8_t other);
int  qlatin1char_compare_byte        (uint8_t self, uint8_t other);

/* ================================================================
 * QChar — accessors
 * ================================================================ */

uint16_t qchar_unicode(uint16_t self);
uint8_t  qchar_to_latin1(uint16_t self);

/* ================================================================
 * QChar — instance predicates
 * ================================================================ */

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

/* ================================================================
 * QChar — instance classification
 * ================================================================ */

int      qchar_category(uint16_t self);
int      qchar_direction(uint16_t self);
int      qchar_joining_type(uint16_t self);
uint8_t  qchar_combining_class(uint16_t self);
uint16_t qchar_mirrored_char(uint16_t self);
int      qchar_decomposition_tag(uint16_t self);
int      qchar_digit_value(uint16_t self);
int      qchar_script(uint16_t self);
int      qchar_unicode_version(uint16_t self);

/* ================================================================
 * QChar — instance case conversion
 * ================================================================ */

uint16_t qchar_to_lower(uint16_t self);
uint16_t qchar_to_upper(uint16_t self);
uint16_t qchar_to_title(uint16_t self);
uint16_t qchar_to_case_folded(uint16_t self);

/* ================================================================
 * QChar — static helpers (UCS-4)
 * ================================================================ */

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

bool qchar_static_is_print(uint32_t ucs4);
bool qchar_static_is_space(uint32_t ucs4);
bool qchar_static_is_mark(uint32_t ucs4);
bool qchar_static_is_letter(uint32_t ucs4);
bool qchar_static_is_number(uint32_t ucs4);
bool qchar_static_is_letter_or_number(uint32_t ucs4);
bool qchar_static_is_digit(uint32_t ucs4);
bool qchar_static_is_lower(uint32_t ucs4);
bool qchar_static_is_upper(uint32_t ucs4);
bool qchar_static_is_title_case(uint32_t ucs4);
bool qchar_static_is_punct(uint32_t ucs4);
bool qchar_static_is_symbol(uint32_t ucs4);
bool qchar_static_is_non_character(uint32_t ucs4);
bool qchar_static_is_high_surrogate(uint32_t ucs4);
bool qchar_static_is_low_surrogate(uint32_t ucs4);
bool qchar_static_is_surrogate(uint32_t ucs4);
bool qchar_requires_surrogates(uint32_t ucs4);

/* ================================================================
 * QChar — surrogates
 * ================================================================ */

uint16_t qchar_high_surrogate(uint32_t ucs4);
uint16_t qchar_low_surrogate(uint32_t ucs4);
uint32_t qchar_surrogate_to_ucs4(uint16_t high, uint16_t low);

/* ================================================================
 * QChar — static constructors
 * ================================================================ */

uint16_t qchar_from_latin1(char c);
uint16_t qchar_from_ucs2(uint16_t ucs2);
/* Returns 1 for BMP (out[0] set, out[1]=0), 2 for supplementary planes. */
int      qchar_from_ucs4(uint32_t ucs4, uint16_t out[2]);

/* ================================================================
 * QChar — decomposition (returns new QString*, caller owns)
 * ================================================================ */

QString* qchar_decomposition(uint32_t ucs4);

#ifdef __cplusplus
} /* extern "C" */
#endif

#endif /* QT6CORE_QCHAR_C_H */