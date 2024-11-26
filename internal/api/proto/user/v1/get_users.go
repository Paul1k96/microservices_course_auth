package userv1

import (
	"context"
	"fmt"

	"github.com/Paul1k96/microservices_course_auth/internal/mapper"
	desc "github.com/Paul1k96/microservices_course_auth/pkg/proto/gen/user_v1"
)

// List returns list of users.
func (u *Implementation) List(ctx context.Context, request *desc.GetListRequest) (*desc.GetListResponse, error) {
	logger := u.logger.
		With("method", "GetList").
		With("ids", request.GetIds())

	users, err := u.userService.GetListByIDs(ctx, request.GetIds())
	if err != nil {
		logger.Error("failed to get users", "error", err)
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	return mapper.ToGetListResponseFromUserService(users), nil
}
