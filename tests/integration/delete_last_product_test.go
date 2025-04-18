//go:build integration

package integration


// import (
// 	"context"
// 	"database/sql"

// 	"github.com/Ranik23/avito-tech-spring/internal/models/domain"
// )

// func (s *TestSuite) TestAddProductSuccess() {

// 	pvz, err := s.pvzService.CreatePVZ(context.Background(), "Moscow")
// 	s.Require().NoError(err)


// 	reception, err := s.pvzService.StartReception(context.Background(), pvzID)
// 	s.Require().NoError(err)


// 	productID, err := s.pvzService.AddProduct(context.Background(), pvzID, "box")
// 	s.Require().NoError(err)

// 	db, err := sql.Open("postgres", s.psqlContainer.GetDSN())
// 	defer db.Close()
// 	s.Require().NoError(err)


// 	var product domain.Product

// 	err = db.QueryRow("SELECT * FROM products").Scan(&product.ID, &product.ReceptionID, &product.Type, &product.DateTime)
// 	s.Require().NoError(err)


// 	var reception domain.Reception

// 	err = db.QueryRow("SELECT * FROM receptions WHERE status = 'open'").Scan(
// 		&reception.ID, &reception.PvzID, &reception.Status, &reception.DateTime)
// 	s.Require().NoError(err)
	
// 	s.Require().Equal(pvzID, reception.PvzID, "PVZ ID должен совпадать")
//     s.Require().Equal(receptionID, product.ReceptionID, "ID приема продукта должен совпадать")
//     s.Require().Equal("box", product.Type, "Тип продукта должен быть 'box'")
//     s.Require().Equal(productID, product.ID, "ID продукта должен совпадать")
// }	


