package config

type VariableBlock struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type EventBlock struct {
	Generators int             `yaml:"generators"`
	Vars       []VariableBlock `yaml:"vars"`
	Template   string          `yaml:"template"`
	Raw        string          `yaml:"raw"`
	Sinks      []string        `yaml:"sink"`
}

func (e *EventBlock) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type EventBlockDefaults struct {
		Generators int             `yaml:"generators"`
		Vars       []VariableBlock `yaml:"vars"`
		Template   string          `yaml:"template"`
		Raw        string          `yaml:"raw"`
		Sinks      []string        `yaml:"sink"`
	}

	var defaults = EventBlockDefaults{
		Generators: 1,
	}
	out := defaults

	if err := unmarshal(&out); err != nil {
		return err
	}

	if out.Sinks == nil {
		out.Sinks = []string{"stdout"}
	}

	evt := EventBlock(out)

	*e = evt
	return nil
}

type EventsBlock map[string]EventBlock
