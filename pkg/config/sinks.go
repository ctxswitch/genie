package config

type SinksBlock struct {
	Http map[string]HttpBlock `yaml:"http"`
}

type HttpBlock struct {
	Url     string `yaml:"url"`
	Headers []struct {
		Name     string `yaml:"name"`
		Value    string `yaml:"value"`
		Resource string `yaml:"resource"`
	}
	Method string `yaml:"method"`
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
