package sinks

import (
	"fmt"
	"strings"

	"ctx.sh/strata"
	"github.com/go-logr/logr"
	"stvz.io/genie/pkg/sinks/http"
	"stvz.io/genie/pkg/sinks/kafka"
	"stvz.io/genie/pkg/sinks/stdout"
)

type Sink interface {
	Init() error
	SendChannel() chan<- []byte
	Start()
	Stop()
}

type Options struct {
	Logger  logr.Logger
	Metrics *strata.Metrics
}

type Sinks struct {
	Stdout Sink
	HTTP   map[string]Sink
	Kafka  map[string]Sink

	logger  logr.Logger
	metrics *strata.Metrics
}

func New(cfg Config, opts *Options) *Sinks {
	httpSinks := parseHTTPSinks(cfg, opts)
	kafkaSinks := parseKafkaSinks(cfg, opts)

	stdoutSink := stdout.New()
	_ = stdoutSink.Init()

	return &Sinks{
		Stdout: stdoutSink,
		HTTP:   httpSinks,
		Kafka:  kafkaSinks,

		logger:  opts.Logger.WithName("sinks"),
		metrics: opts.Metrics.WithPrefix("sinks"),
	}
}

func parseHTTPSinks(cfg Config, opts *Options) map[string]Sink {
	sinks := make(map[string]Sink)

	for k, v := range cfg.HTTP {
		sink := http.New(v, &http.Options{
			Logger:  opts.Logger.WithValues("type", "http", "name", k),
			Metrics: opts.Metrics.WithPrefix("http"),
		})
		sinks[k] = sink
	}

	return sinks
}

func parseKafkaSinks(cfg Config, opts *Options) map[string]Sink {
	sinks := make(map[string]Sink)

	for k, v := range cfg.Kafka {
		sink := kafka.New(v, &kafka.Options{
			Logger:  opts.Logger.WithValues("type", "kafka", "name", k),
			Metrics: opts.Metrics.WithPrefix("kafka"),
		})
		sinks[k] = sink
	}

	return sinks
}

// TODO: no more passing sinks around, we just pass the send channel back.
func (s *Sinks) Get(sink string) (Sink, error) {
	var kind, name string
	var snk Sink
	var ok bool

	parts := strings.SplitN(sink, ".", 2)
	kind = parts[0]
	if len(parts) == 2 {
		name = parts[1]
	} else {
		name = "stdout"
	}

	// TODO: add better checks to give back errors on invalid sinks

	switch kind {
	case "http":
		if snk, ok = s.HTTP[name]; !ok {
			return nil, fmt.Errorf("sink not found: %s", name)
		}
	case "kafka":
		if snk, ok = s.Kafka[name]; !ok {
			return nil, fmt.Errorf("sink not found: %s", name)
		}
	default:
		snk = s.Stdout
	}

	s.logger.Info("using sink", "kind", kind, "name", name)
	return snk, nil
}
