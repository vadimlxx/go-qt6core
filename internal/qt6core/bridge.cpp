/* Order matters: version check first, then real Qt headers. */
#include "check_qt_version.h"

#include <QString>
#include <QByteArray>

#include <cstdlib>
#include <cstring>

#include "bridge.h"

namespace {
inline const QString    *Q(const void *p) { return static_cast<const QString    *>(p); }
inline const QByteArray *B(const void *p) { return static_cast<const QByteArray *>(p); }
inline void *toVoid(QString *q)    { return static_cast<void *>(q); }
inline void *toVoid(QByteArray *b) { return static_cast<void *>(b); }
} // namespace

extern "C" {

/* ---- QString constructors ---- */

void *qt6core_qstring_from_utf8(const char *data, ptrdiff_t size) {
    return toVoid(new QString(QString::fromUtf8(data, static_cast<qsizetype>(size))));
}
void *qt6core_qstring_from_utf16(const uint16_t *data, ptrdiff_t units) {
    return toVoid(new QString(QString::fromUtf16(
        reinterpret_cast<const char16_t *>(data), static_cast<qsizetype>(units))));
}
void *qt6core_qstring_from_latin1(const char *data, ptrdiff_t size) {
    return toVoid(new QString(QString::fromLatin1(data, static_cast<qsizetype>(size))));
}

/* ---- QString extractors ---- */

void qt6core_qstring_to_utf8(const void *q, char **out_data, ptrdiff_t *out_size) {
    if (!q) { *out_data = nullptr; *out_size = 0; return; }
    const QByteArray utf8 = Q(q)->toUtf8();
    const qsizetype n = utf8.size();
    *out_size = static_cast<ptrdiff_t>(n);
    if (n == 0) { *out_data = nullptr; return; }
    char *p = static_cast<char *>(std::malloc(static_cast<size_t>(n)));
    if (p) std::memcpy(p, utf8.constData(), static_cast<size_t>(n));
    *out_data = p;
}

void qt6core_qstring_to_utf16(const void *q, uint16_t **out_data, ptrdiff_t *out_units) {
    if (!q) { *out_data = nullptr; *out_units = 0; return; }
    /* QString data() is UTF-16; we copy the units. */
    const QString *s = Q(q);
    const qsizetype n = s->size();
    *out_units = static_cast<ptrdiff_t>(n);
    if (n == 0) { *out_data = nullptr; return; }
    auto *p = static_cast<uint16_t *>(std::malloc(sizeof(uint16_t) * static_cast<size_t>(n)));
    if (p) std::memcpy(p, s->utf16(), sizeof(uint16_t) * static_cast<size_t>(n));
    *out_data = p;
}

void qt6core_qstring_to_latin1(const void *q, char **out_data, ptrdiff_t *out_size) {
    if (!q) { *out_data = nullptr; *out_size = 0; return; }
    const QByteArray l1 = Q(q)->toLatin1();
    const qsizetype n = l1.size();
    *out_size = static_cast<ptrdiff_t>(n);
    if (n == 0) { *out_data = nullptr; return; }
    char *p = static_cast<char *>(std::malloc(static_cast<size_t>(n)));
    if (p) std::memcpy(p, l1.constData(), static_cast<size_t>(n));
    *out_data = p;
}

/* ---- QByteArray <-> QString ---- */

void *qt6core_qbytearray_to_qstring(const void *ba) {
    if (!ba) return nullptr;
    return toVoid(new QString(QString::fromUtf8(*B(ba))));
}
void *qt6core_qstring_to_qbytearray(const void *q) {
    if (!q) return nullptr;
    return toVoid(new QByteArray(Q(q)->toUtf8()));
}

/* ---- Destructors ---- */

void qt6core_qstring_delete(void *q)    { delete static_cast<QString    *>(q); }
void qt6core_qbytearray_delete(void *ba){ delete static_cast<QByteArray *>(ba); }

/* ---- Version info ---- */

void qt6core_get_version(int *major, int *minor, int *patch) {
    *major = (QT_VERSION >> 16) & 0xFF;
    *minor = (QT_VERSION >>  8) & 0xFF;
    *patch =  QT_VERSION        & 0xFF;
}

} // extern "C"