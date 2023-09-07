package config

import (
	"fmt"
	"os"
	"path/filepath"

	"ctx.sh/genie/pkg/events"
	"ctx.sh/genie/pkg/resources"
	"ctx.sh/genie/pkg/sinks"
	"ctx.sh/strata"
	"github.com/go-logr/logr"
	"gopkg.in/yaml.v2"
)

type LoadOptions struct {
	Logger  logr.Logger
	Metrics *strata.Metrics
	Paths   []string
}

type Config struct {
	Events    events.Events
	Resources *resources.Resources
	Sinks     *sinks.Sinks
}

// TODO: move sinks and events?  Events won't actually need anything special, just
// the config parsing.  Sinks will be set up like resources.
func Load(opts *LoadOptions) (*Config, error) {
	var config struct {
		Events    events.Config    `yaml:"events"`
		Resources resources.Config `yaml:"resources"`
		Sinks     sinks.Config     `yaml:"sinks"`
	}

	for _, path := range opts.Paths {
		// TODO: right now we are only supporting yaml extensions.  This is mostly
		// due to Glob not supporting anything other than single character wildcards.
		// (i.e. patterns like `*.y[a]?ml` are not supported).  It's not a huge deal
		// but it is something to keep in mind.  We can always add support for the other
		// extensions later if we need to.
		files, err := filepath.Glob(fmt.Sprintf("%s/*.yaml", path))
		if err != nil {
			opts.Logger.Error(err, "unable to access files", "path", path)
			return nil, err
		}

		for _, file := range files {
			in, ferr := os.ReadFile(file)
			if ferr != nil {
				opts.Logger.Error(ferr, "unable to read file", "file", file, "path", path)
				return nil, err
			}
			yerr := yaml.Unmarshal(in, &config)
			if yerr != nil {
				opts.Logger.Error(yerr, "unable to parse yaml", "file", file, "path", path)
				return nil, err
			}
		}
	}

	res, err := resources.Parse(config.Resources, &resources.Options{})
	if err != nil {
		opts.Logger.Error(err, "unable to parse resources")
		return nil, err
	}

	evts, err := events.Parse(config.Events, &events.Options{
		Logger:    opts.Logger,
		Metrics:   opts.Metrics,
		Resources: res,
		Paths:     opts.Paths,
	})
	if err != nil {
		opts.Logger.Error(err, "unable to parse events")
		return nil, err
	}

	snks, err := sinks.Parse(config.Sinks, res, &sinks.Options{})
	if err != nil {
		opts.Logger.Error(err, "unable to parse sinks")
		return nil, err
	}

	out := &Config{
		Events:    evts,
		Resources: res,
		Sinks:     snks,
	}

	return out, nil
}
