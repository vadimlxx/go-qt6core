package qdata

// ContrastPreference mirrors Qt::ContrastPreference (enum class) from
// qnamespace.h (Qt 6.11).
type ContrastPreference int

const (
	NoPreference ContrastPreference = 0
	HighContrast ContrastPreference = 1
)