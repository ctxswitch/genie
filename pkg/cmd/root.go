package cmd

import (
	"context"
	"fmt"
	"os"

	"ctx.sh/apex"
	"ctx.sh/genie/pkg/config"
	"github.com/go-logr/logr"
	"github.com/spf13/cobra"
)

var (
	version string = "v0.1.0"
)

type GlobalOpts struct {
	Logger      logr.Logger
	Metrics     *apex.Metrics
	BaseContext context.Context
	CancelFunc  context.CancelFunc
	Config      config.ConfigBlock
}

type Root struct {
	logger  logr.Logger
	metrics *apex.Metrics
	ctx     context.Context
	cancel  context.CancelFunc
	config  config.ConfigBlock
}

func NewRoot(opts *GlobalOpts) *Root {
	return &Root{
		logger:  opts.Logger,
		metrics: opts.Metrics,
		ctx:     opts.BaseContext,
		cancel:  opts.CancelFunc,
		config:  opts.Config,
	}
}

func (r *Root) Execute() {
	if err := r.Command().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func (r *Root) Command() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "genie",
		Short: "Genie is a event payload generator.",
		Long: `An event payload generator used for interacting with services.  It provides
				a flexible templating solution to build out predictable payloads matching
				values for testing and validation`,
		Version: version,
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: output help
		},
	}

	// TODO: This needs to change.  Just keep the options around to repass
	opts := &GlobalOpts{
		Logger:      r.logger,
		Metrics:     r.metrics,
		BaseContext: r.ctx,
		CancelFunc:  r.cancel,
		Config:      r.config,
	}

	rootCmd.AddCommand(NewGenerate(opts).Command())

	return rootCmd
}
