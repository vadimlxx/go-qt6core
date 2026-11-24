// Package qdatetime provides the QCollatorSortKey binding: precomputed keys
// returned by QCollator::sortKey().
//
// Comparing sort keys is cheaper than repeatedly comparing strings with the
// collator. QCollatorSortKey is a heap-allocated Qt class; in Go it is a handle
// that owns the C++ object. The object is released by Close() or by
// runtime.AddCleanup when the wrapper becomes unreachable. Close is optional,
// idempotent and nil-safe.
//
// # Type mapping
//
//	Qt                                  Go
//	QCollatorSortKey                    *qdatetime.QCollatorSortKey
//
// # Construction
//
// QCollatorSortKey has no default constructor in Qt 6.11
// (QCollatorSortKey() = delete); the only public constructor is the copy
// constructor:
//
//	NewQCollatorSortKey(*QCollatorSortKey)  -> QCollatorSortKey(const QCollatorSortKey &)
//	NewQCollatorSortKeyCopy(*QCollatorSortKey)
//
// NewQCollatorSortKey is the generic constructor that dispatches on the
// argument type, following the NewQString pattern. Calling it without an
// argument panics, because Qt has no default-constructed sort key.
//
// Keys are normally obtained from (*QCollator).SortKey / SortKeyString or from
// the package-level DefaultSortKey / DefaultSortKeyString helpers.
//
// # API summary
//
// Lifecycle:
//
//	Close, Clone, Swap
//
// Comparison:
//
//	Compare(*QCollatorSortKey) int
//	Less(*QCollatorSortKey) bool
//
// # Not implemented yet
//
//   - The protected constructor QCollatorSortKey(QCollatorSortKeyPrivate *)
//     takes a private implementation type and cannot be exposed to Go.
//
//   - QCollatorSortKey::operator= (copy assignment): use Clone() or Swap().
//
//   - Move constructor / move assignment: not needed for GC-managed handles.
//
//   - Structured bindings / std::tuple comparison protocol: no Go equivalent;
//     use Compare() and Less().
//
//   - QDataStream / qHash / QDebug integration: not bound.
//
// # C bridge
//
// QCollator and QCollatorSortKey share qcollator_c.h / qcollator_c.cpp,
// mirroring QtCore's qcollator.h/.cpp. QCollatorSortKey is passed to C as an
// opaque heap pointer, so no handle casts are needed.
//
// # Requirements
//
// Qt >= 6.11, C++17, pkg-config Qt6Core. The version check is performed at
// compile time in qcollator_c.cpp via internal/qt6core/check_qt_version.h.
//
// This package is a subpackage of github.com/vadimlxx/go-qt6core.
package qdatetime
