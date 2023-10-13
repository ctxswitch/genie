package events

import (
	"context"
	"sync"

	"ctx.sh/strata"
	"github.com/go-logr/logr"
	"stvz.io/genie/pkg/sinks"
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

// TODO: I don't like this.  Remove the sink discovery from the manager
// and pass in the sinks that are specified on the command line, taking
// out all of the logic for sink configuration in the config file to simplify
// use.
func (m *Manager) Enable(sink string, names ...string) {
	for _, name := range names {
		if _, ok := m.events[name]; !ok {
			m.logger.Info("event not found", "name", name)
			continue
		}

		var sendChan chan<- []byte
		sendChan, err := m.sinks.Get(sink)
		if err != nil {
			sendChan = m.sinks.Stdout.SendChannel()
		}

		m.events[name].
			WithSendChannel(sendChan).
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
		// TODO: not a fan of this, we should initially filter out the
		// disabled events.
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
			// TODO: have a started map that we can check here
			if g.enabled {
				g.logger.Info("stopping event generator", "name", g.name)
				g.Stop()
			}
		}

		m.wg.Wait()
	})
}
