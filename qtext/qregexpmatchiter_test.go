package qtext

import "testing"

func TestQRegularExpressionGlobalMatch(t *testing.T) {
	re := NewQRegularExpressionFromString(`\d+`, NoPatternOption)
	subject := QStringFromString("a1 b22 c333")
	it := re.GlobalMatchString(subject, 0, NormalMatch, NoMatchOption)
	if !it.IsValid() {
		t.Fatal("invalid iterator")
	}
	count := 0
	for it.HasNext() {
		m := it.Next()
		if !m.HasMatch() {
			t.Fatal("expected match")
		}
		count++
	}
	if count != 3 {
		t.Fatalf("expected 3 matches, got %d", count)
	}
}

func TestQRegularExpressionGlobalMatchView(t *testing.T) {
	re := NewQRegularExpressionFromString(`\d+`, NoPatternOption)
	subject := QStringViewFromString("a1 b22")
	it := re.GlobalMatchView(subject, 0, NormalMatch, NoMatchOption)
	count := 0
	for it.HasNext() {
		m := it.Next()
		if m.CapturedInt(0).String() == "" {
			t.Fatal("empty capture")
		}
		count++
	}
	if count != 2 {
		t.Fatalf("expected 2 matches, got %d", count)
	}
}

func TestQRegularExpressionPeekNext(t *testing.T) {
	re := NewQRegularExpressionFromString(`\d+`, NoPatternOption)
	subject := QStringFromString("a1 b2")
	it := re.GlobalMatchString(subject, 0, NormalMatch, NoMatchOption)
	first := it.PeekNext()
	if !first.HasMatch() {
		t.Fatal("expected peek to yield a match")
	}
	if first.CapturedInt(0).String() != "1" {
		t.Fatalf("unexpected peek capture: %s", first.CapturedInt(0).String())
	}
	real := it.Next()
	if real.CapturedInt(0).String() != "1" {
		t.Fatalf("unexpected first capture after peek: %s", real.CapturedInt(0).String())
	}
	second := it.PeekNext()
	if second.CapturedInt(0).String() != "2" {
		t.Fatalf("unexpected second capture: %s", second.CapturedInt(0).String())
	}
}

func TestQRegularExpressionIteratorMeta(t *testing.T) {
	re := NewQRegularExpressionFromString(`\d+`, NoPatternOption)
	subject := QStringFromString("a1 b2")
	it := re.GlobalMatchString(subject, 0, NormalMatch, NoMatchOption)
	if it.MatchType() != NormalMatch {
		t.Fatalf("unexpected match type: %d", it.MatchType())
	}
	if it.MatchOptions() != NoMatchOption {
		t.Fatalf("unexpected match options: %d", it.MatchOptions())
	}
	if it.RegularExpression().Pattern().String() != `\d+` {
		t.Fatalf("unexpected iterator regex: %s", it.RegularExpression().Pattern().String())
	}
}