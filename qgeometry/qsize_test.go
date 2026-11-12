package qgeometry

import "testing"

func TestQSizeConstruction(t *testing.T) {
	s := NewQSize(10, 20)
	if s.Width() != 10 || s.Height() != 20 {
		t.Fatalf("fields: %+v", s)
	}
	if !NewQSizeZero().IsNull() || !NewQSizeZero().IsEmpty() || !NewQSizeZero().IsValid() {
		t.Fatal("(0,0) must be null+empty+valid")
	}
	if NewQSizeDefault().IsValid() {
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
		s := NewQSize(c.w, c.h)
		if s.IsNull() != c.null || s.IsEmpty() != c.empty || s.IsValid() != c.valid {
			t.Errorf("(%d,%d): null=%v empty=%v valid=%v", c.w, c.h, s.IsNull(), s.IsEmpty(), s.IsValid())
		}
	}
}

func TestQSizeTransposed(t *testing.T) {
	s := NewQSize(10, 20)
	if got := s.Transposed(); got != NewQSize(20, 10) {
		t.Fatalf("Transposed: %+v", got)
	}
	s.Transpose()
	if s != NewQSize(20, 10) {
		t.Fatalf("after Transpose: %+v", s)
	}
}

func TestQSizeExpandedBounded(t *testing.T) {
	a := NewQSize(10, 20)
	b := NewQSize(15, 10)
	if got := a.ExpandedTo(b); got != NewQSize(15, 20) {
		t.Fatalf("ExpandedTo: %+v", got)
	}
	if got := a.BoundedTo(b); got != NewQSize(10, 10) {
		t.Fatalf("BoundedTo: %+v", got)
	}
}

func TestQSizeGrownShrunk(t *testing.T) {
	s := NewQSize(100, 50)
	m := NewQMargins(10, 5, 10, 5)
	if got := s.GrownBy(m); got != NewQSize(120, 60) {
		t.Fatalf("GrownBy: %+v", got)
	}
	if got := s.ShrunkBy(m); got != NewQSize(80, 40) {
		t.Fatalf("ShrunkBy: %+v", got)
	}
}

func TestQSizeArithmetic(t *testing.T) {
	a := NewQSize(10, 20)
	b := NewQSize(3, 5)
	if got := a.Add(b); got != NewQSize(13, 25) {
		t.Fatalf("Add: %+v", got)
	}
	if got := a.Sub(b); got != NewQSize(7, 15) {
		t.Fatalf("Sub: %+v", got)
	}
	if got := a.MulFloat(2); got != NewQSize(20, 40) {
		t.Fatalf("MulFloat: %+v", got)
	}
	if got := a.DivFloat(2); got != NewQSize(5, 10) {
		t.Fatalf("DivFloat: %+v", got)
	}
}

func TestQSizeScaled(t *testing.T) {
	s := NewQSize(100, 50)
	if got := s.Scaled(200, 200, IgnoreAspectRatio); got != NewQSize(200, 200) {
		t.Fatalf("Ignore: %+v", got)
	}
	if got := s.Scaled(200, 200, KeepAspectRatio); got != NewQSize(200, 100) {
		t.Fatalf("Keep: %+v", got)
	}
	if got := s.Scaled(200, 200, KeepAspectRatioByExpanding); got != NewQSize(400, 200) {
		t.Fatalf("Expand: %+v", got)
	}
}

func TestQSizeToSizeF(t *testing.T) {
	f := NewQSize(10, 20).ToSizeF()
	if f != NewQSizeF(10, 20) {
		t.Fatalf("ToSizeF: %+v", f)
	}
}

func TestQSizeFConstruction(t *testing.T) {
	s := NewQSizeF(1.5, 2.5)
	if s.Width() != 1.5 || s.Height() != 2.5 {
		t.Fatalf("fields: %+v", s)
	}
	if !NewQSizeFZero().IsNull() || !NewQSizeFZero().IsValid() {
		t.Fatal("(0,0) must be null+valid")
	}
	if NewQSizeFDefault().IsValid() {
		t.Fatal("default (-1,-1) must be invalid")
	}
	if got := QSizeFFromSize(NewQSize(3, 4)); got != NewQSizeF(3, 4) {
		t.Fatalf("QSizeFFromSize: %+v", got)
	}
}

func TestQSizeFArithmetic(t *testing.T) {
	a := NewQSizeF(10, 20)
	b := NewQSizeF(0.5, 0.25)
	if got := a.Add(b); got != NewQSizeF(10.5, 20.25) {
		t.Fatalf("Add: %+v", got)
	}
	if got := a.MulFloat(2); got != NewQSizeF(20, 40) {
		t.Fatalf("MulFloat: %+v", got)
	}
	if got := a.DivFloat(2); got != NewQSizeF(5, 10) {
		t.Fatalf("DivFloat: %+v", got)
	}
}

func TestQSizeFToSize(t *testing.T) {
	if got := NewQSizeF(1.5, 2.5).ToSize(); got != NewQSize(2, 3) {
		t.Fatalf("ToSize: %+v", got)
	}
}

func TestQSizeFEqual(t *testing.T) {
	a := NewQSizeF(10, 20)
	if !a.Equal(NewQSizeF(10, 20)) {
		t.Fatal("Equal")
	}
	if !a.EqualQSize(NewQSize(10, 20)) {
		t.Fatal("EqualQSize")
	}
	if a.NotEqual(NewQSizeF(10, 20)) {
		t.Fatal("NotEqual")
	}
}