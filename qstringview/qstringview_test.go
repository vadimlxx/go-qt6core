package qstringview

import (
	"bytes"
	"testing"
	"unicode/utf16"

	"github.com/vadimlxx/go-qt6core/qbytearray"
	"github.com/vadimlxx/go-qt6core/qchar"
)

func fromStr(s string) QStringView {
	return FromUtf16(utf16.Encode([]rune(s)))
}

func TestEmpty(t *testing.T) {
	v := Empty()
	if !v.IsNull() || !v.IsEmpty() || v.Size() != 0 {
		t.Fatalf("Empty: null=%v empty=%v size=%d", v.IsNull(), v.IsEmpty(), v.Size())
	}
	if v.GoString() != "" || v.Utf16() != nil {
		t.Fatal("Empty content")
	}
}

func TestFromUtf16(t *testing.T) {
	v := fromStr("hello")
	if v.Size() != 5 || v.IsNull() || v.IsEmpty() {
		t.Fatalf("Size=%d IsNull=%v IsEmpty=%v", v.Size(), v.IsNull(), v.IsEmpty())
	}
	if v.GoString() != "hello" {
		t.Fatalf("GoString: %q", v.GoString())
	}
}

func TestFromQChars(t *testing.T) {
	chars := []qchar.QChar{'h', 'e', 'l', 'l', 'o'}
	v := FromQChars(chars)
	if v.GoString() != "hello" {
		t.Fatalf("FromQChars: %q", v.GoString())
	}

	// Zero-copy: mutating chars[0] reflects in view.
	chars[0] = 'H'
	if v.GoString() != "Hello" {
		t.Fatalf("zero-copy broken: %q", v.GoString())
	}
}

func TestAtFrontBack(t *testing.T) {
	v := fromStr("abc")

	if got := v.At(0); got != qchar.New('a') {
		t.Fatalf("At(0) = %#x, want 'a'", got.Unicode())
	}
	if got := v.At(2); got != qchar.New('c') {
		t.Fatalf("At(2) = %#x", got.Unicode())
	}
	if v.Front() != qchar.New('a') || v.Back() != qchar.New('c') {
		t.Fatal("Front/Back")
	}
}

func TestSlicing(t *testing.T) {
	v := fromStr("the quick brown fox")

	check := func(name string, got QStringView, want string) {
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

func TestSearch(t *testing.T) {
	v := fromStr("a,b,c,,d")

	if v.IndexOf(qchar.New(','), 0, true) != 1 {
		t.Fatalf("IndexOf(','): %d", v.IndexOf(qchar.New(','), 0, true))
	}
	if v.LastIndexOf(qchar.New(','), -1, true) != 6 {
		t.Fatalf("LastIndexOf(','): %d", v.LastIndexOf(qchar.New(','), -1, true))
	}
	if v.Count(qchar.New(','), true) != 4 {
		t.Fatalf("Count(','): %d", v.Count(qchar.New(','), true))
	}

	// "b,c" is a real substring of "a,b,c,,d".
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

func TestCaseInsensitive(t *testing.T) {
	v := fromStr("Hello, World!")

	if v.IndexOf(qchar.New('w'), 0, true) != -1 {
		t.Fatal("case-sensitive IndexOf('w') must be -1")
	}
	if v.IndexOf(qchar.New('w'), 0, false) != 7 {
		t.Fatalf("case-insensitive IndexOf('w'): %d", v.IndexOf(qchar.New('w'), 0, false))
	}
	if !v.ContainsView(fromStr("WORLD"), false) {
		t.Fatal("case-insensitive ContainsView(WORLD)")
	}
}

func TestCompare(t *testing.T) {
	a := fromStr("hello")
	b := fromStr("hello")
	c := fromStr("world")

	if a.Compare(b, true) != 0 {
		t.Fatal("equal views must compare 0")
	}
	if a.Compare(c, true) >= 0 {
		t.Fatal("'hello' must compare < 'world'")
	}
	if a.CompareQChar(qchar.New('h'), true) != 1 {
		// "hello" vs 'h' → 1 (size > 1, first char matches)
		t.Fatalf("CompareQChar('h'): %d", a.CompareQChar(qchar.New('h'), true))
	}
}

func TestPredicates(t *testing.T) {
	if !fromStr("HELLO").IsUpper() {
		t.Fatal("HELLO must be upper")
	}
	if !fromStr("hello").IsLower() {
		t.Fatal("hello must be lower")
	}
	if !fromStr("hello").IsValidUtf16() {
		t.Fatal("hello must be valid utf16")
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

func TestConversions(t *testing.T) {
	v := fromStr("héllo")

	u := v.ToUtf8()
	defer u.Close()
	if u.String() != "héllo" {
		t.Fatalf("ToUtf8: %q", u.String())
	}

	l := v.ToLatin1()
	defer l.Close()
	// ToLatin1 produces Latin-1 bytes; 'é' = 0xE9 (not the UTF-8 pair 0xC3 0xA9).
	// "héllo" as Latin-1 is exactly 5 bytes.
	want := []byte{0x68, 0xE9, 0x6C, 0x6C, 0x6F}
	if !bytes.Equal(l.Bytes(), want) {
		t.Fatalf("ToLatin1: got %v, want %v", l.Bytes(), want)
	}

	// Non-Latin-1 chars collapse to '?'.
	j := fromStr("日本").ToLatin1()
	defer j.Close()
	if j.String() != "??" {
		t.Fatalf("ToLatin1(日本): %q", j.String())
	}
}

func TestEmoji(t *testing.T) {
	// U+1F600 → 2 UTF-16 code units (surrogate pair).
	v := fromStr("a😀b")
	if v.Size() != 4 {
		t.Fatalf("Size: %d", v.Size())
	}
	if v.GoString() != "a😀b" {
		t.Fatalf("GoString: %q", v.GoString())
	}

	// QChar from At(1) is the high surrogate.
	hi := v.At(1)
	lo := v.At(2)
	if !hi.IsHighSurrogate() || !lo.IsLowSurrogate() {
		t.Fatalf("surrogates: %#x %#x", hi.Unicode(), lo.Unicode())
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

func TestMaxSize(t *testing.T) {
	if MaxSize() <= 0 {
		t.Fatalf("MaxSize: %d", MaxSize())
	}
}

// Sanity: ToUtf8 returns an independently owned QByteArray.
func TestToUtf8IsOwning(t *testing.T) {
	units := utf16.Encode([]rune("hello"))
	v := FromUtf16(units)
	ba := v.ToUtf8()
	defer ba.Close()

	units[0] = uint16('X') // mutate source
	if ba.String() != "hello" {
		t.Fatalf("ToUtf8 must be an owning copy; got %q", ba.String())
	}
	_ = qbytearray.New // keep import
}
