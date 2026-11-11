package qt6core

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC
#cgo LDFLAGS: -lstdc++
#include "bridge.h"
*/
import "C"

// Version returns the Qt version the binary was linked against.
//
// It is the only remaining public function in this package. All
// QString/QByteArray bridge helpers were moved into package qtext
// (which now contains its own direct C bindings).
func Version() (major, minor, patch int) {
	var M, m, p C.int
	C.qt6core_get_version(&M, &m, &p)
	return int(M), int(m), int(p)
}
