package qgeometry

/*
 QSize — биндинг Qt 6.11 QtCore/QSize как HANDLE.

 QSize — value-тип в Qt (Q_RELOCATABLE_TYPE, 2×int). В Go хранится handle на heap-размещённый
 C++ QSize. Все геттеры/сеттеры идут напрямую в C++.

 ВАЖНОЕ ОГРАНИЧЕНИЕ:
   Qt хранит QSize внутри других классов ПО ЗНАЧЕНИЮ. Handle на QSize не видит изменений,
   сделанных C++ во встроенном QSize другого класса. Это ограничение API Qt, а не биндинга.

 ВАЖНО: Qt-конструктор QSize() возвращает QSize(-1, -1), НЕ (0, 0).
   - NewQSizeDefault()  → (-1, -1)  — соответствует QSize()
   - NewQSizeZero()     → ( 0,  0)  — «нулевой» QSize
   - NewQSize(w, h)     → ( w,  h)  — явный

 Семантика значений vs ссылок:
   - QSize — ссылочный тип (handle). Мутации через один handle видны через все handle'ы на
     тот же объект. Для независимой копии — Clone().
   - Go-нулевое значение QSize{} невалидно. Используйте NewQSize / NewQSizeZero / NewQSizeDefault.

 Управление памятью:
   - NewQSize/NewQSizeZero/NewQSizeDefault/арифметика/Clone/ToSizeF/Transposed/Scaled/… —
     создают НОВЫЙ heap-объект. Close() освобождает немедленно; иначе срабатывает finalizer.

 Соглашения:
   - Геттеры: Width(), Height(), IsNull(), IsEmpty(), IsValid().
   - Мутирующие (in place, void): SetWidth, SetHeight, Transpose, Scale, ScaleSize.
   - Трансформации (новый handle): Transposed, Scaled, ScaledSize, ExpandedTo, BoundedTo,
     GrownBy(*QMargins), ShrunkBy(*QMargins).
   - Арифметика (новый handle): Add, Sub, MulFloat, DivFloat.
   - Присваивающие формы (in place): AddAssign, SubAssign, MulFloatAssign, DivFloatAssign.
   - Сравнение: Equal(other *QSize), NotEqual(other).
   - Конвертация: ToSizeF() → *QSizeF.
   - Клонирование: Clone() → *QSize.

 AspectRatioMode (Qt::AspectRatioMode) определён в этом файле и используется QSize/QSizeF:

   - IgnoreAspectRatio           — размер = запрошенному, аспект не сохраняется
   - KeepAspectRatio             — вписать в запрошенный, сохраняя аспект
   - KeepAspectRatioByExpanding  — покрыть запрошенный, сохраняя аспект

 Особенности:
   - MulFloat/DivFloat используют Qt's qSaturateRound — идентично operator*(QSize, qreal).
   - DivFloat(0): в release — saturate(inf); в debug — Q_ASSERT. Go — как release.
   - Protected-секций у QSize нет.

 Что не реализовано (см. qsize_doc.go):

   - QDataStream operator<< / operator>> — пакет qdatastream пока не готов.
   - QDebug operator<< — пакет qdebug пока не готов.
   - qHash(QSize) — не нужно в Go (пользователь может хэшировать Width()/Height()).
   - rwidth()/rheight() — покрыто SetWidth/SetHeight.
   - std::tuple protocol — не нужно в Go.
   - toCGSize (Darwin) — платформо-специфично.
*/

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC -Wno-sfinae-incomplete
#cgo LDFLAGS: -lstdc++
#include "qsize_c.h"
*/
import "C"

import "runtime"

// AspectRatioMode mirrors Qt::AspectRatioMode.
type AspectRatioMode int

const (
	// IgnoreAspectRatio — размер устанавливается равным указанному, аспект не сохраняется.
	IgnoreAspectRatio AspectRatioMode = 0
	// KeepAspectRatio — размер масштабируется, чтобы вписаться в указанный, сохраняя аспект.
	KeepAspectRatio AspectRatioMode = 1
	// KeepAspectRatioByExpanding — размер масштабируется, чтобы покрыть указанный, сохраняя аспект.
	KeepAspectRatioByExpanding AspectRatioMode = 2
)

// QSize wraps a heap-allocated Qt6 QSize.
type QSize struct {
	p C.go_qsize_t
}

func wrapQSize(p C.go_qsize_t) *QSize {
	if p == nil {
		return nil
	}
	s := &QSize{p: p}
	runtime.SetFinalizer(s, finalizeQSize)
	return s
}

func finalizeQSize(s *QSize) {
	if s.p != nil {
		C.qsize_delete(s.p)
		s.p = nil
	}
}

// Close releases the underlying C++ object immediately and disables the finalizer.
func (s *QSize) Close() {
	if s.p != nil {
		C.qsize_delete(s.p)
		s.p = nil
		runtime.SetFinalizer(s, nil)
	}
}

// ---- construction --------------------------------------------------------

// NewQSize returns a handle to QSize(w, h).
func NewQSize(w, h int) *QSize {
	return wrapQSize(C.qsize_new(C.int(w), C.int(h)))
}

// NewQSizeZero returns a handle to QSize(0, 0).
func NewQSizeZero() *QSize { return NewQSize(0, 0) }

// NewQSizeDefault returns a handle to QSize(-1, -1), matching Qt's default constructor QSize().
func NewQSizeDefault() *QSize { return wrapQSize(C.qsize_new_default()) }

// Clone returns a heap copy of s (independent C++ object).
func (s *QSize) Clone() *QSize { return wrapQSize(C.qsize_clone(s.p)) }

// ---- accessors -----------------------------------------------------------

func (s *QSize) IsNull() bool  { return bool(C.qsize_is_null(s.p)) }
func (s *QSize) IsEmpty() bool { return bool(C.qsize_is_empty(s.p)) }
func (s *QSize) IsValid() bool { return bool(C.qsize_is_valid(s.p)) }
func (s *QSize) Width() int    { return int(C.qsize_width(s.p)) }
func (s *QSize) Height() int   { return int(C.qsize_height(s.p)) }

// ---- mutating setters (in place) ----------------------------------------

func (s *QSize) SetWidth(w int)  { C.qsize_set_width(s.p, C.int(w)) }
func (s *QSize) SetHeight(h int) { C.qsize_set_height(s.p, C.int(h)) }

// Transpose swaps width and height in place (mirrors QSize::transpose()).
func (s *QSize) Transpose() { C.qsize_transpose(s.p) }

// Scale scales s in place to fit (w, h) with the given aspect ratio mode.
func (s *QSize) Scale(w, h int, mode AspectRatioMode) {
	C.qsize_scale(s.p, C.int(w), C.int(h), C.int(mode))
}

// ScaleSize scales s in place to fit other with the given aspect ratio mode.
func (s *QSize) ScaleSize(other *QSize, mode AspectRatioMode) {
	C.qsize_scale_size(s.p, other.p, C.int(mode))
}

// ---- transforms (return new handles) ------------------------------------

func (s *QSize) Transposed() *QSize { return wrapQSize(C.qsize_transposed(s.p)) }

func (s *QSize) Scaled(w, h int, mode AspectRatioMode) *QSize {
	return wrapQSize(C.qsize_scaled(s.p, C.int(w), C.int(h), C.int(mode)))
}

func (s *QSize) ScaledSize(other *QSize, mode AspectRatioMode) *QSize {
	return wrapQSize(C.qsize_scaled_size(s.p, other.p, C.int(mode)))
}

func (s *QSize) ExpandedTo(other *QSize) *QSize {
	return wrapQSize(C.qsize_expanded_to(s.p, other.p))
}

func (s *QSize) BoundedTo(other *QSize) *QSize {
	return wrapQSize(C.qsize_bounded_to(s.p, other.p))
}

// GrownBy grows each dimension by the corresponding QMargins components (mirrors QSize::grownBy).
func (s *QSize) GrownBy(m *QMargins) *QSize {
	return wrapQSize(C.qsize_grown_by(s.p, m.p))
}

// ShrunkBy shrinks each dimension by the corresponding QMargins components (mirrors QSize::shrunkBy).
func (s *QSize) ShrunkBy(m *QMargins) *QSize {
	return wrapQSize(C.qsize_shrunk_by(s.p, m.p))
}

// ---- arithmetic (return new handles) ------------------------------------

func (s *QSize) Add(other *QSize) *QSize      { return wrapQSize(C.qsize_add(s.p, other.p)) }
func (s *QSize) Sub(other *QSize) *QSize      { return wrapQSize(C.qsize_sub(s.p, other.p)) }
func (s *QSize) MulFloat(c float64) *QSize    { return wrapQSize(C.qsize_mul_float(s.p, C.double(c))) }
func (s *QSize) DivFloat(c float64) *QSize    { return wrapQSize(C.qsize_div_float(s.p, C.double(c))) }

// ---- in-place assign (mutate receiver) ----------------------------------

func (s *QSize) AddAssign(other *QSize)    { C.qsize_add_assign(s.p, other.p) }
func (s *QSize) SubAssign(other *QSize)    { C.qsize_sub_assign(s.p, other.p) }
func (s *QSize) MulFloatAssign(c float64)  { C.qsize_mul_float_assign(s.p, C.double(c)) }
func (s *QSize) DivFloatAssign(c float64)  { C.qsize_div_float_assign(s.p, C.double(c)) }

// ---- comparison ----------------------------------------------------------

// Equal returns true if width and height are identical (exact int comparison).
func (s *QSize) Equal(other *QSize) bool { return bool(C.qsize_equal(s.p, other.p)) }

// NotEqual returns !Equal(other).
func (s *QSize) NotEqual(other *QSize) bool { return !s.Equal(other) }

// ---- conversions ---------------------------------------------------------

// ToSizeF returns a new QSizeF handle (exact int→double per component).
func (s *QSize) ToSizeF() *QSizeF { return wrapQSizeF(C.qsize_to_sizef(s.p)) }
