package postgresql

import (
	"context"
	"log/slog"

	"github.com/Ranik23/avito-tech-spring/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CtxKey struct{}

type txManager struct {
	pool   		*pgxpool.Pool
	logger 		*slog.Logger
    ctxManager  repository.CtxManager
}

func NewTxManager(pool *pgxpool.Pool, log *slog.Logger, ctxManager repository.CtxManager) repository.TxManager {
	return &txManager{
		pool: pool,
		logger: log,
		ctxManager: ctxManager,
	}
}

func (p *txManager) Do(ctx context.Context, fn func(context.Context) error) error {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return err
	}

	newCtx := context.WithValue(ctx, p.ctxManager.CtxKey(), tx)

	if err := fn(newCtx); err != nil {
		rollBackErr := tx.Rollback(ctx)
		for rollBackErr != nil {
			rollBackErr = tx.Rollback(ctx)
		}
		return err
	}

	return tx.Commit(ctx)
}