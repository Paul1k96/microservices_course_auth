package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Paul1k96/microservices_course_auth/internal/api/user_v1"
	"github.com/Paul1k96/microservices_course_platform_common/pkg/client/cache"
	"github.com/Paul1k96/microservices_course_platform_common/pkg/client/db/pg"
	"github.com/Paul1k96/microservices_course_platform_common/pkg/client/db/transaction"

	"github.com/Paul1k96/microservices_course_auth/internal/config"
	"github.com/Paul1k96/microservices_course_auth/internal/config/env"
	"github.com/Paul1k96/microservices_course_auth/internal/repository"
	userpg "github.com/Paul1k96/microservices_course_auth/internal/repository/user/pg"
	userRedis "github.com/Paul1k96/microservices_course_auth/internal/repository/user/redis"
	"github.com/Paul1k96/microservices_course_auth/internal/service"
	userSvc "github.com/Paul1k96/microservices_course_auth/internal/service/user"
	commonRedis "github.com/Paul1k96/microservices_course_platform_common/pkg/client/cache/redis"
	"github.com/Paul1k96/microservices_course_platform_common/pkg/client/db"
	"github.com/Paul1k96/microservices_course_platform_common/pkg/closer"
	redigo "github.com/gomodule/redigo/redis"
)

type serviceProvider struct {
	pgConfig    config.PGConfig
	grpcConfig  config.GRPCConfig
	redisConfig config.RedisConfig
	logger      *slog.Logger

	dbClient        db.Client
	redisClient     cache.RedisClient
	redisPool       *redigo.Pool
	txManager       db.TxManager
	usersRepository repository.UsersRepository
	usersCache      repository.UsersCache

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

// RedisConfig returns an instance of config.RedisConfig.
func (s *serviceProvider) RedisConfig() (config.RedisConfig, error) {
	if s.redisConfig == nil {
		cfg, err := env.NewRedisConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to get redis config: %w", err)
		}

		s.redisConfig = cfg
	}

	return s.redisConfig, nil
}

// RedisPool returns an instance of redigo.Pool.
func (s *serviceProvider) RedisPool() (*redigo.Pool, error) {
	if s.redisPool == nil {
		cfg, err := s.RedisConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to get redis config: %w", err)
		}

		s.redisPool = &redigo.Pool{
			MaxIdle:     cfg.GetMaxIdle(),
			IdleTimeout: cfg.GetIdleTimeout(),
			DialContext: func(ctx context.Context) (redigo.Conn, error) {
				return redigo.DialContext(ctx, "tcp", cfg.GetAddress())
			},
		}
	}

	return s.redisPool, nil
}

// CacheClient returns an instance of cache.RedisClient.
func (s *serviceProvider) CacheClient(ctx context.Context) (cache.RedisClient, error) {
	if s.redisClient == nil {
		pool, err := s.RedisPool()
		if err != nil {
			return nil, fmt.Errorf("failed to get redis pool: %w", err)
		}

		cfg, err := s.RedisConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to get redis config: %w", err)
		}

		cl := commonRedis.NewClient(pool, cfg)

		err = cl.Ping(ctx)
		if err != nil {
			return nil, fmt.Errorf("ping error: %w", err)
		}

		s.redisClient = cl
	}

	return s.redisClient, nil
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
		s.usersRepository = userpg.NewRepository(dbClient.DB())
	}

	return s.usersRepository, nil
}

// UsersCache returns an instance of repository.UsersCache.
func (s *serviceProvider) UsersCache(ctx context.Context) (repository.UsersCache, error) {
	if s.usersCache == nil {
		cfg, err := s.RedisConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to get redis config: %w", err)
		}
		cacheClient, err := s.CacheClient(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get cache client: %w", err)
		}
		s.usersCache = userRedis.NewRepository(cacheClient, cfg.GetUserTTL())
	}

	return s.usersCache, nil
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
		userCache, err := s.UsersCache(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get users cache: %w", err)
		}
		s.usersService = userSvc.NewService(userRepository, userCache, txManager)
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
