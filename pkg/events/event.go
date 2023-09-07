package events

import (
	"context"
	"sync"
	"time"

	"ctx.sh/genie/pkg/resources"
	"ctx.sh/genie/pkg/template"
	"ctx.sh/genie/pkg/variables"
	"ctx.sh/strata"
	"github.com/go-logr/logr"
)

type EventOptions struct {
	Logger    logr.Logger
	Metrics   *strata.Metrics
	Resources *resources.Resources
	Paths     []string
}

// TODO: revice config to array of events and not map
type Event struct {
	name       string
	generators int
	enabled    bool
	rate       float64
	vars       *variables.Variables
	template   *template.Template

	logger  logr.Logger
	metrics *strata.Metrics

	resources *resources.Resources

	sendChan chan<- []byte
	stopChan chan struct{}
	stopOnce sync.Once
}

func ParseEvent(cfg EventConfig, opts *EventOptions) (*Event, error) {
	tmpl := template.NewTemplate().WithPaths(opts.Paths)

	var err error
	if cfg.Template != "" {
		err = tmpl.CompileFrom(cfg.Template)
	} else {
		err = tmpl.Compile(cfg.Raw)
	}
	if err != nil {
		return nil, err
	}

	vars, err := variables.Parse(cfg.Vars)
	if err != nil {
		return nil, err
	}

	return &Event{
		name:       cfg.Name,
		generators: cfg.Generators,
		rate:       cfg.RateSeconds,
		vars:       vars,
		template:   tmpl,
		enabled:    false,

		resources: opts.Resources,
		logger:    opts.Logger,
		metrics:   opts.Metrics,
	}, nil
}

func (e *Event) Enable() {
	e.enabled = true
}

func (e *Event) WithSendChannel(send chan<- []byte) *Event {
	e.sendChan = send
	return e
}

// func (e *Event) WithLogger(logger logr.Logger) *Event {
// 	e.logger = logger.V(3).WithValues("event", e.name)
// 	return e
// }

// func (e *Event) WithMetrics(metrics *strata.Metrics) *Event {
// 	e.metrics = metrics.WithPrefix("event")
// 	return e
// }

// func (e *Event) WithResources(res *resources.Resources) *Event {
// 	e.resources = res
// 	return e
// }

func (e *Event) run(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(e.rate) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			e.metrics.CounterInc("run")
			p := e.template.Execute(e.resources, e.vars)
			e.sendChan <- []byte(p)
		case <-e.stopChan:
			return
		case <-ctx.Done():
			return
		}
	}
}

func (e *Event) Start(ctx context.Context) error {
	go e.run(ctx)
	return nil
}

func (e *Event) Stop() {
	e.stopOnce.Do(func() {
		close(e.stopChan)
	})
}
