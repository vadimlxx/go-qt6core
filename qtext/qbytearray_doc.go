// Class QByteArray provides a Go binding for Qt 6.11 QtCore's QByteArray.
//
// QByteArray is a COW (copy-on-write) byte container. In Go, each *QByteArray is
// a handle to a C++ QByteArray object; cloning is cheap (refcount++), and
// mutating methods trigger detach as in Qt.
//
// # Type mapping
//
//	Qt                                  Go
//	QByteArray                          *qtext.QByteArray
//	QByteArray::Base64Option            qtext.Base64Option / Base64Options
//	QByteArray::Base64DecodingStatus    qtext.Base64DecodingStatus
//	const char *                        string (UTF-8 bytes) / []byte
//	QByteArrayView                      []byte
//	qsizetype                           int
//
// # Two variants per operation
//
// Every member operation whose Qt signature takes or returns QByteArray or
// QByteArrayView has two entry points (the only exceptions are IsSharedWith and
// Swap, for which []byte counterparts would be meaningless):
//
//   - a handle variant operating on *QByteArray (ownership, COW, no copy);
//   - a Go variant operating on []byte. Methods returning a new QByteArray have
//     a ...Bytes twin that returns a copy, and methods taking *QByteArray have a
//     []byte twin (AppendBytes, PrependBytes, InsertBytes, ...).
//
// Parameters typed const char * in Qt additionally have ...String conveniences:
// NewQByteArrayFromString, AppendString, PrependString, InsertString,
// AssignString and ReplaceAllString.
//
// QByteArray's Qt 6.11 API takes no QString, QStringView, QLatin1StringView or
// QStringList parameters or return values (the QString comparison operators were
// removed in Qt 6.8), so no string/[]string overloads are required for them.
//
// # Construction
//
//	NewQByteArray(v ...any)             // dispatch by Go type
//	NewQByteArrayFromBytes([]byte)      // nil -> null, non-nil empty -> empty
//	NewQByteArrayFromString(string)     // UTF-8, embedded NULs preserved
//	NewQByteArrayFromFill(n, c)         // n bytes filled with c
//	NewQByteArrayNumberLL/ULL/Int/UInt/Double
//	NewQByteArrayFromBase64 / FromBase64Encoding / FromHex / FromPercentEncoding
//	NewQByteArrayFromQCharDecomposition / ...Latin1
//
// # Lifecycle
//
//	Close, Clone, CloneBytes, Swap
//	AssignBytes, AssignString, AssignFill, AssignQByteArray
//
// A *QByteArray is freed by Close() (immediate, idempotent, nil-safe) or by a
// runtime cleanup attached when the wrapper becomes unreachable. Close and the
// cleanup race safely: exactly one of them claims the C++ object, so there is no
// double free.
//
// # Size / capacity
//
//	Size, Len, Capacity, QByteArrayMaxSize,
//	IsEmpty, IsNull, IsDetached, Detach, IsSharedWith
//
// # Modifiers
//
//	Resize, ResizeFill, ResizeForOverwrite, Fill, Clear, Reserve, Squeeze,
//	Truncate, Chop, NullTerminate, NullTerminated(+Bytes)
//
// # Element access
//
//	At, SetAt, Front, SetFront, Back, SetBack
//
// # Search
//
//	IndexOf, IndexOfBytes, LastIndexOf, LastIndexOfBytes,
//	Count, CountBytes, Contains, ContainsBytes, Compare,
//	StartsWith, StartsWithBytes, EndsWith, EndsWithBytes
//
// # Slicing
//
//	Left(+Bytes), Right(+Bytes), Mid(+Bytes), First(+Bytes), Last(+Bytes),
//	Sliced(+Bytes), SlicedN(+Bytes), Chopped(+Bytes)  // return a new array
//	Slice, SliceN                                    // modify in place
//
// # Transforms
//
//	ToLower(+Bytes), ToUpper(+Bytes), Trimmed(+Bytes), Simplified(+Bytes),
//	LeftJustified(+Bytes), RightJustified(+Bytes), Repeated(+Bytes),
//	Concat(+Bytes)
//
// # Mutation
//
//	Append(char), AppendBytes, AppendString, AppendCount, AppendQByteArray
//	Prepend(char), PrependBytes, PrependString, PrependCount, PrependQByteArray
//	InsertBytes, InsertString, InsertByte, InsertCount, InsertQByteArray
//	Remove, RemoveAt, RemoveFirst, RemoveLast, RemoveIf
//	ReplaceChar, ReplaceCharBytes, ReplaceRange, ReplaceAll, ReplaceAllString
//
// # Predicates
//
//	IsUpper, IsLower, IsValidUtf8
//
// # Numeric
//
//	ToShort, ToUShort, ToInt, ToUInt, ToLong, ToULong,
//	ToLongLong, ToULongLong, ToFloat, ToDouble
//	SetNumLL, SetNumULL, SetNumInt, SetNumUInt, SetNumFloat, SetNumDouble
//
// # Encoding
//
//	ToBase64(+Bytes), ToHex(+Bytes), ToPercentEncoding(+Bytes),
//	PercentDecoded(+Bytes)
//
// # Compression and splitting
//
//	Compress(+Bytes), Uncompress(+Bytes)
//	Split(char) []*QByteArray, SplitBytes(char) [][]byte
//
// # Interop
//
//	UnsafePtr, ConstDataUnsafe, QByteArrayFromUnsafePtr (interop.go)
//
// These exchange borrowed/owned handles with sibling packages without exposing
// cgo types.
//
// # Not bound, or exposed differently than in Qt 6.11
//
//   - fromRawData / setRawData: would require runtime.Pinner and manual lifetime
//     management. Use NewQByteArrayFromBytes, which copies.
//
//   - Iterators (begin/end/rbegin/rend), mutable data(), operator[] references:
//     replaced by At/SetAt/Bytes. Qt's operator[] is a reference; Go has no
//     operator overloading.
//
//   - operator const char * / operator const void * / operator std::string_view:
//     use Bytes(), String() or ConstDataUnsafe().
//
//   - QByteArray::FromBase64Result: replaced by the two-value return of
//     NewQByteArrayFromBase64Encoding.
//
//   - QByteArray::DataPointer (QArrayDataPointer<char>): internal; only
//     IsDetached/IsSharedWith/Capacity are exposed.
//
//   - QDataStream operators and qHash(QByteArray): the qdatastream/qhash parts
//     are not implemented yet.
//
//   - removeIf / erase_if / erase: RemoveIf(pred) implements the predicate-based
//     removal in Go and returns the number of removed bytes.
//
//   - push_back / push_front: covered by Append / Prepend.
//
//   - assign(InputIterator, InputIterator): use AssignBytes or
//     NewQByteArrayFromBytes.
//
//   - setNum(short/ushort/long/ulong) and number(short/ushort/long/ulong):
//     Qt delegates these overloads to the qlonglong/qulonglong versions. Use
//     SetNumLL/SetNumULL (or SetNumInt/SetNumUInt) and
//     NewQByteArrayNumberLL/ULL/Int/UInt.
//
//   - fromStdString / toStdString: covered by NewQByteArrayFromString / String.
//
//   - operator+ : Concat / ConcatBytes. operator+=: Append / Prepend.
//
//   - Comparison operators (==, !=, <, ...): Compare returns the Qt-style
//     three-way result.
//
//   - QByteArrayLiteral and the ""_ba / ""_qba literals: use
//     NewQByteArrayFromString.
//
//   - QByteArray(qsizetype, Qt::Initialization): use ResizeForOverwrite on a
//     fresh array.
//
//   - Platform-specific CFData/NSData (Darwin) and emscripten::val (WASM)
//     conversions: not portable, skipped.
//
//   - Deprecated no-argument count(): use Size or Len.
//
// # Design notes
//
//   - COW: Clone is O(1); the first mutation of a shared clone detaches (O(n)).
//
//   - Out-of-range indices panic (At/SetAt/First/Last/Sliced/SlicedN/Chopped/
//     Slice/SliceN/Remove/Insert/ReplaceRange). Qt only asserts, which is
//     undefined behaviour in release builds; a Go panic is safer.
//
//   - Left/Right/Mid keep Qt's clamping semantics and do not panic.
//
//   - Null vs empty: a null QByteArray has no data; an empty one has size 0 but
//     is not null. NewQByteArrayFromBytes(nil) and AssignBytes(nil) produce a
//     null array; non-nil empty slices produce an empty, non-null array.
//
//   - Thread safety: QByteArray is not safe for concurrent mutation, same as Qt.
//     Sharing a clone across goroutines is safe (COW).
//
// # Requirements
//
// Qt >= 6.11, C++17, pkg-config Qt6Core. The version check is performed at
// compile time in qbytearray_c.cpp via internal/qt6core/check_qt_version.h.
//
// This package is a subpackage of github.com/vadimlxx/go-qt6core.
package qtext
