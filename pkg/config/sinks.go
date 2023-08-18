package config

import "fmt"

type SinksBlock struct {
	Http map[string]HttpBlock `yaml:"http"`
}

type HttpHeaderBlock struct {
	Name     string `yaml:"name"`
	Value    string `yaml:"value"`
	Resource string `yaml:"resource"`
}

func (h *HttpHeaderBlock) validate() error {
	if h.Value != "" && h.Resource != "" {
		return fmt.Errorf("value and resource are mutually exclusive")
	} else if h.Value == "" && h.Resource == "" {
		return fmt.Errorf("value or resource must contain the header value")
	}

	return nil
}

func (h *HttpHeaderBlock) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type HttpHeaderBlockDefaulted HttpHeaderBlock
	var defaults = HttpHeaderBlockDefaulted{}

	out := defaults
	if err := unmarshal(&out); err != nil {
		return err
	}

	header := HttpHeaderBlock(out)
	if err := header.validate(); err != nil {
		return err
	}

	*h = header
	return nil
}

type HttpBlock struct {
	Url     string `yaml:"url"`
	Headers []HttpHeaderBlock
	Method  string `yaml:"method"`
}

func (h *HttpBlock) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type HttpBlockDefaulted HttpBlock
	// TODO: make const defaults
	var defaults = HttpBlockDefaulted{
		Url:    "http://localhost",
		Method: "POST",
	}

	out := defaults
	if err := unmarshal(&out); err != nil {
		return err
	}

	http := HttpBlock(out)

	*h = http
	return nil
}
