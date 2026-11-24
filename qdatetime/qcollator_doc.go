// Package qdatetime provides the QCollator binding: locale-aware comparison of
// UTF-16 strings.
//
// QCollator is a heap-allocated Qt class (pimpl). In Go it is represented by a
// handle that owns the C++ object; the object is released by Close() or by
// runtime.AddCleanup when the wrapper becomes unreachable. Close is optional,
// idempotent and nil-safe.
//
// # Type mapping
//
//	Qt                                  Go
//	QCollator                           *qdatetime.QCollator
//	QLocale                             *qdatetime.QLocale
//	QString                             *qtext.QString / string
//	QStringView                         qtext.QStringView / string
//	Qt::CaseSensitivity                 qdata.CaseSensitivity,
//	                                    qdata.CaseInsensitive,
//	                                    qdata.CaseSensitive
//
// Qt::CaseSensitivity is reused from package qdata; this package does not
// declare a second copy. QCollator itself has no enums, so no extra enum file
// is needed.
//
// # Constructors
//
//	NewQCollator()                    -> QCollator()
//	NewQCollator(*QLocale)            -> QCollator(const QLocale &)
//	NewQCollator(*QCollator)          -> QCollator(const QCollator &)
//	NewQCollatorFromLocale(*QLocale)  -> QCollator(const QLocale &)
//	NewQCollatorCopy(*QCollator)      -> QCollator(const QCollator &)
//
// NewQCollator is the generic constructor that dispatches on the argument
// types, following the NewQString pattern.
//
// # API summary
//
// Lifecycle:
//
//	Close, Clone, Swap
//
// Locale and options:
//
//	SetLocale, Locale,
//	CaseSensitivity, SetCaseSensitivity,
//	NumericMode, SetNumericMode,
//	IgnorePunctuation, SetIgnorePunctuation
//
// Comparison:
//
//	Compare(*QString, *QString), CompareString(string, string),
//	CompareView(QStringView, QStringView), CompareViewString(string, string),
//	Less, LessString, LessView, LessViewString
//
// Sort keys:
//
//	SortKey(*QString), SortKeyString(string)
//
// Static helpers:
//
//	DefaultCompare(QStringView, QStringView), DefaultCompareString(string, string),
//	DefaultSortKey(QStringView), DefaultSortKeyString(string)
//
// # Not implemented yet
//
//   - compare(const QChar *s1, qsizetype len1, const QChar *s2, qsizetype len2):
//     Go users can pass UTF-16 through qtext.QString / qtext.QStringView, so the
//     raw-pointer overload is intentionally left out.
//
//   - QCollator::operator= (copy assignment): use Clone() or Swap().
//
//   - Move constructor / move assignment: not needed for GC-managed handles.
//
//   - QCollator::detach() is private and is not exposed.
//
//   - QDataStream / qHash / QDebug integration: not bound.
//
// # C bridge
//
// QCollator and QCollatorSortKey share one C ABI (qcollator_c.h) and one
// implementation (qcollator_c.cpp), mirroring QtCore's qcollator.h/.cpp. The
// handle types are the Qt classes themselves, so no handle casts are required;
// the only conversions are UTF-16 (uint16_t* -> QStringView/QString) and
// int -> Qt::CaseSensitivity, both centralized in qcollator_c.cpp.
//
// # Requirements
//
// Qt >= 6.11, C++17, pkg-config Qt6Core. The version check is performed at
// compile time in qcollator_c.cpp via internal/qt6core/check_qt_version.h.
//
// This package is a subpackage of github.com/vadimlxx/go-qt6core.
package qdatetime
