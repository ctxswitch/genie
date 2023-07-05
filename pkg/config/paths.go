package config

import (
	"fmt"
	"os"
)

type Paths struct {
	Resources string
	Sinks     string
	Events    string
	Templates string
}

func (p *Paths) validate() (bool, error) {
	exists := func(f string) bool {
		_, err := os.Stat(p.Resources)
		return err == nil
	}

	if p.Resources != "" && !exists(p.Resources) {
		return false, fmt.Errorf("Resources directory %s does not exist or is inaccessable", p.Resources)
	}

	if p.Sinks != "" && !exists(p.Sinks) {
		return false, fmt.Errorf("Sinks directory %s does not exist or is inaccessable", p.Sinks)
	}

	if p.Events != "" && !exists(p.Events) {
		return false, fmt.Errorf("Events directory %s does not exist or is inaccessable", p.Events)
	}

	if p.Templates != "" && !exists(p.Templates) {
		return false, fmt.Errorf("Templates directory %s does not exist or is inaccessable", p.Templates)
	}

	return true, nil
}

func (p *Paths) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type PathsDefaulted Paths
	var defaults = PathsDefaulted{}

	out := defaults
	if err := unmarshal(&out); err != nil {
		return err
	}

	tmpl := Paths(out)
	if valid, err := p.validate(); !valid {
		return err
	}

	*p = tmpl
	return nil
}
