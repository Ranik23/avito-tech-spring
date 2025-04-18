package repository

import (
	"context"
	"time"

	"github.com/Ranik23/avito-tech-spring/internal/models/domain"
)

type ReceptionRepository interface {
	FindOpen(ctx context.Context, pvzID string) (*domain.Reception, error)
	CreateReception(ctx context.Context, pvzID string) (*domain.Reception, error)
	UpdateReceptionStatus(ctx context.Context, receptionID string, newStatus string) error
	GetReceptionsFiltered(ctx context.Context, pvzID string, startTime time.Time, endTime time.Time) ([]*domain.Reception, error)
}
