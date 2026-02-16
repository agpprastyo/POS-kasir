package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	Querier
	ExecTx(ctx context.Context, fn func(Querier) error) error
}

type SQLStore struct {
	*Queries
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) Store {
	return &SQLStore{
		Queries: New(db),
		db:      db,
	}
}

func (store *SQLStore) ExecTx(ctx context.Context, fn func(Querier) error) error {
	tx, err := store.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		} else if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	q := store.WithTx(tx)
	err = fn(q)
	return err
}
