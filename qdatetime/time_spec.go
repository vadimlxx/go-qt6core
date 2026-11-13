package qdatetime

// TimeSpec mirrors Qt::TimeSpec.
type TimeSpec int

const (
	LocalTime     TimeSpec = 0
	UTC           TimeSpec = 1
	OffsetFromUTC TimeSpec = 2
	TimeZone      TimeSpec = 3
)