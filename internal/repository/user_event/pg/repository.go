package pg

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/Paul1k96/microservices_course_auth/internal/model"
	"github.com/Paul1k96/microservices_course_platform_common/pkg/client/db"
	"github.com/jackc/pgtype"
)

const (
	userEventTable = "user_events"

	idColumn     = "id"
	userIDColumn = "user_id"
	eventType    = "event_type"
	eventValue   = "value"
	entityID     = "entity_id"
	createdAt    = "created_at"
)

// Repository is a user event repository.
type Repository struct {
	db db.DB
}

// NewRepository creates a new instance of Repository.
func NewRepository(db db.DB) *Repository {
	return &Repository{db: db}
}

// Save saves user event.
func (r *Repository) Save(ctx context.Context, event *model.UserEvent) error {
	var rawValue pgtype.JSONB

	if err := rawValue.Set(event.Value.Value()); err != nil {
		return fmt.Errorf("set raw value: %w", err)
	}

	queryBuilder := sq.Insert(userEventTable).
		PlaceholderFormat(sq.Dollar).
		Columns(idColumn, userIDColumn, eventType, eventValue, entityID, createdAt).
		Values(event.ID, event.UserID, event.Type, rawValue, event.EntityID, event.CreatedAt)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}

	q := db.Query{
		Name:     "user_event_repository.Save",
		QueryRaw: query,
	}

	_, err = r.db.ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	return nil
}
