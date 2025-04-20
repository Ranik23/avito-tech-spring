//go:build unit

package http

import (
	"testing"
	"time"

	"github.com/Ranik23/avito-tech-spring/internal/models/domain"
	"github.com/stretchr/testify/require"
)

func TestFromDomainProductToDtoPostProductResp(t *testing.T) {
	product := &domain.Product{
		ID:         "productID123",
		ReceptionID: "receptionID123",
		Type:       "TypeA",
		DateTime:   time.Date(2023, time.March, 5, 14, 30, 0, 0, time.UTC),
	}

	result := FromDomainProductToDtoPostProductResp(product)

	require.Equal(t, product.ID, result.Id)
	require.Equal(t, product.ReceptionID, result.ReceptionID)
	require.Equal(t, product.Type, result.Type)
	require.Equal(t, product.DateTime.String(), result.DateTime)
}


func TestFromDomainReceptionToCreateReceptionResp(t *testing.T) {
	reception := &domain.Reception{
		ID:        "receptionID123",
		PvzID:     "pvzID123",
		Status:    "Active",
		DateTime:  time.Date(2023, time.April, 15, 12, 0, 0, 0, time.UTC),
	}

	result := FromDomainReceptionToCreateReceptionResp(reception)

	require.Equal(t, reception.ID, result.Id)
	require.Equal(t, reception.PvzID, result.PvzId)
	require.Equal(t, reception.Status, result.Status)
	require.Equal(t, reception.DateTime.String(), result.DateTime)
}

func TestFromDomainPVZToCreatePvzResp(t *testing.T) {
	pvz := &domain.Pvz{
		ID:              "pvzID123",
		City:            "Moscow",
		RegistrationDate: time.Date(2023, time.January, 10, 9, 30, 0, 0, time.UTC),
	}

	result := FromDomainPVZToCreatePvzResp(pvz)

	require.Equal(t, pvz.ID, result.Id)
	require.Equal(t, pvz.City, result.City)
	require.Equal(t, pvz.RegistrationDate.String(), result.RegistrationDate)
}

func TestFromDomainReceptionToCloseReseptionResp(t *testing.T) {
	reception := &domain.Reception{
		ID:        "receptionID123",
		PvzID:     "pvzID123",
		Status:    "Closed",
		DateTime:  time.Date(2023, time.May, 5, 16, 0, 0, 0, time.UTC),
	}

	result := FromDomainReceptionToCloseReseptionResp(reception)

	require.Equal(t, reception.ID, result.Id)
	require.Equal(t, reception.PvzID, result.PvzId)
	require.Equal(t, reception.Status, result.Status)
	require.Equal(t, reception.DateTime.String(), result.DateTime)
}

