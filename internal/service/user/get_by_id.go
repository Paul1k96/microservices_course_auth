package user

import (
	"context"
	"fmt"

	"github.com/Paul1k96/microservices_course_auth/internal/model"
)

func (s *service) GetByID(ctx context.Context, id int64) (*model.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	if user == nil {
		return nil, fmt.Errorf("user with id %d not found", id)
	}

	return user, nil
}
