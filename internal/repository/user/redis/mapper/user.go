package mapper

import (
	"time"

	"github.com/Paul1k96/microservices_course_auth/internal/model"
	modelRepo "github.com/Paul1k96/microservices_course_auth/internal/repository/user/redis/model"
)

// ToUserFromRepo converts user from repository model to service model.
func ToUserFromRepo(user *modelRepo.User) *model.User {
	var serviceUser model.User

	serviceUser.ID = user.ID
	serviceUser.Name = user.Name
	serviceUser.Email = user.Email
	serviceUser.Role = ToRoleFromRepo(user.Role)
	serviceUser.CreatedAt = time.Unix(0, user.CreatedAt)
	if user.UpdatedAt != nil {
		updateTime := time.Unix(0, *user.UpdatedAt)
		serviceUser.UpdatedAt = &updateTime
	}

	return &serviceUser
}

// ToRoleFromRepo converts role from repository model to service model.
func ToRoleFromRepo(role string) model.Role {
	modelRole, err := model.RoleString(role)
	if err != nil {
		return model.RoleUnknown
	}

	return modelRole
}
