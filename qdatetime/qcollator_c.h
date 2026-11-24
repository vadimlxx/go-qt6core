/* C ABI for Qt 6.11 QtCore/QCollator and QCollatorSortKey.
 *
 * Both classes are declared in QtCore's qcollator.h and implemented in
 * qcollator.cpp, so this bridge follows the same layout: one header and one
 * implementation file for the whole pair.
 *
 * The handle types are the Qt classes themselves (heap-allocated):
 *   - go_qcollator_t          — QCollator*
 *   - go_qcollatorsortkey_t   — QCollatorSortKey*
 *
 * That keeps both ABI layers cast-free: the C functions receive the object
 * pointer directly and call its methods. QLocale is passed as its existing
 * go_qlocale_t handle (a heap QLocale*), see qlocale_c.h.
 *
 * Ownership:
 *   - every qcollator_*_new_* / qcollatorsortkey_*_new_* function returns a
 *     handle owned by the caller; release it with the matching *_delete;
 *   - qcollator_locale() returns a new QLocale handle owned by the caller;
 *   - qcollator_sort_key() and the default/factory helpers return a new
 *     QCollatorSortKey handle owned by the caller.
 */
#ifndef GO_QT6CORE_QDATETIME_QCOLLATOR_C_H
#define GO_QT6CORE_QDATETIME_QCOLLATOR_C_H

#include <stddef.h>
#include <stdint.h>
#include <stdbool.h>

#include "qlocale_c.h"   /* go_qlocale_t */

#ifdef __cplusplus
class QCollator;
class QCollatorSortKey;
typedef QCollator         go_qcollator_t;
typedef QCollatorSortKey  go_qcollatorsortkey_t;
extern "C" {
#else
typedef struct QCollator        go_qcollator_t;
typedef struct QCollatorSortKey go_qcollatorsortkey_t;
#endif

/* ==== QCollator ==== */

/* Qt: QCollator(), QCollator(const QLocale &), QCollator(const QCollator &). */
go_qcollator_t *qcollator_new_default(void);
go_qcollator_t *qcollator_new_from_locale(const go_qlocale_t locale);
go_qcollator_t *qcollator_new_copy(const go_qcollator_t *other);

/* Qt: QCollator::~QCollator(). */
void qcollator_delete(go_qcollator_t *self);

/* Qt: QCollator::swap(QCollator &). */
void qcollator_swap(go_qcollator_t *self, go_qcollator_t *other);

/* Qt: QCollator::setLocale(const QLocale &) / locale() const. */
void            qcollator_set_locale(go_qcollator_t *self, const go_qlocale_t locale);
go_qlocale_t    qcollator_locale(const go_qcollator_t *self);

/* Qt: QCollator::caseSensitivity() const / setCaseSensitivity(Qt::CaseSensitivity). */
int  qcollator_case_sensitivity(const go_qcollator_t *self);
void qcollator_set_case_sensitivity(go_qcollator_t *self, int cs);

/* Qt: QCollator::numericMode() const / setNumericMode(bool). */
bool qcollator_numeric_mode(const go_qcollator_t *self);
void qcollator_set_numeric_mode(go_qcollator_t *self, bool on);

/* Qt: QCollator::ignorePunctuation() const / setIgnorePunctuation(bool). */
bool qcollator_ignore_punctuation(const go_qcollator_t *self);
void qcollator_set_ignore_punctuation(go_qcollator_t *self, bool on);

/* Qt: QCollator::compare(const QString &, const QString &) const and
 * QCollator::compare(QStringView, QStringView) const.
 * The UTF-16 payload is passed as (data, length) code units. */
int qcollator_compare(const go_qcollator_t *self,
    const uint16_t *s1, ptrdiff_t len1,
    const uint16_t *s2, ptrdiff_t len2);

/* Qt: QCollator::sortKey(const QString &) const and
 * QCollator::sortKey(QStringView) const. Returns a new QCollatorSortKey. */
go_qcollatorsortkey_t *qcollator_sort_key(const go_qcollator_t *self,
    const uint16_t *s, ptrdiff_t len);

/* Qt: static QCollator::defaultCompare(QStringView, QStringView). */
int qcollator_default_compare(
    const uint16_t *s1, ptrdiff_t len1,
    const uint16_t *s2, ptrdiff_t len2);

/* Qt: static QCollator::defaultSortKey(QStringView). */
go_qcollatorsortkey_t *qcollator_default_sort_key(
    const uint16_t *s, ptrdiff_t len);

/* ==== QCollatorSortKey ==== */

/* QCollatorSortKey has no public default constructor; only the copy
 * constructor is available through the C ABI.
 * Qt: QCollatorSortKey(const QCollatorSortKey &). */
go_qcollatorsortkey_t *qcollatorsortkey_new_copy(const go_qcollatorsortkey_t *other);

/* Qt: QCollatorSortKey::~QCollatorSortKey(). */
void qcollatorsortkey_delete(go_qcollatorsortkey_t *self);

/* Qt: QCollatorSortKey::swap(QCollatorSortKey &). */
void qcollatorsortkey_swap(go_qcollatorsortkey_t *self, go_qcollatorsortkey_t *other);

/* Qt: QCollatorSortKey::compare(const QCollatorSortKey &) const. */
int qcollatorsortkey_compare(const go_qcollatorsortkey_t *self,
    const go_qcollatorsortkey_t *other);

/* Qt: operator<(const QCollatorSortKey &, const QCollatorSortKey &). */
bool qcollatorsortkey_less(const go_qcollatorsortkey_t *self,
    const go_qcollatorsortkey_t *other);

#ifdef __cplusplus
}
#endif

#endif
