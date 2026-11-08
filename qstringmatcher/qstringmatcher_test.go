package qstringmatcher

import (
	"testing"
	"unicode/utf16"

	"github.com/vadimlxx/go-qt6core/qstring"
	"github.com/vadimlxx/go-qt6core/qstringview"
)

func TestDefault(t *testing.T) {
	m := New()
	defer m.Close()
	// No crash, no panic. PatternView of default (empty pattern) is fine.
	v := m.PatternView()
	if v.Size() != 0 {
		t.Fatalf("default PatternView size: %d, want 0", v.Size())
	}
}

func TestFromQStringBasic(t *testing.T) {
	pattern := qstring.FromString("quick")
	defer pattern.Close()
	m := FromQString(pattern, true)
	defer m.Close()

	haystack := qstring.FromString("the quick brown fox")
	defer haystack.Close()

	if got := m.IndexIn(haystack, 0); got != 4 {
		t.Fatalf("IndexIn: %d, want 4", got)
	}
	if got := m.IndexIn(haystack, 5); got != -1 {
		t.Fatalf("IndexIn from 5: %d, want -1", got)
	}
}

func TestFromQStringView(t *testing.T) {
	sv := qstringview.FromUtf16(utf16.Encode([]rune("quick")))
	m := FromQStringView(sv, true)
	defer m.Close()

	haystack := qstringview.FromUtf16(utf16.Encode([]rune("the quick brown fox")))
	if got := m.IndexInView(haystack, 0); got != 4 {
		t.Fatalf("IndexInView: %d, want 4", got)
	}
}

func TestFromUtf16(t *testing.T) {
	m := FromUtf16(utf16.Encode([]rune("quick")), true)
	defer m.Close()

	haystack := utf16.Encode([]rune("the quick brown fox"))
	if got := m.IndexInUtf16(haystack, 0); got != 4 {
		t.Fatalf("IndexInUtf16: %d, want 4", got)
	}
}

func TestCaseSensitivity(t *testing.T) {
	pattern := qstring.FromString("QUICK")
	defer pattern.Close()

	// case-sensitive: no match
	cs := FromQString(pattern, true)
	defer cs.Close()
	haystack := qstring.FromString("the quick brown fox")
	defer haystack.Close()
	if got := cs.IndexIn(haystack, 0); got != -1 {
		t.Fatalf("case-sensitive IndexIn: %d, want -1", got)
	}

	// case-insensitive: match at 4
	ci := FromQString(pattern, false)
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

func TestSetPattern(t *testing.T) {
	m := New()
	defer m.Close()

	haystack := qstring.FromString("the quick brown fox")
	defer haystack.Close()

	m.SetPattern(qstring.FromString("brown"))
	if got := m.IndexIn(haystack, 0); got != 10 {
		t.Fatalf("IndexIn(brown): %d, want 10", got)
	}

	m.SetPattern(qstring.FromString("zzz"))
	if got := m.IndexIn(haystack, 0); got != -1 {
		t.Fatalf("IndexIn(zzz): %d, want -1", got)
	}
}

func TestSetCaseSensitivity(t *testing.T) {
	pattern := qstring.FromString("QUICK")
	defer pattern.Close()
	m := FromQString(pattern, true)
	defer m.Close()

	haystack := qstring.FromString("the quick brown fox")
	defer haystack.Close()

	if got := m.IndexIn(haystack, 0); got != -1 {
		t.Fatalf("case-sensitive: %d, want -1", got)
	}

	m.SetCaseSensitivity(false)
	if got := m.IndexIn(haystack, 0); got != 4 {
		t.Fatalf("after SetCaseSensitivity(false): %d, want 4", got)
	}
}

func TestClone(t *testing.T) {
	pattern := qstring.FromString("fox")
	defer pattern.Close()
	orig := FromQString(pattern, true)
	defer orig.Close()

	cpy := orig.Clone()
	defer cpy.Close()

	haystack := qstring.FromString("the quick brown fox")
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

func TestAssign(t *testing.T) {
	src := FromQString(qstring.FromString("fox"), true)
	defer src.Close()
	dst := FromQString(qstring.FromString("zzz"), false)
	defer dst.Close()

	dst.Assign(src)

	haystack := qstring.FromString("the quick brown fox")
	defer haystack.Close()

	if got := dst.IndexIn(haystack, 0); got != 16 {
		t.Fatalf("after Assign: %d, want 16", got)
	}
	if !dst.CaseSensitivity() {
		t.Fatal("after Assign, CaseSensitivity must be true")
	}
}

func TestPattern(t *testing.T) {
	pattern := qstring.FromString("hello")
	defer pattern.Close()
	m := FromQString(pattern, true)
	defer m.Close()

	got := m.Pattern()
	defer got.Close()
	if got.GoString() != "hello" {
		t.Fatalf("Pattern: %q", got.GoString())
	}
}

func TestPatternView(t *testing.T) {
	pattern := qstring.FromString("hello")
	defer pattern.Close()
	m := FromQString(pattern, true)
	defer m.Close()

	v := m.PatternView()
	if v.GoString() != "hello" {
		t.Fatalf("PatternView: %q", v.GoString())
	}
}

func TestIndexInMultiByte(t *testing.T) {
	// "😀" = U+1F600 = 2 UTF-16 code units
	pattern := qstring.FromString("😀")
	defer pattern.Close()
	m := FromQString(pattern, true)
	defer m.Close()

	haystack := qstring.FromString("a😀b")
	defer haystack.Close()

	if got := m.IndexIn(haystack, 0); got != 1 {
		t.Fatalf("IndexIn(😀): %d, want 1", got)
	}
}

func TestEmptyPattern(t *testing.T) {
	m := New()
	defer m.Close()

	haystack := qstring.FromString("hello")
	defer haystack.Close()

	// Empty pattern matches at position 0.
	if got := m.IndexIn(haystack, 0); got != 0 {
		t.Fatalf("empty pattern IndexIn: %d, want 0", got)
	}
}
