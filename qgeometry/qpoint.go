package qgeometry

/*
 QPoint — биндинг Qt 6.11 QtCore/QPoint.

 QPoint — value-тип (Q_PRIMITIVE_TYPE): 2 int-поля (x, y). В Go поля неэкспортируемые,
 доступ через геттеры X()/Y(). Через C ABI передаётся плоской структурой go_qpoint_t
 (2×int, 8 байт).

 ОТЛИЧИЕ ОТ QSize: Qt-конструктор QPoint() возвращает QPoint(0, 0) — «нулевую» точку,
 а не (-1, -1). Поэтому:
   - Go-нулевое значение QPoint{} — уже валидная «нулевая» точка, эквивалентная QPoint()
   - NewQPoint(x, y)          — явный конструктор
   - NewQPointZero()          — удобная обёртка для (0, 0)
   - NewQPointDefault() НЕ нужен (совпадает с NewQPointZero / нулевым значением)

 Соглашения:

   - Геттеры: X(), Y(), IsNull().
   - Мутирующие сеттеры: SetX(x), SetY(y).
   - Трансформация: Transposed() — возвращает новый QPoint с обменом координат.
   - ManhattanLength() — |x| + |y| (int, wraps на overflow как Qt release).
   - Арифметика: Add(other), Sub(other), MulInt(f), MulFloat(f), DivFloat(d), Neg(), UnaryPlus().
     В Qt есть operator*=(int/float/double); в Go — только MulInt (int) и MulFloat (float64).
   - Присваивающие формы: AddAssign, SubAssign, MulIntAssign, MulFloatAssign, DivFloatAssign.
   - Скалярное произведение: DotProduct(other) → int (Qt: static QPoint::dotProduct).
   - Сравнение: Equal(other), NotEqual(other) — точное (comparesEqual).
   - Конвертация: ToPointF() → QPointF.

 Особенности:

   - MulInt переполняется как C++ int (QCheckedInt в debug ассертит, в release wraps).
   - MulFloat/DivFloat используют Qt's qSaturateRound — идентично operator*(QPoint, float/double)
     и friend operator/(QPoint, qreal).
   - DotProduct: int(qpoint.x*other.x + qpoint.y*other.y) — может переполниться (как Qt).
   - Protected-секций у QPoint нет — открывать нечего.

 Что не реализовано (см. qpoint_doc.go):

   - QDataStream operator<< / operator>> — пакет qdatastream пока не готов.
   - QDebug operator<< — пакет qdebug пока не готов.
   - qHash(QPoint) — не нужно в Go.
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

// QPoint wraps Qt6 QPoint (Q_PRIMITIVE_TYPE, 2×int value).
// Fields are unexported; use X()/Y() to read, SetX()/SetY() to mutate.
type QPoint struct {
	x int
	y int
}

// ---- construction --------------------------------------------------------

// NewQPoint returns QPoint(x, y).
func NewQPoint(x, y int) QPoint {
	return qpointFromC(C.qpoint_new(C.int(x), C.int(y)))
}

// NewQPointZero returns QPoint(0, 0) — same as Qt's default QPoint() and the Go zero value QPoint{}.
func NewQPointZero() QPoint { return NewQPoint(0, 0) }

// ---- accessors -----------------------------------------------------------

func (p QPoint) IsNull() bool          { return bool(C.qpoint_is_null(qpointToC(p))) }
func (p QPoint) X() int                { return int(C.qpoint_x(qpointToC(p))) }
func (p QPoint) Y() int                { return int(C.qpoint_y(qpointToC(p))) }
func (p QPoint) ManhattanLength() int  { return int(C.qpoint_manhattan_length(qpointToC(p))) }

// ---- modifiers -----------------------------------------------------------

func (p *QPoint) SetX(x int) { *p = qpointFromC(C.qpoint_set_x(qpointToC(*p), C.int(x))) }
func (p *QPoint) SetY(y int) { *p = qpointFromC(C.qpoint_set_y(qpointToC(*p), C.int(y))) }

// ---- transforms ----------------------------------------------------------

// Transposed returns a QPoint with x and y swapped.
func (p QPoint) Transposed() QPoint { return qpointFromC(C.qpoint_transposed(qpointToC(p))) }

// ---- dot product ---------------------------------------------------------

// DotProduct returns p·other = p.X()*other.X() + p.Y()*other.Y().
// Mirrors the static QPoint::dotProduct(p1, p2). May overflow int as in Qt.
func (p QPoint) DotProduct(other QPoint) int {
	return int(C.qpoint_dot_product(qpointToC(p), qpointToC(other)))
}

// ---- arithmetic ----------------------------------------------------------

// Add returns p + other (component-wise).
func (p QPoint) Add(other QPoint) QPoint {
	return qpointFromC(C.qpoint_add(qpointToC(p), qpointToC(other)))
}

// Sub returns p - other (component-wise).
func (p QPoint) Sub(other QPoint) QPoint {
	return qpointFromC(C.qpoint_sub(qpointToC(p), qpointToC(other)))
}

// MulInt returns p * f (int factor, QCheckedInt arithmetic — may wrap on overflow in release builds).
func (p QPoint) MulInt(f int) QPoint {
	return qpointFromC(C.qpoint_mul_int(qpointToC(p), C.int(f)))
}

// MulFloat returns p * f (rounded via Qt's qSaturateRound).
func (p QPoint) MulFloat(f float64) QPoint {
	return qpointFromC(C.qpoint_mul_float(qpointToC(p), C.double(f)))
}

// DivFloat returns p / d (rounded via Qt's qSaturateRound).
func (p QPoint) DivFloat(d float64) QPoint {
	return qpointFromC(C.qpoint_div_float(qpointToC(p), C.double(d)))
}

// Neg returns -p (component-wise).
func (p QPoint) Neg() QPoint { return qpointFromC(C.qpoint_neg(qpointToC(p))) }

// UnaryPlus returns +p (copy).
func (p QPoint) UnaryPlus() QPoint { return qpointFromC(C.qpoint_unary_plus(qpointToC(p))) }

// ---- assign forms --------------------------------------------------------

func (p *QPoint) AddAssign(other QPoint)   { *p = p.Add(other) }
func (p *QPoint) SubAssign(other QPoint)   { *p = p.Sub(other) }
func (p *QPoint) MulIntAssign(f int)       { *p = p.MulInt(f) }
func (p *QPoint) MulFloatAssign(f float64) { *p = p.MulFloat(f) }
func (p *QPoint) DivFloatAssign(d float64) { *p = p.DivFloat(d) }

// ---- comparison ----------------------------------------------------------

// Equal returns true if both coordinates are identical (exact int comparison).
func (p QPoint) Equal(other QPoint) bool {
	return bool(C.qpoint_equal(qpointToC(p), qpointToC(other)))
}

// NotEqual returns !Equal(other).
func (p QPoint) NotEqual(other QPoint) bool { return !p.Equal(other) }

// ---- conversions ---------------------------------------------------------

// ToPointF converts p to QPointF (exact int→double per component). Mirrors QPoint::toPointF().
func (p QPoint) ToPointF() QPointF { return qpointfFromC(C.qpoint_to_pointf(qpointToC(p))) }

// ---- internal helpers ----------------------------------------------------

func qpointToC(p QPoint) C.go_qpoint_t {
	return C.go_qpoint_t{x: C.int(p.x), y: C.int(p.y)}
}

func qpointFromC(c C.go_qpoint_t) QPoint {
	return QPoint{x: int(c.x), y: int(c.y)}
}
