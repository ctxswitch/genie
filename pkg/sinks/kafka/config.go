package kafka

import "fmt"

// Config is the configuration for the Kafka sink.
type Config struct {
	Brokers []string `yaml:"brokers"`
	Topic   string   `yaml:"topic"`
}

// validate ensures that the Kafka config is valid.
func (c *Config) validate() error {
	if len(c.Brokers) == 0 {
		return fmt.Errorf("at least one broker must be provided")
	}

	if c.Topic == "" {
		return fmt.Errorf("topic must be provided")
	}

	return nil
}

// UnmarshalYAML implements yaml.Unmarshaler for defaulting the Kafka config.
func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type ConfigDefaulted Config

	defaults := ConfigDefaulted{
		Topic: DefaultTopic,
		Brokers: []string{
			DefaultBroker,
		},
	}

	out := defaults
	if err := unmarshal(&out); err != nil {
		return err
	}

	kafka := Config(out)
	if err := kafka.validate(); err != nil {
		return err
	}

	*c = kafka
	return nil
}
