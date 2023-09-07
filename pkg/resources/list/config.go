package list

import "fmt"

type Config []string

func (l *Config) validate() error {
	if len(*l) == 0 {
		return fmt.Errorf("items in list cannot be empty")
	}

	return nil
}

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
