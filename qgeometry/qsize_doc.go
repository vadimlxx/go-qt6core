// Package qgeometry — QSize and QSizeF bindings for Qt 6.11 QtCore.
//
// # Handle semantics
//
// Both QSize and QSizeF are exposed as HANDLES (Go structs holding a C++ pointer). All accessors
// and mutators go directly to the heap-allocated C++ object:
//
//	s := NewQSize(100, 50)     // allocates a C++ QSize
//	w := s.Width()              // s->width() — reads from C++
//	s.SetWidth(200)             // s->setWidth(200) — mutates the C++ object
//	s.Close()                   // deletes the C++ QSize
//
// To obtain a value-semantics copy, use Clone(). To share mutations, pass the same *QSize around
// — all handles referring to the same C++ object see the same state.
//
// Memory management:
//
//   - NewQSize / NewQSizeF / NewQSizeZero / NewQSizeFZero / NewQSizeDefault / NewQSizeFDefault /
//     arithmetic results / Clone / ToSizeF / ToSize / Transposed / Scaled / ExpandedTo / BoundedTo /
//     GrownBy / ShrunkBy — allocate a new heap C++ object.
//   - The finalizer frees it; Close() frees immediately and disables the finalizer.
//   - Do not copy *QSize by value expecting independent copies — copy the pointer; use Clone()
//     for a truly independent copy.
//
// # Default constructors
//
// Qt's default constructors are NON-zero:
//
//	QSize()  → QSize(-1, -1)     // invalid
//	QSizeF() → QSizeF(-1., -1.)  // invalid
//
// In Go this is mirrored by NewQSizeDefault() / NewQSizeFDefault(). For explicit zero values use
// NewQSizeZero() / NewQSizeFZero() or NewQSize(0, 0) / NewQSizeF(0, 0).
//
// # Aspect ratio modes
//
// Aspect ratio modes are exposed as the AspectRatioMode type with three constants matching
// Qt's Qt::AspectRatioMode:
//
//	IgnoreAspectRatio            — set to requested size, ignore aspect
//	KeepAspectRatio              — fit inside, keep aspect
//	KeepAspectRatioByExpanding   — cover, keep aspect
//
// # Fundamental limitation with embedded QSize / QSizeF
//
// Qt stores QSize / QSizeF **by value** inside other classes. A Go *QSize handle CANNOT
// reference such internal storage; it always owns a separate heap copy. If C++ code later
// mutates an embedded QSize, the previously-obtained *QSize will NOT see the change. Use the
// containing class's fresh accessors to read the current state. This is a limitation of Qt's
// public API, not of the binding.
//
// # Operator mapping (QSize / QSizeF)
//
//	a + b        → a.Add(b)
//	a - b        → a.Sub(b)
//	a * f        → a.MulFloat(f)       (qreal only — no integer variant in Qt)
//	a / f        → a.DivFloat(f)       (qreal only)
//	a == b       → a.Equal(b) / a.NotEqual(b)   (QSize: exact; QSizeF: fuzzy)
//	a += b       → a.AddAssign(b)
//	a -= b       → a.SubAssign(b)
//	a *= f       → a.MulFloatAssign(f)
//	a /= f       → a.DivFloatAssign(f)
//
// # Protected access
//
// QSize and QSizeF have no protected members — only public API and private implementation.
// Nothing needs to be exposed via a helper subclass.
//
// # Not yet implemented
//
//   - QDataStream operator<< / operator>>  — will appear with package qdatastream.
//   - QDebug operator<<                    — will appear with package qdebug.
//   - qHash(QSize) / qHash(QSizeF)         — user can hash Width()/Height() directly.
//   - rwidth() / rheight()                 — covered by SetWidth / SetHeight.
//   - std::tuple protocol / structured bindings
//                                          — no Go equivalent; use Width()/Height().
//   - get<I>(QSize) / get<I>(QSizeF)       — same; use Width()/Height().
//   - QSize::toCGSize / QSizeF::fromCGSize / toCGSize (Darwin)
//                                          — platform-specific.
//   - Free functions qFuzzyCompare(QSizeF, QSizeF) and qFuzzyIsNull(QSizeF)
//                                          — covered indirectly by Equal() / IsNull().
//   - Direct field access (s.Width, s.Height, m.Left, …)
//                                          — Go does not allow a field and a method with the
//                                            same name; use the accessor methods instead.
//
// # Future work
//
// Future additions will include:
//
//   - QRect / QRectF in this same package, with cross-type conversions
//     (e.g. QRect::size(), QRect::margins(), QRect::topLeft()).
//   - QDataStream serialization once qdatastream exists.
//   - QDebug formatting once qdebug exists.
//
// Requires Qt >= 6.11, C++17, pkg-config Qt6Core.
package qgeometry
