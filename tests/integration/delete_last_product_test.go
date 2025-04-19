//go:build integration

package integration

import (
	"context"
	"database/sql"
)




func (s *TestSuite) TestDeleteLastProduct() {
	ctx := context.Background()

	pvz, err := s.service.CreatePVZ(ctx, "Moscow")
	s.Require().NoError(err)

	_, err = s.service.StartReception(ctx, pvz.ID)
	s.Require().NoError(err)

	_, err = s.service.AddProduct(ctx, pvz.ID, "box")
	s.Require().NoError(err)

	err = s.service.DeleteLastProduct(ctx, pvz.ID)
	s.Require().NoError(err)

	db, err := sql.Open("postgres", s.psqlContainer.GetDSN())
	s.Require().NoError(err)
	defer db.Close()

	var count int
	err = db.QueryRow(`SELECT COUNT(*) FROM product`).Scan(&count)
	s.Require().NoError(err)
	s.Require().Equal(0, count)
}