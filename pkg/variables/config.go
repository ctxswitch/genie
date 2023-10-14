package variables

// Config is the configuration for a single variable.
type Config struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}
