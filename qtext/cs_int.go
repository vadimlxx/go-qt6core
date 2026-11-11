package qtext

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC -Wno-sfinae-incomplete
#cgo LDFLAGS: -lstdc++
#include "qstringview_c.h"
*/
import "C"

func csInt(cs bool) C.int {
	if cs {
		return 0
	}
	return 1
}
