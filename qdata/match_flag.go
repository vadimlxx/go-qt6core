package qdata

// MatchFlag mirrors Qt::MatchFlag from qnamespace.h (Qt 6.11).
type MatchFlag int

const (
	MatchExactly           MatchFlag = 0
	MatchContains          MatchFlag = 1
	MatchStartsWith        MatchFlag = 2
	MatchEndsWith          MatchFlag = 3
	MatchRegularExpression MatchFlag = 4
	MatchWildcard          MatchFlag = 5
	MatchFixedString       MatchFlag = 8
	MatchTypeMask          MatchFlag = 0x0F
	MatchCaseSensitive     MatchFlag = 16
	MatchWrap              MatchFlag = 32
	MatchRecursive         MatchFlag = 64
)

// MatchFlags is a bitmask of MatchFlag flags.
//
// In C++, Q_DECLARE_FLAGS(MatchFlags, MatchFlag) creates a distinct
// QFlags<MatchFlag> type with an implicit constructor taking MatchFlag, so
// individual flags can be passed anywhere a MatchFlags is expected. Go has
// no implicit conversions, so MatchFlags is declared as a type alias for
// MatchFlag — the ergonomic equivalent of Qt's QFlags behavior.
type MatchFlags = MatchFlag