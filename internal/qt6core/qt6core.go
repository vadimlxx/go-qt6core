package qt6core

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC -Wno-sfinae-incomplete
#cgo LDFLAGS: -lstdc++
#include <stdint.h>
#include <stddef.h>
#include <stdlib.h>
#include "bridge.h"
*/
import "C"

import (
	"runtime"
	"unsafe"
)

// ---- QString construction (returns new QString*, caller owns) ------------

// QStringFromUtf8 allocates a new C++ QString from s.
// Caller owns the returned pointer; use QStringDelete to free it
// (typically via a finalizer in the qstring package).
func QStringFromUtf8(s string) unsafe.Pointer {
	if len(s) == 0 {
		return C.qt6core_qstring_from_utf8(nil, 0)
	}
	p := unsafe.StringData(s)
	r := C.qt6core_qstring_from_utf8((*C.char)(unsafe.Pointer(p)), C.ptrdiff_t(len(s)))
	runtime.KeepAlive(s)
	return r
}

// QStringFromUtf8Bytes is QStringFromUtf8 for []byte.
func QStringFromUtf8Bytes(b []byte) unsafe.Pointer {
	var p *C.char
	if len(b) > 0 {
		p = (*C.char)(unsafe.Pointer(&b[0]))
	}
	r := C.qt6core_qstring_from_utf8(p, C.ptrdiff_t(len(b)))
	runtime.KeepAlive(b)
	return r
}

// QStringFromUtf16 allocates a new C++ QString from UTF-16 code units.
func QStringFromUtf16(units []uint16) unsafe.Pointer {
	var p *C.uint16_t
	if len(units) > 0 {
		p = (*C.uint16_t)(unsafe.Pointer(&units[0]))
	}
	r := C.qt6core_qstring_from_utf16(p, C.ptrdiff_t(len(units)))
	runtime.KeepAlive(units)
	return r
}

// QStringFromLatin1 allocates a new C++ QString from Latin-1 bytes.
func QStringFromLatin1(b []byte) unsafe.Pointer {
	var p *C.char
	if len(b) > 0 {
		p = (*C.char)(unsafe.Pointer(&b[0]))
	}
	r := C.qt6core_qstring_from_latin1(p, C.ptrdiff_t(len(b)))
	runtime.KeepAlive(b)
	return r
}

// ---- QString extraction (returns Go-owned copies) -----------------------

// QStringToUtf8 returns a copy of q's contents as UTF-8 bytes.
func QStringToUtf8(q unsafe.Pointer) []byte {
	var data *C.char
	var size C.ptrdiff_t
	C.qt6core_qstring_to_utf8(q, &data, &size)
	if size == 0 {
		return nil
	}
	out := C.GoBytes(unsafe.Pointer(data), C.int(size))
	C.free(unsafe.Pointer(data))
	return out
}

// QStringToUtf8String returns q as a Go string.
func QStringToUtf8String(q unsafe.Pointer) string {
	b := QStringToUtf8(q)
	if len(b) == 0 {
		return ""
	}
	return string(b)
}

// QStringToUtf16 returns q as UTF-16 code units.
func QStringToUtf16(q unsafe.Pointer) []uint16 {
	var data *C.uint16_t
	var n C.ptrdiff_t
	C.qt6core_qstring_to_utf16(q, &data, &n)
	if n == 0 {
		return nil
	}
	out := make([]uint16, int(n))
	copy(out, unsafe.Slice((*uint16)(unsafe.Pointer(data)), int(n)))
	C.free(unsafe.Pointer(data))
	return out
}

// QStringToLatin1 returns q as Latin-1 bytes (chars outside Latin-1 become '?').
func QStringToLatin1(q unsafe.Pointer) []byte {
	var data *C.char
	var size C.ptrdiff_t
	C.qt6core_qstring_to_latin1(q, &data, &size)
	if size == 0 {
		return nil
	}
	out := C.GoBytes(unsafe.Pointer(data), C.int(size))
	C.free(unsafe.Pointer(data))
	return out
}

// ---- QByteArray <-> QString ---------------------------------------------

// QStringFromQByteArray returns a new QString* decoded from ba as UTF-8.
func QStringFromQByteArray(ba unsafe.Pointer) unsafe.Pointer {
	return C.qt6core_qbytearray_to_qstring(ba)
}
func QStringFromQByteArrayLatin1(ba unsafe.Pointer) unsafe.Pointer {
	return C.qt6core_qbytearray_to_qstring_latin1(ba)
}
func QStringFromQByteArrayLocal8Bit(ba unsafe.Pointer) unsafe.Pointer {
	return C.qt6core_qbytearray_to_qstring_local8bit(ba)
}

// QByteArrayFromQString returns a new QByteArray* with q encoded as UTF-8.
func QByteArrayFromQString(q unsafe.Pointer) unsafe.Pointer {
	return C.qt6core_qstring_to_qbytearray(q)
}
func QByteArrayFromQStringLatin1(q unsafe.Pointer) unsafe.Pointer {
	return C.qt6core_qstring_to_qbytearray_latin1(q)
}
func QByteArrayFromQStringLocal8Bit(q unsafe.Pointer) unsafe.Pointer {
	return C.qt6core_qstring_to_qbytearray_local8bit(q)
}

// ---- Destructors --------------------------------------------------------

// QStringDelete frees a QString* created by any QStringFrom* function.
func QStringDelete(q unsafe.Pointer) { C.qt6core_qstring_delete(q) }

// QByteArrayDelete frees a QByteArray* created by any bridge function.
//
// Most packages should NOT call this directly — wrap the result with
// qbytearray.AdoptUnsafePtr and let qbytearray's finalizer own the lifetime.
// QByteArrayDelete exists for the rare case where a package needs to
// construct a QByteArray via the bridge and free it without going through
// the qbytearray package (e.g. inside a C-callback in qstring).
func QByteArrayDelete(ba unsafe.Pointer) { C.qt6core_qbytearray_delete(ba) }

// ---- Diagnostics --------------------------------------------------------

// Version returns the Qt version the binary was linked against.
func Version() (major, minor, patch int) {
	var M, m, p C.int
	C.qt6core_get_version(&M, &m, &p)
	return int(M), int(m), int(p)
}
