// Package qdatetime provides Go bindings for Qt 6.11 QtCore date/time-related types.
//
// # Currently supported
//
//   - QLocale (handle-based; COW via QSharedDataPointer)
//
// # Planned additions (same package)
//
//   - QDate
//   - QTime
//   - QDateTime
//   - QCalendar
//
// Because these classes will live in the same package, QLocale methods that take or return
// QDate / QTime / QDateTime (toString(QDate), toDate(), toTime(), toDateTime(), and the
// QCalendar-aware overloads) can be added without import cycles or unsafe.Pointer bridges.
//
// # Handle semantics (QLocale)
//
// QLocale is exposed as a HANDLE (Go struct holding a C++ pointer). All accessors and
// mutators go directly to the heap-allocated C++ object. Use Clone() for an independent copy.
// Call Close() to release the C++ object immediately; otherwise a finalizer handles it.
//
// # Dependencies
//
//   - qtext.QString for all string-related methods.
//   - qtext.QStringList (future) for uiLanguages / createSeparatedList.
//
// # Not yet implemented
//
//   - All methods involving QDate / QTime / QDateTime — will be added together with
//     those classes in this same package.
//   - QCalendar-dependent methods — will be added together with QCalendar.
//   - QList<Qt::DayOfWeek> weekdays() — waiting for container bindings.
//   - QList<QLocale> matchingLocales() — waiting for container bindings.
//   - QDataStream / QDebug operators.
//   - Static methods: languageToCode, codeToLanguage, territoryToCode, codeToTerritory,
//     scriptToCode, codeToScript, languageToString, territoryToString, scriptToString.
//   - formattedDataSize, toCurrencyString, currencySymbol (low priority).
//   - quoteString, createSeparatedList (low priority).
//   - NumberOptions and DataSizeFormats flags are defined but not yet integrated
//     into method signatures.
//
// # Protected access
//
// QLocale has no protected members — only public API and private implementation.
//
// Requires Qt >= 6.11, C++17, pkg-config Qt6Core.
package qdatetime
