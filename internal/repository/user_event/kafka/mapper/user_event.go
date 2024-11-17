package mapper

import (
	"github.com/Paul1k96/microservices_course_auth/internal/model"
	modelKafka "github.com/Paul1k96/microservices_course_auth/internal/repository/user_event/kafka/model"
	"github.com/pkg/errors"
)

// NewUserEvent creates user event.
func NewUserEvent(userEvent *model.UserEvent) (*modelKafka.UserEvent, error) {
	var data modelKafka.Data

	switch val := userEvent.Value.(type) {
	case *model.CreateUserEventValue:
		data = NewCreateUserEventData(val)
	case *model.UpdateUserEventValue:
		data = NewUpdateUserEventData(val)
	case *model.DeleteUserEventValue:
		data = NewDeleteUserEventData(val)
	default:
		return nil, errors.New("invalid user event data")
	}

	return &modelKafka.UserEvent{
		ID:       userEvent.ID.String(),
		UserID:   userEvent.UserID,
		Type:     userEvent.Type,
		EntityID: userEvent.EntityID,
		Data:     data,
	}, nil
}

// NewCreateUserEventData creates create user event data.
func NewCreateUserEventData(userVal *model.CreateUserEventValue) *modelKafka.CreateUserEventData {
	return &modelKafka.CreateUserEventData{
		User: &modelKafka.User{
			ID:        userVal.User.ID,
			Name:      userVal.User.Name,
			Email:     userVal.User.Email,
			Role:      userVal.User.Role.String(),
			CreatedAt: userVal.User.CreatedAt,
			UpdatedAt: userVal.User.UpdatedAt,
		},
	}
}

// NewUpdateUserEventData creates update user event data.
func NewUpdateUserEventData(userVal *model.UpdateUserEventValue) *modelKafka.UpdateUserEventData {
	return &modelKafka.UpdateUserEventData{
		User: &modelKafka.User{
			ID:        userVal.User.ID,
			Name:      userVal.User.Name,
			Email:     userVal.User.Email,
			Role:      userVal.User.Role.String(),
			CreatedAt: userVal.User.CreatedAt,
			UpdatedAt: userVal.User.UpdatedAt,
		},
	}
}

// NewDeleteUserEventData creates delete user event data.
func NewDeleteUserEventData(_ *model.DeleteUserEventValue) *modelKafka.DeleteUserEventData {
	return &modelKafka.DeleteUserEventData{}
}
