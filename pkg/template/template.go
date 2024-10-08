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
}

func NewTemplate() *Template {
	return &Template{}
}

func (t *Template) WithPaths(paths []string) *Template {
	t.paths = paths
	return t
}

func (t *Template) Compile(input string) error {
	parser := NewParser(input)
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
	if filepath.IsAbs(file) { // nolint:nestif
		data, err = os.ReadFile(file)
		if err != nil {
			return err
		}
		return t.Compile(string(data))
	} else {
		for _, path := range t.paths {
			f := fmt.Sprintf("%s/%s", path, file)
			_, err := os.Stat(f)
			if err == nil {
				data, err = os.ReadFile(f)
				if err != nil {
					return err
				}
				return t.Compile(string(data))
			}
		}
	}

	return fmt.Errorf("template does not exist in search path: %s", file)
}

// Will we have any errors on execute?  Can I pass resources and variables here?
// It would decouple them from the complile stage which would mean that we could
// compile before we have the resources and variables (variables would be pulled
// in with events).  It's going to be required if we allow our templates to be used
// in other configurations (i.e. sinks).
func (t *Template) Execute(res *resources.Resources, vars *variables.Variables) string {
	// TODO: concurrent map writes on the variables here...
	return t.eval(t.root, res, variables.NewScopedVariables(vars))
}

func (t *Template) eval(root Root, res *resources.Resources, vars *variables.ScopedVariables) string {
	var out strings.Builder

	for _, node := range root.Nodes {
		switch n := node.(type) {
		case *Text:
			out.WriteString(n.String())
		case *Comment:
			// Do nothing right now, but I'm thinking that I want to potentially
			// use those for log points.
		case *Expression:
			s := n.WithVariables(vars).WithResources(res).String()
			out.WriteString(s)
		case *LetStatement:
			exp := n.Expression.(*Expression).WithVariables(n.Vars).WithResources(res)
			e := exp.String()
			// TODO: handle errors (or at least log them)
			_ = vars.Set(n.Identifier, e)
		default:
			out.WriteString(n.String())
		}
	}

	return strings.TrimSpace(out.String())
}
