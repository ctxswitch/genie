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

	stopChan chan struct{}
	stopOnce sync.Once
}

func NewGenerator(name string, tmpl *template.Template, sink sinks.Sink) *Generator {
	return &Generator{
		Template: tmpl,
		// fix me
		Sink: sink,
		Rate: time.Second,
		Name: name,
	}
}

func (g *Generator) WithLogger(logger logr.Logger) *Generator {
	g.logger = logger.V(3).WithValues("sink", g.Sink.Name(), "template", g.Name)
	return g
}

func (g *Generator) WithMetrics(metrics *strata.Metrics) *Generator {
	g.metrics = metrics.WithPrefix("generator").WithLabels("name", "sink")
	return g
}

func (g *Generator) run(ctx context.Context) {
	ticker := time.NewTicker(g.Rate)
	defer ticker.Stop()

	// Connect needs a error returned and then we need to
	// figure out what the default behavior is.
	g.Sink.Connect()

	for {
		select {
		case <-ticker.C:
			g.metrics.CounterInc("run", g.Name, g.Sink.Name())
			p := g.Template.Execute()
			err := g.Sink.Send([]byte(p))
			if err != nil {
				g.metrics.CounterInc("send_error", g.Name, g.Sink.Name())
				g.logger.Error(err, "send error", g.Name, g.Sink.Name())
			} else {
				g.metrics.CounterInc("send_success", g.Name, g.Sink.Name())
			}
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
