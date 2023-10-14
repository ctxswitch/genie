package integer_range //nolint:revive

import (
	"fmt"
	"math/rand"
)

// IntegerRange is a resource that generates a random integer between a minimum
type IntegerRange struct {
	min   int64
	max   int64
	step  int64
	pad   uint32
	cache []string
}

// New returns a new integer_range resource initialized from the given config.
func New(cfg Config) *IntegerRange {
	return &IntegerRange{
		min:  cfg.Min,
		max:  cfg.Max,
		step: cfg.Step,
		pad:  cfg.Pad,
	}
}

// Cache creates a cache of all possible values for this resource.  This is a
// simple, non-optimal from a memory perspective, but it's a good first step.
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

// Get implements the Resource interface.
func (i *IntegerRange) Get() string {
	if i.cache == nil {
		i.cache = i.Cache()
	}

	return i.cache[rand.Intn(len(i.cache))] //nolint:gosec
}
