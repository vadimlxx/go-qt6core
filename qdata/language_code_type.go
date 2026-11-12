package qdata

// LanguageCodeType mirrors QLocale::LanguageCodeType.
type LanguageCodeType int

const (
	ISO639Part1        LanguageCodeType = 1 << 0
	ISO639Part2B       LanguageCodeType = 1 << 1
	ISO639Part2T       LanguageCodeType = 1 << 2
	ISO639Part3        LanguageCodeType = 1 << 3
	LegacyLanguageCode LanguageCodeType = 1 << 15

	// Composite flag combinations.
	ISO639Part2  LanguageCodeType = ISO639Part2B | ISO639Part2T
	ISO639Alpha2 LanguageCodeType = ISO639Part1
	ISO639Alpha3 LanguageCodeType = ISO639Part2 | ISO639Part3
	ISO639       LanguageCodeType = ISO639Alpha2 | ISO639Alpha3

	// AnyLanguageCode is the sentinel value (-1) used by Qt to request any
	// available code type. Note: this is intentionally -1 (all bits set),
	// matching the C++ enum value.
	AnyLanguageCode LanguageCodeType = -1
)

// LanguageCodeTypes is a bitmask of LanguageCodeType flags.
//
// In C++, Q_DECLARE_FLAGS(LanguageCodeTypes, LanguageCodeType) creates a
// distinct QFlags<LanguageCodeType> type with an implicit constructor taking
// LanguageCodeType. In Go, LanguageCodeTypes is a type alias for
// LanguageCodeType — so individual flags and composites (ISO639, ISO639Alpha2,
// …) can be passed directly wherever LanguageCodeTypes is expected, and
// combinations can be built at the call site:
//
//	codeTypes := ISO639Alpha2 | ISO639Part2B
//	QLocaleLanguageToCode(lang, codeTypes)   // LanguageCodeTypes
//
// No explicit conversions are required.
type LanguageCodeTypes = LanguageCodeType
