package model

import "time"

// User represents user model.
type User struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Role      string     `json:"role"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}
