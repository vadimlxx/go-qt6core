package qtext

import (
	"bytes"
	"testing"
)

func TestQBitArrayNewEmpty(t *testing.T) {
	a := NewQBitArray()
	defer a.Close()

	if a.Size() != 0 {
		t.Fatalf("Size: %d, want 0", a.Size())
	}
	if !a.IsEmpty() {
		t.Fatal("IsEmpty must be true")
	}
	if !a.IsNull() {
		t.Fatal("IsNull must be true for default-constructed QBitArray")
	}
	if a.Bits() != nil {
		t.Fatalf("Bits: %v, want nil", a.Bits())
	}
}

func TestNewQBitArraySized(t *testing.T) {
	a := NewQBitArraySized(10, false)
	defer a.Close()

	if a.Size() != 10 {
		t.Fatalf("Size: %d, want 10", a.Size())
	}
	if a.IsEmpty() {
		t.Fatal("IsEmpty must be false")
	}
	if a.IsNull() {
		t.Fatal("IsNull must be false")
	}
	if a.CountSet(true) != 0 {
		t.Fatalf("CountSet(true): %d, want 0", a.CountSet(true))
	}

	b := NewQBitArraySized(7, true)
	defer b.Close()
	if b.CountSet(true) != 7 {
		t.Fatalf("CountSet(true): %d, want 7", b.CountSet(true))
	}
}

func TestQBitArraySetClearToggle(t *testing.T) {
	a := NewQBitArraySized(8, false)
	defer a.Close()

	if a.TestBit(3) {
		t.Fatal("bit 3 should be false initially")
	}
	a.SetBit(3)
	if !a.TestBit(3) {
		t.Fatal("bit 3 should be true after SetBit")
	}
	if old := a.ToggleBit(3); !old {
		t.Fatalf("ToggleBit returned %v, want true (old value)", old)
	}
	if a.TestBit(3) {
		t.Fatal("bit 3 should be false after ToggleBit")
	}
	a.SetBitVal(5, true)
	a.ClearBit(5)
	if a.TestBit(5) {
		t.Fatal("bit 5 should be false after ClearBit")
	}
}

func TestQBitArrayResize(t *testing.T) {
	a := NewQBitArraySized(4, true)
	defer a.Close()

	a.Resize(8)
	if a.Size() != 8 {
		t.Fatalf("Size: %d, want 8", a.Size())
	}
	if !a.TestBit(0) || a.TestBit(4) {
		t.Fatal("existing bits must be preserved, new bits must be false")
	}

	a.Resize(2)
	if a.Size() != 2 {
		t.Fatalf("Size: %d, want 2", a.Size())
	}
}

func TestQBitArrayFill(t *testing.T) {
	a := NewQBitArraySized(5, false)
	defer a.Close()

	a.Fill(true)
	if a.CountSet(true) != 5 {
		t.Fatalf("Fill(true): %d, want 5", a.CountSet(true))
	}

	a.FillSized(false, 8)
	if a.Size() != 8 || a.CountSet(true) != 0 {
		t.Fatalf("FillSized: size=%d count=%d", a.Size(), a.CountSet(true))
	}
}

func TestQBitArrayFillRange(t *testing.T) {
	a := NewQBitArraySized(10, false)
	defer a.Close()

	a.FillRange(true, 2, 7)
	if a.CountSet(true) != 5 {
		t.Fatalf("FillRange count: %d, want 5", a.CountSet(true))
	}
	for i := 0; i < 10; i++ {
		want := i >= 2 && i < 7
		if got := a.TestBit(i); got != want {
			t.Fatalf("bit %d: got %v, want %v", i, got, want)
		}
	}
}

func TestQBitArrayAndOrXorNot(t *testing.T) {
	a := NewQBitArraySized(8, false)
	defer a.Close()
	b := NewQBitArraySized(8, false)
	defer b.Close()

	a.SetBit(0)
	a.SetBit(1)
	b.SetBit(1)
	b.SetBit(2)

	and := And(a, b)
	defer and.Close()
	if and.Size() != 8 || and.CountSet(true) != 1 || !and.TestBit(1) {
		t.Fatalf("And: %v", and.Bits())
	}

	or := Or(a, b)
	defer or.Close()
	if or.CountSet(true) != 3 {
		t.Fatalf("Or count: %d, want 3", or.CountSet(true))
	}

	xor := Xor(a, b)
	defer xor.Close()
	if xor.CountSet(true) != 2 || !xor.TestBit(0) || !xor.TestBit(2) {
		t.Fatalf("Xor: %v", xor.Bits())
	}

	not := Not(a)
	defer not.Close()
	if not.CountSet(true) != 6 || not.TestBit(0) || not.TestBit(1) {
		t.Fatalf("Not: %v", not.Bits())
	}
}

func TestQBitArrayAndAssign(t *testing.T) {
	a := NewQBitArraySized(8, true)
	defer a.Close()
	b := NewQBitArraySized(8, false)
	defer b.Close()
	b.SetBit(3)

	a.AndAssign(b)
	if a.CountSet(true) != 1 || !a.TestBit(3) {
		t.Fatalf("AndAssign: %v", a.Bits())
	}
}

func TestQBitArrayBitsRoundTrip(t *testing.T) {
	a := NewQBitArraySized(12, false)
	defer a.Close()
	a.SetBit(0)
	a.SetBit(5)
	a.SetBit(11)

	raw := a.Bits()
	if len(raw) != 2 {
		t.Fatalf("Bits length: %d, want 2", len(raw))
	}

	b := NewQBitArrayFromBits(raw, 12)
	defer b.Close()
	if b.Size() != 12 {
		t.Fatalf("NewQBitArrayFromBits size: %d, want 12", b.Size())
	}
	if !b.Equal(a) {
		t.Fatal("NewQBitArrayFromBits round-trip mismatch")
	}
}

func TestQBitArrayBitsFormat(t *testing.T) {
	// Qt bit layout: bit i = (data[i/8] >> (i % 8)) & 1.
	a := NewQBitArraySized(8, false)
	defer a.Close()
	for i := 0; i < 8; i++ {
		a.SetBit(i)
	}
	raw := a.Bits()
	if !bytes.Equal(raw, []byte{0xFF}) {
		t.Fatalf("Bits: %v, want [0xFF]", raw)
	}
}

func TestQBitArrayToUInt32(t *testing.T) {
	a := NewQBitArraySized(32, false)
	defer a.Close()
	for i := 0; i < 8; i++ {
		a.SetBit(i)
	}

	le, ok := a.ToUInt32(true)
	if !ok {
		t.Fatal("ToUInt32 ok must be true")
	}
	if le != 0x000000FF {
		t.Fatalf("LE value: 0x%08X, want 0x000000FF", le)
	}

	be, ok := a.ToUInt32(false)
	if !ok {
		t.Fatal("ToUInt32 ok must be true")
	}
	if be != 0xFF000000 {
		t.Fatalf("BE value: 0x%08X, want 0xFF000000", be)
	}
}

func TestQBitArrayEqual(t *testing.T) {
	a := NewQBitArraySized(8, false)
	defer a.Close()
	a.SetBit(2)

	b := a.Clone()
	defer b.Close()
	if !a.Equal(b) {
		t.Fatal("equal arrays")
	}

	c := NewQBitArraySized(8, false)
	defer c.Close()
	if a.Equal(c) {
		t.Fatal("different arrays must not be equal")
	}
}

func TestQBitArraySwap(t *testing.T) {
	a := NewQBitArraySized(4, true)
	defer a.Close()
	b := NewQBitArraySized(4, false)
	defer b.Close()
	b.SetBit(0)

	a.Swap(b)
	if a.Size() != 4 || a.CountSet(true) != 1 || !a.TestBit(0) {
		t.Fatalf("after swap, a: %v", a.Bits())
	}
	if b.CountSet(true) != 4 {
		t.Fatalf("after swap, b: %v", b.Bits())
	}
}

func TestBitRef(t *testing.T) {
	a := NewQBitArraySized(8, false)
	defer a.Close()

	r := a.BitRef(3)
	if r.Bool() {
		t.Fatal("initial Bool must be false")
	}
	if !r.Not() {
		t.Fatal("initial Not must be true")
	}
	r.Set(true)
	if !a.TestBit(3) {
		t.Fatal("Set via BitRef must be visible in parent")
	}
	r.Assign(a.BitRef(5)) // bit 5 is false
	if a.TestBit(3) {
		t.Fatal("after Assign from bit 5 (false), bit 3 must be false")
	}

	a.SetBit(5)
	r.Assign(a.BitRef(5)) // bit 5 is now true
	if !a.TestBit(3) {
		t.Fatal("after Assign from bit 5 (true), bit 3 must be true")
	}
}

func TestCloneIsIndependent(t *testing.T) {
	a := NewQBitArraySized(8, true)
	defer a.Close()
	c := a.Clone()
	defer c.Close()

	c.ClearBit(0)
	if !a.TestBit(0) {
		t.Fatal("original must remain true after clone mutation")
	}
	if c.TestBit(0) {
		t.Fatal("clone bit 0 must be false")
	}
}

func TestPanicsOnOutOfRange(t *testing.T) {
	a := NewQBitArraySized(4, false)
	defer a.Close()

	funcs := []func(){
		func() { a.TestBit(-1) },
		func() { a.TestBit(4) },
		func() { a.SetBit(4) },
		func() { a.ClearBit(-1) },
		func() { a.ToggleBit(100) },
		func() { a.BitRef(4) },
	}
	for i, fn := range funcs {
		func() {
			defer func() {
				if r := recover(); r == nil {
					t.Fatalf("funcs[%d] must panic", i)
				}
			}()
			fn()
		}()
	}
}

func TestQBitArrayUseAfterClosePanics(t *testing.T) {
	a := NewQBitArraySized(4, false)
	a.Close()
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Size() after Close() must panic")
		}
	}()
	_ = a.Size()
}

func TestQBitArrayCloseNilSafe(t *testing.T) {
	var a *QBitArray
	a.Close() // must not panic
}

func TestQBitArrayCloseIdempotent(t *testing.T) {
	a := NewQBitArraySized(4, false)
	a.Close()
	a.Close() // second Close must not panic
}

func TestQBitArrayNilOperands(t *testing.T) {
	a := NewQBitArraySized(8, false)
	defer a.Close()
	a.SetBit(1)

	// And/Or/Xor with nil return a clone of a (null operand = zero-length).
	andNil := And(a, nil)
	defer andNil.Close()
	if !andNil.Equal(a) {
		t.Fatal("And(a, nil) must equal clone of a")
	}

	orNil := Or(a, nil)
	defer orNil.Close()
	if !orNil.Equal(a) {
		t.Fatal("Or(a, nil) must equal clone of a")
	}

	xorNil := Xor(a, nil)
	defer xorNil.Close()
	if !xorNil.Equal(a) {
		t.Fatal("Xor(a, nil) must equal clone of a")
	}

	// Assign/Swap with nil are no-ops.
	b := a.Clone()
	defer b.Close()
	b.Assign(nil)
	if !b.Equal(a) {
		t.Fatal("Assign(nil) must be a no-op")
	}
	b.Swap(nil)
	if !b.Equal(a) {
		t.Fatal("Swap(nil) must be a no-op")
	}

	// Equal(nil) is false.
	if a.Equal(nil) {
		t.Fatal("Equal(nil) must be false")
	}
}
