package qdatetime

import (
	"runtime"
	"testing"

	"github.com/vadimlxx/go-qt6core/qdata"
	"github.com/vadimlxx/go-qt6core/qtext"
)

// newTestCollator returns a collator bound to en_US, which makes the ASCII
// comparison expectations below deterministic.
func newTestCollator(t *testing.T) *QCollator {
	t.Helper()
	loc := NewQLocaleFromName("en_US")
	if loc == nil {
		t.Fatal("NewQLocaleFromName(en_US) returned nil")
	}
	defer loc.Close()
	c := NewQCollator(loc)
	if c == nil {
		t.Fatal("NewQCollator(locale) returned nil")
	}
	t.Cleanup(c.Close)
	return c
}

func TestQCollatorDefault(t *testing.T) {
	c := NewQCollator()
	if c == nil {
		t.Fatal("NewQCollator() returned nil")
	}
	defer c.Close()

	if got := c.CaseSensitivity(); got != qdata.CaseSensitive {
		t.Fatalf("default CaseSensitivity = %v, want CaseSensitive", got)
	}
	if c.NumericMode() {
		t.Fatal("default NumericMode must be false")
	}
	if c.IgnorePunctuation() {
		t.Fatal("default IgnorePunctuation must be false")
	}
}

func TestQCollatorFromLocale(t *testing.T) {
	loc := NewQLocaleFromName("de_DE")
	if loc == nil {
		t.Fatal("NewQLocaleFromName(de_DE) returned nil")
	}
	defer loc.Close()

	c := NewQCollatorFromLocale(loc)
	if c == nil {
		t.Fatal("NewQCollatorFromLocale returned nil")
	}
	defer c.Close()

	got := c.Locale()
	if got == nil {
		t.Fatal("Locale() returned nil")
	}
	defer got.Close()
	if got.Language() != German {
		t.Fatalf("Locale().Language() = %v, want German", got.Language())
	}
}

func TestQCollatorCopyAndSwap(t *testing.T) {
	a := newTestCollator(t)
	a.SetNumericMode(true)

	b := NewQCollator(a)
	if b == nil {
		t.Fatal("NewQCollator(collator) returned nil")
	}
	defer b.Close()
	if !b.NumericMode() {
		t.Fatal("copied collator lost NumericMode")
	}

	c := NewQCollatorCopy(a)
	if c == nil {
		t.Fatal("NewQCollatorCopy returned nil")
	}
	defer c.Close()

	// Swap must exchange option state.
	c.SetNumericMode(false)
	a.Swap(c)
	if a.NumericMode() {
		t.Fatal("after Swap, a.NumericMode must be false")
	}
	if !c.NumericMode() {
		t.Fatal("after Swap, c.NumericMode must be true")
	}
}

func TestQCollatorCaseSensitivity(t *testing.T) {
	c := newTestCollator(t)

	c.SetCaseSensitivity(qdata.CaseInsensitive)
	if got := c.CaseSensitivity(); got != qdata.CaseInsensitive {
		t.Fatalf("CaseSensitivity = %v, want CaseInsensitive", got)
	}
	if got := c.CompareString("A", "a"); got != 0 {
		t.Fatalf("case-insensitive Compare(A, a) = %d, want 0", got)
	}

	c.SetCaseSensitivity(qdata.CaseSensitive)
	if got := c.CompareString("A", "a"); got == 0 {
		t.Fatal("case-sensitive Compare(A, a) must be non-zero")
	}
}

func TestQCollatorNumericMode(t *testing.T) {
	c := newTestCollator(t)

	if got := c.CompareString("file10", "file2"); got >= 0 {
		t.Fatalf("lexical Compare(file10, file2) = %d, want < 0", got)
	}

	c.SetNumericMode(true)
	if got := c.CompareString("file10", "file2"); got <= 0 {
		t.Fatalf("numeric Compare(file10, file2) = %d, want > 0", got)
	}
}

func TestQCollatorIgnorePunctuation(t *testing.T) {
	c := newTestCollator(t)
	c.SetIgnorePunctuation(true)

	if got := c.CompareString("ab", "a-b"); got != 0 {
		t.Fatalf("ignore-punctuation Compare(ab, a-b) = %d, want 0", got)
	}
}

func TestQCollatorCompareVariants(t *testing.T) {
	c := newTestCollator(t)

	q1 := qtext.NewQStringFromString("apple")
	defer q1.Close()
	q2 := qtext.NewQStringFromString("banana")
	defer q2.Close()

	want := c.Compare(q1, q2)
	if want >= 0 {
		t.Fatalf("Compare(apple, banana) = %d, want < 0", want)
	}
	if got := c.CompareString("apple", "banana"); got != want {
		t.Fatalf("CompareString = %d, want %d", got, want)
	}
	if got := c.CompareView(q1.ToQStringView(), q2.ToQStringView()); got != want {
		t.Fatalf("CompareView = %d, want %d", got, want)
	}
	if got := c.CompareViewString("apple", "banana"); got != want {
		t.Fatalf("CompareViewString = %d, want %d", got, want)
	}
}

func TestQCollatorLess(t *testing.T) {
	c := newTestCollator(t)

	if !c.LessString("apple", "banana") {
		t.Fatal("LessString(apple, banana) = false")
	}
	if c.LessString("banana", "apple") {
		t.Fatal("LessString(banana, apple) = true")
	}
	if c.LessViewString("banana", "apple") {
		t.Fatal("LessViewString(banana, apple) = true")
	}

	q1 := qtext.NewQStringFromString("apple")
	defer q1.Close()
	q2 := qtext.NewQStringFromString("banana")
	defer q2.Close()
	if !c.Less(q1, q2) {
		t.Fatal("Less(apple, banana) = false")
	}
	if !c.LessView(q1.ToQStringView(), q2.ToQStringView()) {
		t.Fatal("LessView(apple, banana) = false")
	}
}

func TestQCollatorSortKey(t *testing.T) {
	c := newTestCollator(t)

	k1 := c.SortKeyString("apple")
	if k1 == nil {
		t.Fatal("SortKeyString returned nil")
	}
	defer k1.Close()
	k2 := c.SortKeyString("banana")
	if k2 == nil {
		t.Fatal("SortKeyString returned nil")
	}
	defer k2.Close()

	if !k1.Less(k2) {
		t.Fatal("apple sort key must be less than banana sort key")
	}
	if k1.Compare(k2) >= 0 {
		t.Fatalf("Compare(appleKey, bananaKey) = %d, want < 0", k1.Compare(k2))
	}
	equal := c.SortKeyString("apple")
	if equal == nil {
		t.Fatal("SortKeyString returned nil")
	}
	defer equal.Close()
	if got := k1.Compare(equal); got != 0 {
		t.Fatalf("equal strings produced compare = %d, want 0", got)
	}
}

func TestQCollatorStaticHelpers(t *testing.T) {
	if got := DefaultCompareString("apple", "banana"); got >= 0 {
		t.Fatalf("DefaultCompareString = %d, want < 0", got)
	}

	q1 := qtext.NewQStringFromString("apple")
	defer q1.Close()
	q2 := qtext.NewQStringFromString("banana")
	defer q2.Close()
	if got := DefaultCompare(q1.ToQStringView(), q2.ToQStringView()); got >= 0 {
		t.Fatalf("DefaultCompare = %d, want < 0", got)
	}

	k1 := DefaultSortKeyString("apple")
	defer k1.Close()
	k2 := DefaultSortKeyString("banana")
	defer k2.Close()
	if !k1.Less(k2) {
		t.Fatal("DefaultSortKey: apple key must be less than banana key")
	}

	kv := DefaultSortKey(q1.ToQStringView())
	defer kv.Close()
	if !kv.Less(k2) {
		t.Fatal("DefaultSortKey(view): apple key must be less than banana key")
	}
}

func TestQCollatorClose(t *testing.T) {
	c := NewQCollator()
	c.Close()
	c.Close() // idempotent

	defer func() {
		if recover() == nil {
			t.Fatal("Compare on a closed collator must panic")
		}
	}()
	c.CompareString("a", "b")
}

func TestQCollatorGenericConstructorPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("NewQCollator with unsupported argument must panic")
		}
	}()
	NewQCollator("not a collator")
}

func TestQCollatorCleanup(t *testing.T) {
	for i := 0; i < 32; i++ {
		c := NewQCollator()
		if c == nil {
			t.Fatalf("iteration %d: nil collator", i)
		}
		c.SetNumericMode(i%2 == 0)
	}
	runtime.GC()
	runtime.GC()
}
