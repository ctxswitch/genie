package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type ConfigBlock struct {
	Events    EventsBlock    `yaml:"events"`
	Resources ResourcesBlock `yaml:"resources"`
	Sinks     SinksBlock     `yaml:"sinks"`
}

func Load(paths []string) (ConfigBlock, error) {
	var out ConfigBlock

	for _, path := range paths {
		files, err := filepath.Glob(fmt.Sprintf("%s/*.y[a]*ml", path))
		if err != nil {
			return out, err
		}
		for _, file := range files {
			fmt.Println(file)
			in, ferr := os.ReadFile(file)
			if ferr != nil {
				return out, ferr
			}
			yerr := yaml.Unmarshal(in, &out)
			if yerr != nil {
				return out, yerr
			}
		}
	}

	return out, nil
}
