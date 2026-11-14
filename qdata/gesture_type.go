package qdata

// GestureType mirrors Qt::GestureType from qnamespace.h (Qt 6.11).
//
// LastGestureType is the sentinel upper bound (~0u in C++, i.e. -1 in
// signed Go int); CustomGesture is the first value available for
// application-defined gesture types.
type GestureType int

const (
	TapGesture        GestureType = 1
	TapAndHoldGesture GestureType = 2
	PanGesture        GestureType = 3
	PinchGesture      GestureType = 4
	SwipeGesture      GestureType = 5

	CustomGesture   GestureType = 0x0100
	LastGestureType GestureType = -1
)