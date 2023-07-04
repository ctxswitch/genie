package template

type ControlType string

type TokenType string

const (
	TokenComment  TokenType = "Comment"
	TokenKeyword  TokenType = "Keyword"
	TokenResource TokenType = "Resource"
	TokenText     TokenType = "TokenText"
	TokenFilter   TokenType = "Filter"

	TokenString TokenType = "TokenString"
	TokenChar   TokenType = "Char"
	TokenInt    TokenType = "Int"
	TokenFloat  TokenType = "Float"
	TokenBool   TokenType = "Bool"

	TokenIdentifier TokenType = "TokenIdentifier"
	TokenPeriod     TokenType = "TokenPeriod"
	TokenEquals     TokenType = "Equals"
	TokenPipe       TokenType = "Pipe"

	TokenEOF   TokenType = "EOF"
	TokenError TokenType = "Error"
)
