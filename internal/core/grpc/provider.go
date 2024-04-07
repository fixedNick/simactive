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

	timeout         time.Duration
	providerService ProviderService
}

func NewGRPCProviderService(ss ProviderService, timeout time.Duration) GRPCProviderService {
	return GRPCProviderService{
		providerService: ss,
		timeout:         timeout,
	}
}

// GetProviderList retrieves a list of providers.
func (gps GRPCProviderService) GetProviderList(ctx context.Context, _ *pb.Empty) (*pb.ProviderList, error) {
	ctx, cancel := context.WithTimeout(ctx, gps.timeout)
	defer cancel()

	providerList, err := gps.providerService.GetProviderList(ctx)
	if err != nil {
		return nil, err
	}

	// if providerList is nil, return empty list
	if providerList == nil {
		return &pb.ProviderList{}, nil
	}

	providers := make([]*pb.ProviderData, 0, len(*providerList))
	for _, provider := range *providerList {
		providers = append(providers, &pb.ProviderData{
			Id:   int32(provider.Id()),
			Name: provider.Name(),
		})
	}
	return &pb.ProviderList{Providers: providers}, nil
}
