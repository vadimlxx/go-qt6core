package qutf8stringview

import (
	"bytes"
	"testing"
	"unicode/utf16"

	"github.com/vadimlxx/go-qt6core/qbytearray"
	"github.com/vadimlxx/go-qt6core/qchar"
	"github.com/vadimlxx/go-qt6core/qlatin1stringview"
	"github.com/vadimlxx/go-qt6core/qstring"
	"github.com/vadimlxx/go-qt6core/qstringview"
)

func fromStr(s string) QUtf8StringView { return FromString(s) }

func TestEmpty(t *testing.T) {
	v := Empty()
	if !v.IsNull() || !v.IsEmpty() || v.Size() != 0 {
		t.Fatalf("Empty: null=%v empty=%v size=%d", v.IsNull(), v.IsEmpty(), v.Size())
	}
	if v.Bytes() != nil || v.String() != "" {
		t.Fatal("Empty content")
	}
}

func TestFromString(t *testing.T) {
	v := fromStr("hello")
	if v.Size() != 5 || v.IsNull() {
		t.Fatalf("Size=%d IsNull=%v", v.Size(), v.IsNull())
	}
	if v.String() != "hello" {
		t.Fatalf("String: %q", v.String())
	}
	if !v.IsValidUtf8() {
		t.Fatal("hello must be valid UTF-8")
	}
}

func TestFromBytes(t *testing.T) {
	b := []byte("héllo") // UTF-8
	v := FromBytes(b)
	if v.String() != "héllo" {
		t.Fatalf("String: %q", v.String())
	}
	if !v.IsValidUtf8() {
		t.Fatal("héllo must be valid UTF-8")
	}

	// Invalid UTF-8 sequence.
	bad := FromBytes([]byte{0xFF, 0xFE, 0xFD})
	if bad.IsValidUtf8() {
		t.Fatal("expected invalid UTF-8")
	}
}

func TestFromQByteArray(t *testing.T) {
	q := qbytearray.FromString("hello")
	defer q.Close()
	v := FromQByteArray(q)
	if v.String() != "hello" {
		t.Fatalf("FromQByteArray: %q", v.String())
	}
	if !FromQByteArray(nil).IsNull() {
		t.Fatal("FromQByteArray(nil) must be null")
	}
}

func TestFromQString(t *testing.T) {
	s := qstring.FromString("héllo")
	defer s.Close()
	v := FromQString(s)
	if v.String() != "héllo" {
		t.Fatalf("FromQString: %q", v.String())
	}
	if !v.IsValidUtf8() {
		t.Fatal("valid UTF-8 expected")
	}
}

func TestAtFrontBack(t *testing.T) {
	v := fromStr("abc")
	if v.At(0) != 'a' || v.At(2) != 'c' {
		t.Fatalf("At: %c %c", v.At(0), v.At(2))
	}
	if v.Front() != 'a' || v.Back() != 'c' {
		t.Fatal("Front/Back")
	}
}

func TestSlicing(t *testing.T) {
	v := fromStr("the quick brown fox")

	check := func(name string, got QUtf8StringView, want string) {
		t.Helper()
		if got.String() != want {
			t.Fatalf("%s: got %q, want %q", name, got.String(), want)
		}
	}

	check("Left", v.Left(3), "the")
	check("Right", v.Right(3), "fox")
	check("Mid", v.Mid(4, 5), "quick")
	check("First", v.First(3), "the")
	check("Last", v.Last(3), "fox")
	check("Sliced", v.Sliced(4), "quick brown fox")
	check("SlicedN", v.SlicedN(4, 5), "quick")
	check("Chopped", v.Chopped(4), "the quick brown")
}

func TestCompareView(t *testing.T) {
	a := fromStr("hello")
	b := fromStr("hello")
	c := fromStr("world")

	if a.Compare(b, true) != 0 {
		t.Fatal("equal must compare 0")
	}
	if a.Compare(c, true) >= 0 {
		t.Fatal("hello < world")
	}
	if !a.Equal(b) {
		t.Fatal("Equal(hello, hello)")
	}
	if a.Equal(c) {
		t.Fatal("Equal(hello, world) must be false")
	}
}

func TestCompareQByteArray(t *testing.T) {
	v := fromStr("hello")
	q := qbytearray.FromString("hello")
	defer q.Close()

	if v.CompareQByteArray(q, true) != 0 {
		t.Fatal("CompareQByteArray equal")
	}
	if !v.EqualQByteArray(q) {
		t.Fatal("EqualQByteArray")
	}
}

func TestCompareQChar(t *testing.T) {
	v := fromStr("a")
	if v.CompareQChar(qchar.New('a'), true) != 0 {
		t.Fatalf("CompareQChar('a'): %d", v.CompareQChar(qchar.New('a'), true))
	}
	if !v.EqualQChar(qchar.New('a')) {
		t.Fatal("EqualQChar('a')")
	}
	if v.EqualQChar(qchar.New('b')) {
		t.Fatal("EqualQChar('b') must be false")
	}
}

func TestCompareQStringView(t *testing.T) {
	v := fromStr("hello")
	sv := qstringview.FromUtf16(utf16.Encode([]rune("hello")))

	if v.CompareQStringView(sv, true) != 0 {
		t.Fatal("CompareQStringView equal")
	}
	if !v.EqualQStringView(sv) {
		t.Fatal("EqualQStringView")
	}
}

func TestCompareQLatin1StringView(t *testing.T) {
	v := fromStr("hello")
	l1 := qlatin1stringview.FromString("hello")

	if v.CompareQLatin1StringView(l1, true) != 0 {
		t.Fatal("CompareQLatin1StringView equal")
	}
	if !v.EqualQLatin1StringView(l1) {
		t.Fatal("EqualQLatin1StringView")
	}
}

func TestCaseInsensitive(t *testing.T) {
	v := fromStr("Hello")
	sv := qstringview.FromUtf16(utf16.Encode([]rune("HELLO")))

	if v.CompareQStringView(sv, true) == 0 {
		t.Fatal("case-sensitive must differ")
	}
	if v.CompareQStringView(sv, false) != 0 {
		t.Fatal("case-insensitive must be 0")
	}
}

func TestToString(t *testing.T) {
	v := fromStr("héllo")
	s := v.ToString()
	defer s.Close()
	if s.GoString() != "héllo" {
		t.Fatalf("ToString: %q", s.GoString())
	}
}

func TestMutatorsInPlace(t *testing.T) {
	v := fromStr("abcdef")

	v.Chop(2)
	if v.String() != "abcd" {
		t.Fatalf("after Chop: %q", v.String())
	}
	v.Truncate(2)
	if v.String() != "ab" {
		t.Fatalf("after Truncate: %q", v.String())
	}
	v.Slice(1)
	if v.String() != "b" {
		t.Fatalf("after Slice: %q", v.String())
	}
}

func TestBytesCopyIndependent(t *testing.T) {
	b := []byte("hello")
	v := FromBytes(b)

	out := v.Bytes()
	out[0] = 'X'
	if b[0] != 'h' {
		t.Fatal("Bytes() must return a copy")
	}
	if !bytes.Equal(b, []byte("hello")) {
		t.Fatal("source mutated")
	}
}

func TestMaxSize(t *testing.T) {
	if MaxSize() <= 0 {
		t.Fatalf("MaxSize: %d", MaxSize())
	}
}

func TestEmoji(t *testing.T) {
	v := fromStr("a😀b")
	if v.String() != "a😀b" {
		t.Fatalf("String: %q", v.String())
	}
	if !v.IsValidUtf8() {
		t.Fatal("a😀b is valid UTF-8")
	}
	// 😀 = 4 UTF-8 bytes, so total = 1+4+1 = 6 bytes.
	if v.Size() != 6 {
		t.Fatalf("Size: %d, want 6", v.Size())
	}
}
