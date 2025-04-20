//go:build integration

package integration

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Ranik23/avito-tech-spring/internal/models/domain"
	"github.com/Ranik23/avito-tech-spring/internal/repository/mock"
	"github.com/Ranik23/avito-tech-spring/internal/service"
	"go.uber.org/mock/gomock"
)



var fn func(ctx context.Context) error

func (s *TestSuite) TestRollBackRegister() {
	ctrl := gomock.NewController(s.T())

	exampleError := errors.New("do error")
	exampleEmail := "email"
	examplePassword := "password"
	exampleRole := "employee"

	txManager := mock.NewMockTxManager(ctrl)

	authService := service.NewAuthService(s.userRepo, txManager, s.token, s.hasher, s.logger)
	pvzService := service.NewPVZService(s.pvzRepo, s.receptionRepo, s.cities, s.productRepo, txManager, s.logger)
	service := service.NewService(authService, pvzService)

	txManager.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, fn func(context.Context) error) error {
			fn(ctx)
			return exampleError
		},
	)

	_, err := service.Register(context.Background(), exampleEmail, examplePassword, exampleRole)
	s.Require().Error(exampleError)
	
	db, err := sql.Open("postgres", s.psqlContainer.GetDSN())
	s.Require().NoError(err)
	defer db.Close()


	var user domain.User

	err = db.QueryRow(`
		SELECT id FROM users;
	`).Scan(&user.ID)

	s.Require().Equal(sql.ErrNoRows, err)
}

