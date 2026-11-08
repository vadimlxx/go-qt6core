package qstringdecoder

// Package qstringdecoder: QStringDecoder — конкретный декодер байты → UTF-16.

/*
 ВАЖНО (Qt 6.11.2):
   - QStringDecoder наследует QStringConverter (тот — QStringConverterBase).
     В Go структура *QStringDecoder встраивает *qstringconverter.QStringConverter,
     который встраивает *qstringconverter.QStringConverterBase — точное воспроизведение
     иерархии Qt. Все публичные методы базы доступны через promotion.
   - decode(...) в Qt возвращает EncodedData<T>, лениво конвертируемый в QString.
     В Go Decode/DecodeView сразу возвращают *qstring.QString — семантически эквивалентно.
   - finalize(char16_t*, qsizetype) возвращает FinalizeResultChar<char16_t>; в Go — тройка
     (next int, invalid int16, err FinalizeResultError). QChar* и char16_t* перегрузки
     finalize/appendToBuffer для Go идентичны, поскольку QChar == uint16.
   - qstringconverter.AdoptUnsafePtr используется для того, чтобы встроенный
     *QStringConverter не владел C++-объектом: владелец — *QStringDecoder.

 не реализовано:
   - operator() / decode(const QByteArray&) vs decode(QByteArrayView) — в Go разные методы
     (Decode / DecodeView), без перегрузки оператора ().
   - implicit conversion EncodedData<T> → QString — в Go не нужен.
   - QStringBuilder integration:
       QConcatenable<QStringDecoder::EncodedData<T>> и operator+=(QString&, …) —
       используйте s.AppendQString(dec.DecodeView(bv)) или
       qstringbuilder.Concat(s, dec.DecodeView(bv)).ToQString().
   - QString::assign_helper_char8 — внутренняя деталь QString, не относится к
     QStringDecoder; не экспонируется.
   - FinalizeResultChar<Char> — в публичном API только для char16_t / QChar; в Go —
     тройка значений (next, invalid, err).
   - decoderForHtml(QByteArrayView) — экспонирован как пакетная функция DecoderForHtml.
   - availableCodecs() (inherited static) — требует qstringlist; TODO в qstringconverter.
   - encodingForData/encodingForHtml — доступны через qstringconverter.EncodingForData /
     EncodingForHtml.
*/

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC -Wno-sfinae-incomplete
#cgo LDFLAGS: -lstdc++
#include <stdlib.h>
#include <stdint.h>
#include <stdbool.h>
#include <stddef.h>
#include "qstringdecoder_c.h"
*/
import "C"

import (
	"runtime"
	"unsafe"

	"github.com/vadimlxx/go-qt6core/qanystringview"
	"github.com/vadimlxx/go-qt6core/qbytearray"
	"github.com/vadimlxx/go-qt6core/qbytearrayview"
	"github.com/vadimlxx/go-qt6core/qstring"
	"github.com/vadimlxx/go-qt6core/qstringconverter"
)

// FinalizeResultError mirrors QStringConverter::FinalizeResultError.
type FinalizeResultError int

const (
	FinalizeNoError FinalizeResultError = iota
	FinalizeInvalidCharacters
	FinalizeNotEnoughSpace
)

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

// QStringDecoder mirrors Qt 6.11.2 QStringDecoder.
//
// It embeds *qstringconverter.QStringConverter (which itself embeds
// *qstringconverter.QStringConverterBase), reproducing Qt's class hierarchy:
//
//	QStringDecoder → QStringConverter → QStringConverterBase
//
// Every public method of the base chain (IsValid, Name, ResetState, HasError, State, SetState,
// Iface, Swap, Encoding, Flags, …) is available on *QStringDecoder via Go method promotion.
type QStringDecoder struct {
	*qstringconverter.QStringConverter
	p *C.go_qsd_t
}

func newQStringDecoderFromC(c *C.go_qsd_t) *QStringDecoder {
	if c == nil {
		return nil
	}
	// AdoptUnsafePtr creates a non-owning *QStringConverter alias over the same C++ object.
	// Ownership stays with the *QStringDecoder wrapper (its finalizer calls qstringdecoder_delete).
	base := qstringconverter.AdoptUnsafePtr(unsafe.Pointer(c))
	dec := &QStringDecoder{
		QStringConverter: base,
		p:                c,
	}
	runtime.SetFinalizer(dec, finalizeQStringDecoder)
	return dec
}

func finalizeQStringDecoder(dec *QStringDecoder) {
	if dec.p != nil {
		C.qstringdecoder_delete(dec.p)
		dec.p = nil
		// Drop the alias; its C++ memory is gone.
		dec.QStringConverter = nil
	}
}

// NewQStringDecoder creates an invalid decoder (null iface, FlagDefault).
// Mirrors QStringDecoder().
func NewQStringDecoder() *QStringDecoder {
	return newQStringDecoderFromC(C.qstringdecoder_new_default())
}

// NewQStringDecoderFromEncoding creates a decoder for the given encoding and flags.
// Mirrors QStringDecoder(Encoding e, Flags f = Flag::Default).
func NewQStringDecoderFromEncoding(e qstringconverter.Encoding, f qstringconverter.Flags) *QStringDecoder {
	return newQStringDecoderFromC(C.qstringdecoder_new_from_encoding(C.int(int(e)), C.int(int(f))))
}

// NewQStringDecoderFromName creates a decoder by encoding name (e.g. "UTF-8", "ISO-8859-1").
// Mirrors QStringDecoder(QAnyStringView name, Flags f) with a Go string.
// Returns an invalid decoder if the name is unknown.
func NewQStringDecoderFromName(name string, f qstringconverter.Flags) *QStringDecoder {
	var p unsafe.Pointer
	if len(name) > 0 {
		p = unsafe.Pointer(unsafe.StringData(name))
	}
	dec := newQStringDecoderFromC(C.qstringdecoder_new_from_any(
		p, C.ptrdiff_t(len(name)), 0, C.int(int(f))))
	runtime.KeepAlive(name)
	return dec
}

// NewQStringDecoderFromAny creates a decoder from a QAnyStringView name and flags.
// Mirrors QStringDecoder(QAnyStringView name, Flags f).
func NewQStringDecoderFromAny(name qanystringview.QAnyStringView, f qstringconverter.Flags) *QStringDecoder {
	dec := newQStringDecoderFromC(C.qstringdecoder_new_from_any(
		name.DataUnsafe(), C.ptrdiff_t(name.Size()), C.int(int(name.Kind())), C.int(int(f))))
	runtime.KeepAlive(name)
	return dec
}

// NewQStringDecoderFromInterface wraps the protected QStringDecoder(const Interface*) ctor.
// If iface is nil, the decoder is invalid.
func NewQStringDecoderFromInterface(i qstringconverter.Interface) *QStringDecoder {
	return newQStringDecoderFromC(C.qstringdecoder_new_from_iface(i.UnsafePtr()))
}

// Close frees the underlying C++ object immediately and disables the finalizer.
// After Close, all methods on d (including promoted base methods) are UB.
func (d *QStringDecoder) Close() {
	if d.p != nil {
		C.qstringdecoder_delete(d.p)
		d.p = nil
		d.QStringConverter = nil
		runtime.SetFinalizer(d, nil)
	}
}

// Decode decodes the contents of ba as a fresh *qstring.QString.
// Equivalent to Qt's implicit conversion of QStringDecoder::decode(const QByteArray&) to QString.
func (d *QStringDecoder) Decode(ba *qbytearray.QByteArray) *qstring.QString {
	if d.p == nil || ba == nil {
		return qstring.New()
	}
	n := ba.Size()
	var data *C.char
	if n > 0 {
		data = (*C.char)(ba.ConstDataUnsafe())
	}
	var out unsafe.Pointer
	C.qstringdecoder_decode(d.p, data, C.ptrdiff_t(n), &out)
	runtime.KeepAlive(ba)
	return qstring.AdoptUnsafePtr(out)
}

// DecodeView decodes v as a fresh *qstring.QString.
// Equivalent to Qt's implicit conversion of QStringDecoder::decode(QByteArrayView) to QString.
func (d *QStringDecoder) DecodeView(v qbytearrayview.QByteArrayView) *qstring.QString {
	if d.p == nil {
		return qstring.New()
	}
	n := v.Size()
	var data *C.char
	if n > 0 {
		data = (*C.char)(v.DataUnsafe())
	}
	var out unsafe.Pointer
	C.qstringdecoder_decode(d.p, data, C.ptrdiff_t(n), &out)
	runtime.KeepAlive(v)
	return qstring.AdoptUnsafePtr(out)
}

// RequiredSpace returns the maximum number of UTF-16 code units needed to decode
// inputLength bytes (mirrors QStringDecoder::requiredSpace). Returns 0 for an invalid decoder.
func (d *QStringDecoder) RequiredSpace(inputLength int) int {
	if d.p == nil {
		return 0
	}
	return int(C.qstringdecoder_required_space(d.p, C.ptrdiff_t(inputLength)))
}

// AppendToBuffer decodes v into buf (at buf[0]) and returns the number of UTF-16 code units
// written. The caller must ensure len(buf) >= RequiredSpace(v.Size()).
// For an invalid decoder, nothing is written, the decoder's HasError becomes true, and 0 is
// returned. Mirrors QStringDecoder::appendToBuffer.
func (d *QStringDecoder) AppendToBuffer(buf []uint16, v qbytearrayview.QByteArrayView) int {
	if d.p == nil {
		return 0
	}
	var bufPtr unsafe.Pointer
	if len(buf) > 0 {
		bufPtr = unsafe.Pointer(&buf[0])
	}
	n := v.Size()
	var data *C.char
	if n > 0 {
		data = (*C.char)(v.DataUnsafe())
	}
	written := C.qstringdecoder_append_to_buffer(d.p, bufPtr, C.ptrdiff_t(len(buf)),
		data, C.ptrdiff_t(n))
	runtime.KeepAlive(buf)
	runtime.KeepAlive(v)
	return int(written)
}

// FinalizeBuffer flushes pending decoder state into buf (capacity len(buf), UTF-16 units).
// Returns (next int, invalid int16, err FinalizeResultError) where:
//   - next is the number of UTF-16 code units written into buf,
//   - invalid is the number of invalid characters seen during this finalize call,
//   - err is NoError / InvalidCharacters / NotEnoughSpace.
// Mirrors QStringDecoder::finalize(char16_t *out, qsizetype maxlen).
func (d *QStringDecoder) FinalizeBuffer(buf []uint16) (int, int16, FinalizeResultError) {
	if d.p == nil {
		return 0, 0, FinalizeNotEnoughSpace
	}
	var bufPtr unsafe.Pointer
	if len(buf) > 0 {
		bufPtr = unsafe.Pointer(&buf[0])
	}
	var next C.ptrdiff_t
	var invalid C.int16_t
	var errCode C.int
	C.qstringdecoder_finalize(d.p, bufPtr, C.ptrdiff_t(len(buf)), &next, &invalid, &errCode)
	runtime.KeepAlive(buf)
	return int(next), int16(invalid), FinalizeResultError(errCode)
}

// Finalize flushes pending decoder state without providing a buffer.
// Returns (invalid int16, err FinalizeResultError).
// Mirrors QStringDecoder::finalize() (the no-argument overload).
func (d *QStringDecoder) Finalize() (int16, FinalizeResultError) {
	if d.p == nil {
		return 0, FinalizeNotEnoughSpace
	}
	var next C.ptrdiff_t
	var invalid C.int16_t
	var errCode C.int
	C.qstringdecoder_finalize(d.p, nil, 0, &next, &invalid, &errCode)
	return int16(invalid), FinalizeResultError(errCode)
}

// DecoderForHtml returns a decoder for the encoding declared in the HTML <meta> charset
// of data (or UTF-8 if none is found).
// Mirrors the static QStringDecoder::decoderForHtml(QByteArrayView).
func DecoderForHtml(data qbytearrayview.QByteArrayView) *QStringDecoder {
	n := data.Size()
	var p *C.char
	if n > 0 {
		p = (*C.char)(data.DataUnsafe())
	}
	c := C.qstringdecoder_decoder_for_html(p, C.ptrdiff_t(n))
	runtime.KeepAlive(data)
	return newQStringDecoderFromC(c)
}
