package manager

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CtxKey struct{}

type TxManager interface {
	Do(context.Context, func(context.Context) error) error
}

type txManager struct {
	pool   		*pgxpool.Pool
	logger 		*slog.Logger
    ctxManager   CtxManager
}

func NewTxManager(pool *pgxpool.Pool, log *slog.Logger) TxManager {
	return &txManager{
		pool: pool,
		logger: log,
	}
}

func (p *txManager) Do(ctx context.Context, fn func(context.Context) error) error {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return err
	}

	newCtx := context.WithValue(ctx, p.ctxManager.CtxKey(), tx)

	if err := fn(newCtx); err != nil {
		tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}