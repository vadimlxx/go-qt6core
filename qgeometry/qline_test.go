package qgeometry

import "testing"

// helpers ------------------------------------------------------------------

func mkLine(t *testing.T, x1, y1, x2, y2 int) QLine {
	t.Helper()
	return NewQLine(x1, y1, x2, y2)
}

func mkLineF(t *testing.T, x1, y1, x2, y2 float64) QLineF {
	t.Helper()
	return NewQLineF(x1, y1, x2, y2)
}

// ---- QLine tests ---------------------------------------------------------

func TestQLineConstruction(t *testing.T) {
	l := NewQLine(1, 2, 3, 4)
	if l.X1() != 1 || l.Y1() != 2 || l.X2() != 3 || l.Y2() != 4 {
		t.Fatalf("fields: %d,%d,%d,%d", l.X1(), l.Y1(), l.X2(), l.Y2())
	}

	z := NewQLineZero()
	if !z.IsNull() {
		t.Fatal("(0,0)-(0,0) must be null")
	}
	if l.IsNull() {
		t.Fatal("(1,2)-(3,4) must not be null")
	}
}

func TestQLineValueSemantics(t *testing.T) {
	a := NewQLine(1, 2, 3, 4)
	b := a // value copy

	b = b.SetP1(NewQPoint(100, 200))
	if a.X1() != 1 || a.Y1() != 2 {
		t.Fatalf("original changed: %d,%d", a.X1(), a.Y1())
	}
	if b.X1() != 100 || b.Y1() != 200 {
		t.Fatalf("copy not updated: %d,%d", b.X1(), b.Y1())
	}
}

func TestQLineFromPoints(t *testing.T) {
	p1 := NewQPoint(1, 2)
	p2 := NewQPoint(3, 4)
	l := NewQLineFromPoints(p1, p2)
	if l.X1() != 1 || l.Y1() != 2 || l.X2() != 3 || l.Y2() != 4 {
		t.Fatalf("FromPoints: %d,%d,%d,%d", l.X1(), l.Y1(), l.X2(), l.Y2())
	}
}

func TestQLineDxDy(t *testing.T) {
	l := NewQLine(1, 2, 4, 6)
	if l.Dx() != 3 || l.Dy() != 4 {
		t.Fatalf("Dx/Dy: %d/%d", l.Dx(), l.Dy())
	}
}

func TestQLineP1P2AreValues(t *testing.T) {
	l := NewQLine(1, 2, 3, 4)

	p1 := l.P1()
	p2 := l.P2()

	if p1.X() != 1 || p1.Y() != 2 || p2.X() != 3 || p2.Y() != 4 {
		t.Fatalf("P1/P2: (%d,%d) (%d,%d)", p1.X(), p1.Y(), p2.X(), p2.Y())
	}

	// QPoint is a value type with a pointer-receiver SetX — calling it mutates p1 in place.
	// (SetX returns nothing, so it must NOT be used in an assignment.)
	p1.SetX(100)
	if p1.X() != 100 {
		t.Fatalf("p1.SetX(100) did not mutate p1: X=%d", p1.X())
	}
	// The line itself is unaffected — p1 was a copy.
	if l.X1() != 1 {
		t.Fatalf("line changed after mutating returned point: X1=%d", l.X1())
	}
}

func TestQLineCenter(t *testing.T) {
	l := NewQLine(0, 0, 4, 4)
	c := l.Center()
	if c.X() != 2 || c.Y() != 2 {
		t.Fatalf("Center: %d,%d", c.X(), c.Y())
	}
}

func TestQLineTranslate(t *testing.T) {
	l := NewQLine(1, 2, 3, 4)
	p := NewQPoint(10, 20)

	l2 := l.Translate(p)
	if l2.X1() != 11 || l2.Y1() != 22 || l2.X2() != 13 || l2.Y2() != 24 {
		t.Fatalf("Translate: %d,%d,%d,%d", l2.X1(), l2.Y1(), l2.X2(), l2.Y2())
	}
	// original unchanged
	if l.X1() != 1 || l.Y1() != 2 {
		t.Fatal("Translate must not mutate original")
	}
}

func TestQLineTranslated(t *testing.T) {
	l := NewQLine(1, 2, 3, 4)
	p := NewQPoint(10, 20)

	tr := l.Translated(p)
	if tr.X1() != 11 || tr.Y1() != 22 || tr.X2() != 13 || tr.Y2() != 24 {
		t.Fatalf("Translated: %d,%d,%d,%d", tr.X1(), tr.Y1(), tr.X2(), tr.Y2())
	}
	if l.X1() != 1 || l.Y1() != 2 {
		t.Fatal("Translated must not mutate")
	}
}

func TestQLineSetLine(t *testing.T) {
	l := NewQLine(0, 0, 0, 0)
	l = l.SetLine(1, 2, 3, 4)
	if l.X1() != 1 || l.Y2() != 4 {
		t.Fatalf("SetLine: %d,%d,%d,%d", l.X1(), l.Y1(), l.X2(), l.Y2())
	}
}

func TestQLineEqual(t *testing.T) {
	a := NewQLine(1, 2, 3, 4)
	b := NewQLine(1, 2, 3, 4)
	c := NewQLine(1, 2, 3, 5)

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

func TestQLineToLineF(t *testing.T) {
	l := NewQLine(1, 2, 3, 4)
	f := l.ToLineF()
	if f.X1() != 1 || f.Y1() != 2 || f.X2() != 3 || f.Y2() != 4 {
		t.Fatalf("ToLineF: %f,%f,%f,%f", f.X1(), f.Y1(), f.X2(), f.Y2())
	}
}

// ---- QLineF tests --------------------------------------------------------

func TestQLineFConstruction(t *testing.T) {
	l := NewQLineF(1.5, 2.5, 3.5, 4.5)
	if l.X1() != 1.5 || l.Y1() != 2.5 || l.X2() != 3.5 || l.Y2() != 4.5 {
		t.Fatalf("fields: %f,%f,%f,%f", l.X1(), l.Y1(), l.X2(), l.Y2())
	}

	z := NewQLineFZero()
	if !z.IsNull() {
		t.Fatal("(0,0)-(0,0) must be null")
	}
}

func TestQLineFValueSemantics(t *testing.T) {
	a := NewQLineF(1, 2, 3, 4)
	b := a

	b = b.SetP1(NewQPointF(100, 200))
	if a.X1() != 1 || a.Y1() != 2 {
		t.Fatalf("original changed: %f,%f", a.X1(), a.Y1())
	}
	if b.X1() != 100 || b.Y1() != 200 {
		t.Fatalf("copy not updated: %f,%f", b.X1(), b.Y1())
	}
}

func TestQLineFFromPoints(t *testing.T) {
	p1 := NewQPointF(1, 2)
	p2 := NewQPointF(3, 4)
	l := NewQLineFFromPoints(p1, p2)
	if l.X1() != 1 || l.Y2() != 4 {
		t.Fatalf("FromPoints: %f,%f,%f,%f", l.X1(), l.Y1(), l.X2(), l.Y2())
	}
}

func TestQLineFFromLine(t *testing.T) {
	gl := NewQLine(1, 2, 3, 4)
	f := QLineFFromLine(gl)
	if f.X1() != 1 || f.Y2() != 4 {
		t.Fatalf("FromLine: %f,%f,%f,%f", f.X1(), f.Y1(), f.X2(), f.Y2())
	}
}

func TestQLineFFromPolar(t *testing.T) {
	// angle=0 → right
	l := QLineFFromPolar(5, 0)
	if l.X2() < 4.999 || l.X2() > 5.001 || l.Y2() > 0.001 || l.Y2() < -0.001 {
		t.Fatalf("FromPolar(5,0): %f,%f,%f,%f", l.X1(), l.Y1(), l.X2(), l.Y2())
	}

	// angle=90 → up on screen (negative y)
	l2 := QLineFFromPolar(5, 90)
	if l2.X2() > 0.001 || l2.X2() < -0.001 || l2.Y2() > -4.999 || l2.Y2() < -5.001 {
		t.Fatalf("FromPolar(5,90): %f,%f,%f,%f", l2.X1(), l2.Y1(), l2.X2(), l2.Y2())
	}

	// angle=180 → left
	l3 := QLineFFromPolar(5, 180)
	if l3.X2() > -4.999 || l3.Y2() > 0.001 || l3.Y2() < -0.001 {
		t.Fatalf("FromPolar(5,180): %f,%f,%f,%f", l3.X1(), l3.Y1(), l3.X2(), l3.Y2())
	}

	// angle=270 → down on screen (positive y)
	l4 := QLineFFromPolar(5, 270)
	if l4.X2() > 0.001 || l4.X2() < -0.001 || l4.Y2() < 4.999 || l4.Y2() > 5.001 {
		t.Fatalf("FromPolar(5,270): %f,%f,%f,%f", l4.X1(), l4.Y1(), l4.X2(), l4.Y2())
	}
}

func TestQLineFLengthAngle(t *testing.T) {
	l := NewQLineF(0, 0, 5, 0)
	if l.Length() != 5 {
		t.Fatalf("Length: %f", l.Length())
	}
	if l.Angle() != 0 {
		t.Fatalf("Angle: %f", l.Angle())
	}

	l = l.SetLength(10)
	if l.Length() < 9.999 || l.Length() > 10.001 {
		t.Fatalf("after SetLength: %f", l.Length())
	}
	if l.X2() < 9.999 || l.X2() > 10.001 || l.Y2() > 0.001 || l.Y2() < -0.001 {
		t.Fatalf("SetLength P2: %f,%f", l.X2(), l.Y2())
	}

	l = l.SetAngle(90)
	if l.Y2() > -9.999 || l.X2() > 0.001 || l.X2() < -0.001 {
		t.Fatalf("after SetAngle(90): %f,%f", l.X2(), l.Y2())
	}

	l = l.SetAngle(0)
	if l.X2() < 9.999 || l.X2() > 10.001 || l.Y2() > 0.001 || l.Y2() < -0.001 {
		t.Fatalf("after SetAngle(0): %f,%f", l.X2(), l.Y2())
	}

	l = l.SetAngle(270)
	if l.Y2() < 9.999 || l.X2() > 0.001 || l.X2() < -0.001 {
		t.Fatalf("after SetAngle(270): %f,%f", l.X2(), l.Y2())
	}
}

func TestQLineFAngleTo(t *testing.T) {
	h := NewQLineF(0, 0, 1, 0) // angle 0
	vUp := NewQLineF(0, 0, 0, -1) // angle 90 (up on screen)
	vDown := NewQLineF(0, 0, 0, 1) // angle 270 (down)

	if got := h.AngleTo(vUp); got < 89.999 || got > 90.001 {
		t.Fatalf("AngleTo (h→up): %f, want 90", got)
	}
	if got := h.AngleTo(vDown); got < 269.999 || got > 270.001 {
		t.Fatalf("AngleTo (h→down): %f, want 270", got)
	}
	if got := vDown.AngleTo(h); got < 89.999 || got > 90.001 {
		t.Fatalf("AngleTo (down→h): %f, want 90", got)
	}
	if got := h.AngleTo(h); got > 0.001 || got < -0.001 {
		t.Fatalf("AngleTo (h→h): %f, want 0", got)
	}
}

func TestQLineFUnitVector(t *testing.T) {
	l := NewQLineF(0, 0, 3, 4) // length 5
	u := l.UnitVector()
	if u.Length() < 0.999 || u.Length() > 1.001 {
		t.Fatalf("UnitVector length: %f", u.Length())
	}
}

func TestQLineFNormalVector(t *testing.T) {
	l := NewQLineF(0, 0, 3, 4)
	n := l.NormalVector()
	if n.Length() < 4.999 || n.Length() > 5.001 {
		t.Fatalf("NormalVector length: %f", n.Length())
	}
	p := n.P1()
	if p.X() != 0 || p.Y() != 0 {
		t.Fatalf("NormalVector P1: %f,%f", p.X(), p.Y())
	}
}

func TestQLineFPointAt(t *testing.T) {
	l := NewQLineF(0, 0, 10, 0)
	cases := []struct {
		t      float64
		wx, wy float64
	}{
		{0, 0, 0},
		{0.5, 5, 0},
		{1, 10, 0},
	}
	for _, c := range cases {
		p := l.PointAt(c.t)
		if p.X() != c.wx || p.Y() != c.wy {
			t.Errorf("PointAt(%v) = %f,%f, want %f,%f", c.t, p.X(), p.Y(), c.wx, c.wy)
		}
	}
}

func TestQLineFIntersects(t *testing.T) {
	a := NewQLineF(0, 0, 10, 10)
	b := NewQLineF(0, 10, 10, 0)

	if got := a.Intersects(b); got != BoundedIntersection {
		t.Fatalf("Intersects: %v", got)
	}

	typ, pt, ok := a.IntersectsWithPoint(b)
	if !ok || typ != BoundedIntersection {
		t.Fatalf("IntersectsWithPoint: typ=%v ok=%v", typ, ok)
	}
	if pt.X() < 4.999 || pt.X() > 5.001 || pt.Y() < 4.999 || pt.Y() > 5.001 {
		t.Fatalf("IntersectsWithPoint: %f,%f", pt.X(), pt.Y())
	}
}

func TestQLineFIntersectsParallel(t *testing.T) {
	a := NewQLineF(0, 0, 10, 0)
	b := NewQLineF(0, 1, 10, 1)

	if got := a.Intersects(b); got != NoIntersection {
		t.Fatalf("Intersects parallel: %v", got)
	}

	typ, _, ok := a.IntersectsWithPoint(b)
	if ok || typ != NoIntersection {
		t.Fatalf("IntersectsWithPoint parallel: typ=%v ok=%v", typ, ok)
	}
}

func TestQLineFEqual(t *testing.T) {
	a := NewQLineF(1, 2, 3, 4)
	b := NewQLineF(1, 2, 3, 4)
	if !a.Equal(b) {
		t.Fatal("Equal")
	}
	if a.NotEqual(b) {
		t.Fatal("NotEqual")
	}
}

func TestQLineFEqualLine(t *testing.T) {
	a := NewQLineF(1, 2, 3, 4)
	b := NewQLine(1, 2, 3, 4)
	if !a.EqualLine(b) {
		t.Fatal("EqualLine")
	}
}

func TestQLineFToLine(t *testing.T) {
	cases := []struct {
		x1, y1, x2, y2     float64
		wx1, wy1, wx2, wy2 int
	}{
		{1.4, 2.6, 3.4, 4.6, 1, 3, 3, 5},
		{-1.4, -2.6, -3.4, -4.6, -1, -3, -3, -5},
		{1, 2, 3, 4, 1, 2, 3, 4},
	}
	for _, c := range cases {
		f := NewQLineF(c.x1, c.y1, c.x2, c.y2)
		l := f.ToLine()
		if l.X1() != c.wx1 || l.Y1() != c.wy1 || l.X2() != c.wx2 || l.Y2() != c.wy2 {
			t.Errorf("ToLine(%f,%f,%f,%f) = %d,%d,%d,%d, want %d,%d,%d,%d",
				c.x1, c.y1, c.x2, c.y2, l.X1(), l.Y1(), l.X2(), l.Y2(),
				c.wx1, c.wy1, c.wx2, c.wy2)
		}
	}
}
