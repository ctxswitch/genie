package list

import (
	"math/rand"

	"ctx.sh/genie/pkg/config"
)

type List struct {
	items []string
}

func FromConfig(options config.List) *List {
	return &List{items: options}
}

func (l List) Get() string {
	return l.items[rand.Intn(len(l.items))]
}
