/* C ABI for Qt 6.11 QtCore/QStringConverter, QStringEncoder and QStringDecoder.
 *
 * This mirrors QtCore's own layout: qstringconverter.h declares QStringConverter,
 * QStringEncoder and QStringDecoder, and qstringconverter.cpp implements them
 * together. QStringConverterBase keeps its own header/implementation pair
 * (qstringconverterbase_c.h/.cpp), the same way QtCore separates
 * qstringconverter_base.h from qstringconverter.h.
 *
 * Handles are opaque to C:
 *   - go_qscbase_t — a QStringConverterBase object;
 *   - go_qsc_t     — a QStringConverter object;
 *   - go_qse_t     — a QStringEncoder object;
 *   - go_qsd_t     — a QStringDecoder object.
 *
 * In qstringconverter_c.cpp the handle types are real subclasses of helper
 * adapters over the Qt classes, so the ABI functions work with the objects
 * directly and do not cast between a handle and a Qt pointer.
 *
 * The ABI speaks only in C types: int for enums/flags, ptrdiff_t for qsizetype,
 * void* for QAnyStringView/QByteArrayView payloads and for the protected
 * Interface*. Conversions from those C types to Qt types are centralized in the
 * cvt namespace in qstringconverter_c.cpp; this header mentions no C++ type, so
 * cgo can include it directly.
 *
 * Ownership:
 *   - every *_new_* function returns a handle owned by the caller; release it
 *     with the matching *_delete;
 *   - Interface pointers returned by qstringconverter_iface() are borrowed from
 *     the converter and stay valid while the converter is alive;
 *   - qstringconverter_available_codecs() returns a malloc'ed array that must be
 *     released with qstringconverter_available_codecs_free().
 */
#ifndef GO_QT6CORE_QSTRINGCONVERTER_C_H
#define GO_QT6CORE_QSTRINGCONVERTER_C_H

#include <stddef.h>
#include <stdint.h>
#include <stdbool.h>

#include "qstringconverterbase_c.h"

#ifdef __cplusplus
class QByteArray;
extern "C" {
#else
typedef struct QByteArray QByteArray;
#endif

typedef struct go_qsc_s     go_qsc_t;
typedef struct go_qse_s     go_qse_t;
typedef struct go_qsd_s     go_qsd_t;

/* ==== QStringConverter ==== */

/* Constructors (all protected in Qt 6.11). */
go_qsc_t *qstringconverter_new_default(void);
go_qsc_t *qstringconverter_new_from_encoding(int encoding, int flags);
go_qsc_t *qstringconverter_new_from_any(const void *data, ptrdiff_t size, int kind, int flags);
go_qsc_t *qstringconverter_new_from_iface(const void *iface);

/* Lifecycle. */
void qstringconverter_delete(go_qsc_t *self);
void qstringconverter_swap(go_qsc_t *self, go_qsc_t *other);

/* Public methods. */
bool        qstringconverter_is_valid    (const go_qsc_t *self);
const char *qstringconverter_name        (const go_qsc_t *self);
void        qstringconverter_reset_state (go_qsc_t *self);
bool        qstringconverter_has_error   (const go_qsc_t *self);

/* Protected state accessors (State is move-only; fields are read/written
 * individually). */
void qstringconverter_get_state(const go_qsc_t *self,
    int *out_flags, int *out_internal,
    int64_t *out_remaining, int64_t *out_invalid, uint32_t out_data[4]);
void qstringconverter_set_state(go_qsc_t *self,
    int flags, int internal, int64_t remaining, int64_t invalid, const uint32_t data[4]);

/* Protected nested Interface: opaque accessors.
 * The returned function pointers are raw addresses; calling them from Go is not
 * supported (cgo has no portable function-pointer invocation). */
const void *qstringconverter_iface                   (const go_qsc_t *self);
const char *qstringconverter_interface_name          (const void *iface);
void       *qstringconverter_interface_to_utf16      (const void *iface);
void       *qstringconverter_interface_to_utf16_len  (const void *iface);
void       *qstringconverter_interface_from_utf16    (const void *iface);
void       *qstringconverter_interface_from_utf16_len(const void *iface);

/* Static helpers. */
const char *qstringconverter_name_for_encoding(int encoding);
/* Returns Encoding as int, or -1 if the name is not recognized.
 * kind: 0 = Utf8, 1 = Latin1, 2 = Utf16 (matching qanystringview.Kind values). */
int         qstringconverter_encoding_for_name(const void *data, ptrdiff_t size, int kind);
/* Returns Encoding as int, or -1 if not detected. */
int         qstringconverter_encoding_for_data(const void *data, ptrdiff_t size,
                                               uint16_t expectedFirstChar);
int         qstringconverter_encoding_for_html(const void *data, ptrdiff_t size);

/* Returns a malloc'ed array of NUL-terminated UTF-8 names (Qt's QStringList),
 * or NULL when the list is empty. Free with
 * qstringconverter_available_codecs_free(). */
char **qstringconverter_available_codecs(ptrdiff_t *out_len);
void   qstringconverter_available_codecs_free(char **names, ptrdiff_t len);

/* ==== QStringEncoder (inherits QStringConverter) ==== */

go_qse_t *qstringencoder_new_default(void);
go_qse_t *qstringencoder_new_from_encoding(int encoding, int flags);
go_qse_t *qstringencoder_new_from_any(const void *data, ptrdiff_t size, int kind, int flags);
go_qse_t *qstringencoder_new_from_iface(const void *iface);

void qstringencoder_delete(go_qse_t *self);

/* Encodes in_size UTF-16 code units into a newly allocated QByteArray.
 * If the encoder is invalid, state.invalidChars is set and the result is an
 * empty QByteArray (never NULL). */
QByteArray *qstringencoder_encode(go_qse_t *self,
    const uint16_t *in_data, ptrdiff_t in_size);

/* Maximum number of bytes needed to encode input_length UTF-16 code units. */
ptrdiff_t qstringencoder_required_space(const go_qse_t *self, ptrdiff_t input_length);

/* Writes encoded bytes into buf (capacity buf_len) and returns the number of
 * bytes written. If the encoder is invalid, nothing is written and
 * state.invalidChars is set. */
ptrdiff_t qstringencoder_append_to_buffer(go_qse_t *self,
    void *buf, ptrdiff_t buf_len,
    const uint16_t *in_data, ptrdiff_t in_size);

/* Flushes pending encoder state into buf (capacity buf_len, bytes).
 * *out_next = bytes written, *out_invalid = invalidChars delta,
 * *out_err = FinalizeResultError (0=NoError, 1=InvalidCharacters,
 * 2=NotEnoughSpace). */
void qstringencoder_finalize(go_qse_t *self,
    void *buf, ptrdiff_t buf_len,
    ptrdiff_t *out_next, int16_t *out_invalid, int *out_err);

/* ==== QStringDecoder (inherits QStringConverter) ==== */

go_qsd_t *qstringdecoder_new_default(void);
go_qsd_t *qstringdecoder_new_from_encoding(int encoding, int flags);
go_qsd_t *qstringdecoder_new_from_any(const void *data, ptrdiff_t size, int kind, int flags);
go_qsd_t *qstringdecoder_new_from_iface(const void *iface);

void qstringdecoder_delete(go_qsd_t *self);

/* Decodes in_size bytes into a newly allocated QString returned through
 * *out_qstr. If the decoder is invalid, state.invalidChars is set and the
 * result is an empty QString (never NULL). */
void qstringdecoder_decode(go_qsd_t *self,
    const char *in_data, ptrdiff_t in_size,
    void **out_qstr);

/* Maximum number of UTF-16 code units needed to decode input_length bytes. */
ptrdiff_t qstringdecoder_required_space(const go_qsd_t *self, ptrdiff_t input_length);

/* Writes decoded UTF-16 code units into buf (capacity buf_len) and returns the
 * number of code units written. If the decoder is invalid, nothing is written
 * and state.invalidChars is set. */
ptrdiff_t qstringdecoder_append_to_buffer(go_qsd_t *self,
    void *buf, ptrdiff_t buf_len,
    const char *in_data, ptrdiff_t in_size);

/* Flushes pending decoder state into buf (capacity buf_len, UTF-16 units).
 * *out_next = code units written, *out_invalid = invalidChars delta,
 * *out_err = FinalizeResultError. */
void qstringdecoder_finalize(go_qsd_t *self,
    void *buf, ptrdiff_t buf_len,
    ptrdiff_t *out_next, int16_t *out_invalid, int *out_err);

/* Static factory: QStringDecoder::decoderForHtml(QByteArrayView). */
go_qsd_t *qstringdecoder_decoder_for_html(const char *data, ptrdiff_t size);

#ifdef __cplusplus
}
#endif

#endif
