package qbytearraylist

import (
	"testing"

	"github.com/vadimlxx/go-qt6core/qbytearray"
)

func mustBA(s string) *qbytearray.QByteArray {
	return qbytearray.FromString(s)
}

func TestEmpty(t *testing.T) {
	l := New()
	defer l.Close()
	if !l.IsEmpty() || l.Size() != 0 {
		t.Fatalf("New(): empty=%v size=%d", l.IsEmpty(), l.Size())
	}
}

func TestFromQByteArrays(t *testing.T) {
	a := mustBA("one")
	b := mustBA("two")
	c := mustBA("three")
	defer a.Close()
	defer b.Close()
	defer c.Close()

	l := FromQByteArrays([]*qbytearray.QByteArray{a, b, c})
	defer l.Close()

	if got := l.Size(); got != 3 {
		t.Fatalf("size=%d", got)
	}

	// At() returns copies, so mutating them doesn't affect the list.
	got := l.At(1)
	defer got.Close()
	if got.String() != "two" {
		t.Fatalf("At(1)=%q", got.String())
	}
}

func TestFromStrings(t *testing.T) {
	l := FromStrings([]string{"a", "b", "c"})
	defer l.Close()
	if l.Size() != 3 {
		t.Fatalf("size=%d", l.Size())
	}
	if l.At(0).String() != "a" || l.At(2).String() != "c" {
		t.Fatal("FromStrings content mismatch")
	}
}

func TestCloneCOW(t *testing.T) {
	l := FromStrings([]string{"x", "y"})
	defer l.Close()
	c := l.Clone()
	defer c.Close()

	// Mutation of clone must not affect original (Qt COW detach).
	c.RemoveAt(0)
	if l.Size() != 2 || c.Size() != 1 {
		t.Fatalf("COW broken: orig=%d clone=%d", l.Size(), c.Size())
	}
	if c.At(0).String() != "y" {
		t.Fatalf("clone[0]=%q", c.At(0).String())
	}
}

func TestAppendPrependInsertSet(t *testing.T) {
	l := New()
	defer l.Close()

	x := mustBA("X")
	y := mustBA("Y")
	z := mustBA("Z")
	defer x.Close()
	defer y.Close()
	defer z.Close()

	l.Append(x)
	l.Append(y)
	l.Prepend(z) // Z X Y
	if l.Size() != 3 || l.At(0).String() != "Z" || l.At(2).String() != "Y" {
		t.Fatal("prepend/append wrong")
	}

	w := mustBA("W")
	defer w.Close()
	l.Insert(1, w) // Z W X Y
	if l.At(1).String() != "W" || l.At(2).String() != "X" {
		t.Fatal("insert wrong")
	}

	v := mustBA("V")
	defer v.Close()
	l.Set(0, v) // V W X Y
	if l.At(0).String() != "V" {
		t.Fatalf("Set wrong: %q", l.At(0).String())
	}
}

func TestRemove(t *testing.T) {
	l := FromStrings([]string{"a", "b", "c"})
	defer l.Close()

	l.RemoveAt(1)
	if l.Size() != 2 || l.At(0).String() != "a" || l.At(1).String() != "c" {
		t.Fatal("RemoveAt wrong")
	}
	l.RemoveFirst()
	if l.Size() != 1 || l.At(0).String() != "c" {
		t.Fatal("RemoveFirst wrong")
	}
	l.RemoveLast()
	if !l.IsEmpty() {
		t.Fatal("RemoveLast wrong")
	}
}

func TestSearch(t *testing.T) {
	l := FromStrings([]string{"a", "b", "a", "c"})
	defer l.Close()

	a := mustBA("a")
	defer a.Close()
	c := mustBA("c")
	defer c.Close()
	missing := mustBA("zzz")
	defer missing.Close()

	if l.IndexOf(a) != 0 || l.LastIndexOf(a) != 2 {
		t.Fatalf("idx(%q): first=%d last=%d", "a", l.IndexOf(a), l.LastIndexOf(a))
	}
	if l.IndexOf(c) != 3 || !l.Contains(c) {
		t.Fatal("Contains/IndexOf wrong")
	}
	if l.IndexOf(missing) != -1 || l.Contains(missing) {
		t.Fatal("missing wrongly found")
	}
}

func TestJoin(t *testing.T) {
	l := FromStrings([]string{"aa", "bb", "cc"})
	defer l.Close()

	j := l.JoinString(", ")
	defer j.Close()
	if j.String() != "aa, bb, cc" {
		t.Fatalf("JoinString: %q", j.String())
	}

	jb := l.JoinByte('-')
	defer jb.Close()
	if jb.String() != "aa-bb-cc" {
		t.Fatalf("JoinByte: %q", jb.String())
	}

	// Empty separator: concatenation.
	j0 := l.Join(nil)
	defer j0.Close()
	if j0.String() != "aabbcc" {
		t.Fatalf("Join(nil): %q", j0.String())
	}

	// Empty list -> empty result.
	empty := New()
	defer empty.Close()
	je := empty.JoinString(",")
	defer je.Close()
	if je.String() != "" {
		t.Fatalf("Join(empty): %q", je.String())
	}
}

func TestToSlice(t *testing.T) {
	l := FromStrings([]string{"1", "2", "3"})
	defer l.Close()

	slice := l.ToSlice()
	if len(slice) != 3 {
		t.Fatalf("len(slice)=%d", len(slice))
	}
	for _, q := range slice {
		defer q.Close()
	}
	if slice[0].String() != "1" || slice[2].String() != "3" {
		t.Fatal("ToSlice content mismatch")
	}

	// Mutating a returned element does not affect the list.
	slice[1].SetAt(0, 'X')
	if l.At(1).String() != "2" {
		t.Fatal("ToSlice must return copies")
	}
}
