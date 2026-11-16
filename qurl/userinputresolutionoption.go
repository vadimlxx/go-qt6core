package qurl

// UserInputResolutionOption mirrors QUrl::UserInputResolutionOption (Qt 6.11).
type UserInputResolutionOption int

const (
	DefaultResolution UserInputResolutionOption = 0
	AssumeLocalFile   UserInputResolutionOption = 1
)

// UserInputResolutionOptions is a bitmask of UserInputResolutionOption.
type UserInputResolutionOptions = UserInputResolutionOption