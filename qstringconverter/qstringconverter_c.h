/* C ABI for Qt 6.11 QtCore/QStringConverter (and the mostly-empty QStringConverterBase).
 *
 * QStringConverterBase and QStringConverter are exposed as opaque pointers. The protected
 * nested Interface, protected iface/state fields, and protected ctors are all wrapped here
 * via a helper subclass on the C++ side.
 */
#ifndef GO_QT6CORE_QSTRINGCONVERTER_C_H
#define GO_QT6CORE_QSTRINGCONVERTER_C_H

#include <stddef.h>
#include <stdint.h>
#include <stdbool.h>

#ifdef __cplusplus
extern "C" {
#endif

typedef struct go_qscbase_s go_qscbase_t;
typedef struct go_qsc_s     go_qsc_t;

/* ---- QStringConverterBase (minimal in Qt 6.11: protected ctor/dtor + State) ---- */
go_qscbase_t *qstringconverterbase_new_default(void);
void          qstringconverterbase_delete(go_qscbase_t *self);

/* ---- QStringConverter ctors (all protected in Qt 6.11) ---- */
go_qsc_t *qstringconverter_new_default(void);
go_qsc_t *qstringconverter_new_from_encoding(int encoding, int flags);
go_qsc_t *qstringconverter_new_from_any(const void *data, ptrdiff_t size, int kind, int flags);
go_qsc_t *qstringconverter_new_from_iface(const void *iface);

/* ---- Lifecycle ---- */
void qstringconverter_delete(go_qsc_t *self);
void qstringconverter_swap(go_qsc_t *self, go_qsc_t *other);

/* ---- Public methods ---- */
bool        qstringconverter_is_valid    (const go_qsc_t *self);
const char *qstringconverter_name        (const go_qsc_t *self);
void        qstringconverter_reset_state (go_qsc_t *self);
bool        qstringconverter_has_error   (const go_qsc_t *self);

/* ---- Protected state accessors (State is move-only; we set/read fields, not bytes) ---- */
void qstringconverter_get_state(const go_qsc_t *self,
    int *out_flags, int *out_internal,
    int64_t *out_remaining, int64_t *out_invalid, uint32_t out_data[4]);
void qstringconverter_set_state(go_qsc_t *self,
    int flags, int internal, int64_t remaining, int64_t invalid, const uint32_t data[4]);

/* ---- Protected nested Interface: opaque accessors ---- */
const void *qstringconverter_iface                   (const go_qsc_t *self);
const char *qstringconverter_interface_name          (const void *iface);
void       *qstringconverter_interface_to_utf16      (const void *iface);
void       *qstringconverter_interface_to_utf16_len  (const void *iface);
void       *qstringconverter_interface_from_utf16    (const void *iface);
void       *qstringconverter_interface_from_utf16_len(const void *iface);

/* ---- Static helpers ---- */
const char *qstringconverter_name_for_encoding(int encoding);
/* Returns Encoding as int, or -1 if the name is not recognized.
 * kind: 0 = Utf8, 1 = Latin1, 2 = Utf16 (matching qanystringview.Kind values). */
int         qstringconverter_encoding_for_name(const void *data, ptrdiff_t size, int kind);
/* Returns Encoding as int, or -1 if not detected. */
int         qstringconverter_encoding_for_data(const void *data, ptrdiff_t size,
                                                uint16_t expectedFirstChar);
int         qstringconverter_encoding_for_html(const void *data, ptrdiff_t size);

/* TODO: availableCodecs() — returns QStringList, needs qstringlist package. */

#ifdef __cplusplus
}
#endif

#endif