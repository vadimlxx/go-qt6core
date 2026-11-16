package quuid

// Variant mirrors QUuid::Variant (Qt 6.11).
type Variant int

const (
	VarUnknown Variant = -1
	NCS        Variant = 0
	DCE        Variant = 2
	Microsoft  Variant = 6
	Reserved   Variant = 7
)