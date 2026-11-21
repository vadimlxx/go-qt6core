// Package qtext provides Go bindings for Qt 6.11 QtCore QString.
//
// QString is a UTF-16, COW (copy-on-write) string container. In Go, each
// *QString is a handle to a heap-allocated C++ QString; cloning is cheap
// (refcount++), and mutating methods detach as in Qt.
//
// Memory management
//
// The C++ QString is released by one of:
//   - Close(): eager, deterministic release. Use this in long-running
//     processes when you know a QString is no longer needed.
//   - runtime.AddCleanup: safety net. If Close is never called, the C++
//     object is freed once the *QString becomes unreachable.
//
// Calling Close is optional. Calling Close twice, or on nil, is safe.
// Using a QString after Close panics with "qstring: use after Close".
//
// Construction interop:
//   - FromString / FromUtf8       → qstring_new_from_utf8      (QString::fromUtf8)
//   - FromUtf16 / FromQStringView → qstring_new_from_utf16     (QString::fromUtf16)
//   - FromLatin1                  → qstring_new_from_latin1    (QString::fromLatin1)
//   - FromLocal8Bit               → qstring_new_from_local8bit (QString::fromLocal8Bit)
//   - FromQByteArray              → qstring_new_from_qbytearray (QString::fromUtf8(ba))
//   - QStringFromUcs4                     → QString::fromUcs4
//   - QStringFromPercentEncoding          → emulates removed QString::fromPercentEncoding
//                                           (via QByteArray::fromPercentEncoding)
//
// Sibling interop (all in package qtext):
//   - At / Front / Back return QChar.
//   - Utf8 / Latin1 / Local8Bit return *QByteArray.
//   - ToQStringView returns a zero-copy QStringView that keeps the *QString
//     alive via the owner field.
//   - ToUcs4 returns []uint32 (fresh copy of code points).
//   - IndexOfQRegularExpression / LastIndexOfQRegularExpression /
//     ContainsQRegularExpression / CountQRegularExpression /
//     SectionQRegularExpression / ReplaceQRegularExpression /
//     RemoveQRegularExpression / SplitQRegularExpression accept
//     *QRegularExpression (defined in qregexp.go).
//   - ToPercentEncoding returns *QByteArray (percent-encoded UTF-8).
//     Emulates removed QString::toPercentEncoding via toUtf8().toPercentEncoding.
//   - Swap, Clone operate on COW refcounts only, no data copy.
//
// Cross-package shared types:
//   - qdata.SplitBehavior (Qt::SplitBehavior) is used by QString.SplitQString,
//     SplitQChar and SplitQRegularExpression.
//
// Not yet bound:
//   - arg() as a single variadic call chaining %1..%99 — use explicit
//     per-type Arg* methods.
//   - Vasprintf / Asprintf / Sprintf.
//   - toWCharArray / fromWCharArray — wchar_t width mismatch.
//   - toStdString / toStdU16String / toStdU32String — std::basic_string
//     interop needs a separate cgo layer.
//   - toPercentEncoding / fromPercentEncoding.
//   - Reverse-iterators (rbegin/rend/crbegin/crend).
//   - QChar* constructors and setRawData / fromRawData.
//   - QLatin1StringView overloads for compare/startsWith/…/contains.
//   - operator+ / operator+= — use Clone + Append* instead.
//   - QString(nullptr) null constructor — use NewEmptyQString (empty, non-null).
//
// This package is a subpackage of github.com/vadimlxx/go-qt6core.
// Requires Qt >= 6.11 and a C++17 compiler; links via pkg-config Qt6Core.
package qtext
