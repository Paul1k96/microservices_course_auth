package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/Paul1k96/microservices_course_auth/internal/model"
	"github.com/Paul1k96/microservices_course_auth/internal/repository/user_event/kafka/mapper"
)

// Producer is a producer.
type Producer struct {
	config sarama.SyncProducer
	topic  string
}

// NewProducer creates a new producer.
func NewProducer(config sarama.SyncProducer, topic string) *Producer {
	return &Producer{
		config: config,
		topic:  topic,
	}
}

// Save saves user event.
func (p *Producer) Save(_ context.Context, event *model.UserEvent) error {
	kafkaEvent, err := mapper.NewUserEvent(event)
	if err != nil {
		return fmt.Errorf("failed to create user event: %w", err)
	}

	rawEvent, err := json.Marshal(kafkaEvent)
	if err != nil {
		return fmt.Errorf("failed to marshal user event: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder(event.ID.String()),
		Value: sarama.ByteEncoder(rawEvent),
	}

	_, _, err = p.config.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}
