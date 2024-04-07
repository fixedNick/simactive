package services

import (
	"context"
	"simactive/internal/core"
)

type ServiceService struct {
	repository Repository[*core.Service]
}

func NewServiceService(repo Repository[*core.Service]) *ServiceService {
	ss := &ServiceService{
		repository: repo,
	}
	return ss
}

func (ss *ServiceService) Add(ctx context.Context, s *core.Service) (int, error) {
	return ss.repository.Save(ctx, s)
}
func (ss *ServiceService) Remove(ctx context.Context, id int) error {
	return ss.repository.Remove(ctx, id)
}
func (ss *ServiceService) GetServiceList(ctx context.Context) (*core.List[*core.Service], error) {
	return ss.repository.GetList(ctx)
}
