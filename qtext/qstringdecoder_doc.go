// Package qstringdecoder provides a Go binding for Qt 6.11 QtCore's
// QStringDecoder.
//
// QStringDecoder is a concrete QStringConverter that decodes a byte sequence
// (using a selected encoding: UTF-8, UTF-16LE/BE, Latin-1, System, ...) into
// UTF-16 (QString / QChar).
//
// In Go, QStringDecoder embeds *QStringConverter, which itself embeds
// *QStringConverterBase, reproducing Qt's class hierarchy:
//
//	QStringDecoder -> QStringConverter -> QStringConverterBase
//
// Every public method of the base chain (IsValid, Name, ResetState, HasError,
// State, SetState, Iface, Swap, Encoding) is available on *QStringDecoder via
// Go method promotion.
//
// # Type mapping
//
//	Qt                                  Go
//	QStringDecoder                      *qtext.QStringDecoder
//	QStringDecoder::EncodedData<T>      Decode / DecodeView -> *qtext.QString
//	QStringDecoder::FinalizeResult      (next int, invalid int16, err FinalizeResultError)
//	QStringDecoder::FinalizeResultQChar same tuple (QChar == uint16)
//	QByteArray                          *qtext.QByteArray / []byte
//	QByteArrayView                      qtext.QByteArrayView / []byte
//	QString                             *qtext.QString / string
//	QAnyStringView                      qtext.QAnyStringView / string
//	qsizetype                           int
//
// # Ownership and lifecycle
//
// A *QStringDecoder owns its C++ object. The object is released by Close()
// (eager) or by runtime.AddCleanup when the wrapper becomes unreachable; both
// paths race safely through an atomic state slot, so there is no double free.
//
// The embedded *QStringConverter is a non-owning view over the same C++ object
// and the same state slot: its Close is a no-op, and after the decoder is
// closed the converter methods observe a nil handle instead of a dangling
// pointer.
//
// # Go-native variants
//
// In addition to the Qt-flavoured methods, the binding provides twins that use
// plain Go types:
//
//	DecodeBytes([]byte) *QString           // QByteArray input
//	DecodeString(*QByteArray) string       // QString result
//	DecodeViewBytes([]byte) *QString       // QByteArrayView input
//	DecodeViewString(QByteArrayView) string
//	AppendToBufferBytes([]uint16, []byte) int
//	NewQStringDecoderForHtmlBytes([]byte) *QStringDecoder
//
// # API summary
//
// Constructors:
//
//	NewQStringDecoder, NewQStringDecoderFromEncoding,
//	NewQStringDecoderFromName, NewQStringDecoderFromAny,
//	NewQStringDecoderFromInterface, NewQStringDecoderForHtml,
//	NewQStringDecoderForHtmlBytes
//
// Decoding:
//
//	Decode, DecodeBytes, DecodeString,
//	DecodeView, DecodeViewBytes, DecodeViewString,
//	RequiredSpace, AppendToBuffer, AppendToBufferBytes
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
//   - QStringDecoder::operator() has no Go equivalent: use Decode / DecodeView.
//
//   - EncodedData<T> is not exposed: Decode / DecodeView materialize the QString
//     immediately (one C++ pass with the current State, then truncated to the
//     actual output size).
//
//   - decode(const QByteArray &) and decode(QByteArrayView) become separate Go
//     methods; Go has no overload sets.
//
//   - QStringBuilder integration (QConcatenable<EncodedData<T>>,
//     operator+=(QString &, ...)) is not exposed. Use AppendQString or Concat on
//     the materialized result.
//
//   - FinalizeResultChar<Char> is a template; char16_t and QChar overloads are
//     identical in Go (QChar == uint16), and the result is the three-value tuple
//     (next, invalid, err).
//
//   - decoderForHtml(QByteArrayView) is exposed as NewQStringDecoderForHtml
//     (and the []byte form NewQStringDecoderForHtmlBytes); DecoderForHtml is a
//     deprecated alias.
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
