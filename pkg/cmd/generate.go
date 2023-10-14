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

	g.logger.Info("starting sink", "args", args, "sink", g.sink)
	// TODO: Maybe we can add a start to the sinks struct that will take a
	// name and encasulate the logic below.
	sink, err := cfg.Sinks.Get(g.sink)
	if err != nil {
		g.logger.Error(err, "unable to get sink", "name", g.sink)
		os.Exit(1)
	}

	if err := sink.Init(); err != nil {
		g.logger.Error(err, "unable to initialize sink", "name", g.sink)
		os.Exit(1)
	}

	go sink.Start()

	g.logger.Info("starting event generators", "args", args, "sink", g.sink, "once", g.once)
	manager := events.NewManager()

	// TODO: pull me out into another function.  Right now we start all
	// events that are configured, but we should also allow the user to
	// specify which events to start on the command line.  That will
	// happen later and should be pretty simple.
	for name, event := range cfg.Events {
		g.logger.Info("starting event generator", "name", name)
		if g.once {
			event.Run(sink.SendChannel())
		} else {
			manager.Start(ctx, event, sink.SendChannel())
		}
	}

	if g.once {
		// TODO: I still don't like the goto.
		goto shutdown
	}

	<-ctx.Done()
	g.logger.Info("shutting down event generators")
	manager.Stop()

shutdown:
	sink.Stop()

	return nil
}

func (g *Generate) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   usage,
		Short: shortDesc,
		Long:  longDesc,
		RunE:  g.RunE,
	}
	cmd.PersistentFlags().StringVarP(&g.sink, "sink", "s", "stdout", "Override the configured sinks with the sinks provided.")
	cmd.PersistentFlags().BoolVar(&g.once, "run-once", false, "Run the generator one time and exit.  Logging is disabled by default, when run-once is enabled.  Use --enable-logs to enable.")
	cmd.PersistentFlags().BoolVar(&g.enableLogs, "enable-logs", false, "Enable log output.")
	cmd.PersistentFlags().BoolVar(&g.disableLogs, "disable-logs", false, "Disable log output.")
	cmd.MarkFlagsMutuallyExclusive("enable-logs", "disable-logs")
	return cmd
}
