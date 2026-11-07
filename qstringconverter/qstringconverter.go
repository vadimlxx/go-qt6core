package qstringconverter

// Package qstringconverter: QStringConverter — конкретный класс.

/*
 QStringConverter наследует QStringConverterBase (через embedding, как в Qt 6.11).
 В Qt 6.11 QStringConverter содержит:
   - protected: Interface (nested struct), iface (const Interface*), state (State),
     конструкторы (default, Encoding+Flags, QAnyStringView+Flags, const Interface*),
     деструктор.
   - public: move ctor/assign, isValid(), resetState(), hasError(), name(),
     статические nameForEncoding(Encoding), encodingForName(QAnyStringView),
     encodingForData(QByteArrayView, char16_t), encodingForHtml(QByteArrayView),
     availableCodecs().

 не реализовано:
   - availableCodecs() — возвращает QStringList; нужен пакет qstringlist.
   - FinalizeResultError / FinalizeResultChar<Char> — не экспонируются; понадобятся для QStringDecoder / QStringEncoder.
   - QStringDecoder / QStringEncoder — наследники QStringConverter; отдельные пакеты.
   - Вызов функциональных указателей Interface (toUtf16/toUtf16Len/fromUtf16/fromUtf16Len)
     — требует cgo function-pointer support; сейчас они отдаются как unsafe.Pointer.
   - encoding() как метод QStringConverter — отсутствует в Qt 6.11; наш Encoding()
     выводит его из name() через encodingForName().
   - swap() как метод QStringConverter — отсутствует в Qt 6.11; реализован через  std::swap(iface, state) в C++.
   - Копирование QStringConverter — в Qt 6.11 копирующий конструктор удалён (только move);
     пересоздавайте через FromEncoding/FromName (состояние теряется).
*/

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC -Wno-sfinae-incomplete
#cgo LDFLAGS: -lstdc++
#include <stdlib.h>
#include <stdint.h>
#include <stdbool.h>
#include <stddef.h>
#include "qstringconverter_c.h"
*/
import "C"

import (
	"runtime"
	"unsafe"

	"github.com/vadimlxx/go-qt6core/qanystringview"
	"github.com/vadimlxx/go-qt6core/qbytearrayview"
)

// Encoding mirrors QStringConverter::Encoding.
// System = 8 in Qt 6.11 (Latin1 = 7, System = 8).
type Encoding int

const (
	Utf8         Encoding = 0
	Utf16        Encoding = 1
	Utf16LE      Encoding = 2
	Utf16BE      Encoding = 3
	Utf32        Encoding = 4
	Utf32LE      Encoding = 5
	Utf32BE      Encoding = 6
	Latin1       Encoding = 7
	System       Encoding = 8
	LastEncoding          = System
)

// Interface wraps a const QStringConverter::Interface* (protected nested struct).
//
// Interface is protected in Qt 6.11; obtain one via (*QStringConverter).Iface().
type Interface struct {
	p unsafe.Pointer
}

// IsNil reports whether i wraps a null pointer.
func (i Interface) IsNil() bool { return i.p == nil }

// Name returns the encoding name ("UTF-16", "ISO-8859-1", etc.), or "" for a null Interface.
func (i Interface) Name() string {
	if i.p == nil {
		return ""
	}
	cname := C.qstringconverter_interface_name(i.p)
	if cname == nil {
		return ""
	}
	return C.GoString(cname)
}

// ToUtf16Ptr returns the raw decoder function pointer (QChar* (*)(QChar*, QByteArrayView, State*)),
// or nil. TODO: invoking it from Go requires cgo function-pointer support (not implemented).
func (i Interface) ToUtf16Ptr() unsafe.Pointer {
	if i.p == nil {
		return nil
	}
	return C.qstringconverter_interface_to_utf16(i.p)
}

// ToUtf16LenPtr returns the raw decoder length function pointer (qsizetype (*)(qsizetype)),
// or nil.
func (i Interface) ToUtf16LenPtr() unsafe.Pointer {
	if i.p == nil {
		return nil
	}
	return C.qstringconverter_interface_to_utf16_len(i.p)
}

// FromUtf16Ptr returns the raw encoder function pointer (char* (*)(char*, QStringView, State*)),
// or nil.
func (i Interface) FromUtf16Ptr() unsafe.Pointer {
	if i.p == nil {
		return nil
	}
	return C.qstringconverter_interface_from_utf16(i.p)
}

// FromUtf16LenPtr returns the raw encoder length function pointer (qsizetype (*)(qsizetype)),
// or nil.
func (i Interface) FromUtf16LenPtr() unsafe.Pointer {
	if i.p == nil {
		return nil
	}
	return C.qstringconverter_interface_from_utf16_len(i.p)
}

// QStringConverter mirrors Qt 6.11's QStringConverter.
//
// It embeds *QStringConverterBase (matching Qt's class hierarchy). The embedded base is a
// non-owning view over the same C++ object; the *QStringConverter wrapper owns it.
type QStringConverter struct {
	*QStringConverterBase
	p *C.go_qsc_t
}

func newQStringConverterFromC(c *C.go_qsc_t) *QStringConverter {
	if c == nil {
		return nil
	}
	base := &QStringConverterBase{
		p:         (*C.go_qscbase_t)(unsafe.Pointer(c)),
		ownedByGo: false,
	}
	qsc := &QStringConverter{
		QStringConverterBase: base,
		p:                     c,
	}
	runtime.SetFinalizer(qsc, finalizeQStringConverter)
	return qsc
}

func finalizeQStringConverter(qsc *QStringConverter) {
	if qsc.p != nil {
		C.qstringconverter_delete(qsc.p)
		qsc.p = nil
		if qsc.QStringConverterBase != nil {
			qsc.QStringConverterBase.p = nil
		}
	}
}

// NewQStringConverter creates an invalid converter (null iface). Mirrors QStringConverter().
func NewQStringConverter() *QStringConverter {
	return newQStringConverterFromC(C.qstringconverter_new_default())
}

// NewQStringConverterFromEncoding creates a converter for the given encoding and flags.
// Mirrors QStringConverter(Encoding e, Flags f = Flag::Default).
func NewQStringConverterFromEncoding(e Encoding, f Flags) *QStringConverter {
	return newQStringConverterFromC(C.qstringconverter_new_from_encoding(C.int(int(e)), C.int(int(f))))
}

// NewQStringConverterFromName creates a converter by encoding name (e.g. "UTF-16",
// "ISO-8859-1"). Mirrors QStringConverter(QAnyStringView name, Flags f).
// Returns an invalid converter if the name is unknown.
func NewQStringConverterFromName(name string, f Flags) *QStringConverter {
	var p unsafe.Pointer
	if len(name) > 0 {
		p = unsafe.Pointer(unsafe.StringData(name))
	}
	qsc := newQStringConverterFromC(C.qstringconverter_new_from_any(
		p, C.ptrdiff_t(len(name)), 0, C.int(int(f))))
	runtime.KeepAlive(name)
	return qsc
}

// NewQStringConverterFromAny creates a converter from a QAnyStringView name and flags.
// Mirrors QStringConverter(QAnyStringView name, Flags f).
func NewQStringConverterFromAny(name qanystringview.QAnyStringView, f Flags) *QStringConverter {
	qsc := newQStringConverterFromC(C.qstringconverter_new_from_any(
		name.DataUnsafe(), C.ptrdiff_t(name.Size()), C.int(int(name.Kind())), C.int(int(f))))
	runtime.KeepAlive(name)
	return qsc
}

// NewQStringConverterFromIface creates a converter from an existing Interface pointer
// (protected ctor exposed). If i is null, the converter is invalid.
func NewQStringConverterFromIface(i Interface) *QStringConverter {
	return newQStringConverterFromC(C.qstringconverter_new_from_iface(i.p))
}

// Close frees the underlying C++ object immediately and disables the finalizer.
func (q *QStringConverter) Close() {
	if q.p != nil {
		C.qstringconverter_delete(q.p)
		q.p = nil
		if q.QStringConverterBase != nil {
			q.QStringConverterBase.p = nil
		}
		runtime.SetFinalizer(q, nil)
	}
}

// Swap exchanges the contents of q and other.
// Qt 6.11 has no QStringConverter::swap; implemented via std::swap on protected iface/state.
func (q *QStringConverter) Swap(other *QStringConverter) {
	if q.p == nil || other == nil || other.p == nil {
		return
	}
	C.qstringconverter_swap(q.p, other.p)
}

// Encoding returns the converter's encoding, derived from Name() via EncodingForName().
// Qt 6.11 has no public QStringConverter::encoding() method.
// Returns Encoding(-1) if the converter is invalid.
func (q *QStringConverter) Encoding() Encoding {
	if q.p == nil {
		return Encoding(-1)
	}
	name := q.Name()
	if name == "" {
		return Encoding(-1)
	}
	e, ok := EncodingForName(name)
	if !ok {
		return Encoding(-1)
	}
	return e
}

// IsValid mirrors QStringConverter::isValid().
func (q *QStringConverter) IsValid() bool {
	if q.p == nil {
		return false
	}
	return bool(C.qstringconverter_is_valid(q.p))
}

// HasError mirrors QStringConverter::hasError().
func (q *QStringConverter) HasError() bool {
	if q.p == nil {
		return true
	}
	return bool(C.qstringconverter_has_error(q.p))
}

// Name mirrors QStringConverter::name().
func (q *QStringConverter) Name() string {
	if q.p == nil {
		return ""
	}
	cname := C.qstringconverter_name(q.p)
	if cname == nil {
		return ""
	}
	return C.GoString(cname)
}

// ResetState mirrors QStringConverter::resetState().
func (q *QStringConverter) ResetState() {
	if q.p == nil {
		return
	}
	C.qstringconverter_reset_state(q.p)
}

// State returns a copy of the converter's internal State (protected `state` field).
func (q *QStringConverter) State() State {
	var s State
	if q.p == nil {
		return s
	}
	var flags C.int
	var internal C.int
	var remaining C.int64_t
	var invalid C.int64_t
	var data [4]C.uint32_t
	C.qstringconverter_get_state(q.p, &flags, &internal, &remaining, &invalid, &data[0])
	s.Flags = Flags(flags)
	s.InternalState = int32(internal)
	s.RemainingChars = int64(remaining)
	s.InvalidChars = int64(invalid)
	for i := 0; i < 4; i++ {
		s.Data[i] = uint32(data[i])
	}
	return s
}

// SetState overwrites the converter's internal State (protected `state` field).
// Uses move-assignment in C++ (Qt State is move-only), which correctly resets any
// associated state via clearFn if the current encoding uses one.
func (q *QStringConverter) SetState(s State) {
	if q.p == nil {
		return
	}
	var data [4]C.uint32_t
	for i := 0; i < 4; i++ {
		data[i] = C.uint32_t(s.Data[i])
	}
	C.qstringconverter_set_state(q.p,
		C.int(int(s.Flags)), C.int(s.InternalState),
		C.int64_t(s.RemainingChars), C.int64_t(s.InvalidChars), &data[0])
}

// Iface returns the Interface pointer of the converter (protected `iface` field),
// or a null Interface.
func (q *QStringConverter) Iface() Interface {
	if q.p == nil {
		return Interface{}
	}
	return Interface{p: C.qstringconverter_iface(q.p)}
}

// NameForEncoding returns the canonical name for e, or "" if unknown
// (mirrors QStringConverter::nameForEncoding).
func NameForEncoding(e Encoding) string {
	cname := C.qstringconverter_name_for_encoding(C.int(int(e)))
	if cname == nil {
		return ""
	}
	return C.GoString(cname)
}

// EncodingForName returns the Encoding for the given name, and true if recognized;
// otherwise (Utf8, false). Mirrors static QStringConverter::encodingForName
// (which returns std::optional<Encoding> in Qt 6.11).
func EncodingForName(name string) (Encoding, bool) {
	var p unsafe.Pointer
	if len(name) > 0 {
		p = unsafe.Pointer(unsafe.StringData(name))
	}
	v := int(C.qstringconverter_encoding_for_name(p, C.ptrdiff_t(len(name)), 0))
	runtime.KeepAlive(name)
	if v < 0 {
		return Utf8, false
	}
	return Encoding(v), true
}

// EncodingForNameAny is EncodingForName accepting a qanystringview.QAnyStringView directly.
func EncodingForNameAny(name qanystringview.QAnyStringView) (Encoding, bool) {
	v := int(C.qstringconverter_encoding_for_name(
		name.DataUnsafe(), C.ptrdiff_t(name.Size()), C.int(int(name.Kind()))))
	runtime.KeepAlive(name)
	if v < 0 {
		return Utf8, false
	}
	return Encoding(v), true
}

// EncodingForData returns the Encoding detected from the given byte data, and true if
// recognized; otherwise (Utf8, false). If expectedFirstChar != 0, the detected encoding
// is accepted only when the decoded data starts with that code unit.
// Mirrors static QStringConverter::encodingForData.
func EncodingForData(data qbytearrayview.QByteArrayView, expectedFirstChar uint16) (Encoding, bool) {
	bs := data.Bytes()
	var p unsafe.Pointer
	if len(bs) > 0 {
		p = unsafe.Pointer(&bs[0])
	}
	v := int(C.qstringconverter_encoding_for_data(
		p, C.ptrdiff_t(len(bs)), C.uint16_t(expectedFirstChar)))
	runtime.KeepAlive(bs)
	if v < 0 {
		return Utf8, false
	}
	return Encoding(v), true
}

// EncodingForHtml returns the Encoding detected from the HTML <meta> charset declaration
// in the given byte data, and true if recognized; otherwise (Utf8, false).
// Mirrors static QStringConverter::encodingForHtml.
func EncodingForHtml(data qbytearrayview.QByteArrayView) (Encoding, bool) {
	bs := data.Bytes()
	var p unsafe.Pointer
	if len(bs) > 0 {
		p = unsafe.Pointer(&bs[0])
	}
	v := int(C.qstringconverter_encoding_for_html(p, C.ptrdiff_t(len(bs))))
	runtime.KeepAlive(bs)
	if v < 0 {
		return Utf8, false
	}
	return Encoding(v), true
}

// TODO: AvailableCodecs() []string — needs qstringlist package for QStringList result.
