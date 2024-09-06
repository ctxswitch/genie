package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"ctx.sh/genie/pkg/config"
	"ctx.sh/genie/pkg/events"
	"ctx.sh/strata"
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	Usage     = `generate [NAME...] [ARG...]`
	ShortDesc = `Start the generator for one or many events.`
	LongDesc  = `
Start the generator for one or many events. By default all configured
event generators will be run on startup. Individual events can be
specified by using the event name. Generate specific arguments can
be added as after the event name.`

	DefaultMetricsPort = 9090
)

// Generate is the command that starts the event generators.
type Generate struct {
	once        bool
	enableLogs  bool
	disableLogs bool
	sink        string
}

// NewGenerate returns a new Generate command.
func NewGenerate() *Generate {
	return &Generate{}
}

// RunE is the main entry point for the generate command which
// returns an error.
func (g *Generate) RunE(cmd *cobra.Command, args []string) error { // nolint:revive,funlen,gocognit
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Logging
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.RFC3339NanoTimeEncoder
	zapCfg := zap.Config{
		// TODO: make me configurable
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false,
		// TODO: enable this for debug mode
		DisableStacktrace: true,
		Sampling:          nil,
		Encoding:          "console",
		EncoderConfig:     encoderCfg,
		OutputPaths: []string{
			"stderr",
		},
		ErrorOutputPaths: []string{
			"stderr",
		},
	}
	zl := zap.Must(zapCfg.Build())
	defer zl.Sync() //nolint:errcheck

	// Metrics
	logger := zapr.NewLogger(zl)
	if (g.once && !g.enableLogs) || g.disableLogs {
		logger = logr.Discard()
	}

	metrics := strata.New(strata.MetricsOpts{
		// Enable later
		Logger:       logr.Discard(),
		Prefix:       []string{"genie"},
		PanicOnError: true,
	})

	var obs sync.WaitGroup
	obs.Add(1)
	go func() {
		defer obs.Done()
		err := metrics.Start(ctx, strata.ServerOpts{
			Port: DefaultMetricsPort,
		})
		if err != nil && err != http.ErrServerClosed {
			logger.Error(err, "metrics server start failed")
			os.Exit(1)
		}
	}()

	path := cmd.Flag("config").Value.String()
	cfg, err := config.Load(&config.LoadOptions{
		Paths:   []string{path},
		Logger:  logger,
		Metrics: metrics,
	})
	if err != nil {
		logger.Error(err, "unable to load configuration")
		return err
	}

	// TODO: Maybe we can add a start to the sinks struct that will take a
	// name and encasulate the logic below.
	sink, err := cfg.Sinks.Get(g.sink)
	if err != nil {
		logger.Error(err, "unable to get sink", "name", g.sink)
		os.Exit(1)
	}

	if err := sink.Init(); err != nil {
		logger.Error(err, "unable to initialize sink", "name", g.sink)
		os.Exit(1)
	}

	go sink.Start()

	evt := "all"
	if len(args) > 0 {
		if !cfg.HasEvent(args[0]) {
			logger.Error(fmt.Errorf("event %s not found", args[0]), "starting events")
			os.Exit(1)
		}
		evt = args[0]
	}

	if !cfg.HasEvents() {
		logger.Error(fmt.Errorf("no events were found in the configuration"), "starting events", "config", path)
		os.Exit(1)
	}

	logger.Info("starting event generators", "event", evt)

	manager := events.NewManager()
	// TODO: pull me out into another function.  Right now we start all
	// events that are configured, but we should also allow the user to
	// specify which events to start on the command line.  That will
	// happen later and should be pretty simple.
	for name, event := range cfg.Events {
		if evt != "all" && name != evt {
			continue
		}

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
	logger.Info("shutting down event generators")
	manager.Stop()

shutdown:
	sink.Stop()
	obs.Wait()
	return nil
}

// Command returns the cobra command for the generate command.
func (g *Generate) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   Usage,
		Short: ShortDesc,
		Long:  LongDesc,
		RunE:  g.RunE,
		Args:  cobra.MaximumNArgs(1),
	}
	cmd.PersistentFlags().StringVarP(&g.sink, "sink", "s", "stdout", "Override the configured sinks with the sinks provided.")
	cmd.PersistentFlags().BoolVar(&g.once, "run-once", false, "Run the generator one time and exit.  Logging is disabled by default, when run-once is enabled.  Use --enable-logs to enable.")
	cmd.PersistentFlags().BoolVar(&g.enableLogs, "enable-logs", false, "Enable log output.")
	cmd.PersistentFlags().BoolVar(&g.disableLogs, "disable-logs", false, "Disable log output.")
	cmd.MarkFlagsMutuallyExclusive("enable-logs", "disable-logs")
	return cmd
}
