package model

import (
	"time"
)

//go:generate ../../bin/enumer -transform snake-upper -trimprefix Role -type Role -output user_role_string.go user.go

// User roles.
const (
	RoleUnknown Role = iota // default value
	RoleAdmin
	RoleUser
)

// Role represents user role.
type Role int32

// User represents user model.
type User struct {
	ID        int64
	Name      string
	Email     string
	Password  string
	Role      Role
	CreatedAt time.Time
	UpdatedAt *time.Time
}
