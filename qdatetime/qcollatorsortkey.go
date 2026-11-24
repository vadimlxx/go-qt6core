package qdatetime

/*
	QCollatorSortKey — биндинг Qt 6.11 QtCore/QCollatorSortKey.

	Ключ сортировки, который возвращает QCollator::sortKey(); сравнение ключей
	дешевле повторного сравнения строк. Владеет heap-объектом; освобождается
	Close() или runtime.AddCleanup (если Close не вызывался).

	Реализовано:

	- копирующий конструктор QCollatorSortKey(const QCollatorSortKey &)
		(в Qt конструктор по умолчанию удалён: QCollatorSortKey() = delete);
	- compare(const QCollatorSortKey &) и operator< (Less);
	- swap(QCollatorSortKey &);
	- Clone()/Close().

	Не реализовано (оставлено на будущее):

	- protected-конструктор QCollatorSortKey(QCollatorSortKeyPrivate *) —
		принимает приватный тип, из Go недоступен;
	- operator= (copy assignment) — используйте Clone()/Swap();
	- move-конструктор/присваивание — не нужны для GC-управляемых handle;
	- сравнения с QCollatorSortKey через std::tuple/structured bindings —
		в Go не применимы;
	- QDataStream/qHash/QDebug-интеграция.

	Собственных enum у класса нет. Константы Qt::CaseSensitivity берутся из qdata.
*/

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC -Wno-sfinae-incomplete
#cgo LDFLAGS: -lstdc++
#include <stdint.h>
#include <stdbool.h>
#include <stddef.h>
#include "qcollator_c.h"
*/
import "C"

import (
	"runtime"
	"sync/atomic"
	"unsafe"
)

// qcollatorsortkeyState holds the C++ QCollatorSortKey pointer in an atomic
// slot so Close and the runtime cleanup can race safely.
//
// Qt: no direct counterpart — internal ownership helper.
type qcollatorsortkeyState struct {
	p atomic.Pointer[byte]
}

// cpp returns the stored QCollatorSortKey pointer, or nil if the handle was
// claimed. Qt: no direct counterpart — internal helper.
func (s *qcollatorsortkeyState) cpp() *C.go_qcollatorsortkey_t {
	if bp := s.p.Load(); bp != nil {
		return (*C.go_qcollatorsortkey_t)(unsafe.Pointer(bp))
	}
	return nil
}

// claim atomically clears the slot and returns the pointer stored there.
// Qt: no direct counterpart — internal helper.
func (s *qcollatorsortkeyState) claim() *C.go_qcollatorsortkey_t {
	if bp := s.p.Swap(nil); bp != nil {
		return (*C.go_qcollatorsortkey_t)(unsafe.Pointer(bp))
	}
	return nil
}

// QCollatorSortKey wraps a heap-allocated Qt6 QCollatorSortKey.
//
// Qt: QCollatorSortKey.
type QCollatorSortKey struct {
	state   *qcollatorsortkeyState
	cleanup runtime.Cleanup
}

// wrapQCollatorSortKey takes ownership of a raw QCollatorSortKey pointer.
// Qt: no direct counterpart — internal constructor helper.
func wrapQCollatorSortKey(p *C.go_qcollatorsortkey_t) *QCollatorSortKey {
	if p == nil {
		return nil
	}
	st := &qcollatorsortkeyState{}
	st.p.Store((*byte)(unsafe.Pointer(p)))
	k := &QCollatorSortKey{state: st}
	k.cleanup = runtime.AddCleanup(k, func(st *qcollatorsortkeyState) {
		if cpp := st.claim(); cpp != nil {
			C.qcollatorsortkey_delete(cpp)
		}
	}, st)
	return k
}

// handle returns the underlying C++ pointer, panicking if the handle is nil or
// closed. Qt: no direct counterpart — internal handle helper.
func (k *QCollatorSortKey) handle() *C.go_qcollatorsortkey_t {
	if p, ok := k.tryHandle(); ok {
		return p
	}
	panic("qcollatorsortkey: use after Close (or nil receiver)")
}

// tryHandle is the non-panicking variant of handle.
// Qt: no direct counterpart — internal handle helper.
func (k *QCollatorSortKey) tryHandle() (*C.go_qcollatorsortkey_t, bool) {
	if k == nil || k.state == nil {
		return nil, false
	}
	if p := k.state.cpp(); p != nil {
		return p, true
	}
	return nil, false
}

// Close releases the C++ object immediately and cancels the cleanup. Close is
// optional, idempotent and nil-safe.
// Qt: QCollatorSortKey::~QCollatorSortKey().
func (k *QCollatorSortKey) Close() {
	if k == nil || k.state == nil {
		return
	}
	k.cleanup.Stop()
	if p := k.state.claim(); p != nil {
		C.qcollatorsortkey_delete(p)
	}
}

// ---- construction -------------------------------------------------------

// NewQCollatorSortKey is the generic constructor. QCollatorSortKey has no
// default constructor in Qt (QCollatorSortKey() = delete), so the only accepted
// form is the copy constructor:
//
//	NewQCollatorSortKey(*QCollatorSortKey) -> QCollatorSortKey(const QCollatorSortKey &)
//
// Any other argument combination panics.
// Qt: QCollatorSortKey(const QCollatorSortKey &other).
func NewQCollatorSortKey(arg ...any) *QCollatorSortKey {
	if len(arg) == 1 {
		if other, ok := arg[0].(*QCollatorSortKey); ok {
			return NewQCollatorSortKeyCopy(other)
		}
	}
	panic("NewQCollatorSortKey: QCollatorSortKey has no default constructor; pass a *QCollatorSortKey")
}

// NewQCollatorSortKeyCopy creates an independent copy of other.
// Qt: QCollatorSortKey::QCollatorSortKey(const QCollatorSortKey &other).
func NewQCollatorSortKeyCopy(other *QCollatorSortKey) *QCollatorSortKey {
	k := wrapQCollatorSortKey(C.qcollatorsortkey_new_copy(other.handle()))
	runtime.KeepAlive(other)
	return k
}

// Clone returns an independent copy of the key.
// Qt: QCollatorSortKey::QCollatorSortKey(const QCollatorSortKey &other).
func (k *QCollatorSortKey) Clone() *QCollatorSortKey {
	return NewQCollatorSortKeyCopy(k)
}

// Swap exchanges the contents of k and other.
// Qt: QCollatorSortKey::swap(QCollatorSortKey &other).
func (k *QCollatorSortKey) Swap(other *QCollatorSortKey) {
	C.qcollatorsortkey_swap(k.handle(), other.handle())
	runtime.KeepAlive(k)
	runtime.KeepAlive(other)
}

// Compare compares k with other, returning <0, 0 or >0.
// Qt: QCollatorSortKey::compare(const QCollatorSortKey &key) const.
func (k *QCollatorSortKey) Compare(other *QCollatorSortKey) int {
	r := int(C.qcollatorsortkey_compare(k.handle(), other.handle()))
	runtime.KeepAlive(k)
	runtime.KeepAlive(other)
	return r
}

// Less reports whether k sorts before other.
// Qt: operator<(const QCollatorSortKey &lhs, const QCollatorSortKey &rhs).
func (k *QCollatorSortKey) Less(other *QCollatorSortKey) bool {
	r := bool(C.qcollatorsortkey_less(k.handle(), other.handle()))
	runtime.KeepAlive(k)
	runtime.KeepAlive(other)
	return r
}
