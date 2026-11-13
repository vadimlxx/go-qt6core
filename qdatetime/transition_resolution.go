package qdatetime

// TransitionResolution mirrors QDateTime::TransitionResolution.
type TransitionResolution int

const (
	Reject               TransitionResolution = 0
	RelativeToBefore     TransitionResolution = 1
	RelativeToAfter      TransitionResolution = 2
	PreferBefore         TransitionResolution = 3
	PreferAfter          TransitionResolution = 4
	PreferStandard       TransitionResolution = 5
	PreferDaylightSaving TransitionResolution = 6
	LegacyBehavior       TransitionResolution = RelativeToBefore
)