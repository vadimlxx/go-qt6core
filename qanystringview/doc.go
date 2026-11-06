// Package qanystringview provides a Go binding for Qt 6.11 QtCore's QAnyStringView.
//
// QAnyStringView is a type-erased read-only string view over one of three encodings:
//
//	KindUtf8   — UTF-8 bytes        (matches QUtf8StringView)
//	KindLatin1 — Latin-1 bytes      (matches QLatin1StringView)
//	KindUtf16  — UTF-16 code units  (matches QStringView)
//
// Qt stores the tag in the low bits of an internal size_t field; the Go wrapper stores
// (data, size, kind, keepAlive) explicitly. All views are zero-copy and carry a keepAlive
// reference to the underlying buffer.
//
// # Construction
//
//	FromBytes([]byte)                     → KindUtf8
//	FromString(string)                    → KindUtf8
//	FromQUtf8StringView(QUtf8StringView)   → KindUtf8
//	FromQLatin1StringView(QLatin1StringView) → KindLatin1
//	FromQStringView(QStringView)          → KindUtf16
//	FromQByteArray(*QByteArray)           → KindUtf8
//	NewUnsafe(data, size, kind)           → advanced zero-copy
//
// # Interop with sibling packages
//
//	qanystringview → qstringview, qlatin1stringview, qutf8stringview, qstring, qbytearray, qchar
//
// No import cycles: none of the imported packages import qanystringview. If qstring ever
// needs to construct a QAnyStringView, use the unsafe.Pointer bridge pattern (as with
// qchar.Decompose / qbytearray.FromQStringUnsafePtr).
//
// # Bound API (summary)
//
// Construction:
//   FromBytes, FromString, FromQByteArray,
//   FromQStringView, FromQLatin1StringView, FromQUtf8StringView,
//   NewUnsafe, Empty
//
// Introspection:
//   Kind, Size, Len, SizeBytes, IsNull, IsEmpty, IsUtf8, IsLatin1, IsUtf16
//
// Element access:
//   Front, Back (return qchar.QChar)
//
// Slicing (return new views sharing the same buffer; Kind is preserved except Mid which may
// return a null view):
//   Left, Right, Mid, First, Last, Sliced, SlicedN, Chopped, Truncate, Chop, Slice, SliceN
//
// Dispatch to concrete type (matches QAnyStringView::visit):
//   AsQStringView, AsQLatin1StringView, AsQUtf8StringView  — return (view, ok)
//
// Comparison:
//   Compare, Equal (both static-like; accept two QAnyStringView)
//
// Conversion:
//   ToString() *qstring.QString   (dispatches on Kind; equivalent to QAnyStringView::toString())
//
// # TODO (not yet bound, or exposed differently vs Qt 6.11)
//
//   - visit() — Go has no idiomatic equivalent of the C++ visitor pattern. Use Kind() +
//     As*() instead.
//
//   - arg() — variadic template QString::arg; will appear with qstring.Arg.
//
//   - Single-character constructors QAnyStringView(Char), QAnyStringView(Char, QCharContainer&&),
//     QAnyStringView(Char32, Container&&) — use FromString(string(c)).
//
//   - fromArray() — use FromBytes / FromString.
//
//   - Implicit QString/QByteArray constructors with capacity-hint
//     (QAnyStringView(Container&&, wrapped_t<Container, QString>&&)) — use FromQString /
//     FromQByteArray directly.
//
//   - Comparison operators (==, !=, <, <=, >, >=) — use Compare / Equal.
//
//   - Implicit conversions operator QString / operator QByteArray — use ToString / As*.
//
//   - operator std::string_view / operator std::u8string_view — use Bytes().
//
//   - Iterators (begin/end/rbegin/rend/cbegin/cend/crbegin/crend) — use the concrete
//     As*() view types, which expose their own iterator-like helpers, or convert to []byte.
//
//   - front()/back() as constexpr QChar — exposed as Front()/Back() (runtime).
//
//   - charSize() / max_size() — private / rarely useful; not exposed.
//
//   - detects_US_ASCII_at_compile_time — a C++ constant-evaluation optimization; irrelevant
//     in Go.
//
//   - QT_NO_DEBUG_STREAM operator<<(QDebug) — not applicable.
//
//   - QCharContainer — internal helper type; not exposed.
//
//   - Comparison against QChar / QStringView / QLatin1StringView / QUtf8StringView directly
//     (operator overloads). Use Compare with the corresponding From* wrapper.
//
//   - QAnyStringView from other compatible containers (e.g. std::u16string, std::string_view,
//     std::vector<char>) — covered by FromBytes / FromString / FromQStringView etc.
//
// # Design notes
//
//   - QAnyStringView is Q_PRIMITIVE_TYPE. No heap allocation, no finalizer, no Close().
//
//   - Slicing preserves Kind (matching Qt: sliced/first/last/chopped/left/right all copy the
//     tag through setSize()). Mid() may return a null view (tag = KindUtf8) when the requested
//     range is out of bounds; that's what QContainerImplHelper::mid() returns for Null.
//
//   - Truncate/Chop/Slice/SliceN mutate the receiver in place (Qt: non-const methods).
//
//   - ToString() dispatches to QString::fromUtf8 / fromLatin1 / fromUtf16 based on Kind
//     (internally via QAnyStringView::toString / visit); the result is a fresh
//     *qstring.QString owned by the caller.
//
//   - Comparison (Compare) returns Qt's int result (negative/zero/positive); Equal returns
//     bool. Both are always case-sensitive per QAnyStringView::equal semantics; Compare takes
//     a caseSensitive flag matching Qt::CaseSensitivity.
//
// # Requirements
//
// Qt >= 6.11, C++17, pkg-config Qt6Core. Version check is performed at compile time in
// qanystringview_c.cpp via internal/qt6core/check_qt_version.h.
//
// This package is a subpackage of github.com/vadimlxx/go-qt6core.
package qanystringview
