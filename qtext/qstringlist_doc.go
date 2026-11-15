// Package qtext provides Go bindings for Qt 6.11 QtCore QStringList.
//
// QStringList is a QList<QString> with extra convenience methods (sort,
// join, filter, replaceInStrings, indexOf/lastIndexOf with case sensitivity
// or QRegularExpression). In Go, each *QStringList is a handle to a
// heap-allocated C++ QStringList; call Close() for immediate release,
// otherwise a finalizer handles cleanup.
//
// Construction:
//   - NewQStringList()                    → QStringList()
//   - QStringListFromQString(s)           → QStringList(const QString&)
//   - QStringListFromGoString(s)          → QStringList(const QString&) [fromUtf8]
//   - QStringListFromQStrings(items)      → QStringList(const QList<QString>&) [copy]
//   - QStringListFromGoStrings(items)     → QStringList(const QList<QString>&) [fromUtf8]
//   - Clone()                             → copy constructor
//   - Swap(other)                         → QStringList::swap
//
// Element access (each returns a new *QString; ownership is the caller's):
//   - At(i) / Front() / Back() / ToSlice()
//   - TakeAt(i) / TakeFirst() / TakeLast()
//
// Case sensitivity: passed as bool (true == Qt::CaseSensitive), matching the
// rest of the qtext package (QString, QStringView). The qdata.CaseSensitivity
// constants are the underlying Qt values; csInt() converts bool → int for the
// C ABI.
//
// Cross-package types used by this binding:
//   - *QString           (qtext)            — every element
//   - QChar              (qtext)            — JoinChar
//   - *QRegularExpression (qtext)           — FilterQRegularExpression,
//                                             ReplaceInStringsQRegularExpression,
//                                             IndexOfQRegularExpression,
//                                             LastIndexOfQRegularExpression
//
// # Not implemented
//
// The following QStringList APIs are not exposed yet:
//
//   - QListIterator<QString> / QMutableListIterator<QString> — Java-style
//     iterators; use ToSlice() and index-based access.
//   - STL-style iterators (begin / end / rbegin / rend) — not applicable
//     in Go; use ToSlice().
//   - filter(const QStringMatcher&) and filter(const QLatin1StringMatcher&) —
//     require QStringMatcher / QLatin1StringMatcher which are not yet bound.
//   - Overloads taking QStringView / QLatin1StringView — the current API
//     accepts *QString (or Go string as a convenience); the view overloads
//     can be added with suffixes when callers need zero-copy comparisons.
//   - operator+ (concatenation) and operator<< (chainable append) — replaced
//     by Append / AppendList.
//   - QDataStream operator<< / operator>> — not idiomatic in Go.
//   - QDebug operator<< — not idiomatic in Go.
//   - qHash(const QStringList&) and QHashSpecialization — not exposed;
//     Go users can hash via elements themselves.
//   - QList::resize, QList::fill, QList::mid, QList::sliced, QList::first,
//     QList::last (returning lists, not elements) — pending a QList<T> generic
//     binding; use ToSlice() and QStringListFromQStrings for now.
//
// # QtPrivate helpers
//
// QtPrivate::QStringList_sort, QStringList_join, QStringList_filter, etc. are
// internal free functions (declared in qstringlist.h under QtPrivate) that
// back the public member functions. They are not part of the public Qt API and
// are therefore not bound directly — the corresponding QStringList member
// functions are used instead.
//
// This file is part of package qtext of github.com/vadimlxx/go-qt6core.
// Requires Qt >= 6.11 and a C++17 compiler; links via pkg-config Qt6Core.
package qtext