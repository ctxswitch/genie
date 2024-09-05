package template

import (
	"fmt"

	"ctx.sh/genie/pkg/filter"
)

type ParserFunc func() (Node, error)

type Parser struct {
	l       *Lexer
	parsers map[TokenType]ParserFunc
	peek    Token
	curr    Token
}

func NewParser(input string) *Parser {
	p := &Parser{
		l:       NewLexer(input),
		parsers: make(map[TokenType]ParserFunc),
	}

	p.registerParser(TokenText, p.parseText)
	p.registerParser(TokenIdentifier, p.parseExpression)
	p.registerParser(TokenString, p.parseExpression)
	p.registerParser(TokenResource, p.parseExpression)
	p.registerParser(TokenKeyword, p.parseStatement)
	p.registerParser(TokenComment, p.parseComment)

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) Parse() (Root, error) {
	root := NewRoot()

	for p.curr.Type != TokenEOF {
		fn, ok := p.parsers[p.curr.Type]
		if !ok {
			return Root{}, fmt.Errorf("unexpected token: %v", p.curr.Literal)
		}
		node, err := fn()
		if err != nil {
			return Root{}, err
		}

		// In certain cases, such as the whitespace control tokens, the node
		// is ignored and will return nil after the next node to be processed
		// is mutated.
		if node != nil {
			root.Nodes = append(root.Nodes, node)
		}
		p.nextToken()
	}

	return root, nil
}

func (p *Parser) registerParser(t TokenType, fn ParserFunc) {
	p.parsers[t] = fn
}

func (p *Parser) parseComment() (Node, error) {
	return &Comment{Token: p.curr}, nil
}

func (p *Parser) parseText() (Node, error) {
	return &Text{Token: p.curr}, nil
}

func (p *Parser) parseExpression() (Node, error) {
	var fn filter.Func

	tok := p.curr

	var e ExpressionNode

	switch tok.Type {
	case TokenResource:
		if err := p.nextExpect(TokenPeriod, TokenIdentifier); err != nil {
			return nil, SyntaxError
		}
		e = &Expression{
			Token: tok,
			Name:  p.curr.Literal,
		}
	default:
		e = &Expression{Token: tok}
	}

	if err := p.nextExpect(TokenPipe, TokenFilter); err == nil {
		if fn, err = filter.Lookup(p.curr.Literal); err != nil {
			return nil, InternalError
		}
		e.WithFilter(fn)
	}

	// Probably just need to switch here based on type?
	return e, nil
}

func (p *Parser) parseStatement() (Node, error) {
	// Probably just need to switch here based on type?
	tok := p.curr
	switch p.curr.Literal { //nolint:gocritic
	case "let":
		if err := p.nextExpect(TokenIdentifier); err != nil {
			return nil, Error("identifier")
		}
		id := p.curr.Literal

		if err := p.nextExpect(TokenEquals); err != nil {
			return nil, Error("equals")
		}

		p.nextToken()

		expNode, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		return &LetStatement{Token: tok, Identifier: id, Expression: expNode.(ExpressionNode)}, nil
	}
	return nil, InternalError
}

func (p *Parser) nextToken() {
	p.curr = p.peek
	p.peek = p.l.Next()
}

func (p *Parser) nextExpect(expected ...TokenType) error {
	for _, ex := range expected {
		if ex != p.peek.Type {
			return UnexpectedTokenError
		}

		p.nextToken()
	}

	return nil
}

// Once whitespace control gets back in here, we need to be able to move forward and
// backward either before or after the delimiter.
