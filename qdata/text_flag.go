package qdata

// TextFlag mirrors Qt::TextFlag from qnamespace.h (Qt 6.11).
//
// In C++ TextFlag is combined with AlignmentFlag and TextElideMode via
// Q_DECLARE_MIXED_ENUM_OPERATORS_SYMMETRIC. There is no separate
// TextFlags QFlags type — the constants are used as plain int flags.
type TextFlag int

const (
	TextSingleLine            TextFlag = 0x0100
	TextDontClip              TextFlag = 0x0200
	TextExpandTabs            TextFlag = 0x0400
	TextShowMnemonic          TextFlag = 0x0800
	TextWordWrap              TextFlag = 0x1000
	TextWrapAnywhere          TextFlag = 0x2000
	TextDontPrint             TextFlag = 0x4000
	TextIncludeTrailingSpaces TextFlag = 0x08000000
	TextHideMnemonic          TextFlag = 0x8000
	TextJustificationForced   TextFlag = 0x10000
	TextForceLeftToRight      TextFlag = 0x20000
	TextForceRightToLeft      TextFlag = 0x40000
	TextLongestVariant        TextFlag = 0x80000
)