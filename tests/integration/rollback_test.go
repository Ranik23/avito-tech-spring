//go:build integration

package integration

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Ranik23/avito-tech-spring/internal/models/domain"
	"github.com/Ranik23/avito-tech-spring/internal/repository/mock"
	"go.uber.org/mock/gomock"
)



var fn func(ctx context.Context) error

func (s *TestSuite) TestRollBack() {
	ctrl := gomock.NewController(s.T())
	txManager := mock.NewMockTxManager(ctrl)
	exampleError := errors.New("do error")
	exampleEmail := "email"
	examplePassword := "password"
	exampleRole := "employee"

	txManager.EXPECT().Do(gomock.Any(), gomock.AssignableToTypeOf(fn)).DoAndReturn(
		func(ctx context.Context, fn func(context.Context) error) error {
			fn(ctx)
			return exampleError
		},
	)
	_, err := s.service.Register(context.Background(), exampleEmail, examplePassword, exampleRole)
	s.Require().Error(exampleError)
	
	db, err := sql.Open("postgres", s.psqlContainer.GetDSN())
	s.Require().NoError(err)
	defer db.Close()


	var user domain.User

	err = db.QueryRow(`
		SELECT id FROM users;
	`).Scan(&user.ID)

	s.Require().Error(err)
}
