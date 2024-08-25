package sinks

import (
	"ctx.sh/genie/pkg/sinks/http"
	"ctx.sh/genie/pkg/sinks/kafka"
)

// Config is the top-level configuration for a collection of sinks.
type Config struct {
	HTTP  map[string]http.Config  `yaml:"http"`
	Kafka map[string]kafka.Config `yaml:"kafka"`
}
