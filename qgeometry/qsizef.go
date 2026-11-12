package qgeometry

/*
 QSizeF — биндинг Qt 6.11 QtCore/QSizeF.

 QSizeF — value-тип (Q_RELOCATABLE_TYPE): 2 qreal-поля (width, height). В Go поля неэкспортируемые,
 доступ через геттеры Width()/Height(). Через C ABI передаётся плоской структурой go_qsizef_t
 (2×double, 16 байт).

 ВАЖНО: Qt-конструктор QSizeF() возвращает QSizeF(-1., -1.), НЕ (0., 0.).
   - NewQSizeFDefault() → (-1., -1.) — соответствует QSizeF()
   - NewQSizeFZero()    → ( 0.,  0.) — «нулевой» QSizeF
   - NewQSizeF(w, h)    → ( w,   h)  — явный

 Соглашения:

   - Геттеры: Width(), Height(), IsNull(), IsEmpty(), IsValid().
     IsNull() использует qIsNull → qFuzzyIsNull; IsEmpty() — wd <= 0 || ht <= 0;
     IsValid() — wd >= 0 && ht >= 0.
   - Мутирующие сеттеры: SetWidth(w), SetHeight(h), Transpose(), Scale(w, h, mode),
     ScaleSize(sz, mode).
   - Трансформации (возвращают новый QSizeF): Transposed(), Scaled(w, h, mode), ScaledSize(sz, mode),
     ExpandedTo(other), BoundedTo(other), GrownBy(QMarginsF), ShrunkBy(QMarginsF).
   - Арифметика: Add(other), Sub(other), MulFloat(c), DivFloat(c).
   - Присваивающие формы: AddAssign, SubAssign, MulFloatAssign, DivFloatAssign.
   - Сравнение: Equal(other) — fuzzy (comparesEqual → qFuzzyCompare); NotEqual(other);
     EqualQSize(m) — сравнение с QSize (comparesEqual(QSizeF, QSize)).
   - Конвертация: ToSize() → QSize (qSaturateRound).

 Особенности:

   - QSizeF — fuzzy-compared. Два QSizeF с крошечной разницей в double могут быть "равны". Это Qt.
   - DivFloat(0): в release — inf/nan; в debug — Q_ASSERT. Go — IEEE 754 inf/nan без паники.
   - ToSize() использует qSaturateRound — округление и saturating к int-диапазону.
   - Protected-секций у QSizeF нет — открывать нечего.

 Что не реализовано (см. qsize_doc.go):

   - QDataStream operator<< / operator>> — пакет qdatastream пока не готов.
   - QDebug operator<< — пакет qdebug пока не готов.
   - qHash(QSizeF) — не нужно в Go.
   - rwidth()/rheight() — покрыто SetWidth/SetHeight.
   - std::tuple protocol — не нужно в Go.
   - fromCGSize / toCGSize (Darwin) — платформо-специфично.
   - qFuzzyCompare(QSizeF, QSizeF) / qFuzzyIsNull(QSizeF) как свободные функции —
     покрыты через Equal / IsNull.
*/

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC -Wno-sfinae-incomplete
#cgo LDFLAGS: -lstdc++
#include "qsize_c.h"
*/
import "C"

import "structs"

// QSizeF wraps Qt6 QSizeF (Q_RELOCATABLE_TYPE, 2×qreal value).
// Fields are unexported; use Width()/Height() to read, SetWidth()/SetHeight() to mutate.
type QSizeF struct {
	_      structs.HostLayout
	width  float64
	height float64
}

// ---- construction --------------------------------------------------------

// NewQSizeF returns QSizeF(w, h).
func NewQSizeF(w, h float64) QSizeF {
	return qsizefFromC(C.qsizef_new(C.double(w), C.double(h)))
}

// NewQSizeFZero returns QSizeF(0., 0.).
func NewQSizeFZero() QSizeF { return NewQSizeF(0, 0) }

// NewQSizeFDefault returns QSizeF(-1., -1.), matching Qt's default constructor QSizeF().
func NewQSizeFDefault() QSizeF { return qsizefFromC(C.qsizef_default()) }

// QSizeFFromSize builds QSizeF from QSize (exact int→double). Mirrors QSizeF(const QSize&).
func QSizeFFromSize(s QSize) QSizeF { return qsizefFromC(C.qsizef_from_size(qsizeToC(s))) }

// ---- accessors -----------------------------------------------------------

func (s QSizeF) IsNull() bool    { return bool(C.qsizef_is_null(qsizefToC(s))) }
func (s QSizeF) IsEmpty() bool   { return bool(C.qsizef_is_empty(qsizefToC(s))) }
func (s QSizeF) IsValid() bool   { return bool(C.qsizef_is_valid(qsizefToC(s))) }
func (s QSizeF) Width() float64  { return float64(C.qsizef_width(qsizefToC(s))) }
func (s QSizeF) Height() float64 { return float64(C.qsizef_height(qsizefToC(s))) }

// ---- modifiers -----------------------------------------------------------

func (s *QSizeF) SetWidth(w float64) {
	*s = qsizefFromC(C.qsizef_set_width(qsizefToC(*s), C.double(w)))
}
func (s *QSizeF) SetHeight(h float64) {
	*s = qsizefFromC(C.qsizef_set_height(qsizefToC(*s), C.double(h)))
}

// Transpose swaps width and height in place (mirrors QSizeF::transpose()).
func (s *QSizeF) Transpose() { *s = s.Transposed() }

// Scale scales s in place to fit (w, h) with the given aspect ratio mode
// (mirrors QSizeF::scale(qreal, qreal, Qt::AspectRatioMode)).
func (s *QSizeF) Scale(w, h float64, mode AspectRatioMode) { *s = s.Scaled(w, h, mode) }

// ScaleSize scales s in place to fit other with the given aspect ratio mode
// (mirrors QSizeF::scale(const QSizeF&, Qt::AspectRatioMode)).
func (s *QSizeF) ScaleSize(other QSizeF, mode AspectRatioMode) { *s = s.ScaledSize(other, mode) }

// ---- transforms (return new QSizeF) -------------------------------------

// Transposed returns a QSizeF with width and height swapped.
func (s QSizeF) Transposed() QSizeF { return qsizefFromC(C.qsizef_transposed(qsizefToC(s))) }

// Scaled returns a new QSizeF scaled to fit (w, h) with the given aspect ratio mode.
func (s QSizeF) Scaled(w, h float64, mode AspectRatioMode) QSizeF {
	return qsizefFromC(C.qsizef_scaled(qsizefToC(s), C.double(w), C.double(h), C.int(mode)))
}

// ScaledSize returns a new QSizeF scaled to fit other with the given aspect ratio mode.
func (s QSizeF) ScaledSize(other QSizeF, mode AspectRatioMode) QSizeF {
	return qsizefFromC(C.qsizef_scaled_size(qsizefToC(s), qsizefToC(other), C.int(mode)))
}

// ExpandedTo returns the component-wise max of s and other (mirrors QSizeF::expandedTo).
func (s QSizeF) ExpandedTo(other QSizeF) QSizeF {
	return qsizefFromC(C.qsizef_expanded_to(qsizefToC(s), qsizefToC(other)))
}

// BoundedTo returns the component-wise min of s and other (mirrors QSizeF::boundedTo).
func (s QSizeF) BoundedTo(other QSizeF) QSizeF {
	return qsizefFromC(C.qsizef_bounded_to(qsizefToC(s), qsizefToC(other)))
}

// GrownBy grows each dimension by the corresponding QMarginsF components (mirrors QSizeF::grownBy).
func (s QSizeF) GrownBy(m QMarginsF) QSizeF {
	return qsizefFromC(C.qsizef_grown_by(qsizefToC(s), marginsFToC(m)))
}

// ShrunkBy shrinks each dimension by the corresponding QMarginsF components (mirrors QSizeF::shrunkBy).
func (s QSizeF) ShrunkBy(m QMarginsF) QSizeF {
	return qsizefFromC(C.qsizef_shrunk_by(qsizefToC(s), marginsFToC(m)))
}

// ---- arithmetic ----------------------------------------------------------

// Add returns s + other (component-wise).
func (s QSizeF) Add(other QSizeF) QSizeF {
	return qsizefFromC(C.qsizef_add(qsizefToC(s), qsizefToC(other)))
}

// Sub returns s - other (component-wise).
func (s QSizeF) Sub(other QSizeF) QSizeF {
	return qsizefFromC(C.qsizef_sub(qsizefToC(s), qsizefToC(other)))
}

// MulFloat returns s * c on each component.
func (s QSizeF) MulFloat(c float64) QSizeF {
	return qsizefFromC(C.qsizef_mul_float(qsizefToC(s), C.double(c)))
}

// DivFloat returns s / c on each component.
// IEEE 754 semantics: c == 0 yields inf/nan (Qt asserts only in debug builds).
func (s QSizeF) DivFloat(c float64) QSizeF {
	return qsizefFromC(C.qsizef_div_float(qsizefToC(s), C.double(c)))
}

// ---- assign forms --------------------------------------------------------

func (s *QSizeF) AddAssign(other QSizeF)   { *s = s.Add(other) }
func (s *QSizeF) SubAssign(other QSizeF)   { *s = s.Sub(other) }
func (s *QSizeF) MulFloatAssign(c float64) { *s = s.MulFloat(c) }
func (s *QSizeF) DivFloatAssign(c float64) { *s = s.DivFloat(c) }

// ---- comparison ----------------------------------------------------------

// Equal returns true if both components are fuzzy-equal (Qt comparesEqual → qFuzzyCompare).
func (s QSizeF) Equal(other QSizeF) bool {
	return bool(C.qsizef_equal(qsizefToC(s), qsizefToC(other)))
}

// NotEqual returns !Equal(other).
func (s QSizeF) NotEqual(other QSizeF) bool { return !s.Equal(other) }

// EqualQSize returns true if s is fuzzy-equal to other (as QSizeF). Mirrors comparesEqual(QSizeF, QSize).
func (s QSizeF) EqualQSize(other QSize) bool {
	return bool(C.qsizef_equal_size(qsizefToC(s), qsizeToC(other)))
}

// ---- conversions ---------------------------------------------------------

// ToSize converts s to QSize via qSaturateRound on each component. Mirrors QSizeF::toSize().
func (s QSizeF) ToSize() QSize { return qsizeFromC(C.qsizef_to_size(qsizefToC(s))) }

// ---- internal helpers ----------------------------------------------------

func qsizefToC(s QSizeF) C.go_qsizef_t {
	return C.go_qsizef_t{width: C.double(s.width), height: C.double(s.height)}
}

func qsizefFromC(c C.go_qsizef_t) QSizeF {
	return QSizeF{width: float64(c.width), height: float64(c.height)}
}
