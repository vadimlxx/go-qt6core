package qgeometry

/*
 QPointF — биндинг Qt 6.11 QtCore/QPointF.

 QPointF — value-тип (Q_PRIMITIVE_TYPE): 2 qreal-поля (x, y). В Go поля неэкспортируемые,
 доступ через геттеры X()/Y(). Через C ABI передаётся плоской структурой go_qpointf_t
 (2×double, 16 байт).

 ОТЛИЧИЕ ОТ QSizeF: Qt-конструктор QPointF() возвращает QPointF(0., 0.), а не (-1., -1.).
 Поэтому:
   - Go-нулевое значение QPointF{} — уже валидная «нулевая» точка
   - NewQPointF(x, y)          — явный конструктор
   - NewQPointFZero()          — удобная обёртка для (0., 0.)
   - NewQPointFDefault() НЕ нужен (совпадает с NewQPointFZero / нулевым значением)

 Соглашения:

   - Геттеры: X(), Y(), IsNull() — IsNull использует qIsNull → qFuzzyIsNull (как Qt).
   - Мутирующие сеттеры: SetX(x), SetY(y).
   - Трансформация: Transposed() — возвращает новый QPointF с обменом координат.
   - ManhattanLength() — |x| + |y| (float64).
   - Арифметика: Add(other), Sub(other), MulFloat(f), DivFloat(d), Neg(), UnaryPlus().
   - Присваивающие формы: AddAssign, SubAssign, MulFloatAssign, DivFloatAssign.
   - Скалярное произведение: DotProduct(other) → float64 (Qt: static QPointF::dotProduct).
   - Сравнение: Equal(other) — fuzzy (comparesEqual → qFuzzyCompare); NotEqual(other);
     EqualQPoint(other QPoint) — сравнение с QPoint (comparesEqual(QPointF, QPoint)).
   - Конвертация: ToPoint() → QPoint (qSaturateRound); QPointFFromPoint(p) — QPoint → QPointF.

 Особенности:

   - QPointF — fuzzy-compared. Два QPointF с крошечной разницей могут быть "равны". Это Qt.
   - DivFloat(0): в release — inf/nan; в debug — Q_ASSERT. Go — IEEE 754 inf/nan без паники.
   - DotProduct: qpoint.x*other.x + qpoint.y*other.y — может дать inf/nan как Qt.
   - ToPoint() использует qSaturateRound на каждой координате.
   - Protected-секций у QPointF нет — открывать нечего.

 Что не реализовано (см. qpoint_doc.go):

   - QDataStream operator<< / operator>> — пакет qdatastream пока не готов.
   - QDebug operator<< — пакет qdebug пока не готов.
   - qHash(QPointF) — в Qt удалён (QPointF не хэшируется как ключ).
   - rx()/ry() — покрыто SetX/SetY.
   - std::tuple protocol — не нужно в Go.
   - fromCGPoint / toCGPoint (Darwin) — платформо-специфично.
   - qFuzzyCompare(QPointF, QPointF) / qFuzzyIsNull(QPointF) как свободные функции —
     покрыты через Equal / IsNull.
*/

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC -Wno-sfinae-incomplete
#cgo LDFLAGS: -lstdc++
#include "qpoint_c.h"
*/
import "C"

import "structs"

// QPointF wraps Qt6 QPointF (Q_PRIMITIVE_TYPE, 2×qreal value).
// Fields are unexported; use X()/Y() to read, SetX()/SetY() to mutate.
type QPointF struct {
	_ structs.HostLayout
	x float64
	y float64
}

// ---- construction --------------------------------------------------------

// NewQPointF returns QPointF(x, y).
func NewQPointF(x, y float64) QPointF {
	return qpointfFromC(C.qpointf_new(C.double(x), C.double(y)))
}

// NewQPointFZero returns QPointF(0., 0.) — same as Qt's default QPointF() and the Go zero value QPointF{}.
func NewQPointFZero() QPointF { return NewQPointF(0, 0) }

// QPointFFromPoint builds QPointF from QPoint (exact int→double). Mirrors QPointF(const QPoint&).
func QPointFFromPoint(p QPoint) QPointF {
	return qpointfFromC(C.qpointf_from_point(qpointToC(p)))
}

// ---- accessors -----------------------------------------------------------

func (p QPointF) IsNull() bool           { return bool(C.qpointf_is_null(qpointfToC(p))) }
func (p QPointF) X() float64             { return float64(C.qpointf_x(qpointfToC(p))) }
func (p QPointF) Y() float64             { return float64(C.qpointf_y(qpointfToC(p))) }
func (p QPointF) ManhattanLength() float64 { return float64(C.qpointf_manhattan_length(qpointfToC(p))) }

// ---- modifiers -----------------------------------------------------------

func (p *QPointF) SetX(x float64) { *p = qpointfFromC(C.qpointf_set_x(qpointfToC(*p), C.double(x))) }
func (p *QPointF) SetY(y float64) { *p = qpointfFromC(C.qpointf_set_y(qpointfToC(*p), C.double(y))) }

// ---- transforms ----------------------------------------------------------

// Transposed returns a QPointF with x and y swapped.
func (p QPointF) Transposed() QPointF { return qpointfFromC(C.qpointf_transposed(qpointfToC(p))) }

// ---- dot product ---------------------------------------------------------

// DotProduct returns p·other = p.X()*other.X() + p.Y()*other.Y().
// Mirrors the static QPointF::dotProduct(p1, p2).
func (p QPointF) DotProduct(other QPointF) float64 {
	return float64(C.qpointf_dot_product(qpointfToC(p), qpointfToC(other)))
}

// ---- arithmetic ----------------------------------------------------------

// Add returns p + other (component-wise).
func (p QPointF) Add(other QPointF) QPointF {
	return qpointfFromC(C.qpointf_add(qpointfToC(p), qpointfToC(other)))
}

// Sub returns p - other (component-wise).
func (p QPointF) Sub(other QPointF) QPointF {
	return qpointfFromC(C.qpointf_sub(qpointfToC(p), qpointfToC(other)))
}

// MulFloat returns p * f (component-wise).
func (p QPointF) MulFloat(f float64) QPointF {
	return qpointfFromC(C.qpointf_mul_float(qpointfToC(p), C.double(f)))
}

// DivFloat returns p / d (component-wise).
// IEEE 754 semantics: d == 0 yields inf/nan (Qt asserts only in debug builds).
func (p QPointF) DivFloat(d float64) QPointF {
	return qpointfFromC(C.qpointf_div_float(qpointfToC(p), C.double(d)))
}

// Neg returns -p (component-wise).
func (p QPointF) Neg() QPointF { return qpointfFromC(C.qpointf_neg(qpointfToC(p))) }

// UnaryPlus returns +p (copy).
func (p QPointF) UnaryPlus() QPointF { return qpointfFromC(C.qpointf_unary_plus(qpointfToC(p))) }

// ---- assign forms --------------------------------------------------------

func (p *QPointF) AddAssign(other QPointF) { *p = p.Add(other) }
func (p *QPointF) SubAssign(other QPointF) { *p = p.Sub(other) }
func (p *QPointF) MulFloatAssign(f float64) { *p = p.MulFloat(f) }
func (p *QPointF) DivFloatAssign(d float64) { *p = p.DivFloat(d) }

// ---- comparison ----------------------------------------------------------

// Equal returns true if both components are fuzzy-equal (Qt comparesEqual → qFuzzyCompare).
func (p QPointF) Equal(other QPointF) bool {
	return bool(C.qpointf_equal(qpointfToC(p), qpointfToC(other)))
}

// NotEqual returns !Equal(other).
func (p QPointF) NotEqual(other QPointF) bool { return !p.Equal(other) }

// EqualQPoint returns true if p is fuzzy-equal to other (as QPointF).
// Mirrors comparesEqual(QPointF, QPoint).
func (p QPointF) EqualQPoint(other QPoint) bool {
	return bool(C.qpointf_equal_point(qpointfToC(p), qpointToC(other)))
}

// ---- conversions ---------------------------------------------------------

// ToPoint converts p to QPoint via qSaturateRound on each component. Mirrors QPointF::toPoint().
func (p QPointF) ToPoint() QPoint { return qpointFromC(C.qpointf_to_point(qpointfToC(p))) }

// ---- internal helpers ----------------------------------------------------

func qpointfToC(p QPointF) C.go_qpointf_t {
	return C.go_qpointf_t{x: C.double(p.x), y: C.double(p.y)}
}

func qpointfFromC(c C.go_qpointf_t) QPointF {
	return QPointF{x: float64(c.x), y: float64(c.y)}
}
