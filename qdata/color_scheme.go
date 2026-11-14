package qdata

// ColorScheme mirrors Qt::ColorScheme (enum class) from qnamespace.h (Qt 6.11).
type ColorScheme int

const (
	ColorSchemeUnknown ColorScheme = 0
	ColorSchemeLight   ColorScheme = 1
	ColorSchemeDark    ColorScheme = 2
)