package qanystringview

import (
	"testing"
	"unicode/utf16"
	"unsafe"

	"github.com/vadimlxx/go-qt6core/qbytearray"
	"github.com/vadimlxx/go-qt6core/qchar"
	"github.com/vadimlxx/go-qt6core/qlatin1stringview"
	"github.com/vadimlxx/go-qt6core/qstringview"
	"github.com/vadimlxx/go-qt6core/qutf8stringview"
)

func TestEmpty(t *testing.T) {
	v := Empty()
	if !v.IsNull() || !v.IsEmpty() || v.Size() != 0 {
		t.Fatalf("Empty: null=%v empty=%v size=%d", v.IsNull(), v.IsEmpty(), v.Size())
	}
	if v.Kind() != KindUtf8 {
		t.Fatalf("Empty kind = %v, want KindUtf8", v.Kind())
	}
	if v.ToString().GoString() != "" {
		t.Fatal("ToString on empty must be empty")
	}
}

func TestFromBytesUtf8(t *testing.T) {
	v := FromBytes([]byte("hello"))
	if v.Kind() != KindUtf8 || v.Size() != 5 || v.SizeBytes() != 5 {
		t.Fatalf("kind=%v size=%d sizeBytes=%d", v.Kind(), v.Size(), v.SizeBytes())
	}
	if v.ToString().GoString() != "hello" {
		t.Fatalf("ToString: %q", v.ToString().GoString())
	}

	sv, ok := v.AsQUtf8StringView()
	if !ok || sv.String() != "hello" {
		t.Fatalf("AsQUtf8StringView: ok=%v str=%q", ok, sv.String())
	}
	if _, ok := v.AsQStringView(); ok {
		t.Fatal("AsQStringView must be !ok for UTF-8")
	}
	if _, ok := v.AsQLatin1StringView(); ok {
		t.Fatal("AsQLatin1StringView must be !ok for UTF-8")
	}
}

func TestFromQStringView(t *testing.T) {
	sv := qstringview.FromUtf16(utf16.Encode([]rune("héllo")))
	v := FromQStringView(sv)
	if v.Kind() != KindUtf16 || v.Size() != 5 || v.SizeBytes() != 10 {
		t.Fatalf("kind=%v size=%d sizeBytes=%d", v.Kind(), v.Size(), v.SizeBytes())
	}
	if v.ToString().GoString() != "héllo" {
		t.Fatalf("ToString: %q", v.ToString().GoString())
	}

	got, ok := v.AsQStringView()
	if !ok || got.GoString() != "héllo" {
		t.Fatalf("AsQStringView: ok=%v str=%q", ok, got.GoString())
	}
}

func TestFromQLatin1StringView(t *testing.T) {
	l1 := qlatin1stringview.FromBytes([]byte{'h', 0xE9, 'l', 'l', 'o'}) // "héllo" in Latin-1
	v := FromQLatin1StringView(l1)
	if v.Kind() != KindLatin1 || v.Size() != 5 || v.SizeBytes() != 5 {
		t.Fatalf("kind=%v size=%d sizeBytes=%d", v.Kind(), v.Size(), v.SizeBytes())
	}
	if v.ToString().GoString() != "héllo" {
		t.Fatalf("ToString: %q", v.ToString().GoString())
	}

	got, ok := v.AsQLatin1StringView()
	if !ok || got.Size() != 5 || got.At(1).ToLatin1() != 0xE9 {
		t.Fatalf("AsQLatin1StringView: ok=%v size=%d at1=%#x", ok, got.Size(), got.At(1).ToLatin1())
	}
}

func TestFromQUtf8StringView(t *testing.T) {
	u8 := qutf8stringview.FromString("héllo")
	v := FromQUtf8StringView(u8)
	if v.Kind() != KindUtf8 {
		t.Fatalf("kind=%v", v.Kind())
	}
	if v.ToString().GoString() != "héllo" {
		t.Fatalf("ToString: %q", v.ToString().GoString())
	}
}

func TestFromQByteArray(t *testing.T) {
	q := qbytearray.FromString("hello")
	defer q.Close()
	v := FromQByteArray(q)
	if v.Kind() != KindUtf8 || v.ToString().GoString() != "hello" {
		t.Fatalf("kind=%v str=%q", v.Kind(), v.ToString().GoString())
	}
}

func TestFrontBack(t *testing.T) {
	v := FromBytes([]byte("abc"))
	if v.Front() != qchar.New('a') || v.Back() != qchar.New('c') {
		t.Fatalf("Front/Back: %#x %#x", v.Front().Unicode(), v.Back().Unicode())
	}
}

func TestSlicing(t *testing.T) {
	v := FromBytes([]byte("the quick brown fox"))

	check := func(name string, got QAnyStringView, want string) {
		t.Helper()
		if got.ToString().GoString() != want {
			t.Fatalf("%s: got %q, want %q", name, got.ToString().GoString(), want)
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

func TestSlicingPreservesKind(t *testing.T) {
	sv := qstringview.FromUtf16(utf16.Encode([]rune("hello")))
	v := FromQStringView(sv)

	if v.Left(3).Kind() != KindUtf16 {
		t.Fatal("Left must preserve Kind")
	}
	if v.Sliced(1).Kind() != KindUtf16 {
		t.Fatal("Sliced must preserve Kind")
	}
}

func TestMutatorsInPlace(t *testing.T) {
	v := FromBytes([]byte("abcdef"))
	v.Chop(2)
	if v.ToString().GoString() != "abcd" {
		t.Fatalf("after Chop: %q", v.ToString().GoString())
	}
	v.Truncate(2)
	if v.ToString().GoString() != "ab" {
		t.Fatalf("after Truncate: %q", v.ToString().GoString())
	}
	v.Slice(1)
	if v.ToString().GoString() != "b" {
		t.Fatalf("after Slice: %q", v.ToString().GoString())
	}
}

func TestCompareEqual(t *testing.T) {
	a := FromBytes([]byte("hello"))
	b := FromBytes([]byte("hello"))
	c := FromBytes([]byte("world"))

	if Compare(a, b, true) != 0 {
		t.Fatal("Compare(equal) must be 0")
	}
	if Compare(a, c, true) >= 0 {
		t.Fatal("Compare(hello, world) must be < 0")
	}
	if !Equal(a, b) {
		t.Fatal("Equal(hello, hello)")
	}
	if Equal(a, c) {
		t.Fatal("Equal(hello, world) must be false")
	}
}

func TestCompareCrossEncoding(t *testing.T) {
	u8 := FromBytes([]byte("hello"))
	utf16v := FromQStringView(qstringview.FromUtf16(utf16.Encode([]rune("hello"))))

	if Compare(u8, utf16v, true) != 0 {
		t.Fatal("UTF-8 'hello' must equal UTF-16 'hello'")
	}
	if !Equal(u8, utf16v) {
		t.Fatal("Equal across encodings")
	}
}

func TestCaseSensitivity(t *testing.T) {
	a := FromBytes([]byte("Hello"))
	b := FromBytes([]byte("HELLO"))

	if Compare(a, b, true) == 0 {
		t.Fatal("case-sensitive compare must differ")
	}
	if Compare(a, b, false) != 0 {
		t.Fatal("case-insensitive compare must be 0")
	}
}

func TestNewUnsafe(t *testing.T) {
	buf := []byte("abc")
	v := NewUnsafe(unsafePointerOf(&buf[0]), 3, KindUtf8)
	if v.ToString().GoString() != "abc" {
		t.Fatalf("NewUnsafe: %q", v.ToString().GoString())
	}
}

// unsafePointerOf is a helper for the test; it accepts a *byte and returns unsafe.Pointer.
func unsafePointerOf(p *byte) unsafe.Pointer { return unsafe.Pointer(p) }
