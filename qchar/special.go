package qchar

// SpecialCharacter constants mirror QChar::SpecialCharacter.
// Note: QChar::LastValidCodePoint (0x10FFFF) does NOT fit in a UTF-16
// code unit, so it is exposed separately as a rune.
const (
    Null                       QChar = 0x0000
    Tabulation                 QChar = 0x0009
    LineFeed                   QChar = 0x000A
    FormFeed                   QChar = 0x000C
    CarriageReturn             QChar = 0x000D
    Space                      QChar = 0x0020
    Nbsp                       QChar = 0x00A0
    SoftHyphen                 QChar = 0x00AD
    ReplacementCharacter       QChar = 0xFFFD
    ObjectReplacementCharacter QChar = 0xFFFC
    ByteOrderMark              QChar = 0xFEFF
    ByteOrderSwapped           QChar = 0xFFFE
    ParagraphSeparator         QChar = 0x2029
    LineSeparator              QChar = 0x2028
    VisualTabCharacter         QChar = 0x2192
)

// LastValidCodePoint is the highest legal Unicode code point (U+10FFFF).
// It cannot be represented as a single QChar and is therefore a rune.
const LastValidCodePoint rune = 0x10FFFF

// Surrogate range constants (UTF-16).
const (
    HighSurrogateMin uint16 = 0xD800
    HighSurrogateMax uint16 = 0xDBFF
    LowSurrogateMin  uint16 = 0xDC00
    LowSurrogateMax  uint16 = 0xDFFF
)
