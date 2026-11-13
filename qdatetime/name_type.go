package qdatetime

// NameType mirrors QTimeZone::NameType.
type NameType int

const (
	DefaultName NameType = 0
	LongName    NameType = 1
	ShortName   NameType = 2
	OffsetName  NameType = 3
)