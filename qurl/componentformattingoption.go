package qurl

// ComponentFormattingOption mirrors QUrl::ComponentFormattingOption (Qt 6.11).
//
// See UrlFormattingOption for the rationale on the alias-of-int pattern:
// all flag types in this package are aliases for int, mirroring the
// QUrlTwoFlags-permissive behaviour of Qt's QUrl::FormattingOptions.
type ComponentFormattingOption = int

const (
	PrettyDecoded    = 0x000000
	EncodeSpaces     = 0x100000
	EncodeUnicode    = 0x200000
	EncodeDelimiters = 0x400000 | 0x800000
	EncodeReserved   = 0x1000000
	DecodeReserved   = 0x2000000

	FullyEncoded = EncodeSpaces | EncodeUnicode | EncodeDelimiters | EncodeReserved
	FullyDecoded = FullyEncoded | DecodeReserved | 0x4000000
)

// ComponentFormattingOptions mirrors QUrl::ComponentFormattingOptions
// (Q_DECLARE_FLAGS(ComponentFormattingOptions, ComponentFormattingOption)).
//
// Since ComponentFormattingOption is already an alias for int, this alias
// is effectively equivalent; it exists to preserve the API shape from
// qurl.h where both spellings appear.
type ComponentFormattingOptions = ComponentFormattingOption
