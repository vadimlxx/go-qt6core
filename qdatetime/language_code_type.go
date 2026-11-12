package qdatetime

// LanguageCodeType mirrors QLocale::LanguageCodeType.
type LanguageCodeType int

const (
	ISO639Part1        LanguageCodeType = 1 << 0
	ISO639Part2B       LanguageCodeType = 1 << 1
	ISO639Part2T       LanguageCodeType = 1 << 2
	ISO639Part3        LanguageCodeType = 1 << 3
	LegacyLanguageCode LanguageCodeType = 1 << 15

	ISO639Part2    LanguageCodeType = ISO639Part2B | ISO639Part2T
	ISO639Alpha2   LanguageCodeType = ISO639Part1
	ISO639Alpha3   LanguageCodeType = ISO639Part2 | ISO639Part3
	ISO639         LanguageCodeType = ISO639Alpha2 | ISO639Alpha3

	AnyLanguageCode LanguageCodeType = -1
)

// LanguageCodeTypes is a bitmask of LanguageCodeType flags.
type LanguageCodeTypes int