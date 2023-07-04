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

func (t *Token) WithMetadata(l *Lexer) *Token {
	t.line = l.line
	t.col = l.col
	t.file = ""

	return t
}
