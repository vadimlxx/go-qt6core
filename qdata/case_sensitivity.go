package qdata

// CaseSensitivity mirrors Qt::CaseSensitivity from qnamespace.h (Qt 6.11).
//
// Used by many QtCore APIs that compare strings. Note the value order:
// CaseInsensitive is the zero value, matching Qt.
type CaseSensitivity int

const (
	CaseInsensitive CaseSensitivity = 0
	CaseSensitive   CaseSensitivity = 1
)