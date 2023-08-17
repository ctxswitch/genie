package config

// func TestResources(t *testing.T) {
// 	in := `
// integer_ranges:
//   range1:
//     min: 1
//     max: 5
//     pad: 3
// lists:
//   numbers:
//     - one
//     - two
//     - three
// random_strings:
//   hax:
//     chars: hex
// timestamps:
//   now:
//     format: rfc3339nano
// uuids:
//   id:
//     type: uuid4
// `
// 	cfg := &Resources{}
// 	err := yaml.Unmarshal([]byte(in), cfg)
// 	assert.Nil(t, err)

// 	expected := &Resources{
// 		IntegerRanges: map[string]IntegerRange{
// 			"range1": {
// 				Min:  1,
// 				Max:  5,
// 				Pad:  3,
// 				Step: DefaultIntegerRangeStep,
// 			},
// 		},
// 		Lists: map[string]List{
// 			"numbers": []string{"one", "two", "three"},
// 		},
// 		RandomStrings: map[string]RandomString{
// 			"hax": {
// 				Size:    DefaultRandomStringSize,
// 				Chars:   RandomStringHexChars,
// 				Uniques: DefaultRandomStringUniques,
// 			},
// 		},
// 		Timestamps: map[string]Timestamp{
// 			"now": {
// 				Format:    "rfc3339nano",
// 				Timestamp: "",
// 			},
// 		},
// 		Uuids: map[string]Uuid{
// 			"id": {
// 				Type:    "uuid4",
// 				Uniques: DefaultUuidUniques,
// 			},
// 		},
// 	}

// 	assert.Equal(t, expected, cfg)

// }

// func TestResourcesEnsureMergeOverwrite(t *testing.T) {
// 	in := `
// lists:
//   numbers:
//     - one
//     - two
//     - three
//   numbers:
//     - two
//     - three
//     - four
// `
// 	cfg := &Resources{}
// 	err := yaml.Unmarshal([]byte(in), cfg)
// 	assert.Nil(t, err)

// 	expected := &Resources{
// 		IntegerRanges: nil,
// 		Lists: map[string]List{
// 			"numbers": []string{"two", "three", "four"},
// 		},
// 		RandomStrings: nil,
// 		Timestamps:    nil,
// 		Uuids:         nil,
// 	}

// 	assert.Equal(t, expected, cfg)

// }

// func TestDefaultedIntegerRange(t *testing.T) {
// 	in := "min: 0"
// 	cfg := &IntegerRange{}
// 	err := yaml.Unmarshal([]byte(in), cfg)
// 	assert.Nil(t, err)

// 	assert.Equal(t, DefaultIntegerRangeMax, cfg.Max)
// 	assert.Equal(t, DefaultIntegerRangeMin, cfg.Min)
// 	assert.Equal(t, DefaultIntegerRangePad, cfg.Pad)
// 	assert.Equal(t, DefaultIntegerRangeStep, cfg.Step)
// }

// func TestIntegerRange(t *testing.T) {
// 	tests := []struct {
// 		input string
// 		valid bool
// 	}{
// 		{"max: -1", false},
// 		{"max: 0", false},
// 		{"max: 1", true},
// 		{"max: 1000", true},
// 		{"max: 1000\nmin: 1000", false},
// 		{"max: 1000\nmin: 1001", false},
// 		{"max: 9223372036854775807", true},
// 		{"max: 9223372036854775808", false},
// 		// {"min: -9223372036854775809", false},
// 		{"min: 0", true},
// 		{"min: 1", true},
// 		// 10 is the defaulted value for max
// 		{"min: 10", false},
// 	}

// 	for _, tt := range tests {
// 		cfg := &IntegerRange{}
// 		err := yaml.Unmarshal([]byte(tt.input), cfg)
// 		if tt.valid {
// 			assert.Nil(t, err, tt.input)
// 		} else {
// 			assert.NotNil(t, err, tt.input)
// 		}
// 	}
// }

// func TestList(t *testing.T) {
// 	tests := []struct {
// 		input string
// 		items int
// 		valid bool
// 	}{
// 		{"- item1\n- item2", 2, true},
// 	}

// 	for _, tt := range tests {
// 		var cfg List
// 		err := yaml.Unmarshal([]byte(tt.input), &cfg)
// 		if tt.valid {
// 			assert.Nil(t, err, tt.input)
// 			assert.Len(t, cfg, tt.items)
// 		} else {
// 			assert.NotNil(t, err, tt.input)
// 		}
// 	}
// }

// func TestRandomStringSizeUnique(t *testing.T) {
// 	tests := []struct {
// 		input string
// 		valid bool
// 	}{
// 		{"size: -1", false},
// 		{"size: 0", false},
// 		{"size: 1", true},
// 		{"size: 100", true},
// 		{"size: 101", false},
// 		{"uniques: 0", true},
// 		{"uniques: 10", true},
// 		{"uniques: 100001", false},
// 	}

// 	for _, tt := range tests {
// 		cfg := &RandomString{}
// 		err := yaml.Unmarshal([]byte(tt.input), cfg)
// 		if tt.valid {
// 			assert.Nil(t, err, tt.input)
// 		} else {
// 			assert.NotNil(t, err, tt.input)
// 		}
// 	}
// }

// func TestRandomStringChars(t *testing.T) {
// 	tests := []struct {
// 		input string
// 		chars []rune
// 	}{
// 		{"chars: alpha", RandomStringAlphaChars},
// 		{"chars: numeric", RandomStringNumericChars},
// 		{"chars: hex", RandomStringHexChars},
// 		{"chars: alphanum", append(RandomStringAlphaChars, RandomStringNumericChars...)},
// 		{"chars: abcdef!@#$12345", []rune("abcdef!@#$12345")},
// 	}

// 	for _, tt := range tests {
// 		cfg := &RandomString{}
// 		err := yaml.Unmarshal([]byte(tt.input), cfg)
// 		assert.Nil(t, err, tt.input)
// 		assert.EqualValues(t, tt.chars, cfg.Chars)
// 	}
// }

// func TestUuid(t *testing.T) {
// 	tests := []struct {
// 		input string
// 		valid bool
// 	}{
// 		{"type: uuid1", true},
// 		{"type: uuid4", true},
// 		{"type: uuid", false},
// 	}

// 	for _, tt := range tests {
// 		cfg := &Uuid{}
// 		err := yaml.Unmarshal([]byte(tt.input), cfg)
// 		if tt.valid {
// 			assert.Nil(t, err, tt.input)
// 		} else {
// 			assert.NotNil(t, err, tt.input)
// 		}
// 	}
// }
