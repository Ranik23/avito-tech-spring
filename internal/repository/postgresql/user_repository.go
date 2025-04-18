package postgresql

import (
	"context"
	"strconv"

	"github.com/Masterminds/squirrel"
	"github.com/Ranik23/avito-tech-spring/internal/models/domain"
	"github.com/Ranik23/avito-tech-spring/internal/repository"
	"github.com/jackc/pgx/v5"
)


type postgresUserRepository struct {
	ctxManager 	repository.CtxManager
}

func NewPostgresUserRepository(ctxManager repository.CtxManager) repository.UserRepository {
	return &postgresUserRepository{
		ctxManager: ctxManager,
	}
}


func (p *postgresUserRepository) CreateUser(ctx context.Context, email string, hashedPassword string, role string) (userID string, err error) {
	tr := p.ctxManager.ByKey(ctx, p.ctxManager.CtxKey())
	if tr == nil {
		tr = p.ctxManager.Default(ctx)
	}

	exec := tr.(pgx.Tx)

	query, args, err := squirrel.Insert("users").
 			Columns("email", "hashed_password", "role").
    		Values(email, hashedPassword, role).
    		PlaceholderFormat(squirrel.Dollar).
			Suffix("RETURNING id").
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


func (p *postgresUserRepository) GetUser(ctx context.Context, email string) (*domain.User, error) {
	tr := p.ctxManager.ByKey(ctx, p.ctxManager.CtxKey())
	if tr == nil {
		tr = p.ctxManager.Default(ctx)
	}

	exec := tr.(pgx.Tx)


	query, args, err := squirrel.Select("id, email, role, created_at").
			From("users").
			Where(squirrel.Eq{"email": email}).
			PlaceholderFormat(squirrel.Dollar).
			ToSql()
	if err != nil {
		return nil, err
	}

	var user domain.User

	err = exec.QueryRow(ctx, query, args...).Scan(&user.ID, &user.Email, &user.Role, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

