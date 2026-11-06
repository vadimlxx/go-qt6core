package qstring

import (
	"bytes"
	"testing"
	"unicode/utf16"

	"github.com/vadimlxx/go-qt6core/qbytearray"
	"github.com/vadimlxx/go-qt6core/qchar"
	"github.com/vadimlxx/go-qt6core/qstringview"
)

func TestEmpty(t *testing.T) {
	s := New()
	defer s.Close()
	if !s.IsEmpty() || !s.IsNull() || s.Size() != 0 {
		t.Fatalf("New: empty=%v null=%v size=%d", s.IsEmpty(), s.IsNull(), s.Size())
	}
	if s.GoString() != "" {
		t.Fatalf("GoString: %q", s.GoString())
	}
}

func TestFromString(t *testing.T) {
	cases := []string{"", "hello", "héllo", "Привет", "日本語", "a😀b"}
	for _, in := range cases {
		s := FromString(in)
		if s.GoString() != in {
			t.Fatalf("in=%q out=%q", in, s.GoString())
		}
		s.Close()
	}
}

func TestFromUtf16(t *testing.T) {
	units := utf16.Encode([]rune("a😀b"))
	s := FromUtf16(units)
	defer s.Close()
	if s.GoString() != "a😀b" {
		t.Fatalf("GoString: %q", s.GoString())
	}
	if s.Size() != 4 {
		t.Fatalf("Size: %d", s.Size())
	}
}

func TestFromLatin1(t *testing.T) {
	in := []byte{0x68, 0xE9, 0x6C, 0x6C, 0x6F} // "héllo" in Latin-1
	s := FromLatin1(in)
	defer s.Close()
	if s.GoString() != "héllo" {
		t.Fatalf("GoString: %q", s.GoString())
	}
}

func TestFromQByteArray(t *testing.T) {
	ba := qbytearray.FromString("héllo")
	defer ba.Close()
	s := FromQByteArray(ba)
	defer s.Close()
	if s.GoString() != "héllo" {
		t.Fatalf("GoString: %q", s.GoString())
	}
}

func TestFromQStringView(t *testing.T) {
	v := qstringview.FromString("héllo")
	s := FromQStringView(v)
	defer s.Close()
	if s.GoString() != "héllo" {
		t.Fatalf("GoString: %q", s.GoString())
	}
}

func TestToQStringView(t *testing.T) {
	s := FromString("hello")
	defer s.Close()
	v := s.ToQStringView()
	if v.GoString() != "hello" {
		t.Fatalf("ToQStringView: %q", v.GoString())
	}
}

func TestAtFrontBack(t *testing.T) {
	s := FromString("abc")
	defer s.Close()

	if s.At(0) != qchar.New('a') {
		t.Fatalf("At(0): %#x", s.At(0).Unicode())
	}
	if s.Front() != qchar.New('a') || s.Back() != qchar.New('c') {
		t.Fatal("Front/Back")
	}
}

func TestAppendPrependInsert(t *testing.T) {
	s := FromString("bc")
	defer s.Close()
	s.PrependString("a")
	s.AppendString("d")
	s.InsertString(2, "X")
	if s.GoString() != "abXcd" {
		t.Fatalf("got %q, want abXcd", s.GoString())
	}
}

func TestCaseConversion(t *testing.T) {
	s := FromString("Hello, World!")
	defer s.Close()
	if s.ToLower().GoString() != "hello, world!" {
		t.Fatalf("ToLower: %q", s.ToLower().GoString())
	}
	if s.ToUpper().GoString() != "HELLO, WORLD!" {
		t.Fatalf("ToUpper: %q", s.ToUpper().GoString())
	}
}

func TestTrimmedSimplified(t *testing.T) {
	s := FromString("  Hello,  World!  ")
	defer s.Close()
	if s.Trimmed().GoString() != "Hello,  World!" {
		t.Fatalf("Trimmed: %q", s.Trimmed().GoString())
	}
	if s.Simplified().GoString() != "Hello, World!" {
		t.Fatalf("Simplified: %q", s.Simplified().GoString())
	}
}

func TestSlicing(t *testing.T) {
	s := FromString("the quick brown fox")
	defer s.Close()

	check := func(name string, got *QString, want string) {
		t.Helper()
		defer got.Close()
		if got.GoString() != want {
			t.Fatalf("%s: got %q, want %q", name, got.GoString(), want)
		}
	}

	check("Left", s.Left(3), "the")
	check("Right", s.Right(3), "fox")
	check("Mid", s.Mid(4, 5), "quick")
	check("First", s.First(3), "the")
	check("Last", s.Last(3), "fox")
	check("Sliced", s.Sliced(4), "quick brown fox")
	check("SlicedN", s.SlicedN(4, 5), "quick")
	check("Chopped", s.Chopped(4), "the quick brown")
}

func TestSearch(t *testing.T) {
	s := FromString("a,b,c,,d")
	defer s.Close()

	if s.IndexOfQChar(qchar.New(','), 0, true) != 1 {
		t.Fatalf("IndexOfQChar(','): %d", s.IndexOfQChar(qchar.New(','), 0, true))
	}
	if s.LastIndexOfQChar(qchar.New(','), -1, true) != 6 {
		t.Fatalf("LastIndexOfQChar(','): %d", s.LastIndexOfQChar(qchar.New(','), -1, true))
	}
	if s.CountQChar(qchar.New(','), true) != 4 {
		t.Fatalf("CountQChar(','): %d", s.CountQChar(qchar.New(','), true))
	}
	if !s.ContainsString("b,c", true) {
		t.Fatal("ContainsString(b,c)")
	}
	if s.IndexOfString("zz", 0, true) != -1 {
		t.Fatal("missing found")
	}
	if !s.StartsWithString("a,", true) || !s.EndsWithString("d", true) {
		t.Fatal("StartsWith/EndsWith")
	}
}

func TestReplace(t *testing.T) {
	s := FromString("hello world")
	defer s.Close()

	s.ReplaceStringString("world", "there")
	if s.GoString() != "hello there" {
		t.Fatalf("got %q", s.GoString())
	}
	s.ReplaceQCharQChar('o', '0')
	if s.GoString() != "hell0 there" {
		t.Fatalf("got %q", s.GoString())
	}
}

func TestCompare(t *testing.T) {
	a := FromString("hello")
	b := FromString("hello")
	c := FromString("world")
	defer a.Close()
	defer b.Close()
	defer c.Close()

	if a.CompareQString(b, true) != 0 {
		t.Fatal("equal should compare 0")
	}
	if a.CompareQString(c, true) >= 0 {
		t.Fatal("hello < world")
	}
	if a.CompareString("HELLO", false) != 0 {
		t.Fatal("case-insensitive compare")
	}
}

func TestNumeric(t *testing.T) {
	s := FromString("42")
	defer s.Close()
	if v, ok := s.ToInt(10); !ok || v != 42 {
		t.Fatalf("ToInt: %d %v", v, ok)
	}

	f := FromString("3.14")
	defer f.Close()
	if v, ok := f.ToDouble(); !ok || v < 3.13 || v > 3.15 {
		t.Fatalf("ToDouble: %f %v", v, ok)
	}

	n := NumberLL(-123, 10)
	defer n.Close()
	if n.GoString() != "-123" {
		t.Fatalf("NumberLL: %q", n.GoString())
	}
}

func TestUtf8Latin1(t *testing.T) {
	s := FromString("héllo")
	defer s.Close()

	u := s.Utf8()
	defer u.Close()
	if !bytes.Equal(u.Bytes(), []byte("héllo")) {
		t.Fatalf("Utf8: %v", u.Bytes())
	}

	l := s.Latin1()
	defer l.Close()
	want := []byte{0x68, 0xE9, 0x6C, 0x6C, 0x6F}
	if !bytes.Equal(l.Bytes(), want) {
		t.Fatalf("Latin1: %v, want %v", l.Bytes(), want)
	}
}

func TestMaxSize(t *testing.T) {
	if MaxSize() <= 0 {
		t.Fatalf("MaxSize: %d", MaxSize())
	}
}

func TestCloneCOW(t *testing.T) {
	a := FromString("abc")
	defer a.Close()
	b := a.Clone()
	defer b.Close()

	b.AppendString("d")
	if a.GoString() != "abc" || b.GoString() != "abcd" {
		t.Fatalf("COW: a=%q b=%q", a.GoString(), b.GoString())
	}
}
