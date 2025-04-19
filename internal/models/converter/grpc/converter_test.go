//go:build unit

package grpc

import (
	"testing"
	"time"

	"github.com/Ranik23/avito-tech-spring/internal/models/domain"
	"github.com/stretchr/testify/require"
)



func TestConvert(t *testing.T) {
	exampleDomainPvz := []domain.Pvz{
		{
			ID: "1",
			RegistrationDate: time.Date(2010, time.April, 10, 12, 45, 23, 33, time.Local),
			City: "Moscow",
		},
	}
	pvz_v1 := FromDomainPvzListToGRPCList(exampleDomainPvz)

	for i, pvz := range pvz_v1 {
		require.Equal(t, exampleDomainPvz[i].ID, pvz.Id)
		require.Equal(t, exampleDomainPvz[i].City, pvz.City)
		require.Equal(t, exampleDomainPvz[i].RegistrationDate, pvz.RegistrationDate.AsTime().Local())
	}
}