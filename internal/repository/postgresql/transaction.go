package postgresql

import (
	"context"

	"github.com/Ranik23/avito-tech-spring/internal/repository"
	"github.com/jackc/pgx/v5"
)


type postgresTransaction struct {
	Tx 	pgx.Tx
}


func NewTransaction(tx pgx.Tx) repository.Transaction {
	return &postgresTransaction{
		Tx: tx,
	}
}

func (t *postgresTransaction) Transaction() interface{} {
	return t.Tx
}


func (t *postgresTransaction) Commit(ctx context.Context) error {
	return t.Tx.Commit(ctx)
}

func (t *postgresTransaction) Rollback(ctx context.Context) error {
	return t.Tx.Rollback(ctx)
}

