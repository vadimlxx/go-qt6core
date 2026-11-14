package qdata

// CheckState mirrors Qt::CheckState from qnamespace.h (Qt 6.11).
//
// Used by QVariant, QAbstractItemModel, QCheckBox, QStandardItem.
type CheckState int

const (
	Unchecked         CheckState = 0
	PartiallyChecked  CheckState = 1
	Checked           CheckState = 2
)