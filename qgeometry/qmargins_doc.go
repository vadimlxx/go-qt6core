// Package qgeometry provides Go bindings for Qt 6.11 QtCore geometry value types.
//
// # Currently supported
//
//   - QMargins  — 4×int margins (left, top, right, bottom)
//   - QMarginsF — 4×qreal margins
//
// Both are Q_RELOCATABLE_TYPE value types in Qt. In Go they are structs
// with unexported fields and Qt-style accessor methods (Left(), Top(),
// Right(), Bottom()). Construction goes through NewQMargins / NewQMarginsF
// (or the zero-value / NewQMarginsZero / NewQMarginsFZero); mutation goes
// through SetLeft / SetTop / SetRight / SetBottom or the Assign methods.
// Passing through the C ABI uses flat 4-field structs (go_qmargins_t /
// go_qmarginsf_t) in C mode, typedef'd to real QMargins / QMarginsF in
// C++ mode, so the layout stays ABI-compatible with Qt.
//
// # Operator mapping
//
// Because Go has no operator overloading, Qt operators map to methods:
//
//	a + b        → a.Add(b)
//	a - b        → a.Sub(b)
//	a + i        → a.AddInt(i)            (QMargins only)
//	a - i        → a.SubInt(i)            (QMargins only)
//	a * i        → a.MulInt(i)            (QMargins only)
//	a / i        → a.DivInt(i)            (QMargins only)
//	a * f        → a.MulFloat(f)
//	a / f        → a.DivFloat(f)
//	a | b        → a.Or(b)                (component-wise max)
//	-a           → a.Neg()
//	+a           → a.UnaryPlus()
//	a == b       → a.Equal(b)  /  a.NotEqual(b)
//	a += b       → a.AddAssign(b)
//	a -= b       → a.SubAssign(b)
//	a += i       → a.AddIntAssign(i)      (QMargins only)
//	a -= i       → a.SubIntAssign(i)      (QMargins only)
//	a *= i       → a.MulIntAssign(i)      (QMargins only)
//	a /= i       → a.DivIntAssign(i)      (QMargins only)
//	a += f       → a.AddFloatAssign(f)    (QMarginsF only)
//	a -= f       → a.SubFloatAssign(f)    (QMarginsF only)
//	a *= f       → a.MulFloatAssign(f)
//	a /= f       → a.DivFloatAssign(f)
//
// # QMargins ↔ QMarginsF
//
//   - QMargins.ToMarginsF()           — mirrors QMargins::toMarginsF()
//   - QMarginsF.ToMargins()           — mirrors QMarginsF::toMargins()  (qSaturateRound)
//   - QMarginsFFromMargins(m)         — mirrors QMarginsF(const QMargins&)
//   - (QMarginsF).EqualQMargins(m)    — mirrors comparesEqual(QMarginsF, QMargins)
//
// # Semantic deviations from Qt
//
//   - QMarginsF.IsNull() uses qFuzzyIsNull on each component. In Go this
//     delegates to the C++ method, so semantics match exactly.
//   - QMarginsF.Equal uses qFuzzyCompare on each component. Two QMarginsF
//     values with tiny double differences may compare as Equal — matching Qt.
//   - Integer arithmetic (AddInt, MulInt, DivInt, …) on QMargins wraps like
//     C++ int; it does NOT use saturating arithmetic. This matches Qt 6.11
//     release builds (QCheckedInt asserts only in debug builds).
//   - Division by zero: QMargins.DivInt(0) triggers a Go panic (C++ UB →
//     SIGFPE). QMarginsF.DivFloat(0) produces IEEE-754 inf/nan (Qt asserts
//     only in debug builds).
//   - Floating-point forms (MulFloat, DivFloat) use Qt's qSaturateRound for
//     rounding to int — identical to C++ operator*(QMargins, qreal) etc.
//
// # Protected access
//
// QMargins and QMarginsF have no protected members — only public API and
// private implementation. Nothing needs to be exposed via a helper subclass.
//
// # Not yet implemented
//
//   - QDataStream operator<< / operator>>  — will appear with package qdatastream.
//   - QDebug operator<<                    — will appear with package qdebug.
//   - std::tuple protocol / structured bindings
//                                          — no Go equivalent; access public fields directly.
//   - get<I>(QMargins) / get<I>(QMarginsF) — same; access public fields directly.
//   - Free functions qFuzzyCompare(QMarginsF, QMarginsF) and
//     qFuzzyIsNull(QMarginsF)              — covered indirectly by Equal() and IsNull().
//   - Direct field access (m.Left, m.Top, …)   — Go does not allow a field
//     and a method with the same name; use the accessor methods instead.
//
// # Future work
//
// Future additions will include:
//
//   - QDataStream serialization once qdatastream exists.
//   - QDebug formatting once qdebug exists.
//   - QRect / QRectF / QPoint / QPointF / QSize / QSizeF in this same package,
//     with cross-type conversions (e.g. QRect::margins(), QRectF::margins()).
//
// This package is a subpackage of github.com/vadimlxx/go-qt6core.
// Requires Qt >= 6.11, C++17, pkg-config Qt6Core.
package qgeometry
