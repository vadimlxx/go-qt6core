/* C ABI for Qt 6.11 QtCore/QStringView.
 *
 * QStringView is a trivial value type; we pass it across the ABI as
 * (const uint16_t *data, ptrdiff_t size). Functions that return a new
 * view write through (const uint16_t **out_data, ptrdiff_t *out_size).
 */
#ifndef GO_QT6CORE_QSTRINGVIEW_C_H
#define GO_QT6CORE_QSTRINGVIEW_C_H

#include <stddef.h>
#include <stdint.h>
#include <stdbool.h>

#ifdef __cplusplus
extern "C" {
#endif

/* --- accessors --- */
uint16_t qstringview_at    (const uint16_t *data, ptrdiff_t size, ptrdiff_t i);
uint16_t qstringview_front (const uint16_t *data, ptrdiff_t size);
uint16_t qstringview_back  (const uint16_t *data, ptrdiff_t size);

/* --- predicates --- */
bool qstringview_is_null         (const uint16_t *data, ptrdiff_t size);
bool qstringview_is_empty        (const uint16_t *data, ptrdiff_t size);
bool qstringview_is_upper        (const uint16_t *data, ptrdiff_t size);
bool qstringview_is_lower        (const uint16_t *data, ptrdiff_t size);
bool qstringview_is_right_to_left(const uint16_t *data, ptrdiff_t size);
bool qstringview_is_valid_utf16  (const uint16_t *data, ptrdiff_t size);

/* --- slicing (new view via out params) --- */
void qstringview_left   (const uint16_t *data, ptrdiff_t size, ptrdiff_t n,
                         const uint16_t **out_data, ptrdiff_t *out_size);
void qstringview_right  (const uint16_t *data, ptrdiff_t size, ptrdiff_t n,
                         const uint16_t **out_data, ptrdiff_t *out_size);
void qstringview_mid    (const uint16_t *data, ptrdiff_t size, ptrdiff_t pos, ptrdiff_t n,
                         const uint16_t **out_data, ptrdiff_t *out_size);
void qstringview_first  (const uint16_t *data, ptrdiff_t size, ptrdiff_t n,
                         const uint16_t **out_data, ptrdiff_t *out_size);
void qstringview_last   (const uint16_t *data, ptrdiff_t size, ptrdiff_t n,
                         const uint16_t **out_data, ptrdiff_t *out_size);
void qstringview_sliced (const uint16_t *data, ptrdiff_t size, ptrdiff_t pos,
                         const uint16_t **out_data, ptrdiff_t *out_size);
void qstringview_sliced_n(const uint16_t *data, ptrdiff_t size, ptrdiff_t pos, ptrdiff_t n,
                          const uint16_t **out_data, ptrdiff_t *out_size);
void qstringview_chopped(const uint16_t *data, ptrdiff_t size, ptrdiff_t n,
                         const uint16_t **out_data, ptrdiff_t *out_size);
void qstringview_trimmed(const uint16_t *data, ptrdiff_t size,
                         const uint16_t **out_data, ptrdiff_t *out_size);

/* --- search (cs: 0=CaseSensitive, 1=CaseInsensitive) --- */
ptrdiff_t qstringview_index_of_qchar(const uint16_t *data, ptrdiff_t size,
                                     uint16_t c, ptrdiff_t from, int cs);
ptrdiff_t qstringview_index_of_view (const uint16_t *data, ptrdiff_t size,
                                     const uint16_t *needle, ptrdiff_t needle_size,
                                     ptrdiff_t from, int cs);
ptrdiff_t qstringview_last_index_of_qchar(const uint16_t *data, ptrdiff_t size,
                                          uint16_t c, ptrdiff_t from, int cs);
ptrdiff_t qstringview_last_index_of_view (const uint16_t *data, ptrdiff_t size,
                                          const uint16_t *needle, ptrdiff_t needle_size,
                                          ptrdiff_t from, int cs);
ptrdiff_t qstringview_count_qchar(const uint16_t *data, ptrdiff_t size,
                                  uint16_t c, int cs);
ptrdiff_t qstringview_count_view (const uint16_t *data, ptrdiff_t size,
                                  const uint16_t *needle, ptrdiff_t needle_size, int cs);
bool qstringview_contains_qchar(const uint16_t *data, ptrdiff_t size,
                                uint16_t c, int cs);
bool qstringview_contains_view (const uint16_t *data, ptrdiff_t size,
                                const uint16_t *needle, ptrdiff_t needle_size, int cs);
bool qstringview_starts_with_qchar(const uint16_t *data, ptrdiff_t size,
                                   uint16_t c, int cs);
bool qstringview_starts_with_view (const uint16_t *data, ptrdiff_t size,
                                   const uint16_t *prefix, ptrdiff_t prefix_size, int cs);
bool qstringview_ends_with_qchar(const uint16_t *data, ptrdiff_t size,
                                 uint16_t c, int cs);
bool qstringview_ends_with_view (const uint16_t *data, ptrdiff_t size,
                                 const uint16_t *suffix, ptrdiff_t suffix_size, int cs);

/* --- compare --- */
int qstringview_compare_view (const uint16_t *data, ptrdiff_t size,
                              const uint16_t *other, ptrdiff_t other_size, int cs);
int qstringview_compare_qchar(const uint16_t *data, ptrdiff_t size,
                              uint16_t c, int cs);

/* --- conversions (return new C++ QByteArray* via out pointer) --- */
void qstringview_to_utf8     (const uint16_t *data, ptrdiff_t size, void **out_qbytearray);
void qstringview_to_latin1   (const uint16_t *data, ptrdiff_t size, void **out_qbytearray);
void qstringview_to_local8bit(const uint16_t *data, ptrdiff_t size, void **out_qbytearray);

/* --- numeric conversion --- */
int16_t  qstringview_to_short    (const uint16_t *data, ptrdiff_t size, bool *ok, int base);
uint16_t qstringview_to_ushort   (const uint16_t *data, ptrdiff_t size, bool *ok, int base);
int32_t  qstringview_to_int      (const uint16_t *data, ptrdiff_t size, bool *ok, int base);
uint32_t qstringview_to_uint     (const uint16_t *data, ptrdiff_t size, bool *ok, int base);
int64_t  qstringview_to_long     (const uint16_t *data, ptrdiff_t size, bool *ok, int base);
uint64_t qstringview_to_ulong    (const uint16_t *data, ptrdiff_t size, bool *ok, int base);
int64_t  qstringview_to_longlong (const uint16_t *data, ptrdiff_t size, bool *ok, int base);
uint64_t qstringview_to_ulonglong(const uint16_t *data, ptrdiff_t size, bool *ok, int base);
float    qstringview_to_float    (const uint16_t *data, ptrdiff_t size, bool *ok);
double   qstringview_to_double   (const uint16_t *data, ptrdiff_t size, bool *ok);

/* --- static --- */
ptrdiff_t qstringview_max_size(void);

#ifdef __cplusplus
}
#endif

#endif