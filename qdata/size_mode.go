package qdata

// SizeMode mirrors Qt::SizeMode from qnamespace.h (Qt 6.11).
type SizeMode int

const (
	AbsoluteSize SizeMode = 0
	RelativeSize SizeMode = 1
)