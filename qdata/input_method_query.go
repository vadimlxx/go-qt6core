package qdata

// InputMethodQuery mirrors Qt::InputMethodQuery from qnamespace.h (Qt 6.11).
type InputMethodQuery int

const (
	ImEnabled           InputMethodQuery = 0x1
	ImCursorRectangle   InputMethodQuery = 0x2
	ImFont              InputMethodQuery = 0x4
	ImCursorPosition    InputMethodQuery = 0x8
	ImSurroundingText   InputMethodQuery = 0x10
	ImCurrentSelection  InputMethodQuery = 0x20
	ImMaximumTextLength InputMethodQuery = 0x40
	ImAnchorPosition    InputMethodQuery = 0x80
	ImHints             InputMethodQuery = 0x100
	ImPreferredLanguage InputMethodQuery = 0x200

	ImAbsolutePosition       InputMethodQuery = 0x400
	ImTextBeforeCursor       InputMethodQuery = 0x800
	ImTextAfterCursor        InputMethodQuery = 0x1000
	ImEnterKeyType           InputMethodQuery = 0x2000
	ImAnchorRectangle        InputMethodQuery = 0x4000
	ImInputItemClipRectangle InputMethodQuery = 0x8000

	ImReadOnly     InputMethodQuery = 0x10000
	ImPlatformData InputMethodQuery = 0x80000000

	ImQueryInput InputMethodQuery = ImCursorRectangle | ImCursorPosition |
		ImSurroundingText | ImCurrentSelection | ImAnchorRectangle | ImAnchorPosition
	ImQueryAll InputMethodQuery = 0xffffffff
)

// InputMethodQueries is a bitmask of InputMethodQuery flags.
//
// In C++, Q_DECLARE_FLAGS(InputMethodQueries, InputMethodQuery) creates a
// distinct QFlags<InputMethodQuery> type with an implicit constructor taking
// InputMethodQuery, so individual flags can be passed anywhere an
// InputMethodQueries is expected. Go has no implicit conversions, so
// InputMethodQueries is declared as a type alias for InputMethodQuery — the
// ergonomic equivalent of Qt's QFlags behavior.
type InputMethodQueries = InputMethodQuery