//go:build integration

package integration

import (
	"context"
	"database/sql"

	"github.com/Ranik23/avito-tech-spring/internal/models/domain"
)


func (s *TestSuite) TestRegisterSuccess() {
	exampleEmail := "lol"
	exampleRole := "lol"
	examplePassword := "employee"

	userID, err := s.service.Register(context.Background(), exampleEmail, examplePassword, exampleRole)
	s.Require().NoError(err)
	

	db, err := sql.Open("postgres", s.psqlContainer.GetDSN())
	defer db.Close()
	s.Require().NoError(err)


	var user domain.User

	err = db.QueryRow(`
		SELECT id, email, role from users;
	`).Scan(&user.ID, &user.Email, &user.Role)
	s.Require().NoError(err)

	s.Require().Equal(userID, user.ID)
	s.Require().Equal(exampleEmail, user.Email)
	s.Require().Equal(exampleRole, user.Role)
}