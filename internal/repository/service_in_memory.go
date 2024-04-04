package repository

import (
	"context"
	"simactive/internal/core"
)

type ServiceInMemoryRepo struct {
	list core.ServiceList
}

func NewServiceInMemoryRepository() ServiceInMemoryRepo {
	return ServiceInMemoryRepo{
		list: core.NewServiceList(),
	}
}

func (sr ServiceInMemoryRepo) Save(ctx context.Context, service core.Service) (id int, err error) {
	panic("implement")
}
func (sr ServiceInMemoryRepo) GetServiceList(ctx context.Context) (list *core.ServiceList, err error) {
	panic("implement")
}
func (sr ServiceInMemoryRepo) Remove(ctx context.Context, id int) (err error) {
	panic("implement")
}
func (sr ServiceInMemoryRepo) Update(ctx context.Context, s *core.Service) (err error) {
	panic("implement")
}
