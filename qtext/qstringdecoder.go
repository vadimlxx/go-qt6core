package qtext

/*
QStringDecoder — конкретный декодер Qt 6.11 QtCore (байты -> UTF-16).

Иерархия воспроизведена встраиванием, как в Qt:

	QStringDecoder -> QStringConverter -> QStringConverterBase

*QStringDecoder владеет C++ объектом: он освобождается Close() или
runtime.AddCleanup (когда Go-обёртка становится недостижимой). Встроенная
*QStringConverter — не владеющая проекция на тот же объект и тот же qscState:
после Close()/cleanup конвертер видит nil, а не висячий указатель.

C++-мост общий для QStringConverter/QStringEncoder/QStringDecoder
(qstringconverter_c.cpp/.h), как qstringconverter.cpp/.h в самом QtCore.

Go-варианты:

  - DecodeBytes / DecodeViewBytes / AppendToBufferBytes принимают []byte вместо
    QByteArray/QByteArrayView;
  - DecodeString / DecodeViewString возвращают string вместо *QString;
  - NewQStringDecoderForHtmlBytes принимает []byte.

Не реализовано:

  - operator() — в Go нет перегрузки операторов; используйте Decode/DecodeView.
  - EncodedData<T> как публичный тип — не экспонируется; Decode сразу
    материализует QString.
  - QStringBuilder-интеграция (QConcatenable, operator+=) — используйте
    AppendQString/Concat.
  - FinalizeResultChar<Char> — шаблон; QChar и char16_t перегрузки в Go
    совпадают (QChar == uint16), результат — тройка (next, invalid, err).
  - Перегрузка finalize(QChar*) отдельно не нужна: используется char16_t*.
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

// QStringDecoder mirrors Qt 6.11's QStringDecoder.
//
// It embeds *QStringConverter (which embeds *QStringConverterBase), matching
// Qt's class hierarchy. The embedded converter is a non-owning view over the
// same C++ object and the same qscState; the *QStringDecoder wrapper owns the
// object.
//
// Qt: QStringDecoder.
type QStringDecoder struct {
	*QStringConverter
	state   *qscState
	cleanup runtime.Cleanup
}

// newQStringDecoderFromC takes ownership of a C++ QStringDecoder object, builds
// the embedded non-owning converter view over the same state, and registers a
// runtime cleanup for the decoder.
//
// Qt: no direct counterpart — internal constructor helper.
func newQStringDecoderFromC(c *C.go_qsd_t) *QStringDecoder {
	if c == nil {
		return nil
	}
	st := &qscState{}
	st.p.Store((*byte)(unsafe.Pointer(c)))
	dec := &QStringDecoder{
		QStringConverter: newQStringConverterAlias(st),
		state:            st,
	}
	dec.cleanup = runtime.AddCleanup(dec, func(st *qscState) {
		if p := st.claim(); p != nil {
			C.qstringdecoder_delete((*C.go_qsd_t)(p))
		}
	}, st)
	return dec
}

// tryHandle returns the underlying C++ QStringDecoder pointer and whether it is
// still alive. It does not panic: closed and nil decoders are reported as false.
//
// Qt: no direct counterpart — internal handle helper.
func (d *QStringDecoder) tryHandle() (*C.go_qsd_t, bool) {
	if d == nil || d.state == nil {
		return nil, false
	}
	if p := d.state.cpp(); p != nil {
		return (*C.go_qsd_t)(p), true
	}
	return nil, false
}

// NewQStringDecoder creates an invalid decoder (null iface, FlagDefault).
//
// Qt: QStringDecoder::QStringDecoder().
func NewQStringDecoder() *QStringDecoder {
	return newQStringDecoderFromC(C.qstringdecoder_new_default())
}

// NewQStringDecoderFromEncoding creates a decoder for the given encoding and
// flags.
//
// Qt: QStringDecoder::QStringDecoder(Encoding e, Flags f).
func NewQStringDecoderFromEncoding(e Encoding, f Flags) *QStringDecoder {
	return newQStringDecoderFromC(C.qstringdecoder_new_from_encoding(C.int(e), C.int(f)))
}

// NewQStringDecoderFromName creates a decoder by encoding name (for example
// "UTF-8" or "ISO-8859-1"). The Go string is passed as UTF-8 QAnyStringView.
// Returns an invalid decoder if the name is unknown.
//
// Qt: QStringDecoder::QStringDecoder(QAnyStringView name, Flags f) — Go string form.
func NewQStringDecoderFromName(name string, f Flags) *QStringDecoder {
	var p unsafe.Pointer
	if len(name) > 0 {
		p = unsafe.Pointer(unsafe.StringData(name))
	}
	dec := newQStringDecoderFromC(C.qstringdecoder_new_from_any(
		p, C.ptrdiff_t(len(name)), 0, C.int(f)))
	runtime.KeepAlive(name)
	return dec
}

// NewQStringDecoderFromAny creates a decoder from a QAnyStringView name and
// flags.
//
// Qt: QStringDecoder::QStringDecoder(QAnyStringView name, Flags f) — QAnyStringView form.
func NewQStringDecoderFromAny(name QAnyStringView, f Flags) *QStringDecoder {
	dec := newQStringDecoderFromC(C.qstringdecoder_new_from_any(
		name.DataUnsafe(), C.ptrdiff_t(name.Size()), C.int(name.Kind()), C.int(f)))
	runtime.KeepAlive(name)
	return dec
}

// NewQStringDecoderFromInterface creates a decoder from an existing Interface
// pointer (the protected constructor is exposed). If i is null, the decoder is
// invalid.
//
// Qt: QStringDecoder::QStringDecoder(const Interface *i) (protected).
func NewQStringDecoderFromInterface(i Interface) *QStringDecoder {
	return newQStringDecoderFromC(C.qstringdecoder_new_from_iface(i.UnsafePtr()))
}

// Close frees the underlying C++ object immediately and cancels the cleanup.
// Close is optional, idempotent and nil-safe. The embedded converter view
// shares the state slot, so after Close its methods observe a nil handle.
//
// Qt: QStringDecoder::~QStringDecoder().
func (d *QStringDecoder) Close() {
	if d == nil || d.state == nil {
		return
	}
	d.cleanup.Stop()
	if p := d.state.claim(); p != nil {
		C.qstringdecoder_delete((*C.go_qsd_t)(p))
	}
}

// Decode decodes ba and returns a fresh *QString. For an invalid decoder the
// result is empty and HasError becomes true.
//
// Qt: QStringDecoder::decode(const QByteArray &ba) (EncodedData<const QByteArray &>
// materialized as QString).
func (d *QStringDecoder) Decode(ba *QByteArray) *QString {
	p, ok := d.tryHandle()
	if !ok || ba == nil {
		return NewQString()
	}
	n := ba.Size()
	var data *C.char
	if n > 0 {
		data = (*C.char)(ba.ConstDataUnsafe())
	}
	var out unsafe.Pointer
	C.qstringdecoder_decode(p, data, C.ptrdiff_t(n), &out)
	runtime.KeepAlive(ba)
	return wrapQString(out)
}

// DecodeBytes decodes b and returns a fresh *QString.
//
// Qt: QStringDecoder::decode(const QByteArray &ba) — []byte form.
func (d *QStringDecoder) DecodeBytes(b []byte) *QString {
	ba := NewQByteArrayFromBytes(b)
	defer ba.Close()
	return d.Decode(ba)
}

// DecodeString decodes ba and returns the result as a Go string.
//
// Qt: QStringDecoder::decode(const QByteArray &ba) — Go string result.
func (d *QStringDecoder) DecodeString(ba *QByteArray) string {
	r := d.Decode(ba)
	defer r.Close()
	return r.GoString()
}

// DecodeView decodes the view v and returns a fresh *QString.
//
// Qt: QStringDecoder::decode(QByteArrayView ba) (EncodedData<QByteArrayView>
// materialized as QString).
func (d *QStringDecoder) DecodeView(v QByteArrayView) *QString {
	p, ok := d.tryHandle()
	if !ok {
		return NewQString()
	}
	n := v.Size()
	var data *C.char
	if n > 0 {
		data = (*C.char)(v.DataUnsafe())
	}
	var out unsafe.Pointer
	C.qstringdecoder_decode(p, data, C.ptrdiff_t(n), &out)
	runtime.KeepAlive(v)
	return wrapQString(out)
}

// DecodeViewBytes decodes b as a QByteArrayView and returns a fresh *QString.
//
// Qt: QStringDecoder::decode(QByteArrayView ba) — []byte form.
func (d *QStringDecoder) DecodeViewBytes(b []byte) *QString {
	v := NewQByteArrayViewFromBytes(b)
	r := d.DecodeView(v)
	runtime.KeepAlive(b)
	return r
}

// DecodeViewString decodes the view v and returns the result as a Go string.
//
// Qt: QStringDecoder::decode(QByteArrayView ba) — Go string result.
func (d *QStringDecoder) DecodeViewString(v QByteArrayView) string {
	r := d.DecodeView(v)
	defer r.Close()
	return r.GoString()
}

// RequiredSpace returns the maximum number of UTF-16 code units needed to
// decode inputLength bytes. Returns 0 for an invalid decoder.
//
// Qt: QStringDecoder::requiredSpace(qsizetype inputLength) const.
func (d *QStringDecoder) RequiredSpace(inputLength int) int {
	p, ok := d.tryHandle()
	if !ok {
		return 0
	}
	return int(C.qstringdecoder_required_space(p, C.ptrdiff_t(inputLength)))
}

// AppendToBuffer decodes v into buf and returns the number of UTF-16 code units
// written. The caller must ensure len(buf) >= RequiredSpace(v.Size()). For an
// invalid decoder nothing is written, HasError becomes true, and 0 is returned.
//
// Qt: QStringDecoder::appendToBuffer(char16_t *out, QByteArrayView ba).
func (d *QStringDecoder) AppendToBuffer(buf []uint16, v QByteArrayView) int {
	p, ok := d.tryHandle()
	if !ok {
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
	written := C.qstringdecoder_append_to_buffer(p, bufPtr, C.ptrdiff_t(len(buf)),
		data, C.ptrdiff_t(n))
	runtime.KeepAlive(buf)
	runtime.KeepAlive(v)
	return int(written)
}

// AppendToBufferBytes decodes b and writes the UTF-16 code units into buf,
// returning the number of code units written.
//
// Qt: QStringDecoder::appendToBuffer(char16_t *out, QByteArrayView ba) — []byte form.
func (d *QStringDecoder) AppendToBufferBytes(buf []uint16, b []byte) int {
	v := NewQByteArrayViewFromBytes(b)
	n := d.AppendToBuffer(buf, v)
	runtime.KeepAlive(b)
	return n
}

// FinalizeBuffer flushes pending decoder state into buf (capacity len(buf),
// UTF-16 code units). Returns (next int, invalid int16, err
// FinalizeResultError), where next is the number of code units written.
//
// Qt: QStringDecoder::finalize(char16_t *out, qsizetype maxlen).
func (d *QStringDecoder) FinalizeBuffer(buf []uint16) (int, int16, FinalizeResultError) {
	p, ok := d.tryHandle()
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
	C.qstringdecoder_finalize(p, bufPtr, C.ptrdiff_t(len(buf)), &next, &invalid, &errCode)
	runtime.KeepAlive(buf)
	return int(next), int16(invalid), FinalizeResultError(errCode)
}

// Finalize flushes pending decoder state without providing a buffer. Returns
// (invalid int16, err FinalizeResultError).
//
// Qt: QStringDecoder::finalize().
func (d *QStringDecoder) Finalize() (int16, FinalizeResultError) {
	p, ok := d.tryHandle()
	if !ok {
		return 0, FinalizeNotEnoughSpace
	}
	var next C.ptrdiff_t
	var invalid C.int16_t
	var errCode C.int
	C.qstringdecoder_finalize(p, nil, 0, &next, &invalid, &errCode)
	return int16(invalid), FinalizeResultError(errCode)
}

// NewQStringDecoderForHtml returns a decoder for the encoding declared in the
// HTML <meta> charset of data (or UTF-8 if none is found).
//
// Qt: QStringDecoder::decoderForHtml(QByteArrayView data) static — QByteArrayView form.
func NewQStringDecoderForHtml(data QByteArrayView) *QStringDecoder {
	n := data.Size()
	var p *C.char
	if n > 0 {
		p = (*C.char)(data.DataUnsafe())
	}
	c := C.qstringdecoder_decoder_for_html(p, C.ptrdiff_t(n))
	runtime.KeepAlive(data)
	return newQStringDecoderFromC(c)
}

// NewQStringDecoderForHtmlBytes is NewQStringDecoderForHtml accepting a []byte
// directly.
//
// Qt: QStringDecoder::decoderForHtml(QByteArrayView data) static — []byte form.
func NewQStringDecoderForHtmlBytes(data []byte) *QStringDecoder {
	v := NewQByteArrayViewFromBytes(data)
	d := NewQStringDecoderForHtml(v)
	runtime.KeepAlive(data)
	return d
}

// DecoderForHtml is a deprecated alias for NewQStringDecoderForHtml.
//
// Deprecated: use NewQStringDecoderForHtml.
//
// Qt: QStringDecoder::decoderForHtml(QByteArrayView data) static.
func DecoderForHtml(data QByteArrayView) *QStringDecoder {
	return NewQStringDecoderForHtml(data)
}
