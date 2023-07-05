package config

import (
	"fmt"
	"strings"
)

type Resource struct {
	ID string
}

// Resources represents the top level config organization for the
// resources block.
type Resources struct {
	IntegerRanges map[string]IntegerRange `yaml:"integer_ranges"`
	Lists         map[string]List         `yaml:"lists"`
	RandomStrings map[string]RandomString `yaml:"random_strings"`
	Timestamps    map[string]Timestamp    `yaml:"timestamps"`
	Uuids         map[string]Uuid         `yaml:"uuids"`
}

// IntegerRange is a resource that represents a sequential list of numbers.  This
// range has a minimum and maximum.  It also can be controlled by defining the step,
// and leading zero padding.
type IntegerRange struct {
	Min  int64  `yaml:"min"`
	Max  int64  `yaml:"max"`
	Step int64  `yaml:"step"`
	Pad  uint32 `yaml:"pad"`
}

func (i *IntegerRange) validate() (bool, error) {
	if i.Max < 1 || i.Max <= i.Min {
		return false, fmt.Errorf("max (%d) in integer_range must be greater than zero and the minimum value", i.Max)
	}

	return true, nil
}

// UnmarshalYAML sets the IntegerRange defaults and parses an integer_range
// / block.
func (i *IntegerRange) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type IntegerRangeDefaulted IntegerRange
	var defaults = IntegerRangeDefaulted{
		Min:  DefaultIntegerRangeMin,
		Max:  DefaultIntegerRangeMax,
		Step: DefaultIntegerRangeStep,
		Pad:  DefaultIntegerRangePad,
	}

	out := defaults
	if err := unmarshal(&out); err != nil {
		return err
	}

	tmpl := IntegerRange(out)
	if valid, err := tmpl.validate(); !valid {
		return err
	}

	*i = tmpl
	return nil
}

// List represents a resource that includes user defined lists of strings
type List []string

func (l *List) validate() (bool, error) {
	if len(*l) == 0 {
		return false, fmt.Errorf("items in list cannot be empty")
	}

	return true, nil
}

// UnmarshalYAML parses a list block.  Though it currently does not
// have any defaults, the structure is there for future work.
func (l *List) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type ListDefaulted List
	var defaults = ListDefaulted{}

	out := defaults
	if err := unmarshal(&out); err != nil {
		return err
	}

	tmpl := List(out)
	if valid, err := tmpl.validate(); !valid {
		return err
	}

	*l = tmpl
	return nil
}

// Random strings is a resource that represents a random strings of any
// given size.
type RandomString struct {
	Size    uint32 `yaml:"size"`
	Chars   []rune `yaml:"chars"`
	Uniques uint32 `yaml:"uniques"`
}

// Validate checks to ensure that attributes are valid.
// TODO: Actually validate something.
func (r *RandomString) validate() (bool, error) {
	if r.Size < 1 || r.Size > MaxRandomStringSize {
		return false, fmt.Errorf("size (%d) in random_string must be greater than zero and less than or equal to %d",
			r.Size,
			MaxRandomStringSize,
		)
	}

	if r.Uniques > MaxRandomStringUniques {
		return false, fmt.Errorf("uniques (%d) in random_string must be less than or equal to %d",
			r.Uniques,
			MaxRandomStringUniques,
		)
	}
	return true, nil
}

// convertChars converts strings from yaml into runes for either the predefined sets
// of characters or any character in a string if it does not match anything.
func (r *RandomString) convertChars(charstr string) []rune {
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

// UnmarshalYAML sets the RandomString defaults and parses an random_string
// / block.
func (r *RandomString) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// OMG this is a pretty hacky way to do on-the-fly string to rune conversions
	// but it appears to work great and doesn't add a ton of overhead.
	type RandomStringTmp struct {
		Size    uint32
		Chars   string
		Uniques uint32
	}
	var defaults = RandomStringTmp{
		Size:    DefaultRandomStringSize,
		Chars:   DefaultRandomStringChars,
		Uniques: DefaultRandomStringUniques,
	}
	out := defaults

	if err := unmarshal(&out); err != nil {
		return err
	}

	tmpl := RandomString{
		Size:    out.Size,
		Chars:   r.convertChars(out.Chars),
		Uniques: out.Uniques,
	}

	if valid, err := tmpl.validate(); !valid {
		return err
	}

	*r = tmpl
	return nil
}

// Time defines the attributes for a time configuration element.
type Timestamp struct {
	// Format is a textual representation of the time value formatted
	// according to the layout defined by the argument.  It follows
	// the go representation for time formats.  There are several mapped
	// values that will convert to the go equivalent including: rfc3339,
	// rfc3999nano, unix, and unixnano.
	Format    string `yaml:"format,omitempty"`
	Timestamp string `yaml:"timestamp,omitempty"`
}

// Validate returns true or false if the provided attributes are valid.
func (t *Timestamp) Validate() (bool, error) {
	return true, nil
}

// UnmarshalYAML provides a custom interface to set defaults when parsing
// yaml bytes.
func (t *Timestamp) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type TimestampDefaulted Timestamp
	var defaults = TimestampDefaulted{
		Format:    DefaultTimeFormat,
		Timestamp: DefaultTimeTimestamp,
	}

	out := defaults
	if err := unmarshal(&out); err != nil {
		return err
	}

	tmpl := Timestamp(out)
	if valid, err := tmpl.Validate(); !valid {
		return err
	}

	*t = tmpl
	return nil
}

// Produce a random value according to some format.
type Uuid struct {
	Type    string `yaml:"type"`
	Uniques int    `yaml:"uniques"`
}

// Validate checks to ensure that our format is meaningful.
func (u *Uuid) validate() (bool, error) {
	if !(u.Type == "uuid1" || u.Type == "uuid4") {
		return false, fmt.Errorf("unsupported UUID type %s for uuid", u.Type)
	}
	return true, nil
}

// UnmarshalYAML provides a custom interface when parsing UUID resource blocks
func (u *Uuid) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type UuidDefaulted Uuid
	var defaults = UuidDefaulted{
		Type:    DefaultUuidType,
		Uniques: DefaultUuidUniques,
	}

	out := defaults
	if err := unmarshal(&out); err != nil {
		return err
	}

	tmpl := Uuid(out)

	// Ensure we have a lowercase type
	tmpl.Type = strings.ToLower(tmpl.Type)

	if valid, err := tmpl.validate(); !valid {
		return err
	}

	*u = tmpl
	return nil
}
