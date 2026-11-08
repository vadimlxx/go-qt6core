package qstringdecoder

import (
	"bytes"
	"testing"

	"github.com/vadimlxx/go-qt6core/qbytearray"
	"github.com/vadimlxx/go-qt6core/qbytearrayview"
	// "github.com/vadimlxx/go-qt6core/qstring"
	"github.com/vadimlxx/go-qt6core/qstringconverter"
)

func TestDefaultInvalid(t *testing.T) {
	d := NewQStringDecoder()
	defer d.Close()

	if d.IsValid() {
		t.Fatal("default QStringDecoder must be invalid")
	}
	if d.Name() != "" {
		t.Fatalf("Name: %q", d.Name())
	}

	// Decoding on an invalid decoder must return empty and set HasError.
	ba := qbytearray.FromString("hello")
	defer ba.Close()
	out := d.Decode(ba)
	defer out.Close()
	if out.Size() != 0 {
		t.Fatalf("Decode on invalid: size=%d, want 0", out.Size())
	}
	if !d.HasError() {
		t.Fatal("HasError must be true after decoding with invalid decoder")
	}
}

func TestFromEncodingUtf8(t *testing.T) {
	d := NewQStringDecoderFromEncoding(qstringconverter.Utf8, qstringconverter.FlagDefault)
	defer d.Close()

	if !d.IsValid() {
		t.Fatal("UTF-8 decoder must be valid")
	}
	if d.Name() != "UTF-8" {
		t.Fatalf("Name: %q", d.Name())
	}

	ba := qbytearray.FromString("héllo")
	defer ba.Close()
	out := d.Decode(ba)
	defer out.Close()

	if out.GoString() != "héllo" {
		t.Fatalf("Decode: %q, want %q", out.GoString(), "héllo")
	}
}

func TestFromNameUtf8(t *testing.T) {
	d := NewQStringDecoderFromName("UTF-8", qstringconverter.FlagDefault)
	defer d.Close()
	if !d.IsValid() || d.Name() != "UTF-8" {
		t.Fatalf("valid=%v name=%q", d.IsValid(), d.Name())
	}

	bad := NewQStringDecoderFromName("nonexistent", qstringconverter.FlagDefault)
	defer bad.Close()
	if bad.IsValid() {
		t.Fatal("nonexistent name must be invalid")
	}
}

func TestDecodeView(t *testing.T) {
	d := NewQStringDecoderFromEncoding(qstringconverter.Utf8, qstringconverter.FlagDefault)
	defer d.Close()

	bv := qbytearrayview.FromString("日本")
	out := d.DecodeView(bv)
	defer out.Close()

	if out.GoString() != "日本" {
		t.Fatalf("DecodeView: %q", out.GoString())
	}
}

func TestDecodeLatin1(t *testing.T) {
	d := NewQStringDecoderFromEncoding(qstringconverter.Latin1, qstringconverter.FlagDefault)
	defer d.Close()

	// "héllo" as Latin-1 bytes: h=0x68, é=0xE9, l, l, o
	ba := qbytearray.FromBytes([]byte{'h', 0xE9, 'l', 'l', 'o'})
	defer ba.Close()
	out := d.Decode(ba)
	defer out.Close()

	if out.GoString() != "héllo" {
		t.Fatalf("Decode(Latin1): %q", out.GoString())
	}
}

func TestDecodeUtf16LE(t *testing.T) {
	d := NewQStringDecoderFromEncoding(qstringconverter.Utf16LE, qstringconverter.FlagDefault)
	defer d.Close()

	// "hi" in UTF-16LE: 68 00 69 00
	ba := qbytearray.FromBytes([]byte{0x68, 0x00, 0x69, 0x00})
	defer ba.Close()
	out := d.Decode(ba)
	defer out.Close()

	if out.GoString() != "hi" {
		t.Fatalf("Decode(Utf16LE): %q", out.GoString())
	}
}

func TestRequiredSpace(t *testing.T) {
	d := NewQStringDecoderFromEncoding(qstringconverter.Utf8, qstringconverter.FlagDefault)
	defer d.Close()

	// RequiredSpace must be >= input length for any encoding.
	n := d.RequiredSpace(10)
	if n < 10 {
		t.Fatalf("RequiredSpace(10) = %d, want >= 10", n)
	}
}

func TestAppendToBuffer(t *testing.T) {
	d := NewQStringDecoderFromEncoding(qstringconverter.Utf8, qstringconverter.FlagDefault)
	defer d.Close()

	bv := qbytearrayview.FromString("hello")
	buf := make([]uint16, d.RequiredSpace(bv.Size()))
	n := d.AppendToBuffer(buf, bv)
	if n != 5 {
		t.Fatalf("AppendToBuffer: n=%d, want 5", n)
	}
	if string(rune16ToRunes(buf[:n])) != "hello" {
		t.Fatalf("AppendToBuffer: %q", string(rune16ToRunes(buf[:n])))
	}
}

// rune16ToRunes converts UTF-16 code units to []rune for test comparison.
func rune16ToRunes(u []uint16) []rune {
	out := make([]rune, 0, len(u))
	for i := 0; i < len(u); i++ {
		c := u[i]
		if c >= 0xD800 && c < 0xDC00 && i+1 < len(u) {
			lo := u[i+1]
			if lo >= 0xDC00 && lo < 0xE000 {
				out = append(out, rune(c-0xD800)<<10|rune(lo-0xDC00)+0x10000)
				i++
				continue
			}
		}
		out = append(out, rune(c))
	}
	return out
}

func TestFinalizeEmpty(t *testing.T) {
	d := NewQStringDecoderFromEncoding(qstringconverter.Utf8, qstringconverter.FlagDefault)
	defer d.Close()

	// Fresh UTF-8 decoder has no pending state; Finalize should report NoError.
	invalid, errCode := d.Finalize()
	if errCode != FinalizeNoError {
		t.Fatalf("Finalize: err=%v, want NoError", errCode)
	}
	if invalid != 0 {
		t.Fatalf("Finalize: invalid=%d, want 0", invalid)
	}
}

func TestFinalizeBufferEmpty(t *testing.T) {
	d := NewQStringDecoderFromEncoding(qstringconverter.Utf8, qstringconverter.FlagDefault)
	defer d.Close()

	next, invalid, errCode := d.FinalizeBuffer(make([]uint16, 16))
	if errCode != FinalizeNoError || next != 0 || invalid != 0 {
		t.Fatalf("FinalizeBuffer: next=%d invalid=%d err=%v", next, invalid, errCode)
	}
}

func TestEmbeddedBase(t *testing.T) {
	d := NewQStringDecoderFromEncoding(qstringconverter.Utf8, qstringconverter.FlagDefault)
	defer d.Close()

	// The embedded *QStringConverter must not be nil and must reflect the same C++ state.
	if d.QStringConverter == nil {
		t.Fatal("embedded *QStringConverter is nil")
	}
	if !d.QStringConverter.IsValid() {
		t.Fatal("embedded converter must be valid")
	}
	if d.QStringConverter.Name() != "UTF-8" {
		t.Fatalf("embedded Name: %q", d.QStringConverter.Name())
	}

	// Promoted base methods work directly on *QStringDecoder:
	if !d.IsValid() {
		t.Fatal("promoted IsValid must be true")
	}
	if d.Name() != "UTF-8" {
		t.Fatalf("promoted Name: %q", d.Name())
	}
}

func TestFromInterface(t *testing.T) {
	// Build a valid converter to get its Interface, then construct a decoder from it.
	src := qstringconverter.NewQStringConverterFromEncoding(qstringconverter.Utf8, qstringconverter.FlagDefault)
	defer src.Close()

	iface := src.Iface()
	if iface.IsNil() {
		t.Fatal("Iface must be non-nil for valid converter")
	}

	d := NewQStringDecoderFromInterface(iface)
	defer d.Close()
	if !d.IsValid() || d.Name() != "UTF-8" {
		t.Fatalf("FromInterface: valid=%v name=%q", d.IsValid(), d.Name())
	}

	// Null interface → invalid decoder.
	nullDec := NewQStringDecoderFromInterface(qstringconverter.Interface{})
	defer nullDec.Close()
	if nullDec.IsValid() {
		t.Fatal("FromInterface(null) must be invalid")
	}
}

func TestDecodeMultiByte(t *testing.T) {
	d := NewQStringDecoderFromEncoding(qstringconverter.Utf8, qstringconverter.FlagDefault)
	defer d.Close()

	// "a😀b" — 😀 = U+1F600 = 4 UTF-8 bytes, 2 UTF-16 code units
	ba := qbytearray.FromString("a😀b")
	defer ba.Close()

	out := d.Decode(ba)
	defer out.Close()

	if out.GoString() != "a😀b" {
		t.Fatalf("Decode: %q, want %q", out.GoString(), "a😀b")
	}
}

func TestDecodeReusesState(t *testing.T) {
	// Multiple decode calls reuse the same decoder instance and its state.
	d := NewQStringDecoderFromEncoding(qstringconverter.Utf8, qstringconverter.FlagDefault)
	defer d.Close()

	for _, s := range []string{"hello", "world", "!"} {
		ba := qbytearray.FromString(s)
		out := d.Decode(ba)
		ba.Close()
		if out.GoString() != s {
			t.Fatalf("Decode(%q): %q", s, out.GoString())
		}
		out.Close()
	}
}

func TestDecoderForHtml(t *testing.T) {
	html := "<html><head><meta charset=\"UTF-8\"><title>x</title></head></html>"
	bv := qbytearrayview.FromString(html)
	d := DecoderForHtml(bv)
	defer d.Close()
	if !d.IsValid() {
		t.Fatal("DecoderForHtml must return a valid decoder for UTF-8 HTML")
	}

	// Sanity: the returned decoder is usable.
	ba := qbytearray.FromString("hello")
	defer ba.Close()
	out := d.Decode(ba)
	defer out.Close()
	if out.GoString() != "hello" {
		t.Fatalf("Decode via DecoderForHtml: %q", out.GoString())
	}
}

// sanity check that bytes.Equal is used (avoid unused import in future edits).
var _ = bytes.Equal
