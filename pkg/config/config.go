package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Configs struct {
	Events    map[string]Event `yaml:"events"`
	Resources Resources        `yaml:"resources"`
	Sinks     Sinks            `yaml:"sinks"`
}

func LoadAll(dir string) (*Configs, error) {
	configs := &Configs{}

	files, _ := filepath.Glob(fmt.Sprintf("%s/*.yaml", dir))
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}

		if err := yaml.Unmarshal(data, &configs); err != nil {
			return nil, err
		}

		if err != nil {
			fmt.Printf("error loading file: %v", err)
			continue
		}
	}

	return configs, nil
}
