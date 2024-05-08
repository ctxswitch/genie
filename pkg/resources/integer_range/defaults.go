package integer_range //nolint:revive

const (
	// DefaultIntegerRangeMin is the default minimum value for the integer_range resource.
	DefaultIntegerRangeMin int64 = 0
	// DefaultIntegerRangeMax is the default maximum value for the integer_range resource.
	DefaultIntegerRangeMax int64 = 10
	// DefaultIntegerRangePad is the default pad value for the integer_range resource.
	DefaultIntegerRangePad uint32 = 0
	// DefaultIntegerRangeDistribution is the default distribution
	DefaultIntegerRangeDistribution string = "uniform"
	// DefaultStandardDeviation is the default standard deviation for the normal distribution
	DefaultStandardDeviation float64 = 1.0
	// DefaultMean is the default mean for the normal distribution
	DefaultMean float64 = 0.0
)
