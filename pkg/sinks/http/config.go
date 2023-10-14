package http

import "fmt"

// HeaderConfig is a configuration for a single HTTP header.
type HeaderConfig struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

// validate ensures that the header config is valid.
func (h *HeaderConfig) validate() error {
	if h.Name == "" || h.Value == "" {
		return fmt.Errorf("name and value must be provided")
	}

	return nil
}

// UnmarshalYAML implements yaml.Unmarshaler for defaulting the header config.
func (h *HeaderConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type HeaderConfigDefaulted HeaderConfig
	var defaults = HeaderConfigDefaulted{}

	out := defaults
	if err := unmarshal(&out); err != nil {
		return err
	}

	header := HeaderConfig(out)
	if err := header.validate(); err != nil {
		return err
	}

	*h = header
	return nil
}

// Config is the configuration for an HTTP sink.
type Config struct {
	URL     string `yaml:"url"`
	Headers []HeaderConfig
	Method  string `yaml:"method"`
}

// TODO: validation

// UnmarshalYAML implements yaml.Unmarshaler for defaulting the HTTP config.
func (h *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type ConfigDefaulted Config
	// TODO: make const defaults
	var defaults = ConfigDefaulted{
		URL:    DefaultHTTPUrl,
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
