package template

var builtins = map[string]TokenType{
	"integer_range": TokenResource,
	"list":          TokenResource,
	"random_string": TokenResource,
	"timestamp":     TokenResource,
	"uuid":          TokenResource,
	"map":           TokenResource,

	"let": TokenKeyword,

	"urlencode":  TokenFilter,
	"capitalize": TokenFilter,
	"upper":      TokenFilter,
	"lower":      TokenFilter,
	"join":       TokenFilter,
	"env":        TokenFilter,
	"escape":     TokenFilter,
	"default":    TokenFilter, // Maybe, this seems unnecessary
	"minimize":   TokenFilter,
	"replace":    TokenFilter,
	"length":     TokenFilter,
	"truncate":   TokenFilter,
	"reverse":    TokenFilter,
	"wordwrap":   TokenFilter,
	"trim":       TokenFilter,
	"tojson":     TokenFilter,
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
