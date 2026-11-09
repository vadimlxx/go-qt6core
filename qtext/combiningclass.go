package qtext

// CombiningClass mirrors QChar::CombiningClass named constants.
// The instance method QChar.CombiningClass() returns a raw unsigned char,
// which you can compare against these values.
type CombiningClass uint8

const (
    Combining_BelowLeftAttached  CombiningClass = 200
    Combining_BelowAttached      CombiningClass = 202
    Combining_BelowRightAttached CombiningClass = 204
    Combining_LeftAttached       CombiningClass = 208
    Combining_RightAttached      CombiningClass = 210
    Combining_AboveLeftAttached  CombiningClass = 212
    Combining_AboveAttached      CombiningClass = 214
    Combining_AboveRightAttached CombiningClass = 216

    Combining_BelowLeft  CombiningClass = 218
    Combining_Below      CombiningClass = 220
    Combining_BelowRight CombiningClass = 222
    Combining_Left       CombiningClass = 224
    Combining_Right      CombiningClass = 226
    Combining_AboveLeft  CombiningClass = 228
    Combining_Above      CombiningClass = 230
    Combining_AboveRight CombiningClass = 232

    Combining_DoubleBelow   CombiningClass = 233
    Combining_DoubleAbove   CombiningClass = 234
    Combining_IotaSubscript CombiningClass = 240
)
