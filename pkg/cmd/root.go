package cmd

import (
	"context"
	"fmt"
	"os"

	"ctx.sh/genie/pkg/config"
	"ctx.sh/genie/pkg/generator"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "genie",
	Short: "Genie is a event payload generator.",
	Long: `An event payload generator used for interacting with services.  It provides
			a flexible templating solution to build out predictable payloads matching
			values for testing and validation`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		// load config
		_, err := config.LoadAll("./genie.d")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// make resources

		// merge vars

		// make templates

		// sink := stdout.Stdout{}

		// create event generators
		m := generator.NewManager(ctx)

		if err := m.Start(ctx); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
