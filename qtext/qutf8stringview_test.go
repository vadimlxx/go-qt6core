package qtext

import (
	"bytes"
	"testing"
	"unicode/utf16"
)

func TestQUtf8StringViewEmpty(t *testing.T) {
	v := NewQUtf8StringView()
	if !v.IsNull() || !v.IsEmpty() || v.Size() != 0 {
		t.Fatalf("Empty: null=%v empty=%v size=%d", v.IsNull(), v.IsEmpty(), v.Size())
	}
	if v.Bytes() != nil || v.String() != "" {
		t.Fatal("Empty content")
	}
}

func TestQUtf8StringViewFromString(t *testing.T) {
	v := NewQUtf8StringView("hello")
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

func TestQUtf8StringViewFromBytes(t *testing.T) {
	b := []byte("héllo") // UTF-8
	v := NewQUtf8StringView(b)
	if v.String() != "héllo" {
		t.Fatalf("String: %q", v.String())
	}
	if !v.IsValidUtf8() {
		t.Fatal("héllo must be valid UTF-8")
	}

	// Invalid UTF-8 sequence.
	bad := NewQUtf8StringView([]byte{0xFF, 0xFE, 0xFD})
	if bad.IsValidUtf8() {
		t.Fatal("expected invalid UTF-8")
	}
}

func TestQUtf8StringViewFromQByteArray(t *testing.T) {
	q := NewQByteArray("hello")
	defer q.Close()
	v := NewQUtf8StringView(q)
	if v.String() != "hello" {
		t.Fatalf("FromQByteArray: %q", v.String())
	}
	if !NewQUtf8StringView(nil).IsNull() {
		t.Fatal("FromQByteArray(nil) must be null")
	}
}

func TestQUtf8StringViewFromQString(t *testing.T) {
	s := QStringFromString("héllo")
	defer s.Close()
	v := NewQUtf8StringView(s)
	if v.String() != "héllo" {
		t.Fatalf("FromQString: %q", v.String())
	}
	if !v.IsValidUtf8() {
		t.Fatal("valid UTF-8 expected")
	}
}

func TestQUtf8StringViewAtFrontBack(t *testing.T) {
	v := NewQUtf8StringView("abc")
	if v.At(0) != 'a' || v.At(2) != 'c' {
		t.Fatalf("At: %c %c", v.At(0), v.At(2))
	}
	if v.Front() != 'a' || v.Back() != 'c' {
		t.Fatal("Front/Back")
	}
}

func TestQUtf8StringViewSlicing(t *testing.T) {
	v := NewQUtf8StringView("the quick brown fox")

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

func TestQUtf8StringViewCompareView(t *testing.T) {
	a := NewQUtf8StringView("hello")
	b := NewQUtf8StringView("hello")
	c := NewQUtf8StringView("world")

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

func TestQUtf8StringViewCompareQByteArray(t *testing.T) {
	v := NewQUtf8StringView("hello")
	q := NewQByteArray("hello")
	defer q.Close()

	if v.CompareQByteArray(q, true) != 0 {
		t.Fatal("CompareQByteArray equal")
	}
	if !v.EqualQByteArray(q) {
		t.Fatal("EqualQByteArray")
	}
}

func TestQUtf8StringViewCompareQChar(t *testing.T) {
	v := NewQUtf8StringView("a")
	if v.CompareQChar(NewQChar('a'), true) != 0 {
		t.Fatalf("CompareQChar('a'): %d", v.CompareQChar(NewQChar('a'), true))
	}
	if !v.EqualQChar(NewQChar('a')) {
		t.Fatal("EqualQChar('a')")
	}
	if v.EqualQChar(NewQChar('b')) {
		t.Fatal("EqualQChar('b') must be false")
	}
}

func TestQUtf8StringViewCompareQStringView(t *testing.T) {
	v := NewQUtf8StringView("hello")
	sv := QStringViewFromUtf16(utf16.Encode([]rune("hello")))

	if v.CompareQStringView(sv, true) != 0 {
		t.Fatal("CompareQStringView equal")
	}
	if !v.EqualQStringView(sv) {
		t.Fatal("EqualQStringView")
	}
}

func TestQUtf8StringViewCompareQLatin1StringView(t *testing.T) {
	v := NewQUtf8StringView("hello")
	l1 := QLatin1StringViewFromString("hello")

	if v.CompareQLatin1StringView(l1, true) != 0 {
		t.Fatal("CompareQLatin1StringView equal")
	}
	if !v.EqualQLatin1StringView(l1) {
		t.Fatal("EqualQLatin1StringView")
	}
}

func TestQUtf8StringViewCaseInsensitive(t *testing.T) {
	v := NewQUtf8StringView("Hello")
	sv := QStringViewFromUtf16(utf16.Encode([]rune("HELLO")))

	if v.CompareQStringView(sv, true) == 0 {
		t.Fatal("case-sensitive must differ")
	}
	if v.CompareQStringView(sv, false) != 0 {
		t.Fatal("case-insensitive must be 0")
	}
}

func TestQUtf8StringViewToString(t *testing.T) {
	v := NewQUtf8StringView("héllo")
	s := v.ToString()
	defer s.Close()
	if s.GoString() != "héllo" {
		t.Fatalf("ToString: %q", s.GoString())
	}
}

func TestQUtf8StringViewMutatorsInPlace(t *testing.T) {
	v := NewQUtf8StringView("abcdef")

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

func TestQUtf8StringViewBytesCopyIndependent(t *testing.T) {
	b := []byte("hello")
	v := NewQUtf8StringView(b)

	out := v.Bytes()
	out[0] = 'X'
	if b[0] != 'h' {
		t.Fatal("Bytes() must return a copy")
	}
	if !bytes.Equal(b, []byte("hello")) {
		t.Fatal("source mutated")
	}
}

func TestQUtf8StringViewMaxSize(t *testing.T) {
	if QUtf8StringViewMaxSize() <= 0 {
		t.Fatalf("MaxSize: %d", QUtf8StringViewMaxSize())
	}
}

func TestQUtf8StringViewEmoji(t *testing.T) {
	v := NewQUtf8StringView("a😀b")
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
