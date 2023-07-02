package template

import (
	"testing"

	"ctx.sh/genie/pkg/resources"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	input = `
<# This is a comprehensive test of the templating system #>
Hello World
<< name >>
<< list.name >>
<< "Pam Beesly" >>

<< list.greeting >> World!
<% minimize %>this that and the << other >><% endminimize %>
`
	expected = `

Hello World
Dwight Schrute
Jim Halpert
Pam Beesly

Hello World!
`
)

func TestTemplateParse(t *testing.T) {
	var err error

	vars := make(map[string]string)
	vars["name"] = "Dwight Schrute"

	tmpl := NewTemplate().
		WithResources(resources.MockResources()).
		WithVars(vars)

	// Something to think about.  We can look for compile time unknown variable
	// issues by keeping track of any new variable that is set when we parse
	// and then using that to test the existence.  (or with the global variables
	// that are set)

	err = tmpl.Compile(input)
	require.NoError(t, err)

	out := tmpl.Execute()
	assert.Equal(t, expected, out)
}
