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
	ByID(ctx context.Context, id int) (core.Sim, error)
	Update(ctx context.Context, s *core.Sim) (err error)
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
func (ss *SimService) ActivateSim(ctx context.Context, id int) error {
	// get sim from list
	// in-mem -> sql

	receivedSim, err := ss.getByID(ctx, id)
	if err != nil {
		return err
	}

	receivedSim.SetActivated(true)

	// try to update repositories
	if imUpdErr := ss.inMemoryRepo.Update(ctx, receivedSim); imUpdErr != nil {
		return fmt.Errorf("(in-memory repo) cannot update data about sim with id `%d`[%v]. Error: %v", id, receivedSim, imUpdErr)
	}
	if sqlUpdErr := ss.sqlRepo.Update(ctx, receivedSim); sqlUpdErr != nil {
		return fmt.Errorf("(sqlrepo) cannot update data about sim with id `%d`[%v]. Error: %v", id, receivedSim, sqlUpdErr)
	}
	return nil
}
func (ss *SimService) BlockSim(ctx context.Context, id int) error {

	receivedSim, err := ss.getByID(ctx, id)
	if err != nil {
		return err
	}

	receivedSim.SetBlocked(true)

	if imUpdErr := ss.inMemoryRepo.Update(ctx, receivedSim); imUpdErr != nil {
		return fmt.Errorf("(in-memory repo) cannot update data about sim with id `%d`[%v]. Error: %v", id, receivedSim, imUpdErr)
	}
	if sqlUpdErr := ss.sqlRepo.Update(ctx, receivedSim); sqlUpdErr != nil {
		return fmt.Errorf("(sqlrepo) cannot update data about sim with id `%d`[%v]. Error: %v", id, receivedSim, sqlUpdErr)
	}
	return nil
}
func (ss *SimService) getByID(ctx context.Context, id int) (*core.Sim, error) {
	var receivedSim core.Sim

	receivedSim, imErr := ss.inMemoryRepo.ByID(ctx, id)
	if imErr != nil {
		// try to get from sql
		sim, sErr := ss.sqlRepo.ByID(ctx, id)
		if sErr != nil {
			return nil, fmt.Errorf("cannot get sim with id `%d` from in-memory err: [%v] and from sql err: [%v]", id, imErr, sErr)
		}

		// setup outer var
		receivedSim = sim
		// try to save into local repository
		_, err := ss.inMemoryRepo.Save(ctx, sim)
		if err != nil {
			return nil, fmt.Errorf("received sim with id `%d` from sql but got erros with in-memory.\nById error: [%v]\nSave into memory error: [%v]", id, imErr, err)
		}
	}

	return &receivedSim, nil
}
