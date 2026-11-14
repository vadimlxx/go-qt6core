package qtext

import (
	"bytes"
	"testing"
	"unicode/utf16"

	"github.com/vadimlxx/go-qt6core/qdata"
)

func TestEmptyQString(t *testing.T) {
	s := NewQString()
	defer s.Close()
	if !s.IsEmpty() || !s.IsNull() || s.Size() != 0 {
		t.Fatalf("New: empty=%v null=%v size=%d", s.IsEmpty(), s.IsNull(), s.Size())
	}
	if s.GoString() != "" {
		t.Fatalf("GoString: %q", s.GoString())
	}
}

func TestQStringFromString(t *testing.T) {
	cases := []string{"", "hello", "héllo", "Привет", "日本語", "a😀b"}
	for _, in := range cases {
		s := NewQString(in)
		if s.GoString() != in {
			t.Fatalf("in=%q out=%q", in, s.GoString())
		}
		s.Close()
	}
}

func TestQStringFromUtf16(t *testing.T) {
	units := utf16.Encode([]rune("a😀b"))
	s := NewQString(units)
	defer s.Close()
	if s.GoString() != "a😀b" {
		t.Fatalf("GoString: %q", s.GoString())
	}
	if s.Size() != 4 {
		t.Fatalf("Size: %d", s.Size())
	}
}

func TestQStringFromLatin1(t *testing.T) {
	in := []byte{0x68, 0xE9, 0x6C, 0x6C, 0x6F}
	s := QStringFromLatin1(in)
	defer s.Close()
	if s.GoString() != "héllo" {
		t.Fatalf("GoString: %q", s.GoString())
	}
}

func TestQStringFromLocal8Bit(t *testing.T) {
	// ASCII-only to avoid locale dependence in tests.
	in := []byte("hello")
	s := QStringFromLocal8Bit(in)
	defer s.Close()
	if s.GoString() != "hello" {
		t.Fatalf("GoString: %q", s.GoString())
	}
}

func TestQStringFromQByteArray(t *testing.T) {
	ba := NewQByteArray("héllo")
	defer ba.Close()
	s := NewQString(ba)
	defer s.Close()
	if s.GoString() != "héllo" {
		t.Fatalf("GoString: %q", s.GoString())
	}
}

func TestQStringFromQStringView(t *testing.T) {
	v := QStringViewFromString("héllo")
	s := NewQString(v)
	defer s.Close()
	if s.GoString() != "héllo" {
		t.Fatalf("GoString: %q", s.GoString())
	}
}

func TestQStringToQStringView(t *testing.T) {
	s := QStringFromString("hello")
	defer s.Close()
	v := s.ToQStringView()
	if v.GoString() != "hello" {
		t.Fatalf("ToQStringView: %q", v.GoString())
	}
}

func TestQStringAtFrontBack(t *testing.T) {
	s := NewQString("abc")
	defer s.Close()
	if s.At(0) != NewQChar('a') {
		t.Fatalf("At(0): %#x", s.At(0).Unicode())
	}
	if s.Front() != NewQChar('a') || s.Back() != NewQChar('c') {
		t.Fatal("Front/Back")
	}
}

func TestQStringAppendPrependInsert(t *testing.T) {
	s := QStringFromString("bc")
	defer s.Close()
	s.PrependString("a")
	s.AppendString("d")
	s.InsertString(2, "X")
	if s.GoString() != "abXcd" {
		t.Fatalf("got %q, want abXcd", s.GoString())
	}
}

func TestQStringCaseConversion(t *testing.T) {
	s := QStringFromString("Hello, World!")
	defer s.Close()
	if s.ToLower().GoString() != "hello, world!" {
		t.Fatalf("ToLower: %q", s.ToLower().GoString())
	}
	if s.ToUpper().GoString() != "HELLO, WORLD!" {
		t.Fatalf("ToUpper: %q", s.ToUpper().GoString())
	}
}

func TestQStringTrimmedSimplified(t *testing.T) {
	s := QStringFromString("  Hello,  World!  ")
	defer s.Close()
	if s.Trimmed().GoString() != "Hello,  World!" {
		t.Fatalf("Trimmed: %q", s.Trimmed().GoString())
	}
	if s.Simplified().GoString() != "Hello, World!" {
		t.Fatalf("Simplified: %q", s.Simplified().GoString())
	}
}

func TestQStringSlicing(t *testing.T) {
	s := QStringFromString("the quick brown fox")
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

func TestQStringSearch(t *testing.T) {
	s := QStringFromString("a,b,c,,d")
	defer s.Close()
	if s.IndexOfQChar(NewQChar(','), 0, true) != 1 {
		t.Fatalf("IndexOfQChar(','): %d", s.IndexOfQChar(NewQChar(','), 0, true))
	}
	if s.LastIndexOfQChar(NewQChar(','), -1, true) != 6 {
		t.Fatalf("LastIndexOfQChar(','): %d", s.LastIndexOfQChar(NewQChar(','), -1, true))
	}
	if s.CountQChar(NewQChar(','), true) != 4 {
		t.Fatalf("CountQChar(','): %d", s.CountQChar(NewQChar(','), true))
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

func TestQStringReplace(t *testing.T) {
	s := QStringFromString("hello world")
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

func TestQStringReplaceCS(t *testing.T) {
	s := QStringFromString("Hello hello HELLO")
	defer s.Close()
	s.ReplaceStringStringCS("hello", "X", false)
	if s.GoString() != "X X X" {
		t.Fatalf("case-insensitive replace: %q", s.GoString())
	}
}

func TestQStringRemove(t *testing.T) {
	s := QStringFromString("a-b-c-d")
	defer s.Close()
	s.RemoveQChar('-', true)
	if s.GoString() != "abcd" {
		t.Fatalf("RemoveQChar: %q", s.GoString())
	}
}

func TestQStringCompare(t *testing.T) {
	a := QStringFromString("hello")
	b := QStringFromString("hello")
	c := QStringFromString("world")
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

func TestQStringLocaleAwareCompare(t *testing.T) {
	a := QStringFromString("a")
	b := QStringFromString("b")
	defer a.Close()
	defer b.Close()
	if a.LocaleAwareCompare(b) >= 0 {
		t.Fatal("expected a < b")
	}
}

func TestQStringNumeric(t *testing.T) {
	s := QStringFromString("42")
	defer s.Close()
	if v, ok := s.ToInt(10); !ok || v != 42 {
		t.Fatalf("ToInt: %d %v", v, ok)
	}
	f := QStringFromString("3.14")
	defer f.Close()
	if v, ok := f.ToDouble(); !ok || v < 3.13 || v > 3.15 {
		t.Fatalf("ToDouble: %f %v", v, ok)
	}
	n := QStringNumberLL(-123, 10)
	defer n.Close()
	if n.GoString() != "-123" {
		t.Fatalf("NumberLL: %q", n.GoString())
	}
}

func TestQStringSetNum(t *testing.T) {
	s := NewQString()
	defer s.Close()
	s.SetNumInt(42, 10)
	if s.GoString() != "42" {
		t.Fatalf("SetNumInt: %q", s.GoString())
	}
	s.SetNumUInt(7, 10)
	if s.GoString() != "7" {
		t.Fatalf("SetNumUInt: %q", s.GoString())
	}
	s.SetNumFloat(1.5, 'f', 2)
	if s.GoString() != "1.50" {
		t.Fatalf("SetNumFloat: %q", s.GoString())
	}
}

func TestQStringNumberInt(t *testing.T) {
	s := QStringNumberInt(-7, 10)
	defer s.Close()
	if s.GoString() != "-7" {
		t.Fatalf("NumberInt: %q", s.GoString())
	}
	s2 := QStringNumberUInt(255, 16)
	defer s2.Close()
	if s2.GoString() != "ff" {
		t.Fatalf("NumberUInt: %q", s2.GoString())
	}
}

func TestQStringUtf8Latin1(t *testing.T) {
	s := QStringFromString("héllo")
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

func TestQStringToUcs4(t *testing.T) {
	s := QStringFromString("a😀b")
	defer s.Close()
	cps := s.ToUcs4()
	if len(cps) != 3 {
		t.Fatalf("ToUcs4 len: %d", len(cps))
	}
	if cps[0] != 'a' || cps[1] != 0x1F600 || cps[2] != 'b' {
		t.Fatalf("ToUcs4: %v", cps)
	}
}

func TestQStringToHtmlEscaped(t *testing.T) {
	s := QStringFromString("<a href=\"x\">&</a>")
	defer s.Close()
	e := s.ToHtmlEscaped()
	defer e.Close()
	want := "&lt;a href=&quot;x&quot;&gt;&amp;&lt;/a&gt;"
	if e.GoString() != want {
		t.Fatalf("ToHtmlEscaped: %q", e.GoString())
	}
}

func TestQStringRepeated(t *testing.T) {
	s := QStringFromString("ab")
	defer s.Close()
	r := s.Repeated(3)
	defer r.Close()
	if r.GoString() != "ababab" {
		t.Fatalf("Repeated: %q", r.GoString())
	}
}

func TestQStringNormalized(t *testing.T) {
	// U+00E9 (é) vs U+0065 U+0301 (e + combining acute).
	decomposed := QStringFromString("e\u0301")
	defer decomposed.Close()
	composed := decomposed.Normalized(NormalizationForm_C)
	defer composed.Close()
	if composed.Size() != 1 {
		t.Fatalf("Normalized NFC expected size 1, got %d", composed.Size())
	}
}

func TestQStringSplitQChar(t *testing.T) {
	s := QStringFromString("a,b,c")
	defer s.Close()
	parts := s.SplitQChar(NewQChar(','), qdata.KeepEmptyParts, true)
	if len(parts) != 3 {
		t.Fatalf("SplitQChar len: %d", len(parts))
	}
	for _, p := range parts {
		defer p.Close()
	}
	if parts[0].GoString() != "a" || parts[1].GoString() != "b" || parts[2].GoString() != "c" {
		t.Fatalf("SplitQChar values: %q %q %q", parts[0].GoString(), parts[1].GoString(), parts[2].GoString())
	}
}

func TestQStringSplitQString(t *testing.T) {
	s := QStringFromString("a::b::c")
	sep := QStringFromString("::")
	defer s.Close()
	defer sep.Close()
	parts := s.SplitQString(sep, qdata.KeepEmptyParts, true)
	if len(parts) != 3 {
		t.Fatalf("SplitQString len: %d", len(parts))
	}
	for _, p := range parts {
		defer p.Close()
	}
	if parts[2].GoString() != "c" {
		t.Fatalf("SplitQString last: %q", parts[2].GoString())
	}
}

func TestQStringSplitSkipEmpty(t *testing.T) {
	s := QStringFromString("a,,b")
	defer s.Close()
	parts := s.SplitQChar(NewQChar(','), qdata.SkipEmptyParts, true)
	if len(parts) != 2 {
		t.Fatalf("SkipEmptyParts len: %d", len(parts))
	}
	for _, p := range parts {
		defer p.Close()
	}
}

func TestQStringSection(t *testing.T) {
	s := QStringFromString("a,b,c,d")
	defer s.Close()
	sec := s.SectionQChar(NewQChar(','), 1, 2, SectionDefault)
	defer sec.Close()
	if sec.GoString() != "b,c" {
		t.Fatalf("SectionQChar: %q", sec.GoString())
	}
}

func TestQStringArgLl(t *testing.T) {
	s := QStringFromString("n=%1")
	defer s.Close()
	r := s.ArgLl(42, 0, 10, NewQChar(' '))
	defer r.Close()
	if r.GoString() != "n=42" {
		t.Fatalf("ArgLl: %q", r.GoString())
	}
}

func TestQStringArgString(t *testing.T) {
	s := QStringFromString("%1 %2")
	defer s.Close()
	r := s.ArgString("hello", 0, NewQChar(' '))
	defer r.Close()
	if r.GoString() != "hello %2" {
		t.Fatalf("ArgString step 1: %q", r.GoString())
	}
	r2 := r.ArgString("world", 0, NewQChar(' '))
	defer r2.Close()
	if r2.GoString() != "hello world" {
		t.Fatalf("ArgString step 2: %q", r2.GoString())
	}
}

func TestQStringArgQChar(t *testing.T) {
	s := QStringFromString("x=%1")
	defer s.Close()
	r := s.ArgQChar(NewQChar('7'), 0, NewQChar(' '))
	defer r.Close()
	if r.GoString() != "x=7" {
		t.Fatalf("ArgQChar: %q", r.GoString())
	}
}

func TestQStringArgDouble(t *testing.T) {
	s := QStringFromString("pi=%1")
	defer s.Close()
	r := s.ArgDouble(3.14159, 0, 'f', 2, NewQChar(' '))
	defer r.Close()
	if r.GoString() != "pi=3.14" {
		t.Fatalf("ArgDouble: %q", r.GoString())
	}
}

func TestQStringMaxSize(t *testing.T) {
	if QStringMaxSize() <= 0 {
		t.Fatalf("MaxSize: %d", QStringMaxSize())
	}
}

func TestQStringCloneCOW(t *testing.T) {
	a := QStringFromString("abc")
	defer a.Close()
	b := a.Clone()
	defer b.Close()
	b.AppendString("d")
	if a.GoString() != "abc" || b.GoString() != "abcd" {
		t.Fatalf("COW: a=%q b=%q", a.GoString(), b.GoString())
	}
	if !a.IsDetached() || !b.IsDetached() {
		t.Fatal("expected both detached after mutation")
	}
}

func TestQStringCapacity(t *testing.T) {
	s := QStringFromString("hello")
	defer s.Close()
	if s.Capacity() < s.Size() {
		t.Fatalf("Capacity < Size: %d < %d", s.Capacity(), s.Size())
	}
}

func TestQStringIsSharedWith(t *testing.T) {
	a := QStringFromString("abc")
	defer a.Close()
	b := a.Clone()
	defer b.Close()
	if !a.IsSharedWith(b) {
		t.Fatal("expected clone to be shared with original")
	}
	b.AppendString("d")
	if a.IsSharedWith(b) {
		t.Fatal("expected detach after mutation")
	}
}

func TestQStringResizeForOverwrite(t *testing.T) {
	s := NewQString()
	defer s.Close()
	s.ResizeForOverwrite(5)
	if s.Size() != 5 {
		t.Fatalf("ResizeForOverwrite: size=%d", s.Size())
	}
}
