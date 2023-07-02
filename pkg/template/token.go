package template

var builtins = map[string]TokenType{
	"integer_range": TokenResource,
	"list":          TokenResource,
	"random_string": TokenResource,
	"timestamp":     TokenResource,
	"uuid":          TokenResource,
	"map":           TokenResource,

	"let": TokenKeyword,

	"abs":         TokenFilter,
	"select":      TokenFilter,
	"unique":      TokenFilter,
	"attr":        TokenFilter,
	"max":         TokenFilter,
	"upper":       TokenFilter,
	"batch":       TokenFilter,
	"min":         TokenFilter,
	"urlencode":   TokenFilter,
	"capitalize":  TokenFilter,
	"sort":        TokenFilter,
	"urlize":      TokenFilter,
	"int":         TokenFilter,
	"wordcount":   TokenFilter,
	"default":     TokenFilter,
	"reject":      TokenFilter,
	"striptags":   TokenFilter,
	"passthrough": TokenFilter,
	"wordwrap":    TokenFilter,
	"dictsort":    TokenFilter,
	"join":        TokenFilter,
	"groupattr":   TokenFilter,
	"acceptattr":  TokenFilter,
	"rejectattr":  TokenFilter,
	"sum":         TokenFilter,
	"escape":      TokenFilter,
	"replace":     TokenFilter,
	"length":      TokenFilter,
	"reverse":     TokenFilter,
	"tojson":      TokenFilter,
	"round":       TokenFilter,
	"trim":        TokenFilter,
	"float":       TokenFilter,
	"lower":       TokenFilter,
	"safe":        TokenFilter,
	"truncate":    TokenFilter,
}

type Token struct {
	Type    TokenType
	Literal string

	file string
	line int
	col  int
}

func NewToken(t TokenType, l string) Token {
	return Token{Type: t, Literal: l}
}

// Can I make things like text and strings fit in here as well?  Strings I can
// do just by keeping the quotes intact.  Text is not viable.
func TokenLookup(literal string) Token {
	if tok, ok := builtins[literal]; ok {
		return NewToken(tok, literal)
	}
	return NewToken(TokenIdentifier, literal)
}

// func (t *Token) WithControl(cntl ControlType) Token {
// 	t.Control = cntl
// 	return t
// }

// func (t *Token) WithMetadata(l *Lexer) *Token {
// 	t.line = l.line
// 	t.col = l.col
// 	t.file = l.file

// 	return t
// }

// func (t *Token) Append(tok *Token) *Token {
// 	// Since we are the only ones using this we should be able
// 	// to ensure that we always pass the last token in the chain,
// 	// but just to be sure, fast forward.
// 	last := t.Last()
// 	last.next = tok
// 	// Now fast forward to the last of the token chain we
// 	// just added.
// 	return t.Last()
// }

// func (t *Token) Next() *Token {
// 	if t.next != nil {
// 		return t.next
// 	}

// 	return nil
// }

// func (t *Token) Prev() *Token {
// 	if t.prev != nil {
// 		return t.prev
// 	}

// 	return nil
// }

// func (t *Token) First() *Token {
// 	tok := t
// 	for tok.prev != nil {
// 		tok = tok.prev
// 	}
// 	return tok
// }

// func (t *Token) Last() *Token {
// 	tok := t
// 	for tok.next != nil {
// 		tok = tok.next
// 	}
// 	return tok
// }
