package qgeometry

import "testing"

// helper: create a *QMargins and register cleanup.
func mkMargins(t *testing.T, l, top, r, b int) *QMargins {
	t.Helper()
	m := NewQMargins(l, top, r, b)
	t.Cleanup(m.Close)
	return m
}

func mkMarginsF(t *testing.T, l, top, r, b float64) *QMarginsF {
	t.Helper()
	m := NewQMarginsF(l, top, r, b)
	t.Cleanup(m.Close)
	return m
}

// ---- QMargins tests ------------------------------------------------------

func TestQMarginsConstruction(t *testing.T) {
	m := mkMargins(t, 1, 2, 3, 4)
	if m.Left() != 1 || m.Top() != 2 || m.Right() != 3 || m.Bottom() != 4 {
		t.Fatalf("fields: %d,%d,%d,%d", m.Left(), m.Top(), m.Right(), m.Bottom())
	}

	z := NewQMarginsZero()
	t.Cleanup(z.Close)
	if !z.IsNull() {
		t.Fatal("(0,0,0,0) must be null")
	}
	if m.IsNull() {
		t.Fatal("(1,2,3,4) must not be null")
	}
}

func TestQMarginsHandleSharing(t *testing.T) {
	m1 := mkMargins(t, 1, 2, 3, 4)
	m2 := m1 // shared pointer

	m2.SetLeft(100)
	if m1.Left() != 100 {
		t.Fatalf("shared handle: m1.Left = %d, want 100", m1.Left())
	}
}

func TestQMarginsCloneIndependence(t *testing.T) {
	a := mkMargins(t, 1, 2, 3, 4)
	b := a.Clone()
	t.Cleanup(b.Close)

	b.SetLeft(100)
	if a.Left() != 1 || b.Left() != 100 {
		t.Fatalf("clone independence: a.Left=%d b.Left=%d", a.Left(), b.Left())
	}
}

func TestQMarginsSetters(t *testing.T) {
	m := mkMargins(t, 0, 0, 0, 0)
	m.SetLeft(1)
	m.SetTop(2)
	m.SetRight(3)
	m.SetBottom(4)
	if m.Left() != 1 || m.Top() != 2 || m.Right() != 3 || m.Bottom() != 4 {
		t.Fatalf("after setters: %d,%d,%d,%d", m.Left(), m.Top(), m.Right(), m.Bottom())
	}
}

func TestQMarginsArithmetic(t *testing.T) {
	a := mkMargins(t, 1, 2, 3, 4)
	b := mkMargins(t, 10, 20, 30, 40)

	sum := a.Add(b)
	t.Cleanup(sum.Close)
	if sum.Left() != 11 || sum.Top() != 22 || sum.Right() != 33 || sum.Bottom() != 44 {
		t.Fatalf("Add: %d,%d,%d,%d", sum.Left(), sum.Top(), sum.Right(), sum.Bottom())
	}

	diff := b.Sub(a)
	t.Cleanup(diff.Close)
	if diff.Left() != 9 || diff.Top() != 18 || diff.Right() != 27 || diff.Bottom() != 36 {
		t.Fatalf("Sub: %d,%d,%d,%d", diff.Left(), diff.Top(), diff.Right(), diff.Bottom())
	}

	// inputs unchanged
	if a.Left() != 1 || b.Left() != 10 {
		t.Fatal("Add/Sub must not mutate inputs")
	}
}

func TestQMarginsScalarArithmetic(t *testing.T) {
	a := mkMargins(t, 1, 2, 3, 4)

	ai := a.AddInt(5)
	t.Cleanup(ai.Close)
	if ai.Left() != 6 || ai.Top() != 7 || ai.Right() != 8 || ai.Bottom() != 9 {
		t.Fatalf("AddInt: %d,%d,%d,%d", ai.Left(), ai.Top(), ai.Right(), ai.Bottom())
	}

	mi := a.MulInt(2)
	t.Cleanup(mi.Close)
	if mi.Left() != 2 || mi.Top() != 4 || mi.Right() != 6 || mi.Bottom() != 8 {
		t.Fatalf("MulInt: %d,%d,%d,%d", mi.Left(), mi.Top(), mi.Right(), mi.Bottom())
	}

	di := a.DivInt(2)
	t.Cleanup(di.Close)
	// 1/2 = 0 (int trunc), 2/2 = 1, 3/2 = 1, 4/2 = 2
	if di.Left() != 0 || di.Top() != 1 || di.Right() != 1 || di.Bottom() != 2 {
		t.Fatalf("DivInt: %d,%d,%d,%d", di.Left(), di.Top(), di.Right(), di.Bottom())
	}
}

func TestQMarginsNegUnaryPlus(t *testing.T) {
	a := mkMargins(t, 1, 2, 3, 4)

	n := a.Neg()
	t.Cleanup(n.Close)
	if n.Left() != -1 || n.Top() != -2 || n.Right() != -3 || n.Bottom() != -4 {
		t.Fatalf("Neg: %d,%d,%d,%d", n.Left(), n.Top(), n.Right(), n.Bottom())
	}

	u := a.UnaryPlus()
	t.Cleanup(u.Close)
	if u.Left() != 1 || u.Top() != 2 || u.Right() != 3 || u.Bottom() != 4 {
		t.Fatalf("UnaryPlus: %d,%d,%d,%d", u.Left(), u.Top(), u.Right(), u.Bottom())
	}
}

func TestQMarginsOr(t *testing.T) {
	a := mkMargins(t, 1, 5, 3, 7)
	b := mkMargins(t, 4, 2, 6, 1)

	r := a.Or(b)
	t.Cleanup(r.Close)
	if r.Left() != 4 || r.Top() != 5 || r.Right() != 6 || r.Bottom() != 7 {
		t.Fatalf("Or: %d,%d,%d,%d", r.Left(), r.Top(), r.Right(), r.Bottom())
	}
}

func TestQMarginsAssign(t *testing.T) {
	a := mkMargins(t, 1, 2, 3, 4)
	b := mkMargins(t, 10, 10, 10, 10)

	a.AddAssign(b)
	if a.Left() != 11 || a.Top() != 12 || a.Right() != 13 || a.Bottom() != 14 {
		t.Fatalf("AddAssign: %d,%d,%d,%d", a.Left(), a.Top(), a.Right(), a.Bottom())
	}

	a.MulIntAssign(2)
	if a.Left() != 22 || a.Top() != 24 || a.Right() != 26 || a.Bottom() != 28 {
		t.Fatalf("MulIntAssign: %d,%d,%d,%d", a.Left(), a.Top(), a.Right(), a.Bottom())
	}
}

func TestQMarginsEqual(t *testing.T) {
	a := mkMargins(t, 1, 2, 3, 4)
	b := mkMargins(t, 1, 2, 3, 4)
	c := mkMargins(t, 1, 2, 3, 5)

	if !a.Equal(b) {
		t.Fatal("Equal")
	}
	if a.NotEqual(b) {
		t.Fatal("NotEqual false")
	}
	if a.Equal(c) {
		t.Fatal("differing must not be equal")
	}
}

func TestQMarginsToMarginsF(t *testing.T) {
	m := mkMargins(t, 1, 2, 3, 4)
	f := m.ToMarginsF()
	t.Cleanup(f.Close)

	if f.Left() != 1 || f.Top() != 2 || f.Right() != 3 || f.Bottom() != 4 {
		t.Fatalf("ToMarginsF: %v,%v,%v,%v", f.Left(), f.Top(), f.Right(), f.Bottom())
	}
}

// ---- QMarginsF tests -----------------------------------------------------

func TestQMarginsFConstruction(t *testing.T) {
	m := mkMarginsF(t, 1.5, 2.5, 3.5, 4.5)
	if m.Left() != 1.5 || m.Top() != 2.5 || m.Right() != 3.5 || m.Bottom() != 4.5 {
		t.Fatalf("fields: %v,%v,%v,%v", m.Left(), m.Top(), m.Right(), m.Bottom())
	}

	z := NewQMarginsFZero()
	t.Cleanup(z.Close)
	if !z.IsNull() {
		t.Fatal("(0,0,0,0) must be null")
	}
}

func TestQMarginsFFromMargins(t *testing.T) {
	m := mkMargins(t, 1, 2, 3, 4)
	f := QMarginsFFromMargins(m)
	t.Cleanup(f.Close)
	if f.Left() != 1 || f.Top() != 2 || f.Right() != 3 || f.Bottom() != 4 {
		t.Fatalf("QMarginsFFromMargins: %v,%v,%v,%v", f.Left(), f.Top(), f.Right(), f.Bottom())
	}
}

func TestQMarginsFHandleSharing(t *testing.T) {
	m1 := mkMarginsF(t, 1, 2, 3, 4)
	m2 := m1

	m2.SetLeft(100)
	if m1.Left() != 100 {
		t.Fatalf("shared handle: m1.Left = %v", m1.Left())
	}
}

func TestQMarginsFCloneIndependence(t *testing.T) {
	a := mkMarginsF(t, 1, 2, 3, 4)
	b := a.Clone()
	t.Cleanup(b.Close)

	b.SetLeft(100)
	if a.Left() != 1 || b.Left() != 100 {
		t.Fatalf("clone independence: a.Left=%v b.Left=%v", a.Left(), b.Left())
	}
}

func TestQMarginsFArithmetic(t *testing.T) {
	a := mkMarginsF(t, 1, 2, 3, 4)
	b := mkMarginsF(t, 0.5, 0.5, 0.5, 0.5)

	sum := a.Add(b)
	t.Cleanup(sum.Close)
	if sum.Left() != 1.5 || sum.Top() != 2.5 || sum.Right() != 3.5 || sum.Bottom() != 4.5 {
		t.Fatalf("Add: %v,%v,%v,%v", sum.Left(), sum.Top(), sum.Right(), sum.Bottom())
	}

	mf := a.MulFloat(2)
	t.Cleanup(mf.Close)
	if mf.Left() != 2 || mf.Top() != 4 || mf.Right() != 6 || mf.Bottom() != 8 {
		t.Fatalf("MulFloat: %v,%v,%v,%v", mf.Left(), mf.Top(), mf.Right(), mf.Bottom())
	}

	d := a.DivFloat(2)
	t.Cleanup(d.Close)
	if d.Left() != 0.5 || d.Top() != 1 || d.Right() != 1.5 || d.Bottom() != 2 {
		t.Fatalf("DivFloat: %v,%v,%v,%v", d.Left(), d.Top(), d.Right(), d.Bottom())
	}

	n := a.Neg()
	t.Cleanup(n.Close)
	if n.Left() != -1 || n.Top() != -2 || n.Right() != -3 || n.Bottom() != -4 {
		t.Fatalf("Neg: %v,%v,%v,%v", n.Left(), n.Top(), n.Right(), n.Bottom())
	}
}

func TestQMarginsFEqual(t *testing.T) {
	a := mkMarginsF(t, 1, 2, 3, 4)
	b := mkMarginsF(t, 1, 2, 3, 4)
	if !a.Equal(b) {
		t.Fatal("Equal")
	}
	if a.NotEqual(b) {
		t.Fatal("NotEqual false")
	}
}

func TestQMarginsFEqualQMargins(t *testing.T) {
	f := mkMarginsF(t, 1, 2, 3, 4)
	m := mkMargins(t, 1, 2, 3, 4)
	if !f.EqualQMargins(m) {
		t.Fatal("EqualQMargins")
	}
}

func TestQMarginsFToMargins(t *testing.T) {
	cases := []struct {
		l, tp, r, b       float64
		wl, wtp, wr, wb   int
	}{
		{1.4, 2.6, 3.4, 4.6, 1, 3, 3, 5},
		{-1.4, -2.6, -3.4, -4.6, -1, -3, -3, -5},
		{1, 2, 3, 4, 1, 2, 3, 4},
	}
	for _, c := range cases {
		f := mkMarginsF(t, c.l, c.tp, c.r, c.b)
		m := f.ToMargins()
		t.Cleanup(m.Close)
		if m.Left() != c.wl || m.Top() != c.wtp || m.Right() != c.wr || m.Bottom() != c.wb {
			t.Errorf("ToMargins(%v,%v,%v,%v) = %d,%d,%d,%d, want %d,%d,%d,%d",
				c.l, c.tp, c.r, c.b, m.Left(), m.Top(), m.Right(), m.Bottom(),
				c.wl, c.wtp, c.wr, c.wb)
		}
	}
}

func TestQMarginsFAssign(t *testing.T) {
	a := mkMarginsF(t, 1, 2, 3, 4)
	b := mkMarginsF(t, 0.5, 0.5, 0.5, 0.5)

	a.AddAssign(b)
	if a.Left() != 1.5 || a.Top() != 2.5 || a.Right() != 3.5 || a.Bottom() != 4.5 {
		t.Fatalf("AddAssign: %v,%v,%v,%v", a.Left(), a.Top(), a.Right(), a.Bottom())
	}

	a.MulFloatAssign(2)
	if a.Left() != 3 || a.Top() != 5 || a.Right() != 7 || a.Bottom() != 9 {
		t.Fatalf("MulFloatAssign: %v,%v,%v,%v", a.Left(), a.Top(), a.Right(), a.Bottom())
	}
}
