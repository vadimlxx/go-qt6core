package qdata

// ClipOperation mirrors Qt::ClipOperation from qnamespace.h (Qt 6.11).
type ClipOperation int

const (
	NoClip        ClipOperation = 0
	ReplaceClip   ClipOperation = 1
	IntersectClip ClipOperation = 2
)