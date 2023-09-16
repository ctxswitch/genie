package template

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestScan(t *testing.T) {
	tests := []struct {
		input    string
		expected []Token
	}{
		{`Hello World`, []Token{
			NewToken(TokenText, "Hello World"),
		}},

		// Expressions
		{`<< name >>`, []Token{
			NewToken(TokenIdentifier, "name"),
		}},
		{`<< "name" >>`, []Token{
			NewToken(TokenString, "name"),
		}},
		{`<< list.name >>`, []Token{
			NewToken(TokenResource, "list"),
			NewToken(TokenPeriod, "."),
			NewToken(TokenIdentifier, "name"),
		}},
		{`<< integer_range.name >>`, []Token{
			NewToken(TokenResource, "integer_range"),
			NewToken(TokenPeriod, "."),
			NewToken(TokenIdentifier, "name"),
		}},
		{`<< random_string.name >>`, []Token{
			NewToken(TokenResource, "random_string"),
			NewToken(TokenPeriod, "."),
			NewToken(TokenIdentifier, "name"),
		}},
		{`<< map.name >>`, []Token{
			NewToken(TokenResource, "map"),
			NewToken(TokenPeriod, "."),
			NewToken(TokenIdentifier, "name"),
		}},
		{`<< map.names.first >>`, []Token{
			NewToken(TokenResource, "map"),
			NewToken(TokenPeriod, "."),
			NewToken(TokenIdentifier, "names"),
			NewToken(TokenPeriod, "."),
			NewToken(TokenIdentifier, "first"),
		}},
		{`<< timestamp.name >>`, []Token{
			NewToken(TokenResource, "timestamp"),
			NewToken(TokenPeriod, "."),
			NewToken(TokenIdentifier, "name"),
		}},
		{`<< uuid.name >>`, []Token{
			NewToken(TokenResource, "uuid"),
			NewToken(TokenPeriod, "."),
			NewToken(TokenIdentifier, "name"),
		}},
		{`<< map.name | tojson >>`, []Token{
			NewToken(TokenResource, "map"),
			NewToken(TokenPeriod, "."),
			NewToken(TokenIdentifier, "name"),
			NewToken(TokenPipe, "|"),
			NewToken(TokenFilter, "tojson"),
		}},

		// More expression cases
		{`<< uuid#name >>`, []Token{
			NewToken(TokenResource, "uuid"),
			NewToken(TokenError, "unexpected token encountered"),
		}},
		{`<< uuid. name >>`, []Token{
			NewToken(TokenResource, "uuid"),
			NewToken(TokenPeriod, "."),
			NewToken(TokenError, "unexpected token encountered"),
		}},
		{`<< uuid .name >>`, []Token{
			NewToken(TokenResource, "uuid"),
			NewToken(TokenError, "unexpected token encountered"),
		}},
		{`<<
			uuid.name
		  >>`, []Token{
			NewToken(TokenResource, "uuid"),
			NewToken(TokenPeriod, "."),
			NewToken(TokenIdentifier, "name"),
		}},
		{`<<ipaddr.aws_us_east_1>> <<list.left_names>> <<list.right_names>>`, []Token{
			NewToken(TokenResource, "ipaddr"),
			NewToken(TokenPeriod, "."),
			NewToken(TokenIdentifier, "aws_us_east_1"),
		}},

		// Statement
		{`<% let greeting = "hello" %>`, []Token{
			NewToken(TokenKeyword, "let"),
			NewToken(TokenIdentifier, "greeting"),
			NewToken(TokenEquals, "="),
			NewToken(TokenString, "hello"),
		}},
		{`<% let greeting = value %>`, []Token{
			NewToken(TokenKeyword, "let"),
			NewToken(TokenIdentifier, "greeting"),
			NewToken(TokenEquals, "="),
			NewToken(TokenIdentifier, "value"),
		}},
		{`<% let greeting = list.value %>`, []Token{
			NewToken(TokenKeyword, "let"),
			NewToken(TokenIdentifier, "greeting"),
			NewToken(TokenEquals, "="),
			NewToken(TokenResource, "list"),
			NewToken(TokenPeriod, "."),
			NewToken(TokenIdentifier, "value"),
		}},

		// Comments
		{`<# I'm a comment #>`, []Token{
			NewToken(TokenComment, "I'm a comment"),
		}},

		// Raw
		{`<* << name >> <% let foo="bar" %> <# comment #> *>`, []Token{
			NewToken(TokenText, "<< name >> <% let foo=\"bar\" %> <# comment #>"),
		}},

		// Compound
		{`<# I'm a comment #><% let name = list.name %>Hello << name >>`, []Token{
			NewToken(TokenComment, "I'm a comment"),
			NewToken(TokenKeyword, "let"),
			NewToken(TokenIdentifier, "name"),
			NewToken(TokenEquals, "="),
			NewToken(TokenResource, "list"),
			NewToken(TokenPeriod, "."),
			NewToken(TokenIdentifier, "name"),
			NewToken(TokenText, "Hello "),
			NewToken(TokenIdentifier, "name"),
		}},
	}

	for _, tt := range tests {
		// lexer := NewLexer(tt.input)
		// tokens, err := lexer.Scan()

		// if tt.err != nil {
		// 	require.Errorf(t, err, "token: %v", tokens)
		// 	require.Len(t, tokens, 1)
		// 	require.EqualExportedValuesf(t, tt.expected[0], tokens[0], "test: %s", tt.input)
		// } else {
		// 	require.NoErrorf(t, err, "token: %v", tokens)
		// 	for i, expected := range tt.expected {
		// 		tok := tokens[i]
		// 		require.EqualExportedValuesf(t, expected, tok, "test[%d]: %s", i, tt.input)
		// 	}
		// }

		lexer := NewLexer(tt.input)
		for i, exp := range tt.expected {
			tok := lexer.Next()
			require.EqualExportedValuesf(t, exp, tok, "test[%d]: %s", i, tt.input)
		}

		// if tt.err != nil {
		// 	token, err := lexer.Next()
		// 	require.Errorf(t, err, "token: %v", token)
		// 	require.Len(t, token, 1)
		// 	require.EqualExportedValuesf(t, tt.expected[0], token, "test: %s", tt.input)
		// } else {
		// 	require.NoErrorf(t, err, "token: %v", token)
		// 	for i, expected := range tt.expected {
		// 		tok := tokens[i]
		// 		require.EqualExportedValuesf(t, expected, tok, "test[%d]: %s", i, tt.input)
		// 	}
		// }
	}
}
