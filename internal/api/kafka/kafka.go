package kafka

import "context"

// UserCreateConsumer is a user consumer.
type UserCreateConsumer interface {
	RunConsumer(ctx context.Context) error
}
