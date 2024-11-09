package repository

import (
	"context"

	"github.com/Paul1k96/microservices_course_auth/internal/model"
)

// UsersRepository represents user repository.
type UsersRepository interface {
	Create(ctx context.Context, user *model.User) (int64, error)
	GetByID(ctx context.Context, id int64) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id int64) error
}

// UsersCache represents user cache repository.
type UsersCache interface {
	Set(ctx context.Context, user *model.User) error
	Get(ctx context.Context, id int64) (*model.User, error)
	Delete(ctx context.Context, id int64) error
}
