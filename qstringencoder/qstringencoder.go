package qstringencoder

// Package qstringencoder: QStringEncoder — конкретный энкодер UTF-16 → байты.

/*
 ВАЖНО (Qt 6.11.2):
   - QStringEncoder наследует QStringConverter (тот — QStringConverterBase).
     В Go структура *QStringEncoder встраивает *qstringconverter.QStringConverter,
     который встраивает *qstringconverter.QStringConverterBase — точное воспроизведение
     иерархии Qt. Все публичные методы базы доступны через promotion.
   - encode(...) в Qt возвращает DecodedData<T>, лениво конвертируемый в QByteArray.
     В Go Encode/EncodeView сразу возвращают *qbytearray.QByteArray — семантически
     эквивалентно (одна конверсия в C++ с текущим State и последующим truncate).
   - finalize(char*, qsizetype) возвращает FinalizeResultChar<char>; в Go — тройка
     (next int, invalid int16, err FinalizeResultError).
   - qstringconverter.AdoptUnsafePtr используется для того, чтобы встроенный
     *QStringConverter не владел C++-объектом: владелец — *QStringEncoder.

 не реализовано:
   - operator() / encode(const QString&) vs encode(QStringView) — в Go разные методы
     (Encode / EncodeView), без перегрузки оператора ().
   - implicit conversion DecodedData<T> → QByteArray — в Go не нужен.
   - QStringBuilder integration:
       QConcatenable<QStringEncoder::DecodedData<T>> и operator+=(QByteArray&, …) —
       используйте ba.AppendQByteArray(enc.EncodeView(sv)) или
       qstringbuilder.FromQByteArray(enc.EncodeView(sv)).
   - QString::assign_helper_char8 — внутренняя деталь QString, не относится к
     QStringEncoder; не экспонируется.
   - FinalizeResultChar<Char> — в публичном API только для char; в Go — тройка значений.
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
#include "qstringencoder_c.h"
*/
import "C"

import (
	"runtime"
	"unsafe"

	"github.com/vadimlxx/go-qt6core/qanystringview"
	"github.com/vadimlxx/go-qt6core/qbytearray"
	"github.com/vadimlxx/go-qt6core/qstring"
	"github.com/vadimlxx/go-qt6core/qstringconverter"
	"github.com/vadimlxx/go-qt6core/qstringview"
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

// QStringEncoder mirrors Qt 6.11.2 QStringEncoder.
//
// It embeds *qstringconverter.QStringConverter (which itself embeds
// *qstringconverter.QStringConverterBase), reproducing Qt's class hierarchy:
//
//	QStringEncoder → QStringConverter → QStringConverterBase
//
// Every public method of the base chain (IsValid, Name, ResetState, HasError, State, SetState,
// Iface, Swap, Encoding, Flags, …) is available on *QStringEncoder via Go method promotion.
type QStringEncoder struct {
	*qstringconverter.QStringConverter
	p *C.go_qse_t
}

func newQStringEncoderFromC(c *C.go_qse_t) *QStringEncoder {
	if c == nil {
		return nil
	}
	// AdoptUnsafePtr creates a non-owning *QStringConverter alias over the same C++ object.
	// Ownership stays with the *QStringEncoder wrapper (its finalizer calls qstringencoder_delete).
	base := qstringconverter.AdoptUnsafePtr(unsafe.Pointer(c))
	enc := &QStringEncoder{
		QStringConverter: base,
		p:                c,
	}
	runtime.SetFinalizer(enc, finalizeQStringEncoder)
	return enc
}

func finalizeQStringEncoder(enc *QStringEncoder) {
	if enc.p != nil {
		C.qstringencoder_delete(enc.p)
		enc.p = nil
		// Drop the alias; its C++ memory is gone.
		enc.QStringConverter = nil
	}
}

// NewQStringEncoder creates an invalid encoder (null iface, FlagDefault).
// Mirrors QStringEncoder().
func NewQStringEncoder() *QStringEncoder {
	return newQStringEncoderFromC(C.qstringencoder_new_default())
}

// NewQStringEncoderFromEncoding creates an encoder for the given encoding and flags.
// Mirrors QStringEncoder(Encoding e, Flags f = Flag::Default).
func NewQStringEncoderFromEncoding(e qstringconverter.Encoding, f qstringconverter.Flags) *QStringEncoder {
	return newQStringEncoderFromC(C.qstringencoder_new_from_encoding(C.int(int(e)), C.int(int(f))))
}

// NewQStringEncoderFromName creates an encoder by encoding name (e.g. "UTF-8", "ISO-8859-1").
// Mirrors QStringEncoder(QAnyStringView name, Flags f) with a Go string.
// Returns an invalid encoder if the name is unknown.
func NewQStringEncoderFromName(name string, f qstringconverter.Flags) *QStringEncoder {
	var p unsafe.Pointer
	if len(name) > 0 {
		p = unsafe.Pointer(unsafe.StringData(name))
	}
	enc := newQStringEncoderFromC(C.qstringencoder_new_from_any(
		p, C.ptrdiff_t(len(name)), 0, C.int(int(f))))
	runtime.KeepAlive(name)
	return enc
}

// NewQStringEncoderFromAny creates an encoder from a QAnyStringView name and flags.
// Mirrors QStringEncoder(QAnyStringView name, Flags f).
func NewQStringEncoderFromAny(name qanystringview.QAnyStringView, f qstringconverter.Flags) *QStringEncoder {
	enc := newQStringEncoderFromC(C.qstringencoder_new_from_any(
		name.DataUnsafe(), C.ptrdiff_t(name.Size()), C.int(int(name.Kind())), C.int(int(f))))
	runtime.KeepAlive(name)
	return enc
}

// NewQStringEncoderFromInterface wraps the protected QStringEncoder(const Interface*) ctor.
// If iface is nil, the encoder is invalid.
func NewQStringEncoderFromInterface(i qstringconverter.Interface) *QStringEncoder {
	return newQStringEncoderFromC(C.qstringencoder_new_from_iface(i.UnsafePtr()))
}

// Close frees the underlying C++ object immediately and disables the finalizer.
// After Close, all methods on e (including promoted base methods) are UB.
func (e *QStringEncoder) Close() {
	if e.p != nil {
		C.qstringencoder_delete(e.p)
		e.p = nil
		e.QStringConverter = nil
		runtime.SetFinalizer(e, nil)
	}
}

// Encode encodes the contents of q as a fresh *qbytearray.QByteArray.
// Equivalent to Qt's implicit conversion of QStringEncoder::encode(const QString&) to QByteArray.
func (e *QStringEncoder) Encode(q *qstring.QString) *qbytearray.QByteArray {
	if e.p == nil || q == nil {
		return qbytearray.New()
	}
	n := q.Size()
	var data *C.uint16_t
	if n > 0 {
		data = (*C.uint16_t)(q.Utf16Unsafe())
	}
	var out unsafe.Pointer
	var outSize C.ptrdiff_t
	C.qstringencoder_encode(e.p, data, C.ptrdiff_t(n), &out, &outSize)
	runtime.KeepAlive(q)
	return qbytearray.AdoptUnsafePtr(out)
}

// EncodeView encodes v as a fresh *qbytearray.QByteArray.
// Equivalent to Qt's implicit conversion of QStringEncoder::encode(QStringView) to QByteArray.
func (e *QStringEncoder) EncodeView(v qstringview.QStringView) *qbytearray.QByteArray {
	if e.p == nil {
		return qbytearray.New()
	}
	n := v.Size()
	var data *C.uint16_t
	if n > 0 {
		data = (*C.uint16_t)(v.Utf16Unsafe())
	}
	var out unsafe.Pointer
	var outSize C.ptrdiff_t
	C.qstringencoder_encode(e.p, data, C.ptrdiff_t(n), &out, &outSize)
	runtime.KeepAlive(v)
	return qbytearray.AdoptUnsafePtr(out)
}

// RequiredSpace returns the maximum number of bytes needed to encode inputLength UTF-16
// code units (mirrors QStringEncoder::requiredSpace). Returns 0 for an invalid encoder.
func (e *QStringEncoder) RequiredSpace(inputLength int) int {
	if e.p == nil {
		return 0
	}
	return int(C.qstringencoder_required_space(e.p, C.ptrdiff_t(inputLength)))
}

// AppendToBuffer writes the encoded bytes of v into buf (at buf[0]) and returns the number
// of bytes written. The caller must ensure len(buf) >= RequiredSpace(v.Size()).
// For an invalid encoder, nothing is written, the encoder's HasError becomes true, and 0
// is returned. Mirrors QStringEncoder::appendToBuffer.
func (e *QStringEncoder) AppendToBuffer(buf []byte, v qstringview.QStringView) int {
	if e.p == nil {
		return 0
	}
	var bufPtr unsafe.Pointer
	if len(buf) > 0 {
		bufPtr = unsafe.Pointer(&buf[0])
	}
	n := v.Size()
	var data *C.uint16_t
	if n > 0 {
		data = (*C.uint16_t)(v.Utf16Unsafe())
	}
	written := C.qstringencoder_append_to_buffer(e.p, bufPtr, C.ptrdiff_t(len(buf)),
		data, C.ptrdiff_t(n))
	runtime.KeepAlive(buf)
	runtime.KeepAlive(v)
	return int(written)
}

// FinalizeBuffer flushes pending encoder state into buf (capacity len(buf)).
// Returns (next int, invalid int16, err FinalizeResultError) where:
//   - next is the number of bytes written into buf,
//   - invalid is the number of invalid characters seen during this finalize call,
//   - err is NoError / InvalidCharacters / NotEnoughSpace.
// Mirrors QStringEncoder::finalize(char *out, qsizetype maxlen).
func (e *QStringEncoder) FinalizeBuffer(buf []byte) (int, int16, FinalizeResultError) {
	if e.p == nil {
		return 0, 0, FinalizeNotEnoughSpace
	}
	var bufPtr unsafe.Pointer
	if len(buf) > 0 {
		bufPtr = unsafe.Pointer(&buf[0])
	}
	var next C.ptrdiff_t
	var invalid C.int16_t
	var errCode C.int
	C.qstringencoder_finalize(e.p, bufPtr, C.ptrdiff_t(len(buf)), &next, &invalid, &errCode)
	runtime.KeepAlive(buf)
	return int(next), int16(invalid), FinalizeResultError(errCode)
}

// Finalize flushes pending encoder state without providing a buffer.
// Returns (invalid int16, err FinalizeResultError).
// Mirrors QStringEncoder::finalize() (the no-argument overload).
func (e *QStringEncoder) Finalize() (int16, FinalizeResultError) {
	if e.p == nil {
		return 0, FinalizeNotEnoughSpace
	}
	var next C.ptrdiff_t
	var invalid C.int16_t
	var errCode C.int
	C.qstringencoder_finalize(e.p, nil, 0, &next, &invalid, &errCode)
	return int16(invalid), FinalizeResultError(errCode)
}
