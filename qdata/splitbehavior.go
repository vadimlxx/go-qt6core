// Package qdata contains shared Qt namespace types from qnamespace.h:
// enums, flags, and lightweight value types used across multiple packages
// of go-qt6core.
//
// qdata must remain free of dependencies on cgo-bound packages (qtext,
// qbytes, ...) to avoid import cycles. Anything that needs QString or
// QByteArray belongs elsewhere.
package qdata

// SplitBehavior mirrors Qt::SplitBehavior (SplitBehaviorFlags) from
// qnamespace.h (Qt 6.11).
//
// Used by QString.Split*, QStringView.Split* and QByteArray.Split*.
type SplitBehavior int

const (
	KeepEmptyParts SplitBehavior = 0x0000
	SkipEmptyParts SplitBehavior = 0x0001
)