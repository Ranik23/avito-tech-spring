package repository

import (
	"context"
	
	"github.com/Ranik23/avito-tech-spring/internal/models/domain"
)

type UserRepository interface {
	CreateUser(ctx context.Context, email string, hashedPassword string, role string) (userID string, err error)
	GetUser(ctx context.Context, email string) (*domain.User, error)
}
