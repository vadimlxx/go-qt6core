package qtext

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC -Wno-sfinae-incomplete
#cgo LDFLAGS: -lstdc++
#include <stdint.h>
#include <stdbool.h>
#include <stddef.h>
#include "qstring_c.h"
#include "qbytearray_c.h"
*/
import "C"

import "unsafe"

// ---- Cross-package interop ----------------------------------------------
//
// The functions below are the ONLY supported way to exchange handle types
// (*QString, *QByteArray) with sibling go-qt6core packages (qio, qjson,
// qsettings, …).
//
// cgo isolates C types per package: *C.QString in this package and
// *C.go_qt6core_qstring_t in qio are DIFFERENT Go types, even though they
// point at the same C++ object. unsafe.Pointer is the only portable bridge.
//
// Ownership:
//   - UnsafePtr:        borrowed pointer; qtext retains ownership.
//   - *FromUnsafePtr:   takes ownership; the Go wrapper's finalizer frees
//                       the C++ object.

// ---- QByteArray ----------------------------------------------------------

// UnsafePtr returns the underlying C++ QByteArray* as unsafe.Pointer.
// Borrowed: the caller does NOT own; qtext remains responsible for freeing.
// The C++ object is valid until this wrapper is Close().
func (q *QByteArray) UnsafePtr() unsafe.Pointer {
	if q == nil || q.handle() == nil {
		return nil
	}
	return unsafe.Pointer(q.handle())
}

func QByteArrayFromUnsafePtr(p unsafe.Pointer) *QByteArray {
	return wrapQByteArray(p)
}

// ConstDataUnsafe returns a borrowed pointer to q's internal buffer.
// The pointer is valid until q is mutated, resized, or Close().
//
// Intended for zero-copy QByteArrayView construction. The caller must keep q
// alive and unmutated for as long as the pointer is used.
func (q *QByteArray) ConstDataUnsafe() unsafe.Pointer {
	return unsafe.Pointer(C.qbytearray_data(q.handle()))
}

// ---- QString -------------------------------------------------------------

// UnsafePtr returns the underlying C++ QString* as unsafe.Pointer.
// Borrowed: the caller does NOT own; qtext remains responsible for freeing.
// The C++ object is valid until this wrapper is Close().
func (q *QString) UnsafePtr() unsafe.Pointer {
	return unsafe.Pointer(q.handle())
}

func QStringFromUnsafePtr(p unsafe.Pointer) *QString {
	return wrapQString(p)
}

// ---- QStringConverter -----------------------------------------------------

/*
	Интероп-часть QStringConverter: обмен «сырыми» указателями между Go и C++.

	Функции вынесены в отдельный файл (по образцу interop.go), потому что это не
	обычное API конвертера, а мост для передачи C++ объекта между обёртками
	(например, QStringEncoder и QStringDecoder) через unsafe.Pointer.

	Правила владения:

	- UnsafePtr — заимствованный указатель; C++ объектом владеет текущая обёртка,
		и он живёт, пока она не закрыта;
	- NewQStringConverterFromUnsafePtr — обёртка-алиас БЕЗ владения: cleanup не
		регистрируется, Close() ничего не делает. Вызывающий обязан держать C++
		объект живым (обычно им владеет QStringEncoder/QStringDecoder);
	- QStringConverterAdoptUnsafePtr — устаревшее имя той же функции, оставлено
		для совместимости с кодом, написанным до переименования конструкторов
		в форму New*.
*/

// UnsafePtr returns the underlying C++ QStringConverter* as unsafe.Pointer.
//
// The caller does not take ownership: the C++ object stays alive as long as this
// Go wrapper (or the owning wrapper that adopted it, e.g. *QStringEncoder) is
// not closed. The pointer is borrowed and must not outlive the wrapper.
//
// Qt: no direct counterpart — raw pointer access to the converter (interop).
func (q *QStringConverter) UnsafePtr() unsafe.Pointer {
	p, ok := q.tryHandle()
	if !ok {
		return nil
	}
	return unsafe.Pointer(p)
}

// newQStringConverterAlias builds a non-owning *QStringConverter view over an
// existing qscState. It is used by NewQStringConverterFromUnsafePtr and by the
// QStringEncoder/QStringDecoder constructors, so an embedded converter view
// shares the state slot with the owning wrapper: once the owner closes or its
// cleanup runs, the view observes a nil handle instead of a dangling pointer.
//
// Qt: no direct counterpart — internal helper.
func newQStringConverterAlias(st *qscState) *QStringConverter {
	return &QStringConverter{
		QStringConverterBase: &QStringConverterBase{state: st},
		state:                st,
	}
}

// NewQStringConverterFromUnsafePtr wraps an existing C++ QStringConverter* as a
// *QStringConverter that does NOT own the object: no cleanup is registered and
// Close is a no-op.
//
// The caller must keep the C++ object alive (typically by owning it through a
// derived wrapper, e.g. *QStringEncoder or *QStringDecoder, which embed this
// *QStringConverter for API parity with Qt's class hierarchy).
//
// Qt: no direct counterpart — adopts an already constructed QStringConverter.
func NewQStringConverterFromUnsafePtr(p unsafe.Pointer) *QStringConverter {
	if p == nil {
		return nil
	}
	st := &qscState{}
	st.p.Store((*byte)(p))
	return newQStringConverterAlias(st)
}

// QStringConverterAdoptUnsafePtr is the historical name of
// NewQStringConverterFromUnsafePtr, kept as a source-compatibility alias for
// code written before the constructors were renamed to the New* form.
//
// Deprecated: use NewQStringConverterFromUnsafePtr.
//
// Qt: no direct counterpart — compatibility alias for the adoption helper.
func QStringConverterAdoptUnsafePtr(p unsafe.Pointer) *QStringConverter {
	return NewQStringConverterFromUnsafePtr(p)
}
