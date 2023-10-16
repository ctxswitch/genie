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

// EventOptions are the options for an event generator.
type EventOptions struct {
	Logger    logr.Logger
	Metrics   *strata.Metrics
	Resources *resources.Resources
	Paths     []string
}

// Event describes an individual event generator.
// TODO: revice config to array of events and not map
type Event struct {
	name       string
	generators int
	interval   float64
	vars       *variables.Variables
	template   *template.Template

	logger  logr.Logger
	metrics *strata.Metrics

	resources *resources.Resources

	wg       sync.WaitGroup
	stopChan chan struct{}
	stopOnce sync.Once
}

// ParseEvent parses an event config and returns an event generator.  If the
// template does not compile then an error is returned.
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
		interval:   cfg.IntervalSeconds,
		vars:       vars,
		template:   tmpl,

		resources: opts.Resources,
		logger:    opts.Logger,
		metrics:   opts.Metrics,

		stopChan: make(chan struct{}),
	}, nil
}

// Run executes the template and sends the result to the send channel.
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

// Start starts the event generator.
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

// generate is the main loop for an event generator.  It will run the event
// at the configured interval.
func (e *Event) generate(ctx context.Context, sendChan chan<- []byte) { // nolint:unparam,revive
	ticker := time.NewTicker(time.Duration(e.interval) * time.Second)
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

// Stop stops the event generator.
func (e *Event) Stop() {
	e.stopOnce.Do(func() {
		close(e.stopChan)
	})

	e.wg.Wait()
}
