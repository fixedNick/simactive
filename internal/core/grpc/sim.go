package grpc

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	pb "simactive/api/generated/github.com/fixedNick/SimHelper"
	"simactive/internal/core"
	"simactive/internal/infrastructure/repoerrors"
	"strconv"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SimService interface {
	Add(ctx context.Context, s *core.Sim) (int, error)
	Remove(ctx context.Context, id int) error
	GetSimList(ctx context.Context) (*core.List[*core.Sim], error)
	ActivateSim(ctx context.Context, id int) error
	BlockSim(ctx context.Context, id int) error
	GetUsedServiceList(ctx context.Context, id int) (core.List[*core.Used], error)
}

type GRPCSimService struct {
	pb.UnimplementedSimServer

	timeout    time.Duration
	simService SimService
	logger     *slog.Logger
}

func NewGRPCSimService(logger *slog.Logger, ss SimService, timeout time.Duration) GRPCSimService {
	return GRPCSimService{
		simService: ss,
		timeout:    timeout,
		logger:     logger,
	}
}
func (gs GRPCSimService) AddSim(ctx context.Context, req *pb.AddSimRequest) (*pb.AddSimResponse, error) {

	gs.logger.Info("AddSim request", slog.Any("req", req))

	number := req.SimData.Number

	// Validates the length of phone number. Length of phone number should be in range from 11 to 15.
	// And all the digits in the phone number should be in range from 0 to 9
	// Example: 1 999 888 77 66
	if !validatePhoneNumber(number) {
		return nil, status.Errorf(codes.InvalidArgument, "Bad phone number. Please use correct phone number. Example: 1 999 888 77 66")
	}

	if req.SimData.ProviderName == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Provider name is required. Example: Vodafone, Beeline, Tele2, etc.")
	}

	ctx, cancel := context.WithTimeout(ctx, gs.timeout)
	defer cancel()

	provider := core.Provider{}
	provider.SetName(req.SimData.GetProviderName())
	sim := core.NewSim(0, req.SimData.Number, &provider, req.SimData.IsActivated, req.SimData.ActivateUntil, req.SimData.IsBlocked)
	id, err := gs.simService.Add(ctx, &sim)
	if err != nil {
		if errors.Is(err, repoerrors.ErrAlreadyExists) {
			return nil, status.Errorf(codes.AlreadyExists, "sim card with number %s already exists", sim.Number())
		}

		gs.logger.Error("Failed to add sim card", slog.Any("sim", sim), "err", err)
		return nil, ErrInternal
	}
	return &pb.AddSimResponse{
		Id:      int32(id),
		Message: fmt.Sprintf("sim card with number %s added. Sim id is `%d`", sim.Number(), id),
	}, nil
}

// validatePhoneNumber checks if the phone number is within a valid length range.
//
// It takes a string parameter 'number' representing the phone number.
// It returns a boolean value indicating if the phone number is valid.
func validatePhoneNumber(number string) bool {
	if len(number) < 8 || len(number) > 15 {
		return false
	}

	_, err := strconv.Atoi(number)
	return err == nil
}

func (gs GRPCSimService) DeleteSim(ctx context.Context, req *pb.DeleteSimRequest) (*pb.DeleteSimResponse, error) {

	if req.GetId() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid id, id must be greater than 0")
	}

	ctx, cancel := context.WithTimeout(ctx, gs.timeout)
	defer cancel()
	if err := gs.simService.Remove(ctx, int(req.GetId())); err != nil {

		if errors.Is(err, repoerrors.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "sim card with id %d not found", req.GetId())
		}
		return nil, ErrInternal
	}
	return &pb.DeleteSimResponse{
		Id: req.GetId(),
	}, nil
}
func (gs GRPCSimService) GetSimList(ctx context.Context, req *pb.Empty) (*pb.SimList, error) {
	ctx, cancel := context.WithTimeout(ctx, gs.timeout)
	defer cancel()

	list, err := gs.simService.GetSimList(ctx)
	if err != nil {
		return nil, ErrInternal
	}

	if list == nil {
		return &pb.SimList{}, nil
	}

	var response pb.SimList
	response.SimList = make([]*pb.SimData, 0, len(*list))
	for _, sim := range *list {
		p := sim.Provider()
		response.SimList = append(response.SimList, &pb.SimData{
			ID:     int32(sim.Id()),
			Number: sim.Number(),
			Provider: &pb.ProviderData{
				Id:   int32(p.Id()),
				Name: p.Name(),
			},
			IsActivated:   sim.IsActivated(),
			IsBlocked:     sim.IsBlocked(),
			ActivateUntil: sim.ActivateUntil(),
		})
	}

	return &response, nil
}
func (gs GRPCSimService) ActivateSim(ctx context.Context, req *pb.ActivateSimRequest) (*pb.ActivateSimResponse, error) {

	if req.GetId() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid id, id must be greater than 0")
	}

	ctx, cancel := context.WithTimeout(ctx, gs.timeout)
	defer cancel()

	if err := gs.simService.ActivateSim(ctx, int(req.Id)); err != nil {
		fmt.Println("ERROR FROM UPDATE", err)
		if errors.Is(err, repoerrors.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "sim card with id %d not found", req.Id)
		}

		return nil, ErrInternal
	}

	return &pb.ActivateSimResponse{
		IsActivated: true,
	}, nil
}
func (gs GRPCSimService) SetSimBlocked(ctx context.Context, req *pb.SSBRequest) (*pb.SSBResponse, error) {
	if req.GetId() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid id, id must be greater than 0")
	}

	ctx, cancel := context.WithTimeout(ctx, gs.timeout)
	defer cancel()

	if err := gs.simService.BlockSim(ctx, int(req.Id)); err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "sim card with id %d not found", req.Id)
		}
		return nil, ErrInternal
	}

	return &pb.SSBResponse{
		IsBlocked: true,
	}, nil
}

func (gs GRPCSimService) GetFreeServices(ctx context.Context, req *pb.GetFreeServRequest) (*pb.GetFreeServResponse, error) {
	panic("implement")
}

func (gs GRPCSimService) GetUsedServices(ctx context.Context, req *pb.GetUsedServRequest) (*pb.GetUsedServResponse, error) {
	// validate that sim id is greater than 0
	if req.GetSimId() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid sim id, sim id must be greater than 0")
	}

	ctx, cancel := context.WithTimeout(ctx, gs.timeout)
	defer cancel()

	list, err := gs.simService.GetUsedServiceList(ctx, int(req.GetSimId()))
	if err != nil {
		if err == repoerrors.ErrNotFound {
			return nil, status.Errorf(codes.NotFound, "sim card with id %d not found", req.GetSimId())
		}
		return nil, ErrInternal
	}

	if list == nil {
		return &pb.GetUsedServResponse{}, nil
	}

	var response pb.GetUsedServResponse
	response.UsedServices = make([]*pb.UsedService, 0, len(list))
	for _, used := range list {
		response.UsedServices = append(response.UsedServices, &pb.UsedService{
			ServiceId:   int32(used.ServiceID()),
			IsBlocked:   used.IsBlocked(),
			BlockedInfo: used.BlockedInfo(),
		})
	}

	return &response, nil
}
