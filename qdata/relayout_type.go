package qdata

// RelayoutType mirrors QInternal::RelayoutType.
type RelayoutType int

const (
	RelayoutNormal   RelayoutType = 0
	RelayoutDragging RelayoutType = 1
	RelayoutDropped  RelayoutType = 2
)
