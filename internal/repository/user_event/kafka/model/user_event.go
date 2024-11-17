package model

import (
	"github.com/Paul1k96/microservices_course_auth/internal/model"
)

// Data represents event data.
type Data interface {
	isEventData()
}

// UserEvent represents user event model.
type UserEvent struct {
	ID        string              `json:"id"`
	UserID    int64               `json:"user_id"`
	Type      model.UserEventType `json:"type"`
	Data      Data                `json:"data"`
	EntityID  int64               `json:"entity_id"`
	CreatedAt int64               `json:"created_at"`
}

// CreateUserEventData represents create user event data.
type CreateUserEventData struct {
	User *User `json:"user"`
}

func (CreateUserEventData) isEventData() {}

// UpdateUserEventData represents update user event data.
type UpdateUserEventData struct {
	User *User `json:"user"`
}

func (UpdateUserEventData) isEventData() {}

// DeleteUserEventData represents delete user event data.
type DeleteUserEventData struct {
}

func (DeleteUserEventData) isEventData() {}
