package integer_range //nolint:revive

import (
	"fmt"
	"math/rand"
)

type IntegerRange struct {
	min   int64
	max   int64
	step  int64
	pad   uint32
	cache []string
}

func New(cfg Config) *IntegerRange {
	return &IntegerRange{
		min:  cfg.Min,
		max:  cfg.Max,
		step: cfg.Step,
		pad:  cfg.Pad,
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
