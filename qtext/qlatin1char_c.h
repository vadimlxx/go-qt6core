/* C ABI for Qt 6.11 QtCore/QLatin1Char.
 *
 * QLatin1Char is a trivial value type: a single Latin-1 byte.
 * It is passed across the C ABI as uint8_t.
 */
#ifndef GO_QT6CORE_QLATIN1CHAR_C_H
#define GO_QT6CORE_QLATIN1CHAR_C_H

#include <stdint.h>
#include <stdbool.h>

#ifdef __cplusplus
extern "C" {
#endif

/* --- accessors --- */
uint8_t  qlatin1char_to_latin1(uint8_t self);
uint16_t qlatin1char_unicode  (uint8_t self);

/* --- comparisons (QLatin1Char vs QLatin1Char) --- */
bool qlatin1char_compares_equal(uint8_t a, uint8_t b);
int  qlatin1char_compare        (uint8_t a, uint8_t b); /* returns -1, 0, +1 */

/* --- comparisons (QLatin1Char vs raw byte) --- */
bool qlatin1char_compares_equal_byte(uint8_t self, uint8_t other);
int  qlatin1char_compare_byte        (uint8_t self, uint8_t other);

#ifdef __cplusplus
}
#endif

#endif