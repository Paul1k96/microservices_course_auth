package v1

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Paul1k96/microservices_course_auth/pkg/proto/gen/user_v1"
	desc "github.com/Paul1k96/microservices_course_auth/pkg/proto/gen/user_v1"
)

// Delete user by id.
func (u *Implementation) Delete(ctx context.Context, request *desc.DeleteRequest) (*desc.DeleteResponse, error) {
	logger := u.logger.
		With("method", "Delete").
		With("user_id", request.Id)

	err := u.userService.Delete(ctx, request.Id)
	if err != nil {
		logger.Error("failed to delete user", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to delete user: %w", err)
	}

	return &user_v1.DeleteResponse{}, nil
}
