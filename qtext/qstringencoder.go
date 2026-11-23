package qtext

/*
QStringEncoder — конкретный кодировщик Qt 6.11 QtCore (UTF-16 -> байты).

Иерархия воспроизведена встраиванием, как в Qt:

	QStringEncoder -> QStringConverter -> QStringConverterBase

*QStringEncoder владеет C++ объектом: он освобождается Close() или
runtime.AddCleanup (когда Go-обёртка становится недостижимой). Встроенная
*QStringConverter — не владеющая проекция на тот же объект и на тот же
qscState: после Close()/cleanup конвертер видит nil, а не висячий указатель.

C++-мост общий для QStringConverter/QStringEncoder/QStringDecoder
(qstringconverter_c.cpp/.h), как qstringconverter.cpp/.h в самом QtCore.

Go-варианты:

  - EncodeString / EncodeViewString / AppendToBufferString принимают string
    вместо QString/QStringView;
  - EncodeBytes / EncodeViewBytes возвращают []byte вместо *QByteArray.

Не реализовано:

  - operator() — в Go нет перегрузки операторов; используйте Encode/EncodeView.
  - DecodedData<T> как публичный тип — не экспонируется; Encode сразу
    материализует QByteArray.
  - QStringBuilder-интеграция (QConcatenable, operator+=) — используйте
    AppendQByteArray/Concat.
  - FinalizeResultChar<Char> — шаблон; в Go результат представлен тройкой
    (next, invalid, err).
*/

/*
#cgo pkg-config: Qt6Core
#cgo CXXFLAGS: -std=c++17 -fPIC -Wno-sfinae-incomplete
#cgo LDFLAGS: -lstdc++
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

// QStringEncoder mirrors Qt 6.11's QStringEncoder.
//
// It embeds *QStringConverter (which embeds *QStringConverterBase), matching
// Qt's class hierarchy. The embedded converter is a non-owning view over the
// same C++ object and the same qscState; the *QStringEncoder wrapper owns the
// object.
//
// Qt: QStringEncoder.
type QStringEncoder struct {
	*QStringConverter
	state   *qscState
	cleanup runtime.Cleanup
}

// newQStringEncoderFromC takes ownership of a C++ QStringEncoder object,
// builds the embedded non-owning converter view over the same state, and
// registers a runtime cleanup for the encoder.
//
// Qt: no direct counterpart — internal constructor helper.
func newQStringEncoderFromC(c *C.go_qse_t) *QStringEncoder {
	if c == nil {
		return nil
	}
	st := &qscState{}
	st.p.Store((*byte)(unsafe.Pointer(c)))
	enc := &QStringEncoder{
		QStringConverter: newQStringConverterAlias(st),
		state:            st,
	}
	enc.cleanup = runtime.AddCleanup(enc, func(st *qscState) {
		if p := st.claim(); p != nil {
			C.qstringencoder_delete((*C.go_qse_t)(p))
		}
	}, st)
	return enc
}

// tryHandle returns the underlying C++ QStringEncoder pointer and whether it is
// still alive. It does not panic: closed and nil encoders are reported as false.
//
// Qt: no direct counterpart — internal handle helper.
func (e *QStringEncoder) tryHandle() (*C.go_qse_t, bool) {
	if e == nil || e.state == nil {
		return nil, false
	}
	if p := e.state.cpp(); p != nil {
		return (*C.go_qse_t)(p), true
	}
	return nil, false
}

// NewQStringEncoder creates an invalid encoder (null iface, FlagDefault).
//
// Qt: QStringEncoder::QStringEncoder().
func NewQStringEncoder() *QStringEncoder {
	return newQStringEncoderFromC(C.qstringencoder_new_default())
}

// NewQStringEncoderFromEncoding creates an encoder for the given encoding and
// flags.
//
// Qt: QStringEncoder::QStringEncoder(Encoding e, Flags f).
func NewQStringEncoderFromEncoding(e Encoding, f Flags) *QStringEncoder {
	return newQStringEncoderFromC(C.qstringencoder_new_from_encoding(C.int(e), C.int(f)))
}

// NewQStringEncoderFromName creates an encoder by encoding name (for example
// "UTF-8" or "ISO-8859-1"). The Go string is passed as UTF-8 QAnyStringView.
// Returns an invalid encoder if the name is unknown.
//
// Qt: QStringEncoder::QStringEncoder(QAnyStringView name, Flags f) — Go string form.
func NewQStringEncoderFromName(name string, f Flags) *QStringEncoder {
	var p unsafe.Pointer
	if len(name) > 0 {
		p = unsafe.Pointer(unsafe.StringData(name))
	}
	enc := newQStringEncoderFromC(C.qstringencoder_new_from_any(
		p, C.ptrdiff_t(len(name)), 0, C.int(f)))
	runtime.KeepAlive(name)
	return enc
}

// NewQStringEncoderFromAny creates an encoder from a QAnyStringView name and
// flags.
//
// Qt: QStringEncoder::QStringEncoder(QAnyStringView name, Flags f) — QAnyStringView form.
func NewQStringEncoderFromAny(name QAnyStringView, f Flags) *QStringEncoder {
	enc := newQStringEncoderFromC(C.qstringencoder_new_from_any(
		name.DataUnsafe(), C.ptrdiff_t(name.Size()), C.int(name.Kind()), C.int(f)))
	runtime.KeepAlive(name)
	return enc
}

// NewQStringEncoderFromInterface creates an encoder from an existing Interface
// pointer (the protected constructor is exposed). If i is null, the encoder is
// invalid.
//
// Qt: QStringEncoder::QStringEncoder(const Interface *i) (protected).
func NewQStringEncoderFromInterface(i Interface) *QStringEncoder {
	return newQStringEncoderFromC(C.qstringencoder_new_from_iface(i.UnsafePtr()))
}

// Close frees the underlying C++ object immediately and cancels the cleanup.
// Close is optional, idempotent and nil-safe. The embedded converter view
// shares the state slot, so after Close its methods observe a nil handle.
//
// Qt: QStringEncoder::~QStringEncoder().
func (e *QStringEncoder) Close() {
	if e == nil || e.state == nil {
		return
	}
	e.cleanup.Stop()
	if p := e.state.claim(); p != nil {
		C.qstringencoder_delete((*C.go_qse_t)(p))
	}
}

// Encode encodes q and returns a fresh *QByteArray. For an invalid encoder the
// result is empty and HasError becomes true.
//
// Qt: QStringEncoder::encode(const QString &str) (DecodedData<const QString &>
// materialized as QByteArray).
func (e *QStringEncoder) Encode(q *QString) *QByteArray {
	p, ok := e.tryHandle()
	if !ok || q == nil {
		return NewQByteArray()
	}
	n := q.Size()
	var data *C.uint16_t
	if n > 0 {
		data = (*C.uint16_t)(q.Utf16Unsafe())
	}
	out := C.qstringencoder_encode(p, data, C.ptrdiff_t(n))
	runtime.KeepAlive(q)
	return wrapQByteArray(unsafe.Pointer(out))
}

// EncodeString encodes s (UTF-8 Go string decoded to UTF-16) and returns a
// fresh *QByteArray.
//
// Qt: QStringEncoder::encode(const QString &str) — Go string form.
func (e *QStringEncoder) EncodeString(s string) *QByteArray {
	q := NewQStringFromString(s)
	defer q.Close()
	return e.Encode(q)
}

// EncodeBytes encodes q and returns the result as a []byte copy.
//
// Qt: QStringEncoder::encode(const QString &str), result copied into []byte.
func (e *QStringEncoder) EncodeBytes(q *QString) []byte {
	r := e.Encode(q)
	defer r.Close()
	return r.Bytes()
}

// EncodeView encodes the view v and returns a fresh *QByteArray.
//
// Qt: QStringEncoder::encode(QStringView in) (DecodedData<QStringView>
// materialized as QByteArray).
func (e *QStringEncoder) EncodeView(v QStringView) *QByteArray {
	p, ok := e.tryHandle()
	if !ok {
		return NewQByteArray()
	}
	n := v.Size()
	var data *C.uint16_t
	if n > 0 {
		data = (*C.uint16_t)(v.Utf16Unsafe())
	}
	out := C.qstringencoder_encode(p, data, C.ptrdiff_t(n))
	runtime.KeepAlive(v)
	return wrapQByteArray(unsafe.Pointer(out))
}

// EncodeViewString encodes s as a QStringView and returns a fresh *QByteArray.
//
// Qt: QStringEncoder::encode(QStringView in) — Go string form.
func (e *QStringEncoder) EncodeViewString(s string) *QByteArray {
	q := NewQStringFromString(s)
	defer q.Close()
	return e.EncodeView(q.ToQStringView())
}

// EncodeViewBytes encodes the view v and returns the result as a []byte copy.
//
// Qt: QStringEncoder::encode(QStringView in), result copied into []byte.
func (e *QStringEncoder) EncodeViewBytes(v QStringView) []byte {
	r := e.EncodeView(v)
	defer r.Close()
	return r.Bytes()
}

// RequiredSpace returns the maximum number of bytes needed to encode
// inputLength UTF-16 code units. Returns 0 for an invalid encoder.
//
// Qt: QStringEncoder::requiredSpace(qsizetype inputLength) const.
func (e *QStringEncoder) RequiredSpace(inputLength int) int {
	p, ok := e.tryHandle()
	if !ok {
		return 0
	}
	return int(C.qstringencoder_required_space(p, C.ptrdiff_t(inputLength)))
}

// AppendToBuffer writes the encoded bytes of v into buf and returns the number
// of bytes written. The caller must ensure len(buf) >= RequiredSpace(v.Size()).
// For an invalid encoder nothing is written, HasError becomes true, and 0 is
// returned.
//
// Qt: QStringEncoder::appendToBuffer(char *out, QStringView in).
func (e *QStringEncoder) AppendToBuffer(buf []byte, v QStringView) int {
	p, ok := e.tryHandle()
	if !ok {
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
	written := C.qstringencoder_append_to_buffer(p, bufPtr, C.ptrdiff_t(len(buf)),
		data, C.ptrdiff_t(n))
	runtime.KeepAlive(buf)
	runtime.KeepAlive(v)
	return int(written)
}

// AppendToBufferString encodes s and writes the bytes into buf, returning the
// number of bytes written.
//
// Qt: QStringEncoder::appendToBuffer(char *out, QStringView in) — Go string form.
func (e *QStringEncoder) AppendToBufferString(buf []byte, s string) int {
	q := NewQStringFromString(s)
	defer q.Close()
	return e.AppendToBuffer(buf, q.ToQStringView())
}

// FinalizeBuffer flushes pending encoder state into buf (capacity len(buf)).
// Returns (next int, invalid int16, err FinalizeResultError), where next is the
// number of bytes written, invalid is the number of invalid characters seen
// during this call, and err is NoError, InvalidCharacters or NotEnoughSpace.
//
// Qt: QStringEncoder::finalize(char *out, qsizetype maxlen).
func (e *QStringEncoder) FinalizeBuffer(buf []byte) (int, int16, FinalizeResultError) {
	p, ok := e.tryHandle()
	if !ok {
		return 0, 0, FinalizeNotEnoughSpace
	}
	var bufPtr unsafe.Pointer
	if len(buf) > 0 {
		bufPtr = unsafe.Pointer(&buf[0])
	}
	var next C.ptrdiff_t
	var invalid C.int16_t
	var errCode C.int
	C.qstringencoder_finalize(p, bufPtr, C.ptrdiff_t(len(buf)), &next, &invalid, &errCode)
	runtime.KeepAlive(buf)
	return int(next), int16(invalid), FinalizeResultError(errCode)
}

// Finalize flushes pending encoder state without providing a buffer. Returns
// (invalid int16, err FinalizeResultError).
//
// Qt: QStringEncoder::finalize().
func (e *QStringEncoder) Finalize() (int16, FinalizeResultError) {
	p, ok := e.tryHandle()
	if !ok {
		return 0, FinalizeNotEnoughSpace
	}
	var next C.ptrdiff_t
	var invalid C.int16_t
	var errCode C.int
	C.qstringencoder_finalize(p, nil, 0, &next, &invalid, &errCode)
	return int16(invalid), FinalizeResultError(errCode)
}
