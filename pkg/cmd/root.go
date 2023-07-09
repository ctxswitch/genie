package cmd

import (
	"context"
	"fmt"
	"os"

	"ctx.sh/genie/pkg/config"
	"ctx.sh/genie/pkg/generator"
	"ctx.sh/genie/pkg/resources"
	"ctx.sh/genie/pkg/sinks/stdout"
	"ctx.sh/genie/pkg/template"
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
		cfg, err := config.LoadAll("./genie.d")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		m := generator.NewManager(ctx)

		res, err := resources.FromConfig(cfg)

		for k, v := range cfg.Templates {
			// TODO: fix me.  move the map back into the config parsing
			vars := make(map[string]string)
			for _, v := range v.Vars {
				vars[v.Name] = v.Value
			}

			// TODO: Global variables (which we don't have yet), need to to be
			// merged in.
			tmpl := template.NewTemplate().WithResources(res).WithVars(vars)
			tmpl.Compile(v.Raw)

			m.Add(k, tmpl, &stdout.Stdout{})
		}

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
