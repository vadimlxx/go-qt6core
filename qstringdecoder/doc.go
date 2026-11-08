// Package qstringdecoder provides a Go binding for Qt 6.11.2 QtCore's QStringDecoder.
//
// QStringDecoder is a concrete QStringConverter that decodes a byte sequence (using a
// selected encoding — UTF-8, UTF-16LE/BE, Latin-1, System, …) into UTF-16
// (QString / QStringView / QChar).
//
// In Go, QStringDecoder embeds *qstringconverter.QStringConverter, which itself embeds
// *qstringconverter.QStringConverterBase — reproducing Qt's class hierarchy:
//
//	QStringDecoder → QStringConverter → QStringConverterBase
//
// Every public method of QStringConverter (IsValid, Name, ResetState, HasError, State,
// SetState, Iface, Swap, Encoding) and of QStringConverterBase (Flags, …) is available on
// *QStringDecoder via Go method promotion.
//
// # Type mapping
//
//	Qt                                    Go
//	QStringDecoder                        *qstringdecoder.QStringDecoder
//	QStringDecoder::EncodedData<T>        (*QStringDecoder).Decode / DecodeView → *qstring.QString
//	QStringDecoder::FinalizeResult        (next int, invalid int16, err FinalizeResultError)
//	QStringDecoder::FinalizeResultQChar   same tuple (QChar == uint16 in Go)
//	QStringDecoder::FinalizeResultError   qstringdecoder.FinalizeResultError
//
// # Interop with sibling packages
//
// qstringdecoder imports:
//   - qstringconverter: to embed *QStringConverter / *QStringConverterBase.
//   - qstring:          for the Decode / DecodeView result.
//   - qstringview:      available for users converting the result to a view.
//   - qbytearray:       for Decode(*qbytearray.QByteArray).
//   - qbytearrayview:   for DecodeView / AppendToBuffer / DecoderForHtml.
//   - qanystringview:   for NewQStringDecoderFromAny.
//
// No import cycle: qstringconverter imports none of the other packages, so depending on it
// is safe. The unsafe.Pointer bridge qstringconverter.AdoptUnsafePtr() lets this package
// share the same C++ QStringDecoder* between the owning *QStringDecoder wrapper and its
// embedded *QStringConverter alias.
//
// # Protected members exposed
//
// The following protected Qt members are exposed (as requested):
//
//	protected QStringDecoder(const Interface *i) noexcept  → NewQStringDecoderFromInterface
//	protected ~QStringConverter()                          → (*QStringDecoder).Close / finalizer
//
// All other ctors (default / Encoding+Flags / QAnyStringView+Flags) are public in Qt 6.11.2
// and exposed as NewQStringDecoder / FromEncoding / FromName / FromAny.
//
// # Decoding result caching
//
// In Qt, `decode(...)` returns EncodedData<T>, which materializes lazily when implicitly
// converted to QString. In Go, Decode / DecodeView return *qstring.QString directly —
// semantically identical (one C++ pass with the decoder's current State, then truncate to
// the actual output size). An invalid decoder (iface == nullptr) produces an empty QString
// and sets HasError, mirroring QStringDecoder::invalidateAndReturnNull().
//
// # TODO (not yet bound, or API differences vs Qt 6.11.2)
//
//   - QStringDecoder::operator() Go has no operator(); use Decode / DecodeView.
//
//   - QStringDecoder::decode(const QByteArray&) vs decode(QByteArrayView):
//                              Go has one method per input type; no overload set.
//
//   - EncodedData<T> as a public type: not exposed; Decode / DecodeView return the
//                              materialized QString directly.
//
//   - QStringBuilder integration: QConcatenable<QStringDecoder::EncodedData<T>> and
//                              operator+=(QString&, const EncodedData<T>&) are not exposed.
//                              Use the equivalent:
//                                  s.AppendQString(dec.DecodeView(bv))
//                              or:
//                                  qstringbuilder.Concat(s, dec.DecodeView(bv)).ToQString()
//
//   - QString::assign_helper_char8 (C++ helper using QStringDecoder with UTF-8 iterators):
//                              not related to QStringDecoder's public API; not exposed.
//
//   - FinalizeResultChar<Char> template: monomorphized only for char16_t / QChar in Qt's
//                              public API; exposed as the Go tuple
//                              (next int, invalid int16, err FinalizeResultError).
//
//   - QStringDecoder::decoderForHtml(QByteArrayView): exposed as the package-level
//                              function DecoderForHtml(qbytearrayview.QByteArrayView).
//
//   - QStringConverter::availableCodecs() (inherited static): requires qstringlist; TODO
//                              in the qstringconverter package.
//
//   - encodingForData / encodingForHtml (inherited static): available via
//                              qstringconverter.EncodingForData / EncodingForHtml.
//
// # Requirements
//
// Qt >= 6.11.2, C++17, pkg-config Qt6Core. Version check is performed at compile time in
// qstringdecoder_c.cpp via internal/qt6core/check_qt_version.h.
//
// This package is a subpackage of github.com/vadimlxx/go-qt6core.
package qstringdecoder
