package events

import (
	"context"
	"sync"

	"ctx.sh/genie/pkg/sinks"
	"ctx.sh/strata"
	"github.com/go-logr/logr"
)

type ManagerOptions struct {
	Logger  logr.Logger
	Metrics *strata.Metrics
	Sinks   *sinks.Sinks
}

type Manager struct {
	// add logging
	events Events

	logger  logr.Logger
	metrics *strata.Metrics
	sinks   *sinks.Sinks

	sync.RWMutex
}

func NewManager(events Events, opts *ManagerOptions) *Manager {
	return &Manager{
		events:  events,
		logger:  opts.Logger,
		metrics: opts.Metrics,
		sinks:   opts.Sinks,
	}
}

func (m *Manager) Enable(names ...string) {
	for _, name := range names {
		if _, ok := m.events[name]; !ok {
			m.logger.Info("event not found", "name", name)
			continue
		}

		// TODO: move this to the event creation, it's harder to intuit
		// what's going on here.  I know that I didn't want to couple this
		// but it would make more sense here.  SendChannel probably stays here
		// so we can override.
		m.events[name].
			WithSendChannel(m.sinks.Stdout.SendChannel()).
			Enable()
	}
}

func (m *Manager) Start(ctx context.Context) error {
	m.Lock()
	defer m.Unlock()

	metrics := m.metrics.WithLabels("event")
	for n, g := range m.events {
		if !g.enabled {
			metrics.CounterInc("event_disabled", n)
			continue
		}

		if err := g.Start(ctx); err != nil {
			m.Stop()
			metrics.CounterInc("event_start_error", n)
			return err
		} else {
			metrics.CounterInc("event_start", n)
		}
	}

	<-ctx.Done()
	return nil
}

func (m *Manager) Stop() {
	m.Lock()
	defer m.Unlock()

	for n, g := range m.events {
		g.Stop()
		delete(m.events, n)
	}
}
