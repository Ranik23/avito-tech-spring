package postgresql

import (
	"context"
	"log/slog"

	"github.com/Masterminds/squirrel"
	"github.com/Ranik23/avito-tech-spring/internal/models/domain"
	"github.com/Ranik23/avito-tech-spring/internal/repository"
	"github.com/jackc/pgx/v5"
)

type postgresPvzRepository struct {
	ctxManager repository.CtxManager
	logger     *slog.Logger
}

func NewPostgresPvzRepository(manager repository.CtxManager, logger *slog.Logger) repository.PvzRepository {
	return &postgresPvzRepository{
		ctxManager: manager,
		logger:     logger,
	}
}

func (p *postgresPvzRepository) GetListOfPVZS(ctx context.Context) ([]domain.Pvz, error) {
	tr := p.ctxManager.ByKey(ctx, p.ctxManager.CtxKey())
	if tr == nil {
		tr = p.ctxManager.Default(ctx)
	}

	exec := tr.Transaction().(pgx.Tx)

	query, args, err := squirrel.
		Select("*").
		From("pvz").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		p.logger.Error("Failed to build SQL query for GetListOfPVZS",
			slog.String("error", err.Error()))
		return nil, err
	}

	rows, err := exec.Query(ctx, query, args...)
	if err != nil {
		p.logger.Error("Failed to execute SQL query for GetListOfPVZS",
			slog.String("error", err.Error()))
		return nil, err
	}
	defer rows.Close()

	var result []domain.Pvz
	for rows.Next() {
		var pvz domain.Pvz
		err = rows.Scan(&pvz.ID, &pvz.RegistrationDate, &pvz.City)
		if err != nil {
			p.logger.Error("Failed to scan row in GetListOfPVZS",
				slog.String("error", err.Error()))
			return nil, err
		}
		result = append(result, pvz)
	}

	p.logger.Info("Successfully retrieved list of PVZ",
		slog.Int("count", len(result)))

	return result, nil
}

func (p *postgresPvzRepository) GetPVZ(ctx context.Context, id string) (*domain.Pvz, error) {
	tr := p.ctxManager.ByKey(ctx, p.ctxManager.CtxKey())
	if tr == nil {
		tr = p.ctxManager.Default(ctx)
	}
	exec := tr.Transaction().(pgx.Tx)

	query, args, err := squirrel.
		Select("id", "registration_date", "city").
		From("pvz").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		p.logger.Error("Failed to build SQL query for GetPVZ",
			slog.String("id", id),
			slog.String("error", err.Error()))
		return nil, err
	}

	var pvz domain.Pvz
	err = exec.QueryRow(ctx, query, args...).Scan(&pvz.ID, &pvz.RegistrationDate, &pvz.City)
	if err != nil {
		p.logger.Error("Failed to execute SQL query for GetPVZ",
			slog.String("id", id),
			slog.String("error", err.Error()))
		return nil, err
	}

	p.logger.Info("Successfully retrieved PVZ",
		slog.String("id", pvz.ID),
		slog.String("city", pvz.City))

	return &pvz, nil
}

func (p *postgresPvzRepository) GetPVZS(ctx context.Context, offset int, limit int) ([]domain.Pvz, error) {
	tr := p.ctxManager.ByKey(ctx, p.ctxManager.CtxKey())
	if tr == nil {
		tr = p.ctxManager.Default(ctx)
	}
	exec := tr.Transaction().(pgx.Tx)

	query, args, err := squirrel.
		Select("id", "registration_date", "city").
		From("pvz").
		OrderBy("id").
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		p.logger.Error("Failed to build SQL query for GetPVZS",
			slog.Int("offset", offset),
			slog.Int("limit", limit),
			slog.String("error", err.Error()))
		return nil, err
	}

	rows, err := exec.Query(ctx, query, args...)
	if err != nil {
		p.logger.Error("Failed to execute SQL query for GetPVZS",
			slog.Int("offset", offset),
			slog.Int("limit", limit),
			slog.String("error", err.Error()))
		return nil, err
	}
	defer rows.Close()

	var result []domain.Pvz
	for rows.Next() {
		var pvz domain.Pvz
		err = rows.Scan(&pvz.ID, &pvz.RegistrationDate, &pvz.City)
		if err != nil {
			p.logger.Error("Failed to scan row in GetPVZS",
				slog.String("error", err.Error()))
			return nil, err
		}
		result = append(result, pvz)
	}

	p.logger.Info("Successfully retrieved list of PVZ",
		slog.Int("count", len(result)),
		slog.Int("offset", offset),
		slog.Int("limit", limit))

	return result, nil
}

func (p *postgresPvzRepository) CreatePVZ(ctx context.Context, city string) (*domain.Pvz, error) {
	tr := p.ctxManager.ByKey(ctx, p.ctxManager.CtxKey())
	if tr == nil {
		tr = p.ctxManager.Default(ctx)
	}
	exec := tr.Transaction().(pgx.Tx)

	query, args, err := squirrel.
		Insert("pvz").
		Columns("city").
		Values(city).
		Suffix("RETURNING id, city, registration_date").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		p.logger.Error("Failed to build SQL query for CreatePVZ",
			slog.String("city", city),
			slog.String("error", err.Error()))
		return nil, err
	}

	var pvz domain.Pvz

	err = exec.QueryRow(ctx, query, args...).Scan(&pvz.ID, &pvz.City, &pvz.RegistrationDate)
	if err != nil {
		p.logger.Error("Failed to execute SQL query for CreatePVZ",
			slog.String("city", city),
			slog.String("error", err.Error()))
		return nil, err
	}

	p.logger.Info("Successfully created new PVZ",
		slog.String("id", pvz.ID),
		slog.String("city", pvz.City),
		slog.Time("registration_date", pvz.RegistrationDate))

	return &pvz, nil
}
