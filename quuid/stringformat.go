package quuid

// StringFormat mirrors QUuid::StringFormat (Qt 6.11).
type StringFormat int

const (
	WithBraces    StringFormat = 0
	WithoutBraces StringFormat = 1
	Id128         StringFormat = 3
)