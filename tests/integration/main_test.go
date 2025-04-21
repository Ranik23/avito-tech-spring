//go:build integration

package integration

import (
	"context"
	"database/sql"

	"github.com/Ranik23/avito-tech-spring/internal/models/domain"
)



func (s *TestSuite) TestMain() {

	examplePvz, err := s.service.CreatePVZ(context.Background(), "Moscow")
	s.Require().NoError(err)

	exampleReception, err := s.service.StartReception(context.Background(), examplePvz.ID)
	s.Require().NoError(err)

	for i := 0; i < 50; i++ {
		_, err := s.service.AddProduct(context.Background(), examplePvz.ID, "box")
		s.Require().NoError(err)
	}

	_, err = s.service.CloseReception(context.Background(), examplePvz.ID)
	s.Require().NoError(err)


	db, err := sql.Open("postgres", s.psqlContainer.GetDSN())
	s.Require().NoError(err)
	defer db.Close()

	var count int
	// проверяем, что там 50 продуктов
	err = db.QueryRow(`SELECT COUNT(*) FROM product;`).Scan(&count)
	s.Require().NoError(err)

	s.Require().Equal(50, count)


	var recepion domain.Reception
	// проверяем, что приемка закрыта и она та, которую положили
	err = db.QueryRow(`
		SELECT id, pvz_id, status FROM reception;
	`).Scan(&recepion.ID, recepion.PvzID, recepion.Status)
	s.Require().NoError(err)

	s.Require().Equal("closed", recepion.Status)
	s.Require().Equal(examplePvz.ID, recepion.PvzID)
	s.Require().Equal(exampleReception.ID, recepion.ID)

	var pvz domain.Pvz

	err = db.QueryRow(`
		SELECT id, city FROM pvz;
	`).Scan(&pvz.ID, &pvz.City)
	s.Require().NoError(err)

	s.Require().Equal(examplePvz.ID, pvz.ID)
	s.Require().Equal(examplePvz.City, pvz.City)
}