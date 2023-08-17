package config

import (
	"strings"
)

type ResourcesBlock struct {
	IntegerRanges map[string]IntegerRangeBlock `yaml:"integer_ranges"`
	Lists         map[string]ListBlock         `yaml:"lists"`
	RandomStrings map[string]RandomStringBlock `yaml:"random_strings"`
	Timestamps    map[string]TimestampBlock    `yaml:"timestamps"`
	Uuids         map[string]UuidBlock         `yaml:"uuids"`
}

type IntegerRangeBlock struct {
	Min  int64  `yaml:"min"`
	Max  int64  `yaml:"max"`
	Step int64  `yaml:"step"`
	Pad  uint32 `yaml:"pad"`
}

func (i *IntegerRangeBlock) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type IntegerRangeBlockDefaulted IntegerRangeBlock
	var defaults = IntegerRangeBlockDefaulted{
		Min:  DefaultIntegerRangeMin,
		Max:  DefaultIntegerRangeMax,
		Step: DefaultIntegerRangeStep,
		Pad:  DefaultIntegerRangePad,
	}

	out := defaults
	if err := unmarshal(&out); err != nil {
		return err
	}

	tmpl := IntegerRangeBlock(out)

	*i = tmpl
	return nil
}

type ListBlock []string

type RandomStringBlock struct {
	Size    uint32 `yaml:"size"`
	Chars   []rune `yaml:"chars"`
	Uniques uint32 `yaml:"uniques"`
}

func (r *RandomStringBlock) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type RandomStringBlockDefaults struct {
		Size uint32 `yaml:"size"`
		// The config uses a string to signify a char group id or string
		// of characters to use.  We'll convert it to a rune slice later.
		Chars   string `yaml:"chars"`
		Uniques uint32 `yaml:"uniques"`
	}

	var defaults = RandomStringBlockDefaults{
		Size:    DefaultRandomStringSize,
		Chars:   DefaultRandomStringChars,
		Uniques: DefaultRandomStringUniques,
	}
	out := defaults

	if err := unmarshal(&out); err != nil {
		return err
	}

	tmpl := RandomStringBlock{
		Size:    out.Size,
		Chars:   convertChars(out.Chars),
		Uniques: out.Uniques,
	}

	*r = tmpl
	return nil
}

type TimestampBlock struct {
	Format    string `yaml:"format"`
	Timestamp string `yaml:"timestamp"`
}

func (t *TimestampBlock) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type TimestampBlockDefaulted TimestampBlock
	var defaults = TimestampBlockDefaulted{
		Format:    DefaultTimeFormat,
		Timestamp: DefaultTimeTimestamp,
	}

	out := defaults
	if err := unmarshal(&out); err != nil {
		return err
	}

	tmpl := TimestampBlock(out)

	*t = tmpl
	return nil
}

type UuidBlock struct {
	Type    string `yaml:"type"`
	Uniques int    `yaml:"uniques"`
}

func (u *UuidBlock) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type UuidBlockDefaulted UuidBlock
	var defaults = UuidBlockDefaulted{
		Type:    DefaultUuidType,
		Uniques: DefaultUuidUniques,
	}

	out := defaults
	if err := unmarshal(&out); err != nil {
		return err
	}

	tmpl := UuidBlock(out)

	// Ensure we have a lowercase type
	tmpl.Type = strings.ToLower(tmpl.Type)

	*u = tmpl
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
