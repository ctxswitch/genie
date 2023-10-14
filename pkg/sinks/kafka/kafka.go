package kafka

import (
	"context"
	"sync"

	"ctx.sh/strata"
	"github.com/go-logr/logr"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/kversion"
)

// Options are the options for a Kafka sink.
// TODO: this is just a very basic implementation to get things working
// for another project I'm working on.  I'll need to come back and make
// this more stable and configurable.
type Options struct {
	Logger  logr.Logger
	Metrics *strata.Metrics
}

// Kafka is the representation of a kafka sink.
type Kafka struct {
	client  *kgo.Client
	topic   string
	seeds   []string
	logger  logr.Logger
	metrics *strata.Metrics
	// TODO: look at making this a buffered channel when we allow
	// for multiple sink workers.
	sendChan chan []byte
	stopOnce sync.Once
}

// New returns a new Kafka sink.
func New(cfg Config, opts *Options) *Kafka {
	return &Kafka{
		topic:    cfg.Topic,
		seeds:    cfg.Brokers,
		sendChan: make(chan []byte),
		logger:   opts.Logger,
		metrics:  opts.Metrics,
	}
}

// Init initializes the Kafka sink, setting up the client and checking
// that the brokers are reachable and the topic exists.
func (k *Kafka) Init() (err error) {
	k.client, err = kgo.NewClient(
		kgo.SeedBrokers(k.seeds...),
		kgo.DefaultProduceTopic(k.topic),
	)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), DefaultKafkaTimeout)
	defer cancel()
	if err = k.client.Ping(ctx); err != nil {
		k.logger.Error(err, "failed to ping kafka")
		return
	}

	err = k.createTopic()
	return
}

// createTopic creates the configured topic if it does not already exist.
func (k *Kafka) createTopic() error {
	// TODO: this should only create the topic if it's configured to do so
	// I'll add that later.

	var adminClient *kadm.Client
	{
		client, err := kgo.NewClient(
			kgo.SeedBrokers(k.seeds...),
			kgo.MaxVersions(kversion.V2_4_0()),
		)
		if err != nil {
			return err
		}
		defer client.Close()
		adminClient = kadm.NewClient(client)
	}

	ctx, cancel := context.WithTimeout(context.Background(), DefaultKafkaTimeout)
	defer cancel()
	topics, err := adminClient.ListTopics(ctx)
	if err != nil {
		return err
	}

	if topics.Has(k.topic) {
		return nil
	}

	// TODO: Partitions, replication factor, and config should be set in the sink config.
	resp, err := adminClient.CreateTopic(ctx, -1, -1, nil, k.topic)
	if err != nil {
		return err
	}
	k.logger.Info("created topic", "topic", k.topic, "response", resp)

	return nil
}

// Start starts the Kafka sink.
func (k *Kafka) Start() {
	ctx := context.Background()
	for data := range k.sendChan {
		k.send(ctx, data)
	}
}

// send sends the data asynchronously to the configured topic.
// TODO: The sink is very basic at the moment and should be improved.
func (k *Kafka) send(ctx context.Context, data []byte) {
	record := &kgo.Record{
		Topic: k.topic,
		Value: data,
	}

	k.client.Produce(ctx, record, func(r *kgo.Record, err error) {
		if err != nil {
			k.logger.Error(err, "failed to produce record")
			k.metrics.CounterInc("produce_error")
			return
		}

		k.metrics.CounterInc("produce_success")
	})
}

// Stop stops the Kafka sink.
func (k *Kafka) Stop() {
	defer k.client.Close()

	k.stopOnce.Do(func() {
		close(k.sendChan)
	})
}

// SendChannel returns the channel that the event generators will send to.
func (k *Kafka) SendChannel() chan<- []byte {
	return k.sendChan
}
