/* C ABI for Qt 6.11 QtCore/QStringMatcher. */
#ifndef GO_QT6CORE_QSTRINGMATCHER_C_H
#define GO_QT6CORE_QSTRINGMATCHER_C_H

#include <stdint.h>
#include <stdbool.h>
#include <stddef.h>
#include "qstring_c.h"  /* QString forward decl */

#ifdef __cplusplus
class QStringMatcher;
extern "C" {
#else
typedef struct QStringMatcher QStringMatcher;
#endif

/* Constructors */
QStringMatcher* qstringmatcher_new(void);
QStringMatcher* qstringmatcher_new_from_qstring(const QString* pattern, int cs);
QStringMatcher* qstringmatcher_new_from_view(const uint16_t* uc, ptrdiff_t len, int cs);
QStringMatcher* qstringmatcher_new_copy(const QStringMatcher* other);

/* Destructor */
void qstringmatcher_delete(QStringMatcher* self);

/* Setters */
void qstringmatcher_set_pattern(QStringMatcher* self, const QString* pattern);
void qstringmatcher_set_case_sensitivity(QStringMatcher* self, int cs);
void qstringmatcher_assign(QStringMatcher* self, const QStringMatcher* other);

/* indexIn */
ptrdiff_t qstringmatcher_index_in_qstring(const QStringMatcher* self,
                                          const QString* str, ptrdiff_t from);
ptrdiff_t qstringmatcher_index_in_view(const QStringMatcher* self,
                                       const uint16_t* uc, ptrdiff_t len,
                                       ptrdiff_t from);

/* Getters */
QString*         qstringmatcher_pattern(const QStringMatcher* self);
const uint16_t*  qstringmatcher_pattern_view(const QStringMatcher* self,
                                             ptrdiff_t* out_len);
int              qstringmatcher_case_sensitivity(const QStringMatcher* self);

#ifdef __cplusplus
}
#endif

#endif /* GO_QT6CORE_QSTRINGMATCHER_C_H */