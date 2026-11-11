package qgeometry

import "testing"

func TestQMarginsConstruction(t *testing.T) {
	m := NewQMargins(1, 2, 3, 4)
	if m.Left() != 1 || m.Top() != 2 || m.Right() != 3 || m.Bottom() != 4 {
		t.Fatalf("fields: %+v", m)
	}
	if NewQMarginsZero().IsNull() != true {
		t.Fatal("zero should be null")
	}
	if m.IsNull() {
		t.Fatal("non-zero should not be null")
	}
}

func TestQMarginsArithmetic(t *testing.T) {
	a := NewQMargins(1, 2, 3, 4)
	b := NewQMargins(10, 20, 30, 40)

	if got := a.Add(b); got != NewQMargins(11, 22, 33, 44) {
		t.Fatalf("Add: %+v", got)
	}
	if got := b.Sub(a); got != NewQMargins(9, 18, 27, 36) {
		t.Fatalf("Sub: %+v", got)
	}
	if got := a.AddInt(5); got != NewQMargins(6, 7, 8, 9) {
		t.Fatalf("AddInt: %+v", got)
	}
	if got := a.MulInt(2); got != NewQMargins(2, 4, 6, 8) {
		t.Fatalf("MulInt: %+v", got)
	}
	if got := a.DivInt(2); got != NewQMargins(0, 1, 1, 2) {
		t.Fatalf("DivInt: %+v", got)
	}
	if got := a.Neg(); got != NewQMargins(-1, -2, -3, -4) {
		t.Fatalf("Neg: %+v", got)
	}
	if got := a.Or(b); got != NewQMargins(10, 20, 30, 40) {
		t.Fatalf("Or: %+v", got)
	}
}

func TestQMarginsAssign(t *testing.T) {
	a := NewQMargins(1, 2, 3, 4)
	a.AddAssign(NewQMargins(10, 10, 10, 10))
	if a != NewQMargins(11, 12, 13, 14) {
		t.Fatalf("AddAssign: %+v", a)
	}
	a.MulIntAssign(2)
	if a != NewQMargins(22, 24, 26, 28) {
		t.Fatalf("MulIntAssign: %+v", a)
	}
}

func TestQMarginsToMarginsF(t *testing.T) {
	m := NewQMargins(1, 2, 3, 4)
	f := m.ToMarginsF()
	if f != NewQMarginsF(1, 2, 3, 4) {
		t.Fatalf("ToMarginsF: %+v", f)
	}
}

func TestQMarginsFConstruction(t *testing.T) {
	f := NewQMarginsF(1.5, 2.5, 3.5, 4.5)
	if f.Left() != 1.5 || f.Top() != 2.5 || f.Right() != 3.5 || f.Bottom() != 4.5 {
		t.Fatalf("fields: %+v", f)
	}
	if !NewQMarginsFZero().IsNull() {
		t.Fatal("zero should be null")
	}
}

func TestQMarginsFArithmetic(t *testing.T) {
	a := NewQMarginsF(1, 2, 3, 4)
	b := NewQMarginsF(0.5, 0.5, 0.5, 0.5)

	if got := a.Add(b); got != NewQMarginsF(1.5, 2.5, 3.5, 4.5) {
		t.Fatalf("Add: %+v", got)
	}
	if got := a.MulFloat(2); got != NewQMarginsF(2, 4, 6, 8) {
		t.Fatalf("MulFloat: %+v", got)
	}
	if got := a.DivFloat(2); got != NewQMarginsF(0.5, 1, 1.5, 2) {
		t.Fatalf("DivFloat: %+v", got)
	}
	if got := a.Neg(); got != NewQMarginsF(-1, -2, -3, -4) {
		t.Fatalf("Neg: %+v", got)
	}
}

func TestQMarginsFToMargins(t *testing.T) {
	f := NewQMarginsF(1.4, 2.6, 3.4, 4.6)
	if got := f.ToMargins(); got != NewQMargins(1, 3, 3, 5) {
		t.Fatalf("ToMargins: %+v", got)
	}
}

func TestMarginsEquality(t *testing.T) {
	a := NewQMargins(1, 2, 3, 4)
	if !a.Equal(NewQMargins(1, 2, 3, 4)) {
		t.Fatal("Equal")
	}
	if a.NotEqual(NewQMargins(1, 2, 3, 5)) == false {
		t.Fatal("NotEqual")
	}
}

func TestMarginsFEquality(t *testing.T) {
	a := NewQMarginsF(1, 2, 3, 4)
	if !a.Equal(NewQMarginsF(1, 2, 3, 4)) {
		t.Fatal("Equal (fuzzy)")
	}
	if !a.EqualQMargins(NewQMargins(1, 2, 3, 4)) {
		t.Fatal("EqualQMargins")
	}
}