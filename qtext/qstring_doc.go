// Package qtext provides Go bindings for Qt 6.11 QtCore QString.
//
// QString is a UTF-16, COW (copy-on-write) string container. In Go, each
// *QString is a handle to a heap-allocated C++ QString; cloning is cheap
// (refcount++), and mutating methods detach as in Qt.
//
// Memory management: runtime.SetFinalizer; call Close() for immediate
// release.
//
// Construction interop:
//   - FromString / FromUtf8       → qstring_new_from_utf8      (QString::fromUtf8)
//   - FromUtf16 / FromQStringView → qstring_new_from_utf16     (QString::fromUtf16)
//   - FromLatin1                  → qstring_new_from_latin1    (QString::fromLatin1)
//   - FromLocal8Bit               → qstring_new_from_local8bit (QString::fromLocal8Bit)
//   - FromQByteArray              → qstring_new_from_qbytearray (QString::fromUtf8(ba))
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
//
// Cross-package shared types:
//   - qdata.SplitBehavior (Qt::SplitBehavior) is used by QString.SplitQString,
//     SplitQChar and SplitQRegularExpression.
//
// Implemented (as of this revision):
//   - arg: ArgLl, ArgUll, ArgDouble, ArgQChar, ArgQString, ArgString.
//   - compare: CompareQString, CompareQChar, CompareString, LocaleAwareCompare.
//   - normalization: Normalized, NormalizedWithVersion.
//   - split: SplitQString, SplitQChar, SplitQRegularExpression.
//   - section: SectionQChar, SectionQString, SectionQRegularExpression.
//   - escape: ToHtmlEscaped.
//   - repeat: Repeated.
//   - low-level accessors: Capacity, IsSharedWith, ResizeForOverwrite.
//   - remove with cs: RemoveQChar, RemoveQString, RemoveString,
//     RemoveQRegularExpression.
//   - replace with cs: ReplaceQCharQCharCS, ReplaceQCharQStringCS,
//     ReplaceQStringQStringCS, ReplaceStringStringCS;
//     ReplaceQRegularExpression, ReplaceQRegularExpressionString.
//   - setNum/number for int/uint/float: SetNumInt, SetNumUInt, SetNumFloat,
//     QStringNumberInt, QStringNumberUInt.
//   - FromLocal8Bit constructor.
//
// TODO (not yet bound):
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