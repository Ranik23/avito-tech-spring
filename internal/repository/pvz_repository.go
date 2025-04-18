package repository

import (
	"context"

	"github.com/Ranik23/avito-tech-spring/internal/models/domain"
)

type PvzRepository interface {
	CreatePVZ(ctx context.Context, city string) (*domain.Pvz, error)
	GetPVZS(ctx context.Context, offset int, limit int) ([]domain.Pvz, error)
	GetPVZ(ctx context.Context, id string) (*domain.Pvz, error)
	GetListPVZ(ctx context.Context) ([]domain.Pvz, error)
}
