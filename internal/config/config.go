package config

import (
	"time"

	"github.com/IBM/sarama"
)

// PGConfig represents configuration for PostgreSQL.
type PGConfig interface {
	GetDSN() string
}

// GRPCConfig represents configuration for gRPC.
type GRPCConfig interface {
	GetAddress() string
}

// HTTPConfig represents configuration for HTTP.
type HTTPConfig interface {
	GetAddress() string
	GetReadHeaderTimeout() time.Duration
	GetGracefulShutdownTimeout() time.Duration
}

// RedisConfig represents configuration for Redis.
type RedisConfig interface {
	GetAddress() string
	GetConnectionTimeout() time.Duration
	GetMaxIdle() int
	GetIdleTimeout() time.Duration
	GetUserTTL() time.Duration
}

// KafkaConsumerConfig represents configuration for Kafka consumer.
type KafkaConsumerConfig interface {
	Brokers() []string
	GroupID() string
	Config() *sarama.Config
	Topic() string
}

// KafkaProducerConfig represents configuration for Kafka producer.
type KafkaProducerConfig interface {
	Brokers() []string
	Topic() string
	Config() *sarama.Config
}
