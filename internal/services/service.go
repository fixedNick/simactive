package services

import (
	"context"
	"simactive/internal/core"
	repository "simactive/internal/infrastructure"
)

type ServiceService struct {
	repository *repository.Repository
}

func NewServiceService(repo *repository.Repository) *ServiceService {
	ss := &ServiceService{
		repository: repo,
	}
	return ss
}

func (ss *ServiceService) Add(ctx context.Context, s *core.Service) (int, error) {
	return ss.repository.ServiceRepository.Add(ctx, s.Name())
}
func (ss *ServiceService) Remove(ctx context.Context, id int) error {
	return ss.repository.ServiceRepository.Remove(ctx, id)
}
func (ss *ServiceService) GetServiceList(ctx context.Context) (*core.List[*core.Service], error) {
	return ss.repository.ServiceRepository.GetList(ctx)
}
