package random_string // nolint:revive

import "fmt"

type Config struct {
	Size    uint32 `yaml:"size"`
	Chars   []rune `yaml:"chars"`
	Uniques uint32 `yaml:"uniques"`
}

func (r *Config) validate() error {
	if r.Size < 1 || r.Size > MaxRandomStringSize {
		return fmt.Errorf("size (%d) in random_string must be greater than zero and less than or equal to %d",
			r.Size,
			MaxRandomStringSize,
		)
	}

	if r.Uniques > MaxRandomStringUniques {
		return fmt.Errorf("uniques (%d) in random_string must be less than or equal to %d",
			r.Uniques,
			MaxRandomStringUniques,
		)
	}
	return nil
}

func (r *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type ConfigDefaults struct {
		Size uint32 `yaml:"size"`
		// The config uses a string to signify a char group id or string
		// of characters to use.  We'll convert it to a rune slice later.
		Chars   string `yaml:"chars"`
		Uniques uint32 `yaml:"uniques"`
	}

	var defaults = ConfigDefaults{
		Size:    DefaultRandomStringSize,
		Chars:   DefaultRandomStringChars,
		Uniques: DefaultRandomStringUniques,
	}
	out := defaults

	if err := unmarshal(&out); err != nil {
		return err
	}

	rs := Config{
		Size:    out.Size,
		Chars:   convertChars(out.Chars),
		Uniques: out.Uniques,
	}
	if err := rs.validate(); err != nil {
		return err
	}

	*r = rs
	return nil
}

// convertChars converts strings from yaml into runes for either the predefined sets
// of characters or any character in a string if it does not match anything.
func convertChars(charstr string) []rune {
	switch charstr {
	case "alpha":
		return RandomStringAlphaChars
	case "numeric":
		return RandomStringNumericChars
	case "hex":
		return RandomStringHexChars
	case "alphanum":
		return append(RandomStringAlphaChars, RandomStringNumericChars...)
	default:
		return []rune(charstr)
	}
}
