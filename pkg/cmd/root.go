package cmd

import (
	"fmt"
	"os"

	"ctx.sh/genie/pkg/build"
	"github.com/spf13/cobra"
)

// Root is the root command for the genie CLI.
type Root struct{}

// NewRoot returns a new root command.
func NewRoot() *Root {
	return &Root{}
}

// Execute runs the root command.
func (r *Root) Execute() {
	if err := r.Command().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Command returns the root command.
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
			_ = cmd.Help()
		},
	}

	rootCmd.AddCommand(NewGenerate().Command())
	rootCmd.PersistentFlags().StringP("config", "c", "./genie.d", "config file (default is $HOME/.genie.yaml)")

	return rootCmd
}
