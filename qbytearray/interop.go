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

import (
	"unsafe"

	qt6core "github.com/vadimlxx/go-qt6core/internal/qt6core"
)

// ---- Cross-package interop ----------------------------------------------
//
// These functions are the ONLY supported way to hand a C++ QByteArray to / from another
// go-qt6core subpackage (qbytearraylist, qbytearrayview, qstring, qjson, ...).
//
// Rationale: cgo isolates C types per package, so *C.QByteArray in this package and
// *C.go_qt6core_qbytearray_t in qbytearraylist are DIFFERENT Go types even though they point at
// the same C++ object. unsafe.Pointer is the only portable bridge.
//
// Ownership rules:
//   - UnsafePtr: returns a borrowed pointer; caller does NOT own; the QByteArray wrapper
//     remains responsible for freeing the C++ object.
//   - AdoptUnsafePtr: takes ownership of a heap-allocated C++ QByteArray*; the resulting
//     *QByteArray finalizer will delete it.
//   - ConstDataUnsafe: returns a borrowed pointer to the internal buffer; valid only while the
//     QByteArray is alive and unmutated.

// UnsafePtr returns the underlying C++ QByteArray* as unsafe.Pointer. The caller does not take
// ownership. The C++ object stays alive as long as this Go wrapper is not Close()d.
func (q *QByteArray) UnsafePtr() unsafe.Pointer {
	return unsafe.Pointer(q.p)
}

// AdoptUnsafePtr wraps a heap-allocated C++ QByteArray* into a Go wrapper. The Go wrapper takes
// ownership: its finalizer (and explicit Close()) will delete the C++ object.
//
// p must come from a function that returns a freshly new'd QByteArray*, e.g.
// internal/qt6core.QByteArrayFromQString. Passing a borrowed pointer (e.g. one obtained from
// UnsafePtr) causes a double free.
func AdoptUnsafePtr(p unsafe.Pointer) *QByteArray {
	if p == nil {
		return nil
	}
	return wrap((*C.QByteArray)(p))
}

// ConstDataUnsafe returns a borrowed pointer to the internal buffer. The pointer is valid until
// the QByteArray is mutated, resized, or Close()d.
//
// Intended for zero-copy QByteArrayView construction. The caller must keep q alive and unmutated
// for as long as the pointer is used.
func (q *QByteArray) ConstDataUnsafe() unsafe.Pointer {
	return unsafe.Pointer(C.qbytearray_data(q.p))
}

// ---- QString interop (unsafe.Pointer as wire type) ----------------------
//
// Because qstring imports qbytearray (for Utf8/Latin1/Local8Bit), this package cannot import
// qstring without a cycle. The bridge therefore returns raw C++ QString* via unsafe.Pointer;
// users wrap the result with qstring.AdoptUnsafePtr:
//
//   s := qstring.AdoptUnsafePtr(q.ToQString())
//   defer s.Close()
//
// or, more idiomatically:
//
//   s := qstring.FromQByteArray(q)  // lives in qstring

// ToQString returns a freshly allocated C++ QString* (decoded from UTF-8), wrapped in
// unsafe.Pointer.
//
// The caller owns the returned pointer and must wrap it via qstring.AdoptUnsafePtr.
// Equivalent to Qt's QString::fromUtf8(QByteArrayView).
func (q *QByteArray) ToQString() unsafe.Pointer {
	return qt6core.QStringFromQByteArray(q.UnsafePtr())
}

// ToQStringLatin1 is ToQString with Latin-1 decoding
// (equivalent to QString::fromLatin1(QByteArrayView)).
func (q *QByteArray) ToQStringLatin1() unsafe.Pointer {
	return qt6core.QStringFromQByteArrayLatin1(q.UnsafePtr())
}

// ToQStringLocal8Bit is ToQString with local-8-bit decoding
// (equivalent to QString::fromLocal8Bit(QByteArrayView)).
func (q *QByteArray) ToQStringLocal8Bit() unsafe.Pointer {
	return qt6core.QStringFromQByteArrayLocal8Bit(q.UnsafePtr())
}

// FromQStringUnsafePtr wraps a C++ QString* as a UTF-8 QByteArray (equivalent to
// QString::toUtf8()). Takes ownership of the QString*: it is deleted after conversion, so the
// caller must NOT Close() the source QString separately if they obtained it via UnsafePtr.
//
// Typical usage:
//
//   d := qbytearray.FromQCharDecomposition(qchar.New(0x00C0))
//   defer d.Close()
//
// Or from a QString obtained elsewhere:
//
//   q := qbytearray.FromQStringUnsafePtr(c.Decompose())
//   defer q.Close()
func FromQStringUnsafePtr(q unsafe.Pointer) *QByteArray {
	if q == nil {
		return nil
	}
	return AdoptUnsafePtr(qt6core.QByteArrayFromQString(q))
}

// FromQStringLatin1UnsafePtr wraps a C++ QString* as a Latin-1 QByteArray (QString::toLatin1()).
func FromQStringLatin1UnsafePtr(q unsafe.Pointer) *QByteArray {
	if q == nil {
		return nil
	}
	return AdoptUnsafePtr(qt6core.QByteArrayFromQStringLatin1(q))
}

// FromQStringLocal8BitUnsafePtr wraps a C++ QString* as a local-8-bit QByteArray
// (QString::toLocal8Bit()).
func FromQStringLocal8BitUnsafePtr(q unsafe.Pointer) *QByteArray {
	if q == nil {
		return nil
	}
	return AdoptUnsafePtr(qt6core.QByteArrayFromQStringLocal8Bit(q))
}