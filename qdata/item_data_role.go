package qdata

// ItemDataRole mirrors Qt::ItemDataRole from qnamespace.h (Qt 6.11).
//
// The values below UserRole are reserved for Qt's standard roles; custom
// roles start at UserRole (0x100).
type ItemDataRole int

const (
	DisplayRole             ItemDataRole = 0
	DecorationRole          ItemDataRole = 1
	EditRole                ItemDataRole = 2
	ToolTipRole             ItemDataRole = 3
	StatusTipRole           ItemDataRole = 4
	WhatsThisRole           ItemDataRole = 5

	// Metadata
	FontRole                ItemDataRole = 6
	TextAlignmentRole       ItemDataRole = 7
	BackgroundRole          ItemDataRole = 8
	ForegroundRole          ItemDataRole = 9
	CheckStateRole          ItemDataRole = 10

	// Accessibility
	AccessibleTextRole        ItemDataRole = 11
	AccessibleDescriptionRole ItemDataRole = 12

	// More general purpose
	SizeHintRole            ItemDataRole = 13
	InitialSortOrderRole    ItemDataRole = 14

	// Internal UiLib roles
	DisplayPropertyRole     ItemDataRole = 27
	DecorationPropertyRole  ItemDataRole = 28
	ToolTipPropertyRole     ItemDataRole = 29
	StatusTipPropertyRole   ItemDataRole = 30
	WhatsThisPropertyRole   ItemDataRole = 31

	// QRangeModel support for QML's required property var modelData
	RangeModelDataRole      ItemDataRole = 40
	RangeModelAdapterRole   ItemDataRole = 41

	// UserRole is the first value reserved for application-defined roles.
	UserRole ItemDataRole = 0x0100
)