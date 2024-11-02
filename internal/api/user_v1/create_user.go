package user_v1

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Paul1k96/microservices_course_auth/internal/mapper"
	desc "github.com/Paul1k96/microservices_course_auth/pkg/proto/gen/user_v1"
)

// Create new user.
func (u *Implementation) Create(ctx context.Context, request *desc.CreateRequest) (*desc.CreateResponse, error) {
	logger := u.logger.
		With("method", "Create").
		With("name", request.Name).
		With("email", request.Email).
		With("role", request.Role)

	err := u.checkPasswordConfirm(request)
	if err != nil {
		logger.Error("failed to check password confirm", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to check password confirm: %w", err)
	}

	userID, err := u.userService.Create(ctx, mapper.ToUserFromCreateRequest(request))
	if err != nil {
		logger.Error("failed to create user", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return mapper.ToCreateResponseFromUserService(userID), nil
}

func (u *Implementation) checkPasswordConfirm(request *desc.CreateRequest) error {
	if request.Password != request.PasswordConfirm {
		return fmt.Errorf("password and confirm password do not match")
	}

	return nil
}
