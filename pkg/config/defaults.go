package config

import "time"

const (
	DefaultIntegerRangeMin  uint32 = 0
	DefaultIntegerRangeMax  uint32 = 10
	DefaultIntegerRangeStep uint32 = 1
	DefaultIntegerRangePad  uint32 = 0

	DefaultRandomStringSize    uint32 = 10
	DefaultRandomStringChars   string = "alphanum"
	DefaultRandomStringUniques uint32 = 0

	MaxRandomStringSize    = 100
	MaxRandomStringUniques = 100000

	// The default time format to use if no format is provided
	DefaultTimeFormat    = time.RFC3339Nano
	DefaultTimeTimestamp = ""

	DefaultUuidType    = "uuid4"
	DefaultUuidUniques = 0
)

var (
	RandomStringAlphaChars   = []rune("abcdefghijklmnopqrstuvwxyz")
	RandomStringNumericChars = []rune("1234567890")
	RandomStringHexChars     = []rune("0123456789abcdef")
)
