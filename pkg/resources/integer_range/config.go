package integer_range //nolint:revive

import "fmt"

// Config is the configuration for the integer_range resource.
type Config struct {
	Min          int64  `yaml:"min"`
	Max          int64  `yaml:"max"`
	Pad          uint32 `yaml:"pad"`
	Distribution string `yaml:"distribution"`
	StdDev       *float64
	Mean         *int64
}

// validate ensures that the configuration is valid.
func (i *Config) validate() error {
	// Fix me now that we allow negative values
	if i.Max <= i.Min {
		return fmt.Errorf("max (%d) in integer_range must be greater than zero and the minimum value", i.Max)
	}

	if i.Distribution != "uniform" && i.Distribution != "normal" {
		return fmt.Errorf("distribution (%s) in integer_range must be either 'uniform' or 'normal'", i.Distribution)
	}

	if i.StdDev != nil && *i.StdDev > float64(i.Max-i.Min) {
		return fmt.Errorf("standard deviation (%f) in integer_range must be less than the range (%d)", *i.StdDev, i.Max-i.Min)
	}

	if i.Mean != nil && (*i.Mean < i.Min || *i.Mean > i.Max) {
		return fmt.Errorf("mean (%d) in integer_range must be in the range of (%d, %d)", *i.Mean, i.Min, i.Max)
	}

	return nil
}

// UnmarshalYAML implements yaml.Unmarshaler for defaulting the integer_range
func (i *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type ConfigDefaulted Config
	var defaults = ConfigDefaulted{
		Min:          DefaultIntegerRangeMin,
		Max:          DefaultIntegerRangeMax,
		Pad:          DefaultIntegerRangePad,
		Distribution: DefaultIntegerRangeDistribution,
	}

	out := defaults
	if err := unmarshal(&out); err != nil {
		return err
	}

	tmpl := Config(out)
	if err := tmpl.validate(); err != nil {
		return err
	}

	*i = tmpl
	return nil
}
