package resources

import (
	"testing"

	"ctx.sh/genie/pkg/resources/float_range"
	"ctx.sh/genie/pkg/resources/integer_range"
	"ctx.sh/genie/pkg/resources/list"
	"ctx.sh/genie/pkg/resources/random_string"
	"ctx.sh/genie/pkg/resources/timestamp"
	"ctx.sh/genie/pkg/resources/uuid"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestConfigLoad(t *testing.T) {
	in := `
integer_ranges:
  range1:
    min: 1
    max: 5
    pad: 3
float_ranges:
  range1:
    min: 1.0
    max: 5.0
    distribution: exp
    rate: 1.0
lists:
  numbers:
    - one
    - two
    - three
random_strings:
  hax:
    chars: hex
timestamps:
  now:
    format: rfc3339nano
uuids:
  id:
    type: uuid4
`
	cfg := &Config{}
	err := yaml.Unmarshal([]byte(in), cfg)
	assert.Nil(t, err)

	expected := &Config{
		IntegerRanges: map[string]integer_range.Config{
			"range1": {
				Min:          1,
				Max:          5,
				Pad:          3,
				Distribution: integer_range.DefaultIntegerRangeDistribution,
				StdDev:       &[]float64{0.5}[0],
				Mean:         &[]int64{2}[0],
			},
		},
		FloatRanges: map[string]float_range.Config{
			"range1": {
				Min:          1.0,
				Max:          5.0,
				Distribution: "exp",
				Rate:         1.0,
				StdDev:       &[]float64{0.5}[0],
				Mean:         &[]float64{2.0}[0],
				Format:       float_range.DefaultFloatRangeFormat,
				Precision:    float_range.DefaultFloatRangePrecision,
			},
		},
		Lists: map[string]list.Config{
			"numbers": []string{"one", "two", "three"},
		},
		RandomStrings: map[string]random_string.Config{
			"hax": {
				Size:    random_string.DefaultRandomStringSize,
				Chars:   random_string.RandomStringHexChars,
				Uniques: random_string.DefaultRandomStringUniques,
			},
		},
		Timestamps: map[string]timestamp.Config{
			"now": {
				Format:    "rfc3339nano",
				Timestamp: "",
			},
		},
		UUIDs: map[string]uuid.Config{
			"id": {
				Type:    "uuid4",
				Uniques: uuid.DefaultUUIDUniques,
			},
		},
	}

	assert.Equal(t, expected, cfg)

}

func TestResourcesEnsureMergeOverwrite(t *testing.T) {
	in := `
lists:
  numbers:
    - one
    - two
    - three
  numbers:
    - two
    - three
    - four
`
	cfg := &Config{}
	err := yaml.Unmarshal([]byte(in), cfg)
	assert.Nil(t, err)

	expected := &Config{
		IntegerRanges: nil,
		Lists: map[string]list.Config{
			"numbers": []string{"two", "three", "four"},
		},
		RandomStrings: nil,
		Timestamps:    nil,
		UUIDs:         nil,
	}

	assert.Equal(t, expected, cfg)

}

func TestDefaultedFloatRange(t *testing.T) {
	in := "{}"
	cfg := &float_range.Config{}
	err := yaml.Unmarshal([]byte(in), cfg)
	assert.Nil(t, err)

	assert.Equal(t, float_range.DefaultFloatRangeMax, cfg.Max)
	assert.Equal(t, float_range.DefaultFloatRangeMin, cfg.Min)
}

func TestFloatRange(t *testing.T) {
	tests := []struct {
		input string
		valid bool
	}{
		// min is defaulted to 0
		{"max: -1.0", false},
		{"max: 0.0", false},
		{"max: 1.0", true},
		{"max: 1000.0", true},
		{"max: 1000.0\nmin: 1000.0", false},
		{"max: 1000.0\nmin: 1001.0", false},
		{"max: 1.7e+308", true},
		{"max: 1.8e+308", false},
		{"min: -1.0\nmax: 1.0", true},
	}

	for _, tt := range tests {
		cfg := &float_range.Config{}
		err := yaml.Unmarshal([]byte(tt.input), cfg)
		if tt.valid {
			assert.Nil(t, err, tt.input)
		} else {
			assert.NotNil(t, err, tt.input)
		}
	}
}

func TestDefaultedIntegerRange(t *testing.T) {
	in := "{}"
	cfg := &integer_range.Config{}
	err := yaml.Unmarshal([]byte(in), cfg)
	assert.Nil(t, err)

	assert.Equal(t, integer_range.DefaultIntegerRangeMax, cfg.Max)
	assert.Equal(t, integer_range.DefaultIntegerRangeMin, cfg.Min)
	assert.Equal(t, integer_range.DefaultIntegerRangePad, cfg.Pad)
}

func TestIntegerRange(t *testing.T) {
	tests := []struct {
		input string
		valid bool
	}{
		// min is defaulted to 0
		{"max: -1", false},
		{"max: 0", false},
		{"max: 1", true},
		{"max: 1000", true},
		{"max: 1000\nmin: 1000", false},
		{"max: 1000\nmin: 1001", false},
		{"max: 9223372036854775807", true},
		{"max: 9223372036854775808", false},
		// {"min: -9223372036854775809", false},
		{"min: 0", true},
		{"min: 1", true},
		// 10 is the defaulted value for max
		{"min: 10", false},
	}

	for _, tt := range tests {
		cfg := &integer_range.Config{}
		err := yaml.Unmarshal([]byte(tt.input), cfg)
		if tt.valid {
			assert.Nil(t, err, tt.input)
		} else {
			assert.NotNil(t, err, tt.input)
		}
	}
}

func TestList(t *testing.T) {
	tests := []struct {
		input string
		items int
		valid bool
	}{
		{"- item1\n- item2", 2, true},
	}

	for _, tt := range tests {
		var cfg list.Config
		err := yaml.Unmarshal([]byte(tt.input), &cfg)
		if tt.valid {
			assert.Nil(t, err, tt.input)
			assert.Len(t, cfg, tt.items)
		} else {
			assert.NotNil(t, err, tt.input)
		}
	}
}

func TestRandomStringSizeUnique(t *testing.T) {
	tests := []struct {
		input string
		valid bool
	}{
		{"size: -1", false},
		{"size: 0", false},
		{"size: 1", true},
		{"size: 100", true},
		{"size: 101", false},
		{"uniques: 0", true},
		{"uniques: 10", true},
		{"uniques: 100001", false},
	}

	for _, tt := range tests {
		cfg := &random_string.Config{}
		err := yaml.Unmarshal([]byte(tt.input), cfg)
		if tt.valid {
			assert.Nil(t, err, tt.input)
		} else {
			assert.NotNil(t, err, tt.input)
		}
	}
}

func TestRandomStringChars(t *testing.T) {
	tests := []struct {
		input string
		chars []rune
	}{
		{"chars: alpha", random_string.RandomStringAlphaChars},
		{"chars: numeric", random_string.RandomStringNumericChars},
		{"chars: hex", random_string.RandomStringHexChars},
		{"chars: alphanum", append(random_string.RandomStringAlphaChars, random_string.RandomStringNumericChars...)},
		{"chars: abcdef!@#$12345", []rune("abcdef!@#$12345")},
	}

	for _, tt := range tests {
		cfg := &random_string.Config{}
		err := yaml.Unmarshal([]byte(tt.input), cfg)
		assert.Nil(t, err, tt.input)
		assert.EqualValues(t, tt.chars, cfg.Chars)
	}
}

func TestUUID(t *testing.T) {
	tests := []struct {
		input string
		valid bool
	}{
		{"type: uuid1", true},
		{"type: uuid4", true},
		{"type: uuid", false},
	}

	for _, tt := range tests {
		cfg := &uuid.Config{}
		err := yaml.Unmarshal([]byte(tt.input), cfg)
		if tt.valid {
			assert.Nil(t, err, tt.input)
		} else {
			assert.NotNil(t, err, tt.input)
		}
	}
}
