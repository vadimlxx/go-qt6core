package qdata

// SortOrder mirrors Qt::SortOrder from qnamespace.h (Qt 6.11).
//
// Used by QAbstractItemModel, QSortFilterProxyModel, QHeaderView, etc.
type SortOrder int

const (
	AscendingOrder  SortOrder = 0
	DescendingOrder SortOrder = 1
)