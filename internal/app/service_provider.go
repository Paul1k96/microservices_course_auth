package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Paul1k96/microservices_course_auth/internal/api/user_v1"
	"github.com/Paul1k96/microservices_course_auth/internal/client/db"
	"github.com/Paul1k96/microservices_course_auth/internal/client/db/pg"
	"github.com/Paul1k96/microservices_course_auth/internal/client/db/transaction"
	"github.com/Paul1k96/microservices_course_auth/internal/closer"
	"github.com/Paul1k96/microservices_course_auth/internal/config"
	"github.com/Paul1k96/microservices_course_auth/internal/config/env"
	"github.com/Paul1k96/microservices_course_auth/internal/repository"
	"github.com/Paul1k96/microservices_course_auth/internal/repository/user"
	"github.com/Paul1k96/microservices_course_auth/internal/service"
	userSvc "github.com/Paul1k96/microservices_course_auth/internal/service/user"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig
	logger     *slog.Logger

	dbClient        db.Client
	txManager       db.TxManager
	usersRepository repository.UsersRepository

	usersService service.UserService

	userV1Impl *user_v1.Implementation
}

func newServiceProvider(logger *slog.Logger) *serviceProvider {
	return &serviceProvider{logger: logger}
}

// PGConfig returns an instance of config.PGConfig.
func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		s.pgConfig = env.NewPGConfig()
	}

	return s.pgConfig
}

// GRPCConfig returns an instance of config.GRPCConfig.
func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		s.grpcConfig = env.NewGRPCConfig()
	}

	return s.grpcConfig
}

// DBClient returns an instance of db.Client.
func (s *serviceProvider) DBClient(ctx context.Context) (db.Client, error) {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().GetDSN())
		if err != nil {
			return nil, fmt.Errorf("failed to create db client: %w", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			return nil, fmt.Errorf("ping error: %w", err)
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient, nil
}

// TxManager returns an instance of db.TxManager.
func (s *serviceProvider) TxManager(ctx context.Context) (db.TxManager, error) {
	if s.txManager == nil {
		dbClient, err := s.DBClient(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get db client: %w", err)
		}
		s.txManager = transaction.NewTransactionManager(dbClient.DB())
	}

	return s.txManager, nil
}

// UsersRepository returns an instance of repository.UsersRepository.
func (s *serviceProvider) UsersRepository(ctx context.Context) (repository.UsersRepository, error) {
	if s.usersRepository == nil {
		dbClient, err := s.DBClient(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get db client: %w", err)
		}
		s.usersRepository = user.NewRepository(dbClient.DB())
	}

	return s.usersRepository, nil
}

// UsersService returns an instance of service.UserService.
func (s *serviceProvider) UsersService(ctx context.Context) (service.UserService, error) {
	if s.usersService == nil {
		txManager, err := s.TxManager(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get tx manager: %w", err)
		}
		userRepository, err := s.UsersRepository(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get users repository: %w", err)
		}
		s.usersService = userSvc.NewService(userRepository, txManager)
	}

	return s.usersService, nil
}

// UserV1Impl returns an instance of user_v1.Implementation.
func (s *serviceProvider) UserV1Impl(ctx context.Context) (*user_v1.Implementation, error) {
	if s.userV1Impl == nil {
		usersService, err := s.UsersService(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get users service: %w", err)
		}
		s.userV1Impl = user_v1.NewImplementation(s.logger, usersService)
	}

	return s.userV1Impl, nil
}
