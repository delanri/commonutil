package kafka_sarama

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/delanri/commonutil/messaging"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func (l *Kafka) AddTopicListener(topic string, callback messaging.CallbackFunc) {
	l.mu.Lock()
	defer func() {
		l.mu.Unlock()
	}()
	functions := l.CallbackFunctions[topic]
	functions = append(functions, callback)
	l.CallbackFunctions[topic] = functions
	l.Option.ListTopics = append(l.Option.ListTopics, topic)
}

func (l *Kafka) Listen() {

	l.Option.Log.Println("Starting a new Sarama consumer")

	if l.Option.MessagingLogVerbose {
		sarama.Logger = l.Option.Log
	}

	version, err := sarama.ParseKafkaVersion(l.Option.KafkaVersion)
	if err != nil {
		log.Panicf("Error parsing Kafka version: %v", err)
	}

	/**
	 * Construct a new Sarama configuration.
	 * The Kafka cluster version has to be defined before the consumer/producer is initialized.
	 */
	config := sarama.NewConfig()
	config.Version = version

	switch l.Option.Strategy {
	case "sticky":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
	case "roundrobin":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	case "range":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	default:
		log.Panicf("Unrecognized consumer group partition assignor: %s", l.Option.Strategy)
	}

	if l.Option.ConsumerOffsetOldest {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	/**
	 * Setup a new Sarama consumer group
	 */
	consumer := Consumer{
		ready:             make(chan bool),
		Option:            l.Option,
		CallbackFunctions: l.CallbackFunctions,
	}

	ctx, cancel := context.WithCancel(context.Background())
	client, err := sarama.NewConsumerGroup(l.Option.Host, l.Option.ConsumerGroup, config)
	if err != nil {
		log.Panicf("Error creating consumer group client: %v", err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := client.Consume(ctx, l.Option.ListTopics, &consumer); err != nil {
				log.Panicf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}()

	<-consumer.ready // Await till the consumer has been set up
	l.Option.Log.Println("Sarama consumer up and running!...")

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		l.Option.Log.Println("terminating: context cancelled")
	case <-sigterm:
		l.Option.Log.Println("terminating: via signal")
	}
	cancel()
	wg.Wait()
	if err = client.Close(); err != nil {
		log.Panicf("Error closing client: %v", err)
	}
}

// Consumer represents a Sarama consumer group consumer
type Consumer struct {
	ready             chan bool
	CallbackFunctions map[string][]messaging.CallbackFunc
	Option            *Option
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/master/consumer_group.go#L27-L29
	for message := range claim.Messages() {
		for _, callback := range consumer.CallbackFunctions[message.Topic] {
			err := callback(message.Value)
			if err != nil {
				consumer.Option.Log.Error(err)
			}
		}
		session.MarkMessage(message, "")
	}

	return nil
}
