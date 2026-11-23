// Package qstringconverter provides a Go binding for Qt 6.11 QtCore's
// QStringConverter and QStringConverterBase.
//
// In Qt 6.11, QStringConverterBase is almost empty: it declares Flag, Flags,
// State and protected ctors/dtor. All functional API (Encoding, Interface,
// protected iface/state, isValid/name/resetState/hasError, static helpers) lives
// on QStringConverter, which inherits from QStringConverterBase. In Go the class
// hierarchy is reproduced via struct embedding: QStringConverter embeds
// *QStringConverterBase.
//
// # Type mapping
//
//	Qt                                       Go
//	QStringConverterBase::Flag               qtext.Flag
//	QStringConverterBase::Flags              qtext.Flags
//	QStringConverterBase::State              qtext.State
//	QStringConverter::Encoding               qtext.Encoding
//	QStringConverter::Interface              qtext.Interface
//	QStringConverter::FinalizeResultError    qtext.FinalizeResultError
//	QStringConverterBase                     *qtext.QStringConverterBase
//	QStringConverter                         *qtext.QStringConverter
//	QAnyStringView                           qtext.QAnyStringView / string
//	QByteArrayView                           qtext.QByteArrayView / []byte
//	QStringList                              *qtext.QStringList / []string
//
// # Ownership and lifecycle
//
// A *QStringConverter owns its C++ object. The object is released either by
// Close() (eager, deterministic) or by runtime.AddCleanup when the Go wrapper
// becomes unreachable. Close and the cleanup race safely through an atomic
// state slot: exactly one of them claims the object, so there is no double
// free.
//
// The embedded *QStringConverterBase is a non-owning view over the same C++
// object (and the same state slot). Its Close is a no-op; after the converter
// is closed, methods reached through the embedded base observe a nil handle
// instead of a dangling pointer.
//
// A *QStringConverterBase created directly by NewQStringConverterBase owns its
// own C++ object and follows the same Close/cleanup rules.
//
// NewQStringConverterFromUnsafePtr adopts an existing C++ QStringConverter
// without taking ownership: no cleanup is registered and Close is a no-op. It
// exists so that QStringEncoder and QStringDecoder can embed a converter view
// while owning a derived C++ object themselves.
//
// # Protected members exposed
//
//	protected QStringConverter()                      -> NewQStringConverter
//	protected QStringConverter(Encoding, Flags)       -> NewQStringConverterFromEncoding
//	protected QStringConverter(QAnyStringView, Flags) -> NewQStringConverterFromAny (view)
//	                                                     NewQStringConverterFromName (Go string)
//	protected QStringConverter(const Interface *)     -> NewQStringConverterFromIface
//	protected iface                                   -> (*QStringConverter).Iface
//	protected state                                   -> (*QStringConverter).State / SetState
//
// # Go-native variants
//
// Where Qt uses a non-trivial string or container type, the binding provides a
// Go-native twin in addition to the Qt-flavoured form:
//
//	string          NewQStringConverterFromName, EncodingForName
//	QAnyStringView  NewQStringConverterFromAny, EncodingForNameAny
//	[]byte          EncodingForDataBytes, EncodingForHtmlBytes
//	QByteArrayView  EncodingForData, EncodingForHtml
//	QStringList     AvailableCodecs
//	[]string        AvailableCodecsStrings
//
// # API summary
//
// Constructors:
//
//	NewQStringConverter, NewQStringConverterFromEncoding,
//	NewQStringConverterFromName, NewQStringConverterFromAny,
//	NewQStringConverterFromIface, NewQStringConverterFromUnsafePtr,
//	NewQStringConverterBase
//
// Lifecycle:
//
//	Close, Swap, UnsafePtr (see qstringconverter_interop.go)
//
// State and validity:
//
//	IsValid, HasError, Name, Encoding, ResetState, State, SetState, Iface
//
// Static helpers:
//
//	NameForEncoding, EncodingForName, EncodingForNameAny,
//	EncodingForData, EncodingForDataBytes, EncodingForHtml,
//	EncodingForHtmlBytes, AvailableCodecs, AvailableCodecsStrings
//
// Interface (protected nested type):
//
//	IsNil, UnsafePtr, Name, ToUtf16Ptr, ToUtf16LenPtr, FromUtf16Ptr,
//	FromUtf16LenPtr
//
// # C ABI and conversions
//
// The C bridge is deliberately cast-light. QStringConverter, QStringEncoder and
// QStringDecoder share qstringconverter_c.cpp/.h (mirroring QtCore's
// qstringconverter.cpp/.h), while QStringConverterBase has its own
// qstringconverterbase_c.cpp/.h pair, mirroring qstringconverter_base.h:
//
//   - go_qsc_t and go_qscbase_t are real subclasses of helper adapters over the
//     Qt classes, so the ABI functions never cast between a handle and a
//     QStringConverter or QStringConverterBase pointer;
//   - every conversion between the C ABI types (int, ptrdiff_t, void *) and Qt
//     types is centralized in one cvt namespace in the implementation;
//   - outside cvt only the four Interface callback conversions
//     (function pointer -> void *) remain, and ISO C++ makes those
//     conditionally supported, so they cannot be removed portably.
//
// The header qstringconverter_c.h stays strictly C-compatible: it mentions no
// C++ type and no Qt type, which is what keeps cgo's preamble simple.
//
// # Interop
//
// The unsafe.Pointer bridge lives in qstringconverter_interop.go, mirroring
// interop.go for QString/QByteArray:
//
//	(*QStringConverter).UnsafePtr       borrowed C++ pointer, owned by the wrapper
//	NewQStringConverterFromUnsafePtr    non-owning wrapper (no cleanup, Close no-op)
//	QStringConverterAdoptUnsafePtr      deprecated alias of the above
//
// QStringEncoder and QStringDecoder use the internal newQStringConverterAlias
// helper so that their embedded converter view shares both the derived C++
// object they own and the same state slot (closing the owner invalidates the
// view safely).
//
// # Differences from Qt 6.11
//
//   - QStringConverter::encoding() does not exist in Qt 6.11. The binding's
//     Encoding() derives the value from Name() via EncodingForName and returns
//     Encoding(-1) for an invalid converter.
//
//   - QStringConverter::swap() does not exist in Qt 6.11. Swap() is implemented
//     with std::swap on the protected iface/state fields.
//
//   - The copy constructor is deleted in Qt 6.11 (the class is move-only).
//     Recreate a converter with NewQStringConverterFromEncoding or
//     NewQStringConverterFromName; partial conversion state is not carried over.
//
//   - FinalizeResultChar<Char> is a template used by QStringEncoder and
//     QStringDecoder and is not exposed. FinalizeResultError is exposed as an
//     enum.
//
//   - Calling Interface::{toUtf16, toUtf16Len, fromUtf16, fromUtf16Len} from Go
//     requires cgo function-pointer support. The pointers are exposed as
//     unsafe.Pointer, but calling them is not implemented.
//
//   - State::clearFn (function pointer) is not exposed. SetState clears the
//     previous state and copies the fields, so any clearFn of the current
//     encoding is still honoured.
//
//   - State::d[2] (union alternative to state_data[4]) is not a separate field:
//     State.Data holds the same bytes.
//
//   - Q_DECLARE_OPERATORS_FOR_FLAGS(QStringConverter::Flags) (operator|,
//     operator&) has no Go equivalent. Use Flags.Has(flag) and explicit casts.
//
//   - QStringConverterBase(QStringConverterBase&&) is not exposed: Go has no
//     move constructors.
//
// # Design notes
//
//   - Nil and closed converters are not fatal: methods such as IsValid, Name and
//     Encoding return their zero values, HasError returns true, and Close is a
//     no-op. This mirrors Qt's "invalid converter" handling.
//
//   - A default-constructed converter is invalid (null iface) but is a valid
//     object; Close must still be called or left to the cleanup.
//
//   - State is move-only in C++. Go copies it by value; SetState performs the
//     C++ clear-then-assign sequence, which is safe with respect to clearFn.
//
// # Requirements
//
// Qt >= 6.11, C++17, pkg-config Qt6Core. The version check is performed at
// compile time in qstringconverter_c.cpp via internal/qt6core/check_qt_version.h.
//
// This package is a subpackage of github.com/vadimlxx/go-qt6core.
package qtext
