package template

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"ctx.sh/genie/pkg/resources"
	"ctx.sh/genie/pkg/variables"
)

type Template struct {
	root  Root
	paths []string
	// This could end up being any in the future
	vars      *variables.ScopedVariables
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

func (t *Template) CompileFrom(file string) error {
	// Option 1) If the file is relative, tack it on to the paths (which should be absolute?)
	// Option 2) If the file is absolute, just use that.
	var data []byte
	var err error
	if filepath.IsAbs(file) {
		data, err = os.ReadFile(file)
		if err != nil {
			return err
		} else {
			return t.Compile(string(data))
		}
	} else {
		for _, path := range t.paths {
			file := fmt.Sprintf("%s/%s", path, file)
			_, err := os.Stat(file)
			if err == nil {
				data, err = os.ReadFile(file)
				if err != nil {
					return err
				} else {
					return t.Compile(string(data))
				}
			}
		}
	}

	return fmt.Errorf("template does not exist in search path: %s", file)
}

func (t *Template) WithPaths(p []string) *Template {
	t.paths = p
	return t
}

func (t *Template) WithResources(r *resources.Resources) *Template {
	t.resources = r
	return t
}

func (t *Template) WithVariables(vars *variables.Variables) *Template {
	t.vars = variables.NewScopedVariables(vars)
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
			s := n.WithVars(t.vars).String()
			if n.Filter != nil {
				s = n.Filter(s)
			}
			out.WriteString(s)
		case *LetStatement:
			exp := n.Expression.(*Expression)
			e := exp.WithVars(t.vars).String()
			if exp.Filter != nil {
				e = exp.Filter(e)
			}
			t.vars.Set(n.Identifier, e)
		default:
			out.WriteString(n.String())
		}
	}

	return out.String()
}
