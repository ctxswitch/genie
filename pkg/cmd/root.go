package cmd

import (
	"context"
	"fmt"
	"os"

	"ctx.sh/genie/pkg/build"
	"ctx.sh/strata"
	"github.com/go-logr/logr"
	"github.com/spf13/cobra"
)

// GlobalOpts are the global options for all commands
type GlobalOpts struct {
	Logger      logr.Logger
	Metrics     *strata.Metrics
	BaseContext context.Context
	CancelFunc  context.CancelFunc
}

// Root is the root command for the genie CLI
type Root struct {
	logger  logr.Logger
	metrics *strata.Metrics
	ctx     context.Context
	cancel  context.CancelFunc
}

// NewRoot returns a new root command
func NewRoot(opts *GlobalOpts) *Root {
	return &Root{
		logger:  opts.Logger,
		metrics: opts.Metrics,
		ctx:     opts.BaseContext,
		cancel:  opts.CancelFunc,
	}
}

// Execute runs the root command
func (r *Root) Execute() {
	if err := r.Command().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Command returns the root command
func (r *Root) Command() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "genie",
		Short: "Genie is a event payload generator.",
		Long: `
An event payload generator used for interacting with services.  It provides
a flexible templating solution to build out predictable payloads matching
values for the testing and validation of event pipelines.`,
		Version: build.Version,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	// TODO: This needs to change.  Just keep the options around to repass
	opts := &GlobalOpts{
		Logger:      r.logger,
		Metrics:     r.metrics,
		BaseContext: r.ctx,
	}

	rootCmd.AddCommand(NewGenerate(opts).Command())
	rootCmd.PersistentFlags().StringP("config", "c", "./genie.d", "config file (default is $HOME/.genie.yaml)")

	return rootCmd
}
