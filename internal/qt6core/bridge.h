/* Internal C ABI for cross-package QString/QByteArray conversions.
 *
 * This header is NOT part of the public API. Only internal/qt6core's
 * own .cpp includes it directly. Other packages call Go wrappers from
 * internal/qt6core and never see these symbols.
 *
 * Ownership rules:
 *   - All *_from_* functions return a heap-allocated object (new ...).
 *     The caller owns it and must eventually delete it (via the Go
 *     wrapper's finalizer / Close).
 *   - All *_to_* functions malloc a buffer for output. The caller
 *     copies it (C.GoBytes) and frees it with C.free.
 */
#ifndef GO_QT6CORE_BRIDGE_H
#define GO_QT6CORE_BRIDGE_H

#include <stddef.h>
#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

/* ---- QString constructors (return new QString*) ---- */
void* qt6core_qstring_from_utf8   (const char *data, ptrdiff_t size);
void* qt6core_qstring_from_utf16  (const uint16_t *data, ptrdiff_t units);
void* qt6core_qstring_from_latin1 (const char *data, ptrdiff_t size);

/* ---- QString extractors (malloc buffer, caller frees) ---- */
void qt6core_qstring_to_utf8   (const void *q, char **out_data, ptrdiff_t *out_size);
void qt6core_qstring_to_utf16  (const void *q, uint16_t **out_data, ptrdiff_t *out_units);
void qt6core_qstring_to_latin1 (const void *q, char **out_data, ptrdiff_t *out_size);

/* ---- QByteArray <-> QString (return new object of target type) ---- */
void* qt6core_qbytearray_to_qstring (const void *ba);  /* new QString*    */
void* qt6core_qstring_to_qbytearray (const void *q);   /* new QByteArray* */

/* ---- Destructors for objects created above ---- */
void qt6core_qstring_delete    (void *q);
void qt6core_qbytearray_delete (void *ba);

/* ---- Runtime version info (useful for diagnostics) ---- */
void qt6core_get_version(int *major, int *minor, int *patch);

#ifdef __cplusplus
}
#endif

#endif