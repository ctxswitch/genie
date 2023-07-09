package generator

import (
	"context"
	"sync"

	"ctx.sh/genie/pkg/sinks"
	"ctx.sh/genie/pkg/template"
)

type Manager struct {
	// add logging
	generators map[string]*Generator
	ctx        context.Context
	sync.RWMutex
}

func NewManager(ctx context.Context) *Manager {
	return &Manager{
		generators: make(map[string]*Generator),
		ctx:        ctx,
	}
}

func (m *Manager) Add(name string, tmpl *template.Template, sink sinks.Sink) {
	g := NewGenerator(tmpl, sink)
	m.generators[name] = g
}

func (m *Manager) Start(ctx context.Context) error {
	m.Lock()
	defer m.Unlock()

	for _, g := range m.generators {
		if err := g.Start(ctx); err != nil {
			m.Stop()
			return err
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
