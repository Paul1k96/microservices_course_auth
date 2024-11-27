package user

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Paul1k96/microservices_course_auth/internal/model"
)

// Delete deletes user by id.
func (s *service) Delete(ctx context.Context, id int64) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get user by id: %w", err)
	}

	err = s.repo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	err = s.cache.Delete(ctx, id)
	if err != nil {
		s.logger.Error("failed to delete user from cache:", slog.String("error", err.Error()))
	}

	// first argument must be user, which do delete
	err = s.events.Save(ctx, model.NewDeleteUserEvent(id, id))
	if err != nil {
		s.logger.Error("failed to save user event:", slog.String("error", err.Error()))
	}

	return nil
}
