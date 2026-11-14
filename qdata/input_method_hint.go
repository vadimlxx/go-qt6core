package qdata

// InputMethodHint mirrors Qt::InputMethodHint from qnamespace.h (Qt 6.11).
type InputMethodHint int

const (
	ImhNone InputMethodHint = 0x0

	ImhHiddenText       InputMethodHint = 0x1
	ImhSensitiveData    InputMethodHint = 0x2
	ImhNoAutoUppercase  InputMethodHint = 0x4
	ImhPreferNumbers    InputMethodHint = 0x8
	ImhPreferUppercase  InputMethodHint = 0x10
	ImhPreferLowercase  InputMethodHint = 0x20
	ImhNoPredictiveText InputMethodHint = 0x40

	ImhDate InputMethodHint = 0x80
	ImhTime InputMethodHint = 0x100

	ImhPreferLatin InputMethodHint = 0x200

	ImhMultiLine InputMethodHint = 0x400

	ImhNoEditMenu    InputMethodHint = 0x800
	ImhNoTextHandles InputMethodHint = 0x1000

	ImhDigitsOnly             InputMethodHint = 0x10000
	ImhFormattedNumbersOnly   InputMethodHint = 0x20000
	ImhUppercaseOnly          InputMethodHint = 0x40000
	ImhLowercaseOnly          InputMethodHint = 0x80000
	ImhDialableCharactersOnly InputMethodHint = 0x100000
	ImhEmailCharactersOnly    InputMethodHint = 0x200000
	ImhUrlCharactersOnly      InputMethodHint = 0x400000
	ImhLatinOnly              InputMethodHint = 0x800000

	ImhExclusiveInputMask InputMethodHint = 0xffff0000
)

// InputMethodHints is a bitmask of InputMethodHint flags.
//
// In C++, Q_DECLARE_FLAGS(InputMethodHints, InputMethodHint) creates a
// distinct QFlags<InputMethodHint> type with an implicit constructor taking
// InputMethodHint, so individual flags can be passed anywhere an
// InputMethodHints is expected. Go has no implicit conversions, so
// InputMethodHints is declared as a type alias for InputMethodHint — the
// ergonomic equivalent of Qt's QFlags behavior.
type InputMethodHints = InputMethodHint