package mapper

import (
	"github.com/Paul1k96/microservices_course_auth/internal/model"
	modelRepo "github.com/Paul1k96/microservices_course_auth/internal/repository/user/pg/model"
)

// ToUsersFromRepo converts users from repository model to service model.
func ToUsersFromRepo(users []*modelRepo.User) []*model.User {
	serviceUsers := make([]*model.User, 0, len(users))

	for _, user := range users {
		serviceUsers = append(serviceUsers, ToUserFromRepo(user))
	}

	return serviceUsers
}

// ToUserFromRepo converts user from repository model to service model.
func ToUserFromRepo(user *modelRepo.User) *model.User {
	var serviceUser model.User

	serviceUser.ID = user.ID
	serviceUser.Name = user.Name
	serviceUser.Email = user.Email
	serviceUser.Password = user.Password
	serviceUser.Role = ToRoleFromRepo(user.Role)
	serviceUser.CreatedAt = user.CreatedAt
	if user.UpdatedAt.Valid {
		serviceUser.UpdatedAt = &user.UpdatedAt.Time
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
