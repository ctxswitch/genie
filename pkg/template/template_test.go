package template

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	tests := []struct {
		input    string
		vars     map[string]string
		expected string
	}{
		{`hello`, map[string]string{}, "hello"},
		{`{{ greeting }}`, map[string]string{"greeting": "hello"}, "hello"},
		{`{% let greeting = "hello" %}{{ greeting }}`, map[string]string{"greeting": "hello"}, "hello"},
		// {"{{ list:greetings }}", map[string]string{}, "hello"},
	}

	for i, tt := range tests {
		tmpl := New().WithVars(tt.vars)
		err := tmpl.Compile(tt.input)
		assert.NoError(t, err, "test[%d]: %s", i, err)
	}
}
