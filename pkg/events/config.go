package events

import "ctx.sh/genie/pkg/variables"

// EventConfig is the configuration for a single event.
type EventConfig struct {
	Name            string             `yaml:"name"`
	Generators      int                `yaml:"generators"`
	IntervalSeconds float64            `yaml:"intervalSeconds"`
	Vars            []variables.Config `yaml:"vars"`
	Template        string             `yaml:"template"`
	Raw             string             `yaml:"raw"`
}

// UnmarshalYAML implements yaml.Unmarshaler for defaulting the event config.
func (e *EventConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type EventConfigDefaults struct {
		Name            string             `yaml:"name"`
		Generators      int                `yaml:"generators"`
		IntervalSeconds float64            `yaml:"intervalSeconds"`
		Vars            []variables.Config `yaml:"vars"`
		Template        string             `yaml:"template"`
		Raw             string             `yaml:"raw"`
	}

	var defaults = EventConfigDefaults{
		Generators:      1,
		IntervalSeconds: 1.0,
	}
	out := defaults

	if err := unmarshal(&out); err != nil {
		return err
	}

	evt := EventConfig(out)

	*e = evt
	return nil
}

// Config is a collection of event configs.
type Config []EventConfig
