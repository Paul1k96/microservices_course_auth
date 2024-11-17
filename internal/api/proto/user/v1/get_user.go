package v1

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Paul1k96/microservices_course_auth/internal/mapper"
	desc "github.com/Paul1k96/microservices_course_auth/pkg/proto/gen/user_v1"
)

// Get user by id.
func (u *Implementation) Get(ctx context.Context, request *desc.GetRequest) (*desc.GetResponse, error) {
	logger := u.logger.
		With("method", "Get").
		With("user_id", request.Id)

	user, err := u.userService.GetByID(ctx, request.Id)
	if err != nil {
		logger.Error("failed to get user", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return mapper.ToGetResponseFromUserService(user), nil
}
