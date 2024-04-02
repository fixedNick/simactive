package simrepo

import (
	"fmt"
	"simactive/internal/core/sim"
)

type SimRepository interface {
	Save(s sim.Sim) error
	Remove(u sim.Sim) error
	SimList() (*sim.SimList, error)
	ById(id int) (sim.Sim, error)
	ByNumber(number string) (sim.Sim, error)
}

type SimInMemoryRepo struct {
	list sim.SimList
}

// Creates a new In-Memory repository
func NewInMemoryRepository(s ...sim.Sim) *SimInMemoryRepo {
	return &SimInMemoryRepo{
		list: sim.NewSimList(s...),
	}
}

// Saving [s Sim] into in-memory repository
// Returns [error] if sim with same number already in repo
func (simRepo SimInMemoryRepo) Save(s sim.Sim) error {
	if _, ok := simRepo.list[s.Number()]; ok {
		return fmt.Errorf("sim with number %s already exists", s.Number())
	}

	simRepo.list[s.Number()] = s
	return nil
}

// Removing [s Sim] from in-memory repository
// Returns [error] if sim does not exist in repo
func (simRepo SimInMemoryRepo) Remove(s sim.Sim) error {
	if _, ok := simRepo.list[s.Number()]; ok {
		delete(simRepo.list, s.Number())
		return nil
	}
	return fmt.Errorf("sim with number %s does not exist on in-memory repo", s.Number())
}

// Returns [sim.SimList - map[string]Sim] where key is sim.Number()
// Returns [error] if list is not initialized
func (simRepo SimInMemoryRepo) SimList() (*sim.SimList, error) {
	if simRepo.list == nil {
		return nil, fmt.Errorf("sim list is not initialized")
	}
	return &simRepo.list, nil
}

// calls SimList of current repo
// returns [s Sim] found by id
// returns [error] if sim not found
// returns [outer error] if List is not initialized
func (simRepo SimInMemoryRepo) ById(id int) (sim.Sim, error) {
	list, err := simRepo.SimList()
	if err != nil {
		return sim.Sim{}, err
	}

	for _, s := range *list {
		if id == s.Id() {
			return s, nil
		}
	}

	return sim.Sim{}, fmt.Errorf("sim with id %d not found", id)
}

// calls SimList of current repo
// returns [s Sim] found by id
// returns [error] if sim not found
// returns [error] if List is not initialized
func (simRepo SimInMemoryRepo) ByNumber(number string) (sim.Sim, error) {
	list, err := simRepo.SimList()
	if err != nil {
		return sim.Sim{}, err
	}

	if s, ok := (*list)[number]; ok {
		return s, nil
	}
	return sim.Sim{}, fmt.Errorf("sim with number %s not found at in-memory repo", number)
}
