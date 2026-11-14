package qdata

// MouseEventFlag mirrors Qt::MouseEventFlag from qnamespace.h (Qt 6.11).
type MouseEventFlag int

const (
	NoMouseEventFlag             MouseEventFlag = 0x00
	MouseEventCreatedDoubleClick MouseEventFlag = 0x01
	MouseEventFlagMask           MouseEventFlag = 0xFF
)

// MouseEventFlags is a bitmask of MouseEventFlag flags.
//
// In C++, Q_DECLARE_FLAGS(MouseEventFlags, MouseEventFlag) creates a
// distinct QFlags<MouseEventFlag> type with an implicit constructor taking
// MouseEventFlag, so individual flags can be passed anywhere a
// MouseEventFlags is expected. Go has no implicit conversions, so
// MouseEventFlags is declared as a type alias for MouseEventFlag — the
// ergonomic equivalent of Qt's QFlags behavior.
type MouseEventFlags = MouseEventFlag