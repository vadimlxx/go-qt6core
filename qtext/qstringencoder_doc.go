// Package qstringencoder provides a Go binding for Qt 6.11.2 QtCore's QStringEncoder.
//
// QStringEncoder is a concrete QStringConverter that encodes UTF-16 (QString / QStringView)
// into a byte sequence using a selected encoding (UTF-8, UTF-16LE/BE, Latin-1, System, …).
//
// In Go, QStringEncoder embeds *qstringconverter.QStringConverter, which itself embeds
// *qstringconverter.QStringConverterBase — reproducing Qt's class hierarchy:
//
//	QStringEncoder → QStringConverter → QStringConverterBase
//
// Every public method of QStringConverter (IsValid, Name, ResetState, HasError, State,
// SetState, Iface, Swap, Encoding) and QStringConverterBase (Flags, …) is available on
// *QStringEncoder via Go method promotion, without indirection.
//
// # Type mapping
//
//	Qt                                       Go
//	QStringEncoder                           *qstringencoder.QStringEncoder
//	QStringEncoder::DecodedData<T>           (*QStringEncoder).Encode / EncodeView → *qbytearray.QByteArray
//	QStringEncoder::FinalizeResult           (next int, invalid int16, err FinalizeResultError)
//	QStringEncoder::FinalizeResultError      qstringencoder.FinalizeResultError
//
// # Interop with sibling packages
//
// qstringencoder imports:
//   - qstringconverter: to embed *QStringConverter / *QStringConverterBase.
//   - qstring:          for Encode(*qstring.QString).
//   - qstringview:      for EncodeView / AppendToBuffer (QStringView).
//   - qbytearray:       for the result of Encode / EncodeView.
//   - qanystringview:   for NewQStringEncoderFromAny.
//
// # Protected members exposed
//
// The following protected Qt members are exposed (as requested):
//
//	protected QStringEncoder(const Interface *i) noexcept  → NewQStringEncoderFromInterface
//	protected ~QStringConverter()                          → (*QStringEncoder).Close / finalizer
//
// All other "protected" members of QStringEncoder (default / Encoding+Flags / QAnyStringView
// ctors) are public in Qt 6.11.2 and exposed as NewQStringEncoder / FromEncoding / FromName /
// FromAny.
//
// # Encoding result caching
//
// In Qt, `encode(...)` returns DecodedData<T>, which materializes lazily when implicitly
// converted to QByteArray. In Go, Encode / EncodeView return *qbytearray.QByteArray directly
// — semantically identical (both do one C++ pass to encode the whole input using the
// encoder's current state and then truncate to the actual output size).
//
// # TODO (not yet bound, or API differences vs Qt 6.11.2)
//
//   - QStringEncoder::operator()  Go has no operator(); use Encode / EncodeView.
//
//   - QStringEncoder::encode(const QString&) vs encode(QStringView):
//                                 Go has one method per input type; no overload set.
//
//   - DecodedData<T> as a public type: not exposed; Encode/EncodeView return the materialized
//                                 QByteArray directly.
//
//   - QStringBuilder integration: QConcatenable<QStringEncoder::DecodedData<T>> and
//                                 operator+=(QByteArray&, const DecodedData<T>&) are not
//                                 exposed. Use the equivalent:
//                                     ba.AppendQByteArray(enc.EncodeView(sv))
//                                 or:
//                                     qstringbuilder.FromQByteArray(enc.EncodeView(sv))
//
//   - QString::assign_helper_char8 (C++ helper for QString::assign with UTF-8 iterators):
//                                 not related to QStringEncoder's public API; not exposed.
//
//   - FinalizeResultChar<Char> template: monomorphized only for char in Qt's public API;
//                                 exposed as the Go tuple
//                                 (next int, invalid int16, err FinalizeResultError).
//
//   - QStringEncoder::availableCodecs() (inherited static): requires package qstringlist;
//                                 TODO in the qstringconverter package, not here.
//
//   - encodingForData / encodingForHtml (inherited static): available via
//                                 qstringconverter.EncodingForData / EncodingForHtml.
//
// # Requirements
//
// Qt >= 6.11.2, C++17, pkg-config Qt6Core. Version check is performed at compile time in
// qstringencoder_c.cpp via internal/qt6core/check_qt_version.h.
//
// This package is a subpackage of github.com/vadimlxx/go-qt6core.
package qtext
