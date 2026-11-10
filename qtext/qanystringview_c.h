/* C ABI for Qt 6.11 QtCore/QAnyStringView.
 *
 * QAnyStringView is passed as (const void *data, ptrdiff_t size, int tag), where:
 *   tag 0 = Utf8, 1 = Latin1, 2 = Utf16.
 * Slicing results and As*() conversions write through out-parameters.
 */
#ifndef GO_QT6CORE_QANYSTRINGVIEW_C_H
#define GO_QT6CORE_QANYSTRINGVIEW_C_H

#include <stddef.h>
#include <stdint.h>
#include <stdbool.h>
#include "qstring_c.h"

#ifdef __cplusplus
extern "C" {
#endif

/* --- default (empty) constructor --- */
void qanystringview_default(const void **out_data, ptrdiff_t *out_size, int *out_tag);

/* --- introspection --- */
ptrdiff_t qanystringview_size     (const void *data, ptrdiff_t size, int tag);
bool      qanystringview_is_null  (const void *data, ptrdiff_t size, int tag);
bool      qanystringview_is_empty (const void *data, ptrdiff_t size, int tag);
bool      qanystringview_is_utf8  (const void *data, ptrdiff_t size, int tag);
bool      qanystringview_is_latin1(const void *data, ptrdiff_t size, int tag);
bool      qanystringview_is_utf16 (const void *data, ptrdiff_t size, int tag);

/* --- element access --- */
uint16_t qanystringview_front(const void *data, ptrdiff_t size, int tag);
uint16_t qanystringview_back (const void *data, ptrdiff_t size, int tag);

/* --- slicing (new view via out params; tag may change only for Mid→Null) --- */
void qanystringview_left    (const void *data, ptrdiff_t size, int tag, ptrdiff_t n,
                             const void **out_data, ptrdiff_t *out_size, int *out_tag);
void qanystringview_right   (const void *data, ptrdiff_t size, int tag, ptrdiff_t n,
                             const void **out_data, ptrdiff_t *out_size, int *out_tag);
void qanystringview_mid     (const void *data, ptrdiff_t size, int tag, ptrdiff_t pos, ptrdiff_t n,
                             const void **out_data, ptrdiff_t *out_size, int *out_tag);
void qanystringview_first   (const void *data, ptrdiff_t size, int tag, ptrdiff_t n,
                             const void **out_data, ptrdiff_t *out_size, int *out_tag);
void qanystringview_last    (const void *data, ptrdiff_t size, int tag, ptrdiff_t n,
                             const void **out_data, ptrdiff_t *out_size, int *out_tag);
void qanystringview_sliced  (const void *data, ptrdiff_t size, int tag, ptrdiff_t pos,
                             const void **out_data, ptrdiff_t *out_size, int *out_tag);
void qanystringview_sliced_n(const void *data, ptrdiff_t size, int tag, ptrdiff_t pos, ptrdiff_t n,
                             const void **out_data, ptrdiff_t *out_size, int *out_tag);
void qanystringview_chopped (const void *data, ptrdiff_t size, int tag, ptrdiff_t n,
                             const void **out_data, ptrdiff_t *out_size, int *out_tag);

/* --- dispatch to concrete type (out params; null data + 0 size when tag doesn't match) --- */
void qanystringview_as_qstringview      (const void *data, ptrdiff_t size, int tag,
                                         const uint16_t **out_data, ptrdiff_t *out_size);
void qanystringview_as_qlatin1stringview(const void *data, ptrdiff_t size, int tag,
                                         const char **out_data, ptrdiff_t *out_size);
void qanystringview_as_qutf8stringview  (const void *data, ptrdiff_t size, int tag,
                                         const char **out_data, ptrdiff_t *out_size);

/* --- comparison (cs: 0 = CaseSensitive, 1 = CaseInsensitive) --- */
int  qanystringview_compare(const void *ldata, ptrdiff_t lsize, int ltag,
                            const void *rdata, ptrdiff_t rsize, int rtag, int cs);
bool qanystringview_equal  (const void *ldata, ptrdiff_t lsize, int ltag,
                            const void *rdata, ptrdiff_t rsize, int rtag);

/* --- conversion to QString (returns new QString* via out pointer) --- */
go_qt6core_qstring_t* qanystringview_to_qstring(const void *data, ptrdiff_t size, int kind);

#ifdef __cplusplus
}
#endif

#endif
