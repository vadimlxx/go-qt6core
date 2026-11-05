package qbytearraylist

/*
 QList::at() возвращает const T&. Копия изолирует Go-пользователя от мутаций списка; поведение совпадает с семантикой Qt (const ref).

 Что сознательно не вошло
 Java-style итераторы (QByteArrayListIterator, QMutableByteArrayListIterator) — не идиоматичны в Go; используйте ToSlice() или методы по индексу.

 STL-итераторы QList::begin/end — Go их не поддерживает; ToSlice() — замена.

 operator<</operator>> (QDataStream) — появятся вместе с пакетом qdatastream.

 operator+ (конкатенация списков) — не идиоматичен в Go; реализуется пользователем через ToSlice + FromQByteArrays.

 QList-методы, специфичные для Qt6 (emplaceBack, emplace, assign, replace, swapItemsAt, move и т.п.) — покрывают 90% случаев; могут быть добавлены позже без ломки API.
*/

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC -Wno-sfinae-incomplete
#cgo LDFLAGS: -lstdc++
#include <stdint.h>
#include <stdbool.h>
#include <stddef.h>
#include "qbytearraylist_c.h"
*/
import "C"

import (
	"runtime"
	"unsafe"

	"github.com/vadimlxx/go-qt6core/qbytearray"
)

// QByteArrayList wraps a Qt6 QList<QByteArray> (a.k.a. QByteArrayList).
// The zero value is not usable; use New / FromQByteArrays / FromStrings.
type QByteArrayList struct {
	p *C.go_qt6core_qbytearraylist_t
}

func wrap(p *C.go_qt6core_qbytearraylist_t) *QByteArrayList {
	if p == nil {
		return nil
	}
	l := &QByteArrayList{p: p}
	runtime.SetFinalizer(l, finalize)
	return l
}

func finalize(l *QByteArrayList) {
	if l.p != nil {
		C.qbytearraylist_delete(l.p)
		l.p = nil
	}
}

// Close releases the underlying C++ object immediately and disables the finalizer.
// Do not call Close concurrently with other methods.
func (l *QByteArrayList) Close() {
	if l.p != nil {
		C.qbytearraylist_delete(l.p)
		l.p = nil
		runtime.SetFinalizer(l, nil)
	}
}

// ---- construction --------------------------------------------------------

// New returns a new empty QByteArrayList.
func New() *QByteArrayList { return wrap(C.qbytearraylist_new()) }

// FromQByteArrays builds a QByteArrayList by copying each item.
// The items themselves are not adopted; caller still owns them.
func FromQByteArrays(items []*qbytearray.QByteArray) *QByteArrayList {
	n := len(items)
	if n == 0 {
		return New()
	}
	ptrs := make([]*C.go_qt6core_qbytearray_t, n)
	for i, ba := range items {
		if ba == nil {
			panic("qbytearraylist: nil *QByteArray in FromQByteArrays")
		}
		ptrs[i] = (*C.go_qt6core_qbytearray_t)(ba.UnsafePtr())
	}
	l := wrap(C.qbytearraylist_from_array(
		(**C.go_qt6core_qbytearray_t)(unsafe.Pointer(&ptrs[0])),
		C.ptrdiff_t(n)))
	runtime.KeepAlive(items)
	return l
}

// FromStrings builds a QByteArrayList from Go strings.
func FromStrings(items []string) *QByteArrayList {
	l := New()
	l.Reserve(len(items))
	for _, s := range items {
		q := qbytearray.FromString(s)
		l.Append(q)
		q.Close()
	}
	return l
}

// Clone returns a shallow COW copy (Qt refcount++; cheap).
func (l *QByteArrayList) Clone() *QByteArrayList {
	r := wrap(C.qbytearraylist_new_copy(l.p))
	runtime.KeepAlive(l)
	return r
}

// ---- size / capacity -----------------------------------------------------

func (l *QByteArrayList) Size() int {
	n := C.qbytearraylist_size(l.p)
	runtime.KeepAlive(l)
	return int(n)
}

func (l *QByteArrayList) Len() int { return l.Size() }

func (l *QByteArrayList) IsEmpty() bool {
	v := C.qbytearraylist_is_empty(l.p)
	runtime.KeepAlive(l)
	return bool(v)
}

func (l *QByteArrayList) Capacity() int {
	v := C.qbytearraylist_capacity(l.p)
	runtime.KeepAlive(l)
	return int(v)
}

func (l *QByteArrayList) Clear()   { C.qbytearraylist_clear(l.p); runtime.KeepAlive(l) }
func (l *QByteArrayList) Reserve(n int) { C.qbytearraylist_reserve(l.p, C.ptrdiff_t(n)); runtime.KeepAlive(l) }
func (l *QByteArrayList) Squeeze() { C.qbytearraylist_squeeze(l.p); runtime.KeepAlive(l) }

// ---- append / prepend / insert ------------------------------------------

func (l *QByteArrayList) Append(ba *qbytearray.QByteArray) {
	C.qbytearraylist_append(l.p, (*C.go_qt6core_qbytearray_t)(ba.UnsafePtr()))
	runtime.KeepAlive(l)
	runtime.KeepAlive(ba)
}

func (l *QByteArrayList) Prepend(ba *qbytearray.QByteArray) {
	C.qbytearraylist_prepend(l.p, (*C.go_qt6core_qbytearray_t)(ba.UnsafePtr()))
	runtime.KeepAlive(l)
	runtime.KeepAlive(ba)
}

func (l *QByteArrayList) Insert(i int, ba *qbytearray.QByteArray) {
	if i < 0 || i > l.Size() {
		panic("qbytearraylist: Insert out of range")
	}
	C.qbytearraylist_insert(l.p, C.ptrdiff_t(i), (*C.go_qt6core_qbytearray_t)(ba.UnsafePtr()))
	runtime.KeepAlive(l)
	runtime.KeepAlive(ba)
}

// ---- element access (returns a fresh copy) ------------------------------

// At returns a new *qbytearray.QByteArray with a copy of the i-th element.
// Modifying the returned object does not affect the list.
func (l *QByteArrayList) At(i int) *qbytearray.QByteArray {
	if i < 0 || i >= l.Size() {
		panic("qbytearraylist: At out of range")
	}
	var out *C.go_qt6core_qbytearray_t
	C.qbytearraylist_at_new(l.p, C.ptrdiff_t(i), &out)
	runtime.KeepAlive(l)
	return qbytearray.AdoptUnsafePtr(unsafe.Pointer(out))
}

func (l *QByteArrayList) Set(i int, ba *qbytearray.QByteArray) {
	if i < 0 || i >= l.Size() {
		panic("qbytearraylist: Set out of range")
	}
	C.qbytearraylist_set(l.p, C.ptrdiff_t(i), (*C.go_qt6core_qbytearray_t)(ba.UnsafePtr()))
	runtime.KeepAlive(l)
	runtime.KeepAlive(ba)
}

// ---- remove --------------------------------------------------------------

func (l *QByteArrayList) RemoveAt(i int) {
	if i < 0 || i >= l.Size() {
		panic("qbytearraylist: RemoveAt out of range")
	}
	C.qbytearraylist_remove_at(l.p, C.ptrdiff_t(i))
	runtime.KeepAlive(l)
}

func (l *QByteArrayList) RemoveFirst() {
	C.qbytearraylist_remove_first(l.p)
	runtime.KeepAlive(l)
}

func (l *QByteArrayList) RemoveLast() {
	C.qbytearraylist_remove_last(l.p)
	runtime.KeepAlive(l)
}

// ---- search --------------------------------------------------------------

func (l *QByteArrayList) IndexOf(ba *qbytearray.QByteArray) int {
	v := C.qbytearraylist_index_of(l.p, (*C.go_qt6core_qbytearray_t)(ba.UnsafePtr()))
	runtime.KeepAlive(l)
	runtime.KeepAlive(ba)
	return int(v)
}

func (l *QByteArrayList) LastIndexOf(ba *qbytearray.QByteArray) int {
	v := C.qbytearraylist_last_index_of(l.p, (*C.go_qt6core_qbytearray_t)(ba.UnsafePtr()))
	runtime.KeepAlive(l)
	runtime.KeepAlive(ba)
	return int(v)
}

func (l *QByteArrayList) Contains(ba *qbytearray.QByteArray) bool {
	v := C.qbytearraylist_contains(l.p, (*C.go_qt6core_qbytearray_t)(ba.UnsafePtr()))
	runtime.KeepAlive(l)
	runtime.KeepAlive(ba)
	return bool(v)
}

// ---- join (QByteArrayList-specific) -------------------------------------

// Join concatenates all elements using sep as separator
// (mirrors QByteArrayList::join(QByteArrayView)).
func (l *QByteArrayList) Join(sep []byte) *qbytearray.QByteArray {
	var p *C.char
	if len(sep) > 0 {
		p = (*C.char)(unsafe.Pointer(&sep[0]))
	}
	var out *C.go_qt6core_qbytearray_t
	C.qbytearraylist_join(l.p, p, C.ptrdiff_t(len(sep)), &out)
	runtime.KeepAlive(l)
	runtime.KeepAlive(sep)
	return qbytearray.AdoptUnsafePtr(unsafe.Pointer(out))
}

// JoinString is Join with a string separator.
func (l *QByteArrayList) JoinString(sep string) *qbytearray.QByteArray {
	return l.Join([]byte(sep))
}

// JoinByte is Join with a single-byte separator.
func (l *QByteArrayList) JoinByte(sep byte) *qbytearray.QByteArray {
	return l.Join([]byte{sep})
}

// ---- bulk extraction -----------------------------------------------------

// ToSlice returns all elements as fresh *qbytearray.QByteArray values.
// Caller owns each returned value (must Close or let finalizer run).
func (l *QByteArrayList) ToSlice() []*qbytearray.QByteArray {
	var arr **C.go_qt6core_qbytearray_t
	var n C.ptrdiff_t
	C.qbytearraylist_to_array_new(l.p, &arr, &n)
	runtime.KeepAlive(l)
	if arr == nil || n == 0 {
		return nil
	}
	defer C.qbytearraylist_free_array(arr)

	cslice := unsafe.Slice(arr, int(n))
	out := make([]*qbytearray.QByteArray, int(n))
	for i, p := range cslice {
		out[i] = qbytearray.AdoptUnsafePtr(unsafe.Pointer(p))
	}
	return out
}
