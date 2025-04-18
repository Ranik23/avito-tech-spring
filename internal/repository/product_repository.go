package repository

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/Ranik23/avito-tech-spring/internal/models/domain"
	"github.com/Ranik23/avito-tech-spring/internal/repository/manager"
	"github.com/jackc/pgx/v5"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, productType string, receptionID string) (*domain.Product, error)
	DeleteProduct(ctx context.Context, productID string) error
	FindTheLastProduct(ctx context.Context, pvzID string) (product *domain.Product, err error)
	GetProducts(ctx context.Context, receptionID string) ([]domain.Product, error)
}

type postgresProductRepository struct {
	ctxManager manager.CtxManager
}

func NewPostgresProductRepository(manager manager.CtxManager) ProductRepository {
	return &postgresProductRepository{
		ctxManager: manager,
	}
}


func (p *postgresProductRepository) CreateProduct(ctx context.Context, productType string, receptionID string) (*domain.Product, error) {
	tr := p.ctxManager.ByKey(ctx, p.ctxManager.CtxKey())
	if tr == nil {
		tr = p.ctxManager.Default(ctx)
	}
	exec := tr.(pgx.Tx)

	query, args, err := squirrel.
		Insert("product").
		Columns("type", "reception_id").
		Values(productType, receptionID).
		Suffix("RETURNING id, type, reception_id, date_time").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return nil, err
	}

	var product domain.Product

	err = exec.QueryRow(ctx, query, args...).Scan(&product.ID, &product.Type,
					&product.ReceptionID, &product.DateTime)
	
	return &product, err
}


func (p *postgresProductRepository) DeleteProduct(ctx context.Context, productID string) error {
	tr := p.ctxManager.ByKey(ctx, p.ctxManager.CtxKey())
	if tr == nil {
		tr = p.ctxManager.Default(ctx)
	}
	exec := tr.(pgx.Tx)

	query, args, err := squirrel.
		Delete("product").
		Where(squirrel.Eq{"id": productID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = exec.Exec(ctx, query, args...)
	return err
}


func (p *postgresProductRepository) FindTheLastProduct(ctx context.Context, pvzID string) (*domain.Product, error) {
	tr := p.ctxManager.ByKey(ctx, p.ctxManager.CtxKey())
	if tr == nil {
		tr = p.ctxManager.Default(ctx)
	}
	exec := tr.(pgx.Tx)

	query, args, err := squirrel.
		Select("id", "date_time", "type", "reception_id").
		From("product").
		Join("reception on product.reception_id = reception_id").
		Where(squirrel.Eq{"reception.pvz_id": pvzID}).
		OrderBy("date_time DESC").
		Limit(1).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
		
	if err != nil {
		return nil, err
	}

	var prod domain.Product
	
	err = exec.QueryRow(ctx, query, args...).Scan(&prod.ID, &prod.DateTime, &prod.Type, &prod.ReceptionID)
	if err != nil {
		return nil, err
	}

	return &prod, nil
}


func (p *postgresProductRepository) GetProducts(ctx context.Context, receptionID string) ([]domain.Product, error) {
	tr := p.ctxManager.ByKey(ctx, p.ctxManager.CtxKey())
	if tr == nil {
		tr = p.ctxManager.Default(ctx)
	}
	exec := tr.(pgx.Tx)

	query, args, err := squirrel.
		Select("id", "date_time", "type", "reception_id").
		From("product").
		Where(squirrel.Eq{"reception_id": receptionID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := exec.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var prod domain.Product
		if err := rows.Scan(&prod.ID, &prod.DateTime, &prod.Type, &prod.ReceptionID); err != nil {
			return nil, err
		}
		products = append(products, prod)
	}

	return products, nil
}
