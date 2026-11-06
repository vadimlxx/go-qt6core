// Package qchar provides a Go binding for Qt 6.11 QtCore's QChar class.
//
// QChar is a UTF-16 code unit. It is a value type — the binding stores it
// as uint16 and never allocates. All methods are safe for concurrent use.
//
// This package is a subpackage of github.com/example/qt6core. It links
// against Qt6Core via pkg-config; requires Qt >= 6.11 and a C++17 compiler.
//
// Not yet bound: QChar::decomposition() (returns QString, requires qstring
// package which is not implemented yet). All other public methods of
// QChar from Qt 6.11 qchar.h are exposed.
package qchar
