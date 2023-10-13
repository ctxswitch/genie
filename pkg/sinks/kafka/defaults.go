package kafka

import "time"

const (
	DefaultBroker       = "localhost:9092"
	DefaultTopic        = "events"
	DefaultKafkaTimeout = 5 * time.Second
)
