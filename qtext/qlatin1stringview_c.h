/* C ABI for Qt 6.11 QtCore/QLatin1StringView.
 *
 * QLatin1StringView is a trivial value type; we pass it across the ABI
 * as (const char *data, ptrdiff_t size). Views returned by value are
 * written through (const char **out_data, ptrdiff_t *out_size) out params.
 */
#ifndef GO_QT6CORE_QLATIN1STRINGVIEW_C_H
#define GO_QT6CORE_QLATIN1STRINGVIEW_C_H

#include <stddef.h>
#include <stdint.h>
#include <stdbool.h>
#include "qbytearray_c.h"

#ifdef __cplusplus
extern "C" {
#endif

/* --- accessors --- */
uint8_t qlatin1stringview_at    (const char *data, ptrdiff_t size, ptrdiff_t i);
uint8_t qlatin1stringview_front (const char *data, ptrdiff_t size);
uint8_t qlatin1stringview_back  (const char *data, ptrdiff_t size);

/* --- predicates --- */
bool qlatin1stringview_is_null (const char *data, ptrdiff_t size);
bool qlatin1stringview_is_empty(const char *data, ptrdiff_t size);

/* --- slicing (returns new view) --- */
void qlatin1stringview_mid     (const char *data, ptrdiff_t size, ptrdiff_t pos, ptrdiff_t n,
                                const char **out_data, ptrdiff_t *out_size);
void qlatin1stringview_left    (const char *data, ptrdiff_t size, ptrdiff_t n,
                                const char **out_data, ptrdiff_t *out_size);
void qlatin1stringview_right   (const char *data, ptrdiff_t size, ptrdiff_t n,
                                const char **out_data, ptrdiff_t *out_size);
void qlatin1stringview_sliced  (const char *data, ptrdiff_t size, ptrdiff_t pos,
                                const char **out_data, ptrdiff_t *out_size);
void qlatin1stringview_sliced_n(const char *data, ptrdiff_t size, ptrdiff_t pos, ptrdiff_t n,
                                const char **out_data, ptrdiff_t *out_size);
void qlatin1stringview_first   (const char *data, ptrdiff_t size, ptrdiff_t n,
                                const char **out_data, ptrdiff_t *out_size);
void qlatin1stringview_last    (const char *data, ptrdiff_t size, ptrdiff_t n,
                                const char **out_data, ptrdiff_t *out_size);
void qlatin1stringview_chopped (const char *data, ptrdiff_t size, ptrdiff_t n,
                                const char **out_data, ptrdiff_t *out_size);
void qlatin1stringview_trimmed (const char *data, ptrdiff_t size,
                                const char **out_data, ptrdiff_t *out_size);

/* --- compare (cs: 0 = CaseSensitive, 1 = CaseInsensitive) --- */
int qlatin1stringview_compare_view        (const char *data, ptrdiff_t size,
                                           const char *other, ptrdiff_t other_size, int cs);
int qlatin1stringview_compare_qstringview (const char *data, ptrdiff_t size,
                                           const uint16_t *sv_data, ptrdiff_t sv_size, int cs);
int qlatin1stringview_compare_qchar       (const char *data, ptrdiff_t size,
                                           uint16_t c, int cs);

/* --- startsWith --- */
bool qlatin1stringview_starts_with_view        (const char *data, ptrdiff_t size,
                                                const char *other, ptrdiff_t other_size, int cs);
bool qlatin1stringview_starts_with_qstringview (const char *data, ptrdiff_t size,
                                                const uint16_t *sv_data, ptrdiff_t sv_size, int cs);
bool qlatin1stringview_starts_with_qchar       (const char *data, ptrdiff_t size,
                                                uint16_t c, int cs);

/* --- endsWith --- */
bool qlatin1stringview_ends_with_view        (const char *data, ptrdiff_t size,
                                              const char *other, ptrdiff_t other_size, int cs);
bool qlatin1stringview_ends_with_qstringview (const char *data, ptrdiff_t size,
                                              const uint16_t *sv_data, ptrdiff_t sv_size, int cs);
bool qlatin1stringview_ends_with_qchar       (const char *data, ptrdiff_t size,
                                              uint16_t c, int cs);

/* --- indexOf / lastIndexOf --- */
ptrdiff_t qlatin1stringview_index_of_view        (const char *data, ptrdiff_t size,
                                                  const char *other, ptrdiff_t other_size,
                                                  ptrdiff_t from, int cs);
ptrdiff_t qlatin1stringview_index_of_qstringview (const char *data, ptrdiff_t size,
                                                  const uint16_t *sv_data, ptrdiff_t sv_size,
                                                  ptrdiff_t from, int cs);
ptrdiff_t qlatin1stringview_index_of_qchar       (const char *data, ptrdiff_t size,
                                                  uint16_t c, ptrdiff_t from, int cs);
ptrdiff_t qlatin1stringview_last_index_of_view        (const char *data, ptrdiff_t size,
                                                       const char *other, ptrdiff_t other_size,
                                                       ptrdiff_t from, int cs);
ptrdiff_t qlatin1stringview_last_index_of_qstringview (const char *data, ptrdiff_t size,
                                                       const uint16_t *sv_data, ptrdiff_t sv_size,
                                                       ptrdiff_t from, int cs);
ptrdiff_t qlatin1stringview_last_index_of_qchar       (const char *data, ptrdiff_t size,
                                                       uint16_t c, ptrdiff_t from, int cs);

/* --- contains --- */
bool qlatin1stringview_contains_view        (const char *data, ptrdiff_t size,
                                             const char *other, ptrdiff_t other_size, int cs);
bool qlatin1stringview_contains_qstringview (const char *data, ptrdiff_t size,
                                             const uint16_t *sv_data, ptrdiff_t sv_size, int cs);
bool qlatin1stringview_contains_qchar       (const char *data, ptrdiff_t size,
                                             uint16_t c, int cs);

/* --- count --- */
ptrdiff_t qlatin1stringview_count_view        (const char *data, ptrdiff_t size,
                                               const char *other, ptrdiff_t other_size, int cs);
ptrdiff_t qlatin1stringview_count_qstringview (const char *data, ptrdiff_t size,
                                               const uint16_t *sv_data, ptrdiff_t sv_size, int cs);
ptrdiff_t qlatin1stringview_count_qchar       (const char *data, ptrdiff_t size,
                                               uint16_t c, int cs);

/* --- numeric conversion --- */
int16_t  qlatin1stringview_to_short    (const char *data, ptrdiff_t size, bool *ok, int base);
uint16_t qlatin1stringview_to_ushort   (const char *data, ptrdiff_t size, bool *ok, int base);
int32_t  qlatin1stringview_to_int      (const char *data, ptrdiff_t size, bool *ok, int base);
uint32_t qlatin1stringview_to_uint     (const char *data, ptrdiff_t size, bool *ok, int base);
int64_t  qlatin1stringview_to_long     (const char *data, ptrdiff_t size, bool *ok, int base);
uint64_t qlatin1stringview_to_ulong    (const char *data, ptrdiff_t size, bool *ok, int base);
int64_t  qlatin1stringview_to_longlong (const char *data, ptrdiff_t size, bool *ok, int base);
uint64_t qlatin1stringview_to_ulonglong(const char *data, ptrdiff_t size, bool *ok, int base);
float    qlatin1stringview_to_float    (const char *data, ptrdiff_t size, bool *ok);
double   qlatin1stringview_to_double   (const char *data, ptrdiff_t size, bool *ok);

/* --- encoding (returns new C++ QByteArray* via out pointer) --- */
QByteArray* qlatin1stringview_to_utf8(const char *data, ptrdiff_t size);

/* --- static --- */
ptrdiff_t qlatin1stringview_max_size(void);

#ifdef __cplusplus
}
#endif

#endif
