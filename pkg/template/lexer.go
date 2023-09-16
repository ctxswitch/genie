package template

import "strings"

type scannerFunc func() (Token, error)

type Lexer struct {
	input string
	index int
	next  int
	ch    byte

	// Track metadata while scanning
	// file string
	line int
	col  int

	// Scanner mode
	state   LexerMode
	scanner map[mode]scannerFunc
}

func NewLexer(input string) *Lexer {
	l := &Lexer{
		input:   input,
		line:    1,
		col:     0,
		scanner: make(map[mode]scannerFunc),
	}

	l.state.Start(TextMode)
	l.registerScanner(TextMode, l.scanText)
	l.registerScanner(ExpressionMode, l.scanExpression)
	l.registerScanner(StatementMode, l.scanStatement)
	l.registerScanner(CommentMode, l.scanComment)
	l.registerScanner(RawMode, l.scanRaw)

	l.readByte()
	return l
}

// Come back through and get rid of recursion
func (l *Lexer) Next() Token {
	if l.ch == 0 {
		return NewToken(TokenEOF, "EOF")
	}

	delim := l.peekDelimiter()
	switch delim {
	// I think we should move the token scanning to the scanners
	case "<<":
		l.shift(2)
		l.state.Start(ExpressionMode)
		l.skipWhitespace()
	case ">>":
		l.shift(2)
		if err := l.state.End(ExpressionMode); err != nil {
			return NewToken(TokenError, delim)
		}
		return l.Next()
	case "<%":
		l.shift(2)
		l.state.Start(StatementMode)
		l.skipWhitespace()
	case "%>":
		l.shift(2)
		// Automatically trim
		if l.ch == '\n' {
			l.readByte()
		}
		if err := l.state.End(StatementMode); err != nil {
			return NewToken(TokenError, delim)
		}
		return l.Next()
	case "<*":
		l.shift(2)
		l.state.Start(RawMode)
	case "*>":
		l.shift(2)
		// Automatically trim
		if l.ch == '\n' {
			l.readByte()
		}
		if err := l.state.End(RawMode); err != nil {
			return NewToken(TokenError, delim)
		}
		return l.Next()
	case "<#":
		l.shift(2)
		l.state.Start(CommentMode)
	case "#>":
		l.shift(2)
		// Automatically trim
		if l.ch == '\n' {
			l.readByte()
		}
		if err := l.state.End(CommentMode); err != nil {
			return NewToken(TokenError, delim)
		}
		return l.Next()
	}

	mode := l.state.Mode()
	fn, ok := l.scanner[mode]
	if !ok {
		return NewToken(TokenError, string(mode))
	}

	tok, err := fn()
	if err != nil {
		return NewToken(TokenError, err.Error())
	}

	return tok

}

func (l *Lexer) scanComment() (Token, error) {
	return NewToken(TokenComment, l.readComment()), nil
}

func (l *Lexer) scanExpression() (Token, error) {
	var tok Token

	switch ch := l.ch; {
	case ch == '.':
		l.readByte()
		tok = NewToken(TokenPeriod, ".")
	case l.ch == '"':
		tok = NewToken(TokenString, l.readString())
	case l.ch == '|':
		l.readByte()
		tok = NewToken(TokenPipe, "|")
	case 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_':
		tok = TokenLookup(l.readIdentifier())
	default:
		return NewToken(TokenError, string(l.ch)), UnexpectedTokenError
	}

	if tok.Type != TokenPeriod && tok.Type != TokenResource {
		l.skipWhitespace()
	}

	return tok, nil
}

func (l *Lexer) scanRaw() (Token, error) {
	return NewToken(TokenText, l.readRaw()), nil
}

func (l *Lexer) scanStatement() (Token, error) {
	var tok Token

	l.skipWhitespace()

	switch ch := l.ch; {
	case ch == '=':
		l.readByte()
		tok = NewToken(TokenEquals, "=")
	case ch == '.':
		l.readByte()
		tok = NewToken(TokenPeriod, ".")
	case l.ch == '"':
		tok = NewToken(TokenString, l.readString())
	case l.ch == '|':
		l.readByte()
		tok = NewToken(TokenPipe, "|")
	case 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_':
		tok = TokenLookup(l.readIdentifier())
	default:
		return NewToken(TokenError, string(l.ch)), UnexpectedTokenError
	}

	if tok.Type != TokenPeriod && tok.Type != TokenResource {
		l.skipWhitespace()
	}

	return tok, nil
}

func (l *Lexer) scanText() (Token, error) {
	literal := l.readText()
	return NewToken(TokenText, literal), nil
}

func (l *Lexer) peekDelimiter() string {
	if l.peek() != 0 && l.ch != 0 {
		return l.input[l.index : l.index+2]
	}
	return l.input[l.index:len(l.input)]
}

func (l *Lexer) readComment() string {
	start := l.index
	for l.index < len(l.input) {
		peek := l.peek()
		if l.ch == '#' && peek == '>' {
			break
		}

		l.readByte()
	}
	return strings.TrimSpace(l.input[start:l.index])
}

func (l *Lexer) readByte() {
	if l.ch == '\n' {
		l.col = 0
		l.line++
	}
	l.ch = l.peek()
	l.index = l.next
	l.next++
	l.col++
}

func (l *Lexer) readIdentifier() string {
	start := l.index
	for l.index < len(l.input) && 'a' <= l.ch && l.ch <= 'z' || 'A' <= l.ch && l.ch <= 'Z' || l.ch == '_' || '0' <= l.ch && l.ch <= '9' {
		l.readByte()
	}

	return l.input[start:l.index]
}

func (l *Lexer) readRaw() string {
	start := l.index
	for l.index < len(l.input) {
		peek := l.peek()
		if l.ch == '*' && peek == '>' {
			break
		}

		l.readByte()
	}
	return strings.TrimSpace(l.input[start:l.index])
}

func (l *Lexer) readText() string {
	start := l.index

	for l.index < len(l.input) {
		peek := l.peek()
		if l.ch == '<' && (peek == '<' || peek == '%' || peek == '#' || peek == '*') {
			break
		}

		l.readByte()
	}
	return l.input[start:l.index]
}

func (l *Lexer) readString() string {
	l.readByte() // move past the first quote

	start := l.index
	for l.ch != '"' {
		l.readByte()
	}

	str := l.input[start:l.index]
	l.readByte() // move past the last quote

	return str
}

func (l *Lexer) registerScanner(m mode, fn scannerFunc) {
	l.scanner[m] = fn
}

func (l *Lexer) peek() byte {
	if l.next >= len(l.input) {
		return 0
	}
	return l.input[l.next]
}

func (l *Lexer) shift(n int) {
	for i := 0; i < n; i++ {
		l.readByte()
	}
}

func (l *Lexer) skipWhitespace() {
	for l.index < len(l.input) && (l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r') {
		l.readByte()
	}
}
