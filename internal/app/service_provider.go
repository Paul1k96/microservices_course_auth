package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/IBM/sarama"
	userv1 "github.com/Paul1k96/microservices_course_auth/internal/api/proto/user/v1"
	"github.com/Paul1k96/microservices_course_platform_common/pkg/client/cache"
	"github.com/Paul1k96/microservices_course_platform_common/pkg/client/db/pg"
	"github.com/Paul1k96/microservices_course_platform_common/pkg/client/db/transaction"
	"github.com/Paul1k96/microservices_course_platform_common/pkg/client/kafka"

	kafkaApi "github.com/Paul1k96/microservices_course_auth/internal/api/kafka"
	userKafkaV1 "github.com/Paul1k96/microservices_course_auth/internal/api/kafka/user/v1"
	"github.com/Paul1k96/microservices_course_auth/internal/config"
	"github.com/Paul1k96/microservices_course_auth/internal/config/env"
	"github.com/Paul1k96/microservices_course_auth/internal/repository"
	userpg "github.com/Paul1k96/microservices_course_auth/internal/repository/user/pg"
	userRedis "github.com/Paul1k96/microservices_course_auth/internal/repository/user/redis"
	usereventsproducer "github.com/Paul1k96/microservices_course_auth/internal/repository/user_event/kafka"
	usereventspg "github.com/Paul1k96/microservices_course_auth/internal/repository/user_event/pg"
	"github.com/Paul1k96/microservices_course_auth/internal/service"
	userSvc "github.com/Paul1k96/microservices_course_auth/internal/service/user"
	commonRedis "github.com/Paul1k96/microservices_course_platform_common/pkg/client/cache/redis"
	"github.com/Paul1k96/microservices_course_platform_common/pkg/client/db"
	kafkaConsumer "github.com/Paul1k96/microservices_course_platform_common/pkg/client/kafka/consumer"
	"github.com/Paul1k96/microservices_course_platform_common/pkg/closer"
	redigo "github.com/gomodule/redigo/redis"
)

type serviceProvider struct {
	pgConfig                      config.PGConfig
	grpcConfig                    config.GRPCConfig
	httpConfig                    config.HTTPConfig
	redisConfig                   config.RedisConfig
	kafkaCreateUserConsumerConfig config.KafkaConsumerConfig
	kafkaUserEventsProducerConfig config.KafkaProducerConfig
	logger                        *slog.Logger

	consumerGroupHandler *kafkaConsumer.GroupHandler
	consumerGroup        sarama.ConsumerGroup
	consumer             kafka.Consumer
	userCreateConsumer   kafkaApi.UserCreateConsumer
	userEventsProducer   repository.UserEventsRepository

	dbClient             db.Client
	redisClient          cache.RedisClient
	redisPool            *redigo.Pool
	txManager            db.TxManager
	usersRepository      repository.UsersRepository
	userEventsRepository repository.UserEventsRepository
	usersCache           repository.UsersCache

	usersService service.UserService

	userV1Impl *userv1.Implementation
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

// HTTPConfig returns an instance of config.HTTPConfig.
func (s *serviceProvider) HTTPConfig() (config.HTTPConfig, error) {
	if s.httpConfig == nil {
		cfg, err := env.NewHTTPConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to get http config: %w", err)
		}

		s.httpConfig = cfg
	}

	return s.httpConfig, nil
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

// KafkaCreateUserConsumerConfig returns an instance of config.KafkaConsumerConfig.
func (s *serviceProvider) KafkaCreateUserConsumerConfig() (config.KafkaConsumerConfig, error) {
	if s.kafkaCreateUserConsumerConfig == nil {
		cfg, err := env.NewKafkaUserCreateConsumerConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to get kafka consumer config: %w", err)
		}

		s.kafkaCreateUserConsumerConfig = cfg
	}

	return s.kafkaCreateUserConsumerConfig, nil
}

// KafkaUserEventsProducerConfig returns an instance of config.KafkaProducerConfig.
func (s *serviceProvider) KafkaUserEventsProducerConfig() (config.KafkaProducerConfig, error) {
	if s.kafkaUserEventsProducerConfig == nil {
		cfg, err := env.NewKafkaUserEventProducerConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to get kafka producer config: %w", err)
		}

		s.kafkaUserEventsProducerConfig = cfg
	}

	return s.kafkaUserEventsProducerConfig, nil
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

// ConsumerGroup returns an instance of kafka.ConsumerGroup.
func (s *serviceProvider) ConsumerGroup() (sarama.ConsumerGroup, error) {
	if s.consumerGroup == nil {
		cfg, err := s.KafkaCreateUserConsumerConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to get kafka consumer config: %w", err)
		}

		consumerGroup, err := sarama.NewConsumerGroup(
			cfg.Brokers(),
			cfg.GroupID(),
			cfg.Config(),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create consumer group: %w", err)
		}

		s.consumerGroup = consumerGroup
	}

	return s.consumerGroup, nil
}

// Consumer returns an instance of kafka.Consumer.
func (s *serviceProvider) Consumer() (kafka.Consumer, error) {
	if s.consumer == nil {
		group, err := s.ConsumerGroup()
		if err != nil {
			return nil, fmt.Errorf("failed to get consumer group: %w", err)
		}

		s.consumer = kafkaConsumer.NewConsumer(
			group,
			s.ConsumerGroupHandler(),
			s.logger,
		)

		closer.Add(s.consumer.Close)
	}

	return s.consumer, nil
}

// ConsumerGroupHandler returns an instance of kafka.ConsumerGroupHandler.
func (s *serviceProvider) ConsumerGroupHandler() *kafkaConsumer.GroupHandler {
	if s.consumerGroupHandler == nil {
		s.consumerGroupHandler = kafkaConsumer.NewGroupHandler(s.logger)
	}

	return s.consumerGroupHandler
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

// UserEventsRepository returns an instance of repository.UserEventsRepository.
func (s *serviceProvider) UserEventsRepository(ctx context.Context) (repository.UserEventsRepository, error) {
	if s.userEventsRepository == nil {
		dbClient, err := s.DBClient(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get db client: %w", err)
		}

		s.userEventsRepository = usereventspg.NewRepository(dbClient.DB())
	}

	return s.userEventsRepository, nil
}

// UserEventsProducer returns an instance of repository.UserEventsRepository.
func (s *serviceProvider) UserEventsProducer() (repository.UserEventsRepository, error) {
	if s.userEventsProducer == nil {
		cfg, err := s.KafkaUserEventsProducerConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to get kafka producer config: %w", err)
		}

		producer, err := sarama.NewSyncProducer(cfg.Brokers(), cfg.Config())
		if err != nil {
			return nil, fmt.Errorf("failed to create producer: %w", err)
		}

		userEventProducer := usereventsproducer.NewProducer(producer, cfg.Topic())

		s.userEventsProducer = userEventProducer
	}

	return s.userEventsProducer, nil
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

		userEventsProducer, err := s.UserEventsProducer()
		if err != nil {
			return nil, fmt.Errorf("failed to get user events producer: %w", err)
		}

		userCache, err := s.UsersCache(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get users cache: %w", err)
		}

		s.usersService = userSvc.NewService(s.logger, txManager, userRepository, userEventsProducer, userCache)
	}

	return s.usersService, nil
}

// UserV1Impl returns an instance of user_v1.Implementation.
func (s *serviceProvider) UserV1Impl(ctx context.Context) (*userv1.Implementation, error) {
	if s.userV1Impl == nil {
		usersService, err := s.UsersService(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get users service: %w", err)
		}

		s.userV1Impl = userv1.NewImplementation(s.logger, usersService)
	}

	return s.userV1Impl, nil
}

// UserCreateConsumer returns an instance of kafka.UserCreateConsumer.
func (s *serviceProvider) UserCreateConsumer(ctx context.Context) (kafkaApi.UserCreateConsumer, error) {
	if s.userCreateConsumer == nil {
		cfg, err := s.KafkaCreateUserConsumerConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to get kafka consumer config: %w", err)
		}

		userEventsRepo, err := s.UserEventsRepository(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get users repository: %w", err)
		}

		consumer, err := s.Consumer()
		if err != nil {
			return nil, fmt.Errorf("failed to get consumer: %w", err)
		}

		s.userCreateConsumer = userKafkaV1.NewConsumer(userEventsRepo, consumer, cfg.Topic())
	}

	return s.userCreateConsumer, nil
}
