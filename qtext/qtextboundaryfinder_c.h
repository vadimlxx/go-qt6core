/* C ABI for Qt 6.11 QtCore/QTextBoundaryFinder. */
#ifndef GO_QT6CORE_QTEXTBOUNDARYFINDER_C_H
#define GO_QT6CORE_QTEXTBOUNDARYFINDER_C_H

#include <stddef.h>
#include <stdint.h>
#include <stdbool.h>

#ifdef __cplusplus
class QString;
class QTextBoundaryFinder;
extern "C" {
#else
typedef struct QString QString;
typedef struct QTextBoundaryFinder QTextBoundaryFinder;
#endif

/* --- lifecycle --- */
QTextBoundaryFinder* qtextboundaryfinder_new_empty(void);
QTextBoundaryFinder* qtextboundaryfinder_new_copy(const QTextBoundaryFinder* other);
QTextBoundaryFinder* qtextboundaryfinder_new_from_qstring(int type, const QString* s);
QTextBoundaryFinder* qtextboundaryfinder_new_from_view(int type, const uint16_t* data, ptrdiff_t size);
void qtextboundaryfinder_delete(QTextBoundaryFinder* self);
void qtextboundaryfinder_swap(QTextBoundaryFinder* self, QTextBoundaryFinder* other);

/* --- accessors --- */
bool     qtextboundaryfinder_is_valid(const QTextBoundaryFinder* self);
int      qtextboundaryfinder_type(const QTextBoundaryFinder* self);
QString* qtextboundaryfinder_string(const QTextBoundaryFinder* self);

/* --- navigation --- */
void      qtextboundaryfinder_to_start(QTextBoundaryFinder* self);
void      qtextboundaryfinder_to_end(QTextBoundaryFinder* self);
ptrdiff_t qtextboundaryfinder_position(const QTextBoundaryFinder* self);
void      qtextboundaryfinder_set_position(QTextBoundaryFinder* self, ptrdiff_t pos);
ptrdiff_t qtextboundaryfinder_to_next_boundary(QTextBoundaryFinder* self);
ptrdiff_t qtextboundaryfinder_to_previous_boundary(QTextBoundaryFinder* self);

/* --- predicates --- */
bool qtextboundaryfinder_is_at_boundary(const QTextBoundaryFinder* self);
int  qtextboundaryfinder_boundary_reasons(const QTextBoundaryFinder* self);

#ifdef __cplusplus
}
#endif

#endif