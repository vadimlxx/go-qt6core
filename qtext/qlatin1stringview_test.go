package qtext

import (
	"bytes"
	"testing"
	"unicode/utf16"
)

func TestEmptyQLatin1StringView(t *testing.T) {
	var v QLatin1StringView
	if !v.IsNull() || !v.IsEmpty() || v.Size() != 0 {
		t.Fatalf("Empty: null=%v empty=%v size=%d", v.IsNull(), v.IsEmpty(), v.Size())
	}
	if v.Bytes() != nil || v.GoString() != "" {
		t.Fatal("Empty content")
	}
}

func TestQLatin1StringViewFromBytes(t *testing.T) {
	b := []byte("hello")
	v := QLatin1StringViewFromBytes(b)
	if v.Size() != 5 || v.IsNull() || v.IsEmpty() {
		t.Fatalf("Size=%d IsNull=%v IsEmpty=%v", v.Size(), v.IsNull(), v.IsEmpty())
	}
	if v.GoString() != "hello" {
		t.Fatalf("GoString: %q", v.GoString())
	}
}

func TestQLatin1StringViewFromBytesLatin1(t *testing.T) {
	// "héllo" as Latin-1 bytes: h=0x68, é=0xE9, l, l, o
	b := []byte{0x68, 0xE9, 0x6C, 0x6C, 0x6F}
	v := QLatin1StringViewFromBytes(b)
	if v.Size() != 5 {
		t.Fatalf("Size: %d", v.Size())
	}
	if v.At(1).ToLatin1() != 0xE9 {
		t.Fatalf("At(1): %#x", v.At(1).ToLatin1())
	}
}

func TestQLatin1StringViewFromString(t *testing.T) {
	v := QLatin1StringViewFromString("hello")
	if v.GoString() != "hello" {
		t.Fatalf("GoString: %q", v.GoString())
	}
}

func TestQLatin1StringViewFromQByteArray(t *testing.T) {
	q := NewQByteArray("hello")
	defer q.Close()

	v := QLatin1StringViewFromQByteArray(q)
	if v.GoString() != "hello" {
		t.Fatalf("FromQByteArray: %q", v.GoString())
	}

	sub := v.SlicedN(1, 3)
	if sub.GoString() != "ell" {
		t.Fatalf("sliced: %q", sub.GoString())
	}

	// Null QByteArray → null view.
	qnull := NewQByteArray()
	defer qnull.Close()
	if !QLatin1StringViewFromQByteArray(qnull).IsNull() {
		t.Fatal("FromQByteArray(null) must be null view")
	}
}

func TestQLatin1StringViewAtFrontBack(t *testing.T) {
	v := QLatin1StringViewFromString("abc")

	if v.At(0) != NewQLatin1Char('a') {
		t.Fatalf("At(0): %#x", uint8(v.At(0)))
	}
	if v.At(2) != NewQLatin1Char('c') {
		t.Fatalf("At(2): %#x", uint8(v.At(2)))
	}
	if v.Front() != NewQLatin1Char('a') || v.Back() != NewQLatin1Char('c') {
		t.Fatal("Front/Back")
	}
}

func TestQLatin1StringViewSlicing(t *testing.T) {
	v := QLatin1StringViewFromString("the quick brown fox")

	check := func(name string, got QLatin1StringView, want string) {
		t.Helper()
		if got.GoString() != want {
			t.Fatalf("%s: got %q, want %q", name, got.GoString(), want)
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

func TestQLatin1StringViewTrimmed(t *testing.T) {
	v := QLatin1StringViewFromString("   hello   ").Trimmed()
	if v.GoString() != "hello" {
		t.Fatalf("Trimmed: %q", v.GoString())
	}
}

func TestQLatin1StringViewSearchView(t *testing.T) {
	v := QLatin1StringViewFromString("a,b,c,,d")

	if v.IndexOfQChar(NewQChar(','), 0, true) != 1 {
		t.Fatalf("IndexOfQChar(','): %d", v.IndexOfQChar(NewQChar(','), 0, true))
	}
	if v.LastIndexOfQChar(NewQChar(','), true) != 6 {
		t.Fatalf("LastIndexOfQChar(','): %d", v.LastIndexOfQChar(NewQChar(','), true))
	}
	if v.CountQChar(NewQChar(','), true) != 4 {
		t.Fatalf("CountQChar(','): %d", v.CountQChar(NewQChar(','), true))
	}

	if !v.ContainsView(QLatin1StringViewFromString("b,c"), true) {
		t.Fatal("ContainsView(b,c)")
	}
	if v.IndexOfView(QLatin1StringViewFromString("zz"), 0, true) != -1 {
		t.Fatal("missing needle found")
	}
	if !v.StartsWithView(QLatin1StringViewFromString("a,"), true) ||
		!v.EndsWithView(QLatin1StringViewFromString("d"), true) {

		t.Fatal("StartsWithView/EndsWithView")
	}
}

func TestSearchQStringView(t *testing.T) {
	v := QLatin1StringViewFromString("a,b,c,,d")

	// QStringView interop: search for "," as UTF-16 code unit.
	sv := QStringViewFromUtf16(utf16.Encode([]rune(",")))
	if v.IndexOfQStringView(sv, 0, true) != 1 {
		t.Fatalf("IndexOfQStringView(','): %d", v.IndexOfQStringView(sv, 0, true))
	}

	// "b,c" as QStringView
	svBC := QStringViewFromUtf16(utf16.Encode([]rune("b,c")))
	if !v.ContainsQStringView(svBC, true) {
		t.Fatal("ContainsQStringView(b,c)")
	}

	// StartsWith/EndsWith
	svA := QStringViewFromUtf16(utf16.Encode([]rune("a,")))
	svD := QStringViewFromUtf16(utf16.Encode([]rune("d")))
	if !v.StartsWithQStringView(svA, true) || !v.EndsWithQStringView(svD, true) {
		t.Fatal("StartsWithQStringView/EndsWithQStringView")
	}
}

func TestQLatin1StringViewCompare(t *testing.T) {
	a := QLatin1StringViewFromString("hello")
	b := QLatin1StringViewFromString("hello")
	c := QLatin1StringViewFromString("world")

	if a.CompareView(b, true) != 0 {
		t.Fatal("equal views must compare 0")
	}
	if a.CompareView(c, true) >= 0 {
		t.Fatal("'hello' must compare < 'world'")
	}
	if a.CompareQChar(NewQChar('h'), true) != 1 {
		// "hello" vs 'h' → 1 (first char matches, size > 1)
		t.Fatalf("CompareQChar('h'): %d", a.CompareQChar(NewQChar('h'), true))
	}
}

func TestCompareQStringView(t *testing.T) {
	a := QLatin1StringViewFromString("hello")
	sv := QStringViewFromUtf16(utf16.Encode([]rune("hello")))
	if a.CompareQStringView(sv, true) != 0 {
		t.Fatal("CompareQStringView equal views must compare 0")
	}
}

func TestQLatin1StringViewCaseInsensitive(t *testing.T) {
	v := QLatin1StringViewFromString("Hello, World!")

	if v.IndexOfQChar(NewQChar('w'), 0, true) != -1 {
		t.Fatal("case-sensitive IndexOf('w') must be -1")
	}
	if v.IndexOfQChar(NewQChar('w'), 0, false) != 7 {
		t.Fatalf("case-insensitive IndexOf('w'): %d", v.IndexOfQChar(NewQChar('w'), 0, false))
	}
}

func TestQLatin1StringViewNumeric(t *testing.T) {
	v := QLatin1StringViewFromString("42")
	if r, ok := v.ToInt(10); !ok || r != 42 {
		t.Fatalf("ToInt: %d %v", r, ok)
	}
	if _, ok := QLatin1StringViewFromString("x").ToInt(10); ok {
		t.Fatal("ToInt should fail")
	}
	if r, ok := QLatin1StringViewFromString("3.14").ToDouble(); !ok || r < 3.13 || r > 3.15 {
		t.Fatalf("ToDouble: %f %v", r, ok)
	}
}

func TestQLatin1StringViewToUtf8(t *testing.T) {
	// Latin-1 bytes 0x68 0xE9 0x6C 0x6C 0x6F → "héllo" in UTF-8 (6 bytes).
	v := QLatin1StringViewFromBytes([]byte{0x68, 0xE9, 0x6C, 0x6C, 0x6F})
	u := v.ToUtf8()
	defer u.Close()

	if want := []byte("héllo"); !bytes.Equal(u.Bytes(), want) {
		t.Fatalf("ToUtf8: got %v, want %v", u.Bytes(), want)
	}
}

func TestQLatin1StringViewMutatorsInPlace(t *testing.T) {
	v := QLatin1StringViewFromString("abcdef")

	v.Chop(2)
	if v.GoString() != "abcd" {
		t.Fatalf("after Chop: %q", v.GoString())
	}
	v.Truncate(2)
	if v.GoString() != "ab" {
		t.Fatalf("after Truncate: %q", v.GoString())
	}
	v.Slice(1)
	if v.GoString() != "b" {
		t.Fatalf("after Slice: %q", v.GoString())
	}
}

func TestQLatin1StringAlias(t *testing.T) {
	// Compile-time check: QLatin1String and QLatin1StringView are the same type.
	var a QLatin1String = QLatin1StringViewFromString("abc")
	var b QLatin1StringView = a
	if b.GoString() != "abc" {
		t.Fatalf("alias: %q", b.GoString())
	}
}

func TestQLatin1StringViewBytesCopyIndependent(t *testing.T) {
	b := []byte("hello")
	v := QLatin1StringViewFromBytes(b)

	out := v.Bytes()
	out[0] = 'X'
	if b[0] != 'h' {
		t.Fatal("Bytes() must return a copy")
	}
	if !bytes.Equal(b, []byte("hello")) {
		t.Fatal("source mutated")
	}
}

func TestQLatin1StringViewMaxSize(t *testing.T) {
	if QLatin1StringViewMaxSize() <= 0 {
		t.Fatalf("MaxSize: %d", MaxSize())
	}
}
