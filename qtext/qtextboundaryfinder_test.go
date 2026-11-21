package qtext

import (
	"testing"
)

// ---- construction --------------------------------------------------------

func TestQTextBoundaryFinderDefault(t *testing.T) {
	f := NewQTextBoundaryFinder()
	defer f.Close()
	if f.IsValid() {
		t.Fatal("default-constructed finder should be invalid")
	}
}

func TestQTextBoundaryFinderFromQString(t *testing.T) {
	s := NewQStringFromString("hello world")
	defer s.Close()
	f := QTextBoundaryFinderFromQString(Word, s)
	defer f.Close()
	if !f.IsValid() {
		t.Fatal("expected valid finder")
	}
	if f.Type() != Word {
		t.Fatalf("Type: %d, want %d", f.Type(), Word)
	}
	got := f.String()
	defer got.Close()
	if got.String() != "hello world" {
		t.Fatalf("String: %q", got.String())
	}
}

func TestQTextBoundaryFinderFromString(t *testing.T) {
	f := QTextBoundaryFinderFromString(Grapheme, "héllo")
	defer f.Close()
	if !f.IsValid() {
		t.Fatal("expected valid finder")
	}
	if f.Type() != Grapheme {
		t.Fatalf("Type: %d", f.Type())
	}
}

func TestQTextBoundaryFinderFromView(t *testing.T) {
	v := QStringViewFromString("a b c")
	f := QTextBoundaryFinderFromView(Word, v)
	defer f.Close()
	if !f.IsValid() {
		t.Fatal("expected valid finder")
	}
	if f.Type() != Word {
		t.Fatalf("Type: %d", f.Type())
	}
}

func TestQTextBoundaryFinderClone(t *testing.T) {
	s := NewQStringFromString("ab cd")
	defer s.Close()
	f := QTextBoundaryFinderFromQString(Word, s)
	defer f.Close()

	f.ToNextBoundary() // now at pos 2 ("ab")
	c := f.Clone()
	defer c.Close()
	if !c.IsValid() {
		t.Fatal("clone should be valid")
	}
	if c.Position() != f.Position() {
		t.Fatalf("clone position %d != original %d", c.Position(), f.Position())
	}
}

func TestQTextBoundaryFinderSwap(t *testing.T) {
	a := QTextBoundaryFinderFromString(Word, "ab cd")
	b := QTextBoundaryFinderFromString(Line, "xy z")
	defer a.Close()
	defer b.Close()
	a.Swap(b)
	if a.Type() != Line || b.Type() != Word {
		t.Fatalf("swap failed: a=%d b=%d", a.Type(), b.Type())
	}
}

// ---- navigation ----------------------------------------------------------

func TestQTextBoundaryFinderToStartToEnd(t *testing.T) {
	s := NewQStringFromString("hello")
	defer s.Close()
	f := QTextBoundaryFinderFromQString(Word, s)
	defer f.Close()

	f.ToStart()
	if f.Position() != 0 {
		t.Fatalf("after ToStart: %d", f.Position())
	}
	f.ToEnd()
	if f.Position() != 5 {
		t.Fatalf("after ToEnd: %d", f.Position())
	}
}

func TestQTextBoundaryFinderSetPosition(t *testing.T) {
	f := QTextBoundaryFinderFromString(Word, "hello world")
	defer f.Close()
	f.SetPosition(3)
	if f.Position() != 3 {
		t.Fatalf("SetPosition: %d", f.Position())
	}
}

func TestQTextBoundaryFinderNextPrevBoundary(t *testing.T) {
	f := QTextBoundaryFinderFromString(Word, "a bb ccc")
	defer f.Close()
	f.ToStart()

	// Word boundaries at 0, 1, 2, 4, 5, 6, 7, 9 (varies by UAX #29 rules;
	// only sanity-check that positions increase and are >= 0).
	prev := f.Position()
	for {
		next := f.ToNextBoundary()
		if next == -1 {
			break
		}
		if next < prev {
			t.Fatalf("next(%d) < prev(%d)", next, prev)
		}
		prev = next
	}

	f.ToEnd()
	// Walk back; positions must strictly decrease until -1.
	prev = f.Position()
	for {
		next := f.ToPreviousBoundary()
		if next == -1 {
			break
		}
		if next > prev {
			t.Fatalf("prev(%d) > prev(%d)", next, prev)
		}
		prev = next
	}
}

// ---- predicates ----------------------------------------------------------

func TestQTextBoundaryFinderIsAtBoundary(t *testing.T) {
	f := QTextBoundaryFinderFromString(Word, "a b")
	defer f.Close()

	f.ToStart()
	if !f.IsAtBoundary() {
		t.Fatal("start should be at boundary")
	}
	f.SetPosition(1) // space; word boundary per UAX #29
	if !f.IsAtBoundary() {
		t.Fatal("position after 'a' should be a boundary")
	}
}

func TestQTextBoundaryFinderBoundaryReasons(t *testing.T) {
	f := QTextBoundaryFinderFromString(Word, "abc\ndef")
	defer f.Close()

	f.ToStart()
	r := f.BoundaryReasons()
	if r == NotAtBoundary {
		t.Fatal("expected boundary at start")
	}
}
