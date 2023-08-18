package config

import (
	"fmt"
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

func (i *IntegerRangeBlock) validate() error {
	// Fix me now that we allow negative values
	if i.Max <= i.Min {
		return fmt.Errorf("max (%d) in integer_range must be greater than zero and the minimum value", i.Max)
	}

	return nil
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
	if err := tmpl.validate(); err != nil {
		return err
	}

	*i = tmpl
	return nil
}

type ListBlock []string

func (l *ListBlock) validate() error {
	if len(*l) == 0 {
		return fmt.Errorf("items in list cannot be empty")
	}

	return nil
}

func (l *ListBlock) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type ListBlockDefaulted ListBlock
	var defaults = ListBlockDefaulted{}

	out := defaults
	if err := unmarshal(&out); err != nil {
		return err
	}

	list := ListBlock(out)
	if err := list.validate(); err != nil {
		return err
	}

	*l = list
	return nil
}

type RandomStringBlock struct {
	Size    uint32 `yaml:"size"`
	Chars   []rune `yaml:"chars"`
	Uniques uint32 `yaml:"uniques"`
}

func (r *RandomStringBlock) validate() error {
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

	rs := RandomStringBlock{
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

type TimestampBlock struct {
	Format    string `yaml:"format"`
	Timestamp string `yaml:"timestamp"`
}

func (t *TimestampBlock) validate() error {
	return nil
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

	ts := TimestampBlock(out)
	if err := ts.validate(); err != nil {
		return err
	}

	*t = ts
	return nil
}

type UuidBlock struct {
	Type    string `yaml:"type"`
	Uniques int    `yaml:"uniques"`
}

func (u *UuidBlock) validate() error {
	if !(u.Type == "uuid1" || u.Type == "uuid4") {
		return fmt.Errorf("unsupported UUID type %s for uuid", u.Type)
	}
	return nil
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

	uuid := UuidBlock(out)

	// Ensure we have a lowercase type
	uuid.Type = strings.ToLower(uuid.Type)
	if err := uuid.validate(); err != nil {
		return err
	}

	*u = uuid
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
