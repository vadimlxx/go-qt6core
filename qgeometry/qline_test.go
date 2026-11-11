package qgeometry

import "testing"

// helper: create a QLine and register cleanup.
func mkLine(t *testing.T, x1, y1, x2, y2 int) *QLine {
	t.Helper()
	l := NewQLine(x1, y1, x2, y2)
	t.Cleanup(l.Close)
	return l
}

func mkLineF(t *testing.T, x1, y1, x2, y2 float64) *QLineF {
	t.Helper()
	l := NewQLineF(x1, y1, x2, y2)
	t.Cleanup(l.Close)
	return l
}

// ---- QLine tests ---------------------------------------------------------

func TestQLineConstruction(t *testing.T) {
	l := mkLine(t, 1, 2, 3, 4)
	if l.X1() != 1 || l.Y1() != 2 || l.X2() != 3 || l.Y2() != 4 {
		t.Fatalf("fields: %d,%d,%d,%d", l.X1(), l.Y1(), l.X2(), l.Y2())
	}

	z := NewQLineZero()
	t.Cleanup(z.Close)
	if !z.IsNull() {
		t.Fatal("(0,0)-(0,0) must be null")
	}
	if l.IsNull() {
		t.Fatal("(1,2)-(3,4) must not be null")
	}
}

func TestQLineFromPoints(t *testing.T) {
	p1 := mkPoint(t, 1, 2)
	p2 := mkPoint(t, 3, 4)
	l := NewQLineFromPoints(p1, p2)
	t.Cleanup(l.Close)

	if l.X1() != 1 || l.Y1() != 2 || l.X2() != 3 || l.Y2() != 4 {
		t.Fatalf("FromPoints: %d,%d,%d,%d", l.X1(), l.Y1(), l.X2(), l.Y2())
	}
}

func TestQLineHandleSharing(t *testing.T) {
	l1 := mkLine(t, 1, 2, 3, 4)
	l2 := l1 // shared pointer

	tr := mkPoint(t, 10, 20)
	l2.Translate(tr)

	if l1.X1() != 11 || l1.Y1() != 22 || l1.X2() != 13 || l1.Y2() != 24 {
		t.Fatalf("shared handle: l1=(%d,%d,%d,%d)", l1.X1(), l1.Y1(), l1.X2(), l1.Y2())
	}
}

func TestQLineCloneIndependence(t *testing.T) {
	a := mkLine(t, 1, 2, 3, 4)
	b := a.Clone()
	t.Cleanup(b.Close)

	tr := mkPoint(t, 100, 200)
	b.Translate(tr)

	if a.X1() != 1 || a.Y1() != 2 {
		t.Fatalf("original changed: %d,%d", a.X1(), a.Y1())
	}
	if b.X1() != 101 || b.Y1() != 202 {
		t.Fatalf("clone not translated: %d,%d", b.X1(), b.Y1())
	}
}

func TestQLineDxDy(t *testing.T) {
	l := mkLine(t, 1, 2, 4, 6)
	if l.Dx() != 3 || l.Dy() != 4 {
		t.Fatalf("Dx/Dy: %d/%d", l.Dx(), l.Dy())
	}
}

func TestQLineP1P2AreSnapshots(t *testing.T) {
	l := mkLine(t, 1, 2, 3, 4)

	p1 := l.P1()
	t.Cleanup(p1.Close)
	p2 := l.P2()
	t.Cleanup(p2.Close)

	if p1.X() != 1 || p1.Y() != 2 || p2.X() != 3 || p2.Y() != 4 {
		t.Fatalf("P1/P2: (%d,%d) (%d,%d)", p1.X(), p1.Y(), p2.X(), p2.Y())
	}

	// Mutating line does not affect already-returned snapshots.
	tr := mkPoint(t, 10, 20)
	l.Translate(tr)

	if p1.X() != 1 || p1.Y() != 2 {
		t.Fatalf("snapshot changed: %d,%d", p1.X(), p1.Y())
	}
	// But X1 reads current value.
	if l.X1() != 11 || l.Y1() != 22 {
		t.Fatalf("current: %d,%d", l.X1(), l.Y1())
	}
}

func TestQLineCenter(t *testing.T) {
	l := mkLine(t, 0, 0, 4, 4)
	c := l.Center()
	t.Cleanup(c.Close)
	if c.X() != 2 || c.Y() != 2 {
		t.Fatalf("Center: %d,%d", c.X(), c.Y())
	}
}

func TestQLineTranslateInPlace(t *testing.T) {
	l := mkLine(t, 1, 2, 3, 4)
	p := mkPoint(t, 10, 20)

	l.Translate(p)
	if l.X1() != 11 || l.Y1() != 22 || l.X2() != 13 || l.Y2() != 24 {
		t.Fatalf("Translate: %d,%d,%d,%d", l.X1(), l.Y1(), l.X2(), l.Y2())
	}
}

func TestQLineTranslatedReturnsNew(t *testing.T) {
	l := mkLine(t, 1, 2, 3, 4)
	p := mkPoint(t, 10, 20)

	tr := l.Translated(p)
	t.Cleanup(tr.Close)
	if tr.X1() != 11 || tr.Y1() != 22 || tr.X2() != 13 || tr.Y2() != 24 {
		t.Fatalf("Translated: %d,%d,%d,%d", tr.X1(), tr.Y1(), tr.X2(), tr.Y2())
	}
	// original unchanged
	if l.X1() != 1 || l.Y1() != 2 {
		t.Fatal("Translated must not mutate")
	}
}

func TestQLineSetP1P2(t *testing.T) {
	l := mkLine(t, 1, 2, 3, 4)

	p := mkPoint(t, 10, 20)
	l.SetP1(p)
	if l.X1() != 10 || l.Y1() != 20 {
		t.Fatalf("SetP1: %d,%d", l.X1(), l.Y1())
	}

	p2 := mkPoint(t, 30, 40)
	l.SetP2(p2)
	if l.X2() != 30 || l.Y2() != 40 {
		t.Fatalf("SetP2: %d,%d", l.X2(), l.Y2())
	}
}

func TestQLineSetLine(t *testing.T) {
	l := mkLine(t, 0, 0, 0, 0)
	l.SetLine(1, 2, 3, 4)
	if l.X1() != 1 || l.Y2() != 4 {
		t.Fatalf("SetLine: %d,%d,%d,%d", l.X1(), l.Y1(), l.X2(), l.Y2())
	}
}

func TestQLineEqual(t *testing.T) {
	a := mkLine(t, 1, 2, 3, 4)
	b := mkLine(t, 1, 2, 3, 4)
	c := mkLine(t, 1, 2, 3, 5)

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
	l := mkLine(t, 1, 2, 3, 4)
	f := l.ToLineF()
	t.Cleanup(f.Close)
	if f.X1() != 1 || f.Y1() != 2 || f.X2() != 3 || f.Y2() != 4 {
		t.Fatalf("ToLineF: %f,%f,%f,%f", f.X1(), f.Y1(), f.X2(), f.Y2())
	}
}

// ---- QLineF tests --------------------------------------------------------

func TestQLineFConstruction(t *testing.T) {
	l := mkLineF(t, 1.5, 2.5, 3.5, 4.5)
	if l.X1() != 1.5 || l.Y1() != 2.5 || l.X2() != 3.5 || l.Y2() != 4.5 {
		t.Fatalf("fields: %f,%f,%f,%f", l.X1(), l.Y1(), l.X2(), l.Y2())
	}

	z := NewQLineFZero()
	t.Cleanup(z.Close)
	if !z.IsNull() {
		t.Fatal("(0,0)-(0,0) must be null")
	}
}

func TestQLineFFromPoints(t *testing.T) {
	p1 := mkPointF(t, 1, 2)
	p2 := mkPointF(t, 3, 4)
	l := NewQLineFFromPoints(p1, p2)
	t.Cleanup(l.Close)
	if l.X1() != 1 || l.Y2() != 4 {
		t.Fatalf("FromPoints: %f,%f,%f,%f", l.X1(), l.Y1(), l.X2(), l.Y2())
	}
}

func TestQLineFFromLine(t *testing.T) {
	gl := mkLine(t, 1, 2, 3, 4)
	f := QLineFFromLine(gl)
	t.Cleanup(f.Close)
	if f.X1() != 1 || f.Y2() != 4 {
		t.Fatalf("FromLine: %f,%f,%f,%f", f.X1(), f.Y1(), f.X2(), f.Y2())
	}
}

func TestQLineFFromPolar(t *testing.T) {
	// Qt uses screen coordinates (y grows downward). The angle is measured counter-clockwise
	// from +x, but "up on screen" corresponds to NEGATIVE y.
	//
	//   angle=0   → (0,0) → ( 5,  0)   right
	//   angle=90  → (0,0) → ( 0, -5)   up on screen
	//   angle=180 → (0,0) → (-5,  0)   left
	//   angle=270 → (0,0) → ( 0,  5)   down on screen

	// angle=0
	l := QLineFFromPolar(5, 0)
	t.Cleanup(l.Close)
	if l.X1() != 0 || l.Y1() != 0 ||
		l.X2() < 4.999 || l.X2() > 5.001 || l.Y2() > 0.001 || l.Y2() < -0.001 {
		t.Fatalf("FromPolar(5,0): %f,%f,%f,%f", l.X1(), l.Y1(), l.X2(), l.Y2())
	}

	// angle=90 → up on screen → negative y
	l2 := QLineFFromPolar(5, 90)
	t.Cleanup(l2.Close)
	if l2.X2() > 0.001 || l2.X2() < -0.001 || l2.Y2() > -4.999 || l2.Y2() < -5.001 {
		t.Fatalf("FromPolar(5,90): %f,%f,%f,%f", l2.X1(), l2.Y1(), l2.X2(), l2.Y2())
	}

	// angle=180 → left
	l3 := QLineFFromPolar(5, 180)
	t.Cleanup(l3.Close)
	if l3.X2() > -4.999 || l3.Y2() > 0.001 || l3.Y2() < -0.001 {
		t.Fatalf("FromPolar(5,180): %f,%f,%f,%f", l3.X1(), l3.Y1(), l3.X2(), l3.Y2())
	}

	// angle=270 → down on screen → positive y
	l4 := QLineFFromPolar(5, 270)
	t.Cleanup(l4.Close)
	if l4.X2() > 0.001 || l4.X2() < -0.001 || l4.Y2() < 4.999 || l4.Y2() > 5.001 {
		t.Fatalf("FromPolar(5,270): %f,%f,%f,%f", l4.X1(), l4.Y1(), l4.X2(), l4.Y2())
	}
}

func TestQLineFHandleSharing(t *testing.T) {
	l1 := mkLineF(t, 1, 2, 3, 4)
	l2 := l1

	tr := mkPointF(t, 10, 20)
	l2.Translate(tr)

	if l1.X1() != 11 || l1.Y1() != 22 {
		t.Fatalf("shared handle: l1=(%f,%f,%f,%f)", l1.X1(), l1.Y1(), l1.X2(), l1.Y2())
	}
}

func TestQLineFCloneIndependence(t *testing.T) {
	a := mkLineF(t, 1, 2, 3, 4)
	b := a.Clone()
	t.Cleanup(b.Close)

	tr := mkPointF(t, 100, 200)
	b.Translate(tr)

	if a.X1() != 1 || b.X1() != 101 {
		t.Fatalf("clone independence: a.X1=%f b.X1=%f", a.X1(), b.X1())
	}
}

func TestQLineFLengthAngle(t *testing.T) {
	// Horizontal line, length 5. Angle is 0 (points right along +x).
	l := mkLineF(t, 0, 0, 5, 0)
	if l.Length() != 5 {
		t.Fatalf("Length: %f", l.Length())
	}
	if l.Angle() != 0 {
		t.Fatalf("Angle: %f", l.Angle())
	}

	// SetLength preserves direction, updates P2. Length becomes 10.
	l.SetLength(10)
	if l.Length() < 9.999 || l.Length() > 10.001 {
		t.Fatalf("after SetLength: %f", l.Length())
	}
	if l.X2() < 9.999 || l.X2() > 10.001 || l.Y2() > 0.001 || l.Y2() < -0.001 {
		t.Fatalf("SetLength P2: %f,%f", l.X2(), l.Y2())
	}

	// SetAngle(90): Qt uses screen coords (y-down), so 90° points "up on screen" → negative y.
	// P2 becomes (0, -10).
	l.SetAngle(90)
	if l.Y2() > -9.999 || l.X2() > 0.001 || l.X2() < -0.001 {
		t.Fatalf("after SetAngle(90): %f,%f", l.X2(), l.Y2())
	}

	// SetAngle(0): back to horizontal, P2 = (10, 0).
	l.SetAngle(0)
	if l.X2() < 9.999 || l.X2() > 10.001 || l.Y2() > 0.001 || l.Y2() < -0.001 {
		t.Fatalf("after SetAngle(0): %f,%f", l.X2(), l.Y2())
	}

	// SetAngle(270): down on screen → positive y. P2 = (0, 10).
	l.SetAngle(270)
	if l.Y2() < 9.999 || l.X2() > 0.001 || l.X2() < -0.001 {
		t.Fatalf("after SetAngle(270): %f,%f", l.X2(), l.Y2())
	}
}

func TestQLineFAngleTo(t *testing.T) {
	// AngleTo returns the CCW angle from this line to other, normalised to [0, 360).
	// Both lines start at origin.

	// h = points right, angle = 0.
	h := mkLineF(t, 0, 0, 1, 0)

	// vUp = points "up on screen" (negative y), angle = 90.
	// h.AngleTo(vUp) = 90 - 0 = 90.
	vUp := mkLineF(t, 0, 0, 0, -1)
	if got := h.AngleTo(vUp); got < 89.999 || got > 90.001 {
		t.Fatalf("AngleTo (h→up): %f, want 90", got)
	}

	// vDown = points "down on screen" (positive y), angle = 270.
	// h.AngleTo(vDown) = 270 - 0 = 270.
	vDown := mkLineF(t, 0, 0, 0, 1)
	if got := h.AngleTo(vDown); got < 269.999 || got > 270.001 {
		t.Fatalf("AngleTo (h→down): %f, want 270", got)
	}

	// Reverse direction: vDown.AngleTo(h) = (0 - 270 + 360) = 90.
	if got := vDown.AngleTo(h); got < 89.999 || got > 90.001 {
		t.Fatalf("AngleTo (down→h): %f, want 90", got)
	}

	// Self: h.AngleTo(h) = 0.
	if got := h.AngleTo(h); got > 0.001 || got < -0.001 {
		t.Fatalf("AngleTo (h→h): %f, want 0", got)
	}
}

func TestQLineFUnitVector(t *testing.T) {
	l := mkLineF(t, 0, 0, 3, 4) // length 5
	u := l.UnitVector()
	t.Cleanup(u.Close)
	if u.Length() < 0.999 || u.Length() > 1.001 {
		t.Fatalf("UnitVector length: %f", u.Length())
	}
}

func TestQLineFNormalVector(t *testing.T) {
	l := mkLineF(t, 0, 0, 3, 4) // length 5
	n := l.NormalVector()
	t.Cleanup(n.Close)
	if n.Length() < 4.999 || n.Length() > 5.001 {
		t.Fatalf("NormalVector length: %f", n.Length())
	}
	// P1 preserved.
	p := n.P1()
	t.Cleanup(p.Close)
	if p.X() != 0 || p.Y() != 0 {
		t.Fatalf("NormalVector P1: %f,%f", p.X(), p.Y())
	}
}

func TestQLineFPointAt(t *testing.T) {
	l := mkLineF(t, 0, 0, 10, 0)
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
		t.Cleanup(p.Close)
		if p.X() != c.wx || p.Y() != c.wy {
			t.Errorf("PointAt(%v) = %f,%f, want %f,%f", c.t, p.X(), p.Y(), c.wx, c.wy)
		}
	}
}

func TestQLineFIntersects(t *testing.T) {
	// Crossing lines.
	a := mkLineF(t, 0, 0, 10, 10)
	b := mkLineF(t, 0, 10, 10, 0)

	if got := a.Intersects(b); got != BoundedIntersection {
		t.Fatalf("Intersects: %v", got)
	}

	typ, pt := a.IntersectsWithPoint(b)
	if typ != BoundedIntersection {
		t.Fatalf("IntersectsWithPoint type: %v", typ)
	}
	if pt == nil {
		t.Fatal("IntersectsWithPoint point: nil")
	}
	t.Cleanup(pt.Close)
	if pt.X() < 4.999 || pt.X() > 5.001 || pt.Y() < 4.999 || pt.Y() > 5.001 {
		t.Fatalf("IntersectsWithPoint: %f,%f", pt.X(), pt.Y())
	}
}

func TestQLineFIntersectsParallel(t *testing.T) {
	a := mkLineF(t, 0, 0, 10, 0)
	b := mkLineF(t, 0, 1, 10, 1)

	if got := a.Intersects(b); got != NoIntersection {
		t.Fatalf("Intersects parallel: %v", got)
	}

	typ, pt := a.IntersectsWithPoint(b)
	if typ != NoIntersection {
		t.Fatalf("IntersectsWithPoint type: %v", typ)
	}
	if pt != nil {
		t.Fatal("IntersectsWithPoint must return nil point for parallel")
	}
}

func TestQLineFEqual(t *testing.T) {
	a := mkLineF(t, 1, 2, 3, 4)
	b := mkLineF(t, 1, 2, 3, 4)
	if !a.Equal(b) {
		t.Fatal("Equal")
	}
	if a.NotEqual(b) {
		t.Fatal("NotEqual")
	}
}

func TestQLineFEqualLine(t *testing.T) {
	a := mkLineF(t, 1, 2, 3, 4)
	b := mkLine(t, 1, 2, 3, 4)
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
		f := mkLineF(t, c.x1, c.y1, c.x2, c.y2)
		l := f.ToLine()
		t.Cleanup(l.Close)
		if l.X1() != c.wx1 || l.Y1() != c.wy1 || l.X2() != c.wx2 || l.Y2() != c.wy2 {
			t.Errorf("ToLine(%f,%f,%f,%f) = %d,%d,%d,%d, want %d,%d,%d,%d",
				c.x1, c.y1, c.x2, c.y2, l.X1(), l.Y1(), l.X2(), l.Y2(), c.wx1, c.wy1, c.wx2, c.wy2)
		}
	}
}
