package http

import "fmt"

type HttpHeaderConfig struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

func (h *HttpHeaderConfig) validate() error {
	if h.Name == "" || h.Value == "" {
		return fmt.Errorf("name and value must be provided")
	}

	return nil
}

func (h *HttpHeaderConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type HttpHeaderConfigDefaulted HttpHeaderConfig
	var defaults = HttpHeaderConfigDefaulted{}

	out := defaults
	if err := unmarshal(&out); err != nil {
		return err
	}

	header := HttpHeaderConfig(out)
	if err := header.validate(); err != nil {
		return err
	}

	*h = header
	return nil
}

type Config struct {
	Url     string `yaml:"url"`
	Headers []HttpHeaderConfig
	Method  string `yaml:"method"`
}

func (h *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type ConfigDefaulted Config
	// TODO: make const defaults
	var defaults = ConfigDefaulted{
		Url:    DefaultHttpUrl,
		Method: DefaultMethod,
	}

	out := defaults
	if err := unmarshal(&out); err != nil {
		return err
	}

	http := Config(out)

	*h = http
	return nil
}
