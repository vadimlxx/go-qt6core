// QLocale is a Go binding for Qt 6.11 QtCore's QLocale.
//
// QLocale encapsulates a language/territory/script combination along with
// a large set of locale-specific formatting rules: number formatting,
// date/time formatting, currency, quotation, and month/day names.
//
// # Type semantics
//
// QLocale is a handle to a heap-allocated C++ QLocale. Each *QLocale owns
// the underlying C++ object. Copying the Go pointer does NOT create a new
// locale; use Clone() for an independent copy.
//
// # Ownership
//
// The C++ QLocale is released by one of:
//   - Close(): eager, deterministic release.
//   - runtime.AddCleanup: safety net when *QLocale becomes unreachable.
//
// Close is optional, idempotent, and nil-safe. Using a QLocale after Close
// panics with "qlocale: use after Close (or nil receiver)".
//
// # Construction
//
//   - NewQLocale()                                         — QLocale()
//   - NewQLocaleFromName(name string)                      — QLocale(QString)
//   - NewQLocaleFromLanguageTerritory(lang, territory)     — QLocale(Language, Territory)
//   - NewQLocaleFromLanguageScriptTerritory(l, s, t)       — QLocale(Language, Script, Territory)
//   - (q).Clone()                                          — QLocale(const QLocale&)
//   - CLocale() / SystemLocale()                           — QLocale::c() / QLocale::system()
//
// # Value accessors
//
//   - Language() Language
//   - Script() Script
//   - Territory() Territory
//   - MeasurementSystem() MeasurementSystem
//   - TextDirection() int
//   - FirstDayOfWeek() int
//   - NumberOptions() NumberOptions
//   - SetNumberOptions(opts)
//
// # Number parsing
//
//   - ToShort/ToUShort/ToInt/ToUInt/ToLong/ToULong/ToLongLong/
//     ToULongLong/ToFloat/ToDouble(s string) (T, bool)
//     Each maps to the corresponding QLocale::toXxx(QStringView, bool*) const.
//
// # Number formatting
//
//   - ToStringLL(v int64) *qtext.QString
//   - ToStringULL(v uint64) *qtext.QString
//   - ToStringLLW(v, fieldWidth, fillChar rune) *qtext.QString
//   - ToStringULLW(v, fieldWidth, fillChar rune) *qtext.QString
//   - ToStringDouble(v float64, format byte, precision int) *qtext.QString
//   - ToStringFloat(v float32, format byte, precision int) *qtext.QString
//     Each XxxString() variant returns a Go string (closes the intermediate QString).
//
// # Locale-specific strings
//
//   - DecimalPoint / GroupSeparator / Percent / ZeroDigit
//   - NegativeSign / PositiveSign / Exponential
//   - AmText / PmText
//   - MonthName / StandaloneMonthName / DayName / StandaloneDayName
//     Each returns *qtext.QString; XxxString() variants return Go strings.
//
// # Case conversion
//
//   - ToUpper(s string) / ToLower(s string) *qtext.QString
//   - XxxString variants
//     Matches QLocale::toUpper/toLower(const QString&) const.
//
// # Currency
//
//   - CurrencySymbol(format CurrencySymbolFormat) *qtext.QString
//   - ToCurrencyStringLL(v int64, symbol *qtext.QString) *qtext.QString
//   - ToCurrencyStringULL(v uint64, symbol *qtext.QString) *qtext.QString
//   - ToCurrencyStringDouble(v float64, symbol *qtext.QString, precision int)
//     Each maps to QLocale::currencySymbol / toCurrencyString.
//
// # Data size
//
//   - FormattedDataSize(bytes int64, precision int, format DataSizeFormats) *qtext.QString
//     Matches QLocale::formattedDataSize(qint64, int, DataSizeFormats).
//
// # Quotation
//
//   - QuoteString(s string, style QuotationStyle) *qtext.QString
//     Matches QLocale::quoteString(QStringView, QuotationStyle) const.
//
// # Separated list
//
//   - CreateSeparatedList(strings []string) *qtext.QString
//   - CreateSeparatedListQStrings(qs []*qtext.QString) *qtext.QString
//     Matches QLocale::createSeparatedList(const QStringList&).
//
// # Collation
//
//   - Collation() *QLocale
//     Matches QLocale::collation() const.
//
// # UI languages
//
//   - UILanguages(sep TagSeparator) []*qtext.QString
//   - UILanguagesStrings(sep TagSeparator) []string
//     Matches QLocale::uiLanguages(TagSeparator).
//
// # Date/time formatting
//
//   - ToStringDateFormat(d QDate, format FormatType) *qtext.QString
//   - ToStringDateString(d QDate, format string) *qtext.QString
//   - ToStringDateFormatCal(d QDate, format FormatType, cal *QCalendar) *qtext.QString
//   - ToStringDateStringCal(d QDate, format string, cal *QCalendar) *qtext.QString
//   - ToStringTimeFormat(t QTime, format FormatType) *qtext.QString
//   - ToStringTimeString(t QTime, format string) *qtext.QString
//   - ToStringDateTimeFormat(dt *QDateTime, format FormatType) *qtext.QString
//   - ToStringDateTimeString(dt *QDateTime, format string) *qtext.QString
//   - ToStringDateTimeFormatCal(dt *QDateTime, format FormatType, cal *QCalendar) *qtext.QString
//   - ToStringDateTimeStringCal(dt *QDateTime, format string, cal *QCalendar) *qtext.QString
//     Each maps to the corresponding QLocale::toString overload. XxxString()
//     variants return Go strings and close the intermediate QString.
//
// # Date/time patterns
//
//   - DateFormat(format FormatType) *qtext.QString
//   - TimeFormat(format FormatType) *qtext.QString
//   - DateTimeFormat(format FormatType) *qtext.QString
//     Matches QLocale::dateFormat / timeFormat / dateTimeFormat.
//
// # Parsing dates and times
//
//   - ToDateFormat(s string, format FormatType, baseYear int) QDate
//
//   - ToDateString(s string, format string, baseYear int) QDate
//
//   - ToDateFormatCal(s string, format FormatType, cal *QCalendar, baseYear int) QDate
//
//   - ToDateStringCal(s string, format string, cal *QCalendar, baseYear int) QDate
//
//   - ToTimeFormat(s string, format FormatType) QTime
//
//   - ToTimeString(s string, format string) QTime
//
//   - ToDateTimeFormat(s string, format FormatType, baseYear int) *QDateTime
//
//   - ToDateTimeString(s string, format string, baseYear int) *QDateTime
//
//   - ToDateTimeFormatCal(s string, format FormatType, cal *QCalendar, baseYear int) *QDateTime
//
//   - ToDateTimeStringCal(s string, format string, cal *QCalendar, baseYear int) *QDateTime
//     These map to QLocale::toDate / toTime / toDateTime overloads.
//
//     IMPORTANT: Qt's toDate/toDateTime return a NULL QDate/QDateTime when
//     parsing fails (no `bool *ok` out-parameter). Check IsValid() on the
//     result:
//
//     d := loc.ToDateFormat("not a date", LongFormat, DefaultTwoDigitBaseYear)
//     if !d.IsValid() { /* parse failed */ }
//
//     For ToTime*, the returned QTime is also null on failure (IsValid() == false).
//
// # baseYear
//
// DefaultTwoDigitBaseYear = 1900 mirrors QLocale::DefaultTwoDigitBaseYear.
// It affects only how two-digit years in the input are interpreted; pass
// DefaultTwoDigitBaseYear for Qt's default behaviour.
//
// # Date/Time/QCalendar interop
//
// QDate and QTime are value types (their zero values are valid: QDate{}
// is 1 January 4713 BCE, QTime{} is midnight). QDateTime and QCalendar are
// handle types; pass *QDateTime / *QCalendar. All four are defined in this
// same qdatetime package.
//
// # Enum types
//
// QLocale enums are exposed as Go types in separate files:
//   - language.go          — type Language int
//   - script.go            — type Script int
//   - country.go           — type Territory int
//   - formattype.go        — type FormatType int
//   - measurement_system.go— type MeasurementSystem int
//   - number_option.go     — type NumberOption int + type NumberOptions (bitmask)
//   - quotation_style.go   — type QuotationStyle int
//   - currency_format.go   — type CurrencySymbolFormat int
//   - data_size_format.go  — type DataSizeFormat int + type DataSizeFormats (bitmask)
//   - tag_separator.go     — type TagSeparator int
//   - language_code_type.go— type LanguageCodeType int + type LanguageCodeTypes (bitmask)
//
// # Not yet implemented
//
//   - QDataStream / QDebug operators.
//   - qHash(QLocale).
//   - Deprecated Qt 6.6 members (country(), countryToString(), countryToCode(),
//     codeToCountry(), countriesForLanguage()) are intentionally not exposed;
//     use Territory(), TerritoryToCode(), CodeToTerritory() and
//     MatchingLocales() instead.
//   - Deprecated Qt 6.7 members (name() / bcp47Name() without a TagSeparator)
//     are not exposed; use Name(sep) / Bcp47Name(sep).
//
// # Requirements
//
// Qt >= 6.11, C++17, pkg-config Qt6Core.
//
// This package is a subpackage of github.com/vadimlxx/go-qt6core.
package qdatetime
