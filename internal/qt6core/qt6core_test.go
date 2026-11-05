package qt6core

import (
	"bytes"
	"strings"
	"testing"
	"unicode/utf16"
	"unsafe"
)

// mkQString is a helper: builds a QString from s and registers cleanup.
func mkQString(t *testing.T, s string) unsafe.Pointer {
	t.Helper()
	p := QStringFromUtf8(s)
	if p == nil {
		t.Fatalf("QStringFromUtf8(%q) returned nil", s)
	}
	t.Cleanup(func() { QStringDelete(p) })
	return p
}

func TestVersion(t *testing.T) {
	maj, min, pat := Version()
	t.Logf("Qt %d.%d.%d", maj, min, pat)

	if maj != 6 {
		t.Fatalf("expected Qt 6.x, got %d.%d.%d", maj, min, pat)
	}
	if min < 11 {
		t.Fatalf("expected Qt >= 6.11, got %d.%d.%d", maj, min, pat)
	}
}

func TestUtf8RoundTrip(t *testing.T) {
	cases := []struct {
		name string
		in   string
	}{
		{"empty", ""},
		{"ascii", "hello, world"},
		{"latin1", "café résumé"},
		{"cyrillic", "Привет, мир"},
		{"cjk", "日本語テスト"},
		{"emoji", "a😀b🎉c"},          // non-BMP → surrogate pairs
		{"embedded_nul", "a\x00b\x00c"}, // QString holds embedded NULs
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			q := mkQString(t, tc.in)

			// Both extractors must agree and round-trip.
			asString := QStringToUtf8String(q)
			if asString != tc.in {
				t.Fatalf("QStringToUtf8String: in=%q out=%q", tc.in, asString)
			}

			asBytes := QStringToUtf8(q)
			if !bytes.Equal(asBytes, []byte(tc.in)) {
				t.Fatalf("QStringToUtf8: in=%q out=%q", tc.in, asBytes)
			}
		})
	}
}

func TestUtf8FromBytes(t *testing.T) {
	in := []byte("héllo, 世界")
	q := QStringFromUtf8Bytes(in)
	if q == nil {
		t.Fatal("QStringFromUtf8Bytes returned nil")
	}
	defer QStringDelete(q)

	if got := QStringToUtf8(q); !bytes.Equal(got, in) {
		t.Fatalf("in=%v out=%v", in, got)
	}
}

func TestUtf16RoundTrip(t *testing.T) {
	// "a😀b": 😀 = U+1F600 needs a surrogate pair → exercises the
	// reinterpret_cast<char16_t*> path in bridge.cpp.
	s := "a😀b"
	want := utf16.Encode([]rune(s))

	q := QStringFromUtf16(want)
	if q == nil {
		t.Fatal("QStringFromUtf16 returned nil")
	}
	defer QStringDelete(q)

	got := QStringToUtf16(q)
	if len(got) != len(want) {
		t.Fatalf("length: want %d, got %d", len(want), len(got))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("unit[%d]: want %#04x, got %#04x", i, want[i], got[i])
		}
	}
}

func TestLatin1RoundTrip(t *testing.T) {
	// All bytes below 0x100 are valid Latin-1.
	in := []byte("caf\xe9 na\xefve") // "café naïve" in Latin-1
	q := QStringFromLatin1(in)
	if q == nil {
		t.Fatal("QStringFromLatin1 returned nil")
	}
	defer QStringDelete(q)

	if got := QStringToLatin1(q); !bytes.Equal(got, in) {
		t.Fatalf("in=%v out=%v", in, got)
	}

	// Non-Latin-1 chars collapse to '?'.
	q2 := mkQString(t, "日本")
	out := QStringToLatin1(q2)
	if !strings.Contains(string(out), "?") {
		t.Fatalf("expected '?' substitution, got %q", string(out))
	}
}

func TestQStringQByteArrayInterop(t *testing.T) {
	s := "héllo, 世界"

	q := mkQString(t, s)

	ba := QByteArrayFromQString(q)
	if ba == nil {
		t.Fatal("QByteArrayFromQString returned nil")
	}
	defer QByteArrayDelete(ba)

	q2 := QStringFromQByteArray(ba)
	if q2 == nil {
		t.Fatal("QStringFromQByteArray returned nil")
	}
	defer QStringDelete(q2)

	if got := QStringToUtf8String(q2); got != s {
		t.Fatalf("round-trip: in=%q out=%q", s, got)
	}
}

func TestNilSafety(t *testing.T) {
	// Destructors on nil must be no-ops.
	QStringDelete(nil)
	QByteArrayDelete(nil)

	// Extractors on nil must return empty values.
	if got := QStringToUtf8(nil); len(got) != 0 {
		t.Fatalf("QStringToUtf8(nil) = %v, want empty", got)
	}
	if got := QStringToUtf8String(nil); got != "" {
		t.Fatalf("QStringToUtf8String(nil) = %q, want empty", got)
	}
	if got := QStringToUtf16(nil); len(got) != 0 {
		t.Fatalf("QStringToUtf16(nil) = %v, want empty", got)
	}
	if got := QStringToLatin1(nil); len(got) != 0 {
		t.Fatalf("QStringToLatin1(nil) = %v, want empty", got)
	}

	// From-conversions on nil pointer are also safe (return nil).
	if got := QStringFromQByteArray(nil); got != nil {
		t.Fatalf("QStringFromQByteArray(nil) = %p, want nil", got)
	}
	if got := QByteArrayFromQString(nil); got != nil {
		t.Fatalf("QByteArrayFromQString(nil) = %p, want nil", got)
	}
}