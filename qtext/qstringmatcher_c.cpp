/* Order: version check first, then real Qt headers. */
#include "../internal/qt6core/check_qt_version.h"

/* qstringmatcher.h transitively includes qstring.h and qstringview.h. */
#include <QtCore/qstringmatcher.h>
#include <QtCore/qstringview.h>
#include <QtCore/qstring.h>

#include "qstringmatcher_c.h"

namespace {

inline QStringMatcher       *M(go_qsm_t       *p) { return reinterpret_cast<QStringMatcher       *>(p); }
inline const QStringMatcher *M(const go_qsm_t *p) { return reinterpret_cast<const QStringMatcher *>(p); }
inline go_qsm_t *toM(QStringMatcher *p)           { return reinterpret_cast<go_qsm_t *>(p); }

inline Qt::CaseSensitivity CS(int cs) noexcept {
    return cs == 0 ? Qt::CaseInsensitive : Qt::CaseSensitive;
}
inline int csToInt(Qt::CaseSensitivity cs) noexcept {
    return cs == Qt::CaseSensitive ? 1 : 0;
}

} // namespace

extern "C" {

/* ---- Constructors ---- */
go_qsm_t *qstringmatcher_new_default(void) {
    return toM(new QStringMatcher());
}
go_qsm_t *qstringmatcher_new_from_qstring(const void *qstr, int cs) {
    const QString &s = qstr ? *static_cast<const QString *>(qstr) : QString();
    return toM(new QStringMatcher(s, CS(cs)));
}
go_qsm_t *qstringmatcher_new_from_utf16(const uint16_t *data, ptrdiff_t size, int cs) {
    QStringView v(reinterpret_cast<const char16_t *>(data), static_cast<qsizetype>(size));
    return toM(new QStringMatcher(v, CS(cs)));
}
go_qsm_t *qstringmatcher_new_copy(const go_qsm_t *other) {
    return toM(new QStringMatcher(*M(other)));
}

/* ---- Lifecycle ---- */
void qstringmatcher_delete(go_qsm_t *self) {
    if (self) delete M(self);
}
void qstringmatcher_assign(go_qsm_t *self, const go_qsm_t *other) {
    if (self && other) *M(self) = *M(other);
}

/* ---- Pattern / sensitivity ---- */
void qstringmatcher_set_pattern(go_qsm_t *self, const void *qstr) {
    if (!self) return;
    if (qstr) {
        M(self)->setPattern(*static_cast<const QString *>(qstr));
    } else {
        M(self)->setPattern(QString());
    }
}
void qstringmatcher_set_case_sensitivity(go_qsm_t *self, int cs) {
    if (self) M(self)->setCaseSensitivity(CS(cs));
}
int qstringmatcher_case_sensitivity(const go_qsm_t *self) {
    return self ? csToInt(M(self)->caseSensitivity()) : 1;
}

/* ---- Search ---- */
ptrdiff_t qstringmatcher_index_in_utf16(const go_qsm_t *self,
    const uint16_t *data, ptrdiff_t size, ptrdiff_t from) {
    QStringView v(reinterpret_cast<const char16_t *>(data), static_cast<qsizetype>(size));
    return static_cast<ptrdiff_t>(M(self)->indexIn(v, static_cast<qsizetype>(from)));
}

/* ---- Pattern accessors ---- */
void qstringmatcher_pattern(const go_qsm_t *self, void **out_qstring) {
    *out_qstring = static_cast<void *>(new QString(M(self)->pattern()));
}
void qstringmatcher_pattern_view(const go_qsm_t *self,
    const uint16_t **out_data, ptrdiff_t *out_size) {
    QStringView v = M(self)->patternView();
    *out_data = reinterpret_cast<const uint16_t *>(v.utf16());
    *out_size = static_cast<ptrdiff_t>(v.size());
}

} // extern "C"