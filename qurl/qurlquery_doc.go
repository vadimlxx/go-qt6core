// Package qurl provides Go bindings for Qt 6.11 QtCore QUrlQuery.
//
// QUrlQuery represents the query part of a URL as a list of key/value pairs
// and provides parsing, serialization, and item manipulation.
//
// In Go, each *QUrlQuery is a handle to a heap-allocated C++ QUrlQuery; call
// Close() for immediate release, otherwise a finalizer handles cleanup.
//
// # Construction
//
//   - NewQUrlQuery()                         → QUrlQuery()
//   - QUrlQueryFromUrl(url)                  → QUrlQuery(const QUrl&)
//   - QUrlQueryFromString(s)                 → QUrlQuery(const QString&) [fromUtf8]
//   - QUrlQueryFromQString(qs)               → QUrlQuery(const QString&)
//   - Clone()                                → copy constructor
//   - Swap(other)                            → QUrlQuery::swap
//
// # State
//
//   - IsEmpty() / IsDetached() / Clear()     → isEmpty() / isDetached() / clear()
//
// # Query string
//
//   - QueryString(enc) string                → query(ComponentFormattingOptions) → QString
//   - Query(enc) *qtext.QString              → query(ComponentFormattingOptions)
//   - ToString(enc) string                   → toString(ComponentFormattingOptions) → QString
//   - ToQString(enc) *qtext.QString          → toString(ComponentFormattingOptions)
//   - SetQueryQString(qs)                    → setQuery(const QString&)
//   - SetQueryString(s)                      → setQuery(const QString&) [fromUtf8]
//
// # Delimiters
//
//   - SetQueryDelimiters(value, pair)        → setQueryDelimiters(QChar, QChar)
//   - QueryValueDelimiter() QChar            → queryValueDelimiter()
//   - QueryPairDelimiter() QChar             → queryPairDelimiter()
//   - QueryValueDelimiterDefault             → defaultQueryValueDelimiter() (u'=')
//   - QueryPairDelimiterDefault              → defaultQueryPairDelimiter() (u'&')
//
// # Item list (map)
//
//   - SetQueryItems(map[string]string)       → setQueryItems(const QList<std::pair<QString,QString>>&)
//   - QueryItems(enc) map[string]string      → queryItems(ComponentFormattingOptions)
//
// Note: the order of items is not preserved in map form; when duplicate keys
// are present, map construction keeps the last value encountered. For
// order-sensitive access, iterate the query string or use AllQueryItemValues.
//
// # Single-item operations
//
//   - HasQueryItemQString(qs)                → hasQueryItem(const QString&)
//   - HasQueryItem(s)                        → hasQueryItem(const QString&) [fromUtf8]
//   - AddQueryItemQString(k, v)              → addQueryItem(const QString&, const QString&)
//   - AddQueryItem(k, v)                     → addQueryItem(const QString&, const QString&) [fromUtf8]
//   - RemoveQueryItemQString(k)              → removeQueryItem(const QString&)
//   - RemoveQueryItem(k)                     → removeQueryItem(const QString&) [fromUtf8]
//   - QueryItemValueQString(k, enc) string   → queryItemValue(const QString&, ComponentFormattingOptions)
//   - QueryItemValue(k, enc) string          → queryItemValue(const QString&, ComponentFormattingOptions) [fromUtf8]
//   - AllQueryItemValuesQString(k, enc) []string → allQueryItemValues(const QString&, ComponentFormattingOptions)
//   - AllQueryItemValues(k, enc) []string    → allQueryItemValues(const QString&, ComponentFormattingOptions) [fromUtf8]
//   - RemoveAllQueryItemsQString(k)          → removeAllQueryItems(const QString&)
//   - RemoveAllQueryItems(k)                 → removeAllQueryItems(const QString&) [fromUtf8]
//
// # Not implemented
//
//   - QUrlQuery(std::initializer_list<std::pair<QString,QString>>) — Go
//     map constructor covers this via SetQueryItems.
//   - operator== / operator!= (Q_DECLARE_EQUALITY_COMPARABLE) — compare
//     via QueryString or QueryItems.
//   - qHash(QUrlQuery) — Go can hash the serialized query string.
//   - QUrl::setQuery(const QUrlQuery&) — not yet bound; use
//     QUrl.SetQueryString(q.QueryString(...)).
//   - QUrl::query() returning QUrlQuery — not yet bound.
//
// This file is part of package qurl of github.com/vadimlxx/go-qt6core.
// Requires Qt >= 6.11 and a C++17 compiler; links via pkg-config Qt6Core.
package qurl