package qgeometry

import "testing"

// helper: create a QPoint and register cleanup.
func mkPoint(t *testing.T, x, y int) *QPoint {
	t.Helper()
	p := NewQPoint(x, y)
	t.Cleanup(p.Close)
	return p
}

func mkPointF(t *testing.T, x, y float64) *QPointF {
	t.Helper()
	p := NewQPointF(x, y)
	t.Cleanup(p.Close)
	return p
}

func TestQPointConstruction(t *testing.T) {
	p := mkPoint(t, 3, 4)
	if p.X() != 3 || p.Y() != 4 {
		t.Fatalf("fields: x=%d y=%d", p.X(), p.Y())
	}

	z := NewQPointZero()
	t.Cleanup(z.Close)
	if !z.IsNull() {
		t.Fatal("(0,0) must be null")
	}
	if p.IsNull() {
		t.Fatal("(3,4) must not be null")
	}
}

func TestQPointSetters(t *testing.T) {
	p := mkPoint(t, 1, 2)
	p.SetX(10)
	p.SetY(20)
	if p.X() != 10 || p.Y() != 20 {
		t.Fatalf("after setters: x=%d y=%d", p.X(), p.Y())
	}
}

func TestQPointHandleSharing(t *testing.T) {
	// Two Go variables referring to the same C++ QPoint see each other's mutations.
	p1 := mkPoint(t, 1, 2)
	p2 := p1 // shared pointer, not a copy

	p2.SetX(100)
	if p1.X() != 100 {
		t.Fatalf("shared handle: p1.X = %d, want 100", p1.X())
	}
}

func TestQPointCloneIndependence(t *testing.T) {
	a := mkPoint(t, 1, 2)
	b := a.Clone()
	t.Cleanup(b.Close)

	b.SetX(100)
	if a.X() != 1 || b.X() != 100 {
		t.Fatalf("clone independence: a.X=%d b.X=%d", a.X(), b.X())
	}
}

func TestQPointTransposed(t *testing.T) {
	p := mkPoint(t, 3, 4)
	tp := p.Transposed()
	t.Cleanup(tp.Close)

	if tp.X() != 4 || tp.Y() != 3 {
		t.Fatalf("Transposed: %d,%d", tp.X(), tp.Y())
	}
	// original unchanged
	if p.X() != 3 || p.Y() != 4 {
		t.Fatal("Transposed must not mutate")
	}
}

func TestQPointManhattanLength(t *testing.T) {
	cases := []struct {
		x, y int
		want int
	}{
		{0, 0, 0},
		{3, 4, 7},
		{-3, 4, 7},
		{-3, -4, 7},
	}
	for _, c := range cases {
		p := mkPoint(t, c.x, c.y)
		if got := p.ManhattanLength(); got != c.want {
			t.Errorf("(%d,%d).ManhattanLength() = %d, want %d", c.x, c.y, got, c.want)
		}
	}
}

func TestQPointArithmeticDoesNotMutateInputs(t *testing.T) {
	a := mkPoint(t, 3, 4)
	b := mkPoint(t, 1, 2)

	sum := a.Add(b)
	t.Cleanup(sum.Close)
	if sum.X() != 4 || sum.Y() != 6 {
		t.Fatalf("Add: %d,%d", sum.X(), sum.Y())
	}

	diff := a.Sub(b)
	t.Cleanup(diff.Close)
	if diff.X() != 2 || diff.Y() != 2 {
		t.Fatalf("Sub: %d,%d", diff.X(), diff.Y())
	}

	// inputs unchanged
	if a.X() != 3 || b.X() != 1 {
		t.Fatal("Add/Sub must not mutate inputs")
	}
}

func TestQPointMulDiv(t *testing.T) {
	a := mkPoint(t, 3, 4)

	m := a.MulInt(2)
	t.Cleanup(m.Close)
	if m.X() != 6 || m.Y() != 8 {
		t.Fatalf("MulInt: %d,%d", m.X(), m.Y())
	}

	mf := a.MulFloat(0.5)
	t.Cleanup(mf.Close)
	if mf.X() != 2 || mf.Y() != 2 {
		t.Fatalf("MulFloat: %d,%d", mf.X(), mf.Y())
	}

	d := a.DivFloat(2)
	t.Cleanup(d.Close)
	if d.X() != 2 || d.Y() != 2 {
		t.Fatalf("DivFloat: %d,%d", d.X(), d.Y())
	}

	// input unchanged
	if a.X() != 3 || a.Y() != 4 {
		t.Fatal("Mul/Div must not mutate input")
	}
}

func TestQPointAddAssignMutates(t *testing.T) {
	a := mkPoint(t, 3, 4)
	b := mkPoint(t, 1, 2)

	a.AddAssign(b)
	if a.X() != 4 || a.Y() != 6 {
		t.Fatalf("AddAssign: %d,%d", a.X(), a.Y())
	}
	if b.X() != 1 || b.Y() != 2 {
		t.Fatal("AddAssign must not mutate other")
	}
}

func TestQPointMulIntAssignMutates(t *testing.T) {
	a := mkPoint(t, 3, 4)
	a.MulIntAssign(2)
	if a.X() != 6 || a.Y() != 8 {
		t.Fatalf("MulIntAssign: %d,%d", a.X(), a.Y())
	}
}

func TestQPointNegUnaryPlus(t *testing.T) {
	a := mkPoint(t, 3, 4)

	n := a.Neg()
	t.Cleanup(n.Close)
	if n.X() != -3 || n.Y() != -4 {
		t.Fatalf("Neg: %d,%d", n.X(), n.Y())
	}

	u := a.UnaryPlus()
	t.Cleanup(u.Close)
	if u.X() != 3 || u.Y() != 4 {
		t.Fatalf("UnaryPlus: %d,%d", u.X(), u.Y())
	}

	// a unchanged
	if a.X() != 3 || a.Y() != 4 {
		t.Fatal("Neg/UnaryPlus must not mutate")
	}
}

func TestQPointDotProduct(t *testing.T) {
	a := mkPoint(t, 3, 4)
	b := mkPoint(t, 5, 6)
	// 3*5 + 4*6 = 39
	if got := a.DotProduct(b); got != 39 {
		t.Fatalf("DotProduct: %d", got)
	}
}

func TestQPointEqual(t *testing.T) {
	a := mkPoint(t, 3, 4)
	b := mkPoint(t, 3, 4)
	c := mkPoint(t, 4, 3)

	if !a.Equal(b) {
		t.Fatal("Equal(a, b)")
	}
	if a.NotEqual(b) {
		t.Fatal("NotEqual(a, b)")
	}
	if a.Equal(c) {
		t.Fatal("Equal(a, c) must be false")
	}
	if !a.NotEqual(c) {
		t.Fatal("NotEqual(a, c) must be true")
	}
}

func TestQPointToPointF(t *testing.T) {
	p := mkPoint(t, 3, 4)
	f := p.ToPointF()
	t.Cleanup(f.Close)

	if f.X() != 3 || f.Y() != 4 {
		t.Fatalf("ToPointF: %f,%f", f.X(), f.Y())
	}
}

// ---- QPointF tests -------------------------------------------------------

func TestQPointFConstruction(t *testing.T) {
	p := mkPointF(t, 1.5, 2.5)
	if p.X() != 1.5 || p.Y() != 2.5 {
		t.Fatalf("fields: %f,%f", p.X(), p.Y())
	}

	z := NewQPointFZero()
	t.Cleanup(z.Close)
	if !z.IsNull() {
		t.Fatal("(0,0) must be null")
	}
}

func TestQPointFHandleSharing(t *testing.T) {
	p1 := mkPointF(t, 1, 2)
	p2 := p1

	p2.SetX(100)
	if p1.X() != 100 {
		t.Fatalf("shared handle: p1.X = %f", p1.X())
	}
}

func TestQPointFCloneIndependence(t *testing.T) {
	a := mkPointF(t, 1, 2)
	b := a.Clone()
	t.Cleanup(b.Close)

	b.SetX(100)
	if a.X() != 1 || b.X() != 100 {
		t.Fatalf("clone independence: a.X=%f b.X=%f", a.X(), b.X())
	}
}

func TestQPointFArithmetic(t *testing.T) {
	a := mkPointF(t, 3, 4)
	b := mkPointF(t, 0.5, 1.5)

	sum := a.Add(b)
	t.Cleanup(sum.Close)
	if sum.X() != 3.5 || sum.Y() != 5.5 {
		t.Fatalf("Add: %f,%f", sum.X(), sum.Y())
	}

	mf := a.MulFloat(2)
	t.Cleanup(mf.Close)
	if mf.X() != 6 || mf.Y() != 8 {
		t.Fatalf("MulFloat: %f,%f", mf.X(), mf.Y())
	}

	d := a.DivFloat(2)
	t.Cleanup(d.Close)
	if d.X() != 1.5 || d.Y() != 2 {
		t.Fatalf("DivFloat: %f,%f", d.X(), d.Y())
	}

	n := a.Neg()
	t.Cleanup(n.Close)
	if n.X() != -3 || n.Y() != -4 {
		t.Fatalf("Neg: %f,%f", n.X(), n.Y())
	}
}

func TestQPointFDotProduct(t *testing.T) {
	a := mkPointF(t, 3, 4)
	b := mkPointF(t, 5, 6)
	// 3*5 + 4*6 = 39.0
	if got := a.DotProduct(b); got != 39.0 {
		t.Fatalf("DotProduct: %f", got)
	}
}

func TestQPointFEqual(t *testing.T) {
	a := mkPointF(t, 3, 4)
	b := mkPointF(t, 3, 4)
	if !a.Equal(b) {
		t.Fatal("Equal")
	}
	if a.NotEqual(b) {
		t.Fatal("NotEqual")
	}
}

func TestQPointFEqualQPoint(t *testing.T) {
	f := mkPointF(t, 3, 4)
	p := mkPoint(t, 3, 4)
	if !f.EqualQPoint(p) {
		t.Fatal("EqualQPoint")
	}
}

func TestQPointFToPoint(t *testing.T) {
	cases := []struct {
		fx, fy float64
		wx, wy int
	}{
		{1.4, 2.6, 1, 3},
		{-1.4, -2.6, -1, -3},
		{0.5, 0.5, 1, 1},
	}
	for _, c := range cases {
		f := mkPointF(t, c.fx, c.fy)
		p := f.ToPoint()
		t.Cleanup(p.Close)
		if p.X() != c.wx || p.Y() != c.wy {
			t.Errorf("ToPoint(%f,%f) = %d,%d, want %d,%d", c.fx, c.fy, p.X(), p.Y(), c.wx, c.wy)
		}
	}
}

func TestQPointFFromPoint(t *testing.T) {
	p := mkPoint(t, 3, 4)
	f := QPointFFromPoint(p)
	t.Cleanup(f.Close)
	if f.X() != 3 || f.Y() != 4 {
		t.Fatalf("QPointFFromPoint: %f,%f", f.X(), f.Y())
	}
}
