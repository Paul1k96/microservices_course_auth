package user

import (
	"context"
	"fmt"
	"time"

	"github.com/Paul1k96/microservices_course_auth/internal/model"
)

// Update updates user.
func (s *service) Update(ctx context.Context, user *model.User) error {
	if txErr := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		if err := s.validateUpdateUser(user); err != nil {
			return fmt.Errorf("failed to validate user: %w", err)
		}

		updateTime := time.Now()
		user.UpdatedAt = &updateTime
		if err := s.repo.Update(ctx, user); err != nil {
			return fmt.Errorf("failed to update user: %w", err)
		}

		return nil
	}); txErr != nil {
		return fmt.Errorf("transaction error: %w", txErr)
	}

	_ = s.cache.Set(ctx, user)

	return nil
}

func (s *service) validateUpdateUser(user *model.User) error {
	if user.Name != "" {
		if err := s.validateUserName(user.Name); err != nil {
			return fmt.Errorf("name validation: %w", err)
		}
	}

	if user.Email != "" {
		if err := s.validateUserEmail(user.Email); err != nil {
			return fmt.Errorf("email validation: %w", err)
		}
	}

	err := s.validateUserRole(user.Role)
	if err != nil {
		return fmt.Errorf("role validation: %w", err)
	}

	return nil
}
