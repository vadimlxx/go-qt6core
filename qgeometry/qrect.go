package qgeometry

/*
 QRect — биндинг Qt 6.11 QtCore/QRect как VALUE-тип.

 QRect — value-тип (Q_RELOCATABLE_TYPE): хранит (x1, y1, x2, y2) — координаты
 верхнего-левого и нижнего-правого углов (оба включительно). В Go поля
 неэкспортируемые, доступ через геттеры Left()/Top()/Right()/Bottom()/X()/Y()/
 Width()/Height(). Присваивание Go-переменной создаёт независимую копию.

 ВАЖНО: Go-нулевое значение QRect{} = (0,0,0,0), а Qt-конструктор QRect() даёт
 QRect(0, 0, -1, -1) — null+empty+invalid. Они НЕ эквивалентны:
   - NewQRectDefault() → (0, 0, -1, -1) — null, empty, invalid
   - NewQRectZero()    → (0, 0,  0,  0) — валидный (0×0), IsNull() = false
   - NewQRect(l, t, w, h) — параметры left/top/width/height (height, не bottom!)

 Соглашения:

   - Геттеры (значения): Left/Top/Right/Bottom/X/Y/Width/Height, IsNull/IsEmpty/IsValid.
   - Геттеры-значения-структуры: Size() → QSize, TopLeft/TopRight/BottomLeft/BottomRight/
     Center() → QPoint.
   - Мутирующие сеттеры (возвращают новое значение): SetLeft/SetTop/SetRight/SetBottom/
     SetX/SetY/SetWidth/SetHeight/SetSize/SetTopLeft/SetTopRight/SetBottomLeft/
     SetBottomRight/SetRect/SetCoords.
   - Movers (возвращают новое значение): MoveLeft/MoveTop/MoveRight/MoveBottom/MoveTopLeft/
     MoveTopRight/MoveBottomLeft/MoveBottomRight/MoveCenter/MoveTo/MoveToPoint.
   - Translate/adjust (возвращают новое значение): Translate(dx,dy), TranslatePoint(p),
     Adjust(x1,y1,x2,y2).
   - Трансформации (возвращают новое значение): Normalized, Transposed, Translated(dx,dy),
     TranslatedPoint(p), Adjusted(x1,y1,x2,y2).
   - Объединение/пересечение (возвращают новое значение): United, Intersected.
   - Margins: MarginsAdded(QMargins) → QRect, MarginsRemoved(QMargins) → QRect.
   - Boolean: ContainsRect, ContainsPoint(proper), ContainsXY, Intersects, Equal, NotEqual.
   - Conversions: ToRectF() → QRectF.
   - Static: QRectSpan(p1, p2 QPoint) → QRect.

 Особенности:

   - QRect хранит (x1, y1, x2, y2), а не (x, y, w, h). Width() = x2-x1+1.
     SetWidth(w) пересчитывает right = left + w - 1. Семантика Qt.
   - IsNull() для default = true (x2 = x1-1, y2 = y1-1). Для (0,0,0,0) — false.
   - Center() использует 64-битное сложение (без overflow).
   - Protected-секций у QRect нет.

 Что не реализовано (см. qrect_doc.go):

   - QDataStream operator<< / operator>> — пакет qdatastream пока не готов.
   - QDebug operator<< — пакет qdebug пока не готов.
   - qHash(QRect) — не нужно в Go.
   - std::tuple protocol — не нужно в Go.
   - toCGRect (Darwin) — платформо-специфично.
   - operator+/- с QMargins — покрыто через MarginsAdded/MarginsRemoved.
*/

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC -Wno-sfinae-incomplete
#cgo LDFLAGS: -lstdc++
#include "qrect_c.h"
*/
import "C"

// QRect wraps Qt6 QRect (Q_RELOCATABLE_TYPE, (x1,y1,x2,y2) ints).
type QRect struct {
	x1, y1, x2, y2 int
}

// ---- construction --------------------------------------------------------

// NewQRect returns QRect(left, top, width, height).
// NOTE: fourth parameter is HEIGHT, not bottom. Matches QRect(int,int,int,int).
func NewQRect(left, top, width, height int) QRect {
	return qrectFromC(C.qrect_new(C.int(left), C.int(top), C.int(width), C.int(height)))
}

// NewQRectZero returns QRect(0, 0, 0, 0) — a valid zero-sized rect (IsNull false).
func NewQRectZero() QRect { return NewQRect(0, 0, 0, 0) }

// NewQRectDefault returns QRect(0, 0, -1, -1), matching Qt's default constructor.
func NewQRectDefault() QRect { return qrectFromC(C.qrect_default()) }

// NewQRectFromPoints builds QRect(topLeft, bottomRight). Values are copied from the points.
func NewQRectFromPoints(topLeft, bottomRight QPoint) QRect {
	return qrectFromC(C.qrect_from_points(qpointToC(topLeft), qpointToC(bottomRight)))
}

// NewQRectFromPointSize builds QRect(topLeft, size). Values are copied.
func NewQRectFromPointSize(topLeft QPoint, size QSize) QRect {
	return qrectFromC(C.qrect_from_point_size(qpointToC(topLeft), qsizeToC(size)))
}

// ---- accessors (values) --------------------------------------------------

func (r QRect) IsNull() bool  { return bool(C.qrect_is_null(r.toC())) }
func (r QRect) IsEmpty() bool { return bool(C.qrect_is_empty(r.toC())) }
func (r QRect) IsValid() bool { return bool(C.qrect_is_valid(r.toC())) }
func (r QRect) Left() int     { return int(C.qrect_left(r.toC())) }
func (r QRect) Top() int      { return int(C.qrect_top(r.toC())) }
func (r QRect) Right() int    { return int(C.qrect_right(r.toC())) }
func (r QRect) Bottom() int   { return int(C.qrect_bottom(r.toC())) }
func (r QRect) X() int        { return int(C.qrect_x(r.toC())) }
func (r QRect) Y() int        { return int(C.qrect_y(r.toC())) }
func (r QRect) Width() int    { return int(C.qrect_width(r.toC())) }
func (r QRect) Height() int   { return int(C.qrect_height(r.toC())) }

// ---- accessors (packed values) -------------------------------------------

func (r QRect) Size() QSize        { return qsizeFromC(C.qrect_size(r.toC())) }
func (r QRect) TopLeft() QPoint     { return qpointFromC(C.qrect_top_left(r.toC())) }
func (r QRect) TopRight() QPoint    { return qpointFromC(C.qrect_top_right(r.toC())) }
func (r QRect) BottomLeft() QPoint  { return qpointFromC(C.qrect_bottom_left(r.toC())) }
func (r QRect) BottomRight() QPoint { return qpointFromC(C.qrect_bottom_right(r.toC())) }
func (r QRect) Center() QPoint      { return qpointFromC(C.qrect_center(r.toC())) }

// ---- transforms (return new value) --------------------------------------

func (r QRect) Normalized() QRect { return qrectFromC(C.qrect_normalized(r.toC())) }
func (r QRect) Transposed() QRect { return qrectFromC(C.qrect_transposed(r.toC())) }

func (r QRect) Translated(dx, dy int) QRect {
	return qrectFromC(C.qrect_translated(r.toC(), C.int(dx), C.int(dy)))
}
func (r QRect) TranslatedPoint(p QPoint) QRect {
	return qrectFromC(C.qrect_translated_point(r.toC(), qpointToC(p)))
}
func (r QRect) Adjusted(dx1, dy1, dx2, dy2 int) QRect {
	return qrectFromC(C.qrect_adjusted(r.toC(),
		C.int(dx1), C.int(dy1), C.int(dx2), C.int(dy2)))
}

// ---- setters (return updated value) -------------------------------------

func (r QRect) SetLeft(pos int) QRect   { return qrectFromC(C.qrect_set_left(r.toC(), C.int(pos))) }
func (r QRect) SetTop(pos int) QRect    { return qrectFromC(C.qrect_set_top(r.toC(), C.int(pos))) }
func (r QRect) SetRight(pos int) QRect  { return qrectFromC(C.qrect_set_right(r.toC(), C.int(pos))) }
func (r QRect) SetBottom(pos int) QRect { return qrectFromC(C.qrect_set_bottom(r.toC(), C.int(pos))) }
func (r QRect) SetX(pos int) QRect      { return qrectFromC(C.qrect_set_x(r.toC(), C.int(pos))) }
func (r QRect) SetY(pos int) QRect      { return qrectFromC(C.qrect_set_y(r.toC(), C.int(pos))) }
func (r QRect) SetWidth(w int) QRect    { return qrectFromC(C.qrect_set_width(r.toC(), C.int(w))) }
func (r QRect) SetHeight(h int) QRect   { return qrectFromC(C.qrect_set_height(r.toC(), C.int(h))) }
func (r QRect) SetSize(s QSize) QRect   { return qrectFromC(C.qrect_set_size(r.toC(), qsizeToC(s))) }

func (r QRect) SetTopLeft(p QPoint) QRect {
	return qrectFromC(C.qrect_set_top_left(r.toC(), qpointToC(p)))
}
func (r QRect) SetTopRight(p QPoint) QRect {
	return qrectFromC(C.qrect_set_top_right(r.toC(), qpointToC(p)))
}
func (r QRect) SetBottomLeft(p QPoint) QRect {
	return qrectFromC(C.qrect_set_bottom_left(r.toC(), qpointToC(p)))
}
func (r QRect) SetBottomRight(p QPoint) QRect {
	return qrectFromC(C.qrect_set_bottom_right(r.toC(), qpointToC(p)))
}
func (r QRect) SetRect(x, y, w, h int) QRect {
	return qrectFromC(C.qrect_set_rect(r.toC(), C.int(x), C.int(y), C.int(w), C.int(h)))
}
func (r QRect) SetCoords(x1, y1, x2, y2 int) QRect {
	return qrectFromC(C.qrect_set_coords(r.toC(), C.int(x1), C.int(y1), C.int(x2), C.int(y2)))
}

// ---- movers (return updated value) --------------------------------------

func (r QRect) MoveLeft(pos int) QRect   { return qrectFromC(C.qrect_move_left(r.toC(), C.int(pos))) }
func (r QRect) MoveTop(pos int) QRect    { return qrectFromC(C.qrect_move_top(r.toC(), C.int(pos))) }
func (r QRect) MoveRight(pos int) QRect  { return qrectFromC(C.qrect_move_right(r.toC(), C.int(pos))) }
func (r QRect) MoveBottom(pos int) QRect { return qrectFromC(C.qrect_move_bottom(r.toC(), C.int(pos))) }

func (r QRect) MoveTopLeft(p QPoint) QRect {
	return qrectFromC(C.qrect_move_top_left(r.toC(), qpointToC(p)))
}
func (r QRect) MoveTopRight(p QPoint) QRect {
	return qrectFromC(C.qrect_move_top_right(r.toC(), qpointToC(p)))
}
func (r QRect) MoveBottomLeft(p QPoint) QRect {
	return qrectFromC(C.qrect_move_bottom_left(r.toC(), qpointToC(p)))
}
func (r QRect) MoveBottomRight(p QPoint) QRect {
	return qrectFromC(C.qrect_move_bottom_right(r.toC(), qpointToC(p)))
}
func (r QRect) MoveCenter(p QPoint) QRect {
	return qrectFromC(C.qrect_move_center(r.toC(), qpointToC(p)))
}
func (r QRect) MoveTo(x, y int) QRect {
	return qrectFromC(C.qrect_move_to(r.toC(), C.int(x), C.int(y)))
}
func (r QRect) MoveToPoint(p QPoint) QRect {
	return qrectFromC(C.qrect_move_to_point(r.toC(), qpointToC(p)))
}

// ---- translate/adjust (return updated value) ----------------------------

func (r QRect) Translate(dx, dy int) QRect {
	return qrectFromC(C.qrect_translate(r.toC(), C.int(dx), C.int(dy)))
}
func (r QRect) TranslatePoint(p QPoint) QRect {
	return qrectFromC(C.qrect_translate_point(r.toC(), qpointToC(p)))
}
func (r QRect) Adjust(dx1, dy1, dx2, dy2 int) QRect {
	return qrectFromC(C.qrect_adjust(r.toC(),
		C.int(dx1), C.int(dy1), C.int(dx2), C.int(dy2)))
}

// ---- boolean ops ---------------------------------------------------------

func (r QRect) ContainsRect(other QRect, proper bool) bool {
	return bool(C.qrect_contains_rect(r.toC(), other.toC(), C.bool(proper)))
}
func (r QRect) ContainsPoint(p QPoint, proper bool) bool {
	return bool(C.qrect_contains_point(r.toC(), qpointToC(p), C.bool(proper)))
}
func (r QRect) ContainsXY(x, y int) bool {
	return bool(C.qrect_contains_xy(r.toC(), C.int(x), C.int(y)))
}
func (r QRect) Intersects(other QRect) bool { return bool(C.qrect_intersects(r.toC(), other.toC())) }
func (r QRect) Equal(other QRect) bool       { return bool(C.qrect_equal(r.toC(), other.toC())) }
func (r QRect) NotEqual(other QRect) bool    { return !r.Equal(other) }

// ---- unions/intersections (new values) ----------------------------------

func (r QRect) United(other QRect) QRect      { return qrectFromC(C.qrect_united(r.toC(), other.toC())) }
func (r QRect) Intersected(other QRect) QRect { return qrectFromC(C.qrect_intersected(r.toC(), other.toC())) }

// ---- margins (new values) ----------------------------------------------

func (r QRect) MarginsAdded(m QMargins) QRect {
	return qrectFromC(C.qrect_margins_added(r.toC(), marginsToC(m)))
}
func (r QRect) MarginsRemoved(m QMargins) QRect {
	return qrectFromC(C.qrect_margins_removed(r.toC(), marginsToC(m)))
}

// ---- static --------------------------------------------------------------

// QRectSpan returns the bounding QRect spanning p1 and p2. Mirrors the static QRect::span(p1, p2).
func QRectSpan(p1, p2 QPoint) QRect {
	return qrectFromC(C.qrect_span(qpointToC(p1), qpointToC(p2)))
}

// ---- conversions ---------------------------------------------------------

// ToRectF converts r to QRectF (exact int→double per endpoint).
func (r QRect) ToRectF() QRectF { return qrectfFromC(C.qrect_to_rectf(r.toC())) }

// ---- internal helpers ----------------------------------------------------

func (r QRect) toC() C.go_qrect_t {
	return C.go_qrect_t{
		x1: C.int(r.x1), y1: C.int(r.y1),
		x2: C.int(r.x2), y2: C.int(r.y2),
	}
}

func qrectFromC(c C.go_qrect_t) QRect {
	return QRect{
		x1: int(c.x1), y1: int(c.y1),
		x2: int(c.x2), y2: int(c.y2),
	}
}
