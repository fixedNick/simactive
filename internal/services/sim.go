package services

import (
	"context"
	"errors"
	"simactive/internal/core"
	repository "simactive/internal/infrastructure"
	"simactive/internal/infrastructure/repoerrors"
)

type SimService struct {
	repository *repository.Repository
}

func NewSimService(repository *repository.Repository) *SimService {
	ss := &SimService{
		repository: repository,
	}
	return ss
}
func (ss *SimService) Add(ctx context.Context, s *core.Sim) (int, error) {

	// retrive provider data
	// if provider doest not exist, add it

	provider, err := ss.repository.ProviderRepository.ByName(ctx, s.Provider().Name())
	if err != nil {
		if !errors.Is(err, repoerrors.ErrNotFound) {
			return 0, err
		}

		id, err := ss.repository.ProviderRepository.Add(ctx, s.Provider().Name())
		if err != nil {
			return 0, err
		}
		s.Provider().SetId(id)
	} else {
		s.Provider().SetId(provider.Id())
	}

	return ss.repository.SimRepository.Add(ctx, s.Number(), s.Provider(), s.IsActivated(), s.ActivateUntil(), s.IsBlocked())
}
func (ss *SimService) Remove(ctx context.Context, id int) error {
	return ss.repository.SimRepository.Remove(ctx, id)
}
func (ss *SimService) GetSimList(ctx context.Context) (*core.List[*core.Sim], error) {
	return ss.repository.SimRepository.GetList(ctx)
}
func (ss *SimService) ActivateSim(ctx context.Context, id int) error {
	sim, err := ss.repository.SimRepository.ByID(ctx, id)
	if err != nil {
		return err
	}
	sim.SetActivated(true)
	return ss.repository.SimRepository.Update(ctx, sim)
}
func (ss *SimService) BlockSim(ctx context.Context, id int) error {
	sim, err := ss.repository.SimRepository.ByID(ctx, id)
	if err != nil {
		return err
	}
	sim.SetBlocked(true)
	return ss.repository.SimRepository.Update(ctx, sim)
}

func (ss *SimService) GetUsedServiceList(ctx context.Context, id int) (core.List[*core.Used], error) {
	panic("")
}
