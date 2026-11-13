package qdatetime

// DateFormat mirrors Qt::DateFormat.
type DateFormat int

const (
	TextDate      DateFormat = 0
	ISODate       DateFormat = 1
	RFC2822Date   DateFormat = 8
	ISODateWithMs DateFormat = 9
)