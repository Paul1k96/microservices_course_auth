package pg

import (
	"context"
	"fmt"

	"github.com/Paul1k96/microservices_course_auth/internal/client/db"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type key string

const (
	// TxKey is a key for the transaction
	TxKey key = "tx"
)

type pg struct {
	dbc *pgxpool.Pool
}

// NewDB creates a new db connection
func NewDB(dbc *pgxpool.Pool) db.DB {
	return &pg{
		dbc: dbc,
	}
}

// ScanOneContext scans one row from the db
func (p *pg) ScanOneContext(ctx context.Context, dest interface{}, q db.Query, args ...interface{}) error {
	row, err := p.QueryContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	return pgxscan.ScanOne(dest, row)
}

// ScanAllContext scans all rows from the db
func (p *pg) ScanAllContext(ctx context.Context, dest interface{}, q db.Query, args ...interface{}) error {
	rows, err := p.QueryContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	return pgxscan.ScanAll(dest, rows)
}

// ExecContext executes a query
func (p *pg) ExecContext(ctx context.Context, q db.Query, args ...interface{}) (pgconn.CommandTag, error) {
	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.Exec(ctx, q.QueryRaw, args...)
	}

	return p.dbc.Exec(ctx, q.QueryRaw, args...)
}

// QueryContext executes a query and returns the result
func (p *pg) QueryContext(ctx context.Context, q db.Query, args ...interface{}) (pgx.Rows, error) {
	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.Query(ctx, q.QueryRaw, args...)
	}

	return p.dbc.Query(ctx, q.QueryRaw, args...)
}

// QueryRowContext executes a query and returns one row
func (p *pg) QueryRowContext(ctx context.Context, q db.Query, args ...interface{}) pgx.Row {
	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.QueryRow(ctx, q.QueryRaw, args...)
	}

	return p.dbc.QueryRow(ctx, q.QueryRaw, args...)
}

// BeginTx starts a new transaction
func (p *pg) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	return p.dbc.BeginTx(ctx, txOptions)
}

// Ping checks the connection to the db
func (p *pg) Ping(ctx context.Context) error {
	return p.dbc.Ping(ctx)
}

// Close closes the db connection
func (p *pg) Close() {
	p.dbc.Close()
}

// MakeTxContextTx adds a transaction to the context
func MakeTxContextTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, TxKey, tx)
}
