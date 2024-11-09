package pg

import (
	"context"
	"fmt"
	"time"

	"github.com/Paul1k96/microservices_course_auth/internal/errs"
	"github.com/Paul1k96/microservices_course_auth/internal/model"
	"github.com/Paul1k96/microservices_course_auth/internal/repository/user/redis/mapper"
	modelRepo "github.com/Paul1k96/microservices_course_auth/internal/repository/user/redis/model"
	"github.com/Paul1k96/microservices_course_platform_common/pkg/client/cache"
	"github.com/gomodule/redigo/redis"
)

// Repository represents user repository.
type Repository struct {
	redisCache cache.RedisClient
	ttl        time.Duration
}

// NewRepository creates a new instance of repository.UsersRepository.
func NewRepository(redisCache cache.RedisClient, ttl time.Duration) *Repository {
	return &Repository{redisCache: redisCache, ttl: ttl}
}

// Set user.
func (r *Repository) Set(ctx context.Context, user *model.User) error {
	userToCreate := modelRepo.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role.String(),
		CreatedAt: user.CreatedAt.UnixNano(),
	}

	if user.UpdatedAt != nil {
		updateTime := user.UpdatedAt.UnixNano()
		userToCreate.UpdatedAt = &updateTime
	}

	err := r.redisCache.HSet(ctx, fmt.Sprintf("%d", user.ID), userToCreate)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	err = r.redisCache.Expire(ctx, fmt.Sprintf("%d", user.ID), r.ttl)
	if err != nil {
		return fmt.Errorf("set ttl: %w", err)
	}

	return nil
}

// Get user by id.
func (r *Repository) Get(ctx context.Context, id int64) (*model.User, error) {
	values, err := r.redisCache.HGetAll(ctx, fmt.Sprintf("%d", id))
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	if len(values) == 0 {
		return nil, errs.ErrUserNotFound
	}

	var user modelRepo.User
	err = redis.ScanStruct(values, &user)
	if err != nil {
		return nil, fmt.Errorf("scan user: %w", err)
	}

	return mapper.ToUserFromRepo(&user), nil
}

// Delete user by id.
func (r *Repository) Delete(ctx context.Context, id int64) error {
	err := r.redisCache.Delete(ctx, fmt.Sprintf("%d", id))
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	return nil
}
