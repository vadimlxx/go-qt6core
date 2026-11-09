/* Order: version check first, then real Qt headers. */
#include "../internal/qt6core/check_qt_version.h"

#include <QLatin1StringView>
#include <QStringView>
#include <QChar>
#include <QByteArray>

#include "qlatin1stringview_c.h"

namespace {

inline QLatin1StringView V(const char *data, ptrdiff_t n) noexcept {
    return QLatin1StringView(data, static_cast<qsizetype>(n));
}
inline QStringView SV(const uint16_t *data, ptrdiff_t n) noexcept {
    return QStringView(reinterpret_cast<const char16_t *>(data),
                       static_cast<qsizetype>(n));
}
inline Qt::CaseSensitivity CS(int cs) noexcept {
    return cs == 0 ? Qt::CaseSensitive : Qt::CaseInsensitive;
}
inline void outView(const QLatin1StringView &v, const char **d, ptrdiff_t *n) noexcept {
    *d = v.data();
    *n = static_cast<ptrdiff_t>(v.size());
}

} // namespace

extern "C" {

/* --- accessors --- */
uint8_t qlatin1stringview_at(const char *data, ptrdiff_t size, ptrdiff_t i) {
    return static_cast<uint8_t>(V(data, size).at(static_cast<qsizetype>(i)).toLatin1());
}
uint8_t qlatin1stringview_front(const char *data, ptrdiff_t size) {
    return static_cast<uint8_t>(V(data, size).front().toLatin1());
}
uint8_t qlatin1stringview_back(const char *data, ptrdiff_t size) {
    return static_cast<uint8_t>(V(data, size).back().toLatin1());
}

/* --- predicates --- */
bool qlatin1stringview_is_null(const char *data, ptrdiff_t size) {
    return V(data, size).isNull();
}
bool qlatin1stringview_is_empty(const char *data, ptrdiff_t size) {
    return V(data, size).isEmpty();
}

/* --- slicing --- */
void qlatin1stringview_mid(const char *data, ptrdiff_t size, ptrdiff_t pos, ptrdiff_t n,
                           const char **out_data, ptrdiff_t *out_size) {
    outView(V(data, size).mid(static_cast<qsizetype>(pos), static_cast<qsizetype>(n)),
            out_data, out_size);
}
void qlatin1stringview_left(const char *data, ptrdiff_t size, ptrdiff_t n,
                            const char **out_data, ptrdiff_t *out_size) {
    outView(V(data, size).left(static_cast<qsizetype>(n)), out_data, out_size);
}
void qlatin1stringview_right(const char *data, ptrdiff_t size, ptrdiff_t n,
                             const char **out_data, ptrdiff_t *out_size) {
    outView(V(data, size).right(static_cast<qsizetype>(n)), out_data, out_size);
}
void qlatin1stringview_sliced(const char *data, ptrdiff_t size, ptrdiff_t pos,
                              const char **out_data, ptrdiff_t *out_size) {
    outView(V(data, size).sliced(static_cast<qsizetype>(pos)), out_data, out_size);
}
void qlatin1stringview_sliced_n(const char *data, ptrdiff_t size, ptrdiff_t pos, ptrdiff_t n,
                                const char **out_data, ptrdiff_t *out_size) {
    outView(V(data, size).sliced(static_cast<qsizetype>(pos), static_cast<qsizetype>(n)),
            out_data, out_size);
}
void qlatin1stringview_first(const char *data, ptrdiff_t size, ptrdiff_t n,
                             const char **out_data, ptrdiff_t *out_size) {
    outView(V(data, size).first(static_cast<qsizetype>(n)), out_data, out_size);
}
void qlatin1stringview_last(const char *data, ptrdiff_t size, ptrdiff_t n,
                            const char **out_data, ptrdiff_t *out_size) {
    outView(V(data, size).last(static_cast<qsizetype>(n)), out_data, out_size);
}
void qlatin1stringview_chopped(const char *data, ptrdiff_t size, ptrdiff_t n,
                               const char **out_data, ptrdiff_t *out_size) {
    outView(V(data, size).chopped(static_cast<qsizetype>(n)), out_data, out_size);
}
void qlatin1stringview_trimmed(const char *data, ptrdiff_t size,
                               const char **out_data, ptrdiff_t *out_size) {
    outView(V(data, size).trimmed(), out_data, out_size);
}

/* --- compare --- */
int qlatin1stringview_compare_view(const char *data, ptrdiff_t size,
                                   const char *other, ptrdiff_t other_size, int cs) {
    return V(data, size).compare(V(other, other_size), CS(cs));
}
int qlatin1stringview_compare_qstringview(const char *data, ptrdiff_t size,
                                          const uint16_t *sv_data, ptrdiff_t sv_size, int cs) {
    return V(data, size).compare(SV(sv_data, sv_size), CS(cs));
}
int qlatin1stringview_compare_qchar(const char *data, ptrdiff_t size, uint16_t c, int cs) {
    return V(data, size).compare(QChar(c), CS(cs));
}

/* --- startsWith --- */
bool qlatin1stringview_starts_with_view(const char *data, ptrdiff_t size,
                                        const char *other, ptrdiff_t other_size, int cs) {
    return V(data, size).startsWith(V(other, other_size), CS(cs));
}
bool qlatin1stringview_starts_with_qstringview(const char *data, ptrdiff_t size,
                                               const uint16_t *sv_data, ptrdiff_t sv_size, int cs) {
    return V(data, size).startsWith(SV(sv_data, sv_size), CS(cs));
}
bool qlatin1stringview_starts_with_qchar(const char *data, ptrdiff_t size, uint16_t c, int cs) {
    return V(data, size).startsWith(QChar(c), CS(cs));
}

/* --- endsWith --- */
bool qlatin1stringview_ends_with_view(const char *data, ptrdiff_t size,
                                      const char *other, ptrdiff_t other_size, int cs) {
    return V(data, size).endsWith(V(other, other_size), CS(cs));
}
bool qlatin1stringview_ends_with_qstringview(const char *data, ptrdiff_t size,
                                             const uint16_t *sv_data, ptrdiff_t sv_size, int cs) {
    return V(data, size).endsWith(SV(sv_data, sv_size), CS(cs));
}
bool qlatin1stringview_ends_with_qchar(const char *data, ptrdiff_t size, uint16_t c, int cs) {
    return V(data, size).endsWith(QChar(c), CS(cs));
}

/* --- indexOf / lastIndexOf --- */
ptrdiff_t qlatin1stringview_index_of_view(const char *data, ptrdiff_t size,
                                          const char *other, ptrdiff_t other_size,
                                          ptrdiff_t from, int cs) {
    return static_cast<ptrdiff_t>(
        V(data, size).indexOf(V(other, other_size), static_cast<qsizetype>(from), CS(cs)));
}
ptrdiff_t qlatin1stringview_index_of_qstringview(const char *data, ptrdiff_t size,
                                                 const uint16_t *sv_data, ptrdiff_t sv_size,
                                                 ptrdiff_t from, int cs) {
    return static_cast<ptrdiff_t>(
        V(data, size).indexOf(SV(sv_data, sv_size), static_cast<qsizetype>(from), CS(cs)));
}
ptrdiff_t qlatin1stringview_index_of_qchar(const char *data, ptrdiff_t size,
                                           uint16_t c, ptrdiff_t from, int cs) {
    return static_cast<ptrdiff_t>(
        V(data, size).indexOf(QChar(c), static_cast<qsizetype>(from), CS(cs)));
}
ptrdiff_t qlatin1stringview_last_index_of_view(const char *data, ptrdiff_t size,
                                               const char *other, ptrdiff_t other_size,
                                               ptrdiff_t from, int cs) {
    return static_cast<ptrdiff_t>(
        V(data, size).lastIndexOf(V(other, other_size), static_cast<qsizetype>(from), CS(cs)));
}
ptrdiff_t qlatin1stringview_last_index_of_qstringview(const char *data, ptrdiff_t size,
                                                      const uint16_t *sv_data, ptrdiff_t sv_size,
                                                      ptrdiff_t from, int cs) {
    return static_cast<ptrdiff_t>(
        V(data, size).lastIndexOf(SV(sv_data, sv_size), static_cast<qsizetype>(from), CS(cs)));
}
ptrdiff_t qlatin1stringview_last_index_of_qchar(const char *data, ptrdiff_t size,
                                                uint16_t c, ptrdiff_t from, int cs) {
    return static_cast<ptrdiff_t>(
        V(data, size).lastIndexOf(QChar(c), static_cast<qsizetype>(from), CS(cs)));
}

/* --- contains --- */
bool qlatin1stringview_contains_view(const char *data, ptrdiff_t size,
                                     const char *other, ptrdiff_t other_size, int cs) {
    return V(data, size).contains(V(other, other_size), CS(cs));
}
bool qlatin1stringview_contains_qstringview(const char *data, ptrdiff_t size,
                                            const uint16_t *sv_data, ptrdiff_t sv_size, int cs) {
    return V(data, size).contains(SV(sv_data, sv_size), CS(cs));
}
bool qlatin1stringview_contains_qchar(const char *data, ptrdiff_t size, uint16_t c, int cs) {
    return V(data, size).contains(QChar(c), CS(cs));
}

/* --- count --- */
ptrdiff_t qlatin1stringview_count_view(const char *data, ptrdiff_t size,
                                       const char *other, ptrdiff_t other_size, int cs) {
    return static_cast<ptrdiff_t>(V(data, size).count(V(other, other_size), CS(cs)));
}
ptrdiff_t qlatin1stringview_count_qstringview(const char *data, ptrdiff_t size,
                                              const uint16_t *sv_data, ptrdiff_t sv_size, int cs) {
    return static_cast<ptrdiff_t>(V(data, size).count(SV(sv_data, sv_size), CS(cs)));
}
ptrdiff_t qlatin1stringview_count_qchar(const char *data, ptrdiff_t size, uint16_t c, int cs) {
    return static_cast<ptrdiff_t>(V(data, size).count(QChar(c), CS(cs)));
}

/* --- numeric --- */
int16_t qlatin1stringview_to_short(const char *data, ptrdiff_t size, bool *ok, int base) {
    return static_cast<int16_t>(V(data, size).toShort(ok, base));
}
uint16_t qlatin1stringview_to_ushort(const char *data, ptrdiff_t size, bool *ok, int base) {
    return static_cast<uint16_t>(V(data, size).toUShort(ok, base));
}
int32_t qlatin1stringview_to_int(const char *data, ptrdiff_t size, bool *ok, int base) {
    return static_cast<int32_t>(V(data, size).toInt(ok, base));
}
uint32_t qlatin1stringview_to_uint(const char *data, ptrdiff_t size, bool *ok, int base) {
    return static_cast<uint32_t>(V(data, size).toUInt(ok, base));
}
int64_t qlatin1stringview_to_long(const char *data, ptrdiff_t size, bool *ok, int base) {
    return static_cast<int64_t>(V(data, size).toLong(ok, base));
}
uint64_t qlatin1stringview_to_ulong(const char *data, ptrdiff_t size, bool *ok, int base) {
    return static_cast<uint64_t>(V(data, size).toULong(ok, base));
}
int64_t qlatin1stringview_to_longlong(const char *data, ptrdiff_t size, bool *ok, int base) {
    return static_cast<int64_t>(V(data, size).toLongLong(ok, base));
}
uint64_t qlatin1stringview_to_ulonglong(const char *data, ptrdiff_t size, bool *ok, int base) {
    return static_cast<uint64_t>(V(data, size).toULongLong(ok, base));
}
float qlatin1stringview_to_float(const char *data, ptrdiff_t size, bool *ok) {
    return V(data, size).toFloat(ok);
}
double qlatin1stringview_to_double(const char *data, ptrdiff_t size, bool *ok) {
    return V(data, size).toDouble(ok);
}

/* --- encoding --- */
QByteArray* qlatin1stringview_to_utf8(const char *data, ptrdiff_t size) {
	return new QByteArray(V(data, size).toUtf8());
}

/* --- static --- */
ptrdiff_t qlatin1stringview_max_size(void) {
    return static_cast<ptrdiff_t>(QLatin1StringView::maxSize());
}

} // extern "C"
