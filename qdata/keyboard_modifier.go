package qdata

// KeyboardModifier mirrors Qt::KeyboardModifier.
type KeyboardModifier int

const (
	NoModifier           KeyboardModifier = 0x00000000
	ShiftModifier        KeyboardModifier = 0x02000000
	ControlModifier      KeyboardModifier = 0x04000000
	AltModifier          KeyboardModifier = 0x08000000
	MetaModifier         KeyboardModifier = 0x10000000
	KeypadModifier       KeyboardModifier = 0x20000000
	GroupSwitchModifier  KeyboardModifier = 0x40000000
	KeyboardModifierMask KeyboardModifier = 0xfe000000
)

// KeyboardModifiers is a bitmask of KeyboardModifier flags.
//
// In C++, Q_DECLARE_FLAGS(KeyboardModifiers, KeyboardModifier) creates a distinct
// QFlags<KeyboardModifier> type with an implicit constructor taking KeyboardModifier,
// so individual flags (Qt::ShiftModifier, …) can be passed anywhere a KeyboardModifiers
// is expected. Go has no implicit conversions, so KeyboardModifiers is declared as a
// type alias for KeyboardModifier — the ergonomic equivalent of Qt's QFlags behavior.
//
// Constants like ShiftModifier and NoModifier can therefore be passed directly to
// functions accepting KeyboardModifiers, and values returned as KeyboardModifiers can
// be compared against individual flag constants without explicit casts.
type KeyboardModifiers = KeyboardModifier
