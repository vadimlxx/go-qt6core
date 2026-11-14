package qdata

// TextFormat mirrors Qt::TextFormat from qnamespace.h (Qt 6.11).
type TextFormat int

const (
	PlainText    TextFormat = 0
	RichText     TextFormat = 1
	AutoText     TextFormat = 2
	MarkdownText TextFormat = 3
)