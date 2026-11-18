// QLatin1Char is a Go binding for Qt 6.11 QtCore's QLatin1Char.
//
// QLatin1Char is a trivial value type wrapping a single Latin-1 byte
// (C++ `struct { char ch; }`, Q_PRIMITIVE_TYPE). The Go wrapper is
// `type QLatin1Char byte` — zero allocation, trivially copyable,
// comparable with ==.
//
// Relationship with package qchar:
//   - QLatin1Char → QChar: call (*QLatin1Char).ToQChar().
//   - QChar → QLatin1Char: call NewQLatin1CharFromQChar(q), which returns
//     (c, ok) with ok=false if q is outside U+0000..U+00FF.
//
// Note on performance: because QLatin1Char wraps a single byte, all
// methods are pure bit-for-bit operations that could be implemented
// without cgo. This package nonetheless goes through cgo for consistency
// with the sibling package qchar and to keep one version-check point.
// For maximum throughput in a tight loop, operate on uint8(lc) directly.
//
// # C ABI
//
// The C ABI lives in qlatin1char_c.h as a dual-mode header (C declarations
// for cgo's C preamble + C++ inline definitions compiled via
// qtext_cxx_impl.cpp). Per-function notes live in qlatin1char.go.
//
// This package is a subpackage of github.com/vadimlxx/go-qt6core.
// Requires Qt >= 6.11 and a C++17 compiler; links via pkg-config Qt6Core.
package qtext
