package repository

import (
	"context"

	"github.com/Ranik23/avito-tech-spring/internal/models/domain"

)

type ProductRepository interface {
	CreateProduct(ctx context.Context, productType string, receptionID string) (*domain.Product, error)
	DeleteProduct(ctx context.Context, productID string) error
	FindTheLastProduct(ctx context.Context, pvzID string) (product *domain.Product, err error)
	GetProducts(ctx context.Context, receptionID string) ([]domain.Product, error)
}
