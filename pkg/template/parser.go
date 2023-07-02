package template

import (
	"ctx.sh/genie/pkg/filter"
	"ctx.sh/genie/pkg/resources"
)

type ParserFunc func() (Node, error)

type Parser struct {
	l       *Lexer
	parsers map[TokenType]ParserFunc
	tokens  []Token

	res  *resources.Resources
	peek Token
	curr Token
}

func NewParser(input string, res *resources.Resources) *Parser {
	p := &Parser{
		l:       NewLexer(input),
		res:     res,
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
			return Root{}, UnknownParserError
		}

		node, err := fn()
		if err != nil {
			return Root{}, err
		}
		root.Nodes = append(root.Nodes, node)

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
	var fn filter.FilterFunc

	tok := p.curr
	e := &Expression{Token: tok}

	switch tok.Type {
	case TokenResource:
		rtype := p.curr.Literal
		if err := p.nextExpect(TokenPeriod, TokenIdentifier); err != nil {
			return nil, SyntaxError
		}

		r, err := p.res.Get(rtype, p.curr.Literal)
		if err != nil {
			return nil, UnknownResourceError
		}
		e.Resource = r
	}

	if err := p.nextExpect(TokenPipe, TokenFilter); err == nil {
		if fn, err = filter.Lookup(p.curr.Literal); err != nil {
			return nil, InternalError
		}
		e.Filter = fn
	}

	// Probably just need to switch here based on type?
	return e, nil
}

func (p *Parser) parseStatement() (Node, error) {
	// Probably just need to switch here based on type?
	return &LetStatement{}, nil
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
