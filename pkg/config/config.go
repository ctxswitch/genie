package config

import (
	"fmt"
	"os"
	"path/filepath"

	"ctx.sh/strata"
	"github.com/go-logr/logr"
	"gopkg.in/yaml.v2"
	"stvz.io/genie/pkg/events"
	"stvz.io/genie/pkg/resources"
	"stvz.io/genie/pkg/sinks"
)

// LoadOptions are the options for loading the config.
type LoadOptions struct {
	Logger  logr.Logger
	Metrics *strata.Metrics
	Paths   []string
}

// Config is the config for the genie application.
type Config struct {
	Events    events.Events
	Resources *resources.Resources
	Sinks     *sinks.Sinks
}

// Load loads the config files from the provided paths.  It can take multiple paths
// and will load all files with a yaml extension from each path.  It returns an error
// if any of the files fail to load or parse.
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

	res := resources.New(config.Resources)
	snks := sinks.New(config.Sinks, &sinks.Options{
		Logger:  opts.Logger,
		Metrics: opts.Metrics,
	})

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

	out := &Config{
		Events:    evts,
		Resources: res,
		Sinks:     snks,
	}

	return out, nil
}

func (c *Config) HasEvents() bool {
	return len(c.Events) > 0
}

func (c *Config) HasEvent(name string) bool {
	_, has := c.Events[name]
	return has
}
