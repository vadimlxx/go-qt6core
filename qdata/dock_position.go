package qdata

// DockPosition mirrors QInternal::DockPosition.
type DockPosition int

const (
	LeftDock   DockPosition = 0
	RightDock  DockPosition = 1
	TopDock    DockPosition = 2
	BottomDock DockPosition = 3
	DockCount  DockPosition = 4
)
