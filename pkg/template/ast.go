package template

import (
	"ctx.sh/dynamo/pkg/resources"
	"ctx.sh/dynamo/pkg/template/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Text interface {
	Node
	Text()
}

type Statement interface {
	Node
	StatementNode()
}

type Expression interface {
	Node
	ExpressionNode()
}

type Error interface {
	Node
	Error()
}

type Root struct {
	Nodes []any
}

func (r *Root) Length() int {
	return len(r.Nodes)
}

type TextNode struct {
	Token token.Token
}

func (t *TextNode) TokenLiteral() string { return t.Token.Literal }
func (t *TextNode) String() string       { return t.Token.Literal }

type IdentifierExpressionNode struct {
	Token token.Token
	// GlobalVars map[string]
}

func (i *IdentifierExpressionNode) TokenLiteral() string { return i.Token.Literal }
func (i *IdentifierExpressionNode) String() string       { return i.Token.Literal }
func (i *IdentifierExpressionNode) ExpressionNode()      {}

type ResourceExpressionNode struct {
	Token    token.Token
	Resource resources.Resource
}

func (r *ResourceExpressionNode) TokenLiteral() string { return r.Token.Literal }
func (r *ResourceExpressionNode) String() string       { return r.Resource.Get() }
func (r *ResourceExpressionNode) ExpressionNode()      {}

type LetStatementNode struct {
	Token      token.Token
	Expression Expression
}

func (l *LetStatementNode) TokenLiteral() string { return l.Token.Literal }
func (l *LetStatementNode) String() string       { return l.Token.Literal }
func (l *LetStatementNode) StatementNode()       {}
