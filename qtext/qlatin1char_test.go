package qtext

import "testing"

func TestNewRoundTrip(t *testing.T) {
	for b := 0; b < 256; b++ {
		c := NewQLatin1Char(byte(b))
		if got := c.ToLatin1(); got != byte(b) {
			t.Fatalf("ToLatin1(%d): got %d", b, got)
		}
		if got := c.Unicode(); got != uint16(b) {
			t.Fatalf("Unicode(%d): got %#x", b, got)
		}
	}
}

func TestNewQLatin1CharFromRune(t *testing.T) {
	cases := []struct {
		r    rune
		ok   bool
		want QLatin1Char
	}{
		{'A', true, NewQLatin1Char('A')},
		{'é', true, NewQLatin1Char(0xE9)},
		{0x00, true, 0},
		{0xFF, true, NewQLatin1Char(0xFF)},
		{0x100, false, 0},
		{'世', false, 0},
		{-1, false, 0},
	}
	for _, tc := range cases {
		c, ok := NewQLatin1CharFromRune(tc.r)
		if ok != tc.ok || (tc.ok && c != tc.want) {
			t.Fatalf("NewQLatin1CharFromRune(%U): got (%v,%v), want (%v,%v)",
				tc.r, c, ok, tc.want, tc.ok)
		}
	}
}

func TestToQChar(t *testing.T) {
	for b := 0; b < 256; b++ {
		q := NewQLatin1Char(byte(b)).ToQChar()
		if q.Unicode() != uint16(b) {
			t.Fatalf("NewQLatin1Char(%d).ToQChar().Unicode() = %#x", b, q.Unicode())
		}
	}
}

func TestNewQLatin1CharFromQChar(t *testing.T) {
	for b := 0; b < 256; b++ {
		c, ok := NewQLatin1CharFromQChar(NewQChar(uint16(b)))
		if !ok || byte(c) != byte(b) {
			t.Fatalf("NewQLatin1CharFromQChar(%d): got (%d,%v)", b, c, ok)
		}
	}
	if _, ok := NewQLatin1CharFromQChar(NewQChar(0x100)); ok {
		t.Fatal("NewQLatin1CharFromQChar(0x100) must be !ok")
	}
	if _, ok := NewQLatin1CharFromQChar(NewQChar(0x4E16)); ok {
		t.Fatal("NewQLatin1CharFromQChar(世) must be !ok")
	}
}

func TestComparisons(t *testing.T) {
	a := NewQLatin1Char('a')
	b := NewQLatin1Char('b')
	z := NewQLatin1Char(0x00)
	w := NewQLatin1Char(0xFF)

	if !a.Equal(a) { t.Fatal("a.Equal(a) must be true") }
	if a.Equal(b)  { t.Fatal("a.Equal(b) must be false") }
	if a.Compare(a) != 0 { t.Fatal("a.Compare(a) must be 0") }
	if a.Compare(b) != -1 { t.Fatal("a.Compare(b) must be -1") }
	if b.Compare(a) != 1 { t.Fatal("b.Compare(a) must be 1") }
	if z.Compare(w) != -1 || w.Compare(z) != 1 {
		t.Fatal("0x00 vs 0xFF ordering wrong")
	}
}

func TestByteComparisons(t *testing.T) {
	a := NewQLatin1Char('a')
	if !a.EqualByte('a')     { t.Fatal("EqualByte('a')") }
	if a.EqualByte('b')      { t.Fatal("EqualByte('b') must be false") }
	if a.CompareByte('a') != 0  { t.Fatal("CompareByte('a')") }
	if a.CompareByte('b') != -1 { t.Fatal("CompareByte('b')") }
	if a.CompareByte('A') != 1  { t.Fatal("CompareByte('A')") }
}

func TestString(t *testing.T) {
	if got := NewQLatin1Char('A').String(); got != "A" {
		t.Fatalf("String('A') = %q", got)
	}
	if got := NewQLatin1Char(0xE9).String(); got != "é" {
		t.Fatalf("String(0xE9) = %q", got)
	}
	if got := NewQLatin1Char(0xFF).String(); got != "ÿ" {
		t.Fatalf("String(0xFF) = %q", got)
	}
}

func TestZeroValue(t *testing.T) {
	var c QLatin1Char
	if byte(c) != 0 || c.Unicode() != 0 {
		t.Fatalf("zero value: byte=%d unicode=%d", byte(c), c.Unicode())
	}
}
