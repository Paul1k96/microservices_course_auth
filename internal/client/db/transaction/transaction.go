package transaction

import (
	"context"
	"fmt"

	"github.com/Paul1k96/microservices_course_auth/internal/client/db"
	"github.com/Paul1k96/microservices_course_auth/internal/client/db/pg"
	"github.com/jackc/pgx/v4"
)

type manager struct {
	db db.Transactor
}

// NewTransactionManager creates a new transaction manager that satisfies the db.TxManager interface
func NewTransactionManager(db db.Transactor) db.TxManager {
	return &manager{
		db: db,
	}
}

// ReadCommitted executes the user-specified handler in a read committed transaction
func (m *manager) ReadCommitted(ctx context.Context, fn db.Handler) error {
	txOpts := pgx.TxOptions{IsoLevel: pgx.ReadCommitted}
	return m.transaction(ctx, txOpts, fn)
}

// transaction is the main function that executes the user-specified handler in a transaction
func (m *manager) transaction(ctx context.Context, opts pgx.TxOptions, fn db.Handler) (err error) {
	// If this is a nested transaction, skip initiating a new transaction and execute
	tx, ok := ctx.Value(pg.TxKey).(pgx.Tx)
	if ok {
		return fn(ctx)
	}

	// Start a new transaction
	tx, err = m.db.BeginTx(ctx, opts)
	if err != nil {
		return fmt.Errorf("can't begin transaction: %w", err)
	}

	ctx = pg.MakeTxContextTx(ctx, tx)

	// Set up a deferred function for rolling back or committing the transaction
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic recovered: %v", r)
		}

		if err != nil {
			if errRollback := tx.Rollback(ctx); errRollback != nil {
				err = fmt.Errorf("errRollback: %w", errRollback)
			}

			return
		}

		if nil == err {
			err = tx.Commit(ctx)
			if err != nil {
				err = fmt.Errorf("can't commit transaction: %w", err)
			}
		}
	}()

	if err = fn(ctx); err != nil {
		err = fmt.Errorf("tx handler failed: %w", err)
	}

	return err
}
