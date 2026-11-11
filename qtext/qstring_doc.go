// Package qstring provides a Go binding for Qt 6.11 QtCore's QString.
//
// QString is a UTF-16, COW (copy-on-write) string container. In Go, each
// *QString is a handle to a heap-allocated C++ QString; cloning is cheap
// (refcount++), and mutating methods detach as in Qt.
//
// Memory management: runtime.SetFinalizer; call Close() for immediate
// release.
//
// Construction interop:
//   - FromString / FromUtf8       → qstring_new_from_utf8     (QString::fromUtf8)
//   - FromUtf16 / FromQStringView → qstring_new_from_utf16    (QString::fromUtf16)
//   - FromLatin1                  → qstring_new_from_latin1   (QString::fromLatin1)
//   - FromQByteArray              → qstring_new_from_qbytearray (QString::fromUtf8(ba))
//
// All constructors and methods use direct C++ bindings from qstring_c.cpp;
// this package no longer imports internal/qt6core.
//
// Sibling interop (all in package qtext):
//   - At / Front / Back return QChar.
//   - Utf8 / Latin1 / Local8Bit return *QByteArray (via QByteArrayFromUnsafePtr).
//   - ToQStringView returns a zero-copy QStringView that keeps the *QString
//     alive via the owner field.
//
// TODO (not yet bound):
//   - arg()                       — variadic template.
//   - split() / splitRef()        — return QStringList / QList<QStringView>.
//   - QRegularExpression / QRegularExpressionMatch overloads.
//   - localeAwareCompare()        — requires QLocale.
//   - toHtmlEscaped()             — cross-cutting escaping.
//   - toPercentEncoding() / fromPercentEncoding().
//   - toStdString() / fromStdString()       — std::string interop.
//   - toStdU16String() / fromStdU16String() — std::u16string interop.
//   - toStdU32String() / fromStdU32String() — std::u32string interop.
//   - toWCharArray()              — wchar_t platform size mismatch.
//   - toUcs4()                    — returns QList<uint>.
//   - vasprintf / asprintf / sprintf methods.
//   - reverse_iterator methods (rbegin/rend/crbegin/crend).
//   - QChar* constructors and setRawData / fromRawData.
//   - QLatin1StringView overloads for compare/startsWith/…/contains.
//   - operator+ / operator+= (use Clone+Append* instead).
//
// This package is a subpackage of github.com/vadimlxx/go-qt6core.
// Requires Qt >= 6.11 and a C++17 compiler; links via pkg-config Qt6Core.
package qtext
