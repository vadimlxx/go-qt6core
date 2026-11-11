package qgeometry

/*
 QSizeF — биндинг Qt 6.11 QtCore/QSizeF как HANDLE.

 QSizeF — value-тип в Qt (Q_RELOCATABLE_TYPE, 2×qreal). В Go хранится handle на heap-размещённый
 C++ QSizeF.

 ВАЖНОЕ ОГРАНИЧЕНИЕ:
   Qt хранит QSizeF внутри других классов ПО ЗНАЧЕНИЮ. Handle на QSizeF не видит изменений,
   сделанных C++ во встроенном QSizeF другого класса.

 ВАЖНО: Qt-конструктор QSizeF() возвращает QSizeF(-1., -1.), НЕ (0., 0.).
   - NewQSizeFDefault() → (-1., -1.) — соответствует QSizeF()
   - NewQSizeFZero()    → ( 0.,  0.) — «нулевой» QSizeF
   - NewQSizeF(w, h)    → ( w,   h)  — явный

 Семантика значений vs ссылок:
   - QSizeF — ссылочный тип (handle). Для независимой копии — Clone().
   - Go-нулевое значение QSizeF{} невалидно.

 Управление памятью:
   - NewQSizeF/QSizeFFromSize/арифметика/Clone/ToSize/Transposed/Scaled/… — новый heap-объект.
   - Close() освобождает немедленно; иначе срабатывает finalizer.

 Соглашения:
   - Геттеры: Width(), Height(), IsNull() (fuzzy), IsEmpty(), IsValid().
   - Мутирующие (in place, void): SetWidth, SetHeight, Transpose, Scale, ScaleSize.
   - Трансформации (новый handle): Transposed, Scaled, ScaledSize, ExpandedTo, BoundedTo,
     GrownBy(*QMarginsF), ShrunkBy(*QMarginsF).
   - Арифметика (новый handle): Add, Sub, MulFloat, DivFloat.
   - Присваивающие формы (in place): AddAssign, SubAssign, MulFloatAssign, DivFloatAssign.
   - Сравнение: Equal(other *QSizeF) — fuzzy; NotEqual; EqualQSize(other *QSize) — кросс.
   - Конвертация: ToSize() → *QSize (qSaturateRound); QSizeFFromSize(s *QSize) — конструктор.
   - Клонирование: Clone() → *QSizeF.

 Особенности:
   - QSizeF — fuzzy-compared. Два QSizeF с крошечной разницей могут быть "равны".
   - DivFloat(0): в release — inf/nan; в debug — Q_ASSERT. Go — IEEE 754 inf/nan без паники.
   - ToSize() использует qSaturateRound.
   - Protected-секций у QSizeF нет.

 Что не реализовано (см. qsize_doc.go):

   - QDataStream operator<< / operator>> — пакет qdatastream пока не готов.
   - QDebug operator<< — пакет qdebug пока не готов.
   - qHash(QSizeF) — не нужно в Go.
   - rwidth()/rheight() — покрыто SetWidth/SetHeight.
   - std::tuple protocol — не нужно в Go.
   - fromCGSize / toCGSize (Darwin) — платформо-специфично.
   - qFuzzyCompare / qFuzzyIsNull как свободные функции — покрыты Equal/IsNull.
*/

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC -Wno-sfinae-incomplete
#cgo LDFLAGS: -lstdc++
#include "qsize_c.h"
*/
import "C"

import "runtime"

// QSizeF wraps a heap-allocated Qt6 QSizeF.
type QSizeF struct {
	p C.go_qsizef_t
}

func wrapQSizeF(p C.go_qsizef_t) *QSizeF {
	if p == nil {
		return nil
	}
	s := &QSizeF{p: p}
	runtime.SetFinalizer(s, finalizeQSizeF)
	return s
}

func finalizeQSizeF(s *QSizeF) {
	if s.p != nil {
		C.qsizef_delete(s.p)
		s.p = nil
	}
}

// Close releases the underlying C++ object immediately and disables the finalizer.
func (s *QSizeF) Close() {
	if s.p != nil {
		C.qsizef_delete(s.p)
		s.p = nil
		runtime.SetFinalizer(s, nil)
	}
}

// ---- construction --------------------------------------------------------

// NewQSizeF returns a handle to QSizeF(w, h).
func NewQSizeF(w, h float64) *QSizeF {
	return wrapQSizeF(C.qsizef_new(C.double(w), C.double(h)))
}

// NewQSizeFZero returns a handle to QSizeF(0., 0.).
func NewQSizeFZero() *QSizeF { return NewQSizeF(0, 0) }

// NewQSizeFDefault returns a handle to QSizeF(-1., -1.), matching Qt's default constructor QSizeF().
func NewQSizeFDefault() *QSizeF { return wrapQSizeF(C.qsizef_new_default()) }

// QSizeFFromSize builds QSizeF from QSize (exact int→double). Mirrors QSizeF(const QSize&).
func QSizeFFromSize(s *QSize) *QSizeF { return wrapQSizeF(C.qsizef_from_size(s.p)) }

// Clone returns a heap copy of s (independent C++ object).
func (s *QSizeF) Clone() *QSizeF { return wrapQSizeF(C.qsizef_clone(s.p)) }

// ---- accessors -----------------------------------------------------------

func (s *QSizeF) IsNull() bool    { return bool(C.qsizef_is_null(s.p)) }
func (s *QSizeF) IsEmpty() bool   { return bool(C.qsizef_is_empty(s.p)) }
func (s *QSizeF) IsValid() bool   { return bool(C.qsizef_is_valid(s.p)) }
func (s *QSizeF) Width() float64  { return float64(C.qsizef_width(s.p)) }
func (s *QSizeF) Height() float64 { return float64(C.qsizef_height(s.p)) }

// ---- mutating setters (in place) ----------------------------------------

func (s *QSizeF) SetWidth(w float64)  { C.qsizef_set_width(s.p, C.double(w)) }
func (s *QSizeF) SetHeight(h float64) { C.qsizef_set_height(s.p, C.double(h)) }

// Transpose swaps width and height in place (mirrors QSizeF::transpose()).
func (s *QSizeF) Transpose() { C.qsizef_transpose(s.p) }

// Scale scales s in place to fit (w, h) with the given aspect ratio mode.
func (s *QSizeF) Scale(w, h float64, mode AspectRatioMode) {
	C.qsizef_scale(s.p, C.double(w), C.double(h), C.int(mode))
}

// ScaleSize scales s in place to fit other with the given aspect ratio mode.
func (s *QSizeF) ScaleSize(other *QSizeF, mode AspectRatioMode) {
	C.qsizef_scale_size(s.p, other.p, C.int(mode))
}

// ---- transforms (return new handles) ------------------------------------

func (s *QSizeF) Transposed() *QSizeF { return wrapQSizeF(C.qsizef_transposed(s.p)) }

func (s *QSizeF) Scaled(w, h float64, mode AspectRatioMode) *QSizeF {
	return wrapQSizeF(C.qsizef_scaled(s.p, C.double(w), C.double(h), C.int(mode)))
}

func (s *QSizeF) ScaledSize(other *QSizeF, mode AspectRatioMode) *QSizeF {
	return wrapQSizeF(C.qsizef_scaled_size(s.p, other.p, C.int(mode)))
}

func (s *QSizeF) ExpandedTo(other *QSizeF) *QSizeF {
	return wrapQSizeF(C.qsizef_expanded_to(s.p, other.p))
}

func (s *QSizeF) BoundedTo(other *QSizeF) *QSizeF {
	return wrapQSizeF(C.qsizef_bounded_to(s.p, other.p))
}

// GrownBy grows each dimension by the corresponding QMarginsF components.
func (s *QSizeF) GrownBy(m *QMarginsF) *QSizeF {
	return wrapQSizeF(C.qsizef_grown_by(s.p, m.p))
}

// ShrunkBy shrinks each dimension by the corresponding QMarginsF components.
func (s *QSizeF) ShrunkBy(m *QMarginsF) *QSizeF {
	return wrapQSizeF(C.qsizef_shrunk_by(s.p, m.p))
}

// ---- arithmetic (return new handles) ------------------------------------

func (s *QSizeF) Add(other *QSizeF) *QSizeF  { return wrapQSizeF(C.qsizef_add(s.p, other.p)) }
func (s *QSizeF) Sub(other *QSizeF) *QSizeF  { return wrapQSizeF(C.qsizef_sub(s.p, other.p)) }
func (s *QSizeF) MulFloat(c float64) *QSizeF { return wrapQSizeF(C.qsizef_mul_float(s.p, C.double(c))) }
func (s *QSizeF) DivFloat(c float64) *QSizeF { return wrapQSizeF(C.qsizef_div_float(s.p, C.double(c))) }

// ---- in-place assign (mutate receiver) ----------------------------------

func (s *QSizeF) AddAssign(other *QSizeF)   { C.qsizef_add_assign(s.p, other.p) }
func (s *QSizeF) SubAssign(other *QSizeF)   { C.qsizef_sub_assign(s.p, other.p) }
func (s *QSizeF) MulFloatAssign(c float64)  { C.qsizef_mul_float_assign(s.p, C.double(c)) }
func (s *QSizeF) DivFloatAssign(c float64)  { C.qsizef_div_float_assign(s.p, C.double(c)) }

// ---- comparison ----------------------------------------------------------

// Equal returns true if both components are fuzzy-equal (qFuzzyCompare per component).
func (s *QSizeF) Equal(other *QSizeF) bool { return bool(C.qsizef_equal(s.p, other.p)) }

// NotEqual returns !Equal(other).
func (s *QSizeF) NotEqual(other *QSizeF) bool { return !s.Equal(other) }

// EqualQSize returns true if s is fuzzy-equal to other (as QSizeF).
// Mirrors comparesEqual(QSizeF, QSize).
func (s *QSizeF) EqualQSize(other *QSize) bool {
	return bool(C.qsizef_equal_size(s.p, other.p))
}

// ---- conversions ---------------------------------------------------------

// ToSize converts s to a new QSize handle via qSaturateRound on each component.
func (s *QSizeF) ToSize() *QSize { return wrapQSize(C.qsizef_to_size(s.p)) }
