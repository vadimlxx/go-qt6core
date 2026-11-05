#include "../internal/qt6core/check_qt_version.h"

#include <QChar>
#include <QtCore/QLatin1Char>

#include "qchar_c.h"

namespace {
constexpr QChar C(uint16_t u) noexcept { return QChar::fromUcs2(static_cast<char16_t>(u)); }
constexpr char32_t U(uint32_t u) noexcept { return static_cast<char32_t>(u); }
} // namespace

extern "C" {

/* --- accessors --- */
uint16_t qchar_unicode(uint16_t self) {
    return static_cast<uint16_t>(C(self).unicode());
}
uint8_t qchar_to_latin1(uint16_t self) {
    return static_cast<uint8_t>(C(self).toLatin1());
}

/* --- predicates --- */
bool qchar_is_null(uint16_t s)             { return C(s).isNull(); }
bool qchar_is_print(uint16_t s)            { return C(s).isPrint(); }
bool qchar_is_space(uint16_t s)            { return C(s).isSpace(); }
bool qchar_is_mark(uint16_t s)             { return C(s).isMark(); }
bool qchar_is_letter(uint16_t s)           { return C(s).isLetter(); }
bool qchar_is_number(uint16_t s)           { return C(s).isNumber(); }
bool qchar_is_letter_or_number(uint16_t s) { return C(s).isLetterOrNumber(); }
bool qchar_is_digit(uint16_t s)            { return C(s).isDigit(); }
bool qchar_is_lower(uint16_t s)            { return C(s).isLower(); }
bool qchar_is_upper(uint16_t s)            { return C(s).isUpper(); }
bool qchar_is_title_case(uint16_t s)       { return C(s).isTitleCase(); }
bool qchar_is_punct(uint16_t s)            { return C(s).isPunct(); }
bool qchar_is_symbol(uint16_t s)           { return C(s).isSymbol(); }
bool qchar_is_non_character(uint16_t s)    { return C(s).isNonCharacter(); }
bool qchar_is_high_surrogate(uint16_t s)   { return C(s).isHighSurrogate(); }
bool qchar_is_low_surrogate(uint16_t s)    { return C(s).isLowSurrogate(); }
bool qchar_is_surrogate(uint16_t s)        { return C(s).isSurrogate(); }
bool qchar_has_mirrored(uint16_t s)        { return C(s).hasMirrored(); }

/* --- instance classification --- */
int qchar_category(uint16_t s)          { return static_cast<int>(C(s).category()); }
int qchar_direction(uint16_t s)         { return static_cast<int>(C(s).direction()); }
int qchar_joining_type(uint16_t s)      { return static_cast<int>(C(s).joiningType()); }
uint8_t qchar_combining_class(uint16_t s){ return static_cast<uint8_t>(C(s).combiningClass()); }
uint16_t qchar_mirrored_char(uint16_t s){ return static_cast<uint16_t>(C(s).mirroredChar().unicode()); }
int qchar_decomposition_tag(uint16_t s) { return static_cast<int>(C(s).decompositionTag()); }
int qchar_digit_value(uint16_t s)       { return C(s).digitValue(); }
int qchar_script(uint16_t s)            { return static_cast<int>(C(s).script()); }
int qchar_unicode_version(uint16_t s)   { return static_cast<int>(C(s).unicodeVersion()); }

/* --- instance case conversion (return QChar) --- */
uint16_t qchar_to_lower(uint16_t s)      { return static_cast<uint16_t>(C(s).toLower().unicode()); }
uint16_t qchar_to_upper(uint16_t s)      { return static_cast<uint16_t>(C(s).toUpper().unicode()); }
uint16_t qchar_to_title(uint16_t s)      { return static_cast<uint16_t>(C(s).toTitleCase().unicode()); }
uint16_t qchar_to_case_folded(uint16_t s){ return static_cast<uint16_t>(C(s).toCaseFolded().unicode()); }

/* --- static helpers --- */
int      qchar_static_category(uint32_t u)          { return static_cast<int>(QChar::category(U(u))); }
int      qchar_static_direction(uint32_t u)         { return static_cast<int>(QChar::direction(U(u))); }
int      qchar_static_joining_type(uint32_t u)      { return static_cast<int>(QChar::joiningType(U(u))); }
uint8_t  qchar_static_combining_class(uint32_t u)   { return static_cast<uint8_t>(QChar::combiningClass(U(u))); }
uint32_t qchar_static_mirrored_char(uint32_t u)     { return static_cast<uint32_t>(QChar::mirroredChar(U(u))); }
bool     qchar_static_has_mirrored(uint32_t u)      { return QChar::hasMirrored(U(u)); }
int      qchar_static_decomposition_tag(uint32_t u) { return static_cast<int>(QChar::decompositionTag(U(u))); }
int      qchar_static_digit_value(uint32_t u)       { return QChar::digitValue(U(u)); }
uint32_t qchar_static_to_lower(uint32_t u)          { return static_cast<uint32_t>(QChar::toLower(U(u))); }
uint32_t qchar_static_to_upper(uint32_t u)          { return static_cast<uint32_t>(QChar::toUpper(U(u))); }
uint32_t qchar_static_to_title(uint32_t u)          { return static_cast<uint32_t>(QChar::toTitleCase(U(u))); }
uint32_t qchar_static_to_case_folded(uint32_t u)    { return static_cast<uint32_t>(QChar::toCaseFolded(U(u))); }
int      qchar_static_script(uint32_t u)            { return static_cast<int>(QChar::script(U(u))); }
int      qchar_static_unicode_version(uint32_t u)   { return static_cast<int>(QChar::unicodeVersion(U(u))); }
int      qchar_static_current_unicode_version(void) { return static_cast<int>(QChar::currentUnicodeVersion()); }

bool qchar_static_is_print(uint32_t u)            { return QChar::isPrint(U(u)); }
bool qchar_static_is_space(uint32_t u)            { return QChar::isSpace(U(u)); }
bool qchar_static_is_mark(uint32_t u)             { return QChar::isMark(U(u)); }
bool qchar_static_is_letter(uint32_t u)           { return QChar::isLetter(U(u)); }
bool qchar_static_is_number(uint32_t u)           { return QChar::isNumber(U(u)); }
bool qchar_static_is_letter_or_number(uint32_t u) { return QChar::isLetterOrNumber(U(u)); }
bool qchar_static_is_digit(uint32_t u)            { return QChar::isDigit(U(u)); }
bool qchar_static_is_lower(uint32_t u)            { return QChar::isLower(U(u)); }
bool qchar_static_is_upper(uint32_t u)            { return QChar::isUpper(U(u)); }
bool qchar_static_is_title_case(uint32_t u)       { return QChar::isTitleCase(U(u)); }
bool qchar_static_is_punct(uint32_t u)            { return QChar::isPunct(U(u)); }
bool qchar_static_is_symbol(uint32_t u)           { return QChar::isSymbol(U(u)); }
bool qchar_static_is_non_character(uint32_t u)    { return QChar::isNonCharacter(U(u)); }
bool qchar_static_is_high_surrogate(uint32_t u)   { return QChar::isHighSurrogate(U(u)); }
bool qchar_static_is_low_surrogate(uint32_t u)    { return QChar::isLowSurrogate(U(u)); }
bool qchar_static_is_surrogate(uint32_t u)        { return QChar::isSurrogate(U(u)); }
bool qchar_requires_surrogates(uint32_t u)        { return QChar::requiresSurrogates(U(u)); }

/* --- surrogates --- */
uint16_t qchar_high_surrogate(uint32_t u) {
    return static_cast<uint16_t>(QChar::highSurrogate(U(u)));
}
uint16_t qchar_low_surrogate(uint32_t u) {
    return static_cast<uint16_t>(QChar::lowSurrogate(U(u)));
}
uint32_t qchar_surrogate_to_ucs4(uint16_t hi, uint16_t lo) {
    return static_cast<uint32_t>(
        QChar::surrogateToUcs4(static_cast<char16_t>(hi),
                               static_cast<char16_t>(lo)));
}

/* --- static constructors --- */
uint16_t qchar_from_latin1(char c) {
    return static_cast<uint16_t>(QLatin1Char(c).unicode());
}
uint16_t qchar_from_ucs2(uint16_t ucs2) {
    return static_cast<uint16_t>(QChar::fromUcs2(static_cast<char16_t>(ucs2)).unicode());
}
void qchar_from_ucs4(uint32_t ucs4, uint16_t out[2]) {
    if (ucs4 < 0x10000) {
        out[0] = static_cast<uint16_t>(ucs4);
        out[1] = 0;
    } else {
        out[0] = static_cast<uint16_t>(QChar::highSurrogate(U(ucs4)));
        out[1] = static_cast<uint16_t>(QChar::lowSurrogate(U(ucs4)));
    }
}

} // extern "C"
