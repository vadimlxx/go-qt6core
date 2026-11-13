/* C ABI for Qt 6.11 QtCore/QRegularExpression family. */
#ifndef GO_QT6CORE_QREGEXP_C_H
#define GO_QT6CORE_QREGEXP_C_H

#include <stddef.h>
#include <stdint.h>
#include <stdbool.h>
#include "qstring_c.h"

#ifdef __cplusplus
class QRegularExpression;
class QRegularExpressionMatch;
class QRegularExpressionMatchIterator;
extern "C" {
#else
typedef struct QRegularExpression QRegularExpression;
typedef struct QRegularExpressionMatch QRegularExpressionMatch;
typedef struct QRegularExpressionMatchIterator QRegularExpressionMatchIterator;
#endif

/* ===== QRegularExpression ===== */
QRegularExpression* qregexp_new(void);
QRegularExpression* qregexp_new_from_pattern(const QString* pattern, int options);
QRegularExpression* qregexp_new_copy(const QRegularExpression* other);
void qregexp_delete(QRegularExpression* self);
void qregexp_swap(QRegularExpression* self, QRegularExpression* other);

QString* qregexp_pattern(const QRegularExpression* self);
void qregexp_set_pattern(QRegularExpression* self, const QString* pattern);

bool qregexp_is_valid(const QRegularExpression* self);
ptrdiff_t qregexp_pattern_error_offset(const QRegularExpression* self);
QString* qregexp_error_string(const QRegularExpression* self);

int qregexp_capture_count(const QRegularExpression* self);

int qregexp_pattern_options(const QRegularExpression* self);
void qregexp_set_pattern_options(QRegularExpression* self, int options);

QRegularExpressionMatch* qregexp_match_string(const QRegularExpression* self, const QString* subject,
    ptrdiff_t offset, int matchType, int matchOptions);
QRegularExpressionMatch* qregexp_match_view(const QRegularExpression* self, const uint16_t* subjectData,
    ptrdiff_t subjectSize, ptrdiff_t offset, int matchType, int matchOptions);

QRegularExpressionMatchIterator* qregexp_global_match_string(const QRegularExpression* self,
    const QString* subject, ptrdiff_t offset, int matchType, int matchOptions);
QRegularExpressionMatchIterator* qregexp_global_match_view(const QRegularExpression* self,
    const uint16_t* subjectData, ptrdiff_t subjectSize, ptrdiff_t offset, int matchType, int matchOptions);

QString* qregexp_escape(const QString* str);
QString* qregexp_wildcard_to_regular_expression(const QString* str, int options);
QString* qregexp_anchored_pattern(const QString* expression);
QRegularExpression* qregexp_from_wildcard(const QString* pattern, int cs, int options);

/* ===== QRegularExpressionMatch ===== */
QRegularExpressionMatch* qregexpmatch_new_empty(void);
QRegularExpressionMatch* qregexpmatch_new_copy(const QRegularExpressionMatch* other);
void qregexpmatch_delete(QRegularExpressionMatch* self);
void qregexpmatch_swap(QRegularExpressionMatch* self, QRegularExpressionMatch* other);

QRegularExpression* qregexpmatch_regular_expression(const QRegularExpressionMatch* self);
int qregexpmatch_match_type(const QRegularExpressionMatch* self);
int qregexpmatch_match_options(const QRegularExpressionMatch* self);

bool qregexpmatch_has_match(const QRegularExpressionMatch* self);
bool qregexpmatch_has_partial_match(const QRegularExpressionMatch* self);
bool qregexpmatch_is_valid(const QRegularExpressionMatch* self);
int qregexpmatch_last_captured_index(const QRegularExpressionMatch* self);

bool qregexpmatch_has_captured_int(const QRegularExpressionMatch* self, int nth);
bool qregexpmatch_has_captured_view(const QRegularExpressionMatch* self,
    const uint16_t* nameData, ptrdiff_t nameSize);

QString* qregexpmatch_captured_int(const QRegularExpressionMatch* self, int nth);
QString* qregexpmatch_captured_view(const QRegularExpressionMatch* self,
    const uint16_t* nameData, ptrdiff_t nameSize);

void qregexpmatch_captured_view_int(const QRegularExpressionMatch* self, int nth,
    const uint16_t** outData, ptrdiff_t* outSize);
void qregexpmatch_captured_view_view(const QRegularExpressionMatch* self,
    const uint16_t* nameData, ptrdiff_t nameSize,
    const uint16_t** outData, ptrdiff_t* outSize);

ptrdiff_t qregexpmatch_captured_start_int(const QRegularExpressionMatch* self, int nth);
ptrdiff_t qregexpmatch_captured_length_int(const QRegularExpressionMatch* self, int nth);
ptrdiff_t qregexpmatch_captured_end_int(const QRegularExpressionMatch* self, int nth);

ptrdiff_t qregexpmatch_captured_start_view(const QRegularExpressionMatch* self,
    const uint16_t* nameData, ptrdiff_t nameSize);
ptrdiff_t qregexpmatch_captured_length_view(const QRegularExpressionMatch* self,
    const uint16_t* nameData, ptrdiff_t nameSize);
ptrdiff_t qregexpmatch_captured_end_view(const QRegularExpressionMatch* self,
    const uint16_t* nameData, ptrdiff_t nameSize);

/* ===== QRegularExpressionMatchIterator ===== */
QRegularExpressionMatchIterator* qregexpmatchiter_new_empty(void);
QRegularExpressionMatchIterator* qregexpmatchiter_new_copy(const QRegularExpressionMatchIterator* other);
void qregexpmatchiter_delete(QRegularExpressionMatchIterator* self);
void qregexpmatchiter_swap(QRegularExpressionMatchIterator* self, QRegularExpressionMatchIterator* other);

bool qregexpmatchiter_is_valid(const QRegularExpressionMatchIterator* self);
bool qregexpmatchiter_has_next(const QRegularExpressionMatchIterator* self);
QRegularExpressionMatch* qregexpmatchiter_next(QRegularExpressionMatchIterator* self);
QRegularExpressionMatch* qregexpmatchiter_peek_next(const QRegularExpressionMatchIterator* self);

QRegularExpression* qregexpmatchiter_regular_expression(const QRegularExpressionMatchIterator* self);
int qregexpmatchiter_match_type(const QRegularExpressionMatchIterator* self);
int qregexpmatchiter_match_options(const QRegularExpressionMatchIterator* self);

#ifdef __cplusplus
}
#endif

#endif