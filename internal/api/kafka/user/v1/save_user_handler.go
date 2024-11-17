package v1

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/Paul1k96/microservices_course_auth/internal/api/kafka/user/v1/mapper"
	modelKafka "github.com/Paul1k96/microservices_course_auth/internal/api/kafka/user/v1/model"
)

// SaveEventHandler saves user event.
func (c *Consumer) SaveEventHandler(ctx context.Context, message *sarama.ConsumerMessage) error {
	var kafkaEvent modelKafka.UserEvent

	if err := json.Unmarshal(message.Value, &kafkaEvent); err != nil {
		return fmt.Errorf("failed to unmarshal user: %w", err)
	}

	event, err := mapper.ToUserEventFromKafka(&kafkaEvent)
	if err != nil {
		return fmt.Errorf("failed to map user event: %w", err)
	}

	if err = c.eventsRepo.Save(ctx, event); err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}

	return nil
}
