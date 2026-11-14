package qdata

// ContextMenuTrigger mirrors Qt::ContextMenuTrigger (enum class) from
// qnamespace.h (Qt 6.11).
type ContextMenuTrigger int

const (
	ContextMenuTriggerPress   ContextMenuTrigger = 0
	ContextMenuTriggerRelease ContextMenuTrigger = 1
)