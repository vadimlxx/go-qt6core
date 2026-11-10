/* C ABI for Qt 6.11 QtCore/QUtf8StringView. */
#ifndef GO_QT6CORE_QUTF8STRINGVIEW_C_H
#define GO_QT6CORE_QUTF8STRINGVIEW_C_H

#include <stddef.h>
#include <stdint.h>
#include <stdbool.h>
#include "qstring_c.h"

#ifdef __cplusplus
extern "C" {
#endif

/* --- predicates --- */
bool qutf8stringview_is_null      (const char *data, ptrdiff_t size);
bool qutf8stringview_is_empty     (const char *data, ptrdiff_t size);
bool qutf8stringview_is_valid_utf8(const char *data, ptrdiff_t size);

/* --- element access --- */
char qutf8stringview_at    (const char *data, ptrdiff_t size, ptrdiff_t i);
char qutf8stringview_front (const char *data, ptrdiff_t size);
char qutf8stringview_back  (const char *data, ptrdiff_t size);

/* --- slicing (new view via out params) --- */
void qutf8stringview_left   (const char *data, ptrdiff_t size, ptrdiff_t n,
                             const char **out_data, ptrdiff_t *out_size);
void qutf8stringview_right  (const char *data, ptrdiff_t size, ptrdiff_t n,
                             const char **out_data, ptrdiff_t *out_size);
void qutf8stringview_mid    (const char *data, ptrdiff_t size, ptrdiff_t pos, ptrdiff_t n,
                             const char **out_data, ptrdiff_t *out_size);
void qutf8stringview_sliced (const char *data, ptrdiff_t size, ptrdiff_t pos,
                             const char **out_data, ptrdiff_t *out_size);
void qutf8stringview_sliced_n(const char *data, ptrdiff_t size, ptrdiff_t pos, ptrdiff_t n,
                              const char **out_data, ptrdiff_t *out_size);
void qutf8stringview_first  (const char *data, ptrdiff_t size, ptrdiff_t n,
                             const char **out_data, ptrdiff_t *out_size);
void qutf8stringview_last   (const char *data, ptrdiff_t size, ptrdiff_t n,
                             const char **out_data, ptrdiff_t *out_size);
void qutf8stringview_chopped(const char *data, ptrdiff_t size, ptrdiff_t n,
                             const char **out_data, ptrdiff_t *out_size);

/* --- compare (cs: 0 = CaseSensitive, 1 = CaseInsensitive) --- */
int qutf8stringview_compare_view              (const char *data, ptrdiff_t size,
                                               const char *other, ptrdiff_t other_size, int cs);
int qutf8stringview_compare_qbytearray        (const char *data, ptrdiff_t size,
                                               const void *qba, int cs);
int qutf8stringview_compare_qchar             (const char *data, ptrdiff_t size,
                                               uint16_t c, int cs);
int qutf8stringview_compare_qstringview       (const char *data, ptrdiff_t size,
                                               const uint16_t *sv_data, ptrdiff_t sv_size, int cs);
int qutf8stringview_compare_qlatin1stringview (const char *data, ptrdiff_t size,
                                               const char *l1_data, ptrdiff_t l1_size, int cs);

/* --- equal (always case-sensitive, as in Qt) --- */
bool qutf8stringview_equal_view              (const char *data, ptrdiff_t size,
                                              const char *other, ptrdiff_t other_size);
bool qutf8stringview_equal_qbytearray        (const char *data, ptrdiff_t size,
                                              const void *qba);
bool qutf8stringview_equal_qchar             (const char *data, ptrdiff_t size,
                                              uint16_t c);
bool qutf8stringview_equal_qstringview       (const char *data, ptrdiff_t size,
                                              const uint16_t *sv_data, ptrdiff_t sv_size);
bool qutf8stringview_equal_qlatin1stringview (const char *data, ptrdiff_t size,
                                              const char *l1_data, ptrdiff_t l1_size);

/* --- conversion to owning QString (via void* → qstring.AdoptUnsafePtr) --- */
go_qt6core_qstring_t* qutf8stringview_to_qstring(const char *data, ptrdiff_t size);

/* --- static --- */
ptrdiff_t qutf8stringview_max_size(void);

#ifdef __cplusplus
}
#endif

#endif
