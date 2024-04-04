package repository

import (
	"context"
	"fmt"
	"log"
	"simactive/internal/core"
)

type SimInMemoryRepo struct {
	list core.SimList
}

// Creates a new In-Memory repository
func NewSimInMemoryRepository(s ...core.Sim) *SimInMemoryRepo {
	return &SimInMemoryRepo{
		list: core.NewSimList(s...),
	}
}

// Saving [s Sim] into in-memory repository
// Returns [error] if core with same number already in repo
func (coreRepo SimInMemoryRepo) Save(ctx context.Context, s core.Sim) (int, error) {
	if _, ok := coreRepo.list[s.Id()]; ok {
		return 0, fmt.Errorf("sim with number %s already exists", s.Number())
	}

	coreRepo.list[s.Id()] = s

	log.Printf("Sim with number %s saved into in-memory data with id %d", s.Number(), s.Id())
	return s.Id(), nil
}

// Removing [s Sim] from in-memory repository
// Returns [error] if core does not exist in repo
func (coreRepo SimInMemoryRepo) Remove(ctx context.Context, id int) error {
	if _, ok := coreRepo.list[id]; ok {
		delete(coreRepo.list, id)
		return nil
	}
	return fmt.Errorf("sim with id [%d]does not exist on in-memory repo", id)
}

// Returns [core.SimList - map[string]Sim] where key is core.Number()
// Returns [error] if list is not initialized
func (coreRepo SimInMemoryRepo) GetSimList(ctx context.Context) (*core.SimList, error) {
	if coreRepo.list == nil {
		return nil, fmt.Errorf("sim list is not initialized")
	}
	return &coreRepo.list, nil
}

// calls SimList of current repo
// returns [s Sim] found by id
// returns [error] if core not found
// returns [outer error] if List is not initialized
func (coreRepo SimInMemoryRepo) ById(ctx context.Context, id int) (core.Sim, error) {
	list, err := coreRepo.GetSimList(ctx)
	if err != nil {
		return core.Sim{}, err
	}

	for _, s := range *list {
		if id == s.Id() {
			return s, nil
		}
	}

	return core.Sim{}, fmt.Errorf("sim with id %d not found", id)
}
