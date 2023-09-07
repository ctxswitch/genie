package cmd

import (
	"context"
	"os"

	"ctx.sh/genie/pkg/config"
	"ctx.sh/genie/pkg/events"
	"ctx.sh/strata"
	"github.com/go-logr/logr"
	"github.com/spf13/cobra"
)

var usage string = `generate [NAME...] [ARG...]`
var shortDesc string = `Start the generator for one or many events.`
var longDesc string = `Start the generator for one or many events. By default all configured
event generators will be run on startup. Individual events can be specified by using the event
name. Generate specific arguments can be added as after the event name.`

type Generate struct {
	logger  logr.Logger
	metrics *strata.Metrics
	once    bool
	ctx     context.Context
	sink    string
}

func NewGenerate(opts *GlobalOpts) *Generate {
	return &Generate{
		logger:  opts.Logger,
		metrics: opts.Metrics,
		ctx:     opts.BaseContext,
	}
}

func (g *Generate) RunE(cmd *cobra.Command, args []string) error {
	// TODO: I'm probably going to move the subsequential parsing into
	// the objects into the load stage.
	path := cmd.Flag("config").Value.String()
	cfg, err := config.Load(&config.LoadOptions{
		Paths:   []string{path},
		Logger:  g.logger,
		Metrics: g.metrics,
	})
	if err != nil {
		g.logger.Error(err, "unable to load configuration")
		os.Exit(1)
	}

	g.logger.Info("starting event generators", "args", args, "sink", g.sink, "once", g.once)

	manager := events.NewManager(cfg.Events, &events.ManagerOptions{
		Logger:  g.logger,
		Metrics: g.metrics,
		Sinks:   cfg.Sinks,
	})

	var evts []string
	if len(args) > 0 {
		evts = args
	} else {
		evts = cfg.Events.Names()
	}

	manager.Enable(evts...)

	g.logger.Info("starting sinks")
	cfg.Sinks.StartAll(g.ctx)

	g.logger.Info("starting manager")
	if err := manager.Start(g.ctx); err != nil {
		os.Exit(1)
	}

	return nil
}

func (g *Generate) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   usage,
		Short: shortDesc,
		Long:  longDesc,
		RunE:  g.RunE,
	}
	cmd.PersistentFlags().StringVarP(&g.sink, "sink", "s", "", "Override the configured sinks with the sinks provided.")
	cmd.PersistentFlags().BoolVar(&g.once, "run-once", false, "Run the generator one time and exit.")

	return cmd
}
