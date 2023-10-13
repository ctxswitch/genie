package sinks

import "stvz.io/genie/pkg/sinks/http"

type Config struct {
	Http map[string]http.Config `yaml:"http"`
}
