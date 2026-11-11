/* Order: version check first, then real Qt headers. */
#include "../internal/qt6core/check_qt_version.h"

#include <QString>
#include <QStringView>
#include <QChar>
#include <QByteArray>

#include "qstring_c.h"

namespace {

inline QChar          C(uint16_t u)                    { return QChar(u); }
inline Qt::CaseSensitivity CS(int cs) noexcept         { return cs == 0 ? Qt::CaseSensitive : Qt::CaseInsensitive; }
inline QString        U8(const char* d, ptrdiff_t n)   { return QString::fromUtf8(d, static_cast<qsizetype>(n)); }

} // namespace

extern "C" {

/* --- lifecycle --- */
QString* qstring_new_empty() {
    return new QString();
}
QString* qstring_new_copy(const QString* other) {
    return new QString(*other);
}
QString* qstring_new_from_utf8(const char* data, ptrdiff_t size) {
    return new QString(U8(data, size));
}
QString* qstring_new_from_utf16(const uint16_t* data, ptrdiff_t size) {
    return new QString(QString::fromUtf16(reinterpret_cast<const char16_t*>(data),
                                          static_cast<qsizetype>(size)));
}
QString* qstring_new_from_latin1(const char* data, ptrdiff_t size) {
    return new QString(QString::fromLatin1(data, static_cast<qsizetype>(size)));
}
QString* qstring_new_from_qbytearray(const QByteArray* other) {
    if (!other) return new QString();
    return new QString(QString::fromUtf8(*other));
}
void qstring_delete(QString* self) {
    delete self;
}

/* --- accessors --- */
ptrdiff_t qstring_size(const QString* self)         { return static_cast<ptrdiff_t>(self->size()); }
bool qstring_is_empty(const QString* self)          { return self->isEmpty(); }
bool qstring_is_null(const QString* self)           { return self->isNull(); }
bool qstring_is_detached(const QString* self)       { return self->isDetached(); }
void qstring_detach(QString* self)                  { self->detach(); }
const uint16_t* qstring_utf16(const QString* self)  {
    return reinterpret_cast<const uint16_t*>(self->utf16());
}

/* --- element access --- */
uint16_t qstring_at(const QString* self, ptrdiff_t i) {
    return static_cast<uint16_t>(self->at(static_cast<qsizetype>(i)).unicode());
}
uint16_t qstring_front(const QString* self) {
    return static_cast<uint16_t>(self->front().unicode());
}
uint16_t qstring_back(const QString* self) {
    return static_cast<uint16_t>(self->back().unicode());
}

/* --- modifiers --- */
void qstring_clear(QString* self) { self->clear(); }
void qstring_resize(QString* self, ptrdiff_t n) { self->resize(static_cast<qsizetype>(n)); }
void qstring_resize_fill(QString* self, ptrdiff_t n, uint16_t c) {
    self->resize(static_cast<qsizetype>(n), C(c));
}
void qstring_reserve(QString* self, ptrdiff_t n) { self->reserve(static_cast<qsizetype>(n)); }
void qstring_squeeze(QString* self) { self->squeeze(); }
void qstring_truncate(QString* self, ptrdiff_t pos) { self->truncate(static_cast<qsizetype>(pos)); }
void qstring_chop(QString* self, ptrdiff_t n) { self->chop(static_cast<qsizetype>(n)); }
void qstring_fill(QString* self, uint16_t c, ptrdiff_t n) {
    self->fill(C(c), static_cast<qsizetype>(n));
}

void qstring_append_qchar(QString* self, uint16_t c) { self->append(C(c)); }
void qstring_append_qstring(QString* self, const QString* other) {
    self->append(*other);
}
void qstring_append_utf8(QString* self, const char* data, ptrdiff_t size) {
    self->append(U8(data, size));
}
void qstring_prepend_qchar(QString* self, uint16_t c) { self->prepend(C(c)); }
void qstring_prepend_qstring(QString* self, const QString* other) {
    self->prepend(*other);
}
void qstring_prepend_utf8(QString* self, const char* data, ptrdiff_t size) {
    self->prepend(U8(data, size));
}
void qstring_insert_qchar(QString* self, ptrdiff_t i, uint16_t c) {
    self->insert(static_cast<qsizetype>(i), C(c));
}
void qstring_insert_qstring(QString* self, ptrdiff_t i, const QString* other) {
    self->insert(static_cast<qsizetype>(i), *other);
}
void qstring_insert_utf8(QString* self, ptrdiff_t i, const char* data, ptrdiff_t size) {
    self->insert(static_cast<qsizetype>(i), U8(data, size));
}

void qstring_remove(QString* self, ptrdiff_t i, ptrdiff_t n) {
    self->remove(static_cast<qsizetype>(i), static_cast<qsizetype>(n));
}
void qstring_remove_at(QString* self, ptrdiff_t i) {
    self->removeAt(static_cast<qsizetype>(i));
}
void qstring_remove_first(QString* self) { self->removeFirst(); }
void qstring_remove_last(QString* self) { self->removeLast(); }

void qstring_replace_qchar_qchar(QString* self, uint16_t before, uint16_t after) {
    self->replace(C(before), C(after));
}
void qstring_replace_qchar_qstring(QString* self, uint16_t before, const QString* after) {
    self->replace(C(before), *after);
}
void qstring_replace_qchar_utf8(QString* self, uint16_t before, const char* after, ptrdiff_t after_size) {
    self->replace(C(before), U8(after, after_size));
}
void qstring_replace_qstring_qchar(QString* self, const QString* before, uint16_t after) {
    self->replace(*before, C(after));
}
void qstring_replace_qstring_qstring(QString* self, const QString* before, const QString* after) {
    self->replace(*before, *after);
}
void qstring_replace_qstring_utf8(QString* self, const QString* before, const char* after, ptrdiff_t after_size) {
    self->replace(*before, U8(after, after_size));
}
void qstring_replace_utf8_qchar(QString* self, const char* before, ptrdiff_t before_size, uint16_t after) {
    self->replace(U8(before, before_size), C(after));
}
void qstring_replace_utf8_qstring(QString* self, const char* before, ptrdiff_t before_size, const QString* after) {
    self->replace(U8(before, before_size), *after);
}
void qstring_replace_utf8_utf8(QString* self, const char* before, ptrdiff_t before_size, const char* after, ptrdiff_t after_size) {
    self->replace(U8(before, before_size), U8(after, after_size));
}
void qstring_replace_range_utf8(QString* self, ptrdiff_t i, ptrdiff_t n, const char* after, ptrdiff_t after_size) {
    self->replace(static_cast<qsizetype>(i), static_cast<qsizetype>(n), U8(after, after_size));
}

/* --- slicing --- */
QString* qstring_left(const QString* self, ptrdiff_t n) {
    return new QString(self->left(static_cast<qsizetype>(n)));
}
QString* qstring_right(const QString* self, ptrdiff_t n) {
    return new QString(self->right(static_cast<qsizetype>(n)));
}
QString* qstring_mid(const QString* self, ptrdiff_t pos, ptrdiff_t n) {
    return new QString(self->mid(static_cast<qsizetype>(pos),
                                 static_cast<qsizetype>(n)));
}
QString* qstring_first(const QString* self, ptrdiff_t n) {
    return new QString(self->first(static_cast<qsizetype>(n)));
}
QString* qstring_last(const QString* self, ptrdiff_t n) {
    return new QString(self->last(static_cast<qsizetype>(n)));
}
QString* qstring_sliced(const QString* self, ptrdiff_t pos) {
    return new QString(self->sliced(static_cast<qsizetype>(pos)));
}
QString* qstring_sliced_n(const QString* self, ptrdiff_t pos, ptrdiff_t n) {
    return new QString(self->sliced(static_cast<qsizetype>(pos),
                                    static_cast<qsizetype>(n)));
}
QString* qstring_chopped(const QString* self, ptrdiff_t n) {
    return new QString(self->chopped(static_cast<qsizetype>(n)));
}

/* --- search --- */
ptrdiff_t qstring_index_of_qchar(const QString* self, uint16_t c, ptrdiff_t from, int cs) {
    return static_cast<ptrdiff_t>(self->indexOf(C(c), static_cast<qsizetype>(from), CS(cs)));
}
ptrdiff_t qstring_index_of_qstring(const QString* self, const QString* other, ptrdiff_t from, int cs) {
    return static_cast<ptrdiff_t>(self->indexOf(*other, static_cast<qsizetype>(from), CS(cs)));
}
ptrdiff_t qstring_index_of_utf8(const QString* self, const char* data, ptrdiff_t size, ptrdiff_t from, int cs) {
    return static_cast<ptrdiff_t>(self->indexOf(U8(data, size), static_cast<qsizetype>(from), CS(cs)));
}
ptrdiff_t qstring_last_index_of_qchar(const QString* self, uint16_t c, ptrdiff_t from, int cs) {
    return static_cast<ptrdiff_t>(self->lastIndexOf(C(c), static_cast<qsizetype>(from), CS(cs)));
}
ptrdiff_t qstring_last_index_of_qstring(const QString* self, const QString* other, ptrdiff_t from, int cs) {
    return static_cast<ptrdiff_t>(self->lastIndexOf(*other, static_cast<qsizetype>(from), CS(cs)));
}
ptrdiff_t qstring_last_index_of_utf8(const QString* self, const char* data, ptrdiff_t size, ptrdiff_t from, int cs) {
    return static_cast<ptrdiff_t>(self->lastIndexOf(U8(data, size), static_cast<qsizetype>(from), CS(cs)));
}
bool qstring_contains_qchar(const QString* self, uint16_t c, int cs) {
    return self->contains(C(c), CS(cs));
}
bool qstring_contains_qstring(const QString* self, const QString* other, int cs) {
    return self->contains(*other, CS(cs));
}
bool qstring_contains_utf8(const QString* self, const char* data, ptrdiff_t size, int cs) {
    return self->contains(U8(data, size), CS(cs));
}
ptrdiff_t qstring_count_qchar(const QString* self, uint16_t c, int cs) {
    return static_cast<ptrdiff_t>(self->count(C(c), CS(cs)));
}
ptrdiff_t qstring_count_qstring(const QString* self, const QString* other, int cs) {
    return static_cast<ptrdiff_t>(self->count(*other, CS(cs)));
}
ptrdiff_t qstring_count_utf8(const QString* self, const char* data, ptrdiff_t size, int cs) {
    return static_cast<ptrdiff_t>(self->count(U8(data, size), CS(cs)));
}
bool qstring_starts_with_qchar(const QString* self, uint16_t c, int cs) {
    return self->startsWith(C(c), CS(cs));
}
bool qstring_starts_with_qstring(const QString* self, const QString* other, int cs) {
    return self->startsWith(*other, CS(cs));
}
bool qstring_starts_with_utf8(const QString* self, const char* data, ptrdiff_t size, int cs) {
    return self->startsWith(U8(data, size), CS(cs));
}
bool qstring_ends_with_qchar(const QString* self, uint16_t c, int cs) {
    return self->endsWith(C(c), CS(cs));
}
bool qstring_ends_with_qstring(const QString* self, const QString* other, int cs) {
    return self->endsWith(*other, CS(cs));
}
bool qstring_ends_with_utf8(const QString* self, const char* data, ptrdiff_t size, int cs) {
    return self->endsWith(U8(data, size), CS(cs));
}

/* --- case / trim --- */
QString* qstring_to_lower(const QString* self) {
    return new QString(self->toLower());
}
QString* qstring_to_upper(const QString* self) {
    return new QString(self->toUpper());
}
QString* qstring_to_case_folded(const QString* self) {
    return new QString(self->toCaseFolded());
}
bool qstring_is_lower(const QString* self) { return self->isLower(); }
bool qstring_is_upper(const QString* self) { return self->isUpper(); }
QString* qstring_trimmed(const QString* self) {
    return new QString(self->trimmed());
}
QString* qstring_simplified(const QString* self) {
    return new QString(self->simplified());
}
QString* qstring_left_justified(const QString* self, ptrdiff_t width, uint16_t fill, bool truncate) {
    return new QString(self->leftJustified(static_cast<qsizetype>(width), C(fill), truncate));
}
QString* qstring_right_justified(const QString* self, ptrdiff_t width, uint16_t fill, bool truncate) {
    return new QString(self->rightJustified(static_cast<qsizetype>(width), C(fill), truncate));
}

/* --- compare --- */
int qstring_compare_qstring(const QString* self, const QString* other, int cs) {
    return self->compare(*other, CS(cs));
}
int qstring_compare_qchar(const QString* self, uint16_t c, int cs) {
    return self->compare(C(c), CS(cs));
}
int qstring_compare_utf8(const QString* self, const char* data, ptrdiff_t size, int cs) {
    return self->compare(U8(data, size), CS(cs));
}

/* --- numeric --- */
int16_t qstring_to_short(const QString* self, bool* ok, int base) {
    return static_cast<int16_t>(self->toShort(ok, base));
}
uint16_t qstring_to_ushort(const QString* self, bool* ok, int base) {
    return static_cast<uint16_t>(self->toUShort(ok, base));
}
int32_t qstring_to_int(const QString* self, bool* ok, int base) {
    return static_cast<int32_t>(self->toInt(ok, base));
}
uint32_t qstring_to_uint(const QString* self, bool* ok, int base) {
    return static_cast<uint32_t>(self->toUInt(ok, base));
}
int64_t qstring_to_long(const QString* self, bool* ok, int base) {
    return static_cast<int64_t>(self->toLong(ok, base));
}
uint64_t qstring_to_ulong(const QString* self, bool* ok, int base) {
    return static_cast<uint64_t>(self->toULong(ok, base));
}
int64_t qstring_to_longlong(const QString* self, bool* ok, int base) {
    return static_cast<int64_t>(self->toLongLong(ok, base));
}
uint64_t qstring_to_ulonglong(const QString* self, bool* ok, int base) {
    return static_cast<uint64_t>(self->toULongLong(ok, base));
}
float qstring_to_float(const QString* self, bool* ok)  { return self->toFloat(ok); }
double qstring_to_double(const QString* self, bool* ok){ return self->toDouble(ok); }

void qstring_set_num_ll(QString* self, int64_t n, int base) {
    self->setNum(static_cast<qlonglong>(n), base);
}
void qstring_set_num_ull(QString* self, uint64_t n, int base) {
    self->setNum(static_cast<qulonglong>(n), base);
}
void qstring_set_num_double(QString* self, double v, char format, int precision) {
    self->setNum(v, format, precision);
}

/* --- static --- */
QString* qstring_number_ll(int64_t n, int base) {
    return new QString(QString::number(static_cast<qlonglong>(n), base));
}
QString* qstring_number_ull(uint64_t n, int base) {
    return new QString(QString::number(static_cast<qulonglong>(n), base));
}
QString* qstring_number_double(double v, char format, int precision) {
    return new QString(QString::number(v, format, precision));
}
ptrdiff_t qstring_max_size(void) {
    return static_cast<ptrdiff_t>(QString::maxSize());
}

/* --- predicates --- */
bool qstring_is_right_to_left(const QString* self) { return self->isRightToLeft(); }
bool qstring_is_valid_utf16(const QString* self)   { return self->isValidUtf16(); }

/* --- conversion to QByteArray --- */
void qstring_to_utf8_qba(const QString* self, QByteArray** out) {
    *out = new QByteArray(self->toUtf8());
}
void qstring_to_latin1_qba(const QString* self, QByteArray** out) {
    *out = new QByteArray(self->toLatin1());
}
void qstring_to_local8bit_qba(const QString* self, QByteArray** out) {
    *out = new QByteArray(self->toLocal8Bit());
}

} // extern "C"