/* qtext/qchar_c.cpp
 *
 * Single C++ translation unit for qchar_c.h: defines all C symbols
 * declared there. Non-inline, extern "C" — guaranteed to be emitted
 * into this object file, so the Go linker finds them.
 *
 * Both QLatin1Char and QChar live in qchar_c.h, mirroring Qt 6.11
 * QtCore's qchar.h layout.
 */
#include "../internal/qt6core/check_qt_version.h"

#include <cstdint>
#include <QChar>      /* brings QLatin1Char and QChar */
#include <QString>

#include "qchar_c.h"

namespace {

inline QLatin1Char qt6core_ql1(uint8_t u) noexcept {
    return QLatin1Char(static_cast<char>(u));
}

constexpr QChar qt6core_qchar(uint16_t u) noexcept {
    return QChar::fromUcs2(static_cast<char16_t>(u));
}

constexpr char32_t qt6core_ucs4(uint32_t u) noexcept {
    return static_cast<char32_t>(u);
}

} // namespace

extern "C" {

/* ================================================================
 * QLatin1Char
 * ================================================================ */

uint8_t qlatin1char_to_latin1(uint8_t self) {
    return static_cast<uint8_t>(qt6core_ql1(self).toLatin1());
}
uint16_t qlatin1char_unicode(uint8_t self) {
    return static_cast<uint16_t>(qt6core_ql1(self).unicode());
}
bool qlatin1char_compares_equal(uint8_t a, uint8_t b) {
    return qt6core_ql1(a) == qt6core_ql1(b);
}
int qlatin1char_compare(uint8_t a, uint8_t b) {
    const QLatin1Char ca = qt6core_ql1(a);
    const QLatin1Char cb = qt6core_ql1(b);
    if (ca == cb) return 0;
    return ca < cb ? -1 : +1;
}
bool qlatin1char_compares_equal_byte(uint8_t self, uint8_t other) {
    return qt6core_ql1(self) == static_cast<char>(other);
}
int qlatin1char_compare_byte(uint8_t self, uint8_t other) {
    const QLatin1Char cs = qt6core_ql1(self);
    const char o = static_cast<char>(other);
    if (cs == o) return 0;
    return cs < o ? -1 : +1;
}

/* ================================================================
 * QChar — accessors
 * ================================================================ */

uint16_t qchar_unicode(uint16_t s) {
    return static_cast<uint16_t>(qt6core_qchar(s).unicode());
}
uint8_t qchar_to_latin1(uint16_t s) {
    return static_cast<uint8_t>(qt6core_qchar(s).toLatin1());
}

/* ================================================================
 * QChar — instance predicates
 * ================================================================ */

bool qchar_is_null(uint16_t s)             { return qt6core_qchar(s).isNull(); }
bool qchar_is_print(uint16_t s)            { return qt6core_qchar(s).isPrint(); }
bool qchar_is_space(uint16_t s)            { return qt6core_qchar(s).isSpace(); }
bool qchar_is_mark(uint16_t s)             { return qt6core_qchar(s).isMark(); }
bool qchar_is_letter(uint16_t s)           { return qt6core_qchar(s).isLetter(); }
bool qchar_is_number(uint16_t s)           { return qt6core_qchar(s).isNumber(); }
bool qchar_is_letter_or_number(uint16_t s) { return qt6core_qchar(s).isLetterOrNumber(); }
bool qchar_is_digit(uint16_t s)            { return qt6core_qchar(s).isDigit(); }
bool qchar_is_lower(uint16_t s)            { return qt6core_qchar(s).isLower(); }
bool qchar_is_upper(uint16_t s)            { return qt6core_qchar(s).isUpper(); }
bool qchar_is_title_case(uint16_t s)       { return qt6core_qchar(s).isTitleCase(); }
bool qchar_is_punct(uint16_t s)            { return qt6core_qchar(s).isPunct(); }
bool qchar_is_symbol(uint16_t s)           { return qt6core_qchar(s).isSymbol(); }
bool qchar_is_non_character(uint16_t s)    { return qt6core_qchar(s).isNonCharacter(); }
bool qchar_is_high_surrogate(uint16_t s)   { return qt6core_qchar(s).isHighSurrogate(); }
bool qchar_is_low_surrogate(uint16_t s)    { return qt6core_qchar(s).isLowSurrogate(); }
bool qchar_is_surrogate(uint16_t s)        { return qt6core_qchar(s).isSurrogate(); }
bool qchar_has_mirrored(uint16_t s)        { return qt6core_qchar(s).hasMirrored(); }

/* ================================================================
 * QChar — instance classification
 * ================================================================ */

int qchar_category(uint16_t s) {
    return static_cast<int>(qt6core_qchar(s).category());
}
int qchar_direction(uint16_t s) {
    return static_cast<int>(qt6core_qchar(s).direction());
}
int qchar_joining_type(uint16_t s) {
    return static_cast<int>(qt6core_qchar(s).joiningType());
}
uint8_t qchar_combining_class(uint16_t s) {
    return static_cast<uint8_t>(qt6core_qchar(s).combiningClass());
}
uint16_t qchar_mirrored_char(uint16_t s) {
    return static_cast<uint16_t>(qt6core_qchar(s).mirroredChar().unicode());
}
int qchar_decomposition_tag(uint16_t s) {
    return static_cast<int>(qt6core_qchar(s).decompositionTag());
}
int qchar_digit_value(uint16_t s) {
    return qt6core_qchar(s).digitValue();
}
int qchar_script(uint16_t s) {
    return static_cast<int>(qt6core_qchar(s).script());
}
int qchar_unicode_version(uint16_t s) {
    return static_cast<int>(qt6core_qchar(s).unicodeVersion());
}

/* ================================================================
 * QChar — instance case conversion
 * ================================================================ */

uint16_t qchar_to_lower(uint16_t s) {
    return static_cast<uint16_t>(qt6core_qchar(s).toLower().unicode());
}
uint16_t qchar_to_upper(uint16_t s) {
    return static_cast<uint16_t>(qt6core_qchar(s).toUpper().unicode());
}
uint16_t qchar_to_title(uint16_t s) {
    return static_cast<uint16_t>(qt6core_qchar(s).toTitleCase().unicode());
}
uint16_t qchar_to_case_folded(uint16_t s) {
    return static_cast<uint16_t>(qt6core_qchar(s).toCaseFolded().unicode());
}

/* ================================================================
 * QChar — static helpers (UCS-4)
 * ================================================================ */

int      qchar_static_category(uint32_t u)          { return static_cast<int>(QChar::category(qt6core_ucs4(u))); }
int      qchar_static_direction(uint32_t u)         { return static_cast<int>(QChar::direction(qt6core_ucs4(u))); }
int      qchar_static_joining_type(uint32_t u)      { return static_cast<int>(QChar::joiningType(qt6core_ucs4(u))); }
uint8_t  qchar_static_combining_class(uint32_t u)   { return static_cast<uint8_t>(QChar::combiningClass(qt6core_ucs4(u))); }
uint32_t qchar_static_mirrored_char(uint32_t u)     { return static_cast<uint32_t>(QChar::mirroredChar(qt6core_ucs4(u))); }
bool     qchar_static_has_mirrored(uint32_t u)      { return QChar::hasMirrored(qt6core_ucs4(u)); }
int      qchar_static_decomposition_tag(uint32_t u) { return static_cast<int>(QChar::decompositionTag(qt6core_ucs4(u))); }
int      qchar_static_digit_value(uint32_t u)       { return QChar::digitValue(qt6core_ucs4(u)); }
uint32_t qchar_static_to_lower(uint32_t u)          { return static_cast<uint32_t>(QChar::toLower(qt6core_ucs4(u))); }
uint32_t qchar_static_to_upper(uint32_t u)          { return static_cast<uint32_t>(QChar::toUpper(qt6core_ucs4(u))); }
uint32_t qchar_static_to_title(uint32_t u)          { return static_cast<uint32_t>(QChar::toTitleCase(qt6core_ucs4(u))); }
uint32_t qchar_static_to_case_folded(uint32_t u)    { return static_cast<uint32_t>(QChar::toCaseFolded(qt6core_ucs4(u))); }
int      qchar_static_script(uint32_t u)            { return static_cast<int>(QChar::script(qt6core_ucs4(u))); }
int      qchar_static_unicode_version(uint32_t u)   { return static_cast<int>(QChar::unicodeVersion(qt6core_ucs4(u))); }
int      qchar_static_current_unicode_version(void) { return static_cast<int>(QChar::currentUnicodeVersion()); }

bool qchar_static_is_print(uint32_t u)            { return QChar::isPrint(qt6core_ucs4(u)); }
bool qchar_static_is_space(uint32_t u)            { return QChar::isSpace(qt6core_ucs4(u)); }
bool qchar_static_is_mark(uint32_t u)             { return QChar::isMark(qt6core_ucs4(u)); }
bool qchar_static_is_letter(uint32_t u)           { return QChar::isLetter(qt6core_ucs4(u)); }
bool qchar_static_is_number(uint32_t u)           { return QChar::isNumber(qt6core_ucs4(u)); }
bool qchar_static_is_letter_or_number(uint32_t u) { return QChar::isLetterOrNumber(qt6core_ucs4(u)); }
bool qchar_static_is_digit(uint32_t u)            { return QChar::isDigit(qt6core_ucs4(u)); }
bool qchar_static_is_lower(uint32_t u)            { return QChar::isLower(qt6core_ucs4(u)); }
bool qchar_static_is_upper(uint32_t u)            { return QChar::isUpper(qt6core_ucs4(u)); }
bool qchar_static_is_title_case(uint32_t u)       { return QChar::isTitleCase(qt6core_ucs4(u)); }
bool qchar_static_is_punct(uint32_t u)            { return QChar::isPunct(qt6core_ucs4(u)); }
bool qchar_static_is_symbol(uint32_t u)           { return QChar::isSymbol(qt6core_ucs4(u)); }
bool qchar_static_is_non_character(uint32_t u)    { return QChar::isNonCharacter(qt6core_ucs4(u)); }
bool qchar_static_is_high_surrogate(uint32_t u)   { return QChar::isHighSurrogate(qt6core_ucs4(u)); }
bool qchar_static_is_low_surrogate(uint32_t u)    { return QChar::isLowSurrogate(qt6core_ucs4(u)); }
bool qchar_static_is_surrogate(uint32_t u)        { return QChar::isSurrogate(qt6core_ucs4(u)); }
bool qchar_requires_surrogates(uint32_t u)        { return QChar::requiresSurrogates(qt6core_ucs4(u)); }

/* ================================================================
 * QChar — surrogates
 * ================================================================ */

uint16_t qchar_high_surrogate(uint32_t u) {
    return static_cast<uint16_t>(QChar::highSurrogate(qt6core_ucs4(u)));
}
uint16_t qchar_low_surrogate(uint32_t u) {
    return static_cast<uint16_t>(QChar::lowSurrogate(qt6core_ucs4(u)));
}
uint32_t qchar_surrogate_to_ucs4(uint16_t hi, uint16_t lo) {
    return static_cast<uint32_t>(
        QChar::surrogateToUcs4(static_cast<char16_t>(hi), static_cast<char16_t>(lo)));
}

/* ================================================================
 * QChar — static constructors
 * ================================================================ */

uint16_t qchar_from_latin1(char c) {
    return static_cast<uint16_t>(QChar::fromLatin1(c).unicode());
}
uint16_t qchar_from_ucs2(uint16_t ucs2) {
    return static_cast<uint16_t>(QChar::fromUcs2(static_cast<char16_t>(ucs2)).unicode());
}
int qchar_from_ucs4(uint32_t ucs4, uint16_t out[2]) {
    if (ucs4 < 0x10000u) {
        out[0] = static_cast<uint16_t>(ucs4);
        out[1] = 0;
        return 1;
    }
    out[0] = static_cast<uint16_t>(QChar::highSurrogate(qt6core_ucs4(ucs4)));
    out[1] = static_cast<uint16_t>(QChar::lowSurrogate(qt6core_ucs4(ucs4)));
    return 2;
}

/* ================================================================
 * QChar — decomposition
 * ================================================================ */

QString* qchar_decomposition(uint32_t ucs4) {
    try {
        return new QString(QChar::decomposition(qt6core_ucs4(ucs4)));
    } catch (...) {
        return nullptr;
    }
}

} /* extern "C" */
