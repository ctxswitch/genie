package config

import "fmt"

type Var struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type Event struct {
	Generators int    `yaml:"generators"`
	Vars       []Var  `yaml:"vars"`
	Template   string `yaml:"template"`
	Raw        string `yaml:"raw"`
}

func (e *Event) validate() (bool, error) {
	if e.Template != "" && e.Raw != "" {
		return false, fmt.Errorf("Template and raw are mutually exclusive options")
	}
	return true, nil
}

// UnmarshalYAML sets the Event defaults and parses an event block. The
// event block can either be a string or a transport event.
func (e *Event) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type EventDefaulted Event
	var defaults = EventDefaulted{
		Generators: 1,
	}

	out := defaults
	if err := unmarshal(&out); err != nil {
		return err
	}

	tmpl := Event(out)
	if valid, err := e.validate(); !valid {
		return err
	}

	*e = tmpl
	return nil
}
