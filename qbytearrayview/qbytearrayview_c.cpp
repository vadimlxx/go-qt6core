/* Order: real Qt headers first, then our C ABI wrapper. */
#include <QtCore/qglobal.h>
#if !defined(QT_VERSION) || QT_VERSION < QT_VERSION_CHECK(6, 11, 0)
#  error "go-qt6core requires Qt >= 6.11"
#endif

#include <QByteArray>       // brings QByteArrayView + QByteArrayView::toByteArray()
#include <QByteArrayView>

#include "qbytearrayview_c.h"

namespace {
inline QByteArrayView V(const char *data, ptrdiff_t n) noexcept {
    return QByteArrayView(data, static_cast<qsizetype>(n));
}
} // namespace

extern "C" {

/* --- predicates --- */
bool qbytearrayview_is_null(const char *data, ptrdiff_t size) {
    return V(data, size).isNull();
}
bool qbytearrayview_is_empty(const char *data, ptrdiff_t size) {
    return V(data, size).isEmpty();
}
bool qbytearrayview_is_valid_utf8(const char *data, ptrdiff_t size) {
    return V(data, size).isValidUtf8();
}

/* --- element access --- */
char qbytearrayview_at(const char *data, ptrdiff_t size, ptrdiff_t i) {
    return V(data, size).at(static_cast<qsizetype>(i));
}
char qbytearrayview_front(const char *data, ptrdiff_t size) {
    return V(data, size).front();
}
char qbytearrayview_back(const char *data, ptrdiff_t size) {
    return V(data, size).back();
}

/* --- search --- */
ptrdiff_t qbytearrayview_index_of_char(const char *data, ptrdiff_t size,
                                       char c, ptrdiff_t from) {
    return static_cast<ptrdiff_t>(V(data, size).indexOf(c, static_cast<qsizetype>(from)));
}
ptrdiff_t qbytearrayview_index_of_view(const char *data, ptrdiff_t size,
                                       const char *needle, ptrdiff_t needle_size,
                                       ptrdiff_t from) {
    return static_cast<ptrdiff_t>(
        V(data, size).indexOf(V(needle, needle_size), static_cast<qsizetype>(from)));
}
ptrdiff_t qbytearrayview_last_index_of_char(const char *data, ptrdiff_t size,
                                            char c, ptrdiff_t from) {
    return static_cast<ptrdiff_t>(
        V(data, size).lastIndexOf(c, static_cast<qsizetype>(from)));
}
ptrdiff_t qbytearrayview_last_index_of_view(const char *data, ptrdiff_t size,
                                            const char *needle, ptrdiff_t needle_size,
                                            ptrdiff_t from) {
    return static_cast<ptrdiff_t>(
        V(data, size).lastIndexOf(V(needle, needle_size), static_cast<qsizetype>(from)));
}
ptrdiff_t qbytearrayview_count_char(const char *data, ptrdiff_t size, char c) {
    return static_cast<ptrdiff_t>(V(data, size).count(c));
}
ptrdiff_t qbytearrayview_count_view(const char *data, ptrdiff_t size,
                                    const char *needle, ptrdiff_t needle_size) {
    return static_cast<ptrdiff_t>(V(data, size).count(V(needle, needle_size)));
}
bool qbytearrayview_contains_char(const char *data, ptrdiff_t size, char c) {
    return V(data, size).contains(c);
}
bool qbytearrayview_contains_view(const char *data, ptrdiff_t size,
                                  const char *needle, ptrdiff_t needle_size) {
    return V(data, size).contains(V(needle, needle_size));
}
bool qbytearrayview_starts_with_char(const char *data, ptrdiff_t size, char c) {
    return V(data, size).startsWith(c);
}
bool qbytearrayview_starts_with_view(const char *data, ptrdiff_t size,
                                     const char *prefix, ptrdiff_t prefix_size) {
    return V(data, size).startsWith(V(prefix, prefix_size));
}
bool qbytearrayview_ends_with_char(const char *data, ptrdiff_t size, char c) {
    return V(data, size).endsWith(c);
}
bool qbytearrayview_ends_with_view(const char *data, ptrdiff_t size,
                                   const char *suffix, ptrdiff_t suffix_size) {
    return V(data, size).endsWith(V(suffix, suffix_size));
}

/* --- compare --- */
int qbytearrayview_compare(const char *data, ptrdiff_t size,
                           const char *other, ptrdiff_t other_size,
                           int cs) {
    return V(data, size).compare(V(other, other_size),
                                 cs == 0 ? Qt::CaseSensitive : Qt::CaseInsensitive);
}

/* --- slicing --- */
namespace {
inline void out_view(const QByteArrayView &v, const char **d, ptrdiff_t *n) {
    *d = v.data();
    *n = static_cast<ptrdiff_t>(v.size());
}
} // namespace

void qbytearrayview_left(const char *data, ptrdiff_t size, ptrdiff_t n,
                         const char **out_data, ptrdiff_t *out_size) {
    out_view(V(data, size).left(static_cast<qsizetype>(n)), out_data, out_size);
}
void qbytearrayview_right(const char *data, ptrdiff_t size, ptrdiff_t n,
                          const char **out_data, ptrdiff_t *out_size) {
    out_view(V(data, size).right(static_cast<qsizetype>(n)), out_data, out_size);
}
void qbytearrayview_mid(const char *data, ptrdiff_t size, ptrdiff_t pos, ptrdiff_t n,
                        const char **out_data, ptrdiff_t *out_size) {
    out_view(V(data, size).mid(static_cast<qsizetype>(pos),
                               static_cast<qsizetype>(n)),
             out_data, out_size);
}
void qbytearrayview_first(const char *data, ptrdiff_t size, ptrdiff_t n,
                          const char **out_data, ptrdiff_t *out_size) {
    out_view(V(data, size).first(static_cast<qsizetype>(n)), out_data, out_size);
}
void qbytearrayview_last(const char *data, ptrdiff_t size, ptrdiff_t n,
                         const char **out_data, ptrdiff_t *out_size) {
    out_view(V(data, size).last(static_cast<qsizetype>(n)), out_data, out_size);
}
void qbytearrayview_sliced(const char *data, ptrdiff_t size, ptrdiff_t pos,
                           const char **out_data, ptrdiff_t *out_size) {
    out_view(V(data, size).sliced(static_cast<qsizetype>(pos)), out_data, out_size);
}
void qbytearrayview_sliced_n(const char *data, ptrdiff_t size, ptrdiff_t pos, ptrdiff_t n,
                             const char **out_data, ptrdiff_t *out_size) {
    out_view(V(data, size).sliced(static_cast<qsizetype>(pos),
                                  static_cast<qsizetype>(n)),
             out_data, out_size);
}
void qbytearrayview_chopped(const char *data, ptrdiff_t size, ptrdiff_t n,
                            const char **out_data, ptrdiff_t *out_size) {
    out_view(V(data, size).chopped(static_cast<qsizetype>(n)), out_data, out_size);
}
void qbytearrayview_trimmed(const char *data, ptrdiff_t size,
                            const char **out_data, ptrdiff_t *out_size) {
    out_view(V(data, size).trimmed(), out_data, out_size);
}

/* --- numeric conversion --- */
int16_t qbytearrayview_to_short(const char *data, ptrdiff_t size, bool *ok, int base) {
    return static_cast<int16_t>(V(data, size).toShort(ok, base));
}
uint16_t qbytearrayview_to_ushort(const char *data, ptrdiff_t size, bool *ok, int base) {
    return static_cast<uint16_t>(V(data, size).toUShort(ok, base));
}
int32_t qbytearrayview_to_int(const char *data, ptrdiff_t size, bool *ok, int base) {
    return static_cast<int32_t>(V(data, size).toInt(ok, base));
}
uint32_t qbytearrayview_to_uint(const char *data, ptrdiff_t size, bool *ok, int base) {
    return static_cast<uint32_t>(V(data, size).toUInt(ok, base));
}
int64_t qbytearrayview_to_longlong(const char *data, ptrdiff_t size, bool *ok, int base) {
    return static_cast<int64_t>(V(data, size).toLongLong(ok, base));
}
uint64_t qbytearrayview_to_ulonglong(const char *data, ptrdiff_t size, bool *ok, int base) {
    return static_cast<uint64_t>(V(data, size).toULongLong(ok, base));
}
float qbytearrayview_to_float(const char *data, ptrdiff_t size, bool *ok) {
    return V(data, size).toFloat(ok);
}
double qbytearrayview_to_double(const char *data, ptrdiff_t size, bool *ok) {
    return V(data, size).toDouble(ok);
}

/* --- to owning QByteArray --- */
void qbytearrayview_to_bytearray(const char *data, ptrdiff_t size, void **out_qbytearray) {
    QByteArray *q = new QByteArray(V(data, size).toByteArray());
    *out_qbytearray = static_cast<void *>(q);
}

/* --- static --- */
ptrdiff_t qbytearrayview_max_size(void) {
    return static_cast<ptrdiff_t>(QByteArrayView::maxSize());
}

} // extern "C"