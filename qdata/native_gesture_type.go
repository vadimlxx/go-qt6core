package qdata

// NativeGestureType mirrors Qt::NativeGestureType from qnamespace.h
// (Qt 6.11).
type NativeGestureType int

const (
	BeginNativeGesture     NativeGestureType = 0
	EndNativeGesture       NativeGestureType = 1
	PanNativeGesture       NativeGestureType = 2
	ZoomNativeGesture      NativeGestureType = 3
	SmartZoomNativeGesture NativeGestureType = 4
	RotateNativeGesture    NativeGestureType = 5
	SwipeNativeGesture     NativeGestureType = 6
)