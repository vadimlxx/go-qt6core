/* C ABI for Qt 6.11 QtCore/QString. */
#ifndef GO_QT6CORE_QSTRING_C_H
#define GO_QT6CORE_QSTRING_C_H

#include <stddef.h>
#include <stdint.h>
#include <stdbool.h>

/* Forward-declare Qt classes before entering extern "C". */
#ifdef __cplusplus
class QString;
class QByteArray;
extern "C" {
#else
typedef struct QString QString;
typedef struct QByteArray QByteArray;
#endif

/* --- lifecycle --- */
QString* qstring_new_empty(void);
QString* qstring_new_copy(const QString* other);
QString* qstring_new_from_qstringview(const uint16_t* data, ptrdiff_t size);
void     qstring_delete(QString* self);

/* --- basic accessors --- */
ptrdiff_t qstring_size(const QString* self);
bool      qstring_is_empty(const QString* self);
bool      qstring_is_null(const QString* self);
bool      qstring_is_detached(const QString* self);
void      qstring_detach(QString* self);
const uint16_t* qstring_utf16(const QString* self);

/* --- element access --- */
uint16_t qstring_at(const QString* self, ptrdiff_t i);
uint16_t qstring_front(const QString* self);
uint16_t qstring_back(const QString* self);

/* --- modifiers --- */
void qstring_clear(QString* self);
void qstring_resize(QString* self, ptrdiff_t n);
void qstring_resize_fill(QString* self, ptrdiff_t n, uint16_t c);
void qstring_reserve(QString* self, ptrdiff_t n);
void qstring_squeeze(QString* self);
void qstring_truncate(QString* self, ptrdiff_t pos);
void qstring_chop(QString* self, ptrdiff_t n);
void qstring_fill(QString* self, uint16_t c, ptrdiff_t n);

void qstring_append_qchar(QString* self, uint16_t c);
void qstring_append_qstring(QString* self, const QString* other);
void qstring_append_utf8(QString* self, const char* data, ptrdiff_t size);
void qstring_prepend_qchar(QString* self, uint16_t c);
void qstring_prepend_qstring(QString* self, const QString* other);
void qstring_prepend_utf8(QString* self, const char* data, ptrdiff_t size);
void qstring_insert_qchar(QString* self, ptrdiff_t i, uint16_t c);
void qstring_insert_qstring(QString* self, ptrdiff_t i, const QString* other);
void qstring_insert_utf8(QString* self, ptrdiff_t i, const char* data, ptrdiff_t size);

void qstring_remove(QString* self, ptrdiff_t i, ptrdiff_t n);
void qstring_remove_at(QString* self, ptrdiff_t i);
void qstring_remove_first(QString* self);
void qstring_remove_last(QString* self);
void qstring_replace_qchar_qchar(QString* self, uint16_t before, uint16_t after);
void qstring_replace_qchar_qstring(QString* self, uint16_t before, const QString* after);
void qstring_replace_qchar_utf8(QString* self, uint16_t before, const char* after, ptrdiff_t after_size);
void qstring_replace_qstring_qchar(QString* self, const QString* before, uint16_t after);
void qstring_replace_qstring_qstring(QString* self, const QString* before, const QString* after);
void qstring_replace_qstring_utf8(QString* self, const QString* before, const char* after, ptrdiff_t after_size);
void qstring_replace_utf8_qchar(QString* self, const char* before, ptrdiff_t before_size, uint16_t after);
void qstring_replace_utf8_qstring(QString* self, const char* before, ptrdiff_t before_size, const QString* after);
void qstring_replace_utf8_utf8(QString* self, const char* before, ptrdiff_t before_size, const char* after, ptrdiff_t after_size);
void qstring_replace_range_utf8(QString* self, ptrdiff_t i, ptrdiff_t n, const char* after, ptrdiff_t after_size);

/* --- slicing (returns new QString) --- */
QString* qstring_left(const QString* self, ptrdiff_t n);
QString* qstring_right(const QString* self, ptrdiff_t n);
QString* qstring_mid(const QString* self, ptrdiff_t pos, ptrdiff_t n);
QString* qstring_first(const QString* self, ptrdiff_t n);
QString* qstring_last(const QString* self, ptrdiff_t n);
QString* qstring_sliced(const QString* self, ptrdiff_t pos);
QString* qstring_sliced_n(const QString* self, ptrdiff_t pos, ptrdiff_t n);
QString* qstring_chopped(const QString* self, ptrdiff_t n);

/* --- search (cs: 0 = CaseSensitive, 1 = CaseInsensitive) --- */
ptrdiff_t qstring_index_of_qchar(const QString* self, uint16_t c, ptrdiff_t from, int cs);
ptrdiff_t qstring_index_of_qstring(const QString* self, const QString* other, ptrdiff_t from, int cs);
ptrdiff_t qstring_index_of_utf8(const QString* self, const char* data, ptrdiff_t size, ptrdiff_t from, int cs);
ptrdiff_t qstring_last_index_of_qchar(const QString* self, uint16_t c, ptrdiff_t from, int cs);
ptrdiff_t qstring_last_index_of_qstring(const QString* self, const QString* other, ptrdiff_t from, int cs);
ptrdiff_t qstring_last_index_of_utf8(const QString* self, const char* data, ptrdiff_t size, ptrdiff_t from, int cs);
bool      qstring_contains_qchar(const QString* self, uint16_t c, int cs);
bool      qstring_contains_qstring(const QString* self, const QString* other, int cs);
bool      qstring_contains_utf8(const QString* self, const char* data, ptrdiff_t size, int cs);
ptrdiff_t qstring_count_qchar(const QString* self, uint16_t c, int cs);
ptrdiff_t qstring_count_qstring(const QString* self, const QString* other, int cs);
ptrdiff_t qstring_count_utf8(const QString* self, const char* data, ptrdiff_t size, int cs);
bool      qstring_starts_with_qchar(const QString* self, uint16_t c, int cs);
bool      qstring_starts_with_qstring(const QString* self, const QString* other, int cs);
bool      qstring_starts_with_utf8(const QString* self, const char* data, ptrdiff_t size, int cs);
bool      qstring_ends_with_qchar(const QString* self, uint16_t c, int cs);
bool      qstring_ends_with_qstring(const QString* self, const QString* other, int cs);
bool      qstring_ends_with_utf8(const QString* self, const char* data, ptrdiff_t size, int cs);

/* --- case / trim --- */
QString* qstring_to_lower(const QString* self);
QString* qstring_to_upper(const QString* self);
QString* qstring_to_case_folded(const QString* self);
bool qstring_is_lower(const QString* self);
bool qstring_is_upper(const QString* self);
QString* qstring_trimmed(const QString* self);
QString* qstring_simplified(const QString* self);
QString* qstring_left_justified(const QString* self, ptrdiff_t width, uint16_t fill, bool truncate);
QString* qstring_right_justified(const QString* self, ptrdiff_t width, uint16_t fill, bool truncate);

/* --- compare --- */
int qstring_compare_qstring(const QString* self, const QString* other, int cs);
int qstring_compare_qchar(const QString* self, uint16_t c, int cs);
int qstring_compare_utf8(const QString* self, const char* data, ptrdiff_t size, int cs);

/* --- numeric --- */
int16_t  qstring_to_short(const QString* self, bool* ok, int base);
uint16_t qstring_to_ushort(const QString* self, bool* ok, int base);
int32_t  qstring_to_int(const QString* self, bool* ok, int base);
uint32_t qstring_to_uint(const QString* self, bool* ok, int base);
int64_t  qstring_to_long(const QString* self, bool* ok, int base);
uint64_t qstring_to_ulong(const QString* self, bool* ok, int base);
int64_t  qstring_to_longlong(const QString* self, bool* ok, int base);
uint64_t qstring_to_ulonglong(const QString* self, bool* ok, int base);
float    qstring_to_float(const QString* self, bool* ok);
double   qstring_to_double(const QString* self, bool* ok);
void qstring_set_num_ll(QString* self, int64_t n, int base);
void qstring_set_num_ull(QString* self, uint64_t n, int base);
void qstring_set_num_double(QString* self, double v, char format, int precision);

/* --- static --- */
QString* qstring_number_ll(int64_t n, int base);
QString* qstring_number_ull(uint64_t n, int base);
QString* qstring_number_double(double v, char format, int precision);
ptrdiff_t qstring_max_size(void);

/* --- predicates --- */
bool qstring_is_right_to_left(const QString* self);
bool qstring_is_valid_utf16(const QString* self);

/* --- conversion to QByteArray (typed out-parameter) --- */
void qstring_to_utf8_qba(const QString* self, QByteArray** out);
void qstring_to_latin1_qba(const QString* self, QByteArray** out);
void qstring_to_local8bit_qba(const QString* self, QByteArray** out);

#ifdef __cplusplus
}
#endif

#endif