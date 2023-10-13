package cmd

import (
	"context"
	"os"

	"ctx.sh/strata"
	"github.com/go-logr/logr"
	"github.com/spf13/cobra"
	"stvz.io/genie/pkg/config"
	"stvz.io/genie/pkg/events"
)

var usage string = `generate [NAME...] [ARG...]`
var shortDesc string = `Start the generator for one or many events.`
var longDesc string = `Start the generator for one or many events. By default all configured
event generators will be run on startup. Individual events can be specified by using the event
name. Generate specific arguments can be added as after the event name.`

type Generate struct {
	logger      logr.Logger
	metrics     *strata.Metrics
	once        bool
	enableLogs  bool
	disableLogs bool
	ctx         context.Context
	sink        string
}

func NewGenerate(opts *GlobalOpts) *Generate {
	return &Generate{
		logger:  opts.Logger,
		metrics: opts.Metrics,
		ctx:     opts.BaseContext,
	}
}

func (g *Generate) RunE(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithCancel(g.ctx)
	defer cancel()

	if (g.once && !g.enableLogs) || g.disableLogs {
		g.logger = logr.Discard()
	}

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
		RunOnce: g.once,
	})

	var evts []string
	if len(args) > 0 {
		evts = args
	} else {
		evts = cfg.Events.Names()
	}

	manager.Enable(g.sink, evts...)

	g.logger.Info("starting sinks")
	// TODO: Only start sinks that are in use.
	cfg.Sinks.StartAll()

	// Wait for context if we are not running once.
	if g.once {
		manager.RunOnce()
		goto shutdown
	}

	g.logger.Info("starting events manager")
	manager.StartAll()
	<-ctx.Done()

	g.logger.Info("shutting down event generators")
	manager.StopAll()

shutdown:
	cfg.Sinks.StopAll()

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
	cmd.PersistentFlags().BoolVar(&g.once, "run-once", false, "Run the generator one time and exit.  Logging is disabled by default, when run-once is enabled.  Use --enable-logs to enable.")
	cmd.PersistentFlags().BoolVar(&g.enableLogs, "enable-logs", false, "Enable log output.")
	cmd.PersistentFlags().BoolVar(&g.disableLogs, "disable-logs", false, "Disable log output.")
	cmd.MarkFlagsMutuallyExclusive("enable-logs", "disable-logs")
	return cmd
}
