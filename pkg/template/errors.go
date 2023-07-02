package template

type Error string

func (e Error) Error() string { return string(e) }

const (
	NoModeSelectedError        Error = "no mode selected"
	ModeNotStartedError        Error = "mode requested to end has not been started"
	ModeUnsetError             Error = "mode has not been set"
	InternalError              Error = "internal error"
	UnexpectedTokenError       Error = "unexpected token encountered"
	NotImplementedError        Error = "not implemented"
	InvalidDelimiterCloseError Error = "invalid delimiter close"
	ResourceNotFound           Error = "resource was not found"
	UnknownParserError         Error = "unknown parser encountered"
	UnknownFilterError         Error = "unknown filter error"
	SyntaxError                Error = "syntax error"
	UnknownResourceError       Error = "unknown resource"
)

func TokenizedError(literal string) []Token {
	return []Token{NewToken(TokenError, literal)}
}
