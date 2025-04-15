package repository

import "context"




type ProductRepository interface {
	CreateProduct(ctx context.Context)
}