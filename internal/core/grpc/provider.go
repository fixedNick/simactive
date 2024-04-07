package grpc

import (
	"context"
	pb "simactive/api/generated/github.com/fixedNick/SimHelper"
	"simactive/internal/core"
	"time"
)

type ProviderService interface {
	GetProviderList(ctx context.Context) (*core.List[*core.Provider], error)
}

type GRPCProviderService struct {
	pb.UnimplementedProviderServer

	timeout    time.Duration
	simService ProviderService
}

func NewGRPCProviderService(ss ProviderService, timeout time.Duration) GRPCProviderService {
	return GRPCProviderService{
		simService: ss,
		timeout:    timeout,
	}
}

func GetProviderList(ctx context.Context, req *pb.Empty) (*pb.ProviderList, error) {
	panic("implement")
}
