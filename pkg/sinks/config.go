package sinks

import "ctx.sh/genie/pkg/sinks/http"

type Config struct {
	Http map[string]http.Config `yaml:"http"`
}
