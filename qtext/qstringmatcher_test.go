package qtext

import (
	"testing"

	"github.com/vadimlxx/go-qt6core/qdata"
)

func TestQStringMatcherDefault(t *testing.T) {
	m := NewQStringMatcher()
	defer m.Close()
	if m.PatternString() != "" {
		t.Fatalf("default pattern: %q", m.PatternString())
	}
	if m.CaseSensitivity() != qdata.CaseSensitive {
		t.Fatalf("default cs: %v", m.CaseSensitivity())
	}
}

func TestQStringMatcherFromString(t *testing.T) {
	m := NewQStringMatcherFromString("world", qdata.CaseSensitive)
	defer m.Close()
	if m.PatternString() != "world" {
		t.Fatalf("pattern: %q", m.PatternString())
	}
}

func TestQStringMatcherFromQString(t *testing.T) {
	p := NewQStringFromString("hello")
	defer p.Close()
	m := NewQStringMatcherFromQString(p, qdata.CaseSensitive)
	defer m.Close()
	if m.PatternString() != "hello" {
		t.Fatalf("pattern: %q", m.PatternString())
	}
}

func TestQStringMatcherFromQChars(t *testing.T) {
	units := []uint16{'a', 'b', 'c'}
	m := NewQStringMatcherFromQChars(units, qdata.CaseSensitive)
	defer m.Close()
	if m.PatternString() != "abc" {
		t.Fatalf("pattern: %q", m.PatternString())
	}
}

func TestQStringMatcherFromQStringView(t *testing.T) {
	src := NewQStringFromString("view")
	defer src.Close()
	v := src.ToQStringView()
	m := NewQStringMatcherFromQStringView(v, qdata.CaseSensitive)
	defer m.Close()
	if m.PatternString() != "view" {
		t.Fatalf("pattern: %q", m.PatternString())
	}
}

func TestQStringMatcherDispatcher(t *testing.T) {
	// string with default cs
	m1 := NewQStringMatcher("hello")
	defer m1.Close()
	if m1.CaseSensitivity() != qdata.CaseSensitive {
		t.Fatalf("expected CaseSensitive by default")
	}

	// string + cs
	m2 := NewQStringMatcher("hello", qdata.CaseInsensitive)
	defer m2.Close()
	if m2.CaseSensitivity() != qdata.CaseInsensitive {
		t.Fatalf("expected CaseInsensitive")
	}

	// Clone
	m3 := NewQStringMatcher(m2)
	defer m3.Close()
	if m3.PatternString() != "hello" {
		t.Fatalf("clone pattern: %q", m3.PatternString())
	}
}

func TestQStringMatcherIndexInString(t *testing.T) {
	m := NewQStringMatcherFromString("world", qdata.CaseSensitive)
	defer m.Close()

	subject := NewQStringFromString("hello world!")
	defer subject.Close()

	if got := m.IndexIn(subject, 0); got != 6 {
		t.Fatalf("IndexIn: %d", got)
	}
	if got := m.IndexIn(subject, 7); got != -1 {
		t.Fatalf("IndexIn from 7: %d", got)
	}
}

func TestQStringMatcherIndexInGoString(t *testing.T) {
	m := NewQStringMatcherFromString("world", qdata.CaseSensitive)
	defer m.Close()

	if got := m.IndexInString("hello world", 0); got != 6 {
		t.Fatalf("IndexInString: %d", got)
	}
	if got := m.IndexInString("hello world", 7); got != -1 {
		t.Fatalf("IndexInString from 7: %d", got)
	}
	if got := m.IndexInString("nope", 0); got != -1 {
		t.Fatalf("IndexInString not found: %d", got)
	}
}

func TestQStringMatcherIndexInQChars(t *testing.T) {
	m := NewQStringMatcherFromString("world", qdata.CaseSensitive)
	defer m.Close()

	subject := []uint16{'h', 'e', 'l', 'l', 'o', ' ', 'w', 'o', 'r', 'l', 'd'}
	if got := m.IndexInQChars(subject, 0); got != 6 {
		t.Fatalf("IndexInQChars: %d", got)
	}
}

func TestQStringMatcherIndexInView(t *testing.T) {
	m := NewQStringMatcherFromString("world", qdata.CaseSensitive)
	defer m.Close()

	src := NewQStringFromString("hello world")
	defer src.Close()
	v := src.ToQStringView()

	if got := m.IndexInView(v, 0); got != 6 {
		t.Fatalf("IndexInView: %d", got)
	}
}

func TestQStringMatcherCaseInsensitive(t *testing.T) {
	m := NewQStringMatcherFromString("world", qdata.CaseInsensitive)
	defer m.Close()

	if got := m.IndexInString("Hello WORLD", 0); got != 6 {
		t.Fatalf("IndexInString (CI): %d", got)
	}
}

func TestQStringMatcherCaseSensitive(t *testing.T) {
	m := NewQStringMatcherFromString("world", qdata.CaseSensitive)
	defer m.Close()

	if got := m.IndexInString("Hello WORLD", 0); got != -1 {
		t.Fatalf("IndexInString (CS): %d", got)
	}
}

func TestQStringMatcherSetPattern(t *testing.T) {
	m := NewQStringMatcherFromString("old", qdata.CaseSensitive)
	defer m.Close()

	m.SetPatternString("new")
	if m.PatternString() != "new" {
		t.Fatalf("after SetPatternString: %q", m.PatternString())
	}

	p := NewQStringFromString("third")
	defer p.Close()
	m.SetPattern(p)
	if m.PatternString() != "third" {
		t.Fatalf("after SetPattern: %q", m.PatternString())
	}
}

func TestQStringMatcherSetCaseSensitivity(t *testing.T) {
	m := NewQStringMatcherFromString("world", qdata.CaseSensitive)
	defer m.Close()

	if got := m.IndexInString("Hello WORLD", 0); got != -1 {
		t.Fatalf("expected -1 for CS")
	}

	m.SetCaseSensitivity(qdata.CaseInsensitive)
	if m.CaseSensitivity() != qdata.CaseInsensitive {
		t.Fatalf("SetCaseSensitivity failed")
	}
	if got := m.IndexInString("Hello WORLD", 0); got != 6 {
		t.Fatalf("after CI: %d", got)
	}
}

func TestQStringMatcherClone(t *testing.T) {
	m1 := NewQStringMatcherFromString("abc", qdata.CaseSensitive)
	defer m1.Close()

	m2 := m1.Clone()
	defer m2.Close()

	m2.SetPatternString("xyz")
	if m1.PatternString() != "abc" {
		t.Fatalf("original changed: %q", m1.PatternString())
	}
	if m2.PatternString() != "xyz" {
		t.Fatalf("clone not modified: %q", m2.PatternString())
	}
}

func TestQStringMatcherAssign(t *testing.T) {
	m1 := NewQStringMatcherFromString("abc", qdata.CaseSensitive)
	defer m1.Close()

	m2 := NewQStringMatcherFromString("xyz", qdata.CaseInsensitive)
	defer m2.Close()

	m1.Assign(m2)
	if m1.PatternString() != "xyz" {
		t.Fatalf("after Assign pattern: %q", m1.PatternString())
	}
	if m1.CaseSensitivity() != qdata.CaseInsensitive {
		t.Fatalf("after Assign cs: %v", m1.CaseSensitivity())
	}
}

func TestQStringMatcherPattern(t *testing.T) {
	m := NewQStringMatcherFromString("hello", qdata.CaseSensitive)
	defer m.Close()

	p := m.Pattern()
	defer p.Close()
	if p.GoString() != "hello" {
		t.Fatalf("Pattern: %q", p.GoString())
	}

	if m.PatternString() != "hello" {
		t.Fatalf("PatternString: %q", m.PatternString())
	}

	v := m.PatternView()
	if v.Size() != 5 {
		t.Fatalf("PatternView size: %d", v.Size())
	}
}

func TestQStringMatcherNilSubject(t *testing.T) {
	m := NewQStringMatcherFromString("abc", qdata.CaseSensitive)
	defer m.Close()

	if got := m.IndexIn(nil, 0); got != -1 {
		t.Fatalf("IndexIn(nil): %d", got)
	}
}

func TestQStringMatcherUseAfterClose(t *testing.T) {
	m := NewQStringMatcherFromString("abc", qdata.CaseSensitive)
	m.Close()
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on use after Close")
		}
	}()
	_ = m.PatternString()
}

func TestQStringMatcherCloseNilSafe(t *testing.T) {
	var m *QStringMatcher
	m.Close() // must not panic
}

func TestQStringMatcherCloseIdempotent(t *testing.T) {
	m := NewQStringMatcherFromString("abc", qdata.CaseSensitive)
	m.Close()
	m.Close() // second Close must not panic
}

func TestQStringMatcherUnicodePattern(t *testing.T) {
	// "мир" in Cyrillic, search in "привет мир".
	m := NewQStringMatcherFromString("мир", qdata.CaseSensitive)
	defer m.Close()

	got := m.IndexInString("привет мир", 0)
	// "привет " is 7 UTF-16 code units (each Cyrillic letter is 1).
	if got != 7 {
		t.Fatalf("IndexInString cyrillic: %d, want 7", got)
	}
}

func TestQStringMatcherSurrogatePair(t *testing.T) {
	// 😀 is U+1F600, 2 UTF-16 code units.
	m := NewQStringMatcherFromString("😀", qdata.CaseSensitive)
	defer m.Close()

	got := m.IndexInString("a😀b", 0)
	// "a" is 1 code unit, "😀" starts at index 1.
	if got != 1 {
		t.Fatalf("IndexInString emoji: %d, want 1", got)
	}
}