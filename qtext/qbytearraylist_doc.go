// Package qbytearraylist provides a Go binding for Qt 6.11
// QtCore's QByteArrayList (a.k.a. QList<QByteArray>).
//
// QByteArrayList is a typedef for QList<QByteArray>. It is a COW container;
// each *QByteArrayList is a handle to a heap-allocated C++ QList<QByteArray>.
// Clone() is cheap (refcount++); mutations detach as in Qt.
//
// Memory management: runtime.SetFinalizer; call Close() for immediate release.
//
// This package is a subpackage of github.com/vadimlxx/go-qt6core.
// It interoperates with package qbytearray via *qbytearray.QByteArray.
// Requires Qt >= 6.11 and a C++17 compiler; links via pkg-config Qt6Core.


// Interop with package qbytearray:
//
// Every *QByteArray returned by this package (At, ToSlice, Join) is a
// fully-owned *qbytearray.QByteArray — call Close() or rely on the
// finalizer. Items passed in (Append/Prepend/Insert/Set/IndexOf/...)
// are NOT adopted: the caller retains ownership and may Close them
// immediately after the call.
//
// Internally, cross-package handoff goes through qbytearray.UnsafePtr()
// and qbytearray.AdoptUnsafePtr() — the only supported mechanism for
// passing C++ QByteArray objects between subpackages.
package qtext
