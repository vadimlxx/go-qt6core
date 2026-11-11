package qgeometry

/*
 QLineF — биндинг Qt 6.11 QtCore/QLineF как HANDLE.

 QLineF — value-тип в Qt (Q_PRIMITIVE_TYPE, 2×QPointF). В Go хранится handle на heap-размещённый
 C++ QLineF.

 ВАЖНОЕ ОГРАНИЧЕНИЕ:
   Qt хранит QLineF внутри других классов ПО ЗНАЧЕНИЮ. Handle на QLineF не видит изменений,
   сделанных C++ во встроенном QLineF другого класса.

 Семантика значений vs ссылок:
   - QLineF — ссылочный тип (handle). Для независимой копии — Clone().
   - Go-нулевое значение QLineF{} невалидно. Используйте NewQLineF / NewQLineFZero.

 Управление памятью:
   - NewQLineF / P1() / P2() / Center() / PointAt(t) / UnitVector() / NormalVector() /
     арифметика — создают НОВЫЙ heap-объект.
   - Close() освобождает немедленно; иначе срабатывает finalizer.

 QPointF I/O:
   - Конструкторы/сеттеры принимают *QPointF (значения копируются).
   - Геттеры P1()/P2()/Center()/PointAt() возвращают НОВЫЙ *QPointF (heap copy).
   - Для эффективного доступа — X1()/Y1()/X2()/Y2()/Dx()/Dy()/Length()/Angle() (без аллокации).

 Соглашения:
   - Геттеры: P1()/P2()/Center()/PointAt() (аллоцируют), X1()/Y1()/X2()/Y2()/Dx()/Dy()/
     Length()/Angle()/IsNull() (fuzzy).
   - Геометрия: AngleTo(other) float64, UnitVector() *QLineF, NormalVector() *QLineF,
     Intersects(other) IntersectionType, IntersectsWithPoint(other) (IntersectionType, *QPointF),
     PointAt(t) *QPointF.
   - Мутирующие (in place, void): Translate(*QPointF), TranslateXY(dx, dy), SetLength(len),
     SetAngle(angle), SetP1/SetP2/SetPoints/SetLine.
   - Трансформации (новый handle): Translated(*QPointF), TranslatedXY(dx, dy),
     UnitVector(), NormalVector().
   - Сравнение: Equal(other *QLineF) — fuzzy; NotEqual(other); EqualLine(other *QLine) — кросс.
   - Конвертация: ToLine() → *QLine.
   - Клонирование: Clone() → *QLineF.
   - Конструкторы: NewQLineFFromPoints(p1, p2 *QPointF), QLineFFromLine(l *QLine),
     QLineFFromPolar(len, angle).

 IntersectionType — три константы: NoIntersection, BoundedIntersection, UnboundedIntersection.

 Особенности:
   - IntersectsWithPoint возвращает (type, point): point == nil, если пересечения нет.
   - IsNull() — fuzzy (qFuzzyCompare на обеих точках).
   - SetLength/SetAngle сохраняют P1 и меняют P2.

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

import "runtime"

// IntersectionType mirrors Qt's QLineF::IntersectionType.
type IntersectionType int

const (
	NoIntersection        IntersectionType = 0
	BoundedIntersection   IntersectionType = 1
	UnboundedIntersection IntersectionType = 2
)

// QLineF wraps a heap-allocated Qt6 QLineF.
type QLineF struct {
	p C.go_qlinef_t
}

func wrapQLineF(p C.go_qlinef_t) *QLineF {
	if p == nil {
		return nil
	}
	l := &QLineF{p: p}
	runtime.SetFinalizer(l, finalizeQLineF)
	return l
}

func finalizeQLineF(l *QLineF) {
	if l.p != nil {
		C.qlinef_delete(l.p)
		l.p = nil
	}
}

// Close releases the underlying C++ object immediately and disables the finalizer.
func (l *QLineF) Close() {
	if l.p != nil {
		C.qlinef_delete(l.p)
		l.p = nil
		runtime.SetFinalizer(l, nil)
	}
}

// ---- construction --------------------------------------------------------

func NewQLineF(x1, y1, x2, y2 float64) *QLineF {
	return wrapQLineF(C.qlinef_new(C.double(x1), C.double(y1), C.double(x2), C.double(y2)))
}

func NewQLineFZero() *QLineF { return wrapQLineF(C.qlinef_new_default()) }

func NewQLineFFromPoints(p1, p2 *QPointF) *QLineF {
	return wrapQLineF(C.qlinef_from_points(p1.p, p2.p))
}

// QLineFFromLine builds *QLineF from *QLine (exact int→double).
func QLineFFromLine(l *QLine) *QLineF { return wrapQLineF(C.qlinef_from_line(l.p)) }

// QLineFFromPolar builds *QLineF from polar coordinates (length, angle in degrees).
func QLineFFromPolar(length, angle float64) *QLineF {
	return wrapQLineF(C.qlinef_from_polar(C.double(length), C.double(angle)))
}

// Clone returns a heap copy of l (independent C++ object).
func (l *QLineF) Clone() *QLineF { return wrapQLineF(C.qlinef_clone(l.p)) }

// ---- accessors -----------------------------------------------------------

func (l *QLineF) IsNull() bool   { return bool(C.qlinef_is_null(l.p)) }
func (l *QLineF) X1() float64    { return float64(C.qlinef_x1(l.p)) }
func (l *QLineF) Y1() float64    { return float64(C.qlinef_y1(l.p)) }
func (l *QLineF) X2() float64    { return float64(C.qlinef_x2(l.p)) }
func (l *QLineF) Y2() float64    { return float64(C.qlinef_y2(l.p)) }
func (l *QLineF) Dx() float64    { return float64(C.qlinef_dx(l.p)) }
func (l *QLineF) Dy() float64    { return float64(C.qlinef_dy(l.p)) }
func (l *QLineF) Length() float64 { return float64(C.qlinef_length(l.p)) }
func (l *QLineF) Angle() float64  { return float64(C.qlinef_angle(l.p)) }

func (l *QLineF) AngleTo(other *QLineF) float64 {
	return float64(C.qlinef_angle_to(l.p, other.p))
}

// P1 returns a NEW *QPointF snapshot. Caller owns and must Close().
func (l *QLineF) P1() *QPointF { return wrapQPointF(C.qlinef_p1(l.p)) }

// P2 returns a NEW *QPointF snapshot. Caller owns and must Close().
func (l *QLineF) P2() *QPointF { return wrapQPointF(C.qlinef_p2(l.p)) }

// Center returns a NEW *QPointF snapshot. Caller owns and must Close().
func (l *QLineF) Center() *QPointF { return wrapQPointF(C.qlinef_center(l.p)) }

// PointAt returns a NEW *QPointF at parameter t along the line. Caller owns and must Close().
func (l *QLineF) PointAt(t float64) *QPointF {
	return wrapQPointF(C.qlinef_point_at(l.p, C.double(t)))
}

// ---- mutating (in place) -------------------------------------------------

func (l *QLineF) Translate(p *QPointF)               { C.qlinef_translate(l.p, p.p) }
func (l *QLineF) TranslateXY(dx, dy float64)         { C.qlinef_translate_xy(l.p, C.double(dx), C.double(dy)) }
func (l *QLineF) SetP1(p *QPointF)                   { C.qlinef_set_p1(l.p, p.p) }
func (l *QLineF) SetP2(p *QPointF)                   { C.qlinef_set_p2(l.p, p.p) }
func (l *QLineF) SetPoints(p1, p2 *QPointF)          { C.qlinef_set_points(l.p, p1.p, p2.p) }
func (l *QLineF) SetLength(length float64)           { C.qlinef_set_length(l.p, C.double(length)) }
func (l *QLineF) SetAngle(angle float64)             { C.qlinef_set_angle(l.p, C.double(angle)) }

func (l *QLineF) SetLine(x1, y1, x2, y2 float64) {
	C.qlinef_set_line(l.p, C.double(x1), C.double(y1), C.double(x2), C.double(y2))
}

// ---- transforms (return new handles) ------------------------------------

func (l *QLineF) Translated(p *QPointF) *QLineF {
	return wrapQLineF(C.qlinef_translated(l.p, p.p))
}

func (l *QLineF) TranslatedXY(dx, dy float64) *QLineF {
	return wrapQLineF(C.qlinef_translated_xy(l.p, C.double(dx), C.double(dy)))
}

func (l *QLineF) UnitVector() *QLineF   { return wrapQLineF(C.qlinef_unit_vector(l.p)) }
func (l *QLineF) NormalVector() *QLineF { return wrapQLineF(C.qlinef_normal_vector(l.p)) }

// ---- intersection --------------------------------------------------------

func (l *QLineF) Intersects(other *QLineF) IntersectionType {
	return IntersectionType(C.qlinef_intersects(l.p, other.p, nil))
}

// IntersectsWithPoint returns (type, point). When there is no intersection, point is nil.
func (l *QLineF) IntersectsWithPoint(other *QLineF) (IntersectionType, *QPointF) {
	var out C.go_qpointf_t
	t := IntersectionType(C.qlinef_intersects(l.p, other.p, &out))
	if out == nil {
		return t, nil
	}
	return t, wrapQPointF(out)
}

// ---- comparison ----------------------------------------------------------

func (l *QLineF) Equal(other *QLineF) bool    { return bool(C.qlinef_equal(l.p, other.p)) }
func (l *QLineF) NotEqual(other *QLineF) bool { return !l.Equal(other) }

// EqualLine returns true if l is fuzzy-equal to other *QLine (comparesEqual(QLineF, QLine)).
func (l *QLineF) EqualLine(other *QLine) bool {
	return bool(C.qlinef_equal_line(l.p, other.p))
}

// ---- conversions ---------------------------------------------------------

// ToLine converts l to a new *QLine handle (qSaturateRound per endpoint).
func (l *QLineF) ToLine() *QLine { return wrapQLine(C.qlinef_to_line(l.p)) }
