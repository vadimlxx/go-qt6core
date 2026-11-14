package qdata

// GestureFlag mirrors Qt::GestureFlag from qnamespace.h (Qt 6.11).
type GestureFlag int

const (
	DontStartGestureOnChildren       GestureFlag = 0x01
	ReceivePartialGestures           GestureFlag = 0x02
	IgnoredGesturesPropagateToParent GestureFlag = 0x04
)

// GestureFlags is a bitmask of GestureFlag flags.
//
// In C++, Q_DECLARE_FLAGS(GestureFlags, GestureFlag) creates a distinct
// QFlags<GestureFlag> type with an implicit constructor taking GestureFlag,
// so individual flags can be passed anywhere a GestureFlags is expected.
// Go has no implicit conversions, so GestureFlags is declared as a type
// alias for GestureFlag — the ergonomic equivalent of Qt's QFlags behavior.
type GestureFlags = GestureFlag