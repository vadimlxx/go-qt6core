package qdata

// WindowModality mirrors Qt::WindowModality from qnamespace.h (Qt 6.11).
type WindowModality int

const (
	NonModal         WindowModality = 0
	WindowModal      WindowModality = 1
	ApplicationModal WindowModality = 2
)