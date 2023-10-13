package sinks

import (
	"fmt"
	"strings"
	"sync"

	"ctx.sh/strata"
	"github.com/go-logr/logr"
	"stvz.io/genie/pkg/resources"
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

	logger   logr.Logger
	metrics  *strata.Metrics
	wg       sync.WaitGroup
	stopOnce sync.Once
}

func Parse(cfg Config, res *resources.Resources, opts *Options) (*Sinks, error) {
	httpSinks, err := parseHttpSinks(cfg, res, opts)
	if err != nil {
		return nil, err
	}

	kafkaSinks, err := parseKafkaSinks(cfg, res, opts)
	if err != nil {
		return nil, err
	}

	stdoutSink := stdout.New()
	stdoutSink.Init()

	return &Sinks{
		Stdout: stdoutSink,
		HTTP:   httpSinks,
		Kafka:  kafkaSinks,

		logger:  opts.Logger.WithName("sinks"),
		metrics: opts.Metrics.WithPrefix("sinks"),
	}, nil
}

func parseHttpSinks(cfg Config, res *resources.Resources, opts *Options) (map[string]Sink, error) {
	sinks := make(map[string]Sink)

	for k, v := range cfg.Http {
		sink := http.New(v, &http.HTTPOptions{
			Logger:  opts.Logger.WithValues("type", "http", "name", k),
			Metrics: opts.Metrics.WithPrefix("http"),
		})

		err := sink.Init()
		if err != nil {
			return nil, err
		}
		sinks[k] = sink
	}

	return sinks, nil
}

func parseKafkaSinks(cfg Config, res *resources.Resources, opts *Options) (map[string]Sink, error) {
	sinks := make(map[string]Sink)

	for k, v := range cfg.Kafka {
		sink := kafka.New(v, &kafka.KafkaOpts{
			Logger:  opts.Logger.WithValues("type", "kafka", "name", k),
			Metrics: opts.Metrics.WithPrefix("kafka"),
		})

		err := sink.Init()
		if err != nil {
			return nil, err
		}
		sinks[k] = sink
	}

	return sinks, nil
}

func (s *Sinks) StartAll() {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.Stdout.Start()
	}()

	for _, v := range s.HTTP {
		s.wg.Add(1)
		go func(v Sink) {
			defer s.wg.Done()
			v.Start()
		}(v)
	}

	for _, v := range s.Kafka {
		s.wg.Add(1)
		go func(v Sink) {
			defer s.wg.Done()
			v.Start()
		}(v)
	}
}

func (s *Sinks) StopAll() {
	s.stopOnce.Do(func() {
		s.logger.Info("stopping stdout sink")
		s.Stdout.Stop()

		for n, v := range s.HTTP {
			s.logger.Info("stopping http sink", "name", n)
			v.Stop()
		}

		for n, v := range s.Kafka {
			s.logger.Info("stopping kafka sink", "name", n)
			v.Stop()
		}

		s.wg.Wait()
	})
}

// TODO: no more passing sinks around, we just pass the send channel back.
func (s *Sinks) Get(sink string) (chan<- []byte, error) {
	var kind, name string

	parts := strings.SplitN(sink, ".", 2)
	kind = parts[0]
	if len(parts) == 2 {
		name = parts[1]
	}

	// TODO: add better checks to give back errors on invalid sinks

	switch kind {
	case "http":
		if v, ok := s.HTTP[name]; ok {
			s.logger.Info("using http sink", "name", name)
			return v.SendChannel(), nil
		}
		return nil, fmt.Errorf("sink not found: %s", name)
	case "kafka":
		if v, ok := s.Kafka[name]; ok {
			s.logger.Info("using kafka sink", "name", name)
			return v.SendChannel(), nil
		}
		return nil, fmt.Errorf("sink not found: %s", name)
	}

	s.logger.Info("using default sink: stdout")
	return s.Stdout.SendChannel(), nil
}
