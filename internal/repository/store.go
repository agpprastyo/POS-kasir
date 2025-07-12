package repository

import (
	"POS-kasir/pkg/logger"
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	Querier
	ExecTx(ctx context.Context, fn func(*Queries) error) error
}

type SQLStore struct {
	*Queries
	db  *pgxpool.Pool
	log *logger.Logger
}

func NewStore(db *pgxpool.Pool, log *logger.Logger) Store {
	return &SQLStore{
		Queries: New(db),
		db:      db,
		log:     log,
	}
}

func (store *SQLStore) ExecTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			store.log.Error("Failed to rollback transaction", "error", err)
		}
	}(tx, ctx)

	q := New(tx)
	err = fn(q)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
