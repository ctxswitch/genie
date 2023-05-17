package template

import (
	"strings"

	"ctx.sh/dynamo/pkg/resources"
)

type Template struct {
	root      Root
	resources *resources.Resources
	vars      map[string]string
}

func New() *Template {
	return &Template{
		vars: make(map[string]string),
	}
}

func (t *Template) Compile(input string) error {
	p := NewParser(input).WithResources(t.resources)
	root, err := p.Parse()
	t.root = root

	return err
}

func (t *Template) Execute() string {
	return t.eval(t.root)
}

func (t *Template) WithResources(r *resources.Resources) *Template {
	t.resources = r
	return t
}

func (t *Template) WithVars(vars map[string]string) *Template {
	t.vars = vars
	return t
}

func (t *Template) eval(root Root) string {
	var out strings.Builder

	for _, node := range root.Nodes {
		switch n := node.(type) {
		case TextNode:
			out.WriteString(t.evalText(n))
		case IdentifierExpressionNode:
			out.WriteString(t.evalIdentifierExpression(n))
		}
	}

	return out.String()
}

func (t *Template) evalText(n TextNode) string {
	return n.String()
}

func (t *Template) evalIdentifierExpression(n IdentifierExpressionNode) string {
	if value, ok := t.vars[n.TokenLiteral()]; ok {
		return value
	}

	// should I send an error, warning, or just an empty string?
	return ""
}
