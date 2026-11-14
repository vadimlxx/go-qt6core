package qdata

// LayoutDirection mirrors Qt::LayoutDirection from qnamespace.h (Qt 6.11).
//
// Note: LayoutDirectionAuto has the value 2 (not 0), matching Qt.
type LayoutDirection int

const (
	LeftToRight          LayoutDirection = 0
	RightToLeft          LayoutDirection = 1
	LayoutDirectionAuto  LayoutDirection = 2
)