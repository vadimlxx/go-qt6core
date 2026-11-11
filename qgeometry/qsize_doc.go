// Package qgeometry provides Go bindings for Qt 6.11 QtCore geometry value types.
//
// # Currently supported
//
//   - QMargins  — 4×int margins (left, top, right, bottom)
//   - QMarginsF — 4×qreal margins
//   - QSize     — 2×int size (width, height)
//   - QSizeF    — 2×qreal size
//
// All four are Q_RELOCATABLE_TYPE value types in Qt. In Go they are structs
// with unexported fields and Qt-style accessor methods (Left()/Top()/… for
// margins, Width()/Height() for sizes). Construction goes through the NewXxx
// constructors; mutation through SetXxx / Assign methods. Passing through the
// C ABI uses flat 2- or 4-field structs (go_qmargins_t, go_qmarginsf_t,
// go_qsize_t, go_qsizef_t) in C mode, typedef'd to the real QMargins /
// QMarginsF / QSize / QSizeF in C++ mode, so the layout stays ABI-compatible
// with Qt.
//
// # QSize / QSizeF specifics
//
// Qt's default constructors are non-zero:
//
//	QSize()  → QSize(-1, -1)     // invalid
//	QSizeF() → QSizeF(-1., -1.)  // invalid
//
// In Go this is mirrored by NewQSizeDefault() / NewQSizeFDefault(). For
// explicit zero values use NewQSizeZero() / NewQSizeFZero() or NewQSize(0, 0)
// / NewQSizeF(0, 0).
//
// Aspect ratio modes are exposed as the AspectRatioMode type with three
// constants matching Qt's Qt::AspectRatioMode:
//
//	IgnoreAspectRatio            — set to requested size, ignore aspect
//	KeepAspectRatio              — fit inside, keep aspect
//	KeepAspectRatioByExpanding   — cover, keep aspect
//
// QSize and QSizeF have no protected members — only public API and private
// implementation. Nothing needs to be exposed via a helper subclass.
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
// # Not yet implemented
//
//   - QDataStream operator<< / operator>>  — will appear with package qdatastream.
//   - QDebug operator<<                    — will appear with package qdebug.
//   - qHash(QSize) / qHash(QSizeF)         — user can hash Width()/Height() directly.
//   - rwidth() / rheight()                 — covered by SetWidth / SetHeight.
//   - std::tuple protocol / structured bindings
//                                          — no Go equivalent; access via Width()/Height().
//   - get<I>(QSize) / get<I>(QSizeF)       — same; access via Width()/Height().
//   - QSize::toCGSize / QSizeF::fromCGSize / toCGSize (Darwin)
//                                          — platform-specific.
//   - Direct field access (s.Width, s.Height, m.Left, …)
//                                          — Go does not allow a field and a method with the
//                                            same name; use the accessor methods instead.
//
// # Future work
//
// Future additions will include:
//
//   - QRect / QRectF / QPoint / QPointF in this same package, with cross-type
//     conversions (e.g. QRect::size(), QRect::margins(), QSize::grownBy()).
//   - QDataStream serialization once qdatastream exists.
//   - QDebug formatting once qdebug exists.
//
// This package is a subpackage of github.com/vadimlxx/go-qt6core.
// Requires Qt >= 6.11, C++17, pkg-config Qt6Core.
package qgeometry