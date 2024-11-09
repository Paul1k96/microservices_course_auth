package mapper

import (
	"github.com/Paul1k96/microservices_course_auth/internal/model"
	modelRepo "github.com/Paul1k96/microservices_course_auth/internal/repository/user/pg/model"
)

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

// ToRepoUpdateFromUserService converts user from service model to repository model.
func ToRepoUpdateFromUserService(user *model.User) *modelRepo.User {
	var repoUser modelRepo.User

	repoUser.ID = user.ID
	repoUser.Name = user.Name
	repoUser.Email = user.Email
	repoUser.Role = user.Role.String()
	if user.UpdatedAt != nil {
		repoUser.UpdatedAt.Time = *user.UpdatedAt
		repoUser.UpdatedAt.Valid = true
	}

	return &repoUser
}

// ToRoleFromRepo converts role from repository model to service model.
func ToRoleFromRepo(role string) model.Role {
	modelRole, err := model.RoleString(role)
	if err != nil {
		return model.RoleUnknown
	}

	return modelRole
}
