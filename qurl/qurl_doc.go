// Package qurl provides Go bindings for Qt 6.11 QtCore QUrl.
//
// QUrl represents a URL and provides parsing, manipulation, and
// component-extraction facilities. It supports both tolerant and strict
// parsing modes, percent-encoding/decoding, ACE (Punycode) conversion for
// internationalized domain names, and conversion to/from local file paths.
//
// In Go, each *QUrl is a handle to a heap-allocated C++ QUrl; call Close()
// for immediate release, otherwise a finalizer handles cleanup.
//
// # Enum and flag types
//
// ParsingMode          (qurl_parsingmode.go)              — Tolerant/Strict/Decoded
// UrlFormattingOption  (qurl_urlformattingoption.go)      — Remove*/PreferLocalFile/etc.
// ComponentFormattingOption (qurl_componentformattingoption.go)
//                                                         — PrettyDecoded/Encode*/Fully*
// ComponentFormattingOptions (qurl_componentformattingoption.go)
//                                                         — alias of ComponentFormattingOption
// FormattingOptions    (qurl_formattingoptions.go)        — QUrlTwoFlags in Qt
// UserInputResolutionOption (qurl_userinputresolutionoption.go)
//                                                         — Default/AssumeLocalFile
// UserInputResolutionOptions                              — alias of UserInputResolutionOption
// AceProcessingOption  (qurl_aceprocessingoption.go)      — IgnoreIDNWhitelist/AceTransitional
// AceProcessingOptions                                    — alias of AceProcessingOption
//
// # Flag type mapping (Q_DECLARE_FLAGS / QUrlTwoFlags ↔ Go)
//
// The Qt header declares:
//
//   Q_DECLARE_FLAGS(ComponentFormattingOptions, ComponentFormattingOption)
//   Q_DECLARE_FLAGS(AceProcessingOptions, AceProcessingOption)
//   Q_DECLARE_FLAGS(UserInputResolutionOptions, UserInputResolutionOption)
//   typedef QUrlTwoFlags<UrlFormattingOption, ComponentFormattingOption> FormattingOptions;
//
// The Q_DECLARE_FLAGS trio are represented in Go as aliases for their
// respective enum types (which in turn are aliases for int), so individual
// flag constants pass without conversion.
//
// QUrl::FormattingOptions is a QUrlTwoFlags<UrlFormattingOption, ComponentFormattingOption>.
// QUrlTwoFlags provides implicit constructors from both enum types, so in
// C++ any individual flag from either enum is accepted wherever FormattingOptions
// is expected. Go has no implicit conversions, so in this binding
// UrlFormattingOption, ComponentFormattingOption and FormattingOptions are
// ALL aliases for int. Any flag or flag combination — produced by bitwise OR
// as in C++ — works without explicit conversion in any method signature.
//
// Example (identical to Qt usage):
//
//   u.URLString(RemoveQuery | RemoveFragment)
//   u.ToEncodedBytes(FullyEncoded | StripTrailingSlash)
//   u.Host(FullyDecoded)
//   u.Query(PrettyDecoded)
//
// # Construction
//
//   - NewQUrl()                        → QUrl()
//   - QUrlFromString(s, mode)          → QUrl(const QString&, ParsingMode) [fromUtf8]
//   - QUrlFromQString(qs, mode)        → QUrl(const QString&, ParsingMode)
//   - Clone()                          → copy constructor
//   - Swap(other)                      → QUrl::swap
//
// # String representation
//
//   - URLString(options) string        → QUrl::url(FormattingOptions) → QString
//   - URL(options) *qtext.QString      → QUrl::url(FormattingOptions)
//   - ToString(options) string         → QUrl::toString(FormattingOptions) → QString
//   - ToQString(options) *qtext.QString → QUrl::toString(FormattingOptions)
//   - ToDisplayString(options) string  → QUrl::toDisplayString(FormattingOptions) → QString
//   - ToDisplayQString(options)        → QUrl::toDisplayString(FormattingOptions)
//   - Adjusted(options) *QUrl          → QUrl::adjusted(FormattingOptions)
//   - ToEncodedBytes(options) []byte   → QUrl::toEncoded(FormattingOptions) → QByteArray
//   - ToEncoded(options) *qtext.QByteArray → QUrl::toEncoded(FormattingOptions)
//
// # Components
//
// Each getter that returns a single URL component takes
// ComponentFormattingOptions (Qt: QFlags<ComponentFormattingOption>):
//
//   - Scheme / SchemeString                                 → scheme()
//   - Authority / AuthorityQString                          → authority(ComponentFormattingOptions)
//   - UserInfo / UserInfoQString                            → userInfo(ComponentFormattingOptions)
//   - UserName / UserNameQString                            → userName(ComponentFormattingOptions)
//   - Password / PasswordQString                            → password(ComponentFormattingOptions)
//   - Host / HostQString                                    → host(ComponentFormattingOptions)
//   - Path / PathQString                                    → path(ComponentFormattingOptions)
//   - FileName / FileNameQString                            → fileName(ComponentFormattingOptions)
//   - Query / QueryQString                                  → query(ComponentFormattingOptions)
//   - Fragment / FragmentQString                            → fragment(ComponentFormattingOptions)
//
// Setters have `…String` (Go string input) and — where non-trivial —
// `…QString` / `SetXxxQString` forms.
//
// # Relations
//
//   - Resolved(relative)               → QUrl::resolved(const QUrl&)
//   - IsRelative()                     → QUrl::isRelative()
//   - IsParentOf(url)                  → QUrl::isParentOf(const QUrl&)
//   - Matches(other, options)          → QUrl::matches(const QUrl&, FormattingOptions)
//   - IsLocalFile()                    → QUrl::isLocalFile()
//   - ToLocalFile() string             → QUrl::toLocalFile() → QString
//   - ToLocalFileQString()             → QUrl::toLocalFile()
//   - Detach() / IsDetached()          → QUrl::detach() / isDetached()
//
// # Validity
//
//   - IsValid()                        → QUrl::isValid()
//   - ErrorString() string             → QUrl::errorString() → QString
//   - ErrorQString() *qtext.QString    → QUrl::errorString()
//   - IsEmpty() / Clear()              → QUrl::isEmpty()/clear()
//
// # Static constructors and helpers
//
//   - FromEncoded([]byte, mode)             → QUrl::fromEncoded(QByteArrayView, ParsingMode)
//   - FromLocalFile(path) / FromLocalFileQString(qs)
//                                            → QUrl::fromLocalFile(const QString&)
//   - FromUserInputString(in, cwd, opts)     → QUrl::fromUserInput()
//   - FromUserInputQString(in, cwd, opts)    → QUrl::fromUserInput()
//   - FromPercentEncodingString([]byte)      → QUrl::fromPercentEncoding() → QString
//   - FromPercentEncoding([]byte)            → QUrl::fromPercentEncoding() → QString
//   - ToPercentEncodingBytes(s, excl, incl)  → QUrl::toPercentEncoding() → QByteArray
//   - ToPercentEncoding(qs, excl, incl)      → QUrl::toPercentEncoding() → QByteArray
//   - FromAceString([]byte) / FromAce([]byte) → QUrl::fromAce()
//   - ToAceBytes(s) / ToAce(qs)              → QUrl::toAce()
//   - ToStringListStrings(urls, opts) []string → QUrl::toStringList() → QStringList
//   - ToStringList(urls, opts) *qtext.QStringList → QUrl::toStringList()
//   - FromStringListStrings(strs, mode) []*QUrl → QUrl::fromStringList() → QList<QUrl>
//   - FromStringList(sl, mode) []*QUrl       → QUrl::fromStringList()
//
// # Not implemented
//
//   - QUrlQuery integration (QUrl::setQuery(const QUrlQuery&)) — QUrlQuery is
//     not bound yet. Use SetQueryString with a manually built "key=value" string.
//   - idnWhitelist() / setIdnWhitelist(const QStringList&) — pending QStringList
//     integration at package level.
//   - toCFURL / fromCFURL, toNSURL / fromNSURL — Darwin-only.
//   - QDataStream operator<< / operator>> — not idiomatic in Go.
//   - QDebug operator<< — not idiomatic in Go.
//   - qHash(QUrl) — Go can hash the encoded form directly.
//   - operator<, operator==, operator!= — use Matches / Compare (planned) or
//     compare the string forms returned by ToString / ToEncodedBytes.
//   - QUrlTwoFlags operators — Go uses plain int + bitwise OR.
//   - operator|(UrlFormattingOption, ComponentFormattingOption) free functions
//     in qurl.h — subsumed by plain int arithmetic in Go.
//
// # Note on ToStringList / FromStringList
//
// The C ABI implementation joins the QStringList elements with '\n' and
// splits on the Go side. QUrl normally percent-encodes literal newlines as
// %0A, so this is safe in practice, but a custom-formatted URL fragment that
// contains a literal '\n' (possible with FullyDecoded) would be split
// incorrectly. If that matters, call ToStringList one element at a time.
//
// This file is part of package qurl of github.com/vadimlxx/go-qt6core.
// Requires Qt >= 6.11 and a C++17 compiler; links via pkg-config Qt6Core.
package qurl
