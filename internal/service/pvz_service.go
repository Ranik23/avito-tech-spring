package service

import (
	"context"
	"log/slog"
	"time"
)

type PVZService interface {
	CreatePVZ(ctx context.Context, city string) (pvzID string, err error)
	GetPVZInfo(ctx context.Context, start time.Time, end time.Time, page int, limit int) ([]PVZInfo, error)
	AddProductToReception(ctx context.Context, pvzID string, product_type string) (productID string, err error)
	DeleteLastProduct(ctx context.Context, pvzID string) error
	StartReception(ctx context.Context, pvzID string) (receptionID string, err error)
	CloseReception(ctx context.Context, pvzID string) (receptionID string, err error)
}

type pvzService struct {
	logger *slog.Logger
}

// AddProductToReception implements PVZService.
func (p *pvzService) AddProductToReception(ctx context.Context, pvzID string, product_type string) (productID string, err error) {
	panic("unimplemented")
}

// CloseReception implements PVZService.
func (p *pvzService) CloseReception(ctx context.Context, pvzID string) (receptionID string, err error) {
	panic("unimplemented")
}

// CreatePVZ implements PVZService.
func (p *pvzService) CreatePVZ(ctx context.Context, city string) (pvzID string, err error) {
	panic("unimplemented")
}

// DeleteLastProduct implements PVZService.
func (p *pvzService) DeleteLastProduct(ctx context.Context, pvzID string) error {
	panic("unimplemented")
}

// GetPVZInfo implements PVZService.
func (p *pvzService) GetPVZInfo(ctx context.Context, start time.Time, end time.Time, page int, limit int) ([]PVZInfo, error) {
	panic("unimplemented")
}

// StartReception implements PVZService.
func (p *pvzService) StartReception(ctx context.Context, pvzID string) (receptionID string, err error) {
	panic("unimplemented")
}

func NewPVZService(logger *slog.Logger) PVZService {
	return &pvzService{
		logger: logger,
	}
}
