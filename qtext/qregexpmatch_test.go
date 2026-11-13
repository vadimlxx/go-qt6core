package qtext

import "testing"

func TestQRegularExpressionMatch(t *testing.T) {
	re := NewQRegularExpressionFromString(`(\d+)-(\d+)`, NoPatternOption)
	subject := QStringFromString("abc 123-456 def")
	m := re.MatchString(subject, 0, NormalMatch, NoMatchOption)
	if !m.HasMatch() {
		t.Fatal("expected match")
	}
	if m.LastCapturedIndex() != 2 {
		t.Fatalf("expected 2 captures, got %d", m.LastCapturedIndex())
	}
	c0 := m.CapturedInt(0)
	if c0.String() != "123-456" {
		t.Fatalf("unexpected capture 0: %s", c0.String())
	}
	c1 := m.CapturedInt(1)
	if c1.String() != "123" {
		t.Fatalf("unexpected capture 1: %s", c1.String())
	}
	if m.CapturedStartInt(0) != 4 {
		t.Fatalf("unexpected start: %d", m.CapturedStartInt(0))
	}
	if m.CapturedLengthInt(0) != 7 {
		t.Fatalf("unexpected length: %d", m.CapturedLengthInt(0))
	}
	if m.CapturedEndInt(0) != 11 {
		t.Fatalf("unexpected end: %d", m.CapturedEndInt(0))
	}
}

func TestQRegularExpressionMatchView(t *testing.T) {
	re := NewQRegularExpressionFromString(`(\d+)`, NoPatternOption)
	subject := QStringViewFromString("abc 42 def")
	m := re.MatchView(subject, 0, NormalMatch, NoMatchOption)
	if !m.HasMatch() {
		t.Fatal("expected match")
	}
	c := m.CapturedInt(1)
	if c.String() != "42" {
		t.Fatalf("unexpected capture: %s", c.String())
	}
	v := m.CapturedViewInt(1)
	if v.GoString() != "42" {
		t.Fatalf("unexpected captured view: %s", v.GoString())
	}
}

func TestQRegularExpressionMatchNamedGroup(t *testing.T) {
	re := NewQRegularExpressionFromString(`(?<num>\d+)`, NoPatternOption)
	subject := QStringFromString("abc 42 def")
	m := re.MatchString(subject, 0, NormalMatch, NoMatchOption)
	if !m.HasMatch() {
		t.Fatal("expected match")
	}
	if !m.HasCapturedName("num") {
		t.Fatal("expected named capture")
	}
	c := m.CapturedName("num")
	if c.String() != "42" {
		t.Fatalf("unexpected named capture: %s", c.String())
	}
	if m.CapturedStartName("num") != 4 {
		t.Fatalf("unexpected named start: %d", m.CapturedStartName("num"))
	}
	if m.CapturedLengthName("num") != 2 {
		t.Fatalf("unexpected named length: %d", m.CapturedLengthName("num"))
	}
	if m.CapturedEndName("num") != 6 {
		t.Fatalf("unexpected named end: %d", m.CapturedEndName("num"))
	}
}

func TestQRegularExpressionMatchNoMatch(t *testing.T) {
	re := NewQRegularExpressionFromString(`\d+`, NoPatternOption)
	subject := QStringFromString("abc def")
	m := re.MatchString(subject, 0, NormalMatch, NoMatchOption)
	if m.HasMatch() {
		t.Fatal("expected no match")
	}
	if m.HasCapturedInt(0) {
		t.Fatal("expected no capture")
	}
}