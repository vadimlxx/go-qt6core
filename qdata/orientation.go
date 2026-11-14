package qdata

// Orientation mirrors Qt::Orientation from qnamespace.h (Qt 6.11).
type Orientation int

const (
	Horizontal Orientation = 0x1
	Vertical   Orientation = 0x2
)

// Orientations is a bitmask of Orientation flags.
//
// In C++, Q_DECLARE_FLAGS(Orientations, Orientation) creates a distinct
// QFlags<Orientation> type with an implicit constructor taking Orientation,
// so individual flags can be passed anywhere an Orientations is expected.
// Go has no implicit conversions, so Orientations is declared as a type
// alias for Orientation — the ergonomic equivalent of Qt's QFlags behavior.
type Orientations = Orientation