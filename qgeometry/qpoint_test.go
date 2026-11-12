package qgeometry

import "testing"

func TestQPointConstruction(t *testing.T) {
	p := NewQPoint(3, 4)
	if p.X() != 3 || p.Y() != 4 {
		t.Fatalf("fields: %+v", p)
	}
	if !NewQPointZero().IsNull() {
		t.Fatal("(0,0) must be null")
	}
	if p.IsNull() {
		t.Fatal("(3,4) must not be null")
	}
	// Zero value of QPoint struct is also (0, 0) — matching Qt's default ctor.
	var zero QPoint
	if !zero.IsNull() || zero.X() != 0 || zero.Y() != 0 {
		t.Fatal("zero value must be (0,0)")
	}
}

func TestQPointSetters(t *testing.T) {
	p := NewQPoint(1, 2)
	p.SetX(10)
	p.SetY(20)
	if p != NewQPoint(10, 20) {
		t.Fatalf("after setters: %+v", p)
	}
}

func TestQPointTransposed(t *testing.T) {
	p := NewQPoint(3, 4)
	if got := p.Transposed(); got != NewQPoint(4, 3) {
		t.Fatalf("Transposed: %+v", got)
	}
	if p != NewQPoint(3, 4) {
		t.Fatal("Transposed must not mutate")
	}
}

func TestQPointManhattanLength(t *testing.T) {
	cases := []struct {
		p    QPoint
		want int
	}{
		{NewQPoint(0, 0), 0},
		{NewQPoint(3, 4), 7},
		{NewQPoint(-3, 4), 7},
		{NewQPoint(-3, -4), 7},
	}
	for _, c := range cases {
		if got := c.p.ManhattanLength(); got != c.want {
			t.Errorf("ManhattanLength(%+v) = %d, want %d", c.p, got, c.want)
		}
	}
}

func TestQPointArithmetic(t *testing.T) {
	a := NewQPoint(3, 4)
	b := NewQPoint(1, 2)
	if got := a.Add(b); got != NewQPoint(4, 6) {
		t.Fatalf("Add: %+v", got)
	}
	if got := a.Sub(b); got != NewQPoint(2, 2) {
		t.Fatalf("Sub: %+v", got)
	}
	if got := a.MulInt(2); got != NewQPoint(6, 8) {
		t.Fatalf("MulInt: %+v", got)
	}
	if got := a.MulFloat(0.5); got != NewQPoint(2, 2) {
		t.Fatalf("MulFloat: %+v", got)
	}
	if got := a.DivFloat(2); got != NewQPoint(2, 2) {
		t.Fatalf("DivFloat: %+v", got)
	}
	if got := a.Neg(); got != NewQPoint(-3, -4) {
		t.Fatalf("Neg: %+v", got)
	}
	if got := a.UnaryPlus(); got != a {
		t.Fatalf("UnaryPlus: %+v", got)
	}
}

func TestQPointAssign(t *testing.T) {
	p := NewQPoint(3, 4)
	p.AddAssign(NewQPoint(1, 2))
	if p != NewQPoint(4, 6) {
		t.Fatalf("AddAssign: %+v", p)
	}
	p.MulIntAssign(2)
	if p != NewQPoint(8, 12) {
		t.Fatalf("MulIntAssign: %+v", p)
	}
}

func TestQPointDotProduct(t *testing.T) {
	a := NewQPoint(3, 4)
	b := NewQPoint(5, 6)
	// 3*5 + 4*6 = 15 + 24 = 39
	if got := a.DotProduct(b); got != 39 {
		t.Fatalf("DotProduct: %d", got)
	}
}

func TestQPointEqual(t *testing.T) {
	a := NewQPoint(3, 4)
	if !a.Equal(NewQPoint(3, 4)) {
		t.Fatal("Equal")
	}
	if a.NotEqual(NewQPoint(3, 4)) {
		t.Fatal("NotEqual")
	}
	if a.NotEqual(NewQPoint(4, 3)) == false {
		t.Fatal("(3,4) vs (4,3) must differ")
	}
}

func TestQPointToPointF(t *testing.T) {
	f := NewQPoint(3, 4).ToPointF()
	if f != NewQPointF(3, 4) {
		t.Fatalf("ToPointF: %+v", f)
	}
}

func TestQPointFConstruction(t *testing.T) {
	p := NewQPointF(1.5, 2.5)
	if p.X() != 1.5 || p.Y() != 2.5 {
		t.Fatalf("fields: %+v", p)
	}
	if !NewQPointFZero().IsNull() {
		t.Fatal("(0,0) must be null")
	}
	var zero QPointF
	if !zero.IsNull() || zero.X() != 0 || zero.Y() != 0 {
		t.Fatal("zero value must be (0,0)")
	}
	if got := QPointFFromPoint(NewQPoint(3, 4)); got != NewQPointF(3, 4) {
		t.Fatalf("QPointFFromPoint: %+v", got)
	}
}

func TestQPointFArithmetic(t *testing.T) {
	a := NewQPointF(3, 4)
	b := NewQPointF(0.5, 1.5)
	if got := a.Add(b); got != NewQPointF(3.5, 5.5) {
		t.Fatalf("Add: %+v", got)
	}
	if got := a.MulFloat(2); got != NewQPointF(6, 8) {
		t.Fatalf("MulFloat: %+v", got)
	}
	if got := a.DivFloat(2); got != NewQPointF(1.5, 2) {
		t.Fatalf("DivFloat: %+v", got)
	}
	if got := a.Neg(); got != NewQPointF(-3, -4) {
		t.Fatalf("Neg: %+v", got)
	}
}

func TestQPointFDotProduct(t *testing.T) {
	a := NewQPointF(3, 4)
	b := NewQPointF(5, 6)
	// 3*5 + 4*6 = 39.0
	if got := a.DotProduct(b); got != 39.0 {
		t.Fatalf("DotProduct: %v", got)
	}
}

func TestQPointFToPoint(t *testing.T) {
	cases := []struct {
		f    QPointF
		want QPoint
	}{
		{NewQPointF(1.4, 2.6), NewQPoint(1, 3)},
		{NewQPointF(-1.4, -2.6), NewQPoint(-1, -3)},
		{NewQPointF(0.5, 0.5), NewQPoint(1, 1)},
	}
	for _, c := range cases {
		if got := c.f.ToPoint(); got != c.want {
			t.Errorf("ToPoint(%+v) = %+v, want %+v", c.f, got, c.want)
		}
	}
}

func TestQPointFEqual(t *testing.T) {
	a := NewQPointF(3, 4)
	if !a.Equal(NewQPointF(3, 4)) {
		t.Fatal("Equal")
	}
	if !a.EqualQPoint(NewQPoint(3, 4)) {
		t.Fatal("EqualQPoint")
	}
}