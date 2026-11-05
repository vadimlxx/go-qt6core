package qbytearray

// Base64Option mirrors QByteArray::Base64Option.
type Base64Option int

const (
    Base64Encoding             Base64Option = 0
    Base64UrlEncoding          Base64Option = 1
    KeepTrailingEquals         Base64Option = 0
    OmitTrailingEquals         Base64Option = 2
    IgnoreBase64DecodingErrors Base64Option = 0
    AbortOnBase64DecodingErrors Base64Option = 4
)

// Base64Options is a bitmask of Base64Option values.
type Base64Options int

// Base64DecodingStatus mirrors QByteArray::Base64DecodingStatus.
type Base64DecodingStatus int

const (
    Base64DecodingStatusOk Base64DecodingStatus = iota
    Base64DecodingStatusIllegalInputLength
    Base64DecodingStatusIllegalCharacter
    Base64DecodingStatusIllegalPadding
)

func (s Base64DecodingStatus) String() string {
    switch s {
    case Base64DecodingStatusOk:
        return "Ok"
    case Base64DecodingStatusIllegalInputLength:
        return "IllegalInputLength"
    case Base64DecodingStatusIllegalCharacter:
        return "IllegalCharacter"
    case Base64DecodingStatusIllegalPadding:
        return "IllegalPadding"
    }
    return "Base64DecodingStatus(?)"
}