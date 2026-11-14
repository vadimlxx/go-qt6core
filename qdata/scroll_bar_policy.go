package qdata

// ScrollBarPolicy mirrors Qt::ScrollBarPolicy from qnamespace.h (Qt 6.11).
type ScrollBarPolicy int

const (
	ScrollBarAsNeeded  ScrollBarPolicy = 0
	ScrollBarAlwaysOff ScrollBarPolicy = 1
	ScrollBarAlwaysOn  ScrollBarPolicy = 2
)