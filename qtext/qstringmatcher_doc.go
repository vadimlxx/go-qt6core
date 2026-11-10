// Package qstringmatcher provides a Go binding for Qt 6.11.2 QtCore's QStringMatcher.
//
// QStringMatcher is a utility that preprocesses a pattern string to allow fast repeated
// substring searches. It is NOT a QObject — it is a plain C++ value class with a private
// implementation (d_ptr), a public copy constructor, and a public destructor.
//
// In Go, each *QStringMatcher is a handle to a heap-allocated C++ QStringMatcher. The
// wrapper owns the C++ object; runtime.SetFinalizer + explicit Close() free it.
//
// # Type mapping
//
//	Qt                             Go
//	QStringMatcher                 *qstringmatcher.QStringMatcher
//	Qt::CaseSensitivity            bool (caseSensitive: true = CaseSensitive)
//	QString                        *qstring.QString
//	QStringView                    qstringview.QStringView
//	QChar*/char16_t*               []uint16
//	qsizetype (return of indexIn)  int  (-1 for "not found")
//
// # Interop with sibling packages
//
// qstringmatcher imports only qstring and qstringview. No import cycles: neither of
// those packages imports qstringmatcher.
//
// Construction:
//   - New()                         → default-constructed matcher (empty pattern)
//   - FromQString(pattern, cs)      → QStringMatcher(const QString&, cs)
//   - FromQStringView(v, cs)        → QStringMatcher(QStringView, cs)
//   - FromUtf16(units, cs)          → QStringMatcher(const QChar*, len, cs)
//
// Copy / assign:
//   - Clone()                       → QStringMatcher(const QStringMatcher&)
//   - Assign(other)                 → operator=(const QStringMatcher&)
//
// Pattern / sensitivity:
//   - SetPattern(q)                 → setPattern(const QString&)
//   - SetCaseSensitivity(bool)      → setCaseSensitivity(Qt::CaseSensitivity)
//   - CaseSensitivity() bool        → caseSensitivity()
//   - Pattern() *qstring.QString    → pattern()
//   - PatternView() qstringview.QStringView → patternView()
//
// Search:
//   - IndexIn(q, from)              → indexIn(const QString&, from)
//   - IndexInView(v, from)          → indexIn(QStringView, from)
//   - IndexInUtf16(units, from)     → indexIn(const QChar*, len, from)
//
// # Protected members
//
// QStringMatcher has NO protected members in Qt 6.11.2. Its only non-public member,
// updateSkipTable(), is a private helper — not part of the API and not exposed.
//
// # PatternView lifetime
//
// PatternView() returns a zero-copy qstringview.QStringView into the matcher's internal
// pattern buffer. The view keeps the *QStringMatcher alive. Do NOT call SetPattern or
// Close on the matcher while the view (or any of its slices) is in use — the buffer may
// reallocate or be freed. Use Pattern() if you need an independent copy.
//
// # TODO (not yet bound, or API differences vs Qt 6.11.2)
//
//   - QStringMatcher(const QChar *uc, qsizetype len, Qt::CaseSensitivity cs) is an
//     inline forwarding constructor equivalent to QStringMatcher(QStringView(uc, len), cs).
//     Exposed only as FromUtf16; FromQStringView covers the QStringView overload.
//
//   - updateSkipTable() (private helper): not exposed (not part of the public API).
//
//   - Qt::CaseSensitivity is not a separate Go type: pass `caseSensitive bool`
//     (true = CaseSensitive, false = CaseInsensitive), matching the style used by the
//     qstring / qstringview / qbytearray bindings.
//
//   - QStringMatcher is copyable in C++ (no deleted copy ctor, no move-only). This
//     binding exposes Clone() and Assign(other) — both safe because QString is COW and
//     QStringView inside the matcher points into the matcher's own q_pattern buffer.
//
// # Requirements
//
// Qt >= 6.11.2, C++17, pkg-config Qt6Core. Version check is performed at compile time in
// qstringmatcher_c.cpp via internal/qt6core/check_qt_version.h.
//
// This package is a subpackage of github.com/vadimlxx/go-qt6core.
package qtext
