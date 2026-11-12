package qgeometry

/*
 QMarginsF — биндинг Qt 6.11 QtCore/QMarginsF.

 QMarginsF — value-тип (Q_RELOCATABLE_TYPE): 4 qreal-поля (left, top, right,
 bottom). В Go поля неэкспортируемые, доступ через геттеры Left()/Top()/
 Right()/Bottom(). Через C ABI передаётся плоской структурой go_qmarginsf_t
 (4×double, 32 байта).

 Соглашения:

   - Геттеры: Left(), Top(), Right(), Bottom(), IsNull()
     (IsNull — qFuzzyIsNull на каждой компоненте, Qt semantics).
   - Мутирующие сеттеры: SetLeft(v), SetTop(v), SetRight(v), SetBottom(v).
   - Арифметика: Add(other), Sub(other), AddFloat(v), SubFloat(v),
     MulFloat(f), DivFloat(d), Neg(), UnaryPlus().
   - Компонентный максимум: Or(other).
   - Присваивающие формы: AddAssign, SubAssign, AddFloatAssign,
     SubFloatAssign, MulFloatAssign, DivFloatAssign.
   - Сравнение: Equal(other) — fuzzy (Qt comparesEqual → qFuzzyCompare);
     EqualQMargins(m) — сравнение с QMargins через неявный toMarginsF().
   - Конвертация: ToMargins() → QMargins (qSaturateRound).

 Особенности:

   - Равенство fuzzy: два QMarginsF с крошечной разницей в double могут
     быть "равны". Это соответствует Qt.
   - DivFloat(d) с d ≈ 0: в release-сборке Qt даёт inf/nan, в debug —
     Q_ASSERT. В Go результат — стандартный IEEE 754 (inf/nan), без паники.
   - ToMargins() использует qSaturateRound — округление и saturating
     к int-диапазону, идентично Qt.
   - Protected-секций у QMarginsF нет — открывать нечего.

 Что не реализовано (см. qmargins_doc.go):

   - QDataStream operator<< / operator>> — пакет qdatastream пока не готов.
   - QDebug operator<< — пакет qdebug пока не готов.
   - std::tuple protocol — не нужно в Go.
   - qFuzzyCompare / qFuzzyIsNull как свободные функции —
     покрыты через Equal / IsNull.
*/

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC -Wno-sfinae-incomplete
#cgo LDFLAGS: -lstdc++
#include "qmargins_c.h"
*/
import "C"

import "structs"

// QMarginsF wraps Qt6 QMarginsF (Q_RELOCATABLE_TYPE, 4×qreal value).
// Fields are unexported; use Left()/Top()/Right()/Bottom() to read,
// SetLeft()/… to mutate.
type QMarginsF struct {
	_      structs.HostLayout
	left   float64
	top    float64
	right  float64
	bottom float64
}

// ---- construction --------------------------------------------------------

// NewQMarginsF returns QMarginsF(left, top, right, bottom).
func NewQMarginsF(left, top, right, bottom float64) QMarginsF {
	return marginsFFromC(C.qmarginsf_new(C.double(left), C.double(top), C.double(right), C.double(bottom)))
}

// NewQMarginsFZero returns the default-constructed QMarginsF (0, 0, 0, 0).
func NewQMarginsFZero() QMarginsF {
	return marginsFFromC(C.qmarginsf_default())
}

// QMarginsFFromMargins builds QMarginsF from QMargins (exact int→double).
// Mirrors QMarginsF(const QMargins&).
func QMarginsFFromMargins(m QMargins) QMarginsF {
	return marginsFFromC(C.qmarginsf_from_margins(marginsToC(m)))
}

// ---- accessors -----------------------------------------------------------

func (m QMarginsF) IsNull() bool    { return bool(C.qmarginsf_is_null(marginsFToC(m))) }
func (m QMarginsF) Left() float64   { return float64(C.qmarginsf_left(marginsFToC(m))) }
func (m QMarginsF) Top() float64    { return float64(C.qmarginsf_top(marginsFToC(m))) }
func (m QMarginsF) Right() float64  { return float64(C.qmarginsf_right(marginsFToC(m))) }
func (m QMarginsF) Bottom() float64 { return float64(C.qmarginsf_bottom(marginsFToC(m))) }

// ---- modifiers -----------------------------------------------------------

func (m *QMarginsF) SetLeft(v float64) {
	*m = marginsFFromC(C.qmarginsf_set_left(marginsFToC(*m), C.double(v)))
}
func (m *QMarginsF) SetTop(v float64) {
	*m = marginsFFromC(C.qmarginsf_set_top(marginsFToC(*m), C.double(v)))
}
func (m *QMarginsF) SetRight(v float64) {
	*m = marginsFFromC(C.qmarginsf_set_right(marginsFToC(*m), C.double(v)))
}
func (m *QMarginsF) SetBottom(v float64) {
	*m = marginsFFromC(C.qmarginsf_set_bottom(marginsFToC(*m), C.double(v)))
}

// ---- arithmetic ----------------------------------------------------------

// Add returns m + other (component-wise).
func (m QMarginsF) Add(other QMarginsF) QMarginsF {
	return marginsFFromC(C.qmarginsf_add(marginsFToC(m), marginsFToC(other)))
}

// Sub returns m - other (component-wise).
func (m QMarginsF) Sub(other QMarginsF) QMarginsF {
	return marginsFFromC(C.qmarginsf_sub(marginsFToC(m), marginsFToC(other)))
}

// AddFloat returns m + v on each component.
func (m QMarginsF) AddFloat(v float64) QMarginsF {
	return marginsFFromC(C.qmarginsf_add_float(marginsFToC(m), C.double(v)))
}

// SubFloat returns m - v on each component.
func (m QMarginsF) SubFloat(v float64) QMarginsF {
	return marginsFFromC(C.qmarginsf_sub_float(marginsFToC(m), C.double(v)))
}

// MulFloat returns m * f on each component.
func (m QMarginsF) MulFloat(f float64) QMarginsF {
	return marginsFFromC(C.qmarginsf_mul_float(marginsFToC(m), C.double(f)))
}

// DivFloat returns m / d on each component.
// IEEE 754 semantics: d == 0 yields inf/nan (Qt asserts only in debug builds).
func (m QMarginsF) DivFloat(d float64) QMarginsF {
	return marginsFFromC(C.qmarginsf_div_float(marginsFToC(m), C.double(d)))
}

// Or returns component-wise max of m and other (Qt operator|).
func (m QMarginsF) Or(other QMarginsF) QMarginsF {
	return marginsFFromC(C.qmarginsf_or(marginsFToC(m), marginsFToC(other)))
}

// Neg returns -m (component-wise).
func (m QMarginsF) Neg() QMarginsF { return marginsFFromC(C.qmarginsf_neg(marginsFToC(m))) }

// UnaryPlus returns +m (copy).
func (m QMarginsF) UnaryPlus() QMarginsF { return marginsFFromC(C.qmarginsf_unary_plus(marginsFToC(m))) }

// ---- assign forms --------------------------------------------------------

func (m *QMarginsF) AddAssign(other QMarginsF) { *m = m.Add(other) }
func (m *QMarginsF) SubAssign(other QMarginsF) { *m = m.Sub(other) }
func (m *QMarginsF) AddFloatAssign(v float64)  { *m = m.AddFloat(v) }
func (m *QMarginsF) SubFloatAssign(v float64)  { *m = m.SubFloat(v) }
func (m *QMarginsF) MulFloatAssign(f float64)  { *m = m.MulFloat(f) }
func (m *QMarginsF) DivFloatAssign(d float64)  { *m = m.DivFloat(d) }

// ---- comparison ----------------------------------------------------------

// Equal returns true if all four components are fuzzy-equal (Qt comparesEqual).
func (m QMarginsF) Equal(other QMarginsF) bool {
	return bool(C.qmarginsf_equal(marginsFToC(m), marginsFToC(other)))
}

// NotEqual returns !Equal(other).
func (m QMarginsF) NotEqual(other QMarginsF) bool { return !m.Equal(other) }

// EqualQMargins returns true if m is fuzzy-equal to other (as QMarginsF).
func (m QMarginsF) EqualQMargins(other QMargins) bool {
	return bool(C.qmarginsf_equal_margins(marginsFToC(m), marginsToC(other)))
}

// ---- conversions ---------------------------------------------------------

// ToMargins converts m to QMargins via qSaturateRound on each component.
// Mirrors QMarginsF::toMargins().
func (m QMarginsF) ToMargins() QMargins {
	return marginsFromC(C.qmarginsf_to_margins(marginsFToC(m)))
}

// ---- internal helpers ----------------------------------------------------

func marginsFToC(m QMarginsF) C.go_qmarginsf_t {
	return C.go_qmarginsf_t{
		left:   C.double(m.left),
		top:    C.double(m.top),
		right:  C.double(m.right),
		bottom: C.double(m.bottom),
	}
}

func marginsFFromC(c C.go_qmarginsf_t) QMarginsF {
	return QMarginsF{
		left:   float64(c.left),
		top:    float64(c.top),
		right:  float64(c.right),
		bottom: float64(c.bottom),
	}
}
