package qbytearray

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

import "unsafe"

// ---- Cross-package interop ----------------------------------------------
//
// These functions are the ONLY supported way to hand a C++ QByteArray
// to / from another go-qt6core subpackage (qbytearraylist, qbytearrayview,
// qstring, qjson, ...).
//
// Rationale: cgo isolates C types per package, so *C.QByteArray in this
// package and *C.go_qt6core_qbytearray_t in qbytearraylist are DIFFERENT
// Go types even though they point at the same C++ object. unsafe.Pointer
// is the only portable bridge.
//
// Ownership rules:
//   - UnsafePtr: returns a borrowed pointer; caller does NOT own; the
//     QByteArray wrapper remains responsible for freeing the C++ object.
//   - AdoptUnsafePtr: takes ownership of a heap-allocated C++ QByteArray*;
//     the resulting *QByteArray finalizer will delete it.
//   - ConstDataUnsafe: returns a borrowed pointer to the internal buffer;
//     valid only while the QByteArray is alive and unmutated.

// UnsafePtr returns the underlying C++ QByteArray* as unsafe.Pointer.
// The caller does not take ownership. The C++ object stays alive as long
// as this Go wrapper is not Close()d.
func (q *QByteArray) UnsafePtr() unsafe.Pointer {
	return unsafe.Pointer(q.p)
}

// AdoptUnsafePtr wraps a heap-allocated C++ QByteArray* into a Go wrapper.
// The Go wrapper takes ownership: its finalizer (and explicit Close())
// will delete the C++ object.
//
// p must come from a function that returns a freshly new'd QByteArray*,
// e.g. internal/qt6core.QByteArrayFromQString. Passing a borrowed pointer
// (e.g. one obtained from UnsafePtr) causes a double free.
func AdoptUnsafePtr(p unsafe.Pointer) *QByteArray {
	if p == nil {
		return nil
	}
	return wrap((*C.QByteArray)(p))
}

// ConstDataUnsafe returns a borrowed pointer to the internal buffer.
// The pointer is valid until the QByteArray is mutated, resized, or Close()d.
//
// Intended for zero-copy QByteArrayView construction. The caller must keep
// q alive and unmutated for as long as the pointer is used.
func (q *QByteArray) ConstDataUnsafe() unsafe.Pointer {
	return unsafe.Pointer(C.qbytearray_data(q.p))
}
