// Package qbytearray provides a Go binding for Qt 6.11 QtCore's QByteArray.
//
// QByteArray is a COW (copy-on-write) byte container. In Go, each *QByteArray is a handle to a
// C++ QByteArray object; cloning is cheap (refcount++), and mutating methods trigger detach as in
// Qt.
//
// Memory management: each *QByteArray is finalized by runtime.SetFinalizer; call Close() to
// release the C++ object immediately.
//
// # Type mapping
//
//	Qt                                  Go
//	QByteArray                          *qbytearray.QByteArray
//	QByteArray::Base64Option            qbytearray.Base64Option / Base64Options
//	QByteArray::Base64DecodingStatus    qbytearray.Base64DecodingStatus
//	const char* / char16_t*             []byte / string
//	qsizetype                           int
//
// # Interop with sibling packages
//
// qbytearray is a foundation package imported by many others:
//
//	qstring           → qbytearray (Utf8/Latin1/Local8Bit return *QByteArray)
//	qstringview       → qbytearray (ToUtf8/ToLatin1/ToLocal8Bit)
//	qlatin1stringview → qbytearray (ToUtf8)
//	qutf8stringview   → qbytearray (Compare/Equal with QByteArray)
//	qbytearraylist    → qbytearray (list element type)
//	qbytearrayview    → qbytearray (ToByteArray, FromQByteArray)
//
// Because qbytearray is imported BY these packages, it CANNOT import them back (Go import
// cycle). Cross-package helpers therefore either:
//
//   - live on the OTHER side (qstring.FromQByteArray, qstringview.FromQByteArray, …), or
//   - live here as unsafe.Pointer <-> *QByteArray adapters (ToQString, FromQStringUnsafePtr).
//
// # QString interop (unsafe.Pointer as wire type)
//
//	q := qbytearray.FromString("hello")
//	s := qstring.AdoptUnsafePtr(q.ToQString())       // decode UTF-8 → QString
//	defer s.Close()
//
// and in the other direction:
//
//	s := qstring.FromString("hello")
//	defer s.Close()
//	q := qbytearray.FromQStringUnsafePtr(s.UnsafePtr()) // QString → UTF-8
//	defer q.Close()
//
// or, more idiomatically:
//
//	q := s.Utf8()                                     // QString → QByteArray
//	s := qstring.FromQByteArray(q)                    // QByteArray → QString
//
// Latin-1 and Local8Bit variants follow the same pattern: ToQStringLatin1 /
// FromQStringLatin1UnsafePtr / ToQStringLocal8Bit / FromQStringLocal8BitUnsafePtr.
//
// # QChar decomposition interop
//
// qchar.QChar exposes the Unicode decomposition as unsafe.Pointer (mirroring
// QChar::decomposition() which returns QString). This package provides convenience wrappers that
// decode that QString as UTF-8 / Latin-1:
//
//	c := qchar.New(0x00C0)                              // À
//	d := qbytearray.FromQCharDecomposition(c)           // UTF-8 bytes of "A\u0300"
//	defer d.Close()
//	l := qbytearray.FromQCharDecompositionLatin1(c)     // Latin-1 bytes
//	defer l.Close()
//
// # QArrayData / COW inspection
//
// QByteArray::DataPointer (QArrayDataPointer<char>) is not exposed as a separate Go package. The
// COW lifecycle is managed transparently by *QByteArray's finalizer and Close(). If you need to
// inspect sharing:
//
//	q.IsDetached()                         // true if refcount == 1
//	q.IsSharedWith(other)                  // same internal buffer (ptr + size)
//
// Advanced QArrayData inspection (ref counts, capacity vs size) is not currently exposed; see
// TODO below.
//
// # Bound API (summary)
//
// Constructors:
//
//	New, FromBytes, FromString, FromFill,
//	FromBase64, FromBase64Encoding, FromHex, FromPercentEncoding,
//	FromQStringUnsafePtr, FromQStringLatin1UnsafePtr, FromQStringLocal8BitUnsafePtr,
//	FromQCharDecomposition, FromQCharDecompositionLatin1,
//	NumberLL, NumberULL, NumberDouble
//
// Size / capacity:
//
//	Size, Len, Capacity, MaxSize, IsEmpty, IsNull, IsDetached, Detach, IsSharedWith
//
// Modifiers:
//
//	Resize, ResizeFill, ResizeForOverwrite, Fill, Clear, Reserve, Squeeze, Truncate, Chop,
//	NullTerminate
//
// Element access:
//
//	At, SetAt, Front, Back, Cell/Row (not applicable to QByteArray; use At/SetAt)
//
// Slicing (return new owned QByteArray):
//
//	Left, Right, Mid, First, Last, Sliced, SlicedN, Chopped, NullTerminated
//
// Transforms:
//
//	ToLower, ToUpper, Trimmed, Simplified, Repeated, LeftJustified, RightJustified
//
// Mutation:
//
//	Append*, Prepend*, Insert*, Remove*, Replace* (multiple overload variants)
//
// Search:
//
//	IndexOf*, LastIndexOf*, Contains*, Count*, Compare, StartsWith*, EndsWith*
//
// Numeric:
//
//	ToShort, ToUShort, ToInt, ToUInt, ToLong, ToULong, ToLongLong, ToULongLong,
//	ToFloat, ToDouble, SetNumLL, SetNumULL, SetNumDouble
//
// Encoding:
//
//	ToBase64, ToHex, ToPercentEncoding
//
// Compression:
//
//	Compress, Uncompress
//
// Split:
//
//	Split(char)
//
// Predicates:
//
//	IsUpper, IsLower, IsValidUtf8
//
// QString bridge (return unsafe.Pointer to QString*; wrap via qstring.AdoptUnsafePtr):
//
//	(*QByteArray).ToQString, (*QByteArray).ToQStringLatin1, (*QByteArray).ToQStringLocal8Bit
//
// Lifecycle / views:
//
//	Clone, Close, Bytes, String, UnsafePtr, AdoptUnsafePtr, ConstDataUnsafe
//
// # TODO (not yet bound, or exposed differently vs Qt 6.11)
//
//   - fromRawData / setRawData — requires runtime.Pinner + careful lifetime management; not
//     idiomatic in Go. Use FromBytes (copies) instead.
//
//   - Iterators (begin/end/rbegin/rend, const_iterator). Use Bytes() / At / Go indexing.
//
//   - operator[] mutable reference. Use SetAt. Go has no operator overloading.
//
//   - operator std::string_view / operator std::string. Use Bytes() / String().
//
//   - QByteArrayView cross-package helpers. Use []byte at the boundary, or
//     qbytearrayview.FromQByteArray (lives in qbytearrayview to avoid an import cycle).
//
//   - QString cross-package direct returns (*qstring.QString from this package). Handled via
//     unsafe.Pointer adapters + qstring.AdoptUnsafePtr.
//
//   - assign(InputIterator, InputIterator). Use Clear + AppendBytes, or FromBytes.
//
//   - push_back / push_front. Covered by Append / Prepend.
//
//   - erase / erase_if / removeIf. Not exposed — use Bytes() + Go filtering + FromBytes, or
//     Replace for common patterns.
//
//   - fromStdString / toStdString. Covered by FromString / String().
//
//   - nullTerminated() as a Qt-style member returning a copy. Provided as NullTerminated()
//     (returns a new QByteArray) and NullTerminate() (in place).
//
//   - QByteArray::FromBase64Result as a struct. Replaced by the two-value return of
//     FromBase64Encoding: (*QByteArray, Base64DecodingStatus).
//
//   - QByteArray::DataPointer (QArrayDataPointer<char>) exposure. Handled internally; advanced
//     inspection (ref counts, alloc vs size) not exposed. See "QArrayData" above.
//
//   - operator + (concatenation). Use Clone().Append*() or build a fresh QByteArray.
//
//   - Comparison operators (==, <, …). Use Compare() or Bytes() + bytes.Equal.
//
//   - qCompress / qUncompress free functions. Provided as methods: Compress(level),
//     Uncompress().
//
//   - QDataStream operators (operator<<, operator>>). TODO — qdatastream package not yet
//     implemented.
//
//   - QByteArrayLiteral / ""_ba / ""_qba user-defined literals. Use FromString("...").
//
//   - Implicit operator const char*() / operator const void*(). Use Bytes() / String(), or
//     ConstDataUnsafe() (interop.go) for zero-copy access.
//
//   - QByteArray::operator== with QString. Use qstring.FromQByteArray(q) and .CompareString(s).
//
//   - Darwin-specific fromCFData / toCFData / fromNSData / toNSData / fromRawCFData / toRawCFData.
//     Skip — not portable.
//
//   - WASM-specific fromEcmaUint8Array / toEcmaUint8Array. Skip — not portable.
//
//   - Container STL compat (size_type, const_iterator, push_back, …). Not exposed — use
//     Bytes(), Len(), At().
//
//   - Deprecated count() (no-arg, returns size). Use Size() / Len().
//
//   - Deprecated percentDecoded() as a member. Use FromPercentEncoding().
//
// # Design notes
//
//   - QByteArray is a COW container. Clone() is O(1) (refcount++). Mutating a clone triggers
//     detach on first mutation (O(n) copy).
//
//   - Go-panic on out-of-range for At/SetAt/First/Last/Sliced/SlicedN/Chopped/Remove/Insert/
//     Replace-range. Qt release builds disable Q_ASSERT → undefined behavior; a Go panic with a
//     clear message is safer than silent UB.
//
//   - Left/Right/Mid follow Qt's clamping semantics (no panic).
//
//   - Null vs empty: QByteArray distinguishes isNull (data==nullptr) from isEmpty (size==0).
//     Both preserved. FromBytes(nil) returns a null QByteArray; FromBytes([]byte{}) returns an
//     empty (non-null) one.
//
//   - Thread safety: QByteArray is NOT thread-safe for concurrent mutation, same as Qt.
//     Clone-and-share across goroutines is safe (COW).
//
// # Requirements
//
// Qt >= 6.11, C++17, pkg-config Qt6Core. Version check is performed at compile time in
// qbytearray_c.cpp via internal/qt6core/check_qt_version.h.
//
// This package is a subpackage of github.com/vadimlxx/go-qt6core.
package qbytearray