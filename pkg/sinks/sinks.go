package sinks

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"ctx.sh/genie/pkg/config"
	"ctx.sh/genie/pkg/resources"
	"ctx.sh/genie/pkg/sinks/http"
	"ctx.sh/genie/pkg/sinks/stdout"
)

type Sink interface {
	Init() error
	SendChannel() chan<- []byte
	Start(context.Context)
	Stop()
}

type Sinks struct {
	Stdout Sink
	HTTP   map[string]Sink
	wg     sync.WaitGroup
}

func ParseSinks(block config.SinksBlock, res *resources.Resources) (*Sinks, error) {
	httpSinks, err := parseHttpSinks(block, res)
	if err != nil {
		return nil, err
	}

	stdoutSink := stdout.New()
	stdoutSink.Init()

	return &Sinks{
		Stdout: stdoutSink,
		HTTP:   httpSinks,
	}, nil
}

func parseHttpSinks(out config.SinksBlock, res *resources.Resources) (map[string]Sink, error) {
	sinks := make(map[string]Sink)

	for k, v := range out.Http {
		sink := http.New().
			WithURL(v.Url).
			WithMethod(v.Method)

		for _, h := range v.Headers {
			// this is where I'm going to need my resources for the conversion.
			// there may be a way to decouple this from the config process.
			if h.Resource != "" {
				// Split the resources
				parts := strings.SplitN(h.Resource, ".", 2)
				if len(parts) != 2 {
					return nil, fmt.Errorf("invalid resource provided: %s", h.Resource)
				}
				r, err := res.Get(parts[0], parts[1])
				if err != nil {
					return nil, fmt.Errorf("resource does not exist: %s", h.Resource)
				}
				sink.WithHeader(h.Name, r)
			} else {
				sink.WithHeader(h.Name, h.Value)
			}
		}
		err := sink.Init()
		if err != nil {
			// log
			continue
		}
		sinks[k] = sink
	}

	return sinks, nil
}

func (s *Sinks) StartAll(ctx context.Context) {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.Stdout.Start(ctx)
	}()

	for _, v := range s.HTTP {
		s.wg.Add(1)
		go func(v Sink) {
			defer s.wg.Done()
			v.Start(ctx)
		}(v)
	}
}

func (s *Sinks) StopAll() {
	s.Stdout.Stop()
	for _, v := range s.HTTP {
		v.Stop()
	}

	s.wg.Wait()
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
