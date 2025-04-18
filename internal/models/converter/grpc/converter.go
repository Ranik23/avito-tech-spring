package grpc

import (
	"github.com/Ranik23/avito-tech-spring/api/proto/gen/pvz_v1"
	"github.com/Ranik23/avito-tech-spring/internal/models/domain"
	"google.golang.org/protobuf/types/known/timestamppb"
)


func FromDomainPvzListToGRPCList(pvzs []domain.Pvz) []*pvz_v1.PVZ {

	var response []*pvz_v1.PVZ

	for _, pvz := range pvzs {
		response = append(response, &pvz_v1.PVZ{
			Id: pvz.ID,
			City: pvz.City,
			RegistrationDate: timestamppb.New(pvz.RegistrationDate),
		})
	}
	return response
}