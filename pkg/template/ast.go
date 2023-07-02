package template

import (
	"ctx.sh/genie/pkg/filter"
	"ctx.sh/genie/pkg/resources"
)

type Node interface {
	Literal() string
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
	WithVars(map[string]string) ExpressionNode
}

type ExpressionNode interface {
	Node
	ExpressionNode()
	Filter(string) string
	WithVars(map[string]string) ExpressionNode
}

type Control struct{}

func (n *Control) Literal() string { return "-" }
func (n *Control) String() string  { return "" }

type Text struct {
	Token Token
}

func (n *Text) Literal() string { return n.Token.Literal }
func (n *Text) String() string  { return n.Token.Literal }

// Think about coming back through here and removing the token requirements.
// We don't actually need them for anything as the tokens are useless to us
// other than looking at the type of the expression - and in that case we
// just need to store the type.
type Expression struct {
	Token    Token
	Name     string
	Resource resources.Resource
	Vars     map[string]string // This changes to any once we introduce
	// more types.
	Filter filter.FilterFunc
}

func (n *Expression) Literal() string { return n.Token.Literal }
func (n *Expression) String() string {
	var out string
	var ok bool
	switch n.Token.Type {
	case TokenResource:
		out = n.Resource.Get()
	case TokenIdentifier:
		if out, ok = n.Vars[n.Token.Literal]; !ok {
			// This should change once we start checking for existence
			// when parsing. i.e. we want to make sure it exists...
			out = n.Token.Literal
		}

	default:
		out = n.Token.Literal
	}

	return out
}
func (n *Expression) WithVars(vars map[string]string) *Expression { n.Vars = vars; return n }
func (n *Expression) ExpressionNode()                             {}

// Probably can consolidate this like we did with expression
type LetStatement struct {
	Token      Token
	Expression ExpressionNode
	Vars       map[string]string
}

func (n *LetStatement) Literal() string                               { return n.Token.Literal }
func (n *LetStatement) String() string                                { return n.Token.Literal }
func (n *LetStatement) WithVars(vars map[string]string) *LetStatement { n.Vars = vars; return n }
func (n *LetStatement) StatementNode()                                {}

// I'll convert this to allow for comments as logs later.
type Comment struct {
	Token Token
}

func (n *Comment) Literal() string { return n.Token.Literal }
func (n *Comment) String() string  { return n.Token.Literal }
