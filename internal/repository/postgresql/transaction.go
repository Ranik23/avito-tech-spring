package postgresql

import (
	"context"

	"github.com/Ranik23/avito-tech-spring/internal/repository"
	"github.com/jackc/pgx/v5"
)


type postgresTransaction struct {
	tx 	pgx.Tx
}


func NewTransaction(tx pgx.Tx) repository.Transaction {
	return &postgresTransaction{
		tx: tx,
	}
}

func (t *postgresTransaction) Commit(ctx context.Context) error {
	return t.tx.Commit(ctx)
}

func (t *postgresTransaction) Rollback(ctx context.Context) error {
	return t.tx.Rollback(ctx)
}

