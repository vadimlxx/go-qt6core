/* Order matters: real Qt headers first, then our C ABI wrapper. */
// for not write  -Wno-sfinae-incomplete in go file
#if defined(__GNUC__) && !defined(__clang__) && (__GNUC__ >= 16)
#  pragma GCC diagnostic ignored "-Wsfinae-incomplete"
#endif

#include "../internal/qt6core/check_qt_version.h"

#include <QByteArray>
#include <QByteArrayView>
#include <QList>
#include <QtCore/qflags.h>
#include <QtCore/qnamespace.h>

#include <cstdlib>

#include "qbytearray_c.h"

namespace {

inline QByteArrayView V(const char *data, ptrdiff_t n) noexcept {
    return QByteArrayView(data, static_cast<qsizetype>(n));
}

/* Safe int -> QFlags<Base64Option> conversion. Using QFlag avoids
 * the reinterpret_cast fallback that a plain C-style cast would use. */
inline QByteArray::Base64Options BAOpts(int options) noexcept {
    return QByteArray::Base64Options(QFlag(options));
}

} // namespace

extern "C" {

/* --- lifecycle --- */
QByteArray* qbytearray_new_empty() { return new QByteArray(); }
QByteArray* qbytearray_new_from_data(const char *data, ptrdiff_t size) {
    return new QByteArray(data, static_cast<qsizetype>(size));
}
QByteArray* qbytearray_new_fill(ptrdiff_t size, char c) {
    return new QByteArray(static_cast<qsizetype>(size), c);
}
QByteArray* qbytearray_new_copy(const QByteArray *other) {
    return new QByteArray(*other);
}
void qbytearray_delete(QByteArray *self) { delete self; }

/* --- accessors --- */
ptrdiff_t qbytearray_size(const QByteArray *self)     { return static_cast<ptrdiff_t>(self->size()); }
ptrdiff_t qbytearray_capacity(const QByteArray *self) { return static_cast<ptrdiff_t>(self->capacity()); }
ptrdiff_t qbytearray_max_size()                       { return static_cast<ptrdiff_t>(QByteArray::maxSize()); }
bool      qbytearray_is_empty(const QByteArray *self) { return self->isEmpty(); }
bool      qbytearray_is_null(const QByteArray *self)  { return self->isNull(); }
const char* qbytearray_data(const QByteArray *self)   { return self->constData(); }
bool      qbytearray_is_detached(const QByteArray *self) { return self->isDetached(); }
void      qbytearray_detach(QByteArray *self)         { self->detach(); }

/* --- modifiers --- */
void qbytearray_resize(QByteArray *self, ptrdiff_t size) {
    self->resize(static_cast<qsizetype>(size));
}
void qbytearray_resize_fill(QByteArray *self, ptrdiff_t size, char c) {
    self->resize(static_cast<qsizetype>(size), c);
}
void qbytearray_resize_for_overwrite(QByteArray *self, ptrdiff_t size) {
    self->resizeForOverwrite(static_cast<qsizetype>(size));
}
void qbytearray_fill(QByteArray *self, char c, ptrdiff_t size) {
    self->fill(c, static_cast<qsizetype>(size));
}
void qbytearray_clear(QByteArray *self)   { self->clear(); }
void qbytearray_reserve(QByteArray *self, ptrdiff_t size) {
    self->reserve(static_cast<qsizetype>(size));
}
void qbytearray_squeeze(QByteArray *self) { self->squeeze(); }
void qbytearray_truncate(QByteArray *self, ptrdiff_t pos) {
    self->truncate(static_cast<qsizetype>(pos));
}
void qbytearray_chop(QByteArray *self, ptrdiff_t n) {
    self->chop(static_cast<qsizetype>(n));
}

/* --- element access --- */
char qbytearray_at(const QByteArray *self, ptrdiff_t i) {
    return self->at(static_cast<qsizetype>(i));
}
void qbytearray_set_at(QByteArray *self, ptrdiff_t i, char c) {
    (*self)[static_cast<qsizetype>(i)] = c;
}
char qbytearray_front(const QByteArray *self) { return self->front(); }
char qbytearray_back(const QByteArray *self)  { return self->back(); }

/* --- search --- */
ptrdiff_t qbytearray_index_of_char(const QByteArray *self, char c, ptrdiff_t from) {
    return static_cast<ptrdiff_t>(self->indexOf(c, static_cast<qsizetype>(from)));
}
ptrdiff_t qbytearray_index_of_data(const QByteArray *self, const char *data, ptrdiff_t size, ptrdiff_t from) {
    return static_cast<ptrdiff_t>(self->indexOf(V(data, size), static_cast<qsizetype>(from)));
}
ptrdiff_t qbytearray_last_index_of_char(const QByteArray *self, char c, ptrdiff_t from) {
    return static_cast<ptrdiff_t>(self->lastIndexOf(c, static_cast<qsizetype>(from)));
}
ptrdiff_t qbytearray_last_index_of_data(const QByteArray *self, const char *data, ptrdiff_t size, ptrdiff_t from) {
    return static_cast<ptrdiff_t>(self->lastIndexOf(V(data, size), static_cast<qsizetype>(from)));
}
ptrdiff_t qbytearray_count_char(const QByteArray *self, char c) {
    return static_cast<ptrdiff_t>(self->count(c));
}
ptrdiff_t qbytearray_count_data(const QByteArray *self, const char *data, ptrdiff_t size) {
    return static_cast<ptrdiff_t>(self->count(V(data, size)));
}
bool qbytearray_contains_char(const QByteArray *self, char c) { return self->contains(c); }
bool qbytearray_contains_data(const QByteArray *self, const char *data, ptrdiff_t size) {
    return self->contains(V(data, size));
}
int qbytearray_compare(const QByteArray *self, const char *data, ptrdiff_t size, int cs) {
    return self->compare(V(data, size),
                         cs == 0 ? Qt::CaseSensitive : Qt::CaseInsensitive);
}

/* --- slicing --- */
QByteArray* qbytearray_left(const QByteArray *self, ptrdiff_t n) {
    return new QByteArray(self->left(static_cast<qsizetype>(n)));
}
QByteArray* qbytearray_right(const QByteArray *self, ptrdiff_t n) {
    return new QByteArray(self->right(static_cast<qsizetype>(n)));
}
QByteArray* qbytearray_mid(const QByteArray *self, ptrdiff_t i, ptrdiff_t len) {
    return new QByteArray(self->mid(static_cast<qsizetype>(i),
                                    static_cast<qsizetype>(len)));
}
QByteArray* qbytearray_first(const QByteArray *self, ptrdiff_t n) {
    return new QByteArray(self->first(static_cast<qsizetype>(n)));
}
QByteArray* qbytearray_last(const QByteArray *self, ptrdiff_t n) {
    return new QByteArray(self->last(static_cast<qsizetype>(n)));
}
QByteArray* qbytearray_sliced(const QByteArray *self, ptrdiff_t pos) {
    return new QByteArray(self->sliced(static_cast<qsizetype>(pos)));
}
QByteArray* qbytearray_sliced_n(const QByteArray *self, ptrdiff_t pos, ptrdiff_t n) {
    return new QByteArray(self->sliced(static_cast<qsizetype>(pos),
                                       static_cast<qsizetype>(n)));
}
QByteArray* qbytearray_chopped(const QByteArray *self, ptrdiff_t n) {
    return new QByteArray(self->chopped(static_cast<qsizetype>(n)));
}

/* --- predicates --- */
bool qbytearray_starts_with_char(const QByteArray *self, char c) { return self->startsWith(c); }
bool qbytearray_starts_with_data(const QByteArray *self, const char *data, ptrdiff_t size) {
    return self->startsWith(V(data, size));
}
bool qbytearray_ends_with_char(const QByteArray *self, char c) { return self->endsWith(c); }
bool qbytearray_ends_with_data(const QByteArray *self, const char *data, ptrdiff_t size) {
    return self->endsWith(V(data, size));
}
bool qbytearray_is_upper(const QByteArray *self)      { return self->isUpper(); }
bool qbytearray_is_lower(const QByteArray *self)      { return self->isLower(); }
bool qbytearray_is_valid_utf8(const QByteArray *self) { return self->isValidUtf8(); }

/* --- transforms --- */
QByteArray* qbytearray_to_lower(const QByteArray *self)  { return new QByteArray(self->toLower()); }
QByteArray* qbytearray_to_upper(const QByteArray *self)  { return new QByteArray(self->toUpper()); }
QByteArray* qbytearray_trimmed(const QByteArray *self)   { return new QByteArray(self->trimmed()); }
QByteArray* qbytearray_simplified(const QByteArray *self){ return new QByteArray(self->simplified()); }
QByteArray* qbytearray_left_justified(const QByteArray *self, ptrdiff_t width, char fill, bool truncate) {
    return new QByteArray(self->leftJustified(static_cast<qsizetype>(width), fill, truncate));
}
QByteArray* qbytearray_right_justified(const QByteArray *self, ptrdiff_t width, char fill, bool truncate) {
    return new QByteArray(self->rightJustified(static_cast<qsizetype>(width), fill, truncate));
}
QByteArray* qbytearray_repeated(const QByteArray *self, ptrdiff_t times) {
    return new QByteArray(self->repeated(static_cast<qsizetype>(times)));
}

/* --- append/prepend/insert --- */
void qbytearray_append_char(QByteArray *self, char c) { self->append(c); }
void qbytearray_append_data(QByteArray *self, const char *data, ptrdiff_t size) {
    self->append(V(data, size));
}
void qbytearray_append_count(QByteArray *self, ptrdiff_t count, char c) {
    self->append(static_cast<qsizetype>(count), c);
}
void qbytearray_append_ba(QByteArray *self, const QByteArray *other) { self->append(*other); }

void qbytearray_prepend_char(QByteArray *self, char c) { self->prepend(c); }
void qbytearray_prepend_data(QByteArray *self, const char *data, ptrdiff_t size) {
    self->prepend(V(data, size));
}
void qbytearray_prepend_count(QByteArray *self, ptrdiff_t count, char c) {
    self->prepend(static_cast<qsizetype>(count), c);
}
void qbytearray_prepend_ba(QByteArray *self, const QByteArray *other) { self->prepend(*other); }

void qbytearray_insert_data(QByteArray *self, ptrdiff_t i, const char *data, ptrdiff_t size) {
    self->insert(static_cast<qsizetype>(i), V(data, size));
}
void qbytearray_insert_char(QByteArray *self, ptrdiff_t i, char c) {
    self->insert(static_cast<qsizetype>(i), c);
}
void qbytearray_insert_count(QByteArray *self, ptrdiff_t i, ptrdiff_t count, char c) {
    self->insert(static_cast<qsizetype>(i), static_cast<qsizetype>(count), c);
}

/* --- remove/replace --- */
void qbytearray_remove(QByteArray *self, ptrdiff_t i, ptrdiff_t len) {
    self->remove(static_cast<qsizetype>(i), static_cast<qsizetype>(len));
}
void qbytearray_remove_at(QByteArray *self, ptrdiff_t pos) {
    self->removeAt(static_cast<qsizetype>(pos));
}
void qbytearray_remove_first(QByteArray *self) { self->removeFirst(); }
void qbytearray_remove_last(QByteArray *self)  { self->removeLast(); }
void qbytearray_replace_char_char(QByteArray *self, char before, char after) {
    self->replace(before, after);
}
void qbytearray_replace_range(QByteArray *self, ptrdiff_t i, ptrdiff_t len,
                              const char *data, ptrdiff_t dsize) {
    self->replace(static_cast<qsizetype>(i), static_cast<qsizetype>(len),
                  V(data, dsize));
}
void qbytearray_replace_data_data(QByteArray *self,
                                  const char *before, ptrdiff_t bsize,
                                  const char *after,  ptrdiff_t asize) {
    self->replace(V(before, bsize), V(after, asize));
}

/* --- numeric conversion --- */
int16_t qbytearray_to_short(const QByteArray *self, bool *ok, int base) {
    return static_cast<int16_t>(self->toShort(ok, base));
}
uint16_t qbytearray_to_ushort(const QByteArray *self, bool *ok, int base) {
    return static_cast<uint16_t>(self->toUShort(ok, base));
}
int32_t qbytearray_to_int(const QByteArray *self, bool *ok, int base) {
    return static_cast<int32_t>(self->toInt(ok, base));
}
uint32_t qbytearray_to_uint(const QByteArray *self, bool *ok, int base) {
    return static_cast<uint32_t>(self->toUInt(ok, base));
}
int64_t qbytearray_to_longlong(const QByteArray *self, bool *ok, int base) {
    return static_cast<int64_t>(self->toLongLong(ok, base));
}
uint64_t qbytearray_to_ulonglong(const QByteArray *self, bool *ok, int base) {
    return static_cast<uint64_t>(self->toULongLong(ok, base));
}
float qbytearray_to_float(const QByteArray *self, bool *ok) {
    return self->toFloat(ok);
}
double qbytearray_to_double(const QByteArray *self, bool *ok) {
    return self->toDouble(ok);
}

/* --- static number constructors --- */
QByteArray* qbytearray_number_ll(int64_t n, int base) {
    return new QByteArray(QByteArray::number(static_cast<qlonglong>(n), base));
}
QByteArray* qbytearray_number_ull(uint64_t n, int base) {
    return new QByteArray(QByteArray::number(static_cast<qulonglong>(n), base));
}
QByteArray* qbytearray_number_double(double n, char format, int precision) {
    return new QByteArray(QByteArray::number(n, format, precision));
}
void qbytearray_set_num_ll(QByteArray *self, int64_t n, int base) {
    self->setNum(static_cast<qlonglong>(n), base);
}
void qbytearray_set_num_ull(QByteArray *self, uint64_t n, int base) {
    self->setNum(static_cast<qulonglong>(n), base);
}
void qbytearray_set_num_double(QByteArray *self, double n, char format, int precision) {
    self->setNum(n, format, precision);
}

/* --- encoding / decoding (Base64 via QFlag to avoid reinterpret_cast) --- */
QByteArray* qbytearray_to_base64(const QByteArray *self, int options) {
    return new QByteArray(self->toBase64(BAOpts(options)));
}
QByteArray* qbytearray_from_base64(const QByteArray *self, int options) {
    return new QByteArray(QByteArray::fromBase64(*self, BAOpts(options)));
}
QByteArray* qbytearray_from_base64_encoding(const char *data, ptrdiff_t size,
                                            int options, int *out_status) {
    QByteArray input(data, static_cast<qsizetype>(size));
    auto r = QByteArray::fromBase64Encoding(std::move(input), BAOpts(options));
    if (out_status)
        *out_status = static_cast<int>(r.decodingStatus);
    return new QByteArray(std::move(r.decoded));
}
QByteArray* qbytearray_to_hex(const QByteArray *self, char separator) {
    return new QByteArray(self->toHex(separator));
}
QByteArray* qbytearray_from_hex(const QByteArray *self) {
    return new QByteArray(QByteArray::fromHex(*self));
}
QByteArray* qbytearray_to_percent(const QByteArray *self,
                                  const QByteArray *exclude,
                                  const QByteArray *include,
                                  char percent) {
    return new QByteArray(self->toPercentEncoding(*exclude, *include, percent));
}
QByteArray* qbytearray_from_percent(const QByteArray *self, char percent) {
    return new QByteArray(QByteArray::fromPercentEncoding(*self, percent));
}

/* --- compression --- */
#ifndef QT_NO_COMPRESS
QByteArray* qbytearray_compress(const QByteArray *self, int level) {
    return new QByteArray(qCompress(*self, level));
}
QByteArray* qbytearray_uncompress(const QByteArray *self) {
    return new QByteArray(qUncompress(*self));
}
#else
QByteArray* qbytearray_compress(const QByteArray *, int)   { return new QByteArray(); }
QByteArray* qbytearray_uncompress(const QByteArray *)      { return new QByteArray(); }
#endif

/* --- split --- */
QByteArray** qbytearray_split(const QByteArray *self, char sep, ptrdiff_t *out_len) {
    const QList<QByteArray> parts = self->split(sep);
    *out_len = static_cast<ptrdiff_t>(parts.size());
    if (parts.isEmpty())
        return nullptr;
    auto **arr = static_cast<QByteArray**>(std::malloc(sizeof(QByteArray*) * parts.size()));
    for (qsizetype i = 0; i < parts.size(); ++i)
        arr[i] = new QByteArray(parts[i]);
    return arr;
}
void qbytearray_split_free_array(QByteArray **arr) {
    std::free(arr);
}

void qbytearray_null_terminate(QByteArray *self) {
	self->nullTerminate();
}
QByteArray* qbytearray_null_terminated(const QByteArray *self) {
	return new QByteArray(self->nullTerminated());
}
bool qbytearray_is_shared_with(const QByteArray *self, const QByteArray *other) {
	return self->isSharedWith(*other);
}

} // extern "C"
