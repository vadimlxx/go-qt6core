/* C ABI for Qt 6.11 QtCore/QStringBuilder.
 *
 * QStringBuilder is a C++ template with no single type, so we cannot pass it across the ABI.
 * Instead, each side provides an append function that materializes a list of parts into a
 * target QString or QByteArray in one call.
 *
 * NOTE: the `data` field is uintptr_t (raw pointer value), not void*, to avoid cgo's
 * "Go pointer to Go pointer" check. Callers keep the target memory alive via
 * runtime.KeepAlive; C reinterprets the integer back to a pointer.
 */
#ifndef GO_QT6CORE_QSTRINGBUILDER_C_H
#define GO_QT6CORE_QSTRINGBUILDER_C_H

#include <stddef.h>
#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

enum qsb_kind {
    QSB_KIND_UTF8_STRING        = 0,  /* data = const char*, size = bytes */
    QSB_KIND_QLATIN1CHAR        = 1,  /* tag = Latin-1 byte */
    QSB_KIND_QCHAR              = 2,  /* tag = UTF-16 code unit */
    QSB_KIND_QSTRING            = 3,  /* data = QString* (from qstring.UnsafePtr) */
    QSB_KIND_QSTRINGVIEW        = 4,  /* data = const char16_t*, size = code units */
    QSB_KIND_QLATIN1STRINGVIEW  = 5,  /* data = const char*, size = bytes */
    QSB_KIND_QBYTEARRAY         = 6,  /* data = QByteArray* (from qbytearray.UnsafePtr) */
    QSB_KIND_QBYTEARRAYVIEW     = 7   /* data = const char*, size = bytes */
};

typedef struct qsb_part_s {
    int       kind;
    int       tag;
    uintptr_t data;   /* raw pointer value; not a Go pointer for cgo */
    ptrdiff_t size;
} qsb_part_t;

/* Materialize parts into an existing QString. qstr = qstring.QString.UnsafePtr(). */
void qstringbuilder_append_to_qstring(void *qstr, const qsb_part_t *parts, size_t n);

/* Materialize parts into an existing QByteArray. qba = qbytearray.QByteArray.UnsafePtr(). */
void qstringbuilder_append_to_qbytearray(void *qba, const qsb_part_t *parts, size_t n);

#ifdef __cplusplus
}
#endif

#endif