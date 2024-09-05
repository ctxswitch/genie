package float_range //nolint:revive

import (
	"strconv"
	"time"

	"golang.org/x/exp/rand"
)

// FloatRange is a resource that generates a random float between a minimum.
type FloatRange struct {
	min          float64
	max          float64
	distribution string
	stddev       *float64
	mean         *float64
	rate         float64
	format       string
	precision    int
}

// New returns a new float_range resource initialized from the given config.
func New(cfg Config) *FloatRange {
	rand.Seed(uint64(time.Now().UnixNano()))

	return &FloatRange{
		min:          cfg.Min,
		max:          cfg.Max,
		distribution: cfg.Distribution,
		stddev:       cfg.StdDev,
		mean:         cfg.Mean,
		rate:         cfg.Rate,
		format:       cfg.Format,
		precision:    cfg.Precision,
	}
}

// Get implements the Resource interface.
func (i *FloatRange) Get() string {
	var float float64

	switch i.distribution {
	case "normal":
		float = rand.NormFloat64()*(*i.stddev) + *i.mean
	case "exp", "exponential":
		float = rand.ExpFloat64()/i.rate + i.min
	default:
		float = rand.Float64()*(i.max-i.min) + i.min
	}

	if float < i.min {
		float = i.min
	} else if float > i.max {
		float = i.max
	}

	return strconv.FormatFloat(float, convertFormat(i.format), i.precision, 64)
}
