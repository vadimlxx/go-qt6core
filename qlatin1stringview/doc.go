// Package qlatin1stringview provides a Go binding for Qt 6.11
// QtCore's QLatin1StringView (and its alias QLatin1String).
//
// QLatin1StringView is a non-owning read-only view over a Latin-1
// byte buffer (a (const char *, qsizetype) pair, Q_RELOCATABLE_TYPE).
// The Go wrapper is a value type that carries a keepAlive reference
// to the underlying buffer.
//
// In Qt 6.x, QLatin1String is a synonym (typedef) of QLatin1StringView.
// In Go, package users should prefer QLatin1StringView; QLatin1String
// exists as `type QLatin1String = QLatin1StringView` for Qt source parity.
//
// Interop with sibling packages:
//   - At / Front / Back return qlatin1char.QLatin1Char.
//   - ToUtf8 returns *qbytearray.QByteArray.
//   - Compare / StartsWith / EndsWith / IndexOf / LastIndexOf / Contains /
//     Count accept qstringview.QStringView and qchar.QChar variants.
//
// Lifetime rules (same as QByteArrayView):
//   - Slicing methods (Left/Right/Mid/First/Last/Sliced/SlicedN/Chopped/
//     Trimmed) return views of the same buffer; keepAlive is propagated.
//   - Do not mutate the source slice/string/QByteArray while a view
//     into it is alive.
//
// FromString semantics: s's bytes are interpreted as raw Latin-1.
// For ASCII strings this is exact; for UTF-8 strings with non-ASCII
// code points, each UTF-8 byte becomes one Latin-1 character (the same
// way QByteArray → QLatin1StringView behaves in Qt). To convert a Go
// UTF-8 string to a proper Latin-1 representation, encode manually or
// wait for the qstring package.
//
// TODO (not yet bound):
//   - toString()       — returns QString; requires package qstring.
//   - arg()            — varargs template returning QString.
//   - tokenize()       — returns QStringTokenizer.
//   - toWCharArray()   — wchar_t platform size mismatch.
//   - toUcs4()         — returns QList<uint>.
//   - QLatin1StringView toString comparisons with QString (via qstring).
//   - QRegularExpression / QRegularExpressionMatch overloads.
//
// This package is a subpackage of github.com/vadimlxx/go-qt6core.
// Requires Qt >= 6.11 and a C++17 compiler; links via pkg-config Qt6Core.
package qlatin1stringview
