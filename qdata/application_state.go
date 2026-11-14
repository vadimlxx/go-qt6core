package qdata

// ApplicationState mirrors Qt::ApplicationState from qnamespace.h (Qt 6.11).
type ApplicationState int

const (
	ApplicationSuspended ApplicationState = 0x00000000
	ApplicationHidden    ApplicationState = 0x00000001
	ApplicationInactive  ApplicationState = 0x00000002
	ApplicationActive    ApplicationState = 0x00000004
)

// ApplicationStates is a bitmask of ApplicationState flags.
//
// In C++, Q_DECLARE_FLAGS(ApplicationStates, ApplicationState) creates a
// distinct QFlags<ApplicationState> type with an implicit constructor
// taking ApplicationState, so individual flags can be passed anywhere an
// ApplicationStates is expected. Go has no implicit conversions, so
// ApplicationStates is declared as a type alias for ApplicationState —
// the ergonomic equivalent of Qt's QFlags behavior.
type ApplicationStates = ApplicationState