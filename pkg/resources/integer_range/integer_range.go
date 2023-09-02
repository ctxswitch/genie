package integer_range

import (
	"fmt"
	"math/rand"

	"ctx.sh/genie/pkg/config"
)

type IntegerRange struct {
	min   int64
	max   int64
	step  int64
	pad   uint32
	cache []string
}

func New(settings config.IntegerRangeBlock) *IntegerRange {
	return &IntegerRange{
		min:  settings.Min,
		max:  settings.Max,
		step: settings.Step,
		pad:  settings.Pad,
	}
}

func (i *IntegerRange) Cache() []string {
	c := make([]string, 0)
	num := i.min
	for num <= i.max {
		padded := fmt.Sprintf("%0*d", i.pad, num)
		c = append(c, padded)
		num += i.step
	}
	return c
}

func (i *IntegerRange) Get() string {
	if i.cache == nil {
		i.cache = i.Cache()
	}

	return i.cache[rand.Intn(len(i.cache))] //nolint:gosec
}
