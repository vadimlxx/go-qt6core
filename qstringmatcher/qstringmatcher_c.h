/* C ABI for Qt 6.11.2 QtCore/QStringMatcher. */
#ifndef GO_QT6CORE_QSTRINGMATCHER_C_H
#define GO_QT6CORE_QSTRINGMATCHER_C_H

#include <stddef.h>
#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

typedef struct go_qsm_s go_qsm_t; /* QStringMatcher* */

/* ---- Constructors (cs: 0 = CaseInsensitive, 1 = CaseSensitive) ---- */
go_qsm_t *qstringmatcher_new_default(void);
go_qsm_t *qstringmatcher_new_from_qstring(const void *qstr, int cs);
go_qsm_t *qstringmatcher_new_from_utf16(const uint16_t *data, ptrdiff_t size, int cs);
go_qsm_t *qstringmatcher_new_copy(const go_qsm_t *other);

/* ---- Lifecycle ---- */
void qstringmatcher_delete(go_qsm_t *self);
void qstringmatcher_assign(go_qsm_t *self, const go_qsm_t *other);

/* ---- Pattern / sensitivity ---- */
void qstringmatcher_set_pattern(go_qsm_t *self, const void *qstr);
void qstringmatcher_set_case_sensitivity(go_qsm_t *self, int cs);
int  qstringmatcher_case_sensitivity(const go_qsm_t *self);

/* ---- Search (returns qsizetype as ptrdiff_t; -1 = not found) ---- */
ptrdiff_t qstringmatcher_index_in_utf16(const go_qsm_t *self,
    const uint16_t *data, ptrdiff_t size, ptrdiff_t from);

/* ---- Pattern accessors ---- */
/* qstringmatcher_pattern allocates a new QString* (owning) via *out_qstring. */
void qstringmatcher_pattern(const go_qsm_t *self, void **out_qstring);
/* qstringmatcher_pattern_view returns a borrowed pointer into the matcher's internal
 * buffer. Valid until the matcher is modified or deleted. */
void qstringmatcher_pattern_view(const go_qsm_t *self,
    const uint16_t **out_data, ptrdiff_t *out_size);

#ifdef __cplusplus
}
#endif

#endif