/* C ABI for Qt 6.11 QtCore/QByteArray. */
#ifndef GO_QT6CORE_QBYTEARRAY_C_H
#define GO_QT6CORE_QBYTEARRAY_C_H

#include <stddef.h>    /* ptrdiff_t, size_t */
#include <stdint.h>    /* int16_t, int32_t, int64_t, uint16_t, ... */
#include <stdbool.h>   /* bool */

#ifdef __cplusplus
/*
 * In C++ we rely on the real class QByteArray from <QByteArray> (must be
 * included before this header in the .cpp). The elaborated-type-specifier
 * struct QByteArray below only *refers* to the existing class, provided
 * QByteArray lives in the global namespace (i.e. QT_NAMESPACE is empty,
 * which is the default). If you build Qt with a non-empty QT_NAMESPACE,
 * switch to `typedef void QByteArray;` and cast in the implementation.
 */
extern "C" {
#else
typedef struct QByteArray QByteArray;
#endif

/* --- lifecycle --- */
QByteArray* qbytearray_new_empty(void);
QByteArray* qbytearray_new_from_data(const char *data, ptrdiff_t size);
QByteArray* qbytearray_new_fill(ptrdiff_t size, char c);
QByteArray* qbytearray_new_copy(const QByteArray *other);
void        qbytearray_delete(QByteArray *self);

/* --- basic accessors --- */
ptrdiff_t   qbytearray_size(const QByteArray *self);
ptrdiff_t   qbytearray_capacity(const QByteArray *self);
ptrdiff_t   qbytearray_max_size(void);
bool        qbytearray_is_empty(const QByteArray *self);
bool        qbytearray_is_null(const QByteArray *self);
const char* qbytearray_data(const QByteArray *self);
bool        qbytearray_is_detached(const QByteArray *self);
void        qbytearray_detach(QByteArray *self);

/* --- modifiers --- */
void qbytearray_resize(QByteArray *self, ptrdiff_t size);
void qbytearray_resize_fill(QByteArray *self, ptrdiff_t size, char c);
void qbytearray_resize_for_overwrite(QByteArray *self, ptrdiff_t size);
void qbytearray_fill(QByteArray *self, char c, ptrdiff_t size);
void qbytearray_clear(QByteArray *self);
void qbytearray_reserve(QByteArray *self, ptrdiff_t size);
void qbytearray_squeeze(QByteArray *self);
void qbytearray_truncate(QByteArray *self, ptrdiff_t pos);
void qbytearray_chop(QByteArray *self, ptrdiff_t n);

/* --- element access --- */
char qbytearray_at(const QByteArray *self, ptrdiff_t i);
void qbytearray_set_at(QByteArray *self, ptrdiff_t i, char c);
char qbytearray_front(const QByteArray *self);
char qbytearray_back(const QByteArray *self);

/* --- search --- */
ptrdiff_t qbytearray_index_of_char(const QByteArray *self, char c, ptrdiff_t from);
ptrdiff_t qbytearray_index_of_data(const QByteArray *self, const char *data, ptrdiff_t size, ptrdiff_t from);
ptrdiff_t qbytearray_last_index_of_char(const QByteArray *self, char c, ptrdiff_t from);
ptrdiff_t qbytearray_last_index_of_data(const QByteArray *self, const char *data, ptrdiff_t size, ptrdiff_t from);
ptrdiff_t qbytearray_count_char(const QByteArray *self, char c);
ptrdiff_t qbytearray_count_data(const QByteArray *self, const char *data, ptrdiff_t size);
bool      qbytearray_contains_char(const QByteArray *self, char c);
bool      qbytearray_contains_data(const QByteArray *self, const char *data, ptrdiff_t size);
int       qbytearray_compare(const QByteArray *self, const char *data, ptrdiff_t size, int cs);

/* --- slicing (returns new QByteArray) --- */
QByteArray* qbytearray_left(const QByteArray *self, ptrdiff_t n);
QByteArray* qbytearray_right(const QByteArray *self, ptrdiff_t n);
QByteArray* qbytearray_mid(const QByteArray *self, ptrdiff_t i, ptrdiff_t len);
QByteArray* qbytearray_first(const QByteArray *self, ptrdiff_t n);
QByteArray* qbytearray_last(const QByteArray *self, ptrdiff_t n);
QByteArray* qbytearray_sliced(const QByteArray *self, ptrdiff_t pos);
QByteArray* qbytearray_sliced_n(const QByteArray *self, ptrdiff_t pos, ptrdiff_t n);
QByteArray* qbytearray_chopped(const QByteArray *self, ptrdiff_t n);

/* --- predicates --- */
bool qbytearray_starts_with_char(const QByteArray *self, char c);
bool qbytearray_starts_with_data(const QByteArray *self, const char *data, ptrdiff_t size);
bool qbytearray_ends_with_char(const QByteArray *self, char c);
bool qbytearray_ends_with_data(const QByteArray *self, const char *data, ptrdiff_t size);
bool qbytearray_is_upper(const QByteArray *self);
bool qbytearray_is_lower(const QByteArray *self);
bool qbytearray_is_valid_utf8(const QByteArray *self);

/* --- transforms --- */
QByteArray* qbytearray_to_lower(const QByteArray *self);
QByteArray* qbytearray_to_upper(const QByteArray *self);
QByteArray* qbytearray_trimmed(const QByteArray *self);
QByteArray* qbytearray_simplified(const QByteArray *self);
QByteArray* qbytearray_left_justified(const QByteArray *self, ptrdiff_t width, char fill, bool truncate);
QByteArray* qbytearray_right_justified(const QByteArray *self, ptrdiff_t width, char fill, bool truncate);
QByteArray* qbytearray_repeated(const QByteArray *self, ptrdiff_t times);

/* --- modification: append/prepend/insert --- */
void qbytearray_append_char(QByteArray *self, char c);
void qbytearray_append_data(QByteArray *self, const char *data, ptrdiff_t size);
void qbytearray_append_count(QByteArray *self, ptrdiff_t count, char c);
void qbytearray_append_ba(QByteArray *self, const QByteArray *other);
void qbytearray_prepend_char(QByteArray *self, char c);
void qbytearray_prepend_data(QByteArray *self, const char *data, ptrdiff_t size);
void qbytearray_prepend_count(QByteArray *self, ptrdiff_t count, char c);
void qbytearray_prepend_ba(QByteArray *self, const QByteArray *other);
void qbytearray_insert_data(QByteArray *self, ptrdiff_t i, const char *data, ptrdiff_t size);
void qbytearray_insert_char(QByteArray *self, ptrdiff_t i, char c);
void qbytearray_insert_count(QByteArray *self, ptrdiff_t i, ptrdiff_t count, char c);

/* --- modification: remove/replace --- */
void qbytearray_remove(QByteArray *self, ptrdiff_t i, ptrdiff_t len);
void qbytearray_remove_at(QByteArray *self, ptrdiff_t pos);
void qbytearray_remove_first(QByteArray *self);
void qbytearray_remove_last(QByteArray *self);
void qbytearray_replace_char_char(QByteArray *self, char before, char after);
void qbytearray_replace_range(QByteArray *self, ptrdiff_t i, ptrdiff_t len,
                              const char *data, ptrdiff_t dsize);
void qbytearray_replace_data_data(QByteArray *self,
    const char *before, ptrdiff_t bsize, const char *after, ptrdiff_t asize);

/* --- numeric conversion --- */
int16_t  qbytearray_to_short(const QByteArray *self, bool *ok, int base);
uint16_t qbytearray_to_ushort(const QByteArray *self, bool *ok, int base);
int32_t  qbytearray_to_int(const QByteArray *self, bool *ok, int base);
uint32_t qbytearray_to_uint(const QByteArray *self, bool *ok, int base);
int64_t  qbytearray_to_longlong(const QByteArray *self, bool *ok, int base);
uint64_t qbytearray_to_ulonglong(const QByteArray *self, bool *ok, int base);
float    qbytearray_to_float(const QByteArray *self, bool *ok);
double   qbytearray_to_double(const QByteArray *self, bool *ok);

/* --- static number constructors --- */
QByteArray* qbytearray_number_ll(int64_t n, int base);
QByteArray* qbytearray_number_ull(uint64_t n, int base);
QByteArray* qbytearray_number_double(double n, char format, int precision);
void        qbytearray_set_num_ll(QByteArray *self, int64_t n, int base);
void        qbytearray_set_num_ull(QByteArray *self, uint64_t n, int base);
void        qbytearray_set_num_double(QByteArray *self, double n, char format, int precision);

/* --- encoding / decoding --- */
QByteArray* qbytearray_to_base64(const QByteArray *self, int options);
QByteArray* qbytearray_from_base64(const QByteArray *self, int options);
/* Returns new QByteArray and writes status
 * (0=Ok, 1=IllegalInputLength, 2=IllegalCharacter, 3=IllegalPadding). */
QByteArray* qbytearray_from_base64_encoding(const char *data, ptrdiff_t size, int options, int *out_status);
QByteArray* qbytearray_to_hex(const QByteArray *self, char separator);
QByteArray* qbytearray_from_hex(const QByteArray *self);
QByteArray* qbytearray_to_percent(const QByteArray *self,
    const QByteArray *exclude, const QByteArray *include, char percent);
QByteArray* qbytearray_from_percent(const QByteArray *self, char percent);

/* --- compression (free functions) --- */
QByteArray* qbytearray_compress(const QByteArray *self, int level);
QByteArray* qbytearray_uncompress(const QByteArray *self);

/* --- split --- */
QByteArray** qbytearray_split(const QByteArray *self, char sep, ptrdiff_t *out_len);
void         qbytearray_split_free_array(QByteArray **arr);

/* --- null-termination --- */
void        qbytearray_null_terminate(QByteArray *self);
QByteArray* qbytearray_null_terminated(const QByteArray *self);

/* --- sharing inspection --- */
bool qbytearray_is_shared_with(const QByteArray *self, const QByteArray *other);

#ifdef __cplusplus
}
#endif

#endif
