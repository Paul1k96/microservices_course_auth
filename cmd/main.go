package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/Paul1k96/microservices_course_auth/pkg/proto/gen/user_v1"
	"github.com/brianvoe/gofakeit/v7"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const grpcPort = 50051

// User represents user model.
type User struct {
	ID        int64
	Name      string
	Email     string
	Password  string
	Role      user_v1.Role
	CreatedAt time.Time
	UpdatedAt time.Time
}

// UsersMap represents users repository.
type UsersMap struct {
	users map[int64]User
	mu    sync.RWMutex
}

// NewUsers creates a new users repository.
func NewUsers() *UsersMap {
	return &UsersMap{
		users: make(map[int64]User),
	}
}

// Create user to users.
func (u *UsersMap) Create(_ context.Context, user User) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	u.users[user.ID] = user

	return nil
}

// Get user by id.
func (u *UsersMap) Get(_ context.Context, id int64) (User, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()

	user, ok := u.users[id]
	if !ok {
		return User{}, fmt.Errorf("user with id %d not found", id)
	}

	return user, nil
}

// Delete user by id.
func (u *UsersMap) Delete(_ context.Context, id int64) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	delete(u.users, id)

	return nil
}

// Update user fields.
func (u *UsersMap) Update(_ context.Context, id int64, user User) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	u.users[id] = user

	return nil
}

// UsersRepository represents user repository.
type UsersRepository interface {
	Create(ctx context.Context, user User) error
	Get(ctx context.Context, id int64) (User, error)
	Update(ctx context.Context, id int64, user User) error
	Delete(ctx context.Context, id int64) error
}

// UserAPI represents user service.
type UserAPI struct {
	logger   *slog.Logger
	userRepo UsersRepository
	user_v1.UnimplementedUserServer
}

// NewUserAPI creates a new user service.
func NewUserAPI(logger *slog.Logger, userRepo UsersRepository) *UserAPI {
	return &UserAPI{logger: logger, userRepo: userRepo}
}

// Create new user.
func (u UserAPI) Create(ctx context.Context, request *user_v1.CreateRequest) (*user_v1.CreateResponse, error) {
	logger := u.logger.
		With("method", "Create").
		With("name", request.Name).
		With("email", request.Email).
		With("role", request.Role)

	createTime := time.Now()

	user := User{
		ID:        gofakeit.Int64(),
		Name:      request.Name,
		Email:     request.Email,
		Password:  request.Password,
		Role:      request.Role,
		CreatedAt: createTime,
		UpdatedAt: createTime,
	}

	err := u.userRepo.Create(ctx, user)
	if err != nil {
		logger.Error("failed to create user", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user_v1.CreateResponse{Id: user.ID}, nil
}

// Get user by id.
func (u UserAPI) Get(ctx context.Context, request *user_v1.GetRequest) (*user_v1.GetResponse, error) {
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
func (u UserAPI) Update(ctx context.Context, request *user_v1.UpdateRequest) (*user_v1.UpdateResponse, error) {
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
func (u UserAPI) Delete(ctx context.Context, request *user_v1.DeleteRequest) (*user_v1.DeleteResponse, error) {
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

func main() {
	logger := slog.Default()

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		logger.Error("failed to listen", slog.String("error", err.Error()))
		return
	}

	userDB := NewUsers()

	userAPIv1 := NewUserAPI(logger, userDB)

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	user_v1.RegisterUserServer(grpcServer, userAPIv1)

	logger.Info("server listening at", slog.Any("addr", listen.Addr()))

	if err = grpcServer.Serve(listen); err != nil {
		logger.Error("failed to serve", slog.String("error", err.Error()))
		return
	}
}
