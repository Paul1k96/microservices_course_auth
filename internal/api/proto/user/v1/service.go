package v1

import (
	"log/slog"

	"github.com/Paul1k96/microservices_course_auth/internal/service"
	"github.com/Paul1k96/microservices_course_auth/pkg/proto/gen/user_v1"
)

// Implementation of the user service.
type Implementation struct {
	logger      *slog.Logger
	userService service.UserService
	user_v1.UnimplementedUserServer
}

// NewImplementation creates a new user service implementation.
func NewImplementation(logger *slog.Logger, userService service.UserService) *Implementation {
	return &Implementation{logger: logger, userService: userService}
}
