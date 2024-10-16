package user

import (
	"time"

	"github.com/Paul1k96/microservices_course_auth/pkg/proto/gen/user_v1"
)

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
