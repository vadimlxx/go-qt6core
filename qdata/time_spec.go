package qdata

// TimeSpec mirrors Qt::TimeSpec from qnamespace.h (Qt 6.11).
//
// In Qt 6 QDateTime switched to QTimeZone; TimeSpec is retained for
// source compatibility.
type TimeSpec int

const (
	LocalTime     TimeSpec = 0
	UTC           TimeSpec = 1
	OffsetFromUTC TimeSpec = 2
	TimeZone      TimeSpec = 3
)