//go:build integration

package integration

import (
	"context"
	"time"
)


func (s *TestSuite) TestGetPVZSInfo() {
	ctx := context.Background()

	pvz, err := s.service.CreatePVZ(ctx, "Moscow")
	s.Require().NoError(err)

	_, err = s.service.StartReception(ctx, pvz.ID)
	s.Require().NoError(err)

	_, err = s.service.AddProduct(ctx, pvz.ID, "box")
	s.Require().NoError(err)

	start := time.Now().Add(-1 * time.Hour)
	end := time.Now().Add(1 * time.Hour)

	pvzInfo, err := s.service.GetPVZSInfo(ctx, start, end, 0, 10)
	s.Require().NoError(err)
	s.Require().GreaterOrEqual(len(pvzInfo), 1)
}