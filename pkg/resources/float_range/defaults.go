package float_range //nolint:revive

// TODO: add precision to the float_range resource

const (
	// DefaultFloatRangeMin is the default minimum value for the integer_range resource.
	DefaultFloatRangeMin float64 = 0.0
	// DefaultFloatRangeMax is the default maximum value for the integer_range resource.
	DefaultFloatRangeMax float64 = 10.0
	// DefaultFloatRangeDistribution is the default distribution.
	DefaultFloatRangeDistribution string = "uniform"
	// DefaultFloatRangeStandardDeviation is the default standard deviation for the normal distribution.
	DefaultFloatRangeStandardDeviation float64 = 1.0
	// DefaultFloatRangeMean is the default mean for the normal distribution.
	DefaultFloatRangeMean float64 = 0.0
	// DefaultFloatRangeRate is the default rate of occurrences in an interval for the exponential distribution.
	DefaultFloatRangeRate float64 = 1.0
	// DefaultFloatRangeFormat is the default format for the float_range resource.
	DefaultFloatRangeFormat string = "none"
	// DefaultFloatRangePrecision is the default precision for the float_range resource.
	DefaultFloatRangePrecision int = -1
)
