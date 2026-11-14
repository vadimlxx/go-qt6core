package qdata

// SizeHint mirrors Qt::SizeHint from qnamespace.h (Qt 6.11).
//
// NSizeHints is the number of actual hint values (count sentinel), not a
// hint itself; it exists for array sizing in QtCore.
type SizeHint int

const (
	MinimumSize    SizeHint = 0
	PreferredSize  SizeHint = 1
	MaximumSize    SizeHint = 2
	MinimumDescent SizeHint = 3
	NSizeHints     SizeHint = 4
)