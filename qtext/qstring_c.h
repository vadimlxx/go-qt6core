/* C ABI for Qt 6.11 QtCore/QString. */
#ifndef GO_QT6CORE_QSTRING_C_H
#define GO_QT6CORE_QSTRING_C_H

#include <stddef.h>
#include <stdint.h>
#include <stdbool.h>

#ifdef __cplusplus
extern "C" {
#else
typedef struct QString QString;
#endif

typedef struct go_qt6core_qstring_s go_qt6core_qstring_t;

/* --- lifecycle --- */
go_qt6core_qstring_t* qstring_new_empty(void);
go_qt6core_qstring_t* qstring_new_copy(const go_qt6core_qstring_t *other);
go_qt6core_qstring_t* qstring_new_from_qstringview(const uint16_t *data, ptrdiff_t size);
void qstring_delete(go_qt6core_qstring_t *self);

/* --- basic accessors --- */
ptrdiff_t qstring_size(const go_qt6core_qstring_t *self);
bool      qstring_is_empty(const go_qt6core_qstring_t *self);
bool      qstring_is_null(const go_qt6core_qstring_t *self);
bool      qstring_is_detached(const go_qt6core_qstring_t *self);
void      qstring_detach(go_qt6core_qstring_t *self);
const uint16_t* qstring_utf16(const go_qt6core_qstring_t *self);

/* --- element access --- */
uint16_t qstring_at(const go_qt6core_qstring_t *self, ptrdiff_t i);
uint16_t qstring_front(const go_qt6core_qstring_t *self);
uint16_t qstring_back(const go_qt6core_qstring_t *self);

/* --- modifiers --- */
void qstring_clear(go_qt6core_qstring_t *self);
void qstring_resize(go_qt6core_qstring_t *self, ptrdiff_t n);
void qstring_resize_fill(go_qt6core_qstring_t *self, ptrdiff_t n, uint16_t c);
void qstring_reserve(go_qt6core_qstring_t *self, ptrdiff_t n);
void qstring_squeeze(go_qt6core_qstring_t *self);
void qstring_truncate(go_qt6core_qstring_t *self, ptrdiff_t pos);
void qstring_chop(go_qt6core_qstring_t *self, ptrdiff_t n);
void qstring_fill(go_qt6core_qstring_t *self, uint16_t c, ptrdiff_t n);

void qstring_append_qchar(go_qt6core_qstring_t *self, uint16_t c);
void qstring_append_qstring(go_qt6core_qstring_t *self, const go_qt6core_qstring_t *other);
void qstring_append_utf8(go_qt6core_qstring_t *self, const char *data, ptrdiff_t size);
void qstring_prepend_qchar(go_qt6core_qstring_t *self, uint16_t c);
void qstring_prepend_qstring(go_qt6core_qstring_t *self, const go_qt6core_qstring_t *other);
void qstring_prepend_utf8(go_qt6core_qstring_t *self, const char *data, ptrdiff_t size);
void qstring_insert_qchar(go_qt6core_qstring_t *self, ptrdiff_t i, uint16_t c);
void qstring_insert_qstring(go_qt6core_qstring_t *self, ptrdiff_t i, const go_qt6core_qstring_t *other);
void qstring_insert_utf8(go_qt6core_qstring_t *self, ptrdiff_t i, const char *data, ptrdiff_t size);

void qstring_remove(go_qt6core_qstring_t *self, ptrdiff_t i, ptrdiff_t n);
void qstring_remove_at(go_qt6core_qstring_t *self, ptrdiff_t i);
void qstring_remove_first(go_qt6core_qstring_t *self);
void qstring_remove_last(go_qt6core_qstring_t *self);
void qstring_replace_qchar_qchar(go_qt6core_qstring_t *self, uint16_t before, uint16_t after);
void qstring_replace_qchar_qstring(go_qt6core_qstring_t *self, uint16_t before, const go_qt6core_qstring_t *after);
void qstring_replace_qchar_utf8(go_qt6core_qstring_t *self, uint16_t before, const char *after, ptrdiff_t after_size);
void qstring_replace_qstring_qchar(go_qt6core_qstring_t *self, const go_qt6core_qstring_t *before, uint16_t after);
void qstring_replace_qstring_qstring(go_qt6core_qstring_t *self, const go_qt6core_qstring_t *before, const go_qt6core_qstring_t *after);
void qstring_replace_qstring_utf8(go_qt6core_qstring_t *self, const go_qt6core_qstring_t *before, const char *after, ptrdiff_t after_size);
void qstring_replace_utf8_qchar(go_qt6core_qstring_t *self, const char *before, ptrdiff_t before_size, uint16_t after);
void qstring_replace_utf8_qstring(go_qt6core_qstring_t *self, const char *before, ptrdiff_t before_size, const go_qt6core_qstring_t *after);
void qstring_replace_utf8_utf8(go_qt6core_qstring_t *self, const char *before, ptrdiff_t before_size, const char *after, ptrdiff_t after_size);
void qstring_replace_range_utf8(go_qt6core_qstring_t *self, ptrdiff_t i, ptrdiff_t n, const char *after, ptrdiff_t after_size);

/* --- slicing (returns new QString) --- */
go_qt6core_qstring_t* qstring_left(const go_qt6core_qstring_t *self, ptrdiff_t n);
go_qt6core_qstring_t* qstring_right(const go_qt6core_qstring_t *self, ptrdiff_t n);
go_qt6core_qstring_t* qstring_mid(const go_qt6core_qstring_t *self, ptrdiff_t pos, ptrdiff_t n);
go_qt6core_qstring_t* qstring_first(const go_qt6core_qstring_t *self, ptrdiff_t n);
go_qt6core_qstring_t* qstring_last(const go_qt6core_qstring_t *self, ptrdiff_t n);
go_qt6core_qstring_t* qstring_sliced(const go_qt6core_qstring_t *self, ptrdiff_t pos);
go_qt6core_qstring_t* qstring_sliced_n(const go_qt6core_qstring_t *self, ptrdiff_t pos, ptrdiff_t n);
go_qt6core_qstring_t* qstring_chopped(const go_qt6core_qstring_t *self, ptrdiff_t n);

/* --- search (cs: 0 = CaseSensitive, 1 = CaseInsensitive) --- */
ptrdiff_t qstring_index_of_qchar(const go_qt6core_qstring_t *self, uint16_t c, ptrdiff_t from, int cs);
ptrdiff_t qstring_index_of_qstring(const go_qt6core_qstring_t *self, const go_qt6core_qstring_t *other, ptrdiff_t from, int cs);
ptrdiff_t qstring_index_of_utf8(const go_qt6core_qstring_t *self, const char *data, ptrdiff_t size, ptrdiff_t from, int cs);
ptrdiff_t qstring_last_index_of_qchar(const go_qt6core_qstring_t *self, uint16_t c, ptrdiff_t from, int cs);
ptrdiff_t qstring_last_index_of_qstring(const go_qt6core_qstring_t *self, const go_qt6core_qstring_t *other, ptrdiff_t from, int cs);
ptrdiff_t qstring_last_index_of_utf8(const go_qt6core_qstring_t *self, const char *data, ptrdiff_t size, ptrdiff_t from, int cs);
bool      qstring_contains_qchar(const go_qt6core_qstring_t *self, uint16_t c, int cs);
bool      qstring_contains_qstring(const go_qt6core_qstring_t *self, const go_qt6core_qstring_t *other, int cs);
bool      qstring_contains_utf8(const go_qt6core_qstring_t *self, const char *data, ptrdiff_t size, int cs);
ptrdiff_t qstring_count_qchar(const go_qt6core_qstring_t *self, uint16_t c, int cs);
ptrdiff_t qstring_count_qstring(const go_qt6core_qstring_t *self, const go_qt6core_qstring_t *other, int cs);
ptrdiff_t qstring_count_utf8(const go_qt6core_qstring_t *self, const char *data, ptrdiff_t size, int cs);
bool      qstring_starts_with_qchar(const go_qt6core_qstring_t *self, uint16_t c, int cs);
bool      qstring_starts_with_qstring(const go_qt6core_qstring_t *self, const go_qt6core_qstring_t *other, int cs);
bool      qstring_starts_with_utf8(const go_qt6core_qstring_t *self, const char *data, ptrdiff_t size, int cs);
bool      qstring_ends_with_qchar(const go_qt6core_qstring_t *self, uint16_t c, int cs);
bool      qstring_ends_with_qstring(const go_qt6core_qstring_t *self, const go_qt6core_qstring_t *other, int cs);
bool      qstring_ends_with_utf8(const go_qt6core_qstring_t *self, const char *data, ptrdiff_t size, int cs);

/* --- case / trim --- */
go_qt6core_qstring_t* qstring_to_lower(const go_qt6core_qstring_t *self);
go_qt6core_qstring_t* qstring_to_upper(const go_qt6core_qstring_t *self);
go_qt6core_qstring_t* qstring_to_case_folded(const go_qt6core_qstring_t *self);
bool qstring_is_lower(const go_qt6core_qstring_t *self);
bool qstring_is_upper(const go_qt6core_qstring_t *self);
go_qt6core_qstring_t* qstring_trimmed(const go_qt6core_qstring_t *self);
go_qt6core_qstring_t* qstring_simplified(const go_qt6core_qstring_t *self);
go_qt6core_qstring_t* qstring_left_justified(const go_qt6core_qstring_t *self, ptrdiff_t width, uint16_t fill, bool truncate);
go_qt6core_qstring_t* qstring_right_justified(const go_qt6core_qstring_t *self, ptrdiff_t width, uint16_t fill, bool truncate);

/* --- compare --- */
int qstring_compare_qstring(const go_qt6core_qstring_t *self, const go_qt6core_qstring_t *other, int cs);
int qstring_compare_qchar(const go_qt6core_qstring_t *self, uint16_t c, int cs);
int qstring_compare_utf8(const go_qt6core_qstring_t *self, const char *data, ptrdiff_t size, int cs);

/* --- numeric --- */
int16_t  qstring_to_short(const go_qt6core_qstring_t *self, bool *ok, int base);
uint16_t qstring_to_ushort(const go_qt6core_qstring_t *self, bool *ok, int base);
int32_t  qstring_to_int(const go_qt6core_qstring_t *self, bool *ok, int base);
uint32_t qstring_to_uint(const go_qt6core_qstring_t *self, bool *ok, int base);
int64_t  qstring_to_long(const go_qt6core_qstring_t *self, bool *ok, int base);
uint64_t qstring_to_ulong(const go_qt6core_qstring_t *self, bool *ok, int base);
int64_t  qstring_to_longlong(const go_qt6core_qstring_t *self, bool *ok, int base);
uint64_t qstring_to_ulonglong(const go_qt6core_qstring_t *self, bool *ok, int base);
float    qstring_to_float(const go_qt6core_qstring_t *self, bool *ok);
double   qstring_to_double(const go_qt6core_qstring_t *self, bool *ok);
void qstring_set_num_ll(go_qt6core_qstring_t *self, int64_t n, int base);
void qstring_set_num_ull(go_qt6core_qstring_t *self, uint64_t n, int base);
void qstring_set_num_double(go_qt6core_qstring_t *self, double v, char format, int precision);

/* --- static --- */
go_qt6core_qstring_t* qstring_number_ll(int64_t n, int base);
go_qt6core_qstring_t* qstring_number_ull(uint64_t n, int base);
go_qt6core_qstring_t* qstring_number_double(double v, char format, int precision);
ptrdiff_t qstring_max_size(void);

/* --- predicates --- */
bool qstring_is_right_to_left(const go_qt6core_qstring_t *self);
bool qstring_is_valid_utf16(const go_qt6core_qstring_t *self);

/* --- conversion to QByteArray (returns new QByteArray* via out pointer) --- */
void qstring_to_utf8_qba(const go_qt6core_qstring_t *self, void **out);
void qstring_to_latin1_qba(const go_qt6core_qstring_t *self, void **out);
void qstring_to_local8bit_qba(const go_qt6core_qstring_t *self, void **out);

#ifdef __cplusplus
}
#endif

#endif
