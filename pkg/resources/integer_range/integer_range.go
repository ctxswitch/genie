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

func FromConfig(options config.IntegerRange) *IntegerRange {
	return &IntegerRange{
		min:  options.Min,
		max:  options.Max,
		step: options.Step,
		pad:  options.Pad,
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

	return i.cache[rand.Intn(len(i.cache))]
}
