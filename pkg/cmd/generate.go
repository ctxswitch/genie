package cmd

import (
	"context"
	"os"

	"ctx.sh/genie/pkg/config"
	"ctx.sh/genie/pkg/generator"
	"ctx.sh/genie/pkg/resources"
	"ctx.sh/genie/pkg/sinks"
	"ctx.sh/genie/pkg/template"
	"ctx.sh/genie/pkg/variables"
	"ctx.sh/strata"
	"github.com/go-logr/logr"
	"github.com/spf13/cobra"
)

var usage string = `generate events|event/NAME [ARG...]`
var shortDesc string = `Start the event generator(s)`
var longDesc string = `Start the event generator(s).  If 'events' is specified all configured
event generators will be run on startup.  Individual events can be specified by using
the 'event/NAME' syntax.  Generate specific arguments can be added as after the event
specifiers`

type Generate struct {
	logger  logr.Logger
	metrics *strata.Metrics
	cfg     config.ConfigBlock
	once    bool
	ctx     context.Context
	sinks   []string
}

func NewGenerate(opts *GlobalOpts) *Generate {
	return &Generate{
		logger:  opts.Logger,
		metrics: opts.Metrics,
		ctx:     opts.BaseContext,
		cfg:     opts.Config,
	}
}

func (g *Generate) RunE(cmd *cobra.Command, args []string) error {
	g.logger.Info("starting generator", "args", args, "config", g.cfg)

	g.logger.Info("loading resources")
	res, err := resources.Parse(g.cfg.Resources)
	if err != nil {
		g.logger.Error(err, "unable to load resources")
		os.Exit(1)
	}

	g.logger.Info("loading sinks")
	snks, err := sinks.ParseSinks(g.cfg.Sinks, res)
	if err != nil {
		g.logger.Error(err, "unable to load sinks")
		os.Exit(1)
	}

	m := generator.NewManager(g.ctx).
		WithLogger(g.logger).
		WithMetrics(g.metrics)

	g.logger.Info("loading events", "events", g.cfg.Events)

	for k, v := range g.cfg.Events {
		g.logger.Info("loading event", "event", k, "values", v)

		vars, err := variables.Parse(v.Vars)
		if err != nil {
			g.logger.Error(err, "event load failed, invalid variables", "event", k)
			continue
		}

		tmpl := template.NewTemplate().
			// TODO: configure paths (use several paths + configurable in priority order)
			WithPaths([]string{"./genie.d"}).
			WithResources(res).
			WithVariables(vars)

		if v.Template != "" {
			err = tmpl.CompileFrom(v.Template)
		} else {
			err = tmpl.Compile(v.Raw)
		}
		if err != nil {
			g.logger.Error(err, "event load failed, invalid template", "event", k)
			continue
		}

		if v.Sinks == nil {
			v.Sinks = append(v.Sinks, "stdout")
		}

		for n, s := range v.Sinks {
			sink, err := snks.Get(s)
			if err != nil {
				g.logger.Error(err, "sink load failed", "event", k, "sink", n)
				continue
			}

			// TODO: impl num generators
			m.Add(k, tmpl, sink)
		}
	}

	g.logger.Info("starting manager")
	if err := m.Start(g.ctx); err != nil {
		os.Exit(1)
	}

	return nil
}

func (g *Generate) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   usage,
		Args:                  cobra.MinimumNArgs(1),
		Short:                 shortDesc,
		Long:                  longDesc,
		RunE:                  g.RunE,
		DisableFlagsInUseLine: true,
	}
	cmd.Flags().StringArrayVarP(&g.sinks, "sinks", "s", []string{}, "Override the configured sinks with the sinks provided.")
	cmd.Flags().BoolVar(&g.once, "once", false, "Run the generator one time and exit.")
	cmd.Flags().SetInterspersed(false)

	return cmd
}
