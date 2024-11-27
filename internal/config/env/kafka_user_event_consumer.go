package env

import (
	"os"
	"strings"

	"github.com/IBM/sarama"
	"github.com/pkg/errors"
)

const (
	brokersEnvName           = "KAFKA_BROKERS"
	userCreateGroupIDEnvName = "KAFKA_CONSUMER_USER_EVENTS_GROUP_ID"
	userCreateTopicEnvName   = "KAFKA_CONSUMER_USER_EVENTS_TOPIC"
)

// KafkaUserCreateConsumerConfig represents configuration for Kafka consumer.
type KafkaUserCreateConsumerConfig struct {
	brokers []string
	groupID string
	topic   string
}

// NewKafkaUserCreateConsumerConfig creates a new Kafka consumer configuration.
func NewKafkaUserCreateConsumerConfig() (*KafkaUserCreateConsumerConfig, error) {
	brokersStr := os.Getenv(brokersEnvName)
	if len(brokersStr) == 0 {
		return nil, errors.New("kafka brokers address not found")
	}

	brokers := strings.Split(brokersStr, ",")

	groupID := os.Getenv(userCreateGroupIDEnvName)
	if len(groupID) == 0 {
		return nil, errors.New("kafka group id not found")
	}

	topic := os.Getenv(userCreateTopicEnvName)
	if len(topic) == 0 {
		return nil, errors.New("kafka topic not found")
	}

	return &KafkaUserCreateConsumerConfig{
		brokers: brokers,
		groupID: groupID,
		topic:   topic,
	}, nil
}

// Brokers returns Kafka brokers addresses.
func (cfg *KafkaUserCreateConsumerConfig) Brokers() []string {
	return cfg.brokers
}

// GroupID returns group ID for sarama consumer.
func (cfg *KafkaUserCreateConsumerConfig) GroupID() string {
	return cfg.groupID
}

// Topic returns Kafka topic.
func (cfg *KafkaUserCreateConsumerConfig) Topic() string {
	return cfg.topic
}

// Config returns configuration for sarama consumer.
func (cfg *KafkaUserCreateConsumerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V2_6_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	return config
}
