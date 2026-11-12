package qdata

// CallMode mirrors QInternal::CallMode.
type CallMode int

const (
	Testing CallMode = 0
	Live    CallMode = 1
)
