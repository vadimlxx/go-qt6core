/* Order: include the C ABI header first; it pulls in <QtCore/QStringList>
 * (which transitively provides QString). Then add the headers whose full
 * definitions are needed by the implementations below. */
#include "qstringlist_c.h"
#include <QtCore/QRegularExpression>

extern "C" {

/* --- lifecycle --- */
QStringList* qstringlist_new_empty() { return new QStringList(); }
QStringList* qstringlist_new_copy(const QStringList* other) { return new QStringList(*other); }
QStringList* qstringlist_new_from_string(const QString* s) { return new QStringList(*s); }
void qstringlist_delete(QStringList* self) { delete self; }
void qstringlist_swap(QStringList* self, QStringList* other) { self->swap(*other); }

/* --- basic QList<QString> ops --- */
ptrdiff_t qstringlist_size(const QStringList* self) { return self->size(); }
bool qstringlist_is_empty(const QStringList* self) { return self->isEmpty(); }
void qstringlist_clear(QStringList* self) { self->clear(); }
void qstringlist_reserve(QStringList* self, ptrdiff_t n) { self->reserve(n); }
void qstringlist_squeeze(QStringList* self) { self->squeeze(); }

QString* qstringlist_at(const QStringList* self, ptrdiff_t i) { return new QString(self->at(i)); }
void qstringlist_set_at(QStringList* self, ptrdiff_t i, const QString* s) { (*self)[i] = *s; }
QString* qstringlist_first(const QStringList* self) { return new QString(self->first()); }
QString* qstringlist_last(const QStringList* self) { return new QString(self->last()); }

void qstringlist_append(QStringList* self, const QString* s) { self->append(*s); }
void qstringlist_prepend(QStringList* self, const QString* s) { self->prepend(*s); }
void qstringlist_insert(QStringList* self, ptrdiff_t i, const QString* s) { self->insert(i, *s); }
void qstringlist_remove_at(QStringList* self, ptrdiff_t i) { self->removeAt(i); }
void qstringlist_remove_first(QStringList* self) { self->removeFirst(); }
void qstringlist_remove_last(QStringList* self) { self->removeLast(); }
QString* qstringlist_take_at(QStringList* self, ptrdiff_t i) { return new QString(self->takeAt(i)); }
QString* qstringlist_take_first(QStringList* self) { return new QString(self->takeFirst()); }
QString* qstringlist_take_last(QStringList* self) { return new QString(self->takeLast()); }

QString** qstringlist_to_array(const QStringList* self, ptrdiff_t* out_count) {
    const auto n = static_cast<size_t>(self->size());
    auto** arr = new QString*[n > 0 ? n : 1];
    for (size_t i = 0; i < n; ++i)
        arr[i] = new QString((*self)[static_cast<qsizetype>(i)]);
    *out_count = static_cast<ptrdiff_t>(n);
    return arr;
}
void qstringlist_free_string_array(QString** arr) { delete[] arr; }

QStringList* qstringlist_from_array(const QString* const* arr, ptrdiff_t count) {
    auto* out = new QStringList();
    out->reserve(count);
    for (ptrdiff_t i = 0; i < count; ++i)
        out->append(*arr[i]);
    return out;
}

/* --- QStringList-specific --- */
void qstringlist_sort(QStringList* self, int cs) {
    self->sort(cs == 0 ? Qt::CaseSensitive : Qt::CaseInsensitive);
}
ptrdiff_t qstringlist_remove_duplicates(QStringList* self) { return self->removeDuplicates(); }

QString* qstringlist_join_qstring(const QStringList* self, const QString* sep) {
    return new QString(self->join(*sep));
}
QString* qstringlist_join_qchar(const QStringList* self, uint16_t sep) {
    return new QString(self->join(QChar(sep)));
}

QStringList* qstringlist_filter_qstring(const QStringList* self, const QString* needle, int cs) {
    return new QStringList(self->filter(*needle, cs == 0 ? Qt::CaseSensitive : Qt::CaseInsensitive));
}
void qstringlist_replace_in_strings_qstring(QStringList* self, const QString* before,
                                            const QString* after, int cs) {
    self->replaceInStrings(*before, *after, cs == 0 ? Qt::CaseSensitive : Qt::CaseInsensitive);
}

bool qstringlist_contains_qstring(const QStringList* self, const QString* s, int cs) {
    return self->contains(*s, cs == 0 ? Qt::CaseSensitive : Qt::CaseInsensitive);
}
ptrdiff_t qstringlist_index_of_qstring(const QStringList* self, const QString* s,
                                       ptrdiff_t from, int cs) {
    return self->indexOf(*s, from, cs == 0 ? Qt::CaseSensitive : Qt::CaseInsensitive);
}
ptrdiff_t qstringlist_last_index_of_qstring(const QStringList* self, const QString* s,
                                            ptrdiff_t from, int cs) {
    return self->lastIndexOf(*s, from, cs == 0 ? Qt::CaseSensitive : Qt::CaseInsensitive);
}

QStringList* qstringlist_filter_qregexp(const QStringList* self, const QRegularExpression* re) {
    return new QStringList(self->filter(*re));
}
void qstringlist_replace_in_strings_qregexp(QStringList* self, const QRegularExpression* re,
                                            const QString* after) {
    self->replaceInStrings(*re, *after);
}
ptrdiff_t qstringlist_index_of_qregexp(const QStringList* self, const QRegularExpression* re,
                                       ptrdiff_t from) {
    return self->indexOf(*re, from);
}
ptrdiff_t qstringlist_last_index_of_qregexp(const QStringList* self, const QRegularExpression* re,
                                            ptrdiff_t from) {
    return self->lastIndexOf(*re, from);
}

} // extern "C"