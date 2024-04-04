package services

import (
	"context"
	"log"
	"simactive/internal/core"
)

type ServiceRepository interface {
	Save(ctx context.Context, service core.Service) (id int, err error)
	GetServiceList(ctx context.Context) (list *core.ServiceList, err error)
	Remove(ctx context.Context, id int) (err error)
	Update(ctx context.Context, s *core.Service) (err error)
}

type ServiceService struct {
	inMemoryRepo ServiceRepository
	sqlRepo      ServiceRepository
}

func NewServiceService(inMemoryRepo ServiceRepository, sqlRepo ServiceRepository) *ServiceService {
	ss := &ServiceService{
		inMemoryRepo: inMemoryRepo,
		sqlRepo:      sqlRepo,
	}

	ctx := context.Background()

	serviceList, err := ss.sqlRepo.GetServiceList(ctx)
	if err != nil {
		log.Fatalf("Fatal error on getting sim list from sql: %v", err)
	}

	for _, sim := range *serviceList {
		if _, err := ss.inMemoryRepo.Save(ctx, sim); err != nil {
			log.Fatalf("Error on loading database into memory: %v", err)
		}
	}

	log.Printf("Loaded %d simcards into memory", len(*serviceList))

	return ss
}

func (ss *ServiceService) Add(ctx context.Context, s core.Service) error {
	panic("implement")
}
func (ss *ServiceService) Remove(ctx context.Context, id int) error {
	panic("implement")
}
func (ss *ServiceService) GetServiceList(ctx context.Context) (*core.ServiceList, error) {
	panic("implement")
}
