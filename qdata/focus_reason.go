package qdata

// FocusReason mirrors Qt::FocusReason from qnamespace.h (Qt 6.11).
//
// Reported by QFocusEvent::reason.
type FocusReason int

const (
	MouseFocusReason        FocusReason = 0
	TabFocusReason          FocusReason = 1
	BacktabFocusReason      FocusReason = 2
	ActiveWindowFocusReason FocusReason = 3
	PopupFocusReason        FocusReason = 4
	ShortcutFocusReason     FocusReason = 5
	MenuBarFocusReason      FocusReason = 6
	OtherFocusReason        FocusReason = 7
	NoFocusReason           FocusReason = 8
)