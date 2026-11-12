package qdata

// Callback mirrors QInternal::Callback.
type Callback int

const (
	EventNotifyCallback Callback = 0
	LastCallback        Callback = 1
)
