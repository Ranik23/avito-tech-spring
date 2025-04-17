package manager

import (
	"context"

	"github.com/jackc/pgx/v5"
)


type Transaction interface {
    Commit(context.Context) error 
    Rollback(context.Context) error 
 }

 type transaction struct {
	tx 	pgx.Tx
}

func (t *transaction) Commit(ctx context.Context) error {
	return t.tx.Commit(ctx)
}

func (t *transaction) Rollback(ctx context.Context) error {
	return t.tx.Rollback(ctx)
}

func NewTransaction(tx pgx.Tx) Transaction {
	return &transaction{
		tx: tx,
	}
}