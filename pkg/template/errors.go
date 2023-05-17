package template

type ParserError string

func (e ParserError) Error() string { return string(e) }

const (
	UnexpectedTokenError ParserError = "unexpected token encountered"
	UnknownKeyword       ParserError = "unknown keyword"
	ResourceNotFound     ParserError = "resource was not found"
)
