package qdatetime

// TagSeparator mirrors QLocale::TagSeparator.
type TagSeparator int

const (
	// TagSeparatorDash corresponds to QLocale::TagSeparator::Dash ('-').
	TagSeparatorDash TagSeparator = '-'
	// TagSeparatorUnderscore corresponds to QLocale::TagSeparator::Underscore ('_').
	TagSeparatorUnderscore TagSeparator = '_'
)