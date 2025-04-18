package service

import (
	"context"
	"log/slog"
	"slices"

	"github.com/Ranik23/avito-tech-spring/internal/hasher"
	"github.com/Ranik23/avito-tech-spring/internal/repository/manager"
	"github.com/Ranik23/avito-tech-spring/internal/repository"
	"github.com/Ranik23/avito-tech-spring/internal/token"
	"github.com/Ranik23/avito-tech-spring/pkg/errs"
)

type AuthService interface {
	DummyLogin(ctx context.Context, role string) (token string, err error)
	Register(ctx context.Context, email string, password string, role string) (userID string, err error)
	Login(ctx context.Context, email string, password string) (token string, err error)
}

type authService struct {
	userRepo 	repository.UserRepository
	txManager 	manager.TxManager
	token	    token.Token  
	hasher		hasher.Hasher
	logger 		*slog.Logger
	acceptableRoles []string
}

func NewAuthService(
	userRepo 	repository.UserRepository,
	txManager 	manager.TxManager,
	token 		token.Token,
	hasher		hasher.Hasher,
	logger 		*slog.Logger,
) AuthService {
	return &authService{
		userRepo:  userRepo,
		txManager: txManager,
		token:     token,
		hasher:    hasher,
		logger:    logger,
		acceptableRoles: []string{"Client", "Moderator"},
	}
}


func (a *authService) Login(ctx context.Context, email string, password string) (token string, err error) {
	err = a.txManager.Do(ctx, func(txCtx context.Context) error {
		user, err := a.userRepo.GetUser(txCtx, email)
		if err != nil {
			return errs.Wrap("failed to get the user", err)
		}

		if user == nil {
			return ErrUserNotFound
		}

		if !a.hasher.Equal(user.PasswordHash, password) {
			return ErrInvalidCredentials
		}

		token, err = a.token.GenerateToken(user.ID, user.Role)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		a.logger.Error("Login failed",
			slog.String("email", email),
			slog.String("error", err.Error()),
		)
		return "", err
	}

	a.logger.Info("Login successful",
		slog.String("email", email),
	)
	return token, nil
}


func (a *authService) Register(ctx context.Context, email string, password string, role string) (userID string, err error) {
	err = a.txManager.Do(ctx, func(txCtx context.Context) error {
		user, err := a.userRepo.GetUser(txCtx, email)
		if err != nil {
			return errs.Wrap("failed to get the user", err)
		}

		if user != nil {
			return ErrAlreadyExists
		}

		hashedPassword, err := a.hasher.Hash(password)
		if err != nil {
			return errs.Wrap("failed to hash the password", err)
		}

		userID, err = a.userRepo.CreateUser(txCtx, email, hashedPassword, role)
		if err != nil {
			return errs.Wrap("failed to create the user", err)
		}

		return nil
	})

	if err != nil {
		a.logger.Error("Registration failed",
			slog.String("email", email),
			slog.String("role", role),
			slog.String("error", err.Error()),
		)
		return "", err
	}

	a.logger.Info("Registration successful",
		slog.String("email", email),
		slog.String("role", role),
		slog.String("userID", userID),
	)
	return userID, nil
}


func (a *authService) DummyLogin(ctx context.Context, role string) (token string, err error) {
	if !slices.Contains(a.acceptableRoles, role) {
		a.logger.Error("Dummy login failed: invalid role",
			slog.String("role", role),
		)
		return "", ErrInvalidRole
	}

	token, err = a.token.GenerateToken("dummy", role)
	if err != nil {
		a.logger.Error("Dummy login failed: token generation error",
			slog.String("role", role),
			slog.String("error", err.Error()),
		)
		return "", err
	}
	return token, err
}
