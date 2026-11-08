/* C ABI for Qt 6.11.2 QtCore/QStringEncoder.
 *
 * QStringEncoder derives from QStringConverter, so the base-class state (iface/state) is
 * shared with the qstringconverter package's ABI. This header exposes only the
 * encoder-specific functionality.
 */
#ifndef GO_QT6CORE_QSTRINGENCODER_C_H
#define GO_QT6CORE_QSTRINGENCODER_C_H

#include <stddef.h>
#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

typedef struct go_qse_s go_qse_t; /* QStringEncoder* */

/* ---- Constructors ---- */
go_qse_t *qstringencoder_new_default(void);
go_qse_t *qstringencoder_new_from_encoding(int encoding, int flags);
go_qse_t *qstringencoder_new_from_any(const void *data, ptrdiff_t size, int kind, int flags);
go_qse_t *qstringencoder_new_from_iface(const void *iface);

/* ---- Lifecycle ---- */
void qstringencoder_delete(go_qse_t *self);

/* ---- Encoding ----
 * Encodes in_size UTF-16 code units (in_data treated as const char16_t*) into a newly
 * allocated QByteArray returned via *out_qba. Its size is written to *out_size.
 * If the encoder is invalid, the encoder's state.invalidChars is incremented and
 * *out_qba points to an empty QByteArray (never NULL).
 */
void qstringencoder_encode(go_qse_t *self,
    const uint16_t *in_data, ptrdiff_t in_size,
    void **out_qba, ptrdiff_t *out_size);

/* ---- Required space ---- */
ptrdiff_t qstringencoder_required_space(const go_qse_t *self, ptrdiff_t input_length);

/* ---- Append to buffer ----
 * Writes encoded bytes into buf (capacity buf_len) starting at buf[0].
 * Returns the number of bytes written. If the encoder is invalid, nothing is written and
 * the encoder's state.invalidChars is incremented.
 */
ptrdiff_t qstringencoder_append_to_buffer(go_qse_t *self,
    void *buf, ptrdiff_t buf_len,
    const uint16_t *in_data, ptrdiff_t in_size);

/* ---- Finalize ----
 * Flushes pending encoder state into buf (capacity buf_len).
 * Writes *out_next  = bytes written (pointer delta),
 *        *out_invalid = invalidChars delta (from Qt's FinalizeResultChar<char>),
 *        *out_err     = FinalizeResultError enum (0=NoError, 1=InvalidCharacters,
 *                       2=NotEnoughSpace).
 */
void qstringencoder_finalize(go_qse_t *self,
    void *buf, ptrdiff_t buf_len,
    ptrdiff_t *out_next, int16_t *out_invalid, int *out_err);

#ifdef __cplusplus
}
#endif

#endif