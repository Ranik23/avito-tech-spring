package service

import (
	"context"

	"github.com/Ranik23/avito-tech-spring/internal/repository"
)

type AuthService interface {
	DummyLogin(ctx context.Context, role string)
	Register(ctx context.Context, email string, password string, role string)
	Login(ctx context.Context, email string, password string)
}

type authService struct {
	userRepo repository.UserRepository
}

// DummyLogin implements AuthService.
func (a *authService) DummyLogin(ctx context.Context, role string) {
	panic("unimplemented")
}

// Login implements AuthService.
func (a *authService) Login(ctx context.Context, email string, password string) {
	panic("unimplemented")
}

// Register implements AuthService.
func (a *authService) Register(ctx context.Context, email string, password string, role string) {
	panic("unimplemented")
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{
		userRepo: userRepo,
	}
}
