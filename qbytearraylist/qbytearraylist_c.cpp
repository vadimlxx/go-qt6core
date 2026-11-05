/* Order: real Qt headers first, then our C ABI wrapper. */
#include <QtCore/qglobal.h>
#if !defined(QT_VERSION) || QT_VERSION < QT_VERSION_CHECK(6, 11, 0)
#  error "go-qt6core requires Qt >= 6.11"
#endif

#include <QByteArray>
#include <QByteArrayView>
#include <QList>
#include <QByteArrayList>

#include <cstdlib>

#include "qbytearraylist_c.h"

namespace {

inline QByteArrayList       *L(go_qt6core_qbytearraylist_t       *p) { return reinterpret_cast<QByteArrayList *>(p); }
inline const QByteArrayList *L(const go_qt6core_qbytearraylist_t *p) { return reinterpret_cast<const QByteArrayList *>(p); }

inline QByteArray       *Q(go_qt6core_qbytearray_t       *p) { return reinterpret_cast<QByteArray *>(p); }
inline const QByteArray *Q(const go_qt6core_qbytearray_t *p) { return reinterpret_cast<const QByteArray *>(p); }

inline go_qt6core_qbytearraylist_t *toHandle(QByteArrayList *l) {
    return reinterpret_cast<go_qt6core_qbytearraylist_t *>(l);
}
inline go_qt6core_qbytearray_t *toHandle(QByteArray *q) {
    return reinterpret_cast<go_qt6core_qbytearray_t *>(q);
}

} // namespace

extern "C" {

/* --- lifecycle --- */
go_qt6core_qbytearraylist_t *qbytearraylist_new() {
    return toHandle(new QByteArrayList());
}
go_qt6core_qbytearraylist_t *qbytearraylist_new_copy(const go_qt6core_qbytearraylist_t *other) {
    return toHandle(new QByteArrayList(*L(other)));
}
void qbytearraylist_delete(go_qt6core_qbytearraylist_t *self) {
    delete L(self);
}

/* --- capacity / basic --- */
ptrdiff_t qbytearraylist_size(const go_qt6core_qbytearraylist_t *self) {
    return static_cast<ptrdiff_t>(L(self)->size());
}
bool qbytearraylist_is_empty(const go_qt6core_qbytearraylist_t *self) {
    return L(self)->isEmpty();
}
ptrdiff_t qbytearraylist_capacity(const go_qt6core_qbytearraylist_t *self) {
    return static_cast<ptrdiff_t>(L(self)->capacity());
}
void qbytearraylist_clear(go_qt6core_qbytearraylist_t *self) {
    L(self)->clear();
}
void qbytearraylist_reserve(go_qt6core_qbytearraylist_t *self, ptrdiff_t n) {
    L(self)->reserve(static_cast<qsizetype>(n));
}
void qbytearraylist_squeeze(go_qt6core_qbytearraylist_t *self) {
    L(self)->squeeze();
}

/* --- append / prepend / insert --- */
void qbytearraylist_append(go_qt6core_qbytearraylist_t *self,
                           const go_qt6core_qbytearray_t *ba) {
    L(self)->append(*Q(ba));
}
void qbytearraylist_prepend(go_qt6core_qbytearraylist_t *self,
                            const go_qt6core_qbytearray_t *ba) {
    L(self)->prepend(*Q(ba));
}
void qbytearraylist_insert(go_qt6core_qbytearraylist_t *self,
                           ptrdiff_t i,
                           const go_qt6core_qbytearray_t *ba) {
    L(self)->insert(static_cast<qsizetype>(i), *Q(ba));
}

/* --- element access --- */
void qbytearraylist_at_new(const go_qt6core_qbytearraylist_t *self,
                           ptrdiff_t i,
                           go_qt6core_qbytearray_t **out) {
    *out = toHandle(new QByteArray(L(self)->at(static_cast<qsizetype>(i))));
}
void qbytearraylist_set(go_qt6core_qbytearraylist_t *self,
                        ptrdiff_t i,
                        const go_qt6core_qbytearray_t *ba) {
    (*L(self))[static_cast<qsizetype>(i)] = *Q(ba);
}

/* --- remove --- */
void qbytearraylist_remove_at(go_qt6core_qbytearraylist_t *self, ptrdiff_t i) {
    L(self)->removeAt(static_cast<qsizetype>(i));
}
void qbytearraylist_remove_first(go_qt6core_qbytearraylist_t *self) {
    L(self)->removeFirst();
}
void qbytearraylist_remove_last(go_qt6core_qbytearraylist_t *self) {
    L(self)->removeLast();
}

/* --- search --- */
ptrdiff_t qbytearraylist_index_of(const go_qt6core_qbytearraylist_t *self,
                                  const go_qt6core_qbytearray_t *ba) {
    return static_cast<ptrdiff_t>(L(self)->indexOf(*Q(ba)));
}
ptrdiff_t qbytearraylist_last_index_of(const go_qt6core_qbytearraylist_t *self,
                                       const go_qt6core_qbytearray_t *ba) {
    return static_cast<ptrdiff_t>(L(self)->lastIndexOf(*Q(ba)));
}
bool qbytearraylist_contains(const go_qt6core_qbytearraylist_t *self,
                             const go_qt6core_qbytearray_t *ba) {
    return L(self)->contains(*Q(ba));
}

/* --- join (QByteArrayList specialization) --- */
void qbytearraylist_join(const go_qt6core_qbytearraylist_t *self,
                         const char *sep, ptrdiff_t sep_len,
                         go_qt6core_qbytearray_t **out) {
    QByteArrayView v(sep, static_cast<qsizetype>(sep_len));
    *out = toHandle(new QByteArray(L(self)->join(v)));
}

/* --- bulk construction / extraction --- */
go_qt6core_qbytearraylist_t *qbytearraylist_from_array(
        const go_qt6core_qbytearray_t *const *bas, ptrdiff_t n) {
    auto *l = new QByteArrayList();
    l->reserve(static_cast<qsizetype>(n));
    for (ptrdiff_t i = 0; i < n; ++i)
        l->append(*Q(bas[i]));
    return toHandle(l);
}

void qbytearraylist_to_array_new(const go_qt6core_qbytearraylist_t *self,
                                 go_qt6core_qbytearray_t ***out_arr,
                                 ptrdiff_t *out_len) {
    const qsizetype n = L(self)->size();
    *out_len = static_cast<ptrdiff_t>(n);
    if (n == 0) {
        *out_arr = nullptr;
        return;
    }
    auto **arr = static_cast<go_qt6core_qbytearray_t **>(
        std::malloc(sizeof(go_qt6core_qbytearray_t *) * static_cast<size_t>(n)));
    for (qsizetype i = 0; i < n; ++i)
        arr[i] = toHandle(new QByteArray(L(self)->at(i)));
    *out_arr = arr;
}
void qbytearraylist_free_array(go_qt6core_qbytearray_t **arr) {
    std::free(arr);
}

} // extern "C"