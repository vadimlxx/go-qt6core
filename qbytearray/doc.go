// Package qbytearray provides a Go binding for Qt 6.11 QtCore's QByteArray.
//
// QByteArray is a COW (copy-on-write) byte container. In Go, each *QByteArray
// is a handle to a C++ QByteArray object; cloning is cheap (refcount++),
// and mutating methods trigger detach as in Qt.
//
// Memory management: each *QByteArray is finalized by runtime.SetFinalizer;
// call Close() to release the C++ object immediately.
//
// This package is a subpackage of github.com/vadimlxx/go-qt6core.
// Requires Qt >= 6.11 and a C++17 compiler; links via pkg-config Qt6Core.
package qbytearray
