package qdata

// Modifier mirrors Qt::Modifier — the legacy, all-caps aliases of
// Qt::KeyboardModifier preserved in qnamespace.h for source compatibility
// with pre-Qt-6 code. In Qt 7 they are expected to be deprecated in favor
// of Qt::KeyboardModifier.
//
// Values are identical to the corresponding KeyboardModifier constants.
type Modifier int

const (
	META          Modifier = Modifier(MetaModifier)
	SHIFT         Modifier = Modifier(ShiftModifier)
	CTRL          Modifier = Modifier(ControlModifier)
	ALT           Modifier = Modifier(AltModifier)
	MODIFIER_MASK Modifier = Modifier(KeyboardModifierMask)
)

// Modifiers is a bitmask of Modifier flags.
//
// In C++, Q_DECLARE_FLAGS(Modifiers, Modifier) creates a distinct
// QFlags<Modifier> type with an implicit constructor taking Modifier, so
// individual flags can be passed anywhere a Modifiers is expected. Go has no
// implicit conversions, so Modifiers is declared as a type alias for
// Modifier — the ergonomic equivalent of Qt's QFlags behavior.
//
// Because Qt::Modifier values are numerically identical to
// Qt::KeyboardModifier values, a Modifiers and a KeyboardModifiers can be
// used interchangeably in Go (they are both aliases of KeyboardModifier):
// pass ShiftModifier to Combine() without casts.
type Modifiers = Modifier
