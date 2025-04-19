//go:build integration

package integration

import (
	"context"

	"github.com/Ranik23/avito-tech-spring/internal/service"
)



func (s *TestSuite) TestLoginSuccess() {
	exampleEmail := "test@example.com"
	exampleRole := "Client"
	examplePassword := "strongpassword"

	_, err := s.service.Register(context.Background(), exampleEmail, examplePassword, exampleRole)
	s.Require().NoError(err)

	token, err := s.service.Login(context.Background(), exampleEmail, examplePassword)
	s.Require().NoError(err)
	s.Require().NotEmpty(token)
}


func (s *TestSuite) TestLoginInvalidPassword() {
	exampleEmail := "wrongpass@example.com"
	exampleRole := "Client"
	correctPassword := "correct"
	wrongPassword := "wrong"

	_, err := s.service.Register(context.Background(), exampleEmail, correctPassword, exampleRole)
	s.Require().NoError(err)

	_, err = s.service.Login(context.Background(), exampleEmail, wrongPassword)
	s.Require().ErrorIs(err, service.ErrInvalidCredentials)
}


func (s *TestSuite) TestLoginUserNotFound() {
	_, err := s.service.Login(context.Background(), "ghost@example.com", "whatever")
	s.Require().ErrorIs(err, service.ErrUserNotFound)
}


