package qdatetime

// FormatType mirrors QLocale::FormatType.
type FormatType int

const (
	LongFormat   FormatType = 0
	ShortFormat  FormatType = 1
	NarrowFormat FormatType = 2
)