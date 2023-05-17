package template

import (
	"testing"

	"ctx.sh/dynamo/pkg/resources"
	"ctx.sh/dynamo/pkg/template/token"
	"github.com/stretchr/testify/assert"
)

var mockResource = resources.MockResources()

func TestParse(t *testing.T) {
	tests := []struct {
		input    string
		err      bool
		expected []any
	}{
		{"hello world", false, []any{
			&TextNode{token.Token{Type: token.Text, Literal: "hello world"}},
		}},
		{"{{ name }}", false, []any{
			&IdentifierExpressionNode{token.Token{Type: token.Identifier, Literal: "name"}},
		}},
		{"{{ list.name }}", false, []any{
			&ResourceExpressionNode{
				token.Token{Type: token.Identifier, Literal: "name"},
				mockResource.MustGet("list", "name"),
			},
		}},
		{"{{ list.name", true, []any{}},
		// TODO: This should error, but right now it's just treated as text.
		{"list.name }}", false, []any{
			&TextNode{token.Token{Type: token.Text, Literal: "list.name }}"}},
		}},
		{"{% let name = list.name %}", false, []any{
			&LetStatementNode{
				Token: token.Token{Type: token.Identifier, Literal: "name"},
				Expression: &ResourceExpressionNode{
					token.Token{Type: token.Identifier, Literal: "name"},
					mockResource.MustGet("list", "name"),
				},
			}},
		},
		{"{%let other_name = name %}", false, []any{
			&LetStatementNode{
				Token: token.Token{Type: token.Identifier, Literal: "other_name"},
				Expression: &IdentifierExpressionNode{
					token.Token{Type: token.Identifier, Literal: "name"},
				},
			},
		}},
	}

	for i, tt := range tests {
		parser := NewParser(tt.input).WithResources(resources.MockResources())
		root, err := parser.Parse()
		if tt.err {
			assert.Error(t, err, "test[%d]: error expected but not encountered", i)
		} else {
			assert.NoError(t, err, "test[%d]: %s", i, err)
		}
		assert.EqualValues(t, tt.expected, root.Nodes, "test[%d]: %s", i, parser.Error())
	}
}
