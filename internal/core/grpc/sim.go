package grpc

import (
	"context"
	"fmt"
	pb "simactive/api/generated/github.com/fixedNick/SimHelper"
	"simactive/internal/core"
	"time"
)

type SimService interface {
	Add(ctx context.Context, s core.Sim) (int, error)
	Remove(ctx context.Context, id int) error
	GetSimList(ctx context.Context) (*core.List[core.Sim], error)
	ActivateSim(ctx context.Context, id int) error
	BlockSim(ctx context.Context, id int) error
}

type GRPCSimService struct {
	pb.UnimplementedSimServer

	timeout    time.Duration
	simService SimService
}

func NewGRPCSimService(ss SimService, timeout time.Duration) GRPCSimService {
	return GRPCSimService{
		simService: ss,
		timeout:    timeout,
	}
}
func (gs GRPCSimService) AddSim(ctx context.Context, req *pb.AddSimRequest) (*pb.AddSimResponse, error) {
	// validate
	ctx, cancel := context.WithTimeout(ctx, gs.timeout)
	defer cancel()

	sim := core.NewSim(0, req.SimData.Number, int(req.SimData.ProviderID), req.SimData.IsActivated, req.SimData.ActivateUntil, req.SimData.IsBlocked)
	id, err := gs.simService.Add(ctx, sim)
	if err != nil {
		return nil, ErrInternal
	}
	return &pb.AddSimResponse{
		IsAdded: true,
		Message: fmt.Sprintf("sim card with number %s added. Sim id is `%d`", sim.Number(), id),
	}, nil
}
func (gs GRPCSimService) DeleteSim(ctx context.Context, req *pb.DeleteSimRequest) (*pb.DeleteSimResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, gs.timeout)
	defer cancel()
	if err := gs.simService.Remove(ctx, int(req.GetId())); err != nil {
		return nil, ErrInternal
	}
	return &pb.DeleteSimResponse{
		IsDeleted: true,
	}, nil
}
func (gs GRPCSimService) GetSimList(ctx context.Context, req *pb.Empty) (*pb.SimList, error) {
	ctx, cancel := context.WithTimeout(ctx, gs.timeout)
	defer cancel()

	list, err := gs.simService.GetSimList(ctx)
	if err != nil {
		return nil, ErrInternal
	}

	var response pb.SimList
	response.SimList = make([]*pb.SimData, 0, len(*list))
	for _, sim := range *list {
		response.SimList = append(response.SimList, &pb.SimData{
			ID:            int32(sim.Id()),
			Number:        sim.Number(),
			ProviderID:    int32(sim.ProviderID()),
			IsActivated:   sim.IsActivated(),
			IsBlocked:     sim.IsBlocked(),
			ActivateUntil: sim.ActivateUntil(),
		})
	}

	return &response, nil
}
func (gs GRPCSimService) ActivateSim(ctx context.Context, req *pb.ActivateSimRequest) (*pb.ActivateSimResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, gs.timeout)
	defer cancel()

	if err := gs.simService.ActivateSim(ctx, int(req.Id)); err != nil {
		return nil, ErrInternal
	}

	return &pb.ActivateSimResponse{
		IsActivated: true,
	}, nil
}
func (gs GRPCSimService) SetSimBlocked(ctx context.Context, req *pb.SSBRequest) (*pb.SSBResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, gs.timeout)
	defer cancel()

	if err := gs.simService.BlockSim(ctx, int(req.Id)); err != nil {
		return nil, ErrInternal
	}

	return &pb.SSBResponse{
		IsBlocked: true,
	}, nil
}
