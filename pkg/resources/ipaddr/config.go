package ipaddr

import "net"

type Config struct {
	Cidrs   []string `yaml:"cidrs"`
	Uniques uint32   `yaml:"uniques"`
}

func (i *Config) validate() error {
	for _, cidr := range i.Cidrs {
		_, _, err := net.ParseCIDR(cidr)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type ConfigDefaulted Config
	var defaults = ConfigDefaulted{
		Cidrs:   DefaultIPAddrCidrs,
		Uniques: DefaultIPAddrUniques,
	}

	out := defaults
	if err := unmarshal(&out); err != nil {
		return err
	}

	tmpl := Config(out)
	if err := tmpl.validate(); err != nil {
		return err
	}

	*i = tmpl
	return nil
}
