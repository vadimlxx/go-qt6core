package qdatetime

// NumberOption mirrors QLocale::NumberOption.
type NumberOption int

const (
	DefaultNumberOptions          NumberOption = 0x0
	OmitGroupSeparator            NumberOption = 0x01
	RejectGroupSeparator          NumberOption = 0x02
	OmitLeadingZeroInExponent     NumberOption = 0x04
	RejectLeadingZeroInExponent   NumberOption = 0x08
	IncludeTrailingZeroesAfterDot NumberOption = 0x10
	RejectTrailingZeroesAfterDot  NumberOption = 0x20
)

// NumberOptions is a bitmask of NumberOption flags.
type NumberOptions int
