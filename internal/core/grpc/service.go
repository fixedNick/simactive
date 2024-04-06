package grpc

import (
	"context"
	pb "simactive/api/generated/github.com/fixedNick/SimHelper"
	"simactive/internal/core"
	"time"
)

type ServiceService interface {
	Add(ctx context.Context, s core.Service) error
	Remove(ctx context.Context, id int) error
	GetServiceList(ctx context.Context) (*core.List[*core.Service], error)
}

type GRPCServiceService struct {
	pb.UnimplementedServiceServer

	timeout    time.Duration
	simService ServiceService
}

func NewGRPCServiceService(ss ServiceService, timeout time.Duration) GRPCServiceService {
	return GRPCServiceService{
		simService: ss,
		timeout:    timeout,
	}
}

func (ss GRPCServiceService) AddService(ctx context.Context, req *pb.AddServiceRequest) (*pb.AddServiceResponse, error) {
	panic("implement")
}
func (ss GRPCServiceService) DeleteService(ctx context.Context, req *pb.DeleteServiceRequest) (*pb.DeleteServiceResponse, error) {
	panic("implement")
}
func (ss GRPCServiceService) GetSetviceList(ctx context.Context, req *pb.Empty) (*pb.GSLResponse, error) {
	panic("implement")
}
