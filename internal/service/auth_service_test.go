//go:build unit

package service

import (
	"context"
	"log/slog"
	"testing"

	hashermock "github.com/Ranik23/avito-tech-spring/internal/hasher/mock"
	"github.com/Ranik23/avito-tech-spring/internal/models/domain"
	managermock "github.com/Ranik23/avito-tech-spring/internal/repository/manager/mock"
	repomock "github.com/Ranik23/avito-tech-spring/internal/repository/mock"
	tokenmock "github.com/Ranik23/avito-tech-spring/internal/token/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

var fn func(ctx context.Context) error


func TestLogin_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := repomock.NewMockUserRepository(ctrl)
	mockTxManager := managermock.NewMockTxManager(ctrl)
	mockToken := tokenmock.NewMockToken(ctrl)
	mockHasher := hashermock.NewMockHasher(ctrl)

	email := "test@example.com"
	password := "password"
	expectedToken := "token123"
	user := &domain.User{
		ID:           "userID123",
		Email:        email,
		PasswordHash: "hashedPassword",
		Role:         "Client",
	}

	authService := NewAuthService(mockUserRepo, mockTxManager, mockToken, mockHasher, slog.Default())

	mockUserRepo.EXPECT().GetUser(gomock.Any(), email).Return(user, nil).Times(1)
	mockHasher.EXPECT().Equal(user.PasswordHash, password).Return(true).Times(1)
	mockToken.EXPECT().GenerateToken(user.ID, user.Role).Return(expectedToken, nil).Times(1)
	mockTxManager.EXPECT().Do(gomock.Any(), gomock.AssignableToTypeOf(fn)).DoAndReturn(
		func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		},
	).Times(1)

	token, err := authService.Login(context.Background(), email, password)

	assert.NoError(t, err)
	assert.Equal(t, expectedToken, token)
}

func TestLogin_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := repomock.NewMockUserRepository(ctrl)
	mockTxManager := managermock.NewMockTxManager(ctrl)
	mockToken := tokenmock.NewMockToken(ctrl)
	mockHasher := hashermock.NewMockHasher(ctrl)

	email := "test@example.com"
	password := "password"

	authService := NewAuthService(mockUserRepo, mockTxManager, mockToken, mockHasher, slog.Default())

	mockUserRepo.EXPECT().GetUser(gomock.Any(), email).Return(nil, nil).Times(1)


	mockTxManager.EXPECT().Do(gomock.Any(), gomock.AssignableToTypeOf(fn)).DoAndReturn(
		func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		},
	)

	token, err := authService.Login(context.Background(), email, password)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, ErrNotFound, err)
	assert.Empty(t, token)
}

func TestLogin_InvalidCredentials(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := repomock.NewMockUserRepository(ctrl)
	mockTxManager := managermock.NewMockTxManager(ctrl)
	mockToken := tokenmock.NewMockToken(ctrl)
	mockHasher := hashermock.NewMockHasher(ctrl)

	email := "test@example.com"
	password := "password"
	user := &domain.User{
		ID:           "userID123",
		Email:        email,
		PasswordHash: "hashedPassword",
		Role:         "Client",
	}

	authService := NewAuthService(mockUserRepo, mockTxManager, mockToken, mockHasher, slog.Default())

	mockUserRepo.EXPECT().GetUser(gomock.Any(), email).Return(user, nil).Times(1)
	mockHasher.EXPECT().Equal(user.PasswordHash, password).Return(false).Times(1)


	mockTxManager.EXPECT().Do(gomock.Any(), gomock.AssignableToTypeOf(fn)).DoAndReturn(
		func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		},
	)

	token, err := authService.Login(context.Background(), email, password)

	assert.Error(t, err)
	assert.Equal(t, ErrInvalidCredentials, err)
	assert.Empty(t, token)
}

func TestRegister_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := repomock.NewMockUserRepository(ctrl)
	mockTxManager := managermock.NewMockTxManager(ctrl)
	mockToken := tokenmock.NewMockToken(ctrl)
	mockHasher := hashermock.NewMockHasher(ctrl)

	email := "test@example.com"
	password := "password"
	role := "Client"
	userID := "userID123"
	hashedPassword := "hashedPassword"

	authService := NewAuthService(mockUserRepo, mockTxManager, mockToken, mockHasher, slog.Default())

	mockUserRepo.EXPECT().GetUser(gomock.Any(), email).Return(nil, nil).Times(1)
	mockHasher.EXPECT().Hash(password).Return(hashedPassword, nil).Times(1)
	mockUserRepo.EXPECT().CreateUser(gomock.Any(), email, hashedPassword, role).Return(userID, nil).Times(1)

	mockTxManager.EXPECT().Do(gomock.Any(), gomock.AssignableToTypeOf(fn)).DoAndReturn(
		func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		},
	)

	id, err := authService.Register(context.Background(), email, password, role)

	assert.NoError(t, err)
	assert.Equal(t, userID, id)
}

func TestRegister_UserAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := repomock.NewMockUserRepository(ctrl)
	mockTxManager := managermock.NewMockTxManager(ctrl)
	mockToken := tokenmock.NewMockToken(ctrl)
	mockHasher := hashermock.NewMockHasher(ctrl)

	email := "test@example.com"
	password := "password"
	role := "Client"

	authService := NewAuthService(mockUserRepo, mockTxManager, mockToken, mockHasher, slog.Default())

	mockUserRepo.EXPECT().GetUser(gomock.Any(), email).Return(&domain.User{}, nil).Times(1)

	mockTxManager.EXPECT().Do(gomock.Any(), gomock.AssignableToTypeOf(fn)).DoAndReturn(
		func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		},
	)

	id, err := authService.Register(context.Background(), email, password, role)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, ErrAlreadyExists, err)
	assert.Empty(t, id)
}

func TestDummyLogin_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := repomock.NewMockUserRepository(ctrl)
	mockTxManager := managermock.NewMockTxManager(ctrl)
	mockToken := tokenmock.NewMockToken(ctrl)
	mockHasher := hashermock.NewMockHasher(ctrl)

	role := "Client"
	expectedToken := "dummyToken"

	authService := NewAuthService(mockUserRepo, mockTxManager, mockToken, mockHasher, slog.Default())

	mockToken.EXPECT().GenerateToken("dummy", role).Return(expectedToken, nil).Times(1)

	token, err := authService.DummyLogin(context.Background(), role)

	assert.NoError(t, err)
	assert.Equal(t, expectedToken, token)
}

func TestDummyLogin_InvalidRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := repomock.NewMockUserRepository(ctrl)
	mockTxManager := managermock.NewMockTxManager(ctrl)
	mockToken := tokenmock.NewMockToken(ctrl)
	mockHasher := hashermock.NewMockHasher(ctrl)

	role := "InvalidRole"

	authService := NewAuthService(mockUserRepo, mockTxManager, mockToken, mockHasher, slog.Default())

	token, err := authService.DummyLogin(context.Background(), role)

	assert.Error(t, err)
	assert.Equal(t, ErrInvalidRole, err)
	assert.Empty(t, token)
}
