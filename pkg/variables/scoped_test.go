package variables

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewScopedVariables(t *testing.T) {
	vars := MockVariables()
	scoped := NewScopedVariables(vars)
	assert.Equal(t, 1, scoped.Len())
	assert.Equal(t, "Dwight Schrute", scoped.vars[0].vars["name"])
}

func TestScopedVariables(t *testing.T) {
	vars := MockVariables()
	scoped := NewScopedVariables(vars)
	scoped.NewScope()
	assert.Equal(t, 2, scoped.Len())
	assert.Equal(t, "Dwight Schrute", scoped.vars[0].vars["name"])
	assert.Equal(t, "Dwight Schrute", scoped.vars[1].vars["name"])

	_ = scoped.Set("name", "Jim Halpert")
	assert.Equal(t, "Dwight Schrute", scoped.vars[0].vars["name"])
	assert.Equal(t, "Jim Halpert", scoped.vars[1].vars["name"])

	name, ok := scoped.Get("name")
	assert.True(t, ok)
	assert.Equal(t, "Jim Halpert", name)

	scoped.ExitScope()
	assert.Equal(t, 1, scoped.Len())
	assert.Equal(t, "Dwight Schrute", scoped.vars[0].vars["name"])

	name, ok = scoped.Get("name")
	assert.True(t, ok)
	assert.Equal(t, "Dwight Schrute", name)
}
