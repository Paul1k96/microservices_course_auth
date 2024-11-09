package user

import (
	"github.com/Paul1k96/microservices_course_auth/internal/repository"
	svc "github.com/Paul1k96/microservices_course_auth/internal/service"
	"github.com/Paul1k96/microservices_course_platform_common/pkg/client/db"
)

type service struct {
	repo      repository.UsersRepository
	cache     repository.UsersCache
	txManager db.TxManager
}

// NewService creates a new service.
func NewService(
	repo repository.UsersRepository,
	cache repository.UsersCache,
	txManager db.TxManager,
) svc.UserService {
	return &service{
		repo:      repo,
		cache:     cache,
		txManager: txManager,
	}
}
