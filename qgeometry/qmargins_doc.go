// Package qgeometry — QMargins and QMarginsF bindings for Qt 6.11 QtCore.
//
// # Handle semantics
//
// Both QMargins and QMarginsF are exposed as HANDLES (Go structs holding a C++ pointer). All
// accessors and mutators go directly to the heap-allocated C++ object:
//
//	m := NewQMargins(1, 2, 3, 4)   // allocates a C++ QMargins
//	x := m.Left()                   // m->left() — reads from C++
//	m.SetLeft(10)                   // m->setLeft(10) — mutates the C++ object
//	m.Close()                       // deletes the C++ QMargins
//
// To obtain a value-semantics copy, use Clone(). To share mutations, pass the same *QMargins
// around — all handles referring to the same C++ object see the same state.
//
// Memory management:
//
//   - NewQMargins / NewQMarginsF / arithmetic results / Clone / ToMarginsF / ToMargins /
//     QMarginsFFromMargins allocate a new heap C++ object.
//   - The finalizer frees it; Close() frees immediately and disables the finalizer.
//   - Do not copy *QMargins by value expecting independent copies — copy the pointer; use
//     Clone() for a truly independent copy.
//
// # Fundamental limitation with embedded QMargins
//
// Qt stores QMargins / QMarginsF **by value** inside other classes. A Go *QMargins handle
// CANNOT reference such internal storage; it always owns a separate heap copy. If C++ code
// later mutates an embedded QMargins, the previously-obtained *QMargins will NOT see the
// change. Use the containing class's fresh accessors to read the current state. This is a
// limitation of Qt's public API, not of the binding.
//
// # Operator mapping (QMargins / QMarginsF)
//
// Because Go has no operator overloading, Qt operators map to methods:
//
//	a + b        → a.Add(b)
//	a - b        → a.Sub(b)
//	a + i        → a.AddInt(i)            (QMargins only)
//	a - i        → a.SubInt(i)            (QMargins only)
//	a * i        → a.MulInt(i)            (QMargins only)
//	a / i        → a.DivInt(i)            (QMargins only)
//	a + f        → a.AddFloat(f)          (QMarginsF only)
//	a - f        → a.SubFloat(f)          (QMarginsF only)
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
//   - (QMargins).ToMarginsF()          — mirrors QMargins::toMarginsF()
//   - (QMarginsF).ToMargins()          — mirrors QMarginsF::toMargins()  (qSaturateRound)
//   - QMarginsFFromMargins(m)          — mirrors QMarginsF(const QMargins&)
//   - (QMarginsF).EqualQMargins(m)     — mirrors comparesEqual(QMarginsF, QMargins)
//
// # Semantic deviations from Qt
//
//   - QMarginsF.IsNull() uses qFuzzyIsNull on each component. In Go this delegates to the
//     C++ method, so semantics match exactly.
//   - QMarginsF.Equal uses qFuzzyCompare on each component. Tiny double differences may
//     compare as Equal — matching Qt.
//   - Integer arithmetic (AddInt, MulInt, DivInt, …) on QMargins wraps like C++ int; it does
//     NOT use saturating arithmetic. This matches Qt 6.11 release builds.
//   - Division by zero: QMargins.DivInt(0) triggers a Go panic (C++ UB → SIGFPE).
//     QMarginsF.DivFloat(0) produces IEEE-754 inf/nan (Qt asserts only in debug builds).
//   - Floating-point forms (MulFloat, DivFloat) use Qt's qSaturateRound for rounding to int —
//     identical to C++ operator*(QMargins, qreal) etc.
//
// # Protected access
//
// QMargins and QMarginsF have no protected members — only public API and private
// implementation. Nothing needs to be exposed via a helper subclass.
//
// # Not yet implemented
//
//   - QDataStream operator<< / operator>>  — will appear with package qdatastream.
//   - QDebug operator<<                    — will appear with package qdebug.
//   - std::tuple protocol                  — no Go equivalent.
//   - get<I>(QMargins) / get<I>(QMarginsF) — same; use Left()/Top()/Right()/Bottom().
//   - Free functions qFuzzyCompare / qFuzzyIsNull for QMarginsF
//                                          — covered indirectly by Equal() / IsNull().
//   - Direct field access (m.Left, m.Top, …)
//                                          — Go does not allow a field and a method with the
//                                            same name; use the accessor methods instead.
//
// # Future work
//
// Future additions will include:
//
//   - QRect / QRectF / QPoint / QPointF / QSize / QSizeF cross-type conversions that involve
//     QMargins (e.g. QRect::margins()).
//   - QDataStream serialization once qdatastream exists.
//   - QDebug formatting once qdebug exists.
//
// Requires Qt >= 6.11, C++17, pkg-config Qt6Core.
package qgeometry
