// Package qchar provides a Go binding for Qt 6.11 QtCore's QChar class.
//
// QChar is a UTF-16 code unit (Q_PRIMITIVE_TYPE). The Go wrapper is
// `type QChar uint16`: trivially comparable, hashable, allocation-free,
// no finalizer, no Close(). All methods are safe for concurrent use.
//
// # Type mapping
//
//	Qt                              Go
//	QChar                            qchar.QChar  (uint16)
//	QChar::SpecialCharacter          qchar.QChar  constants (Null, Space, …)
//	QChar::Category                  qchar.Category        (int)
//	QChar::Direction                 qchar.Direction       (int)
//	QChar::JoiningType               qchar.JoiningType     (int)
//	QChar::Decomposition             qchar.Decomposition   (int)
//	QChar::CombiningClass            qchar.CombiningClass  (uint8)
//	QChar::Script                    qchar.Script          (int)
//	QChar::UnicodeVersion            qchar.UnicodeVersion  (int)
//	char16_t                         uint16
//	char32_t (UCS-4)                 rune
//
// # Bound API (summary)
//
// Constructors:
//   - New(uint16)          — from a UTF-16 code unit
//   - FromUcs2(uint16)     — mirrors QChar::fromUcs2
//   - FromLatin1(byte)     — mirrors QChar::fromLatin1
//   - FromRune(rune)       — truncating (see "Design notes")
//   - FromRuneStrict(rune) — returns (QChar, bool); false for non-BMP
//   - FromUcs4(rune)       — []QChar (1 or 2 code units)
//
// Predicates (instance + static by rune):
//   IsNull, IsPrint, IsSpace, IsMark, IsLetter, IsNumber,
//   IsLetterOrNumber, IsDigit, IsLower, IsUpper, IsTitleCase,
//   IsPunct, IsSymbol, IsNonCharacter, IsHighSurrogate, IsLowSurrogate,
//   IsSurrogate, HasMirrored
//
// Classification:
//   Category, Direction, JoiningType, CombiningClass, MirroredChar,
//   DecompositionTag, DigitValue, Script, UnicodeVersion
//
// Case conversion:
//   ToLower / ToUpper / ToTitleCase / ToCaseFolded              (instance, → QChar)
//   ToLowerOf / ToUpperOf / ToTitleCaseOf / ToCaseFoldedOf      (static, → rune)
//
// Accessors:
//   Unicode, SetUnicode, ToLatin1, Rune, Cell, Row, SetCell, SetRow
//
// Surrogates:
//   RequiresSurrogates, HighSurrogate, LowSurrogate, SurrogateToUcs4
//
// Version info:
//   CurrentUnicodeVersion
//
// Decomposition (see "Interop with sibling packages"):
//   Decomposition(QChar) unsafe.Pointer
//   DecompositionOf(rune) unsafe.Pointer
//
// # Interop with sibling packages
//
// qchar is the foundation of the string stack:
//
//	qstring           → imports qchar (At/Front/Back return QChar)
//	qstringview       → imports qchar (At/Front/Back return QChar)
//	qlatin1char       → imports qchar (ToQChar)
//	qlatin1stringview → imports qchar (Compare/StartsWith/… with QChar)
//	qutf8stringview   → imports qchar (CompareQChar/EqualQChar)
//
// # Not yet bound / API differences vs Qt 6.11
//
// The following Qt API is either not exposed or exposed differently:
//
//   - QChar vs const char* comparison (Qt's comparesEqual(QChar, const char*),
//     Q_DECLARE_STRONGLY_ORDERED(QChar, const char*)).
//     Go has no operator overloading. Semantics are non-obvious
//     (Qt's QChar::compare_helper is private and works via QUtf8StringView
//     internally). Tracked as TODO; for now use qstring.FromUtf16 / .Utf8()
//     and compare at the string level.
//
//   - QChar(uchar row, uchar cell) two-byte constructor.
//     Covered by New(uint16(row)<<8 | uint16(cell)). Not exposed separately.
//
//   - QChar::fromUcs4 return type. Qt returns a small struct with
//     `operator QStringView()`; Go returns []QChar (length 1 for BMP,
//     length 2 for supplementary planes). Use qstringview.FromQChars(slice)
//     if a view is needed (qstringview already imports qchar, so no cycle).
//
//   - QChar::SpecialCharacter::LastValidCodePoint (0x10FFFF) does not fit
//     in uint16. Exposed as `const LastValidCodePoint rune`, not as QChar.
//
//   - QChar(QLatin1Char). Covered by FromLatin1(byte) or
//     qlatin1char.QLatin1Char.ToQChar().
//
//   - QChar(char16_t) / QChar(ushort) / QChar(short) / QChar(uint) /
//     QChar(int) / QChar(char) / QChar(uchar).
//     All are covered by New(uint16). The Qt-side explicit/implicit
//     distinction is a C++ compile-time concern; Go has no such issue.
//     Use FromRuneStrict for checked rune→QChar conversion.
//
//   - QChar(wchar_t) on Windows. Go has no wchar_t; use rune on Unix
//     builds, or explicit uint16 conversions on Windows.
//
//   - std::hash<QChar> specialization.
//     `type QChar uint16` is natively hashable in Go (`==` and map keys).
//
//   - QT_CHAR_FASTCALL, Q_DECL_CONST_FUNCTION, Q_DECLARE_TYPEINFO,
//     Q_DECLARE_STRONGLY_ORDERED_LITERAL_TYPE, QT6_ONLY, QT7_ONLY,
//     QT_CORE_REMOVED_SINCE(6,9) — Qt compiler macros and version-guards.
//     Implementation details; not exposed.
//
//   - The private helper namespace (QChar::isSpace_helper,
//     QChar::isLetter_helper, QChar::isNumber_helper,
//     QChar::isLetterOrNumber_helper, QChar::equal_helper,
//     QChar::compare_helper). Internal to Qt; not part of the public API.
//
// # Design notes
//
//   - QChar is a value type: no heap allocation, no finalizer, no Close().
//     `QChar` and `uint16` are interchangeable via a single cast.
//
//   - All instance methods are pure (no mutation) except SetUnicode,
//     SetCell, SetRow, which mirror Qt's non-const accessors.
//
//   - FromRune truncates a rune to 16 bits. For non-BMP code points use
//     FromUcs4 (returns a surrogate pair) or FromRuneStrict (returns false).
//
//   - CurrentUnicodeVersion() exposes the linked Qt's Unicode database
//     version; this is a compile-time constant that depends on how Qt
//     was built.
//
//   - Go-bounds-checks on methods that access by index (none in QChar:
//     all methods are total, since every QChar wraps a valid uint16).
//     The truncated FromRune is the only potential footgun; FromRuneStrict
//     is provided to avoid it.
//
// # Requirements
//
// Qt >= 6.11, C++17, pkg-config Qt6Core. Version check is performed at
// compile time in qchar_c.cpp via internal/qt6core/check_qt_version.h.
//
// This package is a subpackage of github.com/vadimlxx/go-qt6core.
package qtext
