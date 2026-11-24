/* Order: version check first, then real Qt headers. */
#include "../internal/qt6core/check_qt_version.h"

/* qcollator.h declares both QCollator and QCollatorSortKey; it is the
 * canonical header for the pair, mirroring QtCore's qcollator.cpp layout. */
#include <QCollator>
#include <QLocale>
#include <QString>
#include <QStringView>

#include "qcollator_c.h"

namespace {

/* Go UTF-16 strings are passed as (uint16_t*, length). QChar and char16_t are
 * both 16-bit code units; this helper is the single conversion point. */
inline QStringView V(const uint16_t *data, ptrdiff_t len) noexcept {
    return QStringView(reinterpret_cast<const char16_t *>(data),
                       static_cast<qsizetype>(len));
}

/* Qt::CaseSensitivity crosses the ABI as an int; this is the single
 * conversion point (qdata.CaseSensitivity mirrors the same values). */
inline Qt::CaseSensitivity CS(int v) noexcept {
    return static_cast<Qt::CaseSensitivity>(v);
}

/* QCollator::sortKey() only has a const QString & overload (unlike compare()),
 * so the UTF-16 payload must be materialized as a QString. */
inline QString S(const uint16_t *data, ptrdiff_t len) {
    return QString::fromUtf16(reinterpret_cast<const char16_t *>(data),
                              static_cast<qsizetype>(len));
}

} // namespace

extern "C" {

/* ==== QCollator ==== */

go_qcollator_t *qcollator_new_default(void) {
    return new QCollator();
}
go_qcollator_t *qcollator_new_from_locale(const go_qlocale_t locale) {
    return new QCollator(locale ? *locale : QLocale());
}
go_qcollator_t *qcollator_new_copy(const go_qcollator_t *other) {
    return new QCollator(*other);
}
void qcollator_delete(go_qcollator_t *self) {
    delete self;
}
void qcollator_swap(go_qcollator_t *self, go_qcollator_t *other) {
    self->swap(*other);
}

void qcollator_set_locale(go_qcollator_t *self, const go_qlocale_t locale) {
    self->setLocale(locale ? *locale : QLocale());
}
go_qlocale_t qcollator_locale(const go_qcollator_t *self) {
    return new QLocale(self->locale());
}

int qcollator_case_sensitivity(const go_qcollator_t *self) {
    return static_cast<int>(self->caseSensitivity());
}
void qcollator_set_case_sensitivity(go_qcollator_t *self, int cs) {
    self->setCaseSensitivity(CS(cs));
}

bool qcollator_numeric_mode(const go_qcollator_t *self) {
    return self->numericMode();
}
void qcollator_set_numeric_mode(go_qcollator_t *self, bool on) {
    self->setNumericMode(on);
}

bool qcollator_ignore_punctuation(const go_qcollator_t *self) {
    return self->ignorePunctuation();
}
void qcollator_set_ignore_punctuation(go_qcollator_t *self, bool on) {
    self->setIgnorePunctuation(on);
}

int qcollator_compare(const go_qcollator_t *self,
    const uint16_t *s1, ptrdiff_t len1,
    const uint16_t *s2, ptrdiff_t len2) {
    return self->compare(V(s1, len1), V(s2, len2));
}

go_qcollatorsortkey_t *qcollator_sort_key(const go_qcollator_t *self,
    const uint16_t *s, ptrdiff_t len) {
    return new QCollatorSortKey(self->sortKey(S(s, len)));
}

int qcollator_default_compare(
    const uint16_t *s1, ptrdiff_t len1,
    const uint16_t *s2, ptrdiff_t len2) {
    return QCollator::defaultCompare(V(s1, len1), V(s2, len2));
}

go_qcollatorsortkey_t *qcollator_default_sort_key(
    const uint16_t *s, ptrdiff_t len) {
    return new QCollatorSortKey(QCollator::defaultSortKey(V(s, len)));
}

/* ==== QCollatorSortKey ==== */

go_qcollatorsortkey_t *qcollatorsortkey_new_copy(const go_qcollatorsortkey_t *other) {
    return new QCollatorSortKey(*other);
}
void qcollatorsortkey_delete(go_qcollatorsortkey_t *self) {
    delete self;
}
void qcollatorsortkey_swap(go_qcollatorsortkey_t *self,
                           go_qcollatorsortkey_t *other) {
    self->swap(*other);
}

int qcollatorsortkey_compare(const go_qcollatorsortkey_t *self,
    const go_qcollatorsortkey_t *other) {
    return self->compare(*other);
}

bool qcollatorsortkey_less(const go_qcollatorsortkey_t *self,
    const go_qcollatorsortkey_t *other) {
    return self->compare(*other) < 0;
}

} // extern "C"
