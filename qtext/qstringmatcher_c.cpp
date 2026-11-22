/* C ABI for Qt 6.11 QtCore/QStringMatcher — implementation. */
#include "../internal/qt6core/check_qt_version.h"

#include <QString>
#include <QStringMatcher>
#include <QStringView>

#include "qstringmatcher_c.h"

namespace {

// Qt::CaseSensitivity raw values: CaseInsensitive = 0, CaseSensitive = 1.
// qdata.CaseSensitivity mirrors these values, so a direct static_cast is
// correct (no inversion like the one used in qstring_c.cpp for bool-based
// APIs).
inline Qt::CaseSensitivity CS(int cs) {
    return static_cast<Qt::CaseSensitivity>(cs);
}

} // namespace

extern "C" {

/* Constructors */

QStringMatcher* qstringmatcher_new() {
    return new QStringMatcher();
}

QStringMatcher* qstringmatcher_new_from_qstring(const QString* pattern, int cs) {
    if (!pattern) return new QStringMatcher();
    return new QStringMatcher(*pattern, CS(cs));
}

QStringMatcher* qstringmatcher_new_from_view(const uint16_t* uc, ptrdiff_t len, int cs) {
    if (!uc || len <= 0) return new QStringMatcher();
    QStringView v(reinterpret_cast<const char16_t*>(uc), static_cast<qsizetype>(len));
    return new QStringMatcher(v, CS(cs));
}

QStringMatcher* qstringmatcher_new_copy(const QStringMatcher* other) {
    if (!other) return new QStringMatcher();
    return new QStringMatcher(*other);
}

/* Destructor */

void qstringmatcher_delete(QStringMatcher* self) {
    delete self;
}

/* Setters */

void qstringmatcher_set_pattern(QStringMatcher* self, const QString* pattern) {
    if (!self || !pattern) return;
    self->setPattern(*pattern);
}

void qstringmatcher_set_case_sensitivity(QStringMatcher* self, int cs) {
    if (!self) return;
    self->setCaseSensitivity(CS(cs));
}

void qstringmatcher_assign(QStringMatcher* self, const QStringMatcher* other) {
    if (!self || !other) return;
    *self = *other;
}

/* indexIn */

ptrdiff_t qstringmatcher_index_in_qstring(const QStringMatcher* self,
                                          const QString* str, ptrdiff_t from) {
    if (!self || !str) return -1;
    return static_cast<ptrdiff_t>(self->indexIn(*str, static_cast<qsizetype>(from)));
}

ptrdiff_t qstringmatcher_index_in_view(const QStringMatcher* self,
                                       const uint16_t* uc, ptrdiff_t len,
                                       ptrdiff_t from) {
    if (!self) return -1;
    if (!uc || len <= 0) {
        return static_cast<ptrdiff_t>(
            self->indexIn(QStringView(), static_cast<qsizetype>(from)));
    }
    QStringView v(reinterpret_cast<const char16_t*>(uc), static_cast<qsizetype>(len));
    return static_cast<ptrdiff_t>(self->indexIn(v, static_cast<qsizetype>(from)));
}

/* Getters */

QString* qstringmatcher_pattern(const QStringMatcher* self) {
    if (!self) return new QString();
    return new QString(self->pattern());
}

const uint16_t* qstringmatcher_pattern_view(const QStringMatcher* self,
                                            ptrdiff_t* out_len) {
    if (out_len) *out_len = 0;
    if (!self) return nullptr;
    QStringView v = self->patternView();
    if (out_len) *out_len = static_cast<ptrdiff_t>(v.size());
    return reinterpret_cast<const uint16_t*>(v.utf16());
}

int qstringmatcher_case_sensitivity(const QStringMatcher* self) {
    if (!self) return 0;
    return static_cast<int>(self->caseSensitivity());
}

} /* extern "C" */