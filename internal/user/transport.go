package user

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Paul1k96/microservices_course_auth/pkg/proto/gen/user_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// UsersRepository represents user repository.
type UsersRepository interface {
	Create(ctx context.Context, user User) (*int, error)
	Get(ctx context.Context, id int64) (User, error)
	Update(ctx context.Context, id int64, user User) error
	Delete(ctx context.Context, id int64) error
}

// API represents user service.
type API struct {
	logger   *slog.Logger
	userRepo UsersRepository
	user_v1.UnimplementedUserServer
}

// NewUserAPI creates a new user service.
func NewUserAPI(logger *slog.Logger, userRepo UsersRepository) *API {
	return &API{logger: logger, userRepo: userRepo}
}

// Create new user.
func (u API) Create(ctx context.Context, request *user_v1.CreateRequest) (*user_v1.CreateResponse, error) {
	logger := u.logger.
		With("method", "Create").
		With("name", request.Name).
		With("email", request.Email).
		With("role", request.Role)

	createTime := time.Now()

	user := User{
		Name:      request.Name,
		Email:     request.Email,
		Password:  request.Password,
		Role:      request.Role,
		CreatedAt: createTime,
		UpdatedAt: createTime,
	}

	userID, err := u.userRepo.Create(ctx, user)
	if err != nil {
		logger.Error("failed to create user", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user_v1.CreateResponse{Id: int64(*userID)}, nil
}

// Get user by id.
func (u API) Get(ctx context.Context, request *user_v1.GetRequest) (*user_v1.GetResponse, error) {
	logger := u.logger.
		With("method", "Get").
		With("user_id", request.Id)

	user, err := u.userRepo.Get(ctx, request.Id)
	if err != nil {
		logger.Error("failed to get user", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user_v1.GetResponse{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}, nil
}

// Update user fields.
func (u API) Update(ctx context.Context, request *user_v1.UpdateRequest) (*user_v1.UpdateResponse, error) {
	logger := u.logger.
		With("method", "Update").
		With("user_id", request.Id)

	user, err := u.userRepo.Get(ctx, request.Id)
	if err != nil {
		logger.Error("failed to get user", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if request.Name != nil {
		user.Name = request.Name.Value
	}
	if request.Email != nil {
		user.Email = request.Email.Value
	}
	user.Role = request.Role

	err = u.userRepo.Update(ctx, request.Id, user)
	if err != nil {
		logger.Error("failed to update user", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &user_v1.UpdateResponse{}, nil
}

// Delete user by id.
func (u API) Delete(ctx context.Context, request *user_v1.DeleteRequest) (*user_v1.DeleteResponse, error) {
	logger := u.logger.
		With("method", "Delete").
		With("user_id", request.Id)

	err := u.userRepo.Delete(ctx, request.Id)
	if err != nil {
		logger.Error("failed to delete user", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to delete user: %w", err)
	}

	return &user_v1.DeleteResponse{}, nil
}
