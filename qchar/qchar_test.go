package qchar

import "testing"

func TestPredicatesAndCase(t *testing.T) {
    a := New('A')
    if !a.IsLetter() || !a.IsUpper() || a.IsLower() {
        t.Fatalf("'A': letter=%v upper=%v lower=%v", a.IsLetter(), a.IsUpper(), a.IsLower())
    }
    // instance methods return QChar, not rune
    if got := a.ToLower(); got != New('a') {
        t.Fatalf("ToLower('A') = %#x, want %#x", uint16(got), uint16('a'))
    }
    // static helper returns rune
    if got := ToLowerOf('A'); got != 'a' {
        t.Fatalf("ToLowerOf('A') = %q", got)
    }
}

func TestCategoryAndDirection(t *testing.T) {
    if got := New('A').Category(); got != Letter_Uppercase {
        t.Fatalf("category('A') = %v", got)
    }
    if got := New('א').Direction(); got != DirR {
        t.Fatalf("direction('א') = %v", got)
    }
    // '1' is DecimalDigit
    if got := New('1').Category(); got != Number_DecimalDigit {
        t.Fatalf("category('1') = %v", got)
    }
}

func TestSpecialConstants(t *testing.T) {
    if !Null.IsNull()        { t.Fatal("Null should be null") }
    if !Space.IsSpace()      { t.Fatal("Space should be space") }
    if FormFeed != 0x000C    { t.Fatalf("FormFeed = %#x", uint16(FormFeed)) }
    if SoftHyphen != 0x00AD  { t.Fatalf("SoftHyphen = %#x", uint16(SoftHyphen)) }
    if VisualTabCharacter != 0x2192 {
        t.Fatalf("VisualTabCharacter = %#x", uint16(VisualTabCharacter))
    }
}

func TestCombiningClass(t *testing.T) {
    // U+0300 (COMBINING GRAVE ACCENT) has combining class 230 (Above).
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
    if got := New('A').Script(); got != Script_Latin {
        t.Fatalf("Script('A') = %v", got)
    }
}

func TestUnicodeVersion(t *testing.T) {
    if got := CurrentUnicodeVersion(); got == Unicode_Unassigned {
        t.Fatalf("CurrentUnicodeVersion() = %v", got)
    }
    // Latin 'A' has existed since Unicode 1.1
    if got := UnicodeVersionOf('A'); got != Unicode_1_1 {
        t.Fatalf("UnicodeVersionOf('A') = %v", got)
    }
}

// old
// func TestDigitValue(t *testing.T) {
//     if got := DigitValueOf('7'); got != 7 {
//         t.Fatalf("DigitValueOf('7') = %d", got)
//     }
//     // Roman numeral Ⅻ (U+216B) — 12
//     if got := DigitValueOf(0x216B); got != -1 {
//         t.Fatalf("DigitValueOf(U+216B) = %d", got)
//     }
//     if got := DigitValueOf('x'); got != -1 {
//         t.Fatalf("DigitValueOf('x') = %d, want -1", got)
//     }
// }

// new
func TestDigitValue(t *testing.T) {
	// ASCII digits: 0–9
	if got := DigitValueOf('0'); got != 0 {
		t.Fatalf("DigitValueOf('0') = %d, want 0", got)
	}
	if got := DigitValueOf('7'); got != 7 {
		t.Fatalf("DigitValueOf('7') = %d, want 7", got)
	}

	// Non-digits — must return -1
	if got := DigitValueOf('x'); got != -1 {
		t.Fatalf("DigitValueOf('x') = %d, want -1", got)
	}
	if got := DigitValueOf(' '); got != -1 {
		t.Fatalf("DigitValueOf(' ') = %d, want -1", got)
	}

	// Roman numeral Ⅻ (U+216B) is category Nl (Letter Number), NOT Nd.
	// QChar::digitValue() only handles Nd characters -> -1.
	if got := DigitValueOf(0x216B); got != -1 {
		t.Fatalf("DigitValueOf(U+216B) = %d, want -1 (category Nl, not Nd)", got)
	}

	// Non-ASCII decimal digit: ARABIC-INDIC DIGIT FIVE (U+0665), category Nd.
	if got := DigitValueOf(0x0665); got != 5 {
		t.Fatalf("DigitValueOf(U+0665) = %d, want 5", got)
	}
	// Devanagari digit nine (U+096F), category Nd.
	if got := DigitValueOf(0x096F); got != 9 {
		t.Fatalf("DigitValueOf(U+096F) = %d, want 9", got)
	}
}

func TestSurrogates(t *testing.T) {
    r := rune(0x1F600) // 😀
    if !RequiresSurrogates(r) { t.Fatal("must require surrogates") }
    hi, lo := HighSurrogate(r), LowSurrogate(r)
    if !hi.IsHighSurrogate() || !lo.IsLowSurrogate() {
        t.Fatalf("bad surrogate pair: %04X %04X", uint16(hi), uint16(lo))
    }
    if back := SurrogateToUcs4(hi, lo); back != r {
        t.Fatalf("round-trip: %U != %U", back, r)
    }
}

func TestFromUcs4(t *testing.T) {
    // BMP: single element
    if s := FromUcs4('A'); len(s) != 1 || s[0] != New('A') {
        t.Fatalf("FromUcs4('A') = %v", s)
    }
    // Supplementary: two elements
    s := FromUcs4(0x1F600)
    if len(s) != 2 || !s[0].IsHighSurrogate() || !s[1].IsLowSurrogate() {
        t.Fatalf("FromUcs4(0x1F600) = %v", s)
    }
}

func TestCellRow(t *testing.T) {
    c := New(0x1234)
    if c.Cell() != 0x34 || c.Row() != 0x12 {
        t.Fatalf("cell=%#x row=%#x", c.Cell(), c.Row())
    }
    c.SetCell(0xAB)
    if c != 0x12AB { t.Fatalf("after SetCell: %#x", uint16(c)) }
    c.SetRow(0xCD)
    if c != 0xCDAB { t.Fatalf("after SetRow: %#x", uint16(c)) }
}

func TestFromLatin1(t *testing.T) {
    if got := FromLatin1('é'); got != New(0x00E9) {
        t.Fatalf("FromLatin1('é') = %#x", uint16(got))
    }
    // Latin1 chars above 0xFF are not representable — ToLatin1 returns 0
    if got := New(0x20AC /* € */).ToLatin1(); got != 0 {
        t.Fatalf("ToLatin1(€) = %#x, want 0", got)
    }
}
