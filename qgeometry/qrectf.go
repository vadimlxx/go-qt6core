package qgeometry

/*
 QRectF — биндинг Qt 6.11 QtCore/QRectF как VALUE-тип.

 QRectF — value-тип (Q_RELOCATABLE_TYPE): хранит (xp, yp, w, h) — позицию и размер.
 В Go поля неэкспортируемые, доступ через геттеры Left()/Top()/Right()/Bottom()/
 X()/Y()/Width()/Height(). Присваивание Go-переменной создаёт независимую копию.

 Go-нулевое значение QRectF{} = QRectF() = (0,0,0,0) — null+empty+invalid.
 Поэтому:
   - NewQRectFDefault()     → (0, 0, 0, 0)  — соответствует QRectF()
   - NewQRectFZero()        → (0, 0, 0, 0)  — то же самое (для симметрии с QSizeF)
   - NewQRectF(l, t, w, h)  → (l, t) + (w, h)

 Соглашения:

   - Геттеры (значения): Left/Top/Right/Bottom/X/Y/Width/Height,
     IsNull/IsEmpty/IsValid.
   - Геттеры-значения-структуры: Size() → QSizeF, TopLeft/TopRight/BottomLeft/
     BottomRight/Center() → QPointF.
   - Мутирующие сеттеры (возвращают новое значение): SetLeft/SetTop/SetRight/
     SetBottom/SetX/SetY/SetWidth/SetHeight/SetSize/SetTopLeft/SetTopRight/
     SetBottomLeft/SetBottomRight/SetRect/SetCoords.
   - Movers (возвращают новое значение): MoveLeft/MoveTop/MoveRight/MoveBottom/
     MoveTopLeft/MoveTopRight/MoveBottomLeft/MoveBottomRight/MoveCenter/MoveTo/
     MoveToPoint.
   - Translate/adjust (возвращают новое значение): Translate(dx,dy), TranslatePoint(p),
     Adjust(x1,y1,x2,y2).
   - Трансформации (возвращают новое значение): Normalized, Transposed, Translated(dx,dy),
     TranslatedPoint(p), Adjusted(x1,y1,x2,y2).
   - Объединение/пересечение (возвращают новое значение): United, Intersected.
   - Margins: MarginsAdded(QMarginsF) → QRectF, MarginsRemoved(QMarginsF) → QRectF.
   - Boolean: ContainsRect, ContainsPoint, ContainsXY, Intersects,
     Equal (fuzzy), NotEqual, EqualRect (кросс QRectF == QRect).
   - Conversions: ToRect() → QRect (специальное rounding); ToAlignedRect() → QRect
     (floor top-left, ceil bottom-right).
   - Конструкторы: NewQRectFFromPoints(p1, p2 QPointF), NewQRectFFromPointSize(p1, size),
     QRectFFromRect(r QRect).

 Особенности:

   - QRectF хранит (xp, yp, w, h). Width() = w. IsNull() — fuzzy: w==0 && h==0.
   - IsValid(): w > 0 && h > 0.
   - Equal — fuzzy (qFuzzyCompare обеих точек).
   - ToRect() применяет специальное rounding: nw = round(w + (xp - nxp)/2) и т.п.
   - ToAlignedRect() — floor для top-left, ceil для bottom-right.
   - Protected-секций у QRectF нет.

 Что не реализовано (см. qrect_doc.go):

   - QDataStream operator<< / operator>> — пакет qdatastream пока не готов.
   - QDebug operator<< — пакет qdebug пока не готов.
   - std::tuple protocol — не нужно в Go.
   - fromCGRect / toCGRect (Darwin), fromDOMRect / toDOMRect (WASM) — платформо-специфично.
   - qFuzzyCompare / qFuzzyIsNull как свободные функции — покрыты Equal/IsNull.
   - operator+/- с QMarginsF — покрыто через MarginsAdded/MarginsRemoved.
*/

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC -Wno-sfinae-incomplete
#cgo LDFLAGS: -lstdc++
#include "qrect_c.h"
*/
import "C"

import "structs"

// QRectF wraps Qt6 QRectF (Q_RELOCATABLE_TYPE, (xp,yp,w,h) float64).
type QRectF struct {
	_      structs.HostLayout
	xp, yp, w, h float64
}

// ---- construction --------------------------------------------------------

// NewQRectF returns QRectF(left, top, width, height).
func NewQRectF(left, top, width, height float64) QRectF {
	return qrectfFromC(C.qrectf_new(
		C.double(left), C.double(top), C.double(width), C.double(height)))
}

// NewQRectFZero returns QRectF(0, 0, 0, 0) — same as Qt's default QRectF() and the Go zero value.
func NewQRectFZero() QRectF { return NewQRectF(0, 0, 0, 0) }

// NewQRectFDefault returns QRectF() — same as NewQRectFZero.
func NewQRectFDefault() QRectF { return qrectfFromC(C.qrectf_default()) }

// NewQRectFFromPoints builds QRectF(topLeft, bottomRight).
func NewQRectFFromPoints(topLeft, bottomRight QPointF) QRectF {
	return qrectfFromC(C.qrectf_from_points(qpointfToC(topLeft), qpointfToC(bottomRight)))
}

// NewQRectFFromPointSize builds QRectF(topLeft, size).
func NewQRectFFromPointSize(topLeft QPointF, size QSizeF) QRectF {
	return qrectfFromC(C.qrectf_from_point_size(qpointfToC(topLeft), qsizefToC(size)))
}

// QRectFFromRect builds QRectF from QRect (exact int→double).
func QRectFFromRect(r QRect) QRectF { return qrectfFromC(C.qrectf_from_rect(r.toC())) }

// ---- accessors (values) --------------------------------------------------

func (r QRectF) IsNull() bool    { return bool(C.qrectf_is_null(r.toC())) }
func (r QRectF) IsEmpty() bool   { return bool(C.qrectf_is_empty(r.toC())) }
func (r QRectF) IsValid() bool   { return bool(C.qrectf_is_valid(r.toC())) }
func (r QRectF) Left() float64   { return float64(C.qrectf_left(r.toC())) }
func (r QRectF) Top() float64    { return float64(C.qrectf_top(r.toC())) }
func (r QRectF) Right() float64  { return float64(C.qrectf_right(r.toC())) }
func (r QRectF) Bottom() float64 { return float64(C.qrectf_bottom(r.toC())) }
func (r QRectF) X() float64      { return float64(C.qrectf_x(r.toC())) }
func (r QRectF) Y() float64      { return float64(C.qrectf_y(r.toC())) }
func (r QRectF) Width() float64  { return float64(C.qrectf_width(r.toC())) }
func (r QRectF) Height() float64 { return float64(C.qrectf_height(r.toC())) }

// ---- accessors (packed values) -------------------------------------------

func (r QRectF) Size() QSizeF        { return qsizefFromC(C.qrectf_size(r.toC())) }
func (r QRectF) TopLeft() QPointF     { return qpointfFromC(C.qrectf_top_left(r.toC())) }
func (r QRectF) TopRight() QPointF    { return qpointfFromC(C.qrectf_top_right(r.toC())) }
func (r QRectF) BottomLeft() QPointF  { return qpointfFromC(C.qrectf_bottom_left(r.toC())) }
func (r QRectF) BottomRight() QPointF { return qpointfFromC(C.qrectf_bottom_right(r.toC())) }
func (r QRectF) Center() QPointF      { return qpointfFromC(C.qrectf_center(r.toC())) }

// ---- transforms (return new value) --------------------------------------

func (r QRectF) Normalized() QRectF { return qrectfFromC(C.qrectf_normalized(r.toC())) }
func (r QRectF) Transposed() QRectF { return qrectfFromC(C.qrectf_transposed(r.toC())) }

func (r QRectF) Translated(dx, dy float64) QRectF {
	return qrectfFromC(C.qrectf_translated(r.toC(), C.double(dx), C.double(dy)))
}
func (r QRectF) TranslatedPoint(p QPointF) QRectF {
	return qrectfFromC(C.qrectf_translated_point(r.toC(), qpointfToC(p)))
}
func (r QRectF) Adjusted(dx1, dy1, dx2, dy2 float64) QRectF {
	return qrectfFromC(C.qrectf_adjusted(r.toC(),
		C.double(dx1), C.double(dy1), C.double(dx2), C.double(dy2)))
}

// ---- setters (return updated value) -------------------------------------

func (r QRectF) SetLeft(pos float64) QRectF {
	return qrectfFromC(C.qrectf_set_left(r.toC(), C.double(pos)))
}
func (r QRectF) SetTop(pos float64) QRectF {
	return qrectfFromC(C.qrectf_set_top(r.toC(), C.double(pos)))
}
func (r QRectF) SetRight(pos float64) QRectF {
	return qrectfFromC(C.qrectf_set_right(r.toC(), C.double(pos)))
}
func (r QRectF) SetBottom(pos float64) QRectF {
	return qrectfFromC(C.qrectf_set_bottom(r.toC(), C.double(pos)))
}
func (r QRectF) SetX(pos float64) QRectF {
	return qrectfFromC(C.qrectf_set_x(r.toC(), C.double(pos)))
}
func (r QRectF) SetY(pos float64) QRectF {
	return qrectfFromC(C.qrectf_set_y(r.toC(), C.double(pos)))
}
func (r QRectF) SetWidth(w float64) QRectF {
	return qrectfFromC(C.qrectf_set_width(r.toC(), C.double(w)))
}
func (r QRectF) SetHeight(h float64) QRectF {
	return qrectfFromC(C.qrectf_set_height(r.toC(), C.double(h)))
}
func (r QRectF) SetSize(s QSizeF) QRectF {
	return qrectfFromC(C.qrectf_set_size(r.toC(), qsizefToC(s)))
}

func (r QRectF) SetTopLeft(p QPointF) QRectF {
	return qrectfFromC(C.qrectf_set_top_left(r.toC(), qpointfToC(p)))
}
func (r QRectF) SetTopRight(p QPointF) QRectF {
	return qrectfFromC(C.qrectf_set_top_right(r.toC(), qpointfToC(p)))
}
func (r QRectF) SetBottomLeft(p QPointF) QRectF {
	return qrectfFromC(C.qrectf_set_bottom_left(r.toC(), qpointfToC(p)))
}
func (r QRectF) SetBottomRight(p QPointF) QRectF {
	return qrectfFromC(C.qrectf_set_bottom_right(r.toC(), qpointfToC(p)))
}
func (r QRectF) SetRect(x, y, w, h float64) QRectF {
	return qrectfFromC(C.qrectf_set_rect(r.toC(),
		C.double(x), C.double(y), C.double(w), C.double(h)))
}
func (r QRectF) SetCoords(x1, y1, x2, y2 float64) QRectF {
	return qrectfFromC(C.qrectf_set_coords(r.toC(),
		C.double(x1), C.double(y1), C.double(x2), C.double(y2)))
}

// ---- movers (return updated value) --------------------------------------

func (r QRectF) MoveLeft(pos float64) QRectF {
	return qrectfFromC(C.qrectf_move_left(r.toC(), C.double(pos)))
}
func (r QRectF) MoveTop(pos float64) QRectF {
	return qrectfFromC(C.qrectf_move_top(r.toC(), C.double(pos)))
}
func (r QRectF) MoveRight(pos float64) QRectF {
	return qrectfFromC(C.qrectf_move_right(r.toC(), C.double(pos)))
}
func (r QRectF) MoveBottom(pos float64) QRectF {
	return qrectfFromC(C.qrectf_move_bottom(r.toC(), C.double(pos)))
}

func (r QRectF) MoveTopLeft(p QPointF) QRectF {
	return qrectfFromC(C.qrectf_move_top_left(r.toC(), qpointfToC(p)))
}
func (r QRectF) MoveTopRight(p QPointF) QRectF {
	return qrectfFromC(C.qrectf_move_top_right(r.toC(), qpointfToC(p)))
}
func (r QRectF) MoveBottomLeft(p QPointF) QRectF {
	return qrectfFromC(C.qrectf_move_bottom_left(r.toC(), qpointfToC(p)))
}
func (r QRectF) MoveBottomRight(p QPointF) QRectF {
	return qrectfFromC(C.qrectf_move_bottom_right(r.toC(), qpointfToC(p)))
}
func (r QRectF) MoveCenter(p QPointF) QRectF {
	return qrectfFromC(C.qrectf_move_center(r.toC(), qpointfToC(p)))
}
func (r QRectF) MoveTo(x, y float64) QRectF {
	return qrectfFromC(C.qrectf_move_to(r.toC(), C.double(x), C.double(y)))
}
func (r QRectF) MoveToPoint(p QPointF) QRectF {
	return qrectfFromC(C.qrectf_move_to_point(r.toC(), qpointfToC(p)))
}

// ---- translate/adjust (return updated value) ----------------------------

func (r QRectF) Translate(dx, dy float64) QRectF {
	return qrectfFromC(C.qrectf_translate(r.toC(), C.double(dx), C.double(dy)))
}
func (r QRectF) TranslatePoint(p QPointF) QRectF {
	return qrectfFromC(C.qrectf_translate_point(r.toC(), qpointfToC(p)))
}
func (r QRectF) Adjust(dx1, dy1, dx2, dy2 float64) QRectF {
	return qrectfFromC(C.qrectf_adjust(r.toC(),
		C.double(dx1), C.double(dy1), C.double(dx2), C.double(dy2)))
}

// ---- boolean ops ---------------------------------------------------------

func (r QRectF) ContainsRect(other QRectF) bool {
	return bool(C.qrectf_contains_rect(r.toC(), other.toC()))
}
func (r QRectF) ContainsPoint(p QPointF) bool {
	return bool(C.qrectf_contains_point(r.toC(), qpointfToC(p)))
}
func (r QRectF) ContainsXY(x, y float64) bool {
	return bool(C.qrectf_contains_xy(r.toC(), C.double(x), C.double(y)))
}
func (r QRectF) Intersects(other QRectF) bool { return bool(C.qrectf_intersects(r.toC(), other.toC())) }
func (r QRectF) Equal(other QRectF) bool       { return bool(C.qrectf_equal(r.toC(), other.toC())) }
func (r QRectF) NotEqual(other QRectF) bool    { return !r.Equal(other) }

// EqualRect compares r (fuzzy) to a QRect, treating it as QRectF.
func (r QRectF) EqualRect(other QRect) bool {
	return bool(C.qrectf_equal_rect(r.toC(), other.toC()))
}

// ---- unions/intersections (new values) ----------------------------------

func (r QRectF) United(other QRectF) QRectF {
	return qrectfFromC(C.qrectf_united(r.toC(), other.toC()))
}
func (r QRectF) Intersected(other QRectF) QRectF {
	return qrectfFromC(C.qrectf_intersected(r.toC(), other.toC()))
}

// ---- margins (new values) ----------------------------------------------

func (r QRectF) MarginsAdded(m QMarginsF) QRectF {
	return qrectfFromC(C.qrectf_margins_added(r.toC(), marginsFToC(m)))
}
func (r QRectF) MarginsRemoved(m QMarginsF) QRectF {
	return qrectfFromC(C.qrectf_margins_removed(r.toC(), marginsFToC(m)))
}

// ---- conversions ---------------------------------------------------------

// ToRect converts r to QRect using Qt's rounding (minimizes top-left and size error).
func (r QRectF) ToRect() QRect { return qrectFromC(C.qrectf_to_rect(r.toC())) }

// ToAlignedRect converts r to QRect by floor(topLeft) and ceil(bottomRight).
func (r QRectF) ToAlignedRect() QRect { return qrectFromC(C.qrectf_to_aligned_rect(r.toC())) }

// ---- internal helpers ----------------------------------------------------

func (r QRectF) toC() C.go_qrectf_t {
	return C.go_qrectf_t{
		xp: C.double(r.xp), yp: C.double(r.yp),
		w: C.double(r.w), h: C.double(r.h),
	}
}

func qrectfFromC(c C.go_qrectf_t) QRectF {
	return QRectF{
		xp: float64(c.xp), yp: float64(c.yp),
		w: float64(c.w), h: float64(c.h),
	}
}
