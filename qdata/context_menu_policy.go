package qdata

// ContextMenuPolicy mirrors Qt::ContextMenuPolicy from qnamespace.h (Qt 6.11).
type ContextMenuPolicy int

const (
	NoContextMenu      ContextMenuPolicy = 0
	DefaultContextMenu ContextMenuPolicy = 1
	ActionsContextMenu ContextMenuPolicy = 2
	CustomContextMenu  ContextMenuPolicy = 3
	PreventContextMenu ContextMenuPolicy = 4
)