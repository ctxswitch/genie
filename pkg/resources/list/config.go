package list

import "fmt"

// Config is the configuration for lists.
type Config []string

// validate ensures that the list is not empty. All other values are
// acceptable.
func (l *Config) validate() error {
	if len(*l) == 0 {
		return fmt.Errorf("items in list cannot be empty")
	}

	return nil
}

// UnmarshalYAML implements yaml.Unmarshaler for defaulting the list config.
// currently there is no defaulting and it is just here for consistency.
func (l *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type ConfigDefaulted Config
	var defaults = ConfigDefaulted{}

	out := defaults
	if err := unmarshal(&out); err != nil {
		return err
	}

	list := Config(out)
	if err := list.validate(); err != nil {
		return err
	}

	*l = list
	return nil
}
