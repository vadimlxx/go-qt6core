package qgeometry

/*
 QSize — биндинг Qt 6.11 QtCore/QSize.

 QSize — value-тип (Q_RELOCATABLE_TYPE): 2 int-поля (width, height). В Go поля неэкспортируемые,
 доступ через геттеры Width()/Height(). Через C ABI передаётся плоской структурой go_qsize_t
 (2×int, 8 байт).

 ВАЖНО: Qt-конструктор QSize() возвращает QSize(-1, -1), НЕ (0, 0).
   - NewQSizeDefault()  → (-1, -1)  — соответствует QSize()
   - NewQSizeZero()     → ( 0,  0)  — «нулевой» QSize
   - NewQSize(w, h)     → ( w,  h)  — явный

 Соглашения:

   - Геттеры: Width(), Height(), IsNull(), IsEmpty(), IsValid().
   - Мутирующие сеттеры: SetWidth(w), SetHeight(h), Transpose(), Scale(w, h, mode),
     ScaleSize(sz, mode).
   - Трансформации (возвращают новый QSize): Transposed(), Scaled(w, h, mode), ScaledSize(sz, mode),
     ExpandedTo(other), BoundedTo(other), GrownBy(QMargins), ShrunkBy(QMargins).
   - Арифметика: Add(other), Sub(other), MulFloat(c), DivFloat(c).
     В Qt операторы *= и /= принимают только qreal — поэтому MulInt/DivInt отсутствуют.
   - Присваивающие формы: AddAssign, SubAssign, MulFloatAssign, DivFloatAssign.
   - Сравнение: Equal(other), NotEqual(other) — точное (comparesEqual).
   - Конвертация: ToSizeF() → QSizeF.

 AspectRatioMode (Qt::AspectRatioMode) определён в этом файле и используется QSize/QSizeF:

   - IgnoreAspectRatio           — размер = запрошенному, аспект не сохраняется
   - KeepAspectRatio             — вписать в запрошенный, сохраняя аспект
   - KeepAspectRatioByExpanding  — покрыть запрошенный, сохраняя аспект

 Особенности:

   - MulFloat/DivFloat используют Qt's qSaturateRound — идентично operator*(QSize, qreal).
   - DivFloat(0): в release — saturate(inf) = INT_MAX/INT_MIN; в debug — Q_ASSERT. Go — как release.
   - Scale(w,h,mode) и Scaled(w,h,mode) различаются только мутацией/возвратом — оба делегируют Qt.
   - Protected-секций у QSize нет — открывать нечего.

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

// QSize wraps Qt6 QSize (Q_RELOCATABLE_TYPE, 2×int value).
// Fields are unexported; use Width()/Height() to read, SetWidth()/SetHeight() to mutate.
type QSize struct {
	width  int
	height int
}

// ---- construction --------------------------------------------------------

// NewQSize returns QSize(w, h).
func NewQSize(w, h int) QSize {
	return qsizeFromC(C.qsize_new(C.int(w), C.int(h)))
}

// NewQSizeZero returns QSize(0, 0).
func NewQSizeZero() QSize { return NewQSize(0, 0) }

// NewQSizeDefault returns QSize(-1, -1), matching Qt's default constructor QSize().
func NewQSizeDefault() QSize { return qsizeFromC(C.qsize_default()) }

// ---- accessors -----------------------------------------------------------

func (s QSize) IsNull() bool  { return bool(C.qsize_is_null(qsizeToC(s))) }
func (s QSize) IsEmpty() bool { return bool(C.qsize_is_empty(qsizeToC(s))) }
func (s QSize) IsValid() bool { return bool(C.qsize_is_valid(qsizeToC(s))) }
func (s QSize) Width() int    { return int(C.qsize_width(qsizeToC(s))) }
func (s QSize) Height() int   { return int(C.qsize_height(qsizeToC(s))) }

// ---- modifiers -----------------------------------------------------------

func (s *QSize) SetWidth(w int)  { *s = qsizeFromC(C.qsize_set_width(qsizeToC(*s), C.int(w))) }
func (s *QSize) SetHeight(h int) { *s = qsizeFromC(C.qsize_set_height(qsizeToC(*s), C.int(h))) }

// Transpose swaps width and height in place (mirrors QSize::transpose()).
func (s *QSize) Transpose() { *s = s.Transposed() }

// Scale scales s in place to fit (w, h) with the given aspect ratio mode
// (mirrors QSize::scale(int, int, Qt::AspectRatioMode)).
func (s *QSize) Scale(w, h int, mode AspectRatioMode) { *s = s.Scaled(w, h, mode) }

// ScaleSize scales s in place to fit other with the given aspect ratio mode
// (mirrors QSize::scale(const QSize&, Qt::AspectRatioMode)).
func (s *QSize) ScaleSize(other QSize, mode AspectRatioMode) { *s = s.ScaledSize(other, mode) }

// ---- transforms (return new QSize) --------------------------------------

// Transposed returns a QSize with width and height swapped.
func (s QSize) Transposed() QSize { return qsizeFromC(C.qsize_transposed(qsizeToC(s))) }

// Scaled returns a new QSize scaled to fit (w, h) with the given aspect ratio mode.
func (s QSize) Scaled(w, h int, mode AspectRatioMode) QSize {
	return qsizeFromC(C.qsize_scaled(qsizeToC(s), C.int(w), C.int(h), C.int(mode)))
}

// ScaledSize returns a new QSize scaled to fit other with the given aspect ratio mode.
func (s QSize) ScaledSize(other QSize, mode AspectRatioMode) QSize {
	return qsizeFromC(C.qsize_scaled_size(qsizeToC(s), qsizeToC(other), C.int(mode)))
}

// ExpandedTo returns the component-wise max of s and other (mirrors QSize::expandedTo).
func (s QSize) ExpandedTo(other QSize) QSize {
	return qsizeFromC(C.qsize_expanded_to(qsizeToC(s), qsizeToC(other)))
}

// BoundedTo returns the component-wise min of s and other (mirrors QSize::boundedTo).
func (s QSize) BoundedTo(other QSize) QSize {
	return qsizeFromC(C.qsize_bounded_to(qsizeToC(s), qsizeToC(other)))
}

// GrownBy grows each dimension by the corresponding QMargins components (mirrors QSize::grownBy).
func (s QSize) GrownBy(m QMargins) QSize {
	return qsizeFromC(C.qsize_grown_by(qsizeToC(s), marginsToC(m)))
}

// ShrunkBy shrinks each dimension by the corresponding QMargins components (mirrors QSize::shrunkBy).
func (s QSize) ShrunkBy(m QMargins) QSize {
	return qsizeFromC(C.qsize_shrunk_by(qsizeToC(s), marginsToC(m)))
}

// ---- arithmetic ----------------------------------------------------------

// Add returns s + other (component-wise).
func (s QSize) Add(other QSize) QSize {
	return qsizeFromC(C.qsize_add(qsizeToC(s), qsizeToC(other)))
}

// Sub returns s - other (component-wise).
func (s QSize) Sub(other QSize) QSize {
	return qsizeFromC(C.qsize_sub(qsizeToC(s), qsizeToC(other)))
}

// MulFloat returns s * c on each component (rounded via Qt's qSaturateRound).
func (s QSize) MulFloat(c float64) QSize {
	return qsizeFromC(C.qsize_mul_float(qsizeToC(s), C.double(c)))
}

// DivFloat returns s / c on each component (rounded via Qt's qSaturateRound).
func (s QSize) DivFloat(c float64) QSize {
	return qsizeFromC(C.qsize_div_float(qsizeToC(s), C.double(c)))
}

// ---- assign forms --------------------------------------------------------

func (s *QSize) AddAssign(other QSize)    { *s = s.Add(other) }
func (s *QSize) SubAssign(other QSize)    { *s = s.Sub(other) }
func (s *QSize) MulFloatAssign(c float64) { *s = s.MulFloat(c) }
func (s *QSize) DivFloatAssign(c float64) { *s = s.DivFloat(c) }

// ---- comparison ----------------------------------------------------------

// Equal returns true if width and height are identical (exact int comparison).
func (s QSize) Equal(other QSize) bool {
	return bool(C.qsize_equal(qsizeToC(s), qsizeToC(other)))
}

// NotEqual returns !Equal(other).
func (s QSize) NotEqual(other QSize) bool { return !s.Equal(other) }

// ---- conversions ---------------------------------------------------------

// ToSizeF converts s to QSizeF (exact int→double per component).
func (s QSize) ToSizeF() QSizeF { return qsizefFromC(C.qsize_to_sizef(qsizeToC(s))) }

// ---- internal helpers ----------------------------------------------------

func qsizeToC(s QSize) C.go_qsize_t {
	return C.go_qsize_t{width: C.int(s.width), height: C.int(s.height)}
}

func qsizeFromC(c C.go_qsize_t) QSize {
	return QSize{width: int(c.width), height: int(c.height)}
}