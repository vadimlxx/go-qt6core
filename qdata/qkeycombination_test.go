package qdata

import "testing"

func TestQKeyCombinationZeroValue(t *testing.T) {
	var q QKeyCombination
	// Go zero value has raw combination == 0, which is NOT the same as
	// Qt's default-constructed QKeyCombination() (that one holds
	// int(Qt::Key_unknown) == 0x01ffffff). Use NewQKeyCombination() for the
	// C++ default behavior.
	if q.ToCombined() != 0 {
		t.Fatalf("zero value ToCombined() = %d, want 0", q.ToCombined())
	}
	if q.Key() != Key(0) {
		t.Fatalf("zero value Key() = %#x, want 0", int(q.Key()))
	}
	if q.KeyboardModifiers() != NoModifier {
		t.Fatalf("zero value KeyboardModifiers() = %#x, want NoModifier",
			int(q.KeyboardModifiers()))
	}
}

func TestQKeyCombinationDefaultConstructor(t *testing.T) {
	q := NewQKeyCombination()
	// Qt's QKeyCombination() default constructor stores int(Qt::Key_unknown).
	if q.ToCombined() != int(Key_unknown) {
		t.Fatalf("NewQKeyCombination().ToCombined() = %#x, want %#x",
			q.ToCombined(), int(Key_unknown))
	}
	if q.Key() != Key_unknown {
		t.Fatalf("NewQKeyCombination().Key() = %#x, want Key_unknown", int(q.Key()))
	}
	if q.KeyboardModifiers() != NoModifier {
		t.Fatalf("NewQKeyCombination().KeyboardModifiers() = %#x, want NoModifier",
			int(q.KeyboardModifiers()))
	}
}

func TestQKeyCombinationFromKey(t *testing.T) {
	q := NewQKeyCombinationFromKey(Key_A)
	if q.Key() != Key_A {
		t.Fatalf("Key() = %#x, want Key_A", int(q.Key()))
	}
	if q.KeyboardModifiers() != NoModifier {
		t.Fatalf("KeyboardModifiers() = %#x, want NoModifier", int(q.KeyboardModifiers()))
	}
	if q.ToCombined() != int(Key_A) {
		t.Fatalf("ToCombined() = %#x, want %#x", q.ToCombined(), int(Key_A))
	}
}

func TestQKeyCombinationFromModifiersKey(t *testing.T) {
	cases := []struct {
		mods KeyboardModifiers
		key  Key
	}{
		{ShiftModifier, Key_A},
		{ControlModifier, Key_C},
		{AltModifier, Key_F4},
		{MetaModifier | ShiftModifier, Key_Q},
	}
	for _, c := range cases {
		q := NewQKeyCombinationFromModifiersKey(c.mods, c.key)
		if q.Key() != c.key {
			t.Errorf("Key() = %#x, want %#x", int(q.Key()), int(c.key))
		}
		if q.KeyboardModifiers() != c.mods {
			t.Errorf("KeyboardModifiers() = %#x, want %#x",
				int(q.KeyboardModifiers()), int(c.mods))
		}
	}
}

func TestQKeyCombinationCombineAlias(t *testing.T) {
	q1 := Combine(ControlModifier, Key_S)
	q2 := NewQKeyCombinationFromModifiersKey(ControlModifier, Key_S)
	if !q1.Equal(q2) {
		t.Fatal("Combine and NewQKeyCombinationFromModifiersKey must produce equal values")
	}
}

func TestQKeyCombinationFromCombined(t *testing.T) {
	raw := int(ControlModifier) | int(Key_S)
	q := QKeyCombinationFromCombined(raw)
	if q.ToCombined() != raw {
		t.Fatalf("ToCombined() = %#x, want %#x", q.ToCombined(), raw)
	}
	if q.Key() != Key_S {
		t.Fatalf("Key() = %#x, want Key_S", int(q.Key()))
	}
	if q.KeyboardModifiers() != ControlModifier {
		t.Fatalf("KeyboardModifiers() = %#x, want ControlModifier", int(q.KeyboardModifiers()))
	}
}

func TestQKeyCombinationValueSemantics(t *testing.T) {
	a := Combine(ShiftModifier, Key_Z)
	b := a // value copy
	if !a.Equal(b) {
		t.Fatal("value copy must be equal")
	}
	// Assignment of a different value to b must not affect a.
	b = NewQKeyCombinationFromKey(Key_X)
	if a.Key() != Key_Z {
		t.Fatalf("a.Key() changed after b reassignment: %#x", int(a.Key()))
	}
	if b.Key() != Key_X {
		t.Fatalf("b.Key() = %#x, want Key_X", int(b.Key()))
	}
}

func TestQKeyCombinationEquality(t *testing.T) {
	a := Combine(ControlModifier|ShiftModifier, Key_A)
	b := Combine(ControlModifier|ShiftModifier, Key_A)
	c := Combine(ControlModifier, Key_A)
	if !a.Equal(b) {
		t.Fatal("equal combinations must compare equal")
	}
	if a.Equal(c) {
		t.Fatal("combinations with different modifiers must not be equal")
	}
}

func TestQKeyCombinationKeyConstants(t *testing.T) {
	cases := []struct {
		key  Key
		want int
	}{
		{Key_A, 0x41},
		{Key_Z, 0x5a},
		{Key_0, 0x30},
		{Key_Space, 0x20},
		{Key_Escape, 0x01000000},
		{Key_F1, 0x01000030},
		{Key_F12, 0x0100003b},
		{Key_unknown, 0x01ffffff},
	}
	for _, c := range cases {
		if int(c.key) != c.want {
			t.Errorf("Key %#x = %#x, want %#x", c.want, int(c.key), c.want)
		}
	}
}

func TestQKeyCombinationModifierConstants(t *testing.T) {
	cases := []struct {
		mod  KeyboardModifier
		want int
	}{
		{NoModifier, 0x00000000},
		{ShiftModifier, 0x02000000},
		{ControlModifier, 0x04000000},
		{AltModifier, 0x08000000},
		{MetaModifier, 0x10000000},
		{KeypadModifier, 0x20000000},
	}
	for _, c := range cases {
		if int(c.mod) != c.want {
			t.Errorf("Modifier %#x = %#x, want %#x", c.want, int(c.mod), c.want)
		}
	}
}
