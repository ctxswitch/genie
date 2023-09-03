package generator

import (
	"context"
	"sync"
	"time"

	"ctx.sh/genie/pkg/sinks"
	"ctx.sh/genie/pkg/template"
	"ctx.sh/strata"
	"github.com/go-logr/logr"
)

type Generator struct {
	Template *template.Template
	Sink     sinks.Sink
	Rate     time.Duration
	Name     string

	logger  logr.Logger
	metrics *strata.Metrics

	sendChan chan<- []byte
	stopChan chan struct{}
	stopOnce sync.Once
}

func NewGenerator(name string) *Generator {
	return &Generator{
		// TODO: make me configurable
		Rate: time.Second,
		Name: name,
	}
}

func (g *Generator) WithTemplate(tmpl *template.Template) *Generator {
	g.Template = tmpl
	return g
}

func (g *Generator) WithSendChannel(send chan<- []byte) *Generator {
	g.sendChan = send
	return g
}

func (g *Generator) WithLogger(logger logr.Logger) *Generator {
	g.logger = logger.V(3).WithValues("template", g.Name)
	return g
}

func (g *Generator) WithMetrics(metrics *strata.Metrics) *Generator {
	g.metrics = metrics.WithPrefix("generator")
	return g
}

func (g *Generator) run(ctx context.Context) {
	ticker := time.NewTicker(g.Rate)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			g.metrics.CounterInc("run")
			p := g.Template.Execute()
			g.sendChan <- []byte(p)
		case <-g.stopChan:
			return
		case <-ctx.Done():
			return
		}
	}
}

func (g *Generator) Start(ctx context.Context) error {
	go g.run(ctx)
	return nil
}

func (g *Generator) Stop() {
	g.stopOnce.Do(func() {
		close(g.stopChan)
	})
}
