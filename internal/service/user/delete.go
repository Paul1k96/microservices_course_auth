package user

import (
	"context"
	"fmt"
)

// Delete deletes user by id.
func (s *service) Delete(ctx context.Context, id int64) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	_ = s.cache.Delete(ctx, id)

	return nil
}
