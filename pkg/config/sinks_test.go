package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestSinksHttp(t *testing.T) {
	tests := []struct {
		input string
		valid bool
	}{
		{"url: http//localhost", true},
		{"url: http://localhost:3000", true},
	}

	for _, tt := range tests {
		cfg := &Http{}
		err := yaml.Unmarshal([]byte(tt.input), cfg)
		if tt.valid {
			assert.Nil(t, err, tt.input)
		} else {
			assert.NotNil(t, err, tt.input)
		}
	}
}
