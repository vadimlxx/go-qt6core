package qtext

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC -Wno-sfinae-incomplete
#cgo LDFLAGS: -lstdc++
#include <stdint.h>
#include <stdbool.h>
#include "qlatin1char_c.h"
*/
import "C"

// QLatin1Char represents a single Latin-1 character (C++ QLatin1Char).
// It is a value type: the zero value equals QLatin1Char('\x00').
type QLatin1Char byte

// ---- construction --------------------------------------------------------

// New returns a QLatin1Char wrapping the given Latin-1 byte.
// Mirrors QLatin1Char(char).
func NewQLatin1Char(c byte) QLatin1Char { return QLatin1Char(c) }

// FromRune returns a QLatin1Char if r fits in U+0000..U+00FF,
// or (0, false) otherwise. Latin-1 byte values equal their Unicode
// code points for this range.
func QLatin1CharFromRune(r rune) (QLatin1Char, bool) {
	if r < 0 || r > 0xFF {
		return 0, false
	}
	return QLatin1Char(byte(r)), true
}

// FromQChar converts a QChar to QLatin1Char if it lies in
// U+0000..U+00FF, or (0, false) otherwise.
func QLatin1CharFromQChar(q QChar) (QLatin1Char, bool) {
	if q.Unicode() > 0xFF {
		return 0, false
	}
	return QLatin1Char(byte(q.Unicode())), true
}

// ---- accessors -----------------------------------------------------------

// ToLatin1 returns the underlying Latin-1 byte.
// Mirrors QLatin1Char::toLatin1().
func (c QLatin1Char) ToLatin1() byte {
	r := C.qlatin1char_to_latin1(C.uint8_t(c))
	return byte(r)
}

// Unicode returns the Unicode code point (always in 0x00..0xFF).
// Mirrors QLatin1Char::unicode().
func (c QLatin1Char) Unicode() uint16 {
	r := C.qlatin1char_unicode(C.uint8_t(c))
	return uint16(r)
}

// ToQChar widens to QChar (equivalent to QChar(QLatin1Char(c))).
func (c QLatin1Char) ToQChar() QChar {
	return QCharFromLatin1(byte(c))
}

// String returns the character as a Go string.
// Latin-1 bytes 0x00..0x7F become ASCII, 0x80..0xFF become their
// two-byte UTF-8 encoding.
func (c QLatin1Char) String() string {
	return string(byte(c))
}

// ---- comparisons --------------------------------------------------------

// Equal reports whether c == other (same stored byte).
func (c QLatin1Char) Equal(other QLatin1Char) bool {
	r := C.qlatin1char_compares_equal(C.uint8_t(c), C.uint8_t(other))
	return bool(r)
}

// Compare returns -1, 0, or +1 by comparing the two Latin-1 chars
// using Qt's strongly-ordered comparison (unsigned interpretation).
func (c QLatin1Char) Compare(other QLatin1Char) int {
	return int(C.qlatin1char_compare(C.uint8_t(c), C.uint8_t(other)))
}

// EqualByte reports whether c equals the raw Latin-1 byte b.
func (c QLatin1Char) EqualByte(b byte) bool {
	r := C.qlatin1char_compares_equal_byte(C.uint8_t(c), C.uint8_t(b))
	return bool(r)
}

// CompareByte returns -1, 0, or +1 by comparing c against the raw byte b.
func (c QLatin1Char) CompareByte(b byte) int {
	return int(C.qlatin1char_compare_byte(C.uint8_t(c), C.uint8_t(b)))
}
