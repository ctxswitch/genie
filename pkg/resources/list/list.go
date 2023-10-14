package list

import (
	"math/rand"
)

// List is a resource representing a list of strings.
type List struct {
	items []string
}

// New returns a new list resource.
func New(config Config) *List {
	return &List{items: config}
}

// WithItems sets the items on the list.
func (l *List) WithItems(items []string) *List {
	l.items = items
	return l
}

// Get returns a random item from the list.
func (l *List) Get() string {
	if len(l.items) == 0 {
		return ""
	}

	return l.items[rand.Intn(len(l.items))] //nolint:gosec
}
