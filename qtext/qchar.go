package qtext

/*
  QChar — обёртка над одноимённым классом Qt 6.11 (UTF-16 code unit).

  Семантика

  QChar — value-тип (uint16): тривиально сравним, хешируем, без аллокаций
  и финализатора. Нулевое значение эквивалентно QChar(0) (QChar::Null).

  Что сознательно не вошло

  - QChar::compare_helper(QChar, const char*) — приватная семантика Qt.
  - QChar(uchar row, uchar cell) — покрывается NewQChar(uint16(row)<<8 | uint16(cell)).
  - QChar::SpecialCharacter::LastValidCodePoint — доступен как rune LastValidCodePoint.
  - std::hash<QChar>, QT_CHAR_FASTCALL, Q_DECL_CONST_FUNCTION, QT6_ONLY и прочие impl-макросы Qt.
  - QChar(wchar_t) (Windows) — в Go нет wchar_t.

  Особенности дизайна

  - NewQCharFromRune(r) — усечение до uint16; для non-BMP — NewQCharFromRuneStrict
    или NewQCharFromUcs4.
  - NewQCharFromUcs4(r) — принимает только валидные code point-ы (0..0x10FFFF),
    иначе возвращает nil.
  - Decompose / QCharDecompose возвращают *QString (владение у вызывающего);
    DecomposeString / QCharDecomposeString возвращают Go string и закрывают
    промежуточный QString сами.
  - Все статические функции (CategoryOf, ToLowerOf, …) принимают rune и не
    создают Go-объектов.
*/

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC
#cgo LDFLAGS: -lstdc++
#include <stdint.h>
#include "qchar_c.h"
*/
import "C"

import "unsafe"


// SpecialCharacter constants mirror QChar::SpecialCharacter.
// Note: QChar::LastValidCodePoint (0x10FFFF) does NOT fit in a UTF-16
// code unit, so it is exposed separately as a rune.
const (
	Null                       QChar = 0x0000
	Tabulation                 QChar = 0x0009
	LineFeed                   QChar = 0x000A
	FormFeed                   QChar = 0x000C
	CarriageReturn             QChar = 0x000D
	Space                      QChar = 0x0020
	Nbsp                       QChar = 0x00A0
	SoftHyphen                 QChar = 0x00AD
	ReplacementCharacter       QChar = 0xFFFD
	ObjectReplacementCharacter QChar = 0xFFFC
	ByteOrderMark              QChar = 0xFEFF
	ByteOrderSwapped           QChar = 0xFFFE
	ParagraphSeparator         QChar = 0x2029
	LineSeparator              QChar = 0x2028
	VisualTabCharacter         QChar = 0x2192
)

// LastValidCodePoint is the highest legal Unicode code point (U+10FFFF).
// It cannot be represented as a single QChar and is therefore a rune.
const LastValidCodePoint rune = 0x10FFFF

// Surrogate range constants (UTF-16).
const (
	HighSurrogateMin uint16 = 0xD800
	HighSurrogateMax uint16 = 0xDBFF
	LowSurrogateMin  uint16 = 0xDC00
	LowSurrogateMax  uint16 = 0xDFFF
)


// QChar represents a single UTF-16 code unit; mirrors Qt6's QChar.
// It is a value type (Q_PRIMITIVE_TYPE): trivially comparable, hashable,
// allocation-free. The zero value equals QChar(0) (QChar::Null).
type QChar uint16

// ---- Constructors --------------------------------------------------------

// NewQChar returns a QChar for the given UTF-16 code unit.
// Mirrors QChar(char16_t) (direct cast: QChar is uint16 in Go).
func NewQChar(u uint16) QChar { return QChar(u) }

// NewQCharFromLatin1 returns a QChar for the given Latin-1 byte.
// Mirrors QChar::fromLatin1(char).
func NewQCharFromLatin1(c byte) QChar {
	return QChar(uint16(C.qchar_from_latin1(C.char(c))))
}

// NewQCharFromUcs2 returns a QChar for the given UTF-16 code unit.
// Mirrors QChar::fromUcs2(char16_t).
func NewQCharFromUcs2(ucs2 uint16) QChar {
	return QChar(uint16(C.qchar_from_ucs2(C.uint16_t(ucs2))))
}

// NewQCharFromUcs4 returns a slice of QChar (length 1 or 2) for a UCS-4
// code point. Mirrors QChar::fromUcs4(char32_t).
//
// Length 1 for BMP code points, length 2 (surrogate pair) for supplementary
// planes. Returns nil for invalid code points (< 0 or > LastValidCodePoint);
// Qt has no such check and exhibits UB for invalid input.
func NewQCharFromUcs4(ucs4 rune) []QChar {
	if ucs4 < 0 || ucs4 > LastValidCodePoint {
		return nil
	}
	var out [2]C.uint16_t
	n := C.qchar_from_ucs4(C.uint32_t(ucs4), &out[0])
	if n == 1 {
		return []QChar{QChar(out[0])}
	}
	return []QChar{QChar(out[0]), QChar(out[1])}
}

// NewQCharFromRune returns a QChar by truncating r to 16 bits: QChar(uint16(r)).
//
// For BMP code points (U+0000..U+FFFF) this is exact. For non-BMP this is a
// silent truncation; use NewQCharFromUcs4 or NewQCharFromRuneStrict.
//
// Mirrors QChar(char32_t) (explicit; Qt asserts range in debug builds).
func NewQCharFromRune(r rune) QChar { return QChar(uint16(r)) }

// NewQCharFromRuneStrict returns a QChar for r and true if r fits in a single
// UTF-16 code unit (U+0000..U+FFFF), or (0, false) otherwise.
//
// Go-only helper; equivalent to QChar(char32_t) with an explicit range check.
func NewQCharFromRuneStrict(r rune) (QChar, bool) {
	if r < 0 || r > 0xFFFF {
		return 0, false
	}
	return QChar(uint16(r)), true
}

// ---- Accessors -----------------------------------------------------------

// Unicode returns the underlying UTF-16 code unit.
// Mirrors QChar::unicode() const.
func (c QChar) Unicode() uint16 {
	return uint16(C.qchar_unicode(C.uint16_t(c)))
}

// SetUnicode replaces the underlying UTF-16 code unit.
// Mirrors QChar::unicode() (non-const overload).
func (c *QChar) SetUnicode(v uint16) { *c = QChar(v) }

// ToLatin1 returns the Latin-1 byte, or 0 if not representable.
// Mirrors QChar::toLatin1() const.
func (c QChar) ToLatin1() byte {
	return byte(C.qchar_to_latin1(C.uint16_t(c)))
}

// Rune returns the code unit as a Go rune (may be a surrogate half).
// Go-only helper.
func (c QChar) Rune() rune { return rune(c) }

// Cell returns the low byte of the UTF-16 code unit.
// Mirrors QChar::cell() const.
func (c QChar) Cell() byte { return byte(c & 0xFF) }

// Row returns the high byte of the UTF-16 code unit.
// Mirrors QChar::row() const.
func (c QChar) Row() byte { return byte((c >> 8) & 0xFF) }

// SetCell sets the low byte.
// Mirrors QChar::setCell(uchar).
func (c *QChar) SetCell(b byte) { *c = QChar(uint16(*c)&0xFF00 | uint16(b)) }

// SetRow sets the high byte.
// Mirrors QChar::setRow(uchar).
func (c *QChar) SetRow(b byte) { *c = QChar(uint16(b)<<8 | (uint16(*c) & 0x00FF)) }

// ---- Instance predicates -------------------------------------------------

// IsNull — QChar::isNull() const.
func (c QChar) IsNull() bool { return bool(C.qchar_is_null(C.uint16_t(c))) }

// IsPrint — QChar::isPrint() const.
func (c QChar) IsPrint() bool { return bool(C.qchar_is_print(C.uint16_t(c))) }

// IsSpace — QChar::isSpace() const.
func (c QChar) IsSpace() bool { return bool(C.qchar_is_space(C.uint16_t(c))) }

// IsMark — QChar::isMark() const.
func (c QChar) IsMark() bool { return bool(C.qchar_is_mark(C.uint16_t(c))) }

// IsLetter — QChar::isLetter() const.
func (c QChar) IsLetter() bool { return bool(C.qchar_is_letter(C.uint16_t(c))) }

// IsNumber — QChar::isNumber() const.
func (c QChar) IsNumber() bool { return bool(C.qchar_is_number(C.uint16_t(c))) }

// IsLetterOrNumber — QChar::isLetterOrNumber() const.
func (c QChar) IsLetterOrNumber() bool {
	return bool(C.qchar_is_letter_or_number(C.uint16_t(c)))
}

// IsDigit — QChar::isDigit() const.
func (c QChar) IsDigit() bool { return bool(C.qchar_is_digit(C.uint16_t(c))) }

// IsLower — QChar::isLower() const.
func (c QChar) IsLower() bool { return bool(C.qchar_is_lower(C.uint16_t(c))) }

// IsUpper — QChar::isUpper() const.
func (c QChar) IsUpper() bool { return bool(C.qchar_is_upper(C.uint16_t(c))) }

// IsTitleCase — QChar::isTitleCase() const.
func (c QChar) IsTitleCase() bool { return bool(C.qchar_is_title_case(C.uint16_t(c))) }

// IsPunct — QChar::isPunct() const.
func (c QChar) IsPunct() bool { return bool(C.qchar_is_punct(C.uint16_t(c))) }

// IsSymbol — QChar::isSymbol() const.
func (c QChar) IsSymbol() bool { return bool(C.qchar_is_symbol(C.uint16_t(c))) }

// IsNonCharacter — QChar::isNonCharacter() const.
func (c QChar) IsNonCharacter() bool {
	return bool(C.qchar_is_non_character(C.uint16_t(c)))
}

// IsHighSurrogate — QChar::isHighSurrogate() const.
func (c QChar) IsHighSurrogate() bool {
	return bool(C.qchar_is_high_surrogate(C.uint16_t(c)))
}

// IsLowSurrogate — QChar::isLowSurrogate() const.
func (c QChar) IsLowSurrogate() bool {
	return bool(C.qchar_is_low_surrogate(C.uint16_t(c)))
}

// IsSurrogate — QChar::isSurrogate() const.
func (c QChar) IsSurrogate() bool { return bool(C.qchar_is_surrogate(C.uint16_t(c))) }

// HasMirrored — QChar::hasMirrored() const.
func (c QChar) HasMirrored() bool { return bool(C.qchar_has_mirrored(C.uint16_t(c))) }

// ---- Instance classification --------------------------------------------

// Category — QChar::category() const.
func (c QChar) Category() Category {
	return Category(C.qchar_category(C.uint16_t(c)))
}

// Direction — QChar::direction() const.
func (c QChar) Direction() Direction {
	return Direction(C.qchar_direction(C.uint16_t(c)))
}

// JoiningType — QChar::joiningType() const.
func (c QChar) JoiningType() JoiningType {
	return JoiningType(C.qchar_joining_type(C.uint16_t(c)))
}

// CombiningClass — QChar::combiningClass() const.
func (c QChar) CombiningClass() CombiningClass {
	return CombiningClass(C.qchar_combining_class(C.uint16_t(c)))
}

// MirroredChar — QChar::mirroredChar() const.
func (c QChar) MirroredChar() QChar {
	return QChar(C.qchar_mirrored_char(C.uint16_t(c)))
}

// DecompositionTag — QChar::decompositionTag() const.
func (c QChar) DecompositionTag() Decomposition {
	return Decomposition(C.qchar_decomposition_tag(C.uint16_t(c)))
}

// DigitValue — QChar::digitValue() const.
func (c QChar) DigitValue() int {
	return int(C.qchar_digit_value(C.uint16_t(c)))
}

// Script — QChar::script() const.
func (c QChar) Script() Script {
	return Script(C.qchar_script(C.uint16_t(c)))
}

// UnicodeVersion — QChar::unicodeVersion() const.
func (c QChar) UnicodeVersion() UnicodeVersion {
	return UnicodeVersion(C.qchar_unicode_version(C.uint16_t(c)))
}

// ---- Instance case conversion -------------------------------------------

// ToLower — QChar::toLower() const.
func (c QChar) ToLower() QChar { return QChar(C.qchar_to_lower(C.uint16_t(c))) }

// ToUpper — QChar::toUpper() const.
func (c QChar) ToUpper() QChar { return QChar(C.qchar_to_upper(C.uint16_t(c))) }

// ToTitleCase — QChar::toTitleCase() const.
func (c QChar) ToTitleCase() QChar { return QChar(C.qchar_to_title(C.uint16_t(c))) }

// ToCaseFolded — QChar::toCaseFolded() const.
func (c QChar) ToCaseFolded() QChar { return QChar(C.qchar_to_case_folded(C.uint16_t(c))) }

// ---- Static classification (by UCS-4) -----------------------------------

// CategoryOf — QChar::category(char32_t).
func CategoryOf(ucs4 rune) Category {
	return Category(C.qchar_static_category(C.uint32_t(ucs4)))
}

// DirectionOf — QChar::direction(char32_t).
func DirectionOf(ucs4 rune) Direction {
	return Direction(C.qchar_static_direction(C.uint32_t(ucs4)))
}

// JoiningTypeOf — QChar::joiningType(char32_t).
func JoiningTypeOf(ucs4 rune) JoiningType {
	return JoiningType(C.qchar_static_joining_type(C.uint32_t(ucs4)))
}

// CombiningClassOf — QChar::combiningClass(char32_t).
func CombiningClassOf(ucs4 rune) CombiningClass {
	return CombiningClass(C.qchar_static_combining_class(C.uint32_t(ucs4)))
}

// MirroredCharOf — QChar::mirroredChar(char32_t).
func MirroredCharOf(ucs4 rune) rune {
	return rune(C.qchar_static_mirrored_char(C.uint32_t(ucs4)))
}

// HasMirroredOf — QChar::hasMirrored(char32_t).
func HasMirroredOf(ucs4 rune) bool {
	return bool(C.qchar_static_has_mirrored(C.uint32_t(ucs4)))
}

// DecompositionOf returns the decomposition tag of ucs4.
// Mirrors QChar::decompositionTag(char32_t). Does NOT return the
// decomposition string; use QCharDecompose for that.
func DecompositionOf(ucs4 rune) Decomposition {
	return Decomposition(C.qchar_static_decomposition_tag(C.uint32_t(ucs4)))
}

// DigitValueOf — QChar::digitValue(char32_t).
func DigitValueOf(ucs4 rune) int {
	return int(C.qchar_static_digit_value(C.uint32_t(ucs4)))
}

// ScriptOf — QChar::script(char32_t).
func ScriptOf(ucs4 rune) Script {
	return Script(C.qchar_static_script(C.uint32_t(ucs4)))
}

// UnicodeVersionOf — QChar::unicodeVersion(char32_t).
func UnicodeVersionOf(ucs4 rune) UnicodeVersion {
	return UnicodeVersion(C.qchar_static_unicode_version(C.uint32_t(ucs4)))
}

// CurrentUnicodeVersion — QChar::currentUnicodeVersion().
func CurrentUnicodeVersion() UnicodeVersion {
	return UnicodeVersion(C.qchar_static_current_unicode_version())
}

// ---- Static case conversion (by UCS-4; returns rune) ---------------------

// ToLowerOf — QChar::toLower(char32_t).
func ToLowerOf(ucs4 rune) rune { return rune(C.qchar_static_to_lower(C.uint32_t(ucs4))) }

// ToUpperOf — QChar::toUpper(char32_t).
func ToUpperOf(ucs4 rune) rune { return rune(C.qchar_static_to_upper(C.uint32_t(ucs4))) }

// ToTitleCaseOf — QChar::toTitleCase(char32_t).
func ToTitleCaseOf(ucs4 rune) rune { return rune(C.qchar_static_to_title(C.uint32_t(ucs4))) }

// ToCaseFoldedOf — QChar::toCaseFolded(char32_t).
func ToCaseFoldedOf(ucs4 rune) rune {
	return rune(C.qchar_static_to_case_folded(C.uint32_t(ucs4)))
}

// ---- Static predicates (by UCS-4) ---------------------------------------

// IsNull reports whether ucs4 is U+0000.
// Go-only convenience; Qt has no static QChar::isNull(char32_t).
func IsNull(ucs4 rune) bool { return ucs4 == 0 }

// IsPrint — QChar::isPrint(char32_t).
func IsPrint(ucs4 rune) bool { return bool(C.qchar_static_is_print(C.uint32_t(ucs4))) }

// IsSpace — QChar::isSpace(char32_t).
func IsSpace(ucs4 rune) bool { return bool(C.qchar_static_is_space(C.uint32_t(ucs4))) }

// IsMark — QChar::isMark(char32_t).
func IsMark(ucs4 rune) bool { return bool(C.qchar_static_is_mark(C.uint32_t(ucs4))) }

// IsLetter — QChar::isLetter(char32_t).
func IsLetter(ucs4 rune) bool { return bool(C.qchar_static_is_letter(C.uint32_t(ucs4))) }

// IsNumber — QChar::isNumber(char32_t).
func IsNumber(ucs4 rune) bool { return bool(C.qchar_static_is_number(C.uint32_t(ucs4))) }

// IsLetterOrNumber — QChar::isLetterOrNumber(char32_t).
func IsLetterOrNumber(ucs4 rune) bool {
	return bool(C.qchar_static_is_letter_or_number(C.uint32_t(ucs4)))
}

// IsDigit — QChar::isDigit(char32_t).
func IsDigit(ucs4 rune) bool { return bool(C.qchar_static_is_digit(C.uint32_t(ucs4))) }

// IsLower — QChar::isLower(char32_t).
func IsLower(ucs4 rune) bool { return bool(C.qchar_static_is_lower(C.uint32_t(ucs4))) }

// IsUpper — QChar::isUpper(char32_t).
func IsUpper(ucs4 rune) bool { return bool(C.qchar_static_is_upper(C.uint32_t(ucs4))) }

// IsTitleCase — QChar::isTitleCase(char32_t).
func IsTitleCase(ucs4 rune) bool { return bool(C.qchar_static_is_title_case(C.uint32_t(ucs4))) }

// IsPunct — QChar::isPunct(char32_t).
func IsPunct(ucs4 rune) bool { return bool(C.qchar_static_is_punct(C.uint32_t(ucs4))) }

// IsSymbol — QChar::isSymbol(char32_t).
func IsSymbol(ucs4 rune) bool { return bool(C.qchar_static_is_symbol(C.uint32_t(ucs4))) }

// IsNonCharacter — QChar::isNonCharacter(char32_t).
func IsNonCharacter(ucs4 rune) bool {
	return bool(C.qchar_static_is_non_character(C.uint32_t(ucs4)))
}

// IsHighSurrogateRune — QChar::isHighSurrogate(char32_t).
// Named ...Rune to avoid confusion with the instance method IsHighSurrogate.
func IsHighSurrogateRune(ucs4 rune) bool {
	return bool(C.qchar_static_is_high_surrogate(C.uint32_t(ucs4)))
}

// IsLowSurrogateRune — QChar::isLowSurrogate(char32_t).
// Named ...Rune to avoid confusion with the instance method IsLowSurrogate.
func IsLowSurrogateRune(ucs4 rune) bool {
	return bool(C.qchar_static_is_low_surrogate(C.uint32_t(ucs4)))
}

// IsSurrogateRune — QChar::isSurrogate(char32_t).
// Named ...Rune to avoid confusion with the instance method IsSurrogate.
func IsSurrogateRune(ucs4 rune) bool {
	return bool(C.qchar_static_is_surrogate(C.uint32_t(ucs4)))
}

// ---- Surrogate pairs -----------------------------------------------------

// RequiresSurrogates — QChar::requiresSurrogates(char32_t).
func RequiresSurrogates(ucs4 rune) bool {
	return bool(C.qchar_requires_surrogates(C.uint32_t(ucs4)))
}

// HighSurrogate — QChar::highSurrogate(char32_t).
func HighSurrogate(ucs4 rune) QChar {
	return QChar(uint16(C.qchar_high_surrogate(C.uint32_t(ucs4))))
}

// LowSurrogate — QChar::lowSurrogate(char32_t).
func LowSurrogate(ucs4 rune) QChar {
	return QChar(uint16(C.qchar_low_surrogate(C.uint32_t(ucs4))))
}

// SurrogateToUcs4 — QChar::surrogateToUcs4(char16_t, char16_t).
func SurrogateToUcs4(high, low QChar) rune {
	return rune(C.qchar_surrogate_to_ucs4(C.uint16_t(high), C.uint16_t(low)))
}

// ---- Decomposition -------------------------------------------------------

// Decompose returns the Unicode decomposition of c as a fresh *QString.
// Mirrors QChar::decomposition() const.
//
// Ownership: the caller owns the returned *QString — call Close() when
// done, or rely on the finalizer. Returns nil only if Qt fails to allocate
// (C++ exception caught at the C boundary).
func (c QChar) Decompose() *QString {
	p := C.qchar_decomposition(C.uint32_t(c))
	if p == nil {
		return nil
	}
	return wrapQString(unsafe.Pointer(p))
}

// DecomposeString is a Go-string variant of Decompose.
// Returns "" (empty) for code points without a decomposition.
// The intermediate *QString is closed automatically.
func (c QChar) DecomposeString() string {
	q := c.Decompose()
	if q == nil {
		return ""
	}
	defer q.Close()
	return q.String()
}

// QCharDecompose returns the Unicode decomposition of ucs4 as a fresh *QString.
// Mirrors static QChar::decomposition(char32_t).
//
// Returns nil for invalid code points (< 0 or > LastValidCodePoint) and on
// allocation failure. Ownership: the caller owns the *QString.
func QCharDecompose(ucs4 rune) *QString {
	if ucs4 < 0 || ucs4 > LastValidCodePoint {
		return nil
	}
	p := C.qchar_decomposition(C.uint32_t(ucs4))
	if p == nil {
		return nil
	}
	return wrapQString(unsafe.Pointer(p))
}

// QCharDecomposeString is a Go-string variant of QCharDecompose.
// Returns "" (empty) for invalid code points or code points without
// a decomposition. The intermediate *QString is closed automatically.
func QCharDecomposeString(ucs4 rune) string {
	q := QCharDecompose(ucs4)
	if q == nil {
		return ""
	}
	defer q.Close()
	return q.String()
}
