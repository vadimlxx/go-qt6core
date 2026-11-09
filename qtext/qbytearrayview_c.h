/* C ABI for Qt 6.11 QtCore/QByteArrayView.
 *
 * QByteArrayView is a trivial value type. We never pass it across the
 * ABI directly; instead, every function takes (data, size) and constructs
 * a QByteArrayView inside. Functions returning a new view write through
 * (const char **out_data, ptrdiff_t *out_size) out-parameters.
 */
#ifndef GO_QT6CORE_QBYTEARRAYVIEW_C_H
#define GO_QT6CORE_QBYTEARRAYVIEW_C_H

#include <stddef.h>
#include <stdint.h>
#include <stdbool.h>
#include "qbytearray_c.h"


#ifdef __cplusplus
extern "C" {
#endif

/* --- predicates --- */
bool qbytearrayview_is_null(const char *data, ptrdiff_t size);
bool qbytearrayview_is_empty(const char *data, ptrdiff_t size);
bool qbytearrayview_is_valid_utf8(const char *data, ptrdiff_t size);

/* --- element access --- */
char qbytearrayview_at(const char *data, ptrdiff_t size, ptrdiff_t i);
char qbytearrayview_front(const char *data, ptrdiff_t size);
char qbytearrayview_back(const char *data, ptrdiff_t size);

/* --- search --- */
ptrdiff_t qbytearrayview_index_of_char(const char *data, ptrdiff_t size,
                                       char c, ptrdiff_t from);
ptrdiff_t qbytearrayview_index_of_view(const char *data, ptrdiff_t size,
                                       const char *needle, ptrdiff_t needle_size,
                                       ptrdiff_t from);
ptrdiff_t qbytearrayview_last_index_of_char(const char *data, ptrdiff_t size,
                                            char c, ptrdiff_t from);
ptrdiff_t qbytearrayview_last_index_of_view(const char *data, ptrdiff_t size,
                                            const char *needle, ptrdiff_t needle_size,
                                            ptrdiff_t from);
ptrdiff_t qbytearrayview_count_char(const char *data, ptrdiff_t size, char c);
ptrdiff_t qbytearrayview_count_view(const char *data, ptrdiff_t size,
                                    const char *needle, ptrdiff_t needle_size);
bool qbytearrayview_contains_char(const char *data, ptrdiff_t size, char c);
bool qbytearrayview_contains_view(const char *data, ptrdiff_t size,
                                  const char *needle, ptrdiff_t needle_size);
bool qbytearrayview_starts_with_char(const char *data, ptrdiff_t size, char c);
bool qbytearrayview_starts_with_view(const char *data, ptrdiff_t size,
                                     const char *prefix, ptrdiff_t prefix_size);
bool qbytearrayview_ends_with_char(const char *data, ptrdiff_t size, char c);
bool qbytearrayview_ends_with_view(const char *data, ptrdiff_t size,
                                   const char *suffix, ptrdiff_t suffix_size);

/* --- compare --- */
int qbytearrayview_compare(const char *data, ptrdiff_t size,
                           const char *other, ptrdiff_t other_size,
                           int cs /* 0=CaseSensitive, 1=CaseInsensitive */);

/* --- slicing (new view via out-params) --- */
void qbytearrayview_left(const char *data, ptrdiff_t size, ptrdiff_t n,
                         const char **out_data, ptrdiff_t *out_size);
void qbytearrayview_right(const char *data, ptrdiff_t size, ptrdiff_t n,
                          const char **out_data, ptrdiff_t *out_size);
void qbytearrayview_mid(const char *data, ptrdiff_t size, ptrdiff_t pos, ptrdiff_t n,
                        const char **out_data, ptrdiff_t *out_size);
void qbytearrayview_first(const char *data, ptrdiff_t size, ptrdiff_t n,
                          const char **out_data, ptrdiff_t *out_size);
void qbytearrayview_last(const char *data, ptrdiff_t size, ptrdiff_t n,
                         const char **out_data, ptrdiff_t *out_size);
void qbytearrayview_sliced(const char *data, ptrdiff_t size, ptrdiff_t pos,
                           const char **out_data, ptrdiff_t *out_size);
void qbytearrayview_sliced_n(const char *data, ptrdiff_t size, ptrdiff_t pos, ptrdiff_t n,
                             const char **out_data, ptrdiff_t *out_size);
void qbytearrayview_chopped(const char *data, ptrdiff_t size, ptrdiff_t n,
                            const char **out_data, ptrdiff_t *out_size);
void qbytearrayview_trimmed(const char *data, ptrdiff_t size,
                            const char **out_data, ptrdiff_t *out_size);

/* --- numeric conversion --- */
int16_t  qbytearrayview_to_short(const char *data, ptrdiff_t size, bool *ok, int base);
uint16_t qbytearrayview_to_ushort(const char *data, ptrdiff_t size, bool *ok, int base);
int32_t  qbytearrayview_to_int(const char *data, ptrdiff_t size, bool *ok, int base);
uint32_t qbytearrayview_to_uint(const char *data, ptrdiff_t size, bool *ok, int base);
int64_t  qbytearrayview_to_longlong(const char *data, ptrdiff_t size, bool *ok, int base);
uint64_t qbytearrayview_to_ulonglong(const char *data, ptrdiff_t size, bool *ok, int base);
float    qbytearrayview_to_float(const char *data, ptrdiff_t size, bool *ok);
double   qbytearrayview_to_double(const char *data, ptrdiff_t size, bool *ok);

/* --- conversion to owning QByteArray (requires linking with qbytearray) --- */
QByteArray* qbytearrayview_to_bytearray(const char *data, ptrdiff_t size);

/* --- static --- */
ptrdiff_t qbytearrayview_max_size(void);

#ifdef __cplusplus
}
#endif

#endif
