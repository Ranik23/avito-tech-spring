package manager

import (
	"context"
	pool "github.com/jackc/pgx/v5/pgxpool"
)



type CtxManager interface {
    Default(context.Context) Transaction
    ByKey(context.Context, CtxKey) Transaction
	CtxKey() CtxKey
}

 type ctxManager struct {
	pool *pool.Pool
}

func (p *ctxManager) ByKey(ctx context.Context, key CtxKey) Transaction {
	tx, ok := ctx.Value(key).(Transaction)
	if !ok {
		return nil
	}
	return tx
}


func (p *ctxManager) Default(ctx context.Context) Transaction {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return nil
	}
	return NewTransaction(tx)
}

func (p *ctxManager) CtxKey() CtxKey {
	return CtxKey{}
}

func NewCtxManager(pool *pool.Pool) CtxManager {
	return &ctxManager{
		pool: pool,
	}
}