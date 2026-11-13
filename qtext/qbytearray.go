package qtext

/*
 Что сознательно не вошло:

   - fromRawData / setRawData — требуют runtime.Pinner + явного владения.
     Используйте QByteArrayFromBytes (копирует данные).
   - Итераторы (begin/end/rbegin/rend) — не идиоматичны в Go; используйте Bytes().
   - operator[] — Go не поддерживает перегрузку; используйте At/SetAt.
   - Неявные конверсии (operator const char*, std::string_view) —
     используйте Bytes()/String().
   - fromCFData / toNSData / fromEcmaUint8Array — платформо-специфично (Darwin/WASM).

 Унификация конструкторов:

   - NewQByteArray(args ...any) — диспетчер поверх именованных:
       NewQByteArray()               → пустой QByteArray
       NewQByteArray("hello")        → из строки (UTF-8)
       NewQByteArray([]byte{...})    → копия среза
       NewQByteArray(n int, c byte)  → n байт, заполненных c (см. QByteArrayFromFill)
       NewQByteArray(*QByteArray)    → COW-копия (см. Clone)

 Управление памятью:

   - QByteArray — handle-тип, каждый *QByteArray владеет C++ объектом.
   - Close() освобождает немедленно; иначе срабатывает finalizer.
   - COW на стороне Qt (Qt refcount), Clone() дёшев.
*/

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC -Wno-sfinae-incomplete
#cgo LDFLAGS: -lstdc++
#include <stdint.h>
#include <stdbool.h>
#include <stddef.h>
#include "qbytearray_c.h"
*/
import "C"

import (
	"runtime"
	"unsafe"
)

// QByteArray wraps a Qt6 QByteArray.
// Zero value is not usable; use NewQByteArray / QByteArrayFromBytes / etc.
type QByteArray struct {
	p *C.QByteArray
}

// wrapQByteArray wraps an opaque C++ QByteArray* into *QByteArray.
// Accepts unsafe.Pointer for symmetry with wrapQString and to allow the
// C-ABI to return either *C.QByteArray or go_qt6core_qbytearray_t*.
// Takes ownership: the finalizer frees the C++ object.
func wrapQByteArray(p unsafe.Pointer) *QByteArray {
	if p == nil {
		return nil
	}
	q := &QByteArray{p: (*C.QByteArray)(p)}
	runtime.SetFinalizer(q, finalizeQByteArray)
	return q
}

func finalizeQByteArray(q *QByteArray) {
	if q.p != nil {
		C.qbytearray_delete(q.p)
		q.p = nil
	}
}

// ---- construction --------------------------------------------------------

// QByteArrayFromBytes constructs a QByteArray copying b.
func QByteArrayFromBytes(b []byte) *QByteArray {
	var p *C.char
	if len(b) > 0 {
		p = (*C.char)(unsafe.Pointer(&b[0]))
	}
	q := wrapQByteArray(unsafe.Pointer(C.qbytearray_new_from_data(p, C.ptrdiff_t(len(b)))))
	runtime.KeepAlive(b)
	return q
}

// QByteArrayFromFill returns a QByteArray of size n filled with c.
func QByteArrayFromFill(n int, c byte) *QByteArray {
	return wrapQByteArray(unsafe.Pointer(
		C.qbytearray_new_fill(C.ptrdiff_t(n), C.char(c))))
}

// QByteArrayFromQCharDecomposition returns the Unicode decomposition of c
// as UTF-8 bytes (QChar::decomposition().toUtf8()).
func QByteArrayFromQCharDecomposition(c QChar) *QByteArray {
	return c.Decompose().Utf8()
}

// QByteArrayFromQCharDecompositionLatin1 is QByteArrayFromQCharDecomposition
// with Latin-1 encoding.
func QByteArrayFromQCharDecompositionLatin1(c QChar) *QByteArray {
	return c.Decompose().Latin1()
}

func NewQByteArray(args ...any) *QByteArray {
	switch len(args) {
	case 0:
		return wrapQByteArray(unsafe.Pointer(C.qbytearray_new_empty()))
	case 1:
		switch v := args[0].(type) {
		case []byte:
			return QByteArrayFromBytes(v)
		case string:
			return QByteArrayFromBytes([]byte(v))
		case QChar:
			return QByteArrayFromQCharDecomposition(v)
		}
	}

	panic("NewQByteArray accepts 0 or 1 ([]byte|string) argument")
}

// ---- lifecycle -----------------------------------------------------------

// Close releases the underlying C++ object immediately and disables the finalizer.
// Do not call Close concurrently with other methods.
func (q *QByteArray) Close() {
	if q.p != nil {
		C.qbytearray_delete(q.p)
		q.p = nil
		runtime.SetFinalizer(q, nil)
	}
}

// Clone returns a shallow COW copy (refcount++; cheap).
func (q *QByteArray) Clone() *QByteArray {
	r := wrapQByteArray(unsafe.Pointer(C.qbytearray_new_copy(q.p)))
	runtime.KeepAlive(q)
	return r
}

// ---- basic accessors -----------------------------------------------------

func (q *QByteArray) Size() int {
	n := C.qbytearray_size(q.p)
	runtime.KeepAlive(q)
	return int(n)
}
func (q *QByteArray) Len() int { return q.Size() }

func (q *QByteArray) Capacity() int {
	n := C.qbytearray_capacity(q.p)
	runtime.KeepAlive(q)
	return int(n)
}
func (q *QByteArray) IsEmpty() bool {
	v := C.qbytearray_is_empty(q.p)
	runtime.KeepAlive(q)
	return bool(v)
}
func (q *QByteArray) IsNull() bool {
	v := C.qbytearray_is_null(q.p)
	runtime.KeepAlive(q)
	return bool(v)
}
func (q *QByteArray) IsDetached() bool {
	v := C.qbytearray_is_detached(q.p)
	runtime.KeepAlive(q)
	return bool(v)
}
func (q *QByteArray) Detach() {
	C.qbytearray_detach(q.p)
	runtime.KeepAlive(q)
}

// IsSharedWith reports whether q and other share the same internal buffer.
func (q *QByteArray) IsSharedWith(other *QByteArray) bool {
	v := C.qbytearray_is_shared_with(q.p, other.p)
	runtime.KeepAlive(q)
	runtime.KeepAlive(other)
	return bool(v)
}

// QByteArrayMaxSize mirrors QByteArray::maxSize().
func QByteArrayMaxSize() int { return int(C.qbytearray_max_size()) }

// Bytes returns a copy of q's bytes.
func (q *QByteArray) Bytes() []byte {
	n := C.qbytearray_size(q.p)
	if n == 0 {
		runtime.KeepAlive(q)
		return nil
	}
	data := C.qbytearray_data(q.p)
	b := C.GoBytes(unsafe.Pointer(data), C.int(n))
	runtime.KeepAlive(q)
	return b
}

// String returns a copy of q's bytes as a Go string.
func (q *QByteArray) String() string {
	n := C.qbytearray_size(q.p)
	if n == 0 {
		runtime.KeepAlive(q)
		return ""
	}
	data := C.qbytearray_data(q.p)
	s := C.GoStringN(data, C.int(n))
	runtime.KeepAlive(q)
	return s
}

// ---- modifiers -----------------------------------------------------------

func (q *QByteArray) Resize(n int) {
	C.qbytearray_resize(q.p, C.ptrdiff_t(n))
	runtime.KeepAlive(q)
}
func (q *QByteArray) ResizeFill(n int, c byte) {
	C.qbytearray_resize_fill(q.p, C.ptrdiff_t(n), C.char(c))
	runtime.KeepAlive(q)
}
func (q *QByteArray) ResizeForOverwrite(n int) {
	C.qbytearray_resize_for_overwrite(q.p, C.ptrdiff_t(n))
	runtime.KeepAlive(q)
}
func (q *QByteArray) Fill(c byte, n int) {
	C.qbytearray_fill(q.p, C.char(c), C.ptrdiff_t(n))
	runtime.KeepAlive(q)
}
func (q *QByteArray) Clear() {
	C.qbytearray_clear(q.p)
	runtime.KeepAlive(q)
}
func (q *QByteArray) Reserve(n int) {
	C.qbytearray_reserve(q.p, C.ptrdiff_t(n))
	runtime.KeepAlive(q)
}
func (q *QByteArray) Squeeze() {
	C.qbytearray_squeeze(q.p)
	runtime.KeepAlive(q)
}
func (q *QByteArray) Truncate(pos int) {
	C.qbytearray_truncate(q.p, C.ptrdiff_t(pos))
	runtime.KeepAlive(q)
}
func (q *QByteArray) Chop(n int) {
	C.qbytearray_chop(q.p, C.ptrdiff_t(n))
	runtime.KeepAlive(q)
}

// NullTerminate appends a '\0' to the end of q in place (if not already present).
func (q *QByteArray) NullTerminate() {
	C.qbytearray_null_terminate(q.p)
	runtime.KeepAlive(q)
}

// NullTerminated returns a new QByteArray with a trailing '\0' appended.
func (q *QByteArray) NullTerminated() *QByteArray {
	r := wrapQByteArray(unsafe.Pointer(C.qbytearray_null_terminated(q.p)))
	runtime.KeepAlive(q)
	return r
}

// ---- element access ------------------------------------------------------

func (q *QByteArray) At(i int) byte {
	n := q.Size()
	if i < 0 || i >= n {
		panic("qbytearray: At out of range")
	}
	c := C.qbytearray_at(q.p, C.ptrdiff_t(i))
	runtime.KeepAlive(q)
	return byte(c)
}

func (q *QByteArray) SetAt(i int, c byte) {
	n := q.Size()
	if i < 0 || i >= n {
		panic("qbytearray: SetAt out of range")
	}
	C.qbytearray_set_at(q.p, C.ptrdiff_t(i), C.char(c))
	runtime.KeepAlive(q)
}

func (q *QByteArray) Front() byte {
	if q.IsEmpty() {
		panic("qbytearray: Front on empty QByteArray")
	}
	c := C.qbytearray_front(q.p)
	runtime.KeepAlive(q)
	return byte(c)
}

func (q *QByteArray) Back() byte {
	if q.IsEmpty() {
		panic("qbytearray: Back on empty QByteArray")
	}
	c := C.qbytearray_back(q.p)
	runtime.KeepAlive(q)
	return byte(c)
}

// ---- search --------------------------------------------------------------

func (q *QByteArray) IndexOf(c byte, from int) int {
	n := C.qbytearray_index_of_char(q.p, C.char(c), C.ptrdiff_t(from))
	runtime.KeepAlive(q)
	return int(n)
}

func (q *QByteArray) IndexOfBytes(sub []byte, from int) int {
	var p *C.char
	if len(sub) > 0 {
		p = (*C.char)(unsafe.Pointer(&sub[0]))
	}
	n := C.qbytearray_index_of_data(q.p, p, C.ptrdiff_t(len(sub)), C.ptrdiff_t(from))
	runtime.KeepAlive(q)
	runtime.KeepAlive(sub)
	return int(n)
}

func (q *QByteArray) LastIndexOf(c byte, from int) int {
	n := C.qbytearray_last_index_of_char(q.p, C.char(c), C.ptrdiff_t(from))
	runtime.KeepAlive(q)
	return int(n)
}

func (q *QByteArray) LastIndexOfBytes(sub []byte, from int) int {
	var p *C.char
	if len(sub) > 0 {
		p = (*C.char)(unsafe.Pointer(&sub[0]))
	}
	n := C.qbytearray_last_index_of_data(q.p, p, C.ptrdiff_t(len(sub)), C.ptrdiff_t(from))
	runtime.KeepAlive(q)
	runtime.KeepAlive(sub)
	return int(n)
}

func (q *QByteArray) Count(c byte) int {
	n := C.qbytearray_count_char(q.p, C.char(c))
	runtime.KeepAlive(q)
	return int(n)
}

func (q *QByteArray) CountBytes(sub []byte) int {
	var p *C.char
	if len(sub) > 0 {
		p = (*C.char)(unsafe.Pointer(&sub[0]))
	}
	n := C.qbytearray_count_data(q.p, p, C.ptrdiff_t(len(sub)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(sub)
	return int(n)
}

func (q *QByteArray) Contains(c byte) bool {
	v := C.qbytearray_contains_char(q.p, C.char(c))
	runtime.KeepAlive(q)
	return bool(v)
}

func (q *QByteArray) ContainsBytes(sub []byte) bool {
	var p *C.char
	if len(sub) > 0 {
		p = (*C.char)(unsafe.Pointer(&sub[0]))
	}
	v := C.qbytearray_contains_data(q.p, p, C.ptrdiff_t(len(sub)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(sub)
	return bool(v)
}

// Compare returns Qt's compare() result. cs == true → case-sensitive.
func (q *QByteArray) Compare(other []byte, cs bool) int {
	var p *C.char
	if len(other) > 0 {
		p = (*C.char)(unsafe.Pointer(&other[0]))
	}
	csInt := C.int(0)
	if !cs {
		csInt = 1
	}
	r := C.qbytearray_compare(q.p, p, C.ptrdiff_t(len(other)), csInt)
	runtime.KeepAlive(q)
	runtime.KeepAlive(other)
	return int(r)
}

// ---- slicing -------------------------------------------------------------

func (q *QByteArray) Left(n int) *QByteArray {
	r := wrapQByteArray(unsafe.Pointer(C.qbytearray_left(q.p, C.ptrdiff_t(n))))
	runtime.KeepAlive(q)
	return r
}

func (q *QByteArray) Right(n int) *QByteArray {
	r := wrapQByteArray(unsafe.Pointer(C.qbytearray_right(q.p, C.ptrdiff_t(n))))
	runtime.KeepAlive(q)
	return r
}

func (q *QByteArray) Mid(i, len int) *QByteArray {
	r := wrapQByteArray(unsafe.Pointer(C.qbytearray_mid(q.p, C.ptrdiff_t(i), C.ptrdiff_t(len))))
	runtime.KeepAlive(q)
	return r
}

func (q *QByteArray) First(n int) *QByteArray {
	if n < 0 || n > q.Size() {
		panic("qbytearray: First out of range")
	}
	r := wrapQByteArray(unsafe.Pointer(C.qbytearray_first(q.p, C.ptrdiff_t(n))))
	runtime.KeepAlive(q)
	return r
}

func (q *QByteArray) Last(n int) *QByteArray {
	if n < 0 || n > q.Size() {
		panic("qbytearray: Last out of range")
	}
	r := wrapQByteArray(unsafe.Pointer(C.qbytearray_last(q.p, C.ptrdiff_t(n))))
	runtime.KeepAlive(q)
	return r
}

func (q *QByteArray) Sliced(pos int) *QByteArray {
	if pos < 0 || pos > q.Size() {
		panic("qbytearray: Sliced out of range")
	}
	r := wrapQByteArray(unsafe.Pointer(C.qbytearray_sliced(q.p, C.ptrdiff_t(pos))))
	runtime.KeepAlive(q)
	return r
}

func (q *QByteArray) SlicedN(pos, n int) *QByteArray {
	sz := q.Size()
	if pos < 0 || n < 0 || pos+n > sz {
		panic("qbytearray: SlicedN out of range")
	}
	r := wrapQByteArray(unsafe.Pointer(
		C.qbytearray_sliced_n(q.p, C.ptrdiff_t(pos), C.ptrdiff_t(n))))
	runtime.KeepAlive(q)
	return r
}

func (q *QByteArray) Chopped(n int) *QByteArray {
	if n < 0 || n > q.Size() {
		panic("qbytearray: Chopped out of range")
	}
	r := wrapQByteArray(unsafe.Pointer(C.qbytearray_chopped(q.p, C.ptrdiff_t(n))))
	runtime.KeepAlive(q)
	return r
}

// ---- predicates ----------------------------------------------------------

func (q *QByteArray) StartsWith(c byte) bool {
	v := C.qbytearray_starts_with_char(q.p, C.char(c))
	runtime.KeepAlive(q)
	return bool(v)
}

func (q *QByteArray) StartsWithBytes(prefix []byte) bool {
	var p *C.char
	if len(prefix) > 0 {
		p = (*C.char)(unsafe.Pointer(&prefix[0]))
	}
	v := C.qbytearray_starts_with_data(q.p, p, C.ptrdiff_t(len(prefix)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(prefix)
	return bool(v)
}

func (q *QByteArray) EndsWith(c byte) bool {
	v := C.qbytearray_ends_with_char(q.p, C.char(c))
	runtime.KeepAlive(q)
	return bool(v)
}

func (q *QByteArray) EndsWithBytes(suffix []byte) bool {
	var p *C.char
	if len(suffix) > 0 {
		p = (*C.char)(unsafe.Pointer(&suffix[0]))
	}
	v := C.qbytearray_ends_with_data(q.p, p, C.ptrdiff_t(len(suffix)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(suffix)
	return bool(v)
}

func (q *QByteArray) IsUpper() bool {
	v := C.qbytearray_is_upper(q.p)
	runtime.KeepAlive(q)
	return bool(v)
}
func (q *QByteArray) IsLower() bool {
	v := C.qbytearray_is_lower(q.p)
	runtime.KeepAlive(q)
	return bool(v)
}
func (q *QByteArray) IsValidUtf8() bool {
	v := C.qbytearray_is_valid_utf8(q.p)
	runtime.KeepAlive(q)
	return bool(v)
}

// ---- transforms ----------------------------------------------------------

func (q *QByteArray) ToLower() *QByteArray {
	r := wrapQByteArray(unsafe.Pointer(C.qbytearray_to_lower(q.p)))
	runtime.KeepAlive(q)
	return r
}
func (q *QByteArray) ToUpper() *QByteArray {
	r := wrapQByteArray(unsafe.Pointer(C.qbytearray_to_upper(q.p)))
	runtime.KeepAlive(q)
	return r
}
func (q *QByteArray) Trimmed() *QByteArray {
	r := wrapQByteArray(unsafe.Pointer(C.qbytearray_trimmed(q.p)))
	runtime.KeepAlive(q)
	return r
}
func (q *QByteArray) Simplified() *QByteArray {
	r := wrapQByteArray(unsafe.Pointer(C.qbytearray_simplified(q.p)))
	runtime.KeepAlive(q)
	return r
}
func (q *QByteArray) LeftJustified(width int, fill byte, truncate bool) *QByteArray {
	r := wrapQByteArray(unsafe.Pointer(
		C.qbytearray_left_justified(q.p, C.ptrdiff_t(width), C.char(fill), C.bool(truncate))))
	runtime.KeepAlive(q)
	return r
}
func (q *QByteArray) RightJustified(width int, fill byte, truncate bool) *QByteArray {
	r := wrapQByteArray(unsafe.Pointer(
		C.qbytearray_right_justified(q.p, C.ptrdiff_t(width), C.char(fill), C.bool(truncate))))
	runtime.KeepAlive(q)
	return r
}
func (q *QByteArray) Repeated(times int) *QByteArray {
	r := wrapQByteArray(unsafe.Pointer(C.qbytearray_repeated(q.p, C.ptrdiff_t(times))))
	runtime.KeepAlive(q)
	return r
}

// ---- append/prepend/insert ----------------------------------------------

func (q *QByteArray) Append(c byte) {
	C.qbytearray_append_char(q.p, C.char(c))
	runtime.KeepAlive(q)
}

func (q *QByteArray) AppendBytes(b []byte) {
	var p *C.char
	if len(b) > 0 {
		p = (*C.char)(unsafe.Pointer(&b[0]))
	}
	C.qbytearray_append_data(q.p, p, C.ptrdiff_t(len(b)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(b)
}

func (q *QByteArray) AppendString(s string) { q.AppendBytes([]byte(s)) }

func (q *QByteArray) AppendCount(count int, c byte) {
	C.qbytearray_append_count(q.p, C.ptrdiff_t(count), C.char(c))
	runtime.KeepAlive(q)
}

func (q *QByteArray) AppendQByteArray(other *QByteArray) {
	C.qbytearray_append_ba(q.p, other.p)
	runtime.KeepAlive(q)
	runtime.KeepAlive(other)
}

func (q *QByteArray) Prepend(c byte) {
	C.qbytearray_prepend_char(q.p, C.char(c))
	runtime.KeepAlive(q)
}

func (q *QByteArray) PrependBytes(b []byte) {
	var p *C.char
	if len(b) > 0 {
		p = (*C.char)(unsafe.Pointer(&b[0]))
	}
	C.qbytearray_prepend_data(q.p, p, C.ptrdiff_t(len(b)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(b)
}

func (q *QByteArray) PrependString(s string) { q.PrependBytes([]byte(s)) }

func (q *QByteArray) PrependCount(count int, c byte) {
	C.qbytearray_prepend_count(q.p, C.ptrdiff_t(count), C.char(c))
	runtime.KeepAlive(q)
}

func (q *QByteArray) PrependQByteArray(other *QByteArray) {
	C.qbytearray_prepend_ba(q.p, other.p)
	runtime.KeepAlive(q)
	runtime.KeepAlive(other)
}

func (q *QByteArray) InsertBytes(i int, b []byte) {
	if i < 0 || i > q.Size() {
		panic("qbytearray: InsertBytes out of range")
	}
	var p *C.char
	if len(b) > 0 {
		p = (*C.char)(unsafe.Pointer(&b[0]))
	}
	C.qbytearray_insert_data(q.p, C.ptrdiff_t(i), p, C.ptrdiff_t(len(b)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(b)
}

func (q *QByteArray) InsertString(i int, s string) { q.InsertBytes(i, []byte(s)) }

func (q *QByteArray) InsertByte(i int, c byte) {
	if i < 0 || i > q.Size() {
		panic("qbytearray: InsertByte out of range")
	}
	C.qbytearray_insert_char(q.p, C.ptrdiff_t(i), C.char(c))
	runtime.KeepAlive(q)
}

func (q *QByteArray) InsertCount(i, count int, c byte) {
	if i < 0 || i > q.Size() {
		panic("qbytearray: InsertCount out of range")
	}
	C.qbytearray_insert_count(q.p, C.ptrdiff_t(i), C.ptrdiff_t(count), C.char(c))
	runtime.KeepAlive(q)
}

// ---- remove/replace ------------------------------------------------------

func (q *QByteArray) Remove(i, n int) {
	if i < 0 || n < 0 || i+n > q.Size() {
		panic("qbytearray: Remove out of range")
	}
	C.qbytearray_remove(q.p, C.ptrdiff_t(i), C.ptrdiff_t(n))
	runtime.KeepAlive(q)
}

func (q *QByteArray) RemoveAt(pos int) {
	C.qbytearray_remove_at(q.p, C.ptrdiff_t(pos))
	runtime.KeepAlive(q)
}

func (q *QByteArray) RemoveFirst() {
	C.qbytearray_remove_first(q.p)
	runtime.KeepAlive(q)
}

func (q *QByteArray) RemoveLast() {
	C.qbytearray_remove_last(q.p)
	runtime.KeepAlive(q)
}

func (q *QByteArray) ReplaceChar(before, after byte) {
	C.qbytearray_replace_char_char(q.p, C.char(before), C.char(after))
	runtime.KeepAlive(q)
}

func (q *QByteArray) ReplaceRange(i, n int, with []byte) {
	if i < 0 || n < 0 || i+n > q.Size() {
		panic("qbytearray: ReplaceRange out of range")
	}
	var p *C.char
	if len(with) > 0 {
		p = (*C.char)(unsafe.Pointer(&with[0]))
	}
	C.qbytearray_replace_range(q.p, C.ptrdiff_t(i), C.ptrdiff_t(n), p, C.ptrdiff_t(len(with)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(with)
}

func (q *QByteArray) ReplaceAll(before, after []byte) {
	var pb, pa *C.char
	if len(before) > 0 {
		pb = (*C.char)(unsafe.Pointer(&before[0]))
	}
	if len(after) > 0 {
		pa = (*C.char)(unsafe.Pointer(&after[0]))
	}
	C.qbytearray_replace_data_data(q.p, pb, C.ptrdiff_t(len(before)), pa, C.ptrdiff_t(len(after)))
	runtime.KeepAlive(q)
	runtime.KeepAlive(before)
	runtime.KeepAlive(after)
}

// ---- numeric conversion --------------------------------------------------

func (q *QByteArray) ToShort(base int) (int16, bool) {
	var ok C.bool
	v := C.qbytearray_to_short(q.p, &ok, C.int(base))
	runtime.KeepAlive(q)
	return int16(v), bool(ok)
}
func (q *QByteArray) ToUShort(base int) (uint16, bool) {
	var ok C.bool
	v := C.qbytearray_to_ushort(q.p, &ok, C.int(base))
	runtime.KeepAlive(q)
	return uint16(v), bool(ok)
}
func (q *QByteArray) ToInt(base int) (int32, bool) {
	var ok C.bool
	v := C.qbytearray_to_int(q.p, &ok, C.int(base))
	runtime.KeepAlive(q)
	return int32(v), bool(ok)
}
func (q *QByteArray) ToUInt(base int) (uint32, bool) {
	var ok C.bool
	v := C.qbytearray_to_uint(q.p, &ok, C.int(base))
	runtime.KeepAlive(q)
	return uint32(v), bool(ok)
}
func (q *QByteArray) ToLongLong(base int) (int64, bool) {
	var ok C.bool
	v := C.qbytearray_to_longlong(q.p, &ok, C.int(base))
	runtime.KeepAlive(q)
	return int64(v), bool(ok)
}
func (q *QByteArray) ToULongLong(base int) (uint64, bool) {
	var ok C.bool
	v := C.qbytearray_to_ulonglong(q.p, &ok, C.int(base))
	runtime.KeepAlive(q)
	return uint64(v), bool(ok)
}
func (q *QByteArray) ToFloat() (float32, bool) {
	var ok C.bool
	v := C.qbytearray_to_float(q.p, &ok)
	runtime.KeepAlive(q)
	return float32(v), bool(ok)
}
func (q *QByteArray) ToDouble() (float64, bool) {
	var ok C.bool
	v := C.qbytearray_to_double(q.p, &ok)
	runtime.KeepAlive(q)
	return float64(v), bool(ok)
}

func NumberLL(n int64, base int) *QByteArray {
	return wrapQByteArray(unsafe.Pointer(C.qbytearray_number_ll(C.int64_t(n), C.int(base))))
}
func ULL(n uint64, base int) *QByteArray {
	return wrapQByteArray(unsafe.Pointer(C.qbytearray_number_ull(C.uint64_t(n), C.int(base))))
}
func QByteArrayFromNumberDouble(v float64, format byte, precision int) *QByteArray {
	return wrapQByteArray(unsafe.Pointer(
		C.qbytearray_number_double(C.double(v), C.char(format), C.int(precision))))
}

func (q *QByteArray) SetNumLL(n int64, base int) {
	C.qbytearray_set_num_ll(q.p, C.int64_t(n), C.int(base))
	runtime.KeepAlive(q)
}
func (q *QByteArray) SetNumULL(n uint64, base int) {
	C.qbytearray_set_num_ull(q.p, C.uint64_t(n), C.int(base))
	runtime.KeepAlive(q)
}
func (q *QByteArray) SetNumDouble(v float64, format byte, precision int) {
	C.qbytearray_set_num_double(q.p, C.double(v), C.char(format), C.int(precision))
	runtime.KeepAlive(q)
}

// ---- encoding ------------------------------------------------------------

func (q *QByteArray) ToBase64(options Base64Options) *QByteArray {
	r := wrapQByteArray(unsafe.Pointer(C.qbytearray_to_base64(q.p, C.int(options))))
	runtime.KeepAlive(q)
	return r
}

func FromBase64(b []byte, options Base64Options) *QByteArray {
	src := QByteArrayFromBytes(b)
	defer src.Close()
	r := wrapQByteArray(unsafe.Pointer(C.qbytearray_from_base64(src.p, C.int(options))))
	return r
}

// QByteArrayFromBase64Encoding mirrors QByteArray::fromBase64Encoding.
func FromBase64Encoding(b []byte, options Base64Options) (*QByteArray, Base64DecodingStatus) {
	var p *C.char
	if len(b) > 0 {
		p = (*C.char)(unsafe.Pointer(&b[0]))
	}
	var status C.int
	r := wrapQByteArray(unsafe.Pointer(
		C.qbytearray_from_base64_encoding(p, C.ptrdiff_t(len(b)), C.int(options), &status)))
	runtime.KeepAlive(b)
	return r, Base64DecodingStatus(status)
}

func (q *QByteArray) ToHex(separator byte) *QByteArray {
	r := wrapQByteArray(unsafe.Pointer(C.qbytearray_to_hex(q.p, C.char(separator))))
	runtime.KeepAlive(q)
	return r
}

func FromHex(b []byte) *QByteArray {
	src := QByteArrayFromBytes(b)
	defer src.Close()
	return wrapQByteArray(unsafe.Pointer(C.qbytearray_from_hex(src.p)))
}

func (q *QByteArray) ToPercentEncoding(exclude, include []byte, percent byte) *QByteArray {
	exc := QByteArrayFromBytes(exclude)
	inc := QByteArrayFromBytes(include)
	defer exc.Close()
	defer inc.Close()
	r := wrapQByteArray(unsafe.Pointer(
		C.qbytearray_to_percent(q.p, exc.p, inc.p, C.char(percent))))
	runtime.KeepAlive(q)
	return r
}

func FromPercentEncoding(b []byte, percent byte) *QByteArray {
	src := QByteArrayFromBytes(b)
	defer src.Close()
	return wrapQByteArray(unsafe.Pointer(C.qbytearray_from_percent(src.p, C.char(percent))))
}

// ---- compression ---------------------------------------------------------

func (q *QByteArray) Compress(level int) *QByteArray {
	r := wrapQByteArray(unsafe.Pointer(C.qbytearray_compress(q.p, C.int(level))))
	runtime.KeepAlive(q)
	return r
}

func (q *QByteArray) Uncompress() *QByteArray {
	r := wrapQByteArray(unsafe.Pointer(C.qbytearray_uncompress(q.p)))
	runtime.KeepAlive(q)
	return r
}

// ---- split ---------------------------------------------------------------

// Split mirrors QByteArray::split, returning []*QByteArray.
// Each element is a separately owned handle.
func (q *QByteArray) Split(sep byte) []*QByteArray {
	var n C.ptrdiff_t
	arr := C.qbytearray_split(q.p, C.char(sep), &n)
	runtime.KeepAlive(q)
	if arr == nil || n == 0 {
		return nil
	}
	defer C.qbytearray_split_free_array(arr)

	cSlice := unsafe.Slice(arr, int(n))
	out := make([]*QByteArray, int(n))
	for i, p := range cSlice {
		out[i] = wrapQByteArray(unsafe.Pointer(p))
	}
	return out
}
