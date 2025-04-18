//go:build intergation

package integration

import (
	"context"
	"database/sql"

	"github.com/Ranik23/avito-tech-spring/internal/models/domain"
)

func (s *TestSuite) TestAddProduct() {

	pvzID, err := s.pvzService.CreatePVZ(context.Background(), "Moscow")
	s.Require().NoError(err)


	receptionID, err := s.pvzService.StartReception(context.Background(), pvzID)
	s.Require().NoError(err)


	productID, err := s.pvzService.AddProduct(context.Background(), pvzID, "box")
	s.Require().NoError(err)

	db, err := sql.Open("postgres", s.psqlContainer.GetDSN())
	defer db.Close()
	s.Require().NoError(err)


	var product domain.Product

	err = db.QueryRow("SELECT * from products").Scan(&product.ID, &product.ReceptionID, &product.Type, &product.DateTime)
	s.Require().NoError(err)

	s.Require().Equal(receptionID, product.ReceptionID)
	s.Require().Equal("box", product.Type)
	s.Require().Equal(productID, product.ID)

}	