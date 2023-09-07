package list

import (
	"math/rand"
)

type List struct {
	items []string
}

func New(config Config) *List {
	return &List{items: config}
}

func (l *List) WithItems(items []string) *List {
	l.items = items
	return l
}

func (l *List) Get() string {
	if len(l.items) == 0 {
		return ""
	}

	return l.items[rand.Intn(len(l.items))] //nolint:gosec
}
