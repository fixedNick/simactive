package grpc

import (
	"context"
	pb "simactive/api/generated/github.com/fixedNick/SimHelper"
	"time"
)

type UsedService interface {
	UseSimForService(ctx context.Context, simId int, serviceId int) error
}

type GRPCUsedService struct {
	pb.UnimplementedUsedServer

	timeout     time.Duration
	usedService UsedService
}

func NewGRPCUsedService(ss UsedService, timeout time.Duration) GRPCUsedService {
	return GRPCUsedService{
		usedService: ss,
		timeout:     timeout,
	}
}

// UseSimForService is a method that connects sim to a service.
//
// ctx - The context in which the function operates.
// req - The request containing information about the simulated service.
// Returns a response indicating if the service is used successfully, otherwise an error.
func (gus GRPCUsedService) UseSimForService(ctx context.Context, req *pb.USFSRequest) (*pb.USFSResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, gus.timeout)
	defer cancel()

	if err := gus.usedService.UseSimForService(ctx, int(req.GetSimID()), int(req.GetServiceID())); err != nil {
		return nil, err
	}

	return &pb.USFSResponse{
		IsUsed: true,
	}, nil
}
