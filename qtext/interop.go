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
// The C++ object is valid until this wrapper is Close()d.
func (q *QByteArray) UnsafePtr() unsafe.Pointer {
	if q == nil || q.p == nil {
		return nil
	}
	return unsafe.Pointer(q.p)
}

func QByteArrayFromUnsafePtr(p unsafe.Pointer) *QByteArray {
	return wrapQByteArray(p)
}

// ConstDataUnsafe returns a borrowed pointer to q's internal buffer.
// The pointer is valid until q is mutated, resized, or Close()d.
//
// Intended for zero-copy QByteArrayView construction. The caller must keep q
// alive and unmutated for as long as the pointer is used.
func (q *QByteArray) ConstDataUnsafe() unsafe.Pointer {
	return unsafe.Pointer(C.qbytearray_data(q.p))
}

// ---- QString -------------------------------------------------------------

// UnsafePtr returns the underlying C++ QString* as unsafe.Pointer.
// Borrowed: the caller does NOT own; qtext remains responsible for freeing.
// The C++ object is valid until this wrapper is Close()d.
func (q *QString) UnsafePtr() unsafe.Pointer {
	if q == nil || q.p == nil {
		return nil
	}
	return unsafe.Pointer(q.p)
}

func QStringFromUnsafePtr(p unsafe.Pointer) *QString {
	return wrapQString(p)
}
