package qdata

// NavigationMode mirrors Qt::NavigationMode from qnamespace.h (Qt 6.11).
type NavigationMode int

const (
	NavigationModeNone              NavigationMode = 0
	NavigationModeKeypadTabOrder    NavigationMode = 1
	NavigationModeKeypadDirectional NavigationMode = 2
	NavigationModeCursorAuto        NavigationMode = 3
	NavigationModeCursorForceVisible NavigationMode = 4
)