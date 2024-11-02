package repository

import (
	"context"

	"github.com/Paul1k96/microservices_course_auth/internal/model"
	modelRepo "github.com/Paul1k96/microservices_course_auth/internal/repository/user/model"
)

// UsersRepository represents user repository.
type UsersRepository interface {
	Create(ctx context.Context, user *modelRepo.User) (int64, error)
	GetByID(ctx context.Context, id int64) (*model.User, error)
	Update(ctx context.Context, user *modelRepo.User) error
	Delete(ctx context.Context, id int64) error
}
