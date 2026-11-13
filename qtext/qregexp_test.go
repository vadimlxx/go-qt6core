package qtext

import "testing"

func TestQRegularExpressionBasic(t *testing.T) {
	re := NewQRegularExpressionFromString(`\d+`, NoPatternOption)
	if !re.IsValid() {
		t.Fatal("expected valid regex")
	}
	if re.CaptureCount() != 0 {
		t.Fatalf("expected 0 captures, got %d", re.CaptureCount())
	}
	pat := re.Pattern()
	if pat.String() != `\d+` {
		t.Fatalf("unexpected pattern: %s", pat.String())
	}
}

func TestQRegularExpressionDefaultConstructor(t *testing.T) {
	re := NewQRegularExpression()
	if re.Pattern().String() != "" {
		t.Fatal("expected empty pattern")
	}
	// Qt6: an empty pattern is a valid regular expression (it matches
	// the empty string at every position). This mirrors the C++ behavior
	// of QRegularExpression() and QRegularExpression(QString()).
	if !re.IsValid() {
		t.Fatalf("expected empty pattern to be valid, error: %s", re.ErrorString().String())
	}
	subject := QStringFromString("")
	m := re.MatchString(subject, 0, NormalMatch, NoMatchOption)
	if !m.HasMatch() {
		t.Fatal("expected empty pattern to match empty subject")
	}
}

func TestQRegularExpressionClone(t *testing.T) {
	re := NewQRegularExpressionFromString(`abc`, NoPatternOption)
	cl := re.Clone()
	if !cl.IsValid() {
		t.Fatal("expected clone valid")
	}
	if cl.Pattern().String() != "abc" {
		t.Fatalf("unexpected clone pattern: %s", cl.Pattern().String())
	}
}

func TestQRegularExpressionSetPattern(t *testing.T) {
	re := NewQRegularExpression()
	p := QStringFromString(`\w+`)
	defer p.Close()
	re.SetPattern(p)
	if !re.IsValid() {
		t.Fatal("expected valid after SetPattern")
	}
	if re.Pattern().String() != `\w+` {
		t.Fatalf("unexpected pattern: %s", re.Pattern().String())
	}
}

func TestQRegularExpressionPatternOptions(t *testing.T) {
	re := NewQRegularExpressionFromString(`abc`, CaseInsensitiveOption|MultilineOption)
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
	esc := QRegularExpressionEscape(s)
	if esc.String() != `\\d\+` {
		t.Fatalf("unexpected escape: %s", esc.String())
	}
}

func TestQRegularExpressionWildcard(t *testing.T) {
	s := QStringFromString(`*.txt`)
	pattern := QRegularExpressionWildcardToRegularExpression(s, DefaultWildcardConversion)
	if pattern.IsEmpty() {
		t.Fatal("wildcard conversion returned empty pattern")
	}
	re := NewQRegularExpressionFromPattern(pattern, NoPatternOption)
	if !re.IsValid() {
		t.Fatalf("wildcard regex invalid: %s", re.ErrorString().String())
	}
}

func TestQRegularExpressionAnchoredPattern(t *testing.T) {
	s := QStringFromString(`\d+`)
	ap := QRegularExpressionAnchoredPattern(s)
	if ap.String() == "" {
		t.Fatal("expected non-empty anchored pattern")
	}
}

func TestQRegularExpressionFromWildcard(t *testing.T) {
	s := QStringFromString(`*.txt`)
	re := QRegularExpressionFromWildcard(s, true, DefaultWildcardConversion)
	if !re.IsValid() {
		t.Fatal("fromWildcard regex invalid")
	}
}