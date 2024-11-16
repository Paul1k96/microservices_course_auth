package user

import (
	"log/slog"

	"github.com/Paul1k96/microservices_course_auth/internal/repository"
	svc "github.com/Paul1k96/microservices_course_auth/internal/service"
	"github.com/Paul1k96/microservices_course_platform_common/pkg/client/db"
)

type service struct {
	logger    *slog.Logger
	txManager db.TxManager

	repo  repository.UsersRepository
	cache repository.UsersCache
}

// NewService creates a new service.
func NewService(
	logger *slog.Logger,
	txManager db.TxManager,
	repo repository.UsersRepository,
	cache repository.UsersCache,
) svc.UserService {
	return &service{
		logger:    logger,
		txManager: txManager,
		repo:      repo,
		cache:     cache,
	}
}
