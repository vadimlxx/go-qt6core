package qdata

// FillRule mirrors Qt::FillRule from qnamespace.h (Qt 6.11).
type FillRule int

const (
	OddEvenFill FillRule = 0
	WindingFill FillRule = 1
)