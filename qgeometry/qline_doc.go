// Package qgeometry — QLine and QLineF bindings for Qt 6.11 QtCore.
//
// # Value semantics
//
// Both QLine and QLineF are Q_PRIMITIVE_TYPE value types in Qt. In Go they are plain structs
// with unexported fields (x1, y1, x2, y2) holding a full copy of the data. Assignment and
// function-passing copy the value; mutating one copy does not affect another. There are no
// pointers, no finalizers and no Close().
//
// This matches the current state of QPoint / QPointF / QSize / QSizeF / QMargins / QMarginsF,
// which are also value types in this package.
//
// # Construction
//
//	QLine()  → QLine((0,0), (0,0))   // isNull() == true (exact)
//	QLineF() → QLineF((0,0), (0,0))  // isNull() == true (fuzzy)
//
// The Go zero value QLine{} / QLineF{} is therefore already a valid default. Use
// NewQLineZero() / NewQLineFZero() as explicit aliases.
//
// # Operator mapping (QLine / QLineF)
//
// Both types have NO arithmetic operators (+, -, *, /) in Qt — only translation and
// geometric transforms. The mapping for what Qt does expose:
//
//	l == l2      → l.Equal(l2)  /  l.NotEqual(l2)    (QLine: exact; QLineF: fuzzy)
//	translate    → l.Translate(p) / l.TranslateXY(dx, dy)      (returns updated value)
//	translated   → l.Translated(p) / l.TranslatedXY(dx, dy)    (returns new value)
//
// Note: in Go, because QLine / QLineF are value types, methods like Translate return a new
// value rather than mutating in place. Assign back to update:
//
//	l = l.Translate(p)
//
// # QLine ↔ QLineF
//
//   - (QLine).ToLineF()          — mirrors QLine::toLineF()
//   - (QLineF).ToLine()          — mirrors QLineF::toLine()  (QPointF::toPoint per endpoint)
//   - QLineFFromLine(l)          — mirrors QLineF(const QLine&)
//   - (QLineF).EqualLine(l)      — mirrors comparesEqual(QLineF, QLine)
//
// # QLineF-specific geometry
//
//   - Length() / SetLength(len)  — Euclidean length; SetLength keeps P1 fixed.
//   - Angle() / SetAngle(deg)    — angle in degrees (0 = +x direction, CCW); SetAngle keeps P1 fixed.
//   - AngleTo(other)             — angle between lines in degrees (CCW).
//   - UnitVector()               — unit-length vector (P1 = (0,0), P2 = direction).
//   - NormalVector()             — vector perpendicular to l (P1 = original P1, length = Length()).
//   - PointAt(t)                 — interpolation point (t=0 → P1, t=1 → P2).
//   - fromPolar(length, angle)   — QLineFFromPolar(len, deg) static constructor.
//
// # Angle convention (Qt screen coordinates — y grows downward)
//
// Qt uses screen coordinates where +y points DOWN. Angles are counter-clockwise from +x,
// normalised to [0, 360):
//
//	0°   → right (+x)
//	90°  → up on screen (-y)     ← negative y!
//	180° → left (-x)
//	270° → down on screen (+y)
//
// This is NOT the math-textbook convention (y-up). Users coming from math backgrounds
// commonly expect 90° to point toward +y — in Qt it points toward −y.
//
// angleTo(other) = (other.angle() - this.angle() + 360) mod 360. So a rightward line
// to a downward line is 270°, not 90°.
//
// # Intersection (QLineF only)
//
// IntersectionType is exposed as the IntersectionType type with three constants:
//
//	NoIntersection           — lines do not intersect
//	BoundedIntersection      — intersect within both segments
//	UnboundedIntersection    — intersect only when extended past the segment
//
// Two methods:
//
//	l.Intersects(other)                      → IntersectionType
//	l.IntersectsWithPoint(other)             → (IntersectionType, QPointF, bool)
//
// When there is no intersection, IntersectsWithPoint returns (NoIntersection, QPointF{}, false).
// Check the bool before using the point.
//
// # Semantic deviations from Qt
//
//   - QLine.IsNull() uses exact comparison p1 == p2. QLineF.IsNull() uses qFuzzyCompare(pt1, pt2).
//   - QLine.Equal() is exact; QLineF.Equal() is fuzzy (both endpoints compared via qFuzzyCompare).
//   - QLine.Center() uses 64-bit intermediate addition ((x1+x2)/2) to avoid int overflow — matching Qt.
//   - SetLength preserves P1 and moves P2 along the current direction; if the original length is 0,
//     P2 is not modified (matching Qt).
//   - SetAngle preserves P1 and recomputes P2 from length + angle (angle in degrees).
//
// # Not yet implemented
//
//   - QDataStream operator<< / operator>>  — will appear with package qdatastream.
//   - QDebug operator<<                    — will appear with package qdebug.
//   - operator== / operator!=              — use Equal / NotEqual.
//   - Free functions qFuzzyCompare(QLineF, QLineF) and qFuzzyIsNull(QLineF)
//                                          — covered indirectly by Equal() / IsNull().
//   - Direct field access (l.x1, l.y1, …)  — Go does not allow a field and a method with the
//                                            same name; use the accessor methods instead.
//   - std::tuple protocol                  — no equivalent in Go.
//
// # Future work
//
// Future additions will include:
//
//   - QRect / QRectF in this same package, referencing QLine, QPoint, QSize for construction and
//     inspection (e.g. QRect::topLeft(), QRect::size()).
//   - QDataStream serialization once qdatastream exists.
//   - QDebug formatting once qdebug exists.
//
// Requires Qt >= 6.11, C++17, pkg-config Qt6Core.
package qgeometry
