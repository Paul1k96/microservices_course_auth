package user

import (
	"context"
	"fmt"
	"net/mail"

	"github.com/Paul1k96/microservices_course_auth/internal/errs"
	"github.com/Paul1k96/microservices_course_auth/internal/model"
	"github.com/Paul1k96/microservices_course_auth/internal/repository/user/mapper"
	"github.com/pkg/errors"
)

// Create creates a new user.
func (s *service) Create(ctx context.Context, user *model.User) (int64, error) {
	var id int64

	if txErr := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var err error

		if err = s.validateUser(ctx, user); err != nil {
			return fmt.Errorf("create user: %w", err)
		}

		id, err = s.repo.Create(ctx, mapper.ToRepoCreateFromUserService(user))
		if err != nil {
			return fmt.Errorf("create user: %w", err)
		}

		return nil
	}); txErr != nil {
		return 0, fmt.Errorf("transaction error: %w", txErr)
	}

	return id, nil
}

func (s *service) validateUser(ctx context.Context, user *model.User) error {
	err := s.validateUserName(user.Name)
	if err != nil {
		return fmt.Errorf("validate user: %w", err)
	}

	err = s.validateUserEmail(ctx, user.Email)
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
	err := s.checkUserNameLength([]rune(name))
	if err != nil {
		return fmt.Errorf("validate name: %w", err)
	}

	err = s.checkRestrictedSymbols(name)
	if err != nil {
		return fmt.Errorf("validate name: %w", err)
	}

	return nil
}

func (s *service) checkUserNameLength(name []rune) error {
	if len(name) == 0 {
		return fmt.Errorf("name is required")
	}

	if len(name) < 2 {
		return fmt.Errorf("name must be at least 2 characters")
	}

	if len(name) > 100 {
		return fmt.Errorf("name must be at most 100 characters")
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

func (s *service) validateUserEmail(ctx context.Context, email string) error {
	err := s.checkUserEmailLength(email)
	if err != nil {
		return fmt.Errorf("validate email: %w", err)
	}

	err = s.checkUserEmailFormat(email)
	if err != nil {
		return fmt.Errorf("validate email: %w", err)
	}

	err = s.checkUserExistsByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("validate user: %w", err)
	}

	return nil
}

func (s *service) checkUserEmailLength(email string) error {
	if len(email) == 0 {
		return fmt.Errorf("email is required")
	}

	if len(email) < 5 {
		return fmt.Errorf("email must be at least 5 characters")
	}

	if len(email) > 100 {
		return fmt.Errorf("email must be at most 100 characters")
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

func (s *service) checkUserExistsByEmail(ctx context.Context, email string) error {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil && !errors.Is(err, errs.ErrUserNotFound) {
		return fmt.Errorf("check user exists: %w", err)
	}

	if user != nil {
		return fmt.Errorf("user with email %s already exists", email)
	}

	return nil
}
