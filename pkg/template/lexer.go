package template

import "ctx.sh/dynamo/pkg/template/token"

type Lexer struct {
	input    string
	position int
	next     int
	ch       byte
	mode     mode
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input, mode: TextMode}
	l.readByte()
	return l
}

func (l *Lexer) NextToken() token.Token {
	curr := l.ch
	next := l.peek()

	switch curr {
	case '{':
		if next == '{' || next == '%' {
			l.readBytes(2)
			l.setMode(next)
			return token.TokenFromBytes(curr, next)
		}
	case '}':
		if next == '}' {
			l.readBytes(2)
			l.setMode(next)
			return token.TokenFromBytes(curr, next)
		}
	case '%':
		if next == '}' {
			l.readBytes(2)
			l.setMode(next)
			return token.TokenFromBytes(curr, next)
		}
	case 0:
		return token.TokenFromBytes(0)
	}

	switch l.mode {
	case ExpressionMode:
		l.skipWhitespace()
		exp := l.nextExpression()
		l.skipWhitespace()
		return exp
	case StatementMode:
		l.skipWhitespace()
		exp := l.nextStatement()
		l.skipWhitespace()
		return exp
	default:
		return l.nextText()
	}
}

func (l *Lexer) nextText() token.Token {
	return token.TokenFromText(l.readText())
}

func (l *Lexer) nextExpression() token.Token {
	curr := l.ch

	switch curr {
	case ':':
		l.readByte()
		return token.Token{Type: token.Identifier, Literal: l.readIdentifier()}
	case '.':
		l.readByte()
		return token.Token{Type: token.Period, Literal: "."}
	case '"':
		l.readByte()
		return token.Token{Type: token.String, Literal: l.readString()}
	default:
		return token.TokenFromString(l.readIdentifier())
	}
}

func (l *Lexer) nextStatement() token.Token {
	curr := l.ch

	switch curr {
	case '=':
		l.readBytes(2)
		l.mode = ExpressionMode
		return token.Token{Type: token.Equals, Literal: "="}
	default:
		return token.TokenFromString(l.readIdentifier())
	}
}

func (l *Lexer) readByte() {
	l.ch = l.peek()
	l.position = l.next
	l.next++
}

func (l *Lexer) readIdentifier() string {
	start := l.position
	for l.position < len(l.input) && isIdentifier(l.ch) {
		l.readByte()
	}
	return l.input[start:l.position]
}

func (l *Lexer) readText() string {
	start := l.position
	for l.position < len(l.input) {
		curr := l.ch
		next := l.peek()
		if curr == '{' && (next == '%' || next == '{') {
			return l.input[start:l.position]
		}
		l.readByte()
	}
	return l.input[start:l.position]
}

func (l *Lexer) peek() byte {
	if l.next >= len(l.input) {
		return 0
	}
	return l.input[l.next]
}

func (l *Lexer) readBytes(pos int) {
	for i := 0; i < pos; i++ {
		l.readByte()
	}
}

func (l *Lexer) skipWhitespace() {
	for isWhitespace(l.ch) {
		l.readByte()
	}
}

func (l *Lexer) setMode(ch byte) {
	switch ch {
	case '{':
		l.mode = ExpressionMode
	case '%':
		l.mode = StatementMode
	default:
		l.mode = TextMode
	}
}
