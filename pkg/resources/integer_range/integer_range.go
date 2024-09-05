package integer_range //nolint:revive

import (
	"fmt"
	"time"

	"golang.org/x/exp/rand"
)

// IntegerRange is a resource that generates a random integer between a minimum.
type IntegerRange struct {
	min          int64
	max          int64
	pad          uint32
	distribution string
	stddev       *float64
	mean         *int64
	// cache        []string
}

// New returns a new integer_range resource initialized from the given config.
func New(cfg Config) *IntegerRange {
	rand.Seed(uint64(time.Now().UnixNano()))

	return &IntegerRange{
		min:          cfg.Min,
		max:          cfg.Max,
		pad:          cfg.Pad,
		distribution: cfg.Distribution,
		stddev:       cfg.StdDev,
		mean:         cfg.Mean,
	}
}

// Cache creates a cache of all possible values for this resource.  This is a
// simple, non-optimal from a memory perspective, but it's a good first step.
// func (i *IntegerRange) Cache() []string {
// 	c := make([]string, 0)
// 	num := i.min
// 	for num <= i.max {
// 		padded := fmt.Sprintf("%0*d", i.pad, num)
// 		c = append(c, padded)
// 		num += i.step
// 	}
// 	return c
// }

// Get implements the Resource interface.
func (i *IntegerRange) Get() string {
	var integer int64

	switch i.distribution {
	case "normal":
		// TODO: make stddev and mean configurable.  The mean should be allowed to shift.
		// we will also need to floor the i.min and ceil the i.max
		// mean := float64((i.max - i.min) / 2)
		// stddev := float64((i.max - i.min) / 10)

		integer = int64(rand.NormFloat64()*(*i.stddev) + float64(*i.mean))
		if integer < i.min {
			integer = i.min
		} else if integer > i.max {
			integer = i.max
		}
	default:
		integer = rand.Int63n(i.max-i.min) + i.min
	}

	// TODO: Move this to strconv and then move padding to a filter and out of the resource.
	// i.e. <<integer_range.item|pad(5)>>
	return fmt.Sprintf("%0*d", i.pad, integer)
}
