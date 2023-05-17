package token

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLookup(t *testing.T) {
	tests := []struct {
		id       string
		expected Type
	}{
		{"integer_range", Resource},
		{"list", Resource},
		{"random_string", Resource},
		{"timestamp", Resource},
		{"uuid", Resource},
		{"let", Keyword},
		{"notit", Identifier},
	}

	for i, tt := range tests {
		got := Lookup(tt.id)
		assert.Equal(t, tt.expected, got, "tests[%d]: unexpected lookup result", i)

	}
}
