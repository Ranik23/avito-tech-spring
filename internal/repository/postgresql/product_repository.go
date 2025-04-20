package postgresql

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Masterminds/squirrel"
	"github.com/Ranik23/avito-tech-spring/internal/models/domain"
	"github.com/Ranik23/avito-tech-spring/internal/repository"
	"github.com/jackc/pgx/v5"
)

type postgresProductRepository struct {
	ctxManager repository.CtxManager
	logger     *slog.Logger
}

func NewPostgresProductRepository(manager repository.CtxManager, logger *slog.Logger) repository.ProductRepository {
	return &postgresProductRepository{
		ctxManager: manager,
		logger:     logger,
	}
}

func (p *postgresProductRepository) CreateProduct(ctx context.Context, productType string, receptionID string) (*domain.Product, error) {
	tr := p.ctxManager.ByKey(ctx, p.ctxManager.CtxKey())
	if tr == nil {
		tr = p.ctxManager.Default(ctx)
	}

	exec := tr.Begin().(pgx.Tx)

	query, args, err := squirrel.
		Insert("product").
		Columns("type", "reception_id").
		Values(productType, receptionID).
		Suffix("RETURNING id, type, reception_id, date_time").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		p.logger.Error("Failed to build SQL query for CreateProduct",
			slog.String("product_type", productType),
			slog.String("reception_id", receptionID),
			slog.String("error", err.Error()))
		return nil, err
	}

	var product domain.Product
	err = exec.QueryRow(ctx, query, args...).Scan(&product.ID, &product.Type, &product.ReceptionID, &product.DateTime)
	if err != nil {
		p.logger.Error("Failed to execute SQL query for CreateProduct",
			slog.String("product_type", productType),
			slog.String("reception_id", receptionID),
			slog.String("error", err.Error()))
		return nil, err
	}

	p.logger.Info("Successfully created product",
		slog.String("id", product.ID),
		slog.String("type", product.Type),
		slog.String("reception_id", product.ReceptionID),
		slog.Time("date_time", product.DateTime))

	return &product, nil
}

func (p *postgresProductRepository) DeleteProduct(ctx context.Context, productID string) error {
	tr := p.ctxManager.ByKey(ctx, p.ctxManager.CtxKey())
	if tr == nil {
		tr = p.ctxManager.Default(ctx)
	}

	exec := tr.Begin().(pgx.Tx)

	query, args, err := squirrel.
		Delete("product").
		Where(squirrel.Eq{"id": productID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		p.logger.Error("Failed to build SQL query for DeleteProduct",
			slog.String("product_id", productID),
			slog.String("error", err.Error()))
		return err
	}

	_, err = exec.Exec(ctx, query, args...)
	if err != nil {
		p.logger.Error("Failed to execute SQL query for DeleteProduct",
			slog.String("product_id", productID),
			slog.String("error", err.Error()))
		return err
	}

	p.logger.Info("Successfully deleted product",
		slog.String("product_id", productID))

	return nil
}

func (p *postgresProductRepository) FindTheLastProduct(ctx context.Context, pvzID string) (*domain.Product, error) {
	tr := p.ctxManager.ByKey(ctx, p.ctxManager.CtxKey())
	if tr == nil {
		tr = p.ctxManager.Default(ctx)
	}
	exec := tr.Begin().(pgx.Tx)

	query, args, err := squirrel.
		Select("product.id", "product.date_time", "type", "reception_id").
		From("product").
		Join("reception on product.reception_id = reception.id").
		Where(squirrel.Eq{"reception.pvz_id": pvzID}).
		OrderBy("date_time DESC").
		Limit(1).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		p.logger.Error("Failed to build SQL query for FindTheLastProduct",
			slog.String("pvz_id", pvzID),
			slog.String("error", err.Error()))
		return nil, err
	}

	var prod domain.Product
	err = exec.QueryRow(ctx, query, args...).Scan(&prod.ID, &prod.DateTime, &prod.Type, &prod.ReceptionID)
	if err != nil {
		p.logger.Error("Failed to execute SQL query for FindTheLastProduct",
			slog.String("pvz_id", pvzID),
			slog.String("error", err.Error()))

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	p.logger.Info("Successfully retrieved last product",
		slog.String("id", prod.ID),
		slog.String("type", prod.Type),
		slog.String("reception_id", prod.ReceptionID),
		slog.Time("date_time", prod.DateTime))

	return &prod, nil
}

func (p *postgresProductRepository) GetProducts(ctx context.Context, receptionID string) ([]domain.Product, error) {
	tr := p.ctxManager.ByKey(ctx, p.ctxManager.CtxKey())
	if tr == nil {
		tr = p.ctxManager.Default(ctx)
	}
	exec := tr.Begin().(pgx.Tx)

	query, args, err := squirrel.
		Select("id", "date_time", "type", "reception_id").
		From("product").
		Where(squirrel.Eq{"reception_id": receptionID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		p.logger.Error("Failed to build SQL query for GetProducts",
			slog.String("reception_id", receptionID),
			slog.String("error", err.Error()))
		return nil, err
	}

	rows, err := exec.Query(ctx, query, args...)
	if err != nil {
		p.logger.Error("Failed to execute SQL query for GetProducts",
			slog.String("reception_id", receptionID),
			slog.String("error", err.Error()))
		return nil, err
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var prod domain.Product
		if err := rows.Scan(&prod.ID, &prod.DateTime, &prod.Type, &prod.ReceptionID); err != nil {
			p.logger.Error("Failed to scan row in GetProducts",
				slog.String("reception_id", receptionID),
				slog.String("error", err.Error()))
			return nil, err
		}
		products = append(products, prod)
	}

	p.logger.Info("Successfully retrieved list of products",
		slog.Int("count", len(products)),
		slog.String("reception_id", receptionID))

	return products, nil
}
