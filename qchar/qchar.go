package qchar

/*
 Что отсутствует:
 QChar::decomposition(char32_t) → QString	instance + static	Добавить через unsafe.Pointer (Go-цикл qchar ↔ qst*ring нельзя) + удобная обёртка в qstring.
 QChar::compare_helper(QChar, const char*) → int (comparesEqual/compareThreeWay vs const char*)	static	TODO — семантика Qt неочевидна (сравнение с первым UTF-16 code unit строки); требует отдельного решения.
 QChar(uchar row, uchar cell) (двухбайтовый конструктор)	ctor	Покрывается New(uint16(row)<<8 | uint16(cell)) — отдельный хелпер не нужен.
 QChar::fromUcs4 возвращает struct с operator QStringView()	ctor	Go возвращает []QChar. FromUcs4QStringView невозможен (цикл qchar ↔ qstringview).
 QChar::SpecialCharacter::LastValidCodePoint	const	Не влезает в uint16 → const LastValidCodePoint rune = 0x10FFFF.
 std::hash<QChar>, QT_CHAR_FASTCALL, Q_DECL_CONST_FUNCTION, Q_DECLARE_TYPEINFO, QT6_ONLY	impl	Qt-специфика, не для Go.
 QChar(wchar_t) (Windows)	ctor	Go не имеет wchar_t.
 FromRune для non-BMP rune	behavior	Сейчас — truncation (footgun). Добавить FromRuneStrict(rune) (QChar, bool).
*/

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC -Wno-sfinae-incomplete
#cgo LDFLAGS: -lstdc++
#include <stdint.h>
#include "qchar_c.h"
*/
import "C"

import "unsafe"

// QChar represents a single UTF-16 code unit; mirrors Qt6's QChar.
// It is a value type (Q_PRIMITIVE_TYPE): trivially comparable, hashable,
// allocation-free. The zero value equals QChar(0) (Null).
type QChar uint16

// ---- Construction --------------------------------------------------------

// New returns a QChar for the given UTF-16 code unit.
func New(u uint16) QChar { return QChar(u) }

// FromLatin1 constructs a QChar from a Latin-1 byte
// (mirrors QChar::fromLatin1).
func FromLatin1(c byte) QChar {
	return QChar(uint16(C.qchar_from_latin1(C.char(c))))
}

// FromUcs2 constructs a QChar from a UTF-16 code unit
// (mirrors QChar::fromUcs2).
func FromUcs2(ucs2 uint16) QChar {
	return QChar(uint16(C.qchar_from_ucs2(C.uint16_t(ucs2))))
}

// FromUcs4 constructs a QChar slice (1 or 2 elements) from a UCS-4 code
// point, mirroring QChar::fromUcs4. For BMP code points the result has
// one element; for supplementary planes it has two (a high/low surrogate
// pair).
//
// Use qstringview.FromQChars(slice) if you need a QStringView over the
// result.
func FromUcs4(ucs4 rune) []QChar {
	var out [2]C.uint16_t
	C.qchar_from_ucs4(C.uint32_t(ucs4), &out[0])
	if out[1] == 0 {
		return []QChar{QChar(out[0])}
	}
	return []QChar{QChar(out[0]), QChar(out[1])}
}

// FromRune returns a QChar by truncating r to 16 bits: QChar(uint16(r)).
//
// For BMP code points (U+0000..U+FFFF) this is exact.
// For non-BMP code points this is a silent truncation — the result is
// neither the code point nor its surrogate. Use FromUcs4(r) to obtain a
// proper surrogate pair, or FromRuneStrict(r) to detect the case.
func FromRune(r rune) QChar { return QChar(uint16(r)) }

// FromRuneStrict returns a QChar for r and true if r fits in a single
// UTF-16 code unit (U+0000..U+FFFF), or (0, false) otherwise.
//
// For non-BMP code points (U+10000..U+10FFFF), use FromUcs4(r) to get
// the surrogate pair.
func FromRuneStrict(r rune) (QChar, bool) {
	if r < 0 || r > 0xFFFF {
		return 0, false
	}
	return QChar(uint16(r)), true
}

// ---- Accessors -----------------------------------------------------------

// Unicode returns the underlying UTF-16 code unit.
func (c QChar) Unicode() uint16 {
	return uint16(C.qchar_unicode(C.uint16_t(c)))
}

// SetUnicode replaces the underlying UTF-16 code unit
// (mirrors the non-const char16_t& QChar::unicode() overload).
func (c *QChar) SetUnicode(v uint16) { *c = QChar(v) }

// ToLatin1 returns the Latin-1 byte, or 0 if not representable.
func (c QChar) ToLatin1() byte {
	return byte(C.qchar_to_latin1(C.uint16_t(c)))
}

// Rune returns the code unit as a Go rune (may be a surrogate half).
func (c QChar) Rune() rune { return rune(c) }

// Cell returns the low byte of the UTF-16 code unit (mirrors QChar::cell).
func (c QChar) Cell() byte { return byte(c & 0xFF) }

// Row returns the high byte of the UTF-16 code unit (mirrors QChar::row).
func (c QChar) Row() byte { return byte((c >> 8) & 0xFF) }

// SetCell sets the low byte (mirrors QChar::setCell).
func (c *QChar) SetCell(b byte) { *c = QChar(uint16(*c)&0xFF00 | uint16(b)) }

// SetRow sets the high byte (mirrors QChar::setRow).
func (c *QChar) SetRow(b byte) { *c = QChar(uint16(b)<<8 | (uint16(*c) & 0x00FF)) }

// ---- Predicates (instance) ----------------------------------------------

func (c QChar) IsNull() bool           { return bool(C.qchar_is_null(C.uint16_t(c))) }
func (c QChar) IsPrint() bool          { return bool(C.qchar_is_print(C.uint16_t(c))) }
func (c QChar) IsSpace() bool          { return bool(C.qchar_is_space(C.uint16_t(c))) }
func (c QChar) IsMark() bool           { return bool(C.qchar_is_mark(C.uint16_t(c))) }
func (c QChar) IsLetter() bool         { return bool(C.qchar_is_letter(C.uint16_t(c))) }
func (c QChar) IsNumber() bool         { return bool(C.qchar_is_number(C.uint16_t(c))) }
func (c QChar) IsLetterOrNumber() bool { return bool(C.qchar_is_letter_or_number(C.uint16_t(c))) }
func (c QChar) IsDigit() bool          { return bool(C.qchar_is_digit(C.uint16_t(c))) }
func (c QChar) IsLower() bool          { return bool(C.qchar_is_lower(C.uint16_t(c))) }
func (c QChar) IsUpper() bool          { return bool(C.qchar_is_upper(C.uint16_t(c))) }
func (c QChar) IsTitleCase() bool      { return bool(C.qchar_is_title_case(C.uint16_t(c))) }
func (c QChar) IsPunct() bool          { return bool(C.qchar_is_punct(C.uint16_t(c))) }
func (c QChar) IsSymbol() bool         { return bool(C.qchar_is_symbol(C.uint16_t(c))) }
func (c QChar) IsNonCharacter() bool   { return bool(C.qchar_is_non_character(C.uint16_t(c))) }
func (c QChar) IsHighSurrogate() bool  { return bool(C.qchar_is_high_surrogate(C.uint16_t(c))) }
func (c QChar) IsLowSurrogate() bool   { return bool(C.qchar_is_low_surrogate(C.uint16_t(c))) }
func (c QChar) IsSurrogate() bool      { return bool(C.qchar_is_surrogate(C.uint16_t(c))) }
func (c QChar) HasMirrored() bool      { return bool(C.qchar_has_mirrored(C.uint16_t(c))) }

// ---- Classification (instance) ------------------------------------------

func (c QChar) Category() Category {
	return Category(C.qchar_category(C.uint16_t(c)))
}
func (c QChar) Direction() Direction {
	return Direction(C.qchar_direction(C.uint16_t(c)))
}
func (c QChar) JoiningType() JoiningType {
	return JoiningType(C.qchar_joining_type(C.uint16_t(c)))
}
func (c QChar) CombiningClass() CombiningClass {
	return CombiningClass(C.qchar_combining_class(C.uint16_t(c)))
}
func (c QChar) MirroredChar() QChar {
	return QChar(C.qchar_mirrored_char(C.uint16_t(c)))
}
func (c QChar) DecompositionTag() Decomposition {
	return Decomposition(C.qchar_decomposition_tag(C.uint16_t(c)))
}
func (c QChar) DigitValue() int {
	return int(C.qchar_digit_value(C.uint16_t(c)))
}
func (c QChar) Script() Script {
	return Script(C.qchar_script(C.uint16_t(c)))
}
func (c QChar) UnicodeVersion() UnicodeVersion {
	return UnicodeVersion(C.qchar_unicode_version(C.uint16_t(c)))
}

// ---- Case conversion (instance; returns QChar, as in Qt) -----------------

func (c QChar) ToLower() QChar {
	return QChar(C.qchar_to_lower(C.uint16_t(c)))
}
func (c QChar) ToUpper() QChar {
	return QChar(C.qchar_to_upper(C.uint16_t(c)))
}
func (c QChar) ToTitleCase() QChar {
	return QChar(C.qchar_to_title(C.uint16_t(c)))
}
func (c QChar) ToCaseFolded() QChar {
	return QChar(C.qchar_to_case_folded(C.uint16_t(c)))
}

// ---- Static classification (by UCS-4) ------------------------------------

func CategoryOf(ucs4 rune) Category {
	return Category(C.qchar_static_category(C.uint32_t(ucs4)))
}
func DirectionOf(ucs4 rune) Direction {
	return Direction(C.qchar_static_direction(C.uint32_t(ucs4)))
}
func JoiningTypeOf(ucs4 rune) JoiningType {
	return JoiningType(C.qchar_static_joining_type(C.uint32_t(ucs4)))
}
func CombiningClassOf(ucs4 rune) CombiningClass {
	return CombiningClass(C.qchar_static_combining_class(C.uint32_t(ucs4)))
}
func MirroredCharOf(ucs4 rune) rune {
	return rune(C.qchar_static_mirrored_char(C.uint32_t(ucs4)))
}
func HasMirroredOf(ucs4 rune) bool {
	return bool(C.qchar_static_has_mirrored(C.uint32_t(ucs4)))
}

// DecompositionOf returns the decomposition tag of ucs4
// (mirrors static QChar::decompositionTag(char32_t)).
// It does NOT return the decomposition string; see DecomposeOf for that.
func DecompositionOf(ucs4 rune) Decomposition {
	return Decomposition(C.qchar_static_decomposition_tag(C.uint32_t(ucs4)))
}

func DigitValueOf(ucs4 rune) int {
	return int(C.qchar_static_digit_value(C.uint32_t(ucs4)))
}
func ScriptOf(ucs4 rune) Script {
	return Script(C.qchar_static_script(C.uint32_t(ucs4)))
}
func UnicodeVersionOf(ucs4 rune) UnicodeVersion {
	return UnicodeVersion(C.qchar_static_unicode_version(C.uint32_t(ucs4)))
}

// CurrentUnicodeVersion returns the Unicode version supported by the
// linked Qt (mirrors QChar::currentUnicodeVersion()).
func CurrentUnicodeVersion() UnicodeVersion {
	return UnicodeVersion(C.qchar_static_current_unicode_version())
}

// ---- Static case conversion (by UCS-4; returns rune) ---------------------

func ToLowerOf(ucs4 rune) rune {
	return rune(C.qchar_static_to_lower(C.uint32_t(ucs4)))
}
func ToUpperOf(ucs4 rune) rune {
	return rune(C.qchar_static_to_upper(C.uint32_t(ucs4)))
}
func ToTitleCaseOf(ucs4 rune) rune {
	return rune(C.qchar_static_to_title(C.uint32_t(ucs4)))
}
func ToCaseFoldedOf(ucs4 rune) rune {
	return rune(C.qchar_static_to_case_folded(C.uint32_t(ucs4)))
}

// ---- Static predicates (by UCS-4) ---------------------------------------
//
// Note: there is no IsNull(ucs4) that calls into C — Qt has no static
// QChar::isNull(char32_t) overload. The Go-side IsNull below is a pure
// convenience wrapper.

// IsNull reports whether ucs4 is U+0000.
// Convenience only — Qt has no static QChar::isNull(char32_t).
func IsNull(ucs4 rune) bool { return ucs4 == 0 }

func IsPrint(ucs4 rune) bool {
	return bool(C.qchar_static_is_print(C.uint32_t(ucs4)))
}
func IsSpace(ucs4 rune) bool {
	return bool(C.qchar_static_is_space(C.uint32_t(ucs4)))
}
func IsMark(ucs4 rune) bool {
	return bool(C.qchar_static_is_mark(C.uint32_t(ucs4)))
}
func IsLetter(ucs4 rune) bool {
	return bool(C.qchar_static_is_letter(C.uint32_t(ucs4)))
}
func IsNumber(ucs4 rune) bool {
	return bool(C.qchar_static_is_number(C.uint32_t(ucs4)))
}
func IsLetterOrNumber(ucs4 rune) bool {
	return bool(C.qchar_static_is_letter_or_number(C.uint32_t(ucs4)))
}
func IsDigit(ucs4 rune) bool {
	return bool(C.qchar_static_is_digit(C.uint32_t(ucs4)))
}
func IsLower(ucs4 rune) bool {
	return bool(C.qchar_static_is_lower(C.uint32_t(ucs4)))
}
func IsUpper(ucs4 rune) bool {
	return bool(C.qchar_static_is_upper(C.uint32_t(ucs4)))
}
func IsTitleCase(ucs4 rune) bool {
	return bool(C.qchar_static_is_title_case(C.uint32_t(ucs4)))
}
func IsPunct(ucs4 rune) bool {
	return bool(C.qchar_static_is_punct(C.uint32_t(ucs4)))
}
func IsSymbol(ucs4 rune) bool {
	return bool(C.qchar_static_is_symbol(C.uint32_t(ucs4)))
}
func IsNonCharacter(ucs4 rune) bool {
	return bool(C.qchar_static_is_non_character(C.uint32_t(ucs4)))
}

// IsHighSurrogateRune reports whether ucs4 is a UTF-16 high surrogate
// (U+D800..U+DBFF). Named ...Rune to avoid confusion with the instance
// method IsHighSurrogate on QChar.
func IsHighSurrogateRune(ucs4 rune) bool {
	return bool(C.qchar_static_is_high_surrogate(C.uint32_t(ucs4)))
}

// IsLowSurrogateRune reports whether ucs4 is a UTF-16 low surrogate
// (U+DC00..U+DFFF).
func IsLowSurrogateRune(ucs4 rune) bool {
	return bool(C.qchar_static_is_low_surrogate(C.uint32_t(ucs4)))
}

// IsSurrogateRune reports whether ucs4 is a UTF-16 surrogate
// (U+D800..U+DFFF).
func IsSurrogateRune(ucs4 rune) bool {
	return bool(C.qchar_static_is_surrogate(C.uint32_t(ucs4)))
}

// ---- Surrogate pairs -----------------------------------------------------

// RequiresSurrogates reports whether ucs4 needs a surrogate pair in UTF-16.
func RequiresSurrogates(ucs4 rune) bool {
	return bool(C.qchar_requires_surrogates(C.uint32_t(ucs4)))
}

// HighSurrogate returns the high surrogate of ucs4.
func HighSurrogate(ucs4 rune) QChar {
	return QChar(uint16(C.qchar_high_surrogate(C.uint32_t(ucs4))))
}

// LowSurrogate returns the low surrogate of ucs4.
func LowSurrogate(ucs4 rune) QChar {
	return QChar(uint16(C.qchar_low_surrogate(C.uint32_t(ucs4))))
}

// SurrogateToUcs4 combines a high/low surrogate pair into a UCS-4 rune.
func SurrogateToUcs4(high, low QChar) rune {
	return rune(C.qchar_surrogate_to_ucs4(C.uint16_t(high), C.uint16_t(low)))
}

// ---- Decomposition string (QString-producing) ---------------------------
//
// QChar::decomposition() returns a QString. Because qstring imports qchar
// (for At/Front/Back returning QChar), returning *qstring.QString here
// would create an import cycle. Instead, the raw C++ QString* is returned
// as unsafe.Pointer; the caller wraps it via the qstring package.
//
// Naming: the enum type Decomposition and the tag-static DecompositionOf()
// already exist; the string-returning functions are named Decompose /
// DecomposeOf to avoid collision.
//
// Typical usage:
//
//	import (
//	    "github.com/vadimlxx/go-qt6core/qchar"
//	    "github.com/vadimlxx/go-qt6core/qstring"
//	)
//
//	s := qstring.AdoptUnsafePtr(c.Decompose())
//	defer s.Close()
//
// or, more conveniently:
//
//	s := qstring.FromDecomposition(c)
//	defer s.Close()

// Decompose returns the Unicode decomposition of c as a freshly allocated
// C++ QString* wrapped in unsafe.Pointer
// (mirrors QChar::decomposition()).
//
// The caller owns the returned pointer and must wrap it via
// qstring.AdoptUnsafePtr, or use qstring.FromDecomposition directly.
// For code points with no decomposition, the returned QString is empty
// (non-nil).
func (c QChar) Decompose() unsafe.Pointer {
	var out unsafe.Pointer
	C.qchar_decomposition(C.uint32_t(c), &out)
	return out
}

// DecomposeOf is Decompose by UCS-4 code point
// (mirrors static QChar::decomposition(char32_t)).
//
// Same ownership/lifetime contract as Decompose.
func DecomposeOf(ucs4 rune) unsafe.Pointer {
	var out unsafe.Pointer
	C.qchar_decomposition(C.uint32_t(ucs4), &out)
	return out
}
