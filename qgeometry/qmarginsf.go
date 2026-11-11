package qgeometry

/*
 QMarginsF — биндинг Qt 6.11 QtCore/QMarginsF как HANDLE.

 QMarginsF — value-тип в Qt (Q_RELOCATABLE_TYPE, 2×qreal… точнее 4×qreal). В Go хранится handle
 на heap-размещённый C++ QMarginsF.

 ВАЖНОЕ ОГРАНИЧЕНИЕ:
   Qt хранит QMarginsF внутри других классов ПО ЗНАЧЕНИЮ. Handle на QMarginsF не видит
   изменений, сделанных C++ во встроенном QMarginsF другого класса.

 Семантика значений vs ссылок:
   - QMarginsF — ссылочный тип (handle). Для независимой копии — Clone().
   - Go-нулевое значение QMarginsF{} невалидно. Используйте NewQMarginsF / NewQMarginsFZero.

 Управление памятью:
   - NewQMarginsF / QMarginsFFromMargins / арифметика / Clone / ToMargins — новый heap-объект.
   - Close() освобождает немедленно; иначе срабатывает finalizer.

 Соглашения:
   - Геттеры: Left(), Top(), Right(), Bottom(), IsNull() (fuzzy).
   - Мутирующие (in place, void): SetLeft, SetTop, SetRight, SetBottom.
   - Арифметика (новый handle): Add, Sub, AddFloat, SubFloat, MulFloat, DivFloat, Or, Neg,
     UnaryPlus.
   - Присваивающие формы (in place): AddAssign, SubAssign, AddFloatAssign, SubFloatAssign,
     MulFloatAssign, DivFloatAssign.
   - Сравнение: Equal(other *QMarginsF) — fuzzy; NotEqual; EqualQMargins(other *QMargins) — кросс.
   - Конвертация: ToMargins() → *QMargins; QMarginsFFromMargins(m *QMargins) — конструктор.
   - Клонирование: Clone() → *QMarginsF.

 Особенности:
   - Равенство fuzzy: два QMarginsF с крошечной разницей могут быть "равны".
   - DivFloat(0): в release — inf/nan; в debug — Q_ASSERT. Go — IEEE 754 inf/nan без паники.
   - ToMargins() использует qSaturateRound.
   - Protected-секций у QMarginsF нет.

 Что не реализовано (см. qmargins_doc.go):

   - QDataStream operator<< / operator>> — пакет qdatastream пока не готов.
   - QDebug operator<< — пакет qdebug пока не готов.
   - std::tuple protocol — не нужно в Go.
   - qFuzzyCompare / qFuzzyIsNull как свободные функции — покрыты Equal/IsNull.
*/

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC -Wno-sfinae-incomplete
#cgo LDFLAGS: -lstdc++
#include "qmargins_c.h"
*/
import "C"

import "runtime"

// QMarginsF wraps a heap-allocated Qt6 QMarginsF.
type QMarginsF struct {
	p C.go_qmarginsf_t
}

func wrapQMarginsF(p C.go_qmarginsf_t) *QMarginsF {
	if p == nil {
		return nil
	}
	m := &QMarginsF{p: p}
	runtime.SetFinalizer(m, finalizeQMarginsF)
	return m
}

func finalizeQMarginsF(m *QMarginsF) {
	if m.p != nil {
		C.qmarginsf_delete(m.p)
		m.p = nil
	}
}

// Close releases the underlying C++ object immediately and disables the finalizer.
func (m *QMarginsF) Close() {
	if m.p != nil {
		C.qmarginsf_delete(m.p)
		m.p = nil
		runtime.SetFinalizer(m, nil)
	}
}

// ---- construction --------------------------------------------------------

// NewQMarginsF returns a handle to QMarginsF(left, top, right, bottom).
func NewQMarginsF(left, top, right, bottom float64) *QMarginsF {
	return wrapQMarginsF(C.qmarginsf_new(
		C.double(left), C.double(top), C.double(right), C.double(bottom)))
}

// NewQMarginsFZero returns a handle to QMarginsF(0, 0, 0, 0) — same as Qt's default QMarginsF().
func NewQMarginsFZero() *QMarginsF {
	return wrapQMarginsF(C.qmarginsf_new_default())
}

// QMarginsFFromMargins builds a QMarginsF handle from QMargins (exact int→double).
func QMarginsFFromMargins(m *QMargins) *QMarginsF {
	return wrapQMarginsF(C.qmarginsf_from_margins(m.p))
}

// Clone returns a heap copy of m (independent C++ object).
func (m *QMarginsF) Clone() *QMarginsF {
	return wrapQMarginsF(C.qmarginsf_clone(m.p))
}

// ---- accessors -----------------------------------------------------------

func (m *QMarginsF) IsNull() bool    { return bool(C.qmarginsf_is_null(m.p)) }
func (m *QMarginsF) Left() float64   { return float64(C.qmarginsf_left(m.p)) }
func (m *QMarginsF) Top() float64    { return float64(C.qmarginsf_top(m.p)) }
func (m *QMarginsF) Right() float64  { return float64(C.qmarginsf_right(m.p)) }
func (m *QMarginsF) Bottom() float64 { return float64(C.qmarginsf_bottom(m.p)) }

// ---- mutating setters (in place) ----------------------------------------

func (m *QMarginsF) SetLeft(v float64)   { C.qmarginsf_set_left(m.p, C.double(v)) }
func (m *QMarginsF) SetTop(v float64)    { C.qmarginsf_set_top(m.p, C.double(v)) }
func (m *QMarginsF) SetRight(v float64)  { C.qmarginsf_set_right(m.p, C.double(v)) }
func (m *QMarginsF) SetBottom(v float64) { C.qmarginsf_set_bottom(m.p, C.double(v)) }

// ---- arithmetic (return new handles) ------------------------------------

func (m *QMarginsF) Add(other *QMarginsF) *QMarginsF {
	return wrapQMarginsF(C.qmarginsf_add(m.p, other.p))
}
func (m *QMarginsF) Sub(other *QMarginsF) *QMarginsF {
	return wrapQMarginsF(C.qmarginsf_sub(m.p, other.p))
}
func (m *QMarginsF) AddFloat(v float64) *QMarginsF {
	return wrapQMarginsF(C.qmarginsf_add_float(m.p, C.double(v)))
}
func (m *QMarginsF) SubFloat(v float64) *QMarginsF {
	return wrapQMarginsF(C.qmarginsf_sub_float(m.p, C.double(v)))
}
func (m *QMarginsF) MulFloat(f float64) *QMarginsF {
	return wrapQMarginsF(C.qmarginsf_mul_float(m.p, C.double(f)))
}
func (m *QMarginsF) DivFloat(d float64) *QMarginsF {
	return wrapQMarginsF(C.qmarginsf_div_float(m.p, C.double(d)))
}
func (m *QMarginsF) Or(other *QMarginsF) *QMarginsF {
	return wrapQMarginsF(C.qmarginsf_or(m.p, other.p))
}
func (m *QMarginsF) Neg() *QMarginsF       { return wrapQMarginsF(C.qmarginsf_neg(m.p)) }
func (m *QMarginsF) UnaryPlus() *QMarginsF { return wrapQMarginsF(C.qmarginsf_unary_plus(m.p)) }

// ---- in-place assign (mutate receiver) ----------------------------------

func (m *QMarginsF) AddAssign(other *QMarginsF) { C.qmarginsf_add_assign(m.p, other.p) }
func (m *QMarginsF) SubAssign(other *QMarginsF) { C.qmarginsf_sub_assign(m.p, other.p) }
func (m *QMarginsF) AddFloatAssign(v float64)   { C.qmarginsf_add_float_assign(m.p, C.double(v)) }
func (m *QMarginsF) SubFloatAssign(v float64)   { C.qmarginsf_sub_float_assign(m.p, C.double(v)) }
func (m *QMarginsF) MulFloatAssign(f float64)   { C.qmarginsf_mul_float_assign(m.p, C.double(f)) }
func (m *QMarginsF) DivFloatAssign(d float64)   { C.qmarginsf_div_float_assign(m.p, C.double(d)) }

// ---- comparison ----------------------------------------------------------

// Equal returns true if both are fuzzy-equal (Qt comparesEqual → qFuzzyCompare).
func (m *QMarginsF) Equal(other *QMarginsF) bool { return bool(C.qmarginsf_equal(m.p, other.p)) }

// NotEqual returns !Equal(other).
func (m *QMarginsF) NotEqual(other *QMarginsF) bool { return !m.Equal(other) }

// EqualQMargins returns true if m is fuzzy-equal to other (as QMarginsF).
func (m *QMarginsF) EqualQMargins(other *QMargins) bool {
	return bool(C.qmarginsf_equal_margins(m.p, other.p))
}

// ---- conversions ---------------------------------------------------------

// ToMargins converts m to a new QMargins handle (qSaturateRound per component).
func (m *QMarginsF) ToMargins() *QMargins {
	return wrapQMargins(C.qmarginsf_to_margins(m.p))
}
