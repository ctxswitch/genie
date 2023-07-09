package generator

import (
	"context"
	"sync"
	"time"

	"ctx.sh/genie/pkg/sinks"
	"ctx.sh/genie/pkg/template"
)

type Generator struct {
	Template *template.Template
	Sink     sinks.Sink
	Rate     time.Duration

	stopChan chan struct{}
	stopOnce sync.Once
}

func NewGenerator(tmpl *template.Template, sink sinks.Sink) *Generator {
	return &Generator{
		Template: tmpl,
		// fix me
		Sink: sink,
		Rate: time.Second,
	}
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
			g.Sink.Send([]byte(g.Template.Execute()))
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
