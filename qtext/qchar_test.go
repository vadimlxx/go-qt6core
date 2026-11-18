package qtext

import "testing"

func TestPredicatesAndCase(t *testing.T) {
	a := NewQChar('A')
	if !a.IsLetter() || !a.IsUpper() || a.IsLower() {
		t.Fatalf("'A': letter=%v upper=%v lower=%v", a.IsLetter(), a.IsUpper(), a.IsLower())
	}
	if got := a.ToLower(); got != NewQChar('a') {
		t.Fatalf("ToLower('A') = %#x, want %#x", uint16(got), uint16('a'))
	}
	if got := ToLowerOf('A'); got != 'a' {
		t.Fatalf("ToLowerOf('A') = %q", got)
	}
}

func TestCategoryAndDirection(t *testing.T) {
	if got := NewQChar('A').Category(); got != Letter_Uppercase {
		t.Fatalf("category('A') = %v", got)
	}
	if got := NewQChar('א').Direction(); got != DirR {
		t.Fatalf("direction('א') = %v", got)
	}
	if got := NewQChar('1').Category(); got != Number_DecimalDigit {
		t.Fatalf("category('1') = %v", got)
	}
}

func TestSpecialConstants(t *testing.T) {
	if !Null.IsNull()        { t.Fatal("Null should be null") }
	if !Space.IsSpace()      { t.Fatal("Space should be space") }
	if FormFeed != 0x000C    { t.Fatalf("FormFeed = %#x", uint16(FormFeed)) }
	if SoftHyphen != 0x00AD  { t.Fatalf("SoftHyphen = %#x", uint16(SoftHyphen)) }
	if VisualTabCharacter != 0x2192 { t.Fatalf("VisualTabCharacter = %#x", uint16(VisualTabCharacter)) }
}

func TestCombiningClass(t *testing.T) {
	if got := CombiningClassOf(0x0300); got != Combining_Above {
		t.Fatalf("CombiningClassOf(U+0300) = %d, want %d", got, Combining_Above)
	}
}

func TestScript(t *testing.T) {
	if got := ScriptOf('Ж'); got != Script_Cyrillic {
		t.Fatalf("ScriptOf('Ж') = %v", got)
	}
	if got := ScriptOf('א'); got != Script_Hebrew {
		t.Fatalf("ScriptOf('א') = %v", got)
	}
	if got := NewQChar('A').Script(); got != Script_Latin {
		t.Fatalf("Script('A') = %v", got)
	}
}

func TestUnicodeVersion(t *testing.T) {
	if got := CurrentUnicodeVersion(); got == Unicode_Unassigned {
		t.Fatalf("CurrentUnicodeVersion() = %v", got)
	}
	if got := UnicodeVersionOf('A'); got != Unicode_1_1 {
		t.Fatalf("UnicodeVersionOf('A') = %v", got)
	}
}

func TestDigitValue(t *testing.T) {
	if got := DigitValueOf('0'); got != 0 { t.Fatalf("DigitValueOf('0') = %d", got) }
	if got := DigitValueOf('7'); got != 7 { t.Fatalf("DigitValueOf('7') = %d", got) }
	if got := DigitValueOf('x'); got != -1 { t.Fatalf("DigitValueOf('x') = %d", got) }
	if got := DigitValueOf(' '); got != -1 { t.Fatalf("DigitValueOf(' ') = %d", got) }
	if got := DigitValueOf(0x216B); got != -1 { t.Fatalf("DigitValueOf(U+216B) = %d", got) }
	if got := DigitValueOf(0x0665); got != 5 { t.Fatalf("DigitValueOf(U+0665) = %d", got) }
	if got := DigitValueOf(0x096F); got != 9 { t.Fatalf("DigitValueOf(U+096F) = %d", got) }
}

func TestSurrogates(t *testing.T) {
	r := rune(0x1F600)
	if !RequiresSurrogates(r) { t.Fatal("must require surrogates") }
	hi, lo := HighSurrogate(r), LowSurrogate(r)
	if !hi.IsHighSurrogate() || !lo.IsLowSurrogate() {
		t.Fatalf("bad surrogate pair: %04X %04X", uint16(hi), uint16(lo))
	}
	if back := SurrogateToUcs4(hi, lo); back != r {
		t.Fatalf("round-trip: %U != %U", back, r)
	}
}

func TestNewQCharFromUcs4(t *testing.T) {
	if s := NewQCharFromUcs4('A'); len(s) != 1 || s[0] != NewQChar('A') {
		t.Fatalf("NewQCharFromUcs4('A') = %v", s)
	}
	s := NewQCharFromUcs4(0x1F600)
	if len(s) != 2 || !s[0].IsHighSurrogate() || !s[1].IsLowSurrogate() {
		t.Fatalf("NewQCharFromUcs4(0x1F600) = %v", s)
	}
	// New range checks:
	if got := NewQCharFromUcs4(-1); got != nil {
		t.Fatalf("NewQCharFromUcs4(-1) = %v, want nil", got)
	}
	if got := NewQCharFromUcs4(0x110000); got != nil {
		t.Fatalf("NewQCharFromUcs4(0x110000) = %v, want nil", got)
	}
	if got := NewQCharFromUcs4(LastValidCodePoint); len(got) != 2 {
		t.Fatalf("NewQCharFromUcs4(0x10FFFF) = %v", got)
	}
}

func TestNewQCharFromRuneStrict(t *testing.T) {
	if c, ok := NewQCharFromRuneStrict('A'); !ok || c != NewQChar('A') {
		t.Fatalf("NewQCharFromRuneStrict('A') = (%v,%v)", c, ok)
	}
	if _, ok := NewQCharFromRuneStrict(0x1F600); ok {
		t.Fatal("NewQCharFromRuneStrict(0x1F600) must be !ok")
	}
	if _, ok := NewQCharFromRuneStrict(-1); ok {
		t.Fatal("NewQCharFromRuneStrict(-1) must be !ok")
	}
}

func TestCellRow(t *testing.T) {
	c := NewQChar(0x1234)
	if c.Cell() != 0x34 || c.Row() != 0x12 {
		t.Fatalf("cell=%#x row=%#x", c.Cell(), c.Row())
	}
	c.SetCell(0xAB)
	if c != 0x12AB { t.Fatalf("after SetCell: %#x", uint16(c)) }
	c.SetRow(0xCD)
	if c != 0xCDAB { t.Fatalf("after SetRow: %#x", uint16(c)) }
}

func TestNewQCharFromLatin1(t *testing.T) {
	if got := NewQCharFromLatin1('é'); got != NewQChar(0x00E9) {
		t.Fatalf("NewQCharFromLatin1('é') = %#x", uint16(got))
	}
	if got := NewQChar(0x20AC).ToLatin1(); got != 0 {
		t.Fatalf("ToLatin1(€) = %#x, want 0", got)
	}
}

func TestDecomposeGoString(t *testing.T) {
	// U+00C0 (LATIN CAPITAL LETTER A WITH GRAVE) canonical-decomposes to "A\u0300".
	got := QCharDecomposeString(0x00C0)
	if got != "A\u0300" {
		t.Fatalf("QCharDecomposeString(U+00C0) = %q, want %q", got, "A\u0300")
	}
	// Invalid code point returns "".
	if got := QCharDecomposeString(-1); got != "" {
		t.Fatalf("QCharDecomposeString(-1) = %q, want \"\"", got)
	}
	if got := QCharDecomposeString(0x110000); got != "" {
		t.Fatalf("QCharDecomposeString(0x110000) = %q, want \"\"", got)
	}
}
