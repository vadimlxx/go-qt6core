package qurl

// ParsingMode mirrors QUrl::ParsingMode (Qt 6.11).
type ParsingMode int

const (
	TolerantMode ParsingMode = 0
	StrictMode   ParsingMode = 1
	DecodedMode  ParsingMode = 2
)