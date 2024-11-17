package model

import (
	"time"

	"github.com/google/uuid"
)

// Allowed UserEventType.
const (
	UserEventTypeUnspecified UserEventType = iota
	UserEventTypeCreateUser
	UserEventTypeUpdateUser
	UserEventTypeDeleteUser
)

// UserEventType represents user event type.
type UserEventType int64

// UserEventValue represents user event value.
type UserEventValue interface {
	Value() interface{}
}

// CreateUserEventValue represents create user event value.
type CreateUserEventValue struct {
	User *User
}

// Value returns value.
func (v CreateUserEventValue) Value() interface{} {
	return &v.User
}

// UpdateUserEventValue represents update user event value.
type UpdateUserEventValue struct {
	User *User
}

// Value returns value.
func (v UpdateUserEventValue) Value() interface{} {
	return &v.User
}

// DeleteUserEventValue represents delete user event value.
type DeleteUserEventValue struct{}

// Value returns value.
func (v DeleteUserEventValue) Value() interface{} {
	return nil
}

// UserEvent represents user event model.
type UserEvent struct {
	ID        uuid.UUID
	UserID    int64
	Type      UserEventType
	EntityID  int64
	Value     UserEventValue
	CreatedAt time.Time
}

// NewUserEvent creates a new user event.
func NewUserEvent(userID, entityID int64, eventType UserEventType, value UserEventValue) *UserEvent {
	return &UserEvent{
		ID:        uuid.New(),
		UserID:    userID,
		Type:      eventType,
		EntityID:  entityID,
		Value:     value,
		CreatedAt: time.Now(),
	}
}

// NewUnspecifiedUserEvent creates a new unspecified user event.
func NewUnspecifiedUserEvent(userID, entityID int64) *UserEvent {
	return NewUserEvent(userID, entityID, UserEventTypeUnspecified, nil)
}

// NewCreateUserEvent creates a new create user event.
func NewCreateUserEvent(userID int64, user *User) *UserEvent {
	return NewUserEvent(userID, user.ID, UserEventTypeCreateUser, &CreateUserEventValue{User: user})
}

// NewUpdateUserEvent creates a new update user event.
func NewUpdateUserEvent(userID int64, user *User) *UserEvent {
	return NewUserEvent(userID, user.ID, UserEventTypeUpdateUser, &UpdateUserEventValue{User: user})
}

// NewDeleteUserEvent creates a new delete user event.
func NewDeleteUserEvent(userID, entityID int64) *UserEvent {
	return NewUserEvent(userID, entityID, UserEventTypeDeleteUser, &DeleteUserEventValue{})
}
