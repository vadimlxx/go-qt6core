// Package qtext provides Go bindings for Qt 6.11 QtCore QTextBoundaryFinder.
//
// QTextBoundaryFinder locates grapheme, word, sentence, and line boundaries
// inside a UTF-16 string, following the Unicode UAX #29 / UAX #14
// segmentation algorithms. In Go, each *QTextBoundaryFinder is a handle to a
// heap-allocated C++ QTextBoundaryFinder; call Close() for immediate release,
// otherwise a finalizer handles cleanup.
//
// Construction:
//   - NewQTextBoundaryFinder()                    → QTextBoundaryFinder()
//                                                    [invalid until re-assigned]
//   - QTextBoundaryFinderFromString(type, *QString)
//                                                 → QTextBoundaryFinder(BoundaryType, const QString&)
//   - QTextBoundaryFinderFromGoString(type, s)    → same, with QString::fromUtf8(s)
//   - QTextBoundaryFinderFromView(type, QStringView)
//                                                 → QTextBoundaryFinder(BoundaryType, QStringView, ...)
//   - Clone()                                     → copy constructor
//   - Swap(other)                                 → QTextBoundaryFinder::swap
//
// Boundary types (enum BoundaryType):
//   - Grapheme  — user-perceived character
//   - Word      — UAX #29 word boundary
//   - Sentence  — UAX #29 sentence boundary
//   - Line      — UAX #14 line-break opportunity
//
// Boundary reasons (bitmask BoundaryReason, alias BoundaryReasons):
//   - NotAtBoundary (0)            — not on a boundary
//   - BreakOpportunity (0x1f)      — mask of the low 5 bits
//   - StartOfItem (0x20)           — start of a text item
//   - EndOfItem (0x40)             — end of a text item
//   - MandatoryBreak (0x80)        — mandatory break (e.g. newline)
//   - BoundarySoftHyphen (0x100)   — soft-hyphen break
//
// Navigation:
//   - ToStart / ToEnd / Position / SetPosition
//   - ToNextBoundary / ToPreviousBoundary — return the new position or -1
//   - IsAtBoundary / BoundaryReasons
//
// The Class uses QString and QStringView from the same qtext package; the
// boundary-related constants live in qtextboundaryfinder_boundarytype.go and
// qtextboundaryfinder_boundaryreason.go.
//
// Not implemented
//
//   - QTextBoundaryFinder(BoundaryType, const QChar*, qsizetype,
//     unsigned char *buffer, qsizetype bufferSize)
//     — the explicit-buffer constructor is an optimisation for callers that
//     want to reuse an external buffer; in Go this adds unsafe complexity
//     without a measurable benefit (Qt re-allocates internally on copy anyway).
//   - Move constructor and move assignment — Go transfers ownership through
//     the *QTextBoundaryFinder handle instead; no separate binding.
//   - Q_DECLARE_SHARED copy-on-write semantics — in Go use Clone() for an
//     independent copy.
//   - Java-style iterators (QTextBoundaryFinderIterator) and STL-style
//     iterators — not idiomatic in Go; use ToNextBoundary/ToPreviousBoundary
//     and BoundaryReasons directly.
//
// This file is part of package qtext of github.com/vadimlxx/go-qt6core.
// Requires Qt >= 6.11 and a C++17 compiler; links via pkg-config Qt6Core.
package qtext
