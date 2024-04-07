package grpc

import (
	"context"
	pb "simactive/api/generated/github.com/fixedNick/SimHelper"
	"time"
)

type UsedService interface {
	UseSimForService(ctx context.Context, simId int, serviceId int)
}

type GRPCUsedService struct {
	pb.UnimplementedUsedServer

	timeout    time.Duration
	simService UsedService
}

func NewGRPCUsedService(ss UsedService, timeout time.Duration) GRPCUsedService {
	return GRPCUsedService{
		simService: ss,
		timeout:    timeout,
	}
}

func UseSimForService(ctx context.Context, req *pb.USFSRequest) (*pb.USFSResponse, error) {
	panic("implement")
}
