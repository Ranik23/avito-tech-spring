package postgresql

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/Ranik23/avito-tech-spring/internal/models/domain"
	"github.com/Ranik23/avito-tech-spring/internal/repository"
	"github.com/jackc/pgx/v5"
)


type postgresPvzRepository struct {
	ctxManager repository.CtxManager
}

func NewPostgresPvzRepository(manager repository.CtxManager) repository.PvzRepository {
	return &postgresPvzRepository{
		ctxManager: manager,
	}
}

func (p *postgresPvzRepository) GetListPVZ(ctx context.Context) ([]domain.Pvz, error) {
	tr := p.ctxManager.ByKey(ctx, p.ctxManager.CtxKey())
	if tr == nil {
		tr = p.ctxManager.Default(ctx)
	}

	query, args, err := squirrel.
		Select("*").
		From("pvz").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	exec := tr.(pgx.Tx)

	rows, err := exec.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.Pvz
	for rows.Next() {
		var pvz domain.Pvz
		err = rows.Scan(&pvz.ID, &pvz.RegistrationDate, &pvz.City)
		if err != nil {
			return nil, err
		}
		result = append(result, pvz)
	}

	return result, nil
}


func (p *postgresPvzRepository) GetPVZ(ctx context.Context, id string) (*domain.Pvz, error) {
	tr := p.ctxManager.ByKey(ctx, p.ctxManager.CtxKey())
	if tr == nil {
		tr = p.ctxManager.Default(ctx)
	}
	exec := tr.(pgx.Tx)

	query, args, err := squirrel.
		Select("id", "registration_date", "city").
		From("pvz").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	var pvz domain.Pvz
	err = exec.QueryRow(ctx, query, args...).Scan(&pvz.ID, &pvz.RegistrationDate, &pvz.City)
	if err != nil {
		return nil, err
	}

	return &pvz, nil
}

func (p *postgresPvzRepository) GetPVZS(ctx context.Context, offset int, limit int) ([]domain.Pvz, error) {
	tr := p.ctxManager.ByKey(ctx, p.ctxManager.CtxKey())
	if tr == nil {
		tr = p.ctxManager.Default(ctx)
	}
	exec := tr.(pgx.Tx)

	query, args, err := squirrel.
		Select("id", "registration_date", "city").
		From("pvz").
		OrderBy("id").
		Limit(uint64(limit)).
		Offset(uint64(offset)).
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

	var result []domain.Pvz
	for rows.Next() {
		var pvz domain.Pvz
		err = rows.Scan(&pvz.ID, &pvz.RegistrationDate, &pvz.City)
		if err != nil {
			return nil, err
		}
		result = append(result, pvz)
	}

	return result, nil
}

func (p *postgresPvzRepository) CreatePVZ(ctx context.Context, city string) (*domain.Pvz, error) {
	tr := p.ctxManager.ByKey(ctx, p.ctxManager.CtxKey())
	if tr == nil {
		tr = p.ctxManager.Default(ctx)
	}
	exec := tr.(pgx.Tx)

	query, args, err := squirrel.
		Insert("pvz").
		Columns("city").
		Values(city).
		Suffix("RETURNING id, city, registration_date").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return nil, err
	}

	var pvz domain.Pvz

	err = exec.QueryRow(ctx, query, args...).Scan(&pvz.ID, &pvz.City, &pvz.RegistrationDate)
	if err != nil {
		return nil, err
	}

	return &pvz, nil
}
