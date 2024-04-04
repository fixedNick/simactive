package services

import (
	"context"
	"fmt"
	"log"
	"simactive/internal/core"
)

type SimRepository interface {
	Save(ctx context.Context, sim core.Sim) (id int, err error)
	GetSimList(ctx context.Context) (list *core.SimList, err error)
	Remove(ctx context.Context, id int) (err error)
}

type SimService struct {
	inMemoryRepo SimRepository
	sqlRepo      SimRepository
}

func NewSimService(inMemRepository SimRepository, sqlRepository SimRepository) *SimService {
	ss := &SimService{
		inMemoryRepo: inMemRepository,
		sqlRepo:      sqlRepository,
	}

	ctx := context.Background()

	simList, err := ss.sqlRepo.GetSimList(ctx)
	if err != nil {
		log.Fatalf("Fatal error on getting sim list from sql: %v", err)
	}

	for _, sim := range *simList {
		if _, err := ss.inMemoryRepo.Save(ctx, sim); err != nil {
			log.Fatalf("Error on loading database into memory: %v", err)
		}
	}

	log.Printf("Loaded %d simcards into memory", len(*simList))

	return ss
}
func (ss *SimService) Add(ctx context.Context, s core.Sim) error {
	id, err := ss.sqlRepo.Save(ctx, s)
	if err != nil {
		return err
	}

	s.SetID(id)

	if _, err = ss.inMemoryRepo.Save(ctx, s); err != nil {
		return err
	}

	return nil
}
func (ss *SimService) GetByID(ctx context.Context, id int) (sim core.Sim, err error) {
	panic("")
}
func (ss *SimService) Remove(ctx context.Context, id int) error {

	// firstly remove it from local storage
	// then remove from sql

	if err := ss.inMemoryRepo.Remove(ctx, id); err != nil {
		return fmt.Errorf("error ocured on in-memory repository: %v", err)
	}
	if err := ss.sqlRepo.Remove(ctx, id); err != nil {
		return fmt.Errorf("error ocured on sql repository: %v", err)
	}

	return nil
}

func (ss *SimService) GetSimList(ctx context.Context) (*core.SimList, error) {
	return ss.inMemoryRepo.GetSimList(ctx)
}
