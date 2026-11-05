package qbytearrayview

import (
	"bytes"
	"testing"

	"github.com/vadimlxx/go-qt6core/qbytearray"
)

func TestEmpty(t *testing.T) {
	v := Empty()
	if !v.IsNull() || !v.IsEmpty() || v.Size() != 0 {
		t.Fatalf("Empty: null=%v empty=%v size=%d", v.IsNull(), v.IsEmpty(), v.Size())
	}
	if v.Bytes() != nil || v.String() != "" {
		t.Fatal("Empty content")
	}
}

func TestFromBytes(t *testing.T) {
	b := []byte("hello")
	v := FromBytes(b)
	if v.Size() != 5 || v.IsNull() || !v.IsEmpty() == false {
		t.Fatalf("Size=%d IsNull=%v", v.Size(), v.IsNull())
	}
	if v.String() != "hello" {
		t.Fatalf("String: %q", v.String())
	}
	if v.At(0) != 'h' || v.Front() != 'h' || v.Back() != 'o' {
		t.Fatal("At/Front/Back wrong")
	}
}

func TestFromString(t *testing.T) {
	v := FromString("héllo")
	if v.String() != "héllo" {
		t.Fatalf("String: %q", v.String())
	}
	if !v.IsValidUtf8() {
		t.Fatal("expected valid utf8")
	}
}

func TestFromEmptyBytes(t *testing.T) {
	if !FromBytes(nil).IsNull() {
		t.Fatal("FromBytes(nil) must be null")
	}
	if !FromBytes([]byte{}).IsNull() {
		t.Fatal("FromBytes([]byte{}) must be null")
	}
	if !FromString("").IsNull() {
		t.Fatal("FromString(\"\") must be null")
	}
}

func TestSlicing(t *testing.T) {
	v := FromString("the quick brown fox")

	check := func(name string, got QByteArrayView, want string) {
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

func TestTrimmed(t *testing.T) {
	v := FromString("   hello   ").Trimmed()
	if v.String() != "hello" {
		t.Fatalf("Trimmed: %q", v.String())
	}
}

func TestSearch(t *testing.T) {
	v := FromString("a,b,c,,d")

	if v.IndexOf(',', 0) != 1 {
		t.Fatalf("IndexOf(','): %d", v.IndexOf(',', 0))
	}
	if v.LastIndexOf(',', -1) != 6 {
		t.Fatalf("LastIndexOf(','): %d", v.LastIndexOf(',', -1))
	}
	if v.Count(',') != 4 {
		t.Fatalf("Count(','): %d", v.Count(','))
	}
	if !v.ContainsBytes([]byte("b,c")) {
		t.Fatal("ContainsBytes(b,c)")
	}
	if v.IndexOfBytes([]byte("zz"), 0) != -1 {
		t.Fatal("missing needle found")
	}
	if !v.StartsWithBytes([]byte("a,")) || !v.EndsWithBytes([]byte("d")) {
		t.Fatal("StartsWith/EndsWith")
	}
}

func TestCompare(t *testing.T) {
	v := FromString("hello")
	if v.Compare([]byte("hello"), true) != 0 {
		t.Fatal("equal views should compare 0")
	}
	if v.Compare([]byte("HELLO"), true) == 0 {
		t.Fatal("case-sensitive must differ")
	}
	if v.Compare([]byte("HELLO"), false) != 0 {
		t.Fatal("case-insensitive should be 0")
	}
}

func TestNumeric(t *testing.T) {
	v := FromString("42")
	if r, ok := v.ToInt(10); !ok || r != 42 {
		t.Fatalf("ToInt: %d %v", r, ok)
	}
	if _, ok := FromString("x").ToInt(10); ok {
		t.Fatal("ToInt should fail")
	}
	if r, ok := FromString("3.14").ToDouble(); !ok || r < 3.13 || r > 3.15 {
		t.Fatalf("ToDouble: %f %v", r, ok)
	}
}

func TestToByteArray(t *testing.T) {
	v := FromString("hello")
	q := v.ToByteArray()
	defer q.Close()
	if q.String() != "hello" {
		t.Fatalf("ToByteArray: %q", q.String())
	}
	// Independent copy: mutating v's source does not affect q.
	b := []byte("hello")
	v2 := FromBytes(b)
	q2 := v2.ToByteArray()
	defer q2.Close()
	b[0] = 'X'
	if q2.String() != "hello" {
		t.Fatal("ToByteArray should return independent copy")
	}
}

func TestFromQByteArray(t *testing.T) {
	q := qbytearray.FromString("hello world")
	defer q.Close()

	v := FromQByteArray(q)
	if v.String() != "hello world" {
		t.Fatalf("FromQByteArray: %q", v.String())
	}

	sub := v.SlicedN(6, 5)
	if sub.String() != "world" {
		t.Fatalf("sliced view: %q", sub.String())
	}

	// Null QByteArray -> null view.
	qnull := qbytearray.New()
	defer qnull.Close()
	if !FromQByteArray(qnull).IsNull() {
		t.Fatal("FromQByteArray(null) must be null view")
	}
}

func TestMutatorsInPlace(t *testing.T) {
	v := FromString("abcdef")

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