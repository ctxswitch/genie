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
	Paths     *Paths               `yaml:"paths"`
	Templates map[string]*Template `yaml:"events"`
	Resources *Resources           `yaml:"resources"`
	Sinks     *Sinks               `yaml:"sinks"`
}

func LoadAll(dir string) (*Configs, error) {
	configs := &Configs{}

	// Fix me to add custom directories.  At the very least, I want to split out templates
	// from normal config. Not sure about the others.
	templates, err := configs.readTemplates(dir)
	if err != nil {
		return nil, err
	}

	// Fix me to add custom directories.
	if err := configs.readConfigs(dir); err != nil {
		return nil, err
	}

	for _, tmpl := range configs.Templates {
		if tmpl.Template == "" {
			continue
		}

		path := fmt.Sprintf("%s/%s", dir, tmpl.Template)
		abs, err := filepath.Abs(path)
		if err != nil {
			return nil, err
		}

		if content, ok := templates[abs]; !ok {
			return nil, fmt.Errorf("requested template does not exist: %s\n%v", abs, configs.listKeys(templates))
		} else {
			tmpl.Raw = content
		}
	}

	return configs, nil
}

func (c *Configs) listKeys(m map[string]string) string {
	list := make([]string, 0)
	for n := range m {
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

func (c *Configs) readTemplates(dir string) (map[string]string, error) {
	templates := make(map[string]string)
	if templates == nil {
		templates = make(map[string]string)
	}

	files, _ := filepath.Glob(fmt.Sprintf("%s/*.tmpl", dir))
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}

		templates[file] = string(data)
	}

	return templates, nil
}
