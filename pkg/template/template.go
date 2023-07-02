package template

import (
	"strings"

	"ctx.sh/genie/pkg/resources"
)

type Template struct {
	root Root
	// This could end up being any in the future
	vars      map[string]string
	resources *resources.Resources
}

func NewTemplate() *Template {
	return &Template{}
}

func (t *Template) Compile(input string) error {
	parser := NewParser(input, t.resources)
	root, err := parser.Parse()
	if err != nil {
		return err
	}

	t.root = root
	return nil
}

func (t *Template) WithResources(r *resources.Resources) *Template {
	t.resources = r
	return t
}

func (t *Template) WithVars(vars map[string]string) *Template {
	t.vars = vars
	return t
}

// Will we have any errors on execute?
func (t *Template) Execute() string {
	return t.eval(t.root)
}

func (t *Template) eval(root Root) string {
	var out strings.Builder

	for _, node := range root.Nodes {
		switch n := node.(type) {
		case *Text:
			out.WriteString(n.String())
		case *Comment:
			// Do nothing right now, but I'm thinking that I want to potentially
			// use those for log points.
		case *Expression:
			out.WriteString(n.WithVars(t.vars).String())
		default:
			out.WriteString(n.String())
		}
	}

	return out.String()
}
