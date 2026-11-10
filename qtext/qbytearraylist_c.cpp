/* Order: real Qt headers first, then our C ABI wrapper. */
#include "../internal/qt6core/check_qt_version.h"

#include <QByteArray>
#include <QByteArrayView>
#include <QList>
#include <QByteArrayList>

#include <cstdlib>

#include "qbytearraylist_c.h"

extern "C" {

/* --- lifecycle --- */
QByteArrayList* qbytearraylist_new() {
    return new QByteArrayList();
}
QByteArrayList* qbytearraylist_new_copy(const QByteArrayList* other) {
    return new QByteArrayList(*other);
}
void qbytearraylist_delete(QByteArrayList* self) {
    delete self;
}

/* --- capacity / basic --- */
ptrdiff_t qbytearraylist_size(const QByteArrayList* self) {
    return static_cast<ptrdiff_t>(self->size());
}
bool qbytearraylist_is_empty(const QByteArrayList* self) {
    return self->isEmpty();
}
ptrdiff_t qbytearraylist_capacity(const QByteArrayList* self) {
    return static_cast<ptrdiff_t>(self->capacity());
}
void qbytearraylist_clear(QByteArrayList* self) {
    self->clear();
}
void qbytearraylist_reserve(QByteArrayList* self, ptrdiff_t n) {
    self->reserve(static_cast<qsizetype>(n));
}
void qbytearraylist_squeeze(QByteArrayList* self) {
    self->squeeze();
}

/* --- append / prepend / insert --- */
void qbytearraylist_append(QByteArrayList* self, const QByteArray* ba) {
    self->append(*ba);
}
void qbytearraylist_prepend(QByteArrayList* self, const QByteArray* ba) {
    self->prepend(*ba);
}
void qbytearraylist_insert(QByteArrayList* self, ptrdiff_t i, const QByteArray* ba) {
    self->insert(static_cast<qsizetype>(i), *ba);
}

/* --- element access ---
 * at() возвращает const T&; создаём копию, чтобы изолировать
 * Go-пользователя от мутаций списка. out-параметр типизирован
 * как QByteArray**. */
void qbytearraylist_at_new(const QByteArrayList* self, ptrdiff_t i, QByteArray** out) {
    *out = new QByteArray(self->at(static_cast<qsizetype>(i)));
}
void qbytearraylist_set(QByteArrayList* self, ptrdiff_t i, const QByteArray* ba) {
    (*self)[static_cast<qsizetype>(i)] = *ba;
}

/* --- remove --- */
void qbytearraylist_remove_at(QByteArrayList* self, ptrdiff_t i) {
    self->removeAt(static_cast<qsizetype>(i));
}
void qbytearraylist_remove_first(QByteArrayList* self) {
    self->removeFirst();
}
void qbytearraylist_remove_last(QByteArrayList* self) {
    self->removeLast();
}

/* --- search --- */
ptrdiff_t qbytearraylist_index_of(const QByteArrayList* self, const QByteArray* ba) {
    return static_cast<ptrdiff_t>(self->indexOf(*ba));
}
ptrdiff_t qbytearraylist_last_index_of(const QByteArrayList* self, const QByteArray* ba) {
    return static_cast<ptrdiff_t>(self->lastIndexOf(*ba));
}
bool qbytearraylist_contains(const QByteArrayList* self, const QByteArray* ba) {
    return self->contains(*ba);
}

/* --- join (QByteArrayList specialization) --- */
void qbytearraylist_join(const QByteArrayList* self,
                         const char* sep, ptrdiff_t sep_len,
                         QByteArray** out) {
    QByteArrayView v(sep, static_cast<qsizetype>(sep_len));
    *out = new QByteArray(self->join(v));
}

/* --- bulk construction / extraction --- */
QByteArrayList* qbytearraylist_from_array(const QByteArray* const* bas, ptrdiff_t n) {
    auto* l = new QByteArrayList();
    l->reserve(static_cast<qsizetype>(n));
    for (ptrdiff_t i = 0; i < n; ++i)
        l->append(*bas[i]);
    return l;
}

void qbytearraylist_to_array_new(const QByteArrayList* self,
                                 QByteArray*** out_arr,
                                 ptrdiff_t* out_len) {
    const qsizetype n = self->size();
    *out_len = static_cast<ptrdiff_t>(n);
    if (n == 0) {
        *out_arr = nullptr;
        return;
    }
    QByteArray** arr = static_cast<QByteArray**>(
        std::malloc(sizeof(QByteArray*) * static_cast<size_t>(n)));
    for (qsizetype i = 0; i < n; ++i)
        arr[i] = new QByteArray(self->at(i));
    *out_arr = arr;
}
void qbytearraylist_free_array(QByteArray** arr) {
    std::free(arr);
}

} // extern "C"
