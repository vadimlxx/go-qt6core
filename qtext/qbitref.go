package qtext

// Package qbitarray: QBitRef — прокси к одиночному биту внутри QBitArray.

/*
 не реализовано:
   - operator bool() — используйте метод Bool().
   - operator!() — используйте метод Not().
   - operator=(bool) — используйте метод Set(val bool).
   - operator=(const QBitRef&) — используйте метод Assign(other *QBitRef).
   - QBitRef как самостоятельный C++ объект: в Go QBitRef — это {array, index},
     его методы делегируют методам родительского QBitArray. Семантика эквивалентна
     Qt: любая мутация через QBitRef сразу видна в QBitArray.

 ВАЖНО: QBitRef держит *QBitArray живым (через ссылку в поле array), но НЕ защищает
 от изменения размера родителя. Не вызывайте Resize/Clear/Fill/FillSized/Close на
 родительском массиве, пока QBitRef используется. Доступ к QBitRef после изменения
 размера — неопределённое поведение (в Qt — UB; в Go — panic при выходе за границы).
*/

// QBitRef is a proxy for a single bit inside a QBitArray.
//
// Mirrors Qt 6.11.2 QBitRef. In Qt, QBitRef is a lightweight value type returned by the
// mutable QBitArray::operator[]. In Go, it is a {array, index} pair whose methods delegate
// to the parent QBitArray; no separate C++ QBitRef object exists.
//
// Obtain one via (*QBitArray).BitRef(i). The QBitRef keeps the parent alive but does not
// protect against structural modifications of the parent — see the package-level warning
// in doc.go.
type QBitRef struct {
	array *QBitArray
	index int
}

// BitRef returns a proxy for bit i of a. Panics if i is out of range.
// Mirrors non-const QBitArray::operator[](qsizetype i) returning QBitRef.
func (a *QBitArray) BitRef(i int) *QBitRef {
	a.checkIndex(i)
	return &QBitRef{array: a, index: i}
}

// Bool returns the value of the referenced bit. Mirrors operator bool().
func (r *QBitRef) Bool() bool {
	return r.array.TestBit(r.index)
}

// Not returns the negation of the referenced bit. Mirrors operator!().
func (r *QBitRef) Not() bool {
	return !r.array.TestBit(r.index)
}

// Set writes val into the referenced bit. Mirrors operator=(bool).
func (r *QBitRef) Set(val bool) {
	r.array.SetBitVal(r.index, val)
}

// Assign copies the value of other into the referenced bit.
// Mirrors operator=(const QBitRef&).
func (r *QBitRef) Assign(other *QBitRef) {
	if other == nil {
		return
	}
	r.Set(other.Bool())
}
