// Package qstringencoder provides a Go binding for Qt 6.11 QtCore's
// QStringEncoder.
//
// QStringEncoder is a concrete QStringConverter that encodes UTF-16 (QString /
// QStringView) into a byte sequence using a selected encoding (UTF-8,
// UTF-16LE/BE, Latin-1, System, ...).
//
// In Go, QStringEncoder embeds *QStringConverter, which itself embeds
// *QStringConverterBase, reproducing Qt's class hierarchy:
//
//	QStringEncoder -> QStringConverter -> QStringConverterBase
//
// Every public method of the base chain (IsValid, Name, ResetState, HasError,
// State, SetState, Iface, Swap, Encoding) is available on *QStringEncoder via
// Go method promotion.
//
// # Type mapping
//
//	Qt                                  Go
//	QStringEncoder                      *qtext.QStringEncoder
//	QStringEncoder::DecodedData<T>      Encode / EncodeView -> *qtext.QByteArray
//	QStringEncoder::FinalizeResult      (next int, invalid int16, err FinalizeResultError)
//	QString                             *qtext.QString / string
//	QStringView                         qtext.QStringView / string
//	QAnyStringView                      qtext.QAnyStringView / string
//	QByteArray                          *qtext.QByteArray / []byte
//	qsizetype                           int
//
// # Ownership and lifecycle
//
// A *QStringEncoder owns its C++ object. The object is released by Close()
// (eager) or by runtime.AddCleanup when the wrapper becomes unreachable; both
// paths race safely through an atomic state slot, so there is no double free.
//
// The embedded *QStringConverter is a non-owning view over the same C++ object
// and the same state slot: its Close is a no-op, and after the encoder is
// closed the converter methods observe a nil handle instead of a dangling
// pointer.
//
// # Go-native variants
//
// In addition to the Qt-flavoured methods, the binding provides twins that use
// plain Go types:
//
//	EncodeString(string) *QByteArray       // QString input
//	EncodeBytes(*QString) []byte           // QByteArray result
//	EncodeViewString(string) *QByteArray   // QStringView input
//	EncodeViewBytes(QStringView) []byte    // QByteArray result
//	AppendToBufferString([]byte, string) int
//
// # API summary
//
// Constructors:
//
//	NewQStringEncoder, NewQStringEncoderFromEncoding,
//	NewQStringEncoderFromName, NewQStringEncoderFromAny,
//	NewQStringEncoderFromInterface
//
// Encoding:
//
//	Encode, EncodeString, EncodeBytes,
//	EncodeView, EncodeViewString, EncodeViewBytes,
//	RequiredSpace, AppendToBuffer, AppendToBufferString
//
// Finalization:
//
//	FinalizeBuffer, Finalize
//
// Lifecycle and interop:
//
//	Close, UnsafePtr (qstringconverter_interop.go)
//
// # Differences from Qt 6.11
//
//   - QStringEncoder::operator() has no Go equivalent: use Encode / EncodeView.
//
//   - DecodedData<T> is not exposed: Encode / EncodeView materialize the
//     QByteArray immediately (one C++ pass with the current State, then
//     truncated to the actual output size).
//
//   - encode(const QString &) and encode(QStringView) become separate Go
//     methods; Go has no overload sets.
//
//   - QStringBuilder integration (QConcatenable<DecodedData<T>>,
//     operator+=(QByteArray &, ...)) is not exposed. Use AppendQByteArray or
//     Concat on the materialized result.
//
//   - FinalizeResultChar<Char> is a template; the public char specialization is
//     returned as the three-value tuple (next, invalid, err).
//
//   - availableCodecs(), encodingForData() and encodingForHtml() are inherited
//     from QStringConverter and are available as qtext.AvailableCodecsStrings /
//     AvailableCodecs, qtext.EncodingForData / EncodingForDataBytes, and
//     qtext.EncodingForHtml / EncodingForHtmlBytes.
//
// # C bridge
//
// QStringConverter, QStringEncoder and QStringDecoder share one C ABI
// (qstringconverter_c.h) and one implementation (qstringconverter_c.cpp),
// mirroring QtCore's qstringconverter.h / qstringconverter.cpp. The opaque
// handles are real subclasses of helper adapters, and all C-type-to-Qt-type
// conversions are centralized in the cvt namespace, so the bridge contains no
// handle casts.
//
// # Requirements
//
// Qt >= 6.11, C++17, pkg-config Qt6Core. The version check is performed at
// compile time in qstringconverter_c.cpp via internal/qt6core/check_qt_version.h.
//
// This package is a subpackage of github.com/vadimlxx/go-qt6core.
package qtext
