// Package qstringconverter provides a Go binding for Qt 6.11 QtCore's QStringConverter and
// QStringConverterBase.
//
// In Qt 6.11, QStringConverterBase is almost empty: it declares Flag, Flags, State and
// protected ctors/dtor. All functional API (Encoding, Interface, protected iface/state,
// isValid/name/resetState/hasError, static helpers) lives on QStringConverter, which inherits
// from QStringConverterBase. In Go the class hierarchy is reproduced via struct embedding:
// QStringConverter embeds *QStringConverterBase.
//
// # Type mapping
//
//	Qt                                       Go
//	QStringConverterBase::Flag               qstringconverter.Flag
//	QStringConverterBase::Flags              qstringconverter.Flags
//	QStringConverterBase::State              qstringconverter.State
//	QStringConverter::Encoding               qstringconverter.Encoding
//	QStringConverter::Interface              qstringconverter.Interface
//	QStringConverterBase                     *qstringconverter.QStringConverterBase
//	QStringConverter                         *qstringconverter.QStringConverter
//
// # Interop with sibling packages
//
// qstringconverter imports qanystringview (for the QStringConverter(QAnyStringView, Flags)
// ctor and encodingForName) and qbytearrayview (for encodingForData/encodingForHtml).
// No import cycles: qanystringview and qbytearrayview do not import qstringconverter.
//
// # Protected members exposed
//
// The binding makes protected members accessible, matching your request:
//
//	protected QStringConverter()                     → NewQStringConverter
//	protected QStringConverter(Encoding, Flags)      → NewQStringConverterFromEncoding
//	protected QStringConverter(QAnyStringView, Flags)→ NewQStringConverterFromAny / FromName
//	protected QStringConverter(const Interface*)     → NewQStringConverterFromIface
//	protected iface                                  → (*QStringConverter).Iface
//	protected state                                  → (*QStringConverter).State / SetState
//
// # TODO (not yet bound, or API differences vs Qt 6.11)
//
//   - QStringConverter::availableCodecs() — returns QStringList; needs package qstringlist.
//
//   - QStringConverter::encodingForData / encodingForHtml — implemented (take
//     qbytearrayview.QByteArrayView). Optional argument expectedFirstCharacter is exposed
//     as uint16 (0 = not given).
//
//   - QStringConverter::FinalizeResultError and FinalizeResultChar<Char> — these template
//     helpers are used by QStringDecoder/QStringEncoder (see next item) and are not exposed
//     directly.
//
//   - QStringDecoder / QStringEncoder — derived classes of QStringConverter. Would be
//     separate subpackages (qstringdecoder, qstringencoder) that consume the `Interface`
//     function pointers.
//
//   - Invoking Interface::{toUtf16, toUtf16Len, fromUtf16, fromUtf16Len} from Go — requires
//     cgo function-pointer support; the pointers are exposed as unsafe.Pointer, but calling
//     them is not implemented.
//
//   - QStringConverter::State::clearFn (function pointer) — not exposed. State is move-only
//     in C++; SetState uses move-assignment (which correctly invokes clearFn on the old state
//     if the encoding requires it).
//
//   - QStringConverter::State::d[2] (union alternative to state_data[4]) — same bytes as
//     State.Data; not a separate field.
//
//   - QStringConverter::encoding() — does not exist in Qt 6.11. Our (*QStringConverter).Encoding()
//     derives it from Name() via EncodingForName.
//
//   - QStringConverter::swap() — does not exist in Qt 6.11. Our (*QStringConverter).Swap() is
//     implemented via std::swap on protected iface/state.
//
//   - Copy constructor of QStringConverter — deleted in Qt 6.11 (move-only). Recreate via
//     FromEncoding/FromName; the partial-decode state is lost.
//
//   - Q_DECLARE_OPERATORS_FOR_FLAGS(QStringConverter::Flags) (operator|, operator&) — Go has
//     no operator overloading; use Flags.Has(flag), or manually OR Flag values then cast to Flags.
//
// # Requirements
//
// Qt >= 6.11, C++17, pkg-config Qt6Core. Version check is performed at compile time in
// qstringconverter_c.cpp via internal/qt6core/check_qt_version.h.
//
// This package is a subpackage of github.com/vadimlxx/go-qt6core.
package qstringconverter
