package qgeometry

/*
 QLineF — биндинг Qt 6.11 QtCore/QLineF как VALUE-тип.

 QLineF — value-тип в Qt (Q_PRIMITIVE_TYPE): 2 точки (p1, p2), каждая QPointF.
 В Go хранится как плоская структура из 4 float64-полей (x1, y1, x2, y2) — полная
 копия данных. Присваивание Go-переменной создаёт независимую копию (value-семантика).
 Указателей, финализаторов и Close() нет.

 Default ctor: QLineF() = QLineF((0,0), (0,0)) — «нулевая» линия, IsNull() == true (fuzzy).
 Поэтому:
   - Go-нулевое значение QLineF{} — уже валидная нулевая линия
   - NewQLineF(x1, y1, x2, y2)         — явный конструктор (4 float64)
   - NewQLineFZero()                   — удобная обёртка для (0,0)-(0,0)
   - NewQLineFFromPoints(p1, p2 QPointF) — из двух QPointF
   - QLineFFromLine(l QLine)           — из QLine (implicit ctor)
   - QLineFFromPolar(length, angle)    — static QLineF::fromPolar

 Соглашения:

   - Геттеры: P1()/P2()/Center()/PointAt(t) → QPointF (значение),
     X1()/Y1()/X2()/Y2()/Dx()/Dy()/Length()/Angle()/IsNull().
   - Геометрия: AngleTo(other), UnitVector(), NormalVector(), Intersects(other),
     IntersectsWithPoint(other).
   - Мутирующие сеттеры (возвращают обновлённое значение): Translate, TranslateXY,
     SetLength, SetAngle, SetP1, SetP2, SetPoints, SetLine.
   - Трансформации (возвращают новое значение): Translated, TranslatedXY, UnitVector,
     NormalVector.
   - Сравнение: Equal(other) — fuzzy; NotEqual(other); EqualLine(other QLine) — кросс.
   - Конвертация: ToLine() → QLine (qSaturateRound per endpoint).

 IntersectionType — три константы: NoIntersection, BoundedIntersection, UnboundedIntersection.

 Особенности:

   - IntersectsWithPoint возвращает (type, point, ok bool): ok == false, если пересечения нет.
   - IsNull() — fuzzy (qFuzzyCompare на обеих точках).
   - SetLength/SetAngle сохраняют P1 и меняют P2.
   - Углы — в экранных координатах Qt (y вниз), CCW от +x, [0, 360).

 Что не реализовано (см. qline_doc.go):

   - QDataStream operator<< / operator>> — пакет qdatastream пока не готов.
   - QDebug operator<< — пакет qdebug пока не готов.
   - operator==/operator!= как перегруженные операторы.
   - qFuzzyCompare(QLineF, QLineF) / qFuzzyIsNull(QLineF) — покрыты Equal/IsNull.
*/

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC -Wno-sfinae-incomplete
#cgo LDFLAGS: -lstdc++
#include "qline_c.h"
*/
import "C"

import "structs"

// IntersectionType mirrors Qt's QLineF::IntersectionType.
type IntersectionType int

const (
	NoIntersection        IntersectionType = 0
	BoundedIntersection   IntersectionType = 1
	UnboundedIntersection IntersectionType = 2
)

// QLineF wraps Qt6 QLineF (Q_PRIMITIVE_TYPE, 2×QPointF value).
// Fields are unexported; use X1()/Y1()/X2()/Y2() or P1()/P2() to read.
type QLineF struct {
	_      structs.HostLayout
	x1, y1, x2, y2 float64
}

// ---- construction --------------------------------------------------------

func NewQLineF(x1, y1, x2, y2 float64) QLineF {
	return qlinefFromC(C.qlinef_new(C.double(x1), C.double(y1), C.double(x2), C.double(y2)))
}

func NewQLineFZero() QLineF { return NewQLineF(0, 0, 0, 0) }

func NewQLineFFromPoints(p1, p2 QPointF) QLineF {
	return qlinefFromC(C.qlinef_from_points(qlinePointFToC(p1), qlinePointFToC(p2)))
}

// QLineFFromLine builds QLineF from QLine (exact int→double).
func QLineFFromLine(l QLine) QLineF {
	return qlinefFromC(C.qlinef_from_line(l.toC()))
}

// QLineFFromPolar builds QLineF from polar coordinates (length, angle in degrees).
func QLineFFromPolar(length, angle float64) QLineF {
	return qlinefFromC(C.qlinef_from_polar(C.double(length), C.double(angle)))
}

// ---- accessors -----------------------------------------------------------

func (l QLineF) IsNull() bool    { return bool(C.qlinef_is_null(l.toC())) }
func (l QLineF) X1() float64     { return float64(C.qlinef_x1(l.toC())) }
func (l QLineF) Y1() float64     { return float64(C.qlinef_y1(l.toC())) }
func (l QLineF) X2() float64     { return float64(C.qlinef_x2(l.toC())) }
func (l QLineF) Y2() float64     { return float64(C.qlinef_y2(l.toC())) }
func (l QLineF) Dx() float64     { return float64(C.qlinef_dx(l.toC())) }
func (l QLineF) Dy() float64     { return float64(C.qlinef_dy(l.toC())) }
func (l QLineF) Length() float64 { return float64(C.qlinef_length(l.toC())) }
func (l QLineF) Angle() float64  { return float64(C.qlinef_angle(l.toC())) }

func (l QLineF) AngleTo(other QLineF) float64 {
	return float64(C.qlinef_angle_to(l.toC(), other.toC()))
}

func (l QLineF) P1() QPointF     { return qpointfFromC(C.qlinef_p1(l.toC())) }
func (l QLineF) P2() QPointF     { return qpointfFromC(C.qlinef_p2(l.toC())) }
func (l QLineF) Center() QPointF { return qpointfFromC(C.qlinef_center(l.toC())) }

func (l QLineF) PointAt(t float64) QPointF {
	return qpointfFromC(C.qlinef_point_at(l.toC(), C.double(t)))
}

// ---- mutating setters (return updated value) ----------------------------

func (l QLineF) Translate(p QPointF) QLineF {
	return qlinefFromC(C.qlinef_translate(l.toC(), qlinePointFToC(p)))
}
func (l QLineF) TranslateXY(dx, dy float64) QLineF {
	return qlinefFromC(C.qlinef_translate_xy(l.toC(), C.double(dx), C.double(dy)))
}
func (l QLineF) SetP1(p QPointF) QLineF {
	return qlinefFromC(C.qlinef_set_p1(l.toC(), qlinePointFToC(p)))
}
func (l QLineF) SetP2(p QPointF) QLineF {
	return qlinefFromC(C.qlinef_set_p2(l.toC(), qlinePointFToC(p)))
}
func (l QLineF) SetPoints(p1, p2 QPointF) QLineF {
	return qlinefFromC(C.qlinef_set_points(l.toC(), qlinePointFToC(p1), qlinePointFToC(p2)))
}
func (l QLineF) SetLine(x1, y1, x2, y2 float64) QLineF {
	return qlinefFromC(C.qlinef_set_line(l.toC(),
		C.double(x1), C.double(y1), C.double(x2), C.double(y2)))
}
func (l QLineF) SetLength(length float64) QLineF {
	return qlinefFromC(C.qlinef_set_length(l.toC(), C.double(length)))
}
func (l QLineF) SetAngle(angle float64) QLineF {
	return qlinefFromC(C.qlinef_set_angle(l.toC(), C.double(angle)))
}

// ---- transforms (return new value) --------------------------------------

func (l QLineF) Translated(p QPointF) QLineF {
	return qlinefFromC(C.qlinef_translated(l.toC(), qlinePointFToC(p)))
}
func (l QLineF) TranslatedXY(dx, dy float64) QLineF {
	return qlinefFromC(C.qlinef_translated_xy(l.toC(), C.double(dx), C.double(dy)))
}
func (l QLineF) UnitVector() QLineF   { return qlinefFromC(C.qlinef_unit_vector(l.toC())) }
func (l QLineF) NormalVector() QLineF { return qlinefFromC(C.qlinef_normal_vector(l.toC())) }

// ---- intersection --------------------------------------------------------

func (l QLineF) Intersects(other QLineF) IntersectionType {
	return IntersectionType(C.qlinef_intersects(l.toC(), other.toC(), nil))
}

// IntersectsWithPoint returns (type, point, ok). When there is no intersection, ok is false
// and point is the zero QPointF.
func (l QLineF) IntersectsWithPoint(other QLineF) (IntersectionType, QPointF, bool) {
	var out C.go_qpointf_t
	t := IntersectionType(C.qlinef_intersects(l.toC(), other.toC(), &out))
	if t == NoIntersection {
		return t, QPointF{}, false
	}
	return t, qpointfFromC(out), true
}

// ---- comparison ----------------------------------------------------------

func (l QLineF) Equal(other QLineF) bool    { return bool(C.qlinef_equal(l.toC(), other.toC())) }
func (l QLineF) NotEqual(other QLineF) bool { return !l.Equal(other) }

// EqualLine returns true if l is fuzzy-equal to other QLine (comparesEqual(QLineF, QLine)).
func (l QLineF) EqualLine(other QLine) bool {
	return bool(C.qlinef_equal_line(l.toC(), other.toC()))
}

// ---- conversions ---------------------------------------------------------

// ToLine converts l to QLine (qSaturateRound per endpoint).
func (l QLineF) ToLine() QLine { return qlineFromC(C.qlinef_to_line(l.toC())) }

// ---- internal helpers ----------------------------------------------------

// toC converts l to its C ABI representation.
func (l QLineF) toC() C.go_qlinef_t {
	return C.go_qlinef_t{
		x1: C.double(l.x1), y1: C.double(l.y1),
		x2: C.double(l.x2), y2: C.double(l.y2),
	}
}

func qlinefFromC(c C.go_qlinef_t) QLineF {
	return QLineF{
		x1: float64(c.x1), y1: float64(c.y1),
		x2: float64(c.x2), y2: float64(c.y2),
	}
}

// qlinePointFToC converts a QPointF value to C.go_qpointf_t using only its public API.
// Local helper: avoids depending on any specific internal helper name from qpointf.go.
func qlinePointFToC(p QPointF) C.go_qpointf_t {
	return C.go_qpointf_t{x: C.double(p.X()), y: C.double(p.Y())}
}
