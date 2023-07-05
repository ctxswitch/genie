package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type Configs struct {
	// Not using paths yet, still trying to figure out if we even want it.
	Paths     Paths            `yaml:"paths"`
	Events    map[string]Event `yaml:"events"`
	Resources Resources        `yaml:"resources"`
	Sinks     Sinks            `yaml:"sinks"`

	templates map[string]string
}

func LoadAll(dir string) (*Configs, error) {
	configs := &Configs{}

	if err := configs.readTemplates("./genie.d"); err != nil {
		return nil, err
	}

	if err := configs.readConfigs("./genie.d"); err != nil {
		return nil, err
	}

	for _, event := range configs.Events {
		if event.Filename == "" {
			continue
		}

		if tmpl, ok := configs.templates[event.Filename]; !ok {
			return nil, fmt.Errorf("requested template does not exist: %s\n%v", event.Filename, configs.listTemplates())
		} else {
			event.Raw = tmpl
		}
	}

	return configs, nil
}

func (c *Configs) listTemplates() string {
	list := make([]string, 0)
	for n, _ := range c.templates {
		list = append(list, n)
	}

	return strings.Join(list, ", ")
}

func (c *Configs) readConfigs(dir string) error {
	files, _ := filepath.Glob(fmt.Sprintf("%s/*.yaml", dir))
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		if err := yaml.Unmarshal(data, c); err != nil {
			return err
		}

		if err != nil {
			fmt.Printf("error loading file: %v", err)
			continue
		}
	}

	return nil
}

func (c *Configs) readTemplates(dir string) error {
	if c.templates == nil {
		c.templates = make(map[string]string)
	}

	files, _ := filepath.Glob(fmt.Sprintf("%s/*.tmpl", dir))
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		c.templates[file] = string(data)
	}

	return nil
}
