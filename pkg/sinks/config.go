package sinks

import (
	"stvz.io/genie/pkg/sinks/http"
	"stvz.io/genie/pkg/sinks/kafka"
)

type Config struct {
	Http  map[string]http.Config  `yaml:"http"`
	Kafka map[string]kafka.Config `yaml:"kafka"`
}
