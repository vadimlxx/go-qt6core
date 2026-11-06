// Package qbytearrayview provides a Go binding for Qt 6.11
// QtCore's QByteArrayView.
//
// QByteArrayView is a non-owning read-only view over a byte buffer
// (a (data, size) pair). The Go wrapper is a value type that carries
// a keepAlive reference to the underlying Go object ([]byte, string,
// or *qbytearray.QByteArray) to prevent premature GC.
//
// Lifetime rules:
//   - Slicing methods (Left/Right/Mid/First/Last/Sliced/Chopped/Trimmed/
//     Truncate/Chop/Slice) return or mutate views of the same buffer;
//     keepAlive is propagated.
//   - ToByteArray() returns a fresh *qbytearray.QByteArray (owning copy).
//   - Do not Close() or mutate the source QByteArray while a view into it
//     is alive; the C++ pointer becomes dangling.
//
// This package is a subpackage of github.com/vadimlxx/go-qt6core.
// Requires Qt >= 6.11 and a C++17 compiler; links via pkg-config Qt6Core.
package qbytearrayview
