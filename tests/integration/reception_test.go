//go:build integration

package integration

import (
	"context"
	"database/sql"
)

func (s *TestSuite) TestStartAndCloseReception() {
	ctx := context.Background()

	pvz, err := s.service.CreatePVZ(ctx, "Moscow")
	s.Require().NoError(err)

	reception, err := s.service.StartReception(ctx, pvz.ID)
	s.Require().NoError(err)
	s.Require().Equal(pvz.ID, reception.PvzID)

	closedReception, err := s.service.CloseReception(ctx, pvz.ID)
	s.Require().NoError(err)
	s.Require().Equal(reception.ID, closedReception.ID)

	db, err := sql.Open("postgres", s.psqlContainer.GetDSN())
	s.Require().NoError(err)
	defer db.Close()

	var status string
	err = db.QueryRow(`SELECT status FROM reception WHERE id = $1`, reception.ID).Scan(&status)
	s.Require().NoError(err)
	s.Require().Equal("closed", status)
}