package qdata

// ItemFlag mirrors Qt::ItemFlag from qnamespace.h (Qt 6.11).
type ItemFlag int

const (
	NoItemFlags          ItemFlag = 0
	ItemIsSelectable     ItemFlag = 1
	ItemIsEditable       ItemFlag = 2
	ItemIsDragEnabled    ItemFlag = 4
	ItemIsDropEnabled    ItemFlag = 8
	ItemIsUserCheckable  ItemFlag = 16
	ItemIsEnabled        ItemFlag = 32
	ItemIsAutoTristate   ItemFlag = 64
	ItemNeverHasChildren ItemFlag = 128
	ItemIsUserTristate   ItemFlag = 256
)

// ItemFlags is a bitmask of ItemFlag flags.
//
// In C++, Q_DECLARE_FLAGS(ItemFlags, ItemFlag) creates a distinct
// QFlags<ItemFlag> type with an implicit constructor taking ItemFlag, so
// individual flags can be passed anywhere an ItemFlags is expected. Go has
// no implicit conversions, so ItemFlags is declared as a type alias for
// ItemFlag — the ergonomic equivalent of Qt's QFlags behavior.
type ItemFlags = ItemFlag