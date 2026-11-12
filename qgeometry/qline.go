package qgeometry

/*
 QLine — биндинг Qt 6.11 QtCore/QLine как VALUE-тип.

 QLine — value-тип в Qt (Q_PRIMITIVE_TYPE): 2 точки (p1, p2), каждая QPoint.
 В Go хранится как плоская структура из 4 int-полей (x1, y1, x2, y2) — полная копия
 данных. Присваивание Go-переменной создаёт независимую копию (value-семантика).
 Указателей, финализаторов и Close() нет.

 Default ctor: QLine() = QLine((0,0), (0,0)) — «нулевая» линия, IsNull() == true.
 Поэтому:
   - Go-нулевое значение QLine{} — уже валидная нулевая линия, эквивалентная QLine()
   - NewQLine(x1, y1, x2, y2)          — явный конструктор (4 int)
   - NewQLineZero()                    — удобная обёртка для (0,0)-(0,0)
   - NewQLineFromPoints(p1, p2 QPoint) — из двух QPoint (значения копируются)

 Соглашения:

   - Геттеры: P1(), P2() → QPoint (значение), X1(), Y1(), X2(), Y2(), Dx(), Dy(),
     IsNull(), Center() → QPoint (значение).
   - Мутирующие сеттеры (возвращают обновлённое значение): SetP1, SetP2, SetPoints,
     SetLine, Translate, TranslateXY.
   - Трансформации (возвращают новое значение): Translated, TranslatedXY.
   - Сравнение: Equal(other), NotEqual(other) — точное (comparesEqual → operator==).
   - Конвертация: ToLineF() → QLineF (значение).
   - Арифметических операторов (+, -, *, /) у QLine в Qt НЕТ.

 Особенности:

   - IsNull() — точное сравнение p1 == p2 (в отличие от QLineF, где fuzzy).
   - Center() — QPoint((x1+x2)/2, (y1+y2)/2) с промежуточным 64-битным сложением.
   - Dx()/Dy() — x2-x1, y2-y1 (int, wraps на overflow как Qt release).
   - Protected-секций у QLine нет — открывать нечего.

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

// QLine wraps Qt6 QLine (Q_PRIMITIVE_TYPE, 2×QPoint value).
// Fields are unexported; use X1()/Y1()/X2()/Y2() or P1()/P2() to read.
type QLine struct {
	x1, y1, x2, y2 int
}

// ---- construction --------------------------------------------------------

// NewQLine returns QLine(x1, y1, x2, y2).
func NewQLine(x1, y1, x2, y2 int) QLine {
	return qlineFromC(C.qline_new(C.int(x1), C.int(y1), C.int(x2), C.int(y2)))
}

// NewQLineZero returns QLine((0,0), (0,0)) — same as Qt's default QLine() and the Go zero value.
func NewQLineZero() QLine { return NewQLine(0, 0, 0, 0) }

// NewQLineFromPoints returns QLine(p1, p2). Values are copied from p1 and p2.
func NewQLineFromPoints(p1, p2 QPoint) QLine {
	return qlineFromC(C.qline_from_points(qlinePointToC(p1), qlinePointToC(p2)))
}

// ---- accessors -----------------------------------------------------------

func (l QLine) IsNull() bool { return bool(C.qline_is_null(l.toC())) }
func (l QLine) X1() int      { return int(C.qline_x1(l.toC())) }
func (l QLine) Y1() int      { return int(C.qline_y1(l.toC())) }
func (l QLine) X2() int      { return int(C.qline_x2(l.toC())) }
func (l QLine) Y2() int      { return int(C.qline_y2(l.toC())) }
func (l QLine) Dx() int      { return int(C.qline_dx(l.toC())) }
func (l QLine) Dy() int      { return int(C.qline_dy(l.toC())) }

// P1 returns a copy of l's p1 as a QPoint value.
func (l QLine) P1() QPoint { return qpointFromC(C.qline_p1(l.toC())) }

// P2 returns a copy of l's p2 as a QPoint value.
func (l QLine) P2() QPoint { return qpointFromC(C.qline_p2(l.toC())) }

// Center returns a copy of the line's center as a QPoint value.
func (l QLine) Center() QPoint { return qpointFromC(C.qline_center(l.toC())) }

// ---- mutating setters (return updated value) ----------------------------

func (l QLine) SetP1(p QPoint) QLine {
	return qlineFromC(C.qline_set_p1(l.toC(), qlinePointToC(p)))
}
func (l QLine) SetP2(p QPoint) QLine {
	return qlineFromC(C.qline_set_p2(l.toC(), qlinePointToC(p)))
}
func (l QLine) SetPoints(p1, p2 QPoint) QLine {
	return qlineFromC(C.qline_set_points(l.toC(), qlinePointToC(p1), qlinePointToC(p2)))
}
func (l QLine) SetLine(x1, y1, x2, y2 int) QLine {
	return qlineFromC(C.qline_set_line(l.toC(), C.int(x1), C.int(y1), C.int(x2), C.int(y2)))
}
func (l QLine) Translate(p QPoint) QLine {
	return qlineFromC(C.qline_translate(l.toC(), qlinePointToC(p)))
}
func (l QLine) TranslateXY(dx, dy int) QLine {
	return qlineFromC(C.qline_translate_xy(l.toC(), C.int(dx), C.int(dy)))
}

// ---- transforms (return new value) --------------------------------------

func (l QLine) Translated(p QPoint) QLine {
	return qlineFromC(C.qline_translated(l.toC(), qlinePointToC(p)))
}
func (l QLine) TranslatedXY(dx, dy int) QLine {
	return qlineFromC(C.qline_translated_xy(l.toC(), C.int(dx), C.int(dy)))
}

// ---- comparison ----------------------------------------------------------

func (l QLine) Equal(other QLine) bool    { return bool(C.qline_equal(l.toC(), other.toC())) }
func (l QLine) NotEqual(other QLine) bool { return !l.Equal(other) }

// ---- conversions ---------------------------------------------------------

// ToLineF converts l to QLineF (exact int→double per coordinate).
func (l QLine) ToLineF() QLineF { return qlinefFromC(C.qline_to_linef(l.toC())) }

// ---- internal helpers ----------------------------------------------------

// toC converts l to its C ABI representation.
func (l QLine) toC() C.go_qline_t {
	return C.go_qline_t{
		x1: C.int(l.x1), y1: C.int(l.y1),
		x2: C.int(l.x2), y2: C.int(l.y2),
	}
}

func qlineFromC(c C.go_qline_t) QLine {
	return QLine{
		x1: int(c.x1), y1: int(c.y1),
		x2: int(c.x2), y2: int(c.y2),
	}
}

// qlinePointToC converts a QPoint value to C.go_qpoint_t using only its public API.
// Local helper: avoids depending on any specific internal helper name from qpoint.go.
func qlinePointToC(p QPoint) C.go_qpoint_t {
	return C.go_qpoint_t{x: C.int(p.X()), y: C.int(p.Y())}
}