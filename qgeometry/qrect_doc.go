// Package qgeometry — QRect and QRectF bindings for Qt 6.11 QtCore.
//
// # Value semantics
//
// Both QRect and QRectF are Q_RELOCATABLE_TYPE value types in Qt. In Go they are plain
// structs with unexported fields holding a full copy of the data. Assignment and
// function-passing copy the value; mutating one copy does not affect another. There are no
// pointers, no finalizers and no Close(). This matches QPoint / QPointF / QSize / QSizeF /
// QMargins / QMarginsF / QLine / QLineF in this package.
//
// Because Go value methods cannot mutate the receiver, "setter"-style methods return the
// updated value. Assign back to apply:
//
//	r = r.SetLeft(10)
//	r = r.Translate(dx, dy)
//	r = r.Adjust(1, 2, 3, 4)
//	r = r.MoveTo(x, y)
//
// In-place "assign" helpers (UniteAssign / IntersectAssign) are not provided, because
// value semantics make them redundant — just assign the result of United / Intersected.
//
// # Storage layout
//
// Internally:
//   - QRect  stores (x1, y1, x2, y2) — left, top, right, bottom (both endpoints inclusive).
//     width()  = right - left + 1; height() = bottom - top + 1.
//   - QRectF stores (xp, yp, w, h)   — left, top, width, height.
//     right()  = xp + w; bottom() = yp + h.
//
// This is why QRect(10, 20, 100, 50) and QRectF(10, 20, 100, 50) produce consistent
// left()/top()/width()/height() but differ in right()/bottom():
//   - QRect(10, 20, 100, 50): right = 10+100-1 = 109; bottom = 20+50-1 = 69.
//   - QRectF(10, 20, 100, 50): right = 10+100 = 110; bottom = 20+50 = 70.
//
// SetWidth(w) and SetRect(x,y,w,h) follow this convention. Use SetCoords(x1,y1,x2,y2) if you
// want to set inclusive corners directly.
//
// # Default constructors
//
// Qt's default constructors differ:
//
//	QRect()  → (0, 0, -1, -1)  // null, empty, invalid
//	QRectF() → (0, 0,  0,  0)  // null, empty, invalid
//
// The Go zero value QRect{} = (0, 0, 0, 0) is NOT the same as QRect() — it is a valid
// zero-sized rect (IsNull false). Use NewQRectDefault() for the Qt default. For QRectF the
// Go zero value QRectF{} is exactly QRectF().
//
// # Operator mapping (QRect / QRectF)
//
// Because Go has no operator overloading, Qt operators map to methods:
//
//	a | b        → a.United(b)              (bounding rect — new value)
//	a & b        → a.Intersected(b)         (intersection — new value)
//	a + margins  → a.MarginsAdded(m)        (m is QMargins / QMarginsF)
//	a - margins  → a.MarginsRemoved(m)
//	a == b       → a.Equal(b) / a.NotEqual(b)     (QRectF: fuzzy)
//
// QRect has no operator+=/operator-= with QMargins as member functions in Qt 6.11 — only
// free operator+/operator-. We expose those as MarginsAdded/MarginsRemoved.
//
// # QRect ↔ QRectF
//
//   - (QRect).ToRectF()          — mirrors QRect::toRectF()
//   - (QRectF).ToRect()          — mirrors QRectF::toRect()  (special rounding)
//   - (QRectF).ToAlignedRect()   — mirrors QRectF::toAlignedRect()  (floor/ceil)
//   - QRectFFromRect(r)          — mirrors QRectF(const QRect&)
//   - (QRectF).EqualRect(r)      — mirrors comparesEqual(QRectF, QRect)
//
// # QRect.Normalized() vs QRectF.Normalized()
//
// Both preserve |width| and |height|, but produce different corner values because
// QRect uses inclusive coordinates (width = right-left+1) while QRectF uses
// continuous ones (width = right-left):
//
//	QRect(10, 20, 0, 5).SetCoords()  →  Normalized() → (1, 6, 9, 19)
//	QRectF(10, 20, 0, 5).SetCoords() →  Normalized() → (0, 5, 10, 20)
//
// This mirrors Qt's actual behaviour (qtbase/src/corelib/tools/qrect.cpp uses
// `x2 + 1, x1 - 1` for QRect; `xp += w, w = -w` for QRectF) and is not a bug.
//
// # Not yet implemented
//
//   - QDataStream operator<< / operator>>  — will appear with package qdatastream.
//   - QDebug operator<<                    — will appear with package qdebug.
//   - qHash(QRect)                         — user hashes coordinates directly.
//   - std::tuple protocol / structured bindings
//                                          — no Go equivalent; use Left()/Top()/Width()/Height().
//   - get<I>(QRect) / get<I>(QRectF)       — same; use accessors.
//   - QRect::toCGRect / QRectF::fromCGRect / toCGRect (Darwin)
//                                          — platform-specific.
//   - QRectF::fromDOMRect / toDOMRect (WASM) — platform-specific.
//   - Free operator+(QRect, QMargins) / operator+(QMargins, QRect)
//                                          — covered by MarginsAdded (commutative in effect).
//   - Free operator-(QRect, QMargins)      — covered by MarginsRemoved.
//   - qFuzzyCompare(QRectF, QRectF) / qFuzzyIsNull(QRectF)
//                                          — covered indirectly by Equal() / IsNull().
//   - Direct field access (r.Width, r.Height, …)
//                                          — Go does not allow a field and a method with the
//                                            same name; use the accessor methods instead.
//
// # Future work
//
// Future additions will include:
//
//   - QPolygon / QPolygonF (vector of QPoint / QPointF).
//   - QTransform (2D transformations).
//   - QDataStream serialization once qdatastream exists.
//   - QDebug formatting once qdebug exists.
//
// Requires Qt >= 6.11, C++17, pkg-config Qt6Core.
package qgeometry
