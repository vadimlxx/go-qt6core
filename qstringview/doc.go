// Package qstringview provides a Go binding for Qt 6.11
// QtCore's QStringView.
//
// QStringView is a non-owning read-only view over a UTF-16 buffer
// (a (const char16_t *, qsizetype) pair, Q_PRIMITIVE_TYPE). The Go
// wrapper is a value type that carries a keepAlive reference to the
// underlying buffer to prevent premature GC.
//
// Interop with sibling packages:
//   - At / Front / Back return qchar.QChar.
//   - ToUtf8 / ToLatin1 / ToLocal8Bit return *qbytearray.QByteArray.
//   - FromQChars accepts []qchar.QChar (zero-copy).
//
// Lifetime rules (same as QByteArrayView):
//   - Slicing methods (Left/Right/Mid/First/Last/Sliced/SlicedN/
//     Chopped/Trimmed) return views of the same buffer; keepAlive is
//     propagated.
//   - Do not mutate the source slice/string/qchar slice while a view
//     into it is alive.
//
// Not yet bound (documented as TODO):
//   - toString()            — requires package qstring.
//   - split()               — returns QList<QStringView>; will appear
//                             together with a qstringviewlist package.
//   - arg()                 — varargs template.
//   - toUcs4()              — returns QList<uint>.
//   - toWCharArray()        — wchar_t platform size mismatch.
//   - localeAwareCompare()  — requires QLocale.
//   - QLatin1StringView overloads (compare, startsWith, endsWith,
//     indexOf, contains, count, lastIndexOf) — will be added with the
//     qlatin1stringview package.
//   - QRegularExpression / QRegularExpressionMatch overloads.
//   - Darwin-specific toCFString / toNSString.
//
// This package is a subpackage of github.com/vadimlxx/go-qt6core.
// Requires Qt >= 6.11 and a C++17 compiler; links via pkg-config Qt6Core.
package qstringview