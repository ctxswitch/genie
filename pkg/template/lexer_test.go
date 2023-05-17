package template

import (
	"testing"

	"ctx.sh/dynamo/pkg/template/token"
	"github.com/stretchr/testify/assert"
)

func TestNextToken(t *testing.T) {
	tests := []struct {
		input    string
		expected []token.Token
	}{
		{"{}", []token.Token{
			{Type: token.Text, Literal: "{}"},
		}},
		{"{{}}", []token.Token{
			{Type: token.ExpressionStart, Literal: "{{"},
			{Type: token.ExpressionEnd, Literal: "}}"},
		}},
		// TODO: set up escaping chars
		// {"\\{{\\}}", []token.Token{
		// 	{Type: token.Text, Literal: "{{}}"},
		// }},
		{"hello {{}} world", []token.Token{
			{Type: token.Text, Literal: "hello "},
			{Type: token.ExpressionStart, Literal: "{{"},
			{Type: token.ExpressionEnd, Literal: "}}"},
			{Type: token.Text, Literal: " world"},
		}},
		{"hello {{ name }} world", []token.Token{
			{Type: token.Text, Literal: "hello "},
			{Type: token.ExpressionStart, Literal: "{{"},
			{Type: token.Identifier, Literal: "name"},
			{Type: token.ExpressionEnd, Literal: "}}"},
			{Type: token.Text, Literal: " world"},
		}},
		{"hello {{list.name}} world", []token.Token{
			{Type: token.Text, Literal: "hello "},
			{Type: token.ExpressionStart, Literal: "{{"},
			{Type: token.Resource, Literal: "list"},
			{Type: token.Period, Literal: "."},
			{Type: token.Identifier, Literal: "name"},
			{Type: token.ExpressionEnd, Literal: "}}"},
			{Type: token.Text, Literal: " world"},
		}},
		{"hello {{ list.name }} world", []token.Token{
			{Type: token.Text, Literal: "hello "},
			{Type: token.ExpressionStart, Literal: "{{"},
			{Type: token.Resource, Literal: "list"},
			{Type: token.Period, Literal: "."},
			{Type: token.Identifier, Literal: "name"},
			{Type: token.ExpressionEnd, Literal: "}}"},
			{Type: token.Text, Literal: " world"},
		}},
		{"hello {{       list.name       }} world", []token.Token{
			{Type: token.Text, Literal: "hello "},
			{Type: token.ExpressionStart, Literal: "{{"},
			{Type: token.Resource, Literal: "list"},
			{Type: token.Period, Literal: "."},
			{Type: token.Identifier, Literal: "name"},
			{Type: token.ExpressionEnd, Literal: "}}"},
			{Type: token.Text, Literal: " world"},
		}},
		{"hello {{ integer_range.name }} world", []token.Token{
			{Type: token.Text, Literal: "hello "},
			{Type: token.ExpressionStart, Literal: "{{"},
			{Type: token.Resource, Literal: "integer_range"},
			{Type: token.Period, Literal: "."},
			{Type: token.Identifier, Literal: "name"},
			{Type: token.ExpressionEnd, Literal: "}}"},
			{Type: token.Text, Literal: " world"},
		}},
		{"hello {{ random_string.name }} world", []token.Token{
			{Type: token.Text, Literal: "hello "},
			{Type: token.ExpressionStart, Literal: "{{"},
			{Type: token.Resource, Literal: "random_string"},
			{Type: token.Period, Literal: "."},
			{Type: token.Identifier, Literal: "name"},
			{Type: token.ExpressionEnd, Literal: "}}"},
			{Type: token.Text, Literal: " world"},
		}},
		{"hello {{ timestamp.name }} world", []token.Token{
			{Type: token.Text, Literal: "hello "},
			{Type: token.ExpressionStart, Literal: "{{"},
			{Type: token.Resource, Literal: "timestamp"},
			{Type: token.Period, Literal: "."},
			{Type: token.Identifier, Literal: "name"},
			{Type: token.ExpressionEnd, Literal: "}}"},
			{Type: token.Text, Literal: " world"},
		}},
		{"hello {{ uuid.name }} world", []token.Token{
			{Type: token.Text, Literal: "hello "},
			{Type: token.ExpressionStart, Literal: "{{"},
			{Type: token.Resource, Literal: "uuid"},
			{Type: token.Period, Literal: "."},
			{Type: token.Identifier, Literal: "name"},
			{Type: token.ExpressionEnd, Literal: "}}"},
			{Type: token.Text, Literal: " world"},
		}},
		{"{% let name = list.name %}", []token.Token{
			{Type: token.StatementStart, Literal: "{%"},
			{Type: token.Keyword, Literal: "let"},
			{Type: token.Identifier, Literal: "name"},
			{Type: token.Equals, Literal: "="},
			{Type: token.Resource, Literal: "list"},
			{Type: token.Period, Literal: "."},
			{Type: token.Identifier, Literal: "name"},
			{Type: token.StatementEnd, Literal: "%}"},
		}},
		{"{% let name = list.name %}\nhello {{ name }} world", []token.Token{
			{Type: token.StatementStart, Literal: "{%"},
			{Type: token.Keyword, Literal: "let"},
			{Type: token.Identifier, Literal: "name"},
			{Type: token.Equals, Literal: "="},
			{Type: token.Resource, Literal: "list"},
			{Type: token.Period, Literal: "."},
			{Type: token.Identifier, Literal: "name"},
			{Type: token.StatementEnd, Literal: "%}"},
			{Type: token.Text, Literal: "\nhello "},
			{Type: token.ExpressionStart, Literal: "{{"},
			{Type: token.Identifier, Literal: "name"},
			{Type: token.ExpressionEnd, Literal: "}}"},
			{Type: token.Text, Literal: " world"},
		}},
		{"{{{{}}", []token.Token{
			{Type: token.ExpressionStart, Literal: "{{"},
			{Type: token.ExpressionStart, Literal: "{{"},
			{Type: token.ExpressionEnd, Literal: "}}"},
		}},
		{"hello world", []token.Token{
			{Type: token.Text, Literal: "hello world"},
		}},
	}

	for _, tt := range tests {
		lexer := NewLexer(tt.input)
		for _, expc := range tt.expected {
			got := lexer.NextToken()
			assert.Equal(t, expc, got, "test[%s]: %s != %s", tt.input, expc.Literal, got.Literal)
		}
	}
}
