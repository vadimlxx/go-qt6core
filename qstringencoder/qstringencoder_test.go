package qstringencoder

import (
	"bytes"
	"testing"
	"unicode/utf16"

	"github.com/vadimlxx/go-qt6core/qbytearray"
	"github.com/vadimlxx/go-qt6core/qstring"
	"github.com/vadimlxx/go-qt6core/qstringconverter"
	"github.com/vadimlxx/go-qt6core/qstringview"
)

func TestDefaultInvalid(t *testing.T) {
	e := NewQStringEncoder()
	defer e.Close()

	if e.IsValid() {
		t.Fatal("default QStringEncoder must be invalid")
	}
	if e.Name() != "" {
		t.Fatalf("Name: %q", e.Name())
	}

	// Encoding on an invalid encoder must return empty and set HasError.
	q := qstring.FromString("hello")
	defer q.Close()
	out := e.Encode(q)
	defer out.Close()
	if out.Size() != 0 {
		t.Fatalf("Encode on invalid: size=%d, want 0", out.Size())
	}
	if !e.HasError() {
		t.Fatal("HasError must be true after encoding with invalid encoder")
	}
}

func TestFromEncodingUtf8(t *testing.T) {
	e := NewQStringEncoderFromEncoding(qstringconverter.Utf8, qstringconverter.FlagDefault)
	defer e.Close()

	if !e.IsValid() {
		t.Fatal("UTF-8 encoder must be valid")
	}
	if e.Name() != "UTF-8" {
		t.Fatalf("Name: %q", e.Name())
	}

	q := qstring.FromString("héllo")
	defer q.Close()
	out := e.Encode(q)
	defer out.Close()

	if !bytes.Equal(out.Bytes(), []byte("héllo")) {
		t.Fatalf("Encode: %v, want %v", out.Bytes(), []byte("héllo"))
	}
}

func TestFromNameUtf8(t *testing.T) {
	e := NewQStringEncoderFromName("UTF-8", qstringconverter.FlagDefault)
	defer e.Close()
	if !e.IsValid() || e.Name() != "UTF-8" {
		t.Fatalf("valid=%v name=%q", e.IsValid(), e.Name())
	}

	bad := NewQStringEncoderFromName("nonexistent", qstringconverter.FlagDefault)
	defer bad.Close()
	if bad.IsValid() {
		t.Fatal("nonexistent name must be invalid")
	}
}

func TestEncodeView(t *testing.T) {
	e := NewQStringEncoderFromEncoding(qstringconverter.Utf8, qstringconverter.FlagDefault)
	defer e.Close()

	sv := qstringview.FromUtf16(utf16.Encode([]rune("日本")))
	out := e.EncodeView(sv)
	defer out.Close()

	if !bytes.Equal(out.Bytes(), []byte("日本")) {
		t.Fatalf("EncodeView: %v", out.Bytes())
	}
}

func TestEncodeLatin1(t *testing.T) {
	e := NewQStringEncoderFromEncoding(qstringconverter.Latin1, qstringconverter.FlagDefault)
	defer e.Close()

	// "héllo" — Latin1 encodes é as single byte 0xE9.
	q := qstring.FromString("héllo")
	defer q.Close()
	out := e.Encode(q)
	defer out.Close()

	want := []byte{'h', 0xE9, 'l', 'l', 'o'}
	if !bytes.Equal(out.Bytes(), want) {
		t.Fatalf("Encode(Latin1): %v, want %v", out.Bytes(), want)
	}
}

func TestRequiredSpace(t *testing.T) {
	e := NewQStringEncoderFromEncoding(qstringconverter.Utf8, qstringconverter.FlagDefault)
	defer e.Close()

	// requiredSpace must be an upper bound; for UTF-8 it's typically 4*len for non-ASCII.
	n := e.RequiredSpace(10)
	if n < 10 {
		t.Fatalf("RequiredSpace(10) = %d, want >= 10", n)
	}
}

func TestAppendToBuffer(t *testing.T) {
	e := NewQStringEncoderFromEncoding(qstringconverter.Utf8, qstringconverter.FlagDefault)
	defer e.Close()

	sv := qstringview.FromUtf16(utf16.Encode([]rune("hello")))
	buf := make([]byte, e.RequiredSpace(sv.Size()))
	n := e.AppendToBuffer(buf, sv)
	if n != 5 || string(buf[:n]) != "hello" {
		t.Fatalf("AppendToBuffer: n=%d, buf=%q", n, buf[:n])
	}
}

func TestFinalizeEmpty(t *testing.T) {
	e := NewQStringEncoderFromEncoding(qstringconverter.Utf8, qstringconverter.FlagDefault)
	defer e.Close()

	// Fresh UTF-8 encoder has no pending state; Finalize should report NoError.
	invalid, errCode := e.Finalize()
	if errCode != FinalizeNoError {
		t.Fatalf("Finalize: err=%v, want NoError", errCode)
	}
	if invalid != 0 {
		t.Fatalf("Finalize: invalid=%d, want 0", invalid)
	}
}

func TestFinalizeBufferEmpty(t *testing.T) {
	e := NewQStringEncoderFromEncoding(qstringconverter.Utf8, qstringconverter.FlagDefault)
	defer e.Close()

	next, invalid, errCode := e.FinalizeBuffer(make([]byte, 16))
	if errCode != FinalizeNoError || next != 0 || invalid != 0 {
		t.Fatalf("FinalizeBuffer: next=%d invalid=%d err=%v", next, invalid, errCode)
	}
}

func TestEmbeddedBase(t *testing.T) {
	e := NewQStringEncoderFromEncoding(qstringconverter.Utf8, qstringconverter.FlagDefault)
	defer e.Close()

	// The embedded *QStringConverter must not be nil and must reflect the same C++ state.
	if e.QStringConverter == nil {
		t.Fatal("embedded *QStringConverter is nil")
	}
	if !e.QStringConverter.IsValid() {
		t.Fatal("embedded converter must be valid")
	}
	if e.QStringConverter.Name() != "UTF-8" {
		t.Fatalf("embedded Name: %q", e.QStringConverter.Name())
	}

	// Promoted base methods work directly on *QStringEncoder:
	if !e.IsValid() {
		t.Fatal("promoted IsValid must be true")
	}
	if e.Name() != "UTF-8" {
		t.Fatalf("promoted Name: %q", e.Name())
	}
}

func TestFromInterface(t *testing.T) {
	// Build a valid converter to get its Interface, then construct an encoder from it.
	src := qstringconverter.NewQStringConverterFromEncoding(qstringconverter.Utf8, qstringconverter.FlagDefault)
	defer src.Close()

	iface := src.Iface()
	if iface.IsNil() {
		t.Fatal("Iface must be non-nil for valid converter")
	}

	e := NewQStringEncoderFromInterface(iface)
	defer e.Close()
	if !e.IsValid() || e.Name() != "UTF-8" {
		t.Fatalf("FromInterface: valid=%v name=%q", e.IsValid(), e.Name())
	}

	// Null interface → invalid encoder.
	nullEnc := NewQStringEncoderFromInterface(qstringconverter.Interface{})
	defer nullEnc.Close()
	if nullEnc.IsValid() {
		t.Fatal("FromInterface(null) must be invalid")
	}
}

func TestEncodeMultipleChars(t *testing.T) {
	e := NewQStringEncoderFromEncoding(qstringconverter.Utf8, qstringconverter.FlagDefault)
	defer e.Close()

	q := qstring.FromString("a😀b") // 😀 = U+1F600 = 2 UTF-16 units, 4 UTF-8 bytes
	defer q.Close()

	out := e.Encode(q)
	defer out.Close()

	if !bytes.Equal(out.Bytes(), []byte("a😀b")) {
		t.Fatalf("Encode: %v", out.Bytes())
	}
}

func TestEncodeThenEncodeAgain(t *testing.T) {
	// Multiple encode calls reuse the same encoder instance and its state.
	e := NewQStringEncoderFromEncoding(qstringconverter.Utf8, qstringconverter.FlagDefault)
	defer e.Close()

	for _, s := range []string{"hello", "world", "!"} {
		q := qstring.FromString(s)
		out := e.Encode(q)
		q.Close()
		if !bytes.Equal(out.Bytes(), []byte(s)) {
			t.Fatalf("Encode(%q): %v", s, out.Bytes())
		}
		out.Close()
	}
}

// helper alias to avoid unused import warnings in trimmed builds.
var _ = qbytearray.AdoptUnsafePtr
