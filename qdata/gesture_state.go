package qdata

// GestureState mirrors Qt::GestureState from qnamespace.h (Qt 6.11).
type GestureState int

const (
	NoGesture       GestureState = 0
	GestureStarted  GestureState = 1
	GestureUpdated  GestureState = 2
	GestureFinished GestureState = 3
	GestureCanceled GestureState = 4
)