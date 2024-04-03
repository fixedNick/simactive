package repository

import (
	"fmt"
	"simactive/internal/core"
)

type SimInMemoryRepo struct {
	list core.SimList
}

// Creates a new In-Memory repository
func NewInMemoryRepository(s ...core.Sim) *SimInMemoryRepo {
	return &SimInMemoryRepo{
		list: core.NewSimList(s...),
	}
}

// Saving [s Sim] into in-memory repository
// Returns [error] if core with same number already in repo
func (coreRepo SimInMemoryRepo) Save(s core.Sim) error {
	if _, ok := coreRepo.list[s.Number()]; ok {
		return fmt.Errorf("core with number %s already exists", s.Number())
	}

	coreRepo.list[s.Number()] = s
	return nil
}

// Removing [s Sim] from in-memory repository
// Returns [error] if core does not exist in repo
func (coreRepo SimInMemoryRepo) Remove(s core.Sim) error {
	if _, ok := coreRepo.list[s.Number()]; ok {
		delete(coreRepo.list, s.Number())
		return nil
	}
	return fmt.Errorf("core with number %s does not exist on in-memory repo", s.Number())
}

// Returns [core.SimList - map[string]Sim] where key is core.Number()
// Returns [error] if list is not initialized
func (coreRepo SimInMemoryRepo) SimList() (*core.SimList, error) {
	if coreRepo.list == nil {
		return nil, fmt.Errorf("core list is not initialized")
	}
	return &coreRepo.list, nil
}

// calls SimList of current repo
// returns [s Sim] found by id
// returns [error] if core not found
// returns [outer error] if List is not initialized
func (coreRepo SimInMemoryRepo) ById(id int) (core.Sim, error) {
	list, err := coreRepo.SimList()
	if err != nil {
		return core.Sim{}, err
	}

	for _, s := range *list {
		if id == s.Id() {
			return s, nil
		}
	}

	return core.Sim{}, fmt.Errorf("core with id %d not found", id)
}

// calls SimList of current repo
// returns [s Sim] found by id
// returns [error] if core not found
// returns [error] if List is not initialized
func (coreRepo SimInMemoryRepo) ByNumber(number string) (core.Sim, error) {
	list, err := coreRepo.SimList()
	if err != nil {
		return core.Sim{}, err
	}

	if s, ok := (*list)[number]; ok {
		return s, nil
	}
	return core.Sim{}, fmt.Errorf("core with number %s not found at in-memory repo", number)
}
