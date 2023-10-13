package sinks

import (
	"fmt"
	"strings"
	"sync"

	"ctx.sh/strata"
	"github.com/go-logr/logr"
	"stvz.io/genie/pkg/resources"
	"stvz.io/genie/pkg/sinks/http"
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

	stdoutSink := stdout.New()
	stdoutSink.Init()

	return &Sinks{
		Stdout: stdoutSink,
		HTTP:   httpSinks,

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
}

func (s *Sinks) StopAll() {
	s.stopOnce.Do(func() {
		s.logger.Info("stopping stdout sink")
		s.Stdout.Stop()

		for n, v := range s.HTTP {
			s.logger.Info("stopping http sink", "name", n)
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
			return v.SendChannel(), nil
		}
		return nil, fmt.Errorf("sink not found: %s", name)
	default:
		return s.Stdout.SendChannel(), nil
	}
}
