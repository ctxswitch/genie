package variables

import (
	"testing"

	"ctx.sh/genie/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	vars, err := Parse([]config.VariableBlock{
		{Name: "name", Value: "Dwight Schrute"},
	})
	assert.NoError(t, err)
	assert.Equal(t, "Dwight Schrute", vars.vars["name"])
}

func TestGet(t *testing.T) {
	vars := MockVariables()
	val, ok := vars.Get("name")
	assert.True(t, ok)
	assert.Equal(t, "Dwight Schrute", val)
}

func TestSet(t *testing.T) {
	vars := MockVariables()
	err := vars.Set("name", "Jim Halpert")
	assert.NoError(t, err)
	val, ok := vars.Get("name")
	assert.True(t, ok)
	assert.Equal(t, "Jim Halpert", val)
}
