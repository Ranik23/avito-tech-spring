package postgresql

import (
	"context"
	"errors"
	"log/slog"
	"strconv"

	"github.com/Masterminds/squirrel"
	"github.com/Ranik23/avito-tech-spring/internal/models/domain"
	"github.com/Ranik23/avito-tech-spring/internal/repository"
	"github.com/jackc/pgx/v5"
)

type postgresUserRepository struct {
	ctxManager repository.CtxManager
	logger     *slog.Logger
}

func NewPostgresUserRepository(ctxManager repository.CtxManager, logger *slog.Logger) repository.UserRepository {
	return &postgresUserRepository{
		ctxManager: ctxManager,
		logger:     logger,
	}
}

func (p *postgresUserRepository) CreateUser(ctx context.Context, email string, hashedPassword string, role string) (userID string, err error) {
	tr := p.ctxManager.ByKey(ctx, p.ctxManager.CtxKey())
	if tr == nil {
		tr = p.ctxManager.Default(ctx)
	}

	exec := tr.Begin().(pgx.Tx)

	query, args, err := squirrel.Insert("users").
		Columns("email", "hashed_password", "role").
		Values(email, hashedPassword, role).
		PlaceholderFormat(squirrel.Dollar).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		p.logger.Error("Failed to build SQL query for creating user",
			slog.String("email", email),
			slog.String("role", role),
			slog.String("error", err.Error()))
		return "", err
	}

	var id int
	err = exec.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		p.logger.Error("Failed to execute SQL query for creating user",
			slog.String("email", email),
			slog.String("role", role),
			slog.String("error", err.Error()))
		return "", err
	}

	p.logger.Info("Successfully created user",
		slog.String("email", email),
		slog.String("role", role),
		slog.Int("userID", id))

	return strconv.Itoa(id), nil
}

func (p *postgresUserRepository) GetUser(ctx context.Context, email string) (*domain.User, error) {
	tr := p.ctxManager.ByKey(ctx, p.ctxManager.CtxKey())
	if tr == nil {
		tr = p.ctxManager.Default(ctx)
	}

	exec := tr.Begin().(pgx.Tx)

	query, args, err := squirrel.Select("id, email, hashed_password, role, created_at").
		From("users").
		Where(squirrel.Eq{"email": email}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		p.logger.Error("Failed to build SQL query for getting user",
			slog.String("email", email),
			slog.String("error", err.Error()))
		return nil, err
	}

	var user domain.User
	err = exec.QueryRow(ctx, query, args...).Scan(&user.ID, &user.Email,&user.PasswordHash,  &user.Role, &user.CreatedAt)
	if err != nil {
		p.logger.Error("Failed to execute SQL query for getting user",
			slog.String("email", email),
			slog.String("error", err.Error()))
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	p.logger.Info("Successfully retrieved user",
		slog.String("email", email),
		slog.String("userID", user.ID))

	return &user, nil
}
