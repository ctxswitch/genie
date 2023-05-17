package template

import (
	"strings"

	"ctx.sh/dynamo/pkg/resources"
	"ctx.sh/dynamo/pkg/template/token"
)

// TODO: errors should be more descriptive, adding syntax information.

type Parser struct {
	l *Lexer
	e []error

	res  *resources.Resources
	curr token.Token
	peek token.Token
}

func NewParser(input string) *Parser {
	p := &Parser{
		l: NewLexer(input),
		e: make([]error, 0),
	}

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) WithResources(r *resources.Resources) *Parser {
	p.res = r
	return p
}

func (p *Parser) Parse() (Root, error) {
	root := Root{
		Nodes: make([]any, 0),
	}

	for p.curr.Type != token.EOF {
		var node Node
		switch p.curr.Type {
		case token.Text:
			node = p.parseText()
		case token.ExpressionStart:
			node = p.parseExpression()
			if err := p.nextExpect(token.ExpressionEnd); err != nil {
				p.addError(UnexpectedTokenError)
				return root, err
			}
		case token.StatementStart:
			node = p.parseStatement()
			if err := p.nextExpect(token.StatementEnd); err != nil {
				p.addError(UnexpectedTokenError)
				return root, err
			}
		}

		root.Nodes = append(root.Nodes, node)
		p.nextToken()
	}

	return root, nil
}

func (p *Parser) parseText() Node {
	return &TextNode{
		Token: p.curr,
	}
}

func (p *Parser) parseExpression() Expression {
	p.nextToken()

	switch p.curr.Type {
	case token.Resource:
		return p.parseResourceExpression()
	case token.Identifier:
		return p.parseIdentifierExpression()
	default:
		p.addError(UnexpectedTokenError)
		return nil
	}
}

func (p *Parser) parseResourceExpression() Expression {
	resourceToken := p.curr

	if err := p.nextExpect(token.Period); err != nil {
		p.addError(UnexpectedTokenError)
		return nil
	}

	if err := p.nextExpect(token.Identifier); err != nil {
		p.addError(UnexpectedTokenError)
		return nil
	}

	identifierToken := p.curr

	resource, err := p.res.Get(resourceToken.Literal, p.curr.Literal)
	if err != nil {
		p.addError(err)
		return nil
	}

	node := &ResourceExpressionNode{
		Token:    identifierToken,
		Resource: resource,
	}

	return node
}

func (p *Parser) parseIdentifierExpression() Expression {
	node := &IdentifierExpressionNode{
		Token: p.curr,
	}

	return node
}

func (p *Parser) parseStatement() Node {
	p.nextToken()

	switch p.curr.Type {
	case token.Keyword:
		return p.parseKeywordStatement()
	default:
		p.addError(UnexpectedTokenError)
		return nil
	}
}

func (p *Parser) parseKeywordStatement() Node {
	switch p.curr.Literal {
	case "let":
		return p.parseLetKeywordStatement()
	default:
		p.addError(UnknownKeyword)
		return nil
	}
}

func (p *Parser) parseLetKeywordStatement() Node {
	p.nextExpect()

	if err := p.nextExpect(token.Identifier); err != nil {
		p.addError(UnexpectedTokenError)
		return nil
	}

	identifier := p.curr

	if err := p.nextExpect(token.Equals); err != nil {
		p.addError(UnexpectedTokenError)
		return nil
	}

	exp := p.parseExpression()

	node := &LetStatementNode{Token: identifier, Expression: exp}

	return node
}

func (p *Parser) nextToken() {
	p.curr = p.peek
	p.peek = p.l.NextToken()
}

func (p *Parser) nextExpect(expected ...token.Type) error {
	for _, ex := range expected {
		if ex == p.peek.Type {
			p.nextToken()
			return nil
		}
	}

	return UnexpectedTokenError
}

func (p *Parser) addError(err error) {
	p.e = append(p.e, err)
}

func (p *Parser) IsError() bool {
	return len(p.e) > 0
}

func (p *Parser) Error() string {
	var builder strings.Builder
	for _, err := range p.e {
		builder.WriteString(err.Error())
		builder.WriteByte('\n')
	}

	return builder.String()
}
