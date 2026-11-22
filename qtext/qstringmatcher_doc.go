// QStringMatcher is a Go binding for Qt 6.11 QtCore's QStringMatcher.
//
// QStringMatcher is a fast substring matcher using the Boyer-Moore-Horspool
// algorithm with a skip table. It is efficient when the same pattern is
// searched many times across different subject strings.
//
// # Type semantics
//
// QStringMatcher is a handle to a C++ QStringMatcher. Each *QStringMatcher
// owns the underlying C++ object.
//
// # Ownership
//
// The C++ QStringMatcher is released by one of:
//   - Close(): eager, deterministic release;
//   - runtime.AddCleanup: safety net when *QStringMatcher becomes unreachable.
//
// Close is optional, idempotent, and nil-safe. Using a QStringMatcher after
// Close panics with "qstringmatcher: use after Close (or nil receiver)".
//
// # Construction
//
//   - NewQStringMatcher()                                  — QStringMatcher()
//   - NewQStringMatcherFromQString(p *QString, cs)         — QStringMatcher(const QString&, cs)
//   - NewQStringMatcherFromString(s string, cs)            — same, UTF-8 decoded
//   - NewQStringMatcherFromQStringView(v QStringView, cs)  — QStringMatcher(QStringView, cs)
//   - NewQStringMatcherFromQChars(units []uint16, cs)      — QStringMatcher(const QChar*, qsizetype, cs)
//   - (m).Clone()                                          — QStringMatcher(const QStringMatcher&)
//   - (m).Assign(other)                                    — QStringMatcher::operator=
//   - NewQStringMatcher(...)                               — dispatcher on argument type
//
// # Searching
//
//   - IndexIn(str *QString, from int) int          — indexIn(const QString&, qsizetype)
//   - IndexInString(s string, from int) int        — same, UTF-8 decoded
//   - IndexInView(v QStringView, from int) int     — indexIn(QStringView, qsizetype)
//   - IndexInQChars(units []uint16, from int) int  — indexIn(const QChar*, qsizetype, qsizetype)
//
// Returns -1 if the pattern is not found. The `from` offset is in UTF-16
// code units (matching Qt's QString semantics, not byte offset).
//
// # Setters / getters
//
//   - SetPattern(pattern *QString)                — setPattern(const QString&)
//   - SetPatternString(s string)                  — same, UTF-8 decoded
//   - SetPatternView(v QStringView)               — convenience (via QString)
//   - SetCaseSensitivity(cs qdata.CaseSensitivity) — setCaseSensitivity(Qt::CaseSensitivity)
//   - Pattern() *QString                          — pattern() const
//   - PatternString() string                      — same, encoded as UTF-8
//   - PatternView() QStringView                   — patternView() const noexcept
//   - CaseSensitivity() qdata.CaseSensitivity     — caseSensitivity() const
//
// # Case sensitivity
//
// The cs argument uses qdata.CaseSensitivity:
//   - qdata.CaseInsensitive — case-insensitive matching
//   - qdata.CaseSensitive   — case-sensitive matching (default for Qt)
//
// Note: this is the qdata-typed constant set. Some other classes in this
// package (notably QString) use a `caseSensitive bool` parameter instead;
// the two conventions coexist for historical reasons.
//
// # PatternView lifetime
//
// PatternView returns a non-owning QStringView over the matcher's internal
// pattern storage. The view is invalidated by any subsequent SetPattern /
// SetPatternString / SetPatternView / Assign call, or by Close. If you need
// a stable copy of the pattern, use Pattern() or PatternString().
//
// # Not yet bound
//
//   - updateSkipTable() — protected in Qt, called automatically by setPattern;
//     not exposed separately.
//   - Copy-assignment via Go syntax — use Assign method instead.
//   - QStringMatcher(const QChar*, qsizetype, cs) — covered by
//     NewQStringMatcherFromQChars (functionally identical to the QStringView
//     constructor, since Qt's implementation forwards to it).
//   - indexIn(const QChar*, qsizetype, from) — covered by IndexInQChars
//     (also forwards to the QStringView overload in Qt).
//
// This package is a subpackage of github.com/vadimlxx/go-qt6core.
// Requires Qt >= 6.11 and a C++17 compiler; links via pkg-config Qt6Core.
package qtext