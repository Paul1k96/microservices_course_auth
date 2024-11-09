package testmodel

import (
	"time"

	"github.com/Paul1k96/microservices_course_auth/internal/model"
	"github.com/brianvoe/gofakeit/v7"
)

// NewUser creates a new User instance
func NewUser() *model.User {
	m := struct {
		ID        int64
		Name      string
		Email     string `fake:"{email}"`
		Password  string
		Role      model.Role `fake:"{number:0,2}"`
		CreatedAt time.Time
		UpdatedAt *time.Time
	}{}

	_ = gofakeit.Struct(&m)
	res := model.User(m)
	return &res
}

// NewUsers creates a slice of User instances
func NewUsers(quantity int) []*model.User {
	res := make([]*model.User, 0, quantity)
	for i := 0; i < quantity; i++ {
		res = append(res, NewUser())
	}

	return res
}
