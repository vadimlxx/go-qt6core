package qdata

// TimerType mirrors Qt::TimerType from qnamespace.h (Qt 6.11).
type TimerType int

const (
	PreciseTimer    TimerType = 0
	CoarseTimer     TimerType = 1
	VeryCoarseTimer TimerType = 2
)