package qurl

// UrlFormattingOption mirrors QUrl::UrlFormattingOption (Qt 6.11).
//
// In C++ this enum is combined with ComponentFormattingOption via
// QUrlTwoFlags<UrlFormattingOption, ComponentFormattingOption> to form
// FormattingOptions. QUrlTwoFlags provides implicit constructors from both
// enum types, so in C++ any individual flag can be passed wherever
// FormattingOptions is expected, without explicit conversion.
//
// Go has no implicit conversions, so all three types — UrlFormattingOption,
// ComponentFormattingOption and FormattingOptions — are declared as aliases
// for int. This mirrors QUrlTwoFlags' permissiveness: any flag value can be
// passed to any method taking any of the three, and combinations produced by
// bitwise OR (as in C++) work identically.
type UrlFormattingOption = int

const (
	None                  = 0x0
	RemoveScheme          = 0x1
	RemovePassword        = 0x2
	RemoveUserInfo        = RemovePassword | 0x4
	RemovePort            = 0x8
	RemoveAuthority       = RemoveUserInfo | RemovePort | 0x10
	RemovePath            = 0x20
	RemoveQuery           = 0x40
	RemoveFragment        = 0x80
	PreferLocalFile       = 0x200
	StripTrailingSlash    = 0x400
	RemoveFilename        = 0x800
	NormalizePathSegments = 0x1000
)

// FormattingOptions mirrors QUrl::FormattingOptions
// (QUrlTwoFlags<UrlFormattingOption, ComponentFormattingOption>, Qt 6.11).
//
// In C++, QUrlTwoFlags has implicit constructors from both UrlFormattingOption
// and ComponentFormattingOption, so any individual flag is implicitly
// convertible to FormattingOptions. In Go, FormattingOptions is declared as
// an alias for int so that any flag value (whose Go types are also aliases
// for int) can be passed without explicit conversion — matching Qt's
// permissive semantics.
//
// Combinations via bitwise OR work as in C++:
//
//     u.ToString(RemoveQuery | RemoveFragment)
//     u.ToEncodedBytes(FullyEncoded | StripTrailingSlash)
type FormattingOptions = int
