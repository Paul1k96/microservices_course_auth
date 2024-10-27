package user

import (
	"context"
	"fmt"
)

// Delete deletes user by id.
func (s *service) Delete(ctx context.Context, id int64) error {
	if txErr := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		err := s.checkUserExistsByID(ctx, id)
		if err != nil {
			return fmt.Errorf("delete user: %w", err)
		}

		err = s.repo.Delete(ctx, id)
		if err != nil {
			return fmt.Errorf("delete user: %w", err)
		}

		return nil
	}); txErr != nil {
		return fmt.Errorf("transaction error: %w", txErr)
	}

	return nil
}

func (s *service) checkUserExistsByID(ctx context.Context, id int64) error {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("check user exists: %w", err)
	}

	if user == nil {
		return fmt.Errorf("user with id %d not found", id)
	}

	return nil
}
