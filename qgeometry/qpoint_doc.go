// Package qgeometry — QPoint and QPointF bindings for Qt 6.11 QtCore.
//
// # Currently supported in this file
//
//   - QPoint  — 2×int point (x, y)
//   - QPointF — 2×qreal point (x, y)
//
// Both are Q_PRIMITIVE_TYPE value types in Qt. In Go they are structs with
// unexported fields and Qt-style accessor methods X()/Y(). Construction goes
// through NewQPoint / NewQPointF (or the zero value / NewQPointZero /
// NewQPointFZero); mutation through SetX / SetY or the Assign methods.
// Passing through the C ABI uses flat 2-field structs (go_qpoint_t,
// go_qpointf_t) in C mode, typedef'd to real QPoint / QPointF in C++ mode,
// so the layout stays ABI-compatible with Qt.
//
// # Default constructors
//
// Unlike QSize/QSizeF, Qt's default constructors are zero-valued:
//
//	QPoint()  → QPoint(0, 0)
//	QPointF() → QPointF(0., 0.)
//
// Therefore the Go zero value QPoint{} / QPointF{} is already a valid default,
// and NewQPointZero() / NewQPointFZero() are just convenience aliases. There
// are no NewQPointDefault / NewQPointFDefault functions — they would be
// redundant with the zero constructors.
//
// # Operator mapping (QPoint / QPointF)
//
//	p + q        → p.Add(q)
//	p - q        → p.Sub(q)
//	p * i        → p.MulInt(i)        (QPoint only — QCheckedInt arithmetic)
//	p * f        → p.MulFloat(f)      (QPoint: qSaturateRound; QPointF: raw)
//	p / f        → p.DivFloat(f)      (QPoint: qSaturateRound; QPointF: raw)
//	-p           → p.Neg()
//	+p           → p.UnaryPlus()
//	p == q       → p.Equal(q) / p.NotEqual(q)   (QPoint: exact; QPointF: fuzzy)
//	p += q       → p.AddAssign(q)
//	p -= q       → p.SubAssign(q)
//	p *= i       → p.MulIntAssign(i)  (QPoint only)
//	p *= f       → p.MulFloatAssign(f)
//	p /= f       → p.DivFloatAssign(f)
//
// Dot product:
//
//	QPoint::dotProduct(p, q)  → p.DotProduct(q) int
//	QPointF::dotProduct(p, q) → p.DotProduct(q) float64
//
// # QPoint ↔ QPointF
//
//   - (QPoint).ToPointF()        — mirrors QPoint::toPointF()
//   - (QPointF).ToPoint()        — mirrors QPointF::toPoint()   (qSaturateRound per component)
//   - QPointFFromPoint(p)        — mirrors QPointF(const QPoint&)
//   - (QPointF).EqualQPoint(q)   — mirrors comparesEqual(QPointF, QPoint)
//
// # Semantic deviations from Qt
//
//   - QPointF.IsNull() uses qIsNull on each component (fuzzy near-zero).
//   - QPointF.Equal uses qFuzzyCompare on each component. Tiny differences may
//     compare as Equal — matching Qt.
//   - QPoint.MulInt / Add / Sub use QCheckedInt arithmetic; in release builds
//     they wrap on overflow (no saturation), matching Qt 6.11 release behavior.
//   - QPoint.MulFloat / DivFloat and QPointF.ToPoint() use qSaturateRound.
//   - QPoint.DotProduct / QPointF.DotProduct may overflow / produce inf/nan
//     exactly as Qt does; no extra safety is added.
//   - Division by zero: QPoint.DivFloat(0) → saturate(inf) → INT_MAX/INT_MIN
//     (Qt release behavior). QPointF.DivFloat(0) → IEEE inf/nan.
//
// # Protected access
//
// QPoint and QPointF have no protected members — only public API and private
// implementation. Nothing needs to be exposed via a helper subclass.
//
// # Not yet implemented
//
//   - QDataStream operator<< / operator>>  — will appear with package qdatastream.
//   - QDebug operator<<                    — will appear with package qdebug.
//   - qHash(QPoint)                        — user can hash X()/Y() directly.
//                                              qHash(QPointF) is deleted in Qt itself.
//   - rx() / ry()                          — covered by SetX / SetY.
//   - std::tuple protocol / structured bindings
//                                          — no Go equivalent; use X()/Y().
//   - get<I>(QPoint) / get<I>(QPointF)     — same; use X()/Y().
//   - QPoint::toCGPoint / QPointF::fromCGPoint / toCGPoint (Darwin)
//                                          — platform-specific.
//   - Free functions qFuzzyCompare(QPointF, QPointF) and
//     qFuzzyIsNull(QPointF)                — covered indirectly by Equal / IsNull.
//   - Direct field access (p.X, p.Y, …)    — Go does not allow a field and a
//     method with the same name; use the accessor methods instead.
//
// # Future work
//
// Future additions will include:
//
//   - QRect / QRectF in this same package, with cross-type conversions
//     (e.g. QRect::topLeft(), QRect::size(), QRect::margins()).
//   - QDataStream serialization once qdatastream exists.
//   - QDebug formatting once qdebug exists.
//
// Requires Qt >= 6.11, C++17, pkg-config Qt6Core.
package qgeometry