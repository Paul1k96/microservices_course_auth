package env

import (
	"os"
	"strconv"
	"strings"

	"github.com/IBM/sarama"
	"github.com/pkg/errors"
)

const (
	userEventProducerTopicEnvName           = "KAFKA_PRODUCER_USER_EVENTS_TOPIC"
	userEventProducerRequiredAcksEnvName    = "KAFKA_PRODUCER_USER_EVENTS_REQUIRED_ACKS"
	userEventProducerRetryMaxEnvName        = "KAFKA_PRODUCER_USER_EVENTS_RETRY_MAX"
	userEventProducerReturnSuccessesEnvName = "KAFKA_PRODUCER_USER_EVENTS_RETURN_SUCCESSES"
)

// KafkaUserEventProducerConfig represents configuration for Kafka producer.
type KafkaUserEventProducerConfig struct {
	brokers         []string
	topic           string
	requiredAcks    int16
	retryMax        int
	returnSuccesses bool
}

// NewKafkaUserEventProducerConfig creates a new Kafka producer configuration.
func NewKafkaUserEventProducerConfig() (*KafkaUserEventProducerConfig, error) {
	brokersStr := os.Getenv(brokersEnvName)
	if len(brokersStr) == 0 {
		return nil, errors.New("kafka brokers address not found")
	}

	brokers := strings.Split(brokersStr, ",")

	topic := os.Getenv(userEventProducerTopicEnvName)
	if len(topic) == 0 {
		return nil, errors.New("kafka topic not found")
	}

	requiredAcksStr := os.Getenv(userEventProducerRequiredAcksEnvName)
	if len(requiredAcksStr) == 0 {
		return nil, errors.New("kafka required acks not found")
	}
	requiredAcks, err := strconv.ParseInt(requiredAcksStr, 10, 16)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse required acks")
	}

	retryMaxStr := os.Getenv(userEventProducerRetryMaxEnvName)
	if len(retryMaxStr) == 0 {
		return nil, errors.New("kafka retry max not found")
	}
	retryMax, err := strconv.Atoi(retryMaxStr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse retry max")
	}

	returnSuccessesStr := os.Getenv(userEventProducerReturnSuccessesEnvName)
	if len(returnSuccessesStr) == 0 {
		return nil, errors.New("kafka return successes not found")
	}
	returnSuccesses, err := strconv.ParseBool(returnSuccessesStr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse return successes")
	}

	return &KafkaUserEventProducerConfig{
		brokers:         brokers,
		topic:           topic,
		requiredAcks:    int16(requiredAcks),
		retryMax:        retryMax,
		returnSuccesses: returnSuccesses,
	}, nil
}

// Brokers returns Kafka brokers addresses.
func (cfg *KafkaUserEventProducerConfig) Brokers() []string {
	return cfg.brokers
}

// Topic returns Kafka topic.
func (cfg *KafkaUserEventProducerConfig) Topic() string {
	return cfg.topic
}

// Config returns configuration for sarama producer.
func (cfg *KafkaUserEventProducerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.RequiredAcks(cfg.requiredAcks)
	config.Producer.Retry.Max = cfg.retryMax
	config.Producer.Return.Successes = cfg.returnSuccesses

	return config
}
