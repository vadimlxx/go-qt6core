package qgeometry

/*
 QPoint — биндинг Qt 6.11 QtCore/QPoint как HANDLE.

 QPoint — value-тип в Qt (Q_PRIMITIVE_TYPE, 4 int-поля). Однако в Go хранится НЕ копия,
 а handle на heap-размещённый C++ QPoint. Все геттеры/сеттеры идут напрямую в C++:
   - p.X() → p->x()  — всегда читает актуальное значение из C++
   - p.SetX(5)       — мутирует C++-объект, изменения видны через все handle'ы на тот же объект

 ВАЖНОЕ ОГРАНИЧЕНИЕ:
   Qt хранит QPoint внутри других классов (QLine, QRect, …) ПО ЗНАЧЕНИЮ. Если C++-класс
   содержит QPoint и мутирует его, наш *QPoint этого НЕ увидит — handle указывает на
   отдельный heap-объект, а не на внутреннее хранилище того класса. Это ограничение API Qt
   (в Qt 6 нет rpt1()), а не биндинга. Для составных типов (QLine) используйте свежие
   геттеры (line.X1(), line.Y1(), …) — они читают значения из C++ каждый раз.

 Семантика значений vs ссылок:
   - QPoint — ссылочный тип (handle). Мутации через один handle видны через все handle'ы на
     тот же объект. Чтобы получить независимую копию — используйте p.Clone().
   - Go-нулевое значение QPoint{} невалидно (nil handle). Используйте NewQPoint/NewQPointZero.

 Управление памятью:
   - Каждый NewQPoint / QLine.P1() / арифметика (*QPoint-returning метод) создаёт НОВЫЙ
     heap QPoint, за который отвечает вызывающий.
   - Close() освобождает немедленно; иначе срабатывает finalizer.

 Соглашения:
   - Геттеры: X(), Y(), IsNull(), ManhattanLength().
   - Мутирующие: SetX(x), SetY(y) — мутируют C++-объект.
   - Арифметика возвращает НОВЫЙ handle: Add, Sub, MulInt, MulFloat, DivFloat, Neg, UnaryPlus,
     Transposed.
   - Присваивающие формы мутируют получатель: AddAssign, SubAssign, MulIntAssign,
     MulFloatAssign, DivFloatAssign.
   - DotProduct(other) → int.
   - Сравнение: Equal(other), NotEqual(other).
   - Конвертация: ToPointF() → *QPointF (fresh handle).
   - Клонирование: Clone() → *QPoint (fresh heap copy).

 Что не реализовано (см. qpoint_doc.go):

   - QDataStream operator<< / operator>> — пакет qdatastream пока не готов.
   - QDebug operator<< — пакет qdebug пока не готов.
   - qHash(QPoint) — не нужно в Go (пользователь может хэшировать X()/Y()).
   - rx()/ry() — покрыто SetX/SetY.
   - std::tuple protocol — не нужно в Go.
   - toCGPoint (Darwin) — платформо-специфично.
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

// QPoint wraps a heap-allocated Qt6 QPoint. Handlers are QPoint handles; all operations go
// directly to the C++ object.
type QPoint struct {
	p C.go_qpoint_t
}

// wrapQPoint takes ownership of a C++ QPoint* and installs a finalizer.
func wrapQPoint(p C.go_qpoint_t) *QPoint {
	if p == nil {
		return nil
	}
	q := &QPoint{p: p}
	runtime.SetFinalizer(q, finalizeQPoint)
	return q
}

func finalizeQPoint(q *QPoint) {
	if q.p != nil {
		C.qpoint_delete(q.p)
		q.p = nil
	}
}

// Close releases the underlying C++ object immediately and disables the finalizer.
func (q *QPoint) Close() {
	if q.p != nil {
		C.qpoint_delete(q.p)
		q.p = nil
		runtime.SetFinalizer(q, nil)
	}
}

// ---- construction --------------------------------------------------------

// NewQPoint returns a handle to QPoint(x, y).
func NewQPoint(x, y int) *QPoint {
	return wrapQPoint(C.qpoint_new(C.int(x), C.int(y)))
}

// NewQPointZero returns a handle to QPoint(0, 0) — same as Qt's default QPoint().
func NewQPointZero() *QPoint {
	return wrapQPoint(C.qpoint_new_default())
}

// Clone returns a heap copy of q (independent C++ object; safe for value semantics).
func (q *QPoint) Clone() *QPoint {
	return wrapQPoint(C.qpoint_clone(q.p))
}

// ---- accessors -----------------------------------------------------------

func (q *QPoint) IsNull() bool          { return bool(C.qpoint_is_null(q.p)) }
func (q *QPoint) X() int                { return int(C.qpoint_x(q.p)) }
func (q *QPoint) Y() int                { return int(C.qpoint_y(q.p)) }
func (q *QPoint) ManhattanLength() int  { return int(C.qpoint_manhattan_length(q.p)) }

// ---- modifiers -----------------------------------------------------------

func (q *QPoint) SetX(x int) { C.qpoint_set_x(q.p, C.int(x)) }
func (q *QPoint) SetY(y int) { C.qpoint_set_y(q.p, C.int(y)) }

// ---- transforms (return new handles) ------------------------------------

// Transposed returns a new handle to q.Transposed() (x and y swapped).
func (q *QPoint) Transposed() *QPoint { return wrapQPoint(C.qpoint_transposed(q.p)) }

// ---- dot product ---------------------------------------------------------

// DotProduct returns q·other = q.X()*other.X() + q.Y()*other.Y().
// Mirrors the static QPoint::dotProduct(p1, p2). May overflow int as in Qt.
func (q *QPoint) DotProduct(other *QPoint) int {
	return int(C.qpoint_dot_product(q.p, other.p))
}

// ---- arithmetic (return new handles) ------------------------------------

func (q *QPoint) Add(other *QPoint) *QPoint       { return wrapQPoint(C.qpoint_add(q.p, other.p)) }
func (q *QPoint) Sub(other *QPoint) *QPoint       { return wrapQPoint(C.qpoint_sub(q.p, other.p)) }
func (q *QPoint) MulInt(f int) *QPoint            { return wrapQPoint(C.qpoint_mul_int(q.p, C.int(f))) }
func (q *QPoint) MulFloat(f float64) *QPoint      { return wrapQPoint(C.qpoint_mul_float(q.p, C.double(f))) }
func (q *QPoint) DivFloat(d float64) *QPoint      { return wrapQPoint(C.qpoint_div_float(q.p, C.double(d))) }
func (q *QPoint) Neg() *QPoint                    { return wrapQPoint(C.qpoint_neg(q.p)) }
func (q *QPoint) UnaryPlus() *QPoint              { return wrapQPoint(C.qpoint_unary_plus(q.p)) }

// ---- in-place assign (mutate receiver) ----------------------------------

func (q *QPoint) AddAssign(other *QPoint)   { C.qpoint_add_assign(q.p, other.p) }
func (q *QPoint) SubAssign(other *QPoint)   { C.qpoint_sub_assign(q.p, other.p) }
func (q *QPoint) MulIntAssign(f int)        { C.qpoint_mul_int_assign(q.p, C.int(f)) }
func (q *QPoint) MulFloatAssign(f float64)  { C.qpoint_mul_float_assign(q.p, C.double(f)) }
func (q *QPoint) DivFloatAssign(d float64)  { C.qpoint_div_float_assign(q.p, C.double(d)) }

// ---- comparison ----------------------------------------------------------

// Equal returns true if both coordinates are identical (exact int comparison).
func (q *QPoint) Equal(other *QPoint) bool { return bool(C.qpoint_equal(q.p, other.p)) }

// NotEqual returns !Equal(other).
func (q *QPoint) NotEqual(other *QPoint) bool { return !q.Equal(other) }

// ---- conversions ---------------------------------------------------------

// ToPointF returns a new QPointF handle with q's coordinates converted to double.
// Mirrors QPoint::toPointF().
func (q *QPoint) ToPointF() *QPointF { return wrapQPointF(C.qpoint_to_pointf(q.p)) }
