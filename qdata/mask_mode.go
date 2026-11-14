package qdata

// MaskMode mirrors Qt::MaskMode from qnamespace.h (Qt 6.11).
type MaskMode int

const (
	MaskInColor  MaskMode = 0
	MaskOutColor MaskMode = 1
)