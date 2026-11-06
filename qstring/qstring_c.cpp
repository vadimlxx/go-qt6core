/* Order: version check first, then real Qt headers. */
#include "../internal/qt6core/check_qt_version.h"

#include <QString>
#include <QStringView>
#include <QChar>
#include <QByteArray>

#include "qstring_c.h"

namespace {

inline QString       *Q(go_qt6core_qstring_t       *p) { return reinterpret_cast<QString *>(p); }
inline const QString *Q(const go_qt6core_qstring_t *p) { return reinterpret_cast<const QString *>(p); }
inline QChar          C(uint16_t u)                    { return QChar(u); }
inline Qt::CaseSensitivity CS(int cs) noexcept         { return cs == 0 ? Qt::CaseSensitive : Qt::CaseInsensitive; }
inline QString        U8(const char *d, ptrdiff_t n)   { return QString::fromUtf8(d, static_cast<qsizetype>(n)); }

} // namespace

extern "C" {

/* --- lifecycle --- */
go_qt6core_qstring_t* qstring_new_empty() { return reinterpret_cast<go_qt6core_qstring_t*>(new QString()); }
go_qt6core_qstring_t* qstring_new_copy(const go_qt6core_qstring_t *other) {
    return reinterpret_cast<go_qt6core_qstring_t*>(new QString(*Q(other)));
}
go_qt6core_qstring_t* qstring_new_from_qstringview(const uint16_t *data, ptrdiff_t size) {
    return reinterpret_cast<go_qt6core_qstring_t*>(
        new QString(QString::fromUtf16(reinterpret_cast<const char16_t *>(data),
                                       static_cast<qsizetype>(size))));
}
void qstring_delete(go_qt6core_qstring_t *self) { delete Q(self); }

/* --- accessors --- */
ptrdiff_t qstring_size(const go_qt6core_qstring_t *self)         { return static_cast<ptrdiff_t>(Q(self)->size()); }
bool qstring_is_empty(const go_qt6core_qstring_t *self)          { return Q(self)->isEmpty(); }
bool qstring_is_null(const go_qt6core_qstring_t *self)           { return Q(self)->isNull(); }
bool qstring_is_detached(const go_qt6core_qstring_t *self)       { return Q(self)->isDetached(); }
void qstring_detach(go_qt6core_qstring_t *self)                  { Q(self)->detach(); }
const uint16_t* qstring_utf16(const go_qt6core_qstring_t *self)  {
    return reinterpret_cast<const uint16_t *>(Q(self)->utf16());
}

/* --- element access --- */
uint16_t qstring_at(const go_qt6core_qstring_t *self, ptrdiff_t i) {
    return static_cast<uint16_t>(Q(self)->at(static_cast<qsizetype>(i)).unicode());
}
uint16_t qstring_front(const go_qt6core_qstring_t *self) {
    return static_cast<uint16_t>(Q(self)->front().unicode());
}
uint16_t qstring_back(const go_qt6core_qstring_t *self) {
    return static_cast<uint16_t>(Q(self)->back().unicode());
}

/* --- modifiers --- */
void qstring_clear(go_qt6core_qstring_t *self) { Q(self)->clear(); }
void qstring_resize(go_qt6core_qstring_t *self, ptrdiff_t n) { Q(self)->resize(static_cast<qsizetype>(n)); }
void qstring_resize_fill(go_qt6core_qstring_t *self, ptrdiff_t n, uint16_t c) {
    Q(self)->resize(static_cast<qsizetype>(n), C(c));
}
void qstring_reserve(go_qt6core_qstring_t *self, ptrdiff_t n) { Q(self)->reserve(static_cast<qsizetype>(n)); }
void qstring_squeeze(go_qt6core_qstring_t *self) { Q(self)->squeeze(); }
void qstring_truncate(go_qt6core_qstring_t *self, ptrdiff_t pos) { Q(self)->truncate(static_cast<qsizetype>(pos)); }
void qstring_chop(go_qt6core_qstring_t *self, ptrdiff_t n) { Q(self)->chop(static_cast<qsizetype>(n)); }
void qstring_fill(go_qt6core_qstring_t *self, uint16_t c, ptrdiff_t n) {
    Q(self)->fill(C(c), static_cast<qsizetype>(n));
}

void qstring_append_qchar(go_qt6core_qstring_t *self, uint16_t c) { Q(self)->append(C(c)); }
void qstring_append_qstring(go_qt6core_qstring_t *self, const go_qt6core_qstring_t *other) {
    Q(self)->append(*Q(other));
}
void qstring_append_utf8(go_qt6core_qstring_t *self, const char *data, ptrdiff_t size) {
    Q(self)->append(U8(data, size));
}
void qstring_prepend_qchar(go_qt6core_qstring_t *self, uint16_t c) { Q(self)->prepend(C(c)); }
void qstring_prepend_qstring(go_qt6core_qstring_t *self, const go_qt6core_qstring_t *other) {
    Q(self)->prepend(*Q(other));
}
void qstring_prepend_utf8(go_qt6core_qstring_t *self, const char *data, ptrdiff_t size) {
    Q(self)->prepend(U8(data, size));
}
void qstring_insert_qchar(go_qt6core_qstring_t *self, ptrdiff_t i, uint16_t c) {
    Q(self)->insert(static_cast<qsizetype>(i), C(c));
}
void qstring_insert_qstring(go_qt6core_qstring_t *self, ptrdiff_t i, const go_qt6core_qstring_t *other) {
    Q(self)->insert(static_cast<qsizetype>(i), *Q(other));
}
void qstring_insert_utf8(go_qt6core_qstring_t *self, ptrdiff_t i, const char *data, ptrdiff_t size) {
    Q(self)->insert(static_cast<qsizetype>(i), U8(data, size));
}

void qstring_remove(go_qt6core_qstring_t *self, ptrdiff_t i, ptrdiff_t n) {
    Q(self)->remove(static_cast<qsizetype>(i), static_cast<qsizetype>(n));
}
void qstring_remove_at(go_qt6core_qstring_t *self, ptrdiff_t i) {
    Q(self)->removeAt(static_cast<qsizetype>(i));
}
void qstring_remove_first(go_qt6core_qstring_t *self) { Q(self)->removeFirst(); }
void qstring_remove_last(go_qt6core_qstring_t *self) { Q(self)->removeLast(); }

void qstring_replace_qchar_qchar(go_qt6core_qstring_t *self, uint16_t before, uint16_t after) {
    Q(self)->replace(C(before), C(after));
}
void qstring_replace_qchar_qstring(go_qt6core_qstring_t *self, uint16_t before, const go_qt6core_qstring_t *after) {
    Q(self)->replace(C(before), *Q(after));
}
void qstring_replace_qchar_utf8(go_qt6core_qstring_t *self, uint16_t before, const char *after, ptrdiff_t after_size) {
    Q(self)->replace(C(before), U8(after, after_size));
}
void qstring_replace_qstring_qchar(go_qt6core_qstring_t *self, const go_qt6core_qstring_t *before, uint16_t after) {
    Q(self)->replace(*Q(before), C(after));
}
void qstring_replace_qstring_qstring(go_qt6core_qstring_t *self, const go_qt6core_qstring_t *before, const go_qt6core_qstring_t *after) {
    Q(self)->replace(*Q(before), *Q(after));
}
void qstring_replace_qstring_utf8(go_qt6core_qstring_t *self, const go_qt6core_qstring_t *before, const char *after, ptrdiff_t after_size) {
    Q(self)->replace(*Q(before), U8(after, after_size));
}
void qstring_replace_utf8_qchar(go_qt6core_qstring_t *self, const char *before, ptrdiff_t before_size, uint16_t after) {
    Q(self)->replace(U8(before, before_size), C(after));
}
void qstring_replace_utf8_qstring(go_qt6core_qstring_t *self, const char *before, ptrdiff_t before_size, const go_qt6core_qstring_t *after) {
    Q(self)->replace(U8(before, before_size), *Q(after));
}
void qstring_replace_utf8_utf8(go_qt6core_qstring_t *self, const char *before, ptrdiff_t before_size, const char *after, ptrdiff_t after_size) {
    Q(self)->replace(U8(before, before_size), U8(after, after_size));
}
void qstring_replace_range_utf8(go_qt6core_qstring_t *self, ptrdiff_t i, ptrdiff_t n, const char *after, ptrdiff_t after_size) {
    Q(self)->replace(static_cast<qsizetype>(i), static_cast<qsizetype>(n), U8(after, after_size));
}

/* --- slicing --- */
#define QSLICE(NAME, EXPR) \
    go_qt6core_qstring_t* NAME(const go_qt6core_qstring_t *self) { return reinterpret_cast<go_qt6core_qstring_t*>(new QString(EXPR)); }

go_qt6core_qstring_t* qstring_left(const go_qt6core_qstring_t *self, ptrdiff_t n) {
    return reinterpret_cast<go_qt6core_qstring_t*>(new QString(Q(self)->left(static_cast<qsizetype>(n))));
}
go_qt6core_qstring_t* qstring_right(const go_qt6core_qstring_t *self, ptrdiff_t n) {
    return reinterpret_cast<go_qt6core_qstring_t*>(new QString(Q(self)->right(static_cast<qsizetype>(n))));
}
go_qt6core_qstring_t* qstring_mid(const go_qt6core_qstring_t *self, ptrdiff_t pos, ptrdiff_t n) {
    return reinterpret_cast<go_qt6core_qstring_t*>(
        new QString(Q(self)->mid(static_cast<qsizetype>(pos), static_cast<qsizetype>(n))));
}
go_qt6core_qstring_t* qstring_first(const go_qt6core_qstring_t *self, ptrdiff_t n) {
    return reinterpret_cast<go_qt6core_qstring_t*>(new QString(Q(self)->first(static_cast<qsizetype>(n))));
}
go_qt6core_qstring_t* qstring_last(const go_qt6core_qstring_t *self, ptrdiff_t n) {
    return reinterpret_cast<go_qt6core_qstring_t*>(new QString(Q(self)->last(static_cast<qsizetype>(n))));
}
go_qt6core_qstring_t* qstring_sliced(const go_qt6core_qstring_t *self, ptrdiff_t pos) {
    return reinterpret_cast<go_qt6core_qstring_t*>(new QString(Q(self)->sliced(static_cast<qsizetype>(pos))));
}
go_qt6core_qstring_t* qstring_sliced_n(const go_qt6core_qstring_t *self, ptrdiff_t pos, ptrdiff_t n) {
    return reinterpret_cast<go_qt6core_qstring_t*>(
        new QString(Q(self)->sliced(static_cast<qsizetype>(pos), static_cast<qsizetype>(n))));
}
go_qt6core_qstring_t* qstring_chopped(const go_qt6core_qstring_t *self, ptrdiff_t n) {
    return reinterpret_cast<go_qt6core_qstring_t*>(new QString(Q(self)->chopped(static_cast<qsizetype>(n))));
}
#undef QSLICE

/* --- search --- */
ptrdiff_t qstring_index_of_qchar(const go_qt6core_qstring_t *self, uint16_t c, ptrdiff_t from, int cs) {
    return static_cast<ptrdiff_t>(Q(self)->indexOf(C(c), static_cast<qsizetype>(from), CS(cs)));
}
ptrdiff_t qstring_index_of_qstring(const go_qt6core_qstring_t *self, const go_qt6core_qstring_t *other, ptrdiff_t from, int cs) {
    return static_cast<ptrdiff_t>(Q(self)->indexOf(*Q(other), static_cast<qsizetype>(from), CS(cs)));
}
ptrdiff_t qstring_index_of_utf8(const go_qt6core_qstring_t *self, const char *data, ptrdiff_t size, ptrdiff_t from, int cs) {
    return static_cast<ptrdiff_t>(Q(self)->indexOf(U8(data, size), static_cast<qsizetype>(from), CS(cs)));
}
ptrdiff_t qstring_last_index_of_qchar(const go_qt6core_qstring_t *self, uint16_t c, ptrdiff_t from, int cs) {
    return static_cast<ptrdiff_t>(Q(self)->lastIndexOf(C(c), static_cast<qsizetype>(from), CS(cs)));
}
ptrdiff_t qstring_last_index_of_qstring(const go_qt6core_qstring_t *self, const go_qt6core_qstring_t *other, ptrdiff_t from, int cs) {
    return static_cast<ptrdiff_t>(Q(self)->lastIndexOf(*Q(other), static_cast<qsizetype>(from), CS(cs)));
}
ptrdiff_t qstring_last_index_of_utf8(const go_qt6core_qstring_t *self, const char *data, ptrdiff_t size, ptrdiff_t from, int cs) {
    return static_cast<ptrdiff_t>(Q(self)->lastIndexOf(U8(data, size), static_cast<qsizetype>(from), CS(cs)));
}
bool qstring_contains_qchar(const go_qt6core_qstring_t *self, uint16_t c, int cs) {
    return Q(self)->contains(C(c), CS(cs));
}
bool qstring_contains_qstring(const go_qt6core_qstring_t *self, const go_qt6core_qstring_t *other, int cs) {
    return Q(self)->contains(*Q(other), CS(cs));
}
bool qstring_contains_utf8(const go_qt6core_qstring_t *self, const char *data, ptrdiff_t size, int cs) {
    return Q(self)->contains(U8(data, size), CS(cs));
}
ptrdiff_t qstring_count_qchar(const go_qt6core_qstring_t *self, uint16_t c, int cs) {
    return static_cast<ptrdiff_t>(Q(self)->count(C(c), CS(cs)));
}
ptrdiff_t qstring_count_qstring(const go_qt6core_qstring_t *self, const go_qt6core_qstring_t *other, int cs) {
    return static_cast<ptrdiff_t>(Q(self)->count(*Q(other), CS(cs)));
}
ptrdiff_t qstring_count_utf8(const go_qt6core_qstring_t *self, const char *data, ptrdiff_t size, int cs) {
    return static_cast<ptrdiff_t>(Q(self)->count(U8(data, size), CS(cs)));
}
bool qstring_starts_with_qchar(const go_qt6core_qstring_t *self, uint16_t c, int cs) {
    return Q(self)->startsWith(C(c), CS(cs));
}
bool qstring_starts_with_qstring(const go_qt6core_qstring_t *self, const go_qt6core_qstring_t *other, int cs) {
    return Q(self)->startsWith(*Q(other), CS(cs));
}
bool qstring_starts_with_utf8(const go_qt6core_qstring_t *self, const char *data, ptrdiff_t size, int cs) {
    return Q(self)->startsWith(U8(data, size), CS(cs));
}
bool qstring_ends_with_qchar(const go_qt6core_qstring_t *self, uint16_t c, int cs) {
    return Q(self)->endsWith(C(c), CS(cs));
}
bool qstring_ends_with_qstring(const go_qt6core_qstring_t *self, const go_qt6core_qstring_t *other, int cs) {
    return Q(self)->endsWith(*Q(other), CS(cs));
}
bool qstring_ends_with_utf8(const go_qt6core_qstring_t *self, const char *data, ptrdiff_t size, int cs) {
    return Q(self)->endsWith(U8(data, size), CS(cs));
}

/* --- case / trim --- */
go_qt6core_qstring_t* qstring_to_lower(const go_qt6core_qstring_t *self) {
    return reinterpret_cast<go_qt6core_qstring_t*>(new QString(Q(self)->toLower()));
}
go_qt6core_qstring_t* qstring_to_upper(const go_qt6core_qstring_t *self) {
    return reinterpret_cast<go_qt6core_qstring_t*>(new QString(Q(self)->toUpper()));
}
go_qt6core_qstring_t* qstring_to_case_folded(const go_qt6core_qstring_t *self) {
    return reinterpret_cast<go_qt6core_qstring_t*>(new QString(Q(self)->toCaseFolded()));
}
bool qstring_is_lower(const go_qt6core_qstring_t *self) { return Q(self)->isLower(); }
bool qstring_is_upper(const go_qt6core_qstring_t *self) { return Q(self)->isUpper(); }
go_qt6core_qstring_t* qstring_trimmed(const go_qt6core_qstring_t *self) {
    return reinterpret_cast<go_qt6core_qstring_t*>(new QString(Q(self)->trimmed()));
}
go_qt6core_qstring_t* qstring_simplified(const go_qt6core_qstring_t *self) {
    return reinterpret_cast<go_qt6core_qstring_t*>(new QString(Q(self)->simplified()));
}
go_qt6core_qstring_t* qstring_left_justified(const go_qt6core_qstring_t *self, ptrdiff_t width, uint16_t fill, bool truncate) {
    return reinterpret_cast<go_qt6core_qstring_t*>(
        new QString(Q(self)->leftJustified(static_cast<qsizetype>(width), C(fill), truncate)));
}
go_qt6core_qstring_t* qstring_right_justified(const go_qt6core_qstring_t *self, ptrdiff_t width, uint16_t fill, bool truncate) {
    return reinterpret_cast<go_qt6core_qstring_t*>(
        new QString(Q(self)->rightJustified(static_cast<qsizetype>(width), C(fill), truncate)));
}

/* --- compare --- */
int qstring_compare_qstring(const go_qt6core_qstring_t *self, const go_qt6core_qstring_t *other, int cs) {
    return Q(self)->compare(*Q(other), CS(cs));
}
int qstring_compare_qchar(const go_qt6core_qstring_t *self, uint16_t c, int cs) {
    return Q(self)->compare(C(c), CS(cs));
}
int qstring_compare_utf8(const go_qt6core_qstring_t *self, const char *data, ptrdiff_t size, int cs) {
    return Q(self)->compare(U8(data, size), CS(cs));
}

/* --- numeric --- */
int16_t qstring_to_short(const go_qt6core_qstring_t *self, bool *ok, int base) {
    return static_cast<int16_t>(Q(self)->toShort(ok, base));
}
uint16_t qstring_to_ushort(const go_qt6core_qstring_t *self, bool *ok, int base) {
    return static_cast<uint16_t>(Q(self)->toUShort(ok, base));
}
int32_t qstring_to_int(const go_qt6core_qstring_t *self, bool *ok, int base) {
    return static_cast<int32_t>(Q(self)->toInt(ok, base));
}
uint32_t qstring_to_uint(const go_qt6core_qstring_t *self, bool *ok, int base) {
    return static_cast<uint32_t>(Q(self)->toUInt(ok, base));
}
int64_t qstring_to_long(const go_qt6core_qstring_t *self, bool *ok, int base) {
    return static_cast<int64_t>(Q(self)->toLong(ok, base));
}
uint64_t qstring_to_ulong(const go_qt6core_qstring_t *self, bool *ok, int base) {
    return static_cast<uint64_t>(Q(self)->toULong(ok, base));
}
int64_t qstring_to_longlong(const go_qt6core_qstring_t *self, bool *ok, int base) {
    return static_cast<int64_t>(Q(self)->toLongLong(ok, base));
}
uint64_t qstring_to_ulonglong(const go_qt6core_qstring_t *self, bool *ok, int base) {
    return static_cast<uint64_t>(Q(self)->toULongLong(ok, base));
}
float qstring_to_float(const go_qt6core_qstring_t *self, bool *ok)  { return Q(self)->toFloat(ok); }
double qstring_to_double(const go_qt6core_qstring_t *self, bool *ok){ return Q(self)->toDouble(ok); }

void qstring_set_num_ll(go_qt6core_qstring_t *self, int64_t n, int base) {
    Q(self)->setNum(static_cast<qlonglong>(n), base);
}
void qstring_set_num_ull(go_qt6core_qstring_t *self, uint64_t n, int base) {
    Q(self)->setNum(static_cast<qulonglong>(n), base);
}
void qstring_set_num_double(go_qt6core_qstring_t *self, double v, char format, int precision) {
    Q(self)->setNum(v, format, precision);
}

/* --- static --- */
go_qt6core_qstring_t* qstring_number_ll(int64_t n, int base) {
    return reinterpret_cast<go_qt6core_qstring_t*>(new QString(QString::number(static_cast<qlonglong>(n), base)));
}
go_qt6core_qstring_t* qstring_number_ull(uint64_t n, int base) {
    return reinterpret_cast<go_qt6core_qstring_t*>(new QString(QString::number(static_cast<qulonglong>(n), base)));
}
go_qt6core_qstring_t* qstring_number_double(double v, char format, int precision) {
    return reinterpret_cast<go_qt6core_qstring_t*>(new QString(QString::number(v, format, precision)));
}
ptrdiff_t qstring_max_size(void) { return static_cast<ptrdiff_t>(QString::maxSize()); }

/* --- predicates --- */
bool qstring_is_right_to_left(const go_qt6core_qstring_t *self) { return Q(self)->isRightToLeft(); }
bool qstring_is_valid_utf16(const go_qt6core_qstring_t *self)   { return Q(self)->isValidUtf16(); }

/* --- conversion to QByteArray --- */
void qstring_to_utf8_qba(const go_qt6core_qstring_t *self, void **out) {
    *out = static_cast<void *>(new QByteArray(Q(self)->toUtf8()));
}
void qstring_to_latin1_qba(const go_qt6core_qstring_t *self, void **out) {
    *out = static_cast<void *>(new QByteArray(Q(self)->toLatin1()));
}
void qstring_to_local8bit_qba(const go_qt6core_qstring_t *self, void **out) {
    *out = static_cast<void *>(new QByteArray(Q(self)->toLocal8Bit()));
}

} // extern "C"