//go:build integration

package integration

import (
	"context"
	"database/sql"

	"github.com/Ranik23/avito-tech-spring/internal/models/domain"
)

func (s *TestSuite) TestCreatePVZSuccess() {

	exampleCity := "Moscow"

	_, err := s.service.CreatePVZ(context.Background(), exampleCity)
	s.Require().NoError(err)

	db, err := sql.Open("postgres", s.psqlContainer.GetDSN())
	defer db.Close()
	s.Require().NoError(err)

	var pvz domain.Pvz

	err = db.QueryRow(`
		SELECT id, registration_date, city FROM pvz;
	`).Scan(&pvz.ID, &pvz.RegistrationDate, &pvz.City)
	s.Require().NoError(err)

	s.Require().Equal(pvz.City, exampleCity)
}
