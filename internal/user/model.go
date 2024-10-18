package user

import (
	"time"
)

//go:generate go run github.com/dmarkham/enumer -transform snake-upper -trimprefix Role -type Role -output role_string.go model.go

// User represents user model.
type User struct {
	ID        int64
	Name      string
	Email     string
	Password  string
	Role      Role
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Role represents user role.
type Role int32

// User roles.
const (
	RoleUnspecified Role = iota // default value
	RoleAdmin
	RoleUser
)
