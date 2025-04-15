package grpc

import (
	"context"

	"github.com/Ranik23/avito-tech-spring/api/proto/gen/pvz_v1"
	"github.com/Ranik23/avito-tech-spring/internal/service"
)


type PVZServer struct {
	pvz_v1.UnimplementedPVZServiceServer
	srv service.Service
}


func NewPVZServer(srv service.Service) pvz_v1.PVZServiceServer {
	return &PVZServer{
		srv: srv,
	}
}


func (pvz *PVZServer) GetPVZList(ctx context.Context, req *pvz_v1.GetPVZListRequest) (*pvz_v1.GetPVZListResponse, error) {
	return nil, nil
}