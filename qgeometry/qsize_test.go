package qgeometry

import "testing"

// helper: create a *QSize and register cleanup.
func mkSize(t *testing.T, w, h int) *QSize {
	t.Helper()
	s := NewQSize(w, h)
	t.Cleanup(s.Close)
	return s
}

func mkSizeF(t *testing.T, w, h float64) *QSizeF {
	t.Helper()
	s := NewQSizeF(w, h)
	t.Cleanup(s.Close)
	return s
}

func expectSize(t *testing.T, s *QSize, w, h int) {
	t.Helper()
	if s.Width() != w || s.Height() != h {
		t.Fatalf("QSize: got (%d,%d), want (%d,%d)", s.Width(), s.Height(), w, h)
	}
}

func expectSizeF(t *testing.T, s *QSizeF, w, h float64) {
	t.Helper()
	if s.Width() != w || s.Height() != h {
		t.Fatalf("QSizeF: got (%v,%v), want (%v,%v)", s.Width(), s.Height(), w, h)
	}
}

// ---- QSize tests ---------------------------------------------------------

func TestQSizeConstruction(t *testing.T) {
	s := mkSize(t, 10, 20)
	expectSize(t, s, 10, 20)

	z := NewQSizeZero()
	t.Cleanup(z.Close)
	if !z.IsNull() || !z.IsEmpty() || !z.IsValid() {
		t.Fatal("(0,0) must be null+empty+valid")
	}

	d := NewQSizeDefault()
	t.Cleanup(d.Close)
	if d.IsValid() {
		t.Fatal("default (-1,-1) must be invalid")
	}
}

func TestQSizePredicates(t *testing.T) {
	cases := []struct {
		w, h  int
		null  bool
		empty bool
		valid bool
	}{
		{0, 0, true, true, true},
		{10, 20, false, false, true},
		{0, 5, false, true, true},
		{-1, 5, false, true, false},
		{-1, -1, false, true, false},
	}
	for _, c := range cases {
		s := mkSize(t, c.w, c.h)
		if s.IsNull() != c.null || s.IsEmpty() != c.empty || s.IsValid() != c.valid {
			t.Errorf("(%d,%d): null=%v empty=%v valid=%v", c.w, c.h, s.IsNull(), s.IsEmpty(), s.IsValid())
		}
	}
}

func TestQSizeHandleSharing(t *testing.T) {
	a := mkSize(t, 10, 20)
	b := a // shared pointer

	b.SetWidth(100)
	if a.Width() != 100 {
		t.Fatalf("shared handle: a.Width = %d, want 100", a.Width())
	}
}

func TestQSizeCloneIndependence(t *testing.T) {
	a := mkSize(t, 10, 20)
	b := a.Clone()
	t.Cleanup(b.Close)

	b.SetWidth(100)
	if a.Width() != 10 || b.Width() != 100 {
		t.Fatalf("clone independence: a.Width=%d b.Width=%d", a.Width(), b.Width())
	}
}

func TestQSizeTransposed(t *testing.T) {
	s := mkSize(t, 10, 20)
	tp := s.Transposed()
	t.Cleanup(tp.Close)
	expectSize(t, tp, 20, 10)

	// original unchanged
	expectSize(t, s, 10, 20)
}

func TestQSizeTransposeInPlace(t *testing.T) {
	s := mkSize(t, 10, 20)
	s.Transpose()
	expectSize(t, s, 20, 10)
}

func TestQSizeExpandedBounded(t *testing.T) {
	a := mkSize(t, 10, 20)
	b := mkSize(t, 15, 10)

	exp := a.ExpandedTo(b)
	t.Cleanup(exp.Close)
	expectSize(t, exp, 15, 20)

	bnd := a.BoundedTo(b)
	t.Cleanup(bnd.Close)
	expectSize(t, bnd, 10, 10)
}

func TestQSizeGrownShrunk(t *testing.T) {
	s := mkSize(t, 100, 50)
	m := mkMargins(t, 10, 5, 10, 5)

	grown := s.GrownBy(m)
	t.Cleanup(grown.Close)
	expectSize(t, grown, 120, 60)

	shrunk := s.ShrunkBy(m)
	t.Cleanup(shrunk.Close)
	expectSize(t, shrunk, 80, 40)
}

func TestQSizeArithmetic(t *testing.T) {
	a := mkSize(t, 10, 20)
	b := mkSize(t, 3, 5)

	sum := a.Add(b)
	t.Cleanup(sum.Close)
	expectSize(t, sum, 13, 25)

	diff := a.Sub(b)
	t.Cleanup(diff.Close)
	expectSize(t, diff, 7, 15)

	prod := a.MulFloat(2)
	t.Cleanup(prod.Close)
	expectSize(t, prod, 20, 40)

	quot := a.DivFloat(2)
	t.Cleanup(quot.Close)
	expectSize(t, quot, 5, 10)

	// inputs unchanged
	expectSize(t, a, 10, 20)
	expectSize(t, b, 3, 5)
}

func TestQSizeAddAssign(t *testing.T) {
	a := mkSize(t, 10, 20)
	b := mkSize(t, 3, 5)

	a.AddAssign(b)
	expectSize(t, a, 13, 25)
	expectSize(t, b, 3, 5)
}

func TestQSizeScaled(t *testing.T) {
	s := mkSize(t, 100, 50)

	ig := s.Scaled(200, 200, IgnoreAspectRatio)
	t.Cleanup(ig.Close)
	expectSize(t, ig, 200, 200)

	ka := s.Scaled(200, 200, KeepAspectRatio)
	t.Cleanup(ka.Close)
	expectSize(t, ka, 200, 100)

	ke := s.Scaled(200, 200, KeepAspectRatioByExpanding)
	t.Cleanup(ke.Close)
	expectSize(t, ke, 400, 200)
}

func TestQSizeScaleInPlace(t *testing.T) {
	s := mkSize(t, 100, 50)
	s.Scale(200, 200, KeepAspectRatio)
	expectSize(t, s, 200, 100)
}

func TestQSizeToSizeF(t *testing.T) {
	s := mkSize(t, 10, 20)
	f := s.ToSizeF()
	t.Cleanup(f.Close)
	expectSizeF(t, f, 10, 20)
}

// ---- QSizeF tests --------------------------------------------------------

func TestQSizeFConstruction(t *testing.T) {
	s := mkSizeF(t, 1.5, 2.5)
	expectSizeF(t, s, 1.5, 2.5)

	z := NewQSizeFZero()
	t.Cleanup(z.Close)
	if !z.IsNull() || !z.IsValid() {
		t.Fatal("(0,0) must be null+valid")
	}

	d := NewQSizeFDefault()
	t.Cleanup(d.Close)
	if d.IsValid() {
		t.Fatal("default (-1,-1) must be invalid")
	}
}

func TestQSizeFHandleSharing(t *testing.T) {
	a := mkSizeF(t, 10, 20)
	b := a
	b.SetWidth(100)
	if a.Width() != 100 {
		t.Fatalf("shared handle: a.Width = %v", a.Width())
	}
}

func TestQSizeFCloneIndependence(t *testing.T) {
	a := mkSizeF(t, 10, 20)
	b := a.Clone()
	t.Cleanup(b.Close)
	b.SetWidth(100)
	if a.Width() != 10 || b.Width() != 100 {
		t.Fatalf("clone independence: a.Width=%v b.Width=%v", a.Width(), b.Width())
	}
}

func TestQSizeFFromSize(t *testing.T) {
	s := mkSize(t, 3, 4)
	f := QSizeFFromSize(s)
	t.Cleanup(f.Close)
	expectSizeF(t, f, 3, 4)
}

func TestQSizeFArithmetic(t *testing.T) {
	a := mkSizeF(t, 10, 20)
	b := mkSizeF(t, 0.5, 0.25)

	sum := a.Add(b)
	t.Cleanup(sum.Close)
	expectSizeF(t, sum, 10.5, 20.25)

	prod := a.MulFloat(2)
	t.Cleanup(prod.Close)
	expectSizeF(t, prod, 20, 40)

	quot := a.DivFloat(2)
	t.Cleanup(quot.Close)
	expectSizeF(t, quot, 5, 10)

	// inputs unchanged
	expectSizeF(t, a, 10, 20)
	expectSizeF(t, b, 0.5, 0.25)
}

func TestQSizeFToSize(t *testing.T) {
	f := mkSizeF(t, 1.5, 2.5)
	s := f.ToSize()
	t.Cleanup(s.Close)
	expectSize(t, s, 2, 3)
}

func TestQSizeFEqual(t *testing.T) {
	a := mkSizeF(t, 10, 20)
	b := mkSizeF(t, 10, 20)

	if !a.Equal(b) {
		t.Fatal("Equal")
	}
	if a.NotEqual(b) {
		t.Fatal("NotEqual false")
	}

	s := mkSize(t, 10, 20)
	if !a.EqualQSize(s) {
		t.Fatal("EqualQSize")
	}
}

func TestQSizeFGrownShrunk(t *testing.T) {
	s := mkSizeF(t, 100, 50)
	m := mkMarginsF(t, 10, 5, 10, 5)

	grown := s.GrownBy(m)
	t.Cleanup(grown.Close)
	expectSizeF(t, grown, 120, 60)

	shrunk := s.ShrunkBy(m)
	t.Cleanup(shrunk.Close)
	expectSizeF(t, shrunk, 80, 40)
}

func TestQSizeFScaleInPlace(t *testing.T) {
	s := mkSizeF(t, 100, 50)
	s.Scale(200, 200, KeepAspectRatio)
	expectSizeF(t, s, 200, 100)
}

func TestQSizeFTransposeInPlace(t *testing.T) {
	s := mkSizeF(t, 10, 20)
	s.Transpose()
	expectSizeF(t, s, 20, 10)
}
