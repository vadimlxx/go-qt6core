package qtext

import (
	"testing"

	"github.com/vadimlxx/go-qt6core/qdata"
)

func TestQRegularExpressionBasic(t *testing.T) {
	re := NewQRegularExpressionFromString(`\d+`, NoPatternOption)
	defer re.Close()
	if !re.IsValid() {
		t.Fatal("expected valid regex")
	}
	if re.CaptureCount() != 0 {
		t.Fatalf("expected 0 captures, got %d", re.CaptureCount())
	}
	if re.Pattern().String() != `\d+` {
		t.Fatalf("unexpected pattern: %q", re.Pattern().String())
	}
}

func TestQRegularExpressionDefaultConstructor(t *testing.T) {
	re := NewQRegularExpression()
	defer re.Close()
	// Qt6: empty pattern is a valid regular expression.
	if !re.IsValid() {
		t.Fatalf("expected empty pattern to be valid, error: %s", re.ErrorString().String())
	}
}

func TestQRegularExpressionClone(t *testing.T) {
	re := NewQRegularExpressionFromString(`abc`, NoPatternOption)
	defer re.Close()
	cl := re.Clone()
	defer cl.Close()
	if !cl.IsValid() || cl.Pattern().String() != "abc" {
		t.Fatalf("unexpected clone: valid=%v pattern=%q", cl.IsValid(), cl.Pattern().String())
	}
}

func TestQRegularExpressionSetPattern(t *testing.T) {
	re := NewQRegularExpression()
	defer re.Close()
	p := QStringFromString(`\w+`)
	defer p.Close()
	re.SetPattern(p)
	if !re.IsValid() {
		t.Fatal("expected valid after SetPattern")
	}
	if re.Pattern().String() != `\w+` {
		t.Fatalf("unexpected pattern: %q", re.Pattern().String())
	}
}

func TestQRegularExpressionPatternOptions(t *testing.T) {
	re := NewQRegularExpressionFromString(`abc`, CaseInsensitiveOption|MultilineOption)
	defer re.Close()
	opts := re.PatternOptions()
	if opts&CaseInsensitiveOption == 0 {
		t.Fatal("expected CaseInsensitiveOption")
	}
	if opts&MultilineOption == 0 {
		t.Fatal("expected MultilineOption")
	}
	re.SetPatternOptions(DotMatchesEverythingOption)
	if re.PatternOptions() != DotMatchesEverythingOption {
		t.Fatalf("unexpected options: %d", re.PatternOptions())
	}
}

func TestQRegularExpressionInvalidPattern(t *testing.T) {
	re := NewQRegularExpressionFromString(`[`, NoPatternOption)
	defer re.Close()
	if re.IsValid() {
		t.Fatal("expected invalid regex")
	}
	if re.ErrorString().String() == "" {
		t.Fatal("expected non-empty error string")
	}
	if re.PatternErrorOffset() < 0 {
		t.Fatal("expected non-negative error offset")
	}
}

func TestQRegularExpressionEscape(t *testing.T) {
	s := QStringFromString(`\d+`)
	defer s.Close()
	esc := QRegularExpressionEscape(s)
	defer esc.Close()
	if esc.String() != `\\d\+` {
		t.Fatalf("unexpected escape: %q", esc.String())
	}
}

func TestQRegularExpressionWildcard(t *testing.T) {
	s := QStringFromString(`*.txt`)
	defer s.Close()
	pattern := QRegularExpressionWildcardToRegularExpression(s, DefaultWildcardConversion)
	defer pattern.Close()
	if pattern.IsEmpty() {
		t.Fatal("wildcard conversion returned empty pattern")
	}
	re := NewQRegularExpressionFromPattern(pattern, NoPatternOption)
	defer re.Close()
	if !re.IsValid() {
		t.Fatalf("wildcard regex invalid: %s", re.ErrorString().String())
	}
}

func TestQRegularExpressionAnchoredPattern(t *testing.T) {
	s := QStringFromString(`\d+`)
	defer s.Close()
	ap := QRegularExpressionAnchoredPattern(s)
	defer ap.Close()
	if ap.String() == "" {
		t.Fatal("expected non-empty anchored pattern")
	}
}

func TestQRegularExpressionFromWildcard(t *testing.T) {
	s := QStringFromString(`*.txt`)
	defer s.Close()
	re := QRegularExpressionFromWildcard(s, true, DefaultWildcardConversion)
	defer re.Close()
	if !re.IsValid() {
		t.Fatal("fromWildcard regex invalid")
	}
}

// ---- QString ↔ QRegularExpression interop ------------------------------

func TestQStringIndexOfQRegularExpression(t *testing.T) {
	s := QStringFromString("abc 123 def 456")
	// Use word-bounded \b\d+\b rather than bare \d+: QString::count(QRegularExpression)
	// counts OVERLAPPING matches in Qt, so a bare \d+ here would report 6 matches
	// ("123", "23", "3", "456", "56", "6") instead of 2. See the dedicated
	// TestQStringCountQRegularExpressionOverlapping below.
	re := NewQRegularExpressionFromString(`\b\d+\b`, NoPatternOption)
	defer s.Close()
	defer re.Close()

	if got := s.IndexOfQRegularExpression(re, 0); got != 4 {
		t.Fatalf("IndexOfQRegularExpression: %d", got)
	}
	if got := s.LastIndexOfQRegularExpression(re, -1); got != 12 {
		t.Fatalf("LastIndexOfQRegularExpression: %d", got)
	}
	if got := s.CountQRegularExpression(re); got != 2 {
		t.Fatalf("CountQRegularExpression: %d", got)
	}
	if !s.ContainsQRegularExpression(re) {
		t.Fatal("ContainsQRegularExpression: expected true")
	}
}

// TestQStringCountQRegularExpressionOverlapping documents Qt's overlapping-match
// semantics for QString::count(const QRegularExpression&): \d+ over "123" yields
// three matches ("123", "23", "3"), not one.
func TestQStringCountQRegularExpressionOverlapping(t *testing.T) {
	s := QStringFromString("abc 123 def 456")
	re := NewQRegularExpressionFromString(`\d+`, NoPatternOption)
	defer s.Close()
	defer re.Close()

	if got := s.CountQRegularExpression(re); got != 6 {
		t.Fatalf("expected overlapping count 6, got %d", got)
	}

	single := QStringFromString("123")
	singleRe := NewQRegularExpressionFromString(`\d+`, NoPatternOption)
	defer single.Close()
	defer singleRe.Close()

	if got := single.CountQRegularExpression(singleRe); got != 3 {
		t.Fatalf(`expected overlapping count 3 for "123", got %d`, got)
	}
}

func TestQStringIndexOfQRegularExpressionWithMatch(t *testing.T) {
	s := QStringFromString("abc 123 def")
	re := NewQRegularExpressionFromString(`(\d+)`, NoPatternOption)
	defer s.Close()
	defer re.Close()

	m := NewQRegularExpressionMatch()
	defer m.Close()
	idx := s.IndexOfQRegularExpressionWithMatch(re, 0, m)
	if idx != 4 {
		t.Fatalf("index: %d", idx)
	}
	if !m.HasMatch() {
		t.Fatal("expected match")
	}
	c := m.CapturedInt(1)
	defer c.Close()
	if c.String() != "123" {
		t.Fatalf("captured: %q", c.String())
	}
}

func TestQStringContainsQRegularExpressionWithMatch(t *testing.T) {
	s := QStringFromString("abc 42")
	re := NewQRegularExpressionFromString(`\d+`, NoPatternOption)
	defer s.Close()
	defer re.Close()

	m := NewQRegularExpressionMatch()
	defer m.Close()
	if !s.ContainsQRegularExpressionWithMatch(re, m) {
		t.Fatal("expected contains")
	}
	c := m.CapturedInt(0)
	defer c.Close()
	if c.String() != "42" {
		t.Fatalf("captured: %q", c.String())
	}
}

func TestQStringSectionQRegularExpression(t *testing.T) {
	s := QStringFromString("a1b22c333d")
	re := NewQRegularExpressionFromString(`\d+`, NoPatternOption)
	defer s.Close()
	defer re.Close()

	// Fields: "a", "b", "c", "d"; separators: "1", "22", "333".
	// Section 1..2 preserves the original text of the separator between
	// fields 1 and 2, so the result is "b" + "22" + "c" == "b22c".
	sec := s.SectionQRegularExpression(re, 1, 2, SectionDefault)
	defer sec.Close()
	if sec.GoString() != "b22c" {
		t.Fatalf("SectionQRegularExpression: %q", sec.GoString())
	}
}

func TestQStringReplaceQRegularExpression(t *testing.T) {
	s := QStringFromString("a1b22c333")
	re := NewQRegularExpressionFromString(`\d+`, NoPatternOption)
	defer s.Close()
	defer re.Close()

	s.ReplaceQRegularExpressionString(re, "#")
	if s.GoString() != "a#b#c#" {
		t.Fatalf("ReplaceQRegularExpressionString: %q", s.GoString())
	}
}

func TestQStringRemoveQRegularExpression(t *testing.T) {
	s := QStringFromString("a1b22c333d")
	re := NewQRegularExpressionFromString(`\d+`, NoPatternOption)
	defer s.Close()
	defer re.Close()

	s.RemoveQRegularExpression(re)
	if s.GoString() != "abcd" {
		t.Fatalf("RemoveQRegularExpression: %q", s.GoString())
	}
}

func TestQStringSplitQRegularExpression(t *testing.T) {
	s := QStringFromString("a1b22c333d")
	re := NewQRegularExpressionFromString(`\d+`, NoPatternOption)
	defer s.Close()
	defer re.Close()

	parts := s.SplitQRegularExpression(re, qdata.KeepEmptyParts)
	if len(parts) != 4 {
		t.Fatalf("SplitQRegularExpression len: %d", len(parts))
	}
	for _, p := range parts {
		defer p.Close()
	}
	want := []string{"a", "b", "c", "d"}
	for i, w := range want {
		if parts[i].GoString() != w {
			t.Fatalf("part[%d]: got %q want %q", i, parts[i].GoString(), w)
		}
	}
}
