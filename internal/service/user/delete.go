package user

import (
	"context"
	"fmt"
	"log/slog"
)

// Delete deletes user by id.
func (s *service) Delete(ctx context.Context, id int64) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	err = s.cache.Delete(ctx, id)
	if err != nil {
		s.logger.Error("failed to delete user from cache:", slog.String("error", err.Error()))
	}

	return nil
}
