//go:build integration

package integration

import (
	"context"
	"database/sql"

	"github.com/Ranik23/avito-tech-spring/internal/models/domain"
)

func (s *TestSuite) TestAddProduct() {

	pvz, err := s.service.CreatePVZ(context.Background(), "Moscow")
	s.Require().NoError(err)


	reception, err := s.service.StartReception(context.Background(), pvz.ID)
	s.Require().NoError(err)
	

	_, err = s.service.AddProduct(context.Background(), pvz.ID, "box")
	s.Require().NoError(err)


	db, err := sql.Open("postgres", s.psqlContainer.GetDSN())
	defer db.Close()
	s.Require().NoError(err)

	var (
		product domain.Product	
	)

	err = db.QueryRow(`
		SELECT id, type, reception_id FROM product
	`).Scan(&product.ID, &product.Type, &product.ReceptionID)
	s.Require().NoError(err)


	s.Require().Equal(reception.ID, product.ReceptionID)
	s.Require().Equal("box", product.Type)
}




