package list

import (
	"math/rand"

	"ctx.sh/genie/pkg/config"
)

type List struct {
	items []string
}

func New(settings config.ListBlock) *List {
	return &List{items: []string(settings)}
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
