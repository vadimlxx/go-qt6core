/* C ABI for Qt 6.11.2 QtCore/QStringDecoder.
 *
 * QStringDecoder derives from QStringConverter, so the base-class state (iface/state) is
 * shared with the qstringconverter package's ABI. This header exposes only the
 * decoder-specific functionality.
 */
#ifndef GO_QT6CORE_QSTRINGDECODER_C_H
#define GO_QT6CORE_QSTRINGDECODER_C_H

#include <stddef.h>
#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

typedef struct go_qsd_s go_qsd_t; /* QStringDecoder* */

/* ---- Constructors ---- */
go_qsd_t *qstringdecoder_new_default(void);
go_qsd_t *qstringdecoder_new_from_encoding(int encoding, int flags);
go_qsd_t *qstringdecoder_new_from_any(const void *data, ptrdiff_t size, int kind, int flags);
go_qsd_t *qstringdecoder_new_from_iface(const void *iface);

/* ---- Lifecycle ---- */
void qstringdecoder_delete(go_qsd_t *self);

/* ---- Decoding ----
 * Decodes in_size bytes (in_data treated as const char*) into a newly allocated QString
 * returned via *out_qstr. Its length is implicit in the QString object.
 * If the decoder is invalid, the decoder's state.invalidChars is incremented and *out_qstr
 * points to an empty QString (never NULL).
 */
void qstringdecoder_decode(go_qsd_t *self,
    const char *in_data, ptrdiff_t in_size,
    void **out_qstr);

/* ---- Required space (in UTF-16 code units) ---- */
ptrdiff_t qstringdecoder_required_space(const go_qsd_t *self, ptrdiff_t input_length);

/* ---- Append to buffer ----
 * Writes decoded UTF-16 code units into buf (capacity buf_len) starting at buf[0].
 * buf is treated as char16_t* (== QChar* on all supported platforms). Returns the number
 * of code units written. If the decoder is invalid, nothing is written and the decoder's
 * state.invalidChars is incremented.
 */
ptrdiff_t qstringdecoder_append_to_buffer(go_qsd_t *self,
    void *buf, ptrdiff_t buf_len,
    const char *in_data, ptrdiff_t in_size);

/* ---- Finalize ----
 * Flushes pending decoder state into buf (capacity buf_len, char16_t*).
 * Writes *out_next   = number of UTF-16 code units written (pointer delta),
 *        *out_invalid = invalidChars delta (from Qt's FinalizeResultChar<char16_t>),
 *        *out_err     = FinalizeResultError enum (0=NoError, 1=InvalidCharacters,
 *                       2=NotEnoughSpace).
 */
void qstringdecoder_finalize(go_qsd_t *self,
    void *buf, ptrdiff_t buf_len,
    ptrdiff_t *out_next, int16_t *out_invalid, int *out_err);

/* ---- decoderForHtml static ---- */
go_qsd_t *qstringdecoder_decoder_for_html(const char *data, ptrdiff_t size);

#ifdef __cplusplus
}
#endif

#endif