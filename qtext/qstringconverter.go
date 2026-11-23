package qtext

/*
QStringConverter — конкретный класс Qt 6.11 QtCore/QStringConverter.

Иерархия воспроизведена встраиванием, как в Qt:

	QStringConverter -> QStringConverterBase

*QStringConverter владеет C++ объектом: он освобождается Close() или
runtime.AddCleanup (когда Go-обёртка становится недостижимой). Встроенная
*QStringConverterBase — не владеющая проекция на тот же объект; её Close()
ничего не делает. Оба используют общий qscState с атомарным слотом указателя,
поэтому Close() и cleanup безопасно соревнуются, а после закрытия методы видят
nil.

Варианты API:

  - функции, принимающие/возвращающие QAnyStringView, имеют строковую форму:
    NewQStringConverterFromName / EncodingForName (Go string) и
    NewQStringConverterFromAny / EncodingForNameAny (QAnyStringView);
  - функции, принимающие QByteArrayView, имеют форму с []byte:
    EncodingForData / EncodingForDataBytes, EncodingForHtml /
    EncodingForHtmlBytes;
  - availableCodecs() возвращает QStringList, поэтому есть две формы:
    AvailableCodecs (*QStringList) и AvailableCodecsStrings ([]string);
  - методы QByteArray/QByteArrayView в QStringConverter отсутствуют, поэтому
    других []byte-вариантов нет.

# C-часть и приведения типов

C++-мост (qstringconverter_c.cpp) сведён к минимуму приведений:

  - go_qsc_t и go_qscbase_t — реальные наследники helper-классов, поэтому
    приведений между handle и QStringConverter* нет вообще;
  - все преобразования int/ptrdiff_t/void* в типы Qt собраны в одном
    namespace cvt (encoding, flags, size, chars, utf16, iface, ...);
  - вне cvt остаются только четыре неизбежных приведения «указатель на
    функцию -> void*» для callbacks Interface (в C++ это условно
    поддерживаемая операция).

ABI при этом не менялся: заголовок qstringconverter_c.h по-прежнему говорит
только в терминах C-типов.

Интероп-функции для обмена «сырыми» указателями (UnsafePtr,
NewQStringConverterFromUnsafePtr и устаревший алиас
QStringConverterAdoptUnsafePtr) вынесены в qstringconverter_interop.go.

Не реализовано:

  - FinalizeResultChar<Char> — шаблон, используемый QStringEncoder/
    QStringDecoder; не экспонируется. FinalizeResultError экспонируется как
    enum.
  - Вызов функциональных указателей Interface (toUtf16/toUtf16Len/fromUtf16/
    fromUtf16Len) из Go — требует cgo function-pointer support; указатели
    доступны как unsafe.Pointer, но вызов не реализован.
  - QStringConverter::encoding() — в Qt 6.11 такого метода нет; Encoding()
    выводит кодировку из Name() через EncodingForName().
  - QStringConverter::swap() — в Qt 6.11 такого метода нет; Swap() реализован
    через std::swap protected-полей iface/state.
  - Копирование QStringConverter в Qt 6.11 удалено (класс move-only);
    пересоздавайте конвертер через NewQStringConverterFromEncoding/FromName.
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
)

// Encoding mirrors QStringConverter::Encoding.
// System = 8 in Qt 6.11 (Latin1 = 7, System = 8).
// Qt: QStringConverter::Encoding.
type Encoding int

const (
	// Utf8 mirrors QStringConverter::Utf8.
	// Qt: QStringConverter::Utf8.
	Utf8 Encoding = 0
	// Utf16 mirrors QStringConverter::Utf16.
	// Qt: QStringConverter::Utf16.
	Utf16 Encoding = 1
	// Utf16LE mirrors QStringConverter::Utf16LE.
	// Qt: QStringConverter::Utf16LE.
	Utf16LE Encoding = 2
	// Utf16BE mirrors QStringConverter::Utf16BE.
	// Qt: QStringConverter::Utf16BE.
	Utf16BE Encoding = 3
	// Utf32 mirrors QStringConverter::Utf32.
	// Qt: QStringConverter::Utf32.
	Utf32 Encoding = 4
	// Utf32LE mirrors QStringConverter::Utf32LE.
	// Qt: QStringConverter::Utf32LE.
	Utf32LE Encoding = 5
	// Utf32BE mirrors QStringConverter::Utf32BE.
	// Qt: QStringConverter::Utf32BE.
	Utf32BE Encoding = 6
	// Latin1 mirrors QStringConverter::Latin1.
	// Qt: QStringConverter::Latin1.
	Latin1 Encoding = 7
	// System mirrors QStringConverter::System.
	// Qt: QStringConverter::System.
	System Encoding = 8
	// LastEncoding is an alias for System (the last enum value).
	// Qt: QStringConverter::LastEncoding.
	LastEncoding = System
)

// FinalizeResultError mirrors QStringConverter::FinalizeResultError.
// Qt: QStringConverter::FinalizeResultError.
type FinalizeResultError int

const (
	// FinalizeNoError mirrors FinalizeResultError::NoError.
	// Qt: QStringConverter::FinalizeResultError::NoError.
	FinalizeNoError FinalizeResultError = iota
	// FinalizeInvalidCharacters mirrors FinalizeResultError::InvalidCharacters.
	// Qt: QStringConverter::FinalizeResultError::InvalidCharacters.
	FinalizeInvalidCharacters
	// FinalizeNotEnoughSpace mirrors FinalizeResultError::NotEnoughSpace.
	// Qt: QStringConverter::FinalizeResultError::NotEnoughSpace.
	FinalizeNotEnoughSpace
)

// String returns the Qt enumerator name of e.
// Qt: QStringConverter::FinalizeResultError (enum names).
func (e FinalizeResultError) String() string {
	switch e {
	case FinalizeNoError:
		return "NoError"
	case FinalizeInvalidCharacters:
		return "InvalidCharacters"
	case FinalizeNotEnoughSpace:
		return "NotEnoughSpace"
	}
	return "FinalizeResultError(?)"
}

// Interface wraps a const QStringConverter::Interface* (protected nested
// struct).
//
// Obtain one via (*QStringConverter).Iface().
//
// Qt: QStringConverter::Interface (protected nested struct).
type Interface struct {
	p unsafe.Pointer
}

// IsNil reports whether i wraps a null pointer.
// Qt: no direct counterpart — null check for the protected Interface pointer.
func (i Interface) IsNil() bool { return i.p == nil }

// UnsafePtr returns the raw const QStringConverter::Interface* as
// unsafe.Pointer.
//
// Intended for interop with C++ helpers that need the interface pointer; it is
// not part of the regular public API.
// Qt: no direct counterpart — raw pointer access to Interface.
func (i Interface) UnsafePtr() unsafe.Pointer { return i.p }

// Name returns the encoding name ("UTF-16", "ISO-8859-1", etc.), or "" for a
// null Interface.
// Qt: QStringConverter::Interface::name (const char *).
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

// ToUtf16Ptr returns the raw decoder function pointer
// (QChar* (*)(QChar*, QByteArrayView, State*)), or nil.
//
// Calling it from Go requires cgo function-pointer support and is not
// implemented.
// Qt: QStringConverter::Interface::toUtf16 (function pointer field).
func (i Interface) ToUtf16Ptr() unsafe.Pointer {
	if i.p == nil {
		return nil
	}
	return C.qstringconverter_interface_to_utf16(i.p)
}

// ToUtf16LenPtr returns the raw decoder length function pointer
// (qsizetype (*)(qsizetype)), or nil.
// Qt: QStringConverter::Interface::toUtf16Len (function pointer field).
func (i Interface) ToUtf16LenPtr() unsafe.Pointer {
	if i.p == nil {
		return nil
	}
	return C.qstringconverter_interface_to_utf16_len(i.p)
}

// FromUtf16Ptr returns the raw encoder function pointer
// (char* (*)(char*, QStringView, State*)), or nil.
// Qt: QStringConverter::Interface::fromUtf16 (function pointer field).
func (i Interface) FromUtf16Ptr() unsafe.Pointer {
	if i.p == nil {
		return nil
	}
	return C.qstringconverter_interface_from_utf16(i.p)
}

// FromUtf16LenPtr returns the raw encoder length function pointer
// (qsizetype (*)(qsizetype)), or nil.
// Qt: QStringConverter::Interface::fromUtf16Len (function pointer field).
func (i Interface) FromUtf16LenPtr() unsafe.Pointer {
	if i.p == nil {
		return nil
	}
	return C.qstringconverter_interface_from_utf16_len(i.p)
}

// QStringConverter mirrors Qt 6.11's QStringConverter.
//
// It embeds *QStringConverterBase (matching Qt's class hierarchy). The embedded
// base is a non-owning view over the same C++ object; the *QStringConverter
// wrapper owns the object.
//
// Qt: QStringConverter.
type QStringConverter struct {
	*QStringConverterBase
	state   *qscState
	owned   bool
	cleanup runtime.Cleanup
}

// newQStringConverterFromC takes ownership of a C++ QStringConverter object and
// registers a runtime cleanup for it.
//
// The cleanup receives the *qscState, not the wrapper, so the wrapper stays
// collectable; the state is shared with the embedded base view.
// Qt: no direct counterpart — internal constructor helper.
func newQStringConverterFromC(c *C.go_qsc_t) *QStringConverter {
	if c == nil {
		return nil
	}
	st := &qscState{}
	st.p.Store((*byte)(unsafe.Pointer(c)))
	q := &QStringConverter{
		QStringConverterBase: &QStringConverterBase{state: st},
		state:                st,
		owned:                true,
	}
	q.cleanup = runtime.AddCleanup(q, func(st *qscState) {
		if p := st.claim(); p != nil {
			C.qstringconverter_delete((*C.go_qsc_t)(p))
		}
	}, st)
	return q
}

// tryHandle returns the underlying C++ pointer and whether it is still alive.
// It does not panic: invalid, closed and nil converters are reported as false.
// Qt: no direct counterpart — internal handle helper.
func (q *QStringConverter) tryHandle() (*C.go_qsc_t, bool) {
	if q == nil || q.state == nil {
		return nil, false
	}
	if p := q.state.cpp(); p != nil {
		return (*C.go_qsc_t)(p), true
	}
	return nil, false
}

// NewQStringConverter creates an invalid converter (null iface).
// Qt: QStringConverter::QStringConverter() (protected default constructor).
func NewQStringConverter() *QStringConverter {
	return newQStringConverterFromC(C.qstringconverter_new_default())
}

// NewQStringConverterFromEncoding creates a converter for the given encoding and
// flags.
// Qt: QStringConverter::QStringConverter(Encoding e, Flags f) (protected).
func NewQStringConverterFromEncoding(e Encoding, f Flags) *QStringConverter {
	return newQStringConverterFromC(C.qstringconverter_new_from_encoding(C.int(e), C.int(f)))
}

// NewQStringConverterFromName creates a converter by encoding name (for example
// "UTF-16" or "ISO-8859-1"). The Go string is passed as UTF-8 QAnyStringView.
// Returns an invalid converter if the name is unknown.
// Qt: QStringConverter::QStringConverter(QAnyStringView name, Flags f) (protected) — Go string form.
func NewQStringConverterFromName(name string, f Flags) *QStringConverter {
	var p unsafe.Pointer
	if len(name) > 0 {
		p = unsafe.Pointer(unsafe.StringData(name))
	}
	q := newQStringConverterFromC(C.qstringconverter_new_from_any(
		p, C.ptrdiff_t(len(name)), 0, C.int(f)))
	runtime.KeepAlive(name)
	return q
}

// NewQStringConverterFromAny creates a converter from a QAnyStringView name and
// flags.
// Qt: QStringConverter::QStringConverter(QAnyStringView name, Flags f) (protected) — QAnyStringView form.
func NewQStringConverterFromAny(name QAnyStringView, f Flags) *QStringConverter {
	q := newQStringConverterFromC(C.qstringconverter_new_from_any(
		name.DataUnsafe(), C.ptrdiff_t(name.Size()), C.int(name.Kind()), C.int(f)))
	runtime.KeepAlive(name)
	return q
}

// NewQStringConverterFromIface creates a converter from an existing Interface
// pointer (the protected constructor is exposed). If i is null, the converter is
// invalid.
// Qt: QStringConverter::QStringConverter(const Interface *i) (protected).
func NewQStringConverterFromIface(i Interface) *QStringConverter {
	return newQStringConverterFromC(C.qstringconverter_new_from_iface(i.p))
}

// Close frees the underlying C++ object immediately and cancels the cleanup.
// For a non-owning converter (created by NewQStringConverterFromUnsafePtr),
// Close is a no-op.
//
// Close is optional, idempotent and nil-safe. If it races with the runtime
// cleanup, exactly one of them claims the object, so there is no double free.
// Qt: QStringConverter::~QStringConverter() (destructor).
func (q *QStringConverter) Close() {
	if q == nil || q.state == nil || !q.owned {
		return
	}
	q.cleanup.Stop()
	if p := q.state.claim(); p != nil {
		C.qstringconverter_delete((*C.go_qsc_t)(p))
	}
}

// Swap exchanges the contents of q and other. It is a no-op if either side is
// nil or closed.
// Qt: no QStringConverter::swap() in Qt 6.11 — implemented via std::swap on the
// protected iface and state fields.
func (q *QStringConverter) Swap(other *QStringConverter) {
	p, ok := q.tryHandle()
	if !ok {
		return
	}
	op, ok := other.tryHandle()
	if !ok {
		return
	}
	C.qstringconverter_swap(p, op)
}

// Encoding returns the converter's encoding, derived from Name via
// EncodingForName. Returns Encoding(-1) if the converter is invalid.
// Qt: no public QStringConverter::encoding() in Qt 6.11 — derived from name().
func (q *QStringConverter) Encoding() Encoding {
	if _, ok := q.tryHandle(); !ok {
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

// IsValid reports whether the converter has a valid encoding interface.
// Qt: QStringConverter::isValid() const.
func (q *QStringConverter) IsValid() bool {
	p, ok := q.tryHandle()
	if !ok {
		return false
	}
	return bool(C.qstringconverter_is_valid(p))
}

// HasError reports whether the converter has encountered invalid characters.
// An invalid converter (null interface) reports true; a fresh valid converter
// reports false.
// Qt: QStringConverter::hasError() const.
func (q *QStringConverter) HasError() bool {
	p, ok := q.tryHandle()
	if !ok {
		return true
	}
	return bool(C.qstringconverter_has_error(p))
}

// Name returns the encoding name (for example "UTF-16").
// Qt: QStringConverter::name() const.
func (q *QStringConverter) Name() string {
	p, ok := q.tryHandle()
	if !ok {
		return ""
	}
	cname := C.qstringconverter_name(p)
	if cname == nil {
		return ""
	}
	return C.GoString(cname)
}

// ResetState resets the conversion state.
// Qt: QStringConverter::resetState().
func (q *QStringConverter) ResetState() {
	p, ok := q.tryHandle()
	if !ok {
		return
	}
	C.qstringconverter_reset_state(p)
}

// State returns a copy of the converter's internal State (the protected `state`
// field).
// Qt: QStringConverter::state (protected data member).
func (q *QStringConverter) State() State {
	var s State
	p, ok := q.tryHandle()
	if !ok {
		return s
	}
	var flags C.int
	var internal C.int
	var remaining C.int64_t
	var invalid C.int64_t
	var data [4]C.uint32_t
	C.qstringconverter_get_state(p, &flags, &internal, &remaining, &invalid, &data[0])
	s.Flags = Flags(flags)
	s.InternalState = int32(internal)
	s.RemainingChars = int64(remaining)
	s.InvalidChars = int64(invalid)
	for i := 0; i < 4; i++ {
		s.Data[i] = uint32(data[i])
	}
	return s
}

// SetState overwrites the converter's internal State (the protected `state`
// field). The C++ side clears the previous state first and then copies the
// fields, so any clearFn of the current encoding is honoured.
// Qt: QStringConverter::state (protected data member).
func (q *QStringConverter) SetState(s State) {
	p, ok := q.tryHandle()
	if !ok {
		return
	}
	var data [4]C.uint32_t
	for i := 0; i < 4; i++ {
		data[i] = C.uint32_t(s.Data[i])
	}
	C.qstringconverter_set_state(p,
		C.int(s.Flags), C.int(s.InternalState),
		C.int64_t(s.RemainingChars), C.int64_t(s.InvalidChars), &data[0])
}

// Iface returns the converter's Interface pointer (the protected `iface` field),
// or a null Interface.
// Qt: QStringConverter::iface (protected data member).
func (q *QStringConverter) Iface() Interface {
	p, ok := q.tryHandle()
	if !ok {
		return Interface{}
	}
	return Interface{p: C.qstringconverter_iface(p)}
}

// NameForEncoding returns the canonical name for e, or "" if e is unknown.
// Qt: QStringConverter::nameForEncoding(Encoding e) static.
func NameForEncoding(e Encoding) string {
	cname := C.qstringconverter_name_for_encoding(C.int(e))
	if cname == nil {
		return ""
	}
	return C.GoString(cname)
}

// EncodingForName returns the Encoding for the given name and true if it is
// recognized. Otherwise it returns (Utf8, false). This is the Go string form of
// Qt's std::optional<Encoding> result.
// Qt: QStringConverter::encodingForName(QAnyStringView name) static — Go string form.
func EncodingForName(name string) (Encoding, bool) {
	var p unsafe.Pointer
	if len(name) > 0 {
		p = unsafe.Pointer(unsafe.StringData(name))
	}
	v := C.qstringconverter_encoding_for_name(p, C.ptrdiff_t(len(name)), 0)
	runtime.KeepAlive(name)
	if v < 0 {
		return Utf8, false
	}
	return Encoding(v), true
}

// EncodingForNameAny is EncodingForName accepting a QAnyStringView directly.
// Qt: QStringConverter::encodingForName(QAnyStringView name) static — QAnyStringView form.
func EncodingForNameAny(name QAnyStringView) (Encoding, bool) {
	v := C.qstringconverter_encoding_for_name(
		name.DataUnsafe(), C.ptrdiff_t(name.Size()), C.int(name.Kind()))
	runtime.KeepAlive(name)
	if v < 0 {
		return Utf8, false
	}
	return Encoding(v), true
}

// EncodingForData returns the Encoding detected from the given view, and true if
// it is recognized; otherwise (Utf8, false). If expectedFirstChar is non-zero,
// the detected encoding is accepted only when the decoded data starts with that
// code unit.
// Qt: QStringConverter::encodingForData(QByteArrayView data, char16_t expectedFirstCharacter) static — QByteArrayView form.
func EncodingForData(data QByteArrayView, expectedFirstChar uint16) (Encoding, bool) {
	return EncodingForDataBytes(data.Bytes(), expectedFirstChar)
}

// EncodingForDataBytes is EncodingForData accepting a []byte directly.
// Qt: QStringConverter::encodingForData(QByteArrayView data, char16_t expectedFirstCharacter) static — []byte form.
func EncodingForDataBytes(data []byte, expectedFirstChar uint16) (Encoding, bool) {
	var p unsafe.Pointer
	if len(data) > 0 {
		p = unsafe.Pointer(&data[0])
	}
	v := C.qstringconverter_encoding_for_data(
		p, C.ptrdiff_t(len(data)), C.uint16_t(expectedFirstChar))
	runtime.KeepAlive(data)
	if v < 0 {
		return Utf8, false
	}
	return Encoding(v), true
}

// EncodingForHtml returns the Encoding detected from the HTML <meta> charset
// declaration in the given view, and true if it is recognized; otherwise
// (Utf8, false).
// Qt: QStringConverter::encodingForHtml(QByteArrayView data) static — QByteArrayView form.
func EncodingForHtml(data QByteArrayView) (Encoding, bool) {
	return EncodingForHtmlBytes(data.Bytes())
}

// EncodingForHtmlBytes is EncodingForHtml accepting a []byte directly.
// Qt: QStringConverter::encodingForHtml(QByteArrayView data) static — []byte form.
func EncodingForHtmlBytes(data []byte) (Encoding, bool) {
	var p unsafe.Pointer
	if len(data) > 0 {
		p = unsafe.Pointer(&data[0])
	}
	v := C.qstringconverter_encoding_for_html(p, C.ptrdiff_t(len(data)))
	runtime.KeepAlive(data)
	if v < 0 {
		return Utf8, false
	}
	return Encoding(v), true
}

// AvailableCodecsStrings returns the names of all available codecs as Go
// strings.
// Qt: QStringConverter::availableCodecs() static — []string form of the returned QStringList.
func AvailableCodecsStrings() []string {
	var n C.ptrdiff_t
	arr := C.qstringconverter_available_codecs(&n)
	if arr == nil || n == 0 {
		return nil
	}
	defer C.qstringconverter_available_codecs_free(arr, n)

	cSlice := unsafe.Slice(arr, int(n))
	names := make([]string, int(n))
	for i, p := range cSlice {
		names[i] = C.GoString(p)
	}
	return names
}

// AvailableCodecs returns the names of all available codecs as a QStringList.
// Qt: QStringConverter::availableCodecs() static.
func AvailableCodecs() *QStringList {
	return QStringListFromGoStrings(AvailableCodecsStrings())
}
