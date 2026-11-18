package qtext

/*
  QLatin1Char — обёртка над одноимённым классом Qt 6.11 (один Latin-1 байт).

  Семантика

  QLatin1Char — value-тип (byte): тривиально сравним, без аллокаций и
  финализатора. Нулевое значение эквивалентно QLatin1Char('\x00').

  Что сознательно не вошло

  - QLatin1Char(char) — покрывается NewQLatin1Char(byte).
  - QLatin1Char(QChar) — покрывается NewQLatin1CharFromQChar.
  - Сравнения через операторы (==, <) — заменены на Equal / Compare и
    EqualByte / CompareByte.
  - UDL `operator""_L1` — в Go нет user-defined literals.

  Особенности дизайна

  - String() возвращает UTF-8-строку: байты 0x00..0x7F — ASCII,
    0x80..0xFF — двухбайтовый UTF-8 code point (Latin-1 == Unicode
    для 0x00..0xFF).
  - Все методы QLatin1Char работают с одним байтом; cgo используется
    ради единообразия и единственной точки проверки версии Qt.
*/

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC
#cgo LDFLAGS: -lstdc++
#include <stdint.h>
#include <stdbool.h>
#include "qchar_c.h"
*/
import "C"

// QLatin1Char represents a single Latin-1 character (C++ QLatin1Char).
// It is a value type: the zero value equals QLatin1Char('\x00').
type QLatin1Char byte

// ---- Constructors --------------------------------------------------------

// NewQLatin1Char returns a QLatin1Char for the given Latin-1 byte.
// Mirrors QLatin1Char(char).
func NewQLatin1Char(c byte) QLatin1Char { return QLatin1Char(c) }

// NewQLatin1CharFromRune returns a QLatin1Char if r fits in U+0000..U+00FF,
// or (0, false) otherwise. Latin-1 byte values equal their Unicode code points
// in this range.
//
// Go-only helper; equivalent to QLatin1Char::fromLatin1 with a range check.
func NewQLatin1CharFromRune(r rune) (QLatin1Char, bool) {
	if r < 0 || r > 0xFF {
		return 0, false
	}
	return QLatin1Char(byte(r)), true
}

// NewQLatin1CharFromQChar returns a QLatin1Char if q lies in U+0000..U+00FF,
// or (0, false) otherwise.
//
// Mirrors QLatin1Char(QChar) with an explicit range check.
func NewQLatin1CharFromQChar(q QChar) (QLatin1Char, bool) {
	if q.Unicode() > 0xFF {
		return 0, false
	}
	return QLatin1Char(byte(q.Unicode())), true
}

// ---- Accessors -----------------------------------------------------------

// ToLatin1 returns the underlying Latin-1 byte.
// Mirrors QLatin1Char::toLatin1().
func (c QLatin1Char) ToLatin1() byte {
	return byte(C.qlatin1char_to_latin1(C.uint8_t(c)))
}

// Unicode returns the Unicode code point (always in 0x00..0xFF).
// Mirrors QLatin1Char::unicode().
func (c QLatin1Char) Unicode() uint16 {
	return uint16(C.qlatin1char_unicode(C.uint8_t(c)))
}

// ToQChar widens to QChar. Equivalent to QChar(QLatin1Char(c)).
func (c QLatin1Char) ToQChar() QChar {
	return NewQCharFromLatin1(byte(c))
}

// String returns the character as a Go string.
// Bytes 0x00..0x7F become ASCII; 0x80..0xFF become their two-byte UTF-8
// encoding (Latin-1 == Unicode for this range).
//
// Go-only helper.
func (c QLatin1Char) String() string { return string(rune(byte(c))) }

// ---- Comparisons --------------------------------------------------------

// Equal reports whether c == other.
// Mirrors QLatin1Char's generated operator== (bit equality).
func (c QLatin1Char) Equal(other QLatin1Char) bool {
	return bool(C.qlatin1char_compares_equal(C.uint8_t(c), C.uint8_t(other)))
}

// Compare returns -1, 0, or +1 using Qt's strongly-ordered comparison
// (unsigned interpretation of the stored byte).
// Mirrors QLatin1Char's generated operator<=>.
func (c QLatin1Char) Compare(other QLatin1Char) int {
	return int(C.qlatin1char_compare(C.uint8_t(c), C.uint8_t(other)))
}

// EqualByte reports whether c equals the raw byte b.
// Mirrors QLatin1Char's generated operator==(QLatin1Char, char).
func (c QLatin1Char) EqualByte(b byte) bool {
	return bool(C.qlatin1char_compares_equal_byte(C.uint8_t(c), C.uint8_t(b)))
}

// CompareByte returns -1, 0, or +1 by comparing c against the raw byte b
// (unsigned interpretation).
// Mirrors QLatin1Char's generated operator<=>(QLatin1Char, char).
func (c QLatin1Char) CompareByte(b byte) int {
	return int(C.qlatin1char_compare_byte(C.uint8_t(c), C.uint8_t(b)))
}
