package qdata

// ArrowType mirrors Qt::ArrowType from qnamespace.h (Qt 6.11).
type ArrowType int

const (
	NoArrow    ArrowType = 0
	UpArrow    ArrowType = 1
	DownArrow  ArrowType = 2
	LeftArrow  ArrowType = 3
	RightArrow ArrowType = 4
)