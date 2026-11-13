package qdatetime

// TimeType mirrors QTimeZone::TimeType.
type TimeType int

const (
	StandardTime TimeType = 0
	DaylightTime TimeType = 1
	GenericTime  TimeType = 2
)