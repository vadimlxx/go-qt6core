/* C ABI for Qt 6.11.2 QtCore/QBitArray (and QBitRef, which is a Go-side proxy). */
#ifndef GO_QT6CORE_QBITARRAY_C_H
#define GO_QT6CORE_QBITARRAY_C_H

#include <stddef.h>
#include <stdint.h>
#include <stdbool.h>

#ifdef __cplusplus
extern "C" {
#endif

typedef struct go_qbitarray_s go_qbitarray_t;

/* ---- Constructors ---- */
go_qbitarray_t *qbitarray_new_empty(void);
go_qbitarray_t *qbitarray_new_sized(ptrdiff_t size, bool val);
go_qbitarray_t *qbitarray_new_copy(const go_qbitarray_t *other);

/* ---- Lifecycle ---- */
void qbitarray_delete(go_qbitarray_t *self);
void qbitarray_assign(go_qbitarray_t *self, const go_qbitarray_t *other);
void qbitarray_swap(go_qbitarray_t *self, go_qbitarray_t *other);

/* ---- Size / status ---- */
ptrdiff_t qbitarray_size(const go_qbitarray_t *self);
ptrdiff_t qbitarray_count_set(const go_qbitarray_t *self, bool on);
bool      qbitarray_is_empty(const go_qbitarray_t *self);
bool      qbitarray_is_null(const go_qbitarray_t *self);

/* ---- Modifiers ---- */
void qbitarray_resize(go_qbitarray_t *self, ptrdiff_t size);
void qbitarray_detach(go_qbitarray_t *self);
bool qbitarray_is_detached(const go_qbitarray_t *self);
void qbitarray_clear(go_qbitarray_t *self);

/* ---- Bit operations ---- */
bool qbitarray_test_bit(const go_qbitarray_t *self, ptrdiff_t i);
void qbitarray_set_bit(go_qbitarray_t *self, ptrdiff_t i);
void qbitarray_set_bit_val(go_qbitarray_t *self, ptrdiff_t i, bool val);
void qbitarray_clear_bit(go_qbitarray_t *self, ptrdiff_t i);
bool qbitarray_toggle_bit(go_qbitarray_t *self, ptrdiff_t i);

/* ---- Bitwise assign (operator&=, |=, ^= with const QBitArray&) ---- */
void qbitarray_and_assign(go_qbitarray_t *self, const go_qbitarray_t *other);
void qbitarray_or_assign(go_qbitarray_t *self, const go_qbitarray_t *other);
void qbitarray_xor_assign(go_qbitarray_t *self, const go_qbitarray_t *other);

/* ---- Bitwise non-assign (operator&, |, ^, ~ returning new QBitArray) ---- */
go_qbitarray_t *qbitarray_op_and(const go_qbitarray_t *a, const go_qbitarray_t *b);
go_qbitarray_t *qbitarray_op_or (const go_qbitarray_t *a, const go_qbitarray_t *b);
go_qbitarray_t *qbitarray_op_xor(const go_qbitarray_t *a, const go_qbitarray_t *b);
go_qbitarray_t *qbitarray_inverted(const go_qbitarray_t *a);

/* ---- Fill / truncate ---- */
void qbitarray_fill_all(go_qbitarray_t *self, bool val);
void qbitarray_fill_sized(go_qbitarray_t *self, bool val, ptrdiff_t size);
void qbitarray_fill_range(go_qbitarray_t *self, bool val, ptrdiff_t first, ptrdiff_t last);
void qbitarray_truncate(go_qbitarray_t *self, ptrdiff_t pos);

/* ---- Raw bits (LSB-first within each byte; bit i = (data[i/8] >> (i%8)) & 1) ----
 *
 * qbitarray_bits: returns a borrowed pointer into the array's internal buffer, or NULL for
 * an empty array. *out_bytes receives (size()+7)/8. The caller must copy the bytes before
 * calling any mutating method on `self`.
 *
 * qbitarray_from_bits: constructs a new QBitArray from `bit_len` valid bits. The input buffer
 * must contain at least (bit_len+7)/8 bytes.
 */
const char* qbitarray_bits(const go_qbitarray_t *self, ptrdiff_t *out_bytes);
go_qbitarray_t *qbitarray_from_bits(const char *data, ptrdiff_t bit_len);

/* ---- toUInt32 (returns the first 32 bits interpreted per endianness) ----
 *
 * If the array has fewer than 32 bits, *ok is set to false; otherwise to true.
 */
uint32_t qbitarray_to_uint32(const go_qbitarray_t *self, bool littleEndian, bool *ok);

/* ---- Comparison (operator==) ---- */
bool qbitarray_equal(const go_qbitarray_t *self, const go_qbitarray_t *other);

#ifdef __cplusplus
}
#endif

#endif