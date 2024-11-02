package service

import (
	"context"

	"github.com/Paul1k96/microservices_course_auth/internal/model"
)

// UserService represents user service.
type UserService interface {
	Create(ctx context.Context, user *model.User) (int64, error)
	GetByID(ctx context.Context, id int64) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id int64) error
}
