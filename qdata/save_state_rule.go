package qdata

// SaveStateRule mirrors QInternal::SaveStateRule.
type SaveStateRule int

const (
	KeepSavedState  SaveStateRule = 0
	ClearSavedState SaveStateRule = 1
)
