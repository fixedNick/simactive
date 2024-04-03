package services

import (
	"context"
	"simactive/internal/core"
	"simactive/internal/repository"
)

type SimService struct {
	inMemoryRepo *repository.SimInMemoryRepo
	sqlRepo      *repository.SimSqlRepository
}

func NewSimService(inMemRepository *repository.SimInMemoryRepo, sqlRepository *repository.SimSqlRepository) *SimService {
	return &SimService{
		inMemoryRepo: inMemRepository,
		sqlRepo:      sqlRepository,
	}
}

func (ss *SimService) Add(ctx context.Context, s core.Sim) error {
	id, err := ss.sqlRepo.Save(s)
	if err != nil {
		return err
	}

	s.SetID(id)

	if err = ss.inMemoryRepo.Save(s); err != nil {
		return err
	}
	return nil
}
func (ss *SimService) GetByID(ctx context.Context, id int) (sim core.Sim, err error) {
	panic("")
}
func (ss *SimService) Remove(ctx context.Context, s core.Sim) error {
	panic("")
}
