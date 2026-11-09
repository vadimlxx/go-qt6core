/* C ABI for Qt 6.11 QtCore/QByteArrayList (QList<QByteArray>). */
#ifndef GO_QT6CORE_QBYTEARRAYLIST_C_H
#define GO_QT6CORE_QBYTEARRAYLIST_C_H

#include <stddef.h>
#include <stdint.h>
#include <stdbool.h>

#ifdef __cplusplus
extern "C" {
#endif

/* Opaque handle types. The .cpp casts them to real C++ types
 * (QByteArrayList* and QByteArray*) via reinterpret_cast. */
typedef struct go_qt6core_qbytearray_s    go_qt6core_qbytearray_t;
typedef struct go_qt6core_qbytearraylist_s go_qt6core_qbytearraylist_t;

/* --- lifecycle --- */
go_qt6core_qbytearraylist_t* qbytearraylist_new(void);
go_qt6core_qbytearraylist_t* qbytearraylist_new_copy(const go_qt6core_qbytearraylist_t *other);
void                         qbytearraylist_delete(go_qt6core_qbytearraylist_t *self);

/* --- capacity / basic --- */
ptrdiff_t qbytearraylist_size(const go_qt6core_qbytearraylist_t *self);
bool      qbytearraylist_is_empty(const go_qt6core_qbytearraylist_t *self);
ptrdiff_t qbytearraylist_capacity(const go_qt6core_qbytearraylist_t *self);

void qbytearraylist_clear(go_qt6core_qbytearraylist_t *self);
void qbytearraylist_reserve(go_qt6core_qbytearraylist_t *self, ptrdiff_t n);
void qbytearraylist_squeeze(go_qt6core_qbytearraylist_t *self);

/* --- append / prepend / insert --- */
void qbytearraylist_append(go_qt6core_qbytearraylist_t *self,
                           const go_qt6core_qbytearray_t *ba);
void qbytearraylist_prepend(go_qt6core_qbytearraylist_t *self,
                            const go_qt6core_qbytearray_t *ba);
void qbytearraylist_insert(go_qt6core_qbytearraylist_t *self,
                           ptrdiff_t i,
                           const go_qt6core_qbytearray_t *ba);

/* --- element access (returns a fresh copy) --- */
void qbytearraylist_at_new(const go_qt6core_qbytearraylist_t *self,
                           ptrdiff_t i,
                           go_qt6core_qbytearray_t **out);
void qbytearraylist_set(go_qt6core_qbytearraylist_t *self,
                        ptrdiff_t i,
                        const go_qt6core_qbytearray_t *ba);

/* --- remove --- */
void qbytearraylist_remove_at(go_qt6core_qbytearraylist_t *self, ptrdiff_t i);
void qbytearraylist_remove_first(go_qt6core_qbytearraylist_t *self);
void qbytearraylist_remove_last(go_qt6core_qbytearraylist_t *self);

/* --- search (QListSpecialMethodsBase<QByteArray> additions) --- */
ptrdiff_t qbytearraylist_index_of(const go_qt6core_qbytearraylist_t *self,
                                  const go_qt6core_qbytearray_t *ba);
ptrdiff_t qbytearraylist_last_index_of(const go_qt6core_qbytearraylist_t *self,
                                       const go_qt6core_qbytearray_t *ba);
bool      qbytearraylist_contains(const go_qt6core_qbytearraylist_t *self,
                                  const go_qt6core_qbytearray_t *ba);

/* --- QByteArrayList-specific: join --- */
void qbytearraylist_join(const go_qt6core_qbytearraylist_t *self,
                         const char *sep, ptrdiff_t sep_len,
                         go_qt6core_qbytearray_t **out);

/* --- bulk construction / extraction --- */
go_qt6core_qbytearraylist_t* qbytearraylist_from_array(
    const go_qt6core_qbytearray_t *const *bas, ptrdiff_t n);
void qbytearraylist_to_array_new(const go_qt6core_qbytearraylist_t *self,
                                 go_qt6core_qbytearray_t ***out_arr,
                                 ptrdiff_t *out_len);
void qbytearraylist_free_array(go_qt6core_qbytearray_t **arr);

#ifdef __cplusplus
}
#endif

#endif