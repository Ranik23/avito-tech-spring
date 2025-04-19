package postgresql

import (
	"context"

	"github.com/Ranik23/avito-tech-spring/internal/repository"
	pool "github.com/jackc/pgx/v5/pgxpool"
)


type postgresCtxManager struct {
	pool *pool.Pool
}

func NewCtxManager(pool *pool.Pool) repository.CtxManager {
	return &postgresCtxManager{
		pool: pool,
	}
}

func (p *postgresCtxManager) ByKey(ctx context.Context, key repository.CtxKey) repository.Transaction {
	tx, ok := ctx.Value(key).(repository.Transaction)
	if !ok {
		return nil
	}
	return tx
}

func (p *postgresCtxManager) Default(ctx context.Context) repository.Transaction {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return nil
	}
	return NewTransaction(tx)
}

func (p *postgresCtxManager) CtxKey() repository.CtxKey {
	return repository.CtxKey{}
}
