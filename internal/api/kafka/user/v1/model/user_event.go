package model

import (
	"encoding/json"

	"github.com/Paul1k96/microservices_course_auth/internal/model"
)

// UserEvent represents user event model.
type UserEvent struct {
	ID        string              `json:"id"`
	UserID    int64               `json:"user_id"`
	Type      model.UserEventType `json:"type"`
	Data      json.RawMessage     `json:"data"`
	EntityID  int64               `json:"entity_id"`
	CreatedAt int64               `json:"created_at"`
}
