package repository

import (
	"context"
	"strconv"

	"github.com/Masterminds/squirrel"
	"github.com/Ranik23/avito-tech-spring/internal/models/domain"
	"github.com/Ranik23/avito-tech-spring/internal/repository/manager"
	"github.com/jackc/pgx/v5"
)

type PvzRepository interface {
	CreatePVZ(ctx context.Context, city string) (pvzID string, err error)
	GetPVZS(ctx context.Context, offset int, limit int) ([]domain.Pvz, error)
	GetPVZ(ctx context.Context, id string) (*domain.Pvz, error)
}

type postgresPvzRepository struct {
	сtxManager manager.CtxManager
}

func NewPostgresPvzRepository() PvzRepository {
	return &postgresPvzRepository{}
}

func (p *postgresPvzRepository) GetPVZ(ctx context.Context, id string) (*domain.Pvz, error) {
	tr := p.сtxManager.ByKey(ctx, p.сtxManager.CtxKey())
	if tr == nil {
		tr = p.сtxManager.Default(ctx)
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
	tr := p.сtxManager.ByKey(ctx, p.сtxManager.CtxKey())
	if tr == nil {
		tr = p.сtxManager.Default(ctx)
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

func (p *postgresPvzRepository) CreatePVZ(ctx context.Context, city string) (pvzID string, err error) {
	tr := p.сtxManager.ByKey(ctx, p.сtxManager.CtxKey())
	if tr == nil {
		tr = p.сtxManager.Default(ctx)
	}
	exec := tr.(pgx.Tx)

	query, args, err := squirrel.
		Insert("pvz").
		Columns("city").
		Values(city).
		Suffix("RETURNING id").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return "", err
	}

	var id int
	err = exec.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(id), nil
}

