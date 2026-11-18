// QChar is a Go binding for Qt 6.11 QtCore's QChar class.
//
// QChar is a UTF-16 code unit (Q_PRIMITIVE_TYPE). The Go wrapper is
// `type QChar uint16`: trivially comparable, hashable, allocation-free,
// no finalizer, no Close(). All methods are safe for concurrent use.
//
// # Type mapping
//
//	Qt                              Go
//	QChar                            qtext.QChar  (uint16)
//	QChar::SpecialCharacter          qtext.QChar  constants (Null, Space, …)
//	QChar::Category                  qtext.Category        (int)
//	QChar::Direction                 qtext.Direction       (int)
//	QChar::JoiningType               qtext.JoiningType     (int)
//	QChar::Decomposition             qtext.Decomposition   (int)
//	QChar::CombiningClass            qtext.CombiningClass  (uint8)
//	QChar::Script                    qtext.Script          (int)
//	QChar::UnicodeVersion            qtext.UnicodeVersion  (int)
//	char16_t                         uint16
//	char32_t (UCS-4)                 rune
//
// # Bound API (summary)
//
// Constructors:
//   - NewQChar(uint16)               — QChar(char16_t)
//   - NewQCharFromLatin1(byte)       — QChar::fromLatin1(char)
//   - NewQCharFromUcs2(uint16)       — QChar::fromUcs2(char16_t)
//   - NewQCharFromUcs4(rune)         — QChar::fromUcs4(char32_t) → []QChar
//   - NewQCharFromRune(rune)         — truncating
//   - NewQCharFromRuneStrict(rune)   — (QChar, bool); false for non-BMP
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
// Decomposition:
//   (QChar).Decompose()       *QString      — QChar::decomposition() const
//   (QChar).DecomposeString() string        — Go-string variant
//   QCharDecompose(rune)      *QString      — QChar::decomposition(char32_t)
//   QCharDecomposeString(rune) string       — Go-string variant
//
// # Interop with sibling packages
//
// qchar is the foundation of the string stack:
//
//	qstring           → imports qchar (At/Front/Back return QChar)
//	qstringview       → imports qchar
//	qlatin1char       → imports qchar (ToQChar)
//	qlatin1stringview → imports qchar
//	qutf8stringview   → imports qchar
//
// # Not yet bound / API differences vs Qt 6.11
//
//   - QChar vs const char* comparison (Qt's comparesEqual(QChar, const char*),
//     Q_DECLARE_STRONGLY_ORDERED(QChar, const char*)). Go has no operator
//     overloading; semantics of Qt's compare_helper are non-obvious. Use
//     qstring/qstringview and compare at the string level.
//
//   - QChar(uchar row, uchar cell) two-byte constructor.
//     Covered by NewQChar(uint16(row)<<8 | uint16(cell)).
//
//   - QChar::fromUcs4 return type. Qt returns a small struct with
//     `operator QStringView()`; Go returns []QChar (length 1 for BMP,
//     length 2 for supplementary planes). Use qstringview.FromQChars(slice)
//     if a view is needed.
//
//   - QChar::SpecialCharacter::LastValidCodePoint (0x10FFFF) does not fit
//     in uint16. Exposed as `const LastValidCodePoint rune`.
//
//   - QChar(QLatin1Char). Covered by NewQCharFromLatin1(byte) or
//     qlatin1char.QLatin1Char.ToQChar().
//
//   - QChar(char16_t) / QChar(ushort) / QChar(short) / QChar(uint) /
//     QChar(int) / QChar(char) / QChar(uchar). All covered by NewQChar(uint16).
//
//   - QChar(wchar_t) on Windows. Go has no wchar_t.
//
//   - std::hash<QChar>. `type QChar uint16` is natively hashable in Go.
//
//   - QT_CHAR_FASTCALL, Q_DECL_CONST_FUNCTION, Q_DECLARE_TYPEINFO,
//     Q_DECLARE_STRONGLY_ORDERED_LITERAL_TYPE, QT6_ONLY, QT7_ONLY,
//     QT_CORE_REMOVED_SINCE(6,9) — Qt compiler macros / version guards.
//
//   - Private helper namespace (QChar::isSpace_helper, isLetter_helper,
//     isNumber_helper, isLetterOrNumber_helper, equal_helper, compare_helper).
//
// # Design notes
//
//   - QChar is a value type: no heap allocation, no finalizer, no Close().
//     `QChar` and `uint16` are interchangeable via a single cast.
//
//   - All instance methods are pure except SetUnicode, SetCell, SetRow.
//
//   - NewQCharFromRune truncates a rune to 16 bits. Use NewQCharFromUcs4
//     or NewQCharFromRuneStrict for checked conversion.
//
//   - NewQCharFromUcs4 returns nil for out-of-range code points
//     (< 0 or > LastValidCodePoint). Qt's QChar::fromUcs4 has no such
//     check (UB for invalid input).
//
//   - QCharDecompose returns nil for out-of-range code points and on
//     allocation failure. The returned *QString is owned by the caller.
//
//   - CurrentUnicodeVersion() exposes the linked Qt's Unicode database
//     version (compile-time constant depending on Qt build).
//
//   - The C ABI lives in qchar_c.h as a dual-mode header (C declarations
//     for cgo's C preamble + C++ inline definitions compiled via
//     qtext_cxx_impl.cpp). Per-function notes live in qchar.go, not in
//     the C header, because one C++ function may map to a Go function
//     with a different name or be split into two.
//
// # Requirements
//
// Qt >= 6.11, C++17, pkg-config Qt6Core. Version check is performed at
// compile time via internal/qt6core/check_qt_version.h.
//
// This package is a subpackage of github.com/vadimlxx/go-qt6core.
package qtext
