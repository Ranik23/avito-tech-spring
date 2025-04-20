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

func (s *TestSuite) TestRollBackCreatePVZ() {
	ctrl := gomock.NewController(s.T())
	exampleError := errors.New("do error")
	exampleCity := s.cities[0]

	txManager := mock.NewMockTxManager(ctrl)

	authService := service.NewAuthService(s.userRepo, txManager, s.token, s.hasher, s.logger)
	pvzService := service.NewPVZService(s.pvzRepo, s.receptionRepo, s.cities, s.productRepo, txManager, s.logger)
	service := service.NewService(authService, pvzService)

	txManager.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, fn func(txCtx context.Context) error) error {
			fn(ctx)
			return exampleError
		},
	)

	_, err := service.CreatePVZ(context.Background(), exampleCity)
	s.Require().Error(exampleError)


	db, err := sql.Open("postgres", s.psqlContainer.GetDSN())
	s.Require().NoError(err)
	defer db.Close()

	var pvz domain.Pvz

	err = db.QueryRow(`
		SELECT id FROM pvz;
	`).Scan(&pvz.ID)

	s.Require().Equal(sql.ErrNoRows, err)


}