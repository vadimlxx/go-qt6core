package qtext

// Package qbitarray: QBitArray — массив битов с COW-семантикой (handle-модель в Go).

/*
 не реализовано:
   - operator& / operator| / operator^ / operator~ в C++-форме — в Go используйте
     пакетные функции And(a, b) / Or(a, b) / Xor(a, b) / Not(a).
   - operator&= / operator|= / operator^= — используйте AndAssign / OrAssign / XorAssign.
   - QBitArray::data_ptr() (DataPtr = QByteArray::DataPointer) — внутренняя деталь
     реализации, не экспортируется. Для raw-доступа используйте Bits() / FromBits().
   - QDataStream operators (operator<< / operator>>) — требует qdatastream (TODO).
   - qHash(QBitArray) — требует qhash (TODO); в Go используйте хеширование среза из
     Bits() (crypto/sha256, hash/fnv и т.п.).
   - QBitArray(const QBitArray&) / QBitArray(QBitArray&&) — в Go используйте Clone()
     для явной копии; присваивание указателя Go не копирует содержимое.
   - QBitArray QDebug operator<< — неприменимо в Go.
   - QBitArray::count() (deprecated alias size()) — используйте Size().
*/

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC -Wno-sfinae-incomplete
#cgo LDFLAGS: -lstdc++
#include <stdlib.h>
#include <stdint.h>
#include <stdbool.h>
#include <stddef.h>
#include "qbitarray_c.h"
*/
import "C"

import (
	"runtime"
	"unsafe"
)

// QBitArray mirrors Qt 6.11.2 QBitArray.
//
// It is a COW (copy-on-write) array of bits. Each *QBitArray owns a heap-allocated C++
// QBitArray; the finalizer and Close() free it. Clone() is cheap (refcount++).
type QBitArray struct {
	p *C.go_qbitarray_t
}

func wrapQBitArray(p *C.go_qbitarray_t) *QBitArray {
	if p == nil {
		return nil
	}
	a := &QBitArray{p: p}
	runtime.SetFinalizer(a, finalizeQBitArray)
	return a
}

func finalizeQBitArray(a *QBitArray) {
	if a.p != nil {
		C.qbitarray_delete(a.p)
		a.p = nil
	}
}

// ---- construction --------------------------------------------------------

// New returns an empty (null and empty) QBitArray. Mirrors QBitArray().
func NewQBitArray() *QBitArray {
	return wrapQBitArray(C.qbitarray_new_empty())
}

// NewSized returns a QBitArray of the given size, with all bits initialized to val.
// Mirrors QBitArray(qsizetype size, bool val = false).
func NewQBitArraySized(size int, val bool) *QBitArray {
	return wrapQBitArray(C.qbitarray_new_sized(C.ptrdiff_t(size), C.bool(val)))
}

// ---- lifecycle -----------------------------------------------------------

// Close frees the underlying C++ object immediately and disables the finalizer.
// After Close, all methods on a are UB. Invalidates any *QBitRef into a.
func (a *QBitArray) Close() {
	if a.p != nil {
		C.qbitarray_delete(a.p)
		a.p = nil
		runtime.SetFinalizer(a, nil)
	}
}

// Clone returns a shallow COW copy (refcount++; cheap). Mirrors QBitArray(const QBitArray&).
func (a *QBitArray) Clone() *QBitArray {
	r := wrapQBitArray(C.qbitarray_new_copy(a.p))
	runtime.KeepAlive(a)
	return r
}

// Assign copies the contents of other into a (COW; cheap). Mirrors operator=(const QBitArray&).
// Invalidates any *QBitRef into a.
func (a *QBitArray) Assign(other *QBitArray) {
	if other == nil || other.p == nil {
		return
	}
	C.qbitarray_assign(a.p, other.p)
	runtime.KeepAlive(a)
	runtime.KeepAlive(other)
}

// Swap exchanges the contents of a and other (COW swap). Mirrors swap(QBitArray&).
// Invalidates any *QBitRef into a or other.
func (a *QBitArray) Swap(other *QBitArray) {
	if other == nil || other.p == nil {
		return
	}
	C.qbitarray_swap(a.p, other.p)
	runtime.KeepAlive(a)
	runtime.KeepAlive(other)
}

// ---- size / status -------------------------------------------------------

// Size returns the number of bits. Mirrors size().
func (a *QBitArray) Size() int {
	return int(C.qbitarray_size(a.p))
}

// Count returns the number of bits. Deprecated alias of Size() (as in Qt).
func (a *QBitArray) Count() int { return a.Size() }

// CountSet returns the number of bits set to on. Mirrors count(bool on).
func (a *QBitArray) CountSet(on bool) int {
	return int(C.qbitarray_count_set(a.p, C.bool(on)))
}

// IsEmpty reports whether Size() == 0. Mirrors isEmpty().
func (a *QBitArray) IsEmpty() bool {
	return bool(C.qbitarray_is_empty(a.p))
}

// IsNull reports whether the array is null (no underlying storage allocated).
// Mirrors isNull(). Note: null is a subset of empty.
func (a *QBitArray) IsNull() bool {
	return bool(C.qbitarray_is_null(a.p))
}

// ---- modifiers -----------------------------------------------------------

// Resize changes the number of bits to size. New bits are false; removed bits are dropped.
// Mirrors resize(qsizetype size). Invalidates any *QBitRef into a.
func (a *QBitArray) Resize(size int) {
	C.qbitarray_resize(a.p, C.ptrdiff_t(size))
	runtime.KeepAlive(a)
}

// Detach ensures that the array has its own private storage (COW detach).
// Mirrors detach().
func (a *QBitArray) Detach() {
	C.qbitarray_detach(a.p)
	runtime.KeepAlive(a)
}

// IsDetached reports whether the array has its own private storage.
// Mirrors isDetached().
func (a *QBitArray) IsDetached() bool {
	return bool(C.qbitarray_is_detached(a.p))
}

// Clear resets the array to null and empty. Mirrors clear().
// Invalidates any *QBitRef into a.
func (a *QBitArray) Clear() {
	C.qbitarray_clear(a.p)
	runtime.KeepAlive(a)
}

// ---- bit operations ------------------------------------------------------

// checkIndex panics on out-of-range access.
func (a *QBitArray) checkIndex(i int) {
	if i < 0 || i >= a.Size() {
		panic("qbitarray: index out of range")
	}
}

// TestBit returns the value of bit i. Panics if i is out of range.
// Mirrors testBit(qsizetype i).
func (a *QBitArray) TestBit(i int) bool {
	a.checkIndex(i)
	return bool(C.qbitarray_test_bit(a.p, C.ptrdiff_t(i)))
}

// At is TestBit (read-only operator[]). Mirrors at(qsizetype i).
func (a *QBitArray) At(i int) bool { return a.TestBit(i) }

// SetBit sets bit i to true. Panics if i is out of range. Mirrors setBit(qsizetype i).
func (a *QBitArray) SetBit(i int) {
	a.checkIndex(i)
	C.qbitarray_set_bit(a.p, C.ptrdiff_t(i))
	runtime.KeepAlive(a)
}

// SetBitVal sets bit i to val. Panics if i is out of range.
// Mirrors setBit(qsizetype i, bool val).
func (a *QBitArray) SetBitVal(i int, val bool) {
	a.checkIndex(i)
	C.qbitarray_set_bit_val(a.p, C.ptrdiff_t(i), C.bool(val))
	runtime.KeepAlive(a)
}

// ClearBit sets bit i to false. Panics if i is out of range. Mirrors clearBit(qsizetype i).
func (a *QBitArray) ClearBit(i int) {
	a.checkIndex(i)
	C.qbitarray_clear_bit(a.p, C.ptrdiff_t(i))
	runtime.KeepAlive(a)
}

// ToggleBit flips bit i and returns its previous value. Panics if i is out of range.
// Mirrors toggleBit(qsizetype i).
func (a *QBitArray) ToggleBit(i int) bool {
	a.checkIndex(i)
	return bool(C.qbitarray_toggle_bit(a.p, C.ptrdiff_t(i)))
}

// ---- bitwise assign ------------------------------------------------------

// AndAssign performs a &= other. Sizes are extended to the longer of the two, with the
// shorter array treated as zero-padded on the right. Invalidates *QBitRef into a.
// Mirrors operator&=(const QBitArray&).
func (a *QBitArray) AndAssign(other *QBitArray) {
	if other == nil || other.p == nil {
		return
	}
	C.qbitarray_and_assign(a.p, other.p)
	runtime.KeepAlive(a)
	runtime.KeepAlive(other)
}

// OrAssign performs a |= other. Mirrors operator|=(const QBitArray&).
func (a *QBitArray) OrAssign(other *QBitArray) {
	if other == nil || other.p == nil {
		return
	}
	C.qbitarray_or_assign(a.p, other.p)
	runtime.KeepAlive(a)
	runtime.KeepAlive(other)
}

// XorAssign performs a ^= other. Mirrors operator^=(const QBitArray&).
func (a *QBitArray) XorAssign(other *QBitArray) {
	if other == nil || other.p == nil {
		return
	}
	C.qbitarray_xor_assign(a.p, other.p)
	runtime.KeepAlive(a)
	runtime.KeepAlive(other)
}

// ---- bitwise non-assign (package functions) -----------------------------

// And returns a fresh QBitArray equal to a & b. Mirrors operator&(const QBitArray&, const QBitArray&).
func And(a, b *QBitArray) *QBitArray {
	r := wrapQBitArray(C.qbitarray_op_and(a.p, b.p))
	runtime.KeepAlive(a)
	runtime.KeepAlive(b)
	return r
}

// Or returns a fresh QBitArray equal to a | b. Mirrors operator|(const QBitArray&, const QBitArray&).
func Or(a, b *QBitArray) *QBitArray {
	r := wrapQBitArray(C.qbitarray_op_or(a.p, b.p))
	runtime.KeepAlive(a)
	runtime.KeepAlive(b)
	return r
}

// Xor returns a fresh QBitArray equal to a ^ b. Mirrors operator^(const QBitArray&, const QBitArray&).
func Xor(a, b *QBitArray) *QBitArray {
	r := wrapQBitArray(C.qbitarray_op_xor(a.p, b.p))
	runtime.KeepAlive(a)
	runtime.KeepAlive(b)
	return r
}

// Not returns a fresh QBitArray equal to ~a (all bits inverted).
// Mirrors operator~(QBitArray a).
func Not(a *QBitArray) *QBitArray {
	r := wrapQBitArray(C.qbitarray_inverted(a.p))
	runtime.KeepAlive(a)
	return r
}

// ---- fill / truncate -----------------------------------------------------

// Fill sets every bit of a to val (size unchanged). Mirrors fill(bool val) (with asize=-1).
func (a *QBitArray) Fill(val bool) {
	C.qbitarray_fill_all(a.p, C.bool(val))
	runtime.KeepAlive(a)
}

// FillSized resizes a to size and then fills every bit with val.
// Mirrors fill(bool val, qsizetype size). Invalidates *QBitRef into a.
func (a *QBitArray) FillSized(val bool, size int) {
	C.qbitarray_fill_sized(a.p, C.bool(val), C.ptrdiff_t(size))
	runtime.KeepAlive(a)
}

// FillRange fills bits [first, last) with val. Panics if range is invalid.
// Mirrors fill(bool val, qsizetype first, qsizetype last).
func (a *QBitArray) FillRange(val bool, first, last int) {
	if first < 0 || last < first || last > a.Size() {
		panic("qbitarray: FillRange out of range")
	}
	C.qbitarray_fill_range(a.p, C.bool(val), C.ptrdiff_t(first), C.ptrdiff_t(last))
	runtime.KeepAlive(a)
}

// Truncate truncates the array at position pos (no-op if pos >= Size()).
// Mirrors truncate(qsizetype pos). Invalidates *QBitRef into a.
func (a *QBitArray) Truncate(pos int) {
	C.qbitarray_truncate(a.p, C.ptrdiff_t(pos))
	runtime.KeepAlive(a)
}

// ---- raw bits (LSB-first within each byte) -----------------------------

// Bits returns a copy of the raw bit buffer.
// Format: bit i = (out[i/8] >> (i%8)) & 1, len(out) == (Size()+7)/8. For an empty array,
// returns nil. Mirrors bits().
func (a *QBitArray) Bits() []byte {
	var n C.ptrdiff_t
	p := C.qbitarray_bits(a.p, &n)
	runtime.KeepAlive(a)
	if p == nil || n == 0 {
		return nil
	}
	return C.GoBytes(unsafe.Pointer(p), C.int(n))
}

// FromBits constructs a new QBitArray from bitLen valid bits of data.
// data must contain at least (bitLen+7)/8 bytes. Format is the same as Bits().
// Mirrors static QBitArray::fromBits(const char*, qsizetype len).
func FromBits(data []byte, bitLen int) *QBitArray {
	if bitLen < 0 {
		panic("qbitarray: FromBits negative bitLen")
	}
	need := (bitLen + 7) / 8
	if len(data) < need {
		panic("qbitarray: FromBits data too short")
	}
	var p *C.char
	if need > 0 {
		p = (*C.char)(unsafe.Pointer(&data[0]))
	}
	r := wrapQBitArray(C.qbitarray_from_bits(p, C.ptrdiff_t(bitLen)))
	runtime.KeepAlive(data)
	return r
}

// ToUInt32 interprets the first 32 bits of a as uint32 using the specified endianness.
// Requires Size() >= 32; ok=false otherwise (result unspecified, matching Qt).
// Mirrors toUInt32(QSysInfo::Endian endianness, bool *ok).
func (a *QBitArray) ToUInt32(littleEndian bool) (uint32, bool) {
	var ok C.bool
	v := C.qbitarray_to_uint32(a.p, C.bool(littleEndian), &ok)
	runtime.KeepAlive(a)
	return uint32(v), bool(ok)
}

// ---- comparison ----------------------------------------------------------

// Equal reports a == other (bit-by-bit, size and content).
// Mirrors operator==(const QBitArray&).
func (a *QBitArray) Equal(other *QBitArray) bool {
	if other == nil || other.p == nil {
		return false
	}
	r := C.qbitarray_equal(a.p, other.p)
	runtime.KeepAlive(a)
	runtime.KeepAlive(other)
	return bool(r)
}
