package qtext

/*
	QBitArray — handle над C++ QBitArray; владение через runtime.AddCleanup + Close().

	# Семантика

	Каждый *QBitArray — handle на C++ QBitArray (COW-массив битов).
	Копирование (Clone) дёшево (COW refcount++). Мутирующие методы detach-ят
	внутри Qt.

	# Владение

	C++ объект освобождается одним из двух способов:
	- Close() — немедленно, детерминированно;
	- runtime.AddCleanup — когда *QBitArray становится недостижим.

	Close() опционален, идемпотентен и безопасен на nil. Если Close() и cleanup
	гоняются, владение захватывает ровно один участник (атомарный Swap(nil)),
	двойного free нет.

	# Use after Close

	После Close() любой метод, кроме Close(), паникует
	("qbitarray: use after Close (or nil receiver)").

	# QBitRef

	QBitRef — Go-side прокси {array, index}, без отдельного C++ объекта.
	Он держит родительский *QBitArray живым, но не защищает от структурных
	изменений родителя. Не вызывайте Resize/Clear/Fill/FillSized/Close на
	родителе, пока QBitRef в работе. После Close родителя QBitRef паникует
	при любом обращении (через родительский handle()).

	# Соглашения по именам

	Все конструкторы QBitArray начинаются с NewQBitArray* (NewQBitArray,
	NewQBitArraySized, NewQBitArrayFromBits). Пакетные операторы And/Or/Xor/Not
	соответствуют C++-операторам &/|/^/~ и имени New не имеют — это не
	конструкторы, а операторные эквиваленты.

	# Что сознательно не вошло (TODO, не блокирует текущий API)

	- QDataStream operator<< / operator>> — требует qdatastream (TODO).
	- qHash(QBitArray) — требует qhash (TODO); в Go хешируйте срез из Bits()
		(crypto/sha256, hash/fnv и т.п.).
	- QDebug operator<< — неприменимо в Go.
	- QBitArray::count() без аргумента (deprecated alias size()) — используйте Size().
	- data_ptr() (DataPtr = QByteArray::DataPointer) — внутренняя деталь; для
		raw-доступа используйте Bits()/NewQBitArrayFromBits.
	- Перегрузки bitwise-операторов на rvalue (& &, |, ^ с moves) — не
		экспортируются; AndAssign/OrAssign/XorAssign + Clone покрывают те же
		сценарии.
*/

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC
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
	"sync/atomic"
	"unsafe"
)

// qbitarrayState holds the C++ QBitArray pointer in an atomic slot so that
// Close() and the runtime cleanup can race safely: whoever swaps the slot
// to nil owns the C++ object and is responsible for deleting it.
//
// Why *byte and not *C.go_qbitarray_t: cgo exposes the opaque struct
// go_qbitarray_s as an incomplete Go type (forward-declared only under C
// mode), and atomic.Pointer[T] requires T to be a complete type. The
// pointer value is unchanged; only its Go-side static type differs. Cast
// back to *C.go_qbitarray_t at C-call boundaries.
//
// State is a separate allocation (not embedded in QBitArray) because the
// cleanup function must not capture *QBitArray — doing so would keep the
// object reachable forever and the cleanup would never run.
type qbitarrayState struct {
	p atomic.Pointer[byte]
}

// cpp returns the underlying *C.go_qbitarray_t, or nil if the slot was
// claimed (closed) already. Read-only, does not take ownership.
func (s *qbitarrayState) cpp() *C.go_qbitarray_t {
	if bp := s.p.Load(); bp != nil {
		return (*C.go_qbitarray_t)(unsafe.Pointer(bp))
	}
	return nil
}

// claim atomically sets the slot to nil and returns the
// *C.go_qbitarray_t that was there (if any). Exactly one caller gets a
// non-nil result for a given stored pointer; subsequent claimers see nil.
// This is what makes Close() vs cleanup safe against double free.
func (s *qbitarrayState) claim() *C.go_qbitarray_t {
	if bp := s.p.Swap(nil); bp != nil {
		return (*C.go_qbitarray_t)(unsafe.Pointer(bp))
	}
	return nil
}

// QBitArray mirrors Qt 6.11.2 QBitArray.
//
// It is a COW (copy-on-write) array of bits. Each *QBitArray owns a
// heap-allocated C++ QBitArray; the runtime cleanup and Close() free it.
// Clone() is cheap (refcount++).
type QBitArray struct {
	state   *qbitarrayState
	cleanup runtime.Cleanup
}

// wrapQBitArray wraps an opaque *C.go_qbitarray_t into *QBitArray.
// Takes ownership: cleanup (or Close) frees the C++ object.
func wrapQBitArray(p *C.go_qbitarray_t) *QBitArray {
	if p == nil {
		return nil
	}
	st := &qbitarrayState{}
	st.p.Store((*byte)(unsafe.Pointer(p)))
	a := &QBitArray{state: st}
	a.cleanup = runtime.AddCleanup(a, func(st *qbitarrayState) {
		if cpp := st.claim(); cpp != nil {
			C.qbitarray_delete(cpp)
		}
	}, st)
	return a
}

// handle returns the underlying C++ pointer, panicking if a is nil or closed.
func (a *QBitArray) handle() *C.go_qbitarray_t {
	if cpp, ok := a.tryHandle(); ok {
		return cpp
	}
	panic("qbitarray: use after Close (or nil receiver)")
}

// tryHandle is the non-panicking variant of handle. Returns nil, false if
// a is nil, has no state, or has already been closed.
func (a *QBitArray) tryHandle() (*C.go_qbitarray_t, bool) {
	if a == nil || a.state == nil {
		return nil, false
	}
	if cpp := a.state.cpp(); cpp != nil {
		return cpp, true
	}
	return nil, false
}

// ---- construction --------------------------------------------------------

// NewQBitArray dispatches to the matching constructor based on the concrete
// Go argument types. Qt: no single counterpart.
//
//	NewQBitArray()                → empty QBitArray (QBitArray())
//	NewQBitArray(*QBitArray)      → Clone (QBitArray(const QBitArray&))
//	NewQBitArray(int, bool)       → NewQBitArraySized(size, val)
func NewQBitArray(arg ...any) *QBitArray {
	if len(arg) == 0 || arg[0] == nil {
		return wrapQBitArray(C.qbitarray_new_empty())
	}

	if len(arg) == 1 {
		if b, ok := arg[0].(*QBitArray); ok {
			return b.Clone()
		}
	}

	if len(arg) == 2 {
		size, ok1 := arg[0].(int)
		val, ok2 := arg[1].(bool)
		if ok1 && ok2 {
			return NewQBitArraySized(size, val)
		}
	}

	panic("NewQBitArray invalid argument")
}

// NewQBitArraySized returns a QBitArray of the given size, with all bits
// initialized to val.
// Qt: QBitArray(qsizetype size, bool val = false) — explicit constructor.
func NewQBitArraySized(size int, val bool) *QBitArray {
	return wrapQBitArray(C.qbitarray_new_sized(C.ptrdiff_t(size), C.bool(val)))
}

// NewQBitArrayFromBits constructs a QBitArray from bitLen valid bits of data.
// data must contain at least (bitLen+7)/8 bytes. Bit layout is the same as
// Bits(): bit i = (data[i/8] >> (i % 8)) & 1.
// Qt: static QBitArray::fromBits(const char *data, qsizetype len).
func NewQBitArrayFromBits(data []byte, bitLen int) *QBitArray {
	if bitLen < 0 {
		panic("qbitarray: NewQBitArrayFromBits negative bitLen")
	}
	need := (bitLen + 7) / 8
	if len(data) < need {
		panic("qbitarray: NewQBitArrayFromBits data too short")
	}
	var p *C.char
	if need > 0 {
		p = (*C.char)(unsafe.Pointer(&data[0]))
	}
	r := wrapQBitArray(C.qbitarray_from_bits(p, C.ptrdiff_t(bitLen)))
	runtime.KeepAlive(data)
	return r
}

// ---- lifecycle -----------------------------------------------------------

// Close releases the underlying C++ object immediately and cancels the
// runtime cleanup. Close is optional, idempotent, and nil-safe.
//
// If Close races with the runtime cleanup, exactly one of them claims the
// C++ object and deletes it; the other sees nil and does nothing. No
// double free.
//
// Invalidates any *QBitRef into a.
// Qt: QBitArray::~QBitArray().
func (a *QBitArray) Close() {
	if a == nil || a.state == nil {
		return
	}
	a.cleanup.Stop()
	if cpp := a.state.claim(); cpp != nil {
		C.qbitarray_delete(cpp)
	}
}

// Clone returns a shallow COW copy (refcount++; cheap).
// Qt: QBitArray(const QBitArray&) — copy constructor.
func (a *QBitArray) Clone() *QBitArray {
	r := wrapQBitArray(C.qbitarray_new_copy(a.handle()))
	runtime.KeepAlive(a)
	return r
}

// Assign copies the contents of other into a (COW; cheap). No-op if other
// is nil or closed. Invalidates any *QBitRef into a.
// Qt: QBitArray &operator=(const QBitArray&).
func (a *QBitArray) Assign(other *QBitArray) {
	op, ok := other.tryHandle()
	if !ok {
		return
	}
	C.qbitarray_assign(a.handle(), op)
	runtime.KeepAlive(a)
	runtime.KeepAlive(other)
}

// Swap exchanges the contents of a and other (COW swap). No-op if other is
// nil or closed. Invalidates any *QBitRef into a or other.
// Qt: QBitArray::swap(QBitArray&).
func (a *QBitArray) Swap(other *QBitArray) {
	op, ok := other.tryHandle()
	if !ok {
		return
	}
	C.qbitarray_swap(a.handle(), op)
	runtime.KeepAlive(a)
	runtime.KeepAlive(other)
}

// ---- size / status -------------------------------------------------------

// Size returns the number of bits.
// Qt: QBitArray::size() const.
func (a *QBitArray) Size() int {
	r := int(C.qbitarray_size(a.handle()))
	runtime.KeepAlive(a)
	return r
}

// Count is an alias for Size (deprecated in Qt as a no-arg count()).
// Qt: QBitArray::count() const (deprecated alias for size()).
func (a *QBitArray) Count() int { return a.Size() }

// CountSet returns the number of bits set to on.
// Qt: QBitArray::count(bool on) const.
func (a *QBitArray) CountSet(on bool) int {
	r := int(C.qbitarray_count_set(a.handle(), C.bool(on)))
	runtime.KeepAlive(a)
	return r
}

// IsEmpty reports whether Size() == 0.
// Qt: QBitArray::isEmpty() const.
func (a *QBitArray) IsEmpty() bool {
	r := bool(C.qbitarray_is_empty(a.handle()))
	runtime.KeepAlive(a)
	return r
}

// IsNull reports whether the array is null (no underlying storage allocated).
// Null is a subset of empty.
// Qt: QBitArray::isNull() const.
func (a *QBitArray) IsNull() bool {
	r := bool(C.qbitarray_is_null(a.handle()))
	runtime.KeepAlive(a)
	return r
}

// ---- modifiers -----------------------------------------------------------

// Resize changes the number of bits to size. New bits are false; removed
// bits are dropped. Invalidates any *QBitRef into a.
// Qt: QBitArray::resize(qsizetype size).
func (a *QBitArray) Resize(size int) {
	C.qbitarray_resize(a.handle(), C.ptrdiff_t(size))
	runtime.KeepAlive(a)
}

// Detach ensures that the array has its own private storage (COW detach).
// Qt: QBitArray::detach().
func (a *QBitArray) Detach() {
	C.qbitarray_detach(a.handle())
	runtime.KeepAlive(a)
}

// IsDetached reports whether the array has its own private storage.
// Qt: QBitArray::isDetached() const.
func (a *QBitArray) IsDetached() bool {
	r := bool(C.qbitarray_is_detached(a.handle()))
	runtime.KeepAlive(a)
	return r
}

// Clear resets the array to null and empty. Invalidates any *QBitRef into a.
// Qt: QBitArray::clear().
func (a *QBitArray) Clear() {
	C.qbitarray_clear(a.handle())
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
// Qt: QBitArray::testBit(qsizetype i) const.
func (a *QBitArray) TestBit(i int) bool {
	a.checkIndex(i)
	r := bool(C.qbitarray_test_bit(a.handle(), C.ptrdiff_t(i)))
	runtime.KeepAlive(a)
	return r
}

// At is a read-only alias for TestBit (C++ operator[] const).
// Qt: QBitArray::at(qsizetype i) const.
func (a *QBitArray) At(i int) bool { return a.TestBit(i) }

// SetBit sets bit i to true. Panics if i is out of range.
// Qt: QBitArray::setBit(qsizetype i).
func (a *QBitArray) SetBit(i int) {
	a.checkIndex(i)
	C.qbitarray_set_bit(a.handle(), C.ptrdiff_t(i))
	runtime.KeepAlive(a)
}

// SetBitVal sets bit i to val. Panics if i is out of range.
// Qt: QBitArray::setBit(qsizetype i, bool val).
func (a *QBitArray) SetBitVal(i int, val bool) {
	a.checkIndex(i)
	C.qbitarray_set_bit_val(a.handle(), C.ptrdiff_t(i), C.bool(val))
	runtime.KeepAlive(a)
}

// ClearBit sets bit i to false. Panics if i is out of range.
// Qt: QBitArray::clearBit(qsizetype i).
func (a *QBitArray) ClearBit(i int) {
	a.checkIndex(i)
	C.qbitarray_clear_bit(a.handle(), C.ptrdiff_t(i))
	runtime.KeepAlive(a)
}

// ToggleBit flips bit i and returns its previous value. Panics if i is out
// of range.
// Qt: QBitArray::toggleBit(qsizetype i).
func (a *QBitArray) ToggleBit(i int) bool {
	a.checkIndex(i)
	r := bool(C.qbitarray_toggle_bit(a.handle(), C.ptrdiff_t(i)))
	runtime.KeepAlive(a)
	return r
}

// ---- bitwise assign ------------------------------------------------------

// AndAssign performs a &= other. Sizes are extended to the longer of the
// two, with the shorter array treated as zero-padded on the right. No-op if
// other is nil or closed. Invalidates *QBitRef into a.
// Qt: QBitArray &operator&=(const QBitArray&).
func (a *QBitArray) AndAssign(other *QBitArray) {
	op, ok := other.tryHandle()
	if !ok {
		return
	}
	C.qbitarray_and_assign(a.handle(), op)
	runtime.KeepAlive(a)
	runtime.KeepAlive(other)
}

// OrAssign performs a |= other. No-op if other is nil or closed.
// Qt: QBitArray &operator|=(const QBitArray&).
func (a *QBitArray) OrAssign(other *QBitArray) {
	op, ok := other.tryHandle()
	if !ok {
		return
	}
	C.qbitarray_or_assign(a.handle(), op)
	runtime.KeepAlive(a)
	runtime.KeepAlive(other)
}

// XorAssign performs a ^= other. No-op if other is nil or closed.
// Qt: QBitArray &operator^=(const QBitArray&).
func (a *QBitArray) XorAssign(other *QBitArray) {
	op, ok := other.tryHandle()
	if !ok {
		return
	}
	C.qbitarray_xor_assign(a.handle(), op)
	runtime.KeepAlive(a)
	runtime.KeepAlive(other)
}

// ---- bitwise non-assign (package functions) -----------------------------

// And returns a fresh QBitArray equal to a & b. If b is nil or closed, the
// result is a clone of a (Qt treats a null QBitArray as zero-length, and a
// zero-length operand is zero-padded to the length of a — same as cloning).
// Qt: QBitArray operator&(const QBitArray&, const QBitArray&).
func And(a, b *QBitArray) *QBitArray {
	ap := a.handle()
	bp, ok := b.tryHandle()
	if !ok {
		return a.Clone()
	}
	r := wrapQBitArray(C.qbitarray_op_and(ap, bp))
	runtime.KeepAlive(a)
	runtime.KeepAlive(b)
	return r
}

// Or returns a fresh QBitArray equal to a | b. If b is nil or closed, the
// result is a clone of a.
// Qt: QBitArray operator|(const QBitArray&, const QBitArray&).
func Or(a, b *QBitArray) *QBitArray {
	ap := a.handle()
	bp, ok := b.tryHandle()
	if !ok {
		return a.Clone()
	}
	r := wrapQBitArray(C.qbitarray_op_or(ap, bp))
	runtime.KeepAlive(a)
	runtime.KeepAlive(b)
	return r
}

// Xor returns a fresh QBitArray equal to a ^ b. If b is nil or closed, the
// result is a clone of a.
// Qt: QBitArray operator^(const QBitArray&, const QBitArray&).
func Xor(a, b *QBitArray) *QBitArray {
	ap := a.handle()
	bp, ok := b.tryHandle()
	if !ok {
		return a.Clone()
	}
	r := wrapQBitArray(C.qbitarray_op_xor(ap, bp))
	runtime.KeepAlive(a)
	runtime.KeepAlive(b)
	return r
}

// Not returns a fresh QBitArray equal to ~a (all bits inverted).
// Qt: QBitArray operator~(QBitArray a).
func Not(a *QBitArray) *QBitArray {
	r := wrapQBitArray(C.qbitarray_inverted(a.handle()))
	runtime.KeepAlive(a)
	return r
}

// ---- fill / truncate -----------------------------------------------------

// Fill sets every bit of a to val (size unchanged).
// Qt: QBitArray::fill(bool val, qsizetype asize = -1) with asize < 0.
func (a *QBitArray) Fill(val bool) {
	C.qbitarray_fill_all(a.handle(), C.bool(val))
	runtime.KeepAlive(a)
}

// FillSized resizes a to size and then fills every bit with val.
// Invalidates *QBitRef into a.
// Qt: QBitArray::fill(bool val, qsizetype size).
func (a *QBitArray) FillSized(val bool, size int) {
	C.qbitarray_fill_sized(a.handle(), C.bool(val), C.ptrdiff_t(size))
	runtime.KeepAlive(a)
}

// FillRange fills bits [first, last) with val. Panics if range is invalid.
// Qt: QBitArray::fill(bool val, qsizetype first, qsizetype last).
func (a *QBitArray) FillRange(val bool, first, last int) {
	if first < 0 || last < first || last > a.Size() {
		panic("qbitarray: FillRange out of range")
	}
	C.qbitarray_fill_range(a.handle(), C.bool(val), C.ptrdiff_t(first), C.ptrdiff_t(last))
	runtime.KeepAlive(a)
}

// Truncate truncates the array at position pos (no-op if pos >= Size()).
// Invalidates *QBitRef into a.
// Qt: QBitArray::truncate(qsizetype pos).
func (a *QBitArray) Truncate(pos int) {
	C.qbitarray_truncate(a.handle(), C.ptrdiff_t(pos))
	runtime.KeepAlive(a)
}

// ---- raw bits (LSB-first within each byte) -----------------------------

// Bits returns a copy of the raw bit buffer.
// Format: bit i = (out[i/8] >> (i % 8)) & 1, len(out) == (Size()+7)/8.
// For an empty array, returns nil.
// Qt: QBitArray::bits() const.
func (a *QBitArray) Bits() []byte {
	var n C.ptrdiff_t
	p := C.qbitarray_bits(a.handle(), &n)
	runtime.KeepAlive(a)
	if p == nil || n == 0 {
		return nil
	}
	return C.GoBytes(unsafe.Pointer(p), C.int(n))
}

// ToUInt32 interprets the first 32 bits of a as uint32 using the specified
// endianness. Returns (value, false) if Size() < 32 (value is unspecified,
// matching Qt).
// Qt: QBitArray::toUInt32(QSysInfo::Endian endianness, bool *ok) const.
func (a *QBitArray) ToUInt32(littleEndian bool) (uint32, bool) {
	var ok C.bool
	v := C.qbitarray_to_uint32(a.handle(), C.bool(littleEndian), &ok)
	runtime.KeepAlive(a)
	return uint32(v), bool(ok)
}

// ---- comparison ----------------------------------------------------------

// Equal reports a == other (bit-by-bit, size and content). Returns false if
// other is nil or closed.
// Qt: QBitArray operator==(const QBitArray&).
func (a *QBitArray) Equal(other *QBitArray) bool {
	op, ok := other.tryHandle()
	if !ok {
		return false
	}
	r := bool(C.qbitarray_equal(a.handle(), op))
	runtime.KeepAlive(a)
	runtime.KeepAlive(other)
	return r
}
