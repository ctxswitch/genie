package generator

import (
	"context"
	"sync"

	"ctx.sh/genie/pkg/sinks"
	"ctx.sh/genie/pkg/template"
	"ctx.sh/strata"
	"github.com/go-logr/logr"
)

type Manager struct {
	// add logging
	generators map[string]*Generator
	ctx        context.Context
	logger     logr.Logger
	metrics    *strata.Metrics
	sync.RWMutex
}

func NewManager(ctx context.Context) *Manager {
	return &Manager{
		generators: make(map[string]*Generator),
		ctx:        ctx,
	}
}

func (m *Manager) WithLogger(logger logr.Logger) *Manager {
	m.logger = logger
	return m
}

func (m *Manager) WithMetrics(metrics *strata.Metrics) *Manager {
	m.metrics = metrics
	return m
}

func (m *Manager) Add(name string, tmpl *template.Template, sink sinks.Sink) {
	g := NewGenerator(name, tmpl, sink).
		WithLogger(m.logger).
		WithMetrics(m.metrics)
	m.generators[name] = g
}

func (m *Manager) Start(ctx context.Context) error {
	m.Lock()
	defer m.Unlock()

	metrics := m.metrics.WithLabels("generator")
	for _, g := range m.generators {
		if err := g.Start(ctx); err != nil {
			m.Stop()
			metrics.CounterInc("generator_start_error", g.Name)
			return err
		} else {
			metrics.CounterInc("generator_start", g.Name)
		}
	}

	<-m.ctx.Done()
	return nil
}

func (m *Manager) Stop() {
	m.Lock()
	defer m.Unlock()

	for n, g := range m.generators {
		g.Stop()
		delete(m.generators, n)
	}
}
