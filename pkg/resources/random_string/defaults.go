package random_string

const (
	DefaultRandomStringSize    uint32 = 10
	DefaultRandomStringChars   string = "alphanum"
	DefaultRandomStringUniques uint32 = 0

	MaxRandomStringSize    = 100
	MaxRandomStringUniques = 100000
)

var (
	RandomStringAlphaChars   = []rune("abcdefghijklmnopqrstuvwxyz")
	RandomStringNumericChars = []rune("1234567890")
	RandomStringHexChars     = []rune("0123456789abcdef")
)
