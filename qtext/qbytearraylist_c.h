/* C ABI for Qt 6.11 QtCore/QByteArrayList (QList<QByteArray>). */
#ifndef GO_QT6CORE_QBYTEARRAYLIST_C_H
#define GO_QT6CORE_QBYTEARRAYLIST_C_H

#include <stddef.h>
#include <stdint.h>
#include <stdbool.h>

/* Forward-declare Qt types before extern "C". */
#ifdef __cplusplus
class QByteArray;
template <typename T> class QList;
typedef QList<QByteArray> QByteArrayList;
extern "C" {
#else
typedef struct QByteArray QByteArray;
typedef struct QByteArrayList QByteArrayList;
#endif

/* --- lifecycle --- */
QByteArrayList* qbytearraylist_new(void);
QByteArrayList* qbytearraylist_new_copy(const QByteArrayList* other);
void            qbytearraylist_delete(QByteArrayList* self);

/* --- capacity / basic --- */
ptrdiff_t qbytearraylist_size(const QByteArrayList* self);
bool      qbytearraylist_is_empty(const QByteArrayList* self);
ptrdiff_t qbytearraylist_capacity(const QByteArrayList* self);

void qbytearraylist_clear(QByteArrayList* self);
void qbytearraylist_reserve(QByteArrayList* self, ptrdiff_t n);
void qbytearraylist_squeeze(QByteArrayList* self);

/* --- append / prepend / insert --- */
void qbytearraylist_append(QByteArrayList* self, const QByteArray* ba);
void qbytearraylist_prepend(QByteArrayList* self, const QByteArray* ba);
void qbytearraylist_insert(QByteArrayList* self, ptrdiff_t i, const QByteArray* ba);

/* --- element access (returns a fresh copy) --- */
void qbytearraylist_at_new(const QByteArrayList* self, ptrdiff_t i, QByteArray** out);
void qbytearraylist_set(QByteArrayList* self, ptrdiff_t i, const QByteArray* ba);

/* --- remove --- */
void qbytearraylist_remove_at(QByteArrayList* self, ptrdiff_t i);
void qbytearraylist_remove_first(QByteArrayList* self);
void qbytearraylist_remove_last(QByteArrayList* self);

/* --- search --- */
ptrdiff_t qbytearraylist_index_of(const QByteArrayList* self, const QByteArray* ba);
ptrdiff_t qbytearraylist_last_index_of(const QByteArrayList* self, const QByteArray* ba);
bool      qbytearraylist_contains(const QByteArrayList* self, const QByteArray* ba);

/* --- QByteArrayList-specific: join --- */
void qbytearraylist_join(const QByteArrayList* self,
                         const char* sep, ptrdiff_t sep_len,
                         QByteArray** out);

/* --- bulk construction / extraction --- */
QByteArrayList* qbytearraylist_from_array(const QByteArray* const* bas, ptrdiff_t n);
void qbytearraylist_to_array_new(const QByteArrayList* self,
                                 QByteArray*** out_arr,
                                 ptrdiff_t* out_len);
void qbytearraylist_free_array(QByteArray** arr);

#ifdef __cplusplus
}
#endif

#endif