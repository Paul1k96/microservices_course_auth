package user

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

const userTable = "users"

type userRaw struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	Role      string    `db:"role"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (u userRaw) toDomain() User {
	var user User

	user.ID = u.ID
	user.Name = u.Name
	user.Email = u.Email
	user.Password = u.Password
	user.Role = u.getRole()
	user.CreatedAt = u.CreatedAt
	user.UpdatedAt = u.UpdatedAt

	return user
}

func (u userRaw) getRole() Role {
	switch u.Role {
	case "ADMIN":
		return RoleAdmin
	case "USER":
		return RoleUser
	default:
		return RoleUnspecified
	}
}

// Repository represents user repository.
type Repository struct {
	pg *sqlx.DB
}

// NewUserRepository creates a new instance of Repository.
func NewUserRepository(pg *sqlx.DB) *Repository {
	return &Repository{pg: pg}
}

// Create user.
func (u *Repository) Create(ctx context.Context, user User) (*int, error) {
	queryBuilder := sq.Insert(userTable).
		PlaceholderFormat(sq.Dollar).
		Columns("name", "email", "password", "role").
		Values(user.Name, user.Email, user.Password, user.Role.String()).
		Suffix("RETURNING id")

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var userID int
	err = u.pg.QueryRowxContext(ctx, query, args...).Scan(&userID)
	if err != nil {
		return nil, fmt.Errorf("exec query: %w", err)
	}

	return &userID, nil
}

// Get user by id.
func (u *Repository) Get(ctx context.Context, id int64) (User, error) {
	queryBuilder := sq.Select("*").
		PlaceholderFormat(sq.Dollar).
		From(userTable).
		Where(sq.Eq{"id": id})

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return User{}, fmt.Errorf("build query: %w", err)
	}

	var user userRaw
	err = u.pg.GetContext(ctx, &user, query, args...)
	if err != nil {
		return User{}, fmt.Errorf("get user: %w", err)
	}

	return user.toDomain(), nil
}

// Update user by id.
// If user.Name or user.Email is empty, this field will not be updated.
func (u *Repository) Update(ctx context.Context, id int64, user User) error {
	queryBuilder := u.setUserDataForUpdate(id, user)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}

	_, err = u.pg.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	return nil
}

func (u *Repository) setUserDataForUpdate(id int64, user User) sq.UpdateBuilder {
	queryBuilder := sq.Update(userTable).
		PlaceholderFormat(sq.Dollar).
		Set("Role", user.Role.String()).
		Where(sq.Eq{"id": id})

	if user.Name != "" {
		queryBuilder = queryBuilder.Set("Name", user.Name)
	}
	if user.Email != "" {
		queryBuilder = queryBuilder.Set("Email", user.Email)
	}

	return queryBuilder
}

// Delete user by id.
func (u *Repository) Delete(ctx context.Context, id int64) error {
	queryBuilder := sq.Delete(userTable).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"ID": id})

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}

	_, err = u.pg.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	return nil
}
