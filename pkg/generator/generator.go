package generator

import (
	"context"
	"sync"
	"time"

	"ctx.sh/genie/pkg/sinks"
	"ctx.sh/genie/pkg/sinks/stdout"
	"ctx.sh/genie/pkg/template"
)

type Generator struct {
	Template template.Template
	Sink     sinks.Sink
	Rate     time.Duration

	stopChan chan struct{}
	stopOnce sync.Once
}

func NewGenerator(tmpl template.Template) *Generator {
	return &Generator{
		Template: tmpl,
		// fix me
		Sink: &stdout.Stdout{},
		Rate: time.Second,
	}
}

func (g *Generator) run(ctx context.Context) {
	ticker := time.NewTicker(g.Rate)
	defer ticker.Stop()

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
