package http

import (
	"github.com/Ranik23/avito-tech-spring/internal/models/domain"
	"github.com/Ranik23/avito-tech-spring/internal/models/dto"
)




func FromDomainProductToDtoPostProductResp(product *domain.Product) *dto.PostProductResp {
	return &dto.PostProductResp{
		DateTime: product.DateTime.String(),
		Id: product.ID,
		ReceptionID: product.ReceptionID,
		Type: product.Type,
	}
}

func FromDomainReceptionToCreateReceptionResp(reception *domain.Reception) *dto.CreateReceptionResp {
	return &dto.CreateReceptionResp{
		DateTime: reception.DateTime.String(),
		Id: reception.ID,
		PvzId: reception.PvzID,
		Status: reception.Status,
	}
}


func FromDomainPVZToCreatePvzResp(pvz *domain.Pvz) *dto.CreatePvzResp {
	return &dto.CreatePvzResp{
		Id: pvz.ID,
		City: pvz.City,
		RegistrationDate: pvz.RegistrationDate.String(),
	}
}


func FromDomainReceptionToCloseReseptionREsp(reception *domain.Reception) *dto.CloseReceptionResp {
	return &dto.CloseReceptionResp{
		DateTime: reception.DateTime.String(),
		Id: reception.ID,
		PvzId: reception.PvzID,
		Status: reception.Status,
	}
}

