package postgresql

import (
	"context"
	"errors"
	"log/slog"
	"strconv"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/Ranik23/avito-tech-spring/internal/models/domain"
	"github.com/Ranik23/avito-tech-spring/internal/repository"
	"github.com/jackc/pgx/v5"
)

type postgresReceptionRepository struct {
	ctxManager repository.CtxManager
	logger     *slog.Logger
}

func NewPostgresReceptionRepository(manager repository.CtxManager, logger *slog.Logger) repository.ReceptionRepository {
	return &postgresReceptionRepository{
		ctxManager: manager,
		logger:     logger,
	}
}

// CreateReception implements ReceptionRepository.
func (p *postgresReceptionRepository) CreateReception(ctx context.Context, pvzID string) (*domain.Reception, error) {
	tr := p.ctxManager.ByKey(ctx, p.ctxManager.CtxKey())
	if tr == nil {
		tr = p.ctxManager.Default(ctx)
	}
	exec := tr.Begin().(pgx.Tx)

	query, args, err := squirrel.
		Insert("reception").
		Columns("pvz_id", "status").
		Values(pvzID, "open").
		Suffix("RETURNING id, pvz_id, status, date_time").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		p.logger.Error("Failed to build SQL query for creating reception",
			slog.String("pvzID", pvzID),
			slog.String("error", err.Error()))
		return nil, err
	}

	var reception domain.Reception
	err = exec.QueryRow(ctx, query, args...).Scan(&reception.ID, &reception.PvzID, &reception.Status, &reception.DateTime)
	if err != nil {
		p.logger.Error("Failed to execute SQL query for creating reception",
			slog.String("pvzID", pvzID),
			slog.String("error", err.Error()))
		return nil, err
	}

	p.logger.Info("Successfully created reception",
		slog.String("pvzID", pvzID),
		slog.String("receptionID", reception.ID),
		slog.String("status", reception.Status),
		slog.Time("date_time", reception.DateTime))

	return &reception, nil
}

// FindOpen implements ReceptionRepository.
func (p *postgresReceptionRepository) FindOpen(ctx context.Context, pvzID string) (*domain.Reception, error) {
	tr := p.ctxManager.ByKey(ctx, p.ctxManager.CtxKey())
	if tr == nil {
		tr = p.ctxManager.Default(ctx)
	}
	exec := tr.Begin().(pgx.Tx)

	query, args, err := squirrel.
		Select("id", "date_time", "pvz_id", "status").
		From("reception").
		Where(squirrel.Eq{"pvz_id": pvzID, "status": "open"}).
		OrderBy("date_time DESC").
		Limit(1).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		p.logger.Error("Failed to build SQL query for finding open reception",
			slog.String("pvzID", pvzID),
			slog.String("error", err.Error()))
		return nil, err
	}

	var r domain.Reception
	err = exec.QueryRow(ctx, query, args...).Scan(&r.ID, &r.DateTime, &r.PvzID, &r.Status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			p.logger.Info("No open reception found",
				slog.String("pvzID", pvzID))
			return nil, nil
		}
		p.logger.Error("Failed to execute SQL query for finding open reception",
			slog.String("pvzID", pvzID),
			slog.String("error", err.Error()))
		return nil, err
	}

	p.logger.Info("Successfully found open reception",
		slog.String("pvzID", pvzID),
		slog.String("receptionID", r.ID),
		slog.String("status", r.Status),
		slog.Time("date_time", r.DateTime))

	return &r, nil
}

// GetReceptionsFiltered implements ReceptionRepository.
func (p *postgresReceptionRepository) GetReceptionsFiltered(ctx context.Context, pvzID string, startTime time.Time, endTime time.Time) ([]*domain.Reception, error) {
	tr := p.ctxManager.ByKey(ctx, p.ctxManager.CtxKey())
	if tr == nil {
		tr = p.ctxManager.Default(ctx)
	}
	exec := tr.Begin().(pgx.Tx)

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
		p.logger.Error("Failed to build SQL query for getting filtered receptions",
			slog.String("pvzID", pvzID),
			slog.Time("startTime", startTime),
			slog.Time("endTime", endTime),
			slog.String("error", err.Error()))
		return nil, err
	}

	rows, err := exec.Query(ctx, query, args...)
	if err != nil {
		p.logger.Error("Failed to execute SQL query for getting filtered receptions",
			slog.String("pvzID", pvzID),
			slog.Time("startTime", startTime),
			slog.Time("endTime", endTime),
			slog.String("error", err.Error()))
		return nil, err
	}
	defer rows.Close()

	var result []*domain.Reception
	for rows.Next() {
		var r domain.Reception
		err = rows.Scan(&r.ID, &r.DateTime, &r.PvzID, &r.Status)
		if err != nil {
			p.logger.Error("Failed to scan reception data",
				slog.String("pvzID", pvzID),
				slog.String("error", err.Error()))
			return nil, err
		}
		result = append(result, &r)
	}

	p.logger.Info("Successfully retrieved filtered receptions",
		slog.String("pvzID", pvzID),
		slog.Int("receptionsCount", len(result)))

	return result, nil
}

// UpdateReceptionStatus implements ReceptionRepository.
func (p *postgresReceptionRepository) UpdateReceptionStatus(ctx context.Context, receptionID string, newStatus string) error {
	tr := p.ctxManager.ByKey(ctx, p.ctxManager.CtxKey())
	if tr == nil {
		tr = p.ctxManager.Default(ctx)
	}
	exec := tr.Begin().(pgx.Tx)

	id, err := strconv.Atoi(receptionID)
	if err != nil {
		p.logger.Error("Failed to convert receptionID to int",
			slog.String("receptionID", receptionID),
			slog.String("error", err.Error()))
		return err
	}

	query, args, err := squirrel.
		Update("reception").
		Set("status", newStatus).
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		p.logger.Error("Failed to build SQL query for updating reception status",
			slog.String("receptionID", receptionID),
			slog.String("newStatus", newStatus),
			slog.String("error", err.Error()))
		return err
	}

	_, err = exec.Exec(ctx, query, args...)
	if err != nil {
		p.logger.Error("Failed to execute SQL query for updating reception status",
			slog.String("receptionID", receptionID),
			slog.String("newStatus", newStatus),
			slog.String("error", err.Error()))
		return err
	}

	p.logger.Info("Successfully updated reception status",
		slog.String("receptionID", receptionID),
		slog.String("newStatus", newStatus))

	return nil
}
