package events

import (
	"ctx.sh/strata"
	"github.com/go-logr/logr"
	"stvz.io/genie/pkg/resources"
)

// Options are the options for a collection of configured event generators.
type Options struct {
	Logger    logr.Logger
	Metrics   *strata.Metrics
	Resources *resources.Resources
	Paths     []string
}

// Events is a collection of configured event generators.
// TODO: rethink this as an array.
type Events map[string]*Event

func Parse(cfg Config, opts *Options) (Events, error) {
	events := make(Events)
	for _, evt := range cfg {
		event, err := ParseEvent(evt, &EventOptions{
			Logger:    opts.Logger.WithValues("event", evt.Name),
			Metrics:   opts.Metrics.WithPrefix("event"),
			Resources: opts.Resources,
			Paths:     opts.Paths,
		})
		if err != nil {
			return nil, err
		}
		events[evt.Name] = event
	}
	return events, nil
}
