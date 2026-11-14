package qdata

// DateFormat mirrors Qt::DateFormat from qnamespace.h (Qt 6.11).
//
// Note: RFC2822Date has the explicit value 8 (not 2), matching Qt.
type DateFormat int

const (
	TextDate      DateFormat = 0
	ISODate       DateFormat = 1
	RFC2822Date   DateFormat = 8
	ISODateWithMs DateFormat = 9
)