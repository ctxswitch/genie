package timestamp

type Config struct {
	Format    string `yaml:"format"`
	Timestamp string `yaml:"timestamp"`
}

func (t *Config) validate() error {
	return nil
}

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
