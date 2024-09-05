package random_string // nolint:revive

const (
	// DefaultRandomStringSize is the default size of the random string.
	DefaultRandomStringSize uint32 = 10
	// DefaultRandomStringChars is the default character set of the random string.
	DefaultRandomStringChars string = "alphanum"
	// DefaultRandomStringUniques is the default number of unique random strings to generate.
	DefaultRandomStringUniques uint32 = 0

	// MaxRandomStringSize is the maximum size of the random string.
	MaxRandomStringSize = 100
	//
	MaxRandomStringUniques = 100000
)

var (
	// RandomStringAlphaChars is the set of alpha characters.
	RandomStringAlphaChars = []rune("abcdefghijklmnopqrstuvwxyz") // nolint:gochecknoglobals
	// RandomStringAlphaNumChars is the set of alpha numeric characters.
	RandomStringNumericChars = []rune("1234567890") // nolint:gochecknoglobals
	// RandomStringAlphaNumChars is the set of alpha numeric characters.
	RandomStringHexChars = []rune("0123456789abcdef") // nolint:gochecknoglobals
)
