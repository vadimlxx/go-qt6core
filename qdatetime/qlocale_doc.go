// Package qdatetime provides Go bindings for Qt 6.11 QtCore date/time-related types.
//
// # Currently supported
//
//   - QLocale (handle-based; COW via QSharedDataPointer)
//
// # QLocale: handle semantics
//
// QLocale is a value type in Qt that uses Copy-on-Write (COW) via
// QSharedDataPointer<QLocalePrivate>. In Go it is exposed as a HANDLE: a Go struct holding a
// C++ pointer (*QLocale). All accessors go directly to the C++ object.
//
// Value semantics vs reference semantics:
//
//   - QLocale is a reference type (handle). Copying the Go pointer does NOT create a new
//     locale; both variables refer to the same C++ object. Use Clone() to obtain an
//     independent copy.
//   - The Go zero value QLocale{} is invalid (nil pointer). Always construct via NewQLocale /
//     NewQLocaleFromName / NewQLocaleFromLanguageCountry / CLocale / SystemLocale.
//
// Memory management:
//
//   - NewQLocale, NewQLocaleFromName, NewQLocaleFromLanguageCountry, Clone, CLocale and
//     SystemLocale allocate a NEW heap-allocated C++ object.
//   - Close() releases the C++ object immediately and disables the finalizer. If Close() is
//     not called, a finalizer will free the object at some point during GC.
//   - Unlike QString or QByteArray, QLocale does NOT have a "FromUnsafePtr" constructor in
//     the current binding — internal code that needs one should add it explicitly if the
//     need arises.
//
// # Dependencies
//
//   - qtext.QString — all string-related methods return *qtext.QString.
//   - QDate / QTime / QDateTime — future classes in this same package.
//   - qtext.QStringList — for uiLanguages() / createSeparatedList() (not yet implemented in
//     qtext as of this writing).
//
// # Conventions
//
// Getters (values):
//
//   - Language() Language
//   - Script() Script
//   - Country() Country
//   - MeasurementSystem() MeasurementSystem
//   - TextDirection() int (cast to Qt::LayoutDirection in C++)
//
// String getters (each returns a fresh *qtext.QString; caller must Close() it):
//
//   - Name()                — QLocale::name()               (e.g. "en_US")
//   - Bcp47Name()           — QLocale::bcp47Name()          (e.g. "en-US")
//   - NativeLanguageName()  — QLocale::nativeLanguageName()
//   - NativeTerritoryName() — QLocale::nativeTerritoryName()
//   - DecimalPoint()        — QLocale::decimalPoint()
//   - GroupSeparator()      — QLocale::groupSeparator()
//   - Percent()             — QLocale::percent()
//   - ZeroDigit()           — QLocale::zeroDigit()
//   - NegativeSign()        — QLocale::negativeSign()
//   - PositiveSign()        — QLocale::positiveSign()
//   - Exponential()         — QLocale::exponential()
//   - AmText()              — QLocale::amText()
//   - PmText()              — QLocale::pmText()
//
// Number formatting (each returns a fresh *qtext.QString):
//
//   - ToStringInt64(v int64)                          — QLocale::toString(qlonglong)
//   - ToStringUInt64(v uint64)                        — QLocale::toString(qulonglong)
//   - ToStringDouble(v float64, format byte, precision int)
//                                                     — QLocale::toString(double, char, int)
//
// Number parsing (returns (value, ok bool)):
//
//   - ToInt64(s string) (int64, bool)   — QLocale::toLongLong(QString, bool*)
//   - ToDouble(s string) (float64, bool) — QLocale::toDouble(QString, bool*)
//
// Case conversion (returns fresh *qtext.QString):
//
//   - ToUpper(s string) *qtext.QString — QLocale::toUpper(QString)
//   - ToLower(s string) *qtext.QString — QLocale::toLower(QString)
//
// Month and day names (returns fresh *qtext.QString):
//
//   - MonthName(month int, format FormatType)           — QLocale::monthName(int, FormatType)
//   - StandaloneMonthName(month int, format FormatType) — QLocale::standaloneMonthName
//   - DayName(day int, format FormatType)               — QLocale::dayName
//   - StandaloneDayName(day int, format FormatType)     — QLocale::standaloneDayName
//
// Constructors:
//
//   - NewQLocale()                                    — QLocale()
//   - NewQLocaleFromName(name string)                 — QLocale(QString)
//   - NewQLocaleFromLanguageCountry(lang, country)    — QLocale(Language, Country)
//   - Clone() *QLocale                                — copy constructor (independent object)
//
// Static methods:
//
//   - SetDefaultLocale(locale *QLocale) — QLocale::setDefault(QLocale)
//   - CLocale() *QLocale                — QLocale::c()
//   - SystemLocale() *QLocale           — QLocale::system()
//
// # Notes on specific methods
//
// Name and Bcp47Name — Qt 6.7 introduced an overload that takes a TagSeparator argument:
//
//	QString name(QLocale::TagSeparator separator = TagSeparator::Underscore) const;
//	QString bcp47Name(QLocale::TagSeparator separator = TagSeparator::Dash) const;
//
// The current Go binding exposes only the no-argument overloads (default separators). The
// separator-aware overloads are listed under "Not yet implemented" below.
//
// ToStringDouble — format is a single ASCII character ('g', 'f', 'e', …); precision is the
// number of fractional digits, matching QString::number(double, char, int) semantics.
//
// # Enum types
//
// The following QLocale enums are exposed as Go types, each in its own file:
//
//   - language.go          — type Language int        (Language enum)
//   - script.go            — type Script int          (Script enum)
//   - country.go           — type Country int         (Country/Territory enum)
//   - formattype.go        — type FormatType int      (FormatType enum)
//   - measurement_system.go— type MeasurementSystem int
//   - number_option.go     — type NumberOption int + type NumberOptions (bitmask)
//   - quotation_style.go   — type QuotationStyle int
//   - currency_format.go   — type CurrencySymbolFormat int
//   - data_size_format.go  — type DataSizeFormat int + type DataSizeFormats (bitmask)
//   - tag_separator.go     — type TagSeparator int
//   - language_code_type.go— type LanguageCodeType int + type LanguageCodeTypes (bitmask)
//
// All enums are declared as `int`-based types. This keeps them trivially convertible from
// C.int and matches the existing enum conventions in qgeometry (FormatType-like types).
//
// # Not yet implemented
//
// Methods that depend on classes not yet bound:
//
//   - QDate / QTime / QDateTime — will be added together with these classes in this package:
//       QString toString(QDate date, const QString& format)
//       QString toString(QTime time, const QString& format)
//       QString toString(const QDateTime& dt, const QString& format)
//       QString toString(QDate date, FormatType format)
//       QString toString(QTime time, FormatType format)
//       QString toString(const QDateTime& dt, FormatType format)
//       QString dateFormat(FormatType format)
//       QString timeFormat(FormatType format)
//       QString dateTimeFormat(FormatType format)
//       QDate   toDate(const QString& string, ...)
//       QTime   toTime(const QString& string, ...)
//       QDateTime toDateTime(const QString& string, ...)
//
//   - QCalendar — will be added together with the QCalendar class in this package:
//       All QCalendar-aware overloads of toString / toDate / toDateTime.
//
// Container-dependent methods (need QList<...> bindings):
//
//   - QList<Qt::DayOfWeek> weekdays() const
//   - QList<QLocale> matchingLocales(Language, Script, Territory)
//   - QList<Country> countriesForLanguage(Language)  (deprecated in Qt 6.6)
//
// QStringList-dependent methods (need QStringList in qtext):
//
//   - QStringList uiLanguages(TagSeparator separator) const
//   - QString createSeparatedList(const QStringList& strl) const
//
// Stream operators (need qdatastream / qdebug):
//
//   - QDataStream& operator<<(QDataStream&, const QLocale&)
//   - QDataStream& operator>>(QDataStream&, QLocale&)
//   - QDebug operator<<(QDebug, const QLocale&)
//
// Static conversion methods (low priority — code/string ↔ enum tables):
//
//   - QString languageToCode(Language language, LanguageCodeTypes codeTypes)
//   - Language codeToLanguage(QStringView languageCode, LanguageCodeTypes codeTypes)
//   - QString territoryToCode(Territory territory)
//   - Territory codeToTerritory(QStringView territoryCode)
//   - QString scriptToCode(Script script)
//   - Script codeToScript(QStringView scriptCode)
//   - QString languageToString(Language language)
//   - QString territoryToString(Territory territory)
//   - QString scriptToString(Script script)
//   - QString countryToCode(Country) / Country codeToCountry(QStringView)  (deprecated Qt 6.6)
//   - QString countryToString(Country)                                    (deprecated Qt 6.6)
//
// Currency and data-size formatting (need extra overloads and/or flags integration):
//
//   - QString currencySymbol(CurrencySymbolFormat format = CurrencySymbol)
//   - QString toCurrencyString(qlonglong value, const QString& symbol)
//   - QString toCurrencyString(qulonglong value, const QString& symbol)
//   - QString toCurrencyString(double value, const QString& symbol, int precision)
//   - QString formattedDataSize(qint64 bytes, int precision, DataSizeFormats format)
//
// Quotation helpers:
//
//   - QString quoteString(const QString& str, QuotationStyle style)
//   - QString quoteString(QStringView str, QuotationStyle style)
//
// Number options (accessors not yet wired up):
//
//   - void setNumberOptions(NumberOptions options)
//   - NumberOptions numberOptions() const
//
// Collation:
//
//   - QLocale collation() const
//
// # Flags
//
// NumberOptions, DataSizeFormats, LanguageCodeTypes are declared as int-based bitmask types,
// but they are NOT yet integrated into the method signatures that take them in Qt
// (setNumberOptions, numberOptions, formattedDataSize, languageToCode/codeToLanguage, …).
// They are ready for future use — no signature changes will be needed when those methods
// are added.
//
// # Requirements
//
// Qt >= 6.11, C++17, pkg-config Qt6Core.
package qdatetime
