package user

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/Paul1k96/microservices_course_auth/internal/errs"
	"github.com/Paul1k96/microservices_course_auth/internal/model"
	"github.com/Paul1k96/microservices_course_auth/internal/repository/user/mapper"
	modelRepo "github.com/Paul1k96/microservices_course_auth/internal/repository/user/model"
	"github.com/Paul1k96/microservices_course_platform_common/pkg/client/db"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

const (
	userTable = "users"

	idColumn        = "id"
	nameColumn      = "name"
	emailColumn     = "email"
	passwordColumn  = "password"
	roleColumn      = "role"
	createdAtColumn = "created_at"
	updateAtColumn  = "updated_at"
)

// Repository represents user repository.
type Repository struct {
	db db.DB
}

// NewRepository creates a new instance of repository.UsersRepository.
func NewRepository(pg db.DB) *Repository {
	return &Repository{db: pg}
}

// Create user.
func (r *Repository) Create(ctx context.Context, user *modelRepo.User) (int64, error) {
	queryBuilder := sq.Insert(userTable).
		PlaceholderFormat(sq.Dollar).
		Columns(nameColumn, emailColumn, passwordColumn, roleColumn).
		Values(user.Name, user.Email, user.Password, user.Role).
		Suffix("RETURNING id")

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return 0, fmt.Errorf("build query: %w", err)
	}

	q := db.Query{
		Name:     "user_repository.Create",
		QueryRaw: query,
	}

	var userID int64
	err = r.db.QueryRowContext(ctx, q, args...).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("exec query: %w", err)
	}

	return userID, nil
}

// GetByID user by id.
func (r *Repository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	queryBuilder := sq.Select("*").
		PlaceholderFormat(sq.Dollar).
		From(userTable).
		Where(sq.Eq{idColumn: id})

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	q := db.Query{
		Name:     "user_repository.Get",
		QueryRaw: query,
	}

	var user modelRepo.User
	err = r.db.ScanOneContext(ctx, &user, q, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("get user: %w", errs.ErrUserNotFound)
		}

		return nil, fmt.Errorf("get user: %w", err)
	}

	return mapper.ToUserFromRepo(&user), nil
}

// Update user by id.
// If user.Name or user.Email is empty, this field will not be updated.
func (r *Repository) Update(ctx context.Context, user *modelRepo.User) error {
	queryBuilder := r.setUserDataForUpdate(user)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}

	q := db.Query{
		Name:     "user_repository.Update",
		QueryRaw: query,
	}

	_, err = r.db.ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	return nil
}

func (r *Repository) setUserDataForUpdate(user *modelRepo.User) sq.UpdateBuilder {
	queryBuilder := sq.Update(userTable).
		PlaceholderFormat(sq.Dollar).
		Set(roleColumn, user.Role).
		Set(updateAtColumn, user.UpdatedAt).
		Where(sq.Eq{idColumn: user.ID})

	if user.Name != "" {
		queryBuilder = queryBuilder.Set("Name", user.Name)
	}
	if user.Email != "" {
		queryBuilder = queryBuilder.Set("Email", user.Email)
	}

	return queryBuilder
}

// Delete user by id.
func (r *Repository) Delete(ctx context.Context, id int64) error {
	queryBuilder := sq.Delete(userTable).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id})

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}

	q := db.Query{
		Name:     "user_repository.Delete",
		QueryRaw: query,
	}

	_, err = r.db.ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	return nil
}
