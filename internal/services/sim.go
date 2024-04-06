package services

import (
	"context"
	"simactive/internal/core"
)

type SimService struct {
	repository Repository[core.Sim]
}

func NewSimService(repository Repository[core.Sim]) *SimService {
	ss := &SimService{
		repository: repository,
	}
	return ss
}
func (ss *SimService) Add(ctx context.Context, s core.Sim) (int, error) {
	return ss.repository.Save(ctx, s)
}
func (ss *SimService) Remove(ctx context.Context, id int) error {
	return ss.repository.Remove(ctx, id)
}
func (ss *SimService) GetSimList(ctx context.Context) (*core.List[core.Sim], error) {
	return ss.repository.GetList(ctx)
}
func (ss *SimService) ActivateSim(ctx context.Context, id int) error {
	list, err := ss.GetSimList(ctx)
	if err != nil {
		return err
	}

	sim, err := list.PtrByKey(id)
	if err != nil {
		return nil
	}
	sim.SetActivated(true)
	return ss.repository.Update(ctx, sim)
}
func (ss *SimService) BlockSim(ctx context.Context, id int) error {
	list, err := ss.GetSimList(ctx)
	if err != nil {
		return err
	}

	sim, err := list.PtrByKey(id)
	if err != nil {
		return nil
	}
	sim.SetBlocked(true)
	return ss.repository.Update(ctx, sim)
}
