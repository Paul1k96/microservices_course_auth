package v1

import (
	"context"
	"fmt"

	"github.com/Paul1k96/microservices_course_auth/internal/repository"
	"github.com/Paul1k96/microservices_course_platform_common/pkg/client/kafka"
)

// Consumer is a user consumer.
type Consumer struct {
	eventsRepo repository.UserEventsRepository
	consumer   kafka.Consumer
	topic      string
}

// NewConsumer creates a new user consumer.
func NewConsumer(eventsRepo repository.UserEventsRepository, consumer kafka.Consumer, topic string) *Consumer {
	return &Consumer{
		eventsRepo: eventsRepo,
		consumer:   consumer,
		topic:      topic,
	}
}

// RunConsumer runs the user consumer.
func (c *Consumer) RunConsumer(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-c.run(ctx):
			if err != nil {
				return fmt.Errorf("failed to run consumer: %w", err)
			}
		}
	}
}

func (c *Consumer) run(ctx context.Context) <-chan error {
	errChan := make(chan error)

	go func() {
		defer close(errChan)

		errChan <- c.consumer.Consume(ctx, c.topic, c.SaveEventHandler)
	}()

	return errChan
}
