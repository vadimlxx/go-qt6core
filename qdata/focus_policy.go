package qdata

// FocusPolicy mirrors Qt::FocusPolicy from qnamespace.h (Qt 6.11).
//
// Note the combined values (StrongFocus, WheelFocus) have non-contiguous
// bit patterns, matching Qt.
type FocusPolicy int

const (
	NoFocus     FocusPolicy = 0
	TabFocus    FocusPolicy = 0x1
	ClickFocus  FocusPolicy = 0x2
	StrongFocus FocusPolicy = TabFocus | ClickFocus | 0x8
	WheelFocus  FocusPolicy = StrongFocus | 0x4
)