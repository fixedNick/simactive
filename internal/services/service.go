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

func (ss *ServiceService) Add(ctx context.Context, s core.Service) error {
	panic("implement")
}
func (ss *ServiceService) Remove(ctx context.Context, id int) error {
	panic("implement")
}
func (ss *ServiceService) GetServiceList(ctx context.Context) (*core.List[*core.Service], error) {
	panic("implement")
}
