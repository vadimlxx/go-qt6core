package qdata

// NumberOption mirrors QLocale::NumberOption.
type NumberOption int

const (
	DefaultNumberOptions          NumberOption = 0x00
	OmitGroupSeparator            NumberOption = 0x01
	RejectGroupSeparator          NumberOption = 0x02
	OmitLeadingZeroInExponent     NumberOption = 0x04
	RejectLeadingZeroInExponent   NumberOption = 0x08
	IncludeTrailingZeroesAfterDot NumberOption = 0x10
	RejectTrailingZeroesAfterDot  NumberOption = 0x20
)

// NumberOptions is a bitmask of NumberOption flags.
//
// In C++, Q_DECLARE_FLAGS(NumberOptions, NumberOption) creates a distinct
// QFlags<NumberOption> type with an implicit constructor taking NumberOption.
// In Go, NumberOptions is a type alias for NumberOption — so individual flags
// can be passed directly wherever NumberOptions is expected, and combined with
// the | operator at the call site:
//
//	opts := OmitGroupSeparator | RejectLeadingZeroInExponent
//	locale.SetNumberOptions(opts)         // SetNumberOptions(NumberOptions)
//
// No explicit conversions are required.
type NumberOptions = NumberOption
