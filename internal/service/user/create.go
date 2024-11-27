package user

import (
	"context"
	"fmt"
	"log/slog"
	"net/mail"

	"github.com/Paul1k96/microservices_course_auth/internal/model"
)

// Create creates a new user.
func (s *service) Create(ctx context.Context, user *model.User) (int64, error) {
	if err := s.validateUser(user); err != nil {
		return 0, fmt.Errorf("create user: %w", err)
	}

	id, err := s.repo.Create(ctx, user)
	if err != nil {
		return 0, fmt.Errorf("create user: %w", err)
	}

	user.ID = id
	err = s.cache.Set(ctx, user)
	if err != nil {
		s.logger.Error("failed to set user to cache:", slog.String("error", err.Error()))
	}

	err = s.events.Save(ctx, model.NewCreateUserEvent(id, user))
	if err != nil {
		s.logger.Error("failed to save user event:", slog.String("error", err.Error()))
	}

	return id, nil
}

func (s *service) validateUser(user *model.User) error {
	err := s.validateUserName(user.Name)
	if err != nil {
		return fmt.Errorf("validate user: %w", err)
	}

	err = s.validateUserEmail(user.Email)
	if err != nil {
		return fmt.Errorf("validate user: %w", err)
	}

	err = s.validateUserRole(user.Role)
	if err != nil {
		return fmt.Errorf("validate user: %w", err)
	}

	return nil
}

func (s *service) validateUserName(name string) error {
	err := s.checkRestrictedSymbols(name)
	if err != nil {
		return fmt.Errorf("validate name: %w", err)
	}

	return nil
}

func (s *service) checkRestrictedSymbols(name string) error {
	restrictedChars := map[rune]struct{}{
		'!': {}, '@': {}, '#': {}, '$': {}, '%': {}, '^': {}, '&': {}, '*': {}, '(': {}, ')': {},
	}

	for _, char := range name {
		if _, ok := restrictedChars[char]; ok {
			return fmt.Errorf("name contains restricted symbols")
		}
	}

	return nil
}

func (s *service) validateUserEmail(email string) error {
	err := s.checkUserEmailFormat(email)
	if err != nil {
		return fmt.Errorf("validate email: %w", err)
	}

	return nil
}

func (s *service) checkUserEmailFormat(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("email has incorrect format")
	}

	return nil
}

func (s *service) validateUserRole(role model.Role) error {
	if role.IsARole() {
		return nil
	}

	return fmt.Errorf("role is not valid")
}
