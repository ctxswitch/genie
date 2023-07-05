package config

import (
	"fmt"
)

// Sinks represents the top level configuration block for all sink
// types.
type Sinks struct {
	Http map[string]Http `yaml:"http"`
}

type HttpHeader struct {
	Name     string
	Value    string
	Resource string
}

func (h *HttpHeader) validate() (bool, error) {
	if h.Resource != "" && h.Value != "" {
		return false, fmt.Errorf("Resource and value are exclusive")
	}
	return true, nil
}

func (h *HttpHeader) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type HttpHeaderDefaulted HttpHeader
	var defaults = HttpHeaderDefaulted{}

	out := defaults
	if err := unmarshal(&out); err != nil {
		return err
	}

	tmpl := HttpHeader(out)
	if valid, err := tmpl.validate(); !valid {
		return err
	}

	*h = tmpl
	return nil
}

// A sink of type "HTTP"; i.e., an HttpSink.
type Http struct {
	Url     string
	Headers []HttpHeader
}

func (h *Http) validate() (bool, error) {
	if h.Url == "" {
		return false, fmt.Errorf("required URL missing for http definition")
	}

	return true, nil
}

// UnmarshalYAML provides a custom interface to set defaults when parsing
// yaml bytes for the http sink.
func (h *Http) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type HttpDefaulted Http
	var defaults = HttpDefaulted{}

	out := defaults
	if err := unmarshal(&out); err != nil {
		return err
	}

	tmpl := Http(out)
	if valid, err := tmpl.validate(); !valid {
		return err
	}

	*h = tmpl
	return nil
}
