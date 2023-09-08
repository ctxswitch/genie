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
	RunOnce bool
}

type Manager struct {
	// add logging
	events Events
	once   bool

	logger  logr.Logger
	metrics *strata.Metrics
	sinks   *sinks.Sinks

	wg       sync.WaitGroup
	stopOnce sync.Once

	sync.RWMutex
}

func NewManager(events Events, opts *ManagerOptions) *Manager {
	return &Manager{
		events:  events,
		logger:  opts.Logger,
		metrics: opts.Metrics,
		sinks:   opts.Sinks,
		once:    opts.RunOnce,
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

func (m *Manager) RunOnce() {
	for _, g := range m.events {
		if !g.enabled {
			continue
		}

		g.run()
	}
}

func (m *Manager) StartAll() {
	m.Lock()
	defer m.Unlock()

	ctx := context.Background()
	metrics := m.metrics.WithLabels("event")
	for n, g := range m.events {
		if !g.enabled {
			metrics.CounterInc("event_disabled", n)
			continue
		}

		m.wg.Add(1)
		go func(g *Event) {
			defer m.wg.Done()
			g.Start(ctx)
		}(g)
	}
}

func (m *Manager) StopAll() {
	m.Lock()
	defer m.Unlock()

	m.stopOnce.Do(func() {
		for _, g := range m.events {
			g.Stop()
		}

		m.wg.Wait()
	})
}
