/* C ABI for Qt 6.11 QtCore/QStringList. */
#ifndef GO_QT6CORE_QSTRINGLIST_C_H
#define GO_QT6CORE_QSTRINGLIST_C_H

#include <stddef.h>
#include <stdint.h>
#include <stdbool.h>

#ifdef __cplusplus
/*
 * In Qt 6, QStringList is not a class but a type alias:
 *     using QStringList = QList<QString>;
 * Therefore it cannot be forward-declared with `class QStringList;`.
 * Include the real header instead; this also brings in QString (via
 * qstring.h) and forward-declares QRegularExpression.
 */
#include <QtCore/QStringList>
extern "C" {
#else
typedef struct QString QString;
typedef struct QStringList QStringList;
typedef struct QRegularExpression QRegularExpression;
#endif

/* --- lifecycle --- */
QStringList* qstringlist_new_empty(void);
QStringList* qstringlist_new_copy(const QStringList* other);
QStringList* qstringlist_new_from_string(const QString* s);
void         qstringlist_delete(QStringList* self);
void         qstringlist_swap(QStringList* self, QStringList* other);

/* --- basic QList<QString> ops --- */
ptrdiff_t qstringlist_size(const QStringList* self);
bool      qstringlist_is_empty(const QStringList* self);
void      qstringlist_clear(QStringList* self);
void      qstringlist_reserve(QStringList* self, ptrdiff_t n);
void      qstringlist_squeeze(QStringList* self);

QString* qstringlist_at(const QStringList* self, ptrdiff_t i);
void     qstringlist_set_at(QStringList* self, ptrdiff_t i, const QString* s);
QString* qstringlist_first(const QStringList* self);
QString* qstringlist_last(const QStringList* self);

void qstringlist_append(QStringList* self, const QString* s);
void qstringlist_prepend(QStringList* self, const QString* s);
void qstringlist_insert(QStringList* self, ptrdiff_t i, const QString* s);
void qstringlist_remove_at(QStringList* self, ptrdiff_t i);
void qstringlist_remove_first(QStringList* self);
void qstringlist_remove_last(QStringList* self);
QString* qstringlist_take_at(QStringList* self, ptrdiff_t i);
QString* qstringlist_take_first(QStringList* self);
QString* qstringlist_take_last(QStringList* self);

/* ToSlice: returns heap QString*[]; free outer array with qstringlist_free_string_array.
   Each QString* is owned by the caller (wrapped by Go). */
QString** qstringlist_to_array(const QStringList* self, ptrdiff_t* out_count);
void      qstringlist_free_string_array(QString** arr);

/* FromArray: builds a QStringList by copying each element. */
QStringList* qstringlist_from_array(const QString* const* arr, ptrdiff_t count);

/* --- QStringList-specific --- */

/* cs: 0 = CaseSensitive, 1 = CaseInsensitive */
void      qstringlist_sort(QStringList* self, int cs);
ptrdiff_t qstringlist_remove_duplicates(QStringList* self);

/* join overloads */
QString* qstringlist_join_qstring(const QStringList* self, const QString* sep);
QString* qstringlist_join_qchar(const QStringList* self, uint16_t sep);

/* filter / replaceInStrings overloads (QString + cs) */
QStringList* qstringlist_filter_qstring(const QStringList* self, const QString* needle, int cs);
void         qstringlist_replace_in_strings_qstring(QStringList* self, const QString* before,
                                                    const QString* after, int cs);

/* contains / indexOf / lastIndexOf (QString + cs) */
bool      qstringlist_contains_qstring(const QStringList* self, const QString* s, int cs);
ptrdiff_t qstringlist_index_of_qstring(const QStringList* self, const QString* s,
                                       ptrdiff_t from, int cs);
ptrdiff_t qstringlist_last_index_of_qstring(const QStringList* self, const QString* s,
                                            ptrdiff_t from, int cs);

/* QRegularExpression interop (requires QT_CONFIG(regularexpression)) */
QStringList* qstringlist_filter_qregexp(const QStringList* self, const QRegularExpression* re);
void         qstringlist_replace_in_strings_qregexp(QStringList* self, const QRegularExpression* re,
                                                    const QString* after);
ptrdiff_t    qstringlist_index_of_qregexp(const QStringList* self, const QRegularExpression* re,
                                          ptrdiff_t from);
ptrdiff_t    qstringlist_last_index_of_qregexp(const QStringList* self, const QRegularExpression* re,
                                               ptrdiff_t from);

#ifdef __cplusplus
}
#endif

#endif