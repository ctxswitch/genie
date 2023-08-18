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
<< list.greeting | capitalize >> World!
<% let greeting = list.greeting %><< greeting >> World!
<% let greeting = list.greeting | capitalize %><< greeting >> World!
`
	expected = `
Hello World
Dwight Schrute
Jim Halpert
Pam Beesly

Hello World!
HELLO World!
Hello World!
HELLO World!
`
)

// <% minimize %>this that and the << other >><% endminimize %>
// treat these as global filters?

func TestTemplateParse(t *testing.T) {
	var err error

	// TODO: fix test
	tmpl := NewTemplate().
		WithResources(resources.MockResources()).
		WithVar("name", "Dwight Schrute")

	// Something to think about.  We can look for compile time unknown variable
	// issues by keeping track of any new variable that is set when we parse
	// and then using that to test the existence.  (or with the global variables
	// that are set)

	err = tmpl.Compile(input)
	require.NoError(t, err)

	out := tmpl.Execute()
	assert.Equal(t, expected, out)
}
