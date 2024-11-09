package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/Paul1k96/microservices_course_auth/internal/errs"
	"github.com/Paul1k96/microservices_course_auth/internal/model"
)

// GetByID returns user by id.
func (s *service) GetByID(ctx context.Context, id int64) (*model.User, error) {
	var (
		user *model.User
		err  error
	)

	user, err = s.cache.Get(ctx, id)
	if err != nil && !errors.Is(err, errs.ErrUserNotFound) {
		return nil, fmt.Errorf("failed to get user from cache: %w", err)
	}

	if user != nil {
		return user, nil
	}

	user, err = s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	if user == nil {
		return nil, errs.ErrUserNotFound
	}

	_ = s.cache.Set(ctx, user)

	return user, nil
}
