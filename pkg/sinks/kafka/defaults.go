package kafka

import "time"

const (
	// DefaultBroker is the default Kafka broker to connect to.
	DefaultBroker = "localhost:9092"
	// DefaultTopic is the default Kafka topic to send events to.
	DefaultTopic = "events"
	// DefaultKafkaTimeout is the default timeout for Kafka operations.
	DefaultKafkaTimeout = 5 * time.Second
)
