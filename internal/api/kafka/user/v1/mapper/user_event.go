package mapper

import (
	"encoding/json"

	modelKafka "github.com/Paul1k96/microservices_course_auth/internal/api/kafka/user/v1/model"
	"github.com/Paul1k96/microservices_course_auth/internal/model"
	modelRepoKafka "github.com/Paul1k96/microservices_course_auth/internal/repository/user_event/kafka/model"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// ToUserEventFromKafka creates user event from kafka.
func ToUserEventFromKafka(event *modelKafka.UserEvent) (*model.UserEvent, error) {
	var value model.UserEventValue

	eventID, err := uuid.Parse(event.ID)
	if err != nil {
		return nil, err
	}

	switch event.Type {
	case model.UserEventTypeCreateUser:
		value, err = ToCreateUserEventValueFromKafka(event.Data)
		if err != nil {
			return nil, err
		}
	case model.UserEventTypeUpdateUser:
		value, err = ToUpdateUserEventValueFromKafka(event.Data)
		if err != nil {
			return nil, err
		}
	case model.UserEventTypeDeleteUser:
		value = &model.DeleteUserEventValue{}
	default:
		return nil, errors.New("invalid user event data")
	}

	return &model.UserEvent{
		ID:        eventID,
		UserID:    event.UserID,
		Type:      event.Type,
		EntityID:  event.EntityID,
		Value:     value,
		CreatedAt: event.CreatedAt,
	}, nil
}

// ToCreateUserEventValueFromKafka creates create user event value from kafka.
func ToCreateUserEventValueFromKafka(data json.RawMessage) (*model.CreateUserEventValue, error) {
	var kafkaModel modelRepoKafka.CreateUserEventData
	if err := json.Unmarshal(data, &kafkaModel); err != nil {
		return nil, err
	}

	return &model.CreateUserEventValue{
		User: &model.User{
			ID:        kafkaModel.User.ID,
			Name:      kafkaModel.User.Name,
			Email:     kafkaModel.User.Email,
			Role:      ToRoleFromKafka(kafkaModel.User.Role),
			CreatedAt: kafkaModel.User.CreatedAt,
			UpdatedAt: kafkaModel.User.UpdatedAt,
		},
	}, nil
}

// ToUpdateUserEventValueFromKafka creates update user event value from kafka.
func ToUpdateUserEventValueFromKafka(data json.RawMessage) (*model.UpdateUserEventValue, error) {
	var kafkaModel modelRepoKafka.UpdateUserEventData
	if err := json.Unmarshal(data, &kafkaModel); err != nil {
		return nil, err
	}

	return &model.UpdateUserEventValue{
		User: &model.User{
			ID:        kafkaModel.User.ID,
			Name:      kafkaModel.User.Name,
			Email:     kafkaModel.User.Email,
			Role:      ToRoleFromKafka(kafkaModel.User.Role),
			CreatedAt: kafkaModel.User.CreatedAt,
			UpdatedAt: kafkaModel.User.UpdatedAt,
		},
	}, nil
}

// ToRoleFromKafka creates role from kafka.
func ToRoleFromKafka(role string) model.Role {
	switch role {
	case model.RoleAdmin.String():
		return model.RoleAdmin
	case model.RoleUser.String():
		return model.RoleUser
	default:
		return model.RoleUnknown
	}
}
