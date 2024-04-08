package grpc

import (
	"context"
	"errors"
	"fmt"
	pb "simactive/api/generated/github.com/fixedNick/SimHelper"
	"simactive/internal/core"
	"simactive/internal/repository"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ServiceService interface {
	Add(ctx context.Context, s *core.Service) (int, error)
	Remove(ctx context.Context, id int) error
	GetServiceList(ctx context.Context) (*core.List[*core.Service], error)
}

type GRPCServiceService struct {
	pb.UnimplementedServiceServer

	timeout        time.Duration
	serviceService ServiceService
}

func NewGRPCServiceService(ss ServiceService, timeout time.Duration) GRPCServiceService {
	return GRPCServiceService{
		serviceService: ss,
		timeout:        timeout,
	}
}

// AddService adds a new service to the GRPC service.
//
// ctx: The context for the operation.
// req: The request to add a service.
// Returns the response indicating if the service was added successfully or an error.
func (gss GRPCServiceService) AddService(ctx context.Context, req *pb.AddServiceRequest) (*pb.AddServiceResponse, error) {

	name := req.GetName()
	if name == "" {
		return nil, status.Error(codes.InvalidArgument, "service name cannot be empty")
	}

	if len(name) > 64 {
		return nil, status.Error(codes.InvalidArgument, "service name cannot be longer than 64 characters")
	}

	ctx, cancel := context.WithTimeout(ctx, gss.timeout)
	defer cancel()

	service := core.Service{}.WithName(name)

	id, err := gss.serviceService.Add(ctx, &service)
	if err != nil {
		if err == repository.ErrAlreadyExists {
			return nil, status.Error(codes.AlreadyExists, fmt.Sprintf("service with name %s already exists", name))
		}
		return nil, err
	}

	return &pb.AddServiceResponse{
		Id: int32(id),
	}, nil
}

// DeleteService deletes a service.
// ctx context.Context, req *pb.DeleteServiceRequest
// *pb.DeleteServiceResponse, error
func (gss GRPCServiceService) DeleteService(ctx context.Context, req *pb.DeleteServiceRequest) (*pb.DeleteServiceResponse, error) {

	if req.GetID() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid id, id must be greater than 0")
	}

	ctx, cancel := context.WithTimeout(ctx, gss.timeout)
	defer cancel()

	id := int(req.GetID())
	if id <= 0 {
		return nil, fmt.Errorf("invalid id")
	}

	err := gss.serviceService.Remove(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "service with id %d not found", id)
		}
		return nil, err
	}

	return &pb.DeleteServiceResponse{
		Id: req.GetID(),
	}, nil
}

// GetServiceList retrieves a list of services.
//
// ctx: the context for the request
// req: the request parameter
// Returns a GSLResponse whitch contains a list of services and an error.
func (gss GRPCServiceService) GetServiceList(ctx context.Context, req *pb.Empty) (*pb.GSLResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, gss.timeout)
	defer cancel()

	services, err := gss.serviceService.GetServiceList(ctx)
	if err != nil {
		return nil, err
	}

	if services == nil {
		return &pb.GSLResponse{}, nil
	}

	pbServices := make([]*pb.ServiceData, 0, len(*services))
	for _, s := range *services {
		pbServices = append(pbServices, serviceToPB(s))
	}

	return &pb.GSLResponse{
		Services: pbServices,
	}, nil
}

// serviceToPB converts a core.Service to a pb.ServiceData.
//
// s *core.Service - input core.Service
// *pb.ServiceData - returned pb.ServiceData
func serviceToPB(s *core.Service) *pb.ServiceData {
	return &pb.ServiceData{
		Id:   int32(s.Id()),
		Name: s.Name(),
	}
}
