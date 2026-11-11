// Package qgeometry — QPoint and QPointF bindings for Qt 6.11 QtCore.
//
// # Handle semantics
//
// Both QPoint and QPointF are exposed as HANDLES (Go structs holding a C++ pointer). All
// accessors and mutators go directly to the heap-allocated C++ object:
//
//	p := NewQPoint(10, 20)   // allocates a C++ QPoint
//	x := p.X()               // p->x() — reads from C++
//	p.SetX(15)               // p->setX(15) — mutates the C++ object
//	p.Close()                // deletes the C++ QPoint
//
// To obtain a value-semantics copy, use Clone(). To share mutations, pass the same *QPoint
// around — all handles referring to the same C++ object see the same state.
//
// Memory management:
//
//   - NewQPoint / NewQPointF / arithmetic results / Clone allocate a new heap C++ object.
//   - The finalizer frees it; Close() frees immediately and disables the finalizer.
//   - Do not copy *QPoint by value expecting independent copies — copy the pointer
//     (both point to the same C++ object); use Clone() for a truly independent copy.
//
// # Fundamental limitation with embedded QPoints
//
// Qt stores QPoint / QPointF **by value** inside other classes (QLine contains two QPoint
// members; QRect contains two; …). Because those members are not separately heap-allocated,
// a Go *QPoint handle CANNOT reference them. Consequences:
//
//   - QLine.P1() returns a NEW heap QPoint with a copy of the internal value.
//   - If C++ code later mutates the QLine's internal QPoint (e.g. via translate()), the
//     previously-returned *QPoint remains unchanged.
//   - To read the current value, call P1() again (or use X1()/Y1()/X2()/Y2(), which read
//     ints directly from C++ each time).
//
// This is a limitation of Qt's public API (Qt 6 removed QLine::rpt1()), not of the binding.
// A borrowed-pointer "QPointView" type could be added in the future if a Qt class exposes
// `QPoint* someMethod()`, but QLine/QPoint API does not offer that.
//
// # Not yet implemented
//
//   - QDataStream operator<< / operator>>  — will appear with package qdatastream.
//   - QDebug operator<<                    — will appear with package qdebug.
//   - qHash(QPoint) / qHash(QPointF)       — not needed in Go (user hashes X()/Y()).
//   - rx() / ry()                          — covered by SetX / SetY.
//   - std::tuple protocol                  — no Go equivalent; use X()/Y().
//   - QPoint::toCGPoint / QPointF::fromCGPoint / toCGPoint (Darwin) — platform-specific.
//   - Free functions qFuzzyCompare / qFuzzyIsNull — covered by Equal / IsNull.
//   - QPointView (borrowed reference)      — future work; see limitation above.
//
// Requires Qt >= 6.11, C++17, pkg-config Qt6Core.
package qgeometry
