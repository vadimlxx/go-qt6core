package qdata

// DropAction mirrors Qt::DropAction from qnamespace.h (Qt 6.11).
type DropAction int

const (
	CopyAction       DropAction = 0x1
	MoveAction       DropAction = 0x2
	LinkAction       DropAction = 0x4
	ActionMask       DropAction = 0xff
	TargetMoveAction DropAction = 0x8002
	IgnoreAction     DropAction = 0x0
)

// DropActions is a bitmask of DropAction flags.
//
// In C++, Q_DECLARE_FLAGS(DropActions, DropAction) creates a distinct
// QFlags<DropAction> type with an implicit constructor taking DropAction,
// so individual flags can be passed anywhere a DropActions is expected.
// Go has no implicit conversions, so DropActions is declared as a type
// alias for DropAction — the ergonomic equivalent of Qt's QFlags behavior.
type DropActions = DropAction