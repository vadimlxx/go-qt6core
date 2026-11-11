// Package qgeometry — QLine and QLineF bindings for Qt 6.11 QtCore.
//
// # Currently supported in these files
//
//   - QLine  — 2×QPoint line (p1, p2). Integer geometry.
//   - QLineF — 2×QPointF line (p1, p2). Floating-point geometry.
//
// Both are Q_PRIMITIVE_TYPE value types in Qt. In Go they are structs with unexported fields (p1, p2)
// and Qt-style accessor methods (P1(), P2(), X1(), Y1(), X2(), Y2(), Dx(), Dy(), Length(), Angle(),
// etc.). Construction goes through NewQLine / NewQLineF and their variants; mutation through
// Translate / SetP1 / SetP2 / SetLength / SetAngle / … ; transforms return new values via
// Translated / UnitVector / NormalVector. Passing through the C ABI uses flat 2-field structs
// (go_qline_t, go_qlinef_t) built from go_qpoint_t/go_qpointf_t in C mode, typedef'd to real
// QLine / QLineF in C++ mode — layout stays ABI-compatible with Qt.
//
// # Default constructors
//
// Qt's default constructors are zero-valued:
//
//	QLine()  → QLine((0,0), (0,0))   // isNull() == true (exact)
//	QLineF() → QLineF((0,0), (0,0))  // isNull() == true (fuzzy)
//
// Therefore the Go zero value QLine{} / QLineF{} is already a valid default, and
// NewQLineZero() / NewQLineFZero() are convenience aliases. There are no NewQLineDefault /
// NewQLineFDefault functions.
//
// # Operator mapping (QLine / QLineF)
//
// Both types have NO arithmetic operators (+, -, *, /) in Qt — only translation and geometric
// transforms. The mapping for what Qt does expose:
//
//	l == l2      → l.Equal(l2)  /  l.NotEqual(l2)    (QLine: exact; QLineF: fuzzy)
//	translate    → l.Translate(p) / l.TranslateXY(dx, dy)      (in place)
//	translated   → l.Translated(p) / l.TranslatedXY(dx, dy)    (returns new)
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
//	l.IntersectsWithPoint(other)             → (IntersectionType, QPointF)
//
// When there is no intersection, IntersectsWithPoint returns (0, 0) as the point — Qt does not
// write to the out-parameter in that case, so Go's zero-initialised value is used.
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
// # Protected access
//
// QLine and QLineF have no protected members — only public API and private implementation.
// Nothing needs to be exposed via a helper subclass.
//
// # Not yet implemented
//
//   - QDataStream operator<< / operator>>  — will appear with package qdatastream.
//   - QDebug operator<<                    — will appear with package qdebug.
//   - operator== / operator!= as overloaded operators
//                                          — Go has no operator overloading; use Equal / NotEqual.
//   - Free functions qFuzzyCompare(QLineF, QLineF) and qFuzzyIsNull(QLineF)
//                                          — covered indirectly by Equal() / IsNull().
//   - Direct field access (l.P1, l.P2)     — Go does not allow a field and a method with the same
//                                            name; use the accessor methods instead.
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
// # Handle semantics (matching QPoint / QPointF)
//
// QLine and QLineF are exposed as HANDLES (Go structs holding a C++ pointer). Accessors and
// mutators go directly to the heap-allocated C++ object. Mutations through one Go variable are
// visible through all variables pointing to the same C++ object. Clone() produces an independent
// copy.
//
// Memory management:
//
//   - NewQLine / NewQLineF / arithmetic / Clone / P1() / P2() / Center() / PointAt() /
//     UnitVector() / NormalVector() allocate a new heap C++ object.
//   - The finalizer frees it; Close() frees immediately and disables the finalizer.
//   - Do not copy *QLine by value expecting independent copies — copy the pointer; use Clone()
//     for a truly independent copy.
//
// # Fundamental limitation with embedded QLine / QLineF
//
// Qt stores QLine / QLineF **by value** inside other classes. A Go *QLine handle CANNOT
// reference such internal storage; it always owns a separate heap copy. If C++ code later
// mutates an embedded QLine, the previously-obtained *QLine will NOT see the change. Use the
// containing class's fresh accessors (e.g. QPainterPath element getters) to read the current
// state. This is a limitation of Qt's public API, not of the binding.
// Requires Qt >= 6.11, C++17, pkg-config Qt6Core.

package qgeometry
