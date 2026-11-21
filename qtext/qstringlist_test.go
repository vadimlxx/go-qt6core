package qtext

import (
	"testing"
)

// ---- construction --------------------------------------------------------

func TestQStringListEmpty(t *testing.T) {
	l := NewQStringList()
	defer l.Close()
	if !l.IsEmpty() || l.Size() != 0 || l.Len() != 0 {
		t.Fatalf("empty=%v size=%d len=%d", l.IsEmpty(), l.Size(), l.Len())
	}
	if got := l.ToSlice(); got != nil {
		t.Fatalf("ToSlice on empty: %v", got)
	}
}

func TestQStringListFromQString(t *testing.T) {
	s := NewQStringFromString("hello")
	defer s.Close()
	l := QStringListFromQString(s)
	defer l.Close()
	if l.Size() != 1 {
		t.Fatalf("size: %d", l.Size())
	}
	elem := l.At(0)
	defer elem.Close()
	if elem.String() != "hello" {
		t.Fatalf("element: %q", elem.String())
	}
}

func TestQStringListFromGoString(t *testing.T) {
	l := QStringListFromGoString("héllo")
	defer l.Close()
	if l.Size() != 1 {
		t.Fatalf("size: %d", l.Size())
	}
	if l.Front().String() != "héllo" {
		t.Fatalf("element: %q", l.Front().String())
	}
}

func TestQStringListFromQStrings(t *testing.T) {
	a := NewQStringFromString("a")
	b := NewQStringFromString("b")
	c := NewQStringFromString("c")
	defer a.Close()
	defer b.Close()
	defer c.Close()
	l := QStringListFromQStrings([]*QString{a, b, c})
	defer l.Close()
	if l.Size() != 3 {
		t.Fatalf("size: %d", l.Size())
	}
	for i, want := range []string{"a", "b", "c"} {
		e := l.At(i)
		if e.String() != want {
			t.Fatalf("At(%d): got %q want %q", i, e.String(), want)
		}
		e.Close()
	}
}

func TestQStringListFromGoStrings(t *testing.T) {
	l := QStringListFromGoStrings([]string{"x", "y", "z"})
	defer l.Close()
	if l.Size() != 3 {
		t.Fatalf("size: %d", l.Size())
	}
	got := l.ToSlice()
	defer func() {
		for _, s := range got {
			s.Close()
		}
	}()
	for i, want := range []string{"x", "y", "z"} {
		if got[i].String() != want {
			t.Fatalf("element %d: got %q want %q", i, got[i].String(), want)
		}
	}
}

func TestQStringListClone(t *testing.T) {
	l := QStringListFromGoStrings([]string{"a", "b"})
	defer l.Close()
	c := l.Clone()
	defer c.Close()
	c.AppendGoString("c")
	if l.Size() != 2 || c.Size() != 3 {
		t.Fatalf("clone not independent: orig=%d clone=%d", l.Size(), c.Size())
	}
}

func TestQStringListSwap(t *testing.T) {
	a := QStringListFromGoStrings([]string{"a", "b"})
	b := QStringListFromGoStrings([]string{"x", "y", "z"})
	defer a.Close()
	defer b.Close()
	a.Swap(b)
	if a.Size() != 3 || b.Size() != 2 {
		t.Fatalf("swap failed: a=%d b=%d", a.Size(), b.Size())
	}
}

// ---- modifiers -----------------------------------------------------------

func TestQStringListClearReserveSqueeze(t *testing.T) {
	l := QStringListFromGoStrings([]string{"a", "b", "c"})
	defer l.Close()
	l.Reserve(100)
	l.Squeeze()
	l.Clear()
	if !l.IsEmpty() {
		t.Fatal("expected empty after Clear")
	}
}

// ---- element access -----------------------------------------------------

func TestQStringListAtSetAt(t *testing.T) {
	l := QStringListFromGoStrings([]string{"a", "b", "c"})
	defer l.Close()
	e := l.At(1)
	if e.String() != "b" {
		t.Fatalf("At(1): %q", e.String())
	}
	e.Close()

	n := NewQStringFromString("B")
	defer n.Close()
	l.SetAt(1, n)

	got := l.At(1)
	defer got.Close()
	if got.String() != "B" {
		t.Fatalf("after SetAt: %q", got.String())
	}
}

func TestQStringListFrontBack(t *testing.T) {
	l := QStringListFromGoStrings([]string{"a", "b", "c"})
	defer l.Close()
	f := l.Front()
	if f.String() != "a" {
		t.Fatalf("Front: %q", f.String())
	}
	f.Close()
	b := l.Back()
	defer b.Close()
	if b.String() != "c" {
		t.Fatalf("Back: %q", b.String())
	}
}

// ---- append / insert / remove -------------------------------------------

func TestQStringListAppendPrependInsert(t *testing.T) {
	l := NewQStringList()
	defer l.Close()
	s := NewQStringFromString("b")
	defer s.Close()
	l.Append(s)
	l.Prepend(s)
	x := NewQStringFromString("x")
	defer x.Close()
	l.Insert(1, x) // a b x b → actually b x b
	want := []string{"b", "x", "b"}
	got := l.ToSlice()
	defer func() {
		for _, e := range got {
			e.Close()
		}
	}()
	if len(got) != len(want) {
		t.Fatalf("len: %d", len(got))
	}
	for i, w := range want {
		if got[i].String() != w {
			t.Fatalf("index %d: got %q want %q", i, got[i].String(), w)
		}
	}
}

func TestQStringListAppendGoString(t *testing.T) {
	l := NewQStringList()
	defer l.Close()
	l.AppendGoString("one")
	l.AppendGoString("two")
	if l.Size() != 2 {
		t.Fatalf("size: %d", l.Size())
	}
}

func TestQStringListAppendList(t *testing.T) {
	a := QStringListFromGoStrings([]string{"a", "b"})
	b := QStringListFromGoStrings([]string{"c", "d"})
	defer a.Close()
	defer b.Close()
	a.AppendList(b)
	if a.Size() != 4 {
		t.Fatalf("size after AppendList: %d", a.Size())
	}
}

func TestQStringListRemove(t *testing.T) {
	l := QStringListFromGoStrings([]string{"a", "b", "c", "d"})
	defer l.Close()
	l.RemoveAt(1)
	l.RemoveFirst()
	l.RemoveLast()
	if l.Size() != 1 {
		t.Fatalf("size after removes: %d", l.Size())
	}
	got := l.At(0)
	defer got.Close()
	if got.String() != "c" {
		t.Fatalf("remaining: %q", got.String())
	}
}

func TestQStringListTake(t *testing.T) {
	l := QStringListFromGoStrings([]string{"a", "b", "c"})
	defer l.Close()

	f := l.TakeFirst()
	if f.String() != "a" {
		t.Fatalf("TakeFirst: %q", f.String())
	}
	f.Close()

	b := l.TakeLast()
	if b.String() != "c" {
		t.Fatalf("TakeLast: %q", b.String())
	}
	b.Close()

	m := l.TakeAt(0)
	defer m.Close()
	if m.String() != "b" {
		t.Fatalf("TakeAt(0): %q", m.String())
	}
	if !l.IsEmpty() {
		t.Fatal("expected empty after takes")
	}
}

// ---- QStringList-specific ------------------------------------------------

func TestQStringListSort(t *testing.T) {
	l := QStringListFromGoStrings([]string{"banana", "Apple", "cherry", "apple"})
	defer l.Close()

	l.Sort(true) // case-sensitive
	got := l.ToSlice()
	defer func() {
		for _, e := range got {
			e.Close()
		}
	}()
	want := []string{"Apple", "apple", "banana", "cherry"}
	for i, w := range want {
		if got[i].String() != w {
			t.Fatalf("sorted[%d]: got %q want %q", i, got[i].String(), w)
		}
	}
}

func TestQStringListRemoveDuplicates(t *testing.T) {
	l := QStringListFromGoStrings([]string{"a", "a", "b", "b", "c"})
	defer l.Close()
	l.Sort(true)
	n := l.RemoveDuplicates()
	if n != 2 {
		t.Fatalf("RemoveDuplicates returned %d, want 2", n)
	}
	if l.Size() != 3 {
		t.Fatalf("size after dedupe: %d", l.Size())
	}
}

func TestQStringListJoinString(t *testing.T) {
	l := QStringListFromGoStrings([]string{"a", "b", "c"})
	defer l.Close()
	sep := NewQStringFromString(", ")
	defer sep.Close()
	r := l.JoinString(sep)
	defer r.Close()
	if r.String() != "a, b, c" {
		t.Fatalf("JoinString: %q", r.String())
	}
}

func TestQStringListJoinGoString(t *testing.T) {
	l := QStringListFromGoStrings([]string{"x", "y"})
	defer l.Close()
	r := l.JoinGoString("-")
	defer r.Close()
	if r.String() != "x-y" {
		t.Fatalf("JoinGoString: %q", r.String())
	}
}

func TestQStringListJoinChar(t *testing.T) {
	l := QStringListFromGoStrings([]string{"1", "2", "3"})
	defer l.Close()
	r := l.JoinChar(NewQChar('|'))
	defer r.Close()
	if r.String() != "1|2|3" {
		t.Fatalf("JoinChar: %q", r.String())
	}
}

func TestQStringListFilterString(t *testing.T) {
	l := QStringListFromGoStrings([]string{"apple", "banana", "apricot", "cherry"})
	defer l.Close()
	needle := NewQStringFromString("ap")
	defer needle.Close()
	f := l.FilterString(needle, true)
	defer f.Close()
	if f.Size() != 2 {
		t.Fatalf("FilterString size: %d", f.Size())
	}
}

func TestQStringListFilterGoString(t *testing.T) {
	l := QStringListFromGoStrings([]string{"Apple", "apple", "APRICOT"})
	defer l.Close()
	f := l.FilterGoString("ap", false) // case-insensitive
	defer f.Close()
	if f.Size() != 3 {
		t.Fatalf("case-insensitive filter size: %d", f.Size())
	}
}

func TestQStringListFilterQRegularExpression(t *testing.T) {
	l := QStringListFromGoStrings([]string{"file1.txt", "readme.md", "file2.txt", "script.sh"})
	defer l.Close()
	re := NewQRegularExpressionFromString(`\.txt$`, NoPatternOption)
	defer re.Close()
	f := l.FilterQRegularExpression(re)
	defer f.Close()
	if f.Size() != 2 {
		t.Fatalf("regex filter size: %d", f.Size())
	}
}

func TestQStringListReplaceInStringsString(t *testing.T) {
	l := QStringListFromGoStrings([]string{"foo bar", "bar baz", "FOO qux"})
	defer l.Close()
	before := NewQStringFromString("foo")
	after := NewQStringFromString("X")
	defer before.Close()
	defer after.Close()
	l.ReplaceInStringsString(before, after, true)
	got := l.ToSlice()
	defer func() {
		for _, e := range got {
			e.Close()
		}
	}()
	want := []string{"X bar", "bar baz", "FOO qux"}
	for i, w := range want {
		if got[i].String() != w {
			t.Fatalf("replace[%d]: got %q want %q", i, got[i].String(), w)
		}
	}
}

func TestQStringListReplaceInStringsGoString(t *testing.T) {
	l := QStringListFromGoStrings([]string{"Foo", "foo", "FOO"})
	defer l.Close()
	l.ReplaceInStringsGoString("foo", "X", false)
	got := l.ToSlice()
	defer func() {
		for _, e := range got {
			e.Close()
		}
	}()
	for i, w := range []string{"X", "X", "X"} {
		if got[i].String() != w {
			t.Fatalf("replace[%d]: got %q want %q", i, got[i].String(), w)
		}
	}
}

func TestQStringListReplaceInStringsQRegularExpression(t *testing.T) {
	l := QStringListFromGoStrings([]string{"a1b22c", "d3e", "no-digits"})
	defer l.Close()
	re := NewQRegularExpressionFromString(`\d+`, NoPatternOption)
	defer re.Close()
	after := NewQStringFromString("#")
	defer after.Close()
	l.ReplaceInStringsQRegularExpression(re, after)
	got := l.ToSlice()
	defer func() {
		for _, e := range got {
			e.Close()
		}
	}()
	want := []string{"a#b#c", "d#e", "no-digits"}
	for i, w := range want {
		if got[i].String() != w {
			t.Fatalf("replace[%d]: got %q want %q", i, got[i].String(), w)
		}
	}
}

// ---- contains / indexOf / lastIndexOf -----------------------------------

func TestQStringListContainsQString(t *testing.T) {
	l := QStringListFromGoStrings([]string{"Apple", "apple", "banana"})
	defer l.Close()
	needle := NewQStringFromString("apple")
	defer needle.Close()

	if !l.ContainsQString(needle, true) {
		t.Fatal("expected case-sensitive contains")
	}
	upper := NewQStringFromString("APPLE")
	defer upper.Close()
	if l.ContainsQString(upper, true) {
		t.Fatal("did not expect case-sensitive contains of APPLE")
	}
	if !l.ContainsQString(upper, false) {
		t.Fatal("expected case-insensitive contains of APPLE")
	}
}

func TestQStringListContainsGoString(t *testing.T) {
	l := QStringListFromGoStrings([]string{"x", "y", "z"})
	defer l.Close()
	if !l.ContainsGoString("y", true) {
		t.Fatal("expected contains y")
	}
	if l.ContainsGoString("w", true) {
		t.Fatal("did not expect contains w")
	}
}

func TestQStringListIndexOfQString(t *testing.T) {
	l := QStringListFromGoStrings([]string{"a", "b", "a", "c"})
	defer l.Close()
	needle := NewQStringFromString("a")
	defer needle.Close()
	if got := l.IndexOfQString(needle, 0, true); got != 0 {
		t.Fatalf("IndexOfQString from 0: %d", got)
	}
	if got := l.IndexOfQString(needle, 1, true); got != 2 {
		t.Fatalf("IndexOfQString from 1: %d", got)
	}
	if got := l.IndexOfQString(needle, 3, true); got != -1 {
		t.Fatalf("IndexOfQString from 3: %d", got)
	}
}

func TestQStringListIndexOfGoString(t *testing.T) {
	l := QStringListFromGoStrings([]string{"Alpha", "beta", "alpha"})
	defer l.Close()
	if got := l.IndexOfGoString("alpha", 0, false); got != 0 {
		t.Fatalf("case-insensitive IndexOfGoString: %d", got)
	}
	if got := l.IndexOfGoString("alpha", 1, false); got != 2 {
		t.Fatalf("second case-insensitive match: %d", got)
	}
}

func TestQStringListIndexOfQRegularExpression(t *testing.T) {
	l := QStringListFromGoStrings([]string{"item-1", "item-2", "other", "item-3"})
	defer l.Close()
	re := NewQRegularExpressionFromString(`^item-\d+$`, NoPatternOption)
	defer re.Close()
	if got := l.IndexOfQRegularExpression(re, 0); got != 0 {
		t.Fatalf("IndexOfQRegularExpression from 0: %d", got)
	}
	if got := l.IndexOfQRegularExpression(re, 1); got != 1 {
		t.Fatalf("IndexOfQRegularExpression from 1: %d", got)
	}
	if got := l.IndexOfQRegularExpression(re, 2); got != 3 {
		t.Fatalf("IndexOfQRegularExpression from 2: %d", got)
	}
}

func TestQStringListLastIndexOfQString(t *testing.T) {
	l := QStringListFromGoStrings([]string{"a", "b", "a", "c", "a"})
	defer l.Close()
	needle := NewQStringFromString("a")
	defer needle.Close()
	if got := l.LastIndexOfQString(needle, -1, true); got != 4 {
		t.Fatalf("LastIndexOfQString -1: %d", got)
	}
	if got := l.LastIndexOfQString(needle, 3, true); got != 2 {
		t.Fatalf("LastIndexOfQString from 3: %d", got)
	}
	if got := l.LastIndexOfQString(needle, 1, true); got != 0 {
		t.Fatalf("LastIndexOfQString from 1: %d", got)
	}
}

func TestQStringListLastIndexOfGoString(t *testing.T) {
	l := QStringListFromGoStrings([]string{"X", "y", "x", "z"})
	defer l.Close()
	if got := l.LastIndexOfGoString("x", -1, false); got != 2 {
		t.Fatalf("case-insensitive LastIndexOfGoString: %d", got)
	}
}

func TestQStringListLastIndexOfQRegularExpression(t *testing.T) {
	l := QStringListFromGoStrings([]string{"item-1", "item-2", "other", "item-3"})
	defer l.Close()
	re := NewQRegularExpressionFromString(`^item-\d+$`, NoPatternOption)
	defer re.Close()
	if got := l.LastIndexOfQRegularExpression(re, -1); got != 3 {
		t.Fatalf("LastIndexOfQRegularExpression -1: %d", got)
	}
	if got := l.LastIndexOfQRegularExpression(re, 2); got != 1 {
		t.Fatalf("LastIndexOfQRegularExpression from 2: %d", got)
	}
	if got := l.LastIndexOfQRegularExpression(re, 0); got != 0 {
		t.Fatalf("LastIndexOfQRegularExpression from 0: %d", got)
	}
}
