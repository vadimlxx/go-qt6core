package qgeometry

/*
 QLine — биндинг Qt 6.11 QtCore/QLine как HANDLE.

 QLine — value-тип в Qt (Q_PRIMITIVE_TYPE, 2×QPoint). В Go хранится handle на heap-размещённый
 C++ QLine. Все геттеры/сеттеры идут напрямую в C++.

 ВАЖНОЕ ОГРАНИЧЕНИЕ:
   Qt хранит QLine внутри других классов ПО ЗНАЧЕНИЮ. Handle на QLine не видит изменений,
   сделанных C++ во встроенном QLine другого класса. Это ограничение API Qt, а не биндинга.

 Семантика значений vs ссылок:
   - QLine — ссылочный тип (handle). Мутации через один handle видны через все handle'ы на
     тот же объект. Для независимой копии — Clone().
   - Go-нулевое значение QLine{} невалидно. Используйте NewQLine / NewQLineZero.

 Управление памятью:
   - NewQLine / P1() / P2() / Center() / арифметика создают НОВЫЙ heap-объект.
   - Close() освобождает немедленно; иначе срабатывает finalizer.

 QPoint I/O:
   - Конструкторы принимают *QPoint (значения копируются в QLine).
   - Геттеры P1()/P2()/Center() возвращают НОВЫЙ *QPoint (heap copy) — вызывающий владеет.
   - Сеттеры принимают *QPoint (значения копируются).
   - Для эффективного доступа к координатам — X1()/Y1()/X2()/Y2() (без аллокации).

 Соглашения:
   - Геттеры: P1()/P2()/Center() (аллоцируют), X1()/Y1()/X2()/Y2()/Dx()/Dy()/IsNull().
   - Мутирующие (in place, void): Translate(*QPoint), TranslateXY(dx, dy),
     SetP1/SetP2/SetPoints/SetLine.
   - Трансформации (новый handle): Translated(*QPoint), TranslatedXY(dx, dy).
   - Сравнение: Equal(other *QLine), NotEqual(other).
   - Конвертация: ToLineF() → *QLineF.
   - Клонирование: Clone() → *QLine.
   - Арифметических операторов (+, -, *, /) у QLine нет.

 Что не реализовано (см. qline_doc.go):

   - QDataStream operator<< / operator>> — пакет qdatastream пока не готов.
   - QDebug operator<< — пакет qdebug пока не готов.
   - operator==/operator!= как перегруженные операторы — используйте Equal/NotEqual.
*/

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC -Wno-sfinae-incomplete
#cgo LDFLAGS: -lstdc++
#include "qline_c.h"
*/
import "C"

import "runtime"

// QLine wraps a heap-allocated Qt6 QLine.
type QLine struct {
	p C.go_qline_t
}

func wrapQLine(p C.go_qline_t) *QLine {
	if p == nil {
		return nil
	}
	l := &QLine{p: p}
	runtime.SetFinalizer(l, finalizeQLine)
	return l
}

func finalizeQLine(l *QLine) {
	if l.p != nil {
		C.qline_delete(l.p)
		l.p = nil
	}
}

// Close releases the underlying C++ object immediately and disables the finalizer.
func (l *QLine) Close() {
	if l.p != nil {
		C.qline_delete(l.p)
		l.p = nil
		runtime.SetFinalizer(l, nil)
	}
}

// ---- construction --------------------------------------------------------

// NewQLine returns a handle to QLine(x1, y1, x2, y2).
func NewQLine(x1, y1, x2, y2 int) *QLine {
	return wrapQLine(C.qline_new(C.int(x1), C.int(y1), C.int(x2), C.int(y2)))
}

// NewQLineZero returns a handle to QLine((0,0), (0,0)) — same as Qt's default QLine().
func NewQLineZero() *QLine {
	return wrapQLine(C.qline_new_default())
}

// NewQLineFromPoints returns a handle to QLine(p1, p2); p1/p2 values are copied from the handles.
func NewQLineFromPoints(p1, p2 *QPoint) *QLine {
	return wrapQLine(C.qline_from_points(p1.p, p2.p))
}

// Clone returns a heap copy of l (independent C++ object).
func (l *QLine) Clone() *QLine { return wrapQLine(C.qline_clone(l.p)) }

// ---- accessors -----------------------------------------------------------

func (l *QLine) IsNull() bool { return bool(C.qline_is_null(l.p)) }
func (l *QLine) X1() int      { return int(C.qline_x1(l.p)) }
func (l *QLine) Y1() int      { return int(C.qline_y1(l.p)) }
func (l *QLine) X2() int      { return int(C.qline_x2(l.p)) }
func (l *QLine) Y2() int      { return int(C.qline_y2(l.p)) }
func (l *QLine) Dx() int      { return int(C.qline_dx(l.p)) }
func (l *QLine) Dy() int      { return int(C.qline_dy(l.p)) }

// P1 returns a NEW *QPoint snapshot of l's p1. Caller owns and must Close().
func (l *QLine) P1() *QPoint { return wrapQPoint(C.qline_p1(l.p)) }

// P2 returns a NEW *QPoint snapshot of l's p2. Caller owns and must Close().
func (l *QLine) P2() *QPoint { return wrapQPoint(C.qline_p2(l.p)) }

// Center returns a NEW *QPoint snapshot of the line's center. Caller owns and must Close().
func (l *QLine) Center() *QPoint { return wrapQPoint(C.qline_center(l.p)) }

// ---- mutating (in place) -------------------------------------------------

func (l *QLine) Translate(p *QPoint)              { C.qline_translate(l.p, p.p) }
func (l *QLine) TranslateXY(dx, dy int)           { C.qline_translate_xy(l.p, C.int(dx), C.int(dy)) }
func (l *QLine) SetP1(p *QPoint)                  { C.qline_set_p1(l.p, p.p) }
func (l *QLine) SetP2(p *QPoint)                  { C.qline_set_p2(l.p, p.p) }
func (l *QLine) SetPoints(p1, p2 *QPoint)         { C.qline_set_points(l.p, p1.p, p2.p) }
func (l *QLine) SetLine(x1, y1, x2, y2 int) {
	C.qline_set_line(l.p, C.int(x1), C.int(y1), C.int(x2), C.int(y2))
}

// ---- transforms (return new handles) ------------------------------------

func (l *QLine) Translated(p *QPoint) *QLine {
	return wrapQLine(C.qline_translated(l.p, p.p))
}

func (l *QLine) TranslatedXY(dx, dy int) *QLine {
	return wrapQLine(C.qline_translated_xy(l.p, C.int(dx), C.int(dy)))
}

// ---- comparison ----------------------------------------------------------

func (l *QLine) Equal(other *QLine) bool    { return bool(C.qline_equal(l.p, other.p)) }
func (l *QLine) NotEqual(other *QLine) bool { return !l.Equal(other) }

// ---- conversions ---------------------------------------------------------

// ToLineF converts l to a new *QLineF handle (exact int→double per coordinate).
func (l *QLine) ToLineF() *QLineF { return wrapQLineF(C.qline_to_linef(l.p)) }
