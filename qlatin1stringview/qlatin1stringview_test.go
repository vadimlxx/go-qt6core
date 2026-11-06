package qlatin1stringview

import (
	"bytes"
	"testing"
	"unicode/utf16"

	"github.com/vadimlxx/go-qt6core/qbytearray"
	"github.com/vadimlxx/go-qt6core/qchar"
	"github.com/vadimlxx/go-qt6core/qlatin1char"
	"github.com/vadimlxx/go-qt6core/qstringview"
)

func fromStr(s string) QLatin1StringView { return FromString(s) }

func TestEmpty(t *testing.T) {
	v := Empty()
	if !v.IsNull() || !v.IsEmpty() || v.Size() != 0 {
		t.Fatalf("Empty: null=%v empty=%v size=%d", v.IsNull(), v.IsEmpty(), v.Size())
	}
	if v.Bytes() != nil || v.GoString() != "" {
		t.Fatal("Empty content")
	}
}

func TestFromBytes(t *testing.T) {
	b := []byte("hello")
	v := FromBytes(b)
	if v.Size() != 5 || v.IsNull() || v.IsEmpty() {
		t.Fatalf("Size=%d IsNull=%v IsEmpty=%v", v.Size(), v.IsNull(), v.IsEmpty())
	}
	if v.GoString() != "hello" {
		t.Fatalf("GoString: %q", v.GoString())
	}
}

func TestFromBytesLatin1(t *testing.T) {
	// "héllo" as Latin-1 bytes: h=0x68, é=0xE9, l, l, o
	b := []byte{0x68, 0xE9, 0x6C, 0x6C, 0x6F}
	v := FromBytes(b)
	if v.Size() != 5 {
		t.Fatalf("Size: %d", v.Size())
	}
	if v.At(1).ToLatin1() != 0xE9 {
		t.Fatalf("At(1): %#x", v.At(1).ToLatin1())
	}
}

func TestFromString(t *testing.T) {
	v := fromStr("hello")
	if v.GoString() != "hello" {
		t.Fatalf("GoString: %q", v.GoString())
	}
}

func TestFromQByteArray(t *testing.T) {
	q := qbytearray.FromString("hello")
	defer q.Close()

	v := FromQByteArray(q)
	if v.GoString() != "hello" {
		t.Fatalf("FromQByteArray: %q", v.GoString())
	}

	sub := v.SlicedN(1, 3)
	if sub.GoString() != "ell" {
		t.Fatalf("sliced: %q", sub.GoString())
	}

	// Null QByteArray → null view.
	qnull := qbytearray.New()
	defer qnull.Close()
	if !FromQByteArray(qnull).IsNull() {
		t.Fatal("FromQByteArray(null) must be null view")
	}
}

func TestAtFrontBack(t *testing.T) {
	v := fromStr("abc")

	if v.At(0) != qlatin1char.New('a') {
		t.Fatalf("At(0): %#x", uint8(v.At(0)))
	}
	if v.At(2) != qlatin1char.New('c') {
		t.Fatalf("At(2): %#x", uint8(v.At(2)))
	}
	if v.Front() != qlatin1char.New('a') || v.Back() != qlatin1char.New('c') {
		t.Fatal("Front/Back")
	}
}

func TestSlicing(t *testing.T) {
	v := fromStr("the quick brown fox")

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

func TestTrimmed(t *testing.T) {
	v := fromStr("   hello   ").Trimmed()
	if v.GoString() != "hello" {
		t.Fatalf("Trimmed: %q", v.GoString())
	}
}

func TestSearchView(t *testing.T) {
	v := fromStr("a,b,c,,d")

	if v.IndexOfQChar(qchar.New(','), 0, true) != 1 {
		t.Fatalf("IndexOfQChar(','): %d", v.IndexOfQChar(qchar.New(','), 0, true))
	}
	if v.LastIndexOfQChar(qchar.New(','), true) != 6 {
		t.Fatalf("LastIndexOfQChar(','): %d", v.LastIndexOfQChar(qchar.New(','), true))
	}
	if v.CountQChar(qchar.New(','), true) != 4 {
		t.Fatalf("CountQChar(','): %d", v.CountQChar(qchar.New(','), true))
	}

	if !v.ContainsView(fromStr("b,c"), true) {
		t.Fatal("ContainsView(b,c)")
	}
	if v.IndexOfView(fromStr("zz"), 0, true) != -1 {
		t.Fatal("missing needle found")
	}
	if !v.StartsWithView(fromStr("a,"), true) || !v.EndsWithView(fromStr("d"), true) {
		t.Fatal("StartsWithView/EndsWithView")
	}
}

func TestSearchQStringView(t *testing.T) {
	v := fromStr("a,b,c,,d")

	// QStringView interop: search for "," as UTF-16 code unit.
	sv := qstringview.FromUtf16(utf16.Encode([]rune(",")))
	if v.IndexOfQStringView(sv, 0, true) != 1 {
		t.Fatalf("IndexOfQStringView(','): %d", v.IndexOfQStringView(sv, 0, true))
	}

	// "b,c" as QStringView
	svBC := qstringview.FromUtf16(utf16.Encode([]rune("b,c")))
	if !v.ContainsQStringView(svBC, true) {
		t.Fatal("ContainsQStringView(b,c)")
	}

	// StartsWith/EndsWith
	svA := qstringview.FromUtf16(utf16.Encode([]rune("a,")))
	svD := qstringview.FromUtf16(utf16.Encode([]rune("d")))
	if !v.StartsWithQStringView(svA, true) || !v.EndsWithQStringView(svD, true) {
		t.Fatal("StartsWithQStringView/EndsWithQStringView")
	}
}

func TestCompare(t *testing.T) {
	a := fromStr("hello")
	b := fromStr("hello")
	c := fromStr("world")

	if a.CompareView(b, true) != 0 {
		t.Fatal("equal views must compare 0")
	}
	if a.CompareView(c, true) >= 0 {
		t.Fatal("'hello' must compare < 'world'")
	}
	if a.CompareQChar(qchar.New('h'), true) != 1 {
		// "hello" vs 'h' → 1 (first char matches, size > 1)
		t.Fatalf("CompareQChar('h'): %d", a.CompareQChar(qchar.New('h'), true))
	}
}

func TestCompareQStringView(t *testing.T) {
	a := fromStr("hello")
	sv := qstringview.FromUtf16(utf16.Encode([]rune("hello")))
	if a.CompareQStringView(sv, true) != 0 {
		t.Fatal("CompareQStringView equal views must compare 0")
	}
}

func TestCaseInsensitive(t *testing.T) {
	v := fromStr("Hello, World!")

	if v.IndexOfQChar(qchar.New('w'), 0, true) != -1 {
		t.Fatal("case-sensitive IndexOf('w') must be -1")
	}
	if v.IndexOfQChar(qchar.New('w'), 0, false) != 7 {
		t.Fatalf("case-insensitive IndexOf('w'): %d", v.IndexOfQChar(qchar.New('w'), 0, false))
	}
}

func TestNumeric(t *testing.T) {
	v := fromStr("42")
	if r, ok := v.ToInt(10); !ok || r != 42 {
		t.Fatalf("ToInt: %d %v", r, ok)
	}
	if _, ok := fromStr("x").ToInt(10); ok {
		t.Fatal("ToInt should fail")
	}
	if r, ok := fromStr("3.14").ToDouble(); !ok || r < 3.13 || r > 3.15 {
		t.Fatalf("ToDouble: %f %v", r, ok)
	}
}

func TestToUtf8(t *testing.T) {
	// Latin-1 bytes 0x68 0xE9 0x6C 0x6C 0x6F → "héllo" in UTF-8 (6 bytes).
	v := FromBytes([]byte{0x68, 0xE9, 0x6C, 0x6C, 0x6F})
	u := v.ToUtf8()
	defer u.Close()

	if want := []byte("héllo"); !bytes.Equal(u.Bytes(), want) {
		t.Fatalf("ToUtf8: got %v, want %v", u.Bytes(), want)
	}
}

func TestMutatorsInPlace(t *testing.T) {
	v := fromStr("abcdef")

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
	var a QLatin1String = fromStr("abc")
	var b QLatin1StringView = a
	if b.GoString() != "abc" {
		t.Fatalf("alias: %q", b.GoString())
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