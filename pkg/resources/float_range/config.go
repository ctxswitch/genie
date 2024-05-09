package float_range //nolint:revive

import "fmt"

// Config is the configuration for the float_range resource.
type Config struct {
	Min          float64  `yaml:"min"`
	Max          float64  `yaml:"max"`
	Distribution string   `yaml:"distribution"`
	StdDev       *float64 `yaml:"stddev"`
	Mean         *float64 `yaml:"mean"`
	Rate         float64  `yaml:"rate"`
	Format       string   `yaml:"format"`
	Precision    int      `yaml:"precision"`
}

func (i *Config) validateDistribution(d string) error {
	switch d {
	case "uniform", "normal", "exp", "exponential":
		return nil
	default:
		return fmt.Errorf("distribution (%s) in float_range is unknown", d)
	}
}

func (i *Config) validateFormat(f string) error {
	switch f {
	case "binary", "decimal", "decimal_capitalize", "none", "large", "large_capitalize", "hex", "hex_capitalize":
		return nil
	default:
		return fmt.Errorf("format (%s) in float_range is unknown", f)
	}
}

// validate ensures that the configuration is valid.
func (i *Config) validate() error {
	// Fix me now that we allow negative values
	if i.Max <= i.Min {
		return fmt.Errorf("max (%f) in float_range must be greater than the minimum value", i.Max)
	}

	if err := i.validateDistribution(i.Distribution); err != nil {
		return err
	}

	if err := i.validateFormat(i.Format); err != nil {
		return err
	}

	if i.StdDev != nil && *i.StdDev > float64(i.Max-i.Min) {
		return fmt.Errorf("standard deviation (%f) in float_range must be less than the range (%f)", *i.StdDev, i.Max-i.Min)
	}

	if i.Mean != nil && (*i.Mean < i.Min || *i.Mean > i.Max) {
		return fmt.Errorf("mean (%f) in float_range must be in the range of (%f, %f)", *i.Mean, i.Min, i.Max)
	}

	return nil
}

// UnmarshalYAML implements yaml.Unmarshaler for defaulting the float_range
func (i *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type ConfigDefaulted Config
	var defaults = ConfigDefaulted{
		Min:          DefaultFloatRangeMin,
		Max:          DefaultFloatRangeMax,
		Distribution: DefaultFloatRangeDistribution,
		Rate:         DefaultFloatRangeRate,
		Format:       DefaultFloatRangeFormat,
		Precision:    DefaultFloatRangePrecision,
	}

	out := defaults
	if err := unmarshal(&out); err != nil {
		return err
	}

	tmpl := Config(out)

	// Set the variable attributes if they are nil
	if tmpl.StdDev == nil {
		tmpl.StdDev = new(float64)
		*tmpl.StdDev = float64(tmpl.Max-tmpl.Min) / 8
	}

	if tmpl.Mean == nil {
		tmpl.Mean = new(float64)
		*tmpl.Mean = (tmpl.Max - tmpl.Min) / 2
	}

	if err := tmpl.validate(); err != nil {
		return err
	}

	*i = tmpl
	return nil
}
