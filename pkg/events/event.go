package events

import (
	"context"
	"sync"
	"time"

	"ctx.sh/strata"
	"github.com/go-logr/logr"
	"stvz.io/genie/pkg/resources"
	"stvz.io/genie/pkg/template"
	"stvz.io/genie/pkg/variables"
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
	rate       float64
	vars       *variables.Variables
	template   *template.Template

	logger  logr.Logger
	metrics *strata.Metrics

	resources *resources.Resources

	wg       sync.WaitGroup
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

		resources: opts.Resources,
		logger:    opts.Logger,
		metrics:   opts.Metrics,

		stopChan: make(chan struct{}),
	}, nil
}

func (e *Event) Run(sendChan chan<- []byte) {
	// TODO: there's another case that the send channel could be
	// closed, so check and if it is then close the stop channel.
	// If all the events have exited, then stop the manager and
	// and exit.
	// TODO: think about sending the context through and using it
	// as the parent of a timeout context to break out of the run
	// after a configurable amount of time.
	p := e.template.Execute(e.resources, e.vars)
	sendChan <- []byte(p)
}

func (e *Event) Start(ctx context.Context, sendChan chan<- []byte) {
	e.logger.Info("starting event generator", "count", e.generators)
	for i := 0; i < e.generators; i++ {
		e.wg.Add(1)
		go func() {
			defer e.wg.Done()
			e.generate(ctx, sendChan)
		}()
	}
}

func (e *Event) generate(ctx context.Context, sendChan chan<- []byte) {
	ticker := time.NewTicker(time.Duration(e.rate) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			e.metrics.CounterInc("run")
			e.Run(sendChan)
		case <-e.stopChan:
			return
		}
	}
}

func (e *Event) Stop() {
	e.stopOnce.Do(func() {
		close(e.stopChan)
	})

	e.wg.Wait()
}
