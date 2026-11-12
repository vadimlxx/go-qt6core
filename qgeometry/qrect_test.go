package qgeometry

import "testing"

// helpers ------------------------------------------------------------------

func mkRect(t *testing.T, left, top, width, height int) QRect {
	t.Helper()
	return NewQRect(left, top, width, height)
}

func mkRectF(t *testing.T, left, top, width, height float64) QRectF {
	t.Helper()
	return NewQRectF(left, top, width, height)
}

func expectRect(t *testing.T, r QRect, l, tp, rt, b int) {
	t.Helper()
	if r.Left() != l || r.Top() != tp || r.Right() != rt || r.Bottom() != b {
		t.Fatalf("QRect: (%d,%d,%d,%d), want (%d,%d,%d,%d)",
			r.Left(), r.Top(), r.Right(), r.Bottom(), l, tp, rt, b)
	}
}

func expectRectF(t *testing.T, r QRectF, l, tp, rt, b float64) {
	t.Helper()
	if r.Left() != l || r.Top() != tp || r.Right() != rt || r.Bottom() != b {
		t.Fatalf("QRectF: (%v,%v,%v,%v), want (%v,%v,%v,%v)",
			r.Left(), r.Top(), r.Right(), r.Bottom(), l, tp, rt, b)
	}
}

// ---- QRect tests ---------------------------------------------------------

func TestQRectConstruction(t *testing.T) {
	// Default: (0, 0, -1, -1) — null, empty, invalid.
	d := NewQRectDefault()
	if !d.IsNull() || !d.IsEmpty() || d.IsValid() {
		t.Fatalf("default must be null+empty+invalid: null=%v empty=%v valid=%v",
			d.IsNull(), d.IsEmpty(), d.IsValid())
	}

	// QRect(10, 20, 100, 50): left=10, top=20, right=109, bottom=69, width=100, height=50.
	r := mkRect(t, 10, 20, 100, 50)
	expectRect(t, r, 10, 20, 109, 69)
	if r.Width() != 100 || r.Height() != 50 {
		t.Fatalf("Width/Height: %d/%d, want 100/50", r.Width(), r.Height())
	}
	if r.IsNull() || r.IsEmpty() || !r.IsValid() {
		t.Fatal("(10,20,100,50) must be non-null, non-empty, valid")
	}

	// Go zero value is (0,0,0,0) — valid but NOT null.
	var zero QRect
	if zero.IsNull() || zero.IsEmpty() || !zero.IsValid() {
		t.Fatalf("zero value: null=%v empty=%v valid=%v", zero.IsNull(), zero.IsEmpty(), zero.IsValid())
	}
}

func TestQRectFromPoints(t *testing.T) {
	tl := NewQPoint(1, 2)
	br := NewQPoint(10, 20)
	r := NewQRectFromPoints(tl, br)
	expectRect(t, r, 1, 2, 10, 20)
}

func TestQRectFromPointSize(t *testing.T) {
	tl := NewQPoint(5, 5)
	sz := NewQSize(100, 50)
	r := NewQRectFromPointSize(tl, sz)
	expectRect(t, r, 5, 5, 104, 54)
}

func TestQRectValueSemantics(t *testing.T) {
	a := NewQRect(0, 0, 10, 10)
	b := a // value copy
	b = b.SetWidth(100)
	if a.Width() != 10 || b.Width() != 100 {
		t.Fatalf("value copy independence: a.Width=%d b.Width=%d", a.Width(), b.Width())
	}
}

func TestQRectSetRectVsSetCoords(t *testing.T) {
	r := NewQRect(0, 0, 1, 1)
	r = r.SetRect(10, 20, 30, 40)
	expectRect(t, r, 10, 20, 39, 59)

	r2 := NewQRect(0, 0, 1, 1)
	r2 = r2.SetCoords(10, 20, 30, 40)
	expectRect(t, r2, 10, 20, 30, 40)
}

func TestQRectTranslate(t *testing.T) {
	r := NewQRect(1, 2, 10, 20) // (1,2) .. (10,21)
	r2 := r.Translate(100, 200)
	expectRect(t, r2, 101, 202, 110, 221)
	// original unchanged
	expectRect(t, r, 1, 2, 10, 21)
}

func TestQRectMoveTo(t *testing.T) {
	r := NewQRect(10, 20, 100, 50) // (10,20)..(109,69), 100×50
	r2 := r.MoveTo(0, 0)
	expectRect(t, r2, 0, 0, 99, 49)
}

func TestQRectUnited(t *testing.T) {
	a := NewQRect(0, 0, 10, 10) // (0,0)..(9,9)
	b := NewQRect(5, 5, 10, 10) // (5,5)..(14,14)

	u := a.United(b)
	expectRect(t, u, 0, 0, 14, 14)
	expectRect(t, a, 0, 0, 9, 9)
	expectRect(t, b, 5, 5, 14, 14)
}

func TestQRectIntersected(t *testing.T) {
	a := NewQRect(0, 0, 10, 10)
	b := NewQRect(5, 5, 10, 10)
	i := a.Intersected(b)
	expectRect(t, i, 5, 5, 9, 9)
}

func TestQRectContainsRect(t *testing.T) {
	outer := NewQRect(0, 0, 100, 100)   // (0,0)..(99,99)
	inner := NewQRect(10, 10, 20, 20)   // (10,10)..(29,29)
	partial := NewQRect(90, 90, 20, 20) // extends past outer

	if !outer.ContainsRect(inner, false) {
		t.Fatal("outer must contain inner")
	}
	if outer.ContainsRect(partial, false) {
		t.Fatal("outer must not contain partial")
	}
}

func TestQRectContainsPoint(t *testing.T) {
	r := NewQRect(0, 0, 100, 100)
	cases := []struct {
		x, y int
		want bool
	}{
		{50, 50, true},
		{0, 0, true},
		{99, 99, true},
		{-1, 0, false},
		{100, 0, false},
	}
	for _, c := range cases {
		p := NewQPoint(c.x, c.y)
		if got := r.ContainsPoint(p, false); got != c.want {
			t.Errorf("ContainsPoint(%d,%d) = %v, want %v", c.x, c.y, got, c.want)
		}
		if got := r.ContainsXY(c.x, c.y); got != c.want {
			t.Errorf("ContainsXY(%d,%d) = %v, want %v", c.x, c.y, got, c.want)
		}
	}
}

func TestQRectIntersects(t *testing.T) {
	a := NewQRect(0, 0, 10, 10)
	b := NewQRect(5, 5, 10, 10)
	c := NewQRect(100, 100, 10, 10)
	if !a.Intersects(b) {
		t.Fatal("a must intersect b")
	}
	if a.Intersects(c) {
		t.Fatal("a must not intersect c")
	}
}

func TestQRectMarginsAddedRemoved(t *testing.T) {
	r := NewQRect(10, 20, 100, 50) // (10,20)..(109,69)
	m := NewQMargins(5, 10, 15, 20)

	added := r.MarginsAdded(m)
	expectRect(t, added, 5, 10, 124, 89)

	removed := r.MarginsRemoved(m)
	expectRect(t, removed, 15, 30, 94, 49)
}

func TestQRectEqual(t *testing.T) {
	a := NewQRect(1, 2, 10, 20)
	b := NewQRect(1, 2, 10, 20)
	c := NewQRect(1, 2, 10, 21)
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

func TestQRectToRectF(t *testing.T) {
	r := NewQRect(10, 20, 100, 50) // (10,20)..(109,69), width 100, height 50
	f := r.ToRectF()
	// QRectF from QRect: xp=10, yp=20, w=right-left+1=100, h=bottom-top+1=50
	expectRectF(t, f, 10, 20, 110, 70)
}

func TestQRectSpan(t *testing.T) {
	p1 := NewQPoint(10, 5)
	p2 := NewQPoint(2, 20)
	r := QRectSpan(p1, p2)
	// min x=2, min y=5, max x=10, max y=20
	expectRect(t, r, 2, 5, 10, 20)
}

func TestQRectNormalized(t *testing.T) {
	r := NewQRect(10, 20, 1, 1)
	r = r.SetCoords(10, 20, 0, 5) // x1=10 > x2=0, y1=20 > y2=5
	if r.IsValid() {
		t.Fatal("non-normalized must be invalid")
	}
	n := r.Normalized()

	// Qt's QRect::normalized() for backward rects uses (x2+1, x1-1) and (y2+1, y1-1),
	// preserving |width| and |height| in the inclusive-coordinate sense.
	//   original: x1=10, x2=0   → width  = 0-10+1 = -9   → |width|  = 9
	//   normalized: x1=1, x2=9  → width  = 9-1+1  =  9
	//   original: y1=20, y2=5   → height = 5-20+1 = -14  → |height| = 14
	//   normalized: y1=6, y2=19 → height = 19-6+1 =  14
	//
	// NOTE: this differs from QRectF::normalized(), which uses (xp+w, -w) and
	// (yp+h, -h) — see TestQRectFNormalized. Both preserve |size|, but the
	// inclusive vs continuous coordinate models give different corner values.
	expectRect(t, n, 1, 6, 9, 19)
	if !n.IsValid() {
		t.Fatal("normalized must be valid")
	}
}

// ---- QRectF tests --------------------------------------------------------

func TestQRectFConstruction(t *testing.T) {
	d := NewQRectFDefault()
	if !d.IsNull() || !d.IsEmpty() || d.IsValid() {
		t.Fatalf("default must be null+empty+invalid: null=%v empty=%v valid=%v",
			d.IsNull(), d.IsEmpty(), d.IsValid())
	}

	r := mkRectF(t, 10, 20, 100, 50)
	expectRectF(t, r, 10, 20, 110, 70)
	if r.Width() != 100 || r.Height() != 50 {
		t.Fatalf("Width/Height: %v/%v, want 100/50", r.Width(), r.Height())
	}
	if r.IsNull() || r.IsEmpty() || !r.IsValid() {
		t.Fatal("(10,20,100,50) must be non-null, non-empty, valid")
	}

	// Go zero value == QRectF().
	var zero QRectF
	if !zero.Equal(NewQRectFZero()) {
		t.Fatal("zero value must equal NewQRectFZero")
	}
}

func TestQRectFFromRect(t *testing.T) {
	r := NewQRect(10, 20, 100, 50) // (10,20)..(109,69)
	f := QRectFFromRect(r)
	expectRectF(t, f, 10, 20, 110, 70)
}

func TestQRectFValueSemantics(t *testing.T) {
	a := NewQRectF(0, 0, 10, 10)
	b := a
	b = b.SetWidth(100)
	if a.Width() != 10 || b.Width() != 100 {
		t.Fatalf("value copy independence: a.Width=%v b.Width=%v", a.Width(), b.Width())
	}
}

func TestQRectFTranslate(t *testing.T) {
	r := NewQRectF(1, 2, 10, 20) // (1,2)..(11,22)
	r2 := r.Translate(100, 200)
	expectRectF(t, r2, 101, 202, 111, 222)
	expectRectF(t, r, 1, 2, 11, 22)
}

func TestQRectFUnitedIntersected(t *testing.T) {
	a := NewQRectF(0, 0, 10, 10) // (0,0)..(10,10)
	b := NewQRectF(5, 5, 10, 10) // (5,5)..(15,15)

	u := a.United(b)
	expectRectF(t, u, 0, 0, 15, 15)

	i := a.Intersected(b)
	expectRectF(t, i, 5, 5, 10, 10)
}

func TestQRectFContainsPoint(t *testing.T) {
	r := NewQRectF(0, 0, 100, 100)
	cases := []struct {
		x, y float64
		want bool
	}{
		{50, 50, true},
		{0, 0, true},
		{100, 100, true},
		{-0.5, 0, false},
		{100.5, 0, false},
	}
	for _, c := range cases {
		p := NewQPointF(c.x, c.y)
		if got := r.ContainsPoint(p); got != c.want {
			t.Errorf("ContainsPoint(%v,%v) = %v, want %v", c.x, c.y, got, c.want)
		}
		if got := r.ContainsXY(c.x, c.y); got != c.want {
			t.Errorf("ContainsXY(%v,%v) = %v, want %v", c.x, c.y, got, c.want)
		}
	}
}

func TestQRectFEqual(t *testing.T) {
	a := NewQRectF(1, 2, 10, 20)
	b := NewQRectF(1, 2, 10, 20)
	if !a.Equal(b) {
		t.Fatal("Equal")
	}
	if a.NotEqual(b) {
		t.Fatal("NotEqual false")
	}
}

func TestQRectFEqualRect(t *testing.T) {
	f := NewQRectF(10, 20, 100, 50) // (10,20)..(110,70)
	r := NewQRect(10, 20, 100, 50)  // (10,20)..(109,69)

	// QRectF == QRect uses topLeft equality + size equality.
	// topLeft: (10,20) == (10,20) ✓
	// size: QSizeF(100,50) == QSize(100,50) ✓ (comparesEqual is fuzzy but equal here)
	if !f.EqualRect(r) {
		t.Fatal("EqualRect: should match")
	}
}

func TestQRectFToRect(t *testing.T) {
	r := NewQRectF(10, 20, 100, 50)
	qr := r.ToRect()
	expectRect(t, qr, 10, 20, 109, 69) // QRect: right = 10+100-1, bottom = 20+50-1
}

func TestQRectFToAlignedRect(t *testing.T) {
	// QRectF(1.5, 2.5, 3, 4) covers x∈[1.5, 4.5], y∈[2.5, 6.5]
	// floor/ceil → x∈[1, 5], y∈[2, 7] → QRect(1, 2, 4, 5) → (1,2)..(4,6)
	r := NewQRectF(1.5, 2.5, 3, 4)
	qr := r.ToAlignedRect()
	expectRect(t, qr, 1, 2, 4, 6)
}

func TestQRectFMarginsAdded(t *testing.T) {
	r := NewQRectF(10, 20, 100, 50) // (10,20)..(110,70)
	m := NewQMarginsF(5, 10, 15, 20)

	out := r.MarginsAdded(m)
	// left-5, top-10, right+15, bottom+20
	expectRectF(t, out, 5, 10, 125, 90)
}

func TestQRectFNormalized(t *testing.T) {
	r := NewQRectF(10, 20, 1, 1)
	r = r.SetCoords(10, 20, 0, 5)
	if r.IsValid() {
		t.Fatal("non-normalized must be invalid")
	}
	n := r.Normalized()
	// normalized sets left=0, top=5, right=10, bottom=20
	expectRectF(t, n, 0, 5, 10, 20)
	if !n.IsValid() {
		t.Fatal("normalized must be valid")
	}
}
