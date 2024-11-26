package userv1

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Paul1k96/microservices_course_auth/internal/mapper"
	"github.com/Paul1k96/microservices_course_auth/pkg/proto/gen/user_v1"
	desc "github.com/Paul1k96/microservices_course_auth/pkg/proto/gen/user_v1"
)

// Update user fields.
func (u *Implementation) Update(ctx context.Context, request *desc.UpdateRequest) (*desc.UpdateResponse, error) {
	logger := u.logger.
		With("method", "Update").
		With("user_id", request.Id)

	err := u.userService.Update(ctx, mapper.ToUserFromUpdateRequest(request))
	if err != nil {
		logger.Error("failed to update user", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &user_v1.UpdateResponse{}, nil
}
