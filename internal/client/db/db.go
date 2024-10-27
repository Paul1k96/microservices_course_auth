package db

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

// Handler is a function that is executed in a transaction.
type Handler func(ctx context.Context) error

// Client is a client for working with the database.
type Client interface {
	DB() DB
	Close() error
}

// TxManager is a transaction manager that executes the handler specified by the user in a transaction.
type TxManager interface {
	ReadCommitted(ctx context.Context, f Handler) error
}

// Query is a wrapper over a query that stores the query name and the query itself.
type Query struct {
	Name     string
	QueryRaw string
}

// Transactor is an interface for working with transactions.
type Transactor interface {
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

// SQLExecutor combines NamedExecutor and QueryExecutor.
type SQLExecutor interface {
	NamedExecutor
	QueryExecutor
}

// NamedExecutor is an interface for working with named queries using tags in structures.
type NamedExecutor interface {
	ScanOneContext(ctx context.Context, dest interface{}, q Query, args ...interface{}) error
	ScanAllContext(ctx context.Context, dest interface{}, q Query, args ...interface{}) error
}

// QueryExecutor is an interface for working with regular queries.
type QueryExecutor interface {
	ExecContext(ctx context.Context, q Query, args ...interface{}) (pgconn.CommandTag, error)
	QueryContext(ctx context.Context, q Query, args ...interface{}) (pgx.Rows, error)
	QueryRowContext(ctx context.Context, q Query, args ...interface{}) pgx.Row
}

// Pinger is an interface for checking the connection to the database.
type Pinger interface {
	Ping(ctx context.Context) error
}

// DB is an interface for working with the database.
type DB interface {
	SQLExecutor
	Transactor
	Pinger
	Close()
}
