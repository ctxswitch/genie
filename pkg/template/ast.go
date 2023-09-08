package template

import (
	"ctx.sh/genie/pkg/filter"
	"ctx.sh/genie/pkg/resources"
	"ctx.sh/genie/pkg/variables"
)

type Node interface {
	String() string
}

type Root struct {
	Nodes []Node
}

func NewRoot() Root {
	return Root{
		Nodes: make([]Node, 0),
	}
}

func (n *Root) Length() int {
	return len(n.Nodes)
}

type TextNode interface {
	Node
	TokenText()
}

type StatementNode interface {
	Node
	StatementNode()
}

type ExpressionNode interface {
	Node
	ExpressionNode()
	WithFilter(filter.FilterFunc) ExpressionNode
	WithVariables(*variables.ScopedVariables) ExpressionNode
	WithResources(*resources.Resources) ExpressionNode
}

type Control struct{}

func (n *Control) String() string { return "" }

type Text struct {
	Token Token
}

func (n *Text) String() string { return n.Token.Literal }

// Think about coming back through here and removing the token requirements.
// We don't actually need them for anything as the tokens are useless to us
// other than looking at the type of the expression - and in that case we
// just need to store the type.
type Expression struct {
	Token  Token
	Name   string
	Filter filter.FilterFunc

	vars *variables.ScopedVariables
	res  *resources.Resources
}

func (n *Expression) String() string {
	var out string
	var ok bool
	switch n.Token.Type {
	case TokenIdentifier:
		if out, ok = n.vars.Get(n.Token.Literal); !ok {
			out = n.Token.Literal
		}
	case TokenResource:
		res, err := n.res.Get(n.Token.Literal, n.Name)
		if err != nil {
			return ""
		}

		out = res.Get()
	default:
		out = n.Token.Literal
	}

	if n.Filter != nil {
		out = n.Filter(out)
	}

	return out
}

func (n *Expression) WithFilter(fn filter.FilterFunc) ExpressionNode {
	n.Filter = fn
	return n
}
func (n *Expression) WithVariables(vars *variables.ScopedVariables) ExpressionNode {
	n.vars = vars
	return n
}
func (n *Expression) WithResources(res *resources.Resources) ExpressionNode {
	n.res = res
	return n
}
func (n *Expression) ExpressionNode() {}

// Probably can consolidate this like we did with expression
type LetStatement struct {
	Token      Token
	Identifier string
	Expression ExpressionNode
	Vars       *variables.ScopedVariables
}

func (n *LetStatement) String() string { return n.Token.Literal }
func (n *LetStatement) StatementNode() {}

// I'll convert this to allow for comments as logs later.
type Comment struct {
	Token Token
}

func (n *Comment) String() string { return n.Token.Literal }
