package timestamp

// Config is the configuration for the timestamp resource.
type Config struct {
	Format    string `yaml:"format"`
	Timestamp string `yaml:"timestamp"`
}

// validate validates the timestamp resource configuration.  It cureently
// does nothing, but is here for future use.  There are things that we should
// be validating.
func (t *Config) validate() error {
	return nil
}

// UnmarshalYAML implements yaml.Unmarshaler for defaulting the timestamp
func (t *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type ConfigDefaulted Config
	var defaults = ConfigDefaulted{
		Format:    DefaultTimeFormat,
		Timestamp: DefaultTimeTimestamp,
	}

	out := defaults
	if err := unmarshal(&out); err != nil {
		return err
	}

	ts := Config(out)
	if err := ts.validate(); err != nil {
		return err
	}

	*t = ts
	return nil
}
