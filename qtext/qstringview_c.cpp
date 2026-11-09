/* Order: version check first, then real Qt headers. */
#include "../internal/qt6core/check_qt_version.h"

#include <QString>
#include <QStringView>
#include <QChar>
#include <QByteArray>

#include "qstringview_c.h"

namespace {

inline QStringView V(const uint16_t *data, ptrdiff_t n) noexcept {
    return QStringView(reinterpret_cast<const char16_t *>(data),
                       static_cast<qsizetype>(n));
}
inline Qt::CaseSensitivity CS(int cs) noexcept {
    return cs == 0 ? Qt::CaseSensitive : Qt::CaseInsensitive;
}
inline void outView(const QStringView &v, const uint16_t **d, ptrdiff_t *n) noexcept {
    *d = reinterpret_cast<const uint16_t *>(v.utf16());
    *n = static_cast<ptrdiff_t>(v.size());
}

} // namespace

extern "C" {

/* --- accessors --- */
uint16_t qstringview_at(const uint16_t *data, ptrdiff_t size, ptrdiff_t i) {
    return static_cast<uint16_t>(V(data, size).at(static_cast<qsizetype>(i)).unicode());
}
uint16_t qstringview_front(const uint16_t *data, ptrdiff_t size) {
    return static_cast<uint16_t>(V(data, size).front().unicode());
}
uint16_t qstringview_back(const uint16_t *data, ptrdiff_t size) {
    return static_cast<uint16_t>(V(data, size).back().unicode());
}

/* --- predicates --- */
bool qstringview_is_null(const uint16_t *data, ptrdiff_t size) {
    return V(data, size).isNull();
}
bool qstringview_is_empty(const uint16_t *data, ptrdiff_t size) {
    return V(data, size).isEmpty();
}
bool qstringview_is_upper(const uint16_t *data, ptrdiff_t size) {
    return V(data, size).isUpper();
}
bool qstringview_is_lower(const uint16_t *data, ptrdiff_t size) {
    return V(data, size).isLower();
}
bool qstringview_is_right_to_left(const uint16_t *data, ptrdiff_t size) {
    return V(data, size).isRightToLeft();
}
bool qstringview_is_valid_utf16(const uint16_t *data, ptrdiff_t size) {
    return V(data, size).isValidUtf16();
}

/* --- slicing --- */
void qstringview_left(const uint16_t *data, ptrdiff_t size, ptrdiff_t n,
                      const uint16_t **out_data, ptrdiff_t *out_size) {
    outView(V(data, size).left(static_cast<qsizetype>(n)), out_data, out_size);
}
void qstringview_right(const uint16_t *data, ptrdiff_t size, ptrdiff_t n,
                       const uint16_t **out_data, ptrdiff_t *out_size) {
    outView(V(data, size).right(static_cast<qsizetype>(n)), out_data, out_size);
}
void qstringview_mid(const uint16_t *data, ptrdiff_t size, ptrdiff_t pos, ptrdiff_t n,
                     const uint16_t **out_data, ptrdiff_t *out_size) {
    outView(V(data, size).mid(static_cast<qsizetype>(pos), static_cast<qsizetype>(n)),
            out_data, out_size);
}
void qstringview_first(const uint16_t *data, ptrdiff_t size, ptrdiff_t n,
                       const uint16_t **out_data, ptrdiff_t *out_size) {
    outView(V(data, size).first(static_cast<qsizetype>(n)), out_data, out_size);
}
void qstringview_last(const uint16_t *data, ptrdiff_t size, ptrdiff_t n,
                      const uint16_t **out_data, ptrdiff_t *out_size) {
    outView(V(data, size).last(static_cast<qsizetype>(n)), out_data, out_size);
}
void qstringview_sliced(const uint16_t *data, ptrdiff_t size, ptrdiff_t pos,
                        const uint16_t **out_data, ptrdiff_t *out_size) {
    outView(V(data, size).sliced(static_cast<qsizetype>(pos)), out_data, out_size);
}
void qstringview_sliced_n(const uint16_t *data, ptrdiff_t size, ptrdiff_t pos, ptrdiff_t n,
                          const uint16_t **out_data, ptrdiff_t *out_size) {
    outView(V(data, size).sliced(static_cast<qsizetype>(pos), static_cast<qsizetype>(n)),
            out_data, out_size);
}
void qstringview_chopped(const uint16_t *data, ptrdiff_t size, ptrdiff_t n,
                         const uint16_t **out_data, ptrdiff_t *out_size) {
    outView(V(data, size).chopped(static_cast<qsizetype>(n)), out_data, out_size);
}
void qstringview_trimmed(const uint16_t *data, ptrdiff_t size,
                         const uint16_t **out_data, ptrdiff_t *out_size) {
    outView(V(data, size).trimmed(), out_data, out_size);
}

/* --- search --- */
ptrdiff_t qstringview_index_of_qchar(const uint16_t *data, ptrdiff_t size,
                                     uint16_t c, ptrdiff_t from, int cs) {
    return static_cast<ptrdiff_t>(
        V(data, size).indexOf(QChar(c), static_cast<qsizetype>(from), CS(cs)));
}
ptrdiff_t qstringview_index_of_view(const uint16_t *data, ptrdiff_t size,
                                    const uint16_t *needle, ptrdiff_t needle_size,
                                    ptrdiff_t from, int cs) {
    return static_cast<ptrdiff_t>(
        V(data, size).indexOf(V(needle, needle_size),
                              static_cast<qsizetype>(from), CS(cs)));
}
ptrdiff_t qstringview_last_index_of_qchar(const uint16_t *data, ptrdiff_t size,
                                          uint16_t c, ptrdiff_t from, int cs) {
    return static_cast<ptrdiff_t>(
        V(data, size).lastIndexOf(QChar(c), static_cast<qsizetype>(from), CS(cs)));
}
ptrdiff_t qstringview_last_index_of_view(const uint16_t *data, ptrdiff_t size,
                                         const uint16_t *needle, ptrdiff_t needle_size,
                                         ptrdiff_t from, int cs) {
    return static_cast<ptrdiff_t>(
        V(data, size).lastIndexOf(V(needle, needle_size),
                                  static_cast<qsizetype>(from), CS(cs)));
}
ptrdiff_t qstringview_count_qchar(const uint16_t *data, ptrdiff_t size,
                                  uint16_t c, int cs) {
    return static_cast<ptrdiff_t>(V(data, size).count(QChar(c), CS(cs)));
}
ptrdiff_t qstringview_count_view(const uint16_t *data, ptrdiff_t size,
                                 const uint16_t *needle, ptrdiff_t needle_size, int cs) {
    return static_cast<ptrdiff_t>(V(data, size).count(V(needle, needle_size), CS(cs)));
}
bool qstringview_contains_qchar(const uint16_t *data, ptrdiff_t size,
                                uint16_t c, int cs) {
    return V(data, size).contains(QChar(c), CS(cs));
}
bool qstringview_contains_view(const uint16_t *data, ptrdiff_t size,
                               const uint16_t *needle, ptrdiff_t needle_size, int cs) {
    return V(data, size).contains(V(needle, needle_size), CS(cs));
}
bool qstringview_starts_with_qchar(const uint16_t *data, ptrdiff_t size,
                                   uint16_t c, int cs) {
    return V(data, size).startsWith(QChar(c), CS(cs));
}
bool qstringview_starts_with_view(const uint16_t *data, ptrdiff_t size,
                                  const uint16_t *prefix, ptrdiff_t prefix_size, int cs) {
    return V(data, size).startsWith(V(prefix, prefix_size), CS(cs));
}
bool qstringview_ends_with_qchar(const uint16_t *data, ptrdiff_t size,
                                 uint16_t c, int cs) {
    return V(data, size).endsWith(QChar(c), CS(cs));
}
bool qstringview_ends_with_view(const uint16_t *data, ptrdiff_t size,
                                const uint16_t *suffix, ptrdiff_t suffix_size, int cs) {
    return V(data, size).endsWith(V(suffix, suffix_size), CS(cs));
}

/* --- compare --- */
int qstringview_compare_view(const uint16_t *data, ptrdiff_t size,
                             const uint16_t *other, ptrdiff_t other_size, int cs) {
    return V(data, size).compare(V(other, other_size), CS(cs));
}
int qstringview_compare_qchar(const uint16_t *data, ptrdiff_t size,
                              uint16_t c, int cs) {
    return V(data, size).compare(QChar(c), CS(cs));
}

/* --- conversions --- */
QByteArray* qstringview_to_utf8(const uint16_t *data, ptrdiff_t size) {
	return new QByteArray(V(data, size).toUtf8());
}
QByteArray* qstringview_to_latin1(const uint16_t *data, ptrdiff_t size) {
	return new QByteArray(V(data, size).toLatin1());
}
QByteArray* qstringview_to_local8bit(const uint16_t *data, ptrdiff_t size) {
	return new QByteArray(V(data, size).toLocal8Bit());
}

/* --- numeric --- */
int16_t qstringview_to_short(const uint16_t *data, ptrdiff_t size, bool *ok, int base) {
    return static_cast<int16_t>(V(data, size).toShort(ok, base));
}
uint16_t qstringview_to_ushort(const uint16_t *data, ptrdiff_t size, bool *ok, int base) {
    return static_cast<uint16_t>(V(data, size).toUShort(ok, base));
}
int32_t qstringview_to_int(const uint16_t *data, ptrdiff_t size, bool *ok, int base) {
    return static_cast<int32_t>(V(data, size).toInt(ok, base));
}
uint32_t qstringview_to_uint(const uint16_t *data, ptrdiff_t size, bool *ok, int base) {
    return static_cast<uint32_t>(V(data, size).toUInt(ok, base));
}
int64_t qstringview_to_long(const uint16_t *data, ptrdiff_t size, bool *ok, int base) {
    return static_cast<int64_t>(V(data, size).toLong(ok, base));
}
uint64_t qstringview_to_ulong(const uint16_t *data, ptrdiff_t size, bool *ok, int base) {
    return static_cast<uint64_t>(V(data, size).toULong(ok, base));
}
int64_t qstringview_to_longlong(const uint16_t *data, ptrdiff_t size, bool *ok, int base) {
    return static_cast<int64_t>(V(data, size).toLongLong(ok, base));
}
uint64_t qstringview_to_ulonglong(const uint16_t *data, ptrdiff_t size, bool *ok, int base) {
    return static_cast<uint64_t>(V(data, size).toULongLong(ok, base));
}
float qstringview_to_float(const uint16_t *data, ptrdiff_t size, bool *ok) {
    return V(data, size).toFloat(ok);
}
double qstringview_to_double(const uint16_t *data, ptrdiff_t size, bool *ok) {
    return V(data, size).toDouble(ok);
}

/* --- static --- */
ptrdiff_t qstringview_max_size(void) {
    return static_cast<ptrdiff_t>(QStringView::maxSize());
}

} // extern "C"
