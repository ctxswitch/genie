package template

import (
	"reflect"
	"runtime"
	"testing"

	"ctx.sh/genie/pkg/filter"
	"ctx.sh/genie/pkg/resources"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	tests := []struct {
		input    string
		expected []Node
	}{
		{`<# Comment #>`, []Node{
			&Comment{
				NewToken(TokenComment, "Comment"),
			},
		}},
		{`<< name >>`, []Node{
			&Expression{
				Token: NewToken(TokenIdentifier, "name"),
			},
		}},
		{`<< name | capitalize >>`, []Node{
			&Expression{
				Token:  NewToken(TokenIdentifier, "name"),
				Filter: filter.Capitalize,
			},
		}},
		{`<< list.name >>`, []Node{
			&Expression{
				Token:    NewToken(TokenResource, "list"),
				Resource: resources.MockResources().MustGet("list", "name"),
			},
		}},
		{`<% let name = list.name %>`, []Node{
			&LetStatement{
				Token:      NewToken(TokenKeyword, "let"),
				Identifier: "name",
				Expression: &Expression{
					Token:    NewToken(TokenResource, "list"),
					Resource: resources.MockResources().MustGet("list", "name"),
				},
			},
		}},
	}

	for i, tt := range tests {
		root, err := NewParser(tt.input, resources.MockResources()).Parse()
		require.NoError(t, err)
		require.Len(t, tt.expected, root.Length())
		for j, exp := range tt.expected {
			switch n := exp.(type) {
			case *Expression:
				got := root.Nodes[j].(*Expression)
				assert.EqualValues(t, n.Token, got.Token, "test[%d]: %s", i, tt.input)
				assert.EqualValues(t, n.Resource, got.Resource, "test[%d]: %s", i, tt.input)
				// Get around the inability for testify to compare function pointers.
				fn1 := runtime.FuncForPC(reflect.ValueOf(n.Filter).Pointer()).Name()
				fn2 := runtime.FuncForPC(reflect.ValueOf(got.Filter).Pointer()).Name()
				assert.Equal(t, fn1, fn2)
			}
		}
	}
}
