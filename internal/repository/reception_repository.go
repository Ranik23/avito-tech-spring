package repository

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/Ranik23/avito-tech-spring/internal/models/domain"
	"github.com/Ranik23/avito-tech-spring/internal/repository/manager"
	"github.com/jackc/pgx/v5"
)

type ReceptionRepository interface {
	FindOpen(ctx context.Context, pvzID string) (*domain.Reception, error)
	CreateReception(ctx context.Context, pvzID string) (*domain.Reception, error)
	UpdateReceptionStatus(ctx context.Context, receptionID string, newStatus string) error
	GetReceptionsFiltered(ctx context.Context, pvzID string, startTime time.Time, endTime time.Time) ([]*domain.Reception, error)
}

type postgresReceptionRepository struct {
	ctxManager manager.CtxManager
}


func NewPostgresReceptionRepository(manager manager.CtxManager) ReceptionRepository {
	return &postgresReceptionRepository{
		ctxManager: manager,
	}
}

// CreateReception implements ReceptionRepository.
func (p *postgresReceptionRepository) CreateReception(ctx context.Context, pvzID string) (*domain.Reception, error) {
	tr := p.ctxManager.ByKey(ctx, p.ctxManager.CtxKey())
	if tr == nil {
		tr = p.ctxManager.Default(ctx)
	}
	exec := tr.(pgx.Tx)

	query, args, err := squirrel.
		Insert("reception").
		Columns("pvz_id", "status").
		Values(pvzID, "open").
		Suffix("RETURNING id, pvz_id, status, date_time").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	var reception domain.Reception

	err = exec.QueryRow(ctx, query, args...).Scan(&reception.ID, &reception.PvzID, &reception.Status, &reception.DateTime)
	if err != nil {
		return nil, err
	}

	return &reception, nil
}

func (p *postgresReceptionRepository) FindOpen(ctx context.Context, pvzID string) (*domain.Reception, error) {
	tr := p.ctxManager.ByKey(ctx, p.ctxManager.CtxKey())
	if tr == nil {
		tr = p.ctxManager.Default(ctx)
	}
	exec := tr.(pgx.Tx)

	query, args, err := squirrel.
		Select("id", "date_time", "pvz_id", "status").
		From("reception").
		Where(squirrel.Eq{"pvz_id": pvzID, "status": "open"}).
		OrderBy("date_time DESC").
		Limit(1).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	var r domain.Reception
	err = exec.QueryRow(ctx, query, args...).Scan(&r.ID, &r.DateTime, &r.PvzID, &r.Status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &r, nil
}

func (p *postgresReceptionRepository) GetReceptionsFiltered(ctx context.Context, pvzID string, startTime time.Time, endTime time.Time) ([]*domain.Reception, error) {
	tr := p.ctxManager.ByKey(ctx, p.ctxManager.CtxKey())
	if tr == nil {
		tr = p.ctxManager.Default(ctx)
	}
	exec := tr.(pgx.Tx)

	query, args, err := squirrel.
		Select("id", "date_time", "pvz_id", "status").
		From("reception").
		Where(squirrel.Eq{"pvz_id": pvzID}).
		Where(squirrel.And{
			squirrel.GtOrEq{"date_time": startTime},
			squirrel.LtOrEq{"date_time": endTime},
		}).
		OrderBy("date_time DESC").
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

	var result []*domain.Reception
	for rows.Next() {
		var r domain.Reception
		err = rows.Scan(&r.ID, &r.DateTime, &r.PvzID, &r.Status)
		if err != nil {
			return nil, err
		}
		result = append(result, &r)
	}
	return result, nil
}


func (p *postgresReceptionRepository) UpdateReceptionStatus(ctx context.Context, receptionID string, newStatus string) error {
	tr := p.ctxManager.ByKey(ctx, p.ctxManager.CtxKey())
	if tr == nil {
		tr = p.ctxManager.Default(ctx)
	}
	exec := tr.(pgx.Tx)

	id, err := strconv.Atoi(receptionID)
	if err != nil {
		return err
	}

	query, args, err := squirrel.
		Update("reception").
		Set("status", newStatus).
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = exec.Exec(ctx, query, args...)
	return err
}
