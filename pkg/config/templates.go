package config

import "fmt"

type Var struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type Template struct {
	// TODO: implement generators through threading in pkg/generator
	Generators int    `yaml:"generators"`
	Vars       []Var  `yaml:"vars"`
	Template   string `yaml:"template"`
	Raw        string `yaml:"raw"`
}

func (t *Template) validate() (bool, error) {
	if t.Template != "" && t.Raw != "" {
		return false, fmt.Errorf("Template and raw are mutually exclusive options")
	}
	return true, nil
}

// UnmarshalYAML sets the Event defaults and parses an event block. The
// event block can either be a string or a transport event.
func (t *Template) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type TemplateDefaulted Template
	var defaults = TemplateDefaulted{
		Generators: 1,
	}

	out := defaults
	if err := unmarshal(&out); err != nil {
		return err
	}

	tmpl := Template(out)
	if valid, err := t.validate(); !valid {
		return err
	}

	*t = tmpl
	return nil
}
