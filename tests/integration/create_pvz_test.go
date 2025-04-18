//go:build integration

package integration

import (
	"context"
	"database/sql"

	"github.com/Ranik23/avito-tech-spring/internal/models/domain"
)

func (s *TestSuite) TestCreatePVZ() {
	pvzID, err := s.pvzService.CreatePVZ(context.Background(), "Moscow")
	s.Require().NoError(err)

	db, err := sql.Open("postgres", s.psqlContainer.GetDSN())
	defer db.Close()
	s.Require().NoError(err)

	var pvz domain.Pvz

	err = db.QueryRow("SELECT * FROM pvz").Scan(&pvz.ID, &pvz.RegistrationDate, &pvz.City)
	s.Require().NoError(err)

	s.Require().Equal(pvzID, pvz.ID)
	s.Require().Equal("Moscow", pvz.City)
}


func (s *TestSuite) TestCreateMultiplePVZs() {
	pvzID1, err := s.pvzService.CreatePVZ(context.Background(), "Moscow")
	s.Require().NoError(err)

	pvzID2, err := s.pvzService.CreatePVZ(context.Background(), "Kazan")
	s.Require().NoError(err)

	db, err := sql.Open("postgres", s.psqlContainer.GetDSN())
	defer db.Close()
	s.Require().NoError(err)

	var pvz1, pvz2 domain.Pvz

	err = db.QueryRow("SELECT * FROM pvz WHERE city = $1", "Moscow").Scan(&pvz1.ID, &pvz1.RegistrationDate, &pvz1.City)
	s.Require().NoError(err)

	err = db.QueryRow("SELECT * FROM pvz WHERE city = $1", "Kazan").Scan(&pvz2.ID, &pvz2.RegistrationDate, &pvz2.City)
	s.Require().NoError(err)

	s.Require().Equal(pvzID1, pvz1.ID)
	s.Require().Equal("Moscow", pvz1.City)

	s.Require().Equal(pvzID2, pvz2.ID)
	s.Require().Equal("Saint Petersburg", pvz2.City)
}





