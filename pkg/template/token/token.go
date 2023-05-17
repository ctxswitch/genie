package token

type Type string

type Token struct {
	Type    Type
	Literal string
}

const (
	EOF             Type = "EOF"
	Text            Type = "Text"
	Escape          Type = "Escape"
	Identifier      Type = "Identifier"
	Period          Type = "Period"
	String          Type = "String"
	Equals          Type = "Equals"
	ExpressionStart Type = "ExpressionStart"
	ExpressionEnd   Type = "ExpressionEnd"
	StatementStart  Type = "StatementStart"
	StatementEnd    Type = "StatementEnd"
	Resource        Type = "Resource"
	Keyword         Type = "Keyword"
)

var lookup = map[string]Type{
	"\x00":          EOF,
	"integer_range": Resource,
	"list":          Resource,
	"random_string": Resource,
	"timestamp":     Resource,
	"uuid":          Resource,
	"let":           Keyword,
	// Template
	"{{": ExpressionStart,
	"}}": ExpressionEnd,
	"{%": StatementStart,
	"%}": StatementEnd,
}

func Lookup(token string) Type {
	if tok, ok := lookup[token]; ok {
		return tok
	}
	return Identifier
}

func TokenFromString(literal string) Token {
	return Token{Type: Lookup(literal), Literal: literal}
}

func TokenFromBytes(literals ...byte) Token {
	return TokenFromString(string(literals))
}

func TokenFromText(text string) Token {
	return Token{Type: Text, Literal: text}
}
