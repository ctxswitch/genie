package sinks

import (
	"fmt"
	"strings"

	"ctx.sh/genie/pkg/config"
	"ctx.sh/genie/pkg/resources"
	"ctx.sh/genie/pkg/sinks/http"
	"ctx.sh/genie/pkg/sinks/stdout"
)

type Sink interface {
	Send([]byte) error
	Connect()
	Init()
	Name() string
}

type Sinks struct {
	HTTP map[string]Sink
}

func ParseSinks(block config.SinksBlock, res *resources.Resources) (*Sinks, error) {
	httpSinks, err := parseHttpSinks(block, res)
	if err != nil {
		return nil, err
	}

	return &Sinks{
		HTTP: httpSinks,
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

			sinks[k] = sink
		}
	}

	return sinks, nil
}

func (s *Sinks) Get(sink string) (Sink, error) {
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
			return v, nil
		}
		return nil, fmt.Errorf("sink not found: %s", name)
	default:
		return &stdout.Stdout{}, nil
	}
}
