#include "qregexp_c.h"
#include <QtCore/QRegularExpression>
#include <QtCore/QRegularExpressionMatch>
#include <QtCore/QRegularExpressionMatchIterator>
#include <QtCore/QString>
#include <QtCore/QStringView>

extern "C" {

/* ===== QRegularExpression ===== */

QRegularExpression* qregexp_new() {
    return new QRegularExpression();
}

QRegularExpression* qregexp_new_from_pattern(const QString* pattern, int options) {
    return new QRegularExpression(*pattern, QRegularExpression::PatternOptions(options));
}

QRegularExpression* qregexp_new_copy(const QRegularExpression* other) {
    return new QRegularExpression(*other);
}

void qregexp_delete(QRegularExpression* self) {
    delete self;
}

void qregexp_swap(QRegularExpression* self, QRegularExpression* other) {
    self->swap(*other);
}

QString* qregexp_pattern(const QRegularExpression* self) {
    return new QString(self->pattern());
}

void qregexp_set_pattern(QRegularExpression* self, const QString* pattern) {
    self->setPattern(*pattern);
}

bool qregexp_is_valid(const QRegularExpression* self) {
    return self->isValid();
}

ptrdiff_t qregexp_pattern_error_offset(const QRegularExpression* self) {
    return self->patternErrorOffset();
}

QString* qregexp_error_string(const QRegularExpression* self) {
    return new QString(self->errorString());
}

int qregexp_capture_count(const QRegularExpression* self) {
    return self->captureCount();
}

int qregexp_pattern_options(const QRegularExpression* self) {
    return self->patternOptions().toInt();
}

void qregexp_set_pattern_options(QRegularExpression* self, int options) {
    self->setPatternOptions(QRegularExpression::PatternOptions(options));
}

QRegularExpressionMatch* qregexp_match_string(const QRegularExpression* self, const QString* subject,
    ptrdiff_t offset, int matchType, int matchOptions) {
    auto m = self->match(*subject, offset,
                         QRegularExpression::MatchType(matchType),
                         QRegularExpression::MatchOptions(matchOptions));
    return new QRegularExpressionMatch(m);
}

QRegularExpressionMatch* qregexp_match_view(const QRegularExpression* self, const uint16_t* subjectData,
    ptrdiff_t subjectSize, ptrdiff_t offset, int matchType, int matchOptions) {
    QStringView sv(reinterpret_cast<const char16_t*>(subjectData), subjectSize);
    auto m = self->matchView(sv, offset,
                             QRegularExpression::MatchType(matchType),
                             QRegularExpression::MatchOptions(matchOptions));
    return new QRegularExpressionMatch(m);
}

QRegularExpressionMatchIterator* qregexp_global_match_string(const QRegularExpression* self,
    const QString* subject, ptrdiff_t offset, int matchType, int matchOptions) {
    auto it = self->globalMatch(*subject, offset,
                                QRegularExpression::MatchType(matchType),
                                QRegularExpression::MatchOptions(matchOptions));
    return new QRegularExpressionMatchIterator(it);
}

QRegularExpressionMatchIterator* qregexp_global_match_view(const QRegularExpression* self,
    const uint16_t* subjectData, ptrdiff_t subjectSize, ptrdiff_t offset, int matchType, int matchOptions) {
    QStringView sv(reinterpret_cast<const char16_t*>(subjectData), subjectSize);
    auto it = self->globalMatchView(sv, offset,
                                    QRegularExpression::MatchType(matchType),
                                    QRegularExpression::MatchOptions(matchOptions));
    return new QRegularExpressionMatchIterator(it);
}

QString* qregexp_escape(const QString* str) {
    return new QString(QRegularExpression::escape(*str));
}

QString* qregexp_wildcard_to_regular_expression(const QString* str, int options) {
    return new QString(QRegularExpression::wildcardToRegularExpression(
        *str, QRegularExpression::WildcardConversionOptions(options)));
}

QString* qregexp_anchored_pattern(const QString* expression) {
    return new QString(QRegularExpression::anchoredPattern(*expression));
}

QRegularExpression* qregexp_from_wildcard(const QString* pattern, int cs, int options) {
    return new QRegularExpression(QRegularExpression::fromWildcard(*pattern,
        cs == 0 ? Qt::CaseSensitive : Qt::CaseInsensitive,
        QRegularExpression::WildcardConversionOptions(options)));
}

/* ===== QRegularExpressionMatch ===== */

QRegularExpressionMatch* qregexpmatch_new_empty() {
    return new QRegularExpressionMatch();
}

QRegularExpressionMatch* qregexpmatch_new_copy(const QRegularExpressionMatch* other) {
    return new QRegularExpressionMatch(*other);
}

void qregexpmatch_delete(QRegularExpressionMatch* self) {
    delete self;
}

void qregexpmatch_swap(QRegularExpressionMatch* self, QRegularExpressionMatch* other) {
    self->swap(*other);
}

QRegularExpression* qregexpmatch_regular_expression(const QRegularExpressionMatch* self) {
    return new QRegularExpression(self->regularExpression());
}

int qregexpmatch_match_type(const QRegularExpressionMatch* self) {
    return self->matchType();
}

int qregexpmatch_match_options(const QRegularExpressionMatch* self) {
    return self->matchOptions().toInt();
}

bool qregexpmatch_has_match(const QRegularExpressionMatch* self) {
    return self->hasMatch();
}

bool qregexpmatch_has_partial_match(const QRegularExpressionMatch* self) {
    return self->hasPartialMatch();
}

bool qregexpmatch_is_valid(const QRegularExpressionMatch* self) {
    return self->isValid();
}

int qregexpmatch_last_captured_index(const QRegularExpressionMatch* self) {
    return self->lastCapturedIndex();
}

bool qregexpmatch_has_captured_int(const QRegularExpressionMatch* self, int nth) {
    return self->hasCaptured(nth);
}

bool qregexpmatch_has_captured_view(const QRegularExpressionMatch* self,
    const uint16_t* nameData, ptrdiff_t nameSize) {
    QStringView name(reinterpret_cast<const char16_t*>(nameData), nameSize);
    return self->hasCaptured(name);
}

QString* qregexpmatch_captured_int(const QRegularExpressionMatch* self, int nth) {
    return new QString(self->captured(nth));
}

QString* qregexpmatch_captured_view(const QRegularExpressionMatch* self,
    const uint16_t* nameData, ptrdiff_t nameSize) {
    QStringView name(reinterpret_cast<const char16_t*>(nameData), nameSize);
    return new QString(self->captured(name));
}

void qregexpmatch_captured_view_int(const QRegularExpressionMatch* self, int nth,
    const uint16_t** outData, ptrdiff_t* outSize) {
    auto sv = self->capturedView(nth);
    *outData = reinterpret_cast<const uint16_t*>(sv.utf16());
    *outSize = sv.size();
}

void qregexpmatch_captured_view_view(const QRegularExpressionMatch* self,
    const uint16_t* nameData, ptrdiff_t nameSize,
    const uint16_t** outData, ptrdiff_t* outSize) {
    QStringView name(reinterpret_cast<const char16_t*>(nameData), nameSize);
    auto sv = self->capturedView(name);
    *outData = reinterpret_cast<const uint16_t*>(sv.utf16());
    *outSize = sv.size();
}

ptrdiff_t qregexpmatch_captured_start_int(const QRegularExpressionMatch* self, int nth) {
    return self->capturedStart(nth);
}

ptrdiff_t qregexpmatch_captured_length_int(const QRegularExpressionMatch* self, int nth) {
    return self->capturedLength(nth);
}

ptrdiff_t qregexpmatch_captured_end_int(const QRegularExpressionMatch* self, int nth) {
    return self->capturedEnd(nth);
}

ptrdiff_t qregexpmatch_captured_start_view(const QRegularExpressionMatch* self,
    const uint16_t* nameData, ptrdiff_t nameSize) {
    QStringView name(reinterpret_cast<const char16_t*>(nameData), nameSize);
    return self->capturedStart(name);
}

ptrdiff_t qregexpmatch_captured_length_view(const QRegularExpressionMatch* self,
    const uint16_t* nameData, ptrdiff_t nameSize) {
    QStringView name(reinterpret_cast<const char16_t*>(nameData), nameSize);
    return self->capturedLength(name);
}

ptrdiff_t qregexpmatch_captured_end_view(const QRegularExpressionMatch* self,
    const uint16_t* nameData, ptrdiff_t nameSize) {
    QStringView name(reinterpret_cast<const char16_t*>(nameData), nameSize);
    return self->capturedEnd(name);
}

/* ===== QRegularExpressionMatchIterator ===== */

QRegularExpressionMatchIterator* qregexpmatchiter_new_empty() {
    return new QRegularExpressionMatchIterator();
}

QRegularExpressionMatchIterator* qregexpmatchiter_new_copy(const QRegularExpressionMatchIterator* other) {
    return new QRegularExpressionMatchIterator(*other);
}

void qregexpmatchiter_delete(QRegularExpressionMatchIterator* self) {
    delete self;
}

void qregexpmatchiter_swap(QRegularExpressionMatchIterator* self, QRegularExpressionMatchIterator* other) {
    self->swap(*other);
}

bool qregexpmatchiter_is_valid(const QRegularExpressionMatchIterator* self) {
    return self->isValid();
}

bool qregexpmatchiter_has_next(const QRegularExpressionMatchIterator* self) {
    return self->hasNext();
}

QRegularExpressionMatch* qregexpmatchiter_next(QRegularExpressionMatchIterator* self) {
    return new QRegularExpressionMatch(self->next());
}

QRegularExpressionMatch* qregexpmatchiter_peek_next(const QRegularExpressionMatchIterator* self) {
    return new QRegularExpressionMatch(self->peekNext());
}

QRegularExpression* qregexpmatchiter_regular_expression(const QRegularExpressionMatchIterator* self) {
    return new QRegularExpression(self->regularExpression());
}

int qregexpmatchiter_match_type(const QRegularExpressionMatchIterator* self) {
    return self->matchType();
}

int qregexpmatchiter_match_options(const QRegularExpressionMatchIterator* self) {
    return self->matchOptions().toInt();
}

} // extern "C"