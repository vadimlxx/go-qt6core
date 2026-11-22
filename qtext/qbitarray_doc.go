// QBitArray is a Go binding for Qt 6.11.2 QtCore's QBitArray and QBitRef.
//
// QBitArray is a COW (copy-on-write) array of bits, backed by a QByteArray
// internally. In Go, each *QBitArray is a handle to a heap-allocated C++
// QBitArray; the wrapper owns the object and frees it via runtime.AddCleanup
// or explicit Close().
//
// QBitRef is a proxy for a single bit inside a QBitArray. In Qt it is a
// lightweight value type returned by the mutable operator[]; in Go it is a
// {array, index} pair whose methods delegate to the parent QBitArray. No
// separate C++ QBitRef object is created.
//
// # Ownership
//
// The C++ QBitArray is released by one of:
//   - Close(): eager, deterministic release. Use this in long-running
//     processes when you know a QBitArray is no longer needed.
//   - runtime.AddCleanup: safety net. If Close is never called, the C++
//     object is freed once the *QBitArray becomes unreachable.
//
// Close is optional, idempotent, and nil-safe. If Close races with the
// runtime cleanup, exactly one of them claims the C++ object (via atomic
// Swap) and deletes it; the other sees nil and does nothing. No double free.
//
// Using a QBitArray after Close panics with
// "qbitarray: use after Close (or nil receiver)".
//
// # Type mapping
//
//	Qt                       Go
//	QBitArray                *qtext.QBitArray
//	QBitRef                  *qtext.QBitRef
//	qsizetype                int
//	bool                     bool
//	operator& / | / ^ / ~    package functions And / Or / Xor / Not
//	operator&= / |= / ^=     methods AndAssign / OrAssign / XorAssign
//
// # Interop with sibling packages
//
// qbitarray is self-contained and does not import any other go-qt6core
// subpackage. The internal QByteArray storage is an implementation detail of
// Qt's QBitArray and is not exposed in the Go API. To exchange data with
// other packages, use Bits() / NewQBitArrayFromBits() (raw bit buffer in the
// same format as Qt's QBitArray::bits()).
//
// # Protected members
//
// QBitArray has no protected members; QBitRef's only private constructor
// QBitRef(QBitArray&, qsizetype) is invoked implicitly by
// QBitArray::operator[], which the Go QBitRef.BitRef method reproduces (as a
// {array, index} pair). No protected API needs to be exposed.
//
// # Bit layout (Bits / NewQBitArrayFromBits)
//
// The raw bit buffer uses Qt's standard LSB-first packing:
//
//	bit i  is  (data[i/8] >> (i % 8)) & 1
//
// Bits() returns a Go copy of (size()+7)/8 bytes (unused high bits in the
// last byte are unspecified). NewQBitArrayFromBits() expects the same
// format, with bitLen being the number of VALID bits (not bytes), matching
// QBitArray::fromBits semantics.
//
// # Lifecycle of QBitRef
//
// A *QBitRef keeps its parent *QBitArray alive, but does not protect against
// structural modifications of the parent. Do NOT call Resize, Clear, Fill,
// FillSized, FillRange, or Close on the parent while a QBitRef into it is in
// use. Read-only access after Resize to a smaller size panics in Go (bounds
// check). After Close of the parent, any QBitRef method panics (through the
// parent's handle() guard) — it does not exhibit undefined behaviour.
//
// # Not yet bound
//
//   - Operator overloads & / | / ^ / ~ / &= / |= / ^=: Go has no operator
//     overloading. Use package functions And(a, b) / Or(a, b) / Xor(a, b) /
//     Not(a) and methods AndAssign / OrAssign / XorAssign. They are
//     semantically equivalent.
//
//   - QBitRef::operator bool() / operator!() / operator=(bool) /
//     operator=(const QBitRef&): exposed as methods Bool() / Not() /
//     Set(bool) / Assign(*QBitRef).
//
//   - QBitArray::data_ptr() (DataPtr = QByteArray::DataPointer): not exposed.
//     This is a private implementation detail used by Qt internals to slice
//     the backing QByteArray. If you need raw access, use Bits() /
//     NewQBitArrayFromBits().
//
//   - QDataStream operators (operator<< / operator>> for QBitArray): requires
//     the qdatastream package (not yet implemented).
//
//   - qHash(QBitArray): requires the qhash package (not yet implemented). In
//     Go, hash the byte slice from Bits() (e.g. via crypto/sha256 or FNV) for
//     map keys.
//
//   - QBitArray(const QBitArray&) / QBitArray(QBitArray&&): in Go the handle
//     model uses Clone() for explicit copies; Go assignment of a *QBitArray
//     pointer is reference-copy, not value-copy. Do NOT treat Go pointer
//     assignment as Qt copy-construction.
//
//   - QBitArray QDebug operator<<: not applicable in Go.
//
//   - QBitArray::count() (no-arg deprecated alias for size()): not exposed
//     as a separate method — use Size() (or the Count() convenience alias,
//     which just forwards to Size()).
//
//   - The bitwise friend operators that take rvalues (&, |, ^ with moves):
//     not exposed; AndAssign / OrAssign / XorAssign plus Clone cover the same
//     use cases.
//
// # Requirements
//
// Qt >= 6.11.2, C++17, pkg-config Qt6Core. Version check is performed at
// compile time in qbitarray_c.cpp via internal/qt6core/check_qt_version.h.
//
// This package is a subpackage of github.com/vadimlxx/go-qt6core.
package qtext
