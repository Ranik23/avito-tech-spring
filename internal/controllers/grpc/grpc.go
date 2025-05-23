package grpc

import (
	"context"

	"github.com/Ranik23/avito-tech-spring/api/proto/gen/pvz_v1"
	"github.com/Ranik23/avito-tech-spring/internal/models/converter/grpc"
	"github.com/Ranik23/avito-tech-spring/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


type PVZServer struct {
	pvz_v1.UnimplementedPVZServiceServer
	service service.Service
}


func NewPVZServer(service service.Service) *PVZServer {
	return &PVZServer{
		service: service,
	}
}


func (pvz *PVZServer) GetPVZList(ctx context.Context, req *pvz_v1.GetPVZListRequest) (*pvz_v1.GetPVZListResponse, error) {
	pvzs, err := pvz.service.GetPVZList(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal Server Error")
	}
	
	grpcPVZs := grpc.FromDomainPvzListToGRPCList(pvzs)

	return &pvz_v1.GetPVZListResponse{
		Pvzs: grpcPVZs,
	}, nil
}