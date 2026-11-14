package qdata

// WindowState mirrors Qt::WindowState from qnamespace.h (Qt 6.11).
type WindowState int

const (
	WindowNoState    WindowState = 0x00000000
	WindowMinimized  WindowState = 0x00000001
	WindowMaximized  WindowState = 0x00000002
	WindowFullScreen WindowState = 0x00000004
	WindowActive     WindowState = 0x00000008
)

// WindowStates is a bitmask of WindowState flags.
//
// In C++, Q_DECLARE_FLAGS(WindowStates, WindowState) creates a distinct
// QFlags<WindowState> type with an implicit constructor taking WindowState,
// so individual flags can be passed anywhere a WindowStates is expected.
// Go has no implicit conversions, so WindowStates is declared as a type
// alias for WindowState — the ergonomic equivalent of Qt's QFlags behavior.
type WindowStates = WindowState