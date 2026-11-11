package qgeometry

/*
 QMargins — биндинг Qt 6.11 QtCore/QMargins.

 QMargins — value-тип (Q_RELOCATABLE_TYPE): 4 int-поля (left, top, right, bottom).
 В Go поля неэкспортируемые, доступ через геттеры Left()/Top()/Right()/Bottom()
 (Qt-style, без Get-префикса). Через C ABI передаётся плоской структурой
 go_qmargins_t (4×int, 16 байт).

 Соглашения:

   - Геттеры: Left(), Top(), Right(), Bottom(), IsNull().
   - Мутирующие сеттеры: SetLeft(v), SetTop(v), SetRight(v), SetBottom(v) —
     на *QMargins.
   - Арифметика: Add(other), Sub(other), AddInt(v), SubInt(v), MulInt(f),
     DivInt(d), MulFloat(f), DivFloat(d), Neg(), UnaryPlus().
   - Компонентный максимум: Or(other) — соответствует Qt operator|.
   - Присваивающие формы: AddAssign, SubAssign, AddIntAssign, SubIntAssign,
     MulIntAssign, DivIntAssign, MulFloatAssign, DivFloatAssign.
   - Сравнение: Equal(other) / NotEqual(other). Также работает `==` внутри
     пакета qgeometry; для внешних пользователей предпочтителен Equal.
   - Конвертация: ToMarginsF() → QMarginsF.

 Особенности:

   - DivInt при d == 0 — C++ UB (обычно SIGFPE). В Go — паника.
   - MulInt/DivInt переполняются как C++ int (без saturating) — совпадает
     с release-сборкой Qt 6.11. Для saturating используйте MulFloat.
   - MulFloat/DivFloat используют Qt's qSaturateRound — идентично C++.
   - Сравнение QMargins с QMarginsF — через m.ToMarginsF().Equal(other)
     или other.EqualQMargins(m).
   - Protected-секций у QMargins нет — открывать нечего.

 Что не реализовано (см. qmargins_doc.go):

   - QDataStream operator<< / operator>> — пакет qdatastream пока не готов.
   - QDebug operator<< — пакет qdebug пока не готов.
   - std::tuple protocol (structured bindings) — не нужно в Go.
*/

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC -Wno-sfinae-incomplete
#cgo LDFLAGS: -lstdc++
#include "qmargins_c.h"
*/
import "C"

// QMargins wraps Qt6 QMargins (Q_RELOCATABLE_TYPE, 4×int value).
// Fields are unexported; use Left()/Top()/Right()/Bottom() to read,
// SetLeft()/… to mutate.
type QMargins struct {
	left   int
	top    int
	right  int
	bottom int
}

// ---- construction --------------------------------------------------------

// NewQMargins returns QMargins(left, top, right, bottom).
func NewQMargins(left, top, right, bottom int) QMargins {
	return marginsFromC(C.qmargins_new(C.int(left), C.int(top), C.int(right), C.int(bottom)))
}

// NewQMarginsZero returns the default-constructed QMargins (0, 0, 0, 0).
func NewQMarginsZero() QMargins {
	return marginsFromC(C.qmargins_default())
}

// ---- accessors -----------------------------------------------------------

func (m QMargins) IsNull() bool { return bool(C.qmargins_is_null(marginsToC(m))) }
func (m QMargins) Left() int    { return int(C.qmargins_left(marginsToC(m))) }
func (m QMargins) Top() int     { return int(C.qmargins_top(marginsToC(m))) }
func (m QMargins) Right() int   { return int(C.qmargins_right(marginsToC(m))) }
func (m QMargins) Bottom() int  { return int(C.qmargins_bottom(marginsToC(m))) }

// ---- modifiers -----------------------------------------------------------

func (m *QMargins) SetLeft(v int)   { *m = marginsFromC(C.qmargins_set_left(marginsToC(*m), C.int(v))) }
func (m *QMargins) SetTop(v int)    { *m = marginsFromC(C.qmargins_set_top(marginsToC(*m), C.int(v))) }
func (m *QMargins) SetRight(v int)  { *m = marginsFromC(C.qmargins_set_right(marginsToC(*m), C.int(v))) }
func (m *QMargins) SetBottom(v int) { *m = marginsFromC(C.qmargins_set_bottom(marginsToC(*m), C.int(v))) }

// ---- arithmetic ----------------------------------------------------------

// Add returns m + other (component-wise).
func (m QMargins) Add(other QMargins) QMargins {
	return marginsFromC(C.qmargins_add(marginsToC(m), marginsToC(other)))
}

// Sub returns m - other (component-wise).
func (m QMargins) Sub(other QMargins) QMargins {
	return marginsFromC(C.qmargins_sub(marginsToC(m), marginsToC(other)))
}

// AddInt returns m + v on each component.
func (m QMargins) AddInt(v int) QMargins {
	return marginsFromC(C.qmargins_add_int(marginsToC(m), C.int(v)))
}

// SubInt returns m - v on each component.
func (m QMargins) SubInt(v int) QMargins {
	return marginsFromC(C.qmargins_sub_int(marginsToC(m), C.int(v)))
}

// MulInt returns m * f on each component (C++ int overflow semantics).
func (m QMargins) MulInt(f int) QMargins {
	return marginsFromC(C.qmargins_mul_int(marginsToC(m), C.int(f)))
}

// DivInt returns m / d on each component. Panics if d == 0 (C++ UB → Go panic).
func (m QMargins) DivInt(d int) QMargins {
	return marginsFromC(C.qmargins_div_int(marginsToC(m), C.int(d)))
}

// MulFloat returns m * f on each component (rounded via Qt's qSaturateRound).
func (m QMargins) MulFloat(f float64) QMargins {
	return marginsFromC(C.qmargins_mul_float(marginsToC(m), C.double(f)))
}

// DivFloat returns m / d on each component (rounded via Qt's qSaturateRound).
func (m QMargins) DivFloat(d float64) QMargins {
	return marginsFromC(C.qmargins_div_float(marginsToC(m), C.double(d)))
}

// Or returns component-wise max of m and other (Qt operator|).
func (m QMargins) Or(other QMargins) QMargins {
	return marginsFromC(C.qmargins_or(marginsToC(m), marginsToC(other)))
}

// Neg returns -m (component-wise).
func (m QMargins) Neg() QMargins { return marginsFromC(C.qmargins_neg(marginsToC(m))) }

// UnaryPlus returns +m (copy).
func (m QMargins) UnaryPlus() QMargins { return marginsFromC(C.qmargins_unary_plus(marginsToC(m))) }

// ---- assign forms --------------------------------------------------------

func (m *QMargins) AddAssign(other QMargins) { *m = m.Add(other) }
func (m *QMargins) SubAssign(other QMargins) { *m = m.Sub(other) }
func (m *QMargins) AddIntAssign(v int)       { *m = m.AddInt(v) }
func (m *QMargins) SubIntAssign(v int)       { *m = m.SubInt(v) }
func (m *QMargins) MulIntAssign(f int)       { *m = m.MulInt(f) }
func (m *QMargins) DivIntAssign(d int)       { *m = m.DivInt(d) }
func (m *QMargins) MulFloatAssign(f float64) { *m = m.MulFloat(f) }
func (m *QMargins) DivFloatAssign(d float64) { *m = m.DivFloat(d) }

// ---- comparison ----------------------------------------------------------

// Equal returns true if all four int fields are identical (Qt comparesEqual).
func (m QMargins) Equal(other QMargins) bool {
	return bool(C.qmargins_equal(marginsToC(m), marginsToC(other)))
}

// NotEqual returns !Equal(other).
func (m QMargins) NotEqual(other QMargins) bool { return !m.Equal(other) }

// ---- conversions ---------------------------------------------------------

// ToMarginsF converts m to QMarginsF (exact int→double per component).
func (m QMargins) ToMarginsF() QMarginsF {
	return marginsFFromC(C.qmargins_to_marginsf(marginsToC(m)))
}

// ---- internal helpers ----------------------------------------------------

func marginsToC(m QMargins) C.go_qmargins_t {
	return C.go_qmargins_t{
		left:   C.int(m.left),
		top:    C.int(m.top),
		right:  C.int(m.right),
		bottom: C.int(m.bottom),
	}
}

func marginsFromC(c C.go_qmargins_t) QMargins {
	return QMargins{
		left:   int(c.left),
		top:    int(c.top),
		right:  int(c.right),
		bottom: int(c.bottom),
	}
}