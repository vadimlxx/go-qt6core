package qurl

// AceProcessingOption mirrors QUrl::AceProcessingOption (Qt 6.11).
type AceProcessingOption int

const (
	IgnoreIDNWhitelist         AceProcessingOption = 0x1
	AceTransitionalProcessing  AceProcessingOption = 0x2
)

// AceProcessingOptions is a bitmask of AceProcessingOption.
type AceProcessingOptions = AceProcessingOption