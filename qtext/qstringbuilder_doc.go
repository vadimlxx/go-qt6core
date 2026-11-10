// Package qstringbuilder provides a Go binding for Qt 6.11 QtCore's QStringBuilder.
//
// QStringBuilder<A, B> is a C++ template used for efficient concatenation of string-like
// operands. In Qt, `a % b` (or `a + b` with QT_USE_QSTRINGBUILDER) creates a lightweight
// QStringBuilder that lazily builds the final QString or QByteArray when converted. Because
// QStringBuilder is a template (its concrete type is QStringBuilder<A, B>, unique per operand
// pair), there is no C++ type that a cgo handle could point to.
//
// Therefore this package provides a Go-native equivalent: QstringBuilder holds a sequence of
// Parts and materializes them into *qstring.QString or *qbytearray.QByteArray on demand.
// Semantics match Qt's QStringBuilder:
//   - Each operand's "flavor" (QString vs QByteArray) is determined by its type, mirroring
//     QConcatenable<T>::ConvertTo in qstringbuilder.h.
//   - The builder's overall flavor is QString if ANY part is QString-flavored; QByteArray
//     otherwise (matching Qt's ConvertToTypeHelper: QString always wins).
//   - ToQString() / ToQByteArray() materialize the result in one C call.
//
// # Supported operand types (mapping to QConcatenable<T> specializations)
//
//	Qt type                 Go type                              Flavor
//	char                    byte                                 QByteArray
//	const char* / char[N]   string                               QByteArray
//	QByteArray              *qbytearray.QByteArray               QByteArray
//	QByteArrayView          qbytearrayview.QByteArrayView        QByteArray
//	char16_t                uint16                               QString
//	const char16_t*         []uint16                             QString
//	QChar                   qchar.QChar                          QString
//	QLatin1Char             qlatin1char.QLatin1Char              QString
//	QLatin1StringView       qlatin1stringview.QLatin1StringView  QString
//	QString                 *qstring.QString                     QString
//	QStringView             qstringview.QStringView              QString
//	QStringBuilder<A,B>     *QstringBuilder                      aggregated
//
// QUtf8StringView and QAnyStringView are NOT supported by QStringBuilder in Qt 6.11 (no
// QConcatenable<T> specialization in qstringbuilder.h); passing them will panic. To include
// their contents, convert first: e.g. qstringbuilder.FromString(qutf8.ToBytes()).
//
// # Construction
//
//	Concat(parts ...any) *QstringBuilder
//	FromString, FromBytes, FromByte,
//	FromQChar, FromQLatin1Char, FromUint16, FromUtf16,
//	FromQString, FromQStringView,
//	FromQLatin1StringView,
//	FromQByteArray, FromQByteArrayView
//
// # Queries
//
//	Len() int                     — number of parts
//	IsQStringFlavored() bool      — true if the result materializes as QString
//
// # Concatenation
//
//	(*QstringBuilder).Append(other *QstringBuilder) *QstringBuilder
//	(*QstringBuilder).AppendPart(p any) *QstringBuilder
//
// # Materialization
//
//	ToQString()       *qstring.QString
//	ToQByteArray()    *qbytearray.QByteArray
//	ToUtf8()          *qbytearray.QByteArray  — equivalent to ToQString().Utf8()
//	ToLatin1()        *qbytearray.QByteArray  — equivalent to ToQString().Latin1()
//	ToLocal8Bit()     *qbytearray.QByteArray  — equivalent to ToQString().Local8Bit()
//	AppendToQString(q *qstring.QString)            — equivalent to q += builder
//	AppendToQByteArray(ba *qbytearray.QByteArray)  — equivalent to ba += builder
//
// # Interop with sibling packages
//
//	qstringbuilder → qstring, qbytearray, qbytearrayview, qlatin1char, qlatin1stringview,
//	                 qstringview, qchar
//
// No import cycles: none of those packages import qstringbuilder. To pass a QString* or
// QByteArray* into the C-ABI, qstring.UnsafePtr() and qbytearray.UnsafePtr() are used; the
// C++ side reinterprets the opaque pointer as the real Qt type. Parts holding views keep the
// source alive via runtime.KeepAlive on the builder after each C call.
//
// # TODO (not yet bound, or exposed differently vs Qt 6.11)
//
//   - operator% (a % b): Go has no user-defined operators. Use Concat or Append.
//
//   - operator+ (with QT_USE_QSTRINGBUILDER): same; use Concat / AppendPart.
//
//   - QConcatenable<T> template machinery: internal; not exposed.
//
//   - ExactSize optimization: internal; the Go side always resizes correctly after materialization.
//
//   - QStringBuilder<A,B>::size() and ::isNull(): not exposed (semantics are "approximate" in Qt —
//     upper bound rather than exact), use ToQString().Size() / ToQByteArray().Size() instead.
//
//   - QStringBuilderCommon::toUpper() / toLower(): not exposed as builder methods (return type
//     depends on flavor). Use ToQString().ToUpper() / ToLower() or ToQByteArray().ToUpper() / ToLower().
//
//   - operator std::string_view / std::u16string_view: not applicable in Go.
//
//   - QDebug operator<<: not applicable.
//
//   - QByteArrayView-to-QString conversion is handled automatically (Qt interprets via UTF-8);
//     no separate method needed.
//
//   - QStringBuilder with more than two operands at compile time: handled by Append chaining or
//     a single Concat call with N parts.
//
//   - Lazy evaluation across multiple materializations: Go evaluates eagerly per ToQString() /
//     ToQByteArray() call, but in a single C call (equivalent efficiency to Qt's one-shot
//     conversion; multiple materializations do multiple C calls, matching Qt's behavior).
//
// # Requirements
//
// Qt >= 6.11, C++17, pkg-config Qt6Core. Version check is performed at compile time in
// qstringbuilder_c.cpp via internal/qt6core/check_qt_version.h.
//
// This package is a subpackage of github.com/vadimlxx/go-qt6core.
package qtext
