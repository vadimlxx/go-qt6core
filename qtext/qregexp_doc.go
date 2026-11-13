package qtext

// QRegularExpression, QRegularExpressionMatch, QRegularExpressionMatchIterator
//
// # Layout
//
// The C ABI is a single header/source pair (qregexp_c.h / qregexp_c.cpp),
// mirroring the upstream Qt6 qregularexpression.h, which holds all three
// classes together. Go wrappers stay split by type (qregexp.go,
// qregexpmatch.go, qregexpmatchiter.go) to match Go idioms and the
// per-type test files.
//
// # Not implemented
//
// The following APIs are not exposed in this binding yet:
//
//   - QRegularExpression::namedCaptureGroups() — returns QStringList.
//   - QRegularExpressionMatch::capturedTexts() — returns QStringList.
//   - QRegularExpression comparison operators (==, !=) — use Compare instead.
//   - QRegularExpressionMatch comparison operators.
//   - QRegularExpressionMatchIterator range-based for loop — not applicable in Go.
//   - qHash(QRegularExpression).
//   - QRegularExpression::optimize() — no effect in the C ABI wrapper.
//   - Zero-copy captured()/hasCaptured() overloads that accept an already-built
//     *QString name; current Name-suffixed methods take a Go string.
//
// # Enum notes
//
// PatternOption, MatchType, MatchOption, and WildcardConversionOption are
// represented as Go int-based types. Flags are combined with bitwise OR.
//
// # Lifetime
//
// QRegularExpressionMatch keeps a reference to the subject passed to
// MatchString or MatchView so that CapturedView remains valid. Call Close()
// to release resources early; otherwise finalizers handle cleanup.
//
// QRegularExpressionMatchIterator owns a copy of the subject internally.