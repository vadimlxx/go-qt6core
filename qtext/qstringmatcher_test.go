package qtext

import (
	"testing"
	"unicode/utf16"
)

func TestQStringMatcherDefault(t *testing.T) {
	m := NewQStringMatcher()
	defer m.Close()
	// No crash, no panic. PatternView of default (empty pattern) is fine.
	v := m.PatternView()
	if v.Size() != 0 {
		t.Fatalf("default PatternView size: %d, want 0", v.Size())
	}
}

func TestQStringMatcherFromQStringBasic(t *testing.T) {
	pattern := QStringFromString("quick")
	defer pattern.Close()
	m := QStringMatcherFromQString(pattern, true)
	defer m.Close()

	haystack := QStringFromString("the quick brown fox")
	defer haystack.Close()

	if got := m.IndexIn(haystack, 0); got != 4 {
		t.Fatalf("IndexIn: %d, want 4", got)
	}
	if got := m.IndexIn(haystack, 5); got != -1 {
		t.Fatalf("IndexIn from 5: %d, want -1", got)
	}
}

func TestQStringMatcherFromQStringView(t *testing.T) {
	sv := QStringViewFromUtf16(utf16.Encode([]rune("quick")))
	m := QStringMatcherFromQStringView(sv, true)
	defer m.Close()

	haystack := QStringViewFromUtf16(utf16.Encode([]rune("the quick brown fox")))
	if got := m.IndexInView(haystack, 0); got != 4 {
		t.Fatalf("IndexInView: %d, want 4", got)
	}
}

func TestQStringMatcherFromUtf16(t *testing.T) {
	m := QStringMatcherFromUtf16(utf16.Encode([]rune("quick")), true)
	defer m.Close()

	haystack := utf16.Encode([]rune("the quick brown fox"))
	if got := m.IndexInUtf16(haystack, 0); got != 4 {
		t.Fatalf("IndexInUtf16: %d, want 4", got)
	}
}

func TestQStringMatcherCaseSensitivity(t *testing.T) {
	pattern := QStringFromString("QUICK")
	defer pattern.Close()

	// case-sensitive: no match
	cs := QStringMatcherFromQString(pattern, true)
	defer cs.Close()
	haystack := QStringFromString("the quick brown fox")
	defer haystack.Close()
	if got := cs.IndexIn(haystack, 0); got != -1 {
		t.Fatalf("case-sensitive IndexIn: %d, want -1", got)
	}

	// case-insensitive: match at 4
	ci := QStringMatcherFromQString(pattern, false)
	defer ci.Close()
	if got := ci.IndexIn(haystack, 0); got != 4 {
		t.Fatalf("case-insensitive IndexIn: %d, want 4", got)
	}

	if !cs.CaseSensitivity() {
		t.Fatal("CaseSensitivity must be true for cs matcher")
	}
	if ci.CaseSensitivity() {
		t.Fatal("CaseSensitivity must be false for ci matcher")
	}
}

func TestQStringMatcherSetPattern(t *testing.T) {
	m := NewQStringMatcher()
	defer m.Close()

	haystack := QStringFromString("the quick brown fox")
	defer haystack.Close()

	m.SetPattern(QStringFromString("brown"))
	if got := m.IndexIn(haystack, 0); got != 10 {
		t.Fatalf("IndexIn(brown): %d, want 10", got)
	}

	m.SetPattern(QStringFromString("zzz"))
	if got := m.IndexIn(haystack, 0); got != -1 {
		t.Fatalf("IndexIn(zzz): %d, want -1", got)
	}
}

func TestQStringMatcherSetCaseSensitivity(t *testing.T) {
	pattern := QStringFromString("QUICK")
	defer pattern.Close()
	m := QStringMatcherFromQString(pattern, true)
	defer m.Close()

	haystack := QStringFromString("the quick brown fox")
	defer haystack.Close()

	if got := m.IndexIn(haystack, 0); got != -1 {
		t.Fatalf("case-sensitive: %d, want -1", got)
	}

	m.SetCaseSensitivity(false)
	if got := m.IndexIn(haystack, 0); got != 4 {
		t.Fatalf("after SetCaseSensitivity(false): %d, want 4", got)
	}
}

func TestQStringMatcherClone(t *testing.T) {
	pattern := QStringFromString("fox")
	defer pattern.Close()
	orig := QStringMatcherFromQString(pattern, true)
	defer orig.Close()

	cpy := orig.Clone()
	defer cpy.Close()

	haystack := QStringFromString("the quick brown fox")
	defer haystack.Close()

	if got := cpy.IndexIn(haystack, 0); got != 16 {
		t.Fatalf("clone IndexIn: %d, want 16", got)
	}

	// Mutating the clone does not affect the original.
	cpy.SetCaseSensitivity(false)
	if !orig.CaseSensitivity() {
		t.Fatal("orig.CaseSensitivity must remain true after clone mutation")
	}
}

func TestQStringMatcherAssign(t *testing.T) {
	src := QStringMatcherFromQString(QStringFromString("fox"), true)
	defer src.Close()
	dst := QStringMatcherFromQString(QStringFromString("zzz"), false)
	defer dst.Close()

	dst.Assign(src)

	haystack := QStringFromString("the quick brown fox")
	defer haystack.Close()

	if got := dst.IndexIn(haystack, 0); got != 16 {
		t.Fatalf("after Assign: %d, want 16", got)
	}
	if !dst.CaseSensitivity() {
		t.Fatal("after Assign, CaseSensitivity must be true")
	}
}

func TestQStringMatcherPattern(t *testing.T) {
	pattern := QStringFromString("hello")
	defer pattern.Close()
	m := QStringMatcherFromQString(pattern, true)
	defer m.Close()

	got := m.Pattern()
	defer got.Close()
	if got.GoString() != "hello" {
		t.Fatalf("Pattern: %q", got.GoString())
	}
}

func TestQStringMatcherPatternView(t *testing.T) {
	pattern := QStringFromString("hello")
	defer pattern.Close()
	m := QStringMatcherFromQString(pattern, true)
	defer m.Close()

	v := m.PatternView()
	if v.GoString() != "hello" {
		t.Fatalf("PatternView: %q", v.GoString())
	}
}

func TestQStringMatcherIndexInMultiByte(t *testing.T) {
	// "😀" = U+1F600 = 2 UTF-16 code units
	pattern := QStringFromString("😀")
	defer pattern.Close()
	m := QStringMatcherFromQString(pattern, true)
	defer m.Close()

	haystack := QStringFromString("a😀b")
	defer haystack.Close()

	if got := m.IndexIn(haystack, 0); got != 1 {
		t.Fatalf("IndexIn(😀): %d, want 1", got)
	}
}

func TestQStringMatcherEmptyPattern(t *testing.T) {
	m := NewQStringMatcher()
	defer m.Close()

	haystack := QStringFromString("hello")
	defer haystack.Close()

	// Empty pattern matches at position 0.
	if got := m.IndexIn(haystack, 0); got != 0 {
		t.Fatalf("empty pattern IndexIn: %d, want 0", got)
	}
}
