package qgeometry

/*
 QPointF — биндинг Qt 6.11 QtCore/QPointF как HANDLE.

 QPointF — value-тип в Qt (Q_PRIMITIVE_TYPE, 2×qreal). В Go хранится handle на heap-размещённый
 C++ QPointF. Все геттеры/сеттеры идут напрямую в C++.

 ВАЖНОЕ ОГРАНИЧЕНИЕ (аналогично QPoint):
   Qt хранит QPointF внутри других классов (QLineF, QRectF, …) ПО ЗНАЧЕНИЮ. Handle на
   QPointF не видит изменений, сделанных C++ во встроенном QPointF. Для получения
   актуальных значений из составных типов используйте свежие геттеры (linef.X1(), …).

 Семантика значений vs ссылок:
   - QPointF — ссылочный тип (handle). Для независимой копии используйте Clone().
   - Go-нулевое значение QPointF{} невалидно. Используйте NewQPointF/NewQPointFZero.

 Управление памятью:
   - Каждый NewQPointF / QLineF.P1() / арифметика создаёт НОВЫЙ heap QPointF.
   - Close() освобождает немедленно; иначе срабатывает finalizer.

 Соглашения:
   - Геттеры: X(), Y(), IsNull() (fuzzy — qIsNull), ManhattanLength().
   - Мутирующие: SetX(x), SetY(y).
   - Арифметика возвращает НОВЫЙ handle: Add, Sub, MulFloat, DivFloat, Neg, UnaryPlus, Transposed.
   - Присваивающие формы: AddAssign, SubAssign, MulFloatAssign, DivFloatAssign.
   - DotProduct(other) → float64.
   - Сравнение: Equal(other) — fuzzy (qFuzzyCompare), NotEqual(other);
     EqualQPoint(other *QPoint) — кросс-сравнение QPointF == QPoint.
   - Конвертация: ToPoint() → *QPoint (qSaturateRound); QPointFFromPoint(p *QPoint) — конструктор.
   - Клонирование: Clone() → *QPointF.

 Что не реализовано (см. qpoint_doc.go):

   - QDataStream operator<< / operator>> — пакет qdatastream пока не готов.
   - QDebug operator<< — пакет qdebug пока не готов.
   - qHash(QPointF) — в Qt удалён.
   - rx()/ry() — покрыто SetX/SetY.
   - std::tuple protocol — не нужно в Go.
   - fromCGPoint / toCGPoint (Darwin) — платформо-специфично.
*/

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC -Wno-sfinae-incomplete
#cgo LDFLAGS: -lstdc++
#include "qpoint_c.h"
*/
import "C"

import (
	"runtime"
)

// QPointF wraps a heap-allocated Qt6 QPointF. All operations go directly to the C++ object.
type QPointF struct {
	p C.go_qpointf_t
}

func wrapQPointF(p C.go_qpointf_t) *QPointF {
	if p == nil {
		return nil
	}
	q := &QPointF{p: p}
	runtime.SetFinalizer(q, finalizeQPointF)
	return q
}

func finalizeQPointF(q *QPointF) {
	if q.p != nil {
		C.qpointf_delete(q.p)
		q.p = nil
	}
}

// Close releases the underlying C++ object immediately and disables the finalizer.
func (q *QPointF) Close() {
	if q.p != nil {
		C.qpointf_delete(q.p)
		q.p = nil
		runtime.SetFinalizer(q, nil)
	}
}

// ---- construction --------------------------------------------------------

// NewQPointF returns a handle to QPointF(x, y).
func NewQPointF(x, y float64) *QPointF {
	return wrapQPointF(C.qpointf_new(C.double(x), C.double(y)))
}

// NewQPointFZero returns a handle to QPointF(0., 0.) — same as Qt's default QPointF().
func NewQPointFZero() *QPointF {
	return wrapQPointF(C.qpointf_new_default())
}

// QPointFFromPoint builds a QPointF handle from a QPoint handle (exact int→double).
// Mirrors QPointF(const QPoint&).
func QPointFFromPoint(p *QPoint) *QPointF {
	return wrapQPointF(C.qpointf_from_point(p.p))
}

// Clone returns a heap copy of q (independent C++ object).
func (q *QPointF) Clone() *QPointF {
	return wrapQPointF(C.qpointf_clone(q.p))
}

// ---- accessors -----------------------------------------------------------

func (q *QPointF) IsNull() bool           { return bool(C.qpointf_is_null(q.p)) }
func (q *QPointF) X() float64             { return float64(C.qpointf_x(q.p)) }
func (q *QPointF) Y() float64             { return float64(C.qpointf_y(q.p)) }
func (q *QPointF) ManhattanLength() float64 { return float64(C.qpointf_manhattan_length(q.p)) }

// ---- modifiers -----------------------------------------------------------

func (q *QPointF) SetX(x float64) { C.qpointf_set_x(q.p, C.double(x)) }
func (q *QPointF) SetY(y float64) { C.qpointf_set_y(q.p, C.double(y)) }

// ---- transforms ----------------------------------------------------------

// Transposed returns a new handle with x and y swapped.
func (q *QPointF) Transposed() *QPointF { return wrapQPointF(C.qpointf_transposed(q.p)) }

// ---- dot product ---------------------------------------------------------

// DotProduct returns q·other = q.X()*other.X() + q.Y()*other.Y().
func (q *QPointF) DotProduct(other *QPointF) float64 {
	return float64(C.qpointf_dot_product(q.p, other.p))
}

// ---- arithmetic (return new handles) ------------------------------------

func (q *QPointF) Add(other *QPointF) *QPointF   { return wrapQPointF(C.qpointf_add(q.p, other.p)) }
func (q *QPointF) Sub(other *QPointF) *QPointF   { return wrapQPointF(C.qpointf_sub(q.p, other.p)) }
func (q *QPointF) MulFloat(f float64) *QPointF   { return wrapQPointF(C.qpointf_mul_float(q.p, C.double(f))) }
func (q *QPointF) DivFloat(d float64) *QPointF   { return wrapQPointF(C.qpointf_div_float(q.p, C.double(d))) }
func (q *QPointF) Neg() *QPointF                 { return wrapQPointF(C.qpointf_neg(q.p)) }
func (q *QPointF) UnaryPlus() *QPointF           { return wrapQPointF(C.qpointf_unary_plus(q.p)) }

// ---- in-place assign ----------------------------------------------------

func (q *QPointF) AddAssign(other *QPointF)  { C.qpointf_add_assign(q.p, other.p) }
func (q *QPointF) SubAssign(other *QPointF)  { C.qpointf_sub_assign(q.p, other.p) }
func (q *QPointF) MulFloatAssign(f float64)  { C.qpointf_mul_float_assign(q.p, C.double(f)) }
func (q *QPointF) DivFloatAssign(d float64)  { C.qpointf_div_float_assign(q.p, C.double(d)) }

// ---- comparison ----------------------------------------------------------

// Equal returns true if both components are fuzzy-equal (qFuzzyCompare per component).
func (q *QPointF) Equal(other *QPointF) bool { return bool(C.qpointf_equal(q.p, other.p)) }

// NotEqual returns !Equal(other).
func (q *QPointF) NotEqual(other *QPointF) bool { return !q.Equal(other) }

// EqualQPoint returns true if q is fuzzy-equal to other (as QPointF).
// Mirrors comparesEqual(QPointF, QPoint).
func (q *QPointF) EqualQPoint(other *QPoint) bool {
	return bool(C.qpointf_equal_point(q.p, other.p))
}

// ---- conversions ---------------------------------------------------------

// ToPoint returns a new QPoint handle via qSaturateRound on each component.
// Mirrors QPointF::toPoint().
func (q *QPointF) ToPoint() *QPoint { return wrapQPoint(C.qpointf_to_point(q.p)) }
