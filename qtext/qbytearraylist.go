package qtext

/*
 Соглашения:

   - At() возвращает копию элемента (const T& в Qt → копия в Go),
     пользователь не зависит от мутаций списка.
   - ToSlice() — замена STL-итераторов QList::begin/end.
   - Join(sep []byte) — семантика QByteArrayList::join(QByteArrayView).

 Что сознательно не вошло:

   - Java-style итераторы (QByteArrayListIterator, QMutableByteArrayListIterator) —
     не идиоматичны в Go; используйте ToSlice() или индексы.
   - STL-итераторы QList::begin/end — ToSlice() — замена.
   - operator<< / operator>> (QDataStream) — появятся вместе с пакетом qdatastream.
   - operator+ (конкатенация списков) — реализуется пользователем через
     ToSlice + NewQByteArrayList([]*QByteArray{...}).
   - emplaceBack / emplace / assign / replace / swapItemsAt / move и т.п. —
     покрывают 90% случаев; могут быть добавлены позже без ломки API.

 Унификация конструкторов:

   - NewQByteArrayList()                            → пустой список
   - NewQByteArrayList([]string{"a", "b"})          → из []string
   - NewQByteArrayList([]*QByteArray{ba1, ba2})     → из среза указателей
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
	"sync/atomic"
	"unsafe"
)

type qbytearraylistState struct {
	p atomic.Pointer[byte]
}

func (s *qbytearraylistState) cpp() *C.QByteArrayList {
	if bp := s.p.Load(); bp != nil {
		return (*C.QByteArrayList)(unsafe.Pointer(bp))
	}
	return nil
}

func (s *qbytearraylistState) claim() *C.QByteArrayList {
	if bp := s.p.Swap(nil); bp != nil {
		return (*C.QByteArrayList)(unsafe.Pointer(bp))
	}
	return nil
}

// QByteArrayList wraps a Qt6 QList<QByteArray> (a.k.a. QByteArrayList).
// The zero value is not usable; use New / FromQByteArrays / FromStrings.
type QByteArrayList struct {
	// p *C.QByteArrayList
	state   *qbytearraylistState
	cleanup runtime.Cleanup
}

func wrapQByteArrayList(l *C.QByteArrayList) *QByteArrayList {
	if l == nil {
		return nil
	}
	p := unsafe.Pointer(l)

	st := &qbytearraylistState{}
	st.p.Store((*byte)(p))
	q := &QByteArrayList{state: st}
	q.cleanup = runtime.AddCleanup(q, func(st *qbytearraylistState) {
		if cpp := st.claim(); cpp != nil {
			C.qbytearraylist_delete(cpp)
		}
	}, st)
	return q
}

func (q *QByteArrayList) handle() *C.QByteArrayList {
	if cpp, ok := q.tryHandle(); ok {
		return cpp
	}
	panic("qbytearraylist: use after Close (or nil receiver)")
}

func (q *QByteArrayList) tryHandle() (*C.QByteArrayList, bool) {
	if q == nil || q.state == nil {
		return nil, false
	}
	if cpp := q.state.cpp(); cpp != nil {
		return cpp, true
	}
	return nil, false
}

// Close releases the underlying C++ object immediately and disables the finalizer.
// Do not call Close concurrently with other methods.
func (q *QByteArrayList) Close() {
	if q == nil || q.state == nil {
		return
	}
	q.cleanup.Stop()
	if cpp := q.state.claim(); cpp != nil {
		C.qbytearraylist_delete(cpp)
	}
}

// ---- construction --------------------------------------------------------

func NewQByteArrayListFromStrings(items []string) *QByteArrayList {
	l := NewQByteArrayList()
	l.Reserve(len(items))
	for _, s := range items {
		q := NewQByteArray(s)
		l.Append(q)
		q.Close()
	}
	return l
}

func NewQByteArrayListFromQByteArrays(items []*QByteArray) *QByteArrayList {
	n := len(items)
	if n == 0 {
		return wrapQByteArrayList(C.qbytearraylist_new())
	}
	ptrs := make([]*C.QByteArray, n)
	for i, ba := range items {
		if ba == nil {
			panic("qbytearraylist: nil *QByteArray in FromQByteArrays")
		}
		ptrs[i] = ba.handle()
	}
	l := wrapQByteArrayList(C.qbytearraylist_from_array(
		(**C.QByteArray)(unsafe.Pointer(&ptrs[0])), C.ptrdiff_t(n)))
	runtime.KeepAlive(items)
	return l

}

func NewQByteArrayList(arg ...any) *QByteArrayList {
	if len(arg) == 0 || arg == nil {
		return wrapQByteArrayList(C.qbytearraylist_new())
	}

	switch items := arg[0].(type) {
	case []string:
		return NewQByteArrayListFromStrings(items)
	case []*QByteArray:
		return NewQByteArrayListFromQByteArrays(items)
	}

	panic("NewQByteArrayList invalid argumentes")
}

// Clone returns a shallow COW copy (Qt refcount++; cheap).
func (l *QByteArrayList) Clone() *QByteArrayList {
	r := wrapQByteArrayList(C.qbytearraylist_new_copy(l.handle()))
	runtime.KeepAlive(l)
	return r
}

// ---- size / capacity -----------------------------------------------------

func (l *QByteArrayList) Size() int {
	n := C.qbytearraylist_size(l.handle())
	runtime.KeepAlive(l)
	return int(n)
}

func (l *QByteArrayList) Len() int { return l.Size() }

func (l *QByteArrayList) IsEmpty() bool {
	v := C.qbytearraylist_is_empty(l.handle())
	runtime.KeepAlive(l)
	return bool(v)
}

func (l *QByteArrayList) Capacity() int {
	v := C.qbytearraylist_capacity(l.handle())
	runtime.KeepAlive(l)
	return int(v)
}

func (l *QByteArrayList) Clear()   { C.qbytearraylist_clear(l.handle()); runtime.KeepAlive(l) }
func (l *QByteArrayList) Reserve(n int) { C.qbytearraylist_reserve(l.handle(), C.ptrdiff_t(n)); runtime.KeepAlive(l) }
func (l *QByteArrayList) Squeeze() { C.qbytearraylist_squeeze(l.handle()); runtime.KeepAlive(l) }

// ---- append / prepend / insert ------------------------------------------

func (l *QByteArrayList) Append(ba *QByteArray) {
	C.qbytearraylist_append(l.handle(), ba.handle())
	runtime.KeepAlive(l)
	runtime.KeepAlive(ba)
}

func (l *QByteArrayList) Prepend(ba *QByteArray) {
	C.qbytearraylist_prepend(l.handle(), ba.handle())
	runtime.KeepAlive(l)
	runtime.KeepAlive(ba)
}

func (l *QByteArrayList) Insert(i int, ba *QByteArray) {
	if i < 0 || i > l.Size() {
		panic("qbytearraylist: Insert out of range")
	}
	C.qbytearraylist_insert(l.handle(), C.ptrdiff_t(i), ba.handle())
	runtime.KeepAlive(l)
	runtime.KeepAlive(ba)
}

// ---- element access (returns a fresh copy) ------------------------------

// At returns a new *QByteArray with a copy of the i-th element.
// Modifying the returned object does not affect the list.
func (l *QByteArrayList) At(i int) *QByteArray {
	if i < 0 || i >= l.Size() {
		panic("qbytearraylist: At out of range")
	}
	var out *C.QByteArray
	C.qbytearraylist_at_new(l.handle(), C.ptrdiff_t(i), &out)
	runtime.KeepAlive(l)
	return wrapQByteArray(unsafe.Pointer(out))
}

func (l *QByteArrayList) Set(i int, ba *QByteArray) {
	if i < 0 || i >= l.Size() {
		panic("qbytearraylist: Set out of range")
	}
	C.qbytearraylist_set(l.handle(), C.ptrdiff_t(i), ba.handle())
	runtime.KeepAlive(l)
	runtime.KeepAlive(ba)
}

// ---- remove --------------------------------------------------------------

func (l *QByteArrayList) RemoveAt(i int) {
	if i < 0 || i >= l.Size() {
		panic("qbytearraylist: RemoveAt out of range")
	}
	C.qbytearraylist_remove_at(l.handle(), C.ptrdiff_t(i))
	runtime.KeepAlive(l)
}

func (l *QByteArrayList) RemoveFirst() {
	C.qbytearraylist_remove_first(l.handle())
	runtime.KeepAlive(l)
}

func (l *QByteArrayList) RemoveLast() {
	C.qbytearraylist_remove_last(l.handle())
	runtime.KeepAlive(l)
}

// ---- search --------------------------------------------------------------

func (l *QByteArrayList) IndexOf(ba *QByteArray) int {
	v := C.qbytearraylist_index_of(l.handle(), ba.handle())
	runtime.KeepAlive(l)
	runtime.KeepAlive(ba)
	return int(v)
}

func (l *QByteArrayList) LastIndexOf(ba *QByteArray) int {
	v := C.qbytearraylist_last_index_of(l.handle(), ba.handle())
	runtime.KeepAlive(l)
	runtime.KeepAlive(ba)
	return int(v)
}

func (l *QByteArrayList) Contains(ba *QByteArray) bool {
	v := C.qbytearraylist_contains(l.handle(), ba.handle())
	runtime.KeepAlive(l)
	runtime.KeepAlive(ba)
	return bool(v)
}

// ---- join (QByteArrayList-specific) -------------------------------------

// Join concatenates all elements using sep as separator
// (mirrors QByteArrayList::join(QByteArrayView)).
func (l *QByteArrayList) Join(sep []byte) *QByteArray {
	var p *C.char
	if len(sep) > 0 {
		p = (*C.char)(unsafe.Pointer(&sep[0]))
	}
	var out *C.QByteArray
	C.qbytearraylist_join(l.handle(), p, C.ptrdiff_t(len(sep)), &out)
	runtime.KeepAlive(l)
	runtime.KeepAlive(sep)
	return wrapQByteArray(unsafe.Pointer(out))
}

// JoinString is Join with a string separator.
func (l *QByteArrayList) JoinString(sep string) *QByteArray {
	return l.Join([]byte(sep))
}

// JoinByte is Join with a single-byte separator.
func (l *QByteArrayList) JoinByte(sep byte) *QByteArray {
	return l.Join([]byte{sep})
}

// ---- bulk extraction -----------------------------------------------------

// ToSlice returns all elements as fresh *QByteArray values.
// Caller owns each returned value (must Close or let finalizer run).
func (l *QByteArrayList) ToSlice() []*QByteArray {
	var arr **C.QByteArray
	var n C.ptrdiff_t
	C.qbytearraylist_to_array_new(l.handle(), &arr, &n)
	runtime.KeepAlive(l)
	if arr == nil || n == 0 {
		return nil
	}
	defer C.qbytearraylist_free_array(arr)

	cslice := unsafe.Slice(arr, int(n))
	out := make([]*QByteArray, int(n))
	for i, p := range cslice {
		out[i] = wrapQByteArray(unsafe.Pointer(p))
	}
	return out
}
