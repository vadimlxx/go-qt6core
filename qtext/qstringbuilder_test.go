package qtext

import (
	"bytes"
	"testing"
	"unicode/utf16"
)

func TestQStringBuilderEmpty(t *testing.T) {
	b := Concat()
	if b.Len() != 0 {
		t.Fatalf("Len: %d", b.Len())
	}
	if b.IsQStringFlavored() {
		t.Fatal("empty builder must not be QString-flavored")
	}
	qs := b.ToQString()
	defer qs.Close()
	if qs.GoString() != "" {
		t.Fatalf("ToQString empty: %q", qs.GoString())
	}
	ba := b.ToQByteArray()
	defer ba.Close()
	if ba.Size() != 0 {
		t.Fatalf("ToQByteArray empty: %d", ba.Size())
	}
}

func TestNewQStringBuilder(t *testing.T) {
	qs := QStringFromString("world")
	defer qs.Close()

	if qs.GoString() != "world" {
		t.Fatalf("ToQString: %q", qs.GoString())
	}
}

func TestQStringBuilderConcatStrings(t *testing.T) {
	b := Concat("hello", ", ", "world")
	if b.IsQStringFlavored() {
		t.Fatal("string-only builder must not be QString-flavored")
	}
	ba := b.ToQByteArray()
	defer ba.Close()
	if ba.String() != "hello, world" {
		t.Fatalf("ToQByteArray: %q", ba.String())
	}
	qs := b.ToQString()
	defer qs.Close()
	if qs.GoString() != "hello, world" {
		t.Fatalf("ToQString: %q", qs.GoString())
	}
}

func TestQStringBuilderConcatMixedFlavors(t *testing.T) {
	// char is QByteArray-flavored, QString is QString-flavored → result QString.
	qs := QStringFromString("world")
	defer qs.Close()
	b := Concat(byte('A'), "hello, ", qs)
	if !b.IsQStringFlavored() {
		t.Fatal("byte + QString must be QString-flavored")
	}
	out := b.ToQString()
	defer out.Close()
	if out.GoString() != "Ahello, world" {
		t.Fatalf("ToQString: %q", out.GoString())
	}
}

func TestQStringBuilderConcatQChar(t *testing.T) {
	// QChar is QString-flavored.
	b := Concat(NewQChar('H'), "i")
	if !b.IsQStringFlavored() {
		t.Fatal("QChar + string must be QString-flavored")
	}
	out := b.ToQString()
	defer out.Close()
	if out.GoString() != "Hi" {
		t.Fatalf("ToQString: %q", out.GoString())
	}
}

func TestQStringBuilderConcatQLatin1Char(t *testing.T) {
	b := Concat(NewQLatin1Char(0xE9), "tude") // "étude"
	if !b.IsQStringFlavored() {
		t.Fatal("QLatin1Char must be QString-flavored")
	}
	out := b.ToQString()
	defer out.Close()
	if out.GoString() != "étude" {
		t.Fatalf("ToQString: %q", out.GoString())
	}
}

func TestQStringBuilderConcatQStringView(t *testing.T) {
	sv := QStringViewFromUtf16(utf16.Encode([]rune("hello")))
	b := Concat(sv, ", ", "world")
	if !b.IsQStringFlavored() {
		t.Fatal("QStringView must be QString-flavored")
	}
	out := b.ToQString()
	defer out.Close()
	if out.GoString() != "hello, world" {
		t.Fatalf("ToQString: %q", out.GoString())
	}
}

func TestQStringBuilderConcatQLatin1StringView(t *testing.T) {
	l1 := QLatin1StringViewFromBytes([]byte{0xE9}) // "é"
	b := Concat("caf", l1)
	if !b.IsQStringFlavored() {
		t.Fatal("QLatin1StringView must be QString-flavored")
	}
	out := b.ToQString()
	defer out.Close()
	if out.GoString() != "café" {
		t.Fatalf("ToQString: %q", out.GoString())
	}
}

func TestQStringBuilderConcatQByteArray(t *testing.T) {
	ba := NewQByteArray("world")
	defer ba.Close()
	b := Concat("hello, ", ba)
	if b.IsQStringFlavored() {
		t.Fatal("string + QByteArray must NOT be QString-flavored")
	}
	out := b.ToQByteArray()
	defer out.Close()
	if out.String() != "hello, world" {
		t.Fatalf("ToQByteArray: %q", out.String())
	}
}

func TestQStringBuilderConcatQByteArrayView(t *testing.T) {
	bv := NewQByteArrayView("world")
	b := Concat("hello, ", bv)
	if b.IsQStringFlavored() {
		t.Fatal("string + QByteArrayView must NOT be QString-flavored")
	}
	out := b.ToQByteArray()
	defer out.Close()
	if out.String() != "hello, world" {
		t.Fatalf("ToQByteArray: %q", out.String())
	}
}

func TestQStringBuilderConcatUtf16(t *testing.T) {
	units := utf16.Encode([]rune("héllo"))
	b := Concat(units, "!")
	if !b.IsQStringFlavored() {
		t.Fatal("[]uint16 must be QString-flavored")
	}
	out := b.ToQString()
	defer out.Close()
	if out.GoString() != "héllo!" {
		t.Fatalf("ToQString: %q", out.GoString())
	}
}

func TestQStringBuilderNestedBuilders(t *testing.T) {
	inner := Concat("a", "b")
	outer := Concat(inner, "c")
	if outer.Len() != 3 {
		t.Fatalf("Len: %d", outer.Len())
	}
	out := outer.ToQByteArray()
	defer out.Close()
	if out.String() != "abc" {
		t.Fatalf("ToQByteArray: %q", out.String())
	}
}

func TestQStringBuilderAppend(t *testing.T) {
	b := Concat("a")
	b.Append(Concat("b", "c"))
	out := b.ToQByteArray()
	defer out.Close()
	if out.String() != "abc" {
		t.Fatalf("ToQByteArray: %q", out.String())
	}
}

func TestQStringBuilderAppendPart(t *testing.T) {
	b := FromString("a").AppendPart("b").AppendPart("c")
	out := b.ToQByteArray()
	defer out.Close()
	if out.String() != "abc" {
		t.Fatalf("ToQByteArray: %q", out.String())
	}
}

func TestQStringBuilderAppendToQString(t *testing.T) {
	q := QStringFromString("prefix: ")
	defer q.Close()
	b := Concat("a", "b", "c")
	b.AppendToQString(q)
	if q.GoString() != "prefix: abc" {
		t.Fatalf("AppendToQString: %q", q.GoString())
	}
}

func TestQStringBuilderAppendToQByteArray(t *testing.T) {
	ba := NewQByteArray("prefix: ")
	defer ba.Close()
	b := Concat("a", "b", "c")
	b.AppendToQByteArray(ba)
	if ba.String() != "prefix: abc" {
		t.Fatalf("AppendToQByteArray: %q", ba.String())
	}
}

func TestQStringBuilderToUtf8Latin1Local8Bit(t *testing.T) {
	b := Concat("h", NewQLatin1Char(0xE9), "llo") // "héllo" (QString-flavored)
	u := b.ToUtf8()
	defer u.Close()
	if !bytes.Equal(u.Bytes(), []byte("héllo")) {
		t.Fatalf("ToUtf8: %v", u.Bytes())
	}
	l := b.ToLatin1()
	defer l.Close()
	want := []byte{0x68, 0xE9, 0x6C, 0x6C, 0x6F}
	if !bytes.Equal(l.Bytes(), want) {
		t.Fatalf("ToLatin1: %v, want %v", l.Bytes(), want)
	}
}

func TestQStringBuilderUnsupportedType(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for unsupported type")
		}
	}()
	Concat(42) // int is not a supported operand
}
