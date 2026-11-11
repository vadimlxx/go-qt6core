package qgeometry

/*
 QMargins — биндинг Qt 6.11 QtCore/QMargins как HANDLE.

 QMargins — value-тип в Qt (Q_RELOCATABLE_TYPE, 4×int). В Go хранится handle на heap-размещённый
 C++ QMargins. Все геттеры/сеттеры идут напрямую в C++.

 ВАЖНОЕ ОГРАНИЧЕНИЕ:
   Qt хранит QMargins внутри других классов ПО ЗНАЧЕНИЮ. Handle на QMargins не видит изменений,
   сделанных C++ во встроенном QMargins другого класса. Это ограничение API Qt, а не биндинга.

 Семантика значений vs ссылок:
   - QMargins — ссылочный тип (handle). Мутации через один handle видны через все handle'ы на
     тот же объект. Для независимой копии — Clone().
   - Go-нулевое значение QMargins{} невалидно. Используйте NewQMargins / NewQMarginsZero.

 Управление памятью:
   - NewQMargins / NewQMarginsZero / арифметика / Clone / ToMarginsF создают НОВЫЙ heap-объект.
   - Close() освобождает немедленно; иначе срабатывает finalizer.

 Соглашения:
   - Геттеры: Left(), Top(), Right(), Bottom(), IsNull().
   - Мутирующие (in place, void): SetLeft, SetTop, SetRight, SetBottom.
   - Арифметика (новый handle): Add, Sub, AddInt, SubInt, MulInt, DivInt, MulFloat, DivFloat,
     Or, Neg, UnaryPlus.
   - Присваивающие формы (in place): AddAssign, SubAssign, AddIntAssign, SubIntAssign,
     MulIntAssign, DivIntAssign, MulFloatAssign, DivFloatAssign.
   - Сравнение: Equal(other *QMargins), NotEqual(other).
   - Конвертация: ToMarginsF() → *QMarginsF.
   - Клонирование: Clone() → *QMargins.

 Особенности:
   - DivInt при d == 0 — C++ UB (обычно SIGFPE). В Go — паника.
   - MulInt/DivInt переполняются как C++ int (QCheckedInt ассертит в debug, wraps в release).
   - MulFloat/DivFloat используют Qt's qSaturateRound.
   - Protected-секций у QMargins нет.

 Что не реализовано (см. qmargins_doc.go):

   - QDataStream operator<< / operator>> — пакет qdatastream пока не готов.
   - QDebug operator<< — пакет qdebug пока не готов.
   - std::tuple protocol — не нужно в Go.
*/

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC -Wno-sfinae-incomplete
#cgo LDFLAGS: -lstdc++
#include "qmargins_c.h"
*/
import "C"

import "runtime"

// QMargins wraps a heap-allocated Qt6 QMargins.
type QMargins struct {
	p C.go_qmargins_t
}

func wrapQMargins(p C.go_qmargins_t) *QMargins {
	if p == nil {
		return nil
	}
	m := &QMargins{p: p}
	runtime.SetFinalizer(m, finalizeQMargins)
	return m
}

func finalizeQMargins(m *QMargins) {
	if m.p != nil {
		C.qmargins_delete(m.p)
		m.p = nil
	}
}

// Close releases the underlying C++ object immediately and disables the finalizer.
func (m *QMargins) Close() {
	if m.p != nil {
		C.qmargins_delete(m.p)
		m.p = nil
		runtime.SetFinalizer(m, nil)
	}
}

// ---- construction --------------------------------------------------------

// NewQMargins returns a handle to QMargins(left, top, right, bottom).
func NewQMargins(left, top, right, bottom int) *QMargins {
	return wrapQMargins(C.qmargins_new(C.int(left), C.int(top), C.int(right), C.int(bottom)))
}

// NewQMarginsZero returns a handle to QMargins(0, 0, 0, 0) — same as Qt's default QMargins().
func NewQMarginsZero() *QMargins {
	return wrapQMargins(C.qmargins_new_default())
}

// Clone returns a heap copy of m (independent C++ object).
func (m *QMargins) Clone() *QMargins {
	return wrapQMargins(C.qmargins_clone(m.p))
}

// ---- accessors -----------------------------------------------------------

func (m *QMargins) IsNull() bool { return bool(C.qmargins_is_null(m.p)) }
func (m *QMargins) Left() int    { return int(C.qmargins_left(m.p)) }
func (m *QMargins) Top() int     { return int(C.qmargins_top(m.p)) }
func (m *QMargins) Right() int   { return int(C.qmargins_right(m.p)) }
func (m *QMargins) Bottom() int  { return int(C.qmargins_bottom(m.p)) }

// ---- mutating setters (in place) ----------------------------------------

func (m *QMargins) SetLeft(v int)   { C.qmargins_set_left(m.p, C.int(v)) }
func (m *QMargins) SetTop(v int)    { C.qmargins_set_top(m.p, C.int(v)) }
func (m *QMargins) SetRight(v int)  { C.qmargins_set_right(m.p, C.int(v)) }
func (m *QMargins) SetBottom(v int) { C.qmargins_set_bottom(m.p, C.int(v)) }

// ---- arithmetic (return new handles) ------------------------------------

func (m *QMargins) Add(other *QMargins) *QMargins {
	return wrapQMargins(C.qmargins_add(m.p, other.p))
}
func (m *QMargins) Sub(other *QMargins) *QMargins {
	return wrapQMargins(C.qmargins_sub(m.p, other.p))
}
func (m *QMargins) AddInt(v int) *QMargins {
	return wrapQMargins(C.qmargins_add_int(m.p, C.int(v)))
}
func (m *QMargins) SubInt(v int) *QMargins {
	return wrapQMargins(C.qmargins_sub_int(m.p, C.int(v)))
}
func (m *QMargins) MulInt(f int) *QMargins {
	return wrapQMargins(C.qmargins_mul_int(m.p, C.int(f)))
}
func (m *QMargins) DivInt(d int) *QMargins {
	return wrapQMargins(C.qmargins_div_int(m.p, C.int(d)))
}
func (m *QMargins) MulFloat(f float64) *QMargins {
	return wrapQMargins(C.qmargins_mul_float(m.p, C.double(f)))
}
func (m *QMargins) DivFloat(d float64) *QMargins {
	return wrapQMargins(C.qmargins_div_float(m.p, C.double(d)))
}
func (m *QMargins) Or(other *QMargins) *QMargins {
	return wrapQMargins(C.qmargins_or(m.p, other.p))
}
func (m *QMargins) Neg() *QMargins       { return wrapQMargins(C.qmargins_neg(m.p)) }
func (m *QMargins) UnaryPlus() *QMargins { return wrapQMargins(C.qmargins_unary_plus(m.p)) }

// ---- in-place assign (mutate receiver) ----------------------------------

func (m *QMargins) AddAssign(other *QMargins)   { C.qmargins_add_assign(m.p, other.p) }
func (m *QMargins) SubAssign(other *QMargins)   { C.qmargins_sub_assign(m.p, other.p) }
func (m *QMargins) AddIntAssign(v int)          { C.qmargins_add_int_assign(m.p, C.int(v)) }
func (m *QMargins) SubIntAssign(v int)          { C.qmargins_sub_int_assign(m.p, C.int(v)) }
func (m *QMargins) MulIntAssign(f int)          { C.qmargins_mul_int_assign(m.p, C.int(f)) }
func (m *QMargins) DivIntAssign(d int)          { C.qmargins_div_int_assign(m.p, C.int(d)) }
func (m *QMargins) MulFloatAssign(f float64)    { C.qmargins_mul_float_assign(m.p, C.double(f)) }
func (m *QMargins) DivFloatAssign(d float64)    { C.qmargins_div_float_assign(m.p, C.double(d)) }

// ---- comparison ----------------------------------------------------------

func (m *QMargins) Equal(other *QMargins) bool    { return bool(C.qmargins_equal(m.p, other.p)) }
func (m *QMargins) NotEqual(other *QMargins) bool { return !m.Equal(other) }

// ---- conversions ---------------------------------------------------------

// ToMarginsF returns a new QMarginsF handle (exact int→double per component).
func (m *QMargins) ToMarginsF() *QMarginsF {
	return wrapQMarginsF(C.qmargins_to_marginsf(m.p))
}
