// Package qbitarray provides a Go binding for Qt 6.11.2 QtCore's QBitArray and QBitRef.
//
// QBitArray is a COW (copy-on-write) array of bits, backed by a QByteArray internally. In Go,
// each *QBitArray is a handle to a heap-allocated C++ QBitArray; the wrapper owns the object
// and frees it via runtime.SetFinalizer or explicit Close().
//
// QBitRef is a proxy for a single bit inside a QBitArray. In Qt it is a lightweight value
// type returned by the mutable operator[]; in Go it is a {array, index} pair whose methods
// delegate to the parent QBitArray. No separate C++ QBitRef object is created.
//
// # Type mapping
//
//	Qt                       Go
//	QBitArray                *qbitarray.QBitArray
//	QBitRef                  *qbitarray.QBitRef
//	qsizetype                int
//	bool                     bool
//	operator& / | / ^ / ~    package functions And / Or / Xor / Not
//	operator&= / |= / ^=     methods AndAssign / OrAssign / XorAssign
//
// # Interop with sibling packages
//
// qbitarray is self-contained and does not import any other go-qt6core subpackage. The
// internal QByteArray storage is an implementation detail of Qt's QBitArray and is not
// exposed in the Go API. To exchange data with other packages, use Bits() / FromBits()
// (raw bit buffer in the same format as Qt's QBitArray::bits()).
//
// # Protected members
//
// QBitArray has no protected members; QBitRef's only private constructor
// QBitRef(QBitArray&, qsizetype) is invoked implicitly by QBitArray::operator[], which the
// Go QBitRef.BitRef method reproduces (as a {array, index} pair). No protected API needs
// to be exposed.
//
// # Bit layout (Bits / FromBits)
//
// The raw bit buffer uses Qt's standard LSB-first packing:
//
//	bit i  is  (data[i/8] >> (i % 8)) & 1
//
// Bits() returns a Go copy of `(size() + 7) / 8` bytes (unused high bits in the last byte
// are unspecified). FromBits() expects the same format, with `bitLen` being the number of
// VALID bits (not bytes), matching QStringMatcher / QBitArray::fromBits semantics.
//
// # Lifecycle of QBitRef
//
// A *QBitRef keeps its parent *QBitArray alive, but does not protect against structural
// modifications of the parent. Do NOT call Resize, Clear, Fill, FillSized, or Close on the
// parent while a QBitRef into it is in use. Read-only access after Resize to a smaller size
// will panic in Go (bounds check); mutation is undefined.
//
// # TODO (not yet bound, or API differences vs Qt 6.11.2)
//
//   - Operator overloads & / | / ^ / ~ / &= / |= / ^=: Go has no operator overloading.
//     Use package functions And(a, b) / Or(a, b) / Xor(a, b) / Not(a) and methods
//     (*QBitArray).AndAssign(other) / OrAssign / XorAssign. They are semantically equivalent.
//
//   - QBitRef::operator bool() / operator!() / operator=(bool) / operator=(const QBitRef&):
//     exposed as methods Bool() / Not() / Set(bool) / Assign(*QBitRef).
//
//   - QBitArray::data_ptr() (DataPtr = QByteArray::DataPointer): not exposed. This is a
//     private implementation detail used by Qt internals to slice the backing QByteArray.
//     If you need raw access, use Bits() / FromBits().
//
//   - QDataStream operators (operator<< / operator>> for QBitArray): requires qdatastream
//     package (not yet implemented).
//
//   - qHash(QBitArray): requires qhash package (not yet implemented). In Go, hash the byte
//     slice from Bits() (e.g. via crypto/sha256 or FNV) for map keys.
//
//   - QBitArray(const QBitArray&) / QBitArray(QBitArray&&): in Go the handle model uses
//     Clone() for explicit copies; Go assignment of a *QBitArray pointer is reference-copy,
//     not value-copy. Do NOT treat Go pointer assignment as Qt copy-construction.
//
//   - QBitArray QDebug operator<<: not applicable in Go.
//
//   - QBitArray::count() (no-arg deprecated alias for size()): not exposed; use Size().
//
//   - The Bitwise friend operators that take rvalues (& &, |, ^ with moves): not exposed;
//     AndAssign / OrAssign / XorAssign plus Clone cover the same use cases.
//
// # Requirements
//
// Qt >= 6.11.2, C++17, pkg-config Qt6Core. Version check is performed at compile time in
// qbitarray_c.cpp via internal/qt6core/check_qt_version.h.
//
// This package is a subpackage of github.com/vadimlxx/go-qt6core.
package qtext
