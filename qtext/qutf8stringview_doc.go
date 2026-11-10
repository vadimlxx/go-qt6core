// Package qutf8stringview provides a Go binding for Qt 6.11
// QtCore's QUtf8StringView.
//
// QUtf8StringView is a non-owning read-only view over a UTF-8 byte buffer.
// In Qt 6.11, QUtf8StringView is a typedef for QBasicUtf8StringView<false>
// (storage_type = char); the class is a primitive value type
// (Q_PRIMITIVE_TYPE). The Go wrapper is a value type carrying a keepAlive
// reference to the underlying buffer.
//
// Interop with sibling packages:
//   - At / Front / Back return byte.
//   - ToString returns *qstring.QString.
//   - Compare/Equal variants accept:
//       *qbytearray.QByteArray, qchar.QChar, qstringview.QStringView,
//       qlatin1stringview.QLatin1StringView, QUtf8StringView.
//
// Lifetime rules (same as QByteArrayView):
//   - Slicing methods return views of the same buffer; keepAlive is propagated.
//   - Do not mutate the source slice/string/*QByteArray while a view is alive.
//
// TODO (not yet bound):
//   - arg()                 — variadic template.
//   - reverse_iterator methods (rbegin/rend/crbegin/crend).
//   - operator std::string_view / std::u8string_view conversions.
//   - QBasicUtf8StringView<true> (char8_t-based variant) — would need a
//     separate binding if 8-bit-char interop is required.
//   - fromArray()           — trivial; use FromBytes instead.
//   - Container constructors — use FromBytes / FromString / FromQByteArray.
//   - QString overloads of compare/equal (side) — would create an import
//     cycle with qstringview; use qstring.QString → .Utf8() → FromQByteArray,
//     or convert to QStringView / QLatin1StringView first.
//   - Reverse iterators and STL-style algorithms — use Bytes() / Go slices.
//
// This package is a subpackage of github.com/vadimlxx/go-qt6core.
// Requires Qt >= 6.11 and a C++17 compiler; links via pkg-config Qt6Core.
package qtext
