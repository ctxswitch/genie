package sinks

import (
	"stvz.io/genie/pkg/sinks/http"
	"stvz.io/genie/pkg/sinks/kafka"
)

// Config is the top-level configuration for a collection of sinks.
type Config struct {
	HTTP  map[string]http.Config  `yaml:"http"`
	Kafka map[string]kafka.Config `yaml:"kafka"`
}
