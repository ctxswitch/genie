package uuid

import (
	"fmt"
	"strings"
)

// Config is the configuration for the uuid resource
type Config struct {
	Type    string `yaml:"type"`
	Uniques int    `yaml:"uniques"`
}

// validate ensures the config is valid
func (u *Config) validate() error {
	if !(u.Type == "uuid1" || u.Type == "uuid4") {
		return fmt.Errorf("unsupported UUID type %s for uuid", u.Type)
	}
	return nil
}

// UnmarshalYAML implements yaml.Unmarshaler for defaulting the uuid config.
func (u *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type ConfigDefaulted Config
	var defaults = ConfigDefaulted{
		Type:    DefaultUUIDType,
		Uniques: DefaultUUIDUniques,
	}

	out := defaults
	if err := unmarshal(&out); err != nil {
		return err
	}

	uuid := Config(out)

	// Ensure we have a lowercase type
	uuid.Type = strings.ToLower(uuid.Type)
	if err := uuid.validate(); err != nil {
		return err
	}

	*u = uuid
	return nil
}
