package qdata

// DayOfWeek mirrors Qt::DayOfWeek from qnamespace.h (Qt 6.11).
//
// Monday is 1 (not 0), matching Qt.
type DayOfWeek int

const (
	Monday    DayOfWeek = 1
	Tuesday   DayOfWeek = 2
	Wednesday DayOfWeek = 3
	Thursday  DayOfWeek = 4
	Friday    DayOfWeek = 5
	Saturday  DayOfWeek = 6
	Sunday    DayOfWeek = 7
)