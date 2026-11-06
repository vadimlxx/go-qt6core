/* Order: version check first, then real Qt headers. */
#include "../internal/qt6core/check_qt_version.h"

#include <QUtf8StringView>
#include <QString>            /* for compare/equal(QStringView/QLatin1StringView/QChar/QByteArray) and toString() */
#include <QByteArray>
#include <QChar>
#include <QStringView>
#include <QLatin1StringView>

#include "qutf8stringview_c.h"

namespace {

inline QUtf8StringView V(const char *d, ptrdiff_t n) noexcept {
    return QUtf8StringView(d, static_cast<qsizetype>(n));
}
inline QStringView SV(const uint16_t *d, ptrdiff_t n) noexcept {
    return QStringView(reinterpret_cast<const char16_t *>(d), static_cast<qsizetype>(n));
}
inline QLatin1StringView L1(const char *d, ptrdiff_t n) noexcept {
    return QLatin1StringView(d, static_cast<qsizetype>(n));
}
inline Qt::CaseSensitivity CS(int cs) noexcept {
    return cs == 0 ? Qt::CaseSensitive : Qt::CaseInsensitive;
}
inline void outView(const QUtf8StringView &v, const char **d, ptrdiff_t *n) noexcept {
    *d = v.data();
    *n = static_cast<ptrdiff_t>(v.size());
}

} // namespace

extern "C" {

/* --- predicates --- */
bool qutf8stringview_is_null(const char *data, ptrdiff_t size) {
    return V(data, size).isNull();
}
bool qutf8stringview_is_empty(const char *data, ptrdiff_t size) {
    return V(data, size).isEmpty();
}
bool qutf8stringview_is_valid_utf8(const char *data, ptrdiff_t size) {
    return V(data, size).isValidUtf8();
}

/* --- element access --- */
char qutf8stringview_at(const char *data, ptrdiff_t size, ptrdiff_t i) {
    return V(data, size).at(static_cast<qsizetype>(i));
}
char qutf8stringview_front(const char *data, ptrdiff_t size) {
    return V(data, size).front();
}
char qutf8stringview_back(const char *data, ptrdiff_t size) {
    return V(data, size).back();
}

/* --- slicing --- */
void qutf8stringview_left(const char *data, ptrdiff_t size, ptrdiff_t n,
                          const char **out_data, ptrdiff_t *out_size) {
    outView(V(data, size).left(static_cast<qsizetype>(n)), out_data, out_size);
}
void qutf8stringview_right(const char *data, ptrdiff_t size, ptrdiff_t n,
                           const char **out_data, ptrdiff_t *out_size) {
    outView(V(data, size).right(static_cast<qsizetype>(n)), out_data, out_size);
}
void qutf8stringview_mid(const char *data, ptrdiff_t size, ptrdiff_t pos, ptrdiff_t n,
                         const char **out_data, ptrdiff_t *out_size) {
    outView(V(data, size).mid(static_cast<qsizetype>(pos), static_cast<qsizetype>(n)),
            out_data, out_size);
}
void qutf8stringview_sliced(const char *data, ptrdiff_t size, ptrdiff_t pos,
                            const char **out_data, ptrdiff_t *out_size) {
    outView(V(data, size).sliced(static_cast<qsizetype>(pos)), out_data, out_size);
}
void qutf8stringview_sliced_n(const char *data, ptrdiff_t size, ptrdiff_t pos, ptrdiff_t n,
                              const char **out_data, ptrdiff_t *out_size) {
    outView(V(data, size).sliced(static_cast<qsizetype>(pos), static_cast<qsizetype>(n)),
            out_data, out_size);
}
void qutf8stringview_first(const char *data, ptrdiff_t size, ptrdiff_t n,
                           const char **out_data, ptrdiff_t *out_size) {
    outView(V(data, size).first(static_cast<qsizetype>(n)), out_data, out_size);
}
void qutf8stringview_last(const char *data, ptrdiff_t size, ptrdiff_t n,
                          const char **out_data, ptrdiff_t *out_size) {
    outView(V(data, size).last(static_cast<qsizetype>(n)), out_data, out_size);
}
void qutf8stringview_chopped(const char *data, ptrdiff_t size, ptrdiff_t n,
                             const char **out_data, ptrdiff_t *out_size) {
    outView(V(data, size).chopped(static_cast<qsizetype>(n)), out_data, out_size);
}

/* --- compare --- */
int qutf8stringview_compare_view(const char *data, ptrdiff_t size,
                                 const char *other, ptrdiff_t other_size, int cs) {
    return V(data, size).compare(V(other, other_size), CS(cs));
}
int qutf8stringview_compare_qbytearray(const char *data, ptrdiff_t size,
                                       const void *qba, int cs) {
    return V(data, size).compare(*static_cast<const QByteArray *>(qba), CS(cs));
}
int qutf8stringview_compare_qchar(const char *data, ptrdiff_t size, uint16_t c, int cs) {
    return V(data, size).compare(QChar(c), CS(cs));
}
int qutf8stringview_compare_qstringview(const char *data, ptrdiff_t size,
                                        const uint16_t *sv_data, ptrdiff_t sv_size, int cs) {
    return V(data, size).compare(SV(sv_data, sv_size), CS(cs));
}
int qutf8stringview_compare_qlatin1stringview(const char *data, ptrdiff_t size,
                                              const char *l1_data, ptrdiff_t l1_size, int cs) {
    return V(data, size).compare(L1(l1_data, l1_size), CS(cs));
}

/* --- equal --- */
bool qutf8stringview_equal_view(const char *data, ptrdiff_t size,
                                const char *other, ptrdiff_t other_size) {
    return V(data, size) == V(other, other_size);
}
bool qutf8stringview_equal_qbytearray(const char *data, ptrdiff_t size, const void *qba) {
    return V(data, size).equal(*static_cast<const QByteArray *>(qba));
}
bool qutf8stringview_equal_qchar(const char *data, ptrdiff_t size, uint16_t c) {
    return V(data, size).equal(QChar(c));
}
bool qutf8stringview_equal_qstringview(const char *data, ptrdiff_t size,
                                       const uint16_t *sv_data, ptrdiff_t sv_size) {
    return V(data, size).equal(SV(sv_data, sv_size));
}
bool qutf8stringview_equal_qlatin1stringview(const char *data, ptrdiff_t size,
                                             const char *l1_data, ptrdiff_t l1_size) {
    return V(data, size).equal(L1(l1_data, l1_size));
}

/* --- conversion to owning QString --- */
void qutf8stringview_to_qstring(const char *data, ptrdiff_t size, void **out_qstring) {
    *out_qstring = static_cast<void *>(new QString(V(data, size).toString()));
}

/* --- static --- */
ptrdiff_t qutf8stringview_max_size(void) {
    return static_cast<ptrdiff_t>(QUtf8StringView::maxSize());
}

} // extern "C"