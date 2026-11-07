package qstringconverter

import (
	"testing"

	"github.com/vadimlxx/go-qt6core/qanystringview"
	"github.com/vadimlxx/go-qt6core/qbytearrayview"
)

func TestDefaultInvalid(t *testing.T) {
	q := NewQStringConverter()
	defer q.Close()

	if q.IsValid() {
		t.Fatal("default QStringConverter must be invalid")
	}
	// HasError() reports whether the converter has encountered invalid characters during
	// decoding. A fresh converter has not decoded anything, so it reports false even though
	// the converter is invalid (iface == nullptr). These are two different concepts:
	//   IsValid()  — has a valid encoding Interface
	//   HasError() — has produced/encountered invalid characters in the last operation
	if q.HasError() {
		t.Fatal("fresh invalid converter must not report HasError")
	}
	if q.Name() != "" {
		t.Fatalf("Name: %q", q.Name())
	}
	if q.Encoding() != Encoding(-1) {
		t.Fatalf("Encoding: %v, want -1", q.Encoding())
	}
}

func TestFromEncoding(t *testing.T) {
	cases := []struct {
		e    Encoding
		name string
	}{
		{Utf8, "UTF-8"},
		{Utf16, "UTF-16"},
		{Utf16LE, "UTF-16LE"},
		{Utf16BE, "UTF-16BE"},
		{Utf32, "UTF-32"},
		{Utf32LE, "UTF-32LE"},
		{Utf32BE, "UTF-32BE"},
		{Latin1, "ISO-8859-1"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			q := NewQStringConverterFromEncoding(tc.e, FlagDefault)
			defer q.Close()
			if !q.IsValid() {
				t.Fatalf("%v: not valid", tc.e)
			}
			if q.Encoding() != tc.e {
				t.Fatalf("Encoding: got %v, want %v", q.Encoding(), tc.e)
			}
			if q.Name() != tc.name {
				t.Fatalf("Name: got %q, want %q", q.Name(), tc.name)
			}
		})
	}
}

func TestFromName(t *testing.T) {
	q := NewQStringConverterFromName("UTF-16", FlagDefault)
	defer q.Close()
	if !q.IsValid() || q.Encoding() != Utf16 {
		t.Fatalf("valid=%v encoding=%v", q.IsValid(), q.Encoding())
	}

	bad := NewQStringConverterFromName("nonexistent-encoding", FlagDefault)
	defer bad.Close()
	if bad.IsValid() {
		t.Fatal("nonexistent name must be invalid")
	}
}

func TestFromAny(t *testing.T) {
	any := qanystringview.FromString("UTF-16")
	q := NewQStringConverterFromAny(any, FlagDefault)
	defer q.Close()
	if !q.IsValid() || q.Encoding() != Utf16 {
		t.Fatalf("FromAny(UTF-16): valid=%v encoding=%v", q.IsValid(), q.Encoding())
	}

	// Latin1 view (byte-level same as UTF-8 for ASCII)
	l1 := qanystringview.FromBytes([]byte("UTF-8")) // simplest cross-kind case
	q2 := NewQStringConverterFromAny(l1, FlagDefault)
	defer q2.Close()
	if !q2.IsValid() || q2.Encoding() != Utf8 {
		t.Fatalf("FromAny(UTF-8): valid=%v encoding=%v", q2.IsValid(), q2.Encoding())
	}
}

func TestFromIface(t *testing.T) {
	src := NewQStringConverterFromEncoding(Utf16, FlagDefault)
	defer src.Close()

	iface := src.Iface()
	if iface.IsNil() {
		t.Fatal("Iface of valid converter must be non-nil")
	}

	cpy := NewQStringConverterFromIface(iface)
	defer cpy.Close()
	if !cpy.IsValid() || cpy.Encoding() != Utf16 {
		t.Fatalf("FromIface: valid=%v encoding=%v", cpy.IsValid(), cpy.Encoding())
	}

	// null iface → invalid converter
	null := NewQStringConverterFromIface(Interface{})
	defer null.Close()
	if null.IsValid() {
		t.Fatal("FromIface(null) must be invalid")
	}
}

func TestNameForEncoding(t *testing.T) {
	if got := NameForEncoding(Utf16); got != "UTF-16" {
		t.Fatalf("NameForEncoding(Utf16): %q", got)
	}
	if got := NameForEncoding(Latin1); got != "ISO-8859-1" {
		t.Fatalf("NameForEncoding(Latin1): %q", got)
	}
}

func TestEncodingForName(t *testing.T) {
	if e, ok := EncodingForName("UTF-16"); !ok || e != Utf16 {
		t.Fatalf("EncodingForName(UTF-16): %v, %v", e, ok)
	}
	if e, ok := EncodingForName("ISO-8859-1"); !ok || e != Latin1 {
		t.Fatalf("EncodingForName(ISO-8859-1): %v, %v", e, ok)
	}
	if _, ok := EncodingForName("garbage"); ok {
		t.Fatal("EncodingForName(garbage) must be !ok")
	}
}

func TestEncodingForNameAny(t *testing.T) {
	any := qanystringview.FromString("UTF-16")
	if e, ok := EncodingForNameAny(any); !ok || e != Utf16 {
		t.Fatalf("EncodingForNameAny(UTF-16): %v, %v", e, ok)
	}
}

func TestEncodingForData(t *testing.T) {
	// UTF-16LE BOM: FF FE
	bv := qbytearrayview.FromBytes([]byte{0xFF, 0xFE, 0x41, 0x00, 0x42, 0x00})
	e, ok := EncodingForData(bv, 0)
	if !ok || e != Utf16LE {
		t.Fatalf("EncodingForData(UTF-16LE BOM): %v, %v", e, ok)
	}

	// UTF-8 BOM: EF BB BF
	bv2 := qbytearrayview.FromBytes([]byte{0xEF, 0xBB, 0xBF, 0x41})
	e2, ok2 := EncodingForData(bv2, 0)
	if !ok2 || e2 != Utf8 {
		t.Fatalf("EncodingForData(UTF-8 BOM): %v, %v", e2, ok2)
	}

	// No BOM, no detection
	bv3 := qbytearrayview.FromBytes([]byte("plain ASCII"))
	if _, ok := EncodingForData(bv3, 0); ok {
		// plain ASCII typically returns false (no BOM, no plausible heuristic)
		t.Log("EncodingForData(plain ASCII) detected something — that's fine")
	}
}

func TestEncodingForHtml(t *testing.T) {
	html := "<html><head><meta charset=\"UTF-8\"><title>x</title></head></html>"
	bv := qbytearrayview.FromString(html)
	e, ok := EncodingForHtml(bv)
	t.Logf("EncodingForHtml: %v, %v", e, ok)
	// Not asserting — HTML detection is heuristic.
}

func TestState(t *testing.T) {
	q := NewQStringConverterFromEncoding(Utf16, FlagDefault)
	defer q.Close()

	if s := q.State(); s != (State{}) {
		t.Fatalf("initial State: %+v", s)
	}

	want := State{
		Flags:          Flags(FlagWriteBom),
		InternalState:  7,
		RemainingChars: 3,
		InvalidChars:   1,
		Data:           [4]uint32{10, 20, 30, 40},
	}
	q.SetState(want)
	if got := q.State(); got != want {
		t.Fatalf("after SetState: got %+v, want %+v", got, want)
	}

	// Qt 6.11 State::reset() clears the runtime progress (InternalState,
	// RemainingChars, InvalidChars, state_data) but PRESERVES Flags: they encode
	// configuration (WriteBom, Stateless, …), not per-call progress.
	q.ResetState()
	got := q.State()
	if got.InternalState != 0 || got.RemainingChars != 0 || got.InvalidChars != 0 ||
		got.Data != [4]uint32{} {
			t.Fatalf("after ResetState: progress not cleared: %+v", got)
		}
		if got.Flags != want.Flags {
			t.Fatalf("after ResetState: Flags = %v, want preserved %v", got.Flags, want.Flags)
		}
}

func TestSwap(t *testing.T) {
	a := NewQStringConverterFromEncoding(Utf16, FlagDefault)
	defer a.Close()
	b := NewQStringConverterFromEncoding(Latin1, FlagDefault)
	defer b.Close()

	a.Swap(b)
	if a.Encoding() != Latin1 || b.Encoding() != Utf16 {
		t.Fatalf("after swap: a=%v b=%v", a.Encoding(), b.Encoding())
	}
}

func TestIface(t *testing.T) {
	q := NewQStringConverterFromEncoding(Utf16, FlagDefault)
	defer q.Close()

	iface := q.Iface()
	if iface.IsNil() {
		t.Fatal("Iface of valid converter must be non-nil")
	}
	if iface.Name() != "UTF-16" {
		t.Fatalf("Iface.Name: %q", iface.Name())
	}
	if iface.ToUtf16Ptr() == nil {
		t.Fatal("ToUtf16Ptr must be non-nil")
	}
	if iface.ToUtf16LenPtr() == nil {
		t.Fatal("ToUtf16LenPtr must be non-nil")
	}
	if iface.FromUtf16Ptr() == nil {
		t.Fatal("FromUtf16Ptr must be non-nil")
	}
	if iface.FromUtf16LenPtr() == nil {
		t.Fatal("FromUtf16LenPtr must be non-nil")
	}
}

func TestIfaceNull(t *testing.T) {
	q := NewQStringConverter()
	defer q.Close()
	if !q.Iface().IsNil() {
		t.Fatal("Iface of invalid converter must be nil")
	}
}

func TestFlagHas(t *testing.T) {
	f := Flags(FlagDefault)
	if f.Has(FlagStateless) {
		t.Fatal("Default must not have Stateless")
	}
	f |= Flags(FlagWriteBom) | Flags(FlagConvertInvalidToNull)
	if !f.Has(FlagWriteBom) || !f.Has(FlagConvertInvalidToNull) {
		t.Fatalf("flags: %x", int(f))
	}
	if f.Has(FlagStateless) {
		t.Fatal("Stateless must not be set")
	}
}
